import base64
import io
import os
import warnings
from collections import defaultdict
from contextlib import contextmanager
from dataclasses import dataclass, field
from itertools import groupby
from pathlib import Path, PurePosixPath
from typing import TYPE_CHECKING, Any, BinaryIO, Iterable, Iterator, Literal, Optional, Union
from tqdm.contrib.concurrent import thread_map
from . import constants
from .errors import EntryNotFoundError, HfHubHTTPError, XetAuthorizationError, XetRefreshTokenError
from .file_download import hf_hub_url
from .lfs import UploadInfo, lfs_upload, post_lfs_batch_info
from .utils import FORBIDDEN_FOLDERS, XetTokenType, are_progress_bars_disabled, chunk_iterable, fetch_xet_connection_info_from_repo_info, get_session, hf_raise_for_status, http_backoff, logging, sha, tqdm_stream_file, validate_hf_hub_args
from .utils import tqdm as hf_tqdm
from .utils._runtime import is_xet_available
if TYPE_CHECKING:
    from .hf_api import RepoFile
logger = logging.get_logger(__name__)
UploadMode = Literal['lfs', 'regular']
FETCH_LFS_BATCH_SIZE = 500
UPLOAD_BATCH_MAX_NUM_FILES = 256

@dataclass
class CommitOperationDelete:
    path_in_repo: str
    is_folder: Union[bool, Literal['auto']] = 'auto'

    def __post_init__(self):
        self.path_in_repo = _validate_path_in_repo(self.path_in_repo)
        if self.is_folder == 'auto':
            self.is_folder = self.path_in_repo.endswith('/')
        if not isinstance(self.is_folder, bool):
            raise ValueError(f"Wrong value for `is_folder`. Must be one of [`True`, `False`, `'auto'`]. Got '{self.is_folder}'.")

@dataclass
class CommitOperationCopy:
    src_path_in_repo: str
    path_in_repo: str
    src_revision: Optional[str] = None
    _src_oid: Optional[str] = None
    _dest_oid: Optional[str] = None

    def __post_init__(self):
        self.src_path_in_repo = _validate_path_in_repo(self.src_path_in_repo)
        self.path_in_repo = _validate_path_in_repo(self.path_in_repo)

@dataclass
class CommitOperationAdd:
    path_in_repo: str
    path_or_fileobj: Union[str, Path, bytes, BinaryIO]
    upload_info: UploadInfo = field(init=False, repr=False)
    _upload_mode: Optional[UploadMode] = field(init=False, repr=False, default=None)
    _should_ignore: Optional[bool] = field(init=False, repr=False, default=None)
    _remote_oid: Optional[str] = field(init=False, repr=False, default=None)
    _is_uploaded: bool = field(init=False, repr=False, default=False)
    _is_committed: bool = field(init=False, repr=False, default=False)

    def __post_init__(self) -> None:
        self.path_in_repo = _validate_path_in_repo(self.path_in_repo)
        if isinstance(self.path_or_fileobj, Path):
            self.path_or_fileobj = str(self.path_or_fileobj)
        if isinstance(self.path_or_fileobj, str):
            path_or_fileobj = os.path.normpath(os.path.expanduser(self.path_or_fileobj))
            if not os.path.isfile(path_or_fileobj):
                raise ValueError(f"Provided path: '{path_or_fileobj}' is not a file on the local file system")
        elif not isinstance(self.path_or_fileobj, (io.BufferedIOBase, bytes)):
            raise ValueError('path_or_fileobj must be either an instance of str, bytes or io.BufferedIOBase. If you passed a file-like object, make sure it is in binary mode.')
        if isinstance(self.path_or_fileobj, io.BufferedIOBase):
            try:
                self.path_or_fileobj.tell()
                self.path_or_fileobj.seek(0, os.SEEK_CUR)
            except (OSError, AttributeError) as exc:
                raise ValueError('path_or_fileobj is a file-like object but does not implement seek() and tell()') from exc
        if isinstance(self.path_or_fileobj, str):
            self.upload_info = UploadInfo.from_path(self.path_or_fileobj)
        elif isinstance(self.path_or_fileobj, bytes):
            self.upload_info = UploadInfo.from_bytes(self.path_or_fileobj)
        else:
            self.upload_info = UploadInfo.from_fileobj(self.path_or_fileobj)

    @contextmanager
    def as_file(self, with_tqdm: bool=False) -> Iterator[BinaryIO]:
        if isinstance(self.path_or_fileobj, str) or isinstance(self.path_or_fileobj, Path):
            if with_tqdm:
                with tqdm_stream_file(self.path_or_fileobj) as file:
                    yield file
            else:
                with open(self.path_or_fileobj, 'rb') as file:
                    yield file
        elif isinstance(self.path_or_fileobj, bytes):
            yield io.BytesIO(self.path_or_fileobj)
        elif isinstance(self.path_or_fileobj, io.BufferedIOBase):
            prev_pos = self.path_or_fileobj.tell()
            yield self.path_or_fileobj
            self.path_or_fileobj.seek(prev_pos, io.SEEK_SET)

    def b64content(self) -> bytes:
        with self.as_file() as file:
            return base64.b64encode(file.read())

    @property
    def _local_oid(self) -> Optional[str]:
        if self._upload_mode is None:
            return None
        elif self._upload_mode == 'lfs':
            return self.upload_info.sha256.hex()
        else:
            with self.as_file() as file:
                return sha.git_hash(file.read())

