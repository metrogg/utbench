from __future__ import annotations
import io
import logging
import os
import re
from glob import has_magic
from pathlib import Path
from fsspec.caching import BaseCache, BlockCache, BytesCache, MMapCache, ReadAheadCache, caches
from fsspec.compression import compr
from fsspec.config import conf
from fsspec.registry import available_protocols, filesystem, get_filesystem_class
from fsspec.utils import _unstrip_protocol, build_name_function, infer_compression, stringify_path
logger = logging.getLogger('fsspec')

class OpenFile:

    def __init__(self, fs, path, mode='rb', compression=None, encoding=None, errors=None, newline=None):
        self.fs = fs
        self.path = path
        self.mode = mode
        self.compression = get_compression(path, compression)
        self.encoding = encoding
        self.errors = errors
        self.newline = newline
        self.fobjects = []

    def __reduce__(self):
        return (OpenFile, (self.fs, self.path, self.mode, self.compression, self.encoding, self.errors, self.newline))

    def __repr__(self):
        return f"<OpenFile '{self.path}'>"

    def __enter__(self):
        mode = self.mode.replace('t', '').replace('b', '') + 'b'
        try:
            f = self.fs.open(self.path, mode=mode)
        except FileNotFoundError as e:
            if has_magic(self.path):
                raise FileNotFoundError("%s not found. The URL contains glob characters: you maybe needed\nto pass expand=True in fsspec.open() or the storage_options of \nyour library. You can also set the config value 'open_expand'\nbefore import, or fsspec.core.DEFAULT_EXPAND at runtime, to True.", self.path) from e
            raise
        self.fobjects = [f]
        if self.compression is not None:
            compress = compr[self.compression]
            f = compress(f, mode=mode[0])
            self.fobjects.append(f)
        if 'b' not in self.mode:
            f = PickleableTextIOWrapper(f, encoding=self.encoding, errors=self.errors, newline=self.newline)
            self.fobjects.append(f)
        return self.fobjects[-1]

    def __exit__(self, *args):
        self.close()

    @property
    def full_name(self):
        return _unstrip_protocol(self.path, self.fs)

    def open(self):
        return self.__enter__()

    def close(self):
        for f in reversed(self.fobjects):
            if 'r' not in self.mode and (not f.closed):
                f.flush()
            f.close()
        self.fobjects.clear()

class OpenFiles(list):

    def __init__(self, *args, mode='rb', fs=None):
        self.mode = mode
        self.fs = fs
        self.files = []
        super().__init__(*args)

    def __enter__(self):
        if self.fs is None:
            raise ValueError('Context has already been used')
        fs = self.fs
        while True:
            if hasattr(fs, 'open_many'):
                self.files = fs.open_many(self)
                return self.files
            if hasattr(fs, 'fs') and fs.fs is not None:
                fs = fs.fs
            else:
                break
        return [s.__enter__() for s in self]

    def __exit__(self, *args):
        fs = self.fs
        [s.__exit__(*args) for s in self]
        if 'r' not in self.mode:
            while True:
                if hasattr(fs, 'open_many'):
                    fs.commit_many(self.files)
                    return
                if hasattr(fs, 'fs') and fs.fs is not None:
                    fs = fs.fs
                else:
                    break

    def __getitem__(self, item):
        out = super().__getitem__(item)
        if isinstance(item, slice):
            return OpenFiles(out, mode=self.mode, fs=self.fs)
        return out

    def __repr__(self):
        return f'<List of {len(self)} OpenFile instances>'

def open_files(urlpath, mode='rb', compression=None, encoding='utf8', errors=None, name_function=None, num=1, protocol=None, newline=None, auto_mkdir=True, expand=True, **kwargs):
    fs, fs_token, paths = get_fs_token_paths(urlpath, mode, num=num, name_function=name_function, storage_options=kwargs, protocol=protocol, expand=expand)
    if fs.protocol == 'file':
        fs.auto_mkdir = auto_mkdir
    elif 'r' not in mode and auto_mkdir:
        parents = {fs._parent(path) for path in paths}
        for parent in parents:
            try:
                fs.makedirs(parent, exist_ok=True)
            except PermissionError:
                pass
    return OpenFiles([OpenFile(fs, path, mode=mode, compression=compression, encoding=encoding, errors=errors, newline=newline) for path in paths], mode=mode, fs=fs)

