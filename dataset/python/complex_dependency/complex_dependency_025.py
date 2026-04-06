from __future__ import annotations
import enum
import logging
import ssl
import time
import types
import typing
import h11
from .._backends.base import NetworkStream
from .._exceptions import ConnectionNotAvailable, LocalProtocolError, RemoteProtocolError, WriteError, map_exceptions
from .._models import Origin, Request, Response
from .._synchronization import Lock, ShieldCancellation
from .._trace import Trace
from .interfaces import ConnectionInterface
logger = logging.getLogger('httpcore.http11')
H11SendEvent = typing.Union[h11.Request, h11.Data, h11.EndOfMessage]

class HTTPConnectionState(enum.IntEnum):
    NEW = 0
    ACTIVE = 1
    IDLE = 2
    CLOSED = 3

class HTTP11Connection(ConnectionInterface):
    READ_NUM_BYTES = 64 * 1024
    MAX_INCOMPLETE_EVENT_SIZE = 100 * 1024

    def __init__(self, origin: Origin, stream: NetworkStream, keepalive_expiry: float | None=None) -> None:
        self._origin = origin
        self._network_stream = stream
        self._keepalive_expiry: float | None = keepalive_expiry
        self._expire_at: float | None = None
        self._state = HTTPConnectionState.NEW
        self._state_lock = Lock()
        self._request_count = 0
        self._h11_state = h11.Connection(our_role=h11.CLIENT, max_incomplete_event_size=self.MAX_INCOMPLETE_EVENT_SIZE)

    def handle_request(self, request: Request) -> Response:
        if not self.can_handle_request(request.url.origin):
            raise RuntimeError(f'Attempted to send request to {request.url.origin} on connection to {self._origin}')
        with self._state_lock:
            if self._state in (HTTPConnectionState.NEW, HTTPConnectionState.IDLE):
                self._request_count += 1
                self._state = HTTPConnectionState.ACTIVE
                self._expire_at = None
            else:
                raise ConnectionNotAvailable()
        try:
            kwargs = {'request': request}
            try:
                with Trace('send_request_headers', logger, request, kwargs) as trace:
                    self._send_request_headers(**kwargs)
                with Trace('send_request_body', logger, request, kwargs) as trace:
                    self._send_request_body(**kwargs)
            except WriteError:
                pass
            with Trace('receive_response_headers', logger, request, kwargs) as trace:
                http_version, status, reason_phrase, headers, trailing_data = self._receive_response_headers(**kwargs)
                trace.return_value = (http_version, status, reason_phrase, headers)
            network_stream = self._network_stream
            if status == 101 or (request.method == b'CONNECT' and 200 <= status < 300):
                network_stream = HTTP11UpgradeStream(network_stream, trailing_data)
            return Response(status=status, headers=headers, content=HTTP11ConnectionByteStream(self, request), extensions={'http_version': http_version, 'reason_phrase': reason_phrase, 'network_stream': network_stream})
        except BaseException as exc:
            with ShieldCancellation():
                with Trace('response_closed', logger, request) as trace:
                    self._response_closed()
            raise exc

    def _send_request_headers(self, request: Request) -> None:
        timeouts = request.extensions.get('timeout', {})
        timeout = timeouts.get('write', None)
        with map_exceptions({h11.LocalProtocolError: LocalProtocolError}):
            event = h11.Request(method=request.method, target=request.url.target, headers=request.headers)
        self._send_event(event, timeout=timeout)

    def _send_request_body(self, request: Request) -> None:
        timeouts = request.extensions.get('timeout', {})
        timeout = timeouts.get('write', None)
        assert isinstance(request.stream, typing.Iterable)
        for chunk in request.stream:
            event = h11.Data(data=chunk)
            self._send_event(event, timeout=timeout)
        self._send_event(h11.EndOfMessage(), timeout=timeout)

    def _send_event(self, event: h11.Event, timeout: float | None=None) -> None:
        bytes_to_send = self._h11_state.send(event)
        if bytes_to_send is not None:
            self._network_stream.write(bytes_to_send, timeout=timeout)

    def _receive_response_headers(self, request: Request) -> tuple[bytes, int, bytes, list[tuple[bytes, bytes]], bytes]:
        timeouts = request.extensions.get('timeout', {})
        timeout = timeouts.get('read', None)
        while True:
            event = self._receive_event(timeout=timeout)
            if isinstance(event, h11.Response):
                break
            if isinstance(event, h11.InformationalResponse) and event.status_code == 101:
                break
        http_version = b'HTTP/' + event.http_version
        headers = event.headers.raw_items()
        trailing_data, _ = self._h11_state.trailing_data
        return (http_version, event.status_code, event.reason, headers, trailing_data)

    def _receive_response_body(self, request: Request) -> typing.Iterator[bytes]:
        timeouts = request.extensions.get('timeout', {})
        timeout = timeouts.get('read', None)
        while True:
            event = self._receive_event(timeout=timeout)
            if isinstance(event, h11.Data):
                yield bytes(event.data)
            elif isinstance(event, (h11.EndOfMessage, h11.PAUSED)):
                break

    def _receive_event(self, timeout: float | None=None) -> h11.Event | type[h11.PAUSED]:
        while True:
            with map_exceptions({h11.RemoteProtocolError: RemoteProtocolError}):
                event = self._h11_state.next_event()
            if event is h11.NEED_DATA:
                data = self._network_stream.read(self.READ_NUM_BYTES, timeout=timeout)
                if data == b'' and self._h11_state.their_state == h11.SEND_RESPONSE:
                    msg = 'Server disconnected without sending a response.'
                    raise RemoteProtocolError(msg)
                self._h11_state.receive_data(data)
            else:
                return event

    def _response_closed(self) -> None:
        with self._state_lock:
            if self._h11_state.our_state is h11.DONE and self._h11_state.their_state is h11.DONE:
                self._state = HTTPConnectionState.IDLE
                self._h11_state.start_next_cycle()
                if self._keepalive_expiry is not None:
                    now = time.monotonic()
                    self._expire_at = now + self._keepalive_expiry
            else:
                self.close()

    def close(self) -> None:
        self._state = HTTPConnectionState.CLOSED
        self._network_stream.close()

    def can_handle_request(self, origin: Origin) -> bool:
        return origin == self._origin

    def is_available(self) -> bool:
        return self._state == HTTPConnectionState.IDLE

    def has_expired(self) -> bool:
        now = time.monotonic()
        keepalive_expired = self._expire_at is not None and now > self._expire_at
        server_disconnected = self._state == HTTPConnectionState.IDLE and self._network_stream.get_extra_info('is_readable')
        return keepalive_expired or server_disconnected

    def is_idle(self) -> bool:
        return self._state == HTTPConnectionState.IDLE

    def is_closed(self) -> bool:
        return self._state == HTTPConnectionState.CLOSED

    def info(self) -> str:
        origin = str(self._origin)
        return f'{origin!r}, HTTP/1.1, {self._state.name}, Request Count: {self._request_count}'

    def __repr__(self) -> str:
        class_name = self.__class__.__name__
        origin = str(self._origin)
        return f'<{class_name} [{origin!r}, {self._state.name}, Request Count: {self._request_count}]>'

    def __enter__(self) -> HTTP11Connection:
        return self

    def __exit__(self, exc_type: type[BaseException] | None=None, exc_value: BaseException | None=None, traceback: types.TracebackType | None=None) -> None:
        self.close()

