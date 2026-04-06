from __future__ import annotations
__all__ = ('TLSAttribute', 'TLSConnectable', 'TLSListener', 'TLSStream')
import logging
import re
import ssl
import sys
from collections.abc import Callable, Mapping
from dataclasses import dataclass
from functools import wraps
from ssl import SSLContext
from typing import Any, TypeAlias, TypeVar
from .. import BrokenResourceError, EndOfStream, aclose_forcefully, get_cancelled_exc_class, to_thread
from .._core._typedattr import TypedAttributeSet, typed_attribute
from ..abc import AnyByteStream, AnyByteStreamConnectable, ByteStream, ByteStreamConnectable, Listener, TaskGroup
if sys.version_info >= (3, 11):
    from typing import TypeVarTuple, Unpack
else:
    from typing_extensions import TypeVarTuple, Unpack
if sys.version_info >= (3, 12):
    from typing import override
else:
    from typing_extensions import override
T_Retval = TypeVar('T_Retval')
PosArgsT = TypeVarTuple('PosArgsT')
_PCTRTT: TypeAlias = tuple[tuple[str, str], ...]
_PCTRTTT: TypeAlias = tuple[_PCTRTT, ...]

class TLSAttribute(TypedAttributeSet):
    alpn_protocol: str | None = typed_attribute()
    channel_binding_tls_unique: bytes = typed_attribute()
    cipher: tuple[str, str, int] = typed_attribute()
    peer_certificate: None | dict[str, str | _PCTRTTT | _PCTRTT] = typed_attribute()
    peer_certificate_binary: bytes | None = typed_attribute()
    server_side: bool = typed_attribute()
    shared_ciphers: list[tuple[str, str, int]] | None = typed_attribute()
    ssl_object: ssl.SSLObject = typed_attribute()
    standard_compatible: bool = typed_attribute()
    tls_version: str = typed_attribute()

@dataclass(eq=False)
class TLSStream(ByteStream):
    transport_stream: AnyByteStream
    standard_compatible: bool
    _ssl_object: ssl.SSLObject
    _read_bio: ssl.MemoryBIO
    _write_bio: ssl.MemoryBIO

    @classmethod
    async def wrap(cls, transport_stream: AnyByteStream, *, server_side: bool | None=None, hostname: str | None=None, ssl_context: ssl.SSLContext | None=None, standard_compatible: bool=True) -> TLSStream:
        if server_side is None:
            server_side = not hostname
        if not ssl_context:
            purpose = ssl.Purpose.CLIENT_AUTH if server_side else ssl.Purpose.SERVER_AUTH
            ssl_context = ssl.create_default_context(purpose)
            if hasattr(ssl, 'OP_IGNORE_UNEXPECTED_EOF'):
                ssl_context.options &= ~ssl.OP_IGNORE_UNEXPECTED_EOF
        bio_in = ssl.MemoryBIO()
        bio_out = ssl.MemoryBIO()
        if type(ssl_context) is ssl.SSLContext:
            ssl_object = ssl_context.wrap_bio(bio_in, bio_out, server_side=server_side, server_hostname=hostname)
        else:
            ssl_object = await to_thread.run_sync(ssl_context.wrap_bio, bio_in, bio_out, server_side, hostname, None)
        wrapper = cls(transport_stream=transport_stream, standard_compatible=standard_compatible, _ssl_object=ssl_object, _read_bio=bio_in, _write_bio=bio_out)
        await wrapper._call_sslobject_method(ssl_object.do_handshake)
        return wrapper

    async def _call_sslobject_method(self, func: Callable[[Unpack[PosArgsT]], T_Retval], *args: Unpack[PosArgsT]) -> T_Retval:
        while True:
            try:
                result = func(*args)
            except ssl.SSLWantReadError:
                try:
                    if self._write_bio.pending:
                        await self.transport_stream.send(self._write_bio.read())
                    data = await self.transport_stream.receive()
                except EndOfStream:
                    self._read_bio.write_eof()
                except OSError as exc:
                    self._read_bio.write_eof()
                    self._write_bio.write_eof()
                    raise BrokenResourceError from exc
                else:
                    self._read_bio.write(data)
            except ssl.SSLWantWriteError:
                await self.transport_stream.send(self._write_bio.read())
            except ssl.SSLSyscallError as exc:
                self._read_bio.write_eof()
                self._write_bio.write_eof()
                raise BrokenResourceError from exc
            except ssl.SSLError as exc:
                self._read_bio.write_eof()
                self._write_bio.write_eof()
                if isinstance(exc, ssl.SSLEOFError) or (exc.strerror and 'UNEXPECTED_EOF_WHILE_READING' in exc.strerror):
                    if self.standard_compatible:
                        raise BrokenResourceError from exc
                    else:
                        raise EndOfStream from None
                raise
            else:
                if self._write_bio.pending:
                    await self.transport_stream.send(self._write_bio.read())
                return result

    async def unwrap(self) -> tuple[AnyByteStream, bytes]:
        await self._call_sslobject_method(self._ssl_object.unwrap)
        self._read_bio.write_eof()
        self._write_bio.write_eof()
        return (self.transport_stream, self._read_bio.read())

    async def aclose(self) -> None:
        if self.standard_compatible:
            try:
                await self.unwrap()
            except BaseException:
                await aclose_forcefully(self.transport_stream)
                raise
        await self.transport_stream.aclose()

    async def receive(self, max_bytes: int=65536) -> bytes:
        data = await self._call_sslobject_method(self._ssl_object.read, max_bytes)
        if not data:
            raise EndOfStream
        return data

    async def send(self, item: bytes) -> None:
        await self._call_sslobject_method(self._ssl_object.write, item)

    async def send_eof(self) -> None:
        tls_version = self.extra(TLSAttribute.tls_version)
        match = re.match('TLSv(\\d+)(?:\\.(\\d+))?', tls_version)
        if match:
            major, minor = (int(match.group(1)), int(match.group(2) or 0))
            if (major, minor) < (1, 3):
                raise NotImplementedError(f'send_eof() requires at least TLSv1.3; current session uses {tls_version}')
        raise NotImplementedError('send_eof() has not yet been implemented for TLS streams')

    @property
    def extra_attributes(self) -> Mapping[Any, Callable[[], Any]]:
        return {**self.transport_stream.extra_attributes, TLSAttribute.alpn_protocol: self._ssl_object.selected_alpn_protocol, TLSAttribute.channel_binding_tls_unique: self._ssl_object.get_channel_binding, TLSAttribute.cipher: self._ssl_object.cipher, TLSAttribute.peer_certificate: lambda: self._ssl_object.getpeercert(False), TLSAttribute.peer_certificate_binary: lambda: self._ssl_object.getpeercert(True), TLSAttribute.server_side: lambda: self._ssl_object.server_side, TLSAttribute.shared_ciphers: lambda: self._ssl_object.shared_ciphers() if self._ssl_object.server_side else None, TLSAttribute.standard_compatible: lambda: self.standard_compatible, TLSAttribute.ssl_object: lambda: self._ssl_object, TLSAttribute.tls_version: self._ssl_object.version}

