from __future__ import annotations
import inspect
import typing as t
from functools import update_wrapper
from gettext import gettext as _
from .core import Argument
from .core import Command
from .core import Context
from .core import Group
from .core import Option
from .core import Parameter
from .globals import get_current_context
from .utils import echo
if t.TYPE_CHECKING:
    import typing_extensions as te
    P = te.ParamSpec('P')
R = t.TypeVar('R')
T = t.TypeVar('T')
_AnyCallable = t.Callable[..., t.Any]
FC = t.TypeVar('FC', bound='_AnyCallable | Command')

def pass_context(f: t.Callable[te.Concatenate[Context, P], R]) -> t.Callable[P, R]:

    def new_func(*args: P.args, **kwargs: P.kwargs) -> R:
        return f(get_current_context(), *args, **kwargs)
    return update_wrapper(new_func, f)

def pass_obj(f: t.Callable[te.Concatenate[T, P], R]) -> t.Callable[P, R]:

    def new_func(*args: P.args, **kwargs: P.kwargs) -> R:
        return f(get_current_context().obj, *args, **kwargs)
    return update_wrapper(new_func, f)

def make_pass_decorator(object_type: type[T], ensure: bool=False) -> t.Callable[[t.Callable[te.Concatenate[T, P], R]], t.Callable[P, R]]:

    def decorator(f: t.Callable[te.Concatenate[T, P], R]) -> t.Callable[P, R]:

        def new_func(*args: P.args, **kwargs: P.kwargs) -> R:
            ctx = get_current_context()
            obj: T | None
            if ensure:
                obj = ctx.ensure_object(object_type)
            else:
                obj = ctx.find_object(object_type)
            if obj is None:
                raise RuntimeError(f'Managed to invoke callback without a context object of type {object_type.__name__!r} existing.')
            return ctx.invoke(f, obj, *args, **kwargs)
        return update_wrapper(new_func, f)
    return decorator

def pass_meta_key(key: str, *, doc_description: str | None=None) -> t.Callable[[t.Callable[te.Concatenate[T, P], R]], t.Callable[P, R]]:

    def decorator(f: t.Callable[te.Concatenate[T, P], R]) -> t.Callable[P, R]:

        def new_func(*args: P.args, **kwargs: P.kwargs) -> R:
            ctx = get_current_context()
            obj = ctx.meta[key]
            return ctx.invoke(f, obj, *args, **kwargs)
        return update_wrapper(new_func, f)
    if doc_description is None:
        doc_description = f'the {key!r} key from :attr:`click.Context.meta`'
    decorator.__doc__ = f'Decorator that passes {doc_description} as the first argument to the decorated function.'
    return decorator
CmdType = t.TypeVar('CmdType', bound=Command)

@t.overload
def command(name: _AnyCallable) -> Command:
    ...

@t.overload
def command(name: str | None, cls: type[CmdType], **attrs: t.Any) -> t.Callable[[_AnyCallable], CmdType]:
    ...

@t.overload
def command(name: None=None, *, cls: type[CmdType], **attrs: t.Any) -> t.Callable[[_AnyCallable], CmdType]:
    ...

@t.overload
def command(name: str | None=..., cls: None=None, **attrs: t.Any) -> t.Callable[[_AnyCallable], Command]:
    ...

def command(name: str | _AnyCallable | None=None, cls: type[CmdType] | None=None, **attrs: t.Any) -> Command | t.Callable[[_AnyCallable], Command | CmdType]:
    func: t.Callable[[_AnyCallable], t.Any] | None = None
    if callable(name):
        func = name
        name = None
        assert cls is None, "Use 'command(cls=cls)(callable)' to specify a class."
        assert not attrs, "Use 'command(**kwargs)(callable)' to provide arguments."
    if cls is None:
        cls = t.cast('type[CmdType]', Command)

    def decorator(f: _AnyCallable) -> CmdType:
        if isinstance(f, Command):
            raise TypeError('Attempted to convert a callback into a command twice.')
        attr_params = attrs.pop('params', None)
        params = attr_params if attr_params is not None else []
        try:
            decorator_params = f.__click_params__
        except AttributeError:
            pass
        else:
            del f.__click_params__
            params.extend(reversed(decorator_params))
        if attrs.get('help') is None:
            attrs['help'] = f.__doc__
        if t.TYPE_CHECKING:
            assert cls is not None
            assert not callable(name)
        if name is not None:
            cmd_name = name
        else:
            cmd_name = f.__name__.lower().replace('_', '-')
            cmd_left, sep, suffix = cmd_name.rpartition('-')
            if sep and suffix in {'command', 'cmd', 'group', 'grp'}:
                cmd_name = cmd_left
        cmd = cls(name=cmd_name, callback=f, params=params, **attrs)
        cmd.__doc__ = f.__doc__
        return cmd
    if func is not None:
        return decorator(func)
    return decorator
GrpType = t.TypeVar('GrpType', bound=Group)

@t.overload
def group(name: _AnyCallable) -> Group:
    ...

