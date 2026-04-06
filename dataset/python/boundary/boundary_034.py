from __future__ import annotations
import collections.abc as cabc
import os
import re
import sys
import typing as t
from functools import update_wrapper
from types import ModuleType
from types import TracebackType
from ._compat import _default_text_stderr
from ._compat import _default_text_stdout
from ._compat import _find_binary_writer
from ._compat import auto_wrap_for_ansi
from ._compat import binary_streams
from ._compat import open_stream
from ._compat import should_strip_ansi
from ._compat import strip_ansi
from ._compat import text_streams
from ._compat import WIN
from .globals import resolve_color_default
if t.TYPE_CHECKING:
    import typing_extensions as te
    P = te.ParamSpec('P')
R = t.TypeVar('R')

def _posixify(name: str) -> str:
    return '-'.join(name.split()).lower()

def safecall(func: t.Callable[P, R]) -> t.Callable[P, R | None]:

    def wrapper(*args: P.args, **kwargs: P.kwargs) -> R | None:
        try:
            return func(*args, **kwargs)
        except Exception:
            pass
        return None
    return update_wrapper(wrapper, func)

def make_str(value: t.Any) -> str:
    if isinstance(value, bytes):
        try:
            return value.decode(sys.getfilesystemencoding())
        except UnicodeError:
            return value.decode('utf-8', 'replace')
    return str(value)

def make_default_short_help(help: str, max_length: int=45) -> str:
    paragraph_end = help.find('\n\n')
    if paragraph_end != -1:
        help = help[:paragraph_end]
    words = help.split()
    if not words:
        return ''
    if words[0] == '\x08':
        words = words[1:]
    total_length = 0
    last_index = len(words) - 1
    for i, word in enumerate(words):
        total_length += len(word) + (i > 0)
        if total_length > max_length:
            break
        if word[-1] == '.':
            return ' '.join(words[:i + 1])
        if total_length == max_length and i != last_index:
            break
    else:
        return ' '.join(words)
    total_length += len('...')
    while i > 0:
        total_length -= len(words[i]) + (i > 0)
        if total_length <= max_length:
            break
        i -= 1
    return ' '.join(words[:i]) + '...'

class LazyFile:

    def __init__(self, filename: str | os.PathLike[str], mode: str='r', encoding: str | None=None, errors: str | None='strict', atomic: bool=False):
        self.name: str = os.fspath(filename)
        self.mode = mode
        self.encoding = encoding
        self.errors = errors
        self.atomic = atomic
        self._f: t.IO[t.Any] | None
        self.should_close: bool
        if self.name == '-':
            self._f, self.should_close = open_stream(filename, mode, encoding, errors)
        else:
            if 'r' in mode:
                open(filename, mode).close()
            self._f = None
            self.should_close = True

    def __getattr__(self, name: str) -> t.Any:
        return getattr(self.open(), name)

    def __repr__(self) -> str:
        if self._f is not None:
            return repr(self._f)
        return f"<unopened file '{format_filename(self.name)}' {self.mode}>"

    def open(self) -> t.IO[t.Any]:
        if self._f is not None:
            return self._f
        try:
            rv, self.should_close = open_stream(self.name, self.mode, self.encoding, self.errors, atomic=self.atomic)
        except OSError as e:
            from .exceptions import FileError
            raise FileError(self.name, hint=e.strerror) from e
        self._f = rv
        return rv

    def close(self) -> None:
        if self._f is not None:
            self._f.close()

    def close_intelligently(self) -> None:
        if self.should_close:
            self.close()

    def __enter__(self) -> LazyFile:
        return self

    def __exit__(self, exc_type: type[BaseException] | None, exc_value: BaseException | None, tb: TracebackType | None) -> None:
        self.close_intelligently()

    def __iter__(self) -> cabc.Iterator[t.AnyStr]:
        self.open()
        return iter(self._f)

class KeepOpenFile:

    def __init__(self, file: t.IO[t.Any]) -> None:
        self._file: t.IO[t.Any] = file

    def __getattr__(self, name: str) -> t.Any:
        return getattr(self._file, name)

    def __enter__(self) -> KeepOpenFile:
        return self

    def __exit__(self, exc_type: type[BaseException] | None, exc_value: BaseException | None, tb: TracebackType | None) -> None:
        pass

    def __repr__(self) -> str:
        return repr(self._file)

    def __iter__(self) -> cabc.Iterator[t.AnyStr]:
        return iter(self._file)

