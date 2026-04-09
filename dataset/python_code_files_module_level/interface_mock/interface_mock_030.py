import logging
import os
import secrets
import shutil
import tempfile
import uuid
from contextlib import suppress
from datetime import datetime
from urllib.parse import quote
import requests
from ..spec import AbstractBufferedFile, AbstractFileSystem
from ..utils import infer_storage_options, tokenize
logger = logging.getLogger('webhdfs')

class WebHDFS(AbstractFileSystem):
    tempdir = str(tempfile.gettempdir())
    protocol = ('webhdfs', 'webHDFS')

    def __init__(self, host, port=50070, kerberos=False, token=None, user=None, password=None, proxy_to=None, kerb_kwargs=None, data_proxy=None, use_https=False, session_cert=None, session_verify=True, **kwargs):
        if self._cached:
            return
        super().__init__(**kwargs)
        self.url = f"{('https' if use_https else 'http')}://{host}:{port}/webhdfs/v1"
        self.kerb = kerberos
        self.kerb_kwargs = kerb_kwargs or {}
        self.pars = {}
        self.proxy = data_proxy or {}
        if token is not None:
            if user is not None or proxy_to is not None:
                raise ValueError('If passing a delegation token, must not set user or proxy_to, as these are encoded in the token')
            self.pars['delegation'] = token
        self.user = user
        self.password = password
        if password is not None:
            if user is None:
                raise ValueError('If passing a password, the user must also beset in order to set up the basic-auth')
        elif user is not None:
            self.pars['user.name'] = user
        if proxy_to is not None:
            self.pars['doas'] = proxy_to
        if kerberos and user is not None:
            raise ValueError('If using Kerberos auth, do not specify the user, this is handled by kinit.')
        self.session_cert = session_cert
        self.session_verify = session_verify
        self._connect()
        self._fsid = f'webhdfs_{tokenize(host, port)}'

    @property
    def fsid(self):
        return self._fsid

    def _connect(self):
        self.session = requests.Session()
        if self.session_cert:
            self.session.cert = self.session_cert
        self.session.verify = self.session_verify
        if self.kerb:
            from requests_kerberos import HTTPKerberosAuth
            self.session.auth = HTTPKerberosAuth(**self.kerb_kwargs)
        if self.user is not None and self.password is not None:
            from requests.auth import HTTPBasicAuth
            self.session.auth = HTTPBasicAuth(self.user, self.password)

    def _call(self, op, method='get', path=None, data=None, redirect=True, **kwargs):
        path = self._strip_protocol(path) if path is not None else ''
        url = self._apply_proxy(self.url + quote(path, safe='/='))
        args = kwargs.copy()
        args.update(self.pars)
        args['op'] = op.upper()
        logger.debug('sending %s with %s', url, method)
        out = self.session.request(method=method.upper(), url=url, params=args, data=data, allow_redirects=redirect)
        if out.status_code in [400, 401, 403, 404, 500]:
            try:
                err = out.json()
                msg = err['RemoteException']['message']
                exp = err['RemoteException']['exception']
            except (ValueError, KeyError):
                pass
            else:
                if exp in ['IllegalArgumentException', 'UnsupportedOperationException']:
                    raise ValueError(msg)
                elif exp in ['SecurityException', 'AccessControlException']:
                    raise PermissionError(msg)
                elif exp in ['FileNotFoundException']:
                    raise FileNotFoundError(msg)
                else:
                    raise RuntimeError(msg)
        out.raise_for_status()
        return out

    def _open(self, path, mode='rb', block_size=None, autocommit=True, replication=None, permissions=None, **kwargs):
        block_size = block_size or self.blocksize
        return WebHDFile(self, path, mode=mode, block_size=block_size, tempdir=self.tempdir, autocommit=autocommit, replication=replication, permissions=permissions)

    @staticmethod
    def _process_info(info):
        info['type'] = info['type'].lower()
        info['size'] = info['length']
        return info

    @classmethod
    def _strip_protocol(cls, path):
        return infer_storage_options(path)['path']

    @staticmethod
    def _get_kwargs_from_urls(urlpath):
        out = infer_storage_options(urlpath)
        out.pop('path', None)
        out.pop('protocol', None)
        if 'username' in out:
            out['user'] = out.pop('username')
        return out

    def info(self, path):
        out = self._call('GETFILESTATUS', path=path)
        info = out.json()['FileStatus']
        info['name'] = path
        return self._process_info(info)

    def created(self, path):
        info = self.info(path)
        mtime = info.get('modificationTime', None)
        if mtime is not None:
            return datetime.fromtimestamp(mtime / 1000)
        raise RuntimeError('Could not retrieve creation time (modification time).')

    def modified(self, path):
        info = self.info(path)
        mtime = info.get('modificationTime', None)
        if mtime is not None:
            return datetime.fromtimestamp(mtime / 1000)
        raise RuntimeError('Could not retrieve modification time.')

    def ls(self, path, detail=False, **kwargs):
        out = self._call('LISTSTATUS', path=path)
        infos = out.json()['FileStatuses']['FileStatus']
        for info in infos:
            self._process_info(info)
            info['name'] = path.rstrip('/') + '/' + info['pathSuffix']
        if detail:
            return sorted(infos, key=lambda i: i['name'])
        else:
            return sorted((info['name'] for info in infos))

    def content_summary(self, path):
        out = self._call('GETCONTENTSUMMARY', path=path)
        return out.json()['ContentSummary']

    def ukey(self, path):
        out = self._call('GETFILECHECKSUM', path=path, redirect=False)
        if 'Location' in out.headers:
            location = self._apply_proxy(out.headers['Location'])
            out2 = self.session.get(location)
            out2.raise_for_status()
            return out2.json()['FileChecksum']
        else:
            out.raise_for_status()
            return out.json()['FileChecksum']

    def home_directory(self):
        out = self._call('GETHOMEDIRECTORY')
        return out.json()['Path']

    def get_delegation_token(self, renewer=None):
        if renewer:
            out = self._call('GETDELEGATIONTOKEN', renewer=renewer)
        else:
            out = self._call('GETDELEGATIONTOKEN')
        t = out.json()['Token']
        if t is None:
            raise ValueError('No token available for this user/security context')
        return t['urlString']

    def renew_delegation_token(self, token):
        out = self._call('RENEWDELEGATIONTOKEN', method='put', token=token)
        return out.json()['long']

    def cancel_delegation_token(self, token):
        self._call('CANCELDELEGATIONTOKEN', method='put', token=token)

    def chmod(self, path, mod):
        self._call('SETPERMISSION', method='put', path=path, permission=mod)

    def chown(self, path, owner=None, group=None):
        kwargs = {}
        if owner is not None:
            kwargs['owner'] = owner
        if group is not None:
            kwargs['group'] = group
        self._call('SETOWNER', method='put', path=path, **kwargs)

    def set_replication(self, path, replication):
        self._call('SETREPLICATION', path=path, method='put', replication=replication)

    def mkdir(self, path, **kwargs):
        self._call('MKDIRS', method='put', path=path)

    def makedirs(self, path, exist_ok=False):
        if exist_ok is False and self.exists(path):
            raise FileExistsError(path)
        self.mkdir(path)

    def mv(self, path1, path2, **kwargs):
        self._call('RENAME', method='put', path=path1, destination=path2)

    def rm(self, path, recursive=False, **kwargs):
        self._call('DELETE', method='delete', path=path, recursive='true' if recursive else 'false')

    def rm_file(self, path, **kwargs):
        self.rm(path)

    def cp_file(self, lpath, rpath, **kwargs):
        with self.open(lpath) as lstream:
            tmp_fname = '/'.join([self._parent(rpath), f'.tmp.{secrets.token_hex(16)}'])
            try:
                with self.open(tmp_fname, 'wb') as rstream:
                    shutil.copyfileobj(lstream, rstream)
                self.mv(tmp_fname, rpath)
            except BaseException:
                with suppress(FileNotFoundError):
                    self.rm(tmp_fname)
                raise

    def _apply_proxy(self, location):
        if self.proxy and callable(self.proxy):
            location = self.proxy(location)
        elif self.proxy:
            for k, v in self.proxy.items():
                location = location.replace(k, v, 1)
        return location

