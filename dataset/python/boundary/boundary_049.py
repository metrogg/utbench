from __future__ import annotations
import codecs
import io
import typing
import zlib
from ._exceptions import DecodingError
try:
    import brotli
except ImportError:
    try:
        import brotlicffi as brotli
    except ImportError:
        brotli = None
try:
    import zstandard
except ImportError:
    zstandard = None

class ContentDecoder:

    def decode(self, data: bytes) -> bytes:
        raise NotImplementedError()

    def flush(self) -> bytes:
        raise NotImplementedError()

class IdentityDecoder(ContentDecoder):

    def decode(self, data: bytes) -> bytes:
        return data

    def flush(self) -> bytes:
        return b''

class DeflateDecoder(ContentDecoder):

    def __init__(self) -> None:
        self.first_attempt = True
        self.decompressor = zlib.decompressobj()

    def decode(self, data: bytes) -> bytes:
        was_first_attempt = self.first_attempt
        self.first_attempt = False
        try:
            return self.decompressor.decompress(data)
        except zlib.error as exc:
            if was_first_attempt:
                self.decompressor = zlib.decompressobj(-zlib.MAX_WBITS)
                return self.decode(data)
            raise DecodingError(str(exc)) from exc

    def flush(self) -> bytes:
        try:
            return self.decompressor.flush()
        except zlib.error as exc:
            raise DecodingError(str(exc)) from exc

class GZipDecoder(ContentDecoder):

    def __init__(self) -> None:
        self.decompressor = zlib.decompressobj(zlib.MAX_WBITS | 16)

    def decode(self, data: bytes) -> bytes:
        try:
            return self.decompressor.decompress(data)
        except zlib.error as exc:
            raise DecodingError(str(exc)) from exc

    def flush(self) -> bytes:
        try:
            return self.decompressor.flush()
        except zlib.error as exc:
            raise DecodingError(str(exc)) from exc

class BrotliDecoder(ContentDecoder):

    def __init__(self) -> None:
        if brotli is None:
            raise ImportError("Using 'BrotliDecoder', but neither of the 'brotlicffi' or 'brotli' packages have been installed. Make sure to install httpx using `pip install httpx[brotli]`.") from None
        self.decompressor = brotli.Decompressor()
        self.seen_data = False
        self._decompress: typing.Callable[[bytes], bytes]
        if hasattr(self.decompressor, 'decompress'):
            self._decompress = self.decompressor.decompress
        else:
            self._decompress = self.decompressor.process

    def decode(self, data: bytes) -> bytes:
        if not data:
            return b''
        self.seen_data = True
        try:
            return self._decompress(data)
        except brotli.error as exc:
            raise DecodingError(str(exc)) from exc

    def flush(self) -> bytes:
        if not self.seen_data:
            return b''
        try:
            if hasattr(self.decompressor, 'finish'):
                self.decompressor.finish()
            return b''
        except brotli.error as exc:
            raise DecodingError(str(exc)) from exc

class ZStandardDecoder(ContentDecoder):

    def __init__(self) -> None:
        if zstandard is None:
            raise ImportError("Using 'ZStandardDecoder', ...Make sure to install httpx using `pip install httpx[zstd]`.") from None
        self.decompressor = zstandard.ZstdDecompressor().decompressobj()
        self.seen_data = False

    def decode(self, data: bytes) -> bytes:
        assert zstandard is not None
        self.seen_data = True
        output = io.BytesIO()
        try:
            output.write(self.decompressor.decompress(data))
            while self.decompressor.eof and self.decompressor.unused_data:
                unused_data = self.decompressor.unused_data
                self.decompressor = zstandard.ZstdDecompressor().decompressobj()
                output.write(self.decompressor.decompress(unused_data))
        except zstandard.ZstdError as exc:
            raise DecodingError(str(exc)) from exc
        return output.getvalue()

    def flush(self) -> bytes:
        if not self.seen_data:
            return b''
        ret = self.decompressor.flush()
        if not self.decompressor.eof:
            raise DecodingError('Zstandard data is incomplete')
        return bytes(ret)

