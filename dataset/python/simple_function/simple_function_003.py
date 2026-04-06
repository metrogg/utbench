from __future__ import annotations
from decimal import Decimal
from typing import TYPE_CHECKING, overload
import warnings
import numpy as np
from pandas._libs import lib
import pandas._libs.missing as libmissing
from pandas._libs.tslibs import NaT, iNaT
from pandas.util._decorators import set_module
from pandas.core.dtypes.common import DT64NS_DTYPE, TD64NS_DTYPE, ensure_object, is_scalar, is_string_or_object_np_dtype
from pandas.core.dtypes.dtypes import CategoricalDtype, DatetimeTZDtype, ExtensionDtype, IntervalDtype, PeriodDtype
from pandas.core.dtypes.generic import ABCDataFrame, ABCExtensionArray, ABCIndex, ABCMultiIndex, ABCSeries
from pandas.core.dtypes.inference import is_list_like
if TYPE_CHECKING:
    from re import Pattern
    from pandas._libs.missing import NAType
    from pandas._libs.tslibs import NaTType
    from pandas._typing import ArrayLike, DtypeObj, NDFrame, NDFrameT, Scalar, npt
    from pandas import Series
    from pandas.core.indexes.base import Index
isposinf_scalar = libmissing.isposinf_scalar
isneginf_scalar = libmissing.isneginf_scalar
_dtype_object = np.dtype('object')
_dtype_str = np.dtype(str)

@overload
def isna(obj: Scalar | Pattern | NAType | NaTType) -> bool:
    ...

@overload
def isna(obj: ArrayLike | Index | list) -> npt.NDArray[np.bool_]:
    ...

@overload
def isna(obj: NDFrameT) -> NDFrameT:
    ...

@overload
def isna(obj: NDFrameT | ArrayLike | Index | list) -> NDFrameT | npt.NDArray[np.bool_]:
    ...

@overload
def isna(obj: object) -> bool | npt.NDArray[np.bool_] | NDFrame:
    ...

@set_module('pandas')
def isna(obj: object) -> bool | npt.NDArray[np.bool_] | NDFrame:
    return _isna(obj)
isnull = isna

def _isna(obj):
    if is_scalar(obj):
        return libmissing.checknull(obj)
    elif isinstance(obj, ABCMultiIndex):
        raise NotImplementedError('isna is not defined for MultiIndex')
    elif isinstance(obj, type):
        return False
    elif isinstance(obj, (np.ndarray, ABCExtensionArray)):
        return _isna_array(obj)
    elif isinstance(obj, ABCIndex):
        if not obj._can_hold_na:
            return obj.isna()
        return _isna_array(obj._values)
    elif isinstance(obj, ABCSeries):
        result = _isna_array(obj._values)
        result = obj._constructor(result, index=obj.index, name=obj.name, copy=False)
        return result
    elif isinstance(obj, ABCDataFrame):
        return obj.isna()
    elif isinstance(obj, list):
        return _isna_array(np.asarray(obj, dtype=object))
    elif hasattr(obj, '__array__'):
        return _isna_array(np.asarray(obj))
    else:
        return False

def _isna_array(values: ArrayLike) -> npt.NDArray[np.bool_] | NDFrame:
    dtype = values.dtype
    result: npt.NDArray[np.bool_] | NDFrame
    if not isinstance(values, np.ndarray):
        result = values.isna()
    elif isinstance(values, np.rec.recarray):
        result = _isna_recarray_dtype(values)
    elif is_string_or_object_np_dtype(values.dtype):
        result = _isna_string_dtype(values)
    elif dtype.kind in 'mM':
        result = values.view('i8') == iNaT
    else:
        result = np.isnan(values)
    return result

def _isna_string_dtype(values: np.ndarray) -> npt.NDArray[np.bool_]:
    dtype = values.dtype
    if dtype.kind in ('S', 'U'):
        result = np.zeros(values.shape, dtype=bool)
    elif values.ndim in {1, 2}:
        result = libmissing.isnaobj(values)
    else:
        result = libmissing.isnaobj(values.ravel())
        result = result.reshape(values.shape)
    return result

