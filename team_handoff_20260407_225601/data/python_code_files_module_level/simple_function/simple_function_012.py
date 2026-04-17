from __future__ import annotations
import typing
from urllib.parse import parse_qs, unquote, urlencode
import idna
from ._types import QueryParamTypes
from ._urlparse import urlparse
from ._utils import primitive_value_to_str
__all__ = ['URL', 'QueryParams']

class URL:

    def __init__(self, url: URL | str='', **kwargs: typing.Any) -> None:
        if kwargs:
            allowed = {'scheme': str, 'username': str, 'password': str, 'userinfo': bytes, 'host': str, 'port': int, 'netloc': bytes, 'path': str, 'query': bytes, 'raw_path': bytes, 'fragment': str, 'params': object}
            for key, value in kwargs.items():
                if key not in allowed:
                    message = f'{key!r} is an invalid keyword argument for URL()'
                    raise TypeError(message)
                if value is not None and (not isinstance(value, allowed[key])):
                    expected = allowed[key].__name__
                    seen = type(value).__name__
                    message = f'Argument {key!r} must be {expected} but got {seen}'
                    raise TypeError(message)
                if isinstance(value, bytes):
                    kwargs[key] = value.decode('ascii')
            if 'params' in kwargs:
                params = kwargs.pop('params')
                kwargs['query'] = None if not params else str(QueryParams(params))
        if isinstance(url, str):
            self._uri_reference = urlparse(url, **kwargs)
        elif isinstance(url, URL):
            self._uri_reference = url._uri_reference.copy_with(**kwargs)
        else:
            raise TypeError(f'Invalid type for url.  Expected str or httpx.URL, got {type(url)}: {url!r}')

    @property
    def scheme(self) -> str:
        return self._uri_reference.scheme

    @property
    def raw_scheme(self) -> bytes:
        return self._uri_reference.scheme.encode('ascii')

    @property
    def userinfo(self) -> bytes:
        return self._uri_reference.userinfo.encode('ascii')

    @property
    def username(self) -> str:
        userinfo = self._uri_reference.userinfo
        return unquote(userinfo.partition(':')[0])

    @property
    def password(self) -> str:
        userinfo = self._uri_reference.userinfo
        return unquote(userinfo.partition(':')[2])

    @property
    def host(self) -> str:
        host: str = self._uri_reference.host
        if host.startswith('xn--'):
            host = idna.decode(host)
        return host

    @property
    def raw_host(self) -> bytes:
        return self._uri_reference.host.encode('ascii')

    @property
    def port(self) -> int | None:
        return self._uri_reference.port

    @property
    def netloc(self) -> bytes:
        return self._uri_reference.netloc.encode('ascii')

    @property
    def path(self) -> str:
        path = self._uri_reference.path or '/'
        return unquote(path)

    @property
    def query(self) -> bytes:
        query = self._uri_reference.query or ''
        return query.encode('ascii')

    @property
    def params(self) -> QueryParams:
        return QueryParams(self._uri_reference.query)

    @property
    def raw_path(self) -> bytes:
        path = self._uri_reference.path or '/'
        if self._uri_reference.query is not None:
            path += '?' + self._uri_reference.query
        return path.encode('ascii')

    @property
    def fragment(self) -> str:
        return unquote(self._uri_reference.fragment or '')

    @property
    def is_absolute_url(self) -> bool:
        return bool(self._uri_reference.scheme and self._uri_reference.host)

    @property
    def is_relative_url(self) -> bool:
        return not self.is_absolute_url

    def copy_with(self, **kwargs: typing.Any) -> URL:
        return URL(self, **kwargs)

    def copy_set_param(self, key: str, value: typing.Any=None) -> URL:
        return self.copy_with(params=self.params.set(key, value))

    def copy_add_param(self, key: str, value: typing.Any=None) -> URL:
        return self.copy_with(params=self.params.add(key, value))

    def copy_remove_param(self, key: str) -> URL:
        return self.copy_with(params=self.params.remove(key))

    def copy_merge_params(self, params: QueryParamTypes) -> URL:
        return self.copy_with(params=self.params.merge(params))

    def join(self, url: URL | str) -> URL:
        from urllib.parse import urljoin
        return URL(urljoin(str(self), str(URL(url))))

    def __hash__(self) -> int:
        return hash(str(self))

    def __eq__(self, other: typing.Any) -> bool:
        return isinstance(other, (URL, str)) and str(self) == str(URL(other))

    def __str__(self) -> str:
        return str(self._uri_reference)

    def __repr__(self) -> str:
        scheme, userinfo, host, port, path, query, fragment = self._uri_reference
        if ':' in userinfo:
            userinfo = f"{userinfo.split(':')[0]}:[secure]"
        authority = ''.join([f'{userinfo}@' if userinfo else '', f'[{host}]' if ':' in host else host, f':{port}' if port is not None else ''])
        url = ''.join([f'{self.scheme}:' if scheme else '', f'//{authority}' if authority else '', path, f'?{query}' if query is not None else '', f'#{fragment}' if fragment is not None else ''])
        return f'{self.__class__.__name__}({url!r})'

    @property
    def raw(self) -> tuple[bytes, bytes, int, bytes]:
        import collections
        import warnings
        warnings.warn('URL.raw is deprecated.')
        RawURL = collections.namedtuple('RawURL', ['raw_scheme', 'raw_host', 'port', 'raw_path'])
        return RawURL(raw_scheme=self.raw_scheme, raw_host=self.raw_host, port=self.port, raw_path=self.raw_path)

