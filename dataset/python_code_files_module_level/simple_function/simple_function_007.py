__all__ = ['matrix', 'bmat', 'asmatrix']
import ast
import sys
import warnings
import numpy._core.numeric as N
from numpy._core.numeric import concatenate, isscalar
from numpy._utils import set_module
from numpy.linalg import matrix_power

def _convert_from_string(data):
    for char in '[]':
        data = data.replace(char, '')
    rows = data.split(';')
    newdata = []
    for count, row in enumerate(rows):
        trow = row.split(',')
        newrow = []
        for col in trow:
            temp = col.split()
            newrow.extend(map(ast.literal_eval, temp))
        if count == 0:
            Ncols = len(newrow)
        elif len(newrow) != Ncols:
            raise ValueError('Rows not the same size.')
        newdata.append(newrow)
    return newdata

@set_module('numpy')
def asmatrix(data, dtype=None):
    return matrix(data, dtype=dtype, copy=False)

@set_module('numpy')
class matrix(N.ndarray):
    __array_priority__ = 10.0

    def __new__(subtype, data, dtype=None, copy=True):
        warnings.warn('the matrix subclass is not the recommended way to represent matrices or deal with linear algebra (see https://docs.scipy.org/doc/numpy/user/numpy-for-matlab-users.html). Please adjust your code to use regular ndarray.', PendingDeprecationWarning, stacklevel=2)
        if isinstance(data, matrix):
            dtype2 = data.dtype
            if dtype is None:
                dtype = dtype2
            if dtype2 == dtype and (not copy):
                return data
            return data.astype(dtype)
        if isinstance(data, N.ndarray):
            if dtype is None:
                intype = data.dtype
            else:
                intype = N.dtype(dtype)
            new = data.view(subtype)
            if intype != data.dtype:
                return new.astype(intype)
            if copy:
                return new.copy()
            else:
                return new
        if isinstance(data, str):
            data = _convert_from_string(data)
        copy = None if not copy else True
        arr = N.array(data, dtype=dtype, copy=copy)
        ndim = arr.ndim
        shape = arr.shape
        if ndim > 2:
            raise ValueError('matrix must be 2-dimensional')
        elif ndim == 0:
            shape = (1, 1)
        elif ndim == 1:
            shape = (1, shape[0])
        order = 'C'
        if ndim == 2 and arr.flags.fortran:
            order = 'F'
        if not (order or arr.flags.contiguous):
            arr = arr.copy()
        ret = N.ndarray.__new__(subtype, shape, arr.dtype, buffer=arr, order=order)
        return ret

    def __array_finalize__(self, obj):
        self._getitem = False
        if isinstance(obj, matrix) and obj._getitem:
            return
        ndim = self.ndim
        if ndim == 2:
            return
        if ndim > 2:
            newshape = tuple((x for x in self.shape if x > 1))
            ndim = len(newshape)
            if ndim == 2:
                self.shape = newshape
                return
            elif ndim > 2:
                raise ValueError('shape too large to be a matrix.')
        else:
            newshape = self.shape
        if ndim == 0:
            self.shape = (1, 1)
        elif ndim == 1:
            self.shape = (1, newshape[0])
        return

    def __getitem__(self, index):
        self._getitem = True
        try:
            out = N.ndarray.__getitem__(self, index)
        finally:
            self._getitem = False
        if not isinstance(out, N.ndarray):
            return out
        if out.ndim == 0:
            return out[()]
        if out.ndim == 1:
            sh = out.shape[0]
            try:
                n = len(index)
            except Exception:
                n = 0
            if n > 1 and isscalar(index[1]):
                out.shape = (sh, 1)
            else:
                out.shape = (1, sh)
        return out

    def __mul__(self, other):
        if isinstance(other, (N.ndarray, list, tuple)):
            return N.dot(self, asmatrix(other))
        if isscalar(other) or not hasattr(other, '__rmul__'):
            return N.dot(self, other)
        return NotImplemented

    def __rmul__(self, other):
        return N.dot(other, self)

    def __imul__(self, other):
        self[:] = self * other
        return self

    def __pow__(self, other):
        return matrix_power(self, other)

    def __ipow__(self, other):
        self[:] = self ** other
        return self

    def __rpow__(self, other):
        return NotImplemented

    def _align(self, axis):
        if axis is None:
            return self[0, 0]
        elif axis == 0:
            return self
        elif axis == 1:
            return self.transpose()
        else:
            raise ValueError('unsupported axis')

    def _collapse(self, axis):
        if axis is None:
            return self[0, 0]
        else:
            return self

    def tolist(self):
        return self.__array__().tolist()

    def sum(self, axis=None, dtype=None, out=None):
        return N.ndarray.sum(self, axis, dtype, out, keepdims=True)._collapse(axis)

    def squeeze(self, axis=None):
        return N.ndarray.squeeze(self, axis=axis)

    def flatten(self, order='C'):
        return N.ndarray.flatten(self, order=order)

    def mean(self, axis=None, dtype=None, out=None):
        return N.ndarray.mean(self, axis, dtype, out, keepdims=True)._collapse(axis)

    def std(self, axis=None, dtype=None, out=None, ddof=0):
        return N.ndarray.std(self, axis, dtype, out, ddof, keepdims=True)._collapse(axis)

    def var(self, axis=None, dtype=None, out=None, ddof=0):
        return N.ndarray.var(self, axis, dtype, out, ddof, keepdims=True)._collapse(axis)

    def prod(self, axis=None, dtype=None, out=None):
        return N.ndarray.prod(self, axis, dtype, out, keepdims=True)._collapse(axis)

    def any(self, axis=None, out=None):
        return N.ndarray.any(self, axis, out, keepdims=True)._collapse(axis)

    def all(self, axis=None, out=None):
        return N.ndarray.all(self, axis, out, keepdims=True)._collapse(axis)

    def max(self, axis=None, out=None):
        return N.ndarray.max(self, axis, out, keepdims=True)._collapse(axis)

    def argmax(self, axis=None, out=None):
        return N.ndarray.argmax(self, axis, out)._align(axis)

    def min(self, axis=None, out=None):
        return N.ndarray.min(self, axis, out, keepdims=True)._collapse(axis)

    def argmin(self, axis=None, out=None):
        return N.ndarray.argmin(self, axis, out)._align(axis)

    def ptp(self, axis=None, out=None):
        return N.ptp(self, axis, out)._align(axis)

    @property
    def I(self):
        M, N = self.shape
        if M == N:
            from numpy.linalg import inv as func
        else:
            from numpy.linalg import pinv as func
        return asmatrix(func(self))

    @property
    def A(self):
        return self.__array__()

    @property
    def A1(self):
        return self.__array__().ravel()

    def ravel(self, order='C'):
        return N.ndarray.ravel(self, order=order)

    @property
    def T(self):
        return self.transpose()

    @property
    def H(self):
        if issubclass(self.dtype.type, N.complexfloating):
            return self.transpose().conjugate()
        else:
            return self.transpose()
    getT = T.fget
    getA = A.fget
    getA1 = A1.fget
    getH = H.fget
    getI = I.fget

