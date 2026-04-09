import os
import re
from functools import partial
from glob import has_magic
from pathlib import Path, PurePath
from typing import Callable, Optional, Union
import huggingface_hub
from fsspec.core import url_to_fs
from huggingface_hub import HfFileSystem
from packaging import version
from tqdm.contrib.concurrent import thread_map
from . import config
from .download import DownloadConfig
from .naming import _split_re
from .splits import Split
from .utils import logging
from .utils import tqdm as hf_tqdm
from .utils.file_utils import _prepare_path_and_storage_options, is_local_path, is_relative_path, xbasename, xjoin
from .utils.py_utils import string_to_dict
SingleOriginMetadata = Union[tuple[str, str], tuple[str], tuple[()]]
SANITIZED_DEFAULT_SPLIT = str(Split.TRAIN)
logger = logging.get_logger(__name__)

class Url(str):
    pass

class EmptyDatasetError(FileNotFoundError):
    pass
SPLIT_PATTERN_SHARDED = 'data/{split}-[0-9][0-9][0-9][0-9][0-9]-of-[0-9][0-9][0-9][0-9][0-9]*.*'
SPLIT_KEYWORDS = {Split.TRAIN: ['train', 'training'], Split.VALIDATION: ['validation', 'valid', 'dev', 'val'], Split.TEST: ['test', 'testing', 'eval', 'evaluation']}
NON_WORDS_CHARS = '-._ 0-9'
if config.FSSPEC_VERSION < version.parse('2023.9.0'):
    KEYWORDS_IN_FILENAME_BASE_PATTERNS = ['**[{sep}/]{keyword}[{sep}]*', '{keyword}[{sep}]*']
    KEYWORDS_IN_DIR_NAME_BASE_PATTERNS = ['{keyword}/**', '{keyword}[{sep}]*/**', '**[{sep}/]{keyword}/**', '**[{sep}/]{keyword}[{sep}]*/**']
elif config.FSSPEC_VERSION < version.parse('2023.12.0'):
    KEYWORDS_IN_FILENAME_BASE_PATTERNS = ['**/*[{sep}/]{keyword}[{sep}]*', '{keyword}[{sep}]*']
    KEYWORDS_IN_DIR_NAME_BASE_PATTERNS = ['{keyword}/**/*', '{keyword}[{sep}]*/**/*', '**/*[{sep}/]{keyword}/**/*', '**/*[{sep}/]{keyword}[{sep}]*/**/*']
else:
    KEYWORDS_IN_FILENAME_BASE_PATTERNS = ['**/{keyword}[{sep}]*', '**/*[{sep}]{keyword}[{sep}]*']
    KEYWORDS_IN_DIR_NAME_BASE_PATTERNS = ['**/{keyword}/**', '**/{keyword}[{sep}]*/**', '**/*[{sep}]{keyword}/**', '**/*[{sep}]{keyword}[{sep}]*/**']
DEFAULT_SPLITS = [Split.TRAIN, Split.VALIDATION, Split.TEST]
DEFAULT_PATTERNS_SPLIT_IN_FILENAME = {split: [pattern.format(keyword=keyword, sep=NON_WORDS_CHARS) for keyword in SPLIT_KEYWORDS[split] for pattern in KEYWORDS_IN_FILENAME_BASE_PATTERNS] for split in DEFAULT_SPLITS}
DEFAULT_PATTERNS_SPLIT_IN_DIR_NAME = {split: [pattern.format(keyword=keyword, sep=NON_WORDS_CHARS) for keyword in SPLIT_KEYWORDS[split] for pattern in KEYWORDS_IN_DIR_NAME_BASE_PATTERNS] for split in DEFAULT_SPLITS}
DEFAULT_PATTERNS_ALL = {Split.TRAIN: ['**']}
DEFAULT_PATTERNS_LOGS = {'logs': ['**/*.eval']}
ALL_SPLIT_PATTERNS = [SPLIT_PATTERN_SHARDED]
ALL_DEFAULT_PATTERNS = [DEFAULT_PATTERNS_LOGS, DEFAULT_PATTERNS_SPLIT_IN_DIR_NAME, DEFAULT_PATTERNS_SPLIT_IN_FILENAME, DEFAULT_PATTERNS_ALL]
WILDCARD_CHARACTERS = '*[]'
FILES_TO_IGNORE = ['README.md', 'config.json', 'dataset_info.json', 'dataset_infos.json', 'dummy_data.zip', 'dataset_dict.json']

