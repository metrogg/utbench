from __future__ import annotations
import bdb
from collections.abc import Callable
import dataclasses
import os
import sys
import types
from typing import cast
from typing import final
from typing import Generic
from typing import Literal
from typing import TYPE_CHECKING
from typing import TypeVar
from .config import Config
from .reports import BaseReport
from .reports import CollectErrorRepr
from .reports import CollectReport
from .reports import TestReport
from _pytest import timing
from _pytest._code.code import ExceptionChainRepr
from _pytest._code.code import ExceptionInfo
from _pytest._code.code import TerminalRepr
from _pytest.config.argparsing import Parser
from _pytest.deprecated import check_ispytest
from _pytest.nodes import Collector
from _pytest.nodes import Directory
from _pytest.nodes import Item
from _pytest.nodes import Node
from _pytest.outcomes import Exit
from _pytest.outcomes import OutcomeException
from _pytest.outcomes import Skipped
from _pytest.outcomes import TEST_OUTCOME
if sys.version_info < (3, 11):
    from exceptiongroup import BaseExceptionGroup
if TYPE_CHECKING:
    from _pytest.main import Session
    from _pytest.terminal import TerminalReporter

def pytest_addoption(parser: Parser) -> None:
    group = parser.getgroup('terminal reporting', 'Reporting', after='general')
    group.addoption('--durations', action='store', type=int, default=None, metavar='N', help='Show N slowest setup/test durations (N=0 for all)')
    group.addoption('--durations-min', action='store', type=float, default=None, metavar='N', help='Minimal duration in seconds for inclusion in slowest list. Default: 0.005 (or 0.0 if -vv is given).')

def pytest_terminal_summary(terminalreporter: TerminalReporter) -> None:
    durations = terminalreporter.config.option.durations
    durations_min = terminalreporter.config.option.durations_min
    verbose = terminalreporter.config.get_verbosity()
    if durations is None:
        return
    if durations_min is None:
        durations_min = 0.005 if verbose < 2 else 0.0
    tr = terminalreporter
    dlist = []
    for replist in tr.stats.values():
        for rep in replist:
            if hasattr(rep, 'duration'):
                dlist.append(rep)
    if not dlist:
        return
    dlist.sort(key=lambda x: x.duration, reverse=True)
    if not durations:
        tr.write_sep('=', 'slowest durations')
    else:
        tr.write_sep('=', f'slowest {durations} durations')
        dlist = dlist[:durations]
    for i, rep in enumerate(dlist):
        if rep.duration < durations_min:
            tr.write_line('')
            message = f'({len(dlist) - i} durations < {durations_min:g}s hidden.'
            if terminalreporter.config.option.durations_min is None:
                message += '  Use -vv to show these durations.'
            message += ')'
            tr.write_line(message)
            break
        tr.write_line(f'{rep.duration:02.2f}s {rep.when:<8} {rep.nodeid}')

def pytest_sessionstart(session: Session) -> None:
    session._setupstate = SetupState()

def pytest_sessionfinish(session: Session) -> None:
    session._setupstate.teardown_exact(None)

def pytest_runtest_protocol(item: Item, nextitem: Item | None) -> bool:
    ihook = item.ihook
    ihook.pytest_runtest_logstart(nodeid=item.nodeid, location=item.location)
    runtestprotocol(item, nextitem=nextitem)
    ihook.pytest_runtest_logfinish(nodeid=item.nodeid, location=item.location)
    return True

def runtestprotocol(item: Item, log: bool=True, nextitem: Item | None=None) -> list[TestReport]:
    hasrequest = hasattr(item, '_request')
    if hasrequest and (not item._request):
        item._initrequest()
    rep = call_and_report(item, 'setup', log)
    reports = [rep]
    if rep.passed:
        if item.config.getoption('setupshow', False):
            show_test_item(item)
        if not item.config.getoption('setuponly', False):
            reports.append(call_and_report(item, 'call', log))
    if item.session.shouldfail or item.session.shouldstop:
        nextitem = None
    reports.append(call_and_report(item, 'teardown', log, nextitem=nextitem))
    if hasrequest:
        item._request = False
        item.funcargs = None
    return reports

def show_test_item(item: Item) -> None:
    tw = item.config.get_terminal_writer()
    tw.line()
    tw.write(' ' * 8)
    tw.write(item.nodeid)
    used_fixtures = sorted(getattr(item, 'fixturenames', []))
    if used_fixtures:
        tw.write(' (fixtures used: {})'.format(', '.join(used_fixtures)))
    tw.flush()