def _isna_recarray_dtype(values: np.rec.recarray) -> npt.NDArray[np.bool_]:
    result = np.zeros(values.shape, dtype=bool)
    for i, record in enumerate(values):
        record_as_array = np.array(record.tolist())
        does_record_contain_nan = isna_all(record_as_array)
        result[i] = np.any(does_record_contain_nan)
    return result

@overload
def notna(obj: Scalar | Pattern | NAType | NaTType) -> bool:
    ...

@overload
def notna(obj: ArrayLike | Index | list) -> npt.NDArray[np.bool_]:
    ...

@overload
def notna(obj: NDFrameT) -> NDFrameT:
    ...

@overload
def notna(obj: NDFrameT | ArrayLike | Index | list) -> NDFrameT | npt.NDArray[np.bool_]:
    ...

@overload
def notna(obj: object) -> bool | npt.NDArray[np.bool_] | NDFrame:
    ...

@set_module('pandas')
def notna(obj: object) -> bool | npt.NDArray[np.bool_] | NDFrame:
    res = isna(obj)
    if isinstance(res, bool):
        return not res
    return ~res
notnull = notna

def array_equivalent(left, right, strict_nan: bool=False, dtype_equal: bool=False) -> bool:
    left, right = (np.asarray(left), np.asarray(right))
    if left.shape != right.shape:
        return False
    if dtype_equal:
        if left.dtype.kind in 'fc':
            return _array_equivalent_float(left, right)
        elif left.dtype.kind in 'mM':
            return _array_equivalent_datetimelike(left, right)
        elif is_string_or_object_np_dtype(left.dtype):
            return _array_equivalent_object(left, right, strict_nan)
        else:
            return np.array_equal(left, right)
    if left.dtype.kind in 'OSU' or right.dtype.kind in 'OSU':
        return _array_equivalent_object(left, right, strict_nan)
    if left.dtype.kind in 'fc':
        if not (left.size and right.size):
            return True
        return ((left == right) | isna(left) & isna(right)).all()
    elif left.dtype.kind in 'mM' or right.dtype.kind in 'mM':
        if left.dtype != right.dtype:
            return False
        left = left.view('i8')
        right = right.view('i8')
    if (left.dtype.type is np.void or right.dtype.type is np.void) and left.dtype != right.dtype:
        return False
    return np.array_equal(left, right)

def _array_equivalent_float(left: np.ndarray, right: np.ndarray) -> bool:
    return bool(((left == right) | np.isnan(left) & np.isnan(right)).all())

def _array_equivalent_datetimelike(left: np.ndarray, right: np.ndarray) -> bool:
    return np.array_equal(left.view('i8'), right.view('i8'))

def _array_equivalent_object(left: np.ndarray, right: np.ndarray, strict_nan: bool) -> bool:
    left = ensure_object(left)
    right = ensure_object(right)
    mask: npt.NDArray[np.bool_] | None = None
    if strict_nan:
        mask = isna(left) & isna(right)
        if not mask.any():
            mask = None
    try:
        if mask is None:
            return lib.array_equivalent_object(left, right)
        if not lib.array_equivalent_object(left[~mask], right[~mask]):
            return False
        left_remaining = left[mask]
        right_remaining = right[mask]
    except ValueError:
        left_remaining = left
        right_remaining = right
    for left_value, right_value in zip(left_remaining, right_remaining, strict=True):
        if left_value is NaT and right_value is not NaT:
            return False
        elif left_value is libmissing.NA and right_value is not libmissing.NA:
            return False
        elif isinstance(left_value, float) and np.isnan(left_value):
            if not isinstance(right_value, float) or not np.isnan(right_value):
                return False
        else:
            with warnings.catch_warnings():
                warnings.simplefilter('ignore', DeprecationWarning)
                try:
                    if np.any(np.asarray(left_value != right_value)):
                        return False
                except TypeError as err:
                    if 'boolean value of NA is ambiguous' in str(err):
                        return False
                    raise
                except ValueError:
                    return False
    return True

def array_equals(left: ArrayLike, right: ArrayLike) -> bool:
    if left.dtype != right.dtype:
        return False
    elif isinstance(left, ABCExtensionArray):
        return left.equals(right)
    else:
        return array_equivalent(left, right, dtype_equal=True)

