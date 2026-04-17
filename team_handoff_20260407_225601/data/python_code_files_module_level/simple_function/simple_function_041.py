from __future__ import annotations
from functools import wraps
import inspect
from textwrap import dedent
from typing import TYPE_CHECKING, Any, cast
import warnings
from pandas._libs.properties import cache_readonly
from pandas._typing import F, T
from pandas.util._exceptions import find_stack_level
if TYPE_CHECKING:
    from collections.abc import Callable, Mapping
    from pandas.errors import PandasChangeWarning

def deprecate(klass: type[Warning], name: str, alternative: Callable[..., Any], version: str, alt_name: str | None=None, stacklevel: int=2, msg: str | None=None) -> Callable[[F], F]:
    alt_name = alt_name or alternative.__name__
    warning_msg = msg or f'{name} is deprecated, use {alt_name} instead.'

    @wraps(alternative)
    def wrapper(*args, **kwargs) -> Callable[..., Any]:
        warnings.warn(warning_msg, klass, stacklevel=stacklevel)
        return alternative(*args, **kwargs)
    msg = msg or f'Use `{alt_name}` instead.'
    doc_error_msg = f'deprecate needs a correctly formatted docstring in the target function (should have a one liner short summary, and opening quotes should be in their own line). Found:\n{alternative.__doc__}'
    if alternative.__doc__:
        if alternative.__doc__.count('\n') < 3:
            raise AssertionError(doc_error_msg)
        empty1, summary, empty2, doc_string = alternative.__doc__.split('\n', 3)
        if empty1 or (empty2 and (not summary)):
            raise AssertionError(doc_error_msg)
        wrapper.__doc__ = dedent(f'\n        {summary.strip()}\n\n        .. deprecated:: {version}\n            {msg}\n\n        {dedent(doc_string)}')
    return wrapper

def deprecate_kwarg(klass: type[Warning], old_arg_name: str, new_arg_name: str | None, mapping: Mapping[Any, Any] | Callable[[Any], Any] | None=None, stacklevel: int=2) -> Callable[[F], F]:
    if mapping is not None and (not hasattr(mapping, 'get')) and (not callable(mapping)):
        raise TypeError('mapping from old to new argument values must be dict or callable!')

    def _deprecate_kwarg(func: F) -> F:

        @wraps(func)
        def wrapper(*args, **kwargs) -> Callable[..., Any]:
            __tracebackhide__ = True
            old_arg_value = kwargs.pop(old_arg_name, None)
            if old_arg_value is not None:
                if new_arg_name is None:
                    msg = f'the {old_arg_name!r} keyword is deprecated and will be removed in a future version. Please take steps to stop the use of {old_arg_name!r}'
                    warnings.warn(msg, klass, stacklevel=stacklevel)
                    kwargs[old_arg_name] = old_arg_value
                    return func(*args, **kwargs)
                elif mapping is not None:
                    if callable(mapping):
                        new_arg_value = mapping(old_arg_value)
                    else:
                        new_arg_value = mapping.get(old_arg_value, old_arg_value)
                    msg = f'the {old_arg_name}={old_arg_value!r} keyword is deprecated, use {new_arg_name}={new_arg_value!r} instead.'
                else:
                    new_arg_value = old_arg_value
                    msg = f'the {old_arg_name!r} keyword is deprecated, use {new_arg_name!r} instead.'
                warnings.warn(msg, klass, stacklevel=stacklevel)
                if kwargs.get(new_arg_name) is not None:
                    msg = f'Can only specify {old_arg_name!r} or {new_arg_name!r}, not both.'
                    raise TypeError(msg)
                kwargs[new_arg_name] = new_arg_value
            return func(*args, **kwargs)
        return cast(F, wrapper)
    return _deprecate_kwarg

def _format_argument_list(allow_args: list[str]) -> str:
    if 'self' in allow_args:
        allow_args.remove('self')
    if not allow_args:
        return ''
    elif len(allow_args) == 1:
        return f" except for the argument '{allow_args[0]}'"
    else:
        last = allow_args[-1]
        args = ', '.join(["'" + x + "'" for x in allow_args[:-1]])
        return f" except for the arguments {args} and '{last}'"

