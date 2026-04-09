import os
import re
from pathlib import Path
from typing import Any, Literal, Optional, Union
import yaml
from huggingface_hub.file_download import hf_hub_download
from huggingface_hub.hf_api import upload_file
from huggingface_hub.repocard_data import CardData, DatasetCardData, EvalResult, ModelCardData, SpaceCardData, eval_results_to_model_index, model_index_to_eval_results
from huggingface_hub.utils import HfHubHTTPError, get_session, hf_raise_for_status, is_jinja_available, yaml_dump
from . import constants
from .errors import EntryNotFoundError
from .utils import SoftTemporaryDirectory, logging, validate_hf_hub_args
logger = logging.get_logger(__name__)
TEMPLATE_MODELCARD_PATH = Path(__file__).parent / 'templates' / 'modelcard_template.md'
TEMPLATE_DATASETCARD_PATH = Path(__file__).parent / 'templates' / 'datasetcard_template.md'
REGEX_YAML_BLOCK = re.compile('^(\\s*---[\\r\\n]+)([\\S\\s]*?)([\\r\\n]+---(\\r\\n|\\n|$))')

class RepoCard:
    card_data_class = CardData
    default_template_path = TEMPLATE_MODELCARD_PATH
    repo_type = 'model'

    def __init__(self, content: str, ignore_metadata_errors: bool=False):
        self.ignore_metadata_errors = ignore_metadata_errors
        self.content = content

    @property
    def content(self):
        line_break = _detect_line_ending(self._content) or '\n'
        return f'---{line_break}{self.data.to_yaml(line_break=line_break, original_order=self._original_order)}{line_break}---{line_break}{self.text}'

    @content.setter
    def content(self, content: str):
        self._content = content
        match = REGEX_YAML_BLOCK.search(content)
        if match:
            yaml_block = match.group(2)
            self.text = content[match.end():]
            data_dict = yaml.safe_load(yaml_block)
            if data_dict is None:
                data_dict = {}
            if not isinstance(data_dict, dict):
                raise ValueError('repo card metadata block should be a dict')
        else:
            logger.warning('Repo card metadata block was not found. Setting CardData to empty.')
            data_dict = {}
            self.text = content
        self.data = self.card_data_class(**data_dict, ignore_metadata_errors=self.ignore_metadata_errors)
        self._original_order = list(data_dict.keys())

    def __str__(self):
        return self.content

    def save(self, filepath: Union[Path, str]):
        filepath = Path(filepath)
        filepath.parent.mkdir(parents=True, exist_ok=True)
        with open(filepath, mode='w', newline='', encoding='utf-8') as f:
            f.write(str(self))

    @classmethod
    def load(cls, repo_id_or_path: Union[str, Path], repo_type: Optional[str]=None, token: Optional[str]=None, ignore_metadata_errors: bool=False):
        if Path(repo_id_or_path).is_file():
            card_path = Path(repo_id_or_path)
        elif isinstance(repo_id_or_path, str):
            card_path = Path(hf_hub_download(repo_id_or_path, constants.REPOCARD_NAME, repo_type=repo_type or cls.repo_type, token=token))
        else:
            raise ValueError(f'Cannot load RepoCard: path not found on disk ({repo_id_or_path}).')
        with card_path.open(mode='r', newline='', encoding='utf-8') as f:
            return cls(f.read(), ignore_metadata_errors=ignore_metadata_errors)

    def validate(self, repo_type: Optional[str]=None):
        repo_type = repo_type or self.repo_type
        body = {'repoType': repo_type, 'content': str(self)}
        headers = {'Accept': 'text/plain'}
        try:
            response = get_session().post('https://huggingface.co/api/validate-yaml', json=body, headers=headers)
            hf_raise_for_status(response)
        except HfHubHTTPError as exc:
            if response.status_code == 400:
                raise ValueError(response.text)
            else:
                raise exc

    def push_to_hub(self, repo_id: str, token: Optional[str]=None, repo_type: Optional[str]=None, commit_message: Optional[str]=None, commit_description: Optional[str]=None, revision: Optional[str]=None, create_pr: Optional[bool]=None, parent_commit: Optional[str]=None):
        repo_type = repo_type or self.repo_type
        self.validate(repo_type=repo_type)
        with SoftTemporaryDirectory() as tmpdir:
            tmp_path = Path(tmpdir) / constants.REPOCARD_NAME
            tmp_path.write_text(str(self), encoding='utf-8')
            url = upload_file(path_or_fileobj=str(tmp_path), path_in_repo=constants.REPOCARD_NAME, repo_id=repo_id, token=token, repo_type=repo_type, commit_message=commit_message, commit_description=commit_description, create_pr=create_pr, revision=revision, parent_commit=parent_commit)
        return url

    @classmethod
    def from_template(cls, card_data: CardData, template_path: Optional[str]=None, template_str: Optional[str]=None, **template_kwargs):
        if is_jinja_available():
            import jinja2
        else:
            raise ImportError('Using RepoCard.from_template requires Jinja2 to be installed. Please install it with `pip install Jinja2`.')
        kwargs = card_data.to_dict().copy()
        kwargs.update(template_kwargs)
        if template_path is not None:
            template_str = Path(template_path).read_text()
        if template_str is None:
            template_str = Path(cls.default_template_path).read_text()
        template = jinja2.Template(template_str)
        content = template.render(card_data=card_data.to_yaml(), **kwargs)
        return cls(content)

class ModelCard(RepoCard):
    card_data_class = ModelCardData
    default_template_path = TEMPLATE_MODELCARD_PATH
    repo_type = 'model'

    @classmethod
    def from_template(cls, card_data: ModelCardData, template_path: Optional[str]=None, template_str: Optional[str]=None, **template_kwargs):
        return super().from_template(card_data, template_path, template_str, **template_kwargs)