def _un_chain(path, kwargs):
    from fsspec.implementations.chained import ChainedFileSystem
    if '::' in path:
        x = re.compile('.*[^a-z]+.*')
        known_protocols = set(available_protocols())
        bits = []
        for p in path.split('::'):
            if p in known_protocols:
                bits.append(p + '://')
            elif '://' in p or x.match(p):
                bits.append(p)
            else:
                bits.append(p + '://')
    else:
        bits = [path]
    out = []
    previous_bit = None
    kwargs = kwargs.copy()
    for bit in reversed(bits):
        protocol = kwargs.pop('protocol', None) or split_protocol(bit)[0] or 'file'
        cls = get_filesystem_class(protocol)
        extra_kwargs = cls._get_kwargs_from_urls(bit)
        kws = kwargs.pop(protocol, {})
        if bit is bits[0]:
            kws.update(kwargs)
        kw = dict(**{k: v for k, v in extra_kwargs.items() if k not in kws or v != kws[k]}, **kws)
        bit = cls._strip_protocol(bit)
        if 'target_protocol' not in kw and issubclass(cls, ChainedFileSystem) and (not bit):
            bit = previous_bit
        out.append((bit, protocol, kw))
        previous_bit = bit
    out.reverse()
    return out

def url_to_fs(url, **kwargs):
    url = stringify_path(url)
    known_kwargs = {'compression', 'encoding', 'errors', 'expand', 'mode', 'name_function', 'newline', 'num'}
    kwargs = {k: v for k, v in kwargs.items() if k not in known_kwargs}
    chain = _un_chain(url, kwargs)
    inkwargs = {}
    for i, ch in enumerate(reversed(chain)):
        urls, protocol, kw = ch
        if i == len(chain) - 1:
            inkwargs = dict(**kw, **inkwargs)
            continue
        inkwargs['target_options'] = dict(**kw, **inkwargs)
        inkwargs['target_protocol'] = protocol
        inkwargs['fo'] = urls
    urlpath, protocol, _ = chain[0]
    fs = filesystem(protocol, **inkwargs)
    return (fs, urlpath)
DEFAULT_EXPAND = conf.get('open_expand', False)

def open(urlpath, mode='rb', compression=None, encoding='utf8', errors=None, protocol=None, newline=None, expand=None, **kwargs):
    expand = DEFAULT_EXPAND if expand is None else expand
    out = open_files(urlpath=[urlpath], mode=mode, compression=compression, encoding=encoding, errors=errors, protocol=protocol, newline=newline, expand=expand, **kwargs)
    if not out:
        raise FileNotFoundError(urlpath)
    return out[0]

def open_local(url: str | list[str] | Path | list[Path], mode: str='rb', **storage_options: dict) -> str | list[str]:
    if 'r' not in mode:
        raise ValueError('Can only ensure local files when reading')
    of = open_files(url, mode=mode, **storage_options)
    if not getattr(of[0].fs, 'local_file', False):
        raise ValueError('open_local can only be used on a filesystem which has attribute local_file=True')
    with of as files:
        paths = [f.name for f in files]
    if isinstance(url, str) and (not has_magic(url)) or isinstance(url, Path):
        return paths[0]
    return paths

def get_compression(urlpath, compression):
    if compression == 'infer':
        compression = infer_compression(urlpath)
    if compression is not None and compression not in compr:
        raise ValueError(f'Compression type {compression} not supported')
    return compression

def split_protocol(urlpath):
    urlpath = stringify_path(urlpath)
    if '://' in urlpath:
        protocol, path = urlpath.split('://', 1)
        if len(protocol) > 1:
            return (protocol, path)
    if urlpath.startswith('data:'):
        return urlpath.split(':', 1)
    return (None, urlpath)

