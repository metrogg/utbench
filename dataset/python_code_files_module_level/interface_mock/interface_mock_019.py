import asyncio
import contextlib
import inspect
import warnings
from typing import Any, Awaitable, Callable, Dict, Iterator, Optional, Protocol, Union, overload
import pytest
from .test_utils import BaseTestServer, RawTestServer, TestClient, TestServer, loop_context, setup_test_loop, teardown_test_loop, unused_port as _unused_port
from .web import Application, BaseRequest, Request
from .web_protocol import _RequestHandler
try:
    import uvloop
except ImportError:
    uvloop = None

class AiohttpClient(Protocol):

    @overload
    async def __call__(self, __param: Application, *, server_kwargs: Optional[Dict[str, Any]]=None, **kwargs: Any) -> TestClient[Request, Application]:
        ...

    @overload
    async def __call__(self, __param: BaseTestServer, *, server_kwargs: Optional[Dict[str, Any]]=None, **kwargs: Any) -> TestClient[BaseRequest, None]:
        ...

class AiohttpServer(Protocol):

    def __call__(self, app: Application, *, port: Optional[int]=None, **kwargs: Any) -> Awaitable[TestServer]:
        ...

class AiohttpRawServer(Protocol):

    def __call__(self, handler: _RequestHandler, *, port: Optional[int]=None, **kwargs: Any) -> Awaitable[RawTestServer]:
        ...

def pytest_addoption(parser):
    parser.addoption('--aiohttp-fast', action='store_true', default=False, help='run tests faster by disabling extra checks')
    parser.addoption('--aiohttp-loop', action='store', default='pyloop', help='run tests with specific loop: pyloop, uvloop or all')
    parser.addoption('--aiohttp-enable-loop-debug', action='store_true', default=False, help='enable event loop debug mode')

def pytest_fixture_setup(fixturedef):
    func = fixturedef.func
    if inspect.isasyncgenfunction(func):
        is_async_gen = True
    elif inspect.iscoroutinefunction(func):
        is_async_gen = False
    else:
        return
    strip_request = False
    if 'request' not in fixturedef.argnames:
        fixturedef.argnames += ('request',)
        strip_request = True

    def wrapper(*args, **kwargs):
        request = kwargs['request']
        if strip_request:
            del kwargs['request']
        if 'loop' not in request.fixturenames:
            raise Exception("Asynchronous fixtures must depend on the 'loop' fixture or be used in tests depending from it.")
        _loop = request.getfixturevalue('loop')
        if is_async_gen:
            gen = func(*args, **kwargs)

            def finalizer():
                try:
                    return _loop.run_until_complete(gen.__anext__())
                except StopAsyncIteration:
                    pass
            request.addfinalizer(finalizer)
            return _loop.run_until_complete(gen.__anext__())
        else:
            return _loop.run_until_complete(func(*args, **kwargs))
    fixturedef.func = wrapper

@pytest.fixture
def fast(request):
    return request.config.getoption('--aiohttp-fast')

@pytest.fixture
def loop_debug(request):
    return request.config.getoption('--aiohttp-enable-loop-debug')

@contextlib.contextmanager
def _runtime_warning_context():
    with warnings.catch_warnings(record=True) as _warnings:
        yield
        rw = ['{w.filename}:{w.lineno}:{w.message}'.format(w=w) for w in _warnings if w.category == RuntimeWarning]
        if rw:
            raise RuntimeError('{} Runtime Warning{},\n{}'.format(len(rw), '' if len(rw) == 1 else 's', '\n'.join(rw)))

@contextlib.contextmanager
def _passthrough_loop_context(loop, fast=False):
    if loop:
        yield loop
    else:
        loop = setup_test_loop()
        yield loop
        teardown_test_loop(loop, fast=fast)

def pytest_pycollect_makeitem(collector, name, obj):
    if collector.funcnamefilter(name) and inspect.iscoroutinefunction(obj):
        return list(collector._genfunctions(name, obj))

def pytest_pyfunc_call(pyfuncitem):
    fast = pyfuncitem.config.getoption('--aiohttp-fast')
    if inspect.iscoroutinefunction(pyfuncitem.function):
        existing_loop = pyfuncitem.funcargs.get('proactor_loop') or pyfuncitem.funcargs.get('selector_loop') or pyfuncitem.funcargs.get('uvloop_loop') or pyfuncitem.funcargs.get('loop', None)
        with _runtime_warning_context():
            with _passthrough_loop_context(existing_loop, fast=fast) as _loop:
                testargs = {arg: pyfuncitem.funcargs[arg] for arg in pyfuncitem._fixtureinfo.argnames}
                _loop.run_until_complete(pyfuncitem.obj(**testargs))
        return True

