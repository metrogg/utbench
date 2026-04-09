import os
from dataclasses import dataclass, field
from pathlib import Path
from typing import TYPE_CHECKING, Any, ClassVar, Literal, Optional, TypedDict, Union
import numpy as np
import pyarrow as pa
from .. import config
from ..download.download_config import DownloadConfig
from ..table import array_cast
from ..utils.file_utils import is_local_path, xopen
from ..utils.py_utils import no_op_if_value_is_null, string_to_dict
if TYPE_CHECKING:
    import torch
    from torchcodec.decoders import VideoDecoder
    from .features import FeatureType

class Example(TypedDict):
    path: Optional[str]
    bytes: Optional[bytes]

@dataclass
class Video:
    decode: bool = True
    stream_index: Optional[int] = None
    dimension_order: Literal['NCHW', 'NHWC'] = 'NCHW'
    num_ffmpeg_threads: int = 1
    device: Optional[Union[str, 'torch.device']] = 'cpu'
    seek_mode: Literal['exact', 'approximate'] = 'exact'
    id: Optional[str] = field(default=None, repr=False)
    dtype: ClassVar[str] = 'torchcodec.decoders.VideoDecoder'
    pa_type: ClassVar[Any] = pa.struct({'bytes': pa.binary(), 'path': pa.string()})
    _type: str = field(default='Video', init=False, repr=False)

    def __call__(self):
        return self.pa_type

    def encode_example(self, value: Union[str, bytes, bytearray, Example, np.ndarray, 'VideoDecoder']) -> Example:
        if value is None:
            raise ValueError('value must be provided')
        if config.TORCHCODEC_AVAILABLE:
            from torchcodec.decoders import VideoDecoder
        else:
            VideoDecoder = None
        if isinstance(value, list):
            value = np.array(value)
        if isinstance(value, str):
            return {'path': value, 'bytes': None}
        elif isinstance(value, Path):
            return {'path': str(value.absolute()), 'bytes': None}
        elif isinstance(value, (bytes, bytearray)):
            return {'path': None, 'bytes': value}
        elif isinstance(value, np.ndarray):
            return encode_np_array(value)
        elif VideoDecoder is not None and isinstance(value, VideoDecoder):
            return encode_torchcodec_video(value)
        elif isinstance(value, dict):
            path, bytes_ = (value.get('path'), value.get('bytes'))
            if path is not None and os.path.isfile(path):
                return {'bytes': None, 'path': path}
            elif bytes_ is not None or path is not None:
                return {'bytes': bytes_, 'path': path}
            else:
                raise ValueError(f"A video sample should have one of 'path' or 'bytes' but they are missing or None in {value}.")
        else:
            raise TypeError(f'Unsupported encode_example type: {type(value)}')

    def decode_example(self, value: Union[str, Example], token_per_repo_id: Optional[dict[str, Union[bool, str]]]=None) -> 'VideoDecoder':
        if not self.decode:
            raise RuntimeError('Decoding is disabled for this feature. Please use Video(decode=True) instead.')
        if config.TORCHCODEC_AVAILABLE:
            from torchcodec.decoders import VideoDecoder
        else:
            raise ImportError("To support decoding videos, please install 'torchcodec'.")
        if token_per_repo_id is None:
            token_per_repo_id = {}
        if isinstance(value, str):
            path, bytes_ = (value, None)
        else:
            path, bytes_ = (value['path'], value['bytes'])
        if bytes_ is None:
            if path is None:
                raise ValueError(f"A video should have one of 'path' or 'bytes' but both are None in {value}.")
            elif is_local_path(path):
                video = VideoDecoder(path, stream_index=self.stream_index, dimension_order=self.dimension_order, num_ffmpeg_threads=self.num_ffmpeg_threads, device=self.device, seek_mode=self.seek_mode)
            else:
                video = hf_video_reader(path, token_per_repo_id=token_per_repo_id, dimension_order=self.dimension_order, num_ffmpeg_threads=self.num_ffmpeg_threads, device=self.device, seek_mode=self.seek_mode)
        else:
            video = VideoDecoder(bytes_, stream_index=self.stream_index, dimension_order=self.dimension_order, num_ffmpeg_threads=self.num_ffmpeg_threads, device=self.device, seek_mode=self.seek_mode)
        video._hf_encoded = {'path': path, 'bytes': bytes_}
        video.metadata.path = path
        return video

    def flatten(self) -> Union['FeatureType', dict[str, 'FeatureType']]:
        from .features import Value
        return self if self.decode else {'bytes': Value('binary'), 'path': Value('string')}

    def cast_storage(self, storage: Union[pa.StringArray, pa.StructArray, pa.ListArray]) -> pa.StructArray:
        if pa.types.is_string(storage.type):
            bytes_array = pa.array([None] * len(storage), type=pa.binary())
            storage = pa.StructArray.from_arrays([bytes_array, storage], ['bytes', 'path'], mask=storage.is_null())
        elif pa.types.is_large_binary(storage.type):
            storage = array_cast(storage, pa.binary())
            path_array = pa.array([None] * len(storage), type=pa.string())
            storage = pa.StructArray.from_arrays([storage, path_array], ['bytes', 'path'], mask=storage.is_null())
        elif pa.types.is_binary(storage.type):
            path_array = pa.array([None] * len(storage), type=pa.string())
            storage = pa.StructArray.from_arrays([storage, path_array], ['bytes', 'path'], mask=storage.is_null())
        elif pa.types.is_struct(storage.type):
            if storage.type.get_field_index('bytes') >= 0:
                bytes_array = storage.field('bytes')
            else:
                bytes_array = pa.array([None] * len(storage), type=pa.binary())
            if storage.type.get_field_index('path') >= 0:
                path_array = storage.field('path')
            else:
                path_array = pa.array([None] * len(storage), type=pa.string())
            storage = pa.StructArray.from_arrays([bytes_array, path_array], ['bytes', 'path'], mask=storage.is_null())
        elif pa.types.is_list(storage.type):
            bytes_array = pa.array([encode_np_array(np.array(arr))['bytes'] if arr is not None else None for arr in storage.to_pylist()], type=pa.binary())
            path_array = pa.array([None] * len(storage), type=pa.string())
            storage = pa.StructArray.from_arrays([bytes_array, path_array], ['bytes', 'path'], mask=bytes_array.is_null())
        return array_cast(storage, self.pa_type)

    def embed_storage(self, storage: pa.StructArray, token_per_repo_id=None) -> pa.StructArray:
        if token_per_repo_id is None:
            token_per_repo_id = {}

        @no_op_if_value_is_null
        def path_to_bytes(path):
            source_url = path.split('::')[-1]
            pattern = config.HUB_DATASETS_URL if source_url.startswith(config.HF_ENDPOINT) else config.HUB_DATASETS_HFFS_URL
            source_url_fields = string_to_dict(source_url, pattern)
            token = token_per_repo_id.get(source_url_fields['repo_id']) if source_url_fields is not None else None
            download_config = DownloadConfig(token=token)
            with xopen(path, 'rb', download_config=download_config) as f:
                return f.read()
        bytes_array = pa.array([(path_to_bytes(x['path']) if x['bytes'] is None else x['bytes']) if x is not None else None for x in storage.to_pylist()], type=pa.binary())
        path_array = pa.array([os.path.basename(path) if path is not None else None for path in storage.field('path').to_pylist()], type=pa.string())
        storage = pa.StructArray.from_arrays([bytes_array, path_array], ['bytes', 'path'], mask=bytes_array.is_null())
        return array_cast(storage, self.pa_type)