def contains_wildcards(pattern: str) -> bool:
    return any((wildcard_character in pattern for wildcard_character in WILDCARD_CHARACTERS))

def sanitize_patterns(patterns: Union[dict, list, str]) -> dict[str, Union[list[str], 'DataFilesList']]:
    if isinstance(patterns, dict):
        return {str(key): value if isinstance(value, list) else [value] for key, value in patterns.items()}
    elif isinstance(patterns, str):
        return {SANITIZED_DEFAULT_SPLIT: [patterns]}
    elif isinstance(patterns, list):
        if any((isinstance(pattern, dict) for pattern in patterns)):
            for pattern in patterns:
                if not (isinstance(pattern, dict) and len(pattern) == 2 and ('split' in pattern) and isinstance(pattern.get('path'), (str, list))):
                    raise ValueError(f"Invalid format for data_files entry. Each item must be a dictionary with the structure {{'split': <split_name>, 'path': <path_or_list_of_paths>}}.\nReceived: {pattern}")
            splits = [pattern['split'] for pattern in patterns]
            if len(set(splits)) != len(splits):
                raise ValueError(f'Some splits are duplicated in data_files: {splits}')
            return {str(pattern['split']): pattern['path'] if isinstance(pattern['path'], list) else [pattern['path']] for pattern in patterns}
        else:
            return {SANITIZED_DEFAULT_SPLIT: patterns}
    else:
        return sanitize_patterns(list(patterns))

def _is_inside_unrequested_special_dir(matched_rel_path: str, pattern: str) -> bool:
    data_dirs_to_ignore_in_path = [part for part in PurePath(matched_rel_path).parent.parts if part.startswith('__')]
    data_dirs_to_ignore_in_pattern = [part for part in PurePath(pattern).parent.parts if part.startswith('__')]
    return len(data_dirs_to_ignore_in_path) != len(data_dirs_to_ignore_in_pattern)

def _is_unrequested_hidden_file_or_is_inside_unrequested_hidden_dir(matched_rel_path: str, pattern: str) -> bool:
    hidden_directories_in_path = [part for part in PurePath(matched_rel_path).parts if part.startswith('.') and (not set(part) == {'.'})]
    hidden_directories_in_pattern = [part for part in PurePath(pattern).parts if part.startswith('.') and (not set(part) == {'.'})]
    return len(hidden_directories_in_path) != len(hidden_directories_in_pattern)

def _get_data_files_patterns(pattern_resolver: Callable[[str], list[str]]) -> dict[str, list[str]]:
    for split_pattern in ALL_SPLIT_PATTERNS:
        pattern = split_pattern.replace('{split}', '*')
        try:
            data_files = pattern_resolver(pattern)
        except FileNotFoundError:
            continue
        if len(data_files) > 0:
            splits: set[str] = set()
            for p in data_files:
                p_parts = string_to_dict(xbasename(p), xbasename(split_pattern))
                assert p_parts is not None
                splits.add(p_parts['split'])
            if any((not re.match(_split_re, split) for split in splits)):
                raise ValueError(f"Split name should match '{_split_re}'' but got '{splits}'.")
            sorted_splits = [str(split) for split in DEFAULT_SPLITS if split in splits] + sorted(splits - {str(split) for split in DEFAULT_SPLITS})
            return {split: [split_pattern.format(split=split)] for split in sorted_splits}
    for patterns_dict in ALL_DEFAULT_PATTERNS:
        non_empty_splits = []
        for split, patterns in patterns_dict.items():
            for pattern in patterns:
                try:
                    data_files = pattern_resolver(pattern)
                except FileNotFoundError:
                    continue
                if len(data_files) > 0:
                    non_empty_splits.append(split)
                    break
        if non_empty_splits:
            return {split: patterns_dict[split] for split in non_empty_splits}
    raise FileNotFoundError(f"Couldn't resolve pattern {pattern} with resolver {pattern_resolver}")

