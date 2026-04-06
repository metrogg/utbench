from __future__ import annotations
import email
import logging
import random
import re
import time
import typing
from itertools import takewhile
from types import TracebackType
from ..exceptions import ConnectTimeoutError, InvalidHeader, MaxRetryError, ProtocolError, ProxyError, ReadTimeoutError, ResponseError
from .util import reraise
if typing.TYPE_CHECKING:
    from typing_extensions import Self
    from ..connectionpool import ConnectionPool
    from ..response import BaseHTTPResponse
log = logging.getLogger(__name__)

class RequestHistory(typing.NamedTuple):
    method: str | None
    url: str | None
    error: Exception | None
    status: int | None
    redirect_location: str | None

class Retry:
    DEFAULT_ALLOWED_METHODS = frozenset(['HEAD', 'GET', 'PUT', 'DELETE', 'OPTIONS', 'TRACE'])
    RETRY_AFTER_STATUS_CODES = frozenset([413, 429, 503])
    DEFAULT_REMOVE_HEADERS_ON_REDIRECT = frozenset(['Cookie', 'Authorization', 'Proxy-Authorization'])
    DEFAULT_BACKOFF_MAX = 120
    DEFAULT_RETRY_AFTER_MAX: typing.Final[int] = 21600
    DEFAULT: typing.ClassVar[Retry]

    def __init__(self, total: bool | int | None=10, connect: int | None=None, read: int | None=None, redirect: bool | int | None=None, status: int | None=None, other: int | None=None, allowed_methods: typing.Collection[str] | None=DEFAULT_ALLOWED_METHODS, status_forcelist: typing.Collection[int] | None=None, backoff_factor: float=0, backoff_max: float=DEFAULT_BACKOFF_MAX, raise_on_redirect: bool=True, raise_on_status: bool=True, history: tuple[RequestHistory, ...] | None=None, respect_retry_after_header: bool=True, remove_headers_on_redirect: typing.Collection[str]=DEFAULT_REMOVE_HEADERS_ON_REDIRECT, backoff_jitter: float=0.0, retry_after_max: int=DEFAULT_RETRY_AFTER_MAX) -> None:
        self.total = total
        self.connect = connect
        self.read = read
        self.status = status
        self.other = other
        if redirect is False or total is False:
            redirect = 0
            raise_on_redirect = False
        self.redirect = redirect
        self.status_forcelist = status_forcelist or set()
        self.allowed_methods = allowed_methods
        self.backoff_factor = backoff_factor
        self.backoff_max = backoff_max
        self.retry_after_max = retry_after_max
        self.raise_on_redirect = raise_on_redirect
        self.raise_on_status = raise_on_status
        self.history = history or ()
        self.respect_retry_after_header = respect_retry_after_header
        self.remove_headers_on_redirect = frozenset((h.lower() for h in remove_headers_on_redirect))
        self.backoff_jitter = backoff_jitter

    def new(self, **kw: typing.Any) -> Self:
        params = dict(total=self.total, connect=self.connect, read=self.read, redirect=self.redirect, status=self.status, other=self.other, allowed_methods=self.allowed_methods, status_forcelist=self.status_forcelist, backoff_factor=self.backoff_factor, backoff_max=self.backoff_max, retry_after_max=self.retry_after_max, raise_on_redirect=self.raise_on_redirect, raise_on_status=self.raise_on_status, history=self.history, remove_headers_on_redirect=self.remove_headers_on_redirect, respect_retry_after_header=self.respect_retry_after_header, backoff_jitter=self.backoff_jitter)
        params.update(kw)
        return type(self)(**params)

    @classmethod
    def from_int(cls, retries: Retry | bool | int | None, redirect: bool | int | None=True, default: Retry | bool | int | None=None) -> Retry:
        if retries is None:
            retries = default if default is not None else cls.DEFAULT
        if isinstance(retries, Retry):
            return retries
        redirect = bool(redirect) and None
        new_retries = cls(retries, redirect=redirect)
        log.debug('Converted retries value: %r -> %r', retries, new_retries)
        return new_retries

    def get_backoff_time(self) -> float:
        consecutive_errors_len = len(list(takewhile(lambda x: x.redirect_location is None, reversed(self.history))))
        if consecutive_errors_len <= 1:
            return 0
        backoff_value = self.backoff_factor * 2 ** (consecutive_errors_len - 1)
        if self.backoff_jitter != 0.0:
            backoff_value += random.random() * self.backoff_jitter
        return float(max(0, min(self.backoff_max, backoff_value)))

    def parse_retry_after(self, retry_after: str) -> float:
        seconds: float
        if re.match('^\\s*[0-9]+\\s*$', retry_after):
            seconds = int(retry_after)
        else:
            retry_date_tuple = email.utils.parsedate_tz(retry_after)
            if retry_date_tuple is None:
                raise InvalidHeader(f'Invalid Retry-After header: {retry_after}')
            retry_date = email.utils.mktime_tz(retry_date_tuple)
            seconds = retry_date - time.time()
        seconds = max(seconds, 0)
        if seconds > self.retry_after_max:
            seconds = self.retry_after_max
        return seconds

    def get_retry_after(self, response: BaseHTTPResponse) -> float | None:
        retry_after = response.headers.get('Retry-After')
        if retry_after is None:
            return None
        return self.parse_retry_after(retry_after)

    def sleep_for_retry(self, response: BaseHTTPResponse) -> bool:
        retry_after = self.get_retry_after(response)
        if retry_after:
            time.sleep(retry_after)
            return True
        return False

    def _sleep_backoff(self) -> None:
        backoff = self.get_backoff_time()
        if backoff <= 0:
            return
        time.sleep(backoff)

    def sleep(self, response: BaseHTTPResponse | None=None) -> None:
        if self.respect_retry_after_header and response:
            slept = self.sleep_for_retry(response)
            if slept:
                return
        self._sleep_backoff()

    def _is_connection_error(self, err: Exception) -> bool:
        if isinstance(err, ProxyError):
            err = err.original_error
        return isinstance(err, ConnectTimeoutError)

    def _is_read_error(self, err: Exception) -> bool:
        return isinstance(err, (ReadTimeoutError, ProtocolError))

    def _is_method_retryable(self, method: str) -> bool:
        if self.allowed_methods and method.upper() not in self.allowed_methods:
            return False
        return True

    def is_retry(self, method: str, status_code: int, has_retry_after: bool=False) -> bool:
        if not self._is_method_retryable(method):
            return False
        if self.status_forcelist and status_code in self.status_forcelist:
            return True
        return bool(self.total and self.respect_retry_after_header and has_retry_after and (status_code in self.RETRY_AFTER_STATUS_CODES))

    def is_exhausted(self) -> bool:
        retry_counts = [x for x in (self.total, self.connect, self.read, self.redirect, self.status, self.other) if x]
        if not retry_counts:
            return False
        return min(retry_counts) < 0

    def increment(self, method: str | None=None, url: str | None=None, response: BaseHTTPResponse | None=None, error: Exception | None=None, _pool: ConnectionPool | None=None, _stacktrace: TracebackType | None=None) -> Self:
        if self.total is False and error:
            raise reraise(type(error), error, _stacktrace)
        total = self.total
        if total is not None:
            total -= 1
        connect = self.connect
        read = self.read
        redirect = self.redirect
        status_count = self.status
        other = self.other
        cause = 'unknown'
        status = None
        redirect_location = None
        if error and self._is_connection_error(error):
            if connect is False:
                raise reraise(type(error), error, _stacktrace)
            elif connect is not None:
                connect -= 1
        elif error and self._is_read_error(error):
            if read is False or method is None or (not self._is_method_retryable(method)):
                raise reraise(type(error), error, _stacktrace)
            elif read is not None:
                read -= 1
        elif error:
            if other is not None:
                other -= 1
        elif response and response.get_redirect_location():
            if redirect is not None:
                redirect -= 1
            cause = 'too many redirects'
            response_redirect_location = response.get_redirect_location()
            if response_redirect_location:
                redirect_location = response_redirect_location
            status = response.status
        else:
            cause = ResponseError.GENERIC_ERROR
            if response and response.status:
                if status_count is not None:
                    status_count -= 1
                cause = ResponseError.SPECIFIC_ERROR.format(status_code=response.status)
                status = response.status
        history = self.history + (RequestHistory(method, url, error, status, redirect_location),)
        new_retry = self.new(total=total, connect=connect, read=read, redirect=redirect, status=status_count, other=other, history=history)
        if new_retry.is_exhausted():
            reason = error or ResponseError(cause)
            raise MaxRetryError(_pool, url, reason) from reason
        log.debug("Incremented Retry for (url='%s'): %r", url, new_retry)
        return new_retry

    def __repr__(self) -> str:
        return f'{type(self).__name__}(total={self.total}, connect={self.connect}, read={self.read}, redirect={self.redirect}, status={self.status})'
Retry.DEFAULT = Retry(3)
