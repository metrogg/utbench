import typing
import numpy as np
from numpy._core.overrides import array_function_dispatch
from numpy.lib._index_tricks_impl import ndindex
__all__ = ['pad']

def _round_if_needed(arr, dtype):
    if np.issubdtype(dtype, np.integer):
        arr.round(out=arr)

def _slice_at_axis(sl, axis):
    return (slice(None),) * axis + (sl,) + (...,)

def _view_roi(array, original_area_slice, axis):
    axis += 1
    sl = (slice(None),) * axis + original_area_slice[axis:]
    return array[sl]

def _pad_simple(array, pad_width, fill_value=None):
    new_shape = tuple((left + size + right for size, (left, right) in zip(array.shape, pad_width)))
    order = 'F' if array.flags.fnc else 'C'
    padded = np.empty(new_shape, dtype=array.dtype, order=order)
    if fill_value is not None:
        padded.fill(fill_value)
    original_area_slice = tuple((slice(left, left + size) for size, (left, right) in zip(array.shape, pad_width)))
    padded[original_area_slice] = array
    return (padded, original_area_slice)

def _set_pad_area(padded, axis, width_pair, value_pair):
    left_slice = _slice_at_axis(slice(None, width_pair[0]), axis)
    padded[left_slice] = value_pair[0]
    right_slice = _slice_at_axis(slice(padded.shape[axis] - width_pair[1], None), axis)
    padded[right_slice] = value_pair[1]

def _get_edges(padded, axis, width_pair):
    left_index = width_pair[0]
    left_slice = _slice_at_axis(slice(left_index, left_index + 1), axis)
    left_edge = padded[left_slice]
    right_index = padded.shape[axis] - width_pair[1]
    right_slice = _slice_at_axis(slice(right_index - 1, right_index), axis)
    right_edge = padded[right_slice]
    return (left_edge, right_edge)

def _get_linear_ramps(padded, axis, width_pair, end_value_pair):
    edge_pair = _get_edges(padded, axis, width_pair)
    left_ramp, right_ramp = (np.linspace(start=end_value, stop=edge.squeeze(axis), num=width, endpoint=False, dtype=padded.dtype, axis=axis) for end_value, edge, width in zip(end_value_pair, edge_pair, width_pair))
    right_ramp = right_ramp[_slice_at_axis(slice(None, None, -1), axis)]
    return (left_ramp, right_ramp)

def _get_stats(padded, axis, width_pair, length_pair, stat_func):
    left_index = width_pair[0]
    right_index = padded.shape[axis] - width_pair[1]
    max_length = right_index - left_index
    left_length, right_length = length_pair
    if left_length is None or max_length < left_length:
        left_length = max_length
    if right_length is None or max_length < right_length:
        right_length = max_length
    if (left_length == 0 or right_length == 0) and stat_func in {np.amax, np.amin}:
        raise ValueError('stat_length of 0 yields no value for padding')
    left_slice = _slice_at_axis(slice(left_index, left_index + left_length), axis)
    left_chunk = padded[left_slice]
    left_stat = stat_func(left_chunk, axis=axis, keepdims=True)
    _round_if_needed(left_stat, padded.dtype)
    if left_length == right_length == max_length:
        return (left_stat, left_stat)
    right_slice = _slice_at_axis(slice(right_index - right_length, right_index), axis)
    right_chunk = padded[right_slice]
    right_stat = stat_func(right_chunk, axis=axis, keepdims=True)
    _round_if_needed(right_stat, padded.dtype)
    return (left_stat, right_stat)

