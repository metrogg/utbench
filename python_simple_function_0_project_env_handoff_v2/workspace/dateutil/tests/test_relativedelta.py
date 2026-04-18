import pytest
from dateutil.relativedelta import relativedelta, MO
from datetime import datetime


def test_relativedelta_basic():
    dt = datetime(2023, 1, 15)
    delta = relativedelta(months=1)
    result = dt + delta
    assert result.month == 2


def test_relativedelta_with_weekday():
    dt = datetime(2023, 1, 15)
    delta = relativedelta(days=1, weekday=MO(1))
    result = dt + delta
    assert result.weekday() == 0


def test_relativedelta_years():
    dt = datetime(2023, 1, 15)
    delta = relativedelta(years=2)
    result = dt + delta
    assert result.year == 2025