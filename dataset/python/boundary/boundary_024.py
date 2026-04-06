import contextlib
import functools
import operator
import warnings
import numpy as np
from numpy._core import overrides
__all__ = ['histogram', 'histogramdd', 'histogram_bin_edges']
array_function_dispatch = functools.partial(overrides.array_function_dispatch, module='numpy')
_range = range

def _ptp(x):
    return _unsigned_subtract(x.max(), x.min())

def _hist_bin_sqrt(x, range):
    del range
    return _ptp(x) / np.sqrt(x.size)

def _hist_bin_sturges(x, range):
    del range
    return _ptp(x) / (np.log2(x.size) + 1.0)

def _hist_bin_rice(x, range):
    del range
    return _ptp(x) / (2.0 * x.size ** (1.0 / 3))

def _hist_bin_scott(x, range):
    del range
    return (24.0 * np.pi ** 0.5 / x.size) ** (1.0 / 3.0) * np.std(x)

def _hist_bin_stone(x, range):
    n = x.size
    ptp_x = _ptp(x)
    if n <= 1 or ptp_x == 0:
        return 0

    def jhat(nbins):
        hh = ptp_x / nbins
        p_k = np.histogram(x, bins=nbins, range=range)[0] / n
        return (2 - (n + 1) * p_k.dot(p_k)) / hh
    nbins_upper_bound = max(100, int(np.sqrt(n)))
    nbins = min(_range(1, nbins_upper_bound + 1), key=jhat)
    if nbins == nbins_upper_bound:
        warnings.warn('The number of bins estimated may be suboptimal.', RuntimeWarning, stacklevel=3)
    return ptp_x / nbins

def _hist_bin_doane(x, range):
    del range
    if x.size > 2:
        sg1 = np.sqrt(6.0 * (x.size - 2) / ((x.size + 1.0) * (x.size + 3)))
        sigma = np.std(x)
        if sigma > 0.0:
            temp = x - np.mean(x)
            np.true_divide(temp, sigma, temp)
            np.power(temp, 3, temp)
            g1 = np.mean(temp)
            return _ptp(x) / (1.0 + np.log2(x.size) + np.log2(1.0 + np.absolute(g1) / sg1))
    return 0.0

def _hist_bin_fd(x, range):
    del range
    iqr = np.subtract(*np.percentile(x, [75, 25]))
    return 2.0 * iqr * x.size ** (-1.0 / 3.0)

def _hist_bin_auto(x, range):
    fd_bw = _hist_bin_fd(x, range)
    sturges_bw = _hist_bin_sturges(x, range)
    sqrt_bw = _hist_bin_sqrt(x, range)
    fd_bw_corrected = max(fd_bw, sqrt_bw / 2)
    return min(fd_bw_corrected, sturges_bw)
_hist_bin_selectors = {'stone': _hist_bin_stone, 'auto': _hist_bin_auto, 'doane': _hist_bin_doane, 'fd': _hist_bin_fd, 'rice': _hist_bin_rice, 'scott': _hist_bin_scott, 'sqrt': _hist_bin_sqrt, 'sturges': _hist_bin_sturges}

def _ravel_and_check_weights(a, weights):
    a = np.asarray(a)
    if a.dtype == np.bool:
        msg = f'Converting input from {a.dtype} to {np.uint8} for compatibility.'
        warnings.warn(msg, RuntimeWarning, stacklevel=3)
        a = a.astype(np.uint8)
    if weights is not None:
        weights = np.asarray(weights)
        if weights.shape != a.shape:
            raise ValueError('weights should have the same shape as a.')
        weights = weights.ravel()
    a = a.ravel()
    return (a, weights)

def _get_outer_edges(a, range):
    if range is not None:
        first_edge, last_edge = range
        if first_edge > last_edge:
            raise ValueError('max must be larger than min in range parameter.')
        if not (np.isfinite(first_edge) and np.isfinite(last_edge)):
            raise ValueError(f'supplied range of [{first_edge}, {last_edge}] is not finite')
    elif a.size == 0:
        first_edge, last_edge = (0, 1)
    else:
        first_edge, last_edge = (a.min(), a.max())
        if not (np.isfinite(first_edge) and np.isfinite(last_edge)):
            raise ValueError(f'autodetected range of [{first_edge}, {last_edge}] is not finite')
    if first_edge == last_edge:
        first_edge = first_edge - 0.5
        last_edge = last_edge + 0.5
    return (first_edge, last_edge)

