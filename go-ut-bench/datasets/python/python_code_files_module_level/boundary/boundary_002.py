from __future__ import annotations
import argparse
from collections.abc import Callable
from collections.abc import Mapping
from collections.abc import Sequence
import os
import sys
from typing import Any
from typing import final
from typing import Literal
from typing import NoReturn
from .exceptions import UsageError
import _pytest._io
from _pytest.deprecated import check_ispytest
FILE_OR_DIR = 'file_or_dir'

class NotSet:

    def __repr__(self) -> str:
        return '<notset>'
NOT_SET = NotSet()

@final
class Parser:

    def __init__(self, usage: str | None=None, processopt: Callable[[Argument], None] | None=None, *, _ispytest: bool=False) -> None:
        check_ispytest(_ispytest)
        from _pytest._argcomplete import filescompleter
        self._processopt = processopt
        self.extra_info: dict[str, Any] = {}
        self.optparser = PytestArgumentParser(self, usage, self.extra_info)
        anonymous_arggroup = self.optparser.add_argument_group('Custom options')
        self._anonymous = OptionGroup(anonymous_arggroup, '_anonymous', self, _ispytest=True)
        self._groups = [self._anonymous]
        file_or_dir_arg = self.optparser.add_argument(FILE_OR_DIR, nargs='*')
        file_or_dir_arg.completer = filescompleter
        self._inidict: dict[str, tuple[str, str, Any]] = {}
        self._ini_aliases: dict[str, str] = {}

    @property
    def prog(self) -> str:
        return self.optparser.prog

    @prog.setter
    def prog(self, value: str) -> None:
        self.optparser.prog = value

    def processoption(self, option: Argument) -> None:
        if self._processopt:
            if option.dest:
                self._processopt(option)

    def getgroup(self, name: str, description: str='', after: str | None=None) -> OptionGroup:
        for group in self._groups:
            if group.name == name:
                return group
        arggroup = self.optparser.add_argument_group(description or name)
        group = OptionGroup(arggroup, name, self, _ispytest=True)
        i = 0
        for i, grp in enumerate(self._groups):
            if grp.name == after:
                break
        self._groups.insert(i + 1, group)
        self.optparser._action_groups.insert(i + 1, self.optparser._action_groups.pop())
        return group

    def addoption(self, *opts: str, **attrs: Any) -> None:
        self._anonymous.addoption(*opts, **attrs)

    def parse(self, args: Sequence[str | os.PathLike[str]], namespace: argparse.Namespace | None=None) -> argparse.Namespace:
        from _pytest._argcomplete import try_argcomplete
        try_argcomplete(self.optparser)
        strargs = [os.fspath(x) for x in args]
        if namespace is None:
            namespace = argparse.Namespace()
        try:
            namespace._raise_print_help = True
            return self.optparser.parse_intermixed_args(strargs, namespace=namespace)
        finally:
            del namespace._raise_print_help

    def parse_known_args(self, args: Sequence[str | os.PathLike[str]], namespace: argparse.Namespace | None=None) -> argparse.Namespace:
        return self.parse_known_and_unknown_args(args, namespace=namespace)[0]

    def parse_known_and_unknown_args(self, args: Sequence[str | os.PathLike[str]], namespace: argparse.Namespace | None=None) -> tuple[argparse.Namespace, list[str]]:
        strargs = [os.fspath(x) for x in args]
        if sys.version_info < (3, 12, 8) or (3, 13) <= sys.version_info < (3, 13, 1):
            namespace, unknown = self.optparser.parse_known_args(strargs, namespace)
            assert namespace is not None
            file_or_dir = getattr(namespace, FILE_OR_DIR)
            unknown_flags: list[str] = []
            for arg in unknown:
                (unknown_flags if arg.startswith('-') else file_or_dir).append(arg)
            return (namespace, unknown_flags)
        else:
            return self.optparser.parse_known_intermixed_args(strargs, namespace)

    def addini(self, name: str, help: str, type: Literal['string', 'paths', 'pathlist', 'args', 'linelist', 'bool', 'int', 'float'] | None=None, default: Any=NOT_SET, *, aliases: Sequence[str]=()) -> None:
        assert type in (None, 'string', 'paths', 'pathlist', 'args', 'linelist', 'bool', 'int', 'float')
        if type is None:
            type = 'string'
        if default is NOT_SET:
            default = get_ini_default_for_type(type)
        self._inidict[name] = (help, type, default)
        for alias in aliases:
            if alias in self._inidict:
                raise ValueError(f'alias {alias!r} conflicts with existing configuration option')
            if (already := self._ini_aliases.get(alias)) is not None:
                raise ValueError(f'{alias!r} is already an alias of {already!r}')
            self._ini_aliases[alias] = name

