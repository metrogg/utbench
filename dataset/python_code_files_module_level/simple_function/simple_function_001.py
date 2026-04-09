import functools
import numpy as np
from numpy._core import overrides
from numpy._core.multiarray import compare_chararrays
from numpy._core.strings import _join as join, _rsplit as rsplit, _split as split, _splitlines as splitlines
from numpy._utils import set_module
from numpy.strings import *
from numpy.strings import multiply as strings_multiply, partition as strings_partition, rpartition as strings_rpartition
from .numeric import array as narray, asarray as asnarray, ndarray
from .numerictypes import bytes_, character, str_
__all__ = ['equal', 'not_equal', 'greater_equal', 'less_equal', 'greater', 'less', 'str_len', 'add', 'multiply', 'mod', 'capitalize', 'center', 'count', 'decode', 'encode', 'endswith', 'expandtabs', 'find', 'index', 'isalnum', 'isalpha', 'isdigit', 'islower', 'isspace', 'istitle', 'isupper', 'join', 'ljust', 'lower', 'lstrip', 'partition', 'replace', 'rfind', 'rindex', 'rjust', 'rpartition', 'rsplit', 'rstrip', 'split', 'splitlines', 'startswith', 'strip', 'swapcase', 'title', 'translate', 'upper', 'zfill', 'isnumeric', 'isdecimal', 'array', 'asarray', 'compare_chararrays', 'chararray']
array_function_dispatch = functools.partial(overrides.array_function_dispatch, module='numpy.char')

def _binary_op_dispatcher(x1, x2):
    return (x1, x2)

@array_function_dispatch(_binary_op_dispatcher)
def equal(x1, x2):
    return compare_chararrays(x1, x2, '==', True)

@array_function_dispatch(_binary_op_dispatcher)
def not_equal(x1, x2):
    return compare_chararrays(x1, x2, '!=', True)

@array_function_dispatch(_binary_op_dispatcher)
def greater_equal(x1, x2):
    return compare_chararrays(x1, x2, '>=', True)

@array_function_dispatch(_binary_op_dispatcher)
def less_equal(x1, x2):
    return compare_chararrays(x1, x2, '<=', True)

@array_function_dispatch(_binary_op_dispatcher)
def greater(x1, x2):
    return compare_chararrays(x1, x2, '>', True)

@array_function_dispatch(_binary_op_dispatcher)
def less(x1, x2):
    return compare_chararrays(x1, x2, '<', True)

@set_module('numpy.char')
def multiply(a, i):
    try:
        return strings_multiply(a, i)
    except TypeError:
        raise ValueError('Can only multiply by integers')

@set_module('numpy.char')
def partition(a, sep):
    return np.stack(strings_partition(a, sep), axis=-1)

@set_module('numpy.char')
def rpartition(a, sep):
    return np.stack(strings_rpartition(a, sep), axis=-1)

