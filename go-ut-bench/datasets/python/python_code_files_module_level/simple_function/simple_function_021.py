from __future__ import annotations
from typing import TYPE_CHECKING, Any
import numpy as np
from pandas._libs import lib
from pandas.util._decorators import set_module
from pandas.core.dtypes.common import is_array_like, is_bool_dtype, is_integer, is_integer_dtype, is_list_like
from pandas.core.dtypes.dtypes import ExtensionDtype
from pandas.core.dtypes.generic import ABCIndex, ABCSeries
if TYPE_CHECKING:
    from pandas._typing import AnyArrayLike
    from pandas.core.frame import DataFrame
    from pandas.core.indexes.base import Index

def is_valid_positional_slice(slc: slice) -> bool:
    return lib.is_int_or_none(slc.start) and lib.is_int_or_none(slc.stop) and lib.is_int_or_none(slc.step)

def is_list_like_indexer(key) -> bool:
    return is_list_like(key) and (not (isinstance(key, tuple) and type(key) is not tuple))

def is_scalar_indexer(indexer, ndim: int) -> bool:
    if ndim == 1 and is_integer(indexer):
        return True
    if isinstance(indexer, tuple) and len(indexer) == ndim:
        return all((is_integer(x) for x in indexer))
    return False

def is_empty_indexer(indexer) -> bool:
    if is_list_like(indexer) and (not len(indexer)):
        return True
    if not isinstance(indexer, tuple):
        indexer = (indexer,)
    return any((isinstance(idx, np.ndarray) and len(idx) == 0 for idx in indexer))

def check_setitem_lengths(indexer, value, values) -> bool:
    no_op = False
    if isinstance(indexer, (np.ndarray, list)):
        if is_list_like(value):
            if len(indexer) != len(value) and values.ndim == 1:
                if isinstance(indexer, list):
                    indexer = np.array(indexer)
                if not (isinstance(indexer, np.ndarray) and indexer.dtype == np.bool_ and (indexer.sum() == len(value))):
                    raise ValueError('cannot set using a list-like indexer with a different length than the value')
            if not len(indexer):
                no_op = True
    elif isinstance(indexer, slice):
        if is_list_like(value):
            if len(value) != length_of_indexer(indexer, values) and values.ndim == 1:
                raise ValueError('cannot set using a slice indexer with a different length than the value')
            if not len(value):
                no_op = True
    return no_op

def validate_indices(indices: np.ndarray, n: int) -> None:
    if len(indices):
        min_idx = indices.min()
        if min_idx < -1:
            msg = f"'indices' contains values less than allowed ({min_idx} < -1)"
            raise ValueError(msg)
        max_idx = indices.max()
        if max_idx >= n:
            raise IndexError('indices are out-of-bounds')

def maybe_convert_indices(indices, n: int, verify: bool=True) -> np.ndarray:
    if isinstance(indices, list):
        indices = np.array(indices)
        if len(indices) == 0:
            return np.empty(0, dtype=np.intp)
    mask = indices < 0
    if mask.any():
        indices = indices.copy()
        indices[mask] += n
    if verify:
        mask = (indices >= n) | (indices < 0)
        if mask.any():
            raise IndexError('indices are out-of-bounds')
    return indices

def length_of_indexer(indexer, target=None) -> int:
    if target is not None and isinstance(indexer, slice):
        target_len = len(target)
        start = indexer.start
        stop = indexer.stop
        step = indexer.step
        if start is None:
            start = 0
        elif start < 0:
            start += target_len
        if stop is None or stop > target_len:
            stop = target_len
        elif stop < 0:
            stop += target_len
        if step is None:
            step = 1
        elif step < 0:
            start, stop = (stop + 1, start + 1)
            step = -step
        return (stop - start + step - 1) // step
    elif isinstance(indexer, (ABCSeries, ABCIndex, np.ndarray, list)):
        if isinstance(indexer, list):
            indexer = np.array(indexer)
        if indexer.dtype == bool:
            return indexer.sum()
        return len(indexer)
    elif isinstance(indexer, range):
        return (indexer.stop - indexer.start) // indexer.step
    elif not is_list_like_indexer(indexer):
        return 1
    raise AssertionError('cannot find the length of the indexer')

def disallow_ndim_indexing(result) -> None:
    if np.ndim(result) > 1:
        raise ValueError('Multi-dimensional indexing (e.g. `obj[:, None]`) is no longer supported. Convert to a numpy array before indexing instead.')

def unpack_1tuple(tup):
    if len(tup) == 1 and isinstance(tup[0], slice):
        if isinstance(tup, list):
            raise ValueError('Indexing with a single-item list containing a slice is not allowed. Pass a tuple instead.')
        return tup[0]
    return tup

def check_key_length(columns: Index, key, value: DataFrame) -> None:
    if columns.is_unique:
        if len(value.columns) != len(key):
            raise ValueError('Columns must be same length as key')
    elif len(columns.get_indexer_non_unique(key)[0]) != len(value.columns):
        raise ValueError('Columns must be same length as key')

def unpack_tuple_and_ellipses(item: tuple):
    if len(item) > 1:
        if item[0] is Ellipsis:
            item = item[1:]
        elif item[-1] is Ellipsis:
            item = item[:-1]
    if len(item) > 1:
        raise IndexError('too many indices for array.')
    item = item[0]
    return item

def getitem_returns_view(arr, key) -> bool:
    if not isinstance(key, tuple):
        key = (key,)
    key = tuple((k for k in key if k is not Ellipsis and k is not np.newaxis))
    if not key:
        return True
    if arr.ndim == 2 and lib.is_integer(key[0]):
        return True
    if all((isinstance(k, slice) for k in key)):
        return True
    return False

@set_module('pandas.api.indexers')
def check_array_indexer(array: AnyArrayLike, indexer: Any) -> Any:
    from pandas.core.construction import array as pd_array
    if is_list_like(indexer):
        if isinstance(indexer, tuple):
            return indexer
    else:
        return indexer
    if not is_array_like(indexer):
        indexer = pd_array(indexer)
        if len(indexer) == 0:
            indexer = np.array([], dtype=np.intp)
    dtype = indexer.dtype
    if is_bool_dtype(dtype):
        if isinstance(dtype, ExtensionDtype):
            indexer = indexer.to_numpy(dtype=bool, na_value=False)
        else:
            indexer = np.asarray(indexer, dtype=bool)
        if len(indexer) != len(array):
            raise IndexError(f'Boolean index has wrong length: {len(indexer)} instead of {len(array)}')
    elif is_integer_dtype(dtype):
        try:
            indexer = np.asarray(indexer, dtype=np.intp)
        except ValueError as err:
            raise ValueError('Cannot index with an integer indexer containing NA values') from err
    else:
        raise IndexError('arrays used as indices must be of integer or boolean type')
    return indexer
