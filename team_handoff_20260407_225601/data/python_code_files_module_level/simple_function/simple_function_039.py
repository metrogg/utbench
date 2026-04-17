from __future__ import annotations
import sys
from threading import Event, RLock, Thread
from types import TracebackType
from typing import IO, TYPE_CHECKING, Any, Callable, List, Optional, TextIO, Type, cast
from . import get_console
from .console import Console, ConsoleRenderable, Group, RenderableType, RenderHook
from .control import Control
from .file_proxy import FileProxy
from .jupyter import JupyterMixin
from .live_render import LiveRender, VerticalOverflowMethod
from .screen import Screen
from .text import Text
if TYPE_CHECKING:
    from typing_extensions import Self

class _RefreshThread(Thread):

    def __init__(self, live: 'Live', refresh_per_second: float) -> None:
        self.live = live
        self.refresh_per_second = refresh_per_second
        self.done = Event()
        super().__init__(daemon=True)

    def stop(self) -> None:
        self.done.set()

    def run(self) -> None:
        while not self.done.wait(1 / self.refresh_per_second):
            with self.live._lock:
                if not self.done.is_set():
                    self.live.refresh()

class Live(JupyterMixin, RenderHook):

    def __init__(self, renderable: Optional[RenderableType]=None, *, console: Optional[Console]=None, screen: bool=False, auto_refresh: bool=True, refresh_per_second: float=4, transient: bool=False, redirect_stdout: bool=True, redirect_stderr: bool=True, vertical_overflow: VerticalOverflowMethod='ellipsis', get_renderable: Optional[Callable[[], RenderableType]]=None) -> None:
        assert refresh_per_second > 0, 'refresh_per_second must be > 0'
        self._renderable = renderable
        self.console = console if console is not None else get_console()
        self._screen = screen
        self._alt_screen = False
        self._redirect_stdout = redirect_stdout
        self._redirect_stderr = redirect_stderr
        self._restore_stdout: Optional[IO[str]] = None
        self._restore_stderr: Optional[IO[str]] = None
        self._lock = RLock()
        self.ipy_widget: Optional[Any] = None
        self.auto_refresh = auto_refresh
        self._started: bool = False
        self.transient = True if screen else transient
        self._refresh_thread: Optional[_RefreshThread] = None
        self.refresh_per_second = refresh_per_second
        self.vertical_overflow = vertical_overflow
        self._get_renderable = get_renderable
        self._live_render = LiveRender(self.get_renderable(), vertical_overflow=vertical_overflow)
        self._nested = False

    @property
    def is_started(self) -> bool:
        return self._started

    def get_renderable(self) -> RenderableType:
        renderable = self._get_renderable() if self._get_renderable is not None else self._renderable
        return renderable or ''

    def start(self, refresh: bool=False) -> None:
        with self._lock:
            if self._started:
                return
            self._started = True
            if not self.console.set_live(self):
                self._nested = True
                return
            if self._screen:
                self._alt_screen = self.console.set_alt_screen(True)
            self.console.show_cursor(False)
            self._enable_redirect_io()
            self.console.push_render_hook(self)
            if refresh:
                try:
                    self.refresh()
                except Exception:
                    self.stop()
                    raise
            if self.auto_refresh:
                self._refresh_thread = _RefreshThread(self, self.refresh_per_second)
                self._refresh_thread.start()

    def stop(self) -> None:
        with self._lock:
            if not self._started:
                return
            self._started = False
            self.console.clear_live()
            if self._nested:
                if not self.transient:
                    self.console.print(self.renderable)
                return
            if self.auto_refresh and self._refresh_thread is not None:
                self._refresh_thread.stop()
                self._refresh_thread = None
            self.vertical_overflow = 'visible'
            with self.console:
                try:
                    if not self._alt_screen and (not self.console.is_jupyter):
                        self.refresh()
                finally:
                    self._disable_redirect_io()
                    self.console.pop_render_hook()
                    if not self._alt_screen and self.console.is_terminal and self._live_render.last_render_height:
                        self.console.line()
                    self.console.show_cursor(True)
                    if self._alt_screen:
                        self.console.set_alt_screen(False)
                    if self.transient and (not self._alt_screen):
                        self.console.control(self._live_render.restore_cursor())
                    if self.ipy_widget is not None and self.transient:
                        self.ipy_widget.close()

    def __enter__(self) -> Self:
        self.start(refresh=self._renderable is not None)
        return self

    def __exit__(self, exc_type: Optional[Type[BaseException]], exc_val: Optional[BaseException], exc_tb: Optional[TracebackType]) -> None:
        self.stop()

    def _enable_redirect_io(self) -> None:
        if self.console.is_terminal or self.console.is_jupyter:
            if self._redirect_stdout and (not isinstance(sys.stdout, FileProxy)):
                self._restore_stdout = sys.stdout
                sys.stdout = cast('TextIO', FileProxy(self.console, sys.stdout))
            if self._redirect_stderr and (not isinstance(sys.stderr, FileProxy)):
                self._restore_stderr = sys.stderr
                sys.stderr = cast('TextIO', FileProxy(self.console, sys.stderr))

    def _disable_redirect_io(self) -> None:
        if self._restore_stdout:
            sys.stdout = cast('TextIO', self._restore_stdout)
            self._restore_stdout = None
        if self._restore_stderr:
            sys.stderr = cast('TextIO', self._restore_stderr)
            self._restore_stderr = None

    @property
    def renderable(self) -> RenderableType:
        live_stack = self.console._live_stack
        renderable: RenderableType
        if live_stack and self is live_stack[0]:
            renderable = Group(*[live.get_renderable() for live in live_stack])
        else:
            renderable = self.get_renderable()
        return Screen(renderable) if self._alt_screen else renderable

    def update(self, renderable: RenderableType, *, refresh: bool=False) -> None:
        if isinstance(renderable, str):
            renderable = self.console.render_str(renderable)
        with self._lock:
            self._renderable = renderable
            if refresh:
                self.refresh()

    def refresh(self) -> None:
        with self._lock:
            self._live_render.set_renderable(self.renderable)
            if self._nested:
                if self.console._live_stack:
                    self.console._live_stack[0].refresh()
                return
            if self.console.is_jupyter:
                try:
                    from IPython.display import display
                    from ipywidgets import Output
                except ImportError:
                    import warnings
                    warnings.warn('install "ipywidgets" for Jupyter support')
                else:
                    if self.ipy_widget is None:
                        self.ipy_widget = Output()
                        display(self.ipy_widget)
                    with self.ipy_widget:
                        self.ipy_widget.clear_output(wait=True)
                        self.console.print(self._live_render.renderable)
            elif self.console.is_terminal and (not self.console.is_dumb_terminal):
                with self.console:
                    self.console.print(Control())
            elif not self._started and (not self.transient):
                with self.console:
                    self.console.print(Control())

    def process_renderables(self, renderables: List[ConsoleRenderable]) -> List[ConsoleRenderable]:
        self._live_render.vertical_overflow = self.vertical_overflow
        if self.console.is_interactive:
            with self._lock:
                reset = Control.home() if self._alt_screen else self._live_render.position_cursor()
                renderables = [reset, *renderables, self._live_render]
        elif not self._started and (not self.transient):
            renderables = [*renderables, self._live_render]
        return renderables