def echo(message: t.Any | None=None, file: t.IO[t.Any] | None=None, nl: bool=True, err: bool=False, color: bool | None=None) -> None:
    if file is None:
        if err:
            file = _default_text_stderr()
        else:
            file = _default_text_stdout()
        if file is None:
            return
    if message is not None and (not isinstance(message, (str, bytes, bytearray))):
        out: str | bytes | bytearray | None = str(message)
    else:
        out = message
    if nl:
        out = out or ''
        if isinstance(out, str):
            out += '\n'
        else:
            out += b'\n'
    if not out:
        file.flush()
        return
    if isinstance(out, (bytes, bytearray)):
        binary_file = _find_binary_writer(file)
        if binary_file is not None:
            file.flush()
            binary_file.write(out)
            binary_file.flush()
            return
    else:
        color = resolve_color_default(color)
        if should_strip_ansi(file, color):
            out = strip_ansi(out)
        elif WIN:
            if auto_wrap_for_ansi is not None:
                file = auto_wrap_for_ansi(file, color)
            elif not color:
                out = strip_ansi(out)
    file.write(out)
    file.flush()

def get_binary_stream(name: t.Literal['stdin', 'stdout', 'stderr']) -> t.BinaryIO:
    opener = binary_streams.get(name)
    if opener is None:
        raise TypeError(f"Unknown standard stream '{name}'")
    return opener()

def get_text_stream(name: t.Literal['stdin', 'stdout', 'stderr'], encoding: str | None=None, errors: str | None='strict') -> t.TextIO:
    opener = text_streams.get(name)
    if opener is None:
        raise TypeError(f"Unknown standard stream '{name}'")
    return opener(encoding, errors)

def open_file(filename: str | os.PathLike[str], mode: str='r', encoding: str | None=None, errors: str | None='strict', lazy: bool=False, atomic: bool=False) -> t.IO[t.Any]:
    if lazy:
        return t.cast('t.IO[t.Any]', LazyFile(filename, mode, encoding, errors, atomic=atomic))
    f, should_close = open_stream(filename, mode, encoding, errors, atomic=atomic)
    if not should_close:
        f = t.cast('t.IO[t.Any]', KeepOpenFile(f))
    return f

def format_filename(filename: str | bytes | os.PathLike[str] | os.PathLike[bytes], shorten: bool=False) -> str:
    if shorten:
        filename = os.path.basename(filename)
    else:
        filename = os.fspath(filename)
    if isinstance(filename, bytes):
        filename = filename.decode(sys.getfilesystemencoding(), 'replace')
    else:
        filename = filename.encode('utf-8', 'surrogateescape').decode('utf-8', 'replace')
    return filename

def get_app_dir(app_name: str, roaming: bool=True, force_posix: bool=False) -> str:
    if WIN:
        key = 'APPDATA' if roaming else 'LOCALAPPDATA'
        folder = os.environ.get(key)
        if folder is None:
            folder = os.path.expanduser('~')
        return os.path.join(folder, app_name)
    if force_posix:
        return os.path.join(os.path.expanduser(f'~/.{_posixify(app_name)}'))
    if sys.platform == 'darwin':
        return os.path.join(os.path.expanduser('~/Library/Application Support'), app_name)
    return os.path.join(os.environ.get('XDG_CONFIG_HOME', os.path.expanduser('~/.config')), _posixify(app_name))

class PacifyFlushWrapper:

    def __init__(self, wrapped: t.IO[t.Any]) -> None:
        self.wrapped = wrapped

    def flush(self) -> None:
        try:
            self.wrapped.flush()
        except OSError as e:
            import errno
            if e.errno != errno.EPIPE:
                raise

    def __getattr__(self, attr: str) -> t.Any:
        return getattr(self.wrapped, attr)

def _detect_program_name(path: str | None=None, _main: ModuleType | None=None) -> str:
    if _main is None:
        _main = sys.modules['__main__']
    if not path:
        path = sys.argv[0]
    if getattr(_main, '__package__', None) in {None, ''} or (os.name == 'nt' and _main.__package__ == '' and (not os.path.exists(path)) and os.path.exists(f'{path}.exe')):
        return os.path.basename(path)
    py_module = t.cast(str, _main.__package__)
    name = os.path.splitext(os.path.basename(path))[0]
    if name != '__main__':
        py_module = f'{py_module}.{name}'
    return f"python -m {py_module.lstrip('.')}"

def _expand_args(args: cabc.Iterable[str], *, user: bool=True, env: bool=True, glob_recursive: bool=True) -> list[str]:
    from glob import glob
    out = []
    for arg in args:
        if user:
            arg = os.path.expanduser(arg)
        if env:
            arg = os.path.expandvars(arg)
        try:
            matches = glob(arg, recursive=glob_recursive)
        except re.error:
            matches = []
        if not matches:
            out.append(arg)
        else:
            out.extend(matches)
    return out
