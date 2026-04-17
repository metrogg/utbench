import numpy as np
from . import polyutils as pu
from ._polybase import ABCPolyBase
__all__ = ['hermezero', 'hermeone', 'hermex', 'hermedomain', 'hermeline', 'hermeadd', 'hermesub', 'hermemulx', 'hermemul', 'hermediv', 'hermepow', 'hermeval', 'hermeder', 'hermeint', 'herme2poly', 'poly2herme', 'hermefromroots', 'hermevander', 'hermefit', 'hermetrim', 'hermeroots', 'HermiteE', 'hermeval2d', 'hermeval3d', 'hermegrid2d', 'hermegrid3d', 'hermevander2d', 'hermevander3d', 'hermecompanion', 'hermegauss', 'hermeweight']
hermetrim = pu.trimcoef

def poly2herme(pol):
    [pol] = pu.as_series([pol])
    deg = len(pol) - 1
    res = 0
    for i in range(deg, -1, -1):
        res = hermeadd(hermemulx(res), pol[i])
    return res

def herme2poly(c):
    from .polynomial import polyadd, polymulx, polysub
    [c] = pu.as_series([c])
    n = len(c)
    if n == 1:
        return c
    if n == 2:
        return c
    else:
        c0 = c[-2]
        c1 = c[-1]
        for i in range(n - 1, 1, -1):
            tmp = c0
            c0 = polysub(c[i - 2], c1 * (i - 1))
            c1 = polyadd(tmp, polymulx(c1))
        return polyadd(c0, polymulx(c1))
hermedomain = np.array([-1.0, 1.0])
hermezero = np.array([0])
hermeone = np.array([1])
hermex = np.array([0, 1])

def hermeline(off, scl):
    if scl != 0:
        return np.array([off, scl])
    else:
        return np.array([off])

def hermefromroots(roots):
    return pu._fromroots(hermeline, hermemul, roots)

def hermeadd(c1, c2):
    return pu._add(c1, c2)

def hermesub(c1, c2):
    return pu._sub(c1, c2)

def hermemulx(c):
    [c] = pu.as_series([c])
    if len(c) == 1 and c[0] == 0:
        return c
    prd = np.empty(len(c) + 1, dtype=c.dtype)
    prd[0] = c[0] * 0
    prd[1] = c[0]
    for i in range(1, len(c)):
        prd[i + 1] = c[i]
        prd[i - 1] += c[i] * i
    return prd

def hermemul(c1, c2):
    [c1, c2] = pu.as_series([c1, c2])
    if len(c1) > len(c2):
        c = c2
        xs = c1
    else:
        c = c1
        xs = c2
    if len(c) == 1:
        c0 = c[0] * xs
        c1 = 0
    elif len(c) == 2:
        c0 = c[0] * xs
        c1 = c[1] * xs
    else:
        nd = len(c)
        c0 = c[-2] * xs
        c1 = c[-1] * xs
        for i in range(3, len(c) + 1):
            tmp = c0
            nd = nd - 1
            c0 = hermesub(c[-i] * xs, c1 * (nd - 1))
            c1 = hermeadd(tmp, hermemulx(c1))
    return hermeadd(c0, hermemulx(c1))

def hermediv(c1, c2):
    return pu._div(hermemul, c1, c2)

def hermepow(c, pow, maxpower=16):
    return pu._pow(hermemul, c, pow, maxpower)

def hermeder(c, m=1, scl=1, axis=0):
    c = np.array(c, ndmin=1, copy=True)
    if c.dtype.char in '?bBhHiIlLqQpP':
        c = c.astype(np.double)
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
        return c[:1] * 0
    else:
        for i in range(cnt):
            n = n - 1
            c *= scl
            der = np.empty((n,) + c.shape[1:], dtype=c.dtype)
            for j in range(n, 0, -1):
                der[j - 1] = j * c[j]
            c = der
    c = np.moveaxis(c, 0, iaxis)
    return c

def hermeint(c, m=1, k=[], lbnd=0, scl=1, axis=0):
    c = np.array(c, ndmin=1, copy=True)
    if c.dtype.char in '?bBhHiIlLqQpP':
        c = c.astype(np.double)
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
    c = np.moveaxis(c, iaxis, 0)
    k = list(k) + [0] * (cnt - len(k))
    for i in range(cnt):
        n = len(c)
        c *= scl
        if n == 1 and np.all(c[0] == 0):
            c[0] += k[i]
        else:
            tmp = np.empty((n + 1,) + c.shape[1:], dtype=c.dtype)
            tmp[0] = c[0] * 0
            tmp[1] = c[0]
            for j in range(1, n):
                tmp[j + 1] = c[j] / (j + 1)
            tmp[0] += k[i] - hermeval(lbnd, tmp)
            c = tmp
    c = np.moveaxis(c, 0, iaxis)
    return c

