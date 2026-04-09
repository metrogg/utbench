from typing import Any, Callable, cast, Dict, List, Optional, overload, Tuple, Type, Union
from ._events import ConnectionClosed, Data, EndOfMessage, Event, InformationalResponse, Request, Response
from ._headers import get_comma_header, has_expect_100_continue, set_comma_header
from ._readers import READERS, ReadersType
from ._receivebuffer import ReceiveBuffer
from ._state import _SWITCH_CONNECT, _SWITCH_UPGRADE, CLIENT, ConnectionState, DONE, ERROR, MIGHT_SWITCH_PROTOCOL, SEND_BODY, SERVER, SWITCHED_PROTOCOL
from ._util import LocalProtocolError, RemoteProtocolError, Sentinel
from ._writers import WRITERS, WritersType
__all__ = ['Connection', 'NEED_DATA', 'PAUSED']

class NEED_DATA(Sentinel, metaclass=Sentinel):
    pass

class PAUSED(Sentinel, metaclass=Sentinel):
    pass
DEFAULT_MAX_INCOMPLETE_EVENT_SIZE = 16 * 1024

def _keep_alive(event: Union[Request, Response]) -> bool:
    connection = get_comma_header(event.headers, b'connection')
    if b'close' in connection:
        return False
    if getattr(event, 'http_version', b'1.1') < b'1.1':
        return False
    return True

def _body_framing(request_method: bytes, event: Union[Request, Response]) -> Tuple[str, Union[Tuple[()], Tuple[int]]]:
    assert type(event) in (Request, Response)
    if type(event) is Response:
        if event.status_code in (204, 304) or request_method == b'HEAD' or (request_method == b'CONNECT' and 200 <= event.status_code < 300):
            return ('content-length', (0,))
        assert event.status_code >= 200
    transfer_encodings = get_comma_header(event.headers, b'transfer-encoding')
    if transfer_encodings:
        assert transfer_encodings == [b'chunked']
        return ('chunked', ())
    content_lengths = get_comma_header(event.headers, b'content-length')
    if content_lengths:
        return ('content-length', (int(content_lengths[0]),))
    if type(event) is Request:
        return ('content-length', (0,))
    else:
        return ('http/1.0', ())

