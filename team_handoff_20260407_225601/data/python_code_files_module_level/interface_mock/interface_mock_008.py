from typing import cast, Dict, Optional, Set, Tuple, Type, Union
from ._events import *
from ._util import LocalProtocolError, Sentinel
__all__ = ['CLIENT', 'SERVER', 'IDLE', 'SEND_RESPONSE', 'SEND_BODY', 'DONE', 'MUST_CLOSE', 'CLOSED', 'MIGHT_SWITCH_PROTOCOL', 'SWITCHED_PROTOCOL', 'ERROR']

class CLIENT(Sentinel, metaclass=Sentinel):
    pass

class SERVER(Sentinel, metaclass=Sentinel):
    pass

class IDLE(Sentinel, metaclass=Sentinel):
    pass

class SEND_RESPONSE(Sentinel, metaclass=Sentinel):
    pass

class SEND_BODY(Sentinel, metaclass=Sentinel):
    pass

class DONE(Sentinel, metaclass=Sentinel):
    pass

class MUST_CLOSE(Sentinel, metaclass=Sentinel):
    pass

class CLOSED(Sentinel, metaclass=Sentinel):
    pass

class ERROR(Sentinel, metaclass=Sentinel):
    pass

class MIGHT_SWITCH_PROTOCOL(Sentinel, metaclass=Sentinel):
    pass

class SWITCHED_PROTOCOL(Sentinel, metaclass=Sentinel):
    pass

class _SWITCH_UPGRADE(Sentinel, metaclass=Sentinel):
    pass

class _SWITCH_CONNECT(Sentinel, metaclass=Sentinel):
    pass
EventTransitionType = Dict[Type[Sentinel], Dict[Type[Sentinel], Dict[Union[Type[Event], Tuple[Type[Event], Type[Sentinel]]], Type[Sentinel]]]]
EVENT_TRIGGERED_TRANSITIONS: EventTransitionType = {CLIENT: {IDLE: {Request: SEND_BODY, ConnectionClosed: CLOSED}, SEND_BODY: {Data: SEND_BODY, EndOfMessage: DONE}, DONE: {ConnectionClosed: CLOSED}, MUST_CLOSE: {ConnectionClosed: CLOSED}, CLOSED: {ConnectionClosed: CLOSED}, MIGHT_SWITCH_PROTOCOL: {}, SWITCHED_PROTOCOL: {}, ERROR: {}}, SERVER: {IDLE: {ConnectionClosed: CLOSED, Response: SEND_BODY, (Request, CLIENT): SEND_RESPONSE}, SEND_RESPONSE: {InformationalResponse: SEND_RESPONSE, Response: SEND_BODY, (InformationalResponse, _SWITCH_UPGRADE): SWITCHED_PROTOCOL, (Response, _SWITCH_CONNECT): SWITCHED_PROTOCOL}, SEND_BODY: {Data: SEND_BODY, EndOfMessage: DONE}, DONE: {ConnectionClosed: CLOSED}, MUST_CLOSE: {ConnectionClosed: CLOSED}, CLOSED: {ConnectionClosed: CLOSED}, SWITCHED_PROTOCOL: {}, ERROR: {}}}
StateTransitionType = Dict[Tuple[Type[Sentinel], Type[Sentinel]], Dict[Type[Sentinel], Type[Sentinel]]]
STATE_TRIGGERED_TRANSITIONS: StateTransitionType = {(MIGHT_SWITCH_PROTOCOL, SWITCHED_PROTOCOL): {CLIENT: SWITCHED_PROTOCOL}, (CLOSED, DONE): {SERVER: MUST_CLOSE}, (CLOSED, IDLE): {SERVER: MUST_CLOSE}, (ERROR, DONE): {SERVER: MUST_CLOSE}, (DONE, CLOSED): {CLIENT: MUST_CLOSE}, (IDLE, CLOSED): {CLIENT: MUST_CLOSE}, (DONE, ERROR): {CLIENT: MUST_CLOSE}}

class ConnectionState:

    def __init__(self) -> None:
        self.keep_alive = True
        self.pending_switch_proposals: Set[Type[Sentinel]] = set()
        self.states: Dict[Type[Sentinel], Type[Sentinel]] = {CLIENT: IDLE, SERVER: IDLE}

    def process_error(self, role: Type[Sentinel]) -> None:
        self.states[role] = ERROR
        self._fire_state_triggered_transitions()

    def process_keep_alive_disabled(self) -> None:
        self.keep_alive = False
        self._fire_state_triggered_transitions()

    def process_client_switch_proposal(self, switch_event: Type[Sentinel]) -> None:
        self.pending_switch_proposals.add(switch_event)
        self._fire_state_triggered_transitions()

    def process_event(self, role: Type[Sentinel], event_type: Type[Event], server_switch_event: Optional[Type[Sentinel]]=None) -> None:
        _event_type: Union[Type[Event], Tuple[Type[Event], Type[Sentinel]]] = event_type
        if server_switch_event is not None:
            assert role is SERVER
            if server_switch_event not in self.pending_switch_proposals:
                raise LocalProtocolError('Received server _SWITCH_UPGRADE event without a pending proposal')
            _event_type = (event_type, server_switch_event)
        if server_switch_event is None and _event_type is Response:
            self.pending_switch_proposals = set()
        self._fire_event_triggered_transitions(role, _event_type)
        if _event_type is Request:
            assert role is CLIENT
            self._fire_event_triggered_transitions(SERVER, (Request, CLIENT))
        self._fire_state_triggered_transitions()

    def _fire_event_triggered_transitions(self, role: Type[Sentinel], event_type: Union[Type[Event], Tuple[Type[Event], Type[Sentinel]]]) -> None:
        state = self.states[role]
        try:
            new_state = EVENT_TRIGGERED_TRANSITIONS[role][state][event_type]
        except KeyError:
            event_type = cast(Type[Event], event_type)
            raise LocalProtocolError("can't handle event type {} when role={} and state={}".format(event_type.__name__, role, self.states[role])) from None
        self.states[role] = new_state

    def _fire_state_triggered_transitions(self) -> None:
        while True:
            start_states = dict(self.states)
            if self.pending_switch_proposals:
                if self.states[CLIENT] is DONE:
                    self.states[CLIENT] = MIGHT_SWITCH_PROTOCOL
            if not self.pending_switch_proposals:
                if self.states[CLIENT] is MIGHT_SWITCH_PROTOCOL:
                    self.states[CLIENT] = DONE
            if not self.keep_alive:
                for role in (CLIENT, SERVER):
                    if self.states[role] is DONE:
                        self.states[role] = MUST_CLOSE
            joint_state = (self.states[CLIENT], self.states[SERVER])
            changes = STATE_TRIGGERED_TRANSITIONS.get(joint_state, {})
            self.states.update(changes)
            if self.states == start_states:
                return

    def start_next_cycle(self) -> None:
        if self.states != {CLIENT: DONE, SERVER: DONE}:
            raise LocalProtocolError(f'not in a reusable state. self.states={self.states}')
        assert self.keep_alive
        assert not self.pending_switch_proposals
        self.states = {CLIENT: IDLE, SERVER: IDLE}
