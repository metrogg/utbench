import os
from dataclasses import dataclass, field
from pathlib import Path
from typing import TYPE_CHECKING, Any, ClassVar, Dict, Optional, Union
import pyarrow as pa
from .. import config
from ..download.download_config import DownloadConfig
from ..table import array_cast
from ..utils.file_utils import is_local_path, xopen
from ..utils.py_utils import no_op_if_value_is_null, string_to_dict
if TYPE_CHECKING:
    import nibabel as nib
    from .features import FeatureType
if config.NIBABEL_AVAILABLE:
    import nibabel as nib

    class Nifti1ImageWrapper(nib.nifti1.Nifti1Image):

        def __init__(self, nifti_image: nib.nifti1.Nifti1Image):
            super().__init__(dataobj=nifti_image.dataobj, affine=nifti_image.affine, header=nifti_image.header, extra=nifti_image.extra, file_map=nifti_image.file_map, dtype=nifti_image.get_data_dtype())
            self.nifti_image = nifti_image

        def _repr_html_(self):
            from ipyniivue import NiiVue, ShowRender, SliceType, Volume
            from IPython.display import display
            bytes_ = self.nifti_image.to_bytes()
            nv = NiiVue()
            nv.set_slice_type(SliceType.MULTIPLANAR)
            nv.opts.multiplanar_show_render = ShowRender.ALWAYS
            nv.opts.show_3d_crosshair = True
            nv.opts.multiplanar_force_render = True
            name = None
            if hasattr(self.nifti_image, 'file_map'):
                if 'image' in self.nifti_image.file_map and getattr(self.nifti_image.file_map['image'], 'filename', None) is not None:
                    name = self.nifti_image.file_map['image'].filename
            if name is None:
                name = 'volume.nii.gz'
            volume = Volume(name=name, data=bytes_)
            nv.load_volumes([volume])
            display(nv)

@dataclass
class Nifti:
    decode: bool = True
    id: Optional[str] = field(default=None, repr=False)
    dtype: ClassVar[str] = 'nibabel.nifti1.Nifti1Image'
    pa_type: ClassVar[Any] = pa.struct({'bytes': pa.binary(), 'path': pa.string()})
    _type: str = field(default='Nifti', init=False, repr=False)

    def __call__(self):
        return self.pa_type

    def encode_example(self, value: Union[str, bytes, bytearray, dict, 'nib.Nifti1Image']) -> dict:
        if config.NIBABEL_AVAILABLE:
            import nibabel as nib
        else:
            nib = None
        if isinstance(value, str):
            return {'path': value, 'bytes': None}
        elif isinstance(value, Path):
            return {'path': str(value.absolute()), 'bytes': None}
        elif isinstance(value, (bytes, bytearray)):
            return {'path': None, 'bytes': value}
        elif nib is not None and isinstance(value, nib.spatialimages.SpatialImage):
            return encode_nibabel_image(value)
        elif isinstance(value, dict):
            if value.get('path') is not None and os.path.isfile(value['path']):
                return {'bytes': None, 'path': value.get('path')}
            elif value.get('bytes') is not None or value.get('path') is not None:
                return {'bytes': value.get('bytes'), 'path': value.get('path')}
            else:
                raise ValueError(f"A nifti sample should have one of 'path' or 'bytes' but they are missing or None in {value}.")
        else:
            raise ValueError(f'A nifti sample should be a string, bytes, Path, nibabel image, or dict, but got {type(value)}.')

    def decode_example(self, value: dict, token_per_repo_id=None) -> 'Nifti1ImageWrapper':
        if config.NIBABEL_AVAILABLE:
            import nibabel as nib
        else:
            raise ImportError("To support decoding NIfTI files, please install 'nibabel'.")
        if token_per_repo_id is None:
            token_per_repo_id = {}
        path, bytes_ = (value['path'], value['bytes'])
        if bytes_ is None:
            if path is None:
                raise ValueError(f"A nifti should have one of 'path' or 'bytes' but both are None in {value}.")
            else:
                if path.startswith('gzip://') and is_local_path(path.split('::')[-1]):
                    path = path.split('::')[-1]
                if is_local_path(path):
                    nifti = nib.load(path)
                else:
                    source_url = path.split('::')[-1]
                    pattern = config.HUB_DATASETS_URL if source_url.startswith(config.HF_ENDPOINT) else config.HUB_DATASETS_HFFS_URL
                    source_url_fields = string_to_dict(source_url, pattern)
                    token = token_per_repo_id.get(source_url_fields['repo_id']) if source_url_fields is not None else None
                    download_config = DownloadConfig(token=token)
                    with xopen(path, 'rb', download_config=download_config) as f:
                        nifti = nib.load(f)
        else:
            import gzip
            if bytes_[:2] == b'\x1f\x8b':
                bytes_ = gzip.decompress(bytes_)
            nifti = nib.Nifti1Image.from_bytes(bytes_)
        return Nifti1ImageWrapper(nifti)

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

    def flatten(self) -> Union['FeatureType', Dict[str, 'FeatureType']]:
        from .features import Value
        return self if self.decode else {'bytes': Value('binary'), 'path': Value('string')}

    def cast_storage(self, storage: Union[pa.StringArray, pa.StructArray, pa.BinaryArray]) -> pa.StructArray:
        if pa.types.is_string(storage.type):
            bytes_array = pa.array([None] * len(storage), type=pa.binary())
            storage = pa.StructArray.from_arrays([bytes_array, storage], ['bytes', 'path'], mask=storage.is_null())
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
        return array_cast(storage, self.pa_type)

def encode_nibabel_image(img: 'nib.Nifti1Image', force_bytes: bool=False) -> dict[str, Optional[Union[str, bytes]]]:
    if hasattr(img, 'file_map') and img.file_map is not None and (not force_bytes):
        filename = img.file_map['image'].filename
        return {'path': filename, 'bytes': None}
    bytes_data = img.to_bytes()
    return {'path': None, 'bytes': bytes_data}