def pytest_runtest_setup(item: Item) -> None:
    _update_current_test_var(item, 'setup')
    item.session._setupstate.setup(item)

def pytest_runtest_call(item: Item) -> None:
    _update_current_test_var(item, 'call')
    try:
        del sys.last_type
        del sys.last_value
        del sys.last_traceback
        if sys.version_info >= (3, 12, 0):
            del sys.last_exc
    except AttributeError:
        pass
    try:
        item.runtest()
    except Exception as e:
        sys.last_type = type(e)
        sys.last_value = e
        if sys.version_info >= (3, 12, 0):
            sys.last_exc = e
        assert e.__traceback__ is not None
        sys.last_traceback = e.__traceback__.tb_next
        raise

def pytest_runtest_teardown(item: Item, nextitem: Item | None) -> None:
    _update_current_test_var(item, 'teardown')
    item.session._setupstate.teardown_exact(nextitem)
    _update_current_test_var(item, None)

def _update_current_test_var(item: Item, when: Literal['setup', 'call', 'teardown'] | None) -> None:
    var_name = 'PYTEST_CURRENT_TEST'
    if when:
        value = f'{item.nodeid} ({when})'
        value = value.replace('\x00', '(null)')
        os.environ[var_name] = value
    else:
        os.environ.pop(var_name)

def pytest_report_teststatus(report: BaseReport) -> tuple[str, str, str] | None:
    if report.when in ('setup', 'teardown'):
        if report.failed:
            return ('error', 'E', 'ERROR')
        elif report.skipped:
            return ('skipped', 's', 'SKIPPED')
        else:
            return ('', '', '')
    return None

def call_and_report(item: Item, when: Literal['setup', 'call', 'teardown'], log: bool=True, **kwds) -> TestReport:
    ihook = item.ihook
    if when == 'setup':
        runtest_hook: Callable[..., None] = ihook.pytest_runtest_setup
    elif when == 'call':
        runtest_hook = ihook.pytest_runtest_call
    elif when == 'teardown':
        runtest_hook = ihook.pytest_runtest_teardown
    else:
        assert False, f'Unhandled runtest hook case: {when}'
    call = CallInfo.from_call(lambda: runtest_hook(item=item, **kwds), when=when, reraise=get_reraise_exceptions(item.config))
    report: TestReport = ihook.pytest_runtest_makereport(item=item, call=call)
    if log:
        ihook.pytest_runtest_logreport(report=report)
    if check_interactive_exception(call, report):
        ihook.pytest_exception_interact(node=item, call=call, report=report)
    return report

def get_reraise_exceptions(config: Config) -> tuple[type[BaseException], ...]:
    reraise: tuple[type[BaseException], ...] = (Exit,)
    if not config.getoption('usepdb', False):
        reraise += (KeyboardInterrupt,)
    return reraise

def check_interactive_exception(call: CallInfo[object], report: BaseReport) -> bool:
    if call.excinfo is None:
        return False
    if hasattr(report, 'wasxfail'):
        return False
    if isinstance(call.excinfo.value, Skipped | bdb.BdbQuit):
        return False
    return True
TResult = TypeVar('TResult', covariant=True)

@final
@dataclasses.dataclass
class CallInfo(Generic[TResult]):
    _result: TResult | None
    excinfo: ExceptionInfo[BaseException] | None
    start: float
    stop: float
    duration: float
    when: Literal['collect', 'setup', 'call', 'teardown']

    def __init__(self, result: TResult | None, excinfo: ExceptionInfo[BaseException] | None, start: float, stop: float, duration: float, when: Literal['collect', 'setup', 'call', 'teardown'], *, _ispytest: bool=False) -> None:
        check_ispytest(_ispytest)
        self._result = result
        self.excinfo = excinfo
        self.start = start
        self.stop = stop
        self.duration = duration
        self.when = when

    @property
    def result(self) -> TResult:
        if self.excinfo is not None:
            raise AttributeError(f'{self!r} has no valid result')
        return cast(TResult, self._result)

    @classmethod
    def from_call(cls, func: Callable[[], TResult], when: Literal['collect', 'setup', 'call', 'teardown'], reraise: type[BaseException] | tuple[type[BaseException], ...] | None=None) -> CallInfo[TResult]:
        excinfo = None
        instant = timing.Instant()
        try:
            result: TResult | None = func()
        except BaseException:
            excinfo = ExceptionInfo.from_current()
            if reraise is not None and isinstance(excinfo.value, reraise):
                raise
            result = None
        duration = instant.elapsed()
        return cls(start=duration.start.time, stop=duration.stop.time, duration=duration.seconds, when=when, result=result, excinfo=excinfo, _ispytest=True)

    def __repr__(self) -> str:
        if self.excinfo is None:
            return f'<CallInfo when={self.when!r} result: {self._result!r}>'
        return f'<CallInfo when={self.when!r} excinfo={self.excinfo!r}>'

