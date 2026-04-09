from __future__ import annotations
import errno
import os
import socket
import ssl
import stat
import sys
from collections.abc import Awaitable
from dataclasses import dataclass
from ipaddress import IPv4Address, IPv6Address, ip_address
from os import PathLike, chmod
from socket import AddressFamily, SocketKind
from typing import TYPE_CHECKING, Any, Literal, cast, overload
from .. import ConnectionFailed, to_thread
from ..abc import ByteStreamConnectable, ConnectedUDPSocket, ConnectedUNIXDatagramSocket, IPAddressType, IPSockAddrType, SocketListener, SocketStream, UDPSocket, UNIXDatagramSocket, UNIXSocketStream
from ..streams.stapled import MultiListener
from ..streams.tls import TLSConnectable, TLSStream
from ._eventloop import get_async_backend
from ._resources import aclose_forcefully
from ._synchronization import Event
from ._tasks import create_task_group, move_on_after
if TYPE_CHECKING:
    from _typeshed import FileDescriptorLike
else:
    FileDescriptorLike = object
if sys.version_info < (3, 11):
    from exceptiongroup import ExceptionGroup
if sys.version_info >= (3, 12):
    from typing import override
else:
    from typing_extensions import override
if sys.version_info < (3, 13):
    from typing_extensions import deprecated
else:
    from warnings import deprecated
IPPROTO_IPV6 = getattr(socket, 'IPPROTO_IPV6', 41)
AnyIPAddressFamily = Literal[AddressFamily.AF_UNSPEC, AddressFamily.AF_INET, AddressFamily.AF_INET6]
IPAddressFamily = Literal[AddressFamily.AF_INET, AddressFamily.AF_INET6]

@overload
async def connect_tcp(remote_host: IPAddressType, remote_port: int, *, local_host: IPAddressType | None=..., ssl_context: ssl.SSLContext | None=..., tls_standard_compatible: bool=..., tls_hostname: str, happy_eyeballs_delay: float=...) -> TLSStream:
    ...

@overload
async def connect_tcp(remote_host: IPAddressType, remote_port: int, *, local_host: IPAddressType | None=..., ssl_context: ssl.SSLContext, tls_standard_compatible: bool=..., tls_hostname: str | None=..., happy_eyeballs_delay: float=...) -> TLSStream:
    ...

@overload
async def connect_tcp(remote_host: IPAddressType, remote_port: int, *, local_host: IPAddressType | None=..., tls: Literal[True], ssl_context: ssl.SSLContext | None=..., tls_standard_compatible: bool=..., tls_hostname: str | None=..., happy_eyeballs_delay: float=...) -> TLSStream:
    ...

@overload
async def connect_tcp(remote_host: IPAddressType, remote_port: int, *, local_host: IPAddressType | None=..., tls: Literal[False], ssl_context: ssl.SSLContext | None=..., tls_standard_compatible: bool=..., tls_hostname: str | None=..., happy_eyeballs_delay: float=...) -> SocketStream:
    ...

@overload
async def connect_tcp(remote_host: IPAddressType, remote_port: int, *, local_host: IPAddressType | None=..., happy_eyeballs_delay: float=...) -> SocketStream:
    ...

