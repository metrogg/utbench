import base64
import os
import signal
import sys
import threading
import warnings
from collections import deque
import json
from . import spawn
from . import util
__all__ = ['ensure_running', 'register', 'unregister']
_HAVE_SIGMASK = hasattr(signal, 'pthread_sigmask')
_IGNORED_SIGNALS = (signal.SIGINT, signal.SIGTERM)

def cleanup_noop(name):
    raise RuntimeError('noop should never be registered or cleaned up')
_CLEANUP_FUNCS = {'noop': cleanup_noop, 'dummy': lambda name: None}
if os.name == 'posix':
    try:
        import _multiprocess as _multiprocessing
    except ImportError:
        import _multiprocessing
    import _posixshmem
    if hasattr(_multiprocessing, 'sem_unlink'):
        _CLEANUP_FUNCS['semaphore'] = _multiprocessing.sem_unlink
    _CLEANUP_FUNCS['shared_memory'] = _posixshmem.shm_unlink

class ReentrantCallError(RuntimeError):
    pass

class ResourceTracker(object):

    def __init__(self):
        self._lock = threading.RLock()
        self._fd = None
        self._pid = None
        self._exitcode = None
        self._reentrant_messages = deque()
        self._use_simple_format = True

    def _reentrant_call_error(self):
        raise ReentrantCallError('Reentrant call into the multiprocess resource tracker')

    def __del__(self):
        self._stop(use_blocking_lock=False)

    def _stop(self, use_blocking_lock=True):
        if use_blocking_lock:
            with self._lock:
                self._stop_locked()
        else:
            acquired = self._lock.acquire(blocking=False)
            try:
                self._stop_locked()
            finally:
                if acquired:
                    self._lock.release()

    def _stop_locked(self, close=os.close, waitpid=os.waitpid, waitstatus_to_exitcode=os.waitstatus_to_exitcode):
        if self._lock._recursion_count() > 1:
            raise self._reentrant_call_error()
        if self._fd is None:
            return
        if self._pid is None:
            return
        close(self._fd)
        self._fd = None
        try:
            _, status = waitpid(self._pid, 0)
        except ChildProcessError:
            self._pid = None
            self._exitcode = None
            return
        self._pid = None
        try:
            self._exitcode = waitstatus_to_exitcode(status)
        except ValueError:
            self._exitcode = None

    def getfd(self):
        self.ensure_running()
        return self._fd

    def ensure_running(self):
        return self._ensure_running_and_write()

    def _teardown_dead_process(self):
        os.close(self._fd)
        try:
            if self._pid is not None:
                os.waitpid(self._pid, 0)
        except ChildProcessError:
            pass
        self._fd = None
        self._pid = None
        self._exitcode = None
        warnings.warn('resource_tracker: process died unexpectedly, relaunching.  Some resources might leak.')

    def _launch(self):
        fds_to_pass = []
        try:
            fds_to_pass.append(sys.stderr.fileno())
        except Exception:
            pass
        r, w = os.pipe()
        try:
            fds_to_pass.append(r)
            exe = spawn.get_executable()
            args = [exe, *util._args_from_interpreter_flags(), '-c', f'from multiprocess.resource_tracker import main;main({r})']
            prev_sigmask = None
            try:
                if _HAVE_SIGMASK:
                    prev_sigmask = signal.pthread_sigmask(signal.SIG_BLOCK, _IGNORED_SIGNALS)
                pid = util.spawnv_passfds(exe, args, fds_to_pass)
            finally:
                if prev_sigmask is not None:
                    signal.pthread_sigmask(signal.SIG_SETMASK, prev_sigmask)
        except:
            os.close(w)
            raise
        else:
            self._fd = w
            self._pid = pid
        finally:
            os.close(r)

    def _make_probe_message(self):
        if self._use_simple_format:
            return b'PROBE:0:noop\n'
        return (json.dumps({'cmd': 'PROBE', 'rtype': 'noop'}, ensure_ascii=True, separators=(',', ':')) + '\n').encode('ascii')

    def _ensure_running_and_write(self, msg=None):
        with self._lock:
            if self._lock._recursion_count() > 1:
                if msg is None:
                    raise self._reentrant_call_error()
                return self._reentrant_messages.append(msg)
            if self._fd is not None:
                if msg is None:
                    to_send = self._make_probe_message()
                else:
                    to_send = msg
                try:
                    self._write(to_send)
                except OSError:
                    self._teardown_dead_process()
                    self._launch()
                msg = None
            else:
                self._launch()
        while True:
            try:
                reentrant_msg = self._reentrant_messages.popleft()
            except IndexError:
                break
            self._write(reentrant_msg)
        if msg is not None:
            self._write(msg)

    def _check_alive(self):
        try:
            os.write(self._fd, self._make_probe_message())
        except OSError:
            return False
        else:
            return True

    def register(self, name, rtype):
        self._send('REGISTER', name, rtype)

    def unregister(self, name, rtype):
        self._send('UNREGISTER', name, rtype)

    def _write(self, msg):
        nbytes = os.write(self._fd, msg)
        assert nbytes == len(msg), f'nbytes={nbytes!r} != len(msg)={len(msg)!r}'

    def _send(self, cmd, name, rtype):
        if self._use_simple_format and '\n' not in name:
            msg = f'{cmd}:{name}:{rtype}\n'.encode('ascii')
            if len(msg) > 512:
                raise ValueError('msg too long')
            self._ensure_running_and_write(msg)
            return
        b = name.encode('utf-8', 'surrogateescape')
        if len(b) > 255:
            raise ValueError('shared memory name too long (max 255 bytes)')
        b64 = base64.urlsafe_b64encode(b).decode('ascii')
        payload = {'cmd': cmd, 'rtype': rtype, 'base64_name': b64}
        msg = (json.dumps(payload, ensure_ascii=True, separators=(',', ':')) + '\n').encode('ascii')
        assert len(msg) <= 512, f'internal error: message too long ({len(msg)} bytes)'
        assert msg.startswith(b'{')
        self._ensure_running_and_write(msg)