def _unsigned_subtract(a, b):
    signed_to_unsigned = {np.byte: np.ubyte, np.short: np.ushort, np.intc: np.uintc, np.int_: np.uint, np.longlong: np.ulonglong}
    dt = np.result_type(a, b)
    try:
        unsigned_dt = signed_to_unsigned[dt.type]
    except KeyError:
        return np.subtract(a, b, dtype=dt)
    else:
        return np.subtract(np.asarray(a, dtype=dt), np.asarray(b, dtype=dt), casting='unsafe', dtype=unsigned_dt)

def _get_bin_edges(a, bins, range, weights):
    n_equal_bins = None
    bin_edges = None
    if isinstance(bins, str):
        bin_name = bins
        if bin_name not in _hist_bin_selectors:
            raise ValueError(f'{bin_name!r} is not a valid estimator for `bins`')
        if weights is not None:
            raise TypeError('Automated estimation of the number of bins is not supported for weighted data')
        first_edge, last_edge = _get_outer_edges(a, range)
        if range is not None:
            keep = a >= first_edge
            keep &= a <= last_edge
            if not np.logical_and.reduce(keep):
                a = a[keep]
        if a.size == 0:
            n_equal_bins = 1
        else:
            width = _hist_bin_selectors[bin_name](a, (first_edge, last_edge))
            if width:
                if np.issubdtype(a.dtype, np.integer) and width < 1:
                    width = 1
                delta = _unsigned_subtract(last_edge, first_edge)
                n_equal_bins = int(np.ceil(delta / width))
            else:
                n_equal_bins = 1
    elif np.ndim(bins) == 0:
        try:
            n_equal_bins = operator.index(bins)
        except TypeError as e:
            raise TypeError('`bins` must be an integer, a string, or an array') from e
        if n_equal_bins < 1:
            raise ValueError('`bins` must be positive, when an integer')
        first_edge, last_edge = _get_outer_edges(a, range)
    elif np.ndim(bins) == 1:
        bin_edges = np.asarray(bins)
        if np.any(bin_edges[:-1] > bin_edges[1:]):
            raise ValueError('`bins` must increase monotonically, when an array')
    else:
        raise ValueError('`bins` must be 1d, when an array')
    if n_equal_bins is not None:
        bin_type = np.result_type(first_edge, last_edge, a)
        if np.issubdtype(bin_type, np.integer):
            bin_type = np.result_type(bin_type, float)
        bin_edges = np.linspace(first_edge, last_edge, n_equal_bins + 1, endpoint=True, dtype=bin_type)
        if np.any(bin_edges[:-1] >= bin_edges[1:]):
            raise ValueError(f'Too many bins for data range. Cannot create {n_equal_bins} finite-sized bins.')
        return (bin_edges, (first_edge, last_edge, n_equal_bins))
    else:
        return (bin_edges, None)

def _search_sorted_inclusive(a, v):
    return np.concatenate((a.searchsorted(v[:-1], 'left'), a.searchsorted(v[-1:], 'right')))

def _histogram_bin_edges_dispatcher(a, bins=None, range=None, weights=None):
    return (a, bins, weights)

@array_function_dispatch(_histogram_bin_edges_dispatcher)
def histogram_bin_edges(a, bins=10, range=None, weights=None):
    a, weights = _ravel_and_check_weights(a, weights)
    bin_edges, _ = _get_bin_edges(a, bins, range, weights)
    return bin_edges

def _histogram_dispatcher(a, bins=None, range=None, density=None, weights=None):
    return (a, bins, weights)

