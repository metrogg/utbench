import functools
from typing import NamedTuple
import numpy as np
from numpy._core import overrides
from numpy._core._multiarray_umath import _array_converter, _unique_hash
from numpy.lib.array_utils import normalize_axis_index
array_function_dispatch = functools.partial(overrides.array_function_dispatch, module='numpy')
__all__ = ['ediff1d', 'intersect1d', 'isin', 'setdiff1d', 'setxor1d', 'union1d', 'unique', 'unique_all', 'unique_counts', 'unique_inverse', 'unique_values']

def _ediff1d_dispatcher(ary, to_end=None, to_begin=None):
    return (ary, to_end, to_begin)

@array_function_dispatch(_ediff1d_dispatcher)
def ediff1d(ary, to_end=None, to_begin=None):
    conv = _array_converter(ary)
    ary = conv[0].ravel()
    dtype_req = ary.dtype
    if to_begin is None and to_end is None:
        return ary[1:] - ary[:-1]
    if to_begin is None:
        l_begin = 0
    else:
        to_begin = np.asanyarray(to_begin)
        if not np.can_cast(to_begin, dtype_req, casting='same_kind'):
            raise TypeError('dtype of `to_begin` must be compatible with input `ary` under the `same_kind` rule.')
        to_begin = to_begin.ravel()
        l_begin = len(to_begin)
    if to_end is None:
        l_end = 0
    else:
        to_end = np.asanyarray(to_end)
        if not np.can_cast(to_end, dtype_req, casting='same_kind'):
            raise TypeError('dtype of `to_end` must be compatible with input `ary` under the `same_kind` rule.')
        to_end = to_end.ravel()
        l_end = len(to_end)
    l_diff = max(len(ary) - 1, 0)
    result = np.empty_like(ary, shape=l_diff + l_begin + l_end)
    if l_begin > 0:
        result[:l_begin] = to_begin
    if l_end > 0:
        result[l_begin + l_diff:] = to_end
    np.subtract(ary[1:], ary[:-1], result[l_begin:l_begin + l_diff])
    return conv.wrap(result)

def _unpack_tuple(x):
    if len(x) == 1:
        return x[0]
    else:
        return x

def _unique_dispatcher(ar, return_index=None, return_inverse=None, return_counts=None, axis=None, *, equal_nan=None, sorted=True):
    return (ar,)

@array_function_dispatch(_unique_dispatcher)
def unique(ar, return_index=False, return_inverse=False, return_counts=False, axis=None, *, equal_nan=True, sorted=True):
    ar = np.asanyarray(ar)
    if axis is None or ar.ndim == 1:
        if axis is not None:
            normalize_axis_index(axis, ar.ndim)
        ret = _unique1d(ar, return_index, return_inverse, return_counts, equal_nan=equal_nan, inverse_shape=ar.shape, axis=None, sorted=sorted)
        return _unpack_tuple(ret)
    try:
        ar = np.moveaxis(ar, axis, 0)
    except np.exceptions.AxisError:
        raise np.exceptions.AxisError(axis, ar.ndim) from None
    inverse_shape = [1] * ar.ndim
    inverse_shape[axis] = ar.shape[0]
    orig_shape, orig_dtype = (ar.shape, ar.dtype)
    ar = ar.reshape(orig_shape[0], np.prod(orig_shape[1:], dtype=np.intp))
    ar = np.ascontiguousarray(ar)
    dtype = [(f'f{i}', ar.dtype) for i in range(ar.shape[1])]
    try:
        if ar.shape[1] > 0:
            consolidated = ar.view(dtype)
        else:
            consolidated = np.empty(len(ar), dtype=dtype)
    except TypeError as e:
        msg = 'The axis argument to unique is not supported for dtype {dt}'
        raise TypeError(msg.format(dt=ar.dtype)) from e

    def reshape_uniq(uniq):
        n = len(uniq)
        uniq = uniq.view(orig_dtype)
        uniq = uniq.reshape(n, *orig_shape[1:])
        uniq = np.moveaxis(uniq, 0, axis)
        return uniq
    output = _unique1d(consolidated, return_index, return_inverse, return_counts, equal_nan=equal_nan, inverse_shape=inverse_shape, axis=axis, sorted=sorted)
    output = (reshape_uniq(output[0]),) + output[1:]
    return _unpack_tuple(output)

