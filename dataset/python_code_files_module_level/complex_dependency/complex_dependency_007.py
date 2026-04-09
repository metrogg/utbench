import os
import shutil
from collections import defaultdict
from dataclasses import dataclass
from pathlib import Path
from typing import Literal, Optional, Union
from huggingface_hub.errors import CacheNotFound, CorruptedCacheException
from ..constants import HF_HUB_CACHE
from . import logging
from ._parsing import format_timesince
from ._terminal import tabulate
logger = logging.get_logger(__name__)
REPO_TYPE_T = Literal['model', 'dataset', 'space']
FILES_TO_IGNORE = ['.DS_Store']

@dataclass(frozen=True)
class CachedFileInfo:
    file_name: str
    file_path: Path
    blob_path: Path
    size_on_disk: int
    blob_last_accessed: float
    blob_last_modified: float

    @property
    def blob_last_accessed_str(self) -> str:
        return format_timesince(self.blob_last_accessed)

    @property
    def blob_last_modified_str(self) -> str:
        return format_timesince(self.blob_last_modified)

    @property
    def size_on_disk_str(self) -> str:
        return _format_size(self.size_on_disk)

@dataclass(frozen=True)
class CachedRevisionInfo:
    commit_hash: str
    snapshot_path: Path
    size_on_disk: int
    files: frozenset[CachedFileInfo]
    refs: frozenset[str]
    last_modified: float

    @property
    def last_modified_str(self) -> str:
        return format_timesince(self.last_modified)

    @property
    def size_on_disk_str(self) -> str:
        return _format_size(self.size_on_disk)

    @property
    def nb_files(self) -> int:
        return len(self.files)

@dataclass(frozen=True)
class CachedRepoInfo:
    repo_id: str
    repo_type: REPO_TYPE_T
    repo_path: Path
    size_on_disk: int
    nb_files: int
    revisions: frozenset[CachedRevisionInfo]
    last_accessed: float
    last_modified: float

    @property
    def last_accessed_str(self) -> str:
        return format_timesince(self.last_accessed)

    @property
    def last_modified_str(self) -> str:
        return format_timesince(self.last_modified)

    @property
    def size_on_disk_str(self) -> str:
        return _format_size(self.size_on_disk)

    @property
    def cache_id(self) -> str:
        return f'{self.repo_type}/{self.repo_id}'

    @property
    def refs(self) -> dict[str, CachedRevisionInfo]:
        return {ref: revision for revision in self.revisions for ref in revision.refs}

@dataclass(frozen=True)
class DeleteCacheStrategy:
    expected_freed_size: int
    blobs: frozenset[Path]
    refs: frozenset[Path]
    repos: frozenset[Path]
    snapshots: frozenset[Path]

    @property
    def expected_freed_size_str(self) -> str:
        return _format_size(self.expected_freed_size)

    def execute(self) -> None:
        for path in self.repos:
            _try_delete_path(path, path_type='repo')
        for path in self.snapshots:
            _try_delete_path(path, path_type='snapshot')
        for path in self.refs:
            _try_delete_path(path, path_type='ref')
        for path in self.blobs:
            _try_delete_path(path, path_type='blob')
        logger.info(f'Cache deletion done. Saved {self.expected_freed_size_str}.')

