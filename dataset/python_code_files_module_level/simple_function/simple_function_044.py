from __future__ import annotations
import contextlib
import functools
import os
import sys
from collections.abc import Callable, Mapping
from types import FrameType
from typing import Any, TypeVar, cast
from coverage import env
from coverage.core import Core
from coverage.data import CoverageData
from coverage.debug import short_stack
from coverage.exceptions import ConfigError
from coverage.misc import human_sorted_items, isolate_module
from coverage.plugin import CoveragePlugin
from coverage.types import TArc, TCheckIncludeFn, TFileDisposition, Tracer, TShouldStartContextFn, TShouldTraceFn, TTraceData, TTraceFn, TWarnFn
os = isolate_module(os)
T = TypeVar('T')

class Collector:
    _collectors: list[Collector] = []

    def __init__(self, core: Core, should_trace: TShouldTraceFn, check_include: TCheckIncludeFn, should_start_context: TShouldStartContextFn | None, file_mapper: Callable[[str], str], branch: bool, warn: TWarnFn, concurrency: list[str]) -> None:
        self.core = core
        self.should_trace = should_trace
        self.check_include = check_include
        self.should_start_context = should_start_context
        self.file_mapper = file_mapper
        self.branch = branch
        self.warn = warn
        assert isinstance(concurrency, list), f'Expected a list: {concurrency!r}'
        self.pid = os.getpid()
        self.covdata: CoverageData
        self.threading = None
        self.static_context: str | None = None
        self.origin = short_stack()
        self.concur_id_func = None
        do_threading = False
        tried = 'nothing'
        try:
            if 'greenlet' in concurrency:
                tried = 'greenlet'
                import greenlet
                self.concur_id_func = greenlet.getcurrent
            elif 'eventlet' in concurrency:
                tried = 'eventlet'
                import eventlet.greenthread
                self.concur_id_func = eventlet.greenthread.getcurrent
            elif 'gevent' in concurrency:
                tried = 'gevent'
                import gevent
                self.concur_id_func = gevent.getcurrent
            if 'thread' in concurrency:
                do_threading = True
        except ImportError as ex:
            msg = f"Couldn't trace with concurrency={tried}, the module isn't installed."
            raise ConfigError(msg) from ex
        if self.concur_id_func and (not hasattr(core.tracer_class, 'concur_id_func')):
            raise ConfigError("Can't support concurrency={} with {}, only threads are supported.".format(tried, self.tracer_name()))
        if do_threading or not concurrency:
            import threading
            self.threading = threading
        self.reset()

    def __repr__(self) -> str:
        return f'<Collector at {id(self):#x}: {self.tracer_name()}>'

    def use_data(self, covdata: CoverageData, context: str | None) -> None:
        self.covdata = covdata
        self.static_context = context
        self.covdata.set_context(self.static_context)

    def tracer_name(self) -> str:
        return self.core.tracer_class.__name__

    def _clear_data(self) -> None:
        with self.data_lock or contextlib.nullcontext():
            for d in self.data.values():
                d.clear()
        for tracer in self.tracers:
            tracer.reset_activity()

    def reset(self) -> None:
        self.data_lock = self.threading.Lock() if self.threading else None
        self.data: TTraceData = {}
        self.file_tracers: dict[str, str] = {}
        self.disabled_plugins: set[str] = set()
        if env.PYPY:
            import __pypy__
            self.should_trace_cache = __pypy__.newdict('module')
        else:
            self.should_trace_cache = {}
        self.tracers: list[Tracer] = []
        self._clear_data()

    def lock_data(self) -> None:
        if self.data_lock is not None:
            self.data_lock.acquire()

    def unlock_data(self) -> None:
        if self.data_lock is not None:
            self.data_lock.release()

    def _start_tracer(self) -> TTraceFn | None:
        tracer = self.core.tracer_class(**self.core.tracer_kwargs)
        tracer.data = self.data
        tracer.lock_data = self.lock_data
        tracer.unlock_data = self.unlock_data
        tracer.trace_arcs = self.branch
        tracer.should_trace = self.should_trace
        tracer.should_trace_cache = self.should_trace_cache
        tracer.warn = self.warn
        if hasattr(tracer, 'concur_id_func'):
            tracer.concur_id_func = self.concur_id_func
        if hasattr(tracer, 'file_tracers'):
            tracer.file_tracers = self.file_tracers
        if hasattr(tracer, 'threading'):
            tracer.threading = self.threading
        if hasattr(tracer, 'check_include'):
            tracer.check_include = self.check_include
        if hasattr(tracer, 'should_start_context'):
            tracer.should_start_context = self.should_start_context
        if hasattr(tracer, 'switch_context'):
            tracer.switch_context = self.switch_context
        if hasattr(tracer, 'disable_plugin'):
            tracer.disable_plugin = self.disable_plugin
        fn = tracer.start()
        self.tracers.append(tracer)
        return fn

    def _installation_trace(self, frame: FrameType, event: str, arg: Any) -> TTraceFn | None:
        sys.settrace(None)
        fn: TTraceFn | None = self._start_tracer()
        if fn:
            fn = fn(frame, event, arg)
        return fn

    def start(self) -> None:
        keep_collectors = []
        for c in self._collectors:
            if c.pid == self.pid:
                keep_collectors.append(c)
            else:
                c.post_fork()
        self._collectors[:] = keep_collectors
        if self._collectors:
            self._collectors[-1].pause()
        self.tracers = []
        try:
            self._start_tracer()
        except:
            if self._collectors:
                self._collectors[-1].resume()
            raise
        self._collectors.append(self)
        if self.core.systrace and self.threading:
            self.threading.settrace(self._installation_trace)

    def stop(self) -> None:
        assert self._collectors
        if self._collectors[-1] is not self:
            print('self._collectors:')
            for c in self._collectors:
                print(f'  {c!r}\n{c.origin}')
        assert self._collectors[-1] is self, f"Expected current collector to be {self!r}, but it's {self._collectors[-1]!r}"
        self.pause()
        self._collectors.pop()
        if self._collectors:
            self._collectors[-1].resume()

    def pause(self) -> None:
        for tracer in self.tracers:
            tracer.stop()
            stats = tracer.get_stats()
            if stats:
                print(f'\nCoverage.py {tracer.__class__.__name__} stats:')
                for k, v in human_sorted_items(stats.items()):
                    print(f'{k:>20}: {v}')
        if self.threading:
            self.threading.settrace(None)

    def resume(self) -> None:
        for tracer in self.tracers:
            tracer.start()
        if self.core.systrace:
            if self.threading:
                self.threading.settrace(self._installation_trace)
            else:
                self._start_tracer()

    def post_fork(self) -> None:
        for tracer in self.tracers:
            if hasattr(tracer, 'post_fork'):
                tracer.post_fork()

    def _activity(self) -> bool:
        return any((tracer.activity() for tracer in self.tracers))

    def switch_context(self, new_context: str | None) -> None:
        context: str | None
        self.flush_data()
        if self.static_context:
            context = self.static_context
            if new_context:
                context += '|' + new_context
        else:
            context = new_context
        self.covdata.set_context(context)

    def disable_plugin(self, disposition: TFileDisposition) -> None:
        file_tracer = disposition.file_tracer
        assert file_tracer is not None
        plugin = file_tracer._coverage_plugin
        plugin_name = plugin._coverage_plugin_name
        self.warn(f'Disabling plug-in {plugin_name!r} due to previous exception')
        plugin._coverage_enabled = False
        disposition.trace = False

    @functools.cache
    def cached_mapped_file(self, filename: str) -> str:
        return self.file_mapper(filename)

    def mapped_file_dict(self, d: Mapping[str, T]) -> dict[str, T]:
        runtime_err = None
        for _ in range(3):
            try:
                items = list(d.items())
            except RuntimeError as ex:
                runtime_err = ex
            else:
                break
        else:
            assert isinstance(runtime_err, Exception)
            raise runtime_err
        return {self.cached_mapped_file(k): v for k, v in items if v}

    def plugin_was_disabled(self, plugin: CoveragePlugin) -> None:
        self.disabled_plugins.add(plugin._coverage_plugin_name)

    def flush_data(self) -> bool:
        if not self._activity():
            return False
        if self.branch:
            if self.core.packed_arcs:
                arc_data: dict[str, list[TArc]] = {}
                packed_data = cast(dict[str, set[int]], self.data)
                for fname, packeds in list(packed_data.items()):
                    tuples = []
                    for packed in list(packeds):
                        l1 = packed & 268435455
                        l2 = (packed & 268435455 << 28) >> 28
                        if packed & 1 << 56:
                            l1 *= -1
                        if packed & 1 << 57:
                            l2 *= -1
                        tuples.append((l1, l2))
                    arc_data[fname] = tuples
            else:
                arc_data = cast(dict[str, list[TArc]], self.data)
            self.covdata.add_arcs(self.mapped_file_dict(arc_data))
        else:
            line_data = cast(dict[str, set[int]], self.data)
            self.covdata.add_lines(self.mapped_file_dict(line_data))
        file_tracers = {self.cached_mapped_file(k): v for k, v in self.file_tracers.items() if v not in self.disabled_plugins}
        self.covdata.add_file_tracers(file_tracers)
        self._clear_data()
        return True
