import pyarrow as pa
from pyarrow.util import _is_iterable, _stringify_path, _is_path_like
try:
    from pyarrow._dataset import CsvFileFormat, CsvFragmentScanOptions, JsonFileFormat, JsonFragmentScanOptions, Dataset, DatasetFactory, DirectoryPartitioning, FeatherFileFormat, FilenamePartitioning, FileFormat, FileFragment, FileSystemDataset, FileSystemDatasetFactory, FileSystemFactoryOptions, FileWriteOptions, Fragment, FragmentScanOptions, HivePartitioning, IpcFileFormat, IpcFileWriteOptions, InMemoryDataset, Partitioning, PartitioningFactory, Scanner, TaggedRecordBatch, UnionDataset, UnionDatasetFactory, WrittenFile, get_partition_keys, get_partition_keys as _get_partition_keys, _filesystemdataset_write
except ImportError as exc:
    raise ImportError(f"The pyarrow installation is not built with support for 'dataset' ({str(exc)})") from None
from pyarrow.compute import Expression, scalar, field
_orc_available = False
_orc_msg = 'The pyarrow installation is not built with support for the ORC file format.'
try:
    from pyarrow._dataset_orc import OrcFileFormat
    _orc_available = True
except ImportError:
    pass
_parquet_available = False
_parquet_msg = 'The pyarrow installation is not built with support for the Parquet file format.'
try:
    from pyarrow._dataset_parquet import ParquetDatasetFactory, ParquetFactoryOptions, ParquetFileFormat, ParquetFileFragment, ParquetFileWriteOptions, ParquetFragmentScanOptions, ParquetReadOptions, RowGroupInfo
    _parquet_available = True
except ImportError:
    pass
try:
    from pyarrow._dataset_parquet_encryption import ParquetDecryptionConfig, ParquetEncryptionConfig
except ImportError:
    pass

def __getattr__(name):
    if name == 'OrcFileFormat' and (not _orc_available):
        raise ImportError(_orc_msg)
    if name == 'ParquetFileFormat' and (not _parquet_available):
        raise ImportError(_parquet_msg)
    raise AttributeError(f"module 'pyarrow.dataset' has no attribute '{name}'")

def partitioning(schema=None, field_names=None, flavor=None, dictionaries=None):
    if flavor is None:
        if schema is not None:
            if field_names is not None:
                raise ValueError("Cannot specify both 'schema' and 'field_names'")
            if dictionaries == 'infer':
                return DirectoryPartitioning.discover(schema=schema)
            return DirectoryPartitioning(schema, dictionaries)
        elif field_names is not None:
            if isinstance(field_names, list):
                return DirectoryPartitioning.discover(field_names)
            else:
                raise ValueError(f'Expected list of field names, got {type(field_names)}')
        else:
            raise ValueError('For the default directory flavor, need to specify a Schema or a list of field names')
    if flavor == 'filename':
        if schema is not None:
            if field_names is not None:
                raise ValueError("Cannot specify both 'schema' and 'field_names'")
            if dictionaries == 'infer':
                return FilenamePartitioning.discover(schema=schema)
            return FilenamePartitioning(schema, dictionaries)
        elif field_names is not None:
            if isinstance(field_names, list):
                return FilenamePartitioning.discover(field_names)
            else:
                raise ValueError(f'Expected list of field names, got {type(field_names)}')
        else:
            raise ValueError('For the filename flavor, need to specify a Schema or a list of field names')
    elif flavor == 'hive':
        if field_names is not None:
            raise ValueError("Cannot specify 'field_names' for flavor 'hive'")
        elif schema is not None:
            if isinstance(schema, pa.Schema):
                if dictionaries == 'infer':
                    return HivePartitioning.discover(schema=schema)
                return HivePartitioning(schema, dictionaries)
            else:
                raise ValueError(f"Expected Schema for 'schema', got {type(schema)}")
        else:
            return HivePartitioning.discover()
    else:
        raise ValueError('Unsupported flavor')

def _ensure_partitioning(scheme):
    if scheme is None:
        pass
    elif isinstance(scheme, str):
        scheme = partitioning(flavor=scheme)
    elif isinstance(scheme, list):
        scheme = partitioning(field_names=scheme)
    elif isinstance(scheme, (Partitioning, PartitioningFactory)):
        pass
    else:
        raise ValueError(f'Expected Partitioning or PartitioningFactory, got {type(scheme)}')
    return scheme

def _ensure_format(obj):
    if isinstance(obj, FileFormat):
        return obj
    elif obj == 'parquet':
        if not _parquet_available:
            raise ValueError(_parquet_msg)
        return ParquetFileFormat()
    elif obj in {'ipc', 'arrow'}:
        return IpcFileFormat()
    elif obj == 'feather':
        return FeatherFileFormat()
    elif obj == 'csv':
        return CsvFileFormat()
    elif obj == 'orc':
        if not _orc_available:
            raise ValueError(_orc_msg)
        return OrcFileFormat()
    elif obj == 'json':
        return JsonFileFormat()
    else:
        raise ValueError(f"format '{obj}' is not supported")