@dataclass(frozen=True)
class HFCacheInfo:
    size_on_disk: int
    repos: frozenset[CachedRepoInfo]
    warnings: list[CorruptedCacheException]

    @property
    def size_on_disk_str(self) -> str:
        return _format_size(self.size_on_disk)

    def delete_revisions(self, *revisions: str) -> DeleteCacheStrategy:
        hashes_to_delete: set[str] = set(revisions)
        repos_with_revisions: dict[CachedRepoInfo, set[CachedRevisionInfo]] = defaultdict(set)
        for repo in self.repos:
            for revision in repo.revisions:
                if revision.commit_hash in hashes_to_delete:
                    repos_with_revisions[repo].add(revision)
                    hashes_to_delete.remove(revision.commit_hash)
        if len(hashes_to_delete) > 0:
            logger.warning(f"Revision(s) not found - cannot delete them: {', '.join(hashes_to_delete)}")
        delete_strategy_blobs: set[Path] = set()
        delete_strategy_refs: set[Path] = set()
        delete_strategy_repos: set[Path] = set()
        delete_strategy_snapshots: set[Path] = set()
        delete_strategy_expected_freed_size = 0
        for affected_repo, revisions_to_delete in repos_with_revisions.items():
            other_revisions = affected_repo.revisions - revisions_to_delete
            if len(other_revisions) == 0:
                delete_strategy_repos.add(affected_repo.repo_path)
                delete_strategy_expected_freed_size += affected_repo.size_on_disk
                continue
            for revision_to_delete in revisions_to_delete:
                delete_strategy_snapshots.add(revision_to_delete.snapshot_path)
                for ref in revision_to_delete.refs:
                    delete_strategy_refs.add(affected_repo.repo_path / 'refs' / ref)
                for file in revision_to_delete.files:
                    if file.blob_path not in delete_strategy_blobs:
                        is_file_alone = True
                        for revision in other_revisions:
                            for rev_file in revision.files:
                                if file.blob_path == rev_file.blob_path:
                                    is_file_alone = False
                                    break
                            if not is_file_alone:
                                break
                        if is_file_alone:
                            delete_strategy_blobs.add(file.blob_path)
                            delete_strategy_expected_freed_size += file.size_on_disk
        return DeleteCacheStrategy(blobs=frozenset(delete_strategy_blobs), refs=frozenset(delete_strategy_refs), repos=frozenset(delete_strategy_repos), snapshots=frozenset(delete_strategy_snapshots), expected_freed_size=delete_strategy_expected_freed_size)

    def export_as_table(self, *, verbosity: int=0) -> str:
        if verbosity == 0:
            return tabulate(rows=[[repo.repo_id, repo.repo_type, '{:>12}'.format(repo.size_on_disk_str), repo.nb_files, repo.last_accessed_str, repo.last_modified_str, ', '.join(sorted(repo.refs)), str(repo.repo_path)] for repo in sorted(self.repos, key=lambda repo: repo.repo_path)], headers=['REPO ID', 'REPO TYPE', 'SIZE ON DISK', 'NB FILES', 'LAST_ACCESSED', 'LAST_MODIFIED', 'REFS', 'LOCAL PATH'])
        else:
            return tabulate(rows=[[repo.repo_id, repo.repo_type, revision.commit_hash, '{:>12}'.format(revision.size_on_disk_str), revision.nb_files, revision.last_modified_str, ', '.join(sorted(revision.refs)), str(revision.snapshot_path)] for repo in sorted(self.repos, key=lambda repo: repo.repo_path) for revision in sorted(repo.revisions, key=lambda revision: revision.commit_hash)], headers=['REPO ID', 'REPO TYPE', 'REVISION', 'SIZE ON DISK', 'NB FILES', 'LAST_MODIFIED', 'REFS', 'LOCAL PATH'])

def scan_cache_dir(cache_dir: Optional[Union[str, Path]]=None) -> HFCacheInfo:
    if cache_dir is None:
        cache_dir = HF_HUB_CACHE
    cache_dir = Path(cache_dir).expanduser().resolve()
    if not cache_dir.exists():
        raise CacheNotFound(f'Cache directory not found: {cache_dir}. Please use `cache_dir` argument or set `HF_HUB_CACHE` environment variable.', cache_dir=cache_dir)
    if cache_dir.is_file():
        raise ValueError(f'Scan cache expects a directory but found a file: {cache_dir}. Please use `cache_dir` argument or set `HF_HUB_CACHE` environment variable.')
    repos: set[CachedRepoInfo] = set()
    warnings: list[CorruptedCacheException] = []
    for repo_path in cache_dir.iterdir():
        if repo_path.name == '.locks':
            continue
        try:
            repos.add(_scan_cached_repo(repo_path))
        except CorruptedCacheException as e:
            warnings.append(e)
    return HFCacheInfo(repos=frozenset(repos), size_on_disk=sum((repo.size_on_disk for repo in repos)), warnings=warnings)

