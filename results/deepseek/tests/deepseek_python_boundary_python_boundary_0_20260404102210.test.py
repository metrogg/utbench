import pytest
from typing import List


def test_below_zero_empty_list():
    """Test with empty operations list."""
    result = below_zero([])
    assert result is False


def test_below_zero_single_positive():
    """Test with single positive operation."""
    result = below_zero([10])
    assert result is False


def test_below_zero_single_negative():
    """Test with single negative operation."""
    result = below_zero([-5])
    assert result is True


def test_below_zero_all_positive():
    """Test with all positive operations."""
    result = below_zero([1, 2, 3, 4, 5])
    assert result is False


def test_below_zero_all_negative():
    """Test with all negative operations."""
    result = below_zero([-1, -2, -3])
    assert result is True


def test_below_zero_mixed_no_below_zero():
    """Test with mixed operations that never go below zero."""
    result = below_zero([10, -5, 3, -2, 1])
    assert result is False


def test_below_zero_mixed_below_zero():
    """Test with mixed operations that go below zero."""
    result = below_zero([10, -15, 5])
    assert result is True


def test_below_zero_immediate_below_zero():
    """Test where balance goes below zero immediately."""
    result = below_zero([-1, 100, 200])
    assert result is True


def test_below_zero_late_below_zero():
    """Test where balance goes below zero late in sequence."""
    result = below_zero([10, 20, -35, 5])
    assert result is True


def test_below_zero_exact_zero_balance():
    """Test where balance reaches exactly zero but not below."""
    result = below_zero([10, -10, 5, -5])
    assert result is False


def test_below_zero_boundary_negative_one():
    """Test with balance exactly -1 (below zero)."""
    result = below_zero([5, -6])
    assert result is True


def test_below_zero_large_numbers():
    """Test with large positive and negative numbers."""
    result = below_zero([1000000, -999999, -2])
    assert result is True


def test_below_zero_alternating():
    """Test with alternating positive and negative values."""
    result = below_zero([5, -3, 2, -4, 1])
    assert result is True


def test_below_zero_from_docstring_example1():
    """Test first example from docstring."""
    result = below_zero([1, 2, 3])
    assert result is False


def test_below_zero_from_docstring_example2():
    """Test second example from docstring."""
    result = below_zero([1, 2, -4, 5])
    assert result is True


def test_below_zero_zero_operations():
    """Test with zero values in operations."""
    result = below_zero([0, 0, 0])
    assert result is False


def test_below_zero_mixed_with_zeros():
    """Test with zeros mixed with other operations."""
    result = below_zero([5, 0, -6, 0, 2])
    assert result is True


def test_below_zero_early_return():
    """Test that function returns immediately when balance goes below zero."""
    operations = [-1, 100, 200]
    result = below_zero(operations)
    assert result is True


def test_below_zero_no_early_return():
    """Test that function processes all operations when never below zero."""
    operations = [1, 2, 3, 4]
    result = below_zero(operations)
    assert result is False