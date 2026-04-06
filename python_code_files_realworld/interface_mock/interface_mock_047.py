import hashlib
import os
import re
import sys
import time
from typing import Callable, Dict, Final, FrozenSet, List, Literal, Tuple, TypedDict, Union
from yarl import URL
from . import hdrs
from .client_exceptions import ClientError
from .client_middlewares import ClientHandlerType
from .client_reqrep import ClientRequest, ClientResponse
from .payload import Payload

class DigestAuthChallenge(TypedDict, total=False):
    realm: str
    nonce: str
    qop: str
    algorithm: str
    opaque: str
    domain: str
    stale: str
DigestFunctions: Dict[str, Callable[[bytes], 'hashlib._Hash']] = {'MD5': hashlib.md5, 'MD5-SESS': hashlib.md5, 'SHA': hashlib.sha1, 'SHA-SESS': hashlib.sha1, 'SHA256': hashlib.sha256, 'SHA256-SESS': hashlib.sha256, 'SHA-256': hashlib.sha256, 'SHA-256-SESS': hashlib.sha256, 'SHA512': hashlib.sha512, 'SHA512-SESS': hashlib.sha512, 'SHA-512': hashlib.sha512, 'SHA-512-SESS': hashlib.sha512}
_HEADER_PAIRS_PATTERN = re.compile('(?:^|\\s|,\\s*)(\\w+)\\s*=\\s*(?:"((?:[^"\\\\]|\\\\.)*)"|([^\\s,]+))' if sys.version_info < (3, 11) else '(?:^|\\s|,\\s*)((?>\\w+))\\s*=\\s*(?:"((?:[^"\\\\]|\\\\.)*)"|([^\\s,]+))')
CHALLENGE_FIELDS: Final[Tuple[Literal['realm', 'nonce', 'qop', 'algorithm', 'opaque', 'domain', 'stale'], ...]] = ('realm', 'nonce', 'qop', 'algorithm', 'opaque', 'domain', 'stale')
SUPPORTED_ALGORITHMS: Final[Tuple[str, ...]] = tuple(sorted(DigestFunctions.keys()))
QUOTED_AUTH_FIELDS: Final[FrozenSet[str]] = frozenset({'username', 'realm', 'nonce', 'uri', 'response', 'opaque', 'cnonce'})

def escape_quotes(value: str) -> str:
    return value.replace('"', '\\"')

def unescape_quotes(value: str) -> str:
    return value.replace('\\"', '"')

def parse_header_pairs(header: str) -> Dict[str, str]:
    return {stripped_key: unescape_quotes(quoted_val) if quoted_val else unquoted_val for key, quoted_val, unquoted_val in _HEADER_PAIRS_PATTERN.findall(header) if (stripped_key := key.strip())}

