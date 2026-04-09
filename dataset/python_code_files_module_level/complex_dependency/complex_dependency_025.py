from __future__ import annotations
import functools
import logging
import typing
import warnings
from types import TracebackType
from urllib.parse import urljoin
from ._collections import HTTPHeaderDict, RecentlyUsedContainer
from ._request_methods import RequestMethods
from .connection import ProxyConfig
from .connectionpool import HTTPConnectionPool, HTTPSConnectionPool, port_by_scheme
from .exceptions import LocationValueError, MaxRetryError, ProxySchemeUnknown, URLSchemeUnknown
from .response import BaseHTTPResponse
from .util.connection import _TYPE_SOCKET_OPTIONS
from .util.proxy import connection_requires_http_tunnel
from .util.retry import Retry
from .util.timeout import Timeout
from .util.url import Url, parse_url
if typing.TYPE_CHECKING:
    import ssl
    from typing_extensions import Self
__all__ = ['PoolManager', 'ProxyManager', 'proxy_from_url']
log = logging.getLogger(__name__)
SSL_KEYWORDS = ('key_file', 'cert_file', 'cert_reqs', 'ca_certs', 'ca_cert_data', 'ssl_version', 'ssl_minimum_version', 'ssl_maximum_version', 'ca_cert_dir', 'ssl_context', 'key_password', 'server_hostname')
_DEFAULT_BLOCKSIZE = 16384

class PoolKey(typing.NamedTuple):
    key_scheme: str
    key_host: str
    key_port: int | None
    key_timeout: Timeout | float | int | None
    key_retries: Retry | bool | int | None
    key_block: bool | None
    key_source_address: tuple[str, int] | None
    key_key_file: str | None
    key_key_password: str | None
    key_cert_file: str | None
    key_cert_reqs: str | None
    key_ca_certs: str | None
    key_ca_cert_data: str | bytes | None
    key_ssl_version: int | str | None
    key_ssl_minimum_version: ssl.TLSVersion | None
    key_ssl_maximum_version: ssl.TLSVersion | None
    key_ca_cert_dir: str | None
    key_ssl_context: ssl.SSLContext | None
    key_maxsize: int | None
    key_headers: frozenset[tuple[str, str]] | None
    key__proxy: Url | None
    key__proxy_headers: frozenset[tuple[str, str]] | None
    key__proxy_config: ProxyConfig | None
    key_socket_options: _TYPE_SOCKET_OPTIONS | None
    key__socks_options: frozenset[tuple[str, str]] | None
    key_assert_hostname: bool | str | None
    key_assert_fingerprint: str | None
    key_server_hostname: str | None
    key_blocksize: int | None

def _default_key_normalizer(key_class: type[PoolKey], request_context: dict[str, typing.Any]) -> PoolKey:
    context = request_context.copy()
    context['scheme'] = context['scheme'].lower()
    context['host'] = context['host'].lower()
    for key in ('headers', '_proxy_headers', '_socks_options'):
        if key in context and context[key] is not None:
            context[key] = frozenset(context[key].items())
    socket_opts = context.get('socket_options')
    if socket_opts is not None:
        context['socket_options'] = tuple(socket_opts)
    for key in list(context.keys()):
        context['key_' + key] = context.pop(key)
    for field in key_class._fields:
        if field not in context:
            context[field] = None
    if context.get('key_blocksize') is None:
        context['key_blocksize'] = _DEFAULT_BLOCKSIZE
    return key_class(**context)
key_fn_by_scheme = {'http': functools.partial(_default_key_normalizer, PoolKey), 'https': functools.partial(_default_key_normalizer, PoolKey)}
pool_classes_by_scheme = {'http': HTTPConnectionPool, 'https': HTTPSConnectionPool}