def _set_reflect_both(padded, axis, width_pair, method, original_period, include_edge=False):
    left_pad, right_pad = width_pair
    old_length = padded.shape[axis] - right_pad - left_pad
    if include_edge:
        old_length = old_length // original_period * original_period
        edge_offset = 1
    else:
        old_length = (old_length - 1) // (original_period - 1) * (original_period - 1) + 1
        edge_offset = 0
        old_length -= 1
    if left_pad > 0:
        chunk_length = min(old_length, left_pad)
        stop = left_pad - edge_offset
        start = stop + chunk_length
        left_slice = _slice_at_axis(slice(start, stop, -1), axis)
        left_chunk = padded[left_slice]
        if method == 'odd':
            edge_slice = _slice_at_axis(slice(left_pad, left_pad + 1), axis)
            left_chunk = 2 * padded[edge_slice] - left_chunk
        start = left_pad - chunk_length
        stop = left_pad
        pad_area = _slice_at_axis(slice(start, stop), axis)
        padded[pad_area] = left_chunk
        left_pad -= chunk_length
    if right_pad > 0:
        chunk_length = min(old_length, right_pad)
        start = -right_pad + edge_offset - 2
        stop = start - chunk_length
        right_slice = _slice_at_axis(slice(start, stop, -1), axis)
        right_chunk = padded[right_slice]
        if method == 'odd':
            edge_slice = _slice_at_axis(slice(-right_pad - 1, -right_pad), axis)
            right_chunk = 2 * padded[edge_slice] - right_chunk
        start = padded.shape[axis] - right_pad
        stop = start + chunk_length
        pad_area = _slice_at_axis(slice(start, stop), axis)
        padded[pad_area] = right_chunk
        right_pad -= chunk_length
    return (left_pad, right_pad)

def _set_wrap_both(padded, axis, width_pair, original_period):
    left_pad, right_pad = width_pair
    period = padded.shape[axis] - right_pad - left_pad
    period = period // original_period * original_period
    new_left_pad = 0
    new_right_pad = 0
    if left_pad > 0:
        slice_end = left_pad + period
        slice_start = slice_end - min(period, left_pad)
        right_slice = _slice_at_axis(slice(slice_start, slice_end), axis)
        right_chunk = padded[right_slice]
        if left_pad > period:
            pad_area = _slice_at_axis(slice(left_pad - period, left_pad), axis)
            new_left_pad = left_pad - period
        else:
            pad_area = _slice_at_axis(slice(None, left_pad), axis)
        padded[pad_area] = right_chunk
    if right_pad > 0:
        slice_start = -right_pad - period
        slice_end = slice_start + min(period, right_pad)
        left_slice = _slice_at_axis(slice(slice_start, slice_end), axis)
        left_chunk = padded[left_slice]
        if right_pad > period:
            pad_area = _slice_at_axis(slice(-right_pad, -right_pad + period), axis)
            new_right_pad = right_pad - period
        else:
            pad_area = _slice_at_axis(slice(-right_pad, None), axis)
        padded[pad_area] = left_chunk
    return (new_left_pad, new_right_pad)

def _as_pairs(x, ndim, as_index=False):
    if x is None:
        return ((None, None),) * ndim
    x = np.array(x)
    if as_index:
        x = np.round(x).astype(np.intp, copy=False)
    if x.ndim < 3:
        if x.size == 1:
            x = x.ravel()
            if as_index and x < 0:
                raise ValueError("index can't contain negative values")
            return ((x[0], x[0]),) * ndim
        if x.size == 2 and x.shape != (2, 1):
            x = x.ravel()
            if as_index and (x[0] < 0 or x[1] < 0):
                raise ValueError("index can't contain negative values")
            return ((x[0], x[1]),) * ndim
    if as_index and x.min() < 0:
        raise ValueError("index can't contain negative values")
    return np.broadcast_to(x, (ndim, 2)).tolist()

def _pad_dispatcher(array, pad_width, mode=None, **kwargs):
    return (array,)