def future_version_msg(version: str | None) -> str:
    if version is None:
        return 'In a future version of pandas'
    else:
        return f'Starting with pandas version {version}'

def deprecate_nonkeyword_arguments(klass: type[PandasChangeWarning], allowed_args: list[str] | None=None, name: str | None=None) -> Callable[[F], F]:

    def decorate(func):
        old_sig = inspect.signature(func)
        if allowed_args is not None:
            allow_args = allowed_args
        else:
            allow_args = [p.name for p in old_sig.parameters.values() if p.kind in (p.POSITIONAL_ONLY, p.POSITIONAL_OR_KEYWORD) and p.default is p.empty]
        new_params = [p.replace(kind=p.KEYWORD_ONLY) if p.kind in (p.POSITIONAL_ONLY, p.POSITIONAL_OR_KEYWORD) and p.name not in allow_args else p for p in old_sig.parameters.values()]
        new_params.sort(key=lambda p: p.kind)
        new_sig = old_sig.replace(parameters=new_params)
        num_allow_args = len(allow_args)
        msg = f'{future_version_msg(klass.version())} all arguments of {name or func.__qualname__}{{arguments}} will be keyword-only.'

        @wraps(func)
        def wrapper(*args, **kwargs):
            if len(args) > num_allow_args:
                warnings.warn(msg.format(arguments=_format_argument_list(allow_args)), klass, stacklevel=find_stack_level())
            return func(*args, **kwargs)
        wrapper.__signature__ = new_sig
        return wrapper
    return decorate

def doc(*docstrings: None | str | Callable, **params: object) -> Callable[[F], F]:

    def decorator(decorated: F) -> F:
        docstring_components: list[str | Callable] = []
        if decorated.__doc__:
            docstring_components.append(dedent(decorated.__doc__))
        for docstring in docstrings:
            if docstring is None:
                continue
            if hasattr(docstring, '_docstring_components'):
                docstring_components.extend(docstring._docstring_components)
            elif isinstance(docstring, str) or docstring.__doc__:
                docstring_components.append(docstring)
        params_applied = [component.format(**params) if isinstance(component, str) and len(params) > 0 else component for component in docstring_components]
        decorated.__doc__ = ''.join([component if isinstance(component, str) else dedent(component.__doc__ or '') for component in params_applied])
        decorated._docstring_components = docstring_components
        return decorated
    return decorator

class Substitution:

    def __init__(self, *args, **kwargs) -> None:
        if args and kwargs:
            raise AssertionError('Only positional or keyword args are allowed')
        self.params = args or kwargs

    def __call__(self, func: F) -> F:
        func.__doc__ = func.__doc__ and func.__doc__ % self.params
        return func

    def update(self, *args, **kwargs) -> None:
        if isinstance(self.params, dict):
            self.params.update(*args, **kwargs)

class Appender:
    addendum: str | None

    def __init__(self, addendum: str | None, join: str='', indents: int=0) -> None:
        if indents > 0:
            self.addendum = indent(addendum, indents=indents)
        else:
            self.addendum = addendum
        self.join = join

    def __call__(self, func: T) -> T:
        func.__doc__ = func.__doc__ if func.__doc__ else ''
        self.addendum = self.addendum if self.addendum else ''
        docitems = [func.__doc__, self.addendum]
        func.__doc__ = dedent(self.join.join(docitems))
        return func

def indent(text: str | None, indents: int=1) -> str:
    if not text or not isinstance(text, str):
        return ''
    jointext = ''.join(['\n'] + ['    '] * indents)
    return jointext.join(text.split('\n'))
__all__ = ['Appender', 'Substitution', 'cache_readonly', 'deprecate', 'deprecate_kwarg', 'deprecate_nonkeyword_arguments', 'doc', 'future_version_msg']

def set_module(module) -> Callable[[F], F]:

    def decorator(func: F) -> F:
        if module is not None:
            func.__module__ = module
        return func
    return decorator
