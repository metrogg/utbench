from __future__ import annotations
import datetime
from functools import partial
import operator
from typing import TYPE_CHECKING, Any
import numpy as np
from pandas._libs import NaT, Timedelta, Timestamp, lib, ops as libops
from pandas._libs.tslibs import BaseOffset, get_supported_dtype, is_supported_dtype, is_unitless
from pandas.core.dtypes.cast import construct_1d_object_array_from_listlike, find_common_type
from pandas.core.dtypes.common import ensure_object, is_bool_dtype, is_list_like, is_numeric_v_string_like, is_object_dtype, is_scalar
from pandas.core.dtypes.generic import ABCExtensionArray, ABCIndex, ABCSeries
from pandas.core.dtypes.missing import isna, notna
from pandas.core import roperator
from pandas.core.computation import expressions
from pandas.core.construction import ensure_wrapped_if_datetimelike, sanitize_array
from pandas.core.ops import missing
from pandas.core.ops.dispatch import should_extension_dispatch
from pandas.core.ops.invalid import invalid_comparison
if TYPE_CHECKING:
    from pandas._typing import ArrayLike, Shape

def fill_binop(left, right, fill_value):
    if fill_value is not None:
        left_mask = isna(left)
        right_mask = isna(right)
        mask = left_mask ^ right_mask
        if left_mask.any():
            left = left.copy()
            left[left_mask & mask] = fill_value
        if right_mask.any():
            right = right.copy()
            right[right_mask & mask] = fill_value
    return (left, right)

def comp_method_OBJECT_ARRAY(op, x, y):
    if isinstance(y, list):
        y = construct_1d_object_array_from_listlike(y)
    if isinstance(y, (np.ndarray, ABCSeries, ABCIndex)):
        if not is_object_dtype(y.dtype):
            y = y.astype(np.object_)
        if isinstance(y, (ABCSeries, ABCIndex)):
            y = y._values
        if x.shape != y.shape:
            raise ValueError('Shapes must match', x.shape, y.shape)
        result = libops.vec_compare(x.ravel(), y.ravel(), op)
    else:
        result = libops.scalar_compare(x.ravel(), y, op)
    return result.reshape(x.shape)

def _masked_arith_op(x: np.ndarray, y, op) -> np.ndarray:
    xrav = x.ravel()
    if isinstance(y, np.ndarray):
        dtype = find_common_type([x.dtype, y.dtype])
        result = np.empty(x.size, dtype=dtype)
        if len(x) != len(y):
            raise ValueError(x.shape, y.shape)
        ymask = notna(y)
        yrav = y.ravel()
        mask = notna(xrav) & ymask.ravel()
        if mask.any():
            result[mask] = op(xrav[mask], yrav[mask])
    else:
        if not is_scalar(y):
            raise TypeError(f'Cannot broadcast np.ndarray with operand of type {type(y)}')
        result = np.empty(x.size, dtype=x.dtype)
        mask = notna(xrav)
        if op is pow:
            mask = np.where(x == 1, False, mask)
        elif op is roperator.rpow:
            mask = np.where(y == 1, False, mask)
        if mask.any():
            result[mask] = op(xrav[mask], y)
    np.putmask(result, ~mask, np.nan)
    result = result.reshape(x.shape)
    return result

def _na_arithmetic_op(left: np.ndarray, right, op, is_cmp: bool=False):
    if isinstance(right, str):
        func = op
    else:
        func = partial(expressions.evaluate, op)
    try:
        result = func(left, right)
    except TypeError:
        if not is_cmp and (left.dtype == object or getattr(right, 'dtype', None) == object):
            result = _masked_arith_op(left, right, op)
        else:
            raise
    if is_cmp and (is_scalar(result) or result is NotImplemented):
        return invalid_comparison(left, right, op)
    return missing.dispatch_fill_zeros(op, left, right, result)

def arithmetic_op(left: ArrayLike, right: Any, op):
    if isinstance(right, list):
        right = sanitize_array(right, None)
    right = ensure_wrapped_if_datetimelike(right)
    if should_extension_dispatch(left, right) or isinstance(right, (Timedelta, BaseOffset, Timestamp)) or right is NaT:
        res_values = op(left, right)
    else:
        _bool_arith_check(op, left, right)
        res_values = _na_arithmetic_op(left, right, op)
    return res_values

def comparison_op(left: ArrayLike, right: Any, op) -> ArrayLike:
    lvalues = ensure_wrapped_if_datetimelike(left)
    rvalues = ensure_wrapped_if_datetimelike(right)
    rvalues = lib.item_from_zerodim(rvalues)
    rvalues_is_zerodim: bool = getattr(rvalues, 'ndim', None) == 0
    if isinstance(rvalues, list):
        rvalues = sanitize_array(rvalues, None)
    rvalues = ensure_wrapped_if_datetimelike(rvalues)
    if isinstance(rvalues, (np.ndarray, ABCExtensionArray)) and (not rvalues_is_zerodim):
        if len(lvalues) != len(rvalues):
            raise ValueError('Lengths must match to compare', lvalues.shape, rvalues.shape)
    if should_extension_dispatch(lvalues, rvalues) or ((isinstance(rvalues, (Timedelta, BaseOffset, Timestamp)) or right is NaT) and lvalues.dtype != object):
        res_values = op(lvalues, rvalues)
    elif (is_scalar(rvalues) or rvalues_is_zerodim) and isna(rvalues):
        if op is operator.ne:
            res_values = np.ones(lvalues.shape, dtype=bool)
        else:
            res_values = np.zeros(lvalues.shape, dtype=bool)
    elif is_numeric_v_string_like(lvalues, rvalues):
        return invalid_comparison(lvalues, rvalues, op)
    elif lvalues.dtype == object or isinstance(rvalues, str):
        res_values = comp_method_OBJECT_ARRAY(op, lvalues, rvalues)
    else:
        res_values = _na_arithmetic_op(lvalues, rvalues, op, is_cmp=True)
    return res_values

