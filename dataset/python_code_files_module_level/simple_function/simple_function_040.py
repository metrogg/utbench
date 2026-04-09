__all__ = ['atleast_1d', 'atleast_2d', 'atleast_3d', 'block', 'hstack', 'stack', 'unstack', 'vstack']
import functools
import itertools
import operator
from . import fromnumeric as _from_nx, numeric as _nx, overrides
from .multiarray import array, asanyarray, normalize_axis_index
array_function_dispatch = functools.partial(overrides.array_function_dispatch, module='numpy')

def _atleast_1d_dispatcher(*arys):
    return arys

@array_function_dispatch(_atleast_1d_dispatcher)
def atleast_1d(*arys):
    if len(arys) == 1:
        result = asanyarray(arys[0])
        if result.ndim == 0:
            result = result.reshape(1)
        return result
    res = []
    for ary in arys:
        result = asanyarray(ary)
        if result.ndim == 0:
            result = result.reshape(1)
        res.append(result)
    return tuple(res)

def _atleast_2d_dispatcher(*arys):
    return arys

@array_function_dispatch(_atleast_2d_dispatcher)
def atleast_2d(*arys):
    res = []
    for ary in arys:
        ary = asanyarray(ary)
        if ary.ndim == 0:
            result = ary.reshape(1, 1)
        elif ary.ndim == 1:
            result = ary[_nx.newaxis, :]
        else:
            result = ary
        res.append(result)
    if len(res) == 1:
        return res[0]
    else:
        return tuple(res)

def _atleast_3d_dispatcher(*arys):
    return arys

@array_function_dispatch(_atleast_3d_dispatcher)
def atleast_3d(*arys):
    res = []
    for ary in arys:
        ary = asanyarray(ary)
        if ary.ndim == 0:
            result = ary.reshape(1, 1, 1)
        elif ary.ndim == 1:
            result = ary[_nx.newaxis, :, _nx.newaxis]
        elif ary.ndim == 2:
            result = ary[:, :, _nx.newaxis]
        else:
            result = ary
        res.append(result)
    if len(res) == 1:
        return res[0]
    else:
        return tuple(res)

def _arrays_for_stack_dispatcher(arrays):
    if not hasattr(arrays, '__getitem__'):
        raise TypeError('arrays to stack must be passed as a "sequence" type such as list or tuple.')
    return tuple(arrays)

def _vhstack_dispatcher(tup, *, dtype=None, casting=None):
    return _arrays_for_stack_dispatcher(tup)

@array_function_dispatch(_vhstack_dispatcher)
def vstack(tup, *, dtype=None, casting='same_kind'):
    arrs = atleast_2d(*tup)
    if not isinstance(arrs, tuple):
        arrs = (arrs,)
    return _nx.concatenate(arrs, 0, dtype=dtype, casting=casting)

@array_function_dispatch(_vhstack_dispatcher)
def hstack(tup, *, dtype=None, casting='same_kind'):
    arrs = atleast_1d(*tup)
    if not isinstance(arrs, tuple):
        arrs = (arrs,)
    if arrs and arrs[0].ndim == 1:
        return _nx.concatenate(arrs, 0, dtype=dtype, casting=casting)
    else:
        return _nx.concatenate(arrs, 1, dtype=dtype, casting=casting)

def _stack_dispatcher(arrays, axis=None, out=None, *, dtype=None, casting=None):
    arrays = _arrays_for_stack_dispatcher(arrays)
    if out is not None:
        arrays = list(arrays)
        arrays.append(out)
    return arrays

@array_function_dispatch(_stack_dispatcher)
def stack(arrays, axis=0, out=None, *, dtype=None, casting='same_kind'):
    arrays = [asanyarray(arr) for arr in arrays]
    if not arrays:
        raise ValueError('need at least one array to stack')
    shapes = {arr.shape for arr in arrays}
    if len(shapes) != 1:
        raise ValueError('all input arrays must have the same shape')
    result_ndim = arrays[0].ndim + 1
    axis = normalize_axis_index(axis, result_ndim)
    sl = (slice(None),) * axis + (_nx.newaxis,)
    expanded_arrays = [arr[sl] for arr in arrays]
    return _nx.concatenate(expanded_arrays, axis=axis, out=out, dtype=dtype, casting=casting)

def _unstack_dispatcher(x, /, *, axis=None):
    return (x,)

@array_function_dispatch(_unstack_dispatcher)
def unstack(x, /, *, axis=0):
    if x.ndim == 0:
        raise ValueError('Input array must be at least 1-d.')
    return tuple(_nx.moveaxis(x, axis, 0))
_size = getattr(_from_nx.size, '__wrapped__', _from_nx.size)
_ndim = getattr(_from_nx.ndim, '__wrapped__', _from_nx.ndim)
_concatenate = getattr(_from_nx.concatenate, '__wrapped__', _from_nx.concatenate)

def _block_format_index(index):
    idx_str = ''.join((f'[{i}]' for i in index if i is not None))
    return 'arrays' + idx_str