class WebHDFile(AbstractBufferedFile):

    def __init__(self, fs, path, **kwargs):
        super().__init__(fs, path, **kwargs)
        kwargs = kwargs.copy()
        if kwargs.get('permissions', None) is None:
            kwargs.pop('permissions', None)
        if kwargs.get('replication', None) is None:
            kwargs.pop('replication', None)
        self.permissions = kwargs.pop('permissions', 511)
        tempdir = kwargs.pop('tempdir')
        if kwargs.pop('autocommit', False) is False:
            self.target = self.path
            self.path = os.path.join(tempdir, str(uuid.uuid4()))

    def _upload_chunk(self, final=False):
        out = self.fs.session.post(self.location, data=self.buffer.getvalue(), headers={'content-type': 'application/octet-stream'})
        out.raise_for_status()
        return True

    def _initiate_upload(self):
        kwargs = self.kwargs.copy()
        if 'a' in self.mode:
            op, method = ('APPEND', 'POST')
        else:
            op, method = ('CREATE', 'PUT')
            kwargs['overwrite'] = 'true'
        out = self.fs._call(op, method, self.path, redirect=False, **kwargs)
        location = self.fs._apply_proxy(out.headers['Location'])
        if 'w' in self.mode:
            out2 = self.fs.session.put(location, headers={'content-type': 'application/octet-stream'})
            out2.raise_for_status()
            out2 = self.fs._call('APPEND', 'POST', self.path, redirect=False, **kwargs)
            self.location = self.fs._apply_proxy(out2.headers['Location'])

    def _fetch_range(self, start, end):
        start = max(start, 0)
        end = min(self.size, end)
        if start >= end or start >= self.size:
            return b''
        out = self.fs._call('OPEN', path=self.path, offset=start, length=end - start, redirect=False)
        out.raise_for_status()
        if 'Location' in out.headers:
            location = out.headers['Location']
            out2 = self.fs.session.get(self.fs._apply_proxy(location))
            return out2.content
        else:
            return out.content

    def commit(self):
        self.fs.mv(self.path, self.target)

    def discard(self):
        self.fs.rm(self.path)
