from __future__ import annotations
from contextlib import contextmanager
import re
from typing import TYPE_CHECKING, Any, NamedTuple, cast
import warnings
from pandas._typing import F
from pandas.util._exceptions import find_stack_level
if TYPE_CHECKING:
    from collections.abc import Callable, Generator, Sequence

class DeprecatedOption(NamedTuple):
    key: str
    category: type[Warning]
    msg: str | None
    rkey: str | None
    removal_ver: str | None

class RegisteredOption(NamedTuple):
    key: str
    defval: Any
    doc: str
    validator: Callable[[object], Any] | None
    cb: Callable[[str], Any] | None
_deprecated_options: dict[str, DeprecatedOption] = {}
_registered_options: dict[str, RegisteredOption] = {}
_global_config: dict[str, Any] = {}
_reserved_keys: list[str] = ['all']

class OptionError(AttributeError, KeyError):
    __module__ = 'pandas.errors'

def _get_single_key(pat: str) -> str:
    keys = _select_options(pat)
    if len(keys) == 0:
        _warn_if_deprecated(pat)
        raise OptionError(f'No such keys(s): {pat!r}')
    if len(keys) > 1:
        raise OptionError('Pattern matched multiple keys')
    key = keys[0]
    _warn_if_deprecated(key)
    key = _translate_key(key)
    return key

def get_option(pat: str) -> Any:
    key = _get_single_key(pat)
    root, k = _get_root(key)
    return root[k]

def set_option(*args) -> None:
    if len(args) == 1 and isinstance(args[0], dict):
        args = tuple((kv for item in args[0].items() for kv in item))
    nargs = len(args)
    if not nargs or nargs % 2 != 0:
        raise ValueError('Must provide an even number of non-keyword arguments')
    for k, v in zip(args[::2], args[1::2], strict=True):
        key = _get_single_key(k)
        opt = _get_registered_option(key)
        if opt and opt.validator:
            opt.validator(v)
        root, k_root = _get_root(key)
        root[k_root] = v
        if opt.cb:
            opt.cb(key)

def describe_option(pat: str='', _print_desc: bool=True) -> str | None:
    keys = _select_options(pat)
    if len(keys) == 0:
        raise OptionError(f'No such keys(s) for pat={pat!r}')
    s = '\n'.join([_build_option_description(k) for k in keys])
    if _print_desc:
        print(s)
        return None
    return s

def reset_option(pat: str) -> None:
    keys = _select_options(pat)
    if len(keys) == 0:
        raise OptionError(f'No such keys(s) for pat={pat!r}')
    if len(keys) > 1 and len(pat) < 4 and (pat != 'all'):
        raise ValueError('You must specify at least 4 characters when resetting multiple keys, use the special keyword "all" to reset all the options to their default value')
    for k in keys:
        set_option(k, _registered_options[k].defval)

def get_default_val(pat: str):
    key = _get_single_key(pat)
    return _get_registered_option(key).defval

class DictWrapper:
    d: dict[str, Any]

    def __init__(self, d: dict[str, Any], prefix: str='') -> None:
        object.__setattr__(self, 'd', d)
        object.__setattr__(self, 'prefix', prefix)

    def __setattr__(self, key: str, val: Any) -> None:
        prefix = object.__getattribute__(self, 'prefix')
        if prefix:
            prefix += '.'
        prefix += key
        if key in self.d and (not isinstance(self.d[key], dict)):
            set_option(prefix, val)
        else:
            raise OptionError('You can only set the value of existing options')

    def __getattr__(self, key: str):
        prefix = object.__getattribute__(self, 'prefix')
        if prefix:
            prefix += '.'
        prefix += key
        try:
            v = object.__getattribute__(self, 'd')[key]
        except KeyError as err:
            raise OptionError('No such option') from err
        if isinstance(v, dict):
            return DictWrapper(v, prefix)
        else:
            return get_option(prefix)

    def __dir__(self) -> list[str]:
        return list(self.d.keys())
options = DictWrapper(_global_config)
object.__setattr__(options, '__module__', 'pandas')

@contextmanager
def option_context(*args) -> Generator[None]:
    if len(args) == 1 and isinstance(args[0], dict):
        args = tuple((kv for item in args[0].items() for kv in item))
    if len(args) % 2 != 0 or len(args) < 2:
        raise ValueError('Provide an even amount of arguments as option_context(pat, val, pat, val...).')
    ops = tuple(zip(args[::2], args[1::2], strict=True))
    undo: tuple[tuple[Any, Any], ...] = ()
    try:
        undo = tuple(((pat, get_option(pat)) for pat, val in ops))
        for pat, val in ops:
            set_option(pat, val)
        yield
    finally:
        for pat, val in undo:
            set_option(pat, val)

def register_option(key: str, defval: object, doc: str='', validator: Callable[[object], Any] | None=None, cb: Callable[[str], Any] | None=None) -> None:
    import keyword
    import tokenize
    key = key.lower()
    if key in _registered_options:
        raise OptionError(f"Option '{key}' has already been registered")
    if key in _reserved_keys:
        raise OptionError(f"Option '{key}' is a reserved key")
    if validator:
        validator(defval)
    path = key.split('.')
    for k in path:
        if not re.match('^' + tokenize.Name + '$', k):
            raise ValueError(f'{k} is not a valid identifier')
        if keyword.iskeyword(k):
            raise ValueError(f'{k} is a python keyword')
    cursor = _global_config
    msg = "Path prefix to option '{option}' is already an option"
    for i, p in enumerate(path[:-1]):
        if not isinstance(cursor, dict):
            raise OptionError(msg.format(option='.'.join(path[:i])))
        if p not in cursor:
            cursor[p] = {}
        cursor = cursor[p]
    if not isinstance(cursor, dict):
        raise OptionError(msg.format(option='.'.join(path[:-1])))
    cursor[path[-1]] = defval
    _registered_options[key] = RegisteredOption(key=key, defval=defval, doc=doc, validator=validator, cb=cb)