def na_logical_op(x: np.ndarray, y, op):
    try:
        result = op(x, y)
    except TypeError:
        if isinstance(y, np.ndarray):
            assert not (x.dtype.kind == 'b' and y.dtype.kind == 'b')
            x = ensure_object(x)
            y = ensure_object(y)
            result = libops.vec_binop(x.ravel(), y.ravel(), op)
        else:
            assert lib.is_scalar(y)
            if not isna(y):
                y = bool(y)
            try:
                result = libops.scalar_binop(x, y, op)
            except (TypeError, ValueError, AttributeError, OverflowError, NotImplementedError) as err:
                typ = type(y).__name__
                raise TypeError(f"Cannot perform '{op.__name__}' with a dtyped [{x.dtype}] array and scalar of type [{typ}]") from err
    return result.reshape(x.shape)

def logical_op(left: ArrayLike, right: Any, op) -> ArrayLike:

    def fill_bool(x, left=None):
        if x.dtype.kind in 'cfO':
            mask = isna(x)
            if mask.any():
                x = x.astype(object)
                x[mask] = False
        if left is None or left.dtype.kind == 'b':
            x = x.astype(bool)
        return x
    right = lib.item_from_zerodim(right)
    if is_list_like(right) and (not hasattr(right, 'dtype')):
        raise TypeError('Logical ops (and, or, xor) between Pandas objects and dtype-less sequences (e.g. list, tuple) are no longer supported. Wrap the object in a Series, Index, or np.array before operating instead.')
    lvalues = ensure_wrapped_if_datetimelike(left)
    rvalues = right
    if should_extension_dispatch(lvalues, rvalues):
        res_values = op(lvalues, rvalues)
    else:
        if isinstance(rvalues, np.ndarray):
            is_other_int_dtype = rvalues.dtype.kind in 'iu'
            if not is_other_int_dtype:
                rvalues = fill_bool(rvalues, lvalues)
        else:
            is_other_int_dtype = lib.is_integer(rvalues)
        res_values = na_logical_op(lvalues, rvalues, op)
        if not (left.dtype.kind in 'iu' and is_other_int_dtype):
            res_values = fill_bool(res_values)
    return res_values

def get_array_op(op):
    if isinstance(op, partial):
        return op
    op_name = op.__name__.strip('_').lstrip('r')
    if op_name == 'arith_op':
        return op
    if op_name in {'eq', 'ne', 'lt', 'le', 'gt', 'ge'}:
        return partial(comparison_op, op=op)
    elif op_name in {'and', 'or', 'xor', 'rand', 'ror', 'rxor'}:
        return partial(logical_op, op=op)
    elif op_name in {'add', 'sub', 'mul', 'truediv', 'floordiv', 'mod', 'divmod', 'pow'}:
        return partial(arithmetic_op, op=op)
    else:
        raise NotImplementedError(op_name)

def maybe_prepare_scalar_for_op(obj, shape: Shape):
    if type(obj) is datetime.timedelta:
        return Timedelta(obj)
    elif type(obj) is datetime.datetime:
        return Timestamp(obj)
    elif isinstance(obj, np.datetime64):
        if isna(obj):
            from pandas.core.arrays import DatetimeArray
            if is_unitless(obj.dtype):
                obj = obj.astype('datetime64[s]')
            elif not is_supported_dtype(obj.dtype):
                new_dtype = get_supported_dtype(obj.dtype)
                obj = obj.astype(new_dtype)
            right = np.broadcast_to(obj, shape)
            return DatetimeArray._simple_new(right, dtype=right.dtype)
        return Timestamp(obj)
    elif isinstance(obj, np.timedelta64):
        if isna(obj):
            from pandas.core.arrays import TimedeltaArray
            if is_unitless(obj.dtype):
                obj = obj.astype('timedelta64[s]')
            elif not is_supported_dtype(obj.dtype):
                new_dtype = get_supported_dtype(obj.dtype)
                obj = obj.astype(new_dtype)
            right = np.broadcast_to(obj, shape)
            return TimedeltaArray._simple_new(right, dtype=right.dtype)
        return Timedelta(obj)
    elif isinstance(obj, np.integer):
        return int(obj)
    elif isinstance(obj, np.floating):
        return float(obj)
    return obj
_BOOL_OP_NOT_ALLOWED = {operator.truediv, roperator.rtruediv, operator.floordiv, roperator.rfloordiv, operator.pow, roperator.rpow, divmod, roperator.rdivmod}

def _bool_arith_check(op, a: np.ndarray, b) -> None:
    if op in _BOOL_OP_NOT_ALLOWED:
        if a.dtype.kind == 'b' and (is_bool_dtype(b) or lib.is_bool(b)):
            op_name = op.__name__.strip('_').lstrip('r')
            raise NotImplementedError(f"operator '{op_name}' not implemented for bool dtypes")