def _validate_path_in_repo(path_in_repo: str) -> str:
    if path_in_repo.startswith('/'):
        path_in_repo = path_in_repo[1:]
    if path_in_repo == '.' or path_in_repo == '..' or path_in_repo.startswith('../'):
        raise ValueError(f"Invalid `path_in_repo` in CommitOperation: '{path_in_repo}'")
    if path_in_repo.startswith('./'):
        path_in_repo = path_in_repo[2:]
    for forbidden in FORBIDDEN_FOLDERS:
        if any((part == forbidden for part in path_in_repo.split('/'))):
            raise ValueError(f"Invalid `path_in_repo` in CommitOperation: cannot update files under a '{forbidden}/' folder (path: '{path_in_repo}').")
    return path_in_repo
CommitOperation = Union[CommitOperationAdd, CommitOperationCopy, CommitOperationDelete]

def _warn_on_overwriting_operations(operations: list[CommitOperation]) -> None:
    nb_additions_per_path: dict[str, int] = defaultdict(int)
    for operation in operations:
        path_in_repo = operation.path_in_repo
        if isinstance(operation, CommitOperationAdd):
            if nb_additions_per_path[path_in_repo] > 0:
                warnings.warn(f"About to update multiple times the same file in the same commit: '{path_in_repo}'. This can cause undesired inconsistencies in your repo.")
            nb_additions_per_path[path_in_repo] += 1
            for parent in PurePosixPath(path_in_repo).parents:
                nb_additions_per_path[str(parent)] += 1
        if isinstance(operation, CommitOperationDelete):
            if nb_additions_per_path[str(PurePosixPath(path_in_repo))] > 0:
                if operation.is_folder:
                    warnings.warn(f"About to delete a folder containing files that have just been updated within the same commit: '{path_in_repo}'. This can cause undesired inconsistencies in your repo.")
                else:
                    warnings.warn(f"About to delete a file that have just been updated within the same commit: '{path_in_repo}'. This can cause undesired inconsistencies in your repo.")