def pytest_runtest_makereport(item: Item, call: CallInfo[None]) -> TestReport:
    return TestReport.from_item_and_call(item, call)

def pytest_make_collect_report(collector: Collector) -> CollectReport:

    def collect() -> list[Item | Collector]:
        if isinstance(collector, Directory):
            collector.config.pluginmanager._loadconftestmodules(collector.path, collector.config.getoption('importmode'), rootpath=collector.config.rootpath, consider_namespace_packages=collector.config.getini('consider_namespace_packages'))
        return list(collector.collect())
    call = CallInfo.from_call(collect, 'collect', reraise=(KeyboardInterrupt, SystemExit))
    longrepr: None | tuple[str, int, str] | str | TerminalRepr = None
    if not call.excinfo:
        outcome: Literal['passed', 'skipped', 'failed'] = 'passed'
    else:
        skip_exceptions = [Skipped]
        unittest = sys.modules.get('unittest')
        if unittest is not None:
            skip_exceptions.append(unittest.SkipTest)
        if isinstance(call.excinfo.value, tuple(skip_exceptions)):
            outcome = 'skipped'
            r_ = collector._repr_failure_py(call.excinfo, 'line')
            assert isinstance(r_, ExceptionChainRepr), repr(r_)
            r = r_.reprcrash
            assert r
            longrepr = (str(r.path), r.lineno, r.message)
        else:
            outcome = 'failed'
            errorinfo = collector.repr_failure(call.excinfo)
            if not hasattr(errorinfo, 'toterminal'):
                assert isinstance(errorinfo, str)
                errorinfo = CollectErrorRepr(errorinfo)
            longrepr = errorinfo
    result = call.result if not call.excinfo else None
    rep = CollectReport(collector.nodeid, outcome, longrepr, result)
    rep.call = call
    return rep

class SetupState:

    def __init__(self) -> None:
        self.stack: dict[Node, tuple[list[Callable[[], object]], tuple[OutcomeException | Exception, types.TracebackType | None] | None]] = {}

    def setup(self, item: Item) -> None:
        needed_collectors = item.listchain()
        for col, (finalizers, exc) in self.stack.items():
            assert col in needed_collectors, 'previous item was not torn down properly'
            if exc:
                raise exc[0].with_traceback(exc[1])
        for col in needed_collectors[len(self.stack):]:
            assert col not in self.stack
            self.stack[col] = ([col.teardown], None)
            try:
                col.setup()
            except TEST_OUTCOME as exc:
                self.stack[col] = (self.stack[col][0], (exc, exc.__traceback__))
                raise

    def addfinalizer(self, finalizer: Callable[[], object], node: Node) -> None:
        assert node and (not isinstance(node, tuple))
        assert callable(finalizer)
        assert node in self.stack, (node, self.stack)
        self.stack[node][0].append(finalizer)

    def teardown_exact(self, nextitem: Item | None) -> None:
        needed_collectors = nextitem and nextitem.listchain() or []
        exceptions: list[BaseException] = []
        while self.stack:
            if list(self.stack.keys()) == needed_collectors[:len(self.stack)]:
                break
            node, (finalizers, _) = self.stack.popitem()
            these_exceptions = []
            while finalizers:
                fin = finalizers.pop()
                try:
                    fin()
                except TEST_OUTCOME as e:
                    these_exceptions.append(e)
            if len(these_exceptions) == 1:
                exceptions.extend(these_exceptions)
            elif these_exceptions:
                msg = f'errors while tearing down {node!r}'
                exceptions.append(BaseExceptionGroup(msg, these_exceptions[::-1]))
        if len(exceptions) == 1:
            raise exceptions[0]
        elif exceptions:
            raise BaseExceptionGroup('errors during test teardown', exceptions[::-1])
        if nextitem is None:
            assert not self.stack

def collect_one_node(collector: Collector) -> CollectReport:
    ihook = collector.ihook
    ihook.pytest_collectstart(collector=collector)
    rep: CollectReport = ihook.pytest_make_collect_report(collector=collector)
    call = rep.__dict__.pop('call', None)
    if call and check_interactive_exception(call, rep):
        ihook.pytest_exception_interact(node=collector, call=call, report=rep)
    return rep
