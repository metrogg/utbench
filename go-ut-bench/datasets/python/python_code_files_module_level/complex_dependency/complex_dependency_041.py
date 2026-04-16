from __future__ import annotations
import dataclasses
from pathlib import Path
import shlex
import subprocess
from typing import Final
from typing import final
from typing import TYPE_CHECKING
from iniconfig import SectionWrapper
from _pytest.cacheprovider import Cache
from _pytest.compat import LEGACY_PATH
from _pytest.compat import legacy_path
from _pytest.config import Config
from _pytest.config import hookimpl
from _pytest.config import PytestPluginManager
from _pytest.deprecated import check_ispytest
from _pytest.fixtures import fixture
from _pytest.fixtures import FixtureRequest
from _pytest.main import Session
from _pytest.monkeypatch import MonkeyPatch
from _pytest.nodes import Collector
from _pytest.nodes import Item
from _pytest.nodes import Node
from _pytest.pytester import HookRecorder
from _pytest.pytester import Pytester
from _pytest.pytester import RunResult
from _pytest.terminal import TerminalReporter
from _pytest.tmpdir import TempPathFactory
if TYPE_CHECKING:
    import pexpect

@final
class Testdir:
    __test__ = False
    CLOSE_STDIN: Final = Pytester.CLOSE_STDIN
    TimeoutExpired: Final = Pytester.TimeoutExpired

    def __init__(self, pytester: Pytester, *, _ispytest: bool=False) -> None:
        check_ispytest(_ispytest)
        self._pytester = pytester

    @property
    def tmpdir(self) -> LEGACY_PATH:
        return legacy_path(self._pytester.path)

    @property
    def test_tmproot(self) -> LEGACY_PATH:
        return legacy_path(self._pytester._test_tmproot)

    @property
    def request(self):
        return self._pytester._request

    @property
    def plugins(self):
        return self._pytester.plugins

    @plugins.setter
    def plugins(self, plugins):
        self._pytester.plugins = plugins

    @property
    def monkeypatch(self) -> MonkeyPatch:
        return self._pytester._monkeypatch

    def make_hook_recorder(self, pluginmanager) -> HookRecorder:
        return self._pytester.make_hook_recorder(pluginmanager)

    def chdir(self) -> None:
        return self._pytester.chdir()

    def finalize(self) -> None:
        return self._pytester._finalize()

    def makefile(self, ext, *args, **kwargs) -> LEGACY_PATH:
        if ext and (not ext.startswith('.')):
            ext = '.' + ext
        return legacy_path(self._pytester.makefile(ext, *args, **kwargs))

    def makeconftest(self, source) -> LEGACY_PATH:
        return legacy_path(self._pytester.makeconftest(source))

    def makeini(self, source) -> LEGACY_PATH:
        return legacy_path(self._pytester.makeini(source))

    def getinicfg(self, source: str) -> SectionWrapper:
        return self._pytester.getinicfg(source)

    def makepyprojecttoml(self, source) -> LEGACY_PATH:
        return legacy_path(self._pytester.makepyprojecttoml(source))

    def makepyfile(self, *args, **kwargs) -> LEGACY_PATH:
        return legacy_path(self._pytester.makepyfile(*args, **kwargs))

    def maketxtfile(self, *args, **kwargs) -> LEGACY_PATH:
        return legacy_path(self._pytester.maketxtfile(*args, **kwargs))

    def syspathinsert(self, path=None) -> None:
        return self._pytester.syspathinsert(path)

    def mkdir(self, name) -> LEGACY_PATH:
        return legacy_path(self._pytester.mkdir(name))

    def mkpydir(self, name) -> LEGACY_PATH:
        return legacy_path(self._pytester.mkpydir(name))

    def copy_example(self, name=None) -> LEGACY_PATH:
        return legacy_path(self._pytester.copy_example(name))

    def getnode(self, config: Config, arg) -> Item | Collector | None:
        return self._pytester.getnode(config, arg)

    def getpathnode(self, path):
        return self._pytester.getpathnode(path)

    def genitems(self, colitems: list[Item | Collector]) -> list[Item]:
        return self._pytester.genitems(colitems)

    def runitem(self, source):
        return self._pytester.runitem(source)

    def inline_runsource(self, source, *cmdlineargs):
        return self._pytester.inline_runsource(source, *cmdlineargs)

    def inline_genitems(self, *args):
        return self._pytester.inline_genitems(*args)

    def inline_run(self, *args, plugins=(), no_reraise_ctrlc: bool=False):
        return self._pytester.inline_run(*args, plugins=plugins, no_reraise_ctrlc=no_reraise_ctrlc)

    def runpytest_inprocess(self, *args, **kwargs) -> RunResult:
        return self._pytester.runpytest_inprocess(*args, **kwargs)

    def runpytest(self, *args, **kwargs) -> RunResult:
        return self._pytester.runpytest(*args, **kwargs)

    def parseconfig(self, *args) -> Config:
        return self._pytester.parseconfig(*args)

    def parseconfigure(self, *args) -> Config:
        return self._pytester.parseconfigure(*args)

    def getitem(self, source, funcname='test_func'):
        return self._pytester.getitem(source, funcname)

    def getitems(self, source):
        return self._pytester.getitems(source)

    def getmodulecol(self, source, configargs=(), withinit=False):
        return self._pytester.getmodulecol(source, configargs=configargs, withinit=withinit)

    def collect_by_name(self, modcol: Collector, name: str) -> Item | Collector | None:
        return self._pytester.collect_by_name(modcol, name)

    def popen(self, cmdargs, stdout=subprocess.PIPE, stderr=subprocess.PIPE, stdin=CLOSE_STDIN, **kw):
        return self._pytester.popen(cmdargs, stdout, stderr, stdin, **kw)

    def run(self, *cmdargs, timeout=None, stdin=CLOSE_STDIN) -> RunResult:
        return self._pytester.run(*cmdargs, timeout=timeout, stdin=stdin)

    def runpython(self, script) -> RunResult:
        return self._pytester.runpython(script)

    def runpython_c(self, command):
        return self._pytester.runpython_c(command)

    def runpytest_subprocess(self, *args, timeout=None) -> RunResult:
        return self._pytester.runpytest_subprocess(*args, timeout=timeout)

    def spawn_pytest(self, string: str, expect_timeout: float=10.0) -> pexpect.spawn:
        return self._pytester.spawn_pytest(string, expect_timeout=expect_timeout)

    def spawn(self, cmd: str, expect_timeout: float=10.0) -> pexpect.spawn:
        return self._pytester.spawn(cmd, expect_timeout=expect_timeout)

    def __repr__(self) -> str:
        return f'<Testdir {self.tmpdir!r}>'

    def __str__(self) -> str:
        return str(self.tmpdir)