@t.overload
def group(name: str | None, cls: type[GrpType], **attrs: t.Any) -> t.Callable[[_AnyCallable], GrpType]:
    ...

@t.overload
def group(name: None=None, *, cls: type[GrpType], **attrs: t.Any) -> t.Callable[[_AnyCallable], GrpType]:
    ...

@t.overload
def group(name: str | None=..., cls: None=None, **attrs: t.Any) -> t.Callable[[_AnyCallable], Group]:
    ...

def group(name: str | _AnyCallable | None=None, cls: type[GrpType] | None=None, **attrs: t.Any) -> Group | t.Callable[[_AnyCallable], Group | GrpType]:
    if cls is None:
        cls = t.cast('type[GrpType]', Group)
    if callable(name):
        return command(cls=cls, **attrs)(name)
    return command(name, cls, **attrs)

def _param_memo(f: t.Callable[..., t.Any], param: Parameter) -> None:
    if isinstance(f, Command):
        f.params.append(param)
    else:
        if not hasattr(f, '__click_params__'):
            f.__click_params__ = []
        f.__click_params__.append(param)

def argument(*param_decls: str, cls: type[Argument] | None=None, **attrs: t.Any) -> t.Callable[[FC], FC]:
    if cls is None:
        cls = Argument

    def decorator(f: FC) -> FC:
        _param_memo(f, cls(param_decls, **attrs))
        return f
    return decorator

def option(*param_decls: str, cls: type[Option] | None=None, **attrs: t.Any) -> t.Callable[[FC], FC]:
    if cls is None:
        cls = Option

    def decorator(f: FC) -> FC:
        _param_memo(f, cls(param_decls, **attrs))
        return f
    return decorator

def confirmation_option(*param_decls: str, **kwargs: t.Any) -> t.Callable[[FC], FC]:

    def callback(ctx: Context, param: Parameter, value: bool) -> None:
        if not value:
            ctx.abort()
    if not param_decls:
        param_decls = ('--yes',)
    kwargs.setdefault('is_flag', True)
    kwargs.setdefault('callback', callback)
    kwargs.setdefault('expose_value', False)
    kwargs.setdefault('prompt', 'Do you want to continue?')
    kwargs.setdefault('help', 'Confirm the action without prompting.')
    return option(*param_decls, **kwargs)

def password_option(*param_decls: str, **kwargs: t.Any) -> t.Callable[[FC], FC]:
    if not param_decls:
        param_decls = ('--password',)
    kwargs.setdefault('prompt', True)
    kwargs.setdefault('confirmation_prompt', True)
    kwargs.setdefault('hide_input', True)
    return option(*param_decls, **kwargs)

def version_option(version: str | None=None, *param_decls: str, package_name: str | None=None, prog_name: str | None=None, message: str | None=None, **kwargs: t.Any) -> t.Callable[[FC], FC]:
    if message is None:
        message = _('%(prog)s, version %(version)s')
    if version is None and package_name is None:
        frame = inspect.currentframe()
        f_back = frame.f_back if frame is not None else None
        f_globals = f_back.f_globals if f_back is not None else None
        del frame
        if f_globals is not None:
            package_name = f_globals.get('__name__')
            if package_name == '__main__':
                package_name = f_globals.get('__package__')
            if package_name:
                package_name = package_name.partition('.')[0]

    def callback(ctx: Context, param: Parameter, value: bool) -> None:
        if not value or ctx.resilient_parsing:
            return
        nonlocal prog_name
        nonlocal version
        if prog_name is None:
            prog_name = ctx.find_root().info_name
        if version is None and package_name is not None:
            import importlib.metadata
            try:
                version = importlib.metadata.version(package_name)
            except importlib.metadata.PackageNotFoundError:
                raise RuntimeError(f"{package_name!r} is not installed. Try passing 'package_name' instead.") from None
        if version is None:
            raise RuntimeError(f'Could not determine the version for {package_name!r} automatically.')
        echo(message % {'prog': prog_name, 'package': package_name, 'version': version}, color=ctx.color)
        ctx.exit()
    if not param_decls:
        param_decls = ('--version',)
    kwargs.setdefault('is_flag', True)
    kwargs.setdefault('expose_value', False)
    kwargs.setdefault('is_eager', True)
    kwargs.setdefault('help', _('Show the version and exit.'))
    kwargs['callback'] = callback
    return option(*param_decls, **kwargs)

def help_option(*param_decls: str, **kwargs: t.Any) -> t.Callable[[FC], FC]:

    def show_help(ctx: Context, param: Parameter, value: bool) -> None:
        if value and (not ctx.resilient_parsing):
            echo(ctx.get_help(), color=ctx.color)
            ctx.exit()
    if not param_decls:
        param_decls = ('--help',)
    kwargs.setdefault('is_flag', True)
    kwargs.setdefault('expose_value', False)
    kwargs.setdefault('is_eager', True)
    kwargs.setdefault('help', _('Show this message and exit.'))
    kwargs.setdefault('callback', show_help)
    return option(*param_decls, **kwargs)
