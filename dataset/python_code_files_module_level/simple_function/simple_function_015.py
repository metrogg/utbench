from __future__ import annotations
import collections.abc as cabc
import inspect
import io
import itertools
import sys
import typing as t
from contextlib import AbstractContextManager
from gettext import gettext as _
from ._compat import isatty
from ._compat import strip_ansi
from .exceptions import Abort
from .exceptions import UsageError
from .globals import resolve_color_default
from .types import Choice
from .types import convert_type
from .types import ParamType
from .utils import echo
from .utils import LazyFile
if t.TYPE_CHECKING:
    from ._termui_impl import ProgressBar
V = t.TypeVar('V')
visible_prompt_func: t.Callable[[str], str] = input
_ansi_colors = {'black': 30, 'red': 31, 'green': 32, 'yellow': 33, 'blue': 34, 'magenta': 35, 'cyan': 36, 'white': 37, 'reset': 39, 'bright_black': 90, 'bright_red': 91, 'bright_green': 92, 'bright_yellow': 93, 'bright_blue': 94, 'bright_magenta': 95, 'bright_cyan': 96, 'bright_white': 97}
_ansi_reset_all = '\x1b[0m'

def hidden_prompt_func(prompt: str) -> str:
    import getpass
    return getpass.getpass(prompt)

def _build_prompt(text: str, suffix: str, show_default: bool=False, default: t.Any | None=None, show_choices: bool=True, type: ParamType | None=None) -> str:
    prompt = text
    if type is not None and show_choices and isinstance(type, Choice):
        prompt += f" ({', '.join(map(str, type.choices))})"
    if default is not None and show_default:
        prompt = f'{prompt} [{_format_default(default)}]'
    return f'{prompt}{suffix}'

def _format_default(default: t.Any) -> t.Any:
    if isinstance(default, (io.IOBase, LazyFile)) and hasattr(default, 'name'):
        return default.name
    return default

def prompt(text: str, default: t.Any | None=None, hide_input: bool=False, confirmation_prompt: bool | str=False, type: ParamType | t.Any | None=None, value_proc: t.Callable[[str], t.Any] | None=None, prompt_suffix: str=': ', show_default: bool=True, err: bool=False, show_choices: bool=True) -> t.Any:

    def prompt_func(text: str) -> str:
        f = hidden_prompt_func if hide_input else visible_prompt_func
        try:
            echo(text[:-1], nl=False, err=err)
            return f(text[-1:])
        except (KeyboardInterrupt, EOFError):
            if hide_input:
                echo(None, err=err)
            raise Abort() from None
    if value_proc is None:
        value_proc = convert_type(type, default)
    prompt = _build_prompt(text, prompt_suffix, show_default, default, show_choices, type)
    if confirmation_prompt:
        if confirmation_prompt is True:
            confirmation_prompt = _('Repeat for confirmation')
        confirmation_prompt = _build_prompt(confirmation_prompt, prompt_suffix)
    while True:
        while True:
            value = prompt_func(prompt)
            if value:
                break
            elif default is not None:
                value = default
                break
        try:
            result = value_proc(value)
        except UsageError as e:
            if hide_input:
                echo(_('Error: The value you entered was invalid.'), err=err)
            else:
                echo(_('Error: {e.message}').format(e=e), err=err)
            continue
        if not confirmation_prompt:
            return result
        while True:
            value2 = prompt_func(confirmation_prompt)
            is_empty = not value and (not value2)
            if value2 or is_empty:
                break
        if value == value2:
            return result
        echo(_('Error: The two entered values do not match.'), err=err)