class PoolManager(RequestMethods):
    proxy: Url | None = None
    proxy_config: ProxyConfig | None = None

    def __init__(self, num_pools: int=10, headers: typing.Mapping[str, str] | None=None, **connection_pool_kw: typing.Any) -> None:
        super().__init__(headers)
        if 'retries' in connection_pool_kw:
            retries = connection_pool_kw['retries']
            if not isinstance(retries, Retry):
                retries = Retry.from_int(retries)
                connection_pool_kw = connection_pool_kw.copy()
                connection_pool_kw['retries'] = retries
        self.connection_pool_kw = connection_pool_kw
        self.pools: RecentlyUsedContainer[PoolKey, HTTPConnectionPool]
        self.pools = RecentlyUsedContainer(num_pools)
        self.pool_classes_by_scheme = pool_classes_by_scheme
        self.key_fn_by_scheme = key_fn_by_scheme.copy()

    def __enter__(self) -> Self:
        return self

    def __exit__(self, exc_type: type[BaseException] | None, exc_val: BaseException | None, exc_tb: TracebackType | None) -> typing.Literal[False]:
        self.clear()
        return False

    def _new_pool(self, scheme: str, host: str, port: int, request_context: dict[str, typing.Any] | None=None) -> HTTPConnectionPool:
        pool_cls: type[HTTPConnectionPool] = self.pool_classes_by_scheme[scheme]
        if request_context is None:
            request_context = self.connection_pool_kw.copy()
        if request_context.get('blocksize') is None:
            request_context['blocksize'] = _DEFAULT_BLOCKSIZE
        for key in ('scheme', 'host', 'port'):
            request_context.pop(key, None)
        if scheme == 'http':
            for kw in SSL_KEYWORDS:
                request_context.pop(kw, None)
        return pool_cls(host, port, **request_context)

    def clear(self) -> None:
        self.pools.clear()

    def connection_from_host(self, host: str | None, port: int | None=None, scheme: str | None='http', pool_kwargs: dict[str, typing.Any] | None=None) -> HTTPConnectionPool:
        if not host:
            raise LocationValueError('No host specified.')
        request_context = self._merge_pool_kwargs(pool_kwargs)
        request_context['scheme'] = scheme or 'http'
        if not port:
            port = port_by_scheme.get(request_context['scheme'].lower(), 80)
        request_context['port'] = port
        request_context['host'] = host
        return self.connection_from_context(request_context)

    def connection_from_context(self, request_context: dict[str, typing.Any]) -> HTTPConnectionPool:
        if 'strict' in request_context:
            warnings.warn("The 'strict' parameter is no longer needed on Python 3+. This will raise an error in urllib3 v2.1.0.", DeprecationWarning)
            request_context.pop('strict')
        scheme = request_context['scheme'].lower()
        pool_key_constructor = self.key_fn_by_scheme.get(scheme)
        if not pool_key_constructor:
            raise URLSchemeUnknown(scheme)
        pool_key = pool_key_constructor(request_context)
        return self.connection_from_pool_key(pool_key, request_context=request_context)

    def connection_from_pool_key(self, pool_key: PoolKey, request_context: dict[str, typing.Any]) -> HTTPConnectionPool:
        with self.pools.lock:
            pool = self.pools.get(pool_key)
            if pool:
                return pool
            scheme = request_context['scheme']
            host = request_context['host']
            port = request_context['port']
            pool = self._new_pool(scheme, host, port, request_context=request_context)
            self.pools[pool_key] = pool
        return pool

    def connection_from_url(self, url: str, pool_kwargs: dict[str, typing.Any] | None=None) -> HTTPConnectionPool:
        u = parse_url(url)
        return self.connection_from_host(u.host, port=u.port, scheme=u.scheme, pool_kwargs=pool_kwargs)

    def _merge_pool_kwargs(self, override: dict[str, typing.Any] | None) -> dict[str, typing.Any]:
        base_pool_kwargs = self.connection_pool_kw.copy()
        if override:
            for key, value in override.items():
                if value is None:
                    try:
                        del base_pool_kwargs[key]
                    except KeyError:
                        pass
                else:
                    base_pool_kwargs[key] = value
        return base_pool_kwargs

    def _proxy_requires_url_absolute_form(self, parsed_url: Url) -> bool:
        if self.proxy is None:
            return False
        return not connection_requires_http_tunnel(self.proxy, self.proxy_config, parsed_url.scheme)

    def urlopen(self, method: str, url: str, redirect: bool=True, **kw: typing.Any) -> BaseHTTPResponse:
        u = parse_url(url)
        if u.scheme is None:
            warnings.warn("URLs without a scheme (ie 'https://') are deprecated and will raise an error in a future version of urllib3. To avoid this DeprecationWarning ensure all URLs start with 'https://' or 'http://'. Read more in this issue: https://github.com/urllib3/urllib3/issues/2920", category=DeprecationWarning, stacklevel=2)
        conn = self.connection_from_host(u.host, port=u.port, scheme=u.scheme)
        kw['assert_same_host'] = False
        kw['redirect'] = False
        if 'headers' not in kw:
            kw['headers'] = self.headers
        if self._proxy_requires_url_absolute_form(u):
            response = conn.urlopen(method, url, **kw)
        else:
            response = conn.urlopen(method, u.request_uri, **kw)
        redirect_location = redirect and response.get_redirect_location()
        if not redirect_location:
            return response
        redirect_location = urljoin(url, redirect_location)
        if response.status == 303:
            method = 'GET'
            kw['body'] = None
            kw['headers'] = HTTPHeaderDict(kw['headers'])._prepare_for_method_change()
        retries = kw.get('retries', response.retries)
        if not isinstance(retries, Retry):
            retries = Retry.from_int(retries, redirect=redirect)
        if retries.remove_headers_on_redirect and (not conn.is_same_host(redirect_location)):
            new_headers = kw['headers'].copy()
            for header in kw['headers']:
                if header.lower() in retries.remove_headers_on_redirect:
                    new_headers.pop(header, None)
            kw['headers'] = new_headers
        try:
            retries = retries.increment(method, url, response=response, _pool=conn)
        except MaxRetryError:
            if retries.raise_on_redirect:
                response.drain_conn()
                raise
            return response
        kw['retries'] = retries
        kw['redirect'] = redirect
        log.info('Redirecting %s -> %s', url, redirect_location)
        response.drain_conn()
        return self.urlopen(method, redirect_location, **kw)