async def connect_tcp(remote_host: IPAddressType, remote_port: int, *, local_host: IPAddressType | None=None, tls: bool=False, ssl_context: ssl.SSLContext | None=None, tls_standard_compatible: bool=True, tls_hostname: str | None=None, happy_eyeballs_delay: float=0.25) -> SocketStream | TLSStream:
    connected_stream: SocketStream | None = None

    async def try_connect(remote_host: str, event: Event) -> None:
        nonlocal connected_stream
        try:
            stream = await asynclib.connect_tcp(remote_host, remote_port, local_address)
        except OSError as exc:
            oserrors.append(exc)
            return
        else:
            if connected_stream is None:
                connected_stream = stream
                tg.cancel_scope.cancel()
            else:
                await stream.aclose()
        finally:
            event.set()
    asynclib = get_async_backend()
    local_address: IPSockAddrType | None = None
    family = socket.AF_UNSPEC
    if local_host:
        gai_res = await getaddrinfo(str(local_host), None)
        family, *_, local_address = gai_res[0]
    target_host = str(remote_host)
    try:
        addr_obj = ip_address(remote_host)
    except ValueError:
        addr_obj = None
    if addr_obj is not None:
        if isinstance(addr_obj, IPv6Address):
            target_addrs = [(socket.AF_INET6, addr_obj.compressed)]
        else:
            target_addrs = [(socket.AF_INET, addr_obj.compressed)]
    else:
        gai_res = await getaddrinfo(target_host, remote_port, family=family, type=socket.SOCK_STREAM)
        v6_found = v4_found = False
        target_addrs = []
        for af, *_, sa in gai_res:
            if af == socket.AF_INET6 and (not v6_found):
                v6_found = True
                target_addrs.insert(0, (af, sa[0]))
            elif af == socket.AF_INET and (not v4_found) and v6_found:
                v4_found = True
                target_addrs.insert(1, (af, sa[0]))
            else:
                target_addrs.append((af, sa[0]))
    oserrors: list[OSError] = []
    try:
        async with create_task_group() as tg:
            for _af, addr in target_addrs:
                event = Event()
                tg.start_soon(try_connect, addr, event)
                with move_on_after(happy_eyeballs_delay):
                    await event.wait()
        if connected_stream is None:
            cause = oserrors[0] if len(oserrors) == 1 else ExceptionGroup('multiple connection attempts failed', oserrors)
            raise OSError('All connection attempts failed') from cause
    finally:
        oserrors.clear()
    if tls or tls_hostname or ssl_context:
        try:
            return await TLSStream.wrap(connected_stream, server_side=False, hostname=tls_hostname or str(remote_host), ssl_context=ssl_context, standard_compatible=tls_standard_compatible)
        except BaseException:
            await aclose_forcefully(connected_stream)
            raise
    return connected_stream

async def connect_unix(path: str | bytes | PathLike[Any]) -> UNIXSocketStream:
    path = os.fspath(path)
    return await get_async_backend().connect_unix(path)

async def create_tcp_listener(*, local_host: IPAddressType | None=None, local_port: int=0, family: AnyIPAddressFamily=socket.AddressFamily.AF_UNSPEC, backlog: int=65536, reuse_port: bool=False) -> MultiListener[SocketStream]:
    asynclib = get_async_backend()
    backlog = min(backlog, 65536)
    local_host = str(local_host) if local_host is not None else None

    def setup_raw_socket(fam: AddressFamily, bind_addr: tuple[str, int] | tuple[str, int, int, int], *, v6only: bool=True) -> socket.socket:
        sock = socket.socket(fam)
        try:
            sock.setblocking(False)
            if fam == AddressFamily.AF_INET6:
                sock.setsockopt(IPPROTO_IPV6, socket.IPV6_V6ONLY, v6only)
            if sys.platform == 'win32':
                sock.setsockopt(socket.SOL_SOCKET, socket.SO_EXCLUSIVEADDRUSE, 1)
            else:
                sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            if reuse_port:
                sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEPORT, 1)
            if fam == socket.AF_INET6 and '%' in bind_addr[0]:
                addr, scope_id = bind_addr[0].split('%', 1)
                bind_addr = (addr, bind_addr[1], 0, int(scope_id))
            sock.bind(bind_addr)
            sock.listen(backlog)
        except BaseException:
            sock.close()
            raise
        return sock
    gai_res = await getaddrinfo(local_host, local_port, family=family, type=socket.SOCK_STREAM if sys.platform == 'win32' else 0, flags=socket.AI_PASSIVE | socket.AI_ADDRCONFIG)
    sockaddrs = sorted({res for res in gai_res if res[1] == SocketKind.SOCK_STREAM})
    if local_host is None and family == AddressFamily.AF_UNSPEC and socket.has_dualstack_ipv6() and any((fam == AddressFamily.AF_INET6 for fam, *_ in gai_res)):
        raw_socket = setup_raw_socket(AddressFamily.AF_INET6, ('::', local_port), v6only=False)
        listener = asynclib.create_tcp_listener(raw_socket)
        return MultiListener([listener])
    errors: list[OSError] = []
    try:
        for _ in range(len(sockaddrs)):
            listeners: list[SocketListener] = []
            bound_ephemeral_port = local_port
            try:
                for fam, *_, sockaddr in sockaddrs:
                    sockaddr = (sockaddr[0], bound_ephemeral_port, *sockaddr[2:])
                    raw_socket = setup_raw_socket(fam, sockaddr)
                    if local_port == 0 and len(gai_res) > 1:
                        bound_ephemeral_port = raw_socket.getsockname()[1]
                    listeners.append(asynclib.create_tcp_listener(raw_socket))
            except BaseException as exc:
                for listener in listeners:
                    await listener.aclose()
                if isinstance(exc, OSError) and exc.errno == errno.EADDRINUSE and (local_port == 0) and bound_ephemeral_port:
                    errors.append(exc)
                    sockaddrs.append(sockaddrs.pop(0))
                    continue
                raise
            return MultiListener(listeners)
        raise OSError(f'Could not create {len(sockaddrs)} listeners with a consistent port') from ExceptionGroup('Several bind attempts failed', errors)
    finally:
        del errors