def resolve_pattern(pattern: str, base_path: str, allowed_extensions: Optional[list[str]]=None, download_config: Optional[DownloadConfig]=None) -> list[str]:
    if is_relative_path(pattern):
        pattern = xjoin(base_path, pattern)
    elif is_local_path(pattern):
        base_path = os.path.splitdrive(pattern)[0] + os.sep
    else:
        base_path = ''
    pattern, storage_options = _prepare_path_and_storage_options(pattern, download_config=download_config)
    fs, fs_pattern = url_to_fs(pattern, **storage_options)
    files_to_ignore = set(FILES_TO_IGNORE) - {xbasename(pattern)}
    protocol = pattern.split('://')[0] if '://' in pattern else fs.protocol if isinstance(fs.protocol, str) else fs.protocol[0]
    protocol_prefix = protocol + '://' if protocol != 'file' else ''
    glob_kwargs = {}
    if protocol == 'hf':
        glob_kwargs['expand_info'] = False
    _, *rest_hops = pattern.split('::')
    matched_paths = []
    for filepath, info in fs.glob(fs_pattern, detail=True, **glob_kwargs).items():
        if not (info['type'] == 'file' or (info.get('islink') and os.path.isfile(os.path.realpath(filepath)))) or xbasename(filepath) in files_to_ignore:
            continue
        if _is_inside_unrequested_special_dir(filepath, fs_pattern):
            continue
        if _is_unrequested_hidden_file_or_is_inside_unrequested_hidden_dir(filepath, fs_pattern):
            continue
        filepath = filepath if '://' in filepath else protocol_prefix + filepath
        if rest_hops:
            filepath = '::'.join([filepath] + rest_hops)
        matched_paths.append(filepath)
    if allowed_extensions is not None:
        out = [filepath for filepath in matched_paths if any(('.' + suffix in allowed_extensions for suffix in xbasename(filepath).split('.')[1:]))]
        if len(out) < len(matched_paths):
            invalid_matched_files = list(set(matched_paths) - set(out))
            logger.info(f"Some files matched the pattern '{pattern}' but don't have valid data file extensions: {invalid_matched_files}")
    else:
        out = matched_paths
    if not out:
        error_msg = f"Unable to find '{pattern}'"
        if allowed_extensions is not None:
            error_msg += f' with any supported extension {list(allowed_extensions)}'
        raise FileNotFoundError(error_msg)
    return out

def get_data_patterns(base_path: str, download_config: Optional[DownloadConfig]=None) -> dict[str, list[str]]:
    resolver = partial(resolve_pattern, base_path=base_path, download_config=download_config)
    try:
        return _get_data_files_patterns(resolver)
    except FileNotFoundError:
        raise EmptyDatasetError(f"The directory at {base_path} doesn't contain any data files") from None

def _get_single_origin_metadata(data_file: str, download_config: Optional[DownloadConfig]=None) -> SingleOriginMetadata:
    if data_file.startswith(config.HF_ENDPOINT):
        fs = HfFileSystem(endpoint=config.HF_ENDPOINT, token=download_config.token)
        data_file = 'hf://' + data_file[len(config.HF_ENDPOINT) + 1:]
        data_file = data_file.replace('/resolve/', '/' if data_file.startswith('hf://buckets/') else '@', 1)
        fs_path = data_file
    else:
        data_file, storage_options = _prepare_path_and_storage_options(data_file, download_config=download_config)
        fs, fs_path = url_to_fs(data_file, **storage_options)
    if isinstance(fs, HfFileSystem):
        resolved_path = fs.resolve_path(fs_path)
        if hasattr(resolved_path, 'revision'):
            return (resolved_path.repo_id, resolved_path.revision)
    info = fs.info(fs_path)
    for key in ['ETag', 'etag', 'mtime']:
        if key in info:
            return (str(info[key]),)
    return ()