@validate_hf_hub_args
def _upload_files(*, additions: list[CommitOperationAdd], repo_type: str, repo_id: str, headers: dict[str, str], endpoint: Optional[str]=None, num_threads: int=5, revision: Optional[str]=None, create_pr: Optional[bool]=None):
    xet_additions: list[CommitOperationAdd] = []
    lfs_actions: list[dict[str, Any]] = []
    lfs_oid2addop: dict[str, CommitOperationAdd] = {}
    for chunk in chunk_iterable(additions, chunk_size=UPLOAD_BATCH_MAX_NUM_FILES):
        chunk_list = [op for op in chunk]
        transfers: list[str] = ['basic', 'multipart']
        has_buffered_io_data = any((isinstance(op.path_or_fileobj, io.BufferedIOBase) for op in chunk_list))
        if is_xet_available():
            if not has_buffered_io_data:
                transfers.append('xet')
            else:
                logger.warning('Uploading files as a binary IO buffer is not supported by Xet Storage. Falling back to HTTP upload.')
        actions_chunk, errors_chunk, chosen_transfer = post_lfs_batch_info(upload_infos=[op.upload_info for op in chunk_list], repo_id=repo_id, repo_type=repo_type, revision=revision, endpoint=endpoint, headers=headers, token=None, transfers=transfers)
        if errors_chunk:
            message = '\n'.join([f"Encountered error for file with OID {err.get('oid')}: `{err.get('error', {}).get('message')}" for err in errors_chunk])
            raise ValueError(f'LFS batch API returned errors:\n{message}')
        if chosen_transfer == 'xet' and 'xet' in transfers:
            xet_additions.extend(chunk_list)
        else:
            lfs_actions.extend(actions_chunk)
            for op in chunk_list:
                lfs_oid2addop[op.upload_info.sha256.hex()] = op
    if len(lfs_actions) > 0:
        _upload_lfs_files(actions=lfs_actions, oid2addop=lfs_oid2addop, headers=headers, endpoint=endpoint, num_threads=num_threads)
    if len(xet_additions) > 0:
        _upload_xet_files(additions=xet_additions, repo_type=repo_type, repo_id=repo_id, headers=headers, endpoint=endpoint, revision=revision, create_pr=create_pr)

@validate_hf_hub_args
def _upload_lfs_files(*, actions: list[dict[str, Any]], oid2addop: dict[str, CommitOperationAdd], headers: dict[str, str], endpoint: Optional[str]=None, num_threads: int=5):
    filtered_actions = []
    for action in actions:
        if action.get('actions') is None:
            logger.debug(f"Content of file {oid2addop[action['oid']].path_in_repo} is already present upstream - skipping upload.")
        else:
            filtered_actions.append(action)

    def _wrapped_lfs_upload(batch_action) -> None:
        try:
            operation = oid2addop[batch_action['oid']]
            lfs_upload(operation=operation, lfs_batch_action=batch_action, headers=headers, endpoint=endpoint)
        except Exception as exc:
            raise RuntimeError(f"Error while uploading '{operation.path_in_repo}' to the Hub.") from exc
    if len(filtered_actions) == 1:
        logger.debug('Uploading 1 LFS file to the Hub')
        _wrapped_lfs_upload(filtered_actions[0])
    else:
        logger.debug(f'Uploading {len(filtered_actions)} LFS files to the Hub using up to {num_threads} threads concurrently')
        thread_map(_wrapped_lfs_upload, filtered_actions, desc=f'Upload {len(filtered_actions)} LFS files', max_workers=num_threads, tqdm_class=hf_tqdm)

@validate_hf_hub_args
def _upload_xet_files(*, additions: list[CommitOperationAdd], repo_type: str, repo_id: str, headers: dict[str, str], endpoint: Optional[str]=None, revision: Optional[str]=None, create_pr: Optional[bool]=None):
    if len(additions) == 0:
        return
    from hf_xet import upload_bytes, upload_files
    from .utils._xet_progress_reporting import XetProgressReporter
    try:
        xet_connection_info = fetch_xet_connection_info_from_repo_info(token_type=XetTokenType.WRITE, repo_id=repo_id, repo_type=repo_type, revision=revision, headers=headers, endpoint=endpoint, params={'create_pr': '1'} if create_pr else None)
    except HfHubHTTPError as e:
        if e.response.status_code == 401:
            raise XetAuthorizationError(f'You are unauthorized to upload to xet storage for {repo_type}/{repo_id}. Please check that you have configured your access token with write access to the repo.') from e
        raise
    xet_endpoint = xet_connection_info.endpoint
    access_token_info = (xet_connection_info.access_token, xet_connection_info.expiration_unix_epoch)

    def token_refresher() -> tuple[str, int]:
        new_xet_connection = fetch_xet_connection_info_from_repo_info(token_type=XetTokenType.WRITE, repo_id=repo_id, repo_type=repo_type, revision=revision, headers=headers, endpoint=endpoint, params={'create_pr': '1'} if create_pr else None)
        if new_xet_connection is None:
            raise XetRefreshTokenError('Failed to refresh xet token')
        return (new_xet_connection.access_token, new_xet_connection.expiration_unix_epoch)
    if not are_progress_bars_disabled():
        progress = XetProgressReporter()
        progress_callback = progress.update_progress
    else:
        progress, progress_callback = (None, None)
    try:
        all_bytes_ops = [op for op in additions if isinstance(op.path_or_fileobj, bytes)]
        all_paths_ops = [op for op in additions if isinstance(op.path_or_fileobj, (str, Path))]
        xet_headers = headers.copy()
        xet_headers.pop('authorization', None)
        if len(all_paths_ops) > 0:
            all_paths = [str(op.path_or_fileobj) for op in all_paths_ops]
            all_sha256s = [op.upload_info.sha256.hex() for op in all_paths_ops]
            upload_files(all_paths, xet_endpoint, access_token_info, token_refresher, progress_callback, repo_type, request_headers=xet_headers, sha256s=all_sha256s)
        if len(all_bytes_ops) > 0:
            all_bytes = [op.path_or_fileobj for op in all_bytes_ops]
            all_sha256s = [op.upload_info.sha256.hex() for op in all_bytes_ops]
            upload_bytes(all_bytes, xet_endpoint, access_token_info, token_refresher, progress_callback, repo_type, request_headers=xet_headers, sha256s=all_sha256s)
    finally:
        if progress is not None:
            progress.close(False)
    return

