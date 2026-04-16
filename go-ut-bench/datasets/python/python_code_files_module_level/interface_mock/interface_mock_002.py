from __future__ import annotations
import importlib
from typing import TYPE_CHECKING, Literal
from pandas._config import get_option
from pandas.util._decorators import set_module
from pandas.core.dtypes.common import is_integer, is_list_like
from pandas.core.dtypes.generic import ABCDataFrame, ABCSeries
from pandas.core.base import PandasObject
if TYPE_CHECKING:
    from collections.abc import Callable, Hashable, Sequence
    import types
    from matplotlib.axes import Axes
    import numpy as np
    from pandas._typing import IndexLabel
    from pandas import DataFrame, Index, Series
    from pandas.core.groupby.generic import DataFrameGroupBy

def holds_integer(column: Index) -> bool:
    return column.dtype.kind in 'iu'

@set_module('pandas.plotting')
def hist_series(self: Series, by=None, ax=None, grid: bool=True, xlabelsize: int | None=None, xrot: float | None=None, ylabelsize: int | None=None, yrot: float | None=None, figsize: tuple[int, int] | None=None, bins: int | Sequence[int]=10, backend: str | None=None, legend: bool=False, **kwargs):
    plot_backend = _get_plot_backend(backend)
    return plot_backend.hist_series(self, by=by, ax=ax, grid=grid, xlabelsize=xlabelsize, xrot=xrot, ylabelsize=ylabelsize, yrot=yrot, figsize=figsize, bins=bins, legend=legend, **kwargs)

@set_module('pandas.plotting')
def hist_frame(data: DataFrame, column: IndexLabel | None=None, by=None, grid: bool=True, xlabelsize: int | None=None, xrot: float | None=None, ylabelsize: int | None=None, yrot: float | None=None, ax=None, sharex: bool=False, sharey: bool=False, figsize: tuple[int, int] | None=None, layout: tuple[int, int] | None=None, bins: int | Sequence[int]=10, backend: str | None=None, legend: bool=False, **kwargs):
    plot_backend = _get_plot_backend(backend)
    return plot_backend.hist_frame(data, column=column, by=by, grid=grid, xlabelsize=xlabelsize, xrot=xrot, ylabelsize=ylabelsize, yrot=yrot, ax=ax, sharex=sharex, sharey=sharey, figsize=figsize, layout=layout, legend=legend, bins=bins, **kwargs)

@set_module('pandas.plotting')
def boxplot(data: DataFrame, column: str | list[str] | None=None, by: str | list[str] | None=None, ax: Axes | None=None, fontsize: float | str | None=None, rot: int=0, grid: bool=True, figsize: tuple[float, float] | None=None, layout: tuple[int, int] | None=None, return_type: str | None=None, **kwargs):
    plot_backend = _get_plot_backend('matplotlib')
    return plot_backend.boxplot(data, column=column, by=by, ax=ax, fontsize=fontsize, rot=rot, grid=grid, figsize=figsize, layout=layout, return_type=return_type, **kwargs)

@set_module('pandas.plotting')
def boxplot_frame(self: DataFrame, column=None, by=None, ax=None, fontsize: int | None=None, rot: int=0, grid: bool=True, figsize: tuple[float, float] | None=None, layout=None, return_type=None, backend=None, **kwargs):
    plot_backend = _get_plot_backend(backend)
    return plot_backend.boxplot_frame(self, column=column, by=by, ax=ax, fontsize=fontsize, rot=rot, grid=grid, figsize=figsize, layout=layout, return_type=return_type, **kwargs)

@set_module('pandas.plotting')
def boxplot_frame_groupby(grouped: DataFrameGroupBy, subplots: bool=True, column=None, fontsize: int | None=None, rot: int=0, grid: bool=True, ax=None, figsize: tuple[float, float] | None=None, layout=None, sharex: bool=False, sharey: bool=True, backend=None, **kwargs):
    plot_backend = _get_plot_backend(backend)
    return plot_backend.boxplot_frame_groupby(grouped, subplots=subplots, column=column, fontsize=fontsize, rot=rot, grid=grid, ax=ax, figsize=figsize, layout=layout, sharex=sharex, sharey=sharey, **kwargs)

