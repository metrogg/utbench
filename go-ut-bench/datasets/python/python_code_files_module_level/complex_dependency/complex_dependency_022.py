import importlib
import json
import os
import re
from collections import defaultdict, namedtuple
from functools import lru_cache
from pathlib import Path
from typing import TYPE_CHECKING, Any, Iterable, NamedTuple, Optional, Union
from packaging import version
from .. import constants, logging
from ._base import MAX_SHARD_SIZE, StateDictSplit, split_state_dict_into_shards_factory
logger = logging.get_logger(__file__)
if TYPE_CHECKING:
    import torch

def save_torch_model(model: 'torch.nn.Module', save_directory: Union[str, Path], *, filename_pattern: Optional[str]=None, force_contiguous: bool=True, max_shard_size: Union[int, str]=MAX_SHARD_SIZE, metadata: Optional[dict[str, str]]=None, safe_serialization: bool=True, is_main_process: bool=True, shared_tensors_to_discard: Optional[list[str]]=None):
    save_torch_state_dict(state_dict=model.state_dict(), filename_pattern=filename_pattern, force_contiguous=force_contiguous, max_shard_size=max_shard_size, metadata=metadata, safe_serialization=safe_serialization, save_directory=save_directory, is_main_process=is_main_process, shared_tensors_to_discard=shared_tensors_to_discard)

def save_torch_state_dict(state_dict: dict[str, 'torch.Tensor'], save_directory: Union[str, Path], *, filename_pattern: Optional[str]=None, force_contiguous: bool=True, max_shard_size: Union[int, str]=MAX_SHARD_SIZE, metadata: Optional[dict[str, str]]=None, safe_serialization: bool=True, is_main_process: bool=True, shared_tensors_to_discard: Optional[list[str]]=None) -> None:
    save_directory = str(save_directory)
    if filename_pattern is None:
        filename_pattern = constants.SAFETENSORS_WEIGHTS_FILE_PATTERN if safe_serialization else constants.PYTORCH_WEIGHTS_FILE_PATTERN
    if metadata is None:
        metadata = {}
    if safe_serialization:
        try:
            from safetensors.torch import save_file as save_file_fn
        except ImportError as e:
            raise ImportError('Please install `safetensors` to use safe serialization. You can install it with `pip install safetensors`.') from e
        state_dict = _clean_state_dict_for_safetensors(state_dict, metadata, force_contiguous=force_contiguous, shared_tensors_to_discard=shared_tensors_to_discard)
    else:
        from torch import save as save_file_fn
        logger.warning('You are using unsafe serialization. Due to security reasons, it is recommended not to load pickled models from untrusted sources. If you intend to share your model, we strongly recommend using safe serialization by installing `safetensors` with `pip install safetensors`.')
    state_dict_split = split_torch_state_dict_into_shards(state_dict, filename_pattern=filename_pattern, max_shard_size=max_shard_size)
    if is_main_process:
        existing_files_regex = re.compile(filename_pattern.format(suffix='(-\\d{5}-of-\\d{5})?') + '(\\.index\\.json)?')
        for filename in os.listdir(save_directory):
            if existing_files_regex.match(filename):
                try:
                    logger.debug(f"Removing existing file '{filename}' from folder.")
                    os.remove(os.path.join(save_directory, filename))
                except Exception as e:
                    logger.warning(f"Error when trying to remove existing '{filename}' from folder: {e}. Continuing...")
    per_file_metadata = {'format': 'pt'}
    if not state_dict_split.is_sharded:
        per_file_metadata.update(metadata)
    safe_file_kwargs = {'metadata': per_file_metadata} if safe_serialization else {}
    for filename, tensors in state_dict_split.filename_to_tensors.items():
        shard = {tensor: state_dict[tensor] for tensor in tensors}
        save_file_fn(shard, os.path.join(save_directory, filename), **safe_file_kwargs)
        logger.debug(f'Shard saved to {filename}')
    if state_dict_split.is_sharded:
        index_path = filename_pattern.format(suffix='') + '.index.json'
        index = {'metadata': {**state_dict_split.metadata, **metadata}, 'weight_map': state_dict_split.tensor_to_filename}
        with open(os.path.join(save_directory, index_path), 'w') as f:
            json.dump(index, f, indent=2)
        logger.info(f'The model is bigger than the maximum size per checkpoint ({max_shard_size}). Model weighs have been saved in {len(state_dict_split.filename_to_tensors)} checkpoint shards. You can find where each parameters has been saved in the index located at {index_path}.')
    logger.info(f'Model weights successfully saved to {save_directory}!')

