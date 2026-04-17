import functools
import warnings
import numpy as np
import numpy._core.numeric as _nx
from numpy._core import atleast_3d, overrides, vstack
from numpy._core._multiarray_umath import _array_converter
from numpy._core.fromnumeric import reshape, transpose
from numpy._core.multiarray import normalize_axis_index
from numpy._core.numeric import array, asanyarray, asarray, normalize_axis_tuple, zeros, zeros_like
from numpy._core.overrides import set_module
from numpy._core.shape_base import _arrays_for_stack_dispatcher
from numpy.lib._index_tricks_impl import ndindex
from numpy.matrixlib.defmatrix import matrix
__all__ = ['column_stack', 'row_stack', 'dstack', 'array_split', 'split', 'hsplit', 'vsplit', 'dsplit', 'apply_over_axes', 'expand_dims', 'apply_along_axis', 'kron', 'tile', 'take_along_axis', 'put_along_axis']
array_function_dispatch = functools.partial(overrides.array_function_dispatch, module='numpy')

def _make_along_axis_idx(arr_shape, indices, axis):
    if not _nx.issubdtype(indices.dtype, _nx.integer):
        raise IndexError('`indices` must be an integer array')
    if len(arr_shape) != indices.ndim:
        raise ValueError('`indices` and `arr` must have the same number of dimensions')
    shape_ones = (1,) * indices.ndim
    dest_dims = list(range(axis)) + [None] + list(range(axis + 1, indices.ndim))
    fancy_index = []
    for dim, n in zip(dest_dims, arr_shape):
        if dim is None:
            fancy_index.append(indices)
        else:
            ind_shape = shape_ones[:dim] + (-1,) + shape_ones[dim + 1:]
            fancy_index.append(_nx.arange(n).reshape(ind_shape))
    return tuple(fancy_index)

def _take_along_axis_dispatcher(arr, indices, axis=None):
    return (arr, indices)

@array_function_dispatch(_take_along_axis_dispatcher)
def take_along_axis(arr, indices, axis=-1):
    if axis is None:
        if indices.ndim != 1:
            raise ValueError('when axis=None, `indices` must have a single dimension.')
        arr = np.array(arr.flat)
        axis = 0
    else:
        axis = normalize_axis_index(axis, arr.ndim)
    return arr[_make_along_axis_idx(arr.shape, indices, axis)]

def _put_along_axis_dispatcher(arr, indices, values, axis):
    return (arr, indices, values)

@array_function_dispatch(_put_along_axis_dispatcher)
def put_along_axis(arr, indices, values, axis):
    if axis is None:
        if indices.ndim != 1:
            raise ValueError('when axis=None, `indices` must have a single dimension.')
        arr = np.array(arr.flat)
        axis = 0
    else:
        axis = normalize_axis_index(axis, arr.ndim)
    arr[_make_along_axis_idx(arr.shape, indices, axis)] = values

def _apply_along_axis_dispatcher(func1d, axis, arr, *args, **kwargs):
    return (arr,)

@array_function_dispatch(_apply_along_axis_dispatcher)
def apply_along_axis(func1d, axis, arr, *args, **kwargs):
    conv = _array_converter(arr)
    arr = conv[0]
    nd = arr.ndim
    axis = normalize_axis_index(axis, nd)
    in_dims = list(range(nd))
    inarr_view = transpose(arr, in_dims[:axis] + in_dims[axis + 1:] + [axis])
    inds = ndindex(inarr_view.shape[:-1])
    inds = (ind + (Ellipsis,) for ind in inds)
    try:
        ind0 = next(inds)
    except StopIteration:
        raise ValueError('Cannot apply_along_axis when any iteration dimensions are 0') from None
    res = asanyarray(func1d(inarr_view[ind0], *args, **kwargs))
    if not isinstance(res, matrix):
        buff = zeros_like(res, shape=inarr_view.shape[:-1] + res.shape)
    else:
        buff = zeros(inarr_view.shape[:-1] + res.shape, dtype=res.dtype)
    buff_dims = list(range(buff.ndim))
    buff_permute = buff_dims[0:axis] + buff_dims[buff.ndim - res.ndim:buff.ndim] + buff_dims[axis:buff.ndim - res.ndim]
    buff[ind0] = res
    for ind in inds:
        buff[ind] = asanyarray(func1d(inarr_view[ind], *args, **kwargs))
    res = transpose(buff, buff_permute)
    return conv.wrap(res)

