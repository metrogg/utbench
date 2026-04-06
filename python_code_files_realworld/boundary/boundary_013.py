import asyncio
import collections
import contextlib
import functools
import itertools
import socket
from typing import List, Optional, Sequence, Set, Union
from . import _staggered
from .types import AddrInfoType, SocketFactoryType

async def start_connection(addr_infos: Sequence[AddrInfoType], *, local_addr_infos: Optional[Sequence[AddrInfoType]]=None, happy_eyeballs_delay: Optional[float]=None, interleave: Optional[int]=None, loop: Optional[asyncio.AbstractEventLoop]=None, socket_factory: Optional[SocketFactoryType]=None) -> socket.socket:
    if not (current_loop := loop):
        current_loop = asyncio.get_running_loop()
    single_addr_info = len(addr_infos) == 1
    if happy_eyeballs_delay is not None and interleave is None:
        interleave = 1
    if interleave and (not single_addr_info):
        addr_infos = _interleave_addrinfos(addr_infos, interleave)
    sock: Optional[socket.socket] = None
    exceptions: List[List[Union[OSError, RuntimeError]]] = []
    if happy_eyeballs_delay is None or single_addr_info:
        for addrinfo in addr_infos:
            try:
                sock = await _connect_sock(current_loop, exceptions, addrinfo, local_addr_infos, None, socket_factory)
                break
            except (RuntimeError, OSError):
                continue
    else:
        open_sockets: Set[socket.socket] = set()
        try:
            sock, _, _ = await _staggered.staggered_race((functools.partial(_connect_sock, current_loop, exceptions, addrinfo, local_addr_infos, open_sockets, socket_factory) for addrinfo in addr_infos), happy_eyeballs_delay)
        finally:
            for s in open_sockets:
                if s is not sock:
                    with contextlib.suppress(OSError):
                        s.close()
            open_sockets = None
    if sock is None:
        all_exceptions = [exc for sub in exceptions for exc in sub]
        try:
            first_exception = all_exceptions[0]
            if len(all_exceptions) == 1:
                raise first_exception
            else:
                model = str(first_exception)
                if all((str(exc) == model for exc in all_exceptions)):
                    raise first_exception
                msg = 'Multiple exceptions: {}'.format(', '.join((str(exc) for exc in all_exceptions)))
                if isinstance(first_exception, OSError):
                    first_errno = first_exception.errno
                    if all((isinstance(exc, OSError) and exc.errno == first_errno for exc in all_exceptions)):
                        raise OSError(first_errno, msg)
                elif isinstance(first_exception, RuntimeError) and all((isinstance(exc, RuntimeError) for exc in all_exceptions)):
                    raise RuntimeError(msg)
                raise OSError(msg)
        finally:
            all_exceptions = None
            exceptions = None
    return sock

async def _connect_sock(loop: asyncio.AbstractEventLoop, exceptions: List[List[Union[OSError, RuntimeError]]], addr_info: AddrInfoType, local_addr_infos: Optional[Sequence[AddrInfoType]]=None, open_sockets: Optional[Set[socket.socket]]=None, socket_factory: Optional[SocketFactoryType]=None) -> socket.socket:
    my_exceptions: List[Union[OSError, RuntimeError]] = []
    exceptions.append(my_exceptions)
    family, type_, proto, _, address = addr_info
    sock = None
    try:
        if socket_factory is not None:
            sock = socket_factory(addr_info)
        else:
            sock = socket.socket(family=family, type=type_, proto=proto)
        if open_sockets is not None:
            open_sockets.add(sock)
        sock.setblocking(False)
        if local_addr_infos is not None:
            for lfamily, _, _, _, laddr in local_addr_infos:
                if lfamily != family:
                    continue
                try:
                    sock.bind(laddr)
                    break
                except OSError as exc:
                    msg = f"error while attempting to bind on address {laddr!r}: {(exc.strerror or '').lower()}"
                    exc = OSError(exc.errno, msg)
                    my_exceptions.append(exc)
            else:
                if my_exceptions:
                    raise my_exceptions.pop()
                else:
                    raise OSError(f'no matching local address with family={family!r} found')
        await loop.sock_connect(sock, address)
        return sock
    except (RuntimeError, OSError) as exc:
        my_exceptions.append(exc)
        if sock is not None:
            if open_sockets is not None:
                open_sockets.remove(sock)
            try:
                sock.close()
            except OSError as e:
                my_exceptions.append(e)
                raise
        raise
    except:
        if sock is not None:
            if open_sockets is not None:
                open_sockets.remove(sock)
            try:
                sock.close()
            except OSError as e:
                my_exceptions.append(e)
                raise
        raise
    finally:
        exceptions = my_exceptions = None

def _interleave_addrinfos(addrinfos: Sequence[AddrInfoType], first_address_family_count: int=1) -> List[AddrInfoType]:
    addrinfos_by_family: collections.OrderedDict[int, List[AddrInfoType]] = collections.OrderedDict()
    for addr in addrinfos:
        family = addr[0]
        if family not in addrinfos_by_family:
            addrinfos_by_family[family] = []
        addrinfos_by_family[family].append(addr)
    addrinfos_lists = list(addrinfos_by_family.values())
    reordered: List[AddrInfoType] = []
    if first_address_family_count > 1:
        reordered.extend(addrinfos_lists[0][:first_address_family_count - 1])
        del addrinfos_lists[0][:first_address_family_count - 1]
    reordered.extend((a for a in itertools.chain.from_iterable(itertools.zip_longest(*addrinfos_lists)) if a is not None))
    return reordered