class DigestAuthMiddleware:

    def __init__(self, login: str, password: str, preemptive: bool=True) -> None:
        if login is None:
            raise ValueError('None is not allowed as login value')
        if password is None:
            raise ValueError('None is not allowed as password value')
        if ':' in login:
            raise ValueError('A ":" is not allowed in username (RFC 1945#section-11.1)')
        self._login_str: Final[str] = login
        self._login_bytes: Final[bytes] = login.encode('utf-8')
        self._password_bytes: Final[bytes] = password.encode('utf-8')
        self._last_nonce_bytes = b''
        self._nonce_count = 0
        self._challenge: DigestAuthChallenge = {}
        self._preemptive: bool = preemptive
        self._protection_space: List[str] = []

    async def _encode(self, method: str, url: URL, body: Union[Payload, Literal[b'']]) -> str:
        challenge = self._challenge
        if 'realm' not in challenge:
            raise ClientError("Malformed Digest auth challenge: Missing 'realm' parameter")
        if 'nonce' not in challenge:
            raise ClientError("Malformed Digest auth challenge: Missing 'nonce' parameter")
        realm = challenge['realm']
        nonce = challenge['nonce']
        if not nonce:
            raise ClientError("Security issue: Digest auth challenge contains empty 'nonce' value")
        qop_raw = challenge.get('qop', '')
        algorithm_original = challenge.get('algorithm', 'MD5')
        algorithm = algorithm_original.upper()
        opaque = challenge.get('opaque', '')
        nonce_bytes = nonce.encode('utf-8')
        realm_bytes = realm.encode('utf-8')
        path = URL(url).path_qs
        qop = ''
        qop_bytes = b''
        if qop_raw:
            valid_qops = {'auth', 'auth-int'}.intersection({q.strip() for q in qop_raw.split(',') if q.strip()})
            if not valid_qops:
                raise ClientError(f'Digest auth error: Unsupported Quality of Protection (qop) value(s): {qop_raw}')
            qop = 'auth-int' if 'auth-int' in valid_qops else 'auth'
            qop_bytes = qop.encode('utf-8')
        if algorithm not in DigestFunctions:
            raise ClientError(f"Digest auth error: Unsupported hash algorithm: {algorithm}. Supported algorithms: {', '.join(SUPPORTED_ALGORITHMS)}")
        hash_fn: Final = DigestFunctions[algorithm]

        def H(x: bytes) -> bytes:
            return hash_fn(x).hexdigest().encode()

        def KD(s: bytes, d: bytes) -> bytes:
            return H(b':'.join((s, d)))
        A1 = b':'.join((self._login_bytes, realm_bytes, self._password_bytes))
        A2 = f'{method.upper()}:{path}'.encode()
        if qop == 'auth-int':
            if isinstance(body, Payload):
                entity_bytes = await body.as_bytes()
            else:
                entity_bytes = body
            entity_hash = H(entity_bytes)
            A2 = b':'.join((A2, entity_hash))
        HA1 = H(A1)
        HA2 = H(A2)
        if nonce_bytes == self._last_nonce_bytes:
            self._nonce_count += 1
        else:
            self._nonce_count = 1
        self._last_nonce_bytes = nonce_bytes
        ncvalue = f'{self._nonce_count:08x}'
        ncvalue_bytes = ncvalue.encode('utf-8')
        cnonce = hashlib.sha1(b''.join([str(self._nonce_count).encode('utf-8'), nonce_bytes, time.ctime().encode('utf-8'), os.urandom(8)])).hexdigest()[:16]
        cnonce_bytes = cnonce.encode('utf-8')
        if algorithm.upper().endswith('-SESS'):
            HA1 = H(b':'.join((HA1, nonce_bytes, cnonce_bytes)))
        if qop:
            noncebit = b':'.join((nonce_bytes, ncvalue_bytes, cnonce_bytes, qop_bytes, HA2))
            response_digest = KD(HA1, noncebit)
        else:
            response_digest = KD(HA1, b':'.join((nonce_bytes, HA2)))
        header_fields = {'username': escape_quotes(self._login_str), 'realm': escape_quotes(realm), 'nonce': escape_quotes(nonce), 'uri': path, 'response': response_digest.decode(), 'algorithm': algorithm_original}
        if opaque:
            header_fields['opaque'] = escape_quotes(opaque)
        if qop:
            header_fields['qop'] = qop
            header_fields['nc'] = ncvalue
            header_fields['cnonce'] = cnonce
        pairs: List[str] = []
        for field, value in header_fields.items():
            if field in QUOTED_AUTH_FIELDS:
                pairs.append(f'{field}="{value}"')
            else:
                pairs.append(f'{field}={value}')
        return f"Digest {', '.join(pairs)}"

    def _in_protection_space(self, url: URL) -> bool:
        request_str = str(url)
        for space_str in self._protection_space:
            if not request_str.startswith(space_str):
                continue
            if len(request_str) == len(space_str) or space_str[-1] == '/':
                return True
            if request_str[len(space_str)] == '/':
                return True
        return False

    def _authenticate(self, response: ClientResponse) -> bool:
        if response.status != 401:
            return False
        auth_header = response.headers.get('www-authenticate', '')
        if not auth_header:
            return False
        method, sep, headers = auth_header.partition(' ')
        if not sep:
            return False
        if method.lower() != 'digest':
            return False
        if not headers:
            return False
        if not (header_pairs := parse_header_pairs(headers)):
            return False
        self._challenge = {}
        for field in CHALLENGE_FIELDS:
            if (value := header_pairs.get(field)) is not None:
                self._challenge[field] = value
        origin = response.url.origin()
        if (domain := self._challenge.get('domain')):
            self._protection_space = []
            for uri in domain.split():
                uri = uri.strip('"')
                if uri.startswith('/'):
                    self._protection_space.append(str(origin.join(URL(uri))))
                else:
                    self._protection_space.append(str(URL(uri)))
        else:
            self._protection_space = [str(origin)]
        return bool(self._challenge)

    async def __call__(self, request: ClientRequest, handler: ClientHandlerType) -> ClientResponse:
        response = None
        for retry_count in range(2):
            if retry_count > 0 or (self._preemptive and self._challenge and self._in_protection_space(request.url)):
                request.headers[hdrs.AUTHORIZATION] = await self._encode(request.method, request.url, request.body)
            response = await handler(request)
            if not self._authenticate(response):
                break
        assert response is not None
        return response
