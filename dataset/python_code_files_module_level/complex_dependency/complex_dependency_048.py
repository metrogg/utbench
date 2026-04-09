import base64
import io
import json
import logging
import mimetypes
from dataclasses import dataclass
from pathlib import Path
from typing import TYPE_CHECKING, Any, AsyncIterable, BinaryIO, Iterable, Literal, NoReturn, Optional, Union, overload
import httpx
from huggingface_hub.errors import GenerationError, HfHubHTTPError, IncompleteGenerationError, OverloadedError, TextGenerationError, UnknownError, ValidationError
from ..utils import get_session, is_numpy_available, is_pillow_available
from ._generated.types import ChatCompletionStreamOutput, TextGenerationStreamOutput
if TYPE_CHECKING:
    from PIL.Image import Image
UrlT = str
PathT = Union[str, Path]
ContentT = Union[bytes, BinaryIO, PathT, UrlT, 'Image', bytearray, memoryview]
TASKS_EXPECTING_IMAGES = {'text-to-image', 'image-to-image'}
logger = logging.getLogger(__name__)

@dataclass
class RequestParameters:
    url: str
    task: str
    model: Optional[str]
    json: Optional[Union[str, dict, list]]
    data: Optional[bytes]
    headers: dict[str, Any]

class MimeBytes(bytes):
    mime_type: Optional[str]

    def __new__(cls, data: bytes, mime_type: Optional[str]=None):
        obj = super().__new__(cls, data)
        obj.mime_type = mime_type
        if isinstance(data, MimeBytes) and mime_type is None:
            obj.mime_type = data.mime_type
        return obj

def _import_numpy():
    if not is_numpy_available():
        raise ImportError('Please install numpy to use deal with embeddings (`pip install numpy`).')
    import numpy
    return numpy

def _import_pil_image():
    if not is_pillow_available():
        raise ImportError("Please install Pillow to use deal with images (`pip install Pillow`). If you don't want the image to be post-processed, use `client.post(...)` and get the raw response from the server.")
    from PIL import Image
    return Image

@overload
def _open_as_mime_bytes(content: ContentT) -> MimeBytes:
    ...

@overload
def _open_as_mime_bytes(content: Literal[None]) -> Literal[None]:
    ...

def _open_as_mime_bytes(content: Optional[ContentT]) -> Optional[MimeBytes]:
    if content is None:
        return None
    if isinstance(content, bytes):
        return MimeBytes(content)
    if isinstance(content, (bytearray, memoryview)):
        return MimeBytes(bytes(content))
    if hasattr(content, 'read'):
        logger.debug('Reading content from BinaryIO')
        data = content.read()
        mime_type = mimetypes.guess_type(str(content.name))[0] if hasattr(content, 'name') else None
        if isinstance(data, str):
            raise TypeError('Expected binary stream (bytes), but got text stream')
        return MimeBytes(data, mime_type=mime_type)
    if isinstance(content, str):
        if content.startswith('https://') or content.startswith('http://'):
            logger.debug(f'Downloading content from {content}')
            response = get_session().get(content)
            mime_type = response.headers.get('Content-Type')
            if mime_type is None:
                mime_type = mimetypes.guess_type(content)[0]
            return MimeBytes(response.content, mime_type=mime_type)
        content = Path(content)
        if not content.exists():
            raise FileNotFoundError(f'File not found at {content}. If `data` is a string, it must either be a URL or a path to a local file. To pass raw content, please encode it as bytes first.')
    if isinstance(content, Path):
        logger.debug(f'Opening content from {content}')
        return MimeBytes(content.read_bytes(), mime_type=mimetypes.guess_type(content)[0])
    if is_pillow_available():
        from PIL import Image
        if isinstance(content, Image.Image):
            logger.debug('Converting PIL Image to bytes')
            buffer = io.BytesIO()
            format = content.format or 'PNG'
            content.save(buffer, format=format)
            return MimeBytes(buffer.getvalue(), mime_type=f'image/{format.lower()}')
    raise TypeError(f'Unsupported content type: {type(content)}. Expected one of: bytes, bytearray, BinaryIO, memoryview, Path, str (URL or file path), or PIL.Image.Image.')

def _b64_encode(content: ContentT) -> str:
    raw_bytes = _open_as_mime_bytes(content)
    return base64.b64encode(raw_bytes).decode()

def _as_url(content: ContentT, default_mime_type: str) -> str:
    if isinstance(content, str) and content.startswith(('http://', 'https://', 'data:')):
        return content
    raw_bytes = _open_as_mime_bytes(content)
    mime_type = raw_bytes.mime_type or default_mime_type
    encoded_data = base64.b64encode(raw_bytes).decode()
    return f'data:{mime_type};base64,{encoded_data}'

