from __future__ import annotations
import functools
from typing import TYPE_CHECKING, Any
import warnings
from pandas._libs.tslibs import BaseOffset, Period, to_offset
from pandas._libs.tslibs.dtypes import OFFSET_TO_PERIOD_FREQSTR, FreqGroup
from pandas.core.dtypes.generic import ABCDatetimeIndex, ABCPeriodIndex, ABCTimedeltaIndex
from pandas.io.formats.printing import pprint_thing
from pandas.plotting._matplotlib.converter import TimeSeries_DateFormatter, TimeSeries_DateLocator, TimeSeries_TimedeltaFormatter
from pandas.tseries.frequencies import get_period_alias, is_subperiod, is_superperiod
if TYPE_CHECKING:
    from datetime import timedelta
    from matplotlib.axes import Axes
    from pandas._typing import NDFrameT
    from pandas import DatetimeIndex, Index, PeriodIndex, Series

def maybe_resample(series: Series, ax: Axes, kwargs: dict[str, Any]):
    if 'how' in kwargs:
        raise ValueError("'how' is not a valid keyword for plotting functions. If plotting multiple objects on shared axes, resample manually first.")
    freq, ax_freq = _get_freq(ax, series)
    if freq is None:
        raise ValueError('Cannot use dynamic axis without frequency info')
    if isinstance(series.index, ABCDatetimeIndex):
        series = series.to_period(freq=freq)
    if ax_freq is not None and freq != ax_freq:
        if is_superperiod(freq, ax_freq):
            series = series.copy(deep=False)
            series.index = series.index.asfreq(ax_freq, how='s')
            freq = ax_freq
        elif _is_sup(freq, ax_freq):
            how = 'last'
            series = getattr(series.resample('D'), how)().dropna()
            series = getattr(series.resample(ax_freq), how)().dropna()
            freq = ax_freq
        elif is_subperiod(freq, ax_freq) or _is_sub(freq, ax_freq):
            _upsample_others(ax, freq, kwargs)
        else:
            raise ValueError('Incompatible frequency conversion')
    return (freq, series)

def _is_sub(f1: str, f2: str) -> bool:
    return f1.startswith('W') and is_subperiod('D', f2) or (f2.startswith('W') and is_subperiod(f1, 'D'))

def _is_sup(f1: str, f2: str) -> bool:
    return f1.startswith('W') and is_superperiod('D', f2) or (f2.startswith('W') and is_superperiod(f1, 'D'))

def _upsample_others(ax: Axes, freq: BaseOffset, kwargs: dict[str, Any]) -> None:
    legend = ax.get_legend()
    lines, labels = _replot_ax(ax, freq)
    _replot_ax(ax, freq)
    other_ax = None
    if hasattr(ax, 'left_ax'):
        other_ax = ax.left_ax
    if hasattr(ax, 'right_ax'):
        other_ax = ax.right_ax
    if other_ax is not None:
        rlines, rlabels = _replot_ax(other_ax, freq)
        lines.extend(rlines)
        labels.extend(rlabels)
    if legend is not None and kwargs.get('legend', True) and (len(lines) > 0):
        title: str | None = legend.get_title().get_text()
        if title == 'None':
            title = None
        ax.legend(lines, labels, loc='best', title=title)

def _replot_ax(ax: Axes, freq: BaseOffset):
    data = getattr(ax, '_plot_data', None)
    ax._plot_data = []
    ax.clear()
    decorate_axes(ax, freq)
    lines = []
    labels = []
    if data is not None:
        for series, plotf, kwds in data:
            series = series.copy(deep=False)
            idx = series.index.asfreq(freq, how='S')
            series.index = idx
            ax._plot_data.append((series, plotf, kwds))
            if isinstance(plotf, str):
                from pandas.plotting._matplotlib import PLOT_CLASSES
                plotf = PLOT_CLASSES[plotf]._plot
            lines.append(plotf(ax, series.index._mpl_repr(), series.values, **kwds)[0])
            labels.append(pprint_thing(series.name))
    return (lines, labels)

def decorate_axes(ax: Axes, freq: BaseOffset) -> None:
    if not hasattr(ax, '_plot_data'):
        ax._plot_data = []
    ax.freq = freq
    xaxis = ax.get_xaxis()
    xaxis.freq = freq

def _get_ax_freq(ax: Axes):
    ax_freq = getattr(ax, 'freq', None)
    if ax_freq is None:
        if hasattr(ax, 'left_ax'):
            ax_freq = getattr(ax.left_ax, 'freq', None)
        elif hasattr(ax, 'right_ax'):
            ax_freq = getattr(ax.right_ax, 'freq', None)
    if ax_freq is None:
        shared_axes = ax.get_shared_x_axes().get_siblings(ax)
        if len(shared_axes) > 1:
            for shared_ax in shared_axes:
                ax_freq = getattr(shared_ax, 'freq', None)
                if ax_freq is not None:
                    break
    return ax_freq

