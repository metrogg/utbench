import json
import logging
import mmap
import os
import shutil
import zipfile
from contextlib import contextmanager
from dataclasses import dataclass, field
from pathlib import Path
from typing import Any, Generator, Iterable, Union
from ..errors import DDUFCorruptedFileError, DDUFExportError, DDUFInvalidEntryNameError
logger = logging.getLogger(__name__)
DDUF_ALLOWED_ENTRIES = {'.json', '.model', '.safetensors', '.txt'}
DDUF_FOLDER_REQUIRED_ENTRIES = {'config.json', 'tokenizer_config.json', 'preprocessor_config.json', 'scheduler_config.json'}

@dataclass
class DDUFEntry:
    filename: str
    length: int
    offset: int
    dduf_path: Path = field(repr=False)

    @contextmanager
    def as_mmap(self) -> Generator[bytes, None, None]:
        with self.dduf_path.open('rb') as f:
            with mmap.mmap(f.fileno(), length=0, access=mmap.ACCESS_READ) as mm:
                yield mm[self.offset:self.offset + self.length]

    def read_text(self, encoding: str='utf-8') -> str:
        with self.dduf_path.open('rb') as f:
            f.seek(self.offset)
            return f.read(self.length).decode(encoding=encoding)

def read_dduf_file(dduf_path: Union[os.PathLike, str]) -> dict[str, DDUFEntry]:
    entries = {}
    dduf_path = Path(dduf_path)
    logger.info(f'Reading DDUF file {dduf_path}')
    with zipfile.ZipFile(str(dduf_path), 'r') as zf:
        for info in zf.infolist():
            logger.debug(f'Reading entry {info.filename}')
            if info.compress_type != zipfile.ZIP_STORED:
                raise DDUFCorruptedFileError('Data must not be compressed in DDUF file.')
            try:
                _validate_dduf_entry_name(info.filename)
            except DDUFInvalidEntryNameError as e:
                raise DDUFCorruptedFileError(f'Invalid entry name in DDUF file: {info.filename}') from e
            offset = _get_data_offset(zf, info)
            entries[info.filename] = DDUFEntry(filename=info.filename, offset=offset, length=info.file_size, dduf_path=dduf_path)
    if 'model_index.json' not in entries:
        raise DDUFCorruptedFileError("Missing required 'model_index.json' entry in DDUF file.")
    index = json.loads(entries['model_index.json'].read_text())
    _validate_dduf_structure(index, entries.keys())
    logger.info(f'Done reading DDUF file {dduf_path}. Found {len(entries)} entries')
    return entries

def export_entries_as_dduf(dduf_path: Union[str, os.PathLike], entries: Iterable[tuple[str, Union[str, Path, bytes]]]) -> None:
    logger.info(f"Exporting DDUF file '{dduf_path}'")
    filenames = set()
    index = None
    with zipfile.ZipFile(str(dduf_path), 'w', zipfile.ZIP_STORED) as archive:
        for filename, content in entries:
            if filename in filenames:
                raise DDUFExportError(f"Can't add duplicate entry: {filename}")
            filenames.add(filename)
            if filename == 'model_index.json':
                try:
                    index = json.loads(_load_content(content).decode())
                except json.JSONDecodeError as e:
                    raise DDUFExportError("Failed to parse 'model_index.json'.") from e
            try:
                filename = _validate_dduf_entry_name(filename)
            except DDUFInvalidEntryNameError as e:
                raise DDUFExportError(f'Invalid entry name: {filename}') from e
            logger.debug(f"Adding entry '{filename}' to DDUF file")
            _dump_content_in_archive(archive, filename, content)
    if index is None:
        raise DDUFExportError("Missing required 'model_index.json' entry in DDUF file.")
    try:
        _validate_dduf_structure(index, filenames)
    except DDUFCorruptedFileError as e:
        raise DDUFExportError('Invalid DDUF file structure.') from e
    logger.info(f'Done writing DDUF file {dduf_path}')