async def create_unix_listener(path: str | bytes | PathLike[Any], *, mode: int | None=None, backlog: int=65536) -> SocketListener:
    backlog = min(backlog, 65536)
    raw_socket = await setup_unix_local_socket(path, mode, socket.SOCK_STREAM)
    try:
        raw_socket.listen(backlog)
        return get_async_backend().create_unix_listener(raw_socket)
    except BaseException:
        raw_socket.close()
        raise

async def create_udp_socket(family: AnyIPAddressFamily=AddressFamily.AF_UNSPEC, *, local_host: IPAddressType | None=None, local_port: int=0, reuse_port: bool=False) -> UDPSocket:
    if family is AddressFamily.AF_UNSPEC and (not local_host):
        raise ValueError('Either "family" or "local_host" must be given')
    if local_host:
        gai_res = await getaddrinfo(str(local_host), local_port, family=family, type=socket.SOCK_DGRAM, flags=socket.AI_PASSIVE | socket.AI_ADDRCONFIG)
        family = cast(AnyIPAddressFamily, gai_res[0][0])
        local_address = gai_res[0][-1]
    elif family is AddressFamily.AF_INET6:
        local_address = ('::', 0)
    else:
        local_address = ('0.0.0.0', 0)
    sock = await get_async_backend().create_udp_socket(family, local_address, None, reuse_port)
    return cast(UDPSocket, sock)

async def create_connected_udp_socket(remote_host: IPAddressType, remote_port: int, *, family: AnyIPAddressFamily=AddressFamily.AF_UNSPEC, local_host: IPAddressType | None=None, local_port: int=0, reuse_port: bool=False) -> ConnectedUDPSocket:
    local_address = None
    if local_host:
        gai_res = await getaddrinfo(str(local_host), local_port, family=family, type=socket.SOCK_DGRAM, flags=socket.AI_PASSIVE | socket.AI_ADDRCONFIG)
        family = cast(AnyIPAddressFamily, gai_res[0][0])
        local_address = gai_res[0][-1]
    gai_res = await getaddrinfo(str(remote_host), remote_port, family=family, type=socket.SOCK_DGRAM)
    family = cast(AnyIPAddressFamily, gai_res[0][0])
    remote_address = gai_res[0][-1]
    sock = await get_async_backend().create_udp_socket(family, local_address, remote_address, reuse_port)
    return cast(ConnectedUDPSocket, sock)

async def create_unix_datagram_socket(*, local_path: None | str | bytes | PathLike[Any]=None, local_mode: int | None=None) -> UNIXDatagramSocket:
    raw_socket = await setup_unix_local_socket(local_path, local_mode, socket.SOCK_DGRAM)
    return await get_async_backend().create_unix_datagram_socket(raw_socket, None)

async def create_connected_unix_datagram_socket(remote_path: str | bytes | PathLike[Any], *, local_path: None | str | bytes | PathLike[Any]=None, local_mode: int | None=None) -> ConnectedUNIXDatagramSocket:
    remote_path = os.fspath(remote_path)
    raw_socket = await setup_unix_local_socket(local_path, local_mode, socket.SOCK_DGRAM)
    return await get_async_backend().create_unix_datagram_socket(raw_socket, remote_path)

async def getaddrinfo(host: bytes | str | None, port: str | int | None, *, family: int | AddressFamily=0, type: int | SocketKind=0, proto: int=0, flags: int=0) -> list[tuple[AddressFamily, SocketKind, int, str, tuple[str, int]]]:
    if isinstance(host, str):
        try:
            encoded_host: bytes | None = host.encode('ascii')
        except UnicodeEncodeError:
            import idna
            encoded_host = idna.encode(host, uts46=True)
    else:
        encoded_host = host
    gai_res = await get_async_backend().getaddrinfo(encoded_host, port, family=family, type=type, proto=proto, flags=flags)
    return [(family, type, proto, canonname, convert_ipv6_sockaddr(sockaddr)) for family, type, proto, canonname, sockaddr in gai_res if not isinstance(sockaddr[0], int)]

