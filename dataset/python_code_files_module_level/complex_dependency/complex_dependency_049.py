import io
import os
from collections.abc import Iterable
from typing import Optional, Union
from ..utils.file_utils import SINGLE_FILE_COMPRESSION_PROTOCOLS, ArchiveIterable, FilesIterable, _get_extraction_protocol, _get_path_extension, _prepare_path_and_storage_options, is_relative_path, url_or_path_join, xbasename, xdirname, xet_parse, xexists, xgetsize, xglob, xgzip_open, xisdir, xisfile, xjoin, xlistdir, xnumpy_load, xopen, xpandas_read_csv, xpandas_read_excel, xPath, xpyarrow_parquet_read_table, xrelpath, xsio_loadmat, xsplit, xsplitext, xwalk, xxml_dom_minidom_parse
from ..utils.logging import get_logger
from ..utils.py_utils import map_nested
from .download_config import DownloadConfig
logger = get_logger(__name__)

class StreamingDownloadManager:
    is_streaming = True

    def __init__(self, dataset_name: Optional[str]=None, data_dir: Optional[str]=None, download_config: Optional[DownloadConfig]=None, base_path: Optional[str]=None):
        self._dataset_name = dataset_name
        self._data_dir = data_dir
        self._base_path = base_path or os.path.abspath('.')
        self.download_config = download_config or DownloadConfig()
        self.downloaded_size = None
        self.record_checksums = False

    @property
    def manual_dir(self):
        return self._data_dir

    def download(self, url_or_urls):
        url_or_urls = map_nested(self._download_single, url_or_urls, map_tuple=True)
        return url_or_urls

    def _download_single(self, urlpath: str) -> str:
        urlpath = str(urlpath)
        if is_relative_path(urlpath):
            urlpath = url_or_path_join(self._base_path, urlpath)
        return urlpath

    def extract(self, url_or_urls):
        urlpaths = map_nested(self._extract, url_or_urls, map_tuple=True)
        return urlpaths

    def _extract(self, urlpath: str) -> str:
        urlpath = str(urlpath)
        protocol = _get_extraction_protocol(urlpath, download_config=self.download_config)
        path = urlpath.split('::')[0]
        extension = _get_path_extension(path)
        if extension in ['tgz', 'tar'] or path.endswith(('.tar.gz', '.tar.bz2', '.tar.xz')):
            raise NotImplementedError(f"Extraction protocol for TAR archives like '{urlpath}' is not implemented in streaming mode. Please use `dl_manager.iter_archive` instead.\n\nExample usage:\n\n\turl = dl_manager.download(url)\n\ttar_archive_iterator = dl_manager.iter_archive(url)\n\n\tfor filename, file in tar_archive_iterator:\n\t\t...")
        if protocol is None:
            return urlpath
        elif protocol in SINGLE_FILE_COMPRESSION_PROTOCOLS:
            inner_file = os.path.basename(urlpath.split('::')[0])
            inner_file = inner_file[:inner_file.rindex('.')] if '.' in inner_file else inner_file
            return f'{protocol}://{inner_file}::{urlpath}'
        else:
            return f'{protocol}://::{urlpath}'

    def download_and_extract(self, url_or_urls):
        return self.extract(self.download(url_or_urls))

    def iter_archive(self, urlpath_or_buf: Union[str, io.BufferedReader]) -> Iterable[tuple]:
        if hasattr(urlpath_or_buf, 'read'):
            return ArchiveIterable.from_buf(urlpath_or_buf)
        else:
            return ArchiveIterable.from_urlpath(urlpath_or_buf, download_config=self.download_config)

    def iter_files(self, urlpaths: Union[str, list[str]]) -> Iterable[str]:
        return FilesIterable.from_urlpaths(urlpaths, download_config=self.download_config)

    def manage_extracted_files(self):
        pass

    def get_recorded_sizes_checksums(self):
        pass
