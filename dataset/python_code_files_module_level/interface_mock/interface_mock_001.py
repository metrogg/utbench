from functools import lru_cache
from typing import Any, Optional, Union, overload
from huggingface_hub import constants
from huggingface_hub.hf_api import InferenceProviderMapping
from huggingface_hub.inference._common import MimeBytes, RequestParameters
from huggingface_hub.inference._generated.types.chat_completion import ChatCompletionInputMessage
from huggingface_hub.utils import build_hf_headers, get_token, logging
logger = logging.get_logger(__name__)
HARDCODED_MODEL_INFERENCE_MAPPING: dict[str, dict[str, InferenceProviderMapping]] = {'cerebras': {}, 'cohere': {}, 'clarifai': {}, 'fal-ai': {}, 'fireworks-ai': {}, 'groq': {}, 'hf-inference': {}, 'hyperbolic': {}, 'nebius': {}, 'nscale': {}, 'nvidia': {}, 'ovhcloud': {}, 'replicate': {}, 'sambanova': {}, 'scaleway': {}, 'together': {}, 'wavespeed': {}, 'zai-org': {}}

@overload
def filter_none(obj: dict[str, Any]) -> dict[str, Any]:
    ...

@overload
def filter_none(obj: list[Any]) -> list[Any]:
    ...

def filter_none(obj: Union[dict[str, Any], list[Any]]) -> Union[dict[str, Any], list[Any]]:
    if isinstance(obj, dict):
        cleaned: dict[str, Any] = {}
        for k, v in obj.items():
            if v is None:
                continue
            if isinstance(v, (dict, list)):
                v = filter_none(v)
            cleaned[k] = v
        return cleaned
    if isinstance(obj, list):
        return [filter_none(v) if isinstance(v, (dict, list)) else v for v in obj]
    raise ValueError(f'Expected dict or list, got {type(obj)}')

class TaskProviderHelper:

    def __init__(self, provider: str, base_url: str, task: str) -> None:
        self.provider = provider
        self.task = task
        self.base_url = base_url

    def prepare_request(self, *, inputs: Any, parameters: dict[str, Any], headers: dict, model: Optional[str], api_key: Optional[str], extra_payload: Optional[dict[str, Any]]=None) -> RequestParameters:
        api_key = self._prepare_api_key(api_key)
        provider_mapping_info = self._prepare_mapping_info(model)
        headers = self._prepare_headers(headers, api_key)
        url = self._prepare_url(api_key, provider_mapping_info.provider_id)
        payload = self._prepare_payload_as_dict(inputs, parameters, provider_mapping_info=provider_mapping_info)
        if payload is not None:
            payload = recursive_merge(payload, filter_none(extra_payload or {}))
        data = self._prepare_payload_as_bytes(inputs, parameters, provider_mapping_info, extra_payload)
        if payload is not None and data is not None:
            raise ValueError('Both payload and data cannot be set in the same request.')
        if payload is None and data is None:
            raise ValueError('Either payload or data must be set in the request.')
        normalized_headers = self._normalize_headers(headers, payload, data)
        return RequestParameters(url=url, task=self.task, model=provider_mapping_info.provider_id, json=payload, data=data, headers=normalized_headers)

    def get_response(self, response: Union[bytes, dict], request_params: Optional[RequestParameters]=None) -> Any:
        return response

    def _prepare_api_key(self, api_key: Optional[str]) -> str:
        if api_key is None:
            api_key = get_token()
        if api_key is None:
            raise ValueError(f'You must provide an api_key to work with {self.provider} API or log in with `hf auth login`.')
        return api_key

    def _prepare_mapping_info(self, model: Optional[str]) -> InferenceProviderMapping:
        if model is None:
            raise ValueError(f'Please provide an HF model ID supported by {self.provider}.')
        if HARDCODED_MODEL_INFERENCE_MAPPING.get(self.provider, {}).get(model):
            return HARDCODED_MODEL_INFERENCE_MAPPING[self.provider][model]
        provider_mapping = None
        for mapping in _fetch_inference_provider_mapping(model):
            if mapping.provider == self.provider:
                provider_mapping = mapping
                break
        if provider_mapping is None:
            raise ValueError(f'Model {model} is not supported by provider {self.provider}.')
        if provider_mapping.task != self.task:
            raise ValueError(f'Model {model} is not supported for task {self.task} and provider {self.provider}. Supported task: {provider_mapping.task}.')
        if provider_mapping.status == 'staging':
            logger.warning(f'Model {model} is in staging mode for provider {self.provider}. Meant for test purposes only.')
        if provider_mapping.status == 'error':
            logger.warning(f"Our latest automated health check on model '{model}' for provider '{self.provider}' did not complete successfully.  Inference call might fail.")
        return provider_mapping

    def _normalize_headers(self, headers: dict[str, Any], payload: Optional[dict[str, Any]], data: Optional[MimeBytes]) -> dict[str, Any]:
        normalized_headers = {key.lower(): value for key, value in headers.items() if value is not None}
        if normalized_headers.get('content-type') is None:
            if data is not None and data.mime_type is not None:
                normalized_headers['content-type'] = data.mime_type
            elif payload is not None:
                normalized_headers['content-type'] = 'application/json'
        return normalized_headers

    def _prepare_headers(self, headers: dict, api_key: str) -> dict[str, Any]:
        return {**build_hf_headers(token=api_key), **headers}

    def _prepare_url(self, api_key: str, mapped_model: str) -> str:
        base_url = self._prepare_base_url(api_key)
        route = self._prepare_route(mapped_model, api_key)
        return f"{base_url.rstrip('/')}/{route.lstrip('/')}"

    def _prepare_base_url(self, api_key: str) -> str:
        if api_key.startswith('hf_'):
            logger.info(f"Calling '{self.provider}' provider through Hugging Face router.")
            return constants.INFERENCE_PROXY_TEMPLATE.format(provider=self.provider)
        else:
            logger.info(f"Calling '{self.provider}' provider directly.")
            return self.base_url

    def _prepare_route(self, mapped_model: str, api_key: str) -> str:
        return ''

    def _prepare_payload_as_dict(self, inputs: Any, parameters: dict, provider_mapping_info: InferenceProviderMapping) -> Optional[dict]:
        return None

    def _prepare_payload_as_bytes(self, inputs: Any, parameters: dict, provider_mapping_info: InferenceProviderMapping, extra_payload: Optional[dict]) -> Optional[MimeBytes]:
        return None

