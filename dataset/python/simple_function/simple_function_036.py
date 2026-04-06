import importlib.metadata
import os
import platform
import sys
import warnings
from pathlib import Path
from typing import Any, Literal
from .. import __version__, constants
_PY_VERSION: str = sys.version.split()[0].rstrip('+')
_package_versions = {}
_CANDIDATES = {'aiohttp': {'aiohttp'}, 'fastai': {'fastai'}, 'fastapi': {'fastapi'}, 'fastcore': {'fastcore'}, 'gradio': {'gradio'}, 'graphviz': {'graphviz'}, 'hf_xet': {'hf_xet'}, 'jinja': {'Jinja2'}, 'httpx': {'httpx'}, 'keras': {'keras'}, 'numpy': {'numpy'}, 'pillow': {'Pillow'}, 'pydantic': {'pydantic'}, 'pydot': {'pydot'}, 'safetensors': {'safetensors'}, 'tensorboard': {'tensorboardX'}, 'tensorflow': ('tensorflow', 'tensorflow-cpu', 'tensorflow-gpu', 'tf-nightly', 'tf-nightly-cpu', 'tf-nightly-gpu', 'intel-tensorflow', 'intel-tensorflow-avx512', 'tensorflow-rocm', 'tensorflow-macos'), 'torch': {'torch'}}
for candidate_name, package_names in _CANDIDATES.items():
    _package_versions[candidate_name] = 'N/A'
    for name in package_names:
        try:
            _package_versions[candidate_name] = importlib.metadata.version(name)
            break
        except importlib.metadata.PackageNotFoundError:
            pass

def _get_version(package_name: str) -> str:
    return _package_versions.get(package_name, 'N/A')

def is_package_available(package_name: str) -> bool:
    return _get_version(package_name) != 'N/A'

def get_python_version() -> str:
    return _PY_VERSION

def get_hf_hub_version() -> str:
    return __version__

def is_aiohttp_available() -> bool:
    return is_package_available('aiohttp')

def get_aiohttp_version() -> str:
    return _get_version('aiohttp')

def is_fastai_available() -> bool:
    return is_package_available('fastai')

def get_fastai_version() -> str:
    return _get_version('fastai')

def is_fastapi_available() -> bool:
    return is_package_available('fastapi')

def get_fastapi_version() -> str:
    return _get_version('fastapi')

def is_fastcore_available() -> bool:
    return is_package_available('fastcore')

def get_fastcore_version() -> str:
    return _get_version('fastcore')

def is_gradio_available() -> bool:
    return is_package_available('gradio')

def get_gradio_version() -> str:
    return _get_version('gradio')

def is_graphviz_available() -> bool:
    return is_package_available('graphviz')

def get_graphviz_version() -> str:
    return _get_version('graphviz')

def is_httpx_available() -> bool:
    return is_package_available('httpx')

def get_httpx_version() -> str:
    return _get_version('httpx')

def is_xet_available() -> bool:
    if constants.HF_HUB_DISABLE_XET:
        return False
    return is_package_available('hf_xet')

def get_xet_version() -> str:
    return _get_version('hf_xet')

def is_keras_available() -> bool:
    return is_package_available('keras')

def get_keras_version() -> str:
    return _get_version('keras')

def is_numpy_available() -> bool:
    return is_package_available('numpy')

def get_numpy_version() -> str:
    return _get_version('numpy')

def is_jinja_available() -> bool:
    return is_package_available('jinja')

def get_jinja_version() -> str:
    return _get_version('jinja')

def is_pillow_available() -> bool:
    return is_package_available('pillow')

def get_pillow_version() -> str:
    return _get_version('pillow')

def is_pydantic_available() -> bool:
    if not is_package_available('pydantic'):
        return False
    try:
        from pydantic import validator
    except ImportError:
        warnings.warn("Pydantic is installed but cannot be imported. Please check your installation. `huggingface_hub` will default to not using Pydantic. Error message: '{e}'")
        return False
    return True

def get_pydantic_version() -> str:
    return _get_version('pydantic')

def is_pydot_available() -> bool:
    return is_package_available('pydot')

def get_pydot_version() -> str:
    return _get_version('pydot')

def is_tensorboard_available() -> bool:
    return is_package_available('tensorboard')

def get_tensorboard_version() -> str:
    return _get_version('tensorboard')

def is_tf_available() -> bool:
    return is_package_available('tensorflow')

def get_tf_version() -> str:
    return _get_version('tensorflow')

def is_torch_available() -> bool:
    return is_package_available('torch')

def get_torch_version() -> str:
    return _get_version('torch')

