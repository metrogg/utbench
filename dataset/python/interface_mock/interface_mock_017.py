from __future__ import annotations
import hashlib
import hmac
import os
import socket
import sys
import typing
import warnings
from binascii import unhexlify
from ..exceptions import ProxySchemeUnsupported, SSLError
from .url import _BRACELESS_IPV6_ADDRZ_RE, _IPV4_RE
SSLContext = None
SSLTransport = None
HAS_NEVER_CHECK_COMMON_NAME = False
IS_PYOPENSSL = False
ALPN_PROTOCOLS = ['http/1.1']
_TYPE_VERSION_INFO = tuple[int, int, int, str, int]
HASHFUNC_MAP = {length: getattr(hashlib, algorithm, None) for length, algorithm in ((32, 'md5'), (40, 'sha1'), (64, 'sha256'))}

def _is_bpo_43522_fixed(implementation_name: str, version_info: _TYPE_VERSION_INFO, pypy_version_info: _TYPE_VERSION_INFO | None) -> bool:
    if implementation_name == 'pypy':
        return pypy_version_info >= (7, 3, 8)
    elif implementation_name == 'cpython':
        major_minor = version_info[:2]
        micro = version_info[2]
        return major_minor == (3, 9) and micro >= 3 or major_minor >= (3, 10)
    else:
        return False

def _is_has_never_check_common_name_reliable(openssl_version: str, openssl_version_number: int, implementation_name: str, version_info: _TYPE_VERSION_INFO, pypy_version_info: _TYPE_VERSION_INFO | None) -> bool:
    is_openssl = openssl_version.startswith('OpenSSL ')
    is_openssl_issue_14579_fixed = openssl_version_number >= 269488335
    return is_openssl and (is_openssl_issue_14579_fixed or _is_bpo_43522_fixed(implementation_name, version_info, pypy_version_info))
if typing.TYPE_CHECKING:
    from ssl import VerifyMode
    from typing import TypedDict
    from .ssltransport import SSLTransport as SSLTransportType

    class _TYPE_PEER_CERT_RET_DICT(TypedDict, total=False):
        subjectAltName: tuple[tuple[str, str], ...]
        subject: tuple[tuple[tuple[str, str], ...], ...]
        serialNumber: str
_SSL_VERSION_TO_TLS_VERSION: dict[int, int] = {}
try:
    import ssl
    from ssl import CERT_REQUIRED, HAS_NEVER_CHECK_COMMON_NAME, OP_NO_COMPRESSION, OP_NO_TICKET, OPENSSL_VERSION, OPENSSL_VERSION_NUMBER, PROTOCOL_TLS, PROTOCOL_TLS_CLIENT, VERIFY_X509_STRICT, OP_NO_SSLv2, OP_NO_SSLv3, SSLContext, TLSVersion
    PROTOCOL_SSLv23 = PROTOCOL_TLS
    VERIFY_X509_PARTIAL_CHAIN = getattr(ssl, 'VERIFY_X509_PARTIAL_CHAIN', 524288)
    if HAS_NEVER_CHECK_COMMON_NAME and (not _is_has_never_check_common_name_reliable(OPENSSL_VERSION, OPENSSL_VERSION_NUMBER, sys.implementation.name, sys.version_info, sys.pypy_version_info if sys.implementation.name == 'pypy' else None)):
        HAS_NEVER_CHECK_COMMON_NAME = False
    for attr in ('TLSv1', 'TLSv1_1', 'TLSv1_2'):
        try:
            _SSL_VERSION_TO_TLS_VERSION[getattr(ssl, f'PROTOCOL_{attr}')] = getattr(TLSVersion, attr)
        except AttributeError:
            continue
    from .ssltransport import SSLTransport
except ImportError:
    OP_NO_COMPRESSION = 131072
    OP_NO_TICKET = 16384
    OP_NO_SSLv2 = 16777216
    OP_NO_SSLv3 = 33554432
    PROTOCOL_SSLv23 = PROTOCOL_TLS = 2
    PROTOCOL_TLS_CLIENT = 16
    VERIFY_X509_PARTIAL_CHAIN = 524288
    VERIFY_X509_STRICT = 32