def _unique1d(ar, return_index=False, return_inverse=False, return_counts=False, *, equal_nan=True, inverse_shape=None, axis=None, sorted=True):
    ar = np.asanyarray(ar).flatten()
    if len(ar.shape) != 1:
        ar = np.asarray(ar).flatten()
    optional_indices = return_index or return_inverse
    if not optional_indices and (not return_counts) and (not np.ma.is_masked(ar)):
        conv = _array_converter(ar)
        ar_, = conv
        if (hash_unique := _unique_hash(ar_, equal_nan=equal_nan)) is not NotImplemented:
            if sorted:
                hash_unique.sort()
            return (conv.wrap(hash_unique),)
    if optional_indices:
        perm = ar.argsort(kind='mergesort' if return_index else 'quicksort')
        aux = ar[perm]
    else:
        ar.sort()
        aux = ar
    mask = np.empty(aux.shape, dtype=np.bool)
    mask[:1] = True
    if equal_nan and aux.shape[0] > 0 and (aux.dtype.kind in 'cfmM') and np.isnan(aux[-1]):
        if aux.dtype.kind == 'c':
            aux_firstnan = np.searchsorted(np.isnan(aux), True, side='left')
        else:
            aux_firstnan = np.searchsorted(aux, aux[-1], side='left')
        if aux_firstnan > 0:
            mask[1:aux_firstnan] = aux[1:aux_firstnan] != aux[:aux_firstnan - 1]
        mask[aux_firstnan] = True
        mask[aux_firstnan + 1:] = False
    else:
        mask[1:] = aux[1:] != aux[:-1]
    ret = (aux[mask],)
    if return_index:
        ret += (perm[mask],)
    if return_inverse:
        imask = np.cumsum(mask) - 1
        inv_idx = np.empty(mask.shape, dtype=np.intp)
        inv_idx[perm] = imask
        ret += (inv_idx.reshape(inverse_shape) if axis is None else inv_idx,)
    if return_counts:
        idx = np.concatenate(np.nonzero(mask) + ([mask.size],))
        ret += (np.diff(idx),)
    return ret

class UniqueAllResult(NamedTuple):
    values: np.ndarray
    indices: np.ndarray
    inverse_indices: np.ndarray
    counts: np.ndarray

class UniqueCountsResult(NamedTuple):
    values: np.ndarray
    counts: np.ndarray

class UniqueInverseResult(NamedTuple):
    values: np.ndarray
    inverse_indices: np.ndarray

def _unique_all_dispatcher(x, /):
    return (x,)

@array_function_dispatch(_unique_all_dispatcher)
def unique_all(x):
    result = unique(x, return_index=True, return_inverse=True, return_counts=True, equal_nan=False)
    return UniqueAllResult(*result)

def _unique_counts_dispatcher(x, /):
    return (x,)

@array_function_dispatch(_unique_counts_dispatcher)
def unique_counts(x):
    result = unique(x, return_index=False, return_inverse=False, return_counts=True, equal_nan=False)
    return UniqueCountsResult(*result)

def _unique_inverse_dispatcher(x, /):
    return (x,)

@array_function_dispatch(_unique_inverse_dispatcher)
def unique_inverse(x):
    result = unique(x, return_index=False, return_inverse=True, return_counts=False, equal_nan=False)
    return UniqueInverseResult(*result)

def _unique_values_dispatcher(x, /):
    return (x,)

@array_function_dispatch(_unique_values_dispatcher)
def unique_values(x):
    return unique(x, return_index=False, return_inverse=False, return_counts=False, equal_nan=False, sorted=False)

def _intersect1d_dispatcher(ar1, ar2, assume_unique=None, return_indices=None):
    return (ar1, ar2)

@array_function_dispatch(_intersect1d_dispatcher)
def intersect1d(ar1, ar2, assume_unique=False, return_indices=False):
    ar1 = np.asanyarray(ar1)
    ar2 = np.asanyarray(ar2)
    if not assume_unique:
        if return_indices:
            ar1, ind1 = unique(ar1, return_index=True)
            ar2, ind2 = unique(ar2, return_index=True)
        else:
            ar1 = unique(ar1)
            ar2 = unique(ar2)
    else:
        ar1 = ar1.ravel()
        ar2 = ar2.ravel()
    aux = np.concatenate((ar1, ar2))
    if return_indices:
        aux_sort_indices = np.argsort(aux, kind='mergesort')
        aux = aux[aux_sort_indices]
    else:
        aux.sort()
    mask = aux[1:] == aux[:-1]
    int1d = aux[:-1][mask]
    if return_indices:
        ar1_indices = aux_sort_indices[:-1][mask]
        ar2_indices = aux_sort_indices[1:][mask] - ar1.size
        if not assume_unique:
            ar1_indices = ind1[ar1_indices]
            ar2_indices = ind2[ar2_indices]
        return (int1d, ar1_indices, ar2_indices)
    else:
        return int1d

def _setxor1d_dispatcher(ar1, ar2, assume_unique=None):
    return (ar1, ar2)

@array_function_dispatch(_setxor1d_dispatcher)
def setxor1d(ar1, ar2, assume_unique=False):
    if not assume_unique:
        ar1 = unique(ar1)
        ar2 = unique(ar2)
    aux = np.concatenate((ar1, ar2), axis=None)
    if aux.size == 0:
        return aux
    aux.sort()
    flag = np.concatenate(([True], aux[1:] != aux[:-1], [True]))
    return aux[flag[1:] & flag[:-1]]