def confirm(text: str, default: bool | None=False, abort: bool=False, prompt_suffix: str=': ', show_default: bool=True, err: bool=False) -> bool:
    prompt = _build_prompt(text, prompt_suffix, show_default, 'y/n' if default is None else 'Y/n' if default else 'y/N')
    while True:
        try:
            echo(prompt[:-1], nl=False, err=err)
            value = visible_prompt_func(prompt[-1:]).lower().strip()
        except (KeyboardInterrupt, EOFError):
            raise Abort() from None
        if value in ('y', 'yes'):
            rv = True
        elif value in ('n', 'no'):
            rv = False
        elif default is not None and value == '':
            rv = default
        else:
            echo(_('Error: invalid input'), err=err)
            continue
        break
    if abort and (not rv):
        raise Abort()
    return rv

def echo_via_pager(text_or_generator: cabc.Iterable[str] | t.Callable[[], cabc.Iterable[str]] | str, color: bool | None=None) -> None:
    color = resolve_color_default(color)
    if inspect.isgeneratorfunction(text_or_generator):
        i = t.cast('t.Callable[[], cabc.Iterable[str]]', text_or_generator)()
    elif isinstance(text_or_generator, str):
        i = [text_or_generator]
    else:
        i = iter(t.cast('cabc.Iterable[str]', text_or_generator))
    text_generator = (el if isinstance(el, str) else str(el) for el in i)
    from ._termui_impl import pager
    return pager(itertools.chain(text_generator, '\n'), color)

@t.overload
def progressbar(*, length: int, label: str | None=None, hidden: bool=False, show_eta: bool=True, show_percent: bool | None=None, show_pos: bool=False, fill_char: str='#', empty_char: str='-', bar_template: str='%(label)s  [%(bar)s]  %(info)s', info_sep: str='  ', width: int=36, file: t.TextIO | None=None, color: bool | None=None, update_min_steps: int=1) -> ProgressBar[int]:
    ...

@t.overload
def progressbar(iterable: cabc.Iterable[V] | None=None, length: int | None=None, label: str | None=None, hidden: bool=False, show_eta: bool=True, show_percent: bool | None=None, show_pos: bool=False, item_show_func: t.Callable[[V | None], str | None] | None=None, fill_char: str='#', empty_char: str='-', bar_template: str='%(label)s  [%(bar)s]  %(info)s', info_sep: str='  ', width: int=36, file: t.TextIO | None=None, color: bool | None=None, update_min_steps: int=1) -> ProgressBar[V]:
    ...

def progressbar(iterable: cabc.Iterable[V] | None=None, length: int | None=None, label: str | None=None, hidden: bool=False, show_eta: bool=True, show_percent: bool | None=None, show_pos: bool=False, item_show_func: t.Callable[[V | None], str | None] | None=None, fill_char: str='#', empty_char: str='-', bar_template: str='%(label)s  [%(bar)s]  %(info)s', info_sep: str='  ', width: int=36, file: t.TextIO | None=None, color: bool | None=None, update_min_steps: int=1) -> ProgressBar[V]:
    from ._termui_impl import ProgressBar
    color = resolve_color_default(color)
    return ProgressBar(iterable=iterable, length=length, hidden=hidden, show_eta=show_eta, show_percent=show_percent, show_pos=show_pos, item_show_func=item_show_func, fill_char=fill_char, empty_char=empty_char, bar_template=bar_template, info_sep=info_sep, file=file, label=label, width=width, color=color, update_min_steps=update_min_steps)

def clear() -> None:
    if not isatty(sys.stdout):
        return
    echo('\x1b[2J\x1b[1;1H', nl=False)

def _interpret_color(color: int | tuple[int, int, int] | str, offset: int=0) -> str:
    if isinstance(color, int):
        return f'{38 + offset};5;{color:d}'
    if isinstance(color, (tuple, list)):
        r, g, b = color
        return f'{38 + offset};2;{r:d};{g:d};{b:d}'
    return str(_ansi_colors[color] + offset)

