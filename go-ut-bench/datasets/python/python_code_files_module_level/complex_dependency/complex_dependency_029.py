import base64
import re
import requests
from ..spec import AbstractFileSystem
from ..utils import infer_storage_options
from .memory import MemoryFile

class GithubFileSystem(AbstractFileSystem):
    url = 'https://api.github.com/repos/{org}/{repo}/git/trees/{sha}'
    content_url = 'https://api.github.com/repos/{org}/{repo}/contents/{path}?ref={sha}'
    protocol = 'github'
    timeout = (60, 60)

    def __init__(self, org, repo, sha=None, username=None, token=None, timeout=None, **kwargs):
        super().__init__(**kwargs)
        self.org = org
        self.repo = repo
        if (username is None) ^ (token is None):
            raise ValueError('Auth required both username and token')
        self.username = username
        self.token = token
        if timeout is not None:
            self.timeout = timeout
        if sha is None:
            u = 'https://api.github.com/repos/{org}/{repo}'
            r = requests.get(u.format(org=org, repo=repo), timeout=self.timeout, **self.kw)
            r.raise_for_status()
            sha = r.json()['default_branch']
        self.root = sha
        self.ls('')
        try:
            from .http import HTTPFileSystem
            self.http_fs = HTTPFileSystem(**kwargs)
        except ImportError:
            self.http_fs = None

    @property
    def kw(self):
        if self.username:
            return {'auth': (self.username, self.token)}
        return {}

    @classmethod
    def repos(cls, org_or_user, is_org=True):
        r = requests.get(f"https://api.github.com/{['users', 'orgs'][is_org]}/{org_or_user}/repos", timeout=cls.timeout)
        r.raise_for_status()
        return [repo['name'] for repo in r.json()]

    @property
    def tags(self):
        r = requests.get(f'https://api.github.com/repos/{self.org}/{self.repo}/tags', timeout=self.timeout, **self.kw)
        r.raise_for_status()
        return [t['name'] for t in r.json()]

    @property
    def branches(self):
        r = requests.get(f'https://api.github.com/repos/{self.org}/{self.repo}/branches', timeout=self.timeout, **self.kw)
        r.raise_for_status()
        return [t['name'] for t in r.json()]

    @property
    def refs(self):
        return {'tags': self.tags, 'branches': self.branches}

    def ls(self, path, detail=False, sha=None, _sha=None, **kwargs):
        path = self._strip_protocol(path)
        if path == '':
            _sha = sha or self.root
        if _sha is None:
            parts = path.rstrip('/').split('/')
            so_far = ''
            _sha = sha or self.root
            for part in parts:
                out = self.ls(so_far, True, sha=sha, _sha=_sha)
                so_far += '/' + part if so_far else part
                out = [o for o in out if o['name'] == so_far]
                if not out:
                    raise FileNotFoundError(path)
                out = out[0]
                if out['type'] == 'file':
                    if detail:
                        return [out]
                    else:
                        return path
                _sha = out['sha']
        if path not in self.dircache or sha not in [self.root, None]:
            r = requests.get(self.url.format(org=self.org, repo=self.repo, sha=_sha), timeout=self.timeout, **self.kw)
            if r.status_code == 404:
                raise FileNotFoundError(path)
            r.raise_for_status()
            types = {'blob': 'file', 'tree': 'directory'}
            out = [{'name': path + '/' + f['path'] if path else f['path'], 'mode': f['mode'], 'type': types[f['type']], 'size': f.get('size', 0), 'sha': f['sha']} for f in r.json()['tree'] if f['type'] in types]
            if sha in [self.root, None]:
                self.dircache[path] = out
        else:
            out = self.dircache[path]
        if detail:
            return out
        else:
            return sorted([f['name'] for f in out])

    def invalidate_cache(self, path=None):
        self.dircache.clear()

    @classmethod
    def _strip_protocol(cls, path):
        opts = infer_storage_options(path)
        if 'username' not in opts:
            return super()._strip_protocol(path)
        return opts['path'].lstrip('/')

    @staticmethod
    def _get_kwargs_from_urls(path):
        opts = infer_storage_options(path)
        if 'username' not in opts:
            return {}
        out = {'org': opts['username'], 'repo': opts['password']}
        if opts['host']:
            out['sha'] = opts['host']
        return out

    def _open(self, path, mode='rb', block_size=None, cache_options=None, sha=None, **kwargs):
        if mode != 'rb':
            raise NotImplementedError
        url = self.content_url.format(org=self.org, repo=self.repo, path=path, sha=sha or self.root)
        r = requests.get(url, timeout=self.timeout, **self.kw)
        if r.status_code == 404:
            raise FileNotFoundError(path)
        r.raise_for_status()
        content_json = r.json()
        if content_json['content']:
            content = base64.b64decode(content_json['content'])
            if not content.startswith(b'version https://git-lfs.github.com/'):
                return MemoryFile(None, None, content)
        if self.http_fs is None:
            raise ImportError('Please install fsspec[http] to access github files >1 MB or git-lfs tracked files.')
        return self.http_fs.open(content_json['download_url'], mode=mode, block_size=block_size, cache_options=cache_options, **kwargs)

    def rm(self, path, recursive=False, maxdepth=None, message=None):
        path = self.expand_path(path, recursive=recursive, maxdepth=maxdepth)
        for p in reversed(path):
            self.rm_file(p, message=message)

    def rm_file(self, path, message=None, **kwargs):
        if not self.username:
            raise ValueError('Authentication required')
        path = self._strip_protocol(path)
        sha = self._get_sha_from_cache(path)
        if not sha:
            url = self.content_url.format(org=self.org, repo=self.repo, path=path.lstrip('/'), sha=self.root)
            r = requests.get(url, timeout=self.timeout, **self.kw)
            if r.status_code == 404:
                raise FileNotFoundError(path)
            r.raise_for_status()
            sha = r.json()['sha']
        delete_url = self.content_url.format(org=self.org, repo=self.repo, path=path, sha=self.root)
        branch = self.root
        data = {'message': message or f'Delete {path}', 'sha': sha, **({'branch': branch} if branch else {})}
        r = requests.delete(delete_url, json=data, timeout=self.timeout, **self.kw)
        error_message = r.json().get('message', '')
        if re.search('Branch .+ not found', error_message):
            error = 'Remove only works when the filesystem is initialised from a branch or default (None)'
            raise ValueError(error)
        r.raise_for_status()
        self.invalidate_cache(path)

    def _get_sha_from_cache(self, path):
        for entries in self.dircache.values():
            for entry in entries:
                entry_path = entry.get('name')
                if entry_path and entry_path == path and ('sha' in entry):
                    return entry['sha']
        return None
