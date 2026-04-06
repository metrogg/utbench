import os
import re
import typing
from typing import Literal, Optional
ENV_VARS_TRUE_VALUES = {'1', 'ON', 'YES', 'TRUE'}
ENV_VARS_TRUE_AND_AUTO_VALUES = ENV_VARS_TRUE_VALUES.union({'AUTO'})

def _is_true(value: Optional[str]) -> bool:
    if value is None:
        return False
    return value.upper() in ENV_VARS_TRUE_VALUES

def _as_int(value: Optional[str]) -> Optional[int]:
    if value is None:
        return None
    return int(value)
PYTORCH_WEIGHTS_NAME = 'pytorch_model.bin'
TF2_WEIGHTS_NAME = 'tf_model.h5'
TF_WEIGHTS_NAME = 'model.ckpt'
FLAX_WEIGHTS_NAME = 'flax_model.msgpack'
CONFIG_NAME = 'config.json'
REPOCARD_NAME = 'README.md'
EVAL_RESULTS_FOLDER = '.eval_results'
DEFAULT_ETAG_TIMEOUT = 10
DEFAULT_DOWNLOAD_TIMEOUT = 10
DEFAULT_REQUEST_TIMEOUT = 10
DOWNLOAD_CHUNK_SIZE = 10 * 1024 * 1024
MAX_HTTP_DOWNLOAD_SIZE = 50 * 1000 * 1000 * 1000
PYTORCH_WEIGHTS_FILE_PATTERN = 'pytorch_model{suffix}.bin'
SAFETENSORS_WEIGHTS_FILE_PATTERN = 'model{suffix}.safetensors'
TF2_WEIGHTS_FILE_PATTERN = 'tf_model{suffix}.h5'
SAFETENSORS_SINGLE_FILE = 'model.safetensors'
SAFETENSORS_INDEX_FILE = 'model.safetensors.index.json'
SAFETENSORS_MAX_HEADER_LENGTH = 25000000
FILELOCK_LOG_EVERY_SECONDS = 10
DEFAULT_REVISION = 'main'
REGEX_COMMIT_OID = re.compile('[A-Fa-f0-9]{5,40}')
HUGGINGFACE_CO_URL_HOME = 'https://huggingface.co/'
_staging_mode = _is_true(os.environ.get('HUGGINGFACE_CO_STAGING'))
_HF_DEFAULT_ENDPOINT = 'https://huggingface.co'
_HF_DEFAULT_STAGING_ENDPOINT = 'https://hub-ci.huggingface.co'
ENDPOINT = os.getenv('HF_ENDPOINT', _HF_DEFAULT_ENDPOINT).rstrip('/')
HUGGINGFACE_CO_URL_TEMPLATE = ENDPOINT + '/{repo_id}/resolve/{revision}/{filename}'
if _staging_mode:
    ENDPOINT = _HF_DEFAULT_STAGING_ENDPOINT
    HUGGINGFACE_CO_URL_TEMPLATE = _HF_DEFAULT_STAGING_ENDPOINT + '/{repo_id}/resolve/{revision}/{filename}'
DATASETS_SERVER_ENDPOINT = 'https://datasets-server.huggingface.co'
HUGGINGFACE_HEADER_X_REPO_COMMIT = 'X-Repo-Commit'
HUGGINGFACE_HEADER_X_LINKED_ETAG = 'X-Linked-Etag'
HUGGINGFACE_HEADER_X_LINKED_SIZE = 'X-Linked-Size'
HUGGINGFACE_HEADER_X_BILL_TO = 'X-HF-Bill-To'
INFERENCE_ENDPOINT = os.environ.get('HF_INFERENCE_ENDPOINT', 'https://api-inference.huggingface.co')
INFERENCE_ENDPOINTS_ENDPOINT = 'https://api.endpoints.huggingface.cloud/v2'
INFERENCE_CATALOG_ENDPOINT = 'https://endpoints.huggingface.co/api/catalog'
INFERENCE_ENDPOINT_IMAGE_KEYS = ['custom', 'huggingface', 'huggingfaceNeuron', 'llamacpp', 'tei', 'tgi', 'tgiNeuron']
INFERENCE_PROXY_TEMPLATE = 'https://router.huggingface.co/{provider}'
REPO_ID_SEPARATOR = '--'
REPO_TYPE_DATASET = 'dataset'
REPO_TYPE_SPACE = 'space'
REPO_TYPE_MODEL = 'model'
REPO_TYPES = [None, REPO_TYPE_MODEL, REPO_TYPE_DATASET, REPO_TYPE_SPACE]
SPACES_SDK_TYPES = ['gradio', 'streamlit', 'docker', 'static']
REPO_TYPES_URL_PREFIXES = {REPO_TYPE_DATASET: 'datasets/', REPO_TYPE_SPACE: 'spaces/'}
REPO_TYPES_MAPPING = {'datasets': REPO_TYPE_DATASET, 'spaces': REPO_TYPE_SPACE, 'models': REPO_TYPE_MODEL}
DiscussionTypeFilter = Literal['all', 'discussion', 'pull_request']
DISCUSSION_TYPES: tuple[DiscussionTypeFilter, ...] = typing.get_args(DiscussionTypeFilter)
DiscussionStatusFilter = Literal['all', 'open', 'closed']
DISCUSSION_STATUS: tuple[DiscussionTypeFilter, ...] = typing.get_args(DiscussionStatusFilter)
WEBHOOK_DOMAIN_T = Literal['repo', 'discussions']
default_home = os.path.join(os.path.expanduser('~'), '.cache')
HF_HOME = os.path.expandvars(os.path.expanduser(os.getenv('HF_HOME', os.path.join(os.getenv('XDG_CACHE_HOME', default_home), 'huggingface'))))
default_cache_path = os.path.join(HF_HOME, 'hub')
default_assets_cache_path = os.path.join(HF_HOME, 'assets')
HUGGINGFACE_HUB_CACHE = os.getenv('HUGGINGFACE_HUB_CACHE', default_cache_path)
HUGGINGFACE_ASSETS_CACHE = os.getenv('HUGGINGFACE_ASSETS_CACHE', default_assets_cache_path)
HF_HUB_CACHE = os.path.expandvars(os.path.expanduser(os.getenv('HF_HUB_CACHE', HUGGINGFACE_HUB_CACHE)))
HF_ASSETS_CACHE = os.path.expandvars(os.path.expanduser(os.getenv('HF_ASSETS_CACHE', HUGGINGFACE_ASSETS_CACHE)))
HF_HUB_OFFLINE = _is_true(os.environ.get('HF_HUB_OFFLINE') or os.environ.get('TRANSFORMERS_OFFLINE'))

