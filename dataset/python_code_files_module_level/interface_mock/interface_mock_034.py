from __future__ import annotations
import ssl
import sys
import types
import typing
from .._backends.auto import AutoBackend
from .._backends.base import SOCKET_OPTION, AsyncNetworkBackend
from .._exceptions import ConnectionNotAvailable, UnsupportedProtocol
from .._models import Origin, Proxy, Request, Response
from .._synchronization import AsyncEvent, AsyncShieldCancellation, AsyncThreadLock
from .connection import AsyncHTTPConnection
from .interfaces import AsyncConnectionInterface, AsyncRequestInterface

class AsyncPoolRequest:

    def __init__(self, request: Request) -> None:
        self.request = request
        self.connection: AsyncConnectionInterface | None = None
        self._connection_acquired = AsyncEvent()

    def assign_to_connection(self, connection: AsyncConnectionInterface | None) -> None:
        self.connection = connection
        self._connection_acquired.set()

    def clear_connection(self) -> None:
        self.connection = None
        self._connection_acquired = AsyncEvent()

    async def wait_for_connection(self, timeout: float | None=None) -> AsyncConnectionInterface:
        if self.connection is None:
            await self._connection_acquired.wait(timeout=timeout)
        assert self.connection is not None
        return self.connection

    def is_queued(self) -> bool:
        return self.connection is None