def _from_string(str, gdict, ldict):
    rows = str.split(';')
    rowtup = []
    for row in rows:
        trow = row.split(',')
        newrow = []
        for x in trow:
            newrow.extend(x.split())
        trow = newrow
        coltup = []
        for col in trow:
            col = col.strip()
            try:
                thismat = ldict[col]
            except KeyError:
                try:
                    thismat = gdict[col]
                except KeyError as e:
                    raise NameError(f'name {col!r} is not defined') from None
            coltup.append(thismat)
        rowtup.append(concatenate(coltup, axis=-1))
    return concatenate(rowtup, axis=0)

@set_module('numpy')
def bmat(obj, ldict=None, gdict=None):
    if isinstance(obj, str):
        if gdict is None:
            frame = sys._getframe().f_back
            glob_dict = frame.f_globals
            loc_dict = frame.f_locals
        else:
            glob_dict = gdict
            loc_dict = ldict
        return matrix(_from_string(obj, glob_dict, loc_dict))
    if isinstance(obj, (tuple, list)):
        arr_rows = []
        for row in obj:
            if isinstance(row, N.ndarray):
                return matrix(concatenate(obj, axis=-1))
            else:
                arr_rows.append(concatenate(row, axis=-1))
        return matrix(concatenate(arr_rows, axis=0))
    if isinstance(obj, N.ndarray):
        return matrix(obj)
