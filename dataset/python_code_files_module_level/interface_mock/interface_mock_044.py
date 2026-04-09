from __future__ import annotations
import inspect
import logging
import os
import shutil
import uuid
from .asyn import AsyncFileSystem, _run_coros_in_chunks, sync_wrapper
from .callbacks import DEFAULT_CALLBACK
from .core import filesystem, get_filesystem_class, split_protocol, url_to_fs
_generic_fs = {}
logger = logging.getLogger('fsspec.generic')

def set_generic_fs(protocol, **storage_options):
    _generic_fs[protocol] = filesystem(protocol, **storage_options)

def _resolve_fs(url, method, protocol=None, storage_options=None):
    url = url[0] if isinstance(url, (list, tuple)) else url
    protocol = protocol or split_protocol(url)[0]
    storage_options = storage_options or {}
    if method == 'default':
        return filesystem(protocol)
    if method == 'generic':
        return _generic_fs[protocol]
    if method == 'current':
        cls = get_filesystem_class(protocol)
        return cls.current()
    if method == 'options':
        fs, _ = url_to_fs(url, **storage_options.get(protocol, {}))
        return fs
    raise ValueError(f'Unknown FS resolution method: {method}')

def rsync(source, destination, delete_missing=False, source_field='size', dest_field='size', update_cond='different', inst_kwargs=None, fs=None, **kwargs):
    fs = fs or GenericFileSystem(**inst_kwargs or {})
    source = fs._strip_protocol(source)
    destination = fs._strip_protocol(destination)
    allfiles = fs.find(source, withdirs=True, detail=True)
    if not fs.isdir(source):
        raise ValueError('Can only rsync on a directory')
    otherfiles = fs.find(destination, withdirs=True, detail=True)
    dirs = [a for a, v in allfiles.items() if v['type'] == 'directory' and a.replace(source, destination) not in otherfiles]
    logger.debug(f'{len(dirs)} directories to create')
    if dirs:
        fs.make_many_dirs([dirn.replace(source, destination) for dirn in dirs], exist_ok=True)
    allfiles = {a: v for a, v in allfiles.items() if v['type'] == 'file'}
    logger.debug(f'{len(allfiles)} files to consider for copy')
    to_delete = [o for o, v in otherfiles.items() if o.replace(destination, source) not in allfiles and v['type'] == 'file']
    for k, v in allfiles.copy().items():
        otherfile = k.replace(source, destination)
        if otherfile in otherfiles:
            if update_cond == 'always':
                allfiles[k] = otherfile
            elif update_cond == 'never':
                allfiles.pop(k)
            elif update_cond == 'different':
                inf1 = source_field(v) if callable(source_field) else v[source_field]
                v2 = otherfiles[otherfile]
                inf2 = dest_field(v2) if callable(dest_field) else v2[dest_field]
                if inf1 != inf2:
                    allfiles[k] = otherfile
                else:
                    allfiles.pop(k)
        else:
            allfiles[k] = otherfile
    logger.debug(f'{len(allfiles)} files to copy')
    if allfiles:
        source_files, target_files = zip(*allfiles.items())
        fs.cp(source_files, target_files, **kwargs)
    logger.debug(f'{len(to_delete)} files to delete')
    if delete_missing and to_delete:
        fs.rm(to_delete)
    return allfiles

