import asyncio
import random
import sys
from functools import partial
from typing import Final, Optional, Set, Union
from ..base_protocol import BaseProtocol
from ..client_exceptions import ClientConnectionResetError
from ..compression_utils import ZLibBackend, ZLibCompressor
from .helpers import MASK_LEN, MSG_SIZE, PACK_CLOSE_CODE, PACK_LEN1, PACK_LEN2, PACK_LEN3, PACK_RANDBITS, websocket_mask
from .models import WS_DEFLATE_TRAILING, WSMsgType
DEFAULT_LIMIT: Final[int] = 2 ** 16
WS_CONTROL_FRAME_OPCODE: Final[int] = 8
WEBSOCKET_MAX_SYNC_CHUNK_SIZE = 16 * 1024

class WebSocketWriter:

    def __init__(self, protocol: BaseProtocol, transport: asyncio.Transport, *, use_mask: bool=False, limit: int=DEFAULT_LIMIT, random: random.Random=random.Random(), compress: int=0, notakeover: bool=False) -> None:
        self.protocol = protocol
        self.transport = transport
        self.use_mask = use_mask
        self.get_random_bits = partial(random.getrandbits, 32)
        self.compress = compress
        self.notakeover = notakeover
        self._closing = False
        self._limit = limit
        self._output_size = 0
        self._compressobj: Optional[ZLibCompressor] = None
        self._send_lock = asyncio.Lock()
        self._background_tasks: Set[asyncio.Task[None]] = set()

    async def send_frame(self, message: bytes, opcode: int, compress: Optional[int]=None) -> None:
        if self._closing and (not opcode & WSMsgType.CLOSE):
            raise ClientConnectionResetError('Cannot write to closing transport')
        if not (compress or self.compress) or opcode >= WS_CONTROL_FRAME_OPCODE:
            self._write_websocket_frame(message, opcode, 0)
        elif len(message) <= WEBSOCKET_MAX_SYNC_CHUNK_SIZE:
            async with self._send_lock:
                self._send_compressed_frame_sync(message, opcode, compress)
        else:
            loop = asyncio.get_running_loop()
            coro = self._send_compressed_frame_async_locked(message, opcode, compress)
            if sys.version_info >= (3, 12):
                send_task = asyncio.Task(coro, loop=loop, eager_start=True)
            else:
                send_task = loop.create_task(coro)
            self._background_tasks.add(send_task)
            send_task.add_done_callback(self._background_tasks.discard)
            await asyncio.shield(send_task)
        if self._output_size > self._limit:
            self._output_size = 0
            if self.protocol._paused:
                await self.protocol._drain_helper()

    def _write_websocket_frame(self, message: bytes, opcode: int, rsv: int) -> None:
        msg_length = len(message)
        use_mask = self.use_mask
        mask_bit = 128 if use_mask else 0
        first_byte = 128 | rsv | opcode
        if msg_length < 126:
            header = PACK_LEN1(first_byte, msg_length | mask_bit)
            header_len = 2
        elif msg_length < 65536:
            header = PACK_LEN2(first_byte, 126 | mask_bit, msg_length)
            header_len = 4
        else:
            header = PACK_LEN3(first_byte, 127 | mask_bit, msg_length)
            header_len = 10
        if self.transport.is_closing():
            raise ClientConnectionResetError('Cannot write to closing transport')
        if use_mask:
            mask = PACK_RANDBITS(self.get_random_bits())
            message = bytearray(message)
            websocket_mask(mask, message)
            self.transport.write(header + mask + message)
            self._output_size += MASK_LEN
        elif msg_length > MSG_SIZE:
            self.transport.write(header)
            self.transport.write(message)
        else:
            self.transport.write(header + message)
        self._output_size += header_len + msg_length

    def _get_compressor(self, compress: Optional[int]) -> ZLibCompressor:
        if compress:
            return ZLibCompressor(level=ZLibBackend.Z_BEST_SPEED, wbits=-compress, max_sync_chunk_size=WEBSOCKET_MAX_SYNC_CHUNK_SIZE)
        if not self._compressobj:
            self._compressobj = ZLibCompressor(level=ZLibBackend.Z_BEST_SPEED, wbits=-self.compress, max_sync_chunk_size=WEBSOCKET_MAX_SYNC_CHUNK_SIZE)
        return self._compressobj

    def _send_compressed_frame_sync(self, message: bytes, opcode: int, compress: Optional[int]) -> None:
        compressobj = self._get_compressor(compress)
        self._write_websocket_frame((compressobj.compress_sync(message) + compressobj.flush(ZLibBackend.Z_FULL_FLUSH if self.notakeover else ZLibBackend.Z_SYNC_FLUSH)).removesuffix(WS_DEFLATE_TRAILING), opcode, 64)

    async def _send_compressed_frame_async_locked(self, message: bytes, opcode: int, compress: Optional[int]) -> None:
        async with self._send_lock:
            compressobj = self._get_compressor(compress)
            self._write_websocket_frame((await compressobj.compress(message) + compressobj.flush(ZLibBackend.Z_FULL_FLUSH if self.notakeover else ZLibBackend.Z_SYNC_FLUSH)).removesuffix(WS_DEFLATE_TRAILING), opcode, 64)

    async def close(self, code: int=1000, message: Union[bytes, str]=b'') -> None:
        if isinstance(message, str):
            message = message.encode('utf-8')
        try:
            await self.send_frame(PACK_CLOSE_CODE(code) + message, opcode=WSMsgType.CLOSE)
        finally:
            self._closing = True
