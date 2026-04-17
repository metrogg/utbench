from __future__ import annotations
import argparse
from collections.abc import Callable
from collections.abc import Generator
import functools
import sys
import types
from typing import Any
import unittest
from _pytest import outcomes
from _pytest._code import ExceptionInfo
from _pytest.capture import CaptureManager
from _pytest.config import Config
from _pytest.config import ConftestImportFailure
from _pytest.config import hookimpl
from _pytest.config import PytestPluginManager
from _pytest.config.argparsing import Parser
from _pytest.config.exceptions import UsageError
from _pytest.nodes import Node
from _pytest.reports import BaseReport
from _pytest.runner import CallInfo

def _validate_usepdb_cls(value: str) -> tuple[str, str]:
    try:
        modname, classname = value.split(':')
    except ValueError as e:
        raise argparse.ArgumentTypeError(f"{value!r} is not in the format 'modname:classname'") from e
    return (modname, classname)

def pytest_addoption(parser: Parser) -> None:
    group = parser.getgroup('general')
    group.addoption('--pdb', dest='usepdb', action='store_true', help='Start the interactive Python debugger on errors or KeyboardInterrupt')
    group.addoption('--pdbcls', dest='usepdb_cls', metavar='modulename:classname', type=_validate_usepdb_cls, help='Specify a custom interactive Python debugger for use with --pdb.For example: --pdbcls=IPython.terminal.debugger:TerminalPdb')
    group.addoption('--trace', dest='trace', action='store_true', help='Immediately break when running each test')

def pytest_configure(config: Config) -> None:
    import pdb
    if config.getvalue('trace'):
        config.pluginmanager.register(PdbTrace(), 'pdbtrace')
    if config.getvalue('usepdb'):
        config.pluginmanager.register(PdbInvoke(), 'pdbinvoke')
    pytestPDB._saved.append((pdb.set_trace, pytestPDB._pluginmanager, pytestPDB._config))
    pdb.set_trace = pytestPDB.set_trace
    pytestPDB._pluginmanager = config.pluginmanager
    pytestPDB._config = config

    def fin() -> None:
        pdb.set_trace, pytestPDB._pluginmanager, pytestPDB._config = pytestPDB._saved.pop()
    config.add_cleanup(fin)

class pytestPDB:
    _pluginmanager: PytestPluginManager | None = None
    _config: Config | None = None
    _saved: list[tuple[Callable[..., None], PytestPluginManager | None, Config | None]] = []
    _recursive_debug = 0
    _wrapped_pdb_cls: tuple[type[Any], type[Any]] | None = None

    @classmethod
    def _is_capturing(cls, capman: CaptureManager | None) -> str | bool:
        if capman:
            return capman.is_capturing()
        return False

    @classmethod
    def _import_pdb_cls(cls, capman: CaptureManager | None):
        if not cls._config:
            import pdb
            return pdb.Pdb
        usepdb_cls = cls._config.getvalue('usepdb_cls')
        if cls._wrapped_pdb_cls and cls._wrapped_pdb_cls[0] == usepdb_cls:
            return cls._wrapped_pdb_cls[1]
        if usepdb_cls:
            modname, classname = usepdb_cls
            try:
                __import__(modname)
                mod = sys.modules[modname]
                parts = classname.split('.')
                pdb_cls = getattr(mod, parts[0])
                for part in parts[1:]:
                    pdb_cls = getattr(pdb_cls, part)
            except Exception as exc:
                value = ':'.join((modname, classname))
                raise UsageError(f'--pdbcls: could not import {value!r}: {exc}') from exc
        else:
            import pdb
            pdb_cls = pdb.Pdb
        wrapped_cls = cls._get_pdb_wrapper_class(pdb_cls, capman)
        cls._wrapped_pdb_cls = (usepdb_cls, wrapped_cls)
        return wrapped_cls

    @classmethod
    def _get_pdb_wrapper_class(cls, pdb_cls, capman: CaptureManager | None):
        import _pytest.config

        class PytestPdbWrapper(pdb_cls):
            _pytest_capman = capman
            _continued = False

            def do_debug(self, arg):
                cls._recursive_debug += 1
                ret = super().do_debug(arg)
                cls._recursive_debug -= 1
                return ret
            if hasattr(pdb_cls, 'do_debug'):
                do_debug.__doc__ = pdb_cls.do_debug.__doc__

            def do_continue(self, arg):
                ret = super().do_continue(arg)
                if cls._recursive_debug == 0:
                    assert cls._config is not None
                    tw = _pytest.config.create_terminal_writer(cls._config)
                    tw.line()
                    capman = self._pytest_capman
                    capturing = pytestPDB._is_capturing(capman)
                    if capturing:
                        if capturing == 'global':
                            tw.sep('>', 'PDB continue (IO-capturing resumed)')
                        else:
                            tw.sep('>', f'PDB continue (IO-capturing resumed for {capturing})')
                        assert capman is not None
                        capman.resume()
                    else:
                        tw.sep('>', 'PDB continue')
                assert cls._pluginmanager is not None
                cls._pluginmanager.hook.pytest_leave_pdb(config=cls._config, pdb=self)
                self._continued = True
                return ret
            if hasattr(pdb_cls, 'do_continue'):
                do_continue.__doc__ = pdb_cls.do_continue.__doc__
            do_c = do_cont = do_continue

            def do_quit(self, arg):
                ret = super().do_quit(arg)
                if cls._recursive_debug == 0:
                    outcomes.exit('Quitting debugger')
                return ret
            if hasattr(pdb_cls, 'do_quit'):
                do_quit.__doc__ = pdb_cls.do_quit.__doc__
            do_q = do_quit
            do_exit = do_quit

            def setup(self, f, tb):
                ret = super().setup(f, tb)
                if not ret and self._continued:
                    if self._pytest_capman:
                        self._pytest_capman.suspend_global_capture(in_=True)
                return ret

            def get_stack(self, f, t):
                stack, i = super().get_stack(f, t)
                if f is None:
                    i = max(0, len(stack) - 1)
                    while i and stack[i][0].f_locals.get('__tracebackhide__', False):
                        i -= 1
                return (stack, i)
        return PytestPdbWrapper

    @classmethod
    def _init_pdb(cls, method, *args, **kwargs):
        import _pytest.config
        if cls._pluginmanager is None:
            capman: CaptureManager | None = None
        else:
            capman = cls._pluginmanager.getplugin('capturemanager')
        if capman:
            capman.suspend(in_=True)
        if cls._config:
            tw = _pytest.config.create_terminal_writer(cls._config)
            tw.line()
            if cls._recursive_debug == 0:
                header = kwargs.pop('header', None)
                if header is not None:
                    tw.sep('>', header)
                else:
                    capturing = cls._is_capturing(capman)
                    if capturing == 'global':
                        tw.sep('>', f'PDB {method} (IO-capturing turned off)')
                    elif capturing:
                        tw.sep('>', f'PDB {method} (IO-capturing turned off for {capturing})')
                    else:
                        tw.sep('>', f'PDB {method}')
        _pdb = cls._import_pdb_cls(capman)(**kwargs)
        if cls._pluginmanager:
            cls._pluginmanager.hook.pytest_enter_pdb(config=cls._config, pdb=_pdb)
        return _pdb

    @classmethod
    def set_trace(cls, *args, **kwargs) -> None:
        frame = sys._getframe().f_back
        _pdb = cls._init_pdb('set_trace', *args, **kwargs)
        _pdb.set_trace(frame)

