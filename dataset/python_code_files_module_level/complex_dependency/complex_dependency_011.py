from pyarrow.util import _is_path_like, _stringify_path
from pyarrow._fs import FileSelector, FileType, FileInfo, FileSystem, LocalFileSystem, SubTreeFileSystem, _MockFileSystem, FileSystemHandler, PyFileSystem, _copy_files, _copy_files_selector
FileStats = FileInfo
_not_imported = []
try:
    from pyarrow._azurefs import AzureFileSystem
except ImportError:
    _not_imported.append('AzureFileSystem')
try:
    from pyarrow._hdfs import HadoopFileSystem
except ImportError:
    _not_imported.append('HadoopFileSystem')
try:
    from pyarrow._gcsfs import GcsFileSystem
except ImportError:
    _not_imported.append('GcsFileSystem')
try:
    from pyarrow._s3fs import AwsDefaultS3RetryStrategy, AwsStandardS3RetryStrategy, S3FileSystem, S3LogLevel, S3RetryStrategy, ensure_s3_initialized, finalize_s3, ensure_s3_finalized, initialize_s3, resolve_s3_region
except ImportError:
    _not_imported.append('S3FileSystem')
else:
    import atexit
    atexit.register(ensure_s3_finalized)

def __getattr__(name):
    if name in _not_imported:
        raise ImportError(f"The pyarrow installation is not built with support for '{name}'")
    raise AttributeError(f"module 'pyarrow.fs' has no attribute '{name}'")

def _ensure_filesystem(filesystem, *, use_mmap=False):
    if isinstance(filesystem, FileSystem):
        return filesystem
    elif isinstance(filesystem, str):
        if use_mmap:
            raise ValueError('Specifying to use memory mapping not supported for filesystem specified as an URI string')
        fs, path = FileSystem.from_uri(filesystem)
        prefix = fs.normalize_path(path)
        if prefix:
            prefix_info = fs.get_file_info([prefix])[0]
            if prefix_info.type != FileType.Directory:
                raise ValueError(f'The path component of the filesystem URI must point to a directory but it has a type: `{prefix_info.type.name}`. The path component is `{prefix_info.path}` and the given filesystem URI is `{filesystem}`')
            fs = SubTreeFileSystem(prefix, fs)
        return fs
    else:
        try:
            import fsspec
        except ImportError:
            pass
        else:
            if isinstance(filesystem, fsspec.AbstractFileSystem):
                if type(filesystem).__name__ == 'LocalFileSystem':
                    return LocalFileSystem(use_mmap=use_mmap)
                return PyFileSystem(FSSpecHandler(filesystem))
        raise TypeError(f'Unrecognized filesystem: {type(filesystem)}. `filesystem` argument must be a FileSystem instance or a valid file system URI')

def _resolve_filesystem_and_path(path, filesystem=None, *, memory_map=False):
    if not _is_path_like(path):
        if filesystem is not None:
            raise ValueError("'filesystem' passed but the specified path is file-like, so there is nothing to open with 'filesystem'.")
        return (filesystem, path)
    if filesystem is not None:
        filesystem = _ensure_filesystem(filesystem, use_mmap=memory_map)
        if isinstance(filesystem, LocalFileSystem):
            path = _stringify_path(path)
        elif not isinstance(path, str):
            raise TypeError('Expected string path; path-like objects are only allowed with a local filesystem')
        path = filesystem.normalize_path(path)
        return (filesystem, path)
    path = _stringify_path(path)
    filesystem = LocalFileSystem(use_mmap=memory_map)
    try:
        file_info = filesystem.get_file_info(path)
    except ValueError:
        file_info = None
        exists_locally = False
    else:
        exists_locally = file_info.type != FileType.NotFound
    if not exists_locally:
        try:
            filesystem, path = FileSystem.from_uri(path)
        except ValueError as e:
            msg = str(e)
            if 'empty scheme' in msg or 'Cannot parse URI' in msg:
                pass
            else:
                raise e
    else:
        path = filesystem.normalize_path(path)
    return (filesystem, path)

