from __future__ import annotations
import OpenSSL.SSL
from cryptography import x509
try:
    from cryptography.x509 import UnsupportedExtension
except ImportError:

    class UnsupportedExtension(Exception):
        pass
import logging
import ssl
import typing
from io import BytesIO
from socket import socket as socket_cls
from socket import timeout
from .. import util
if typing.TYPE_CHECKING:
    from OpenSSL.crypto import X509
__all__ = ['inject_into_urllib3', 'extract_from_urllib3']
_openssl_versions: dict[int, int] = {util.ssl_.PROTOCOL_TLS: OpenSSL.SSL.SSLv23_METHOD, util.ssl_.PROTOCOL_TLS_CLIENT: OpenSSL.SSL.SSLv23_METHOD, ssl.PROTOCOL_TLSv1: OpenSSL.SSL.TLSv1_METHOD}
if hasattr(ssl, 'PROTOCOL_TLSv1_1') and hasattr(OpenSSL.SSL, 'TLSv1_1_METHOD'):
    _openssl_versions[ssl.PROTOCOL_TLSv1_1] = OpenSSL.SSL.TLSv1_1_METHOD
if hasattr(ssl, 'PROTOCOL_TLSv1_2') and hasattr(OpenSSL.SSL, 'TLSv1_2_METHOD'):
    _openssl_versions[ssl.PROTOCOL_TLSv1_2] = OpenSSL.SSL.TLSv1_2_METHOD
_stdlib_to_openssl_verify = {ssl.CERT_NONE: OpenSSL.SSL.VERIFY_NONE, ssl.CERT_OPTIONAL: OpenSSL.SSL.VERIFY_PEER, ssl.CERT_REQUIRED: OpenSSL.SSL.VERIFY_PEER + OpenSSL.SSL.VERIFY_FAIL_IF_NO_PEER_CERT}
_openssl_to_stdlib_verify = {v: k for k, v in _stdlib_to_openssl_verify.items()}
_OP_NO_SSLv2_OR_SSLv3: int = getattr(OpenSSL.SSL, 'OP_NO_SSLv2', 0) | getattr(OpenSSL.SSL, 'OP_NO_SSLv3', 0)
_OP_NO_TLSv1: int = getattr(OpenSSL.SSL, 'OP_NO_TLSv1', 0)
_OP_NO_TLSv1_1: int = getattr(OpenSSL.SSL, 'OP_NO_TLSv1_1', 0)
_OP_NO_TLSv1_2: int = getattr(OpenSSL.SSL, 'OP_NO_TLSv1_2', 0)
_OP_NO_TLSv1_3: int = getattr(OpenSSL.SSL, 'OP_NO_TLSv1_3', 0)
_openssl_to_ssl_minimum_version: dict[int, int] = {ssl.TLSVersion.MINIMUM_SUPPORTED: _OP_NO_SSLv2_OR_SSLv3, ssl.TLSVersion.TLSv1: _OP_NO_SSLv2_OR_SSLv3, ssl.TLSVersion.TLSv1_1: _OP_NO_SSLv2_OR_SSLv3 | _OP_NO_TLSv1, ssl.TLSVersion.TLSv1_2: _OP_NO_SSLv2_OR_SSLv3 | _OP_NO_TLSv1 | _OP_NO_TLSv1_1, ssl.TLSVersion.TLSv1_3: _OP_NO_SSLv2_OR_SSLv3 | _OP_NO_TLSv1 | _OP_NO_TLSv1_1 | _OP_NO_TLSv1_2, ssl.TLSVersion.MAXIMUM_SUPPORTED: _OP_NO_SSLv2_OR_SSLv3 | _OP_NO_TLSv1 | _OP_NO_TLSv1_1 | _OP_NO_TLSv1_2}
_openssl_to_ssl_maximum_version: dict[int, int] = {ssl.TLSVersion.MINIMUM_SUPPORTED: _OP_NO_SSLv2_OR_SSLv3 | _OP_NO_TLSv1 | _OP_NO_TLSv1_1 | _OP_NO_TLSv1_2 | _OP_NO_TLSv1_3, ssl.TLSVersion.TLSv1: _OP_NO_SSLv2_OR_SSLv3 | _OP_NO_TLSv1_1 | _OP_NO_TLSv1_2 | _OP_NO_TLSv1_3, ssl.TLSVersion.TLSv1_1: _OP_NO_SSLv2_OR_SSLv3 | _OP_NO_TLSv1_2 | _OP_NO_TLSv1_3, ssl.TLSVersion.TLSv1_2: _OP_NO_SSLv2_OR_SSLv3 | _OP_NO_TLSv1_3, ssl.TLSVersion.TLSv1_3: _OP_NO_SSLv2_OR_SSLv3, ssl.TLSVersion.MAXIMUM_SUPPORTED: _OP_NO_SSLv2_OR_SSLv3}
SSL_WRITE_BLOCKSIZE = 16384
orig_util_SSLContext = util.ssl_.SSLContext
log = logging.getLogger(__name__)