def _apply_over_axes_dispatcher(func, a, axes):
    return (a,)

@array_function_dispatch(_apply_over_axes_dispatcher)
def apply_over_axes(func, a, axes):
    val = asarray(a)
    N = a.ndim
    if array(axes).ndim == 0:
        axes = (axes,)
    for axis in axes:
        if axis < 0:
            axis = N + axis
        args = (val, axis)
        res = func(*args)
        if res.ndim == val.ndim:
            val = res
        else:
            res = expand_dims(res, axis)
            if res.ndim == val.ndim:
                val = res
            else:
                raise ValueError('function is not returning an array of the correct shape')
    return val

def _expand_dims_dispatcher(a, axis):
    return (a,)

@array_function_dispatch(_expand_dims_dispatcher)
def expand_dims(a, axis):
    if isinstance(a, matrix):
        a = asarray(a)
    else:
        a = asanyarray(a)
    if not isinstance(axis, (tuple, list)):
        axis = (axis,)
    out_ndim = len(axis) + a.ndim
    axis = normalize_axis_tuple(axis, out_ndim)
    shape_it = iter(a.shape)
    shape = [1 if ax in axis else next(shape_it) for ax in range(out_ndim)]
    return a.reshape(shape)

@set_module('numpy')
def row_stack(tup, *, dtype=None, casting='same_kind'):
    warnings.warn('`row_stack` alias is deprecated. Use `np.vstack` directly.', DeprecationWarning, stacklevel=2)
    return vstack(tup, dtype=dtype, casting=casting)
row_stack.__doc__ = vstack.__doc__

def _column_stack_dispatcher(tup):
    return _arrays_for_stack_dispatcher(tup)

@array_function_dispatch(_column_stack_dispatcher)
def column_stack(tup):
    arrays = []
    for v in tup:
        arr = asanyarray(v)
        if arr.ndim < 2:
            arr = array(arr, copy=None, subok=True, ndmin=2).T
        arrays.append(arr)
    return _nx.concatenate(arrays, 1)

def _dstack_dispatcher(tup):
    return _arrays_for_stack_dispatcher(tup)

@array_function_dispatch(_dstack_dispatcher)
def dstack(tup):
    arrs = atleast_3d(*tup)
    if not isinstance(arrs, tuple):
        arrs = (arrs,)
    return _nx.concatenate(arrs, 2)

def _array_split_dispatcher(ary, indices_or_sections, axis=None):
    return (ary, indices_or_sections)

@array_function_dispatch(_array_split_dispatcher)
def array_split(ary, indices_or_sections, axis=0):
    try:
        Ntotal = ary.shape[axis]
    except AttributeError:
        Ntotal = len(ary)
    try:
        Nsections = len(indices_or_sections) + 1
        div_points = [0] + list(indices_or_sections) + [Ntotal]
    except TypeError:
        Nsections = int(indices_or_sections)
        if Nsections <= 0:
            raise ValueError('number sections must be larger than 0.') from None
        Neach_section, extras = divmod(Ntotal, Nsections)
        section_sizes = [0] + extras * [Neach_section + 1] + (Nsections - extras) * [Neach_section]
        div_points = _nx.array(section_sizes, dtype=_nx.intp).cumsum()
    sub_arys = []
    sary = _nx.swapaxes(ary, axis, 0)
    for i in range(Nsections):
        st = div_points[i]
        end = div_points[i + 1]
        sub_arys.append(_nx.swapaxes(sary[st:end], axis, 0))
    return sub_arys

def _split_dispatcher(ary, indices_or_sections, axis=None):
    return (ary, indices_or_sections)