class PdbInvoke:

    def pytest_exception_interact(self, node: Node, call: CallInfo[Any], report: BaseReport) -> None:
        capman = node.config.pluginmanager.getplugin('capturemanager')
        if capman:
            capman.suspend_global_capture(in_=True)
            out, err = capman.read_global_capture()
            sys.stdout.write(out)
            sys.stdout.write(err)
        assert call.excinfo is not None
        if not isinstance(call.excinfo.value, unittest.SkipTest):
            _enter_pdb(node, call.excinfo, report)

    def pytest_internalerror(self, excinfo: ExceptionInfo[BaseException]) -> None:
        exc_or_tb = _postmortem_exc_or_tb(excinfo)
        post_mortem(exc_or_tb)

class PdbTrace:

    @hookimpl(wrapper=True)
    def pytest_pyfunc_call(self, pyfuncitem) -> Generator[None, object, object]:
        wrap_pytest_function_for_tracing(pyfuncitem)
        return (yield)

def wrap_pytest_function_for_tracing(pyfuncitem) -> None:
    _pdb = pytestPDB._init_pdb('runcall')
    testfunction = pyfuncitem.obj

    @functools.wraps(testfunction)
    def wrapper(*args, **kwargs) -> None:
        func = functools.partial(testfunction, *args, **kwargs)
        _pdb.runcall(func)
    pyfuncitem.obj = wrapper

def maybe_wrap_pytest_function_for_tracing(pyfuncitem) -> None:
    if pyfuncitem.config.getvalue('trace'):
        wrap_pytest_function_for_tracing(pyfuncitem)

def _enter_pdb(node: Node, excinfo: ExceptionInfo[BaseException], rep: BaseReport) -> BaseReport:
    tw = node.config.pluginmanager.getplugin('terminalreporter')._tw
    tw.line()
    showcapture = node.config.option.showcapture
    for sectionname, content in (('stdout', rep.capstdout), ('stderr', rep.capstderr), ('log', rep.caplog)):
        if showcapture in (sectionname, 'all') and content:
            tw.sep('>', 'captured ' + sectionname)
            if content[-1:] == '\n':
                content = content[:-1]
            tw.line(content)
    tw.sep('>', 'traceback')
    rep.toterminal(tw)
    tw.sep('>', 'entering PDB')
    tb_or_exc = _postmortem_exc_or_tb(excinfo)
    rep._pdbshown = True
    post_mortem(tb_or_exc)
    return rep

def _postmortem_exc_or_tb(excinfo: ExceptionInfo[BaseException]) -> types.TracebackType | BaseException:
    from doctest import UnexpectedException
    get_exc = sys.version_info >= (3, 13)
    if isinstance(excinfo.value, UnexpectedException):
        underlying_exc = excinfo.value
        if get_exc:
            return underlying_exc.exc_info[1]
        return underlying_exc.exc_info[2]
    elif isinstance(excinfo.value, ConftestImportFailure):
        cause = excinfo.value.cause
        if get_exc:
            return cause
        assert cause.__traceback__ is not None
        return cause.__traceback__
    else:
        assert excinfo._excinfo is not None
        if get_exc:
            return excinfo._excinfo[1]
        return excinfo._excinfo[2]

def post_mortem(tb_or_exc: types.TracebackType | BaseException) -> None:
    p = pytestPDB._init_pdb('post_mortem')
    p.reset()
    p.interaction(None, tb_or_exc)
    if p.quitting:
        outcomes.exit('Quitting debugger')