def is_safetensors_available() -> bool:
    return is_package_available('safetensors')
try:
    _is_google_colab = 'google.colab' in str(get_ipython())
except NameError:
    _is_google_colab = False

def is_notebook() -> bool:
    try:
        shell_class = get_ipython().__class__
        for parent_class in shell_class.__mro__:
            if parent_class.__name__ == 'ZMQInteractiveShell':
                return True
        return False
    except NameError:
        return False

def is_google_colab() -> bool:
    return _is_google_colab

def is_colab_enterprise() -> bool:
    return os.environ.get('VERTEX_PRODUCT') == 'COLAB_ENTERPRISE'

def installation_method() -> Literal['brew', 'hf_installer', 'unknown']:
    if _is_brew_installation():
        return 'brew'
    elif _is_hf_installer_installation():
        return 'hf_installer'
    else:
        return 'unknown'

def _is_brew_installation() -> bool:
    exe_path = Path(sys.executable).resolve()
    exe_str = str(exe_path)
    return '/Cellar/' in exe_str or '/opt/homebrew/' in exe_str or exe_str.startswith('/usr/local/Cellar/')

def _is_hf_installer_installation() -> bool:
    venv = sys.prefix
    marker = Path(venv) / '.hf_installer_marker'
    return marker.exists()

def dump_environment_info() -> dict[str, Any]:
    from huggingface_hub import get_token, whoami
    from huggingface_hub.utils import list_credential_helpers
    token = get_token()
    info: dict[str, Any] = {'huggingface_hub version': get_hf_hub_version(), 'Platform': platform.platform(), 'Python version': get_python_version()}
    try:
        shell_class = get_ipython().__class__
        info['Running in iPython ?'] = 'Yes'
        info['iPython shell'] = shell_class.__name__
    except NameError:
        info['Running in iPython ?'] = 'No'
    info['Running in notebook ?'] = 'Yes' if is_notebook() else 'No'
    info['Running in Google Colab ?'] = 'Yes' if is_google_colab() else 'No'
    info['Running in Google Colab Enterprise ?'] = 'Yes' if is_colab_enterprise() else 'No'
    info['Token path ?'] = constants.HF_TOKEN_PATH
    info['Has saved token ?'] = token is not None
    if token is not None:
        try:
            info['Who am I ?'] = whoami()['name']
        except Exception:
            pass
    try:
        info['Configured git credential helpers'] = ', '.join(list_credential_helpers())
    except Exception:
        pass
    info['Installation method'] = installation_method()
    info['httpx'] = get_httpx_version()
    info['hf_xet'] = get_xet_version()
    info['gradio'] = get_gradio_version()
    info['tensorboard'] = get_tensorboard_version()
    info['ENDPOINT'] = constants.ENDPOINT
    info['HF_HUB_CACHE'] = constants.HF_HUB_CACHE
    info['HF_ASSETS_CACHE'] = constants.HF_ASSETS_CACHE
    info['HF_TOKEN_PATH'] = constants.HF_TOKEN_PATH
    info['HF_STORED_TOKENS_PATH'] = constants.HF_STORED_TOKENS_PATH
    info['HF_HUB_OFFLINE'] = constants.HF_HUB_OFFLINE
    info['HF_HUB_DISABLE_TELEMETRY'] = constants.HF_HUB_DISABLE_TELEMETRY
    info['HF_HUB_DISABLE_PROGRESS_BARS'] = constants.HF_HUB_DISABLE_PROGRESS_BARS
    info['HF_HUB_DISABLE_SYMLINKS_WARNING'] = constants.HF_HUB_DISABLE_SYMLINKS_WARNING
    info['HF_HUB_DISABLE_EXPERIMENTAL_WARNING'] = constants.HF_HUB_DISABLE_EXPERIMENTAL_WARNING
    info['HF_HUB_DISABLE_IMPLICIT_TOKEN'] = constants.HF_HUB_DISABLE_IMPLICIT_TOKEN
    info['HF_HUB_DISABLE_XET'] = constants.HF_HUB_DISABLE_XET
    info['HF_HUB_ETAG_TIMEOUT'] = constants.HF_HUB_ETAG_TIMEOUT
    info['HF_HUB_DOWNLOAD_TIMEOUT'] = constants.HF_HUB_DOWNLOAD_TIMEOUT
    info['HF_XET_HIGH_PERFORMANCE'] = constants.HF_XET_HIGH_PERFORMANCE
    print('\nCopy-and-paste the text below in your GitHub issue.\n')
    print('\n'.join([f'- {prop}: {val}' for prop, val in info.items()]) + '\n')
    return info