def _get_origin_metadata(data_files: list[str], download_config: Optional[DownloadConfig]=None, max_workers: Optional[int]=None) -> list[SingleOriginMetadata]:
    max_workers = max_workers if max_workers is not None else config.HF_DATASETS_MULTITHREADING_MAX_WORKERS
    if all(('hf://' in data_file for data_file in data_files)):
        return [_get_single_origin_metadata(data_file, download_config=download_config) for data_file in hf_tqdm(data_files, desc='Resolving data files', disable=len(data_files) <= 16 or None)]
    return thread_map(partial(_get_single_origin_metadata, download_config=download_config), data_files, max_workers=max_workers, tqdm_class=hf_tqdm, desc='Resolving data files', disable=len(data_files) <= 16 or None)

class DataFilesList(list[str]):

    def __init__(self, data_files: list[str], origin_metadata: list[SingleOriginMetadata]) -> None:
        super().__init__(data_files)
        self.origin_metadata = origin_metadata

    def __add__(self, other: 'DataFilesList') -> 'DataFilesList':
        return DataFilesList([*self, *other], self.origin_metadata + other.origin_metadata)

    @classmethod
    def from_hf_repo(cls, patterns: list[str], dataset_info: huggingface_hub.hf_api.DatasetInfo, base_path: Optional[str]=None, allowed_extensions: Optional[list[str]]=None, download_config: Optional[DownloadConfig]=None) -> 'DataFilesList':
        base_path = f"hf://datasets/{dataset_info.id}@{dataset_info.sha}/{base_path or ''}".rstrip('/')
        return cls.from_patterns(patterns, base_path=base_path, allowed_extensions=allowed_extensions, download_config=download_config)

    @classmethod
    def from_local_or_remote(cls, patterns: list[str], base_path: Optional[str]=None, allowed_extensions: Optional[list[str]]=None, download_config: Optional[DownloadConfig]=None) -> 'DataFilesList':
        base_path = base_path if base_path is not None else Path().resolve().as_posix()
        return cls.from_patterns(patterns, base_path=base_path, allowed_extensions=allowed_extensions, download_config=download_config)

    @classmethod
    def from_patterns(cls, patterns: list[str], base_path: Optional[str]=None, allowed_extensions: Optional[list[str]]=None, download_config: Optional[DownloadConfig]=None) -> 'DataFilesList':
        base_path = base_path if base_path is not None else Path().resolve().as_posix()
        data_files = []
        for pattern in patterns:
            try:
                data_files.extend(resolve_pattern(pattern, base_path=base_path, allowed_extensions=allowed_extensions, download_config=download_config))
            except FileNotFoundError:
                if not has_magic(pattern):
                    raise
        origin_metadata = _get_origin_metadata(data_files, download_config=download_config)
        return cls(data_files, origin_metadata)

    def filter(self, *, extensions: Optional[list[str]]=None, file_names: Optional[list[str]]=None) -> 'DataFilesList':
        patterns = []
        if extensions:
            ext_pattern = '|'.join((re.escape(ext) for ext in extensions))
            patterns.append(re.compile(f'.*({ext_pattern})(\\..+)?$'))
        if file_names:
            fn_pattern = '|'.join((re.escape(fn) for fn in file_names))
            patterns.append(re.compile(f'.*[\\/]?({fn_pattern})$'))
        if patterns:
            return DataFilesList([data_file for data_file in self if any((pattern.match(data_file) for pattern in patterns))], origin_metadata=self.origin_metadata)
        else:
            return DataFilesList(list(self), origin_metadata=self.origin_metadata)