def _isin(ar1, ar2, assume_unique=False, invert=False, *, kind=None):
    ar1 = np.asarray(ar1).ravel()
    ar2 = np.asarray(ar2).ravel()
    if ar2.dtype == object:
        ar2 = ar2.reshape(-1, 1)
    if kind not in {None, 'sort', 'table'}:
        raise ValueError(f"Invalid kind: '{kind}'. Please use None, 'sort' or 'table'.")
    is_int_arrays = all((ar.dtype.kind in ('u', 'i', 'b') for ar in (ar1, ar2)))
    use_table_method = is_int_arrays and kind in {None, 'table'}
    if use_table_method:
        if ar2.size == 0:
            if invert:
                return np.ones_like(ar1, dtype=bool)
            else:
                return np.zeros_like(ar1, dtype=bool)
        if ar1.dtype == bool:
            ar1 = ar1.astype(np.uint8)
        if ar2.dtype == bool:
            ar2 = ar2.astype(np.uint8)
        ar2_min = int(np.min(ar2))
        ar2_max = int(np.max(ar2))
        ar2_range = ar2_max - ar2_min
        below_memory_constraint = ar2_range <= 6 * (ar1.size + ar2.size)
        range_safe_from_overflow = ar2_range <= np.iinfo(ar2.dtype).max
        if range_safe_from_overflow and (below_memory_constraint or kind == 'table'):
            if invert:
                outgoing_array = np.ones_like(ar1, dtype=bool)
            else:
                outgoing_array = np.zeros_like(ar1, dtype=bool)
            if invert:
                isin_helper_ar = np.ones(ar2_range + 1, dtype=bool)
                isin_helper_ar[ar2 - ar2_min] = 0
            else:
                isin_helper_ar = np.zeros(ar2_range + 1, dtype=bool)
                isin_helper_ar[ar2 - ar2_min] = 1
            basic_mask = (ar1 <= ar2_max) & (ar1 >= ar2_min)
            in_range_ar1 = ar1[basic_mask]
            if in_range_ar1.size == 0:
                return outgoing_array
            try:
                ar2_min = np.array(ar2_min, dtype=np.intp)
                dtype = np.intp
            except OverflowError:
                dtype = ar2.dtype
            out = np.empty_like(in_range_ar1, dtype=np.intp)
            outgoing_array[basic_mask] = isin_helper_ar[np.subtract(in_range_ar1, ar2_min, dtype=dtype, out=out, casting='unsafe')]
            return outgoing_array
        elif kind == 'table':
            raise RuntimeError("You have specified kind='table', but the range of values in `ar2` or `ar1` exceed the maximum integer of the datatype. Please set `kind` to None or 'sort'.")
    elif kind == 'table':
        raise ValueError("The 'table' method is only supported for boolean or integer arrays. Please select 'sort' or None for kind.")
    contains_object = ar1.dtype.hasobject or ar2.dtype.hasobject
    if len(ar2) < 10 * len(ar1) ** 0.145 or contains_object:
        if invert:
            mask = np.ones(len(ar1), dtype=bool)
            for a in ar2:
                mask &= ar1 != a
        else:
            mask = np.zeros(len(ar1), dtype=bool)
            for a in ar2:
                mask |= ar1 == a
        return mask
    if not assume_unique:
        ar1, rev_idx = np.unique(ar1, return_inverse=True)
        ar2 = np.unique(ar2)
    ar = np.concatenate((ar1, ar2))
    order = ar.argsort(kind='mergesort')
    sar = ar[order]
    if invert:
        bool_ar = sar[1:] != sar[:-1]
    else:
        bool_ar = sar[1:] == sar[:-1]
    flag = np.concatenate((bool_ar, [invert]))
    ret = np.empty(ar.shape, dtype=bool)
    ret[order] = flag
    if assume_unique:
        return ret[:len(ar1)]
    else:
        return ret[rev_idx]

def _isin_dispatcher(element, test_elements, assume_unique=None, invert=None, *, kind=None):
    return (element, test_elements)

@array_function_dispatch(_isin_dispatcher)
def isin(element, test_elements, assume_unique=False, invert=False, *, kind=None):
    element = np.asarray(element)
    return _isin(element, test_elements, assume_unique=assume_unique, invert=invert, kind=kind).reshape(element.shape)

def _union1d_dispatcher(ar1, ar2):
    return (ar1, ar2)

@array_function_dispatch(_union1d_dispatcher)
def union1d(ar1, ar2):
    return unique(np.concatenate((ar1, ar2), axis=None))

def _setdiff1d_dispatcher(ar1, ar2, assume_unique=None):
    return (ar1, ar2)

@array_function_dispatch(_setdiff1d_dispatcher)
def setdiff1d(ar1, ar2, assume_unique=False):
    if assume_unique:
        ar1 = np.asarray(ar1).ravel()
    else:
        ar1 = unique(ar1)
        ar2 = unique(ar2)
    return ar1[_isin(ar1, ar2, assume_unique=True, invert=True)]