@set_module('numpy.char')
class chararray(ndarray):

    def __new__(subtype, shape, itemsize=1, unicode=False, buffer=None, offset=0, strides=None, order='C'):
        if unicode:
            dtype = str_
        else:
            dtype = bytes_
        itemsize = int(itemsize)
        if isinstance(buffer, str):
            filler = buffer
            buffer = None
        else:
            filler = None
        if buffer is None:
            self = ndarray.__new__(subtype, shape, (dtype, itemsize), order=order)
        else:
            self = ndarray.__new__(subtype, shape, (dtype, itemsize), buffer=buffer, offset=offset, strides=strides, order=order)
        if filler is not None:
            self[...] = filler
        return self

    def __array_wrap__(self, arr, context=None, return_scalar=False):
        if arr.dtype.char in 'SUbc':
            return arr.view(type(self))
        return arr

    def __array_finalize__(self, obj):
        if self.dtype.char not in 'VSUbc':
            raise ValueError('Can only create a chararray from string data.')

    def __getitem__(self, obj):
        val = ndarray.__getitem__(self, obj)
        if isinstance(val, character):
            return val.rstrip()
        return val

    def __eq__(self, other):
        return equal(self, other)

    def __ne__(self, other):
        return not_equal(self, other)

    def __ge__(self, other):
        return greater_equal(self, other)

    def __le__(self, other):
        return less_equal(self, other)

    def __gt__(self, other):
        return greater(self, other)

    def __lt__(self, other):
        return less(self, other)

    def __add__(self, other):
        return add(self, other)

    def __radd__(self, other):
        return add(other, self)

    def __mul__(self, i):
        return asarray(multiply(self, i))

    def __rmul__(self, i):
        return asarray(multiply(self, i))

    def __mod__(self, i):
        return asarray(mod(self, i))

    def __rmod__(self, other):
        return NotImplemented

    def argsort(self, axis=-1, kind=None, order=None, *, stable=None):
        return self.__array__().argsort(axis, kind, order, stable=stable)
    argsort.__doc__ = ndarray.argsort.__doc__

    def capitalize(self):
        return asarray(capitalize(self))

    def center(self, width, fillchar=' '):
        return asarray(center(self, width, fillchar))

    def count(self, sub, start=0, end=None):
        return count(self, sub, start, end)

    def decode(self, encoding=None, errors=None):
        return decode(self, encoding, errors)

    def encode(self, encoding=None, errors=None):
        return encode(self, encoding, errors)

    def endswith(self, suffix, start=0, end=None):
        return endswith(self, suffix, start, end)

    def expandtabs(self, tabsize=8):
        return asarray(expandtabs(self, tabsize))

    def find(self, sub, start=0, end=None):
        return find(self, sub, start, end)

    def index(self, sub, start=0, end=None):
        return index(self, sub, start, end)

    def isalnum(self):
        return isalnum(self)

    def isalpha(self):
        return isalpha(self)

    def isdigit(self):
        return isdigit(self)

    def islower(self):
        return islower(self)

    def isspace(self):
        return isspace(self)

    def istitle(self):
        return istitle(self)

    def isupper(self):
        return isupper(self)

    def join(self, seq):
        return join(self, seq)

    def ljust(self, width, fillchar=' '):
        return asarray(ljust(self, width, fillchar))

    def lower(self):
        return asarray(lower(self))

    def lstrip(self, chars=None):
        return lstrip(self, chars)

    def partition(self, sep):
        return asarray(partition(self, sep))

    def replace(self, old, new, count=None):
        return replace(self, old, new, count if count is not None else -1)

    def rfind(self, sub, start=0, end=None):
        return rfind(self, sub, start, end)

    def rindex(self, sub, start=0, end=None):
        return rindex(self, sub, start, end)

    def rjust(self, width, fillchar=' '):
        return asarray(rjust(self, width, fillchar))

    def rpartition(self, sep):
        return asarray(rpartition(self, sep))

    def rsplit(self, sep=None, maxsplit=None):
        return rsplit(self, sep, maxsplit)

    def rstrip(self, chars=None):
        return rstrip(self, chars)

    def split(self, sep=None, maxsplit=None):
        return split(self, sep, maxsplit)

    def splitlines(self, keepends=None):
        return splitlines(self, keepends)

    def startswith(self, prefix, start=0, end=None):
        return startswith(self, prefix, start, end)

    def strip(self, chars=None):
        return strip(self, chars)

    def swapcase(self):
        return asarray(swapcase(self))

    def title(self):
        return asarray(title(self))

    def translate(self, table, deletechars=None):
        return asarray(translate(self, table, deletechars))

    def upper(self):
        return asarray(upper(self))

    def zfill(self, width):
        return asarray(zfill(self, width))

    def isnumeric(self):
        return isnumeric(self)

    def isdecimal(self):
        return isdecimal(self)

@set_module('numpy.char')
def array(obj, itemsize=None, copy=True, unicode=None, order=None):
    if isinstance(obj, (bytes, str)):
        if unicode is None:
            if isinstance(obj, str):
                unicode = True
            else:
                unicode = False
        if itemsize is None:
            itemsize = len(obj)
        shape = len(obj) // itemsize
        return chararray(shape, itemsize=itemsize, unicode=unicode, buffer=obj, order=order)
    if isinstance(obj, (list, tuple)):
        obj = asnarray(obj)
    if isinstance(obj, ndarray) and issubclass(obj.dtype.type, character):
        if not isinstance(obj, chararray):
            obj = obj.view(chararray)
        if itemsize is None:
            itemsize = obj.itemsize
            if issubclass(obj.dtype.type, str_):
                itemsize //= 4
        if unicode is None:
            if issubclass(obj.dtype.type, str_):
                unicode = True
            else:
                unicode = False
        if unicode:
            dtype = str_
        else:
            dtype = bytes_
        if order is not None:
            obj = asnarray(obj, order=order)
        if copy or itemsize != obj.itemsize or (not unicode and isinstance(obj, str_)) or (unicode and isinstance(obj, bytes_)):
            obj = obj.astype((dtype, int(itemsize)))
        return obj
    if isinstance(obj, ndarray) and issubclass(obj.dtype.type, object):
        if itemsize is None:
            obj = obj.tolist()
    if unicode:
        dtype = str_
    else:
        dtype = bytes_
    if itemsize is None:
        val = narray(obj, dtype=dtype, order=order, subok=True)
    else:
        val = narray(obj, dtype=(dtype, itemsize), order=order, subok=True)
    return val.view(chararray)

@set_module('numpy.char')
def asarray(obj, itemsize=None, unicode=None, order=None):
    return array(obj, itemsize, copy=False, unicode=unicode, order=order)
