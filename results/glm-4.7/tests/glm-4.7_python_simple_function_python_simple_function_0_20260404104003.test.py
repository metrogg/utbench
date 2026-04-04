import pytest
from simple_function import has_close_elements


def test_has_close_elements_basic_false():
    """Test the example case where elements are not close enough."""
    numbers = [1.0, 2.0, 3.0]
    threshold = 0.5
    assert has_close_elements(numbers, threshold) is False


def test_has_close_elements_basic_true():
    """Test the example case where elements are close enough."""
    numbers = [1.0, 2.8, 3.0, 4.0, 5.0, 2.0]
    threshold = 0.3
    assert has_close_elements(numbers, threshold) is True


def test_has_close_elements_empty_list():
    """Test with an empty list. Should return False."""
    assert has_close_elements([], 1.0) is False


def test_has_close_elements_single_element():
    """Test with a single element list. Cannot have pairs, so False."""
    assert has_close_elements([1.0], 0.5) is False


def test_has_close_elements_two_elements_far():
    """Test with two elements where distance is greater than threshold."""
    assert has_close_elements([1.0, 2.0], 0.9) is False


def test_has_close_elements_two_elements_close():
    """Test with two elements where distance is less than threshold."""
    assert has_close_elements([1.0, 1.5], 0.6) is True


def test_has_close_elements_boundary_equal():
    """Test boundary condition where distance equals threshold.
    Logic uses strict less-than (<), so should be False."""
    assert has_close_elements([1.0, 2.0], 1.0) is False


def test_has_close_elements_negative_threshold():
    """Test with negative threshold.
    Distance is always non-negative, so cannot be less than negative."""
    assert has_close_elements([1.0, 2.0], -1.0) is False


def test_has_close_elements_zero_threshold():
    """Test with zero threshold.
    Distance 0 is not less than 0."""
    assert has_close_elements([1.0, 1.0], 0.0) is False


def test_has_close_elements_negative_numbers():
    """Test with negative numbers."""
    # -1.0 and -1.1 are 0.1 apart. Threshold is 0.2.
    assert has_close_elements([-1.0, -1.1, -5.0], 0.2) is True


def test_has_close_elements_mixed_signs():
    """Test with mixed positive and negative numbers."""
    # 0.0 and 0.5 are 0.5 apart. Threshold is 0.6.
    assert has_close_elements([-1.0, 0.5, 0.0], 0.6) is True