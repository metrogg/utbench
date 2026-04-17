import datetime
import re
import uuid
from stat import S_ISDIR, S_ISLNK
import smbclient
import smbprotocol.exceptions
from .. import AbstractFileSystem
from ..utils import infer_storage_options

class SMBFileSystem(AbstractFileSystem):
    protocol = 'smb'

    def __init__(self, host, port=None, username=None, password=None, timeout=60, encrypt=None, share_access=None, register_session_retries=4, register_session_retry_wait=1, register_session_retry_factor=10, auto_mkdir=False, **kwargs):
        super().__init__(**kwargs)
        self.host = host
        self.port = port
        self.username = username
        self.password = password
        self.timeout = timeout
        self.encrypt = encrypt
        self.temppath = kwargs.pop('temppath', '')
        self.share_access = share_access
        self.register_session_retries = register_session_retries
        if register_session_retry_wait < 0:
            raise ValueError('register_session_retry_wait must be a non-negative integer')
        self.register_session_retry_wait = register_session_retry_wait
        if register_session_retry_factor < 1:
            raise ValueError('register_session_retry_factor must be a positive integer equal to or greater than 1')
        self.register_session_retry_factor = register_session_retry_factor
        self.auto_mkdir = auto_mkdir
        self._connect()

    @property
    def _port(self):
        return 445 if self.port is None else self.port

    def _connect(self):
        import time
        if self.register_session_retries <= -1:
            return
        retried_errors = []
        wait_time = self.register_session_retry_wait
        n_waits = self.register_session_retries - 1
        factor = self.register_session_retry_factor
        wait_times = iter((factor ** (n / n_waits - 1) * wait_time for n in range(0, n_waits + 1)))
        for attempt in range(self.register_session_retries + 1):
            try:
                smbclient.register_session(self.host, username=self.username, password=self.password, port=self._port, encrypt=self.encrypt, connection_timeout=self.timeout)
                return
            except (smbprotocol.exceptions.SMBAuthenticationError, smbprotocol.exceptions.LogonFailure):
                raise
            except ValueError as exc:
                if re.findall('\\[Errno -\\d+]', str(exc)):
                    retried_errors.append(exc)
                else:
                    raise
            except Exception as exc:
                retried_errors.append(exc)
            if attempt < self.register_session_retries:
                time.sleep(next(wait_times))
        raise retried_errors[-1]

    @classmethod
    def _strip_protocol(cls, path):
        return infer_storage_options(path)['path']

    @staticmethod
    def _get_kwargs_from_urls(path):
        out = infer_storage_options(path)
        out.pop('path', None)
        out.pop('protocol', None)
        return out

    def mkdir(self, path, create_parents=True, **kwargs):
        wpath = _as_unc_path(self.host, path)
        if create_parents:
            smbclient.makedirs(wpath, exist_ok=False, port=self._port, **kwargs)
        else:
            smbclient.mkdir(wpath, port=self._port, **kwargs)

    def makedirs(self, path, exist_ok=False):
        if _share_has_path(path):
            wpath = _as_unc_path(self.host, path)
            smbclient.makedirs(wpath, exist_ok=exist_ok, port=self._port)

    def rmdir(self, path):
        if _share_has_path(path):
            wpath = _as_unc_path(self.host, path)
            smbclient.rmdir(wpath, port=self._port)

    def info(self, path, **kwargs):
        wpath = _as_unc_path(self.host, path)
        stats = smbclient.stat(wpath, port=self._port, **kwargs)
        if S_ISDIR(stats.st_mode):
            stype = 'directory'
        elif S_ISLNK(stats.st_mode):
            stype = 'link'
        else:
            stype = 'file'
        res = {'name': path + '/' if stype == 'directory' else path, 'size': stats.st_size, 'type': stype, 'uid': stats.st_uid, 'gid': stats.st_gid, 'time': stats.st_atime, 'mtime': stats.st_mtime}
        return res

    def created(self, path):
        wpath = _as_unc_path(self.host, path)
        stats = smbclient.stat(wpath, port=self._port)
        return datetime.datetime.fromtimestamp(stats.st_ctime, tz=datetime.timezone.utc)

    def modified(self, path):
        wpath = _as_unc_path(self.host, path)
        stats = smbclient.stat(wpath, port=self._port)
        return datetime.datetime.fromtimestamp(stats.st_mtime, tz=datetime.timezone.utc)

    def ls(self, path, detail=True, **kwargs):
        unc = _as_unc_path(self.host, path)
        listed = smbclient.listdir(unc, port=self._port, **kwargs)
        dirs = ['/'.join([path.rstrip('/'), p]) for p in listed]
        if detail:
            dirs = [self.info(d) for d in dirs]
        return dirs

    def _open(self, path, mode='rb', block_size=-1, autocommit=True, cache_options=None, **kwargs):
        if self.auto_mkdir and 'w' in mode:
            self.makedirs(self._parent(path), exist_ok=True)
        bls = block_size if block_size is not None and block_size >= 0 else -1
        wpath = _as_unc_path(self.host, path)
        share_access = kwargs.pop('share_access', self.share_access)
        if 'w' in mode and autocommit is False:
            temp = _as_temp_path(self.host, path, self.temppath)
            return SMBFileOpener(wpath, temp, mode, port=self._port, block_size=bls, **kwargs)
        return smbclient.open_file(wpath, mode, buffering=bls, share_access=share_access, port=self._port, **kwargs)

    def copy(self, path1, path2, **kwargs):
        wpath1 = _as_unc_path(self.host, path1)
        wpath2 = _as_unc_path(self.host, path2)
        if self.auto_mkdir:
            self.makedirs(self._parent(path2), exist_ok=True)
        smbclient.copyfile(wpath1, wpath2, port=self._port, **kwargs)

    def _rm(self, path):
        if _share_has_path(path):
            wpath = _as_unc_path(self.host, path)
            stats = smbclient.stat(wpath, port=self._port)
            if S_ISDIR(stats.st_mode):
                smbclient.rmdir(wpath, port=self._port)
            else:
                smbclient.remove(wpath, port=self._port)

    def mv(self, path1, path2, recursive=None, maxdepth=None, **kwargs):
        wpath1 = _as_unc_path(self.host, path1)
        wpath2 = _as_unc_path(self.host, path2)
        smbclient.rename(wpath1, wpath2, port=self._port, **kwargs)