def split_torch_state_dict_into_shards(state_dict: dict[str, 'torch.Tensor'], *, filename_pattern: str=constants.SAFETENSORS_WEIGHTS_FILE_PATTERN, max_shard_size: Union[int, str]=MAX_SHARD_SIZE) -> StateDictSplit:
    return split_state_dict_into_shards_factory(state_dict, max_shard_size=max_shard_size, filename_pattern=filename_pattern, get_storage_size=get_torch_storage_size, get_storage_id=get_torch_storage_id)

def load_torch_model(model: 'torch.nn.Module', checkpoint_path: Union[str, os.PathLike], *, strict: bool=False, safe: bool=True, weights_only: bool=False, map_location: Optional[Union[str, 'torch.device']]=None, mmap: bool=False, filename_pattern: Optional[str]=None) -> NamedTuple:
    checkpoint_path = Path(checkpoint_path)
    if not checkpoint_path.exists():
        raise ValueError(f'Checkpoint path {checkpoint_path} does not exist')
    if checkpoint_path.is_file():
        state_dict = load_state_dict_from_file(checkpoint_file=checkpoint_path, map_location=map_location, weights_only=weights_only)
        return model.load_state_dict(state_dict, strict=strict)
    if filename_pattern is None:
        filename_pattern = constants.SAFETENSORS_WEIGHTS_FILE_PATTERN
        index_path = checkpoint_path / (filename_pattern.format(suffix='') + '.index.json')
        if not index_path.is_file() and (not safe):
            filename_pattern = constants.PYTORCH_WEIGHTS_FILE_PATTERN
    index_path = checkpoint_path / (filename_pattern.format(suffix='') + '.index.json')
    if index_path.is_file():
        return _load_sharded_checkpoint(model=model, save_directory=checkpoint_path, strict=strict, weights_only=weights_only, filename_pattern=filename_pattern)
    model_files = list(checkpoint_path.glob('*.safetensors' if safe else '*.bin'))
    if len(model_files) == 1:
        state_dict = load_state_dict_from_file(checkpoint_file=model_files[0], map_location=map_location, weights_only=weights_only, mmap=mmap)
        return model.load_state_dict(state_dict, strict=strict)
    raise ValueError(f"Directory '{checkpoint_path}' does not contain a valid checkpoint. Expected either a sharded checkpoint with an index file, or a single model file.")

def _load_sharded_checkpoint(model: 'torch.nn.Module', save_directory: os.PathLike, *, strict: bool=False, weights_only: bool=False, filename_pattern: str=constants.SAFETENSORS_WEIGHTS_FILE_PATTERN) -> NamedTuple:
    index_path = filename_pattern.format(suffix='') + '.index.json'
    index_file = os.path.join(save_directory, index_path)
    with open(index_file, 'r', encoding='utf-8') as f:
        index = json.load(f)
    if strict:
        _validate_keys_for_strict_loading(model, index['weight_map'].keys())
    shard_files = list(set(index['weight_map'].values()))
    for shard_file in shard_files:
        shard_path = os.path.join(save_directory, shard_file)
        state_dict = load_state_dict_from_file(shard_path, map_location='cpu', weights_only=weights_only)
        model.load_state_dict(state_dict, strict=strict)
        del state_dict
    loaded_keys = set(index['weight_map'].keys())
    model_keys = set(model.state_dict().keys())
    return _IncompatibleKeys(missing_keys=list(model_keys - loaded_keys), unexpected_keys=list(loaded_keys - model_keys))

def load_state_dict_from_file(checkpoint_file: Union[str, os.PathLike], map_location: Optional[Union[str, 'torch.device']]=None, weights_only: bool=False, mmap: bool=False) -> Union[dict[str, 'torch.Tensor'], Any]:
    checkpoint_path = Path(checkpoint_file)
    if not checkpoint_path.is_file():
        raise FileNotFoundError(f"No checkpoint file found at '{checkpoint_path}'. Please verify the path is correct and the file has been properly downloaded.")
    if checkpoint_path.suffix == '.safetensors':
        try:
            from safetensors import safe_open
            from safetensors.torch import load_file
        except ImportError as e:
            raise ImportError('Please install `safetensors` to load safetensors checkpoint. You can install it with `pip install safetensors`.') from e
        with safe_open(checkpoint_file, framework='pt') as f:
            metadata = f.metadata()
        if metadata is not None and metadata.get('format') not in ['pt', 'mlx']:
            raise OSError(f'The safetensors archive passed at {checkpoint_file} does not contain the valid metadata. Make sure you save your model with the `save_torch_model` method.')
        device = str(map_location.type) if map_location is not None and hasattr(map_location, 'type') else map_location
        if device == 'meta':
            logger.warning('Meta device is not supported with safetensors. Falling back to CPU device.')
            device = 'cpu'
        return load_file(checkpoint_file, device=device)
    try:
        import torch
        from torch import load
    except ImportError as e:
        raise ImportError('Please install `torch` to load torch tensors. You can install it with `pip install torch`.') from e
    additional_kwargs = {}
    if version.parse(torch.__version__) >= version.parse('2.1.0'):
        additional_kwargs['mmap'] = mmap
    if version.parse(torch.__version__) >= version.parse('1.13.0'):
        additional_kwargs['weights_only'] = weights_only
    return load(checkpoint_file, map_location=map_location, **additional_kwargs)