def inject_into_urllib3() -> None:
    _validate_dependencies_met()
    util.SSLContext = PyOpenSSLContext
    util.ssl_.SSLContext = PyOpenSSLContext
    util.IS_PYOPENSSL = True
    util.ssl_.IS_PYOPENSSL = True

def extract_from_urllib3() -> None:
    util.SSLContext = orig_util_SSLContext
    util.ssl_.SSLContext = orig_util_SSLContext
    util.IS_PYOPENSSL = False
    util.ssl_.IS_PYOPENSSL = False

def _validate_dependencies_met() -> None:
    from cryptography.x509.extensions import Extensions
    if getattr(Extensions, 'get_extension_for_class', None) is None:
        raise ImportError("'cryptography' module missing required functionality.  Try upgrading to v1.3.4 or newer.")
    from OpenSSL.crypto import X509
    x509 = X509()
    if getattr(x509, '_x509', None) is None:
        raise ImportError("'pyOpenSSL' module missing required functionality. Try upgrading to v0.14 or newer.")

def _dnsname_to_stdlib(name: str) -> str | None:

    def idna_encode(name: str) -> bytes | None:
        import idna
        try:
            for prefix in ['*.', '.']:
                if name.startswith(prefix):
                    name = name[len(prefix):]
                    return prefix.encode('ascii') + idna.encode(name)
            return idna.encode(name)
        except idna.core.IDNAError:
            return None
    if ':' in name:
        return name
    encoded_name = idna_encode(name)
    if encoded_name is None:
        return None
    return encoded_name.decode('utf-8')

def get_subj_alt_name(peer_cert: X509) -> list[tuple[str, str]]:
    cert = peer_cert.to_cryptography()
    try:
        ext = cert.extensions.get_extension_for_class(x509.SubjectAlternativeName).value
    except x509.ExtensionNotFound:
        return []
    except (x509.DuplicateExtension, UnsupportedExtension, x509.UnsupportedGeneralNameType, UnicodeError) as e:
        log.warning('A problem was encountered with the certificate that prevented urllib3 from finding the SubjectAlternativeName field. This can affect certificate validation. The error was %s', e)
        return []
    names = [('DNS', name) for name in map(_dnsname_to_stdlib, ext.get_values_for_type(x509.DNSName)) if name is not None]
    names.extend((('IP Address', str(name)) for name in ext.get_values_for_type(x509.IPAddress)))
    return names

