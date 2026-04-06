import numpy as np
from . import polyutils as pu
from ._polybase import ABCPolyBase
__all__ = ['legzero', 'legone', 'legx', 'legdomain', 'legline', 'legadd', 'legsub', 'legmulx', 'legmul', 'legdiv', 'legpow', 'legval', 'legder', 'legint', 'leg2poly', 'poly2leg', 'legfromroots', 'legvander', 'legfit', 'legtrim', 'legroots', 'Legendre', 'legval2d', 'legval3d', 'leggrid2d', 'leggrid3d', 'legvander2d', 'legvander3d', 'legcompanion', 'leggauss', 'legweight']
legtrim = pu.trimcoef

def poly2leg(pol):
    [pol] = pu.as_series([pol])
    deg = len(pol) - 1
    res = 0
    for i in range(deg, -1, -1):
        res = legadd(legmulx(res), pol[i])
    return res

def leg2poly(c):
    from .polynomial import polyadd, polymulx, polysub
    [c] = pu.as_series([c])
    n = len(c)
    if n < 3:
        return c
    else:
        c0 = c[-2]
        c1 = c[-1]
        for i in range(n - 1, 1, -1):
            tmp = c0
            c0 = polysub(c[i - 2], c1 * (i - 1) / i)
            c1 = polyadd(tmp, polymulx(c1) * (2 * i - 1) / i)
        return polyadd(c0, polymulx(c1))
legdomain = np.array([-1.0, 1.0])
legzero = np.array([0])
legone = np.array([1])
legx = np.array([0, 1])

def legline(off, scl):
    if scl != 0:
        return np.array([off, scl])
    else:
        return np.array([off])

def legfromroots(roots):
    return pu._fromroots(legline, legmul, roots)

def legadd(c1, c2):
    return pu._add(c1, c2)

def legsub(c1, c2):
    return pu._sub(c1, c2)

def legmulx(c):
    [c] = pu.as_series([c])
    if len(c) == 1 and c[0] == 0:
        return c
    prd = np.empty(len(c) + 1, dtype=c.dtype)
    prd[0] = c[0] * 0
    prd[1] = c[0]
    for i in range(1, len(c)):
        j = i + 1
        k = i - 1
        s = i + j
        prd[j] = c[i] * j / s
        prd[k] += c[i] * i / s
    return prd

def legmul(c1, c2):
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
            c0 = legsub(c[-i] * xs, c1 * (nd - 1) / nd)
            c1 = legadd(tmp, legmulx(c1) * (2 * nd - 1) / nd)
    return legadd(c0, legmulx(c1))

def legdiv(c1, c2):
    return pu._div(legmul, c1, c2)

def legpow(c, pow, maxpower=16):
    return pu._pow(legmul, c, pow, maxpower)

def legder(c, m=1, scl=1, axis=0):
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
            for j in range(n, 2, -1):
                der[j - 1] = (2 * j - 1) * c[j]
                c[j - 2] += c[j]
            if n > 1:
                der[1] = 3 * c[2]
            der[0] = c[1]
            c = der
    c = np.moveaxis(c, 0, iaxis)
    return c

def legint(c, m=1, k=[], lbnd=0, scl=1, axis=0):
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
            if n > 1:
                tmp[2] = c[1] / 3
            for j in range(2, n):
                t = c[j] / (2 * j + 1)
                tmp[j + 1] = t
                tmp[j - 1] -= t
            tmp[0] += k[i] - legval(lbnd, tmp)
            c = tmp
    c = np.moveaxis(c, 0, iaxis)
    return c

def legval(x, c, tensor=True):
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
            c0 = c[-i] - c1 * ((nd - 1) / nd)
            c1 = tmp + c1 * x * ((2 * nd - 1) / nd)
    return c0 + c1 * x

def legval2d(x, y, c):
    return pu._valnd(legval, c, x, y)

def leggrid2d(x, y, c):
    return pu._gridnd(legval, c, x, y)

def legval3d(x, y, z, c):
    return pu._valnd(legval, c, x, y, z)

def leggrid3d(x, y, z, c):
    return pu._gridnd(legval, c, x, y, z)

def legvander(x, deg):
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
            v[i] = (v[i - 1] * x * (2 * i - 1) - v[i - 2] * (i - 1)) / i
    return np.moveaxis(v, 0, -1)

def legvander2d(x, y, deg):
    return pu._vander_nd_flat((legvander, legvander), (x, y), deg)

def legvander3d(x, y, z, deg):
    return pu._vander_nd_flat((legvander, legvander, legvander), (x, y, z), deg)

def legfit(x, y, deg, rcond=None, full=False, w=None):
    return pu._fit(legvander, x, y, deg, rcond, full, w)

def legcompanion(c):
    [c] = pu.as_series([c])
    if len(c) < 2:
        raise ValueError('Series must have maximum degree of at least 1.')
    if len(c) == 2:
        return np.array([[-c[0] / c[1]]])
    n = len(c) - 1
    mat = np.zeros((n, n), dtype=c.dtype)
    scl = 1.0 / np.sqrt(2 * np.arange(n) + 1)
    top = mat.reshape(-1)[1::n + 1]
    bot = mat.reshape(-1)[n::n + 1]
    top[...] = np.arange(1, n) * scl[:n - 1] * scl[1:n]
    bot[...] = top
    mat[:, -1] -= c[:-1] / c[-1] * (scl / scl[-1]) * (n / (2 * n - 1))
    return mat

def legroots(c):
    [c] = pu.as_series([c])
    if len(c) < 2:
        return np.array([], dtype=c.dtype)
    if len(c) == 2:
        return np.array([-c[0] / c[1]])
    m = legcompanion(c)[::-1, ::-1]
    r = np.linalg.eigvals(m)
    r.sort()
    return r

def leggauss(deg):
    ideg = pu._as_int(deg, 'deg')
    if ideg <= 0:
        raise ValueError('deg must be a positive integer')
    c = np.array([0] * deg + [1])
    m = legcompanion(c)
    x = np.linalg.eigvalsh(m)
    dy = legval(x, c)
    df = legval(x, legder(c))
    x -= dy / df
    fm = legval(x, c[1:])
    fm /= np.abs(fm).max()
    df /= np.abs(df).max()
    w = 1 / (fm * df)
    w = (w + w[::-1]) / 2
    x = (x - x[::-1]) / 2
    w *= 2.0 / w.sum()
    return (x, w)

def legweight(x):
    w = x * 0.0 + 1.0
    return w

class Legendre(ABCPolyBase):
    _add = staticmethod(legadd)
    _sub = staticmethod(legsub)
    _mul = staticmethod(legmul)
    _div = staticmethod(legdiv)
    _pow = staticmethod(legpow)
    _val = staticmethod(legval)
    _int = staticmethod(legint)
    _der = staticmethod(legder)
    _fit = staticmethod(legfit)
    _line = staticmethod(legline)
    _roots = staticmethod(legroots)
    _fromroots = staticmethod(legfromroots)
    domain = np.array(legdomain)
    window = np.array(legdomain)
    basis_name = 'P'
