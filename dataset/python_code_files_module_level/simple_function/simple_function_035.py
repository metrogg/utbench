__all__ = ['polyzero', 'polyone', 'polyx', 'polydomain', 'polyline', 'polyadd', 'polysub', 'polymulx', 'polymul', 'polydiv', 'polypow', 'polyval', 'polyvalfromroots', 'polyder', 'polyint', 'polyfromroots', 'polyvander', 'polyfit', 'polytrim', 'polyroots', 'Polynomial', 'polyval2d', 'polyval3d', 'polygrid2d', 'polygrid3d', 'polyvander2d', 'polyvander3d', 'polycompanion']
import numpy as np
from numpy._core.overrides import array_function_dispatch as _array_function_dispatch
from . import polyutils as pu
from ._polybase import ABCPolyBase
polytrim = pu.trimcoef
polydomain = np.array([-1.0, 1.0])
polyzero = np.array([0])
polyone = np.array([1])
polyx = np.array([0, 1])

def polyline(off, scl):
    if scl != 0:
        return np.array([off, scl])
    else:
        return np.array([off])

def polyfromroots(roots):
    return pu._fromroots(polyline, polymul, roots)

def polyadd(c1, c2):
    return pu._add(c1, c2)

def polysub(c1, c2):
    return pu._sub(c1, c2)

def polymulx(c):
    [c] = pu.as_series([c])
    if len(c) == 1 and c[0] == 0:
        return c
    prd = np.empty(len(c) + 1, dtype=c.dtype)
    prd[0] = c[0] * 0
    prd[1:] = c
    return prd

def polymul(c1, c2):
    [c1, c2] = pu.as_series([c1, c2])
    ret = np.convolve(c1, c2)
    return pu.trimseq(ret)

def polydiv(c1, c2):
    [c1, c2] = pu.as_series([c1, c2])
    if c2[-1] == 0:
        raise ZeroDivisionError
    lc1 = len(c1)
    lc2 = len(c2)
    if lc1 < lc2:
        return (c1[:1] * 0, c1)
    elif lc2 == 1:
        return (c1 / c2[-1], c1[:1] * 0)
    else:
        dlen = lc1 - lc2
        scl = c2[-1]
        c2 = c2[:-1] / scl
        i = dlen
        j = lc1 - 1
        while i >= 0:
            c1[i:j] -= c2 * c1[j]
            i -= 1
            j -= 1
        return (c1[j + 1:] / scl, pu.trimseq(c1[:j + 1]))

def polypow(c, pow, maxpower=None):
    return pu._pow(np.convolve, c, pow, maxpower)

def polyder(c, m=1, scl=1, axis=0):
    c = np.array(c, ndmin=1, copy=True)
    if c.dtype.char in '?bBhHiIlLqQpP':
        c = c + 0.0
    cdt = c.dtype
    cnt = pu._as_int(m, 'the order of derivation')
    iaxis = pu._as_int(axis, 'the axis')
    if cnt < 0:
        raise ValueError('The order of derivation must be non-negative')
    iaxis = np.lib.array_utils.normalize_axis_index(iaxis, c.ndim)
    if cnt == 0:
        return c
    c = np.moveaxis(c, iaxis, 0)
    n = len(c)
    if cnt >= n:
        c = c[:1] * 0
    else:
        for i in range(cnt):
            n = n - 1
            c *= scl
            der = np.empty((n,) + c.shape[1:], dtype=cdt)
            for j in range(n, 0, -1):
                der[j - 1] = j * c[j]
            c = der
    c = np.moveaxis(c, 0, iaxis)
    return c

def polyint(c, m=1, k=[], lbnd=0, scl=1, axis=0):
    c = np.array(c, ndmin=1, copy=True)
    if c.dtype.char in '?bBhHiIlLqQpP':
        c = c + 0.0
    cdt = c.dtype
    if not np.iterable(k):
        k = [k]
    cnt = pu._as_int(m, 'the order of integration')
    iaxis = pu._as_int(axis, 'the axis')
    if cnt < 0:
        raise ValueError('The order of integration must be non-negative')
    if len(k) > cnt:
        raise ValueError('Too many integration constants')
    if np.ndim(lbnd) != 0:
        raise ValueError('lbnd must be a scalar.')
    if np.ndim(scl) != 0:
        raise ValueError('scl must be a scalar.')
    iaxis = np.lib.array_utils.normalize_axis_index(iaxis, c.ndim)
    if cnt == 0:
        return c
    k = list(k) + [0] * (cnt - len(k))
    c = np.moveaxis(c, iaxis, 0)
    for i in range(cnt):
        n = len(c)
        c *= scl
        if n == 1 and np.all(c[0] == 0):
            c[0] += k[i]
        else:
            tmp = np.empty((n + 1,) + c.shape[1:], dtype=cdt)
            tmp[0] = c[0] * 0
            tmp[1] = c[0]
            for j in range(1, n):
                tmp[j + 1] = c[j] / (j + 1)
            tmp[0] += k[i] - polyval(lbnd, tmp)
            c = tmp
    c = np.moveaxis(c, 0, iaxis)
    return c