def infer_fill_value(val):
    if not is_list_like(val):
        val = [val]
    val = np.asarray(val)
    if val.dtype.kind in 'mM':
        return np.array('NaT', dtype=val.dtype)
    elif val.dtype == object:
        dtype = lib.infer_dtype(ensure_object(val), skipna=False)
        if dtype in ['datetime', 'datetime64']:
            return np.array('NaT', dtype=DT64NS_DTYPE)
        elif dtype in ['timedelta', 'timedelta64']:
            return np.array('NaT', dtype=TD64NS_DTYPE)
        return np.array(np.nan, dtype=object)
    elif val.dtype.kind == 'U':
        return np.array(np.nan, dtype=val.dtype)
    return np.nan

def construct_1d_array_from_inferred_fill_value(value: object, length: int) -> ArrayLike:
    from pandas.core.algorithms import take_nd
    from pandas.core.construction import sanitize_array
    from pandas.core.indexes.base import Index
    arr = sanitize_array(value, Index(range(1)), copy=False)
    taker = -1 * np.ones(length, dtype=np.intp)
    return take_nd(arr, taker)

def maybe_fill(arr: np.ndarray) -> np.ndarray:
    if arr.dtype.kind not in 'iub':
        arr.fill(np.nan)
    return arr

def na_value_for_dtype(dtype: DtypeObj, compat: bool=True):
    if isinstance(dtype, ExtensionDtype):
        return dtype.na_value
    elif dtype.kind in 'mM':
        unit = np.datetime_data(dtype)[0]
        return dtype.type('NaT', unit)
    elif dtype.kind in 'fc':
        return np.nan
    elif dtype.kind in 'iu':
        if compat:
            return 0
        return np.nan
    elif dtype.kind == 'b':
        if compat:
            return False
        return np.nan
    return np.nan

def remove_na_arraylike(arr: Series | Index | np.ndarray):
    if isinstance(arr.dtype, ExtensionDtype):
        return arr[notna(arr)]
    else:
        return arr[notna(np.asarray(arr))]

def is_valid_na_for_dtype(obj, dtype: DtypeObj) -> bool:
    if not lib.is_scalar(obj) or not isna(obj):
        return False
    elif dtype.kind == 'M':
        if isinstance(dtype, np.dtype):
            return not isinstance(obj, (np.timedelta64, Decimal))
        return not isinstance(obj, (np.timedelta64, np.datetime64, Decimal))
    elif dtype.kind == 'm':
        return not isinstance(obj, (np.datetime64, Decimal))
    elif dtype.kind in 'iufc':
        return obj is not NaT and (not isinstance(obj, (np.datetime64, np.timedelta64)))
    elif dtype.kind == 'b':
        return lib.is_float(obj) or obj is None or obj is libmissing.NA
    elif dtype == _dtype_str:
        return not isinstance(obj, (np.datetime64, np.timedelta64, Decimal, float))
    elif dtype == _dtype_object:
        return True
    elif isinstance(dtype, PeriodDtype):
        return not isinstance(obj, (np.datetime64, np.timedelta64, Decimal))
    elif isinstance(dtype, IntervalDtype):
        return lib.is_float(obj) or obj is None or obj is libmissing.NA
    elif isinstance(dtype, CategoricalDtype):
        return is_valid_na_for_dtype(obj, dtype.categories.dtype)
    return not isinstance(obj, (np.datetime64, np.timedelta64, Decimal))

def isna_all(arr: ArrayLike) -> bool:
    total_len = len(arr)
    chunk_len = max(total_len // 40, 1000)
    dtype = arr.dtype
    if lib.is_np_dtype(dtype, 'f'):
        checker = np.isnan
    elif lib.is_np_dtype(dtype, 'mM') or isinstance(dtype, (DatetimeTZDtype, PeriodDtype)):
        checker = lambda x: np.asarray(x.view('i8')) == iNaT
    else:
        checker = _isna_array
    return all((checker(arr[i:i + chunk_len]).all() for i in range(0, total_len, chunk_len)))