def is_offline_mode() -> bool:
    return HF_HUB_OFFLINE
CHECK_FOR_UPDATE_DONE_PATH = os.path.join(HF_HOME, '.check_for_update_done')
HF_DEBUG = _is_true(os.environ.get('HF_DEBUG'))
HF_HUB_DISABLE_TELEMETRY = _is_true(os.environ.get('HF_HUB_DISABLE_TELEMETRY')) or _is_true(os.environ.get('DISABLE_TELEMETRY')) or _is_true(os.environ.get('DO_NOT_TRACK'))
HF_TOKEN_PATH = os.path.expandvars(os.path.expanduser(os.getenv('HF_TOKEN_PATH', os.path.join(HF_HOME, 'token'))))
HF_STORED_TOKENS_PATH = os.path.join(os.path.dirname(HF_TOKEN_PATH), 'stored_tokens')
if _staging_mode:
    _staging_home = os.path.join(os.path.expanduser('~'), '.cache', 'huggingface_staging')
    HUGGINGFACE_HUB_CACHE = os.path.join(_staging_home, 'hub')
    HF_TOKEN_PATH = os.path.join(_staging_home, 'token')
__HF_HUB_DISABLE_PROGRESS_BARS = os.environ.get('HF_HUB_DISABLE_PROGRESS_BARS')
HF_HUB_DISABLE_PROGRESS_BARS: Optional[bool] = _is_true(__HF_HUB_DISABLE_PROGRESS_BARS) if __HF_HUB_DISABLE_PROGRESS_BARS is not None else None
HF_HUB_DISABLE_SYMLINKS_WARNING: bool = _is_true(os.environ.get('HF_HUB_DISABLE_SYMLINKS_WARNING'))
HF_HUB_DISABLE_EXPERIMENTAL_WARNING: bool = _is_true(os.environ.get('HF_HUB_DISABLE_EXPERIMENTAL_WARNING'))
HF_HUB_DISABLE_IMPLICIT_TOKEN: bool = _is_true(os.environ.get('HF_HUB_DISABLE_IMPLICIT_TOKEN'))
HF_XET_HIGH_PERFORMANCE: bool = _is_true(os.environ.get('HF_XET_HIGH_PERFORMANCE'))
if _is_true(os.environ.get('HF_HUB_ENABLE_HF_TRANSFER')) and (not HF_XET_HIGH_PERFORMANCE):
    import warnings
    warnings.warn("The `HF_HUB_ENABLE_HF_TRANSFER` environment variable is deprecated as 'hf_transfer' is not used anymore. Please use `HF_XET_HIGH_PERFORMANCE` instead to enable high performance transfer with Xet. Visit https://huggingface.co/docs/huggingface_hub/package_reference/environment_variables#hfxethighperformance for more details.", DeprecationWarning)
HF_HUB_ETAG_TIMEOUT: int = _as_int(os.environ.get('HF_HUB_ETAG_TIMEOUT')) or DEFAULT_ETAG_TIMEOUT
HF_HUB_DOWNLOAD_TIMEOUT: int = _as_int(os.environ.get('HF_HUB_DOWNLOAD_TIMEOUT')) or DEFAULT_DOWNLOAD_TIMEOUT
HF_HUB_USER_AGENT_ORIGIN: Optional[str] = os.environ.get('HF_HUB_USER_AGENT_ORIGIN')
OAUTH_MAX_REDIRECTS = 2
OAUTH_CLIENT_ID = os.environ.get('OAUTH_CLIENT_ID')
OAUTH_CLIENT_SECRET = os.environ.get('OAUTH_CLIENT_SECRET')
OAUTH_SCOPES = os.environ.get('OAUTH_SCOPES')
OPENID_PROVIDER_URL = os.environ.get('OPENID_PROVIDER_URL')
HUGGINGFACE_HEADER_X_XET_ENDPOINT = 'X-Xet-Cas-Url'
HUGGINGFACE_HEADER_X_XET_ACCESS_TOKEN = 'X-Xet-Access-Token'
HUGGINGFACE_HEADER_X_XET_EXPIRATION = 'X-Xet-Token-Expiration'
HUGGINGFACE_HEADER_X_XET_HASH = 'X-Xet-Hash'
HUGGINGFACE_HEADER_X_XET_REFRESH_ROUTE = 'X-Xet-Refresh-Route'
HUGGINGFACE_HEADER_LINK_XET_AUTH_KEY = 'xet-auth'
default_xet_cache_path = os.path.join(HF_HOME, 'xet')
HF_XET_CACHE = os.getenv('HF_XET_CACHE', default_xet_cache_path)
HF_HUB_DISABLE_XET: bool = _is_true(os.environ.get('HF_HUB_DISABLE_XET'))
