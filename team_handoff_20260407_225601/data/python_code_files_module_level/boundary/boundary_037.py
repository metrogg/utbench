import array
import logging
import posixpath
import warnings
from collections.abc import MutableMapping
from functools import cached_property
from fsspec.core import url_to_fs
logger = logging.getLogger('fsspec.mapping')

class FSMap(MutableMapping):

    def __init__(self, root, fs, check=False, create=False, missing_exceptions=None):
        self.fs = fs
        self.root = fs._strip_protocol(root)
        self._root_key_to_str = fs._strip_protocol(posixpath.join(root, 'x'))[:-1]
        if missing_exceptions is None:
            missing_exceptions = (FileNotFoundError, IsADirectoryError, NotADirectoryError)
        self.missing_exceptions = missing_exceptions
        self.check = check
        self.create = create
        if create:
            if not self.fs.exists(root):
                self.fs.mkdir(root)
        if check:
            if not self.fs.exists(root):
                raise ValueError(f'Path {root} does not exist. Create  with the ``create=True`` keyword')
            self.fs.touch(root + '/a')
            self.fs.rm(root + '/a')

    @cached_property
    def dirfs(self):
        from .implementations.dirfs import DirFileSystem
        return DirFileSystem(path=self._root_key_to_str, fs=self.fs)

    def clear(self):
        logger.info('Clear mapping at %s', self.root)
        try:
            self.fs.rm(self.root, True)
            self.fs.mkdir(self.root)
        except:
            pass

    def getitems(self, keys, on_error='raise'):
        keys2 = [self._key_to_str(k) for k in keys]
        oe = on_error if on_error == 'raise' else 'return'
        try:
            out = self.fs.cat(keys2, on_error=oe)
            if isinstance(out, bytes):
                out = {keys2[0]: out}
        except self.missing_exceptions as e:
            raise KeyError from e
        out = {k: KeyError() if isinstance(v, self.missing_exceptions) else v for k, v in out.items()}
        return {key: out[k2] if on_error == 'raise' else out.get(k2, KeyError(k2)) for key, k2 in zip(keys, keys2) if on_error == 'return' or not isinstance(out[k2], BaseException)}

    def setitems(self, values_dict):
        values = {self._key_to_str(k): maybe_convert(v) for k, v in values_dict.items()}
        self.fs.pipe(values)

    def delitems(self, keys):
        self.fs.rm([self._key_to_str(k) for k in keys])

    def _key_to_str(self, key):
        if not isinstance(key, str):
            warnings.warn('from fsspec 2023.5 onward FSMap non-str keys will raise TypeError', DeprecationWarning)
            if isinstance(key, list):
                key = tuple(key)
            key = str(key)
        return f'{self._root_key_to_str}{key}'.rstrip('/')

    def _str_to_key(self, s):
        return s[len(self.root):].lstrip('/')

    def __getitem__(self, key, default=None):
        k = self._key_to_str(key)
        try:
            result = self.fs.cat(k)
        except self.missing_exceptions as exc:
            if default is not None:
                return default
            raise KeyError(key) from exc
        return result

    def pop(self, key, default=None):
        result = self.__getitem__(key, default)
        try:
            del self[key]
        except KeyError:
            pass
        return result

    def __setitem__(self, key, value):
        key = self._key_to_str(key)
        self.fs.mkdirs(self.fs._parent(key), exist_ok=True)
        self.fs.pipe_file(key, maybe_convert(value))

    def __iter__(self):
        return (self._str_to_key(x) for x in self.fs.find(self.root))

    def __len__(self):
        return len(self.fs.find(self.root))

    def __delitem__(self, key):
        try:
            self.fs.rm(self._key_to_str(key))
        except Exception as exc:
            raise KeyError from exc

    def __contains__(self, key):
        path = self._key_to_str(key)
        return self.fs.isfile(path)

    def __reduce__(self):
        return (FSMap, (self.root, self.fs, False, False, self.missing_exceptions))

def maybe_convert(value):
    if isinstance(value, array.array) or hasattr(value, '__array__'):
        if hasattr(value, 'dtype') and value.dtype.kind in 'Mm':
            value = value.view('int64')
        value = bytes(memoryview(value))
    return value

def get_mapper(url='', check=False, create=False, missing_exceptions=None, alternate_root=None, **kwargs):
    fs, urlpath = url_to_fs(url, **kwargs)
    root = alternate_root if alternate_root is not None else urlpath
    return FSMap(root, fs, check, create, missing_exceptions=missing_exceptions)