_resource_tracker = ResourceTracker()
ensure_running = _resource_tracker.ensure_running
register = _resource_tracker.register
unregister = _resource_tracker.unregister
getfd = _resource_tracker.getfd

def _decode_message(line):
    if line.startswith(b'{'):
        try:
            obj = json.loads(line.decode('ascii'))
        except Exception as e:
            raise ValueError('malformed resource_tracker message: %r' % (line,)) from e
        cmd = obj['cmd']
        rtype = obj['rtype']
        b64 = obj.get('base64_name', '')
        if not isinstance(cmd, str) or not isinstance(rtype, str) or (not isinstance(b64, str)):
            raise ValueError('malformed resource_tracker fields: %r' % (obj,))
        try:
            name = base64.urlsafe_b64decode(b64).decode('utf-8', 'surrogateescape')
        except ValueError as e:
            raise ValueError('malformed resource_tracker base64_name: %r' % (b64,)) from e
    else:
        cmd, rest = line.strip().decode('ascii').split(':', maxsplit=1)
        name, rtype = rest.rsplit(':', maxsplit=1)
    return (cmd, rtype, name)

def main(fd):
    signal.signal(signal.SIGINT, signal.SIG_IGN)
    signal.signal(signal.SIGTERM, signal.SIG_IGN)
    if _HAVE_SIGMASK:
        signal.pthread_sigmask(signal.SIG_UNBLOCK, _IGNORED_SIGNALS)
    for f in (sys.stdin, sys.stdout):
        try:
            f.close()
        except Exception:
            pass
    cache = {rtype: set() for rtype in _CLEANUP_FUNCS.keys()}
    exit_code = 0
    try:
        with open(fd, 'rb') as f:
            for line in f:
                try:
                    cmd, rtype, name = _decode_message(line)
                    cleanup_func = _CLEANUP_FUNCS.get(rtype, None)
                    if cleanup_func is None:
                        raise ValueError(f'Cannot register {name} for automatic cleanup: unknown resource type {rtype}')
                    if cmd == 'REGISTER':
                        cache[rtype].add(name)
                    elif cmd == 'UNREGISTER':
                        cache[rtype].remove(name)
                    elif cmd == 'PROBE':
                        pass
                    else:
                        raise RuntimeError('unrecognized command %r' % cmd)
                except Exception:
                    exit_code = 3
                    try:
                        sys.excepthook(*sys.exc_info())
                    except:
                        pass
    finally:
        for rtype, rtype_cache in cache.items():
            if rtype_cache:
                try:
                    exit_code = 1
                    if rtype == 'dummy':
                        pass
                    else:
                        warnings.warn(f'resource_tracker: There appear to be {len(rtype_cache)} leaked {rtype} objects to clean up at shutdown: {rtype_cache}')
                except Exception:
                    pass
            for name in rtype_cache:
                try:
                    try:
                        _CLEANUP_FUNCS[rtype](name)
                    except Exception as e:
                        exit_code = 2
                        warnings.warn('resource_tracker: %r: %s' % (name, e))
                finally:
                    pass
        sys.exit(exit_code)
