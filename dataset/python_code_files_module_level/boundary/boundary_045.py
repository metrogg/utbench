import inspect
import os
import random
import shutil
import tempfile
import weakref
from functools import wraps
from pathlib import Path
from typing import TYPE_CHECKING, Any, Callable, Optional, Union
import numpy as np
import xxhash
from . import config
from .naming import INVALID_WINDOWS_CHARACTERS_IN_PATH
from .utils._dill import dumps
from .utils.logging import get_logger
if TYPE_CHECKING:
    from .arrow_dataset import Dataset
logger = get_logger(__name__)
_CACHING_ENABLED = True
_TEMP_DIR_FOR_TEMP_CACHE_FILES: Optional['_TempCacheDir'] = None
_DATASETS_WITH_TABLE_IN_TEMP_DIR: Optional[weakref.WeakSet] = None

class _TempCacheDir:

    def __init__(self):
        tmpdir = os.environ.get('TMPDIR') or os.environ.get('TEMP') or os.environ.get('TMP')
        if tmpdir:
            tmpdir = os.path.normpath(tmpdir)
            if not os.path.exists(tmpdir):
                try:
                    os.makedirs(tmpdir, exist_ok=True)
                    logger.info(f'Created TMPDIR directory: {tmpdir}')
                except OSError as e:
                    raise OSError(f"TMPDIR is set to '{tmpdir}' but the directory does not exist and could not be created: {e}. Please create it manually or unset TMPDIR to fall back to the default temporary directory.") from e
            elif not os.path.isdir(tmpdir):
                raise OSError(f"TMPDIR is set to '{tmpdir}' but it is not a directory. Please point TMPDIR to a writable directory or unset it to fall back to the default temporary directory.")
        self.name = tempfile.mkdtemp(prefix=config.TEMP_CACHE_DIR_PREFIX, dir=tmpdir)
        self._finalizer = weakref.finalize(self, self._cleanup)

    def _cleanup(self):
        for dset in get_datasets_with_cache_file_in_temp_dir():
            dset.__del__()
        if os.path.exists(self.name):
            try:
                shutil.rmtree(self.name)
            except Exception as e:
                raise OSError(f'An error occurred while trying to delete temporary cache directory {self.name}. Please delete it manually.') from e

    def cleanup(self):
        if self._finalizer.detach():
            self._cleanup()

def maybe_register_dataset_for_temp_dir_deletion(dataset):
    if _TEMP_DIR_FOR_TEMP_CACHE_FILES is None:
        return
    global _DATASETS_WITH_TABLE_IN_TEMP_DIR
    if _DATASETS_WITH_TABLE_IN_TEMP_DIR is None:
        _DATASETS_WITH_TABLE_IN_TEMP_DIR = weakref.WeakSet()
    if any((Path(_TEMP_DIR_FOR_TEMP_CACHE_FILES.name) in Path(cache_file['filename']).parents for cache_file in dataset.cache_files)):
        _DATASETS_WITH_TABLE_IN_TEMP_DIR.add(dataset)

def get_datasets_with_cache_file_in_temp_dir():
    return list(_DATASETS_WITH_TABLE_IN_TEMP_DIR) if _DATASETS_WITH_TABLE_IN_TEMP_DIR is not None else []

def enable_caching():
    global _CACHING_ENABLED
    _CACHING_ENABLED = True

def disable_caching():
    global _CACHING_ENABLED
    _CACHING_ENABLED = False

def is_caching_enabled() -> bool:
    global _CACHING_ENABLED
    return bool(_CACHING_ENABLED)

def get_temporary_cache_files_directory() -> str:
    global _TEMP_DIR_FOR_TEMP_CACHE_FILES
    if _TEMP_DIR_FOR_TEMP_CACHE_FILES is None:
        _TEMP_DIR_FOR_TEMP_CACHE_FILES = _TempCacheDir()
    return _TEMP_DIR_FOR_TEMP_CACHE_FILES.name