def _ensure_multiple_sources(paths, filesystem=None):
    from pyarrow.fs import LocalFileSystem, SubTreeFileSystem, _MockFileSystem, FileType, _ensure_filesystem
    if filesystem is None:
        filesystem = LocalFileSystem()
    else:
        filesystem = _ensure_filesystem(filesystem)
    is_local = isinstance(filesystem, (LocalFileSystem, _MockFileSystem)) or (isinstance(filesystem, SubTreeFileSystem) and isinstance(filesystem.base_fs, LocalFileSystem))
    paths = [filesystem.normalize_path(_stringify_path(p)) for p in paths]
    if is_local:
        for info in filesystem.get_file_info(paths):
            file_type = info.type
            if file_type == FileType.File:
                continue
            elif file_type == FileType.NotFound:
                raise FileNotFoundError(info.path)
            elif file_type == FileType.Directory:
                raise IsADirectoryError(f'Path {info.path} points to a directory, but only file paths are supported. To construct a nested or union dataset pass a list of dataset objects instead.')
            else:
                raise IOError(f'Path {info.path} exists but its type is unknown (could be a special file such as a Unix socket or character device, or Windows NUL / CON / ...)')
    return (filesystem, paths)

def _ensure_single_source(path, filesystem=None):
    from pyarrow.fs import FileType, FileSelector, _resolve_filesystem_and_path
    filesystem, path = _resolve_filesystem_and_path(path, filesystem)
    path = filesystem.normalize_path(path)
    file_info = filesystem.get_file_info(path)
    if file_info.type == FileType.Directory:
        paths_or_selector = FileSelector(path, recursive=True)
    elif file_info.type == FileType.File:
        paths_or_selector = [path]
    else:
        raise FileNotFoundError(path)
    return (filesystem, paths_or_selector)

def _filesystem_dataset(source, schema=None, filesystem=None, partitioning=None, format=None, partition_base_dir=None, exclude_invalid_files=None, selector_ignore_prefixes=None):
    from pyarrow.fs import LocalFileSystem, _ensure_filesystem, FileInfo
    format = _ensure_format(format or 'parquet')
    partitioning = _ensure_partitioning(partitioning)
    if isinstance(source, (list, tuple)):
        if source and isinstance(source[0], FileInfo):
            if filesystem is None:
                fs = LocalFileSystem()
            else:
                fs = _ensure_filesystem(filesystem)
            paths_or_selector = source
        else:
            fs, paths_or_selector = _ensure_multiple_sources(source, filesystem)
    else:
        fs, paths_or_selector = _ensure_single_source(source, filesystem)
    options = FileSystemFactoryOptions(partitioning=partitioning, partition_base_dir=partition_base_dir, exclude_invalid_files=exclude_invalid_files, selector_ignore_prefixes=selector_ignore_prefixes)
    factory = FileSystemDatasetFactory(fs, paths_or_selector, format, options)
    return factory.finish(schema)

def _in_memory_dataset(source, schema=None, **kwargs):
    if any((v is not None for v in kwargs.values())):
        raise ValueError('For in-memory datasets, you cannot pass any additional arguments')
    return InMemoryDataset(source, schema)

def _union_dataset(children, schema=None, **kwargs):
    if any((v is not None for v in kwargs.values())):
        raise ValueError('When passing a list of Datasets, you cannot pass any additional arguments')
    if schema is None:
        schema = pa.unify_schemas([child.schema for child in children])
    for child in children:
        if getattr(child, '_scan_options', None):
            raise ValueError('Creating an UnionDataset from filtered or projected Datasets is currently not supported. Union the unfiltered datasets and apply the filter to the resulting union.')
    children = [child.replace_schema(schema) for child in children]
    return UnionDataset(schema, children)

def parquet_dataset(metadata_path, schema=None, filesystem=None, format=None, partitioning=None, partition_base_dir=None):
    from pyarrow.fs import LocalFileSystem, _ensure_filesystem
    if format is None:
        format = ParquetFileFormat()
    elif not isinstance(format, ParquetFileFormat):
        raise ValueError('format argument must be a ParquetFileFormat')
    if filesystem is None:
        filesystem = LocalFileSystem()
    else:
        filesystem = _ensure_filesystem(filesystem)
    metadata_path = filesystem.normalize_path(_stringify_path(metadata_path))
    options = ParquetFactoryOptions(partition_base_dir=partition_base_dir, partitioning=_ensure_partitioning(partitioning))
    factory = ParquetDatasetFactory(metadata_path, filesystem, format, options=options)
    return factory.finish(schema)

