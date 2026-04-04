import pytest
from typing import List

def test_has_close_elements_no_close_elements():
    """Test when no elements are closer than threshold"""
    numbers = [1.0, 2.0, 3.0, 4.0]
    threshold = 0.5
    result = has_close_elements(numbers, threshold)
    assert result == False

def test_has_close_elements_with_close_elements():
    """Test when there are elements closer than threshold"""
    numbers = [1.0, 1.2, 3.0, 4.0]
    threshold = 0.3
    result = has_close_elements(numbers, threshold)
    assert result == True

def test_has_close_elements_exact_threshold():
    """Test when distance equals threshold (should return False since condition is < threshold)"""
    numbers = [1.0, 1.5, 3.0]
    threshold = 0.5
    result = has_close_elements(numbers, threshold)
    assert result == False

def test_has_close_elements_negative_numbers():
    """Test with negative numbers"""
    numbers = [-2.0, -1.9, 0.0, 1.0]
    threshold = 0.2
    result = has_close_elements(numbers, threshold)
    assert result == True

def test_has_close_elements_mixed_positive_negative():
    """Test with mixed positive and negative numbers"""
    numbers = [-1.0, 0.9, 2.0, 3.0]
    threshold = 0.2
    result = has_close_elements(numbers, threshold)
    assert result == True

def test_has_close_elements_single_element():
    """Test with single element list"""
    numbers = [1.0]
    threshold = 0.1
    result = has_close_elements(numbers, threshold)
    assert result == False

def test_has_close_elements_empty_list():
    """Test with empty list"""
    numbers = []
    threshold = 0.1
    result = has_close_elements(numbers, threshold)
    assert result == False

def test_has_close_elements_identical_elements():
    """Test with identical elements (distance = 0)"""
    numbers = [1.0, 1.0, 2.0, 3.0]
    threshold = 0.1
    result = has_close_elements(numbers, threshold)
    assert result == True

def test_has_close_elements_zero_threshold():
    """Test with zero threshold"""
    numbers = [1.0, 1.0001, 2.0]
    threshold = 0.0
    result = has_close_elements(numbers, threshold)
    assert result == False

def test_has_close_elements_negative_threshold():
    """Test with negative threshold (edge case)"""
    numbers = [1.0, 2.0, 3.0]
    threshold = -0.5
    result = has_close_elements(numbers, threshold)
    assert result == False

def test_has_close_elements_large_numbers():
    """Test with large numbers"""
    numbers = [1000.0, 1000.1, 2000.0]
    threshold = 0.2
    result = has_close_elements(numbers, threshold)
    assert result == True

def test_has_close_elements_decimal_precision():
    """Test with decimal precision"""
    numbers = [1.000001, 1.000002, 2.0]
    threshold = 0.000001
    result = has_close_elements(numbers, threshold)
    assert result == True

def test_has_close_elements_first_and_last_close():
    """Test when first and last elements are close"""
    numbers = [1.0, 5.0, 10.0, 1.1]
    threshold = 0.2
    result = has_close_elements(numbers, threshold)
    assert result == True

def test_has_close_elements_adjacent_close():
    """Test when adjacent elements are close"""
    numbers = [1.0, 1.05, 2.0, 3.0]
    threshold = 0.1
    result = has_close_elements(numbers, threshold)
    assert result == True

def test_has_close_elements_non_adjacent_close():
    """Test when non-adjacent elements are close"""
    numbers = [1.0, 3.0, 2.0, 1.05]
    threshold = 0.1
    result = has_close_elements(numbers, threshold)
    assert result == True

def test_has_close_elements_very_small_threshold():
    """Test with very small threshold"""
    numbers = [1.0, 1.00000001, 2.0]
    threshold = 0.000000001
    result = has_close_elements(numbers, threshold)
    assert result == False