def strip_protocol(urlpath):
    protocol, _ = split_protocol(urlpath)
    cls = get_filesystem_class(protocol)
    return cls._strip_protocol(urlpath)

def expand_paths_if_needed(paths, mode, num, fs, name_function):
    expanded_paths = []
    paths = list(paths)
    if 'w' in mode:
        if sum((1 for p in paths if '*' in p)) > 1:
            raise ValueError('When writing data, only one filename mask can be specified.')
        num = max(num, len(paths))
        for curr_path in paths:
            if '*' in curr_path:
                expanded_paths.extend(_expand_paths(curr_path, name_function, num))
            else:
                expanded_paths.append(curr_path)
        if len(expanded_paths) > num:
            expanded_paths = expanded_paths[:num]
    else:
        for curr_path in paths:
            if has_magic(curr_path):
                expanded_paths.extend(fs.glob(curr_path))
            else:
                expanded_paths.append(curr_path)
    return expanded_paths

def get_fs_token_paths(urlpath, mode='rb', num=1, name_function=None, storage_options=None, protocol=None, expand=True):
    if isinstance(urlpath, (list, tuple, set)):
        if not urlpath:
            raise ValueError('empty urlpath sequence')
        urlpath0 = stringify_path(next(iter(urlpath)))
    else:
        urlpath0 = stringify_path(urlpath)
    storage_options = storage_options or {}
    if protocol:
        storage_options['protocol'] = protocol
    chain = _un_chain(urlpath0, storage_options or {})
    inkwargs = {}
    for i, ch in enumerate(reversed(chain)):
        urls, nested_protocol, kw = ch
        if i == len(chain) - 1:
            inkwargs = dict(**kw, **inkwargs)
            continue
        inkwargs['target_options'] = dict(**kw, **inkwargs)
        inkwargs['target_protocol'] = nested_protocol
        inkwargs['fo'] = urls
    paths, protocol, _ = chain[0]
    fs = filesystem(protocol, **inkwargs)
    if isinstance(urlpath, (list, tuple, set)):
        pchains = [_un_chain(stringify_path(u), storage_options or {})[0] for u in urlpath]
        if len({pc[1] for pc in pchains}) > 1:
            raise ValueError('Protocol mismatch getting fs from %s', urlpath)
        paths = [pc[0] for pc in pchains]
    else:
        paths = fs._strip_protocol(paths)
    if isinstance(paths, (list, tuple, set)):
        if expand:
            paths = expand_paths_if_needed(paths, mode, num, fs, name_function)
        elif not isinstance(paths, list):
            paths = list(paths)
    elif ('w' in mode or 'x' in mode) and expand:
        paths = _expand_paths(paths, name_function, num)
    elif '*' in paths:
        paths = [f for f in sorted(fs.glob(paths)) if not fs.isdir(f)]
    else:
        paths = [paths]
    return (fs, fs._fs_token, paths)

def _expand_paths(path, name_function, num):
    if isinstance(path, str):
        if path.count('*') > 1:
            raise ValueError("Output path spec must contain exactly one '*'.")
        elif '*' not in path:
            path = os.path.join(path, '*.part')
        if name_function is None:
            name_function = build_name_function(num - 1)
        paths = [path.replace('*', name_function(i)) for i in range(num)]
        if paths != sorted(paths):
            logger.warning('In order to preserve order between partitions paths created with ``name_function`` should sort to partition order')
    elif isinstance(path, (tuple, list)):
        assert len(path) == num
        paths = list(path)
    else:
        raise ValueError("Path should be either\n1. A list of paths: ['foo.json', 'bar.json', ...]\n2. A directory: 'foo/\n3. A path with a '*' in it: 'foo.*.json'")
    return paths

class PickleableTextIOWrapper(io.TextIOWrapper):

    def __init__(self, buffer, encoding=None, errors=None, newline=None, line_buffering=False, write_through=False):
        self.args = (buffer, encoding, errors, newline, line_buffering, write_through)
        super().__init__(*self.args)

    def __reduce__(self):
        return (PickleableTextIOWrapper, self.args)
