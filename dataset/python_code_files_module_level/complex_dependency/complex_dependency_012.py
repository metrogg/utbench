import inspect
import json
import os
from dataclasses import Field, asdict, dataclass, is_dataclass
from pathlib import Path
from typing import Any, Callable, ClassVar, Optional, Protocol, Type, TypeVar, Union
import packaging.version
from . import constants
from .errors import EntryNotFoundError, HfHubHTTPError
from .file_download import hf_hub_download
from .hf_api import HfApi
from .repocard import ModelCard, ModelCardData
from .utils import SoftTemporaryDirectory, is_jsonable, is_safetensors_available, is_simple_optional_type, is_torch_available, logging, unwrap_simple_optional_type, validate_hf_hub_args
if is_torch_available():
    import torch
if is_safetensors_available():
    import safetensors
    from safetensors.torch import load_model as load_model_as_safetensor
    from safetensors.torch import save_model as save_model_as_safetensor
logger = logging.get_logger(__name__)

class DataclassInstance(Protocol):
    __dataclass_fields__: ClassVar[dict[str, Field]]
T = TypeVar('T', bound='ModelHubMixin')
ARGS_T = TypeVar('ARGS_T')
ENCODER_T = Callable[[ARGS_T], Any]
DECODER_T = Callable[[Any], ARGS_T]
CODER_T = tuple[ENCODER_T, DECODER_T]
DEFAULT_MODEL_CARD = '\n---\n# For reference on model card metadata, see the spec: https://github.com/huggingface/hub-docs/blob/main/modelcard.md?plain=1\n# Doc / guide: https://huggingface.co/docs/hub/model-cards\n{{ card_data }}\n---\n\nThis model has been pushed to the Hub using the [PytorchModelHubMixin](https://huggingface.co/docs/huggingface_hub/package_reference/mixins#huggingface_hub.PyTorchModelHubMixin) integration:\n- Code: {{ repo_url | default("[More Information Needed]", true) }}\n- Paper: {{ paper_url | default("[More Information Needed]", true) }}\n- Docs: {{ docs_url | default("[More Information Needed]", true) }}\n'

@dataclass
class MixinInfo:
    model_card_template: str
    model_card_data: ModelCardData
    docs_url: Optional[str] = None
    paper_url: Optional[str] = None
    repo_url: Optional[str] = None

