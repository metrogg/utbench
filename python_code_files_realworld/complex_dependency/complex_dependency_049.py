import atexit
import logging
import os
import time
from concurrent.futures import Future
from dataclasses import dataclass
from io import SEEK_END, SEEK_SET, BytesIO
from pathlib import Path
from threading import Lock, Thread
from typing import Optional, Union
from .hf_api import DEFAULT_IGNORE_PATTERNS, CommitInfo, CommitOperationAdd, HfApi
from .utils import filter_repo_objects
logger = logging.getLogger(__name__)

@dataclass(frozen=True)
class _FileToUpload:
    local_path: Path
    path_in_repo: str
    size_limit: int
    last_modified: float

class CommitScheduler:

    def __init__(self, *, repo_id: str, folder_path: Union[str, Path], every: Union[int, float]=5, path_in_repo: Optional[str]=None, repo_type: Optional[str]=None, revision: Optional[str]=None, private: Optional[bool]=None, token: Optional[str]=None, allow_patterns: Optional[Union[list[str], str]]=None, ignore_patterns: Optional[Union[list[str], str]]=None, squash_history: bool=False, hf_api: Optional['HfApi']=None) -> None:
        self.api = hf_api or HfApi(token=token)
        self.folder_path = Path(folder_path).expanduser().resolve()
        self.path_in_repo = path_in_repo or ''
        self.allow_patterns = allow_patterns
        if ignore_patterns is None:
            ignore_patterns = []
        elif isinstance(ignore_patterns, str):
            ignore_patterns = [ignore_patterns]
        self.ignore_patterns = ignore_patterns + DEFAULT_IGNORE_PATTERNS
        if self.folder_path.is_file():
            raise ValueError(f"'folder_path' must be a directory, not a file: '{self.folder_path}'.")
        self.folder_path.mkdir(parents=True, exist_ok=True)
        repo_url = self.api.create_repo(repo_id=repo_id, private=private, repo_type=repo_type, exist_ok=True)
        self.repo_id = repo_url.repo_id
        self.repo_type = repo_type
        self.revision = revision
        self.token = token
        self.last_uploaded: dict[Path, float] = {}
        if not every > 0:
            raise ValueError(f"'every' must be a positive integer, not '{every}'.")
        self.lock = Lock()
        self.every = every
        self.squash_history = squash_history
        logger.info(f"Scheduled job to push '{self.folder_path}' to '{self.repo_id}' every {self.every} minutes.")
        self._scheduler_thread = Thread(target=self._run_scheduler, daemon=True)
        self._scheduler_thread.start()
        atexit.register(self._push_to_hub)
        self.__stopped = False

    def stop(self) -> None:
        self.__stopped = True

    def __enter__(self) -> 'CommitScheduler':
        return self

    def __exit__(self, exc_type, exc_value, traceback) -> None:
        self.trigger().result()
        self.stop()
        return

    def _run_scheduler(self) -> None:
        while True:
            self.last_future = self.trigger()
            time.sleep(self.every * 60)
            if self.__stopped:
                break

    def trigger(self) -> Future:
        return self.api.run_as_future(self._push_to_hub)

    def _push_to_hub(self) -> Optional[CommitInfo]:
        if self.__stopped:
            return None
        logger.info('(Background) scheduled commit triggered.')
        try:
            value = self.push_to_hub()
            if self.squash_history:
                logger.info('(Background) squashing repo history.')
                self.api.super_squash_history(repo_id=self.repo_id, repo_type=self.repo_type, branch=self.revision)
            return value
        except Exception as e:
            logger.error(f'Error while pushing to Hub: {e}')
            raise

    def push_to_hub(self) -> Optional[CommitInfo]:
        with self.lock:
            logger.debug('Listing files to upload for scheduled commit.')
            relpath_to_abspath = {path.relative_to(self.folder_path).as_posix(): path for path in sorted(self.folder_path.glob('**/*')) if path.is_file()}
            prefix = f"{self.path_in_repo.strip('/')}/" if self.path_in_repo else ''
            files_to_upload: list[_FileToUpload] = []
            for relpath in filter_repo_objects(relpath_to_abspath.keys(), allow_patterns=self.allow_patterns, ignore_patterns=self.ignore_patterns):
                local_path = relpath_to_abspath[relpath]
                stat = local_path.stat()
                if self.last_uploaded.get(local_path) is None or self.last_uploaded[local_path] != stat.st_mtime:
                    files_to_upload.append(_FileToUpload(local_path=local_path, path_in_repo=prefix + relpath, size_limit=stat.st_size, last_modified=stat.st_mtime))
        if len(files_to_upload) == 0:
            logger.debug('Dropping schedule commit: no changed file to upload.')
            return None
        logger.debug('Removing unchanged files since previous scheduled commit.')
        add_operations = [CommitOperationAdd(path_or_fileobj=PartialFileIO(file_to_upload.local_path, size_limit=file_to_upload.size_limit), path_in_repo=file_to_upload.path_in_repo) for file_to_upload in files_to_upload]
        logger.debug('Uploading files for scheduled commit.')
        commit_info = self.api.create_commit(repo_id=self.repo_id, repo_type=self.repo_type, operations=add_operations, commit_message='Scheduled Commit', revision=self.revision)
        for file in files_to_upload:
            self.last_uploaded[file.local_path] = file.last_modified
        return commit_info

class PartialFileIO(BytesIO):

    def __init__(self, file_path: Union[str, Path], size_limit: int) -> None:
        self._file_path = Path(file_path)
        self._file = self._file_path.open('rb')
        self._size_limit = min(size_limit, os.fstat(self._file.fileno()).st_size)

    def __del__(self) -> None:
        self._file.close()
        return super().__del__()

    def __repr__(self) -> str:
        return f'<PartialFileIO file_path={self._file_path} size_limit={self._size_limit}>'

    def __len__(self) -> int:
        return self._size_limit

    def __getattribute__(self, name: str):
        if name.startswith('_') or name in ('read', 'tell', 'seek', 'fileno'):
            return super().__getattribute__(name)
        raise NotImplementedError(f"PartialFileIO does not support '{name}'.")

    def fileno(self):
        raise AttributeError('PartialFileIO does not have a fileno.')

    def tell(self) -> int:
        return self._file.tell()

    def seek(self, __offset: int, __whence: int=SEEK_SET) -> int:
        if __whence == SEEK_END:
            __offset = len(self) + __offset
            __whence = SEEK_SET
        pos = self._file.seek(__offset, __whence)
        if pos > self._size_limit:
            return self._file.seek(self._size_limit)
        return pos

    def read(self, __size: Optional[int]=-1) -> bytes:
        current = self._file.tell()
        if __size is None or __size < 0:
            truncated_size = self._size_limit - current
        else:
            truncated_size = min(__size, self._size_limit - current)
        return self._file.read(truncated_size)
