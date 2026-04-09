from __future__ import annotations
import collections.abc
from collections.abc import Callable
from collections.abc import Collection
from collections.abc import Iterable
from collections.abc import Iterator
from collections.abc import Mapping
from collections.abc import MutableMapping
from collections.abc import Sequence
import dataclasses
import enum
import inspect
from typing import Any
from typing import final
from typing import NamedTuple
from typing import overload
from typing import TYPE_CHECKING
from typing import TypeVar
import warnings
from .._code import getfslineno
from ..compat import NOTSET
from ..compat import NotSetType
from _pytest.config import Config
from _pytest.deprecated import check_ispytest
from _pytest.deprecated import MARKED_FIXTURE
from _pytest.outcomes import fail
from _pytest.raises import AbstractRaises
from _pytest.scope import _ScopeName
from _pytest.warning_types import PytestUnknownMarkWarning
if TYPE_CHECKING:
    from ..nodes import Node
EMPTY_PARAMETERSET_OPTION = 'empty_parameter_set_mark'

class _HiddenParam(enum.Enum):
    token = 0
HIDDEN_PARAM = _HiddenParam.token

def istestfunc(func) -> bool:
    return callable(func) and getattr(func, '__name__', '<lambda>') != '<lambda>'

def get_empty_parameterset_mark(config: Config, argnames: Sequence[str], func) -> MarkDecorator:
    from ..nodes import Collector
    argslisting = ', '.join(argnames)
    _fs, lineno = getfslineno(func)
    reason = f'got empty parameter set for ({argslisting})'
    requested_mark = config.getini(EMPTY_PARAMETERSET_OPTION)
    if requested_mark in ('', None, 'skip'):
        mark = MARK_GEN.skip(reason=reason)
    elif requested_mark == 'xfail':
        mark = MARK_GEN.xfail(reason=reason, run=False)
    elif requested_mark == 'fail_at_collect':
        raise Collector.CollectError(f"Empty parameter set in '{func.__name__}' at line {lineno + 1}")
    else:
        raise LookupError(requested_mark)
    return mark

class ParameterSet(NamedTuple):
    values: Sequence[object | NotSetType]
    marks: Collection[MarkDecorator | Mark]
    id: str | _HiddenParam | None

    @classmethod
    def param(cls, *values: object, marks: MarkDecorator | Collection[MarkDecorator | Mark]=(), id: str | _HiddenParam | None=None) -> ParameterSet:
        if isinstance(marks, MarkDecorator):
            marks = (marks,)
        else:
            assert isinstance(marks, collections.abc.Collection)
        if any((i.name == 'usefixtures' for i in marks)):
            raise ValueError('pytest.param cannot add pytest.mark.usefixtures; see https://docs.pytest.org/en/stable/reference/reference.html#pytest-param')
        if id is not None:
            if not isinstance(id, str) and id is not HIDDEN_PARAM:
                raise TypeError(f'Expected id to be a string or a `pytest.HIDDEN_PARAM` sentinel, got {type(id)}: {id!r}')
        return cls(values, marks, id)

    @classmethod
    def extract_from(cls, parameterset: ParameterSet | Sequence[object] | object, force_tuple: bool=False) -> ParameterSet:
        if isinstance(parameterset, cls):
            return parameterset
        if force_tuple:
            return cls.param(parameterset)
        else:
            return cls(parameterset, marks=[], id=None)

    @staticmethod
    def _parse_parametrize_args(argnames: str | Sequence[str], argvalues: Iterable[ParameterSet | Sequence[object] | object], *args, **kwargs) -> tuple[Sequence[str], bool]:
        if isinstance(argnames, str):
            argnames = [x.strip() for x in argnames.split(',') if x.strip()]
            force_tuple = len(argnames) == 1
        else:
            force_tuple = False
        return (argnames, force_tuple)

    @staticmethod
    def _parse_parametrize_parameters(argvalues: Iterable[ParameterSet | Sequence[object] | object], force_tuple: bool) -> list[ParameterSet]:
        return [ParameterSet.extract_from(x, force_tuple=force_tuple) for x in argvalues]

    @classmethod
    def _for_parametrize(cls, argnames: str | Sequence[str], argvalues: Iterable[ParameterSet | Sequence[object] | object], func, config: Config, nodeid: str) -> tuple[Sequence[str], list[ParameterSet]]:
        argnames, force_tuple = cls._parse_parametrize_args(argnames, argvalues)
        parameters = cls._parse_parametrize_parameters(argvalues, force_tuple)
        del argvalues
        if parameters:
            for param in parameters:
                if len(param.values) != len(argnames):
                    msg = '{nodeid}: in "parametrize" the number of names ({names_len}):\n  {names}\nmust be equal to the number of values ({values_len}):\n  {values}'
                    fail(msg.format(nodeid=nodeid, values=param.values, names=argnames, names_len=len(argnames), values_len=len(param.values)), pytrace=False)
        else:
            mark = get_empty_parameterset_mark(config, argnames, func)
            parameters.append(ParameterSet(values=(NOTSET,) * len(argnames), marks=[mark], id='NOTSET'))
        return (argnames, parameters)