class ModelHubMixin:
    _hub_mixin_config: Optional[Union[dict, DataclassInstance]] = None
    _hub_mixin_info: MixinInfo
    _hub_mixin_inject_config: bool
    _hub_mixin_init_parameters: dict[str, inspect.Parameter]
    _hub_mixin_jsonable_default_values: dict[str, Any]
    _hub_mixin_jsonable_custom_types: tuple[Type, ...]
    _hub_mixin_coders: dict[Type, CODER_T]

    def __init_subclass__(cls, *, repo_url: Optional[str]=None, paper_url: Optional[str]=None, docs_url: Optional[str]=None, model_card_template: str=DEFAULT_MODEL_CARD, language: Optional[list[str]]=None, library_name: Optional[str]=None, license: Optional[str]=None, license_name: Optional[str]=None, license_link: Optional[str]=None, pipeline_tag: Optional[str]=None, tags: Optional[list[str]]=None, coders: Optional[dict[Type, CODER_T]]=None) -> None:
        super().__init_subclass__()
        tags = tags or []
        tags.append('model_hub_mixin')
        info = MixinInfo(model_card_template=model_card_template, model_card_data=ModelCardData())
        if hasattr(cls, '_hub_mixin_info'):
            if model_card_template == DEFAULT_MODEL_CARD:
                info.model_card_template = cls._hub_mixin_info.model_card_template
            info.model_card_data = ModelCardData(**cls._hub_mixin_info.model_card_data.to_dict())
            info.docs_url = cls._hub_mixin_info.docs_url
            info.paper_url = cls._hub_mixin_info.paper_url
            info.repo_url = cls._hub_mixin_info.repo_url
        cls._hub_mixin_info = info
        if model_card_template is not None and model_card_template != DEFAULT_MODEL_CARD:
            info.model_card_template = model_card_template
        if repo_url is not None:
            info.repo_url = repo_url
        if paper_url is not None:
            info.paper_url = paper_url
        if docs_url is not None:
            info.docs_url = docs_url
        if language is not None:
            info.model_card_data.language = language
        if library_name is not None:
            info.model_card_data.library_name = library_name
        if license is not None:
            info.model_card_data.license = license
        if license_name is not None:
            info.model_card_data.license_name = license_name
        if license_link is not None:
            info.model_card_data.license_link = license_link
        if pipeline_tag is not None:
            info.model_card_data.pipeline_tag = pipeline_tag
        if tags is not None:
            normalized_tags = list(tags)
            if info.model_card_data.tags is not None:
                info.model_card_data.tags.extend(normalized_tags)
            else:
                info.model_card_data.tags = normalized_tags
        if info.model_card_data.tags is not None:
            info.model_card_data.tags = sorted(set(info.model_card_data.tags))
        cls._hub_mixin_coders = coders or {}
        cls._hub_mixin_jsonable_custom_types = tuple(cls._hub_mixin_coders.keys())
        cls._hub_mixin_init_parameters = dict(inspect.signature(cls.__init__).parameters)
        cls._hub_mixin_jsonable_default_values = {param.name: cls._encode_arg(param.default) for param in cls._hub_mixin_init_parameters.values() if param.default is not inspect.Parameter.empty and cls._is_jsonable(param.default)}
        cls._hub_mixin_inject_config = 'config' in inspect.signature(cls._from_pretrained).parameters

    def __new__(cls: type[T], *args, **kwargs) -> T:
        instance = super().__new__(cls)
        if instance._hub_mixin_config is not None:
            return instance
        passed_values = {**{key: value for key, value in zip(list(cls._hub_mixin_init_parameters)[1:], args)}, **kwargs}
        if is_dataclass(passed_values.get('config')):
            instance._hub_mixin_config = passed_values['config']
            return instance
        init_config = {**cls._hub_mixin_jsonable_default_values, **{key: cls._encode_arg(value) for key, value in passed_values.items() if instance._is_jsonable(value)}}
        passed_config = init_config.pop('config', {})
        if isinstance(passed_config, dict):
            init_config.update(passed_config)
        if init_config != {}:
            instance._hub_mixin_config = init_config
        return instance

    @classmethod
    def _is_jsonable(cls, value: Any) -> bool:
        if is_dataclass(value):
            return True
        if isinstance(value, cls._hub_mixin_jsonable_custom_types):
            return True
        return is_jsonable(value)

    @classmethod
    def _encode_arg(cls, arg: Any) -> Any:
        if is_dataclass(arg):
            return asdict(arg)
        for type_, (encoder, _) in cls._hub_mixin_coders.items():
            if isinstance(arg, type_):
                if arg is None:
                    return None
                return encoder(arg)
        return arg

    @classmethod
    def _decode_arg(cls, expected_type: type[ARGS_T], value: Any) -> Optional[ARGS_T]:
        if is_simple_optional_type(expected_type):
            if value is None:
                return None
            expected_type = unwrap_simple_optional_type(expected_type)
        if is_dataclass(expected_type):
            return _load_dataclass(expected_type, value)
        for type_, (_, decoder) in cls._hub_mixin_coders.items():
            if inspect.isclass(expected_type) and issubclass(expected_type, type_):
                return decoder(value)
        return value

    def save_pretrained(self, save_directory: Union[str, Path], *, config: Optional[Union[dict, DataclassInstance]]=None, repo_id: Optional[str]=None, push_to_hub: bool=False, model_card_kwargs: Optional[dict[str, Any]]=None, **push_to_hub_kwargs) -> Optional[str]:
        save_directory = Path(save_directory)
        save_directory.mkdir(parents=True, exist_ok=True)
        config_path = save_directory / constants.CONFIG_NAME
        config_path.unlink(missing_ok=True)
        self._save_pretrained(save_directory)
        if config is None:
            config = self._hub_mixin_config
        if config is not None:
            if is_dataclass(config):
                config = asdict(config)
            if not config_path.exists():
                config_str = json.dumps(config, sort_keys=True, indent=2)
                config_path.write_text(config_str)
        model_card_path = save_directory / 'README.md'
        model_card_kwargs = model_card_kwargs if model_card_kwargs is not None else {}
        if not model_card_path.exists():
            self.generate_model_card(**model_card_kwargs).save(save_directory / 'README.md')
        if push_to_hub:
            kwargs = push_to_hub_kwargs.copy()
            if config is not None:
                kwargs['config'] = config
            if repo_id is None:
                repo_id = save_directory.name
            return self.push_to_hub(repo_id=repo_id, model_card_kwargs=model_card_kwargs, **kwargs)
        return None

    def _save_pretrained(self, save_directory: Path) -> None:
        raise NotImplementedError

    @classmethod
    @validate_hf_hub_args
    def from_pretrained(cls: type[T], pretrained_model_name_or_path: Union[str, Path], *, force_download: bool=False, token: Optional[Union[str, bool]]=None, cache_dir: Optional[Union[str, Path]]=None, local_files_only: bool=False, revision: Optional[str]=None, **model_kwargs) -> T:
        model_id = str(pretrained_model_name_or_path)
        config_file: Optional[str] = None
        if os.path.isdir(model_id):
            if constants.CONFIG_NAME in os.listdir(model_id):
                config_file = os.path.join(model_id, constants.CONFIG_NAME)
            else:
                logger.warning(f'{constants.CONFIG_NAME} not found in {Path(model_id).resolve()}')
        else:
            try:
                config_file = hf_hub_download(repo_id=model_id, filename=constants.CONFIG_NAME, revision=revision, cache_dir=cache_dir, force_download=force_download, token=token, local_files_only=local_files_only)
            except HfHubHTTPError as e:
                logger.info(f'{constants.CONFIG_NAME} not found on the HuggingFace Hub: {str(e)}')
        config = None
        if config_file is not None:
            with open(config_file, 'r', encoding='utf-8') as f:
                config = json.load(f)
            for key, value in config.items():
                if key in cls._hub_mixin_init_parameters:
                    expected_type = cls._hub_mixin_init_parameters[key].annotation
                    if expected_type is not inspect.Parameter.empty:
                        config[key] = cls._decode_arg(expected_type, value)
            for param in cls._hub_mixin_init_parameters.values():
                if param.name not in model_kwargs and param.name in config:
                    model_kwargs[param.name] = config[param.name]
            if 'config' in cls._hub_mixin_init_parameters and 'config' not in model_kwargs:
                config_annotation = cls._hub_mixin_init_parameters['config'].annotation
                config = cls._decode_arg(config_annotation, config)
                model_kwargs['config'] = config
            if is_dataclass(cls):
                for key in cls.__dataclass_fields__:
                    if key not in model_kwargs and key in config:
                        model_kwargs[key] = config[key]
            elif any((param.kind == inspect.Parameter.VAR_KEYWORD for param in cls._hub_mixin_init_parameters.values())):
                for key, value in config.items():
                    if key not in model_kwargs:
                        model_kwargs[key] = value
            if cls._hub_mixin_inject_config and 'config' not in model_kwargs:
                model_kwargs['config'] = config
        instance = cls._from_pretrained(model_id=str(model_id), revision=revision, cache_dir=cache_dir, force_download=force_download, local_files_only=local_files_only, token=token, **model_kwargs)
        if config is not None and getattr(instance, '_hub_mixin_config', None) in (None, {}):
            instance._hub_mixin_config = config
        return instance

    @classmethod
    def _from_pretrained(cls: type[T], *, model_id: str, revision: Optional[str], cache_dir: Optional[Union[str, Path]], force_download: bool, local_files_only: bool, token: Optional[Union[str, bool]], **model_kwargs) -> T:
        raise NotImplementedError

    @validate_hf_hub_args
    def push_to_hub(self, repo_id: str, *, config: Optional[Union[dict, DataclassInstance]]=None, commit_message: str='Push model using huggingface_hub.', private: Optional[bool]=None, token: Optional[str]=None, branch: Optional[str]=None, create_pr: Optional[bool]=None, allow_patterns: Optional[Union[list[str], str]]=None, ignore_patterns: Optional[Union[list[str], str]]=None, delete_patterns: Optional[Union[list[str], str]]=None, model_card_kwargs: Optional[dict[str, Any]]=None) -> str:
        api = HfApi(token=token)
        repo_id = api.create_repo(repo_id=repo_id, private=private, exist_ok=True).repo_id
        with SoftTemporaryDirectory() as tmp:
            saved_path = Path(tmp) / repo_id
            self.save_pretrained(saved_path, config=config, model_card_kwargs=model_card_kwargs)
            return api.upload_folder(repo_id=repo_id, repo_type='model', folder_path=saved_path, commit_message=commit_message, revision=branch, create_pr=create_pr, allow_patterns=allow_patterns, ignore_patterns=ignore_patterns, delete_patterns=delete_patterns)

    def generate_model_card(self, *args, **kwargs) -> ModelCard:
        card = ModelCard.from_template(card_data=self._hub_mixin_info.model_card_data, template_str=self._hub_mixin_info.model_card_template, repo_url=self._hub_mixin_info.repo_url, paper_url=self._hub_mixin_info.paper_url, docs_url=self._hub_mixin_info.docs_url, **kwargs)
        return card

