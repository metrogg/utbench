from __future__ import annotations
import builtins
from collections import abc, defaultdict
from collections.abc import Callable, Collection, Generator, Hashable, Iterable, Sequence
import contextlib
from functools import partial
import inspect
import sys
from typing import TYPE_CHECKING, Any, Concatenate, TypeVar, cast, overload
import numpy as np
from pandas._libs import lib
from pandas.core.dtypes.cast import construct_1d_object_array_from_listlike
from pandas.core.dtypes.common import is_bool_dtype, is_integer
from pandas.core.dtypes.generic import ABCExtensionArray, ABCIndex, ABCMultiIndex, ABCNumpyExtensionArray, ABCSeries
from pandas.core.dtypes.inference import iterable_not_string
from pandas.core.col import Expression
if TYPE_CHECKING:
    from pandas._typing import AnyArrayLike, ArrayLike, NpDtype, P, RandomState, T
    from pandas import Index

def flatten(line):
    for element in line:
        if iterable_not_string(element):
            yield from flatten(element)
        else:
            yield element

def consensus_name_attr(objs):
    name = objs[0].name
    for obj in objs[1:]:
        try:
            if obj.name != name:
                name = None
                break
        except ValueError:
            name = None
            break
    return name

def is_bool_indexer(key: Any) -> bool:
    if isinstance(key, (ABCSeries, np.ndarray, ABCIndex, ABCExtensionArray, ABCNumpyExtensionArray)) and (not isinstance(key, ABCMultiIndex)):
        if key.dtype == np.object_:
            key_array = np.asarray(key)
            if not lib.is_bool_array(key_array):
                na_msg = 'Cannot mask with non-boolean array containing NA / NaN values'
                if lib.is_bool_array(key_array, skipna=True):
                    raise ValueError(na_msg)
                return False
            return True
        elif is_bool_dtype(key.dtype):
            return True
    elif isinstance(key, list):
        if len(key) > 0:
            if type(key) is not list:
                key = list(key)
            return lib.is_bool_list(key)
    return False

def cast_scalar_indexer(val):
    if lib.is_float(val) and val.is_integer():
        raise IndexError('Indexing with a float is no longer supported. Manually convert to an integer key instead.')
    return val

def not_none(*args):
    return (arg for arg in args if arg is not None)

def any_none(*args) -> bool:
    return any((arg is None for arg in args))

def all_none(*args) -> bool:
    return all((arg is None for arg in args))

def any_not_none(*args) -> bool:
    return any((arg is not None for arg in args))

def all_not_none(*args) -> bool:
    return all((arg is not None for arg in args))

def count_not_none(*args) -> int:
    return sum((x is not None for x in args))

@overload
def asarray_tuplesafe(values: ArrayLike | list | tuple | zip, dtype: NpDtype | None=...) -> np.ndarray:
    ...

@overload
def asarray_tuplesafe(values: Iterable, dtype: NpDtype | None=...) -> ArrayLike:
    ...

def asarray_tuplesafe(values: Iterable, dtype: NpDtype | None=None) -> ArrayLike:
    if not (isinstance(values, (list, tuple)) or hasattr(values, '__array__')):
        values = list(values)
    elif isinstance(values, ABCIndex):
        return values._values
    elif isinstance(values, ABCSeries):
        return values._values
    if isinstance(values, list) and dtype in [np.object_, object]:
        return construct_1d_object_array_from_listlike(values)
    try:
        result = np.asarray(values, dtype=dtype)
    except ValueError:
        return construct_1d_object_array_from_listlike(values)
    if issubclass(result.dtype.type, str):
        result = np.asarray(values, dtype=object)
    if result.ndim == 2:
        values = [tuple(x) for x in values]
        result = construct_1d_object_array_from_listlike(values)
    return result

def index_labels_to_array(labels: np.ndarray | Iterable, dtype: NpDtype | None=None) -> np.ndarray:
    if isinstance(labels, (str, tuple)):
        labels = [labels]
    if not isinstance(labels, (list, np.ndarray)):
        try:
            labels = list(labels)
        except TypeError:
            labels = [labels]
    rlabels = asarray_tuplesafe(labels, dtype=dtype)
    return rlabels

def maybe_make_list(obj):
    if obj is not None and (not isinstance(obj, (tuple, list))):
        return [obj]
    return obj

def maybe_iterable_to_list(obj: Iterable[T] | T) -> Collection[T] | T:
    if isinstance(obj, abc.Iterable) and (not isinstance(obj, abc.Sized)):
        return list(obj)
    obj = cast(Collection, obj)
    return obj

def is_null_slice(obj) -> bool:
    return isinstance(obj, slice) and obj.start is None and (obj.stop is None) and (obj.step is None)

def is_empty_slice(obj) -> bool:
    return isinstance(obj, slice) and obj.start is not None and (obj.stop is not None) and (obj.start == obj.stop)