@array_function_dispatch(_split_dispatcher)
def split(ary, indices_or_sections, axis=0):
    try:
        len(indices_or_sections)
    except TypeError:
        sections = indices_or_sections
        N = ary.shape[axis]
        if N % sections:
            raise ValueError('array split does not result in an equal division') from None
    return array_split(ary, indices_or_sections, axis)

def _hvdsplit_dispatcher(ary, indices_or_sections):
    return (ary, indices_or_sections)

@array_function_dispatch(_hvdsplit_dispatcher)
def hsplit(ary, indices_or_sections):
    if _nx.ndim(ary) == 0:
        raise ValueError('hsplit only works on arrays of 1 or more dimensions')
    if ary.ndim > 1:
        return split(ary, indices_or_sections, 1)
    else:
        return split(ary, indices_or_sections, 0)

@array_function_dispatch(_hvdsplit_dispatcher)
def vsplit(ary, indices_or_sections):
    if _nx.ndim(ary) < 2:
        raise ValueError('vsplit only works on arrays of 2 or more dimensions')
    return split(ary, indices_or_sections, 0)

@array_function_dispatch(_hvdsplit_dispatcher)
def dsplit(ary, indices_or_sections):
    if _nx.ndim(ary) < 3:
        raise ValueError('dsplit only works on arrays of 3 or more dimensions')
    return split(ary, indices_or_sections, 2)

def get_array_wrap(*args):
    warnings.warn('`get_array_wrap` is deprecated. (deprecated in NumPy 2.0)', DeprecationWarning, stacklevel=2)
    wrappers = sorted(((getattr(x, '__array_priority__', 0), -i, x.__array_wrap__) for i, x in enumerate(args) if hasattr(x, '__array_wrap__')))
    if wrappers:
        return wrappers[-1][-1]
    return None

def _kron_dispatcher(a, b):
    return (a, b)

@array_function_dispatch(_kron_dispatcher)
def kron(a, b):
    b = asanyarray(b)
    a = array(a, copy=None, subok=True, ndmin=b.ndim)
    is_any_mat = isinstance(a, matrix) or isinstance(b, matrix)
    ndb, nda = (b.ndim, a.ndim)
    nd = max(ndb, nda)
    if nda == 0 or ndb == 0:
        return _nx.multiply(a, b)
    as_ = a.shape
    bs = b.shape
    if not a.flags.contiguous:
        a = reshape(a, as_)
    if not b.flags.contiguous:
        b = reshape(b, bs)
    as_ = (1,) * max(0, ndb - nda) + as_
    bs = (1,) * max(0, nda - ndb) + bs
    a_arr = expand_dims(a, axis=tuple(range(ndb - nda)))
    b_arr = expand_dims(b, axis=tuple(range(nda - ndb)))
    a_arr = expand_dims(a_arr, axis=tuple(range(1, nd * 2, 2)))
    b_arr = expand_dims(b_arr, axis=tuple(range(0, nd * 2, 2)))
    result = _nx.multiply(a_arr, b_arr, subok=not is_any_mat)
    result = result.reshape(_nx.multiply(as_, bs))
    return result if not is_any_mat else matrix(result, copy=False)

def _tile_dispatcher(A, reps):
    return (A, reps)

@array_function_dispatch(_tile_dispatcher)
def tile(A, reps):
    try:
        tup = tuple(reps)
    except TypeError:
        tup = (reps,)
    d = len(tup)
    if all((x == 1 for x in tup)) and isinstance(A, _nx.ndarray):
        return _nx.array(A, copy=True, subok=True, ndmin=d)
    else:
        c = _nx.array(A, copy=None, subok=True, ndmin=d)
    if d < c.ndim:
        tup = (1,) * (c.ndim - d) + tup
    shape_out = tuple((s * t for s, t in zip(c.shape, tup)))
    n = c.size
    if n > 0:
        for dim_in, nrep in zip(c.shape, tup):
            if nrep != 1:
                c = c.reshape(-1, n).repeat(nrep, 0)
            n //= dim_in
    return c.reshape(shape_out)
