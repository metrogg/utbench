from __future__ import annotations
import datetime
from functools import partial
from typing import TYPE_CHECKING, cast
import numpy as np
from pandas._libs.tslibs import Timedelta
import pandas._libs.window.aggregations as window_aggregations
from pandas.util._decorators import set_module
from pandas.core.dtypes.common import is_datetime64_dtype, is_numeric_dtype
from pandas.core.dtypes.dtypes import DatetimeTZDtype
from pandas.core.dtypes.generic import ABCSeries
from pandas.core.dtypes.missing import isna
from pandas.core import common
from pandas.core.arrays.datetimelike import dtype_to_unit
from pandas.core.indexers.objects import BaseIndexer, ExponentialMovingWindowIndexer, GroupbyIndexer
from pandas.core.util.numba_ import get_jit_arguments, maybe_use_numba
from pandas.core.window.common import zsqrt
from pandas.core.window.numba_ import generate_numba_ewm_func, generate_numba_ewm_table_func
from pandas.core.window.online import EWMMeanState, generate_online_numba_ewma_func
from pandas.core.window.rolling import BaseWindow, BaseWindowGroupby
if TYPE_CHECKING:
    from pandas._typing import TimedeltaConvertibleTypes, TimeUnit, npt
    from pandas import DataFrame, Series
    from pandas.core.generic import NDFrame

def get_center_of_mass(comass: float | None, span: float | None, halflife: float | None, alpha: float | None) -> float:
    valid_count = common.count_not_none(comass, span, halflife, alpha)
    if valid_count > 1:
        raise ValueError('comass, span, halflife, and alpha are mutually exclusive')
    if comass is not None:
        if comass < 0:
            raise ValueError('comass must satisfy: comass >= 0')
    elif span is not None:
        if span < 1:
            raise ValueError('span must satisfy: span >= 1')
        comass = (span - 1) / 2
    elif halflife is not None:
        if halflife <= 0:
            raise ValueError('halflife must satisfy: halflife > 0')
        decay = 1 - np.exp(np.log(0.5) / halflife)
        comass = 1 / decay - 1
    elif alpha is not None:
        if alpha <= 0 or alpha > 1:
            raise ValueError('alpha must satisfy: 0 < alpha <= 1')
        comass = (1 - alpha) / alpha
    else:
        raise ValueError('Must pass one of comass, span, halflife, or alpha')
    return float(comass)

def _calculate_deltas(times: np.ndarray | NDFrame, halflife: float | TimedeltaConvertibleTypes | None) -> npt.NDArray[np.float64]:
    unit = dtype_to_unit(times.dtype)
    unit = cast('TimeUnit', unit)
    if isinstance(times, ABCSeries):
        times = times._values
    _times = np.asarray(times.view(np.int64), dtype=np.float64)
    _halflife = float(Timedelta(halflife).as_unit(unit)._value)
    return np.diff(_times) / _halflife

