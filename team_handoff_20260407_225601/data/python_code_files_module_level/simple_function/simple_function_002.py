from __future__ import annotations
import operator
import re
from typing import TYPE_CHECKING, Self
import numpy as np
from pandas._libs import lib, missing as libmissing
from pandas.compat import HAS_PYARROW, PYARROW_MIN_VERSION, pa_version_under16p0
from pandas.util._decorators import set_module
from pandas.util._validators import validate_na_arg
from pandas.core.dtypes.common import is_scalar, pandas_dtype
from pandas.core.dtypes.missing import isna
from pandas.core.arrays._arrow_string_mixins import ArrowStringArrayMixin
from pandas.core.arrays.arrow import ArrowExtensionArray
from pandas.core.arrays.boolean import BooleanDtype
from pandas.core.arrays.floating import Float64Dtype
from pandas.core.arrays.integer import Int64Dtype
from pandas.core.arrays.numeric import NumericDtype
from pandas.core.arrays.string_ import BaseStringArray, StringDtype
from pandas.core.strings.object_array import ObjectStringArrayMixin
if HAS_PYARROW:
    import pyarrow as pa
    import pyarrow.compute as pc
if TYPE_CHECKING:
    from collections.abc import Callable, Sequence
    from pandas._typing import ArrayLike, Dtype, NpDtype, Scalar, npt
    from pandas.core.dtypes.dtypes import ExtensionDtype
    from pandas import Series

def _check_pyarrow_available() -> None:
    if not HAS_PYARROW:
        msg = f'pyarrow>={PYARROW_MIN_VERSION} is required for PyArrow backed ArrowExtensionArray.'
        raise ImportError(msg)

def _is_string_view(typ):
    return not pa_version_under16p0 and pa.types.is_string_view(typ)

