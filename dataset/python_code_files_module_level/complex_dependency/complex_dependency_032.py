from __future__ import annotations
import base64
import urllib
import requests
from requests.adapters import HTTPAdapter, Retry
from typing_extensions import override
from fsspec import AbstractFileSystem
from fsspec.spec import AbstractBufferedFile

class DatabricksException(Exception):

    def __init__(self, error_code, message, details=None):
        super().__init__(message)
        self.error_code = error_code
        self.message = message
        self.details = details

class DatabricksFileSystem(AbstractFileSystem):

    def __init__(self, instance, token, **kwargs):
        self.instance = instance
        self.token = token
        self.session = requests.Session()
        self.retries = Retry(total=10, backoff_factor=0.05, status_forcelist=[408, 429, 500, 502, 503, 504])
        self.session.mount('https://', HTTPAdapter(max_retries=self.retries))
        self.session.headers.update({'Authorization': f'Bearer {self.token}'})
        super().__init__(**kwargs)

    @override
    def _ls_from_cache(self, path) -> list[dict[str, str | int]] | None:
        self.dircache.pop(path.rstrip('/'), None)
        parent = self._parent(path)
        if parent in self.dircache:
            for entry in self.dircache[parent]:
                if entry['name'] == path.rstrip('/'):
                    if entry['type'] != 'directory':
                        return [entry]
                    return []
            raise FileNotFoundError(path)

    def ls(self, path, detail=True, **kwargs):
        try:
            out = self._ls_from_cache(path)
        except FileNotFoundError:
            self.dircache.pop(self._parent(path), None)
            out = None
        if not out:
            try:
                r = self._send_to_api(method='get', endpoint='list', json={'path': path})
            except DatabricksException as e:
                if e.error_code == 'RESOURCE_DOES_NOT_EXIST':
                    raise FileNotFoundError(e.message) from e
                raise
            files = r.get('files', [])
            out = [{'name': o['path'], 'type': 'directory' if o['is_dir'] else 'file', 'size': o['file_size']} for o in files]
            self.dircache[path] = out
        if detail:
            return out
        return [o['name'] for o in out]

    def makedirs(self, path, exist_ok=True):
        if not exist_ok:
            try:
                self._send_to_api(method='get', endpoint='get-status', json={'path': path})
                raise FileExistsError(f'Path {path} already exists')
            except DatabricksException as e:
                if e.error_code == 'RESOURCE_DOES_NOT_EXIST':
                    pass
        try:
            self._send_to_api(method='post', endpoint='mkdirs', json={'path': path})
        except DatabricksException as e:
            if e.error_code == 'RESOURCE_ALREADY_EXISTS':
                raise FileExistsError(e.message) from e
            raise
        self.invalidate_cache(self._parent(path))

    def mkdir(self, path, create_parents=True, **kwargs):
        if not create_parents:
            raise NotImplementedError
        self.mkdirs(path, **kwargs)

    def rm(self, path, recursive=False, **kwargs):
        try:
            self._send_to_api(method='post', endpoint='delete', json={'path': path, 'recursive': recursive})
        except DatabricksException as e:
            if e.error_code == 'PARTIAL_DELETE':
                self.rm(path=path, recursive=recursive)
            elif e.error_code == 'IO_ERROR':
                raise OSError(e.message) from e
            raise
        self.invalidate_cache(self._parent(path))

    def mv(self, source_path, destination_path, recursive=False, maxdepth=None, **kwargs):
        if recursive:
            raise NotImplementedError
        if maxdepth:
            raise NotImplementedError
        try:
            self._send_to_api(method='post', endpoint='move', json={'source_path': source_path, 'destination_path': destination_path})
        except DatabricksException as e:
            if e.error_code == 'RESOURCE_DOES_NOT_EXIST':
                raise FileNotFoundError(e.message) from e
            elif e.error_code == 'RESOURCE_ALREADY_EXISTS':
                raise FileExistsError(e.message) from e
            raise
        self.invalidate_cache(self._parent(source_path))
        self.invalidate_cache(self._parent(destination_path))

    def _open(self, path, mode='rb', block_size='default', **kwargs):
        return DatabricksFile(self, path, mode=mode, block_size=block_size, **kwargs)

    def _send_to_api(self, method, endpoint, json):
        if method == 'post':
            session_call = self.session.post
        elif method == 'get':
            session_call = self.session.get
        else:
            raise ValueError(f'Do not understand method {method}')
        url = urllib.parse.urljoin(f'https://{self.instance}/api/2.0/dbfs/', endpoint)
        r = session_call(url, json=json)
        try:
            r.raise_for_status()
        except requests.HTTPError as e:
            try:
                exception_json = e.response.json()
            except Exception:
                raise e from None
            raise DatabricksException(**exception_json) from e
        return r.json()

    def _create_handle(self, path, overwrite=True):
        try:
            r = self._send_to_api(method='post', endpoint='create', json={'path': path, 'overwrite': overwrite})
            return r['handle']
        except DatabricksException as e:
            if e.error_code == 'RESOURCE_ALREADY_EXISTS':
                raise FileExistsError(e.message) from e
            raise

    def _close_handle(self, handle):
        try:
            self._send_to_api(method='post', endpoint='close', json={'handle': handle})
        except DatabricksException as e:
            if e.error_code == 'RESOURCE_DOES_NOT_EXIST':
                raise FileNotFoundError(e.message) from e
            raise

    def _add_data(self, handle, data):
        data = base64.b64encode(data).decode()
        try:
            self._send_to_api(method='post', endpoint='add-block', json={'handle': handle, 'data': data})
        except DatabricksException as e:
            if e.error_code == 'RESOURCE_DOES_NOT_EXIST':
                raise FileNotFoundError(e.message) from e
            elif e.error_code == 'MAX_BLOCK_SIZE_EXCEEDED':
                raise ValueError(e.message) from e
            raise

    def _get_data(self, path, start, end):
        try:
            r = self._send_to_api(method='get', endpoint='read', json={'path': path, 'offset': start, 'length': end - start})
            return base64.b64decode(r['data'])
        except DatabricksException as e:
            if e.error_code == 'RESOURCE_DOES_NOT_EXIST':
                raise FileNotFoundError(e.message) from e
            elif e.error_code in ['INVALID_PARAMETER_VALUE', 'MAX_READ_SIZE_EXCEEDED']:
                raise ValueError(e.message) from e
            raise

    def invalidate_cache(self, path=None):
        if path is None:
            self.dircache.clear()
        else:
            self.dircache.pop(path, None)
        super().invalidate_cache(path)