class MultiDecoder(ContentDecoder):

    def __init__(self, children: typing.Sequence[ContentDecoder]) -> None:
        self.children = list(reversed(children))

    def decode(self, data: bytes) -> bytes:
        for child in self.children:
            data = child.decode(data)
        return data

    def flush(self) -> bytes:
        data = b''
        for child in self.children:
            data = child.decode(data) + child.flush()
        return data

class ByteChunker:

    def __init__(self, chunk_size: int | None=None) -> None:
        self._buffer = io.BytesIO()
        self._chunk_size = chunk_size

    def decode(self, content: bytes) -> list[bytes]:
        if self._chunk_size is None:
            return [content] if content else []
        self._buffer.write(content)
        if self._buffer.tell() >= self._chunk_size:
            value = self._buffer.getvalue()
            chunks = [value[i:i + self._chunk_size] for i in range(0, len(value), self._chunk_size)]
            if len(chunks[-1]) == self._chunk_size:
                self._buffer.seek(0)
                self._buffer.truncate()
                return chunks
            else:
                self._buffer.seek(0)
                self._buffer.write(chunks[-1])
                self._buffer.truncate()
                return chunks[:-1]
        else:
            return []

    def flush(self) -> list[bytes]:
        value = self._buffer.getvalue()
        self._buffer.seek(0)
        self._buffer.truncate()
        return [value] if value else []

class TextChunker:

    def __init__(self, chunk_size: int | None=None) -> None:
        self._buffer = io.StringIO()
        self._chunk_size = chunk_size

    def decode(self, content: str) -> list[str]:
        if self._chunk_size is None:
            return [content] if content else []
        self._buffer.write(content)
        if self._buffer.tell() >= self._chunk_size:
            value = self._buffer.getvalue()
            chunks = [value[i:i + self._chunk_size] for i in range(0, len(value), self._chunk_size)]
            if len(chunks[-1]) == self._chunk_size:
                self._buffer.seek(0)
                self._buffer.truncate()
                return chunks
            else:
                self._buffer.seek(0)
                self._buffer.write(chunks[-1])
                self._buffer.truncate()
                return chunks[:-1]
        else:
            return []

    def flush(self) -> list[str]:
        value = self._buffer.getvalue()
        self._buffer.seek(0)
        self._buffer.truncate()
        return [value] if value else []

class TextDecoder:

    def __init__(self, encoding: str='utf-8') -> None:
        self.decoder = codecs.getincrementaldecoder(encoding)(errors='replace')

    def decode(self, data: bytes) -> str:
        return self.decoder.decode(data)

    def flush(self) -> str:
        return self.decoder.decode(b'', True)

class LineDecoder:

    def __init__(self) -> None:
        self.buffer: list[str] = []
        self.trailing_cr: bool = False

    def decode(self, text: str) -> list[str]:
        NEWLINE_CHARS = '\n\r\x0b\x0c\x1c\x1d\x1e\x85\u2028\u2029'
        if self.trailing_cr:
            text = '\r' + text
            self.trailing_cr = False
        if text.endswith('\r'):
            self.trailing_cr = True
            text = text[:-1]
        if not text:
            return []
        trailing_newline = text[-1] in NEWLINE_CHARS
        lines = text.splitlines()
        if len(lines) == 1 and (not trailing_newline):
            self.buffer.append(lines[0])
            return []
        if self.buffer:
            lines = [''.join(self.buffer) + lines[0]] + lines[1:]
            self.buffer = []
        if not trailing_newline:
            self.buffer = [lines.pop()]
        return lines

    def flush(self) -> list[str]:
        if not self.buffer and (not self.trailing_cr):
            return []
        lines = [''.join(self.buffer)]
        self.buffer = []
        self.trailing_cr = False
        return lines
SUPPORTED_DECODERS = {'identity': IdentityDecoder, 'gzip': GZipDecoder, 'deflate': DeflateDecoder, 'br': BrotliDecoder, 'zstd': ZStandardDecoder}
if brotli is None:
    SUPPORTED_DECODERS.pop('br')
if zstandard is None:
    SUPPORTED_DECODERS.pop('zstd')