def _get_period_alias(freq: timedelta | BaseOffset | str) -> str | None:
    if isinstance(freq, BaseOffset):
        freqstr = freq.name
    else:
        freqstr = to_offset(freq, is_period=True).rule_code
    return get_period_alias(freqstr)

def _get_freq(ax: Axes, series: Series):
    freq = getattr(series.index, 'freq', None)
    if freq is None:
        freq = getattr(series.index, 'inferred_freq', None)
        freq = to_offset(freq, is_period=True)
    ax_freq = _get_ax_freq(ax)
    if freq is None:
        freq = ax_freq
    freq = _get_period_alias(freq)
    return (freq, ax_freq)

def use_dynamic_x(ax: Axes, index: Index) -> bool:
    freq = _get_index_freq(index)
    ax_freq = _get_ax_freq(ax)
    if freq is None:
        freq = ax_freq
    elif ax_freq is None and len(ax.get_lines()) > 0:
        return False
    if freq is None:
        return False
    freq_str = _get_period_alias(freq)
    if freq_str is None:
        return False
    if isinstance(index, ABCDatetimeIndex):
        freq_str = OFFSET_TO_PERIOD_FREQSTR.get(freq_str, freq_str)
        base = to_offset(freq_str, is_period=True)._period_dtype_code
        if base <= FreqGroup.FR_DAY.value:
            return index[:1].is_normalized
        period = Period(index[0], freq_str)
        assert isinstance(period, Period)
        return period.to_timestamp().tz_localize(index.tz) == index[0]
    return True

def _get_index_freq(index: Index) -> BaseOffset | None:
    freq = getattr(index, 'freq', None)
    if freq is None:
        freq = getattr(index, 'inferred_freq', None)
        freq = to_offset(freq)
    return freq

def maybe_convert_index(ax: Axes, data: NDFrameT) -> NDFrameT:
    if isinstance(data.index, (ABCDatetimeIndex, ABCPeriodIndex)):
        freq = _get_index_freq(data.index)
        if freq is None:
            freq = _get_ax_freq(ax)
        if freq is None:
            raise ValueError('Could not get frequency alias for plotting')
        freq_str = _get_period_alias(freq)
        with warnings.catch_warnings():
            warnings.filterwarnings('ignore', 'PeriodDtype\\[B\\] is deprecated', category=FutureWarning)
            if isinstance(data.index, ABCDatetimeIndex):
                data = data.tz_localize(None).to_period(freq=freq_str)
            elif isinstance(data.index, ABCPeriodIndex):
                data.index = data.index.asfreq(freq=freq_str, how='start')
    return data

def _format_coord(freq: BaseOffset, t, y) -> str:
    time_period = Period(ordinal=int(t), freq=freq)
    return f't = {time_period}  y = {y:8f}'

def format_dateaxis(subplot, freq: BaseOffset, index: DatetimeIndex | PeriodIndex) -> None:
    import matplotlib.pyplot as plt
    if isinstance(index, ABCPeriodIndex):
        majlocator = TimeSeries_DateLocator(freq, dynamic_mode=True, minor_locator=False, plot_obj=subplot)
        minlocator = TimeSeries_DateLocator(freq, dynamic_mode=True, minor_locator=True, plot_obj=subplot)
        subplot.xaxis.set_major_locator(majlocator)
        subplot.xaxis.set_minor_locator(minlocator)
        majformatter = TimeSeries_DateFormatter(freq, dynamic_mode=True, minor_locator=False, plot_obj=subplot)
        minformatter = TimeSeries_DateFormatter(freq, dynamic_mode=True, minor_locator=True, plot_obj=subplot)
        subplot.xaxis.set_major_formatter(majformatter)
        subplot.xaxis.set_minor_formatter(minformatter)
        subplot.format_coord = functools.partial(_format_coord, freq)
    elif isinstance(index, ABCTimedeltaIndex):
        subplot.xaxis.set_major_formatter(TimeSeries_TimedeltaFormatter(index.unit))
    else:
        raise TypeError('index type not supported')
    plt.draw_if_interactive()

def prepare_ts_data(series: Series, ax: Axes, kwargs: dict[str, Any]) -> tuple[BaseOffset | str, Series]:
    freq, data = maybe_resample(series, ax, kwargs)
    decorate_axes(ax, freq)
    if hasattr(ax, 'left_ax'):
        decorate_axes(ax.left_ax, freq)
    if hasattr(ax, 'right_ax'):
        decorate_axes(ax.right_ax, freq)
    return (freq, data)