@dataclass(eq=False)
class TLSListener(Listener[TLSStream]):
    listener: Listener[Any]
    ssl_context: ssl.SSLContext
    standard_compatible: bool = True
    handshake_timeout: float = 30

    @staticmethod
    async def handle_handshake_error(exc: BaseException, stream: AnyByteStream) -> None:
        await aclose_forcefully(stream)
        if not isinstance(exc, get_cancelled_exc_class()):
            logging.getLogger(__name__).exception('Error during TLS handshake', exc_info=exc)
        if not isinstance(exc, Exception) or isinstance(exc, get_cancelled_exc_class()):
            raise

    async def serve(self, handler: Callable[[TLSStream], Any], task_group: TaskGroup | None=None) -> None:

        @wraps(handler)
        async def handler_wrapper(stream: AnyByteStream) -> None:
            from .. import fail_after
            try:
                with fail_after(self.handshake_timeout):
                    wrapped_stream = await TLSStream.wrap(stream, ssl_context=self.ssl_context, standard_compatible=self.standard_compatible)
            except BaseException as exc:
                await self.handle_handshake_error(exc, stream)
            else:
                await handler(wrapped_stream)
        await self.listener.serve(handler_wrapper, task_group)

    async def aclose(self) -> None:
        await self.listener.aclose()

    @property
    def extra_attributes(self) -> Mapping[Any, Callable[[], Any]]:
        return {TLSAttribute.standard_compatible: lambda: self.standard_compatible}

class TLSConnectable(ByteStreamConnectable):

    def __init__(self, connectable: AnyByteStreamConnectable, *, hostname: str | None=None, ssl_context: ssl.SSLContext | None=None, standard_compatible: bool=True) -> None:
        self.connectable = connectable
        self.ssl_context: SSLContext = ssl_context or ssl.create_default_context(ssl.Purpose.SERVER_AUTH)
        if not isinstance(self.ssl_context, ssl.SSLContext):
            raise TypeError(f'ssl_context must be an instance of ssl.SSLContext, not {type(self.ssl_context).__name__}')
        self.hostname = hostname
        self.standard_compatible = standard_compatible

    @override
    async def connect(self) -> TLSStream:
        stream = await self.connectable.connect()
        try:
            return await TLSStream.wrap(stream, hostname=self.hostname, ssl_context=self.ssl_context, standard_compatible=self.standard_compatible)
        except BaseException:
            await aclose_forcefully(stream)
            raise