class ProxyManager(PoolManager):

    def __init__(self, proxy_url: str, num_pools: int=10, headers: typing.Mapping[str, str] | None=None, proxy_headers: typing.Mapping[str, str] | None=None, proxy_ssl_context: ssl.SSLContext | None=None, use_forwarding_for_https: bool=False, proxy_assert_hostname: None | str | typing.Literal[False]=None, proxy_assert_fingerprint: str | None=None, **connection_pool_kw: typing.Any) -> None:
        if isinstance(proxy_url, HTTPConnectionPool):
            str_proxy_url = f'{proxy_url.scheme}://{proxy_url.host}:{proxy_url.port}'
        else:
            str_proxy_url = proxy_url
        proxy = parse_url(str_proxy_url)
        if proxy.scheme not in ('http', 'https'):
            raise ProxySchemeUnknown(proxy.scheme)
        if not proxy.port:
            port = port_by_scheme.get(proxy.scheme, 80)
            proxy = proxy._replace(port=port)
        self.proxy = proxy
        self.proxy_headers = proxy_headers or {}
        self.proxy_ssl_context = proxy_ssl_context
        self.proxy_config = ProxyConfig(proxy_ssl_context, use_forwarding_for_https, proxy_assert_hostname, proxy_assert_fingerprint)
        connection_pool_kw['_proxy'] = self.proxy
        connection_pool_kw['_proxy_headers'] = self.proxy_headers
        connection_pool_kw['_proxy_config'] = self.proxy_config
        super().__init__(num_pools, headers, **connection_pool_kw)

    def connection_from_host(self, host: str | None, port: int | None=None, scheme: str | None='http', pool_kwargs: dict[str, typing.Any] | None=None) -> HTTPConnectionPool:
        if scheme == 'https':
            return super().connection_from_host(host, port, scheme, pool_kwargs=pool_kwargs)
        return super().connection_from_host(self.proxy.host, self.proxy.port, self.proxy.scheme, pool_kwargs=pool_kwargs)

    def _set_proxy_headers(self, url: str, headers: typing.Mapping[str, str] | None=None) -> typing.Mapping[str, str]:
        headers_ = {'Accept': '*/*'}
        netloc = parse_url(url).netloc
        if netloc:
            headers_['Host'] = netloc
        if headers:
            headers_.update(headers)
        return headers_

    def urlopen(self, method: str, url: str, redirect: bool=True, **kw: typing.Any) -> BaseHTTPResponse:
        u = parse_url(url)
        if not connection_requires_http_tunnel(self.proxy, self.proxy_config, u.scheme):
            headers = kw.get('headers', self.headers)
            kw['headers'] = self._set_proxy_headers(url, headers)
        return super().urlopen(method, url, redirect=redirect, **kw)

def proxy_from_url(url: str, **kw: typing.Any) -> ProxyManager:
    return ProxyManager(proxy_url=url, **kw)