def _validate_preupload_info(preupload_info: dict):
    files = preupload_info.get('files')
    if not isinstance(files, list):
        raise ValueError('preupload_info is improperly formatted')
    for file_info in files:
        if not (isinstance(file_info, dict) and isinstance(file_info.get('path'), str) and isinstance(file_info.get('uploadMode'), str) and (file_info['uploadMode'] in ('lfs', 'regular'))):
            raise ValueError('preupload_info is improperly formatted:')
    return preupload_info

@validate_hf_hub_args
def _fetch_upload_modes(additions: Iterable[CommitOperationAdd], repo_type: str, repo_id: str, headers: dict[str, str], revision: str, endpoint: Optional[str]=None, create_pr: bool=False, gitignore_content: Optional[str]=None) -> None:
    endpoint = endpoint if endpoint is not None else constants.ENDPOINT
    upload_modes: dict[str, UploadMode] = {}
    should_ignore_info: dict[str, bool] = {}
    oid_info: dict[str, Optional[str]] = {}
    for chunk in chunk_iterable(additions, 256):
        payload: dict = {'files': [{'path': op.path_in_repo, 'sample': base64.b64encode(op.upload_info.sample).decode('ascii'), 'size': op.upload_info.size} for op in chunk]}
        if gitignore_content is not None:
            payload['gitIgnore'] = gitignore_content
        resp = http_backoff('POST', f'{endpoint}/api/{repo_type}s/{repo_id}/preupload/{revision}', json=payload, headers=headers, params={'create_pr': '1'} if create_pr else None)
        hf_raise_for_status(resp)
        preupload_info = _validate_preupload_info(resp.json())
        upload_modes.update(**{file['path']: file['uploadMode'] for file in preupload_info['files']})
        should_ignore_info.update(**{file['path']: file['shouldIgnore'] for file in preupload_info['files']})
        oid_info.update(**{file['path']: file.get('oid') for file in preupload_info['files']})
    for addition in additions:
        addition._upload_mode = upload_modes[addition.path_in_repo]
        addition._should_ignore = should_ignore_info[addition.path_in_repo]
        addition._remote_oid = oid_info[addition.path_in_repo]
    for addition in additions:
        if addition.upload_info.size == 0:
            addition._upload_mode = 'regular'