def pytest_generate_tests(metafunc):
    if 'loop_factory' not in metafunc.fixturenames:
        return
    loops = metafunc.config.option.aiohttp_loop
    avail_factories: dict[str, Callable[[], asyncio.AbstractEventLoop]]
    avail_factories = {'pyloop': asyncio.new_event_loop}
    if uvloop is not None:
        avail_factories['uvloop'] = uvloop.new_event_loop
    if loops == 'all':
        loops = 'pyloop,uvloop?'
    factories = {}
    for name in loops.split(','):
        required = not name.endswith('?')
        name = name.strip(' ?')
        if name not in avail_factories:
            if required:
                raise ValueError("Unknown loop '%s', available loops: %s" % (name, list(factories.keys())))
            else:
                continue
        factories[name] = avail_factories[name]
    metafunc.parametrize('loop_factory', list(factories.values()), ids=list(factories.keys()))

@pytest.fixture
def loop(loop_factory: Callable[[], asyncio.AbstractEventLoop], fast: bool, loop_debug: bool) -> Iterator[asyncio.AbstractEventLoop]:
    with loop_context(loop_factory, fast=fast) as _loop:
        if loop_debug:
            _loop.set_debug(True)
        asyncio.set_event_loop(_loop)
        yield _loop

@pytest.fixture
def proactor_loop() -> Iterator[asyncio.AbstractEventLoop]:
    factory = asyncio.ProactorEventLoop
    with loop_context(factory) as _loop:
        asyncio.set_event_loop(_loop)
        yield _loop

@pytest.fixture
def unused_port(aiohttp_unused_port: Callable[[], int]) -> Callable[[], int]:
    warnings.warn('Deprecated, use aiohttp_unused_port fixture instead', DeprecationWarning, stacklevel=2)
    return aiohttp_unused_port

@pytest.fixture
def aiohttp_unused_port() -> Callable[[], int]:
    return _unused_port

@pytest.fixture
def aiohttp_server(loop: asyncio.AbstractEventLoop) -> Iterator[AiohttpServer]:
    servers = []

    async def go(app: Application, *, host: str='127.0.0.1', port: Optional[int]=None, **kwargs: Any) -> TestServer:
        server = TestServer(app, host=host, port=port)
        await server.start_server(loop=loop, **kwargs)
        servers.append(server)
        return server
    yield go

    async def finalize() -> None:
        while servers:
            await servers.pop().close()
    loop.run_until_complete(finalize())

@pytest.fixture
def test_server(aiohttp_server):
    warnings.warn('Deprecated, use aiohttp_server fixture instead', DeprecationWarning, stacklevel=2)
    return aiohttp_server

@pytest.fixture
def aiohttp_raw_server(loop: asyncio.AbstractEventLoop) -> Iterator[AiohttpRawServer]:
    servers = []

    async def go(handler: _RequestHandler, *, port: Optional[int]=None, **kwargs: Any) -> RawTestServer:
        server = RawTestServer(handler, port=port)
        await server.start_server(loop=loop, **kwargs)
        servers.append(server)
        return server
    yield go

    async def finalize() -> None:
        while servers:
            await servers.pop().close()
    loop.run_until_complete(finalize())

@pytest.fixture
def raw_test_server(aiohttp_raw_server):
    warnings.warn('Deprecated, use aiohttp_raw_server fixture instead', DeprecationWarning, stacklevel=2)
    return aiohttp_raw_server

@pytest.fixture
def aiohttp_client(loop: asyncio.AbstractEventLoop) -> Iterator[AiohttpClient]:
    clients = []

    @overload
    async def go(__param: Application, *, server_kwargs: Optional[Dict[str, Any]]=None, **kwargs: Any) -> TestClient[Request, Application]:
        ...

    @overload
    async def go(__param: BaseTestServer, *, server_kwargs: Optional[Dict[str, Any]]=None, **kwargs: Any) -> TestClient[BaseRequest, None]:
        ...

    async def go(__param: Union[Application, BaseTestServer], *args: Any, server_kwargs: Optional[Dict[str, Any]]=None, **kwargs: Any) -> TestClient[Any, Any]:
        if isinstance(__param, Callable) and (not isinstance(__param, (Application, BaseTestServer))):
            __param = __param(loop, *args, **kwargs)
            kwargs = {}
        else:
            assert not args, 'args should be empty'
        if isinstance(__param, Application):
            server_kwargs = server_kwargs or {}
            server = TestServer(__param, loop=loop, **server_kwargs)
            client = TestClient(server, loop=loop, **kwargs)
        elif isinstance(__param, BaseTestServer):
            client = TestClient(__param, loop=loop, **kwargs)
        else:
            raise ValueError('Unknown argument type: %r' % type(__param))
        await client.start_server()
        clients.append(client)
        return client
    yield go

    async def finalize() -> None:
        while clients:
            await clients.pop().close()
    loop.run_until_complete(finalize())

@pytest.fixture
def test_client(aiohttp_client):
    warnings.warn('Deprecated, use aiohttp_client fixture instead', DeprecationWarning, stacklevel=2)
    return aiohttp_client