def video_to_bytes(video: 'VideoDecoder') -> bytes:
    raise NotImplementedError()

def encode_torchcodec_video(video: 'VideoDecoder') -> Example:
    if hasattr(video, '_hf_encoded'):
        return video._hf_encoded
    else:
        raise NotImplementedError("Encoding a VideoDecoder that doesn't come from datasets.Video.decode() is not implemented")

def encode_np_array(array: np.ndarray) -> Example:
    raise NotImplementedError()

def hf_video_reader(path: str, token_per_repo_id: Optional[dict[str, Union[bool, str]]]=None, stream: str='video', dimension_order: Literal['NCHW', 'NHWC']='NCHW', num_ffmpeg_threads: int=1, device: Optional[Union[str, 'torch.device']]='cpu', seek_mode: Literal['exact', 'approximate']='exact') -> 'VideoDecoder':
    from torchcodec.decoders import VideoDecoder
    if token_per_repo_id is None:
        token_per_repo_id = {}
    source_url = path.split('::')[-1]
    pattern = config.HUB_DATASETS_URL if source_url.startswith(config.HF_ENDPOINT) else config.HUB_DATASETS_HFFS_URL
    source_url_fields = string_to_dict(source_url, pattern)
    token = token_per_repo_id.get(source_url_fields['repo_id']) if source_url_fields is not None else None
    download_config = DownloadConfig(token=token)
    f = xopen(path, 'rb', download_config=download_config)
    stream_id = 0 if len(stream.split(':')) == 1 else int(stream.split(':')[1])
    vd = VideoDecoder(f, stream_index=stream_id, dimension_order=dimension_order, num_ffmpeg_threads=num_ffmpeg_threads, device=device, seek_mode=seek_mode)
    return vd
