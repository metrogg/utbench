import pytest
from simple_function_python_simple_function_0 import (
    NumberValidationError,
    PairDistance,
    _coerce_number,
    _normalize_numbers,
    nearest_pair,
    has_close_elements,
    count_close_pairs,
    min_distance,
)


class TestCoerceNumber:
    def test_valid_int(self):
        assert _coerce_number(42, index=0) == 42.0

    def test_valid_float(self):
        assert _coerce_number(3.14, index=1) == 3.14

    def test_bool_raises(self):
        with pytest.raises(NumberValidationError, match="numbers\\[2\\] must be numeric, got bool"):
            _coerce_number(True, index=2)

    def test_non_numeric_type_raises(self):
        with pytest.raises(NumberValidationError, match="numbers\\[3\\] must be numeric, got str"):
            _coerce_number("hello", index=3)

    def test_infinite_raises(self):
        with pytest.raises(NumberValidationError, match="numbers\\[4\\] must be finite"):
            _coerce_number(float("inf"), index=4)

    def test_nan_raises(self):
        with pytest.raises(NumberValidationError, match="numbers\\[5\\] must be finite"):
            _coerce_number(float("nan"), index=5)


class TestNormalizeNumbers:
    def test_empty_iterable(self):
        assert _normalize_numbers([]) == []

    def test_valid_numbers(self):
        input_data = [1, 2.5, -3]
        expected = [1.0, 2.5, -3.0]
        assert _normalize_numbers(input_data) == expected

    def test_mixed_valid_types(self):
        input_data = [10, 20.0]
        expected = [10.0, 20.0]
        assert _normalize_numbers(input_data) == expected

    def test_invalid_element_raises(self):
        with pytest.raises(NumberValidationError):
            _normalize_numbers([1, "invalid", 3])

    def test_infinite_element_raises(self):
        with pytest.raises(NumberValidationError):
            _normalize_numbers([1.0, float("inf")])


class TestNearestPair:
    def test_empty_sequence(self):
        assert nearest_pair([]) is None

    def test_single_element(self):
        assert nearest_pair([5.0]) is None

    def test_two_elements(self):
        result = nearest_pair([1.0, 2.0])
        assert result == PairDistance(
            left_index=0, right_index=1, left_value=1.0, right_value=2.0, distance=1.0
        )

    def test_three_elements_sorted(self):
        result = nearest_pair([1.0, 2.0, 3.0])
        assert result == PairDistance(
            left_index=0, right_index=1, left_value=1.0, right_value=2.0, distance=1.0
        )

    def test_three_elements_unsorted(self):
        result = nearest_pair([3.0, 1.0, 2.0])
        assert result == PairDistance(
            left_index=1, right_index=2, left_value=1.0, right_value=2.0, distance=1.0
        )

    def test_duplicate_values(self):
        result = nearest_pair([5.0, 5.0, 10.0])
        assert result == PairDistance(
            left_index=0, right_index=1, left_value=5.0, right_value=5.0, distance=0.0
        )

    def test_negative_numbers(self):
        result = nearest_pair([-10.0, -5.0, 0.0])
        assert result == PairDistance(
            left_index=0, right_index=1, left_value=-10.0, right_value=-5.0, distance=5.0
        )

    def test_tie_distance_returns_first_encountered(self):
        result = nearest_pair([0.0, 5.0, 10.0, 15.0])
        assert result.distance == 5.0
        assert result.left_index == 0
        assert result.right_index == 1


class TestHasCloseElements:
    def test_empty_sequence(self):
        assert has_close_elements([], 1.0) is False

    def test_single_element(self):
        assert has_close_elements([5.0], 1.0) is False

    def test_no_close_elements(self):
        assert has_close_elements([1.0, 2.0, 3.0], 0.5) is False

    def test_has_close_elements(self):
        assert has_close_elements([1.0, 2.8, 3.0, 4.0, 5.0, 2.0], 0.3) is True

    def test_exact_threshold(self):
        assert has_close_elements([1.0, 2.0], 1.0) is False

    def test_negative_threshold_returns_false(self):
        assert has_close_elements([1.0, 2.0], -0.1) is False

    def test_zero_threshold_with_duplicates(self):
        assert has_close_elements([1.0, 1.0, 2.0], 0.0) is True

    def test_zero_threshold_without_duplicates(self):
        assert has_close_elements([1.0, 2.0], 0.0) is False

    def test_invalid_threshold_type_raises(self):
        with pytest.raises(NumberValidationError, match="threshold must be numeric"):
            has_close_elements([1.0, 2.0], "invalid")

    def test_bool_threshold_raises(self):
        with pytest.raises(NumberValidationError, match="threshold must be numeric"):
            has_close_elements([1.0, 2.0], True)

    def test_invalid_numbers_in_sequence_raises(self):
        with pytest.raises(NumberValidationError):
            has_close_elements([1.0, "invalid"], 0.5)


class TestCountClosePairs:
    def test_empty_sequence(self):
        assert count_close_pairs([], 1.0) == 0

    def test_single_element(self):
        assert count_close_pairs([5.0], 1.0) == 0

    def test_no_close_pairs(self):
        assert count_close_pairs([1.0, 2.0, 3.0], 0.5) == 0

    def test_all_pairs_close(self):
        assert count_close_pairs([1.0, 1.1, 1.2], 0.5) == 3

    def test_some_close_pairs(self):
        assert count_close_pairs([1.0, 2.0, 2.1, 3.0], 0.2) == 1

    def test_duplicate_values(self):
        assert count_close_pairs([5.0, 5.0, 5.0], 0.0) == 3

    def test_negative_threshold_returns_zero(self):
        assert count_close_pairs([1.0, 2.0], -0.1) == 0

    def test_zero_threshold_with_duplicates(self):
        assert count_close_pairs([1.0, 1.0, 2.0], 0.0) == 1

    def test_zero_threshold_without_duplicates(self):
        assert count_close_pairs([1.0, 2.0, 3.0], 0.0) == 0

    def test_large_threshold_includes_all_pairs(self):
        numbers = [1.0, 2.0, 3.0]
        n = len(numbers)
        expected = n * (n - 1) // 2
        assert count_close_pairs(numbers, 10.0) == expected

    def test_invalid_numbers_in_sequence_raises(self):
        with pytest.raises(NumberValidationError):
            count_close_pairs([1.0, "invalid"], 0.5)


class TestMinDistance:
    def test_empty_sequence(self):
        assert min_distance([]) is None

    def test_single_element(self):
        assert min_distance([5.0]) is None

    def test_two_elements(self):
        assert min_distance([1.0, 2.0]) == 1.0

    def test_multiple_elements(self):
        assert min_distance([3.0, 1.0, 4.0, 2.0]) == 1.0

    def test_duplicate_values(self):
        assert min_distance([5.0, 5.0, 10.0]) == 0.0

    def test_negative_numbers(self):
        assert min_distance([-10.0, -5.0, 0.0]) == 5.0

    def test_invalid_numbers_in_sequence_raises(self):
        with pytest.raises(NumberValidationError):
            min_distance([1.0, "invalid"])