def _validate_keys_for_strict_loading(model: 'torch.nn.Module', loaded_keys: Iterable[str]) -> None:
    loaded_keys_set = set(loaded_keys)
    model_keys = set(model.state_dict().keys())
    missing_keys = model_keys - loaded_keys_set
    unexpected_keys = loaded_keys_set - model_keys
    if missing_keys or unexpected_keys:
        error_message = f'Error(s) in loading state_dict for {model.__class__.__name__}'
        if missing_keys:
            str_missing_keys = ','.join([f'"{k}"' for k in sorted(missing_keys)])
            error_message += f'\nMissing key(s): {str_missing_keys}.'
        if unexpected_keys:
            str_unexpected_keys = ','.join([f'"{k}"' for k in sorted(unexpected_keys)])
            error_message += f'\nUnexpected key(s): {str_unexpected_keys}.'
        raise RuntimeError(error_message)

def _get_unique_id(tensor: 'torch.Tensor') -> Union[int, tuple[Any, ...]]:
    try:
        from torch.distributed.tensor import DTensor
        if isinstance(tensor, DTensor):
            local_tensor = tensor.to_local()
            return local_tensor.storage().data_ptr()
    except ImportError:
        pass
    try:
        from torch.utils._python_dispatch import is_traceable_wrapper_subclass
        if is_traceable_wrapper_subclass(tensor):
            attrs, _ = tensor.__tensor_flatten__()
            return tuple((_get_unique_id(getattr(tensor, attr)) for attr in attrs))
    except ImportError:
        pass
    if tensor.device.type == 'xla' and is_torch_tpu_available():
        import torch_xla
        unique_id = torch_xla._XLAC._xla_get_tensor_id(tensor)
    else:
        unique_id = storage_ptr(tensor)
    return unique_id

def get_torch_storage_id(tensor: 'torch.Tensor') -> Optional[tuple['torch.device', Union[int, tuple[Any, ...]], int]]:
    if tensor.device.type == 'meta':
        return None
    else:
        return (tensor.device, _get_unique_id(tensor), get_torch_storage_size(tensor))

def get_torch_storage_size(tensor: 'torch.Tensor') -> int:
    try:
        from torch.distributed.tensor import DTensor
        if isinstance(tensor, DTensor):
            return tensor.nbytes
    except ImportError:
        pass
    try:
        from torch.utils._python_dispatch import is_traceable_wrapper_subclass
        if is_traceable_wrapper_subclass(tensor):
            attrs, _ = tensor.__tensor_flatten__()
            return sum((get_torch_storage_size(getattr(tensor, attr)) for attr in attrs))
    except ImportError:
        pass
    try:
        return tensor.untyped_storage().nbytes()
    except AttributeError:
        try:
            return tensor.storage().size() * _get_dtype_size(tensor.dtype)
        except NotImplementedError:
            return tensor.nelement() * _get_dtype_size(tensor.dtype)

@lru_cache()
def is_torch_tpu_available(check_device=True):
    if importlib.util.find_spec('torch_xla') is not None:
        if check_device:
            try:
                import torch_xla.core.xla_model as xm
                _ = xm.xla_device()
                return True
            except RuntimeError:
                return False
        return True
    return False

def storage_ptr(tensor: 'torch.Tensor') -> Union[int, tuple[Any, ...]]:
    try:
        from torch.utils._python_dispatch import is_traceable_wrapper_subclass
        if is_traceable_wrapper_subclass(tensor):
            return _get_unique_id(tensor)
    except ImportError:
        pass
    try:
        return tensor.untyped_storage().data_ptr()
    except Exception:
        try:
            return tensor.storage().data_ptr()
        except NotImplementedError:
            return 0

def _clean_state_dict_for_safetensors(state_dict: dict[str, 'torch.Tensor'], metadata: dict[str, str], force_contiguous: bool=True, shared_tensors_to_discard: Optional[list[str]]=None):
    to_removes = _remove_duplicate_names(state_dict, discard_names=shared_tensors_to_discard)
    for kept_name, to_remove_group in to_removes.items():
        for to_remove in to_remove_group:
            if metadata is None:
                metadata = {}
            if to_remove not in metadata:
                metadata[to_remove] = kept_name
            del state_dict[to_remove]
    if force_contiguous:
        state_dict = {k: v.contiguous() for k, v in state_dict.items()}
    return state_dict

