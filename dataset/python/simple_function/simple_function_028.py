import numpy as np
from . import polyutils as pu
from ._polybase import ABCPolyBase
__all__ = ['hermzero', 'hermone', 'hermx', 'hermdomain', 'hermline', 'hermadd', 'hermsub', 'hermmulx', 'hermmul', 'hermdiv', 'hermpow', 'hermval', 'hermder', 'hermint', 'herm2poly', 'poly2herm', 'hermfromroots', 'hermvander', 'hermfit', 'hermtrim', 'hermroots', 'Hermite', 'hermval2d', 'hermval3d', 'hermgrid2d', 'hermgrid3d', 'hermvander2d', 'hermvander3d', 'hermcompanion', 'hermgauss', 'hermweight']
hermtrim = pu.trimcoef

def poly2herm(pol):
    [pol] = pu.as_series([pol])
    deg = len(pol) - 1
    res = 0
    for i in range(deg, -1, -1):
        res = hermadd(hermmulx(res), pol[i])
    return res

def herm2poly(c):
    from .polynomial import polyadd, polymulx, polysub
    [c] = pu.as_series([c])
    n = len(c)
    if n == 1:
        return c
    if n == 2:
        c[1] *= 2
        return c
    else:
        c0 = c[-2]
        c1 = c[-1]
        for i in range(n - 1, 1, -1):
            tmp = c0
            c0 = polysub(c[i - 2], c1 * (2 * (i - 1)))
            c1 = polyadd(tmp, polymulx(c1) * 2)
        return polyadd(c0, polymulx(c1) * 2)
hermdomain = np.array([-1.0, 1.0])
hermzero = np.array([0])
hermone = np.array([1])
hermx = np.array([0, 1 / 2])

def hermline(off, scl):
    if scl != 0:
        return np.array([off, scl / 2])
    else:
        return np.array([off])

def hermfromroots(roots):
    return pu._fromroots(hermline, hermmul, roots)

def hermadd(c1, c2):
    return pu._add(c1, c2)

def hermsub(c1, c2):
    return pu._sub(c1, c2)

def hermmulx(c):
    [c] = pu.as_series([c])
    if len(c) == 1 and c[0] == 0:
        return c
    prd = np.empty(len(c) + 1, dtype=c.dtype)
    prd[0] = c[0] * 0
    prd[1] = c[0] / 2
    for i in range(1, len(c)):
        prd[i + 1] = c[i] / 2
        prd[i - 1] += c[i] * i
    return prd

def hermmul(c1, c2):
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
            c0 = hermsub(c[-i] * xs, c1 * (2 * (nd - 1)))
            c1 = hermadd(tmp, hermmulx(c1) * 2)
    return hermadd(c0, hermmulx(c1) * 2)

def hermdiv(c1, c2):
    return pu._div(hermmul, c1, c2)

def hermpow(c, pow, maxpower=16):
    return pu._pow(hermmul, c, pow, maxpower)

def hermder(c, m=1, scl=1, axis=0):
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
        c = c[:1] * 0
    else:
        for i in range(cnt):
            n = n - 1
            c *= scl
            der = np.empty((n,) + c.shape[1:], dtype=c.dtype)
            for j in range(n, 0, -1):
                der[j - 1] = 2 * j * c[j]
            c = der
    c = np.moveaxis(c, 0, iaxis)
    return c

def hermint(c, m=1, k=[], lbnd=0, scl=1, axis=0):
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
            tmp[1] = c[0] / 2
            for j in range(1, n):
                tmp[j + 1] = c[j] / (2 * (j + 1))
            tmp[0] += k[i] - hermval(lbnd, tmp)
            c = tmp
    c = np.moveaxis(c, 0, iaxis)
    return c

def hermval(x, c, tensor=True):
    c = np.array(c, ndmin=1, copy=None)
    if c.dtype.char in '?bBhHiIlLqQpP':
        c = c.astype(np.double)
    if isinstance(x, (tuple, list)):
        x = np.asarray(x)
    if isinstance(x, np.ndarray) and tensor:
        c = c.reshape(c.shape + (1,) * x.ndim)
    x2 = x * 2
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
            c0 = c[-i] - c1 * (2 * (nd - 1))
            c1 = tmp + c1 * x2
    return c0 + c1 * x2

def hermval2d(x, y, c):
    return pu._valnd(hermval, c, x, y)

def hermgrid2d(x, y, c):
    return pu._gridnd(hermval, c, x, y)

