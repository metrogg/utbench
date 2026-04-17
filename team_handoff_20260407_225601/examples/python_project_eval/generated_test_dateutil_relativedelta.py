```python
import datetime as dt

import pytest

from dateutil.relativedelta import FR, MO, relativedelta


def test_relativedelta_diff_between_dates_keeps_month_and_day_components():
    left = dt.datetime(2024, 3, 15, 10, 30, 5)
    right = dt.datetime(2024, 1, 10, 8, 15, 2)

    delta = relativedelta(left, right)

    assert delta.months == 2
    assert delta.days == 5
    assert delta.hours == 2
    assert delta.minutes == 15
    assert delta.seconds == 3


def test_relativedelta_adds_relative_offsets_to_datetime():
    base = dt.datetime(2024, 1, 31, 23, 0, 0)

    result = base + relativedelta(months=1, days=1, hours=2)

    assert result == dt.datetime(2024, 3, 1, 1, 0, 0)


def test_relativedelta_absolute_fields_replace_original_components():
    base = dt.datetime(2024, 8, 20, 14, 45, 10)

    result = base + relativedelta(year=2025, month=2, day=5, hour=9, minute=30, second=0)

    assert result == dt.datetime(2025, 2, 5, 9, 30, 0)


def test_relativedelta_weekday_moves_forward_to_requested_weekday():
    base = dt.datetime(2024, 1, 1, 9, 0, 0)  # Tuesday

    result = base + relativedelta(weekday=FR(+1))

    assert result == dt.datetime(2024, 1, 5, 9, 0, 0)


def test_relativedelta_weekday_can_target_previous_or_same_weekday():
    base = dt.datetime(2024, 1, 3, 9, 0, 0)  # Wednesday

    result = base + relativedelta(weekday=MO(-1))

    assert result == dt.datetime(2024, 1, 1, 9, 0, 0)


def test_relativedelta_rejects_non_integer_months_and_years():
    with pytest.raises(ValueError):
        relativedelta(months=1.5)

    with pytest.raises(ValueError):
        relativedelta(years=2.25)
```