def hermeval(x, c, tensor=True):
    c = np.array(c, ndmin=1, copy=None)
    if c.dtype.char in '?bBhHiIlLqQpP':
        c = c.astype(np.double)
    if isinstance(x, (tuple, list)):
        x = np.asarray(x)
    if isinstance(x, np.ndarray) and tensor:
        c = c.reshape(c.shape + (1,) * x.ndim)
    if len(c) == 1:
        c0 = c[0]
        c1 = 0
    elif len(c) == 2:
        c0 = c[0]
        c1 = c[1]
    else:
        nd = len(c)
        c0 = c[-2]
        c1 = c[-1]
        for i in range(3, len(c) + 1):
            tmp = c0
            nd = nd - 1
            c0 = c[-i] - c1 * (nd - 1)
            c1 = tmp + c1 * x
    return c0 + c1 * x

def hermeval2d(x, y, c):
    return pu._valnd(hermeval, c, x, y)

def hermegrid2d(x, y, c):
    return pu._gridnd(hermeval, c, x, y)

def hermeval3d(x, y, z, c):
    return pu._valnd(hermeval, c, x, y, z)

def hermegrid3d(x, y, z, c):
    return pu._gridnd(hermeval, c, x, y, z)

def hermevander(x, deg):
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
            v[i] = v[i - 1] * x - v[i - 2] * (i - 1)
    return np.moveaxis(v, 0, -1)

def hermevander2d(x, y, deg):
    return pu._vander_nd_flat((hermevander, hermevander), (x, y), deg)

def hermevander3d(x, y, z, deg):
    return pu._vander_nd_flat((hermevander, hermevander, hermevander), (x, y, z), deg)

def hermefit(x, y, deg, rcond=None, full=False, w=None):
    return pu._fit(hermevander, x, y, deg, rcond, full, w)

def hermecompanion(c):
    [c] = pu.as_series([c])
    if len(c) < 2:
        raise ValueError('Series must have maximum degree of at least 1.')
    if len(c) == 2:
        return np.array([[-c[0] / c[1]]])
    n = len(c) - 1
    mat = np.zeros((n, n), dtype=c.dtype)
    scl = np.hstack((1.0, 1.0 / np.sqrt(np.arange(n - 1, 0, -1))))
    scl = np.multiply.accumulate(scl)[::-1]
    top = mat.reshape(-1)[1::n + 1]
    bot = mat.reshape(-1)[n::n + 1]
    top[...] = np.sqrt(np.arange(1, n))
    bot[...] = top
    mat[:, -1] -= scl * c[:-1] / c[-1]
    return mat

def hermeroots(c):
    [c] = pu.as_series([c])
    if len(c) <= 1:
        return np.array([], dtype=c.dtype)
    if len(c) == 2:
        return np.array([-c[0] / c[1]])
    m = hermecompanion(c)[::-1, ::-1]
    r = np.linalg.eigvals(m)
    r.sort()
    return r

def _normed_hermite_e_n(x, n):
    if n == 0:
        return np.full(x.shape, 1 / np.sqrt(np.sqrt(2 * np.pi)))
    c0 = 0.0
    c1 = 1.0 / np.sqrt(np.sqrt(2 * np.pi))
    nd = float(n)
    for i in range(n - 1):
        tmp = c0
        c0 = -c1 * np.sqrt((nd - 1.0) / nd)
        c1 = tmp + c1 * x * np.sqrt(1.0 / nd)
        nd = nd - 1.0
    return c0 + c1 * x

def hermegauss(deg):
    ideg = pu._as_int(deg, 'deg')
    if ideg <= 0:
        raise ValueError('deg must be a positive integer')
    c = np.array([0] * deg + [1])
    m = hermecompanion(c)
    x = np.linalg.eigvalsh(m)
    dy = _normed_hermite_e_n(x, ideg)
    df = _normed_hermite_e_n(x, ideg - 1) * np.sqrt(ideg)
    x -= dy / df
    fm = _normed_hermite_e_n(x, ideg - 1)
    fm /= np.abs(fm).max()
    w = 1 / (fm * fm)
    w = (w + w[::-1]) / 2
    x = (x - x[::-1]) / 2
    w *= np.sqrt(2 * np.pi) / w.sum()
    return (x, w)

def hermeweight(x):
    w = np.exp(-0.5 * x ** 2)
    return w

class HermiteE(ABCPolyBase):
    _add = staticmethod(hermeadd)
    _sub = staticmethod(hermesub)
    _mul = staticmethod(hermemul)
    _div = staticmethod(hermediv)
    _pow = staticmethod(hermepow)
    _val = staticmethod(hermeval)
    _int = staticmethod(hermeint)
    _der = staticmethod(hermeder)
    _fit = staticmethod(hermefit)
    _line = staticmethod(hermeline)
    _roots = staticmethod(hermeroots)
    _fromroots = staticmethod(hermefromroots)
    domain = np.array(hermedomain)
    window = np.array(hermedomain)
    basis_name = 'He'