if __name__ == '__main__':
    import random
    import time
    from itertools import cycle
    from typing import Dict, List, Tuple
    from .align import Align
    from .console import Console
    from .live import Live as Live
    from .panel import Panel
    from .rule import Rule
    from .syntax import Syntax
    from .table import Table
    console = Console()
    syntax = Syntax('def loop_last(values: Iterable[T]) -> Iterable[Tuple[bool, T]]:\n    """Iterate and generate a tuple with a flag for last value."""\n    iter_values = iter(values)\n    try:\n        previous_value = next(iter_values)\n    except StopIteration:\n        return\n    for value in iter_values:\n        yield False, previous_value\n        previous_value = value\n    yield True, previous_value', 'python', line_numbers=True)
    table = Table('foo', 'bar', 'baz')
    table.add_row('1', '2', '3')
    progress_renderables = ['You can make the terminal shorter and taller to see the live table hideText may be printed while the progress bars are rendering.', Panel('In fact, [i]any[/i] renderable will work'), 'Such as [magenta]tables[/]...', table, 'Pretty printed structures...', {'type': 'example', 'text': 'Pretty printed'}, 'Syntax...', syntax, Rule('Give it a try!')]
    examples = cycle(progress_renderables)
    exchanges = ['SGD', 'MYR', 'EUR', 'USD', 'AUD', 'JPY', 'CNH', 'HKD', 'CAD', 'INR', 'DKK', 'GBP', 'RUB', 'NZD', 'MXN', 'IDR', 'TWD', 'THB', 'VND']
    with Live(console=console) as live_table:
        exchange_rate_dict: Dict[Tuple[str, str], float] = {}
        for index in range(100):
            select_exchange = exchanges[index % len(exchanges)]
            for exchange in exchanges:
                if exchange == select_exchange:
                    continue
                time.sleep(0.4)
                if random.randint(0, 10) < 1:
                    console.log(next(examples))
                exchange_rate_dict[select_exchange, exchange] = 200 / (random.random() * 320 + 1)
                if len(exchange_rate_dict) > len(exchanges) - 1:
                    exchange_rate_dict.pop(list(exchange_rate_dict.keys())[0])
                table = Table(title='Exchange Rates')
                table.add_column('Source Currency')
                table.add_column('Destination Currency')
                table.add_column('Exchange Rate')
                for (source, dest), exchange_rate in exchange_rate_dict.items():
                    table.add_row(source, dest, Text(f'{exchange_rate:.4f}', style='red' if exchange_rate < 1.0 else 'green'))
                live_table.update(Align.center(table))