@set_module('pandas.plotting')
class PlotAccessor(PandasObject):
    _common_kinds = ('line', 'bar', 'barh', 'kde', 'density', 'area', 'hist', 'box')
    _series_kinds = ('pie',)
    _dataframe_kinds = ('scatter', 'hexbin')
    _kind_aliases = {'density': 'kde'}
    _all_kinds = _common_kinds + _series_kinds + _dataframe_kinds

    def __init__(self, data: Series | DataFrame) -> None:
        self._parent = data

    @staticmethod
    def _get_call_args(backend_name: str, data: Series | DataFrame, args, kwargs):
        if isinstance(data, ABCSeries):
            arg_def = [('kind', 'line'), ('ax', None), ('figsize', None), ('use_index', True), ('title', None), ('grid', None), ('legend', False), ('style', None), ('logx', False), ('logy', False), ('loglog', False), ('xticks', None), ('yticks', None), ('xlim', None), ('ylim', None), ('rot', None), ('fontsize', None), ('colormap', None), ('table', False), ('yerr', None), ('xerr', None), ('label', None), ('secondary_y', False), ('xlabel', None), ('ylabel', None)]
        elif isinstance(data, ABCDataFrame):
            arg_def = [('x', None), ('y', None), ('kind', 'line'), ('ax', None), ('subplots', False), ('sharex', None), ('sharey', False), ('layout', None), ('figsize', None), ('use_index', True), ('title', None), ('grid', None), ('legend', True), ('style', None), ('logx', False), ('logy', False), ('loglog', False), ('xticks', None), ('yticks', None), ('xlim', None), ('ylim', None), ('rot', None), ('fontsize', None), ('colormap', None), ('table', False), ('yerr', None), ('xerr', None), ('secondary_y', False), ('xlabel', None), ('ylabel', None)]
        else:
            raise TypeError(f'Called plot accessor for type {type(data).__name__}, expected Series or DataFrame')
        if args and isinstance(data, ABCSeries):
            positional_args = str(args)[1:-1]
            keyword_args = ', '.join([f'{name}={value!r}' for (name, _), value in zip(arg_def, args, strict=False)])
            msg = f'`Series.plot()` should not be called with positional arguments, only keyword arguments. The order of positional arguments will change in the future. Use `Series.plot({keyword_args})` instead of `Series.plot({positional_args})`.'
            raise TypeError(msg)
        pos_args = {name: value for (name, _), value in zip(arg_def, args, strict=False)}
        if backend_name == 'pandas.plotting._matplotlib':
            kwargs = dict(arg_def, **pos_args, **kwargs)
        else:
            kwargs = dict(pos_args, **kwargs)
        x = kwargs.pop('x', None)
        y = kwargs.pop('y', None)
        kind = kwargs.pop('kind', 'line')
        return (x, y, kind, kwargs)

    def __call__(self, *args, **kwargs):
        plot_backend = _get_plot_backend(kwargs.pop('backend', None))
        x, y, kind, kwargs = self._get_call_args(plot_backend.__name__, self._parent, args, kwargs)
        kind = self._kind_aliases.get(kind, kind)
        if plot_backend.__name__ != 'pandas.plotting._matplotlib':
            return plot_backend.plot(self._parent, x=x, y=y, kind=kind, **kwargs)
        if kind not in self._all_kinds:
            raise ValueError(f'{kind} is not a valid plot kind Valid plot kinds: {self._all_kinds}')
        data = self._parent
        if isinstance(data, ABCSeries):
            kwargs['reuse_plot'] = True
        if kind in self._dataframe_kinds:
            if isinstance(data, ABCDataFrame):
                return plot_backend.plot(data, x=x, y=y, kind=kind, **kwargs)
            else:
                raise ValueError(f'plot kind {kind} can only be used for data frames')
        elif kind in self._series_kinds:
            if isinstance(data, ABCDataFrame):
                if y is None and kwargs.get('subplots') is False:
                    raise ValueError(f"{kind} requires either y column or 'subplots=True'")
                if y is not None:
                    if is_integer(y) and (not holds_integer(data.columns)):
                        y = data.columns[y]
                    data = data[y].copy(deep=False)
                    data.index.name = y
        elif isinstance(data, ABCDataFrame):
            data_cols = data.columns
            if x is not None:
                if is_integer(x) and (not holds_integer(data.columns)):
                    x = data_cols[x]
                elif not isinstance(data[x], ABCSeries):
                    raise ValueError('x must be a label or position')
                data = data.set_index(x)
            if y is not None:
                int_ylist = is_list_like(y) and all((is_integer(c) for c in y))
                int_y_arg = is_integer(y) or int_ylist
                if int_y_arg and (not holds_integer(data.columns)):
                    y = data_cols[y]
                label_kw = kwargs['label'] if 'label' in kwargs else False
                for kw in ['xerr', 'yerr']:
                    if kw in kwargs and (isinstance(kwargs[kw], str) or is_integer(kwargs[kw])):
                        try:
                            kwargs[kw] = data[kwargs[kw]]
                        except (IndexError, KeyError, TypeError):
                            pass
                data = data[y]
                if isinstance(data, ABCSeries):
                    label_name = label_kw or y
                    data.name = label_name
                else:
                    match = is_list_like(label_kw) and len(label_kw) == len(y)
                    if label_kw and (not match):
                        raise ValueError('label should be list-like and same length as y')
                    label_name = label_kw or data.columns
                    data.columns = label_name
        return plot_backend.plot(data, kind=kind, **kwargs)
    __call__.__doc__ = __doc__

    def line(self, x: Hashable | None=None, y: Hashable | None=None, color: str | Sequence[str] | dict | None=None, **kwargs) -> PlotAccessor:
        if color is not None:
            kwargs['color'] = color
        return self(kind='line', x=x, y=y, **kwargs)

    def bar(self, x: Hashable | None=None, y: Hashable | None=None, color: str | Sequence[str] | dict | None=None, **kwargs) -> PlotAccessor:
        if color is not None:
            kwargs['color'] = color
        return self(kind='bar', x=x, y=y, **kwargs)

    def barh(self, x: Hashable | None=None, y: Hashable | None=None, color: str | Sequence[str] | dict | None=None, **kwargs) -> PlotAccessor:
        if color is not None:
            kwargs['color'] = color
        return self(kind='barh', x=x, y=y, **kwargs)

    def box(self, by: IndexLabel | None=None, **kwargs) -> PlotAccessor:
        return self(kind='box', by=by, **kwargs)

    def hist(self, by: IndexLabel | None=None, bins: int=10, **kwargs) -> PlotAccessor:
        return self(kind='hist', by=by, bins=bins, **kwargs)

    def kde(self, bw_method: Literal['scott', 'silverman'] | float | Callable | None=None, ind: np.ndarray | int | None=None, weights: np.ndarray | None=None, **kwargs) -> PlotAccessor:
        return self(kind='kde', bw_method=bw_method, ind=ind, weights=weights, **kwargs)
    density = kde

    def area(self, x: Hashable | None=None, y: Hashable | None=None, stacked: bool=True, **kwargs) -> PlotAccessor:
        return self(kind='area', x=x, y=y, stacked=stacked, **kwargs)

    def pie(self, y: IndexLabel | None=None, **kwargs) -> PlotAccessor:
        if y is not None:
            kwargs['y'] = y
        if isinstance(self._parent, ABCDataFrame) and kwargs.get('y', None) is None and (not kwargs.get('subplots', False)):
            raise ValueError("pie requires either y column or 'subplots=True'")
        return self(kind='pie', **kwargs)

    def scatter(self, x: Hashable, y: Hashable, s: Hashable | Sequence[Hashable] | None=None, c: Hashable | Sequence[Hashable] | None=None, **kwargs) -> PlotAccessor:
        return self(kind='scatter', x=x, y=y, s=s, c=c, **kwargs)

    def hexbin(self, x: Hashable, y: Hashable, C: Hashable | None=None, reduce_C_function: Callable | None=None, gridsize: int | tuple[int, int] | None=None, **kwargs) -> PlotAccessor:
        if reduce_C_function is not None:
            kwargs['reduce_C_function'] = reduce_C_function
        if gridsize is not None:
            kwargs['gridsize'] = gridsize
        return self(kind='hexbin', x=x, y=y, C=C, **kwargs)