class DatabricksFile(AbstractBufferedFile):
    DEFAULT_BLOCK_SIZE = 1 * 2 ** 20

    def __init__(self, fs, path, mode='rb', block_size='default', autocommit=True, cache_type='readahead', cache_options=None, **kwargs):
        if block_size is None or block_size == 'default':
            block_size = self.DEFAULT_BLOCK_SIZE
        assert block_size == self.DEFAULT_BLOCK_SIZE, f'Only the default block size is allowed, not {block_size}'
        super().__init__(fs, path, mode=mode, block_size=block_size, autocommit=autocommit, cache_type=cache_type, cache_options=cache_options or {}, **kwargs)

    def _initiate_upload(self):
        self.handle = self.fs._create_handle(self.path)

    def _upload_chunk(self, final=False):
        self.buffer.seek(0)
        data = self.buffer.getvalue()
        data_chunks = [data[start:end] for start, end in self._to_sized_blocks(len(data))]
        for data_chunk in data_chunks:
            self.fs._add_data(handle=self.handle, data=data_chunk)
        if final:
            self.fs._close_handle(handle=self.handle)
            return True

    def _fetch_range(self, start, end):
        return_buffer = b''
        length = end - start
        for chunk_start, chunk_end in self._to_sized_blocks(length, start):
            return_buffer += self.fs._get_data(path=self.path, start=chunk_start, end=chunk_end)
        return return_buffer

    def _to_sized_blocks(self, length, start=0):
        end = start + length
        for data_chunk in range(start, end, self.blocksize):
            data_start = data_chunk
            data_end = min(end, data_chunk + self.blocksize)
            yield (data_start, data_end)