class GenericFileSystem(AsyncFileSystem):
    protocol = 'generic'

    def __init__(self, default_method='default', storage_options=None, **kwargs):
        self.method = default_method
        self.st_opts = storage_options
        super().__init__(**kwargs)

    def _parent(self, path):
        fs = _resolve_fs(path, self.method, storage_options=self.st_opts)
        return fs.unstrip_protocol(fs._parent(path))

    def _strip_protocol(self, path):
        fs = _resolve_fs(path, self.method, storage_options=self.st_opts)
        return fs.unstrip_protocol(fs._strip_protocol(path))

    async def _find(self, path, maxdepth=None, withdirs=False, detail=False, **kwargs):
        fs = _resolve_fs(path, self.method, storage_options=self.st_opts)
        if fs.async_impl:
            out = await fs._find(path, maxdepth=maxdepth, withdirs=withdirs, detail=True, **kwargs)
        else:
            out = fs.find(path, maxdepth=maxdepth, withdirs=withdirs, detail=True, **kwargs)
        result = {}
        for k, v in out.items():
            v = v.copy()
            name = fs.unstrip_protocol(k)
            v['name'] = name
            result[name] = v
        if detail:
            return result
        return list(result)

    async def _info(self, url, **kwargs):
        fs = _resolve_fs(url, self.method)
        if fs.async_impl:
            out = await fs._info(url, **kwargs)
        else:
            out = fs.info(url, **kwargs)
        out = out.copy()
        out['name'] = fs.unstrip_protocol(out['name'])
        return out

    async def _ls(self, url, detail=True, **kwargs):
        fs = _resolve_fs(url, self.method)
        if fs.async_impl:
            out = await fs._ls(url, detail=True, **kwargs)
        else:
            out = fs.ls(url, detail=True, **kwargs)
        out = [o.copy() for o in out]
        for o in out:
            o['name'] = fs.unstrip_protocol(o['name'])
        if detail:
            return out
        else:
            return [o['name'] for o in out]

    async def _cat_file(self, url, **kwargs):
        fs = _resolve_fs(url, self.method)
        if fs.async_impl:
            return await fs._cat_file(url, **kwargs)
        else:
            return fs.cat_file(url, **kwargs)

    async def _pipe_file(self, path, value, **kwargs):
        fs = _resolve_fs(path, self.method, storage_options=self.st_opts)
        if fs.async_impl:
            return await fs._pipe_file(path, value, **kwargs)
        else:
            return fs.pipe_file(path, value, **kwargs)

    async def _rm(self, url, **kwargs):
        urls = url
        if isinstance(urls, str):
            urls = [urls]
        fs = _resolve_fs(urls[0], self.method)
        if fs.async_impl:
            await fs._rm(urls, **kwargs)
        else:
            fs.rm(url, **kwargs)

    async def _makedirs(self, path, exist_ok=False):
        logger.debug('Make dir %s', path)
        fs = _resolve_fs(path, self.method, storage_options=self.st_opts)
        if fs.async_impl:
            await fs._makedirs(path, exist_ok=exist_ok)
        else:
            fs.makedirs(path, exist_ok=exist_ok)

    def rsync(self, source, destination, **kwargs):
        rsync(source, destination, fs=self, **kwargs)

    async def _cp_file(self, url, url2, blocksize=2 ** 20, callback=DEFAULT_CALLBACK, tempdir: str | None=None, **kwargs):
        fs = _resolve_fs(url, self.method)
        fs2 = _resolve_fs(url2, self.method)
        if fs is fs2:
            if fs.async_impl:
                return await fs._copy(url, url2, **kwargs)
            else:
                return fs.copy(url, url2, **kwargs)
        await copy_file_op(fs, [url], fs2, [url2], tempdir, 1, on_error='raise')

    async def _make_many_dirs(self, urls, exist_ok=True):
        fs = _resolve_fs(urls[0], self.method)
        if fs.async_impl:
            coros = [fs._makedirs(u, exist_ok=exist_ok) for u in urls]
            await _run_coros_in_chunks(coros)
        else:
            for u in urls:
                fs.makedirs(u, exist_ok=exist_ok)
    make_many_dirs = sync_wrapper(_make_many_dirs)

    async def _copy(self, path1: list[str], path2: list[str], recursive: bool=False, on_error: str='ignore', maxdepth: int | None=None, batch_size: int | None=None, tempdir: str | None=None, **kwargs):
        if recursive:
            raise NotImplementedError('Please use fsspec.generic.rsync')
        path1 = [path1] if isinstance(path1, str) else path1
        path2 = [path2] if isinstance(path2, str) else path2
        fs = _resolve_fs(path1, self.method)
        fs2 = _resolve_fs(path2, self.method)
        if fs is fs2:
            if fs.async_impl:
                return await fs._copy(path1, path2, **kwargs)
            else:
                return fs.copy(path1, path2, **kwargs)
        await copy_file_op(fs, path1, fs2, path2, tempdir, batch_size, on_error=on_error)

async def copy_file_op(fs1, url1, fs2, url2, tempdir=None, batch_size=20, on_error='ignore'):
    import tempfile
    tempdir = tempdir or tempfile.mkdtemp()
    try:
        coros = [_copy_file_op(fs1, u1, fs2, u2, os.path.join(tempdir, uuid.uuid4().hex)) for u1, u2 in zip(url1, url2)]
        out = await _run_coros_in_chunks(coros, batch_size=batch_size, return_exceptions=True)
    finally:
        shutil.rmtree(tempdir)
    if on_error == 'return':
        return out
    elif on_error == 'raise':
        for o in out:
            if isinstance(o, Exception):
                raise o

async def _copy_file_op(fs1, url1, fs2, url2, local, on_error='ignore'):
    if fs1.async_impl:
        await fs1._get_file(url1, local)
    else:
        fs1.get_file(url1, local)
    if fs2.async_impl:
        await fs2._put_file(local, url2)
    else:
        fs2.put_file(local, url2)
    os.unlink(local)
    logger.debug('Copy %s -> %s; done', url1, url2)

async def maybe_await(cor):
    if inspect.iscoroutine(cor):
        return await cor
    else:
        return cor