@validate_hf_hub_args
def _fetch_files_to_copy(copies: Iterable[CommitOperationCopy], repo_type: str, repo_id: str, headers: dict[str, str], revision: str, endpoint: Optional[str]=None) -> dict[tuple[str, Optional[str]], Union['RepoFile', bytes]]:
    from .hf_api import HfApi, RepoFolder
    hf_api = HfApi(endpoint=endpoint, headers=headers)
    files_to_copy: dict[tuple[str, Optional[str]], Union['RepoFile', bytes]] = {}
    oid_info: dict[tuple[str, Optional[str]], Optional[str]] = {}
    dest_paths = [op.path_in_repo for op in copies]
    for offset in range(0, len(dest_paths), FETCH_LFS_BATCH_SIZE):
        dest_repo_files = hf_api.get_paths_info(repo_id=repo_id, paths=dest_paths[offset:offset + FETCH_LFS_BATCH_SIZE], revision=revision, repo_type=repo_type)
        for file in dest_repo_files:
            if not isinstance(file, RepoFolder):
                oid_info[file.path, revision] = file.blob_id
    for src_revision, operations in groupby(copies, key=lambda op: op.src_revision):
        operations = list(operations)
        src_paths = [op.src_path_in_repo for op in operations]
        for offset in range(0, len(src_paths), FETCH_LFS_BATCH_SIZE):
            src_repo_files = hf_api.get_paths_info(repo_id=repo_id, paths=src_paths[offset:offset + FETCH_LFS_BATCH_SIZE], revision=src_revision or revision, repo_type=repo_type)
            for src_repo_file in src_repo_files:
                if isinstance(src_repo_file, RepoFolder):
                    raise NotImplementedError('Copying a folder is not implemented.')
                oid_info[src_repo_file.path, src_revision] = src_repo_file.blob_id
                if src_repo_file.lfs:
                    files_to_copy[src_repo_file.path, src_revision] = src_repo_file
                else:
                    url = hf_hub_url(endpoint=endpoint, repo_type=repo_type, repo_id=repo_id, revision=src_revision or revision, filename=src_repo_file.path)
                    response = get_session().get(url, headers=headers)
                    hf_raise_for_status(response)
                    files_to_copy[src_repo_file.path, src_revision] = response.content
        for operation in operations:
            if (operation.src_path_in_repo, src_revision) not in files_to_copy:
                raise EntryNotFoundError(f'Cannot copy {operation.src_path_in_repo} at revision {src_revision or revision}: file is missing on repo.')
            operation._src_oid = oid_info.get((operation.src_path_in_repo, operation.src_revision))
            operation._dest_oid = oid_info.get((operation.path_in_repo, revision))
    return files_to_copy

def _prepare_commit_payload(operations: Iterable[CommitOperation], files_to_copy: dict[tuple[str, Optional[str]], Union['RepoFile', bytes]], commit_message: str, commit_description: Optional[str]=None, parent_commit: Optional[str]=None) -> Iterable[dict[str, Any]]:
    commit_description = commit_description if commit_description is not None else ''
    header_value = {'summary': commit_message, 'description': commit_description}
    if parent_commit is not None:
        header_value['parentCommit'] = parent_commit
    yield {'key': 'header', 'value': header_value}
    nb_ignored_files = 0
    for operation in operations:
        if isinstance(operation, CommitOperationAdd) and operation._should_ignore:
            logger.debug(f"Skipping file '{operation.path_in_repo}' in commit (ignored by gitignore file).")
            nb_ignored_files += 1
            continue
        if isinstance(operation, CommitOperationAdd) and operation._upload_mode == 'regular':
            yield {'key': 'file', 'value': {'content': operation.b64content().decode(), 'path': operation.path_in_repo, 'encoding': 'base64'}}
        elif isinstance(operation, CommitOperationAdd) and operation._upload_mode == 'lfs':
            yield {'key': 'lfsFile', 'value': {'path': operation.path_in_repo, 'algo': 'sha256', 'oid': operation.upload_info.sha256.hex(), 'size': operation.upload_info.size}}
        elif isinstance(operation, CommitOperationDelete):
            yield {'key': 'deletedFolder' if operation.is_folder else 'deletedFile', 'value': {'path': operation.path_in_repo}}
        elif isinstance(operation, CommitOperationCopy):
            file_to_copy = files_to_copy[operation.src_path_in_repo, operation.src_revision]
            if isinstance(file_to_copy, bytes):
                yield {'key': 'file', 'value': {'content': base64.b64encode(file_to_copy).decode(), 'path': operation.path_in_repo, 'encoding': 'base64'}}
            elif file_to_copy.lfs:
                yield {'key': 'lfsFile', 'value': {'path': operation.path_in_repo, 'algo': 'sha256', 'oid': file_to_copy.lfs.sha256}}
            else:
                raise ValueError('Malformed files_to_copy (should be raw file content as bytes or RepoFile objects with LFS info.')
        else:
            raise ValueError(f"Unknown operation to commit. Operation: {operation}. Upload mode: {getattr(operation, '_upload_mode', None)}")
    if nb_ignored_files > 0:
        logger.info(f'Skipped {nb_ignored_files} file(s) in commit (ignored by gitignore file).')