def is_true_slices(line: abc.Iterable) -> abc.Generator[bool, None, None]:
    for k in line:
        yield (isinstance(k, slice) and (not is_null_slice(k)))

def is_full_slice(obj, line: int) -> bool:
    return isinstance(obj, slice) and obj.start == 0 and (obj.stop == line) and (obj.step is None)

def get_callable_name(obj):
    if hasattr(obj, '__name__'):
        return obj.__name__
    if isinstance(obj, partial):
        return get_callable_name(obj.func)
    if callable(obj):
        return type(obj).__name__
    return None

def apply_if_callable(maybe_callable, obj, **kwargs):
    if isinstance(maybe_callable, Expression):
        return maybe_callable._eval_expression(obj, **kwargs)
    elif callable(maybe_callable):
        return maybe_callable(obj, **kwargs)
    return maybe_callable

def standardize_mapping(into):
    if not inspect.isclass(into):
        if isinstance(into, defaultdict):
            return partial(defaultdict, into.default_factory)
        into = type(into)
    if not issubclass(into, abc.Mapping):
        raise TypeError(f'unsupported type: {into}')
    if into == defaultdict:
        raise TypeError('to_dict() only accepts initialized defaultdicts')
    return into

@overload
def random_state(state: np.random.Generator) -> np.random.Generator:
    ...

@overload
def random_state(state: int | np.ndarray | np.random.BitGenerator | np.random.RandomState | None) -> np.random.RandomState:
    ...

def random_state(state: RandomState | None=None):
    if is_integer(state) or isinstance(state, (np.ndarray, np.random.BitGenerator)):
        return np.random.RandomState(state)
    elif isinstance(state, np.random.RandomState):
        return state
    elif isinstance(state, np.random.Generator):
        return state
    elif state is None:
        return np.random
    else:
        raise ValueError('random_state must be an integer, array-like, a BitGenerator, Generator, a numpy RandomState, or None')
_T = TypeVar('_T')

@overload
def pipe(obj: _T, func: Callable[Concatenate[_T, P], T], *args: P.args, **kwargs: P.kwargs) -> T:
    ...

@overload
def pipe(obj: Any, func: tuple[Callable[..., T], str], *args: Any, **kwargs: Any) -> T:
    ...

def pipe(obj: _T, func: Callable[Concatenate[_T, P], T] | tuple[Callable[..., T], str], *args: Any, **kwargs: Any) -> T:
    if isinstance(func, tuple):
        func_, target = func
        if target in kwargs:
            msg = f'{target} is both the pipe target and a keyword argument'
            raise ValueError(msg)
        kwargs[target] = obj
        return func_(*args, **kwargs)
    else:
        return func(obj, *args, **kwargs)

def get_rename_function(mapper):

    def f(x):
        if x in mapper:
            return mapper[x]
        else:
            return x
    return f if isinstance(mapper, (abc.Mapping, ABCSeries)) else mapper

def convert_to_list_like(values: Hashable | Iterable | AnyArrayLike) -> list | AnyArrayLike:
    if isinstance(values, (list, np.ndarray, ABCIndex, ABCSeries, ABCExtensionArray)):
        return values
    elif isinstance(values, abc.Iterable) and (not isinstance(values, str)):
        return list(values)
    return [values]

@contextlib.contextmanager
def temp_setattr(obj, attr: str, value, condition: bool=True) -> Generator[None]:
    if condition:
        old_value = getattr(obj, attr)
        setattr(obj, attr, value)
    try:
        yield obj
    finally:
        if condition:
            setattr(obj, attr, old_value)

def require_length_match(data, index: Index) -> None:
    if len(data) != len(index):
        raise ValueError(f'Length of values ({len(data)}) does not match length of index ({len(index)})')
_cython_table = {builtins.sum: 'sum', builtins.max: 'max', builtins.min: 'min', np.all: 'all', np.any: 'any', np.sum: 'sum', np.nansum: 'sum', np.mean: 'mean', np.nanmean: 'mean', np.prod: 'prod', np.nanprod: 'prod', np.std: 'std', np.nanstd: 'std', np.var: 'var', np.nanvar: 'var', np.median: 'median', np.nanmedian: 'median', np.max: 'max', np.nanmax: 'max', np.min: 'min', np.nanmin: 'min', np.cumprod: 'cumprod', np.nancumprod: 'cumprod', np.cumsum: 'cumsum', np.nancumsum: 'cumsum'}

def get_cython_func(arg: Callable) -> str | None:
    return _cython_table.get(arg)

def fill_missing_names(names: Sequence[Hashable | None]) -> list[Hashable]:
    return [f'level_{i}' if name is None else name for i, name in enumerate(names)]

def is_local_in_caller_frame(obj):
    frame = sys._getframe(2)
    for v in frame.f_locals.values():
        if v is obj:
            return True
    return False