def _end_ptr(tensor: 'torch.Tensor') -> int:
    if tensor.nelement():
        stop = tensor.view(-1)[-1].data_ptr() + _get_dtype_size(tensor.dtype)
    else:
        stop = tensor.data_ptr()
    return stop

def _filter_shared_not_shared(tensors: list[set[str]], state_dict: dict[str, 'torch.Tensor']) -> list[set[str]]:
    filtered_tensors = []
    for shared in tensors:
        if len(shared) < 2:
            filtered_tensors.append(shared)
            continue
        areas = []
        for name in shared:
            tensor = state_dict[name]
            areas.append((tensor.data_ptr(), _end_ptr(tensor), name))
        areas.sort()
        _, last_stop, last_name = areas[0]
        filtered_tensors.append({last_name})
        for start, stop, name in areas[1:]:
            if start >= last_stop:
                filtered_tensors.append({name})
            else:
                filtered_tensors[-1].add(name)
            last_stop = stop
    return filtered_tensors

def _find_shared_tensors(state_dict: dict[str, 'torch.Tensor']) -> list[set[str]]:
    import torch
    tensors_dict = defaultdict(set)
    for k, v in state_dict.items():
        if v.device != torch.device('meta') and storage_ptr(v) != 0 and (get_torch_storage_size(v) != 0):
            tensors_dict[v.device, storage_ptr(v), get_torch_storage_size(v)].add(k)
    tensors = list(sorted(tensors_dict.values()))
    tensors = _filter_shared_not_shared(tensors, state_dict)
    return tensors

def _is_complete(tensor: 'torch.Tensor') -> bool:
    try:
        from torch.utils._python_dispatch import is_traceable_wrapper_subclass
        if is_traceable_wrapper_subclass(tensor):
            attrs, _ = tensor.__tensor_flatten__()
            return all((_is_complete(getattr(tensor, attr)) for attr in attrs))
    except ImportError:
        pass
    return tensor.data_ptr() == storage_ptr(tensor) and tensor.nelement() * _get_dtype_size(tensor.dtype) == get_torch_storage_size(tensor)

def _remove_duplicate_names(state_dict: dict[str, 'torch.Tensor'], *, preferred_names: Optional[list[str]]=None, discard_names: Optional[list[str]]=None) -> dict[str, list[str]]:
    if preferred_names is None:
        preferred_names = []
    unique_preferred_names = set(preferred_names)
    if discard_names is None:
        discard_names = []
    unique_discard_names = set(discard_names)
    shareds = _find_shared_tensors(state_dict)
    to_remove = defaultdict(list)
    for shared in shareds:
        complete_names = set([name for name in shared if _is_complete(state_dict[name])])
        if not complete_names:
            raise RuntimeError(f'Error while trying to find names to remove to save state dict, but found no suitable name to keep for saving amongst: {shared}. None is covering the entire storage. Refusing to save/load the model since you could be storing much more memory than needed. Please refer to https://huggingface.co/docs/safetensors/torch_shared_tensors for more information. Or open an issue.')
        keep_name = sorted(list(complete_names))[0]
        preferred = complete_names.difference(unique_discard_names)
        if preferred:
            keep_name = sorted(list(preferred))[0]
        if unique_preferred_names:
            preferred = unique_preferred_names.intersection(complete_names)
            if preferred:
                keep_name = sorted(list(preferred))[0]
        for name in sorted(shared):
            if name != keep_name:
                to_remove[keep_name].append(name)
    return to_remove

@lru_cache()
def _get_dtype_size(dtype: 'torch.dtype') -> int:
    import torch
    _float8_e4m3fn = getattr(torch, 'float8_e4m3fn', None)
    _float8_e5m2 = getattr(torch, 'float8_e5m2', None)
    _SIZE = {torch.int64: 8, torch.float32: 4, torch.int32: 4, torch.bfloat16: 2, torch.float16: 2, torch.int16: 2, torch.uint8: 1, torch.int8: 1, torch.bool: 1, torch.float64: 8, _float8_e4m3fn: 1, _float8_e5m2: 1}
    return _SIZE[dtype]

class _IncompatibleKeys(namedtuple('IncompatibleKeys', ['missing_keys', 'unexpected_keys'])):

    def __repr__(self) -> str:
        if not self.missing_keys and (not self.unexpected_keys):
            return '<All keys matched successfully>'
        return super().__repr__()
    __str__ = __repr__