class LegacyTestdirPlugin:

    @staticmethod
    @fixture
    def testdir(pytester: Pytester) -> Testdir:
        return Testdir(pytester, _ispytest=True)

@final
@dataclasses.dataclass
class TempdirFactory:
    _tmppath_factory: TempPathFactory

    def __init__(self, tmppath_factory: TempPathFactory, *, _ispytest: bool=False) -> None:
        check_ispytest(_ispytest)
        self._tmppath_factory = tmppath_factory

    def mktemp(self, basename: str, numbered: bool=True) -> LEGACY_PATH:
        return legacy_path(self._tmppath_factory.mktemp(basename, numbered).resolve())

    def getbasetemp(self) -> LEGACY_PATH:
        return legacy_path(self._tmppath_factory.getbasetemp().resolve())

class LegacyTmpdirPlugin:

    @staticmethod
    @fixture(scope='session')
    def tmpdir_factory(request: FixtureRequest) -> TempdirFactory:
        return request.config._tmpdirhandler

    @staticmethod
    @fixture
    def tmpdir(tmp_path: Path) -> LEGACY_PATH:
        return legacy_path(tmp_path)

def Cache_makedir(self: Cache, name: str) -> LEGACY_PATH:
    return legacy_path(self.mkdir(name))

def FixtureRequest_fspath(self: FixtureRequest) -> LEGACY_PATH:
    return legacy_path(self.path)

def TerminalReporter_startdir(self: TerminalReporter) -> LEGACY_PATH:
    return legacy_path(self.startpath)

def Config_invocation_dir(self: Config) -> LEGACY_PATH:
    return legacy_path(str(self.invocation_params.dir))

def Config_rootdir(self: Config) -> LEGACY_PATH:
    return legacy_path(str(self.rootpath))

def Config_inifile(self: Config) -> LEGACY_PATH | None:
    return legacy_path(str(self.inipath)) if self.inipath else None

def Session_startdir(self: Session) -> LEGACY_PATH:
    return legacy_path(self.startpath)

def Config__getini_unknown_type(self, name: str, type: str, value: str | list[str]):
    if type == 'pathlist':
        assert self.inipath is not None
        dp = self.inipath.parent
        input_values = shlex.split(value) if isinstance(value, str) else value
        return [legacy_path(str(dp / x)) for x in input_values]
    else:
        raise ValueError(f'unknown configuration type: {type}', value)

def Node_fspath(self: Node) -> LEGACY_PATH:
    return legacy_path(self.path)

def Node_fspath_set(self: Node, value: LEGACY_PATH) -> None:
    self.path = Path(value)

@hookimpl(tryfirst=True)
def pytest_load_initial_conftests(early_config: Config) -> None:
    mp = MonkeyPatch()
    early_config.add_cleanup(mp.undo)
    mp.setattr(Cache, 'makedir', Cache_makedir, raising=False)
    mp.setattr(FixtureRequest, 'fspath', property(FixtureRequest_fspath), raising=False)
    mp.setattr(TerminalReporter, 'startdir', property(TerminalReporter_startdir), raising=False)
    mp.setattr(Config, 'invocation_dir', property(Config_invocation_dir), raising=False)
    mp.setattr(Config, 'rootdir', property(Config_rootdir), raising=False)
    mp.setattr(Config, 'inifile', property(Config_inifile), raising=False)
    mp.setattr(Session, 'startdir', property(Session_startdir), raising=False)
    mp.setattr(Config, '_getini_unknown_type', Config__getini_unknown_type)
    mp.setattr(Node, 'fspath', property(Node_fspath, Node_fspath_set), raising=False)

@hookimpl
def pytest_configure(config: Config) -> None:
    if config.pluginmanager.has_plugin('tmpdir'):
        mp = MonkeyPatch()
        config.add_cleanup(mp.undo)
        try:
            tmp_path_factory = config._tmp_path_factory
        except AttributeError:
            pass
        else:
            _tmpdirhandler = TempdirFactory(tmp_path_factory, _ispytest=True)
            mp.setattr(config, '_tmpdirhandler', _tmpdirhandler, raising=False)
        config.pluginmanager.register(LegacyTmpdirPlugin, 'legacypath-tmpdir')

@hookimpl
def pytest_plugin_registered(plugin: object, manager: PytestPluginManager) -> None:
    is_pytester = plugin is manager.get_plugin('pytester')
    if is_pytester and (not manager.is_registered(LegacyTestdirPlugin)):
        manager.register(LegacyTestdirPlugin, 'legacypath-pytester')