@array_function_dispatch(_pad_dispatcher, module='numpy')
def pad(array, pad_width, mode='constant', **kwargs):
    array = np.asarray(array)
    if isinstance(pad_width, dict):
        seq = [(0, 0)] * array.ndim
        for axis, width in pad_width.items():
            match width:
                case int(both):
                    seq[axis] = (both, both)
                case tuple([int(before), int(after)]):
                    seq[axis] = (before, after)
                case _ as invalid:
                    typing.assert_never(invalid)
        pad_width = seq
    pad_width = np.asarray(pad_width)
    if not pad_width.dtype.kind == 'i':
        raise TypeError('`pad_width` must be of integral type.')
    pad_width = _as_pairs(pad_width, array.ndim, as_index=True)
    if callable(mode):
        function = mode
        padded, _ = _pad_simple(array, pad_width, fill_value=0)
        for axis in range(padded.ndim):
            view = np.moveaxis(padded, axis, -1)
            inds = ndindex(view.shape[:-1])
            inds = (ind + (Ellipsis,) for ind in inds)
            for ind in inds:
                function(view[ind], pad_width[axis], axis, kwargs)
        return padded
    allowed_kwargs = {'empty': [], 'edge': [], 'wrap': [], 'constant': ['constant_values'], 'linear_ramp': ['end_values'], 'maximum': ['stat_length'], 'mean': ['stat_length'], 'median': ['stat_length'], 'minimum': ['stat_length'], 'reflect': ['reflect_type'], 'symmetric': ['reflect_type']}
    try:
        unsupported_kwargs = set(kwargs) - set(allowed_kwargs[mode])
    except KeyError:
        raise ValueError(f"mode '{mode}' is not supported") from None
    if unsupported_kwargs:
        raise ValueError(f"unsupported keyword arguments for mode '{mode}': {unsupported_kwargs}")
    stat_functions = {'maximum': np.amax, 'minimum': np.amin, 'mean': np.mean, 'median': np.median}
    padded, original_area_slice = _pad_simple(array, pad_width)
    axes = range(padded.ndim)
    if mode == 'constant':
        values = kwargs.get('constant_values', 0)
        values = _as_pairs(values, padded.ndim)
        for axis, width_pair, value_pair in zip(axes, pad_width, values):
            roi = _view_roi(padded, original_area_slice, axis)
            _set_pad_area(roi, axis, width_pair, value_pair)
    elif mode == 'empty':
        pass
    elif array.size == 0:
        for axis, width_pair in zip(axes, pad_width):
            if array.shape[axis] == 0 and any(width_pair):
                raise ValueError(f"can't extend empty axis {axis} using modes other than 'constant' or 'empty'")
    elif mode == 'edge':
        for axis, width_pair in zip(axes, pad_width):
            roi = _view_roi(padded, original_area_slice, axis)
            edge_pair = _get_edges(roi, axis, width_pair)
            _set_pad_area(roi, axis, width_pair, edge_pair)
    elif mode == 'linear_ramp':
        end_values = kwargs.get('end_values', 0)
        end_values = _as_pairs(end_values, padded.ndim)
        for axis, width_pair, value_pair in zip(axes, pad_width, end_values):
            roi = _view_roi(padded, original_area_slice, axis)
            ramp_pair = _get_linear_ramps(roi, axis, width_pair, value_pair)
            _set_pad_area(roi, axis, width_pair, ramp_pair)
    elif mode in stat_functions:
        func = stat_functions[mode]
        length = kwargs.get('stat_length')
        length = _as_pairs(length, padded.ndim, as_index=True)
        for axis, width_pair, length_pair in zip(axes, pad_width, length):
            roi = _view_roi(padded, original_area_slice, axis)
            stat_pair = _get_stats(roi, axis, width_pair, length_pair, func)
            _set_pad_area(roi, axis, width_pair, stat_pair)
    elif mode in {'reflect', 'symmetric'}:
        method = kwargs.get('reflect_type', 'even')
        include_edge = mode == 'symmetric'
        for axis, (left_index, right_index) in zip(axes, pad_width):
            if array.shape[axis] == 1 and (left_index > 0 or right_index > 0):
                edge_pair = _get_edges(padded, axis, (left_index, right_index))
                _set_pad_area(padded, axis, (left_index, right_index), edge_pair)
                continue
            roi = _view_roi(padded, original_area_slice, axis)
            while left_index > 0 or right_index > 0:
                left_index, right_index = _set_reflect_both(roi, axis, (left_index, right_index), method, array.shape[axis], include_edge)
    elif mode == 'wrap':
        for axis, (left_index, right_index) in zip(axes, pad_width):
            roi = _view_roi(padded, original_area_slice, axis)
            original_period = padded.shape[axis] - right_index - left_index
            while left_index > 0 or right_index > 0:
                left_index, right_index = _set_wrap_both(roi, axis, (left_index, right_index), original_period)
    return padded