class Connection:

    def __init__(self, our_role: Type[Sentinel], max_incomplete_event_size: int=DEFAULT_MAX_INCOMPLETE_EVENT_SIZE) -> None:
        self._max_incomplete_event_size = max_incomplete_event_size
        if our_role not in (CLIENT, SERVER):
            raise ValueError(f'expected CLIENT or SERVER, not {our_role!r}')
        self.our_role = our_role
        self.their_role: Type[Sentinel]
        if our_role is CLIENT:
            self.their_role = SERVER
        else:
            self.their_role = CLIENT
        self._cstate = ConnectionState()
        self._writer = self._get_io_object(self.our_role, None, WRITERS)
        self._reader = self._get_io_object(self.their_role, None, READERS)
        self._receive_buffer = ReceiveBuffer()
        self._receive_buffer_closed = False
        self.their_http_version: Optional[bytes] = None
        self._request_method: Optional[bytes] = None
        self.client_is_waiting_for_100_continue = False

    @property
    def states(self) -> Dict[Type[Sentinel], Type[Sentinel]]:
        return dict(self._cstate.states)

    @property
    def our_state(self) -> Type[Sentinel]:
        return self._cstate.states[self.our_role]

    @property
    def their_state(self) -> Type[Sentinel]:
        return self._cstate.states[self.their_role]

    @property
    def they_are_waiting_for_100_continue(self) -> bool:
        return self.their_role is CLIENT and self.client_is_waiting_for_100_continue

    def start_next_cycle(self) -> None:
        old_states = dict(self._cstate.states)
        self._cstate.start_next_cycle()
        self._request_method = None
        assert not self.client_is_waiting_for_100_continue
        self._respond_to_state_changes(old_states)

    def _process_error(self, role: Type[Sentinel]) -> None:
        old_states = dict(self._cstate.states)
        self._cstate.process_error(role)
        self._respond_to_state_changes(old_states)

    def _server_switch_event(self, event: Event) -> Optional[Type[Sentinel]]:
        if type(event) is InformationalResponse and event.status_code == 101:
            return _SWITCH_UPGRADE
        if type(event) is Response:
            if _SWITCH_CONNECT in self._cstate.pending_switch_proposals and 200 <= event.status_code < 300:
                return _SWITCH_CONNECT
        return None

    def _process_event(self, role: Type[Sentinel], event: Event) -> None:
        old_states = dict(self._cstate.states)
        if role is CLIENT and type(event) is Request:
            if event.method == b'CONNECT':
                self._cstate.process_client_switch_proposal(_SWITCH_CONNECT)
            if get_comma_header(event.headers, b'upgrade'):
                self._cstate.process_client_switch_proposal(_SWITCH_UPGRADE)
        server_switch_event = None
        if role is SERVER:
            server_switch_event = self._server_switch_event(event)
        self._cstate.process_event(role, type(event), server_switch_event)
        if type(event) is Request:
            self._request_method = event.method
        if role is self.their_role and type(event) in (Request, Response, InformationalResponse):
            event = cast(Union[Request, Response, InformationalResponse], event)
            self.their_http_version = event.http_version
        if type(event) in (Request, Response) and (not _keep_alive(cast(Union[Request, Response], event))):
            self._cstate.process_keep_alive_disabled()
        if type(event) is Request and has_expect_100_continue(event):
            self.client_is_waiting_for_100_continue = True
        if type(event) in (InformationalResponse, Response):
            self.client_is_waiting_for_100_continue = False
        if role is CLIENT and type(event) in (Data, EndOfMessage):
            self.client_is_waiting_for_100_continue = False
        self._respond_to_state_changes(old_states, event)

    def _get_io_object(self, role: Type[Sentinel], event: Optional[Event], io_dict: Union[ReadersType, WritersType]) -> Optional[Callable[..., Any]]:
        state = self._cstate.states[role]
        if state is SEND_BODY:
            framing_type, args = _body_framing(cast(bytes, self._request_method), cast(Union[Request, Response], event))
            return io_dict[SEND_BODY][framing_type](*args)
        else:
            return io_dict.get((role, state))

    def _respond_to_state_changes(self, old_states: Dict[Type[Sentinel], Type[Sentinel]], event: Optional[Event]=None) -> None:
        if self.our_state != old_states[self.our_role]:
            self._writer = self._get_io_object(self.our_role, event, WRITERS)
        if self.their_state != old_states[self.their_role]:
            self._reader = self._get_io_object(self.their_role, event, READERS)

    @property
    def trailing_data(self) -> Tuple[bytes, bool]:
        return (bytes(self._receive_buffer), self._receive_buffer_closed)

    def receive_data(self, data: bytes) -> None:
        if data:
            if self._receive_buffer_closed:
                raise RuntimeError('received close, then received more data?')
            self._receive_buffer += data
        else:
            self._receive_buffer_closed = True

    def _extract_next_receive_event(self) -> Union[Event, Type[NEED_DATA], Type[PAUSED]]:
        state = self.their_state
        if state is DONE and self._receive_buffer:
            return PAUSED
        if state is MIGHT_SWITCH_PROTOCOL or state is SWITCHED_PROTOCOL:
            return PAUSED
        assert self._reader is not None
        event = self._reader(self._receive_buffer)
        if event is None:
            if not self._receive_buffer and self._receive_buffer_closed:
                if hasattr(self._reader, 'read_eof'):
                    event = self._reader.read_eof()
                else:
                    event = ConnectionClosed()
        if event is None:
            event = NEED_DATA
        return event

    def next_event(self) -> Union[Event, Type[NEED_DATA], Type[PAUSED]]:
        if self.their_state is ERROR:
            raise RemoteProtocolError("Can't receive data when peer state is ERROR")
        try:
            event = self._extract_next_receive_event()
            if event not in [NEED_DATA, PAUSED]:
                self._process_event(self.their_role, cast(Event, event))
            if event is NEED_DATA:
                if len(self._receive_buffer) > self._max_incomplete_event_size:
                    raise RemoteProtocolError('Receive buffer too long', error_status_hint=431)
                if self._receive_buffer_closed:
                    raise RemoteProtocolError('peer unexpectedly closed connection')
            return event
        except BaseException as exc:
            self._process_error(self.their_role)
            if isinstance(exc, LocalProtocolError):
                exc._reraise_as_remote_protocol_error()
            else:
                raise

    @overload
    def send(self, event: ConnectionClosed) -> None:
        ...

    @overload
    def send(self, event: Union[Request, InformationalResponse, Response, Data, EndOfMessage]) -> bytes:
        ...

    @overload
    def send(self, event: Event) -> Optional[bytes]:
        ...

    def send(self, event: Event) -> Optional[bytes]:
        data_list = self.send_with_data_passthrough(event)
        if data_list is None:
            return None
        else:
            return b''.join(data_list)

    def send_with_data_passthrough(self, event: Event) -> Optional[List[bytes]]:
        if self.our_state is ERROR:
            raise LocalProtocolError("Can't send data when our state is ERROR")
        try:
            if type(event) is Response:
                event = self._clean_up_response_headers_for_sending(event)
            writer = self._writer
            self._process_event(self.our_role, event)
            if type(event) is ConnectionClosed:
                return None
            else:
                assert writer is not None
                data_list: List[bytes] = []
                writer(event, data_list.append)
                return data_list
        except:
            self._process_error(self.our_role)
            raise

    def send_failed(self) -> None:
        self._process_error(self.our_role)

    def _clean_up_response_headers_for_sending(self, response: Response) -> Response:
        assert type(response) is Response
        headers = response.headers
        need_close = False
        method_for_choosing_headers = cast(bytes, self._request_method)
        if method_for_choosing_headers == b'HEAD':
            method_for_choosing_headers = b'GET'
        framing_type, _ = _body_framing(method_for_choosing_headers, response)
        if framing_type in ('chunked', 'http/1.0'):
            headers = set_comma_header(headers, b'content-length', [])
            if self.their_http_version is None or self.their_http_version < b'1.1':
                headers = set_comma_header(headers, b'transfer-encoding', [])
                if self._request_method != b'HEAD':
                    need_close = True
            else:
                headers = set_comma_header(headers, b'transfer-encoding', [b'chunked'])
        if not self._cstate.keep_alive or need_close:
            connection = set(get_comma_header(headers, b'connection'))
            connection.discard(b'keep-alive')
            connection.add(b'close')
            headers = set_comma_header(headers, b'connection', sorted(connection))
        return Response(headers=headers, status_code=response.status_code, http_version=response.http_version, reason=response.reason)
