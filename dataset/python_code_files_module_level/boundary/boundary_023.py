import copy
from collections import defaultdict
from dataclasses import dataclass
from typing import Any, Optional, Union
from huggingface_hub.utils import logging, yaml_dump
logger = logging.get_logger(__name__)

@dataclass
class EvalResult:
    task_type: str
    dataset_type: str
    dataset_name: str
    metric_type: str
    metric_value: Any
    task_name: Optional[str] = None
    dataset_config: Optional[str] = None
    dataset_split: Optional[str] = None
    dataset_revision: Optional[str] = None
    dataset_args: Optional[dict[str, Any]] = None
    metric_name: Optional[str] = None
    metric_config: Optional[str] = None
    metric_args: Optional[dict[str, Any]] = None
    verified: Optional[bool] = None
    verify_token: Optional[str] = None
    source_name: Optional[str] = None
    source_url: Optional[str] = None

    @property
    def unique_identifier(self) -> tuple:
        return (self.task_type, self.dataset_type, self.dataset_config, self.dataset_split, self.dataset_revision)

    def is_equal_except_value(self, other: 'EvalResult') -> bool:
        for key, _ in self.__dict__.items():
            if key == 'metric_value':
                continue
            if key != 'verify_token' and getattr(self, key) != getattr(other, key):
                return False
        return True

    def __post_init__(self) -> None:
        if self.source_name is not None and self.source_url is None:
            raise ValueError('If `source_name` is provided, `source_url` must also be provided.')

@dataclass
class CardData:

    def __init__(self, ignore_metadata_errors: bool=False, **kwargs):
        self.__dict__.update(kwargs)

    def to_dict(self):
        data_dict = copy.deepcopy(self.__dict__)
        self._to_dict(data_dict)
        return {key: value for key, value in data_dict.items() if value is not None}

    def _to_dict(self, data_dict):
        pass

    def to_yaml(self, line_break=None, original_order: Optional[list[str]]=None) -> str:
        if original_order:
            self.__dict__ = {k: self.__dict__[k] for k in original_order + list(set(self.__dict__.keys()) - set(original_order)) if k in self.__dict__}
        return yaml_dump(self.to_dict(), sort_keys=False, line_break=line_break).strip()

    def __repr__(self):
        return repr(self.__dict__)

    def __str__(self):
        return self.to_yaml()

    def get(self, key: str, default: Any=None) -> Any:
        value = self.__dict__.get(key)
        return default if value is None else value

    def pop(self, key: str, default: Any=None) -> Any:
        return self.__dict__.pop(key, default)

    def __getitem__(self, key: str) -> Any:
        return self.__dict__[key]

    def __setitem__(self, key: str, value: Any) -> None:
        self.__dict__[key] = value

    def __contains__(self, key: str) -> bool:
        return key in self.__dict__

    def __len__(self) -> int:
        return len(self.__dict__)

def _validate_eval_results(eval_results: Optional[Union[EvalResult, list[EvalResult]]], model_name: Optional[str]) -> list[EvalResult]:
    if eval_results is None:
        return []
    if isinstance(eval_results, EvalResult):
        eval_results = [eval_results]
    if not isinstance(eval_results, list) or not all((isinstance(r, EvalResult) for r in eval_results)):
        raise ValueError(f'`eval_results` should be of type `EvalResult` or a list of `EvalResult`, got {type(eval_results)}.')
    if model_name is None:
        raise ValueError('Passing `eval_results` requires `model_name` to be set.')
    return eval_results

class ModelCardData(CardData):

    def __init__(self, *, base_model: Optional[Union[str, list[str]]]=None, datasets: Optional[Union[str, list[str]]]=None, eval_results: Optional[list[EvalResult]]=None, language: Optional[Union[str, list[str]]]=None, library_name: Optional[str]=None, license: Optional[str]=None, license_name: Optional[str]=None, license_link: Optional[str]=None, metrics: Optional[list[str]]=None, model_name: Optional[str]=None, pipeline_tag: Optional[str]=None, tags: Optional[list[str]]=None, ignore_metadata_errors: bool=False, **kwargs):
        self.base_model = base_model
        self.datasets = datasets
        self.eval_results = eval_results
        self.language = language
        self.library_name = library_name
        self.license = license
        self.license_name = license_name
        self.license_link = license_link
        self.metrics = metrics
        self.model_name = model_name
        self.pipeline_tag = pipeline_tag
        self.tags = _to_unique_list(tags)
        model_index = kwargs.pop('model-index', None)
        if model_index:
            try:
                model_name, eval_results = model_index_to_eval_results(model_index)
                self.model_name = model_name
                self.eval_results = eval_results
            except (KeyError, TypeError) as error:
                if ignore_metadata_errors:
                    logger.warning('Invalid model-index. Not loading eval results into CardData.')
                else:
                    raise ValueError(f'Invalid `model_index` in metadata cannot be parsed: {error.__class__} {error}. Pass `ignore_metadata_errors=True` to ignore this error while loading a Model Card. Warning: some information will be lost. Use it at your own risk.')
        super().__init__(**kwargs)
        if self.eval_results:
            try:
                self.eval_results = _validate_eval_results(self.eval_results, self.model_name)
            except Exception as e:
                if ignore_metadata_errors:
                    logger.warning(f'Failed to validate eval_results: {e}. Not loading eval results into CardData.')
                else:
                    raise ValueError(f'Failed to validate eval_results: {e}') from e

    def _to_dict(self, data_dict):
        if self.eval_results is not None:
            data_dict['model-index'] = eval_results_to_model_index(self.model_name, self.eval_results)
            del data_dict['eval_results'], data_dict['model_name']

