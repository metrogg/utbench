from __future__ import annotations
import re
import typing
from ..exceptions import LocationParseError
from .util import to_str
_NORMALIZABLE_SCHEMES = ('http', 'https', None)
_PERCENT_RE = re.compile('%[a-fA-F0-9]{2}')
_SCHEME_RE = re.compile('^(?:[a-zA-Z][a-zA-Z0-9+-]*:|/)')
_URI_RE = re.compile('^(?:([a-zA-Z][a-zA-Z0-9+.-]*):)?(?://([^\\\\/?#]*))?([^?#]*)(?:\\?([^#]*))?(?:#(.*))?$', re.UNICODE | re.DOTALL)
_IPV4_PAT = '(?:[0-9]{1,3}\\.){3}[0-9]{1,3}'
_HEX_PAT = '[0-9A-Fa-f]{1,4}'
_LS32_PAT = '(?:{hex}:{hex}|{ipv4})'.format(hex=_HEX_PAT, ipv4=_IPV4_PAT)
_subs = {'hex': _HEX_PAT, 'ls32': _LS32_PAT}
_variations = ['(?:%(hex)s:){6}%(ls32)s', '::(?:%(hex)s:){5}%(ls32)s', '(?:%(hex)s)?::(?:%(hex)s:){4}%(ls32)s', '(?:(?:%(hex)s:)?%(hex)s)?::(?:%(hex)s:){3}%(ls32)s', '(?:(?:%(hex)s:){0,2}%(hex)s)?::(?:%(hex)s:){2}%(ls32)s', '(?:(?:%(hex)s:){0,3}%(hex)s)?::%(hex)s:%(ls32)s', '(?:(?:%(hex)s:){0,4}%(hex)s)?::%(ls32)s', '(?:(?:%(hex)s:){0,5}%(hex)s)?::%(hex)s', '(?:(?:%(hex)s:){0,6}%(hex)s)?::']
_UNRESERVED_PAT = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789._\\-~'
_IPV6_PAT = '(?:' + '|'.join([x % _subs for x in _variations]) + ')'
_ZONE_ID_PAT = '(?:%25|%)(?:[' + _UNRESERVED_PAT + ']|%[a-fA-F0-9]{2})+'
_IPV6_ADDRZ_PAT = '\\[' + _IPV6_PAT + '(?:' + _ZONE_ID_PAT + ')?\\]'
_REG_NAME_PAT = '(?:[^\\[\\]%:/?#]|%[a-fA-F0-9]{2})*'
_TARGET_RE = re.compile('^(/[^?#]*)(?:\\?([^#]*))?(?:#.*)?$')
_IPV4_RE = re.compile('^' + _IPV4_PAT + '$')
_IPV6_RE = re.compile('^' + _IPV6_PAT + '$')
_IPV6_ADDRZ_RE = re.compile('^' + _IPV6_ADDRZ_PAT + '$')
_BRACELESS_IPV6_ADDRZ_RE = re.compile('^' + _IPV6_ADDRZ_PAT[2:-2] + '$')
_ZONE_ID_RE = re.compile('(' + _ZONE_ID_PAT + ')\\]$')
_HOST_PORT_PAT = '^(%s|%s|%s)(?::0*?(|0|[1-9][0-9]{0,4}))?$' % (_REG_NAME_PAT, _IPV4_PAT, _IPV6_ADDRZ_PAT)
_HOST_PORT_RE = re.compile(_HOST_PORT_PAT, re.UNICODE | re.DOTALL)
_UNRESERVED_CHARS = set('ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789._-~')
_SUB_DELIM_CHARS = set("!$&'()*+,;=")
_USERINFO_CHARS = _UNRESERVED_CHARS | _SUB_DELIM_CHARS | {':'}
_PATH_CHARS = _USERINFO_CHARS | {'@', '/'}
_QUERY_CHARS = _FRAGMENT_CHARS = _PATH_CHARS | {'?'}

class Url(typing.NamedTuple('Url', [('scheme', typing.Optional[str]), ('auth', typing.Optional[str]), ('host', typing.Optional[str]), ('port', typing.Optional[int]), ('path', typing.Optional[str]), ('query', typing.Optional[str]), ('fragment', typing.Optional[str])])):

    def __new__(cls, scheme: str | None=None, auth: str | None=None, host: str | None=None, port: int | None=None, path: str | None=None, query: str | None=None, fragment: str | None=None):
        if path and (not path.startswith('/')):
            path = '/' + path
        if scheme is not None:
            scheme = scheme.lower()
        return super().__new__(cls, scheme, auth, host, port, path, query, fragment)

    @property
    def hostname(self) -> str | None:
        return self.host

    @property
    def request_uri(self) -> str:
        uri = self.path or '/'
        if self.query is not None:
            uri += '?' + self.query
        return uri

    @property
    def authority(self) -> str | None:
        userinfo = self.auth
        netloc = self.netloc
        if netloc is None or userinfo is None:
            return netloc
        else:
            return f'{userinfo}@{netloc}'

    @property
    def netloc(self) -> str | None:
        if self.host is None:
            return None
        if self.port:
            return f'{self.host}:{self.port}'
        return self.host

    @property
    def url(self) -> str:
        scheme, auth, host, port, path, query, fragment = self
        url = ''
        if scheme is not None:
            url += scheme + '://'
        if auth is not None:
            url += auth + '@'
        if host is not None:
            url += host
        if port is not None:
            url += ':' + str(port)
        if path is not None:
            url += path
        if query is not None:
            url += '?' + query
        if fragment is not None:
            url += '#' + fragment
        return url

    def __str__(self) -> str:
        return self.url