class QueryParams(typing.Mapping[str, str]):

    def __init__(self, *args: QueryParamTypes | None, **kwargs: typing.Any) -> None:
        assert len(args) < 2, 'Too many arguments.'
        assert not (args and kwargs), 'Cannot mix named and unnamed arguments.'
        value = args[0] if args else kwargs
        if value is None or isinstance(value, (str, bytes)):
            value = value.decode('ascii') if isinstance(value, bytes) else value
            self._dict = parse_qs(value, keep_blank_values=True)
        elif isinstance(value, QueryParams):
            self._dict = {k: list(v) for k, v in value._dict.items()}
        else:
            dict_value: dict[typing.Any, list[typing.Any]] = {}
            if isinstance(value, (list, tuple)):
                for item in value:
                    dict_value.setdefault(item[0], []).append(item[1])
            else:
                dict_value = {k: list(v) if isinstance(v, (list, tuple)) else [v] for k, v in value.items()}
            self._dict = {str(k): [primitive_value_to_str(item) for item in v] for k, v in dict_value.items()}

    def keys(self) -> typing.KeysView[str]:
        return self._dict.keys()

    def values(self) -> typing.ValuesView[str]:
        return {k: v[0] for k, v in self._dict.items()}.values()

    def items(self) -> typing.ItemsView[str, str]:
        return {k: v[0] for k, v in self._dict.items()}.items()

    def multi_items(self) -> list[tuple[str, str]]:
        multi_items: list[tuple[str, str]] = []
        for k, v in self._dict.items():
            multi_items.extend([(k, i) for i in v])
        return multi_items

    def get(self, key: typing.Any, default: typing.Any=None) -> typing.Any:
        if key in self._dict:
            return self._dict[str(key)][0]
        return default

    def get_list(self, key: str) -> list[str]:
        return list(self._dict.get(str(key), []))

    def set(self, key: str, value: typing.Any=None) -> QueryParams:
        q = QueryParams()
        q._dict = dict(self._dict)
        q._dict[str(key)] = [primitive_value_to_str(value)]
        return q

    def add(self, key: str, value: typing.Any=None) -> QueryParams:
        q = QueryParams()
        q._dict = dict(self._dict)
        q._dict[str(key)] = q.get_list(key) + [primitive_value_to_str(value)]
        return q

    def remove(self, key: str) -> QueryParams:
        q = QueryParams()
        q._dict = dict(self._dict)
        q._dict.pop(str(key), None)
        return q

    def merge(self, params: QueryParamTypes | None=None) -> QueryParams:
        q = QueryParams(params)
        q._dict = {**self._dict, **q._dict}
        return q

    def __getitem__(self, key: typing.Any) -> str:
        return self._dict[key][0]

    def __contains__(self, key: typing.Any) -> bool:
        return key in self._dict

    def __iter__(self) -> typing.Iterator[typing.Any]:
        return iter(self.keys())

    def __len__(self) -> int:
        return len(self._dict)

    def __bool__(self) -> bool:
        return bool(self._dict)

    def __hash__(self) -> int:
        return hash(str(self))

    def __eq__(self, other: typing.Any) -> bool:
        if not isinstance(other, self.__class__):
            return False
        return sorted(self.multi_items()) == sorted(other.multi_items())

    def __str__(self) -> str:
        return urlencode(self.multi_items())

    def __repr__(self) -> str:
        class_name = self.__class__.__name__
        query_string = str(self)
        return f'{class_name}({query_string!r})'

    def update(self, params: QueryParamTypes | None=None) -> None:
        raise RuntimeError('QueryParams are immutable since 0.18.0. Use `q = q.merge(...)` to create an updated copy.')

    def __setitem__(self, key: str, value: str) -> None:
        raise RuntimeError('QueryParams are immutable since 0.18.0. Use `q = q.set(key, value)` to create an updated copy.')
