import pytest

# Including the source code to ensure the test block is runnable independently
def below_zero(operations):
    balance = 0
    for op in operations:
        balance += op
        if balance < 0:
            return True
    return False


def test_below_zero_empty_list():
    """Test with an empty list of operations."""
    # Balance starts at 0 and never changes, so it never falls below zero.
    assert below_zero([]) is False


def test_below_zero_all_positive():
    """Test with a list of all positive integers."""
    # Balance strictly increases.
    assert below_zero([1, 2, 3]) is False
    assert below_zero([10, 50, 100]) is False


def test_below_zero_all_negative():
    """Test with a list of all negative integers."""
    # Balance becomes negative immediately on the first operation.
    assert below_zero([-1, -2, -3]) is True
    assert below_zero([-5]) is True


def test_below_zero_mixed_operations_no_negative():
    """Test with mixed operations where balance never drops below zero."""
    # Sequence: 0 -> 10 -> 5 -> 1
    assert below_zero([10, -5, -4]) is False
    # Sequence: 0 -> 5 -> 0 -> 2
    assert below_zero([5, -5, 2]) is False


def test_below_zero_mixed_operations_drops_negative():
    """Test with mixed operations where balance drops below zero."""
    # Sequence: 0 -> 1 -> 3 -> -1
    assert below_zero([1, 2, -4, 5]) is True
    # Sequence: 0 -> 5 -> -2
    assert below_zero([5, -7]) is True


def test_below_zero_boundary_zero():
    """Test cases where the balance hits exactly zero."""
    # Balance starts at 0.
    assert below_zero([0]) is False
    # Balance goes 5 -> 0.
    assert below_zero([5, -5]) is False
    # Balance goes 0 -> 0 -> 0.
    assert below_zero([0, 0, 0]) is False


def test_below_zero_large_input():
    """Test with a large list to ensure performance and determinism."""
    # Create a large list of positive numbers followed by a large negative withdrawal
    # that forces the balance below zero.
    positive_ops = [1] * 10000
    negative_op = [-10001]
    operations = positive_ops + negative_op
    
    assert below_zero(operations) is True


def test_below_zero_recovery_irrelevant():
    """Test that function returns True immediately upon dropping below zero."""
    # Even if the balance recovers later, the function should return True at the first dip.
    # Sequence: 0 -> -5 -> 10
    assert below_zero([-5, 10]) is True