@typing.overload
def _encode_invalid_chars(component: str, allowed_chars: typing.Container[str]) -> str:
    ...

@typing.overload
def _encode_invalid_chars(component: None, allowed_chars: typing.Container[str]) -> None:
    ...

def _encode_invalid_chars(component: str | None, allowed_chars: typing.Container[str]) -> str | None:
    if component is None:
        return component
    component = to_str(component)
    component, percent_encodings = _PERCENT_RE.subn(lambda match: match.group(0).upper(), component)
    uri_bytes = component.encode('utf-8', 'surrogatepass')
    is_percent_encoded = percent_encodings == uri_bytes.count(b'%')
    encoded_component = bytearray()
    for i in range(0, len(uri_bytes)):
        byte = uri_bytes[i:i + 1]
        byte_ord = ord(byte)
        if is_percent_encoded and byte == b'%' or (byte_ord < 128 and byte.decode() in allowed_chars):
            encoded_component += byte
            continue
        encoded_component.extend(b'%' + hex(byte_ord)[2:].encode().zfill(2).upper())
    return encoded_component.decode()

def _remove_path_dot_segments(path: str) -> str:
    segments = path.split('/')
    output = []
    for segment in segments:
        if segment == '.':
            continue
        if segment != '..':
            output.append(segment)
        elif output:
            output.pop()
    if path.startswith('/') and (not output or output[0]):
        output.insert(0, '')
    if path.endswith(('/.', '/..')):
        output.append('')
    return '/'.join(output)

@typing.overload
def _normalize_host(host: None, scheme: str | None) -> None:
    ...

@typing.overload
def _normalize_host(host: str, scheme: str | None) -> str:
    ...

def _normalize_host(host: str | None, scheme: str | None) -> str | None:
    if host:
        if scheme in _NORMALIZABLE_SCHEMES:
            is_ipv6 = _IPV6_ADDRZ_RE.match(host)
            if is_ipv6:
                match = _ZONE_ID_RE.search(host)
                if match:
                    start, end = match.span(1)
                    zone_id = host[start:end]
                    if zone_id.startswith('%25') and zone_id != '%25':
                        zone_id = zone_id[3:]
                    else:
                        zone_id = zone_id[1:]
                    zone_id = _encode_invalid_chars(zone_id, _UNRESERVED_CHARS)
                    return f'{host[:start].lower()}%{zone_id}{host[end:]}'
                else:
                    return host.lower()
            elif not _IPV4_RE.match(host):
                return to_str(b'.'.join([_idna_encode(label) for label in host.split('.')]), 'ascii')
    return host

def _idna_encode(name: str) -> bytes:
    if not name.isascii():
        try:
            import idna
        except ImportError:
            raise LocationParseError("Unable to parse URL without the 'idna' module") from None
        try:
            return idna.encode(name.lower(), strict=True, std3_rules=True)
        except idna.IDNAError:
            raise LocationParseError(f"Name '{name}' is not a valid IDNA label") from None
    return name.lower().encode('ascii')

def _encode_target(target: str) -> str:
    match = _TARGET_RE.match(target)
    if not match:
        raise LocationParseError(f'{target!r} is not a valid request URI')
    path, query = match.groups()
    encoded_target = _encode_invalid_chars(path, _PATH_CHARS)
    if query is not None:
        query = _encode_invalid_chars(query, _QUERY_CHARS)
        encoded_target += '?' + query
    return encoded_target

def parse_url(url: str) -> Url:
    if not url:
        return Url()
    source_url = url
    if not _SCHEME_RE.search(url):
        url = '//' + url
    scheme: str | None
    authority: str | None
    auth: str | None
    host: str | None
    port: str | None
    port_int: int | None
    path: str | None
    query: str | None
    fragment: str | None
    try:
        scheme, authority, path, query, fragment = _URI_RE.match(url).groups()
        normalize_uri = scheme is None or scheme.lower() in _NORMALIZABLE_SCHEMES
        if scheme:
            scheme = scheme.lower()
        if authority:
            auth, _, host_port = authority.rpartition('@')
            auth = auth or None
            host, port = _HOST_PORT_RE.match(host_port).groups()
            if auth and normalize_uri:
                auth = _encode_invalid_chars(auth, _USERINFO_CHARS)
            if port == '':
                port = None
        else:
            auth, host, port = (None, None, None)
        if port is not None:
            port_int = int(port)
            if not 0 <= port_int <= 65535:
                raise LocationParseError(url)
        else:
            port_int = None
        host = _normalize_host(host, scheme)
        if normalize_uri and path:
            path = _remove_path_dot_segments(path)
            path = _encode_invalid_chars(path, _PATH_CHARS)
        if normalize_uri and query:
            query = _encode_invalid_chars(query, _QUERY_CHARS)
        if normalize_uri and fragment:
            fragment = _encode_invalid_chars(fragment, _FRAGMENT_CHARS)
    except (ValueError, AttributeError) as e:
        raise LocationParseError(source_url) from e
    if not path:
        if query is not None or fragment is not None:
            path = ''
        else:
            path = None
    return Url(scheme=scheme, auth=auth, host=host, port=port_int, path=path, query=query, fragment=fragment)
