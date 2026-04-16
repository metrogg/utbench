import os
from numpy._utils import set_module
_open = open

def _check_mode(mode, encoding, newline):
    if 't' in mode:
        if 'b' in mode:
            raise ValueError(f'Invalid mode: {mode!r}')
    else:
        if encoding is not None:
            raise ValueError("Argument 'encoding' not supported in binary mode")
        if newline is not None:
            raise ValueError("Argument 'newline' not supported in binary mode")

class _FileOpeners:

    def __init__(self):
        self._loaded = False
        self._file_openers = {None: open}

    def _load(self):
        if self._loaded:
            return
        try:
            import bz2
            self._file_openers['.bz2'] = bz2.open
        except ImportError:
            pass
        try:
            import gzip
            self._file_openers['.gz'] = gzip.open
        except ImportError:
            pass
        try:
            import lzma
            self._file_openers['.xz'] = lzma.open
            self._file_openers['.lzma'] = lzma.open
        except (ImportError, AttributeError):
            pass
        self._loaded = True

    def keys(self):
        self._load()
        return list(self._file_openers.keys())

    def __getitem__(self, key):
        self._load()
        return self._file_openers[key]
_file_openers = _FileOpeners()

def open(path, mode='r', destpath=os.curdir, encoding=None, newline=None):
    ds = DataSource(destpath)
    return ds.open(path, mode, encoding=encoding, newline=newline)

@set_module('numpy.lib.npyio')
class DataSource:

    def __init__(self, destpath=os.curdir):
        if destpath:
            self._destpath = os.path.abspath(destpath)
            self._istmpdest = False
        else:
            import tempfile
            self._destpath = tempfile.mkdtemp()
            self._istmpdest = True

    def __del__(self):
        if hasattr(self, '_istmpdest') and self._istmpdest:
            import shutil
            shutil.rmtree(self._destpath)

    def _iszip(self, filename):
        fname, ext = os.path.splitext(filename)
        return ext in _file_openers.keys()

    def _iswritemode(self, mode):
        _writemodes = ('w', '+')
        return any((c in _writemodes for c in mode))

    def _splitzipext(self, filename):
        if self._iszip(filename):
            return os.path.splitext(filename)
        else:
            return (filename, None)

    def _possible_names(self, filename):
        names = [filename]
        if not self._iszip(filename):
            for zipext in _file_openers.keys():
                if zipext:
                    names.append(filename + zipext)
        return names

    def _isurl(self, path):
        from urllib.parse import urlparse
        scheme, netloc, upath, uparams, uquery, ufrag = urlparse(path)
        return bool(scheme and netloc)

    def _cache(self, path):
        import shutil
        from urllib.request import urlopen
        upath = self.abspath(path)
        if not os.path.exists(os.path.dirname(upath)):
            os.makedirs(os.path.dirname(upath))
        if self._isurl(path):
            with urlopen(path) as openedurl:
                with _open(upath, 'wb') as f:
                    shutil.copyfileobj(openedurl, f)
        else:
            shutil.copyfile(path, upath)
        return upath

    def _findfile(self, path):
        if not self._isurl(path):
            filelist = self._possible_names(path)
            filelist += self._possible_names(self.abspath(path))
        else:
            filelist = self._possible_names(self.abspath(path))
            filelist = filelist + self._possible_names(path)
        for name in filelist:
            if self.exists(name):
                if self._isurl(name):
                    name = self._cache(name)
                return name
        return None

    def abspath(self, path):
        from urllib.parse import urlparse
        splitpath = path.split(self._destpath, 2)
        if len(splitpath) > 1:
            path = splitpath[1]
        scheme, netloc, upath, uparams, uquery, ufrag = urlparse(path)
        netloc = self._sanitize_relative_path(netloc)
        upath = self._sanitize_relative_path(upath)
        return os.path.join(self._destpath, netloc, upath)

    def _sanitize_relative_path(self, path):
        last = None
        path = os.path.normpath(path)
        while path != last:
            last = path
            path = path.lstrip(os.sep).lstrip('/')
            path = path.lstrip(os.pardir).removeprefix('..')
            drive, path = os.path.splitdrive(path)
        return path

    def exists(self, path):
        if os.path.exists(path):
            return True
        from urllib.error import URLError
        from urllib.request import urlopen
        upath = self.abspath(path)
        if os.path.exists(upath):
            return True
        if self._isurl(path):
            try:
                netfile = urlopen(path)
                netfile.close()
                del netfile
                return True
            except URLError:
                return False
        return False

    def open(self, path, mode='r', encoding=None, newline=None):
        if self._isurl(path) and self._iswritemode(mode):
            raise ValueError('URLs are not writeable')
        found = self._findfile(path)
        if found:
            _fname, ext = self._splitzipext(found)
            if ext == 'bz2':
                mode.replace('+', '')
            return _file_openers[ext](found, mode=mode, encoding=encoding, newline=newline)
        else:
            raise FileNotFoundError(f'{path} not found.')

class Repository(DataSource):

    def __init__(self, baseurl, destpath=os.curdir):
        DataSource.__init__(self, destpath=destpath)
        self._baseurl = baseurl

    def __del__(self):
        DataSource.__del__(self)

    def _fullpath(self, path):
        splitpath = path.split(self._baseurl, 2)
        if len(splitpath) == 1:
            result = os.path.join(self._baseurl, path)
        else:
            result = path
        return result

    def _findfile(self, path):
        return DataSource._findfile(self, self._fullpath(path))

    def abspath(self, path):
        return DataSource.abspath(self, self._fullpath(path))

    def exists(self, path):
        return DataSource.exists(self, self._fullpath(path))

    def open(self, path, mode='r', encoding=None, newline=None):
        return DataSource.open(self, self._fullpath(path), mode, encoding=encoding, newline=newline)

    def listdir(self):
        if self._isurl(self._baseurl):
            raise NotImplementedError('Directory listing of URLs, not supported yet.')
        else:
            return os.listdir(self._baseurl)
