import pytest
from solution import below_zero

class TestBelowZero:
    def test_empty_operations(self):
        assert below_zero([]) is False

    def test_all_deposits(self):
        assert below_zero([1, 2, 3]) is False

    def test_immediate_negative_balance(self):
        assert below_zero([-5]) is True

    def test_negative_balance_in_middle(self):
        assert below_zero([1, 2, -4, 5]) is True

    def test_balance_returns_to_zero(self):
        assert below_zero([10, -10, 5]) is False

    def test_negative_balance_at_end(self):
        assert below_zero([1, 1, 1, -5]) is True

    def test_zero_operations_only(self):
        assert below_zero([0, 0, 0]) is False

    def test_large_values(self):
        assert below_zero([10**9, -10**9 - 1]) is True

    def test_exact_boundary_zero(self):
        assert below_zero([5, -3, -2]) is False

    def test_invalid_input_type(self):
        with pytest.raises(TypeError):
            below_zero("invalid_type")