def _scan_cached_repo(repo_path: Path) -> CachedRepoInfo:
    if not repo_path.is_dir():
        raise CorruptedCacheException(f'Repo path is not a directory: {repo_path}')
    if '--' not in repo_path.name:
        raise CorruptedCacheException(f'Repo path is not a valid HuggingFace cache directory: {repo_path}')
    repo_type, repo_id = repo_path.name.split('--', maxsplit=1)
    repo_type = repo_type[:-1]
    repo_id = repo_id.replace('--', '/')
    if repo_type not in {'dataset', 'model', 'space'}:
        raise CorruptedCacheException(f'Repo type must be `dataset`, `model` or `space`, found `{repo_type}` ({repo_path}).')
    blob_stats: dict[Path, os.stat_result] = {}
    snapshots_path = repo_path / 'snapshots'
    refs_path = repo_path / 'refs'
    if not snapshots_path.exists() or not snapshots_path.is_dir():
        raise CorruptedCacheException(f"Snapshots dir doesn't exist in cached repo: {snapshots_path}")
    refs_by_hash: dict[str, set[str]] = defaultdict(set)
    if refs_path.exists():
        if refs_path.is_file():
            raise CorruptedCacheException(f'Refs directory cannot be a file: {refs_path}')
        for ref_path in refs_path.glob('**/*'):
            if ref_path.is_dir() or ref_path.name in FILES_TO_IGNORE:
                continue
            ref_name = str(ref_path.relative_to(refs_path))
            with ref_path.open() as f:
                commit_hash = f.read()
            refs_by_hash[commit_hash].add(ref_name)
    cached_revisions: set[CachedRevisionInfo] = set()
    for revision_path in snapshots_path.iterdir():
        if revision_path.name in FILES_TO_IGNORE:
            continue
        if revision_path.is_file():
            raise CorruptedCacheException(f'Snapshots folder corrupted. Found a file: {revision_path}')
        cached_files = set()
        for file_path in revision_path.glob('**/*'):
            if file_path.is_dir():
                continue
            blob_path = Path(file_path).resolve()
            if not blob_path.exists():
                raise CorruptedCacheException(f'Blob missing (broken symlink): {blob_path}')
            if blob_path not in blob_stats:
                blob_stats[blob_path] = blob_path.stat()
            cached_files.add(CachedFileInfo(file_name=file_path.name, file_path=file_path, size_on_disk=blob_stats[blob_path].st_size, blob_path=blob_path, blob_last_accessed=blob_stats[blob_path].st_atime, blob_last_modified=blob_stats[blob_path].st_mtime))
        if len(cached_files) > 0:
            revision_last_modified = max((blob_stats[file.blob_path].st_mtime for file in cached_files))
        else:
            revision_last_modified = revision_path.stat().st_mtime
        cached_revisions.add(CachedRevisionInfo(commit_hash=revision_path.name, files=frozenset(cached_files), refs=frozenset(refs_by_hash.pop(revision_path.name, set())), size_on_disk=sum((blob_stats[blob_path].st_size for blob_path in set((file.blob_path for file in cached_files)))), snapshot_path=revision_path, last_modified=revision_last_modified))
    if len(refs_by_hash) > 0:
        raise CorruptedCacheException(f'Reference(s) refer to missing commit hashes: {dict(refs_by_hash)} ({repo_path}).')
    if len(blob_stats) > 0:
        repo_last_accessed = max((stat.st_atime for stat in blob_stats.values()))
        repo_last_modified = max((stat.st_mtime for stat in blob_stats.values()))
    else:
        repo_stats = repo_path.stat()
        repo_last_accessed = repo_stats.st_atime
        repo_last_modified = repo_stats.st_mtime
    return CachedRepoInfo(nb_files=len(blob_stats), repo_id=repo_id, repo_path=repo_path, repo_type=repo_type, revisions=frozenset(cached_revisions), size_on_disk=sum((stat.st_size for stat in blob_stats.values())), last_accessed=repo_last_accessed, last_modified=repo_last_modified)

def _format_size(num: int) -> str:
    num_f = float(num)
    for unit in ['', 'K', 'M', 'G', 'T', 'P', 'E', 'Z']:
        if abs(num_f) < 1000.0:
            return f'{num_f:3.1f}{unit}'
        num_f /= 1000.0
    return f'{num_f:.1f}Y'

def _try_delete_path(path: Path, path_type: str) -> None:
    logger.info(f'Delete {path_type}: {path}')
    try:
        if path.is_file():
            os.remove(path)
        else:
            shutil.rmtree(path)
    except FileNotFoundError:
        logger.warning(f"Couldn't delete {path_type}: file not found ({path})", exc_info=True)
    except PermissionError:
        logger.warning(f"Couldn't delete {path_type}: permission denied ({path})", exc_info=True)
