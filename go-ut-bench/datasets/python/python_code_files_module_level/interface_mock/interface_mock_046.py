from __future__ import annotations
import logging
import ssl
import socksio
from .._backends.sync import SyncBackend
from .._backends.base import NetworkBackend, NetworkStream
from .._exceptions import ConnectionNotAvailable, ProxyError
from .._models import URL, Origin, Request, Response, enforce_bytes, enforce_url
from .._ssl import default_ssl_context
from .._synchronization import Lock
from .._trace import Trace
from .connection_pool import ConnectionPool
from .http11 import HTTP11Connection
from .interfaces import ConnectionInterface
logger = logging.getLogger('httpcore.socks')
AUTH_METHODS = {b'\x00': 'NO AUTHENTICATION REQUIRED', b'\x01': 'GSSAPI', b'\x02': 'USERNAME/PASSWORD', b'\xff': 'NO ACCEPTABLE METHODS'}
REPLY_CODES = {b'\x00': 'Succeeded', b'\x01': 'General SOCKS server failure', b'\x02': 'Connection not allowed by ruleset', b'\x03': 'Network unreachable', b'\x04': 'Host unreachable', b'\x05': 'Connection refused', b'\x06': 'TTL expired', b'\x07': 'Command not supported', b'\x08': 'Address type not supported'}

def _init_socks5_connection(stream: NetworkStream, *, host: bytes, port: int, auth: tuple[bytes, bytes] | None=None) -> None:
    conn = socksio.socks5.SOCKS5Connection()
    auth_method = socksio.socks5.SOCKS5AuthMethod.NO_AUTH_REQUIRED if auth is None else socksio.socks5.SOCKS5AuthMethod.USERNAME_PASSWORD
    conn.send(socksio.socks5.SOCKS5AuthMethodsRequest([auth_method]))
    outgoing_bytes = conn.data_to_send()
    stream.write(outgoing_bytes)
    incoming_bytes = stream.read(max_bytes=4096)
    response = conn.receive_data(incoming_bytes)
    assert isinstance(response, socksio.socks5.SOCKS5AuthReply)
    if response.method != auth_method:
        requested = AUTH_METHODS.get(auth_method, 'UNKNOWN')
        responded = AUTH_METHODS.get(response.method, 'UNKNOWN')
        raise ProxyError(f'Requested {requested} from proxy server, but got {responded}.')
    if response.method == socksio.socks5.SOCKS5AuthMethod.USERNAME_PASSWORD:
        assert auth is not None
        username, password = auth
        conn.send(socksio.socks5.SOCKS5UsernamePasswordRequest(username, password))
        outgoing_bytes = conn.data_to_send()
        stream.write(outgoing_bytes)
        incoming_bytes = stream.read(max_bytes=4096)
        response = conn.receive_data(incoming_bytes)
        assert isinstance(response, socksio.socks5.SOCKS5UsernamePasswordReply)
        if not response.success:
            raise ProxyError('Invalid username/password')
    conn.send(socksio.socks5.SOCKS5CommandRequest.from_address(socksio.socks5.SOCKS5Command.CONNECT, (host, port)))
    outgoing_bytes = conn.data_to_send()
    stream.write(outgoing_bytes)
    incoming_bytes = stream.read(max_bytes=4096)
    response = conn.receive_data(incoming_bytes)
    assert isinstance(response, socksio.socks5.SOCKS5Reply)
    if response.reply_code != socksio.socks5.SOCKS5ReplyCode.SUCCEEDED:
        reply_code = REPLY_CODES.get(response.reply_code, 'UNKOWN')
        raise ProxyError(f'Proxy Server could not connect: {reply_code}.')

class SOCKSProxy(ConnectionPool):

    def __init__(self, proxy_url: URL | bytes | str, proxy_auth: tuple[bytes | str, bytes | str] | None=None, ssl_context: ssl.SSLContext | None=None, max_connections: int | None=10, max_keepalive_connections: int | None=None, keepalive_expiry: float | None=None, http1: bool=True, http2: bool=False, retries: int=0, network_backend: NetworkBackend | None=None) -> None:
        super().__init__(ssl_context=ssl_context, max_connections=max_connections, max_keepalive_connections=max_keepalive_connections, keepalive_expiry=keepalive_expiry, http1=http1, http2=http2, network_backend=network_backend, retries=retries)
        self._ssl_context = ssl_context
        self._proxy_url = enforce_url(proxy_url, name='proxy_url')
        if proxy_auth is not None:
            username, password = proxy_auth
            username_bytes = enforce_bytes(username, name='proxy_auth')
            password_bytes = enforce_bytes(password, name='proxy_auth')
            self._proxy_auth: tuple[bytes, bytes] | None = (username_bytes, password_bytes)
        else:
            self._proxy_auth = None

    def create_connection(self, origin: Origin) -> ConnectionInterface:
        return Socks5Connection(proxy_origin=self._proxy_url.origin, remote_origin=origin, proxy_auth=self._proxy_auth, ssl_context=self._ssl_context, keepalive_expiry=self._keepalive_expiry, http1=self._http1, http2=self._http2, network_backend=self._network_backend)

