from __future__ import annotations
from collections.abc import Iterable, Sequence
from typing import TypeVar, overload
import numpy as np
from pandas._libs import lib
from pandas._libs.missing import NA
from pandas.core.dtypes.common import is_bool, is_integer
BoolishT = TypeVar('BoolishT', bool, int)
BoolishNoneT = TypeVar('BoolishNoneT', bool, int, None)

def _check_arg_length(fname, args, max_fname_arg_count, compat_args) -> None:
    if max_fname_arg_count < 0:
        raise ValueError("'max_fname_arg_count' must be non-negative")
    if len(args) > len(compat_args):
        max_arg_count = len(compat_args) + max_fname_arg_count
        actual_arg_count = len(args) + max_fname_arg_count
        argument = 'argument' if max_arg_count == 1 else 'arguments'
        raise TypeError(f'{fname}() takes at most {max_arg_count} {argument} ({actual_arg_count} given)')

def _check_for_default_values(fname, arg_val_dict, compat_args) -> None:
    for key in arg_val_dict:
        try:
            v1 = arg_val_dict[key]
            v2 = compat_args[key]
            if v1 is not None and v2 is None or (v1 is None and v2 is not None):
                match = False
            else:
                match = v1 == v2
            if not is_bool(match):
                raise ValueError("'match' is not a boolean")
        except ValueError:
            match = arg_val_dict[key] is compat_args[key]
        if not match:
            raise ValueError(f"the '{key}' parameter is not supported in the pandas implementation of {fname}()")

def validate_args(fname, args, max_fname_arg_count, compat_args) -> None:
    _check_arg_length(fname, args, max_fname_arg_count, compat_args)
    kwargs = dict(zip(compat_args, args, strict=False))
    _check_for_default_values(fname, kwargs, compat_args)

def _check_for_invalid_keys(fname, kwargs, compat_args) -> None:
    diff = set(kwargs) - set(compat_args)
    if diff:
        bad_arg = next(iter(diff))
        raise TypeError(f"{fname}() got an unexpected keyword argument '{bad_arg}'")

def validate_kwargs(fname, kwargs, compat_args) -> None:
    kwds = kwargs.copy()
    _check_for_invalid_keys(fname, kwargs, compat_args)
    _check_for_default_values(fname, kwds, compat_args)

def validate_args_and_kwargs(fname, args, kwargs, max_fname_arg_count, compat_args) -> None:
    _check_arg_length(fname, args + tuple(kwargs.values()), max_fname_arg_count, compat_args)
    args_dict = dict(zip(compat_args, args, strict=False))
    for key in args_dict:
        if key in kwargs:
            raise TypeError(f"{fname}() got multiple values for keyword argument '{key}'")
    kwargs.update(args_dict)
    validate_kwargs(fname, kwargs, compat_args)

def validate_bool_kwarg(value: BoolishNoneT, arg_name: str, none_allowed: bool=True, int_allowed: bool=False) -> BoolishNoneT:
    good_value = is_bool(value)
    if none_allowed:
        good_value = good_value or value is None
    if int_allowed:
        good_value = good_value or isinstance(value, int)
    if not good_value:
        raise ValueError(f'For argument "{arg_name}" expected type bool, received type {type(value).__name__}.')
    return value

def validate_na_arg(value, name: str):
    if value is lib.no_default or isinstance(value, bool) or value is None or (value is NA) or (lib.is_float(value) and np.isnan(value)):
        return
    raise ValueError(f'{name} must be None, pd.NA, np.nan, True, or False; got {value}')

def validate_fillna_kwargs(value, method, validate_scalar_dict_value: bool=True):
    from pandas.core.missing import clean_fill_method
    if value is None and method is None:
        raise ValueError("Must specify a fill 'value' or 'method'.")
    if value is None and method is not None:
        method = clean_fill_method(method)
    elif value is not None and method is None:
        if validate_scalar_dict_value and isinstance(value, (list, tuple)):
            raise TypeError(f'"value" parameter must be a scalar or dict, but you passed a "{type(value).__name__}"')
    elif value is not None and method is not None:
        raise ValueError("Cannot specify both 'value' and 'method'.")
    return (value, method)

def validate_percentile(q: float | Iterable[float]) -> np.ndarray:
    q_arr = np.asarray(q)
    msg = 'percentiles should all be in the interval [0, 1]'
    if q_arr.ndim == 0:
        if not 0 <= q_arr <= 1:
            raise ValueError(msg)
    elif not all((0 <= qs <= 1 for qs in q_arr)):
        raise ValueError(msg)
    return q_arr

@overload
def validate_ascending(ascending: BoolishT) -> BoolishT:
    ...

@overload
def validate_ascending(ascending: Sequence[BoolishT]) -> list[BoolishT]:
    ...

def validate_ascending(ascending: bool | int | Sequence[BoolishT]) -> bool | int | list[BoolishT]:
    kwargs = {'none_allowed': False, 'int_allowed': True}
    if not isinstance(ascending, Sequence):
        return validate_bool_kwarg(ascending, 'ascending', **kwargs)
    return [validate_bool_kwarg(item, 'ascending', **kwargs) for item in ascending]

def validate_endpoints(closed: str | None) -> tuple[bool, bool]:
    left_closed = False
    right_closed = False
    if closed is None:
        left_closed = True
        right_closed = True
    elif closed == 'left':
        left_closed = True
    elif closed == 'right':
        right_closed = True
    else:
        raise ValueError("Closed has to be either 'left', 'right' or None")
    return (left_closed, right_closed)

def validate_inclusive(inclusive: str | None) -> tuple[bool, bool]:
    left_right_inclusive: tuple[bool, bool] | None = None
    if isinstance(inclusive, str):
        left_right_inclusive = {'both': (True, True), 'left': (True, False), 'right': (False, True), 'neither': (False, False)}.get(inclusive)
    if left_right_inclusive is None:
        raise ValueError("Inclusive has to be either 'both', 'neither', 'left' or 'right'")
    return left_right_inclusive

def validate_insert_loc(loc: int, length: int) -> int:
    if not is_integer(loc):
        raise TypeError(f'loc must be an integer between -{length} and {length}')
    if loc < 0:
        loc += length
    if not 0 <= loc <= length:
        raise IndexError(f'loc must be an integer between -{length} and {length}')
    return loc

def check_dtype_backend(dtype_backend) -> None:
    if dtype_backend is not lib.no_default:
        if dtype_backend not in ['numpy_nullable', 'pyarrow']:
            raise ValueError(f"dtype_backend {dtype_backend} is invalid, only 'numpy_nullable' and 'pyarrow' are allowed.")
