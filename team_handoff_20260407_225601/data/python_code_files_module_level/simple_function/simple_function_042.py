import numpy as np
from . import polyutils as pu
from ._polybase import ABCPolyBase
__all__ = ['lagzero', 'lagone', 'lagx', 'lagdomain', 'lagline', 'lagadd', 'lagsub', 'lagmulx', 'lagmul', 'lagdiv', 'lagpow', 'lagval', 'lagder', 'lagint', 'lag2poly', 'poly2lag', 'lagfromroots', 'lagvander', 'lagfit', 'lagtrim', 'lagroots', 'Laguerre', 'lagval2d', 'lagval3d', 'laggrid2d', 'laggrid3d', 'lagvander2d', 'lagvander3d', 'lagcompanion', 'laggauss', 'lagweight']
lagtrim = pu.trimcoef

def poly2lag(pol):
    [pol] = pu.as_series([pol])
    res = 0
    for p in pol[::-1]:
        res = lagadd(lagmulx(res), p)
    return res

def lag2poly(c):
    from .polynomial import polyadd, polymulx, polysub
    [c] = pu.as_series([c])
    n = len(c)
    if n == 1:
        return c
    else:
        c0 = c[-2]
        c1 = c[-1]
        for i in range(n - 1, 1, -1):
            tmp = c0
            c0 = polysub(c[i - 2], c1 * (i - 1) / i)
            c1 = polyadd(tmp, polysub((2 * i - 1) * c1, polymulx(c1)) / i)
        return polyadd(c0, polysub(c1, polymulx(c1)))
lagdomain = np.array([0.0, 1.0])
lagzero = np.array([0])
lagone = np.array([1])
lagx = np.array([1, -1])

def lagline(off, scl):
    if scl != 0:
        return np.array([off + scl, -scl])
    else:
        return np.array([off])

def lagfromroots(roots):
    return pu._fromroots(lagline, lagmul, roots)

def lagadd(c1, c2):
    return pu._add(c1, c2)

def lagsub(c1, c2):
    return pu._sub(c1, c2)

def lagmulx(c):
    [c] = pu.as_series([c])
    if len(c) == 1 and c[0] == 0:
        return c
    prd = np.empty(len(c) + 1, dtype=c.dtype)
    prd[0] = c[0]
    prd[1] = -c[0]
    for i in range(1, len(c)):
        prd[i + 1] = -c[i] * (i + 1)
        prd[i] += c[i] * (2 * i + 1)
        prd[i - 1] -= c[i] * i
    return prd

def lagmul(c1, c2):
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
            c0 = lagsub(c[-i] * xs, c1 * (nd - 1) / nd)
            c1 = lagadd(tmp, lagsub((2 * nd - 1) * c1, lagmulx(c1)) / nd)
    return lagadd(c0, lagsub(c1, lagmulx(c1)))

def lagdiv(c1, c2):
    return pu._div(lagmul, c1, c2)

def lagpow(c, pow, maxpower=16):
    return pu._pow(lagmul, c, pow, maxpower)

def lagder(c, m=1, scl=1, axis=0):
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
            for j in range(n, 1, -1):
                der[j - 1] = -c[j]
                c[j - 1] += c[j]
            der[0] = -c[1]
            c = der
    c = np.moveaxis(c, 0, iaxis)
    return c

def lagint(c, m=1, k=[], lbnd=0, scl=1, axis=0):
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
            tmp[0] = c[0]
            tmp[1] = -c[0]
            for j in range(1, n):
                tmp[j] += c[j]
                tmp[j + 1] = -c[j]
            tmp[0] += k[i] - lagval(lbnd, tmp)
            c = tmp
    c = np.moveaxis(c, 0, iaxis)
    return c

def lagval(x, c, tensor=True):
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
            c0 = c[-i] - c1 * (nd - 1) / nd
            c1 = tmp + c1 * (2 * nd - 1 - x) / nd
    return c0 + c1 * (1 - x)

def lagval2d(x, y, c):
    return pu._valnd(lagval, c, x, y)

def laggrid2d(x, y, c):
    return pu._gridnd(lagval, c, x, y)

def lagval3d(x, y, z, c):
    return pu._valnd(lagval, c, x, y, z)

def laggrid3d(x, y, z, c):
    return pu._gridnd(lagval, c, x, y, z)

def lagvander(x, deg):
    ideg = pu._as_int(deg, 'deg')
    if ideg < 0:
        raise ValueError('deg must be non-negative')
    x = np.array(x, copy=None, ndmin=1) + 0.0
    dims = (ideg + 1,) + x.shape
    dtyp = x.dtype
    v = np.empty(dims, dtype=dtyp)
    v[0] = x * 0 + 1
    if ideg > 0:
        v[1] = 1 - x
        for i in range(2, ideg + 1):
            v[i] = (v[i - 1] * (2 * i - 1 - x) - v[i - 2] * (i - 1)) / i
    return np.moveaxis(v, 0, -1)

def lagvander2d(x, y, deg):
    return pu._vander_nd_flat((lagvander, lagvander), (x, y), deg)

def lagvander3d(x, y, z, deg):
    return pu._vander_nd_flat((lagvander, lagvander, lagvander), (x, y, z), deg)

def lagfit(x, y, deg, rcond=None, full=False, w=None):
    return pu._fit(lagvander, x, y, deg, rcond, full, w)

def lagcompanion(c):
    [c] = pu.as_series([c])
    if len(c) < 2:
        raise ValueError('Series must have maximum degree of at least 1.')
    if len(c) == 2:
        return np.array([[1 + c[0] / c[1]]])
    n = len(c) - 1
    mat = np.zeros((n, n), dtype=c.dtype)
    top = mat.reshape(-1)[1::n + 1]
    mid = mat.reshape(-1)[0::n + 1]
    bot = mat.reshape(-1)[n::n + 1]
    top[...] = -np.arange(1, n)
    mid[...] = 2.0 * np.arange(n) + 1.0
    bot[...] = top
    mat[:, -1] += c[:-1] / c[-1] * n
    return mat

def lagroots(c):
    [c] = pu.as_series([c])
    if len(c) <= 1:
        return np.array([], dtype=c.dtype)
    if len(c) == 2:
        return np.array([1 + c[0] / c[1]])
    m = lagcompanion(c)[::-1, ::-1]
    r = np.linalg.eigvals(m)
    r.sort()
    return r

def laggauss(deg):
    ideg = pu._as_int(deg, 'deg')
    if ideg <= 0:
        raise ValueError('deg must be a positive integer')
    c = np.array([0] * deg + [1])
    m = lagcompanion(c)
    x = np.linalg.eigvalsh(m)
    dy = lagval(x, c)
    df = lagval(x, lagder(c))
    x -= dy / df
    fm = lagval(x, c[1:])
    fm /= np.abs(fm).max()
    df /= np.abs(df).max()
    w = 1 / (fm * df)
    w /= w.sum()
    return (x, w)

def lagweight(x):
    w = np.exp(-x)
    return w

class Laguerre(ABCPolyBase):
    _add = staticmethod(lagadd)
    _sub = staticmethod(lagsub)
    _mul = staticmethod(lagmul)
    _div = staticmethod(lagdiv)
    _pow = staticmethod(lagpow)
    _val = staticmethod(lagval)
    _int = staticmethod(lagint)
    _der = staticmethod(lagder)
    _fit = staticmethod(lagfit)
    _line = staticmethod(lagline)
    _roots = staticmethod(lagroots)
    _fromroots = staticmethod(lagfromroots)
    domain = np.array(lagdomain)
    window = np.array(lagdomain)
    basis_name = 'L'