def export_folder_as_dduf(dduf_path: Union[str, os.PathLike], folder_path: Union[str, os.PathLike]) -> None:
    folder_path = Path(folder_path)

    def _iterate_over_folder() -> Iterable[tuple[str, Path]]:
        for path in Path(folder_path).glob('**/*'):
            if not path.is_file():
                continue
            if path.suffix not in DDUF_ALLOWED_ENTRIES:
                logger.debug(f"Skipping file '{path}' (file type not allowed)")
                continue
            path_in_archive = path.relative_to(folder_path)
            if len(path_in_archive.parts) >= 3:
                logger.debug(f"Skipping file '{path}' (nested directories not allowed)")
                continue
            yield (path_in_archive.as_posix(), path)
    export_entries_as_dduf(dduf_path, _iterate_over_folder())

def _dump_content_in_archive(archive: zipfile.ZipFile, filename: str, content: Union[str, os.PathLike, bytes]) -> None:
    with archive.open(filename, 'w', force_zip64=True) as archive_fh:
        if isinstance(content, (str, Path)):
            content_path = Path(content)
            with content_path.open('rb') as content_fh:
                shutil.copyfileobj(content_fh, archive_fh, 1024 * 1024 * 8)
        elif isinstance(content, bytes):
            archive_fh.write(content)
        else:
            raise DDUFExportError(f'Invalid content type for {filename}. Must be str, Path or bytes.')

def _load_content(content: Union[str, Path, bytes]) -> bytes:
    if isinstance(content, (str, Path)):
        return Path(content).read_bytes()
    elif isinstance(content, bytes):
        return content
    else:
        raise DDUFExportError(f'Invalid content type. Must be str, Path or bytes. Got {type(content)}.')

def _validate_dduf_entry_name(entry_name: str) -> str:
    if '.' + entry_name.split('.')[-1] not in DDUF_ALLOWED_ENTRIES:
        raise DDUFInvalidEntryNameError(f'File type not allowed: {entry_name}')
    if '\\' in entry_name:
        raise DDUFInvalidEntryNameError(f"Entry names must use UNIX separators ('/'). Got {entry_name}.")
    entry_name = entry_name.strip('/')
    if entry_name.count('/') > 1:
        raise DDUFInvalidEntryNameError(f'DDUF only supports 1 level of directory. Got {entry_name}.')
    return entry_name

def _validate_dduf_structure(index: Any, entry_names: Iterable[str]) -> None:
    if not isinstance(index, dict):
        raise DDUFCorruptedFileError(f"Invalid 'model_index.json' content. Must be a dictionary. Got {type(index)}.")
    dduf_folders = {entry.split('/')[0] for entry in entry_names if '/' in entry}
    for folder in dduf_folders:
        if folder not in index:
            raise DDUFCorruptedFileError(f"Missing required entry '{folder}' in 'model_index.json'.")
        if not any((f'{folder}/{required_entry}' in entry_names for required_entry in DDUF_FOLDER_REQUIRED_ENTRIES)):
            raise DDUFCorruptedFileError(f"Missing required file in folder '{folder}'. Must contains at least one of {DDUF_FOLDER_REQUIRED_ENTRIES}.")

def _get_data_offset(zf: zipfile.ZipFile, info: zipfile.ZipInfo) -> int:
    if zf.fp is None:
        raise DDUFCorruptedFileError('ZipFile object must be opened in read mode.')
    header_offset = info.header_offset
    zf.fp.seek(header_offset)
    local_file_header = zf.fp.read(30)
    if len(local_file_header) < 30:
        raise DDUFCorruptedFileError('Incomplete local file header.')
    filename_len = int.from_bytes(local_file_header[26:28], 'little')
    extra_field_len = int.from_bytes(local_file_header[28:30], 'little')
    data_offset = header_offset + 30 + filename_len + extra_field_len
    return data_offset