_backends: dict[str, types.ModuleType] = {}

def _load_backend(backend: str) -> types.ModuleType:
    from importlib.metadata import entry_points
    if backend == 'matplotlib':
        try:
            module = importlib.import_module('pandas.plotting._matplotlib')
        except ImportError:
            raise ImportError('matplotlib is required for plotting when the default backend "matplotlib" is selected.') from None
        return module
    found_backend = False
    eps = entry_points()
    key = 'pandas_plotting_backends'
    if hasattr(eps, 'select'):
        entry = eps.select(group=key)
    else:
        entry = eps.get(key, ())
    for entry_point in entry:
        found_backend = entry_point.name == backend
        if found_backend:
            module = entry_point.load()
            break
    if not found_backend:
        try:
            module = importlib.import_module(backend)
            found_backend = True
        except ImportError:
            pass
    if found_backend:
        if hasattr(module, 'plot'):
            return module
    raise ValueError(f"Could not find plotting backend '{backend}'. Ensure that you've installed the package providing the '{backend}' entrypoint, or that the package has a top-level `.plot` method.")

def _get_plot_backend(backend: str | None=None):
    backend_str: str = backend or get_option('plotting.backend')
    if backend_str in _backends:
        return _backends[backend_str]
    module = _load_backend(backend_str)
    _backends[backend_str] = module
    return module