@array_function_dispatch(_histogram_dispatcher)
def histogram(a, bins=10, range=None, density=None, weights=None):
    a, weights = _ravel_and_check_weights(a, weights)
    bin_edges, uniform_bins = _get_bin_edges(a, bins, range, weights)
    if weights is None:
        ntype = np.dtype(np.intp)
    else:
        ntype = weights.dtype
    BLOCK = 65536
    simple_weights = weights is None or np.can_cast(weights.dtype, np.double) or np.can_cast(weights.dtype, complex)
    if uniform_bins is not None and simple_weights:
        first_edge, last_edge, n_equal_bins = uniform_bins
        n = np.zeros(n_equal_bins, ntype)
        norm_numerator = n_equal_bins
        norm_denom = _unsigned_subtract(last_edge, first_edge)
        for i in _range(0, len(a), BLOCK):
            tmp_a = a[i:i + BLOCK]
            if weights is None:
                tmp_w = None
            else:
                tmp_w = weights[i:i + BLOCK]
            keep = tmp_a >= first_edge
            keep &= tmp_a <= last_edge
            if not np.logical_and.reduce(keep):
                tmp_a = tmp_a[keep]
                if tmp_w is not None:
                    tmp_w = tmp_w[keep]
            tmp_a = tmp_a.astype(bin_edges.dtype, copy=False)
            f_indices = _unsigned_subtract(tmp_a, first_edge) / norm_denom * norm_numerator
            indices = f_indices.astype(np.intp)
            indices[indices == n_equal_bins] -= 1
            decrement = tmp_a < bin_edges[indices]
            indices[decrement] -= 1
            increment = (tmp_a >= bin_edges[indices + 1]) & (indices != n_equal_bins - 1)
            indices[increment] += 1
            if ntype.kind == 'c':
                n.real += np.bincount(indices, weights=tmp_w.real, minlength=n_equal_bins)
                n.imag += np.bincount(indices, weights=tmp_w.imag, minlength=n_equal_bins)
            else:
                n += np.bincount(indices, weights=tmp_w, minlength=n_equal_bins).astype(ntype)
    else:
        cum_n = np.zeros(bin_edges.shape, ntype)
        if weights is None:
            for i in _range(0, len(a), BLOCK):
                sa = np.sort(a[i:i + BLOCK])
                cum_n += _search_sorted_inclusive(sa, bin_edges)
        else:
            zero = np.zeros(1, dtype=ntype)
            for i in _range(0, len(a), BLOCK):
                tmp_a = a[i:i + BLOCK]
                tmp_w = weights[i:i + BLOCK]
                sorting_index = np.argsort(tmp_a)
                sa = tmp_a[sorting_index]
                sw = tmp_w[sorting_index]
                cw = np.concatenate((zero, sw.cumsum()))
                bin_index = _search_sorted_inclusive(sa, bin_edges)
                cum_n += cw[bin_index]
        n = np.diff(cum_n)
    if density:
        db = np.array(np.diff(bin_edges), float)
        return (n / db / n.sum(), bin_edges)
    return (n, bin_edges)

def _histogramdd_dispatcher(sample, bins=None, range=None, density=None, weights=None):
    if hasattr(sample, 'shape'):
        yield sample
    else:
        yield from sample
    with contextlib.suppress(TypeError):
        yield from bins
    yield weights

@array_function_dispatch(_histogramdd_dispatcher)
def histogramdd(sample, bins=10, range=None, density=None, weights=None):
    try:
        N, D = sample.shape
    except (AttributeError, ValueError):
        sample = np.atleast_2d(sample).T
        N, D = sample.shape
    nbin = np.empty(D, np.intp)
    edges = D * [None]
    dedges = D * [None]
    if weights is not None:
        weights = np.asarray(weights)
    try:
        M = len(bins)
        if M != D:
            raise ValueError('The dimension of bins must be equal to the dimension of the sample x.')
    except TypeError:
        bins = D * [bins]
    if range is None:
        range = (None,) * D
    elif len(range) != D:
        raise ValueError('range argument must have one entry per dimension')
    for i in _range(D):
        if np.ndim(bins[i]) == 0:
            if bins[i] < 1:
                raise ValueError(f'`bins[{i}]` must be positive, when an integer')
            smin, smax = _get_outer_edges(sample[:, i], range[i])
            try:
                n = operator.index(bins[i])
            except TypeError as e:
                raise TypeError(f'`bins[{i}]` must be an integer, when a scalar') from e
            edges[i] = np.linspace(smin, smax, n + 1)
        elif np.ndim(bins[i]) == 1:
            edges[i] = np.asarray(bins[i])
            if np.any(edges[i][:-1] > edges[i][1:]):
                raise ValueError(f'`bins[{i}]` must be monotonically increasing, when an array')
        else:
            raise ValueError(f'`bins[{i}]` must be a scalar or 1d array')
        nbin[i] = len(edges[i]) + 1
        dedges[i] = np.diff(edges[i])
    Ncount = tuple((np.searchsorted(edges[i], sample[:, i], side='right') for i in _range(D)))
    for i in _range(D):
        on_edge = sample[:, i] == edges[i][-1]
        Ncount[i][on_edge] -= 1
    xy = np.ravel_multi_index(Ncount, nbin)
    hist = np.bincount(xy, weights, minlength=nbin.prod())
    hist = hist.reshape(nbin)
    hist = hist.astype(float, casting='safe')
    core = D * (slice(1, -1),)
    hist = hist[core]
    if density:
        s = hist.sum()
        for i in _range(D):
            shape = np.ones(D, int)
            shape[i] = nbin[i] - 2
            hist = hist / dedges[i].reshape(shape)
        hist /= s
    if (hist.shape != nbin - 2).any():
        raise RuntimeError('Internal Shape Error')
    return (hist, edges)