class PyTorchModelHubMixin(ModelHubMixin):

    def __init_subclass__(cls, *args, tags: Optional[list[str]]=None, **kwargs) -> None:
        tags = tags or []
        tags.append('pytorch_model_hub_mixin')
        kwargs['tags'] = tags
        return super().__init_subclass__(*args, **kwargs)

    def _save_pretrained(self, save_directory: Path) -> None:
        model_to_save = self.module if hasattr(self, 'module') else self
        save_model_as_safetensor(model_to_save, str(save_directory / constants.SAFETENSORS_SINGLE_FILE))

    @classmethod
    def _from_pretrained(cls, *, model_id: str, revision: Optional[str], cache_dir: Optional[Union[str, Path]], force_download: bool, local_files_only: bool, token: Union[str, bool, None], map_location: str='cpu', strict: bool=False, **model_kwargs):
        model = cls(**model_kwargs)
        if os.path.isdir(model_id):
            print('Loading weights from local directory')
            model_file = os.path.join(model_id, constants.SAFETENSORS_SINGLE_FILE)
            return cls._load_as_safetensor(model, model_file, map_location, strict)
        else:
            try:
                model_file = hf_hub_download(repo_id=model_id, filename=constants.SAFETENSORS_SINGLE_FILE, revision=revision, cache_dir=cache_dir, force_download=force_download, token=token, local_files_only=local_files_only)
                return cls._load_as_safetensor(model, model_file, map_location, strict)
            except EntryNotFoundError:
                model_file = hf_hub_download(repo_id=model_id, filename=constants.PYTORCH_WEIGHTS_NAME, revision=revision, cache_dir=cache_dir, force_download=force_download, token=token, local_files_only=local_files_only)
                return cls._load_as_pickle(model, model_file, map_location, strict)

    @classmethod
    def _load_as_pickle(cls, model: T, model_file: str, map_location: str, strict: bool) -> T:
        state_dict = torch.load(model_file, map_location=torch.device(map_location), weights_only=True)
        model.load_state_dict(state_dict, strict=strict)
        model.eval()
        return model

    @classmethod
    def _load_as_safetensor(cls, model: T, model_file: str, map_location: str, strict: bool) -> T:
        if packaging.version.parse(safetensors.__version__) < packaging.version.parse('0.4.3'):
            load_model_as_safetensor(model, model_file, strict=strict)
            if map_location != 'cpu':
                logger.warning("Loading model weights on other devices than 'cpu' is not supported natively in your version of safetensors. This means that the model is loaded on 'cpu' first and then copied to the device. This leads to a slower loading time. Please update safetensors to version 0.4.3 or above for improved performance.")
                model.to(map_location)
        else:
            safetensors.torch.load_model(model, model_file, strict=strict, device=map_location)
        return model

def _load_dataclass(datacls: type[DataclassInstance], data: dict) -> DataclassInstance:
    return datacls(**{k: v for k, v in data.items() if k in datacls.__dataclass_fields__})