class WrappedSocket:

    def __init__(self, connection: OpenSSL.SSL.Connection, socket: socket_cls, suppress_ragged_eofs: bool=True) -> None:
        self.connection = connection
        self.socket = socket
        self.suppress_ragged_eofs = suppress_ragged_eofs
        self._io_refs = 0
        self._closed = False

    def fileno(self) -> int:
        return self.socket.fileno()

    def _decref_socketios(self) -> None:
        if self._io_refs > 0:
            self._io_refs -= 1
        if self._closed:
            self.close()

    def recv(self, *args: typing.Any, **kwargs: typing.Any) -> bytes:
        try:
            data = self.connection.recv(*args, **kwargs)
        except OpenSSL.SSL.SysCallError as e:
            if self.suppress_ragged_eofs and e.args == (-1, 'Unexpected EOF'):
                return b''
            else:
                raise OSError(e.args[0], str(e)) from e
        except OpenSSL.SSL.ZeroReturnError:
            if self.connection.get_shutdown() == OpenSSL.SSL.RECEIVED_SHUTDOWN:
                return b''
            else:
                raise
        except OpenSSL.SSL.WantReadError as e:
            if not util.wait_for_read(self.socket, self.socket.gettimeout()):
                raise timeout('The read operation timed out') from e
            else:
                return self.recv(*args, **kwargs)
        except OpenSSL.SSL.Error as e:
            raise ssl.SSLError(f'read error: {e!r}') from e
        else:
            return data

    def recv_into(self, *args: typing.Any, **kwargs: typing.Any) -> int:
        try:
            return self.connection.recv_into(*args, **kwargs)
        except OpenSSL.SSL.SysCallError as e:
            if self.suppress_ragged_eofs and e.args == (-1, 'Unexpected EOF'):
                return 0
            else:
                raise OSError(e.args[0], str(e)) from e
        except OpenSSL.SSL.ZeroReturnError:
            if self.connection.get_shutdown() == OpenSSL.SSL.RECEIVED_SHUTDOWN:
                return 0
            else:
                raise
        except OpenSSL.SSL.WantReadError as e:
            if not util.wait_for_read(self.socket, self.socket.gettimeout()):
                raise timeout('The read operation timed out') from e
            else:
                return self.recv_into(*args, **kwargs)
        except OpenSSL.SSL.Error as e:
            raise ssl.SSLError(f'read error: {e!r}') from e

    def settimeout(self, timeout: float) -> None:
        return self.socket.settimeout(timeout)

    def _send_until_done(self, data: bytes) -> int:
        while True:
            try:
                return self.connection.send(data)
            except OpenSSL.SSL.WantWriteError as e:
                if not util.wait_for_write(self.socket, self.socket.gettimeout()):
                    raise timeout() from e
                continue
            except OpenSSL.SSL.SysCallError as e:
                raise OSError(e.args[0], str(e)) from e

    def sendall(self, data: bytes) -> None:
        total_sent = 0
        while total_sent < len(data):
            sent = self._send_until_done(data[total_sent:total_sent + SSL_WRITE_BLOCKSIZE])
            total_sent += sent

    def shutdown(self, how: int) -> None:
        try:
            self.connection.shutdown()
        except OpenSSL.SSL.Error as e:
            raise ssl.SSLError(f'shutdown error: {e!r}') from e

    def close(self) -> None:
        self._closed = True
        if self._io_refs <= 0:
            self._real_close()

    def _real_close(self) -> None:
        try:
            return self.connection.close()
        except OpenSSL.SSL.Error:
            return

    def getpeercert(self, binary_form: bool=False) -> dict[str, list[typing.Any]] | None:
        x509 = self.connection.get_peer_certificate()
        if not x509:
            return x509
        if binary_form:
            return OpenSSL.crypto.dump_certificate(OpenSSL.crypto.FILETYPE_ASN1, x509)
        return {'subject': ((('commonName', x509.get_subject().CN),),), 'subjectAltName': get_subj_alt_name(x509)}

    def version(self) -> str:
        return self.connection.get_protocol_version_name()

    def selected_alpn_protocol(self) -> str | None:
        alpn_proto = self.connection.get_alpn_proto_negotiated()
        return alpn_proto.decode() if alpn_proto else None
WrappedSocket.makefile = socket_cls.makefile