def dataset(source, schema=None, format=None, filesystem=None, partitioning=None, partition_base_dir=None, exclude_invalid_files=None, ignore_prefixes=None):
    from pyarrow.fs import FileInfo
    kwargs = dict(schema=schema, filesystem=filesystem, partitioning=partitioning, format=format, partition_base_dir=partition_base_dir, exclude_invalid_files=exclude_invalid_files, selector_ignore_prefixes=ignore_prefixes)
    if _is_path_like(source):
        return _filesystem_dataset(source, **kwargs)
    elif isinstance(source, (tuple, list)):
        if all((_is_path_like(elem) or isinstance(elem, FileInfo) for elem in source)):
            return _filesystem_dataset(source, **kwargs)
        elif all((isinstance(elem, Dataset) for elem in source)):
            return _union_dataset(source, **kwargs)
        elif all((isinstance(elem, (pa.RecordBatch, pa.Table)) for elem in source)):
            return _in_memory_dataset(source, **kwargs)
        else:
            unique_types = set((type(elem).__name__ for elem in source))
            type_names = ', '.join((f'{t}' for t in unique_types))
            raise TypeError(f'Expected a list of path-like or dataset objects, or a list of batches or tables. The given list contains the following types: {type_names}')
    elif isinstance(source, (pa.RecordBatch, pa.Table, pa.RecordBatchReader)):
        return _in_memory_dataset(source, **kwargs)
    else:
        raise TypeError(f'Expected a path-like, list of path-likes or a list of Datasets instead of the given type: {type(source).__name__}')

def _ensure_write_partitioning(part, schema, flavor):
    if isinstance(part, PartitioningFactory):
        raise ValueError('A PartitioningFactory cannot be used. Did you call the partitioning function without supplying a schema?')
    if isinstance(part, Partitioning) and flavor:
        raise ValueError('Providing a partitioning_flavor with a Partitioning object is not supported')
    elif isinstance(part, (tuple, list)):
        part = partitioning(schema=pa.schema([schema.field(f) for f in part]), flavor=flavor)
    elif part is None:
        part = partitioning(pa.schema([]), flavor=flavor)
    if not isinstance(part, Partitioning):
        raise ValueError('partitioning must be a Partitioning object or a list of column names')
    return part

def write_dataset(data, base_dir, *, basename_template=None, format=None, partitioning=None, partitioning_flavor=None, schema=None, filesystem=None, file_options=None, use_threads=True, preserve_order=False, max_partitions=None, max_open_files=None, max_rows_per_file=None, min_rows_per_group=None, max_rows_per_group=None, file_visitor=None, existing_data_behavior='error', create_dir=True):
    from pyarrow.fs import _resolve_filesystem_and_path
    if isinstance(data, (list, tuple)):
        schema = schema or data[0].schema
        data = InMemoryDataset(data, schema=schema)
    elif isinstance(data, (pa.RecordBatch, pa.Table)):
        schema = schema or data.schema
        data = InMemoryDataset(data, schema=schema)
    elif isinstance(data, pa.ipc.RecordBatchReader) or hasattr(data, '__arrow_c_stream__') or _is_iterable(data):
        data = Scanner.from_batches(data, schema=schema)
        schema = None
    elif not isinstance(data, (Dataset, Scanner)):
        raise ValueError('Only Dataset, Scanner, Table/RecordBatch, RecordBatchReader, a list of Tables/RecordBatches, or iterable of batches are supported.')
    if format is None and isinstance(data, FileSystemDataset):
        format = data.format
    else:
        format = _ensure_format(format)
    if file_options is None:
        file_options = format.make_write_options()
    if format != file_options.format:
        raise TypeError(f"Supplied FileWriteOptions have format {format}, which doesn't match supplied FileFormat {file_options}")
    if basename_template is None:
        basename_template = 'part-{i}.' + format.default_extname
    if max_partitions is None:
        max_partitions = 1024
    if max_open_files is None:
        max_open_files = 1024
    if max_rows_per_file is None:
        max_rows_per_file = 0
    if max_rows_per_group is None:
        max_rows_per_group = 1 << 20
    if min_rows_per_group is None:
        min_rows_per_group = 0
    if isinstance(data, Scanner):
        partitioning_schema = data.projected_schema
    else:
        partitioning_schema = data.schema
    partitioning = _ensure_write_partitioning(partitioning, schema=partitioning_schema, flavor=partitioning_flavor)
    filesystem, base_dir = _resolve_filesystem_and_path(base_dir, filesystem)
    if isinstance(data, Dataset):
        scanner = data.scanner(use_threads=use_threads)
    else:
        if schema is not None:
            raise ValueError('Cannot specify a schema when writing a Scanner')
        scanner = data
    _filesystemdataset_write(scanner, base_dir, basename_template, filesystem, partitioning, preserve_order, file_options, max_partitions, file_visitor, existing_data_behavior, max_open_files, max_rows_per_file, min_rows_per_group, max_rows_per_group, create_dir)