class Hasher:
    dispatch: dict = {}

    def __init__(self):
        self.m = xxhash.xxh64()

    @classmethod
    def hash_bytes(cls, value: Union[bytes, list[bytes]]) -> str:
        value = [value] if isinstance(value, bytes) else value
        m = xxhash.xxh64()
        for x in value:
            m.update(x)
        return m.hexdigest()

    @classmethod
    def hash(cls, value: Any) -> str:
        return cls.hash_bytes(dumps(value))

    def update(self, value: Any) -> None:
        header_for_update = f'=={type(value)}=='
        value_for_update = self.hash(value)
        self.m.update(header_for_update.encode('utf8'))
        self.m.update(value_for_update.encode('utf-8'))

    def hexdigest(self) -> str:
        return self.m.hexdigest()
fingerprint_rng = random.Random()
fingerprint_warnings: dict[str, bool] = {}

def generate_fingerprint(dataset: 'Dataset') -> str:
    state = dataset.__dict__
    hasher = Hasher()
    for key in sorted(state):
        if key == '_fingerprint':
            continue
        hasher.update(key)
        hasher.update(state[key])
    for cache_file in dataset.cache_files:
        hasher.update(os.path.getmtime(cache_file['filename']))
    return hasher.hexdigest()

def generate_random_fingerprint(nbits: int=64) -> str:
    return f'{fingerprint_rng.getrandbits(nbits):0{nbits // 4}x}'

def update_fingerprint(fingerprint, transform, transform_args):
    global fingerprint_warnings
    hasher = Hasher()
    hasher.update(fingerprint)
    try:
        hasher.update(transform)
    except:
        if _CACHING_ENABLED:
            if not fingerprint_warnings.get('update_fingerprint_transform_hash_failed', False):
                logger.warning(f"Transform {transform} couldn't be hashed properly, a random hash was used instead. Make sure your transforms and parameters are serializable with pickle or dill for the dataset fingerprinting and caching to work. If you reuse this transform, the caching mechanism will consider it to be different from the previous calls and recompute everything. This warning is only shown once. Subsequent hashing failures won't be shown.")
                fingerprint_warnings['update_fingerprint_transform_hash_failed'] = True
            else:
                logger.info(f"Transform {transform} couldn't be hashed properly, a random hash was used instead.")
        else:
            logger.info(f"Transform {transform} couldn't be hashed properly, a random hash was used instead. This doesn't affect caching since it's disabled.")
        return generate_random_fingerprint()
    for key in sorted(transform_args):
        hasher.update(key)
        try:
            hasher.update(transform_args[key])
        except:
            if _CACHING_ENABLED:
                if not fingerprint_warnings.get('update_fingerprint_transform_hash_failed', False):
                    logger.warning(f"Parameter '{key}'={transform_args[key]} of the transform {transform} couldn't be hashed properly, a random hash was used instead. Make sure your transforms and parameters are serializable with pickle or dill for the dataset fingerprinting and caching to work. If you reuse this transform, the caching mechanism will consider it to be different from the previous calls and recompute everything. This warning is only shown once. Subsequent hashing failures won't be shown.")
                    fingerprint_warnings['update_fingerprint_transform_hash_failed'] = True
                else:
                    logger.info(f"Parameter '{key}'={transform_args[key]} of the transform {transform} couldn't be hashed properly, a random hash was used instead.")
            else:
                logger.info(f"Parameter '{key}'={transform_args[key]} of the transform {transform} couldn't be hashed properly, a random hash was used instead. This doesn't affect caching since it's disabled.")
            return generate_random_fingerprint()
    return hasher.hexdigest()

def validate_fingerprint(fingerprint: str, max_length=64):
    if not isinstance(fingerprint, str) or not fingerprint:
        raise ValueError(f"Invalid fingerprint '{fingerprint}': it should be a non-empty string.")
    for invalid_char in INVALID_WINDOWS_CHARACTERS_IN_PATH:
        if invalid_char in fingerprint:
            raise ValueError(f"Invalid fingerprint. Bad characters from black list '{INVALID_WINDOWS_CHARACTERS_IN_PATH}' found in '{fingerprint}'. They could create issues when creating cache files.")
    if len(fingerprint) > max_length:
        raise ValueError(f"Invalid fingerprint. Maximum lenth is {max_length} but '{fingerprint}' has length {len(fingerprint)}.It could create issues when creating cache files.")

def format_transform_for_fingerprint(func: Callable, version: Optional[str]=None) -> str:
    transform = f'{func.__module__}.{func.__qualname__}'
    if version is not None:
        transform += f'@{version}'
    return transform