def copy_files(source, destination, source_filesystem=None, destination_filesystem=None, *, chunk_size=1024 * 1024, use_threads=True):
    source_fs, source_path = _resolve_filesystem_and_path(source, source_filesystem)
    destination_fs, destination_path = _resolve_filesystem_and_path(destination, destination_filesystem)
    file_info = source_fs.get_file_info(source_path)
    if file_info.type == FileType.Directory:
        source_sel = FileSelector(source_path, recursive=True)
        _copy_files_selector(source_fs, source_sel, destination_fs, destination_path, chunk_size, use_threads)
    else:
        _copy_files(source_fs, source_path, destination_fs, destination_path, chunk_size, use_threads)

class FSSpecHandler(FileSystemHandler):

    def __init__(self, fs):
        self.fs = fs

    def __eq__(self, other):
        if isinstance(other, FSSpecHandler):
            return self.fs == other.fs
        return NotImplemented

    def __ne__(self, other):
        if isinstance(other, FSSpecHandler):
            return self.fs != other.fs
        return NotImplemented

    def get_type_name(self):
        protocol = self.fs.protocol
        if isinstance(protocol, list):
            protocol = protocol[0]
        return f'fsspec+{protocol}'

    def normalize_path(self, path):
        return path

    @staticmethod
    def _create_file_info(path, info):
        size = info['size']
        if info['type'] == 'file':
            ftype = FileType.File
        elif info['type'] == 'directory':
            ftype = FileType.Directory
            size = None
        else:
            ftype = FileType.Unknown
        return FileInfo(path, ftype, size=size, mtime=info.get('mtime', None))

    def get_file_info(self, paths):
        infos = []
        for path in paths:
            try:
                info = self.fs.info(path)
            except FileNotFoundError:
                infos.append(FileInfo(path, FileType.NotFound))
            else:
                infos.append(self._create_file_info(path, info))
        return infos

    def get_file_info_selector(self, selector):
        if not self.fs.isdir(selector.base_dir):
            if self.fs.exists(selector.base_dir):
                raise NotADirectoryError(selector.base_dir)
            elif selector.allow_not_found:
                return []
            else:
                raise FileNotFoundError(selector.base_dir)
        if selector.recursive:
            maxdepth = None
        else:
            maxdepth = 1
        infos = []
        selected_files = self.fs.find(selector.base_dir, maxdepth=maxdepth, withdirs=True, detail=True)
        for path, info in selected_files.items():
            _path = path.strip('/')
            base_dir = selector.base_dir.strip('/')
            if _path != base_dir:
                infos.append(self._create_file_info(path, info))
        return infos

    def create_dir(self, path, recursive):
        try:
            self.fs.mkdir(path, create_parents=recursive)
        except FileExistsError:
            pass

    def delete_dir(self, path):
        self.fs.rm(path, recursive=True)

    def _delete_dir_contents(self, path, missing_dir_ok):
        try:
            subpaths = self.fs.listdir(path, detail=False)
        except FileNotFoundError:
            if missing_dir_ok:
                return
            raise
        for subpath in subpaths:
            if self.fs.isdir(subpath):
                self.fs.rm(subpath, recursive=True)
            elif self.fs.isfile(subpath):
                self.fs.rm(subpath)

    def delete_dir_contents(self, path, missing_dir_ok):
        if path.strip('/') == '':
            raise ValueError("delete_dir_contents called on path '", path, "'")
        self._delete_dir_contents(path, missing_dir_ok)

    def delete_root_dir_contents(self):
        self._delete_dir_contents('/', False)

    def delete_file(self, path):
        if not self.fs.exists(path):
            raise FileNotFoundError(path)
        self.fs.rm(path)

    def move(self, src, dest):
        self.fs.mv(src, dest, recursive=True)

    def copy_file(self, src, dest):
        self.fs.copy(src, dest)

    def open_input_stream(self, path):
        from pyarrow import PythonFile
        if not self.fs.isfile(path):
            raise FileNotFoundError(path)
        return PythonFile(self.fs.open(path, mode='rb'), mode='r')

    def open_input_file(self, path):
        from pyarrow import PythonFile
        if not self.fs.isfile(path):
            raise FileNotFoundError(path)
        return PythonFile(self.fs.open(path, mode='rb'), mode='r')

    def open_output_stream(self, path, metadata):
        from pyarrow import PythonFile
        return PythonFile(self.fs.open(path, mode='wb'), mode='w')

    def open_append_stream(self, path, metadata):
        from pyarrow import PythonFile
        return PythonFile(self.fs.open(path, mode='ab'), mode='w')