_TYPE_PEER_CERT_RET = typing.Union['_TYPE_PEER_CERT_RET_DICT', bytes, None]

def assert_fingerprint(cert: bytes | None, fingerprint: str) -> None:
    if cert is None:
        raise SSLError('No certificate for the peer.')
    fingerprint = fingerprint.replace(':', '').lower()
    digest_length = len(fingerprint)
    if digest_length not in HASHFUNC_MAP:
        raise SSLError(f'Fingerprint of invalid length: {fingerprint}')
    hashfunc = HASHFUNC_MAP.get(digest_length)
    if hashfunc is None:
        raise SSLError(f'Hash function implementation unavailable for fingerprint length: {digest_length}')
    fingerprint_bytes = unhexlify(fingerprint.encode())
    cert_digest = hashfunc(cert).digest()
    if not hmac.compare_digest(cert_digest, fingerprint_bytes):
        raise SSLError(f'Fingerprints did not match. Expected "{fingerprint}", got "{cert_digest.hex()}"')

def resolve_cert_reqs(candidate: None | int | str) -> VerifyMode:
    if candidate is None:
        return CERT_REQUIRED
    if isinstance(candidate, str):
        res = getattr(ssl, candidate, None)
        if res is None:
            res = getattr(ssl, 'CERT_' + candidate)
        return res
    return candidate

def resolve_ssl_version(candidate: None | int | str) -> int:
    if candidate is None:
        return PROTOCOL_TLS
    if isinstance(candidate, str):
        res = getattr(ssl, candidate, None)
        if res is None:
            res = getattr(ssl, 'PROTOCOL_' + candidate)
        return typing.cast(int, res)
    return candidate

def create_urllib3_context(ssl_version: int | None=None, cert_reqs: int | None=None, options: int | None=None, ciphers: str | None=None, ssl_minimum_version: int | None=None, ssl_maximum_version: int | None=None, verify_flags: int | None=None) -> ssl.SSLContext:
    if SSLContext is None:
        raise TypeError("Can't create an SSLContext object without an ssl module")
    if ssl_version not in (None, PROTOCOL_TLS, PROTOCOL_TLS_CLIENT):
        if ssl_minimum_version is not None or ssl_maximum_version is not None:
            raise ValueError("Can't specify both 'ssl_version' and either 'ssl_minimum_version' or 'ssl_maximum_version'")
        else:
            ssl_minimum_version = _SSL_VERSION_TO_TLS_VERSION.get(ssl_version, TLSVersion.MINIMUM_SUPPORTED)
            ssl_maximum_version = _SSL_VERSION_TO_TLS_VERSION.get(ssl_version, TLSVersion.MAXIMUM_SUPPORTED)
            warnings.warn("'ssl_version' option is deprecated and will be removed in urllib3 v2.6.0. Instead use 'ssl_minimum_version'", category=DeprecationWarning, stacklevel=2)
    context = SSLContext(PROTOCOL_TLS_CLIENT)
    if ssl_minimum_version is not None:
        context.minimum_version = ssl_minimum_version
    else:
        context.minimum_version = TLSVersion.TLSv1_2
    if ssl_maximum_version is not None:
        context.maximum_version = ssl_maximum_version
    if ciphers:
        context.set_ciphers(ciphers)
    cert_reqs = ssl.CERT_REQUIRED if cert_reqs is None else cert_reqs
    if options is None:
        options = 0
        options |= OP_NO_SSLv2
        options |= OP_NO_SSLv3
        options |= OP_NO_COMPRESSION
        options |= OP_NO_TICKET
    context.options |= options
    if verify_flags is None:
        verify_flags = 0
        if sys.version_info >= (3, 13):
            verify_flags |= VERIFY_X509_PARTIAL_CHAIN
            verify_flags |= VERIFY_X509_STRICT
    context.verify_flags |= verify_flags
    if getattr(context, 'post_handshake_auth', None) is not None:
        context.post_handshake_auth = True
    if cert_reqs == ssl.CERT_REQUIRED and (not IS_PYOPENSSL):
        context.verify_mode = cert_reqs
        context.check_hostname = True
    else:
        context.check_hostname = False
        context.verify_mode = cert_reqs
    try:
        context.hostname_checks_common_name = False
    except AttributeError:
        pass
    if 'SSLKEYLOGFILE' in os.environ:
        sslkeylogfile = os.path.expandvars(os.environ.get('SSLKEYLOGFILE'))
    else:
        sslkeylogfile = None
    if sslkeylogfile:
        context.keylog_filename = sslkeylogfile
    return context

