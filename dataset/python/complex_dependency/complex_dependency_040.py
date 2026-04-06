import os
from pathlib import Path
from typing import Iterable, List, Literal, Optional, Union, overload
import httpx
from tqdm.auto import tqdm as base_tqdm
from tqdm.contrib.concurrent import thread_map
from . import constants
from .errors import DryRunError, GatedRepoError, HfHubHTTPError, LocalEntryNotFoundError, RepositoryNotFoundError, RevisionNotFoundError
from .file_download import REGEX_COMMIT_HASH, DryRunFileInfo, hf_hub_download, repo_folder_name
from .hf_api import DatasetInfo, HfApi, ModelInfo, RepoFile, SpaceInfo
from .utils import OfflineModeIsEnabled, filter_repo_objects, is_tqdm_disabled, logging, validate_hf_hub_args
from .utils import tqdm as hf_tqdm
logger = logging.get_logger(__name__)
LARGE_REPO_THRESHOLD = 1000

@overload
def snapshot_download(repo_id: str, *, repo_type: Optional[str]=None, revision: Optional[str]=None, cache_dir: Union[str, Path, None]=None, local_dir: Union[str, Path, None]=None, library_name: Optional[str]=None, library_version: Optional[str]=None, user_agent: Optional[Union[dict, str]]=None, etag_timeout: float=constants.DEFAULT_ETAG_TIMEOUT, force_download: bool=False, token: Optional[Union[bool, str]]=None, local_files_only: bool=False, allow_patterns: Optional[Union[list[str], str]]=None, ignore_patterns: Optional[Union[list[str], str]]=None, max_workers: int=8, tqdm_class: Optional[type[base_tqdm]]=None, headers: Optional[dict[str, str]]=None, endpoint: Optional[str]=None, dry_run: Literal[False]=False) -> str:
    ...

@overload
def snapshot_download(repo_id: str, *, repo_type: Optional[str]=None, revision: Optional[str]=None, cache_dir: Union[str, Path, None]=None, local_dir: Union[str, Path, None]=None, library_name: Optional[str]=None, library_version: Optional[str]=None, user_agent: Optional[Union[dict, str]]=None, etag_timeout: float=constants.DEFAULT_ETAG_TIMEOUT, force_download: bool=False, token: Optional[Union[bool, str]]=None, local_files_only: bool=False, allow_patterns: Optional[Union[list[str], str]]=None, ignore_patterns: Optional[Union[list[str], str]]=None, max_workers: int=8, tqdm_class: Optional[type[base_tqdm]]=None, headers: Optional[dict[str, str]]=None, endpoint: Optional[str]=None, dry_run: Literal[True]=True) -> list[DryRunFileInfo]:
    ...

@overload
def snapshot_download(repo_id: str, *, repo_type: Optional[str]=None, revision: Optional[str]=None, cache_dir: Union[str, Path, None]=None, local_dir: Union[str, Path, None]=None, library_name: Optional[str]=None, library_version: Optional[str]=None, user_agent: Optional[Union[dict, str]]=None, etag_timeout: float=constants.DEFAULT_ETAG_TIMEOUT, force_download: bool=False, token: Optional[Union[bool, str]]=None, local_files_only: bool=False, allow_patterns: Optional[Union[list[str], str]]=None, ignore_patterns: Optional[Union[list[str], str]]=None, max_workers: int=8, tqdm_class: Optional[type[base_tqdm]]=None, headers: Optional[dict[str, str]]=None, endpoint: Optional[str]=None, dry_run: bool=False) -> Union[str, list[DryRunFileInfo]]:
    ...