@set_module('pandas.api.typing')
class ExponentialMovingWindow(BaseWindow):
    _attributes = ['com', 'span', 'halflife', 'alpha', 'min_periods', 'adjust', 'ignore_na', 'times', 'method']

    def __init__(self, obj: NDFrame, com: float | None=None, span: float | None=None, halflife: float | TimedeltaConvertibleTypes | None=None, alpha: float | None=None, min_periods: int | None=0, adjust: bool=True, ignore_na: bool=False, times: np.ndarray | NDFrame | None=None, method: str='single', *, selection=None) -> None:
        super().__init__(obj=obj, min_periods=1 if min_periods is None else max(int(min_periods), 1), on=None, center=False, closed=None, method=method, selection=selection)
        self.com = com
        self.span = span
        self.halflife = halflife
        self.alpha = alpha
        self.adjust = adjust
        self.ignore_na = ignore_na
        self.times = times
        if self.times is not None:
            times_dtype = getattr(self.times, 'dtype', None)
            if not (is_datetime64_dtype(times_dtype) or isinstance(times_dtype, DatetimeTZDtype)):
                raise ValueError('times must be datetime64 dtype.')
            if len(self.times) != len(obj):
                raise ValueError('times must be the same length as the object.')
            if not isinstance(self.halflife, (str, datetime.timedelta, np.timedelta64)):
                raise ValueError('halflife must be a timedelta convertible object')
            if isna(self.times).any():
                raise ValueError('Cannot convert NaT values to integer')
            self._deltas = _calculate_deltas(self.times, self.halflife)
            if common.count_not_none(self.com, self.span, self.alpha) > 0:
                if not self.adjust:
                    raise NotImplementedError('None of com, span, or alpha can be specified if times is provided and adjust=False')
                self._com = get_center_of_mass(self.com, self.span, None, self.alpha)
            else:
                self._com = 1.0
        else:
            if self.halflife is not None and isinstance(self.halflife, (str, datetime.timedelta, np.timedelta64)):
                raise ValueError('halflife can only be a timedelta convertible argument if times is not None.')
            self._deltas = np.ones(max(self.obj.shape[0] - 1, 0), dtype=np.float64)
            self._com = get_center_of_mass(self.com, self.span, self.halflife, self.alpha)

    def _check_window_bounds(self, start: np.ndarray, end: np.ndarray, num_vals: int) -> None:
        pass

    def _get_window_indexer(self) -> BaseIndexer:
        return ExponentialMovingWindowIndexer()

    def online(self, engine: str='numba', engine_kwargs=None) -> OnlineExponentialMovingWindow:
        return OnlineExponentialMovingWindow(obj=self.obj, com=self.com, span=self.span, halflife=self.halflife, alpha=self.alpha, min_periods=self.min_periods, adjust=self.adjust, ignore_na=self.ignore_na, times=self.times, engine=engine, engine_kwargs=engine_kwargs, selection=self._selection)

    def aggregate(self, func=None, *args, **kwargs):
        return super().aggregate(func, *args, **kwargs)
    agg = aggregate

    def mean(self, numeric_only: bool=False, engine=None, engine_kwargs=None):
        if maybe_use_numba(engine):
            if self.method == 'single':
                func = generate_numba_ewm_func
            else:
                func = generate_numba_ewm_table_func
            ewm_func = func(**get_jit_arguments(engine_kwargs), com=self._com, adjust=self.adjust, ignore_na=self.ignore_na, deltas=tuple(self._deltas), normalize=True)
            return self._apply(ewm_func, name='mean')
        elif engine in ('cython', None):
            if engine_kwargs is not None:
                raise ValueError('cython engine does not accept engine_kwargs')
            deltas = None if self.times is None else self._deltas
            window_func = partial(window_aggregations.ewm, com=self._com, adjust=self.adjust, ignore_na=self.ignore_na, deltas=deltas, normalize=True)
            return self._apply(window_func, name='mean', numeric_only=numeric_only)
        else:
            raise ValueError("engine must be either 'numba' or 'cython'")

    def sum(self, numeric_only: bool=False, engine=None, engine_kwargs=None):
        if not self.adjust:
            raise NotImplementedError('sum is not implemented with adjust=False')
        if self.times is not None:
            raise NotImplementedError('sum is not implemented with times')
        if maybe_use_numba(engine):
            if self.method == 'single':
                func = generate_numba_ewm_func
            else:
                func = generate_numba_ewm_table_func
            ewm_func = func(**get_jit_arguments(engine_kwargs), com=self._com, adjust=self.adjust, ignore_na=self.ignore_na, deltas=tuple(self._deltas), normalize=False)
            return self._apply(ewm_func, name='sum')
        elif engine in ('cython', None):
            if engine_kwargs is not None:
                raise ValueError('cython engine does not accept engine_kwargs')
            deltas = None if self.times is None else self._deltas
            window_func = partial(window_aggregations.ewm, com=self._com, adjust=self.adjust, ignore_na=self.ignore_na, deltas=deltas, normalize=False)
            return self._apply(window_func, name='sum', numeric_only=numeric_only)
        else:
            raise ValueError("engine must be either 'numba' or 'cython'")

    def std(self, bias: bool=False, numeric_only: bool=False):
        if numeric_only and self._selected_obj.ndim == 1 and (not is_numeric_dtype(self._selected_obj.dtype)):
            raise NotImplementedError(f'{type(self).__name__}.std does not implement numeric_only')
        if self.times is not None:
            raise NotImplementedError('std is not implemented with times')
        return zsqrt(self.var(bias=bias, numeric_only=numeric_only))

    def var(self, bias: bool=False, numeric_only: bool=False):
        if self.times is not None:
            raise NotImplementedError('var is not implemented with times')
        window_func = window_aggregations.ewmcov
        wfunc = partial(window_func, com=self._com, adjust=self.adjust, ignore_na=self.ignore_na, bias=bias)

        def var_func(values, begin, end, min_periods):
            return wfunc(values, begin, end, min_periods, values)
        return self._apply(var_func, name='var', numeric_only=numeric_only)

    def cov(self, other: DataFrame | Series | None=None, pairwise: bool | None=None, bias: bool=False, numeric_only: bool=False):
        if self.times is not None:
            raise NotImplementedError('cov is not implemented with times')
        from pandas import Series
        self._validate_numeric_only('cov', numeric_only)

        def cov_func(x, y):
            x_array = self._prep_values(x)
            y_array = self._prep_values(y)
            window_indexer = self._get_window_indexer()
            min_periods = self.min_periods if self.min_periods is not None else window_indexer.window_size
            start, end = window_indexer.get_window_bounds(num_values=len(x_array), min_periods=min_periods, center=self.center, closed=self.closed, step=self.step)
            result = window_aggregations.ewmcov(x_array, start, end, self.min_periods, y_array, self._com, self.adjust, self.ignore_na, bias)
            return Series(result, index=x.index, name=x.name, copy=False)
        return self._apply_pairwise(self._selected_obj, other, pairwise, cov_func, numeric_only)

    def corr(self, other: DataFrame | Series | None=None, pairwise: bool | None=None, numeric_only: bool=False):
        if self.times is not None:
            raise NotImplementedError('corr is not implemented with times')
        from pandas import Series
        self._validate_numeric_only('corr', numeric_only)

        def cov_func(x, y):
            x_array = self._prep_values(x)
            y_array = self._prep_values(y)
            window_indexer = self._get_window_indexer()
            min_periods = self.min_periods if self.min_periods is not None else window_indexer.window_size
            start, end = window_indexer.get_window_bounds(num_values=len(x_array), min_periods=min_periods, center=self.center, closed=self.closed, step=self.step)

            def _cov(X, Y):
                return window_aggregations.ewmcov(X, start, end, min_periods, Y, self._com, self.adjust, self.ignore_na, True)
            with np.errstate(all='ignore'):
                cov = _cov(x_array, y_array)
                x_var = _cov(x_array, x_array)
                y_var = _cov(y_array, y_array)
                result = cov / zsqrt(x_var * y_var)
            return Series(result, index=x.index, name=x.name, copy=False)
        return self._apply_pairwise(self._selected_obj, other, pairwise, cov_func, numeric_only)

