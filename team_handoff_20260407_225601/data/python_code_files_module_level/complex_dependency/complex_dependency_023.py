import enum
import io
import multiprocessing
import os
from datetime import datetime
from functools import partial
from typing import Optional, Union
import fsspec
from fsspec.core import url_to_fs
from tqdm.contrib.concurrent import thread_map
from .. import config
from ..utils import tqdm as hf_tqdm
from ..utils.file_utils import ArchiveIterable, FilesIterable, cached_path, is_relative_path, stack_multiprocessing_download_progress_bars, url_or_path_join
from ..utils.info_utils import get_size_checksum_dict
from ..utils.logging import get_logger, tqdm
from ..utils.py_utils import NestedDataStructure, map_nested
from ..utils.track import tracked_str
from .download_config import DownloadConfig
logger = get_logger(__name__)

class DownloadMode(enum.Enum):
    REUSE_DATASET_IF_EXISTS = 'reuse_dataset_if_exists'
    REUSE_CACHE_IF_EXISTS = 'reuse_cache_if_exists'
    FORCE_REDOWNLOAD = 'force_redownload'

class DownloadManager:
    is_streaming = False

    def __init__(self, dataset_name: Optional[str]=None, data_dir: Optional[str]=None, download_config: Optional[DownloadConfig]=None, base_path: Optional[str]=None, record_checksums=True):
        self._dataset_name = dataset_name
        self._data_dir = data_dir
        self._base_path = base_path or os.path.abspath('.')
        self._recorded_sizes_checksums: dict[str, dict[str, Optional[Union[int, str]]]] = {}
        self.record_checksums = record_checksums
        self.download_config = download_config or DownloadConfig()
        self.downloaded_paths = {}
        self.extracted_paths = {}

    @property
    def manual_dir(self):
        return self._data_dir

    @property
    def downloaded_size(self):
        return sum((checksums_dict['num_bytes'] for checksums_dict in self._recorded_sizes_checksums.values()))

    def _record_sizes_checksums(self, url_or_urls: NestedDataStructure, downloaded_path_or_paths: NestedDataStructure):
        delay = 5
        for url, path in hf_tqdm(list(zip(url_or_urls.flatten(), downloaded_path_or_paths.flatten())), delay=delay, desc='Computing checksums'):
            self._recorded_sizes_checksums[str(url)] = get_size_checksum_dict(path, record_checksum=self.record_checksums)

    def download(self, url_or_urls):
        download_config = self.download_config.copy()
        download_config.extract_compressed_file = False
        if download_config.download_desc is None:
            download_config.download_desc = 'Downloading data'
        download_func = partial(self._download_batched, download_config=download_config)
        start_time = datetime.now()
        with stack_multiprocessing_download_progress_bars():
            downloaded_path_or_paths = map_nested(download_func, url_or_urls, map_tuple=True, num_proc=download_config.num_proc, desc='Downloading data files', batched=True, batch_size=-1)
        duration = datetime.now() - start_time
        logger.info(f'Downloading took {duration.total_seconds() // 60} min')
        url_or_urls = NestedDataStructure(url_or_urls)
        downloaded_path_or_paths = NestedDataStructure(downloaded_path_or_paths)
        self.downloaded_paths.update(dict(zip(url_or_urls.flatten(), downloaded_path_or_paths.flatten())))
        start_time = datetime.now()
        self._record_sizes_checksums(url_or_urls, downloaded_path_or_paths)
        duration = datetime.now() - start_time
        logger.info(f'Checksum Computation took {duration.total_seconds() // 60} min')
        return downloaded_path_or_paths.data

    def _download_batched(self, url_or_filenames: list[str], download_config: DownloadConfig) -> list[str]:
        if len(url_or_filenames) >= 16:
            download_config = download_config.copy()
            download_config.disable_tqdm = True
            download_func = partial(self._download_single, download_config=download_config)
            fs: fsspec.AbstractFileSystem
            path = str(url_or_filenames[0])
            if is_relative_path(path):
                path = url_or_path_join(self._base_path, path)
            fs, path = url_to_fs(path, **download_config.storage_options)
            size = 0
            try:
                size = fs.info(path).get('size', 0)
            except Exception:
                pass
            max_workers = config.HF_DATASETS_MULTITHREADING_MAX_WORKERS if size < 20 << 20 else 1
            return thread_map(download_func, url_or_filenames, desc=download_config.download_desc or 'Downloading', unit='files', position=multiprocessing.current_process()._identity[-1] if os.environ.get('HF_DATASETS_STACK_MULTIPROCESSING_DOWNLOAD_PROGRESS_BARS') == '1' and multiprocessing.current_process()._identity else None, max_workers=max_workers, tqdm_class=tqdm)
        else:
            return [self._download_single(url_or_filename, download_config=download_config) for url_or_filename in url_or_filenames]

    def _download_single(self, url_or_filename: str, download_config: DownloadConfig) -> str:
        url_or_filename = str(url_or_filename)
        if is_relative_path(url_or_filename):
            url_or_filename = url_or_path_join(self._base_path, url_or_filename)
        out = cached_path(url_or_filename, download_config=download_config)
        out = tracked_str(out)
        out.set_origin(url_or_filename)
        return out

    def iter_archive(self, path_or_buf: Union[str, io.BufferedReader]):
        if hasattr(path_or_buf, 'read'):
            return ArchiveIterable.from_buf(path_or_buf)
        else:
            return ArchiveIterable.from_urlpath(path_or_buf)

    def iter_files(self, paths: Union[str, list[str]]):
        return FilesIterable.from_urlpaths(paths)

    def extract(self, path_or_paths):
        download_config = self.download_config.copy()
        download_config.extract_compressed_file = True
        extract_func = partial(self._download_single, download_config=download_config)
        extracted_paths = map_nested(extract_func, path_or_paths, num_proc=download_config.num_proc, desc='Extracting data files')
        path_or_paths = NestedDataStructure(path_or_paths)
        extracted_paths = NestedDataStructure(extracted_paths)
        self.extracted_paths.update(dict(zip(path_or_paths.flatten(), extracted_paths.flatten())))
        return extracted_paths.data

    def download_and_extract(self, url_or_urls):
        return self.extract(self.download(url_or_urls))

    def get_recorded_sizes_checksums(self):
        return self._recorded_sizes_checksums.copy()

    def delete_extracted_files(self):
        paths_to_delete = set(self.extracted_paths.values()) - set(self.downloaded_paths.values())
        for key, path in list(self.extracted_paths.items()):
            if path in paths_to_delete and os.path.isfile(path):
                os.remove(path)
                del self.extracted_paths[key]

    def manage_extracted_files(self):
        if self.download_config.delete_extracted:
            self.delete_extracted_files()
