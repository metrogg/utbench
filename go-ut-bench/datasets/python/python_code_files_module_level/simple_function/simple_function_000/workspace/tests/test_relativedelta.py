import pytest
from dateutil.relativedelta import relativedelta, MO
from datetime import datetime


def test_relativedelta_months():
    dt = datetime(2023, 1, 15)
    delta = relativedelta(months=1)
    result = dt + delta
    assert result.month == 2
    assert result.day == 15


def test_relativedelta_years():
    dt = datetime(2023, 1, 15)
    delta = relativedelta(years=2)
    result = dt + delta
    assert result.year == 2025


def test_relativedelta_days():
    dt = datetime(2023, 1, 15)
    delta = relativedelta(days=10)
    result = dt + delta
    assert result.day == 25


def test_relativedelta_weekday():
    dt = datetime(2023, 1, 15)
    delta = relativedelta(weekday=MO)
    result = dt + delta
    assert result.weekday() == 0


def test_relativedelta_negative():
    dt = datetime(2023, 3, 15)
    delta = relativedelta(months=-1)
    result = dt + delta
    assert result.month == 2


def test_relativedelta_addition():
    delta1 = relativedelta(years=1, months=2)
    delta2 = relativedelta(months=3)
    result = delta1 + delta2
    assert result.years == 1
    assert result.months == 5


def test_relativedelta_multiplication():
    delta = relativedelta(days=5)
    result = delta * 2
    assert result.days == 10


def test_relativedelta_equality():
    delta1 = relativedelta(years=1, months=2)
    delta2 = relativedelta(years=1, months=2)
    delta3 = relativedelta(years=1, months=3)
    assert delta1 == delta2
    assert delta1 != delta3
