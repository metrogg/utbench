import pytest
from typing import List

# Import the function to test
from module_name import below_zero  # Replace 'module_name' with actual module name


class TestBelowZero:
    """Test suite for below_zero function."""

    def test_normal_case_no_below_zero(self) -> None:
        """Test normal case where balance never goes below zero."""
        operations = [1, 2, 3]
        result = below_zero(operations)
        assert result is False

    def test_normal_case_below_zero_middle(self) -> None:
        """Test normal case where balance goes below zero in the middle."""
        operations = [1, 2, -4, 5]
        result = below_zero(operations)
        assert result is True

    def test_normal_case_below_zero_end(self) -> None:
        """Test normal case where balance goes below zero at the end."""
        operations = [1, 2, 3, -10]
        result = below_zero(operations)
        assert result is True

    def test_normal_case_below_zero_immediate(self) -> None:
        """Test normal case where first operation makes balance negative."""
        operations = [-5, 10, 3]
        result = below_zero(operations)
        assert result is True

    def test_boundary_case_empty_list(self) -> None:
        """Test boundary case with empty operations list."""
        operations: List[int] = []
        result = below_zero(operations)
        assert result is False

    def test_boundary_case_single_positive(self) -> None:
        """Test boundary case with single positive operation."""
        operations = [10]
        result = below_zero(operations)
        assert result is False

    def test_boundary_case_single_negative(self) -> None:
        """Test boundary case with single negative operation."""
        operations = [-10]
        result = below_zero(operations)
        assert result is True

    def test_boundary_case_single_zero(self) -> None:
        """Test boundary case with single zero operation."""
        operations = [0]
        result = below_zero(operations)
        assert result is False

    def test_boundary_case_exactly_zero_balance(self) -> None:
        """Test boundary case where balance becomes exactly zero but not negative."""
        operations = [5, -5, 3, -3]
        result = below_zero(operations)
        assert result is False

    def test_boundary_case_large_numbers(self) -> None:
        """Test boundary case with very large numbers."""
        operations = [1000000, -500000, -500001]
        result = below_zero(operations)
        assert result is True

    def test_boundary_case_alternating_operations(self) -> None:
        """Test boundary case with alternating positive and negative operations."""
        operations = [10, -5, 3, -8, 2, -1]
        result = below_zero(operations)
        assert result is False

    def test_boundary_case_all_negative(self) -> None:
        """Test boundary case where all operations are negative."""
        operations = [-1, -2, -3]
        result = below_zero(operations)
        assert result is True

    def test_boundary_case_all_positive(self) -> None:
        """Test boundary case where all operations are positive."""
        operations = [1, 2, 3, 4, 5]
        result = below_zero(operations)
        assert result is False

    def test_boundary_case_mixed_with_zero(self) -> None:
        """Test boundary case with zero operations mixed in."""
        operations = [0, 5, 0, -3, 0, -2]
        result = below_zero(operations)
        assert result is True

    def test_error_handling_none_input(self) -> None:
        """Test error handling with None input."""
        with pytest.raises(TypeError):
            below_zero(None)  # type: ignore

    def test_error_handling_wrong_type_in_list(self) -> None:
        """Test error handling with wrong type in list."""
        with pytest.raises(TypeError):
            below_zero([1, "2", 3])  # type: ignore

    def test_edge_case_early_return(self) -> None:
        """Test that function returns immediately when balance goes below zero."""
        operations = [1, -2, 3, -4, 5]
        # The function should return True at the second operation (-2)
        result = below_zero(operations)
        assert result is True

    def test_edge_case_large_negative_then_positive(self) -> None:
        """Test case with large negative followed by positive that doesn't recover."""
        operations = [-100, 50, 49]
        result = below_zero(operations)
        assert result is True

    def test_edge_case_precise_negative(self) -> None:
        """Test case where balance goes exactly -1."""
        operations = [10, -11]
        result = below_zero(operations)
        assert result is True

    def test_edge_case_never_negative_with_negatives(self) -> None:
        """Test case with negative operations but balance never goes below zero."""
        operations = [10, -5, 3, -7, 4]
        result = below_zero(operations)
        assert result is False