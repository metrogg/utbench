import os
from dataclasses import dataclass, field
from io import BytesIO
from pathlib import Path
from typing import TYPE_CHECKING, Any, ClassVar, Optional, Union
import numpy as np
import pyarrow as pa
from .. import config
from ..download.download_config import DownloadConfig
from ..table import array_cast
from ..utils.file_utils import is_local_path, xopen
from ..utils.py_utils import no_op_if_value_is_null, string_to_dict
if TYPE_CHECKING:
    from torchcodec.decoders import AudioDecoder
    from .features import FeatureType

@dataclass
class Audio:
    sampling_rate: Optional[int] = None
    decode: bool = True
    num_channels: Optional[int] = None
    stream_index: Optional[int] = None
    id: Optional[str] = field(default=None, repr=False)
    dtype: ClassVar[str] = 'dict'
    pa_type: ClassVar[Any] = pa.struct({'bytes': pa.binary(), 'path': pa.string()})
    _type: str = field(default='Audio', init=False, repr=False)

    def __call__(self):
        return self.pa_type

    def encode_example(self, value: Union[str, bytes, bytearray, dict, 'AudioDecoder']) -> dict:
        try:
            import torch
            from torchcodec.encoders import AudioEncoder
        except ImportError as err:
            raise ImportError("To support encoding audio data, please install 'torchcodec'.") from err
        if value is None:
            raise ValueError('value must be provided')
        if config.TORCHCODEC_AVAILABLE:
            from torchcodec.decoders import AudioDecoder
        else:
            AudioDecoder = None
        if isinstance(value, str):
            return {'bytes': None, 'path': value}
        elif isinstance(value, Path):
            return {'bytes': None, 'path': str(value.absolute())}
        elif isinstance(value, (bytes, bytearray)):
            return {'bytes': value, 'path': None}
        elif AudioDecoder is not None and isinstance(value, AudioDecoder):
            return encode_torchcodec_audio(value)
        elif 'array' in value:
            buffer = BytesIO()
            AudioEncoder(torch.from_numpy(value['array'].astype(np.float32)), sample_rate=value['sampling_rate']).to_file_like(buffer, format='wav', num_channels=self.num_channels)
            return {'bytes': buffer.getvalue(), 'path': None}
        elif value.get('path') is not None and os.path.isfile(value['path']):
            if value['path'].endswith('pcm'):
                if value.get('sampling_rate') is None:
                    raise KeyError("To use PCM files, please specify a 'sampling_rate' in Audio object")
                if value.get('bytes'):
                    bytes_value = np.frombuffer(value['bytes'], dtype=np.int16).astype(np.float32) / 32767
                else:
                    bytes_value = np.memmap(value['path'], dtype='h', mode='r').astype(np.float32) / 32767
                buffer = BytesIO()
                AudioEncoder(torch.from_numpy(bytes_value), sample_rate=value['sampling_rate']).to_file_like(buffer, format='wav', num_channels=self.num_channels)
                return {'bytes': buffer.getvalue(), 'path': None}
            else:
                return {'bytes': None, 'path': value.get('path')}
        elif value.get('bytes') is not None or value.get('path') is not None:
            return {'bytes': value.get('bytes'), 'path': value.get('path')}
        else:
            raise ValueError(f"An audio sample should have one of 'path' or 'bytes' but they are missing or None in {value}.")

    def decode_example(self, value: dict, token_per_repo_id: Optional[dict[str, Union[str, bool, None]]]=None) -> 'AudioDecoder':
        if config.TORCHCODEC_AVAILABLE:
            from ._torchcodec import AudioDecoder
        else:
            raise ImportError("To support decoding audio data, please install 'torchcodec'.")
        if not self.decode:
            raise RuntimeError('Decoding is disabled for this feature. Please use Audio(decode=True) instead.')
        path, bytes = (value['path'], value['bytes']) if value['bytes'] is not None else (value['path'], None)
        if path is None and bytes is None:
            raise ValueError(f"An audio sample should have one of 'path' or 'bytes' but both are None in {value}.")
        if bytes is None and is_local_path(path):
            audio = AudioDecoder(path, stream_index=self.stream_index, sample_rate=self.sampling_rate, num_channels=self.num_channels)
        elif bytes is None:
            token_per_repo_id = token_per_repo_id or {}
            source_url = path.split('::')[-1]
            pattern = config.HUB_DATASETS_URL if source_url.startswith(config.HF_ENDPOINT) else config.HUB_DATASETS_HFFS_URL
            source_url_fields = string_to_dict(source_url, pattern)
            token = token_per_repo_id.get(source_url_fields['repo_id']) if source_url_fields is not None else None
            download_config = DownloadConfig(token=token)
            f = xopen(path, 'rb', download_config=download_config)
            audio = AudioDecoder(f, stream_index=self.stream_index, sample_rate=self.sampling_rate, num_channels=self.num_channels)
        else:
            audio = AudioDecoder(bytes, stream_index=self.stream_index, sample_rate=self.sampling_rate, num_channels=self.num_channels)
        audio._hf_encoded = {'path': path, 'bytes': bytes}
        audio.metadata.path = path
        return audio

    def flatten(self) -> Union['FeatureType', dict[str, 'FeatureType']]:
        from .features import Value
        if self.decode:
            raise ValueError('Cannot flatten a decoded Audio feature.')
        return {'bytes': Value('binary'), 'path': Value('string')}

    def cast_storage(self, storage: Union[pa.StringArray, pa.StructArray]) -> pa.StructArray:
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
        elif pa.types.is_struct(storage.type) and storage.type.get_all_field_indices('array'):
            storage = pa.array([Audio().encode_example(x) if x is not None else None for x in storage.to_numpy(zero_copy_only=False)])
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

def encode_torchcodec_audio(audio: 'AudioDecoder') -> dict:
    if hasattr(audio, '_hf_encoded'):
        return audio._hf_encoded
    else:
        try:
            from torchcodec.encoders import AudioEncoder
        except ImportError as err:
            raise ImportError("To support encoding audio data, please install 'torchcodec'.") from err
        samples = audio.get_all_samples()
        buffer = BytesIO()
        num_channels = samples.data.shape[0]
        AudioEncoder(samples.data.cpu(), sample_rate=samples.sample_rate).to_file_like(buffer, format='wav', num_channels=num_channels)
        return {'bytes': buffer.getvalue(), 'path': None}