def deprecate_option(key: str, category: type[Warning], msg: str | None=None, rkey: str | None=None, removal_ver: str | None=None) -> None:
    key = key.lower()
    if key in _deprecated_options:
        raise OptionError(f"Option '{key}' has already been defined as deprecated.")
    _deprecated_options[key] = DeprecatedOption(key, category, msg, rkey, removal_ver)

def _select_options(pat: str) -> list[str]:
    if pat in _registered_options:
        return [pat]
    keys = sorted(_registered_options.keys())
    if pat == 'all':
        return keys
    return [k for k in keys if re.search(pat, k, re.I)]

def _get_root(key: str) -> tuple[dict[str, Any], str]:
    path = key.split('.')
    cursor = _global_config
    for p in path[:-1]:
        cursor = cursor[p]
    return (cursor, path[-1])

def _get_deprecated_option(key: str):
    try:
        d = _deprecated_options[key]
    except KeyError:
        return None
    else:
        return d

def _get_registered_option(key: str):
    return _registered_options.get(key)

def _translate_key(key: str) -> str:
    d = _get_deprecated_option(key)
    if d:
        return d.rkey or key
    else:
        return key

def _warn_if_deprecated(key: str) -> bool:
    d = _get_deprecated_option(key)
    if d:
        if d.msg:
            warnings.warn(d.msg, d.category, stacklevel=find_stack_level())
        else:
            msg = f"'{key}' is deprecated"
            if d.removal_ver:
                msg += f' and will be removed in {d.removal_ver}'
            if d.rkey:
                msg += f", please use '{d.rkey}' instead."
            else:
                msg += ', please refrain from using it.'
            warnings.warn(msg, d.category, stacklevel=find_stack_level())
        return True
    return False

def _build_option_description(k: str) -> str:
    o = _get_registered_option(k)
    d = _get_deprecated_option(k)
    s = f'{k} '
    if o.doc:
        s += '\n'.join(o.doc.strip().split('\n'))
    else:
        s += 'No description available.'
    if o:
        with warnings.catch_warnings():
            warnings.simplefilter('ignore', FutureWarning)
            warnings.simplefilter('ignore', DeprecationWarning)
            s += f'\n    [default: {o.defval}] [currently: {get_option(k)}]'
    if d:
        rkey = d.rkey or ''
        s += '\n    (Deprecated'
        s += f', use `{rkey}` instead.'
        s += ')'
    return s

@contextmanager
def config_prefix(prefix: str) -> Generator[None]:
    global register_option, get_option, set_option

    def wrap(func: F) -> F:

        def inner(key: str, *args, **kwds):
            pkey = f'{prefix}.{key}'
            return func(pkey, *args, **kwds)
        return cast(F, inner)
    _register_option = register_option
    _get_option = get_option
    _set_option = set_option
    set_option = wrap(set_option)
    get_option = wrap(get_option)
    register_option = wrap(register_option)
    try:
        yield
    finally:
        set_option = _set_option
        get_option = _get_option
        register_option = _register_option

def is_type_factory(_type: type[Any]) -> Callable[[Any], None]:

    def inner(x) -> None:
        if type(x) != _type:
            raise ValueError(f"Value must have type '{_type}'")
    return inner

def is_instance_factory(_type: type | tuple[type, ...]) -> Callable[[Any], None]:
    if isinstance(_type, tuple):
        type_repr = '|'.join(map(str, _type))
    else:
        type_repr = f"'{_type}'"

    def inner(x) -> None:
        if not isinstance(x, _type):
            raise ValueError(f'Value must be an instance of {type_repr}')
    return inner

def is_one_of_factory(legal_values: Sequence) -> Callable[[Any], None]:
    callables = [c for c in legal_values if callable(c)]
    legal_values = [c for c in legal_values if not callable(c)]

    def inner(x) -> None:
        if x not in legal_values:
            if not any((c(x) for c in callables)):
                uvals = [str(lval) for lval in legal_values]
                pp_values = '|'.join(uvals)
                msg = f'Value must be one of {pp_values}'
                if len(callables):
                    msg += ' or a callable'
                raise ValueError(msg)
    return inner

def is_nonnegative_int(value: object) -> None:
    if value is None:
        return
    elif isinstance(value, int):
        if value >= 0:
            return
    msg = 'Value must be a nonnegative integer or None'
    raise ValueError(msg)
is_int = is_type_factory(int)
is_bool = is_type_factory(bool)
is_float = is_type_factory(float)
is_str = is_type_factory(str)
is_text = is_instance_factory((str, bytes))

def is_callable(obj: object) -> bool:
    if not callable(obj):
        raise ValueError('Value must be a callable')
    return True
get_option.__module__ = 'pandas'
set_option.__module__ = 'pandas'
describe_option.__module__ = 'pandas'
reset_option.__module__ = 'pandas'
option_context.__module__ = 'pandas'