@final
@dataclasses.dataclass(frozen=True)
class Mark:
    name: str
    args: tuple[Any, ...]
    kwargs: Mapping[str, Any]
    _param_ids_from: Mark | None = dataclasses.field(default=None, repr=False)
    _param_ids_generated: Sequence[str] | None = dataclasses.field(default=None, repr=False)

    def __init__(self, name: str, args: tuple[Any, ...], kwargs: Mapping[str, Any], param_ids_from: Mark | None=None, param_ids_generated: Sequence[str] | None=None, *, _ispytest: bool=False) -> None:
        check_ispytest(_ispytest)
        object.__setattr__(self, 'name', name)
        object.__setattr__(self, 'args', args)
        object.__setattr__(self, 'kwargs', kwargs)
        object.__setattr__(self, '_param_ids_from', param_ids_from)
        object.__setattr__(self, '_param_ids_generated', param_ids_generated)

    def _has_param_ids(self) -> bool:
        return 'ids' in self.kwargs or len(self.args) >= 4

    def combined_with(self, other: Mark) -> Mark:
        assert self.name == other.name
        param_ids_from: Mark | None = None
        if self.name == 'parametrize':
            if other._has_param_ids():
                param_ids_from = other
            elif self._has_param_ids():
                param_ids_from = self
        return Mark(self.name, self.args + other.args, dict(self.kwargs, **other.kwargs), param_ids_from=param_ids_from, _ispytest=True)
Markable = TypeVar('Markable', bound=Callable[..., object] | type)

@dataclasses.dataclass
class MarkDecorator:
    mark: Mark

    def __init__(self, mark: Mark, *, _ispytest: bool=False) -> None:
        check_ispytest(_ispytest)
        self.mark = mark

    @property
    def name(self) -> str:
        return self.mark.name

    @property
    def args(self) -> tuple[Any, ...]:
        return self.mark.args

    @property
    def kwargs(self) -> Mapping[str, Any]:
        return self.mark.kwargs

    @property
    def markname(self) -> str:
        return self.name

    def with_args(self, *args: object, **kwargs: object) -> MarkDecorator:
        mark = Mark(self.name, args, kwargs, _ispytest=True)
        return MarkDecorator(self.mark.combined_with(mark), _ispytest=True)

    @overload
    def __call__(self, arg: Markable) -> Markable:
        pass

    @overload
    def __call__(self, *args: object, **kwargs: object) -> MarkDecorator:
        pass

    def __call__(self, *args: object, **kwargs: object):
        if args and (not kwargs):
            func = args[0]
            is_class = inspect.isclass(func)
            unwrapped_func = func
            if isinstance(func, staticmethod | classmethod):
                unwrapped_func = func.__func__
            if len(args) == 1 and (istestfunc(unwrapped_func) or is_class):
                store_mark(unwrapped_func, self.mark, stacklevel=3)
                return func
        return self.with_args(*args, **kwargs)

def get_unpacked_marks(obj: object | type, *, consider_mro: bool=True) -> list[Mark]:
    if isinstance(obj, type):
        if not consider_mro:
            mark_lists = [obj.__dict__.get('pytestmark', [])]
        else:
            mark_lists = [x.__dict__.get('pytestmark', []) for x in reversed(obj.__mro__)]
        mark_list = []
        for item in mark_lists:
            if isinstance(item, list):
                mark_list.extend(item)
            else:
                mark_list.append(item)
    else:
        mark_attribute = getattr(obj, 'pytestmark', [])
        if isinstance(mark_attribute, list):
            mark_list = mark_attribute
        else:
            mark_list = [mark_attribute]
    return list(normalize_mark_list(mark_list))

def normalize_mark_list(mark_list: Iterable[Mark | MarkDecorator]) -> Iterable[Mark]:
    for mark in mark_list:
        mark_obj = getattr(mark, 'mark', mark)
        if not isinstance(mark_obj, Mark):
            raise TypeError(f'got {mark_obj!r} instead of Mark')
        yield mark_obj

def store_mark(obj, mark: Mark, *, stacklevel: int=2) -> None:
    assert isinstance(mark, Mark), mark
    from ..fixtures import getfixturemarker
    if getfixturemarker(obj) is not None:
        warnings.warn(MARKED_FIXTURE, stacklevel=stacklevel)
    obj.pytestmark = [*get_unpacked_marks(obj, consider_mro=False), mark]