def _b64_to_image(encoded_image: str) -> 'Image':
    Image = _import_pil_image()
    return Image.open(io.BytesIO(base64.b64decode(encoded_image)))

def _bytes_to_list(content: bytes) -> list:
    return json.loads(content.decode())

def _bytes_to_dict(content: bytes) -> dict:
    return json.loads(content.decode())

def _bytes_to_image(content: bytes) -> 'Image':
    Image = _import_pil_image()
    return Image.open(io.BytesIO(content))

def _as_dict(response: Union[bytes, dict]) -> dict:
    return json.loads(response) if isinstance(response, bytes) else response

def _stream_text_generation_response(output_lines: Iterable[str], details: bool) -> Union[Iterable[str], Iterable[TextGenerationStreamOutput]]:
    for line in output_lines:
        try:
            output = _format_text_generation_stream_output(line, details)
        except StopIteration:
            break
        if output is not None:
            yield output

async def _async_stream_text_generation_response(output_lines: AsyncIterable[str], details: bool) -> Union[AsyncIterable[str], AsyncIterable[TextGenerationStreamOutput]]:
    async for line in output_lines:
        try:
            output = _format_text_generation_stream_output(line, details)
        except StopIteration:
            break
        if output is not None:
            yield output

def _format_text_generation_stream_output(line: str, details: bool) -> Optional[Union[str, TextGenerationStreamOutput]]:
    if not line.startswith('data:'):
        return None
    if line.strip() == 'data: [DONE]':
        raise StopIteration('[DONE] signal received.')
    payload = line.lstrip('data:').rstrip('/n')
    json_payload = json.loads(payload)
    if json_payload.get('error') is not None:
        raise _parse_text_generation_error(json_payload['error'], json_payload.get('error_type'))
    output = TextGenerationStreamOutput.parse_obj_as_instance(json_payload)
    return output.token.text if not details else output

def _stream_chat_completion_response(lines: Iterable[str]) -> Iterable[ChatCompletionStreamOutput]:
    for line in lines:
        try:
            output = _format_chat_completion_stream_output(line)
        except StopIteration:
            break
        if output is not None:
            yield output

async def _async_stream_chat_completion_response(lines: AsyncIterable[str]) -> AsyncIterable[ChatCompletionStreamOutput]:
    async for line in lines:
        try:
            output = _format_chat_completion_stream_output(line)
        except StopIteration:
            break
        if output is not None:
            yield output

def _format_chat_completion_stream_output(line: str) -> Optional[ChatCompletionStreamOutput]:
    if not line.startswith('data:'):
        return None
    if line.strip() == 'data: [DONE]':
        raise StopIteration('[DONE] signal received.')
    json_payload = json.loads(line.lstrip('data:').strip())
    if json_payload.get('error') is not None:
        raise _parse_text_generation_error(json_payload['error'], json_payload.get('error_type'))
    return ChatCompletionStreamOutput.parse_obj_as_instance(json_payload)

async def _async_yield_from(client: httpx.AsyncClient, response: httpx.Response) -> AsyncIterable[str]:
    async for line in response.aiter_lines():
        yield line.strip()
_UNSUPPORTED_TEXT_GENERATION_KWARGS: dict[Optional[str], list[str]] = {}

def _set_unsupported_text_generation_kwargs(model: Optional[str], unsupported_kwargs: list[str]) -> None:
    _UNSUPPORTED_TEXT_GENERATION_KWARGS.setdefault(model, []).extend(unsupported_kwargs)

def _get_unsupported_text_generation_kwargs(model: Optional[str]) -> list[str]:
    return _UNSUPPORTED_TEXT_GENERATION_KWARGS.get(model, [])

def raise_text_generation_error(http_error: HfHubHTTPError) -> NoReturn:
    if http_error.response is None:
        raise http_error
    try:
        payload = getattr(http_error, 'response_error_payload', None) or http_error.response.json()
        error = payload.get('error')
        error_type = payload.get('error_type')
    except Exception:
        raise http_error
    if error_type is not None:
        exception = _parse_text_generation_error(error, error_type)
        raise exception from http_error
    raise http_error

def _parse_text_generation_error(error: Optional[str], error_type: Optional[str]) -> TextGenerationError:
    if error_type == 'generation':
        return GenerationError(error)
    if error_type == 'incomplete_generation':
        return IncompleteGenerationError(error)
    if error_type == 'overloaded':
        return OverloadedError(error)
    if error_type == 'validation':
        return ValidationError(error)
    return UnknownError(error)