def get_ini_default_for_type(type: Literal['string', 'paths', 'pathlist', 'args', 'linelist', 'bool', 'int', 'float']) -> Any:
    if type in ('paths', 'pathlist', 'args', 'linelist'):
        return []
    elif type == 'bool':
        return False
    elif type == 'int':
        return 0
    elif type == 'float':
        return 0.0
    else:
        return ''

class ArgumentError(Exception):

    def __init__(self, msg: str, option: Argument | str) -> None:
        self.msg = msg
        self.option_id = str(option)

    def __str__(self) -> str:
        if self.option_id:
            return f'option {self.option_id}: {self.msg}'
        else:
            return self.msg

class Argument:

    def __init__(self, *names: str, **attrs: Any) -> None:
        self._attrs = attrs
        self._short_opts: list[str] = []
        self._long_opts: list[str] = []
        try:
            self.type = attrs['type']
        except KeyError:
            pass
        try:
            self.default = attrs['default']
        except KeyError:
            pass
        self._set_opt_strings(names)
        dest: str | None = attrs.get('dest')
        if dest:
            self.dest = dest
        elif self._long_opts:
            self.dest = self._long_opts[0][2:].replace('-', '_')
        else:
            try:
                self.dest = self._short_opts[0][1:]
            except IndexError as e:
                self.dest = '???'
                raise ArgumentError('need a long or short option', self) from e

    def names(self) -> list[str]:
        return self._short_opts + self._long_opts

    def attrs(self) -> Mapping[str, Any]:
        for attr in ('default', 'dest', 'help', self.dest):
            try:
                self._attrs[attr] = getattr(self, attr)
            except AttributeError:
                pass
        return self._attrs

    def _set_opt_strings(self, opts: Sequence[str]) -> None:
        for opt in opts:
            if len(opt) < 2:
                raise ArgumentError(f'invalid option string {opt!r}: must be at least two characters long', self)
            elif len(opt) == 2:
                if not (opt[0] == '-' and opt[1] != '-'):
                    raise ArgumentError(f'invalid short option string {opt!r}: must be of the form -x, (x any non-dash char)', self)
                self._short_opts.append(opt)
            else:
                if not (opt[0:2] == '--' and opt[2] != '-'):
                    raise ArgumentError(f'invalid long option string {opt!r}: must start with --, followed by non-dash', self)
                self._long_opts.append(opt)

    def __repr__(self) -> str:
        args: list[str] = []
        if self._short_opts:
            args += ['_short_opts: ' + repr(self._short_opts)]
        if self._long_opts:
            args += ['_long_opts: ' + repr(self._long_opts)]
        args += ['dest: ' + repr(self.dest)]
        if hasattr(self, 'type'):
            args += ['type: ' + repr(self.type)]
        if hasattr(self, 'default'):
            args += ['default: ' + repr(self.default)]
        return 'Argument({})'.format(', '.join(args))

