from __future__ import annotations
import threading
import types
from ._exceptions import ExceptionMapping, PoolTimeout, map_exceptions
try:
    import trio
except (ImportError, NotImplementedError):
    trio = None
try:
    import anyio
except ImportError:
    anyio = None

def current_async_library() -> str:
    try:
        import sniffio
    except ImportError:
        environment = 'asyncio'
    else:
        environment = sniffio.current_async_library()
    if environment not in ('asyncio', 'trio'):
        raise RuntimeError('Running under an unsupported async environment.')
    if environment == 'asyncio' and anyio is None:
        raise RuntimeError("Running with asyncio requires installation of 'httpcore[asyncio]'.")
    if environment == 'trio' and trio is None:
        raise RuntimeError("Running with trio requires installation of 'httpcore[trio]'.")
    return environment

class AsyncLock:

    def __init__(self) -> None:
        self._backend = ''

    def setup(self) -> None:
        self._backend = current_async_library()
        if self._backend == 'trio':
            self._trio_lock = trio.Lock()
        elif self._backend == 'asyncio':
            self._anyio_lock = anyio.Lock()

    async def __aenter__(self) -> AsyncLock:
        if not self._backend:
            self.setup()
        if self._backend == 'trio':
            await self._trio_lock.acquire()
        elif self._backend == 'asyncio':
            await self._anyio_lock.acquire()
        return self

    async def __aexit__(self, exc_type: type[BaseException] | None=None, exc_value: BaseException | None=None, traceback: types.TracebackType | None=None) -> None:
        if self._backend == 'trio':
            self._trio_lock.release()
        elif self._backend == 'asyncio':
            self._anyio_lock.release()

class AsyncThreadLock:

    def __enter__(self) -> AsyncThreadLock:
        return self

    def __exit__(self, exc_type: type[BaseException] | None=None, exc_value: BaseException | None=None, traceback: types.TracebackType | None=None) -> None:
        pass

class AsyncEvent:

    def __init__(self) -> None:
        self._backend = ''

    def setup(self) -> None:
        self._backend = current_async_library()
        if self._backend == 'trio':
            self._trio_event = trio.Event()
        elif self._backend == 'asyncio':
            self._anyio_event = anyio.Event()

    def set(self) -> None:
        if not self._backend:
            self.setup()
        if self._backend == 'trio':
            self._trio_event.set()
        elif self._backend == 'asyncio':
            self._anyio_event.set()

    async def wait(self, timeout: float | None=None) -> None:
        if not self._backend:
            self.setup()
        if self._backend == 'trio':
            trio_exc_map: ExceptionMapping = {trio.TooSlowError: PoolTimeout}
            timeout_or_inf = float('inf') if timeout is None else timeout
            with map_exceptions(trio_exc_map):
                with trio.fail_after(timeout_or_inf):
                    await self._trio_event.wait()
        elif self._backend == 'asyncio':
            anyio_exc_map: ExceptionMapping = {TimeoutError: PoolTimeout}
            with map_exceptions(anyio_exc_map):
                with anyio.fail_after(timeout):
                    await self._anyio_event.wait()

class AsyncSemaphore:

    def __init__(self, bound: int) -> None:
        self._bound = bound
        self._backend = ''

    def setup(self) -> None:
        self._backend = current_async_library()
        if self._backend == 'trio':
            self._trio_semaphore = trio.Semaphore(initial_value=self._bound, max_value=self._bound)
        elif self._backend == 'asyncio':
            self._anyio_semaphore = anyio.Semaphore(initial_value=self._bound, max_value=self._bound)

    async def acquire(self) -> None:
        if not self._backend:
            self.setup()
        if self._backend == 'trio':
            await self._trio_semaphore.acquire()
        elif self._backend == 'asyncio':
            await self._anyio_semaphore.acquire()

    async def release(self) -> None:
        if self._backend == 'trio':
            self._trio_semaphore.release()
        elif self._backend == 'asyncio':
            self._anyio_semaphore.release()

class AsyncShieldCancellation:

    def __init__(self) -> None:
        self._backend = current_async_library()
        if self._backend == 'trio':
            self._trio_shield = trio.CancelScope(shield=True)
        elif self._backend == 'asyncio':
            self._anyio_shield = anyio.CancelScope(shield=True)

    def __enter__(self) -> AsyncShieldCancellation:
        if self._backend == 'trio':
            self._trio_shield.__enter__()
        elif self._backend == 'asyncio':
            self._anyio_shield.__enter__()
        return self

    def __exit__(self, exc_type: type[BaseException] | None=None, exc_value: BaseException | None=None, traceback: types.TracebackType | None=None) -> None:
        if self._backend == 'trio':
            self._trio_shield.__exit__(exc_type, exc_value, traceback)
        elif self._backend == 'asyncio':
            self._anyio_shield.__exit__(exc_type, exc_value, traceback)

class Lock:

    def __init__(self) -> None:
        self._lock = threading.Lock()

    def __enter__(self) -> Lock:
        self._lock.acquire()
        return self

    def __exit__(self, exc_type: type[BaseException] | None=None, exc_value: BaseException | None=None, traceback: types.TracebackType | None=None) -> None:
        self._lock.release()

class ThreadLock:

    def __init__(self) -> None:
        self._lock = threading.Lock()

    def __enter__(self) -> ThreadLock:
        self._lock.acquire()
        return self

    def __exit__(self, exc_type: type[BaseException] | None=None, exc_value: BaseException | None=None, traceback: types.TracebackType | None=None) -> None:
        self._lock.release()

class Event:

    def __init__(self) -> None:
        self._event = threading.Event()

    def set(self) -> None:
        self._event.set()

    def wait(self, timeout: float | None=None) -> None:
        if timeout == float('inf'):
            timeout = None
        if not self._event.wait(timeout=timeout):
            raise PoolTimeout()

class Semaphore:

    def __init__(self, bound: int) -> None:
        self._semaphore = threading.Semaphore(value=bound)

    def acquire(self) -> None:
        self._semaphore.acquire()

    def release(self) -> None:
        self._semaphore.release()

class ShieldCancellation:

    def __enter__(self) -> ShieldCancellation:
        return self

    def __exit__(self, exc_type: type[BaseException] | None=None, exc_value: BaseException | None=None, traceback: types.TracebackType | None=None) -> None:
        pass