def _as_unc_path(host, path):
    rpath = path.replace('/', '\\')
    unc = f'\\\\{host}{rpath}'
    return unc

def _as_temp_path(host, path, temppath):
    share = path.split('/')[1]
    temp_file = f'/{share}{temppath}/{uuid.uuid4()}'
    unc = _as_unc_path(host, temp_file)
    return unc

def _share_has_path(path):
    parts = path.count('/')
    if path.endswith('/'):
        return parts > 2
    return parts > 1

class SMBFileOpener:

    def __init__(self, path, temp, mode, port=445, block_size=-1, **kwargs):
        self.path = path
        self.temp = temp
        self.mode = mode
        self.block_size = block_size
        self.kwargs = kwargs
        self.smbfile = None
        self._incontext = False
        self.port = port
        self._open()

    def _open(self):
        if self.smbfile is None or self.smbfile.closed:
            self.smbfile = smbclient.open_file(self.temp, self.mode, port=self.port, buffering=self.block_size, **self.kwargs)

    def commit(self):
        smbclient.replace(self.temp, self.path, port=self.port)

    def discard(self):
        smbclient.remove(self.temp, port=self.port)

    def __fspath__(self):
        return self.path

    def __iter__(self):
        return self.smbfile.__iter__()

    def __getattr__(self, item):
        return getattr(self.smbfile, item)

    def __enter__(self):
        self._incontext = True
        return self.smbfile.__enter__()

    def __exit__(self, exc_type, exc_value, traceback):
        self._incontext = False
        self.smbfile.__exit__(exc_type, exc_value, traceback)
