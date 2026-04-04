import pytest
from typing import List


# Copy of the source function to make test file self-contained
def below_zero(operations: List[int]) -> bool:
    balance = 0
    for op in operations:
        balance += op
        if balance < 0:
            return True
    return False


def test_empty_operations_list():
    assert below_zero([]) is False


def test_all_positive_operations():
    assert below_zero([1, 2, 3]) is False


def test_go_below_zero_in_middle():
    assert below_zero([1, 2, -4, 5]) is True


def test_first_operation_is_negative():
    assert below_zero([-10, 20, 30]) is True


def test_balance_never_dips_exactly_zero_only():
    assert below_zero([5, -5, 3, -3, 0]) is False


def test_go_below_zero_then_recover_still_returns_true():
    assert below_zero([-1, 100, 200]) is True


def test_single_zero_operation():
    assert below_zero([0]) is False


def test_single_positive_operation():
    assert below_zero([999]) is False


def test_single_negative_operation():
    assert below_zero([-1]) is True


def test_zeros_followed_by_negative_operation():
    assert below_zero([0, 0, 0, -1]) is True


def test_go_below_zero_at_last_operation():
    assert below_zero([2, 3, -6]) is True


def test_balance_never_below_zero_sum_to_zero():
    assert below_zero([4, -2, -1, -1]) is False


def test_large_number_operations_go_below_zero():
    assert below_zero([1_000_000, -1_000_001]) is True


def test_all_zero_operations():
    assert below_zero([0, 0, 0, 0]) is False


def test_negative_zero_operation_returns_false():
    assert below_zero([-0]) is False


def test_non_integer_operation_raises_type_error():
    with pytest.raises(TypeError):
        below_zero([1, "string", -3])