@validate_hf_hub_args
def snapshot_download(repo_id: str, *, repo_type: Optional[str]=None, revision: Optional[str]=None, cache_dir: Union[str, Path, None]=None, local_dir: Union[str, Path, None]=None, library_name: Optional[str]=None, library_version: Optional[str]=None, user_agent: Optional[Union[dict, str]]=None, etag_timeout: float=constants.DEFAULT_ETAG_TIMEOUT, force_download: bool=False, token: Optional[Union[bool, str]]=None, local_files_only: bool=False, allow_patterns: Optional[Union[list[str], str]]=None, ignore_patterns: Optional[Union[list[str], str]]=None, max_workers: int=8, tqdm_class: Optional[type[base_tqdm]]=None, headers: Optional[dict[str, str]]=None, endpoint: Optional[str]=None, dry_run: bool=False) -> Union[str, list[DryRunFileInfo]]:
    if cache_dir is None:
        cache_dir = constants.HF_HUB_CACHE
    if revision is None:
        revision = constants.DEFAULT_REVISION
    if isinstance(cache_dir, Path):
        cache_dir = str(cache_dir)
    if repo_type is None:
        repo_type = 'model'
    if repo_type not in constants.REPO_TYPES:
        raise ValueError(f'Invalid repo type: {repo_type}. Accepted repo types are: {str(constants.REPO_TYPES)}')
    storage_folder = os.path.join(cache_dir, repo_folder_name(repo_id=repo_id, repo_type=repo_type))
    api = HfApi(library_name=library_name, library_version=library_version, user_agent=user_agent, endpoint=endpoint, headers=headers, token=token)
    repo_info: Union[ModelInfo, DatasetInfo, SpaceInfo, None] = None
    api_call_error: Optional[Exception] = None
    if not local_files_only:
        try:
            repo_info = api.repo_info(repo_id=repo_id, repo_type=repo_type, revision=revision)
        except httpx.ProxyError:
            raise
        except (httpx.ConnectError, httpx.TimeoutException, OfflineModeIsEnabled) as error:
            api_call_error = error
            pass
        except RevisionNotFoundError:
            raise
        except HfHubHTTPError as error:
            api_call_error = error
            pass
    if repo_info is None:
        if dry_run:
            raise DryRunError('Dry run cannot be performed as the repository cannot be accessed. Please check your internet connection or authentication token.') from api_call_error
        commit_hash = None
        if REGEX_COMMIT_HASH.match(revision):
            commit_hash = revision
        else:
            ref_path = os.path.join(storage_folder, 'refs', revision)
            if os.path.exists(ref_path):
                with open(ref_path) as f:
                    commit_hash = f.read()
        if commit_hash is not None and local_dir is None:
            snapshot_folder = os.path.join(storage_folder, 'snapshots', commit_hash)
            if os.path.exists(snapshot_folder):
                return snapshot_folder
        if local_dir is not None:
            local_dir = Path(local_dir)
            if local_dir.is_dir() and any(local_dir.iterdir()):
                logger.warning(f'Returning existing local_dir `{local_dir}` as remote repo cannot be accessed in `snapshot_download` ({api_call_error}).')
                return str(local_dir.resolve())
        if local_files_only:
            raise LocalEntryNotFoundError("Cannot find an appropriate cached snapshot folder for the specified revision on the local disk and outgoing traffic has been disabled. To enable repo look-ups and downloads online, pass 'local_files_only=False' as input.")
        elif isinstance(api_call_error, OfflineModeIsEnabled):
            raise LocalEntryNotFoundError("Cannot find an appropriate cached snapshot folder for the specified revision on the local disk and outgoing traffic has been disabled. To enable repo look-ups and downloads online, set 'HF_HUB_OFFLINE=0' as environment variable.") from api_call_error
        elif isinstance(api_call_error, (RepositoryNotFoundError, GatedRepoError)) or (isinstance(api_call_error, HfHubHTTPError) and api_call_error.response.status_code == 401):
            raise api_call_error
        else:
            raise LocalEntryNotFoundError(f'Got: {api_call_error.__class__.__name__}: {api_call_error}\nAn error happened while trying to locate the files on the Hub and we cannot find the appropriate snapshot folder for the specified revision on the local disk. Please check your internet connection and try again.') from api_call_error
    assert repo_info.sha is not None, 'Repo info returned from server must have a revision sha.'
    repo_files: Iterable[str] = [f.rfilename for f in repo_info.siblings] if repo_info.siblings is not None else []
    unreliable_nb_files = repo_info.siblings is None or len(repo_info.siblings) == 0 or len(repo_info.siblings) > LARGE_REPO_THRESHOLD
    if unreliable_nb_files:
        logger.info('Number of files in the repo is unreliable. Using `list_repo_tree` to ensure all files are listed.')
        repo_files = (f.rfilename for f in api.list_repo_tree(repo_id=repo_id, recursive=True, revision=revision, repo_type=repo_type) if isinstance(f, RepoFile))
    filtered_repo_files: Iterable[str] = filter_repo_objects(items=repo_files, allow_patterns=allow_patterns, ignore_patterns=ignore_patterns)
    if not unreliable_nb_files:
        filtered_repo_files = list(filtered_repo_files)
        tqdm_desc = f'Fetching {len(filtered_repo_files)} files'
    else:
        tqdm_desc = 'Fetching ... files'
    if dry_run:
        tqdm_desc = '[dry-run] ' + tqdm_desc
    commit_hash = repo_info.sha
    snapshot_folder = os.path.join(storage_folder, 'snapshots', commit_hash)
    if revision != commit_hash:
        ref_path = os.path.join(storage_folder, 'refs', revision)
        try:
            os.makedirs(os.path.dirname(ref_path), exist_ok=True)
            with open(ref_path, 'w') as f:
                f.write(commit_hash)
        except OSError as e:
            logger.warning(f'Ignored error while writing commit hash to {ref_path}: {e}.')
    results: List[Union[str, DryRunFileInfo]] = []
    tqdm_class = tqdm_class or hf_tqdm
    bytes_progress = tqdm_class(desc='Downloading (incomplete total...)', disable=is_tqdm_disabled(log_level=logger.getEffectiveLevel()), total=0, initial=0, unit='B', unit_scale=True, name='huggingface_hub.snapshot_download')

    class _AggregatedTqdm:

        def __init__(self, *args, **kwargs):
            total = kwargs.pop('total', None)
            if total is not None:
                bytes_progress.total += total
                bytes_progress.refresh()
            initial = kwargs.pop('initial', 0)
            if initial:
                bytes_progress.update(initial)

        def __enter__(self):
            return self

        def __exit__(self, exc_type, exc_value, traceback):
            pass

        def update(self, n: Optional[Union[int, float]]=1) -> None:
            bytes_progress.update(n)

    def _inner_hf_hub_download(repo_file: str) -> None:
        results.append(hf_hub_download(repo_id, filename=repo_file, repo_type=repo_type, revision=commit_hash, endpoint=endpoint, cache_dir=cache_dir, local_dir=local_dir, library_name=library_name, library_version=library_version, user_agent=user_agent, etag_timeout=etag_timeout, force_download=force_download, token=token, headers=headers, tqdm_class=_AggregatedTqdm, dry_run=dry_run))
    thread_map(_inner_hf_hub_download, filtered_repo_files, desc=tqdm_desc, max_workers=max_workers, tqdm_class=tqdm_class)
    bytes_progress.set_description('Download complete')
    if dry_run:
        assert all((isinstance(r, DryRunFileInfo) for r in results))
        return results
    if local_dir is not None:
        return str(os.path.realpath(local_dir))
    return snapshot_folder