@set_module('pandas.arrays')
class ArrowStringArray(ObjectStringArrayMixin, ArrowExtensionArray, BaseStringArray):
    _dtype: StringDtype

    def __init__(self, values, *, dtype: StringDtype | None=None) -> None:
        _check_pyarrow_available()
        if isinstance(values, (pa.Array, pa.ChunkedArray)) and (pa.types.is_string(values.type) or _is_string_view(values.type) or (pa.types.is_dictionary(values.type) and (pa.types.is_string(values.type.value_type) or pa.types.is_large_string(values.type.value_type) or _is_string_view(values.type.value_type)))):
            values = pc.cast(values, pa.large_string())
        super().__init__(values)
        if dtype is None:
            dtype = StringDtype(storage='pyarrow', na_value=libmissing.NA)
        self._dtype = dtype
        if not pa.types.is_large_string(self._pa_array.type):
            raise ValueError('ArrowStringArray requires a PyArrow (chunked) array of large_string type')

    def _from_pyarrow_array(self, pa_array):
        return type(self)(pa_array, dtype=self.dtype)

    @classmethod
    def _box_pa_scalar(cls, value, pa_type: pa.DataType | None=None) -> pa.Scalar:
        pa_scalar = super()._box_pa_scalar(value, pa_type)
        if pa.types.is_string(pa_scalar.type) and pa_type is None:
            pa_scalar = pc.cast(pa_scalar, pa.large_string())
        return pa_scalar

    @classmethod
    def _box_pa_array(cls, value, pa_type: pa.DataType | None=None, copy: bool=False) -> pa.Array | pa.ChunkedArray:
        pa_array = super()._box_pa_array(value, pa_type)
        if pa.types.is_string(pa_array.type) and pa_type is None:
            pa_array = pc.cast(pa_array, pa.large_string())
        return pa_array

    def __len__(self) -> int:
        return len(self._pa_array)

    @classmethod
    def _from_sequence(cls, scalars, *, dtype: Dtype | None=None, copy: bool=False) -> Self:
        from pandas.core.arrays.masked import BaseMaskedArray
        _check_pyarrow_available()
        if dtype and (not (isinstance(dtype, str) and dtype == 'string')):
            dtype = pandas_dtype(dtype)
            assert isinstance(dtype, StringDtype) and dtype.storage == 'pyarrow'
        if isinstance(scalars, BaseMaskedArray):
            na_values = scalars._mask
            result = scalars._data
            result = lib.ensure_string_array(result, copy=copy, convert_na_value=False)
            pa_arr = pa.array(result, mask=na_values, type=pa.large_string())
        elif isinstance(scalars, ArrowExtensionArray):
            pa_type = scalars._pa_array.type
            if pa.types.is_integer(pa_type) or pa.types.is_large_string(pa_type) or pa.types.is_string(pa_type) or pa.types.is_boolean(pa_type):
                pa_arr = pc.cast(scalars._pa_array, pa.large_string())
                if pa.types.is_boolean(pa_type):
                    pa_arr = pc.utf8_capitalize(pa_arr)
            else:
                result = lib.ensure_string_array(scalars, copy=copy)
                pa_arr = pa.array(result, type=pa.large_string(), from_pandas=True)
        elif isinstance(scalars, (pa.Array, pa.ChunkedArray)):
            pa_arr = pc.cast(scalars, pa.large_string())
        else:
            result = lib.ensure_string_array(scalars, copy=copy)
            pa_arr = pa.array(result, type=pa.large_string(), from_pandas=True)
        return cls(pa_arr, dtype=dtype)

    @classmethod
    def _from_sequence_of_strings(cls, strings, *, dtype: ExtensionDtype, copy: bool=False) -> Self:
        return cls._from_sequence(strings, dtype=dtype, copy=copy)

    @property
    def dtype(self) -> StringDtype:
        return self._dtype

    def insert(self, loc: int, item) -> ArrowStringArray:
        if self.dtype.na_value is np.nan and item is np.nan:
            item = libmissing.NA
        if not isinstance(item, str) and item is not libmissing.NA:
            raise TypeError(f"Invalid value '{item}' for dtype 'str'. Value should be a string or missing value, got '{type(item).__name__}' instead.")
        return super().insert(loc, item)

    def _convert_bool_result(self, values, na=lib.no_default, method_name=None):
        validate_na_arg(na, name='na')
        if self.dtype.na_value is np.nan:
            if na is lib.no_default or isna(na):
                values = values.fill_null(False)
            else:
                values = values.fill_null(na)
            return values.to_numpy()
        elif na is not lib.no_default and (not isna(na)):
            values = values.fill_null(na)
        return BooleanDtype().__from_arrow__(values)

    def _maybe_convert_setitem_value(self, value):
        if is_scalar(value):
            if isna(value):
                value = None
            elif not isinstance(value, str):
                raise TypeError(f"Invalid value '{value}' for dtype 'str'. Value should be a string or missing value, got '{type(value).__name__}' instead.")
        else:
            value = np.array(value, dtype=object, copy=True)
            value[isna(value)] = None
            for v in value:
                if not (v is None or isinstance(v, str)):
                    raise TypeError("Invalid value for dtype 'str'. Value should be a string or missing value (or array of those).")
        return super()._maybe_convert_setitem_value(value)

    def isin(self, values: ArrayLike) -> npt.NDArray[np.bool_]:
        value_set = [pa_scalar.as_py() for pa_scalar in [pa.scalar(value, from_pandas=True) for value in values] if pa_scalar.type in (pa.string(), pa.null(), pa.large_string())]
        if not value_set:
            return np.zeros(len(self), dtype=bool)
        result = pc.is_in(self._pa_array, value_set=pa.array(value_set, type=self._pa_array.type))
        return np.array(result, dtype=np.bool_)

    def astype(self, dtype, copy: bool=True):
        dtype = pandas_dtype(dtype)
        if dtype == self.dtype:
            if copy:
                return self.copy()
            return self
        elif isinstance(dtype, NumericDtype):
            data = self._pa_array.cast(pa.from_numpy_dtype(dtype.numpy_dtype))
            return dtype.__from_arrow__(data)
        elif isinstance(dtype, np.dtype) and np.issubdtype(dtype, np.floating):
            return self.to_numpy(dtype=dtype, na_value=np.nan)
        return super().astype(dtype, copy=copy)
    _str_isalnum = ArrowStringArrayMixin._str_isalnum
    _str_isalpha = ArrowStringArrayMixin._str_isalpha
    _str_isdecimal = ArrowStringArrayMixin._str_isdecimal
    _str_isdigit = ArrowStringArrayMixin._str_isdigit
    _str_islower = ArrowStringArrayMixin._str_islower
    _str_isnumeric = ArrowStringArrayMixin._str_isnumeric
    _str_isspace = ArrowStringArrayMixin._str_isspace
    _str_istitle = ArrowStringArrayMixin._str_istitle
    _str_isupper = ArrowStringArrayMixin._str_isupper
    _str_map = BaseStringArray._str_map
    _str_startswith = ArrowStringArrayMixin._str_startswith
    _str_endswith = ArrowStringArrayMixin._str_endswith
    _str_pad = ArrowStringArrayMixin._str_pad
    _str_lower = ArrowStringArrayMixin._str_lower
    _str_upper = ArrowStringArrayMixin._str_upper
    _str_strip = ArrowStringArrayMixin._str_strip
    _str_lstrip = ArrowStringArrayMixin._str_lstrip
    _str_rstrip = ArrowStringArrayMixin._str_rstrip
    _str_removesuffix = ArrowStringArrayMixin._str_removesuffix
    _str_removeprefix = ArrowStringArrayMixin._str_removeprefix
    _str_find = ArrowStringArrayMixin._str_find
    _str_get = ArrowStringArrayMixin._str_get
    _str_capitalize = ArrowStringArrayMixin._str_capitalize
    _str_title = ArrowStringArrayMixin._str_title
    _str_swapcase = ArrowStringArrayMixin._str_swapcase
    _str_slice_replace = ArrowStringArrayMixin._str_slice_replace
    _str_len = ArrowStringArrayMixin._str_len
    _str_slice = ArrowStringArrayMixin._str_slice

    @staticmethod
    def _is_re_pattern_with_flags(pat: str | re.Pattern) -> bool:
        return isinstance(pat, re.Pattern) and pat.flags & ~(re.IGNORECASE | re.UNICODE) != 0

    @staticmethod
    def _preprocess_re_pattern(pat: str | re.Pattern, case: bool, flags: int) -> tuple[str, bool, int]:
        if isinstance(pat, re.Pattern):
            pattern = pat.pattern
            flags = pat.flags
            if flags & re.IGNORECASE:
                case = False
                flags = flags & ~re.IGNORECASE
            flags = flags & ~re.UNICODE
        else:
            pattern = pat
        if pattern.endswith('\\Z') and (len(pattern) - len(pattern[:-1].rstrip('\\')) + 1) % 2 == 1:
            pattern = pattern[:-2] + '\\z'
        return (pattern, case, flags)

    def _str_contains(self, pat, case: bool=True, flags: int=0, na=lib.no_default, regex: bool=True):
        if flags or self._is_re_pattern_with_flags(pat) or (regex and self._has_unsupported_regex(pat)):
            return super()._str_contains(pat, case, flags, na, regex)
        pat, case, flags = self._preprocess_re_pattern(pat, case, flags)
        return ArrowStringArrayMixin._str_contains(self, pat, case, flags, na, regex)

    def _str_match(self, pat: str | re.Pattern, case: bool=True, flags: int=0, na: Scalar | lib.NoDefault=lib.no_default):
        if flags or self._is_re_pattern_with_flags(pat) or self._has_unsupported_regex(pat):
            return super()._str_match(pat, case, flags, na)
        pat, case, flags = self._preprocess_re_pattern(pat, case, flags)
        return ArrowStringArrayMixin._str_match(self, pat, case, flags, na)

    def _str_fullmatch(self, pat: str | re.Pattern, case: bool=True, flags: int=0, na: Scalar | lib.NoDefault=lib.no_default):
        if flags or self._is_re_pattern_with_flags(pat) or self._has_unsupported_regex(pat):
            return super()._str_fullmatch(pat, case, flags, na)
        pat, case, flags = self._preprocess_re_pattern(pat, case, flags)
        return ArrowStringArrayMixin._str_fullmatch(self, pat, case, flags, na)

    def _str_replace(self, pat: str | re.Pattern, repl: str | Callable, n: int=-1, case: bool=True, flags: int=0, regex: bool=True):
        if isinstance(pat, re.Pattern) or callable(repl) or (not case) or flags or (isinstance(repl, str) and '\\g<' in repl) or (regex and self._has_unsupported_regex(pat)):
            return super()._str_replace(pat, repl, n, case, flags, regex)
        if regex:
            pat, case, flags = self._preprocess_re_pattern(pat, case, flags)
        return ArrowStringArrayMixin._str_replace(self, pat, repl, n, case, flags, regex)

    def _str_repeat(self, repeats: int | Sequence[int]):
        if not isinstance(repeats, int):
            return super()._str_repeat(repeats)
        else:
            return ArrowExtensionArray._str_repeat(self, repeats=repeats)

    def _str_count(self, pat: str, flags: int=0):
        if flags or self._has_unsupported_regex(pat):
            return super()._str_count(pat, flags)
        pat, _, _ = self._preprocess_re_pattern(pat, True, 0)
        result = pc.count_substring_regex(self._pa_array, pat)
        return self._convert_int_result(result)

    def _str_get_dummies(self, sep: str='|', dtype: NpDtype | None=None):
        if dtype is None:
            dtype = np.int64
        dummies_pa, labels = ArrowExtensionArray(self._pa_array)._str_get_dummies(sep, dtype)
        if len(labels) == 0:
            return (np.empty(shape=(0, 0), dtype=dtype), labels)
        dummies = np.vstack(dummies_pa.to_numpy())
        _dtype = pandas_dtype(dtype)
        dummies_dtype: NpDtype
        if isinstance(_dtype, np.dtype):
            dummies_dtype = _dtype
        else:
            dummies_dtype = np.bool_
        return (dummies.astype(dummies_dtype, copy=False), labels)

    def _convert_int_result(self, result):
        if self.dtype.na_value is np.nan:
            if isinstance(result, pa.Array):
                result = result.to_numpy(zero_copy_only=False)
            else:
                result = result.to_numpy()
            if result.dtype == np.int32:
                result = result.astype(np.int64)
            return result
        return Int64Dtype().__from_arrow__(result)

    def _convert_rank_result(self, result):
        if self.dtype.na_value is np.nan:
            if isinstance(result, pa.Array):
                result = result.to_numpy(zero_copy_only=False)
            else:
                result = result.to_numpy()
            return result.astype('float64', copy=False)
        return Float64Dtype().__from_arrow__(result)

    def _reduce(self, name: str, *, skipna: bool=True, keepdims: bool=False, **kwargs):
        if self.dtype.na_value is np.nan and name in ['any', 'all']:
            if not skipna:
                nas = pc.is_null(self._pa_array)
                arr = pc.or_kleene(nas, pc.not_equal(self._pa_array, ''))
            else:
                arr = pc.not_equal(self._pa_array, '')
            result = ArrowExtensionArray(arr)._reduce(name, skipna=skipna, keepdims=keepdims, **kwargs)
            if keepdims:
                return result.astype(np.bool_)
            return result
        if name in ('min', 'max', 'sum', 'argmin', 'argmax'):
            result = self._reduce_calc(name, skipna=skipna, keepdims=keepdims, **kwargs)
        else:
            raise TypeError(f"Cannot perform reduction '{name}' with string dtype")
        if name in ('argmin', 'argmax') and isinstance(result, pa.Array):
            return self._convert_int_result(result)
        elif isinstance(result, pa.Array):
            return type(self)(result, dtype=self.dtype)
        else:
            return result

    def value_counts(self, dropna: bool=True) -> Series:
        result = super().value_counts(dropna=dropna)
        if self.dtype.na_value is np.nan:
            res_values = result._values.to_numpy()
            return result._constructor(res_values, index=result.index, name=result.name, copy=False)
        return result

    def _cmp_method(self, other, op):
        if isinstance(other, (BaseStringArray, ArrowExtensionArray)) and self.dtype.na_value is not libmissing.NA and (other.dtype.na_value is libmissing.NA):
            return NotImplemented
        result = super()._cmp_method(other, op)
        if self.dtype.na_value is np.nan:
            if op == operator.ne:
                return result.to_numpy(np.bool_, na_value=True)
            else:
                return result.to_numpy(np.bool_, na_value=False)
        return result

    def __pos__(self) -> Self:
        raise TypeError(f"bad operand type for unary +: '{self.dtype}'")
