import base64
import hashlib
import logging
import os
import time
from dataclasses import dataclass
from pathlib import Path
from typing import Optional
from .utils import WeakFileLock
logger = logging.getLogger(__name__)

@dataclass
class LocalDownloadFilePaths:
    file_path: Path
    lock_path: Path
    metadata_path: Path

    def incomplete_path(self, etag: str) -> Path:
        path = self.metadata_path.parent / f'{_short_hash(self.metadata_path.name)}.{etag}.incomplete'
        resolved_path = str(path.resolve())
        if os.name == 'nt' and len(resolved_path) > 255 and (not resolved_path.startswith('\\\\?\\')):
            path = Path('\\\\?\\' + resolved_path)
        return path

@dataclass(frozen=True)
class LocalUploadFilePaths:
    path_in_repo: str
    file_path: Path
    lock_path: Path
    metadata_path: Path

@dataclass
class LocalDownloadFileMetadata:
    filename: str
    commit_hash: str
    etag: str
    timestamp: float

@dataclass
class LocalUploadFileMetadata:
    size: int
    timestamp: Optional[float] = None
    should_ignore: Optional[bool] = None
    sha256: Optional[str] = None
    upload_mode: Optional[str] = None
    remote_oid: Optional[str] = None
    is_uploaded: bool = False
    is_committed: bool = False

    def save(self, paths: LocalUploadFilePaths) -> None:
        with WeakFileLock(paths.lock_path):
            with paths.metadata_path.open('w') as f:
                new_timestamp = time.time()
                f.write(str(new_timestamp) + '\n')
                f.write(str(self.size))
                f.write('\n')
                if self.should_ignore is not None:
                    f.write(str(int(self.should_ignore)))
                f.write('\n')
                if self.sha256 is not None:
                    f.write(self.sha256)
                f.write('\n')
                if self.upload_mode is not None:
                    f.write(self.upload_mode)
                f.write('\n')
                if self.remote_oid is not None:
                    f.write(self.remote_oid)
                f.write('\n')
                f.write(str(int(self.is_uploaded)) + '\n')
                f.write(str(int(self.is_committed)) + '\n')
            self.timestamp = new_timestamp

def get_local_download_paths(local_dir: Path, filename: str) -> LocalDownloadFilePaths:
    sanitized_filename = os.path.join(*filename.split('/'))
    if os.name == 'nt':
        if sanitized_filename.startswith('..\\') or '\\..\\' in sanitized_filename:
            raise ValueError(f"Invalid filename: cannot handle filename '{sanitized_filename}' on Windows. Please ask the repository owner to rename this file.")
    file_path = local_dir / sanitized_filename
    metadata_path = _huggingface_dir(local_dir) / 'download' / f'{sanitized_filename}.metadata'
    lock_path = metadata_path.with_suffix('.lock')
    if os.name == 'nt':
        if not str(local_dir).startswith('\\\\?\\') and len(os.path.abspath(lock_path)) > 255:
            file_path = Path('\\\\?\\' + os.path.abspath(file_path))
            lock_path = Path('\\\\?\\' + os.path.abspath(lock_path))
            metadata_path = Path('\\\\?\\' + os.path.abspath(metadata_path))
    file_path.parent.mkdir(parents=True, exist_ok=True)
    metadata_path.parent.mkdir(parents=True, exist_ok=True)
    return LocalDownloadFilePaths(file_path=file_path, lock_path=lock_path, metadata_path=metadata_path)

def get_local_upload_paths(local_dir: Path, filename: str) -> LocalUploadFilePaths:
    sanitized_filename = os.path.join(*filename.split('/'))
    if os.name == 'nt':
        if sanitized_filename.startswith('..\\') or '\\..\\' in sanitized_filename:
            raise ValueError(f"Invalid filename: cannot handle filename '{sanitized_filename}' on Windows. Please ask the repository owner to rename this file.")
    file_path = local_dir / sanitized_filename
    metadata_path = _huggingface_dir(local_dir) / 'upload' / f'{sanitized_filename}.metadata'
    lock_path = metadata_path.with_suffix('.lock')
    if os.name == 'nt':
        if not str(local_dir).startswith('\\\\?\\') and len(os.path.abspath(lock_path)) > 255:
            file_path = Path('\\\\?\\' + os.path.abspath(file_path))
            lock_path = Path('\\\\?\\' + os.path.abspath(lock_path))
            metadata_path = Path('\\\\?\\' + os.path.abspath(metadata_path))
    file_path.parent.mkdir(parents=True, exist_ok=True)
    metadata_path.parent.mkdir(parents=True, exist_ok=True)
    return LocalUploadFilePaths(path_in_repo=filename, file_path=file_path, lock_path=lock_path, metadata_path=metadata_path)