@set_module('pandas.api.typing')
class ExponentialMovingWindowGroupby(BaseWindowGroupby, ExponentialMovingWindow):
    _attributes = ExponentialMovingWindow._attributes + BaseWindowGroupby._attributes

    def __init__(self, obj, *args, _grouper=None, **kwargs) -> None:
        super().__init__(obj, *args, _grouper=_grouper, **kwargs)
        if not obj.empty and self.times is not None:
            groupby_order = np.concatenate(list(self._grouper.indices.values()))
            self._deltas = _calculate_deltas(self.times.take(groupby_order), self.halflife)

    def _get_window_indexer(self) -> GroupbyIndexer:
        window_indexer = GroupbyIndexer(groupby_indices=self._grouper.indices, window_indexer=ExponentialMovingWindowIndexer)
        return window_indexer

class OnlineExponentialMovingWindow(ExponentialMovingWindow):

    def __init__(self, obj: NDFrame, com: float | None=None, span: float | None=None, halflife: float | TimedeltaConvertibleTypes | None=None, alpha: float | None=None, min_periods: int | None=0, adjust: bool=True, ignore_na: bool=False, times: np.ndarray | NDFrame | None=None, engine: str='numba', engine_kwargs: dict[str, bool] | None=None, *, selection=None) -> None:
        if times is not None:
            raise NotImplementedError('times is not implemented with online operations.')
        super().__init__(obj=obj, com=com, span=span, halflife=halflife, alpha=alpha, min_periods=min_periods, adjust=adjust, ignore_na=ignore_na, times=times, selection=selection)
        self._mean = EWMMeanState(self._com, self.adjust, self.ignore_na, obj.shape)
        if maybe_use_numba(engine):
            self.engine = engine
            self.engine_kwargs = engine_kwargs
        else:
            raise ValueError("'numba' is the only supported engine")

    def reset(self) -> None:
        self._mean.reset()

    def aggregate(self, func=None, *args, **kwargs):
        raise NotImplementedError('aggregate is not implemented.')

    def std(self, bias: bool=False, *args, **kwargs):
        raise NotImplementedError('std is not implemented.')

    def corr(self, other: DataFrame | Series | None=None, pairwise: bool | None=None, numeric_only: bool=False):
        raise NotImplementedError('corr is not implemented.')

    def cov(self, other: DataFrame | Series | None=None, pairwise: bool | None=None, bias: bool=False, numeric_only: bool=False):
        raise NotImplementedError('cov is not implemented.')

    def var(self, bias: bool=False, numeric_only: bool=False):
        raise NotImplementedError('var is not implemented.')

    def mean(self, *args, update=None, update_times=None, **kwargs):
        result_kwargs = {}
        is_frame = self._selected_obj.ndim == 2
        if update_times is not None:
            raise NotImplementedError('update_times is not implemented.')
        update_deltas = np.ones(max(self._selected_obj.shape[-1] - 1, 0), dtype=np.float64)
        if update is not None:
            if self._mean.last_ewm is None:
                raise ValueError('Must call mean with update=None first before passing update')
            result_from = 1
            result_kwargs['index'] = update.index
            if is_frame:
                last_value = self._mean.last_ewm[np.newaxis, :]
                result_kwargs['columns'] = update.columns
            else:
                last_value = self._mean.last_ewm
                result_kwargs['name'] = update.name
            np_array = np.concatenate((last_value, update.to_numpy()))
        else:
            result_from = 0
            result_kwargs['index'] = self._selected_obj.index
            if is_frame:
                result_kwargs['columns'] = self._selected_obj.columns
            else:
                result_kwargs['name'] = self._selected_obj.name
            np_array = self._selected_obj.astype(np.float64).to_numpy()
        ewma_func = generate_online_numba_ewma_func(**get_jit_arguments(self.engine_kwargs))
        result = self._mean.run_ewm(np_array if is_frame else np_array[:, np.newaxis], update_deltas, self.min_periods, ewma_func)
        if not is_frame:
            result = result.squeeze()
        result = result[result_from:]
        result = self._selected_obj._constructor(result, **result_kwargs)
        return result