class DatasetCardData(CardData):

    def __init__(self, *, language: Optional[Union[str, list[str]]]=None, license: Optional[Union[str, list[str]]]=None, annotations_creators: Optional[Union[str, list[str]]]=None, language_creators: Optional[Union[str, list[str]]]=None, multilinguality: Optional[Union[str, list[str]]]=None, size_categories: Optional[Union[str, list[str]]]=None, source_datasets: Optional[list[str]]=None, task_categories: Optional[Union[str, list[str]]]=None, task_ids: Optional[Union[str, list[str]]]=None, paperswithcode_id: Optional[str]=None, pretty_name: Optional[str]=None, train_eval_index: Optional[dict]=None, config_names: Optional[Union[str, list[str]]]=None, ignore_metadata_errors: bool=False, **kwargs):
        self.annotations_creators = annotations_creators
        self.language_creators = language_creators
        self.language = language
        self.license = license
        self.multilinguality = multilinguality
        self.size_categories = size_categories
        self.source_datasets = source_datasets
        self.task_categories = task_categories
        self.task_ids = task_ids
        self.paperswithcode_id = paperswithcode_id
        self.pretty_name = pretty_name
        self.config_names = config_names
        self.train_eval_index = train_eval_index or kwargs.pop('train-eval-index', None)
        super().__init__(**kwargs)

    def _to_dict(self, data_dict):
        data_dict['train-eval-index'] = data_dict.pop('train_eval_index')

class SpaceCardData(CardData):

    def __init__(self, *, title: Optional[str]=None, sdk: Optional[str]=None, sdk_version: Optional[str]=None, python_version: Optional[str]=None, app_file: Optional[str]=None, app_port: Optional[int]=None, license: Optional[str]=None, duplicated_from: Optional[str]=None, models: Optional[list[str]]=None, datasets: Optional[list[str]]=None, tags: Optional[list[str]]=None, ignore_metadata_errors: bool=False, **kwargs):
        self.title = title
        self.sdk = sdk
        self.sdk_version = sdk_version
        self.python_version = python_version
        self.app_file = app_file
        self.app_port = app_port
        self.license = license
        self.duplicated_from = duplicated_from
        self.models = models
        self.datasets = datasets
        self.tags = _to_unique_list(tags)
        super().__init__(**kwargs)

def model_index_to_eval_results(model_index: list[dict[str, Any]]) -> tuple[str, list[EvalResult]]:
    eval_results = []
    for elem in model_index:
        name = elem['name']
        results = elem['results']
        for result in results:
            task_type = result['task']['type']
            task_name = result['task'].get('name')
            dataset_type = result['dataset']['type']
            dataset_name = result['dataset']['name']
            dataset_config = result['dataset'].get('config')
            dataset_split = result['dataset'].get('split')
            dataset_revision = result['dataset'].get('revision')
            dataset_args = result['dataset'].get('args')
            source_name = result.get('source', {}).get('name')
            source_url = result.get('source', {}).get('url')
            for metric in result['metrics']:
                metric_type = metric['type']
                metric_value = metric['value']
                metric_name = metric.get('name')
                metric_args = metric.get('args')
                metric_config = metric.get('config')
                verified = metric.get('verified')
                verify_token = metric.get('verifyToken')
                eval_result = EvalResult(task_type=task_type, dataset_type=dataset_type, dataset_name=dataset_name, metric_type=metric_type, metric_value=metric_value, task_name=task_name, dataset_config=dataset_config, dataset_split=dataset_split, dataset_revision=dataset_revision, dataset_args=dataset_args, metric_name=metric_name, metric_args=metric_args, metric_config=metric_config, verified=verified, verify_token=verify_token, source_name=source_name, source_url=source_url)
                eval_results.append(eval_result)
    return (name, eval_results)

def _remove_none(obj):
    if isinstance(obj, (list, tuple, set)):
        return type(obj)((_remove_none(x) for x in obj if x is not None))
    elif isinstance(obj, dict):
        return type(obj)(((_remove_none(k), _remove_none(v)) for k, v in obj.items() if k is not None and v is not None))
    else:
        return obj

def eval_results_to_model_index(model_name: str, eval_results: list[EvalResult]) -> list[dict[str, Any]]:
    task_and_ds_types_map: dict[Any, list[EvalResult]] = defaultdict(list)
    for eval_result in eval_results:
        task_and_ds_types_map[eval_result.unique_identifier].append(eval_result)
    model_index_data: list[dict[str, Any]] = []
    for results in task_and_ds_types_map.values():
        sample_result = results[0]
        data: dict[str, Any] = {'task': {'type': sample_result.task_type, 'name': sample_result.task_name}, 'dataset': {'name': sample_result.dataset_name, 'type': sample_result.dataset_type, 'config': sample_result.dataset_config, 'split': sample_result.dataset_split, 'revision': sample_result.dataset_revision, 'args': sample_result.dataset_args}, 'metrics': [{'type': result.metric_type, 'value': result.metric_value, 'name': result.metric_name, 'config': result.metric_config, 'args': result.metric_args, 'verified': result.verified, 'verifyToken': result.verify_token} for result in results]}
        if sample_result.source_url is not None:
            source: dict[str, str] = {'url': sample_result.source_url}
            if sample_result.source_name is not None:
                source['name'] = sample_result.source_name
            data['source'] = source
        model_index_data.append(data)
    model_index = [{'name': model_name, 'results': model_index_data}]
    return _remove_none(model_index)

def _to_unique_list(tags: Optional[list[str]]) -> Optional[list[str]]:
    if tags is None:
        return tags
    unique_tags = []
    for tag in tags:
        if tag not in unique_tags:
            unique_tags.append(tag)
    return unique_tags