if TYPE_CHECKING:

    class _SkipMarkDecorator(MarkDecorator):

        @overload
        def __call__(self, arg: Markable) -> Markable:
            ...

        @overload
        def __call__(self, reason: str=...) -> MarkDecorator:
            ...

    class _SkipifMarkDecorator(MarkDecorator):

        def __call__(self, condition: str | bool=..., *conditions: str | bool, reason: str=...) -> MarkDecorator:
            ...

    class _XfailMarkDecorator(MarkDecorator):

        @overload
        def __call__(self, arg: Markable) -> Markable:
            ...

        @overload
        def __call__(self, condition: str | bool=False, *conditions: str | bool, reason: str=..., run: bool=..., raises: None | type[BaseException] | tuple[type[BaseException], ...] | AbstractRaises[BaseException]=..., strict: bool=...) -> MarkDecorator:
            ...

    class _ParametrizeMarkDecorator(MarkDecorator):

        def __call__(self, argnames: str | Sequence[str], argvalues: Iterable[ParameterSet | Sequence[object] | object], *, indirect: bool | Sequence[str]=..., ids: Iterable[None | str | float | int | bool] | Callable[[Any], object | None] | None=..., scope: _ScopeName | None=...) -> MarkDecorator:
            ...

    class _UsefixturesMarkDecorator(MarkDecorator):

        def __call__(self, *fixtures: str) -> MarkDecorator:
            ...

    class _FilterwarningsMarkDecorator(MarkDecorator):

        def __call__(self, *filters: str) -> MarkDecorator:
            ...

@final
class MarkGenerator:
    if TYPE_CHECKING:
        skip: _SkipMarkDecorator
        skipif: _SkipifMarkDecorator
        xfail: _XfailMarkDecorator
        parametrize: _ParametrizeMarkDecorator
        usefixtures: _UsefixturesMarkDecorator
        filterwarnings: _FilterwarningsMarkDecorator

    def __init__(self, *, _ispytest: bool=False) -> None:
        check_ispytest(_ispytest)
        self._config: Config | None = None
        self._markers: set[str] = set()

    def __getattr__(self, name: str) -> MarkDecorator:
        if name[0] == '_':
            raise AttributeError('Marker name must NOT start with underscore')
        if self._config is not None:
            if name not in self._markers:
                for line in self._config.getini('markers'):
                    marker = line.split(':')[0].split('(')[0].strip()
                    self._markers.add(marker)
            if name not in self._markers:
                if name in ['parameterize', 'parametrise', 'parameterise']:
                    __tracebackhide__ = True
                    fail(f"Unknown '{name}' mark, did you mean 'parametrize'?")
                strict_markers = self._config.getini('strict_markers')
                if strict_markers is None:
                    strict_markers = self._config.getini('strict')
                if strict_markers:
                    fail(f'{name!r} not found in `markers` configuration option', pytrace=False)
                warnings.warn(f'Unknown pytest.mark.{name} - is this a typo?  You can register custom marks to avoid this warning - for details, see https://docs.pytest.org/en/stable/how-to/mark.html', PytestUnknownMarkWarning, 2)
        return MarkDecorator(Mark(name, (), {}, _ispytest=True), _ispytest=True)
MARK_GEN = MarkGenerator(_ispytest=True)

@final
class NodeKeywords(MutableMapping[str, Any]):
    __slots__ = ('_markers', 'node', 'parent')

    def __init__(self, node: Node) -> None:
        self.node = node
        self.parent = node.parent
        self._markers = {node.name: True}

    def __getitem__(self, key: str) -> Any:
        try:
            return self._markers[key]
        except KeyError:
            if self.parent is None:
                raise
            return self.parent.keywords[key]

    def __setitem__(self, key: str, value: Any) -> None:
        self._markers[key] = value

    def __contains__(self, key: object) -> bool:
        return key in self._markers or (self.parent is not None and key in self.parent.keywords)

    def update(self, other: Mapping[str, Any] | Iterable[tuple[str, Any]]=(), **kwds: Any) -> None:
        self._markers.update(other)
        self._markers.update(kwds)

    def __delitem__(self, key: str) -> None:
        raise ValueError('cannot delete key in keywords dict')

    def __iter__(self) -> Iterator[str]:
        yield from self._markers
        if self.parent is not None:
            for keyword in self.parent.keywords:
                if keyword not in self._markers:
                    yield keyword

    def __len__(self) -> int:
        return sum((1 for keyword in self))

    def __repr__(self) -> str:
        return f'<NodeKeywords for node {self.node}>'