def style(text: t.Any, fg: int | tuple[int, int, int] | str | None=None, bg: int | tuple[int, int, int] | str | None=None, bold: bool | None=None, dim: bool | None=None, underline: bool | None=None, overline: bool | None=None, italic: bool | None=None, blink: bool | None=None, reverse: bool | None=None, strikethrough: bool | None=None, reset: bool=True) -> str:
    if not isinstance(text, str):
        text = str(text)
    bits = []
    if fg:
        try:
            bits.append(f'\x1b[{_interpret_color(fg)}m')
        except KeyError:
            raise TypeError(f'Unknown color {fg!r}') from None
    if bg:
        try:
            bits.append(f'\x1b[{_interpret_color(bg, 10)}m')
        except KeyError:
            raise TypeError(f'Unknown color {bg!r}') from None
    if bold is not None:
        bits.append(f'\x1b[{(1 if bold else 22)}m')
    if dim is not None:
        bits.append(f'\x1b[{(2 if dim else 22)}m')
    if underline is not None:
        bits.append(f'\x1b[{(4 if underline else 24)}m')
    if overline is not None:
        bits.append(f'\x1b[{(53 if overline else 55)}m')
    if italic is not None:
        bits.append(f'\x1b[{(3 if italic else 23)}m')
    if blink is not None:
        bits.append(f'\x1b[{(5 if blink else 25)}m')
    if reverse is not None:
        bits.append(f'\x1b[{(7 if reverse else 27)}m')
    if strikethrough is not None:
        bits.append(f'\x1b[{(9 if strikethrough else 29)}m')
    bits.append(text)
    if reset:
        bits.append(_ansi_reset_all)
    return ''.join(bits)

def unstyle(text: str) -> str:
    return strip_ansi(text)

def secho(message: t.Any | None=None, file: t.IO[t.AnyStr] | None=None, nl: bool=True, err: bool=False, color: bool | None=None, **styles: t.Any) -> None:
    if message is not None and (not isinstance(message, (bytes, bytearray))):
        message = style(message, **styles)
    return echo(message, file=file, nl=nl, err=err, color=color)

@t.overload
def edit(text: bytes | bytearray, editor: str | None=None, env: cabc.Mapping[str, str] | None=None, require_save: bool=False, extension: str='.txt') -> bytes | None:
    ...

@t.overload
def edit(text: str, editor: str | None=None, env: cabc.Mapping[str, str] | None=None, require_save: bool=True, extension: str='.txt') -> str | None:
    ...

@t.overload
def edit(text: None=None, editor: str | None=None, env: cabc.Mapping[str, str] | None=None, require_save: bool=True, extension: str='.txt', filename: str | cabc.Iterable[str] | None=None) -> None:
    ...

def edit(text: str | bytes | bytearray | None=None, editor: str | None=None, env: cabc.Mapping[str, str] | None=None, require_save: bool=True, extension: str='.txt', filename: str | cabc.Iterable[str] | None=None) -> str | bytes | bytearray | None:
    from ._termui_impl import Editor
    ed = Editor(editor=editor, env=env, require_save=require_save, extension=extension)
    if filename is None:
        return ed.edit(text)
    if isinstance(filename, str):
        filename = (filename,)
    ed.edit_files(filenames=filename)
    return None

def launch(url: str, wait: bool=False, locate: bool=False) -> int:
    from ._termui_impl import open_url
    return open_url(url, wait=wait, locate=locate)
_getchar: t.Callable[[bool], str] | None = None

def getchar(echo: bool=False) -> str:
    global _getchar
    if _getchar is None:
        from ._termui_impl import getchar as f
        _getchar = f
    return _getchar(echo)

def raw_terminal() -> AbstractContextManager[int]:
    from ._termui_impl import raw_terminal as f
    return f()

def pause(info: str | None=None, err: bool=False) -> None:
    if not isatty(sys.stdin) or not isatty(sys.stdout):
        return
    if info is None:
        info = _('Press any key to continue...')
    try:
        if info:
            echo(info, nl=False, err=err)
        try:
            getchar()
        except (KeyboardInterrupt, EOFError):
            pass
    finally:
        if info:
            echo(err=err)