def polyval(x, c, tensor=True):
    c = np.array(c, ndmin=1, copy=None)
    if c.dtype.char in '?bBhHiIlLqQpP':
        c = c + 0.0
    if isinstance(x, (tuple, list)):
        x = np.asarray(x)
    if isinstance(x, np.ndarray) and tensor:
        c = c.reshape(c.shape + (1,) * x.ndim)
    c0 = c[-1] + x * 0
    for i in range(2, len(c) + 1):
        c0 = c[-i] + c0 * x
    return c0

def polyvalfromroots(x, r, tensor=True):
    r = np.array(r, ndmin=1, copy=None)
    if r.dtype.char in '?bBhHiIlLqQpP':
        r = r.astype(np.double)
    if isinstance(x, (tuple, list)):
        x = np.asarray(x)
    if isinstance(x, np.ndarray):
        if tensor:
            r = r.reshape(r.shape + (1,) * x.ndim)
        elif x.ndim >= r.ndim:
            raise ValueError('x.ndim must be < r.ndim when tensor == False')
    return np.prod(x - r, axis=0)

def _polyval2d_dispatcher(x, y, c):
    return (x, y, c)

def _polygrid2d_dispatcher(x, y, c):
    return (x, y, c)

@_array_function_dispatch(_polyval2d_dispatcher)
def polyval2d(x, y, c):
    return pu._valnd(polyval, c, x, y)

@_array_function_dispatch(_polygrid2d_dispatcher)
def polygrid2d(x, y, c):
    return pu._gridnd(polyval, c, x, y)

def polyval3d(x, y, z, c):
    return pu._valnd(polyval, c, x, y, z)

def polygrid3d(x, y, z, c):
    return pu._gridnd(polyval, c, x, y, z)

def polyvander(x, deg):
    ideg = pu._as_int(deg, 'deg')
    if ideg < 0:
        raise ValueError('deg must be non-negative')
    x = np.array(x, copy=None, ndmin=1) + 0.0
    dims = (ideg + 1,) + x.shape
    dtyp = x.dtype
    v = np.empty(dims, dtype=dtyp)
    v[0] = x * 0 + 1
    if ideg > 0:
        v[1] = x
        for i in range(2, ideg + 1):
            v[i] = v[i - 1] * x
    return np.moveaxis(v, 0, -1)

def polyvander2d(x, y, deg):
    return pu._vander_nd_flat((polyvander, polyvander), (x, y), deg)

def polyvander3d(x, y, z, deg):
    return pu._vander_nd_flat((polyvander, polyvander, polyvander), (x, y, z), deg)

def polyfit(x, y, deg, rcond=None, full=False, w=None):
    return pu._fit(polyvander, x, y, deg, rcond, full, w)

def polycompanion(c):
    [c] = pu.as_series([c])
    if len(c) < 2:
        raise ValueError('Series must have maximum degree of at least 1.')
    if len(c) == 2:
        return np.array([[-c[0] / c[1]]])
    n = len(c) - 1
    mat = np.zeros((n, n), dtype=c.dtype)
    bot = mat.reshape(-1)[n::n + 1]
    bot[...] = 1
    mat[:, -1] -= c[:-1] / c[-1]
    return mat

def polyroots(c):
    [c] = pu.as_series([c])
    if len(c) < 2:
        return np.array([], dtype=c.dtype)
    if len(c) == 2:
        return np.array([-c[0] / c[1]])
    m = polycompanion(c)
    r = np.linalg.eigvals(m)
    r.sort()
    return r

class Polynomial(ABCPolyBase):
    _add = staticmethod(polyadd)
    _sub = staticmethod(polysub)
    _mul = staticmethod(polymul)
    _div = staticmethod(polydiv)
    _pow = staticmethod(polypow)
    _val = staticmethod(polyval)
    _int = staticmethod(polyint)
    _der = staticmethod(polyder)
    _fit = staticmethod(polyfit)
    _line = staticmethod(polyline)
    _roots = staticmethod(polyroots)
    _fromroots = staticmethod(polyfromroots)
    domain = np.array(polydomain)
    window = np.array(polydomain)
    basis_name = None

    @classmethod
    def _str_term_unicode(cls, i, arg_str):
        if i == '1':
            return f'·{arg_str}'
        else:
            return f'·{arg_str}{i.translate(cls._superscript_mapping)}'

    @staticmethod
    def _str_term_ascii(i, arg_str):
        if i == '1':
            return f' {arg_str}'
        else:
            return f' {arg_str}**{i}'

    @staticmethod
    def _repr_latex_term(i, arg_str, needs_parens):
        if needs_parens:
            arg_str = f'\\left({arg_str}\\right)'
        if i == 0:
            return '1'
        elif i == 1:
            return arg_str
        else:
            return f'{arg_str}^{{{i}}}'