class PyOpenSSLContext:

    def __init__(self, protocol: int) -> None:
        self.protocol = _openssl_versions[protocol]
        self._ctx = OpenSSL.SSL.Context(self.protocol)
        self._options = 0
        self.check_hostname = False
        self._minimum_version: int = ssl.TLSVersion.MINIMUM_SUPPORTED
        self._maximum_version: int = ssl.TLSVersion.MAXIMUM_SUPPORTED
        self._verify_flags: int = ssl.VERIFY_X509_TRUSTED_FIRST

    @property
    def options(self) -> int:
        return self._options

    @options.setter
    def options(self, value: int) -> None:
        self._options = value
        self._set_ctx_options()

    @property
    def verify_flags(self) -> int:
        return self._verify_flags

    @verify_flags.setter
    def verify_flags(self, value: int) -> None:
        self._verify_flags = value
        self._ctx.get_cert_store().set_flags(self._verify_flags)

    @property
    def verify_mode(self) -> int:
        return _openssl_to_stdlib_verify[self._ctx.get_verify_mode()]

    @verify_mode.setter
    def verify_mode(self, value: ssl.VerifyMode) -> None:
        self._ctx.set_verify(_stdlib_to_openssl_verify[value], _verify_callback)

    def set_default_verify_paths(self) -> None:
        self._ctx.set_default_verify_paths()

    def set_ciphers(self, ciphers: bytes | str) -> None:
        if isinstance(ciphers, str):
            ciphers = ciphers.encode('utf-8')
        self._ctx.set_cipher_list(ciphers)

    def load_verify_locations(self, cafile: str | None=None, capath: str | None=None, cadata: bytes | None=None) -> None:
        if cafile is not None:
            cafile = cafile.encode('utf-8')
        if capath is not None:
            capath = capath.encode('utf-8')
        try:
            self._ctx.load_verify_locations(cafile, capath)
            if cadata is not None:
                self._ctx.load_verify_locations(BytesIO(cadata))
        except OpenSSL.SSL.Error as e:
            raise ssl.SSLError(f'unable to load trusted certificates: {e!r}') from e

    def load_cert_chain(self, certfile: str, keyfile: str | None=None, password: str | None=None) -> None:
        try:
            self._ctx.use_certificate_chain_file(certfile)
            if password is not None:
                if not isinstance(password, bytes):
                    password = password.encode('utf-8')
                self._ctx.set_passwd_cb(lambda *_: password)
            self._ctx.use_privatekey_file(keyfile or certfile)
        except OpenSSL.SSL.Error as e:
            raise ssl.SSLError(f'Unable to load certificate chain: {e!r}') from e

    def set_alpn_protocols(self, protocols: list[bytes | str]) -> None:
        protocols = [util.util.to_bytes(p, 'ascii') for p in protocols]
        return self._ctx.set_alpn_protos(protocols)

    def wrap_socket(self, sock: socket_cls, server_side: bool=False, do_handshake_on_connect: bool=True, suppress_ragged_eofs: bool=True, server_hostname: bytes | str | None=None) -> WrappedSocket:
        cnx = OpenSSL.SSL.Connection(self._ctx, sock)
        if server_hostname and (not util.ssl_.is_ipaddress(server_hostname)):
            if isinstance(server_hostname, str):
                server_hostname = server_hostname.encode('utf-8')
            cnx.set_tlsext_host_name(server_hostname)
        cnx.set_connect_state()
        while True:
            try:
                cnx.do_handshake()
            except OpenSSL.SSL.WantReadError as e:
                if not util.wait_for_read(sock, sock.gettimeout()):
                    raise timeout('select timed out') from e
                continue
            except OpenSSL.SSL.Error as e:
                raise ssl.SSLError(f'bad handshake: {e!r}') from e
            break
        return WrappedSocket(cnx, sock)

    def _set_ctx_options(self) -> None:
        self._ctx.set_options(self._options | _openssl_to_ssl_minimum_version[self._minimum_version] | _openssl_to_ssl_maximum_version[self._maximum_version])

    @property
    def minimum_version(self) -> int:
        return self._minimum_version

    @minimum_version.setter
    def minimum_version(self, minimum_version: int) -> None:
        self._minimum_version = minimum_version
        self._set_ctx_options()

    @property
    def maximum_version(self) -> int:
        return self._maximum_version

    @maximum_version.setter
    def maximum_version(self, maximum_version: int) -> None:
        self._maximum_version = maximum_version
        self._set_ctx_options()

def _verify_callback(cnx: OpenSSL.SSL.Connection, x509: X509, err_no: int, err_depth: int, return_code: int) -> bool:
    return err_no == 0