class AsyncConnectionPool(AsyncRequestInterface):

    def __init__(self, ssl_context: ssl.SSLContext | None=None, proxy: Proxy | None=None, max_connections: int | None=10, max_keepalive_connections: int | None=None, keepalive_expiry: float | None=None, http1: bool=True, http2: bool=False, retries: int=0, local_address: str | None=None, uds: str | None=None, network_backend: AsyncNetworkBackend | None=None, socket_options: typing.Iterable[SOCKET_OPTION] | None=None) -> None:
        self._ssl_context = ssl_context
        self._proxy = proxy
        self._max_connections = sys.maxsize if max_connections is None else max_connections
        self._max_keepalive_connections = sys.maxsize if max_keepalive_connections is None else max_keepalive_connections
        self._max_keepalive_connections = min(self._max_connections, self._max_keepalive_connections)
        self._keepalive_expiry = keepalive_expiry
        self._http1 = http1
        self._http2 = http2
        self._retries = retries
        self._local_address = local_address
        self._uds = uds
        self._network_backend = AutoBackend() if network_backend is None else network_backend
        self._socket_options = socket_options
        self._connections: list[AsyncConnectionInterface] = []
        self._requests: list[AsyncPoolRequest] = []
        self._optional_thread_lock = AsyncThreadLock()

    def create_connection(self, origin: Origin) -> AsyncConnectionInterface:
        if self._proxy is not None:
            if self._proxy.url.scheme in (b'socks5', b'socks5h'):
                from .socks_proxy import AsyncSocks5Connection
                return AsyncSocks5Connection(proxy_origin=self._proxy.url.origin, proxy_auth=self._proxy.auth, remote_origin=origin, ssl_context=self._ssl_context, keepalive_expiry=self._keepalive_expiry, http1=self._http1, http2=self._http2, network_backend=self._network_backend)
            elif origin.scheme == b'http':
                from .http_proxy import AsyncForwardHTTPConnection
                return AsyncForwardHTTPConnection(proxy_origin=self._proxy.url.origin, proxy_headers=self._proxy.headers, proxy_ssl_context=self._proxy.ssl_context, remote_origin=origin, keepalive_expiry=self._keepalive_expiry, network_backend=self._network_backend)
            from .http_proxy import AsyncTunnelHTTPConnection
            return AsyncTunnelHTTPConnection(proxy_origin=self._proxy.url.origin, proxy_headers=self._proxy.headers, proxy_ssl_context=self._proxy.ssl_context, remote_origin=origin, ssl_context=self._ssl_context, keepalive_expiry=self._keepalive_expiry, http1=self._http1, http2=self._http2, network_backend=self._network_backend)
        return AsyncHTTPConnection(origin=origin, ssl_context=self._ssl_context, keepalive_expiry=self._keepalive_expiry, http1=self._http1, http2=self._http2, retries=self._retries, local_address=self._local_address, uds=self._uds, network_backend=self._network_backend, socket_options=self._socket_options)

    @property
    def connections(self) -> list[AsyncConnectionInterface]:
        return list(self._connections)

    async def handle_async_request(self, request: Request) -> Response:
        scheme = request.url.scheme.decode()
        if scheme == '':
            raise UnsupportedProtocol("Request URL is missing an 'http://' or 'https://' protocol.")
        if scheme not in ('http', 'https', 'ws', 'wss'):
            raise UnsupportedProtocol(f"Request URL has an unsupported protocol '{scheme}://'.")
        timeouts = request.extensions.get('timeout', {})
        timeout = timeouts.get('pool', None)
        with self._optional_thread_lock:
            pool_request = AsyncPoolRequest(request)
            self._requests.append(pool_request)
        try:
            while True:
                with self._optional_thread_lock:
                    closing = self._assign_requests_to_connections()
                await self._close_connections(closing)
                connection = await pool_request.wait_for_connection(timeout=timeout)
                try:
                    response = await connection.handle_async_request(pool_request.request)
                except ConnectionNotAvailable:
                    pool_request.clear_connection()
                else:
                    break
        except BaseException as exc:
            with self._optional_thread_lock:
                self._requests.remove(pool_request)
                closing = self._assign_requests_to_connections()
            await self._close_connections(closing)
            raise exc from None
        assert isinstance(response.stream, typing.AsyncIterable)
        return Response(status=response.status, headers=response.headers, content=PoolByteStream(stream=response.stream, pool_request=pool_request, pool=self), extensions=response.extensions)

    def _assign_requests_to_connections(self) -> list[AsyncConnectionInterface]:
        closing_connections = []
        for connection in list(self._connections):
            if connection.is_closed():
                self._connections.remove(connection)
            elif connection.has_expired():
                self._connections.remove(connection)
                closing_connections.append(connection)
            elif connection.is_idle() and len([connection.is_idle() for connection in self._connections]) > self._max_keepalive_connections:
                self._connections.remove(connection)
                closing_connections.append(connection)
        queued_requests = [request for request in self._requests if request.is_queued()]
        for pool_request in queued_requests:
            origin = pool_request.request.url.origin
            available_connections = [connection for connection in self._connections if connection.can_handle_request(origin) and connection.is_available()]
            idle_connections = [connection for connection in self._connections if connection.is_idle()]
            if available_connections:
                connection = available_connections[0]
                pool_request.assign_to_connection(connection)
            elif len(self._connections) < self._max_connections:
                connection = self.create_connection(origin)
                self._connections.append(connection)
                pool_request.assign_to_connection(connection)
            elif idle_connections:
                connection = idle_connections[0]
                self._connections.remove(connection)
                closing_connections.append(connection)
                connection = self.create_connection(origin)
                self._connections.append(connection)
                pool_request.assign_to_connection(connection)
        return closing_connections

    async def _close_connections(self, closing: list[AsyncConnectionInterface]) -> None:
        with AsyncShieldCancellation():
            for connection in closing:
                await connection.aclose()

    async def aclose(self) -> None:
        with self._optional_thread_lock:
            closing_connections = list(self._connections)
            self._connections = []
        await self._close_connections(closing_connections)

    async def __aenter__(self) -> AsyncConnectionPool:
        return self

    async def __aexit__(self, exc_type: type[BaseException] | None=None, exc_value: BaseException | None=None, traceback: types.TracebackType | None=None) -> None:
        await self.aclose()

    def __repr__(self) -> str:
        class_name = self.__class__.__name__
        with self._optional_thread_lock:
            request_is_queued = [request.is_queued() for request in self._requests]
            connection_is_idle = [connection.is_idle() for connection in self._connections]
            num_active_requests = request_is_queued.count(False)
            num_queued_requests = request_is_queued.count(True)
            num_active_connections = connection_is_idle.count(False)
            num_idle_connections = connection_is_idle.count(True)
        requests_info = f'Requests: {num_active_requests} active, {num_queued_requests} queued'
        connection_info = f'Connections: {num_active_connections} active, {num_idle_connections} idle'
        return f'<{class_name} [{requests_info} | {connection_info}]>'

class PoolByteStream:

    def __init__(self, stream: typing.AsyncIterable[bytes], pool_request: AsyncPoolRequest, pool: AsyncConnectionPool) -> None:
        self._stream = stream
        self._pool_request = pool_request
        self._pool = pool
        self._closed = False

    async def __aiter__(self) -> typing.AsyncIterator[bytes]:
        try:
            async for part in self._stream:
                yield part
        except BaseException as exc:
            await self.aclose()
            raise exc from None

    async def aclose(self) -> None:
        if not self._closed:
            self._closed = True
            with AsyncShieldCancellation():
                if hasattr(self._stream, 'aclose'):
                    await self._stream.aclose()
            with self._pool._optional_thread_lock:
                self._pool._requests.remove(self._pool_request)
                closing = self._pool._assign_requests_to_connections()
            await self._pool._close_connections(closing)