class OptionGroup:

    def __init__(self, arggroup: argparse._ArgumentGroup, name: str, parser: Parser | None, _ispytest: bool=False) -> None:
        check_ispytest(_ispytest)
        self._arggroup = arggroup
        self.name = name
        self.options: list[Argument] = []
        self.parser = parser

    def addoption(self, *opts: str, **attrs: Any) -> None:
        conflict = set(opts).intersection((name for opt in self.options for name in opt.names()))
        if conflict:
            raise ValueError(f'option names {conflict} already added')
        option = Argument(*opts, **attrs)
        self._addoption_instance(option, shortupper=False)

    def _addoption(self, *opts: str, **attrs: Any) -> None:
        option = Argument(*opts, **attrs)
        self._addoption_instance(option, shortupper=True)

    def _addoption_instance(self, option: Argument, shortupper: bool=False) -> None:
        if not shortupper:
            for opt in option._short_opts:
                if opt[0] == '-' and opt[1].islower():
                    raise ValueError('lowercase shortoptions reserved')
        if self.parser:
            self.parser.processoption(option)
        self._arggroup.add_argument(*option.names(), **option.attrs())
        self.options.append(option)

class PytestArgumentParser(argparse.ArgumentParser):

    def __init__(self, parser: Parser, usage: str | None, extra_info: dict[str, str]) -> None:
        self._parser = parser
        super().__init__(usage=usage, add_help=False, formatter_class=DropShorterLongHelpFormatter, allow_abbrev=False, fromfile_prefix_chars='@')
        self.extra_info = extra_info

    def error(self, message: str) -> NoReturn:
        msg = f'{self.prog}: error: {message}'
        if self.extra_info:
            msg += '\n' + '\n'.join((f'  {k}: {v}' for k, v in sorted(self.extra_info.items())))
        raise UsageError(self.format_usage() + msg)

class DropShorterLongHelpFormatter(argparse.HelpFormatter):

    def __init__(self, *args: Any, **kwargs: Any) -> None:
        if 'width' not in kwargs:
            kwargs['width'] = _pytest._io.get_terminal_width()
        super().__init__(*args, **kwargs)

    def _format_action_invocation(self, action: argparse.Action) -> str:
        orgstr = super()._format_action_invocation(action)
        if orgstr and orgstr[0] != '-':
            return orgstr
        res: str | None = getattr(action, '_formatted_action_invocation', None)
        if res:
            return res
        options = orgstr.split(', ')
        if len(options) == 2 and (len(options[0]) == 2 or len(options[1]) == 2):
            action._formatted_action_invocation = orgstr
            return orgstr
        return_list = []
        short_long: dict[str, str] = {}
        for option in options:
            if len(option) == 2 or option[2] == ' ':
                continue
            if not option.startswith('--'):
                raise ArgumentError(f'long optional argument without "--": [{option}]', option)
            xxoption = option[2:]
            shortened = xxoption.replace('-', '')
            if shortened not in short_long or len(short_long[shortened]) < len(xxoption):
                short_long[shortened] = xxoption
        for option in options:
            if len(option) == 2 or option[2] == ' ':
                return_list.append(option)
            if option[2:] == short_long.get(option.replace('-', '')):
                return_list.append(option.replace(' ', '=', 1))
        formatted_action_invocation = ', '.join(return_list)
        action._formatted_action_invocation = formatted_action_invocation
        return formatted_action_invocation

    def _split_lines(self, text, width):
        import textwrap
        lines = []
        for line in text.splitlines():
            lines.extend(textwrap.wrap(line.strip(), width))
        return lines

class OverrideIniAction(argparse.Action):

    def __init__(self, option_strings: Sequence[str], dest: str, nargs: int | str | None=None, *args, ini_option: str, ini_value: str, **kwargs) -> None:
        super().__init__(option_strings, dest, 0, *args, **kwargs)
        self.ini_option = ini_option
        self.ini_value = ini_value

    def __call__(self, parser: argparse.ArgumentParser, namespace: argparse.Namespace, *args, **kwargs) -> None:
        setattr(namespace, self.dest, True)
        current_overrides = getattr(namespace, 'override_ini', None)
        if current_overrides is None:
            current_overrides = []
        current_overrides.append(f'{self.ini_option}={self.ini_value}')
        setattr(namespace, 'override_ini', current_overrides)