def read_download_metadata(local_dir: Path, filename: str) -> Optional[LocalDownloadFileMetadata]:
    paths = get_local_download_paths(local_dir, filename)
    with WeakFileLock(paths.lock_path):
        if paths.metadata_path.exists():
            try:
                with paths.metadata_path.open() as f:
                    commit_hash = f.readline().strip()
                    etag = f.readline().strip()
                    timestamp = float(f.readline().strip())
                    metadata = LocalDownloadFileMetadata(filename=filename, commit_hash=commit_hash, etag=etag, timestamp=timestamp)
            except Exception as e:
                logger.warning(f'Invalid metadata file {paths.metadata_path}: {e}. Removing it from disk and continue.')
                try:
                    paths.metadata_path.unlink()
                except Exception as e:
                    logger.warning(f'Could not remove corrupted metadata file {paths.metadata_path}: {e}')
                return None
            try:
                stat = paths.file_path.stat()
                if stat.st_mtime - 1 <= metadata.timestamp:
                    return metadata
                logger.info(f"Ignored metadata for '{filename}' (outdated). Will re-compute hash.")
            except FileNotFoundError:
                return None
    return None

def read_upload_metadata(local_dir: Path, filename: str) -> LocalUploadFileMetadata:
    paths = get_local_upload_paths(local_dir, filename)
    with WeakFileLock(paths.lock_path):
        if paths.metadata_path.exists():
            try:
                with paths.metadata_path.open() as f:
                    timestamp = float(f.readline().strip())
                    size = int(f.readline().strip())
                    _should_ignore = f.readline().strip()
                    should_ignore = None if _should_ignore == '' else bool(int(_should_ignore))
                    _sha256 = f.readline().strip()
                    sha256 = None if _sha256 == '' else _sha256
                    _upload_mode = f.readline().strip()
                    upload_mode = None if _upload_mode == '' else _upload_mode
                    if upload_mode not in (None, 'regular', 'lfs'):
                        raise ValueError(f'Invalid upload mode in metadata {paths.path_in_repo}: {upload_mode}')
                    _remote_oid = f.readline().strip()
                    remote_oid = None if _remote_oid == '' else _remote_oid
                    is_uploaded = bool(int(f.readline().strip()))
                    is_committed = bool(int(f.readline().strip()))
                    metadata = LocalUploadFileMetadata(timestamp=timestamp, size=size, should_ignore=should_ignore, sha256=sha256, upload_mode=upload_mode, remote_oid=remote_oid, is_uploaded=is_uploaded, is_committed=is_committed)
            except Exception as e:
                logger.warning(f'Invalid metadata file {paths.metadata_path}: {e}. Removing it from disk and continue.')
                try:
                    paths.metadata_path.unlink()
                except Exception as e:
                    logger.warning(f'Could not remove corrupted metadata file {paths.metadata_path}: {e}')
                return LocalUploadFileMetadata(size=paths.file_path.stat().st_size)
            if metadata.timestamp is not None and metadata.is_uploaded and (not metadata.is_committed) and (time.time() - metadata.timestamp > 20 * 3600):
                metadata.is_uploaded = False
            try:
                if metadata.timestamp is not None and paths.file_path.stat().st_mtime <= metadata.timestamp:
                    return metadata
                logger.info(f"Ignored metadata for '{filename}' (outdated). Will re-compute hash.")
            except FileNotFoundError:
                pass
    return LocalUploadFileMetadata(size=paths.file_path.stat().st_size)

def write_download_metadata(local_dir: Path, filename: str, commit_hash: str, etag: str) -> None:
    paths = get_local_download_paths(local_dir, filename)
    with WeakFileLock(paths.lock_path):
        with paths.metadata_path.open('w') as f:
            f.write(f'{commit_hash}\n{etag}\n{time.time()}\n')

def _huggingface_dir(local_dir: Path) -> Path:
    path = local_dir / '.cache' / 'huggingface'
    path.mkdir(exist_ok=True, parents=True)
    gitignore = path / '.gitignore'
    gitignore_lock = path / '.gitignore.lock'
    if not gitignore.exists():
        try:
            with WeakFileLock(gitignore_lock, timeout=0.1):
                gitignore.write_text('*')
        except IndexError:
            pass
        except OSError:
            pass
        try:
            gitignore_lock.unlink()
        except OSError:
            pass
    return path

def _short_hash(filename: str) -> str:
    return base64.urlsafe_b64encode(hashlib.sha1(filename.encode()).digest()).decode()