class BaseConversationalTask(TaskProviderHelper):

    def __init__(self, provider: str, base_url: str):
        super().__init__(provider=provider, base_url=base_url, task='conversational')

    def _prepare_route(self, mapped_model: str, api_key: str) -> str:
        return '/v1/chat/completions'

    def _prepare_payload_as_dict(self, inputs: list[Union[dict, ChatCompletionInputMessage]], parameters: dict, provider_mapping_info: InferenceProviderMapping) -> Optional[dict]:
        return filter_none({'messages': inputs, **parameters, 'model': provider_mapping_info.provider_id})

class AutoRouterConversationalTask(BaseConversationalTask):

    def __init__(self):
        super().__init__(provider='auto', base_url='https://router.huggingface.co')

    def _prepare_base_url(self, api_key: str) -> str:
        if not api_key.startswith('hf_'):
            raise ValueError('Cannot select auto-router when using non-Hugging Face API key.')
        else:
            return self.base_url

    def _prepare_mapping_info(self, model: Optional[str]) -> InferenceProviderMapping:
        if model is None:
            raise ValueError('Please provide an HF model ID.')
        return InferenceProviderMapping(provider='auto', hf_model_id=model, providerId=model, status='live', task='conversational')

class BaseTextGenerationTask(TaskProviderHelper):

    def __init__(self, provider: str, base_url: str):
        super().__init__(provider=provider, base_url=base_url, task='text-generation')

    def _prepare_route(self, mapped_model: str, api_key: str) -> str:
        return '/v1/completions'

    def _prepare_payload_as_dict(self, inputs: Any, parameters: dict, provider_mapping_info: InferenceProviderMapping) -> Optional[dict]:
        return filter_none({'prompt': inputs, **parameters, 'model': provider_mapping_info.provider_id})

@lru_cache(maxsize=None)
def _fetch_inference_provider_mapping(model: str) -> list['InferenceProviderMapping']:
    from huggingface_hub.hf_api import HfApi
    info = HfApi().model_info(model, expand=['inferenceProviderMapping'])
    provider_mapping = info.inference_provider_mapping
    if provider_mapping is None:
        raise ValueError(f'No provider mapping found for model {model}')
    return provider_mapping

def recursive_merge(dict1: dict, dict2: dict) -> dict:
    return {**dict1, **{key: recursive_merge(dict1[key], value) if key in dict1 and isinstance(dict1[key], dict) and isinstance(value, dict) else value for key, value in dict2.items()}}