def _block_check_depths_match(arrays, parent_index=[]):
    if isinstance(arrays, tuple):
        raise TypeError(f'{_block_format_index(parent_index)} is a tuple. Only lists can be used to arrange blocks, and np.block does not allow implicit conversion from tuple to ndarray.')
    elif isinstance(arrays, list) and len(arrays) > 0:
        idxs_ndims = (_block_check_depths_match(arr, parent_index + [i]) for i, arr in enumerate(arrays))
        first_index, max_arr_ndim, final_size = next(idxs_ndims)
        for index, ndim, size in idxs_ndims:
            final_size += size
            if ndim > max_arr_ndim:
                max_arr_ndim = ndim
            if len(index) != len(first_index):
                raise ValueError(f'List depths are mismatched. First element was at depth {len(first_index)}, but there is an element at depth {len(index)} ({_block_format_index(index)})')
            if index[-1] is None:
                first_index = index
        return (first_index, max_arr_ndim, final_size)
    elif isinstance(arrays, list) and len(arrays) == 0:
        return (parent_index + [None], 0, 0)
    else:
        size = _size(arrays)
        return (parent_index, _ndim(arrays), size)

def _atleast_nd(a, ndim):
    return array(a, ndmin=ndim, copy=None, subok=True)

def _accumulate(values):
    return list(itertools.accumulate(values))

def _concatenate_shapes(shapes, axis):
    shape_at_axis = [shape[axis] for shape in shapes]
    first_shape = shapes[0]
    first_shape_pre = first_shape[:axis]
    first_shape_post = first_shape[axis + 1:]
    if any((shape[:axis] != first_shape_pre or shape[axis + 1:] != first_shape_post for shape in shapes)):
        raise ValueError(f'Mismatched array shapes in block along axis {axis}.')
    shape = first_shape_pre + (sum(shape_at_axis),) + first_shape[axis + 1:]
    offsets_at_axis = _accumulate(shape_at_axis)
    slice_prefixes = [(slice(start, end),) for start, end in zip([0] + offsets_at_axis, offsets_at_axis)]
    return (shape, slice_prefixes)

def _block_info_recursion(arrays, max_depth, result_ndim, depth=0):
    if depth < max_depth:
        shapes, slices, arrays = zip(*[_block_info_recursion(arr, max_depth, result_ndim, depth + 1) for arr in arrays])
        axis = result_ndim - max_depth + depth
        shape, slice_prefixes = _concatenate_shapes(shapes, axis)
        slices = [slice_prefix + the_slice for slice_prefix, inner_slices in zip(slice_prefixes, slices) for the_slice in inner_slices]
        arrays = functools.reduce(operator.add, arrays)
        return (shape, slices, arrays)
    else:
        arr = _atleast_nd(arrays, result_ndim)
        return (arr.shape, [()], [arr])

def _block(arrays, max_depth, result_ndim, depth=0):
    if depth < max_depth:
        arrs = [_block(arr, max_depth, result_ndim, depth + 1) for arr in arrays]
        return _concatenate(arrs, axis=-(max_depth - depth))
    else:
        return _atleast_nd(arrays, result_ndim)

def _block_dispatcher(arrays):
    if isinstance(arrays, list):
        for subarrays in arrays:
            yield from _block_dispatcher(subarrays)
    else:
        yield arrays

@array_function_dispatch(_block_dispatcher)
def block(arrays):
    arrays, list_ndim, result_ndim, final_size = _block_setup(arrays)
    if list_ndim * final_size > 2 * 512 * 512:
        return _block_slicing(arrays, list_ndim, result_ndim)
    else:
        return _block_concatenate(arrays, list_ndim, result_ndim)

def _block_setup(arrays):
    bottom_index, arr_ndim, final_size = _block_check_depths_match(arrays)
    list_ndim = len(bottom_index)
    if bottom_index and bottom_index[-1] is None:
        raise ValueError(f'List at {_block_format_index(bottom_index)} cannot be empty')
    result_ndim = max(arr_ndim, list_ndim)
    return (arrays, list_ndim, result_ndim, final_size)

def _block_slicing(arrays, list_ndim, result_ndim):
    shape, slices, arrays = _block_info_recursion(arrays, list_ndim, result_ndim)
    dtype = _nx.result_type(*[arr.dtype for arr in arrays])
    F_order = all((arr.flags['F_CONTIGUOUS'] for arr in arrays))
    C_order = all((arr.flags['C_CONTIGUOUS'] for arr in arrays))
    order = 'F' if F_order and (not C_order) else 'C'
    result = _nx.empty(shape=shape, dtype=dtype, order=order)
    for the_slice, arr in zip(slices, arrays):
        result[(Ellipsis,) + the_slice] = arr
    return result

def _block_concatenate(arrays, list_ndim, result_ndim):
    result = _block(arrays, list_ndim, result_ndim)
    if list_ndim == 0:
        result = result.copy()
    return result
