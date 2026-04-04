# test_bank.py
import pytest
from solution import below_zero


class TestBelowZero:
    """Unit tests for below_zero function."""

    # Normal cases

    def test_all_positive(self):
        """When all operations are positive, balance never goes negative."""
        assert below_zero([1, 2, 3]) is False

    def test_mixed_operations_with_negative_that_triggers(self):
        """If any operation makes the balance drop below zero, function returns True."""
        assert below_zero([1, 2, -4, 5]) is True

    def test_negative_at_start(self):
        """A negative operation as the first one should immediately return True."""
        assert below_zero([-5, 2, 3]) is True

    # Boundary cases

    def test_empty_list(self):
        """With no operations the balance stays at zero, which is not below zero."""
        assert below_zero([]) is False

    def test_single_positive(self):
        """A single positive operation never makes the balance negative."""
        assert below_zero([10]) is False

    def test_single_negative(self):
        """A single negative operation makes the balance negative."""
        assert below_zero([-1]) is True

    def test_balance_ends_at_zero(self):
        """If the balance reaches exactly zero, it's not below zero."""
        assert below_zero([5, -5]) is False

    def test_balance_hits_zero_then_negative(self):
        """If after reaching zero a later operation makes it negative, return True."""
        # Never goes below zero
        assert below_zero([5, -3, -2]) is False
        # Goes below zero after zero
        assert below_zero([5, -5, -1]) is True

    def test_zero_operations(self):
        """Zero-valued operations should not affect the balance."""
        assert below_zero([0, 0, 0]) is False

    def test_mixed_zeros_and_negatives_without_going_below_zero(self):
        """Zero operations interspersed with negatives that never break zero."""
        assert below_zero([5, 0, -3, -2]) is False

    # Edge cases with large numbers

    def test_large_positive_and_negative_balancing(self):
        """Large numbers that exactly balance should return False."""
        assert below_zero([10**20, -10**20]) is False

    def test_large_positive_and_slightly_larger_negative(self):
        """If a large negative tips the balance below zero, return True."""
        assert below_zero([10**20, -10**20 - 1]) is True

    # Error handling / invalid inputs

    def test_non_iterable_input(self):
        """Passing a non‑iterable (e.g., int) should raise TypeError."""
        with pytest.raises(TypeError):
            below_zero(123)

    def test_none_input(self):
        """Passing None should raise TypeError."""
        with pytest.raises(TypeError):
            below_zero(None)

    def test_list_with_non_int_elements(self):
        """If the list contains non‑int elements, a TypeError should be raised."""
        with pytest.raises(TypeError):
            below_zero([1, "two", 3])