def getnameinfo(sockaddr: IPSockAddrType, flags: int=0) -> Awaitable[tuple[str, str]]:
    return get_async_backend().getnameinfo(sockaddr, flags)

@deprecated('This function is deprecated; use `wait_readable` instead')
def wait_socket_readable(sock: socket.socket) -> Awaitable[None]:
    return get_async_backend().wait_readable(sock.fileno())

@deprecated('This function is deprecated; use `wait_writable` instead')
def wait_socket_writable(sock: socket.socket) -> Awaitable[None]:
    return get_async_backend().wait_writable(sock.fileno())

def wait_readable(obj: FileDescriptorLike) -> Awaitable[None]:
    return get_async_backend().wait_readable(obj)

def wait_writable(obj: FileDescriptorLike) -> Awaitable[None]:
    return get_async_backend().wait_writable(obj)

def notify_closing(obj: FileDescriptorLike) -> None:
    get_async_backend().notify_closing(obj)

def convert_ipv6_sockaddr(sockaddr: tuple[str, int, int, int] | tuple[str, int]) -> tuple[str, int]:
    if isinstance(sockaddr, tuple) and len(sockaddr) == 4:
        host, port, flowinfo, scope_id = sockaddr
        if scope_id:
            host = host.split('%')[0]
            return (f'{host}%{scope_id}', port)
        else:
            return (host, port)
    else:
        return sockaddr

async def setup_unix_local_socket(path: None | str | bytes | PathLike[Any], mode: int | None, socktype: int) -> socket.socket:
    path_str: str | None
    if path is not None:
        path_str = os.fsdecode(path)
        if not path_str.startswith('\x00'):
            try:
                stat_result = os.stat(path)
            except OSError as e:
                if e.errno not in (errno.ENOENT, errno.ENOTDIR, errno.EBADF, errno.ELOOP):
                    raise
            else:
                if stat.S_ISSOCK(stat_result.st_mode):
                    os.unlink(path)
    else:
        path_str = None
    raw_socket = socket.socket(socket.AF_UNIX, socktype)
    raw_socket.setblocking(False)
    if path_str is not None:
        try:
            await to_thread.run_sync(raw_socket.bind, path_str, abandon_on_cancel=True)
            if mode is not None:
                await to_thread.run_sync(chmod, path_str, mode, abandon_on_cancel=True)
        except BaseException:
            raw_socket.close()
            raise
    return raw_socket

@dataclass
class TCPConnectable(ByteStreamConnectable):
    host: str | IPv4Address | IPv6Address
    port: int

    def __post_init__(self) -> None:
        if self.port < 1 or self.port > 65535:
            raise ValueError('TCP port number out of range')

    @override
    async def connect(self) -> SocketStream:
        try:
            return await connect_tcp(self.host, self.port)
        except OSError as exc:
            raise ConnectionFailed(f'error connecting to {self.host}:{self.port}: {exc}') from exc

@dataclass
class UNIXConnectable(ByteStreamConnectable):
    path: str | bytes | PathLike[str] | PathLike[bytes]

    @override
    async def connect(self) -> UNIXSocketStream:
        try:
            return await connect_unix(self.path)
        except OSError as exc:
            raise ConnectionFailed(f'error connecting to {self.path!r}: {exc}') from exc

def as_connectable(remote: ByteStreamConnectable | tuple[str | IPv4Address | IPv6Address, int] | str | bytes | PathLike[str], /, *, tls: bool=False, ssl_context: ssl.SSLContext | None=None, tls_hostname: str | None=None, tls_standard_compatible: bool=True) -> ByteStreamConnectable:
    connectable: TCPConnectable | UNIXConnectable | TLSConnectable
    if isinstance(remote, ByteStreamConnectable):
        return remote
    elif isinstance(remote, tuple) and len(remote) == 2:
        connectable = TCPConnectable(*remote)
    elif isinstance(remote, (str, bytes, PathLike)):
        connectable = UNIXConnectable(remote)
    else:
        raise TypeError(f'cannot convert {remote!r} to a connectable')
    if tls:
        if not tls_hostname and isinstance(connectable, TCPConnectable):
            tls_hostname = str(connectable.host)
        connectable = TLSConnectable(connectable, ssl_context=ssl_context, hostname=tls_hostname, standard_compatible=tls_standard_compatible)
    return connectable