def hermval3d(x, y, z, c):
    return pu._valnd(hermval, c, x, y, z)

def hermgrid3d(x, y, z, c):
    return pu._gridnd(hermval, c, x, y, z)

def hermvander(x, deg):
    ideg = pu._as_int(deg, 'deg')
    if ideg < 0:
        raise ValueError('deg must be non-negative')
    x = np.array(x, copy=None, ndmin=1) + 0.0
    dims = (ideg + 1,) + x.shape
    dtyp = x.dtype
    v = np.empty(dims, dtype=dtyp)
    v[0] = x * 0 + 1
    if ideg > 0:
        x2 = x * 2
        v[1] = x2
        for i in range(2, ideg + 1):
            v[i] = v[i - 1] * x2 - v[i - 2] * (2 * (i - 1))
    return np.moveaxis(v, 0, -1)

def hermvander2d(x, y, deg):
    return pu._vander_nd_flat((hermvander, hermvander), (x, y), deg)

def hermvander3d(x, y, z, deg):
    return pu._vander_nd_flat((hermvander, hermvander, hermvander), (x, y, z), deg)

def hermfit(x, y, deg, rcond=None, full=False, w=None):
    return pu._fit(hermvander, x, y, deg, rcond, full, w)

def hermcompanion(c):
    [c] = pu.as_series([c])
    if len(c) < 2:
        raise ValueError('Series must have maximum degree of at least 1.')
    if len(c) == 2:
        return np.array([[-0.5 * c[0] / c[1]]])
    n = len(c) - 1
    mat = np.zeros((n, n), dtype=c.dtype)
    scl = np.hstack((1.0, 1.0 / np.sqrt(2.0 * np.arange(n - 1, 0, -1))))
    scl = np.multiply.accumulate(scl)[::-1]
    top = mat.reshape(-1)[1::n + 1]
    bot = mat.reshape(-1)[n::n + 1]
    top[...] = np.sqrt(0.5 * np.arange(1, n))
    bot[...] = top
    mat[:, -1] -= scl * c[:-1] / (2.0 * c[-1])
    return mat

def hermroots(c):
    [c] = pu.as_series([c])
    if len(c) <= 1:
        return np.array([], dtype=c.dtype)
    if len(c) == 2:
        return np.array([-0.5 * c[0] / c[1]])
    m = hermcompanion(c)[::-1, ::-1]
    r = np.linalg.eigvals(m)
    r.sort()
    return r

def _normed_hermite_n(x, n):
    if n == 0:
        return np.full(x.shape, 1 / np.sqrt(np.sqrt(np.pi)))
    c0 = 0.0
    c1 = 1.0 / np.sqrt(np.sqrt(np.pi))
    nd = float(n)
    for i in range(n - 1):
        tmp = c0
        c0 = -c1 * np.sqrt((nd - 1.0) / nd)
        c1 = tmp + c1 * x * np.sqrt(2.0 / nd)
        nd = nd - 1.0
    return c0 + c1 * x * np.sqrt(2)

def hermgauss(deg):
    ideg = pu._as_int(deg, 'deg')
    if ideg <= 0:
        raise ValueError('deg must be a positive integer')
    c = np.array([0] * deg + [1], dtype=np.float64)
    m = hermcompanion(c)
    x = np.linalg.eigvalsh(m)
    dy = _normed_hermite_n(x, ideg)
    df = _normed_hermite_n(x, ideg - 1) * np.sqrt(2 * ideg)
    x -= dy / df
    fm = _normed_hermite_n(x, ideg - 1)
    fm /= np.abs(fm).max()
    w = 1 / (fm * fm)
    w = (w + w[::-1]) / 2
    x = (x - x[::-1]) / 2
    w *= np.sqrt(np.pi) / w.sum()
    return (x, w)

def hermweight(x):
    w = np.exp(-x ** 2)
    return w

class Hermite(ABCPolyBase):
    _add = staticmethod(hermadd)
    _sub = staticmethod(hermsub)
    _mul = staticmethod(hermmul)
    _div = staticmethod(hermdiv)
    _pow = staticmethod(hermpow)
    _val = staticmethod(hermval)
    _int = staticmethod(hermint)
    _der = staticmethod(hermder)
    _fit = staticmethod(hermfit)
    _line = staticmethod(hermline)
    _roots = staticmethod(hermroots)
    _fromroots = staticmethod(hermfromroots)
    domain = np.array(hermdomain)
    window = np.array(hermdomain)
    basis_name = 'H'