class HTTP11ConnectionByteStream:

    def __init__(self, connection: HTTP11Connection, request: Request) -> None:
        self._connection = connection
        self._request = request
        self._closed = False

    def __iter__(self) -> typing.Iterator[bytes]:
        kwargs = {'request': self._request}
        try:
            with Trace('receive_response_body', logger, self._request, kwargs):
                for chunk in self._connection._receive_response_body(**kwargs):
                    yield chunk
        except BaseException as exc:
            with ShieldCancellation():
                self.close()
            raise exc

    def close(self) -> None:
        if not self._closed:
            self._closed = True
            with Trace('response_closed', logger, self._request):
                self._connection._response_closed()

class HTTP11UpgradeStream(NetworkStream):

    def __init__(self, stream: NetworkStream, leading_data: bytes) -> None:
        self._stream = stream
        self._leading_data = leading_data

    def read(self, max_bytes: int, timeout: float | None=None) -> bytes:
        if self._leading_data:
            buffer = self._leading_data[:max_bytes]
            self._leading_data = self._leading_data[max_bytes:]
            return buffer
        else:
            return self._stream.read(max_bytes, timeout)

    def write(self, buffer: bytes, timeout: float | None=None) -> None:
        self._stream.write(buffer, timeout)

    def close(self) -> None:
        self._stream.close()

    def start_tls(self, ssl_context: ssl.SSLContext, server_hostname: str | None=None, timeout: float | None=None) -> NetworkStream:
        return self._stream.start_tls(ssl_context, server_hostname, timeout)

    def get_extra_info(self, info: str) -> typing.Any:
        return self._stream.get_extra_info(info)