def format_kwargs_for_fingerprint(func: Callable, args: tuple, kwargs: dict[str, Any], use_kwargs: Optional[list[str]]=None, ignore_kwargs: Optional[list[str]]=None, randomized_function: bool=False) -> dict[str, Any]:
    kwargs_for_fingerprint = kwargs.copy()
    if args:
        params = [p.name for p in inspect.signature(func).parameters.values() if p != p.VAR_KEYWORD]
        args = args[1:]
        params = params[1:]
        kwargs_for_fingerprint.update(zip(params, args))
    else:
        del kwargs_for_fingerprint[next(iter(inspect.signature(func).parameters))]
    if use_kwargs:
        kwargs_for_fingerprint = {k: v for k, v in kwargs_for_fingerprint.items() if k in use_kwargs}
    if ignore_kwargs:
        kwargs_for_fingerprint = {k: v for k, v in kwargs_for_fingerprint.items() if k not in ignore_kwargs}
    if randomized_function:
        if kwargs_for_fingerprint.get('seed') is None and kwargs_for_fingerprint.get('generator') is None:
            _, seed, pos, *_ = np.random.get_state()
            seed = seed[pos] if pos < 624 else seed[0]
            kwargs_for_fingerprint['generator'] = np.random.default_rng(seed)
    default_values = {p.name: p.default for p in inspect.signature(func).parameters.values() if p.default != inspect._empty}
    for default_varname, default_value in default_values.items():
        if default_varname in kwargs_for_fingerprint and kwargs_for_fingerprint[default_varname] == default_value:
            kwargs_for_fingerprint.pop(default_varname)
    return kwargs_for_fingerprint

def fingerprint_transform(inplace: bool, use_kwargs: Optional[list[str]]=None, ignore_kwargs: Optional[list[str]]=None, fingerprint_names: Optional[list[str]]=None, randomized_function: bool=False, version: Optional[str]=None):
    if use_kwargs is not None and (not isinstance(use_kwargs, list)):
        raise ValueError(f'use_kwargs is supposed to be a list, not {type(use_kwargs)}')
    if ignore_kwargs is not None and (not isinstance(ignore_kwargs, list)):
        raise ValueError(f'ignore_kwargs is supposed to be a list, not {type(use_kwargs)}')
    if inplace and fingerprint_names:
        raise ValueError('fingerprint_names are only used when inplace is False')
    fingerprint_names = fingerprint_names if fingerprint_names is not None else ['new_fingerprint']

    def _fingerprint(func):
        if not inplace and (not all((name in func.__code__.co_varnames for name in fingerprint_names))):
            raise ValueError(f'function {func} is missing parameters {fingerprint_names} in signature')
        if randomized_function:
            if 'seed' not in func.__code__.co_varnames:
                raise ValueError(f"'seed' must be in {func}'s signature")
            if 'generator' not in func.__code__.co_varnames:
                raise ValueError(f"'generator' must be in {func}'s signature")
        transform = format_transform_for_fingerprint(func, version=version)

        @wraps(func)
        def wrapper(*args, **kwargs):
            kwargs_for_fingerprint = format_kwargs_for_fingerprint(func, args, kwargs, use_kwargs=use_kwargs, ignore_kwargs=ignore_kwargs, randomized_function=randomized_function)
            if args:
                dataset: Dataset = args[0]
                args = args[1:]
            else:
                dataset: Dataset = kwargs.pop(next(iter(inspect.signature(func).parameters)))
            if inplace:
                new_fingerprint = update_fingerprint(dataset._fingerprint, transform, kwargs_for_fingerprint)
            else:
                for fingerprint_name in fingerprint_names:
                    if kwargs.get(fingerprint_name) is None:
                        kwargs_for_fingerprint['fingerprint_name'] = fingerprint_name
                        kwargs[fingerprint_name] = update_fingerprint(dataset._fingerprint, transform, kwargs_for_fingerprint)
                    else:
                        validate_fingerprint(kwargs[fingerprint_name])
            out = func(dataset, *args, **kwargs)
            if inplace:
                dataset._fingerprint = new_fingerprint
            return out
        wrapper._decorator_name_ = 'fingerprint'
        return wrapper
    return _fingerprint
