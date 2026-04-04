import pytest
from typing import List


# Include the function under test for runnability
def below_zero(operations: List[int]) -> bool:
    balance = 0
    for op in operations:
        balance += op
        if balance < 0:
            return True
    return False


def test_empty_operations_returns_false():
    assert below_zero([]) is False


def test_single_positive_operation_returns_false():
    assert below_zero([9]) is False


def test_single_negative_operation_returns_true():
    assert below_zero([-1]) is True


def test_single_zero_operation_returns_false():
    assert below_zero([0]) is False


def test_all_positive_operations_returns_false():
    # Docstring example 1
    assert below_zero([1, 2, 3]) is False


def test_dip_below_zero_mid_sequence_returns_true():
    # Docstring example 2
    assert below_zero([1, 2, -4, 5]) is True


def test_balance_exactly_zero_never_below_returns_false():
    assert below_zero([5, -2, -3, 10, -10]) is False


def test_below_zero_on_last_operation_returns_true():
    assert below_zero([3, 4, -8]) is True


def test_below_zero_on_first_operation_returns_true():
    assert below_zero([-7, 15, 3]) is True


def test_all_zero_operations_returns_false():
    assert below_zero([0, 0, 0, 0]) is False


def test_negative_after_zero_operation_returns_true():
    assert below_zero([0, 0, -1, 6]) is True


def test_large_numbers_never_below_zero_returns_false():
    assert below_zero([1_000_000, -500_000, 400_000, -100_000]) is False


def test_large_numbers_below_zero_returns_true():
    assert below_zero([1_000_000, -2_000_000, 3_000_000]) is True


def test_non_iterable_input_raises_type_error():
    with pytest.raises(TypeError):
        below_zero(12345)


def test_none_input_raises_type_error():
    with pytest.raises(TypeError):
        below_zero(None)


def test_non_numeric_operation_raises_type_error():
    with pytest.raises(TypeError):
        below_zero([1, "2", -3])