class DataFilesDict(dict[str, DataFilesList]):

    @classmethod
    def from_local_or_remote(cls, patterns: dict[str, Union[list[str], DataFilesList]], base_path: Optional[str]=None, allowed_extensions: Optional[list[str]]=None, download_config: Optional[DownloadConfig]=None) -> 'DataFilesDict':
        out = cls()
        for key, patterns_for_key in patterns.items():
            out[key] = patterns_for_key if isinstance(patterns_for_key, DataFilesList) else DataFilesList.from_local_or_remote(patterns_for_key, base_path=base_path, allowed_extensions=allowed_extensions, download_config=download_config)
        return out

    @classmethod
    def from_hf_repo(cls, patterns: dict[str, Union[list[str], DataFilesList]], dataset_info: huggingface_hub.hf_api.DatasetInfo, base_path: Optional[str]=None, allowed_extensions: Optional[list[str]]=None, download_config: Optional[DownloadConfig]=None) -> 'DataFilesDict':
        out = cls()
        for key, patterns_for_key in patterns.items():
            out[key] = patterns_for_key if isinstance(patterns_for_key, DataFilesList) else DataFilesList.from_hf_repo(patterns_for_key, dataset_info=dataset_info, base_path=base_path, allowed_extensions=allowed_extensions, download_config=download_config)
        return out

    @classmethod
    def from_patterns(cls, patterns: dict[str, Union[list[str], DataFilesList]], base_path: Optional[str]=None, allowed_extensions: Optional[list[str]]=None, download_config: Optional[DownloadConfig]=None) -> 'DataFilesDict':
        out = cls()
        for key, patterns_for_key in patterns.items():
            out[key] = patterns_for_key if isinstance(patterns_for_key, DataFilesList) else DataFilesList.from_patterns(patterns_for_key, base_path=base_path, allowed_extensions=allowed_extensions, download_config=download_config)
        return out

    def filter(self, *, extensions: Optional[list[str]]=None, file_names: Optional[list[str]]=None) -> 'DataFilesDict':
        out = type(self)()
        for key, data_files_list in self.items():
            out[key] = data_files_list.filter(extensions=extensions, file_names=file_names)
        return out

class DataFilesPatternsList(list[str]):

    def __init__(self, patterns: list[str], allowed_extensions: list[Optional[list[str]]]):
        super().__init__(patterns)
        self.allowed_extensions = allowed_extensions

    def __add__(self, other):
        return DataFilesList([*self, *other], self.allowed_extensions + other.allowed_extensions)

    @classmethod
    def from_patterns(cls, patterns: list[str], allowed_extensions: Optional[list[str]]=None) -> 'DataFilesPatternsList':
        return cls(patterns, [allowed_extensions] * len(patterns))

    def resolve(self, base_path: str, download_config: Optional[DownloadConfig]=None) -> 'DataFilesList':
        base_path = base_path if base_path is not None else Path().resolve().as_posix()
        data_files = []
        for pattern, allowed_extensions in zip(self, self.allowed_extensions):
            try:
                data_files.extend(resolve_pattern(pattern, base_path=base_path, allowed_extensions=allowed_extensions, download_config=download_config))
            except FileNotFoundError:
                if not has_magic(pattern):
                    raise
        origin_metadata = _get_origin_metadata(data_files, download_config=download_config)
        return DataFilesList(data_files, origin_metadata)

    def filter_extensions(self, extensions: list[str]) -> 'DataFilesPatternsList':
        return DataFilesPatternsList(self, [allowed_extensions + extensions for allowed_extensions in self.allowed_extensions])

class DataFilesPatternsDict(dict[str, DataFilesPatternsList]):

    @classmethod
    def from_patterns(cls, patterns: dict[str, list[str]], allowed_extensions: Optional[list[str]]=None) -> 'DataFilesPatternsDict':
        out = cls()
        for key, patterns_for_key in patterns.items():
            out[key] = patterns_for_key if isinstance(patterns_for_key, DataFilesPatternsList) else DataFilesPatternsList.from_patterns(patterns_for_key, allowed_extensions=allowed_extensions)
        return out

    def resolve(self, base_path: str, download_config: Optional[DownloadConfig]=None) -> 'DataFilesDict':
        out = DataFilesDict()
        for key, data_files_patterns_list in self.items():
            out[key] = data_files_patterns_list.resolve(base_path, download_config)
        return out

    def filter_extensions(self, extensions: list[str]) -> 'DataFilesPatternsDict':
        out = type(self)()
        for key, data_files_patterns_list in self.items():
            out[key] = data_files_patterns_list.filter_extensions(extensions)
        return out
