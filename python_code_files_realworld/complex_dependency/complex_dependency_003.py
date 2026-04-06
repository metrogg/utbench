from __future__ import annotations
import enum
import logging
import time
import types
import typing
import h2.config
import h2.connection
import h2.events
import h2.exceptions
import h2.settings
from .._backends.base import NetworkStream
from .._exceptions import ConnectionNotAvailable, LocalProtocolError, RemoteProtocolError
from .._models import Origin, Request, Response
from .._synchronization import Lock, Semaphore, ShieldCancellation
from .._trace import Trace
from .interfaces import ConnectionInterface
logger = logging.getLogger('httpcore.http2')

def has_body_headers(request: Request) -> bool:
    return any((k.lower() == b'content-length' or k.lower() == b'transfer-encoding' for k, v in request.headers))

class HTTPConnectionState(enum.IntEnum):
    ACTIVE = 1
    IDLE = 2
    CLOSED = 3

class HTTP2Connection(ConnectionInterface):
    READ_NUM_BYTES = 64 * 1024
    CONFIG = h2.config.H2Configuration(validate_inbound_headers=False)

    def __init__(self, origin: Origin, stream: NetworkStream, keepalive_expiry: float | None=None):
        self._origin = origin
        self._network_stream = stream
        self._keepalive_expiry: float | None = keepalive_expiry
        self._h2_state = h2.connection.H2Connection(config=self.CONFIG)
        self._state = HTTPConnectionState.IDLE
        self._expire_at: float | None = None
        self._request_count = 0
        self._init_lock = Lock()
        self._state_lock = Lock()
        self._read_lock = Lock()
        self._write_lock = Lock()
        self._sent_connection_init = False
        self._used_all_stream_ids = False
        self._connection_error = False
        self._events: dict[int, list[h2.events.ResponseReceived | h2.events.DataReceived | h2.events.StreamEnded | h2.events.StreamReset,]] = {}
        self._connection_terminated: h2.events.ConnectionTerminated | None = None
        self._read_exception: Exception | None = None
        self._write_exception: Exception | None = None

    def handle_request(self, request: Request) -> Response:
        if not self.can_handle_request(request.url.origin):
            raise RuntimeError(f'Attempted to send request to {request.url.origin} on connection to {self._origin}')
        with self._state_lock:
            if self._state in (HTTPConnectionState.ACTIVE, HTTPConnectionState.IDLE):
                self._request_count += 1
                self._expire_at = None
                self._state = HTTPConnectionState.ACTIVE
            else:
                raise ConnectionNotAvailable()
        with self._init_lock:
            if not self._sent_connection_init:
                try:
                    sci_kwargs = {'request': request}
                    with Trace('send_connection_init', logger, request, sci_kwargs):
                        self._send_connection_init(**sci_kwargs)
                except BaseException as exc:
                    with ShieldCancellation():
                        self.close()
                    raise exc
                self._sent_connection_init = True
                self._max_streams = 1
                local_settings_max_streams = self._h2_state.local_settings.max_concurrent_streams
                self._max_streams_semaphore = Semaphore(local_settings_max_streams)
                for _ in range(local_settings_max_streams - self._max_streams):
                    self._max_streams_semaphore.acquire()
        self._max_streams_semaphore.acquire()
        try:
            stream_id = self._h2_state.get_next_available_stream_id()
            self._events[stream_id] = []
        except h2.exceptions.NoAvailableStreamIDError:
            self._used_all_stream_ids = True
            self._request_count -= 1
            raise ConnectionNotAvailable()
        try:
            kwargs = {'request': request, 'stream_id': stream_id}
            with Trace('send_request_headers', logger, request, kwargs):
                self._send_request_headers(request=request, stream_id=stream_id)
            with Trace('send_request_body', logger, request, kwargs):
                self._send_request_body(request=request, stream_id=stream_id)
            with Trace('receive_response_headers', logger, request, kwargs) as trace:
                status, headers = self._receive_response(request=request, stream_id=stream_id)
                trace.return_value = (status, headers)
            return Response(status=status, headers=headers, content=HTTP2ConnectionByteStream(self, request, stream_id=stream_id), extensions={'http_version': b'HTTP/2', 'network_stream': self._network_stream, 'stream_id': stream_id})
        except BaseException as exc:
            with ShieldCancellation():
                kwargs = {'stream_id': stream_id}
                with Trace('response_closed', logger, request, kwargs):
                    self._response_closed(stream_id=stream_id)
            if isinstance(exc, h2.exceptions.ProtocolError):
                if self._connection_terminated:
                    raise RemoteProtocolError(self._connection_terminated)
                raise LocalProtocolError(exc)
            raise exc

    def _send_connection_init(self, request: Request) -> None:
        self._h2_state.local_settings = h2.settings.Settings(client=True, initial_values={h2.settings.SettingCodes.ENABLE_PUSH: 0, h2.settings.SettingCodes.MAX_CONCURRENT_STREAMS: 100, h2.settings.SettingCodes.MAX_HEADER_LIST_SIZE: 65536})
        del self._h2_state.local_settings[h2.settings.SettingCodes.ENABLE_CONNECT_PROTOCOL]
        self._h2_state.initiate_connection()
        self._h2_state.increment_flow_control_window(2 ** 24)
        self._write_outgoing_data(request)

    def _send_request_headers(self, request: Request, stream_id: int) -> None:
        end_stream = not has_body_headers(request)
        authority = [v for k, v in request.headers if k.lower() == b'host'][0]
        headers = [(b':method', request.method), (b':authority', authority), (b':scheme', request.url.scheme), (b':path', request.url.target)] + [(k.lower(), v) for k, v in request.headers if k.lower() not in (b'host', b'transfer-encoding')]
        self._h2_state.send_headers(stream_id, headers, end_stream=end_stream)
        self._h2_state.increment_flow_control_window(2 ** 24, stream_id=stream_id)
        self._write_outgoing_data(request)

    def _send_request_body(self, request: Request, stream_id: int) -> None:
        if not has_body_headers(request):
            return
        assert isinstance(request.stream, typing.Iterable)
        for data in request.stream:
            self._send_stream_data(request, stream_id, data)
        self._send_end_stream(request, stream_id)

    def _send_stream_data(self, request: Request, stream_id: int, data: bytes) -> None:
        while data:
            max_flow = self._wait_for_outgoing_flow(request, stream_id)
            chunk_size = min(len(data), max_flow)
            chunk, data = (data[:chunk_size], data[chunk_size:])
            self._h2_state.send_data(stream_id, chunk)
            self._write_outgoing_data(request)

    def _send_end_stream(self, request: Request, stream_id: int) -> None:
        self._h2_state.end_stream(stream_id)
        self._write_outgoing_data(request)

    def _receive_response(self, request: Request, stream_id: int) -> tuple[int, list[tuple[bytes, bytes]]]:
        while True:
            event = self._receive_stream_event(request, stream_id)
            if isinstance(event, h2.events.ResponseReceived):
                break
        status_code = 200
        headers = []
        assert event.headers is not None
        for k, v in event.headers:
            if k == b':status':
                status_code = int(v.decode('ascii', errors='ignore'))
            elif not k.startswith(b':'):
                headers.append((k, v))
        return (status_code, headers)

    def _receive_response_body(self, request: Request, stream_id: int) -> typing.Iterator[bytes]:
        while True:
            event = self._receive_stream_event(request, stream_id)
            if isinstance(event, h2.events.DataReceived):
                assert event.flow_controlled_length is not None
                assert event.data is not None
                amount = event.flow_controlled_length
                self._h2_state.acknowledge_received_data(amount, stream_id)
                self._write_outgoing_data(request)
                yield event.data
            elif isinstance(event, h2.events.StreamEnded):
                break

    def _receive_stream_event(self, request: Request, stream_id: int) -> h2.events.ResponseReceived | h2.events.DataReceived | h2.events.StreamEnded:
        while not self._events.get(stream_id):
            self._receive_events(request, stream_id)
        event = self._events[stream_id].pop(0)
        if isinstance(event, h2.events.StreamReset):
            raise RemoteProtocolError(event)
        return event

    def _receive_events(self, request: Request, stream_id: int | None=None) -> None:
        with self._read_lock:
            if self._connection_terminated is not None:
                last_stream_id = self._connection_terminated.last_stream_id
                if stream_id and last_stream_id and (stream_id > last_stream_id):
                    self._request_count -= 1
                    raise ConnectionNotAvailable()
                raise RemoteProtocolError(self._connection_terminated)
            if stream_id is None or not self._events.get(stream_id):
                events = self._read_incoming_data(request)
                for event in events:
                    if isinstance(event, h2.events.RemoteSettingsChanged):
                        with Trace('receive_remote_settings', logger, request) as trace:
                            self._receive_remote_settings_change(event)
                            trace.return_value = event
                    elif isinstance(event, (h2.events.ResponseReceived, h2.events.DataReceived, h2.events.StreamEnded, h2.events.StreamReset)):
                        if event.stream_id in self._events:
                            self._events[event.stream_id].append(event)
                    elif isinstance(event, h2.events.ConnectionTerminated):
                        self._connection_terminated = event
        self._write_outgoing_data(request)

    def _receive_remote_settings_change(self, event: h2.events.RemoteSettingsChanged) -> None:
        max_concurrent_streams = event.changed_settings.get(h2.settings.SettingCodes.MAX_CONCURRENT_STREAMS)
        if max_concurrent_streams:
            new_max_streams = min(max_concurrent_streams.new_value, self._h2_state.local_settings.max_concurrent_streams)
            if new_max_streams and new_max_streams != self._max_streams:
                while new_max_streams > self._max_streams:
                    self._max_streams_semaphore.release()
                    self._max_streams += 1
                while new_max_streams < self._max_streams:
                    self._max_streams_semaphore.acquire()
                    self._max_streams -= 1

    def _response_closed(self, stream_id: int) -> None:
        self._max_streams_semaphore.release()
        del self._events[stream_id]
        with self._state_lock:
            if self._connection_terminated and (not self._events):
                self.close()
            elif self._state == HTTPConnectionState.ACTIVE and (not self._events):
                self._state = HTTPConnectionState.IDLE
                if self._keepalive_expiry is not None:
                    now = time.monotonic()
                    self._expire_at = now + self._keepalive_expiry
                if self._used_all_stream_ids:
                    self.close()

    def close(self) -> None:
        self._h2_state.close_connection()
        self._state = HTTPConnectionState.CLOSED
        self._network_stream.close()

    def _read_incoming_data(self, request: Request) -> list[h2.events.Event]:
        timeouts = request.extensions.get('timeout', {})
        timeout = timeouts.get('read', None)
        if self._read_exception is not None:
            raise self._read_exception
        try:
            data = self._network_stream.read(self.READ_NUM_BYTES, timeout)
            if data == b'':
                raise RemoteProtocolError('Server disconnected')
        except Exception as exc:
            self._read_exception = exc
            self._connection_error = True
            raise exc
        events: list[h2.events.Event] = self._h2_state.receive_data(data)
        return events

    def _write_outgoing_data(self, request: Request) -> None:
        timeouts = request.extensions.get('timeout', {})
        timeout = timeouts.get('write', None)
        with self._write_lock:
            data_to_send = self._h2_state.data_to_send()
            if self._write_exception is not None:
                raise self._write_exception
            try:
                self._network_stream.write(data_to_send, timeout)
            except Exception as exc:
                self._write_exception = exc
                self._connection_error = True
                raise exc

    def _wait_for_outgoing_flow(self, request: Request, stream_id: int) -> int:
        local_flow: int = self._h2_state.local_flow_control_window(stream_id)
        max_frame_size: int = self._h2_state.max_outbound_frame_size
        flow = min(local_flow, max_frame_size)
        while flow == 0:
            self._receive_events(request)
            local_flow = self._h2_state.local_flow_control_window(stream_id)
            max_frame_size = self._h2_state.max_outbound_frame_size
            flow = min(local_flow, max_frame_size)
        return flow

    def can_handle_request(self, origin: Origin) -> bool:
        return origin == self._origin

    def is_available(self) -> bool:
        return self._state != HTTPConnectionState.CLOSED and (not self._connection_error) and (not self._used_all_stream_ids) and (not self._h2_state.state_machine.state == h2.connection.ConnectionState.CLOSED)

    def has_expired(self) -> bool:
        now = time.monotonic()
        return self._expire_at is not None and now > self._expire_at

    def is_idle(self) -> bool:
        return self._state == HTTPConnectionState.IDLE

    def is_closed(self) -> bool:
        return self._state == HTTPConnectionState.CLOSED

    def info(self) -> str:
        origin = str(self._origin)
        return f'{origin!r}, HTTP/2, {self._state.name}, Request Count: {self._request_count}'

    def __repr__(self) -> str:
        class_name = self.__class__.__name__
        origin = str(self._origin)
        return f'<{class_name} [{origin!r}, {self._state.name}, Request Count: {self._request_count}]>'

    def __enter__(self) -> HTTP2Connection:
        return self

    def __exit__(self, exc_type: type[BaseException] | None=None, exc_value: BaseException | None=None, traceback: types.TracebackType | None=None) -> None:
        self.close()

class HTTP2ConnectionByteStream:

    def __init__(self, connection: HTTP2Connection, request: Request, stream_id: int) -> None:
        self._connection = connection
        self._request = request
        self._stream_id = stream_id
        self._closed = False

    def __iter__(self) -> typing.Iterator[bytes]:
        kwargs = {'request': self._request, 'stream_id': self._stream_id}
        try:
            with Trace('receive_response_body', logger, self._request, kwargs):
                for chunk in self._connection._receive_response_body(request=self._request, stream_id=self._stream_id):
                    yield chunk
        except BaseException as exc:
            with ShieldCancellation():
                self.close()
            raise exc

    def close(self) -> None:
        if not self._closed:
            self._closed = True
            kwargs = {'stream_id': self._stream_id}
            with Trace('response_closed', logger, self._request, kwargs):
                self._connection._response_closed(stream_id=self._stream_id)