class Socks5Connection(ConnectionInterface):

    def __init__(self, proxy_origin: Origin, remote_origin: Origin, proxy_auth: tuple[bytes, bytes] | None=None, ssl_context: ssl.SSLContext | None=None, keepalive_expiry: float | None=None, http1: bool=True, http2: bool=False, network_backend: NetworkBackend | None=None) -> None:
        self._proxy_origin = proxy_origin
        self._remote_origin = remote_origin
        self._proxy_auth = proxy_auth
        self._ssl_context = ssl_context
        self._keepalive_expiry = keepalive_expiry
        self._http1 = http1
        self._http2 = http2
        self._network_backend: NetworkBackend = SyncBackend() if network_backend is None else network_backend
        self._connect_lock = Lock()
        self._connection: ConnectionInterface | None = None
        self._connect_failed = False

    def handle_request(self, request: Request) -> Response:
        timeouts = request.extensions.get('timeout', {})
        sni_hostname = request.extensions.get('sni_hostname', None)
        timeout = timeouts.get('connect', None)
        with self._connect_lock:
            if self._connection is None:
                try:
                    kwargs = {'host': self._proxy_origin.host.decode('ascii'), 'port': self._proxy_origin.port, 'timeout': timeout}
                    with Trace('connect_tcp', logger, request, kwargs) as trace:
                        stream = self._network_backend.connect_tcp(**kwargs)
                        trace.return_value = stream
                    kwargs = {'stream': stream, 'host': self._remote_origin.host.decode('ascii'), 'port': self._remote_origin.port, 'auth': self._proxy_auth}
                    with Trace('setup_socks5_connection', logger, request, kwargs) as trace:
                        _init_socks5_connection(**kwargs)
                        trace.return_value = stream
                    if self._remote_origin.scheme == b'https':
                        ssl_context = default_ssl_context() if self._ssl_context is None else self._ssl_context
                        alpn_protocols = ['http/1.1', 'h2'] if self._http2 else ['http/1.1']
                        ssl_context.set_alpn_protocols(alpn_protocols)
                        kwargs = {'ssl_context': ssl_context, 'server_hostname': sni_hostname or self._remote_origin.host.decode('ascii'), 'timeout': timeout}
                        with Trace('start_tls', logger, request, kwargs) as trace:
                            stream = stream.start_tls(**kwargs)
                            trace.return_value = stream
                    ssl_object = stream.get_extra_info('ssl_object')
                    http2_negotiated = ssl_object is not None and ssl_object.selected_alpn_protocol() == 'h2'
                    if http2_negotiated or (self._http2 and (not self._http1)):
                        from .http2 import HTTP2Connection
                        self._connection = HTTP2Connection(origin=self._remote_origin, stream=stream, keepalive_expiry=self._keepalive_expiry)
                    else:
                        self._connection = HTTP11Connection(origin=self._remote_origin, stream=stream, keepalive_expiry=self._keepalive_expiry)
                except Exception as exc:
                    self._connect_failed = True
                    raise exc
            elif not self._connection.is_available():
                raise ConnectionNotAvailable()
        return self._connection.handle_request(request)

    def can_handle_request(self, origin: Origin) -> bool:
        return origin == self._remote_origin

    def close(self) -> None:
        if self._connection is not None:
            self._connection.close()

    def is_available(self) -> bool:
        if self._connection is None:
            return self._http2 and (self._remote_origin.scheme == b'https' or not self._http1) and (not self._connect_failed)
        return self._connection.is_available()

    def has_expired(self) -> bool:
        if self._connection is None:
            return self._connect_failed
        return self._connection.has_expired()

    def is_idle(self) -> bool:
        if self._connection is None:
            return self._connect_failed
        return self._connection.is_idle()

    def is_closed(self) -> bool:
        if self._connection is None:
            return self._connect_failed
        return self._connection.is_closed()

    def info(self) -> str:
        if self._connection is None:
            return 'CONNECTION FAILED' if self._connect_failed else 'CONNECTING'
        return self._connection.info()

    def __repr__(self) -> str:
        return f'<{self.__class__.__name__} [{self.info()}]>'