@typing.overload
def ssl_wrap_socket(sock: socket.socket, keyfile: str | None=..., certfile: str | None=..., cert_reqs: int | None=..., ca_certs: str | None=..., server_hostname: str | None=..., ssl_version: int | None=..., ciphers: str | None=..., ssl_context: ssl.SSLContext | None=..., ca_cert_dir: str | None=..., key_password: str | None=..., ca_cert_data: None | str | bytes=..., tls_in_tls: typing.Literal[False]=...) -> ssl.SSLSocket:
    ...

@typing.overload
def ssl_wrap_socket(sock: socket.socket, keyfile: str | None=..., certfile: str | None=..., cert_reqs: int | None=..., ca_certs: str | None=..., server_hostname: str | None=..., ssl_version: int | None=..., ciphers: str | None=..., ssl_context: ssl.SSLContext | None=..., ca_cert_dir: str | None=..., key_password: str | None=..., ca_cert_data: None | str | bytes=..., tls_in_tls: bool=...) -> ssl.SSLSocket | SSLTransportType:
    ...

def ssl_wrap_socket(sock: socket.socket, keyfile: str | None=None, certfile: str | None=None, cert_reqs: int | None=None, ca_certs: str | None=None, server_hostname: str | None=None, ssl_version: int | None=None, ciphers: str | None=None, ssl_context: ssl.SSLContext | None=None, ca_cert_dir: str | None=None, key_password: str | None=None, ca_cert_data: None | str | bytes=None, tls_in_tls: bool=False) -> ssl.SSLSocket | SSLTransportType:
    context = ssl_context
    if context is None:
        context = create_urllib3_context(ssl_version, cert_reqs, ciphers=ciphers)
    if ca_certs or ca_cert_dir or ca_cert_data:
        try:
            context.load_verify_locations(ca_certs, ca_cert_dir, ca_cert_data)
        except OSError as e:
            raise SSLError(e) from e
    elif ssl_context is None and hasattr(context, 'load_default_certs'):
        context.load_default_certs()
    if keyfile and key_password is None and _is_key_file_encrypted(keyfile):
        raise SSLError('Client private key is encrypted, password is required')
    if certfile:
        if key_password is None:
            context.load_cert_chain(certfile, keyfile)
        else:
            context.load_cert_chain(certfile, keyfile, key_password)
    context.set_alpn_protocols(ALPN_PROTOCOLS)
    ssl_sock = _ssl_wrap_socket_impl(sock, context, tls_in_tls, server_hostname)
    return ssl_sock

def is_ipaddress(hostname: str | bytes) -> bool:
    if isinstance(hostname, bytes):
        hostname = hostname.decode('ascii')
    return bool(_IPV4_RE.match(hostname) or _BRACELESS_IPV6_ADDRZ_RE.match(hostname))

def _is_key_file_encrypted(key_file: str) -> bool:
    with open(key_file) as f:
        for line in f:
            if 'ENCRYPTED' in line:
                return True
    return False

def _ssl_wrap_socket_impl(sock: socket.socket, ssl_context: ssl.SSLContext, tls_in_tls: bool, server_hostname: str | None=None) -> ssl.SSLSocket | SSLTransportType:
    if tls_in_tls:
        if not SSLTransport:
            raise ProxySchemeUnsupported("TLS in TLS requires support for the 'ssl' module")
        SSLTransport._validate_ssl_context_for_tls_in_tls(ssl_context)
        return SSLTransport(sock, ssl_context, server_hostname)
    return ssl_context.wrap_socket(sock, server_hostname=server_hostname)