class DatasetCard(RepoCard):
    card_data_class = DatasetCardData
    default_template_path = TEMPLATE_DATASETCARD_PATH
    repo_type = 'dataset'

    @classmethod
    def from_template(cls, card_data: DatasetCardData, template_path: Optional[str]=None, template_str: Optional[str]=None, **template_kwargs):
        return super().from_template(card_data, template_path, template_str, **template_kwargs)

class SpaceCard(RepoCard):
    card_data_class = SpaceCardData
    default_template_path = TEMPLATE_MODELCARD_PATH
    repo_type = 'space'

def _detect_line_ending(content: str) -> Literal['\r', '\n', '\r\n', None]:
    cr = content.count('\r')
    lf = content.count('\n')
    crlf = content.count('\r\n')
    if cr + lf == 0:
        return None
    if crlf == cr and crlf == lf:
        return '\r\n'
    if cr > lf:
        return '\r'
    else:
        return '\n'

def metadata_load(local_path: Union[str, Path]) -> Optional[dict]:
    content = Path(local_path).read_text()
    match = REGEX_YAML_BLOCK.search(content)
    if match:
        yaml_block = match.group(2)
        data = yaml.safe_load(yaml_block)
        if data is None or isinstance(data, dict):
            return data
        raise ValueError('repo card metadata block should be a dict')
    else:
        return None

def metadata_save(local_path: Union[str, Path], data: dict) -> None:
    line_break = '\n'
    content = ''
    if os.path.exists(local_path):
        with open(local_path, 'r', newline='', encoding='utf8') as readme:
            content = readme.read()
            if isinstance(readme.newlines, tuple):
                line_break = readme.newlines[0]
            elif isinstance(readme.newlines, str):
                line_break = readme.newlines
    with open(local_path, 'w', newline='', encoding='utf8') as readme:
        data_yaml = yaml_dump(data, sort_keys=False, line_break=line_break)
        match = REGEX_YAML_BLOCK.search(content)
        if match:
            output = content[:match.start()] + f'---{line_break}{data_yaml}---{line_break}' + content[match.end():]
        else:
            output = f'---{line_break}{data_yaml}---{line_break}{content}'
        readme.write(output)
        readme.close()

def metadata_eval_result(*, model_pretty_name: str, task_pretty_name: str, task_id: str, metrics_pretty_name: str, metrics_id: str, metrics_value: Any, dataset_pretty_name: str, dataset_id: str, metrics_config: Optional[str]=None, metrics_verified: bool=False, dataset_config: Optional[str]=None, dataset_split: Optional[str]=None, dataset_revision: Optional[str]=None, metrics_verification_token: Optional[str]=None) -> dict:
    return {'model-index': eval_results_to_model_index(model_name=model_pretty_name, eval_results=[EvalResult(task_name=task_pretty_name, task_type=task_id, metric_name=metrics_pretty_name, metric_type=metrics_id, metric_value=metrics_value, dataset_name=dataset_pretty_name, dataset_type=dataset_id, metric_config=metrics_config, verified=metrics_verified, verify_token=metrics_verification_token, dataset_config=dataset_config, dataset_split=dataset_split, dataset_revision=dataset_revision)])}

@validate_hf_hub_args
def metadata_update(repo_id: str, metadata: dict, *, repo_type: Optional[str]=None, overwrite: bool=False, token: Optional[str]=None, commit_message: Optional[str]=None, commit_description: Optional[str]=None, revision: Optional[str]=None, create_pr: bool=False, parent_commit: Optional[str]=None) -> str:
    commit_message = commit_message if commit_message is not None else 'Update metadata with huggingface_hub'
    card_class: type[RepoCard]
    if repo_type is None or repo_type == 'model':
        card_class = ModelCard
    elif repo_type == 'dataset':
        card_class = DatasetCard
    elif repo_type == 'space':
        card_class = RepoCard
    else:
        raise ValueError(f'Unknown repo_type: {repo_type}')
    try:
        card = card_class.load(repo_id, token=token, repo_type=repo_type)
    except EntryNotFoundError:
        if repo_type == 'space':
            raise ValueError("Cannot update metadata on a Space that doesn't contain a `README.md` file.")
        card = card_class.from_template(CardData())
    for key, value in metadata.items():
        if key == 'model-index':
            if 'name' not in value[0]:
                value[0]['name'] = getattr(card, 'model_name', repo_id)
            model_name, new_results = model_index_to_eval_results(value)
            if card.data.eval_results is None:
                card.data.eval_results = new_results
                card.data.model_name = model_name
            else:
                existing_results = card.data.eval_results
                for new_result in new_results:
                    result_found = False
                    for existing_result in existing_results:
                        if new_result.is_equal_except_value(existing_result):
                            if new_result != existing_result and (not overwrite):
                                raise ValueError(f"You passed a new value for the existing metric 'name: {new_result.metric_name}, type: {new_result.metric_type}'. Set `overwrite=True` to overwrite existing metrics.")
                            result_found = True
                            existing_result.metric_value = new_result.metric_value
                            if existing_result.verified is True:
                                existing_result.verify_token = new_result.verify_token
                    if not result_found:
                        card.data.eval_results.append(new_result)
        elif card.data.get(key) is not None and (not overwrite) and (card.data.get(key) != value):
            raise ValueError(f"You passed a new value for the existing meta data field '{key}'. Set `overwrite=True` to overwrite existing metadata.")
        else:
            card.data[key] = value
    return card.push_to_hub(repo_id, token=token, repo_type=repo_type, commit_message=commit_message, commit_description=commit_description, create_pr=create_pr, revision=revision, parent_commit=parent_commit)
