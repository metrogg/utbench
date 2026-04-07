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

    def test_non_numeric_raises(self):
        with pytest.raises(NumberValidationError, match="numbers\\[3\\] must be numeric, got str"):
            _coerce_number("hello", index=3)

    def test_infinite_raises(self):
        with pytest.raises(NumberValidationError, match="numbers\\[4\\] must be finite"):
            _coerce_number(float("inf"), index=4)

    def test_nan_raises(self):
        with pytest.raises(NumberValidationError, match="numbers\\[5\\] must be finite"):
            _coerce_number(float("nan"), index=5)


class TestNormalizeNumbers:
    def test_empty(self):
        assert _normalize_numbers([]) == []

    def test_valid_numbers(self):
        assert _normalize_numbers([1, 2.5, -3]) == [1.0, 2.5, -3.0]

    def test_mixed_valid(self):
        assert _normalize_numbers([1, "2.5", -3.0]) == [1.0, 2.5, -3.0]

    def test_invalid_raises(self):
        with pytest.raises(NumberValidationError):
            _normalize_numbers([1, True, 3])


class TestNearestPair:
    def test_empty(self):
        assert nearest_pair([]) is None

    def test_single(self):
        assert nearest_pair([1.0]) is None

    def test_two_elements(self):
        result = nearest_pair([5.0, 3.0])
        assert result == PairDistance(
            left_index=1, right_index=0, left_value=3.0, right_value=5.0, distance=2.0
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

    def test_tie_distance_returns_first_found(self):
        result = nearest_pair([0.0, 5.0, 10.0])
        assert result.distance == 5.0
        assert (result.left_index, result.right_index) == (0, 1)

    def test_negative_numbers(self):
        result = nearest_pair([-10.0, -5.0, 0.0])
        assert result == PairDistance(
            left_index=0, right_index=1, left_value=-10.0, right_value=-5.0, distance=5.0
        )

    def test_duplicate_values(self):
        result = nearest_pair([1.0, 1.0, 2.0])
        assert result.distance == 0.0
        assert result.left_value == result.right_value == 1.0


class TestHasCloseElements:
    def test_empty(self):
        assert has_close_elements([], 0.5) is False

    def test_single(self):
        assert has_close_elements([1.0], 0.5) is False

    def test_no_close_pair(self):
        assert has_close_elements([1.0, 2.0, 3.0], 0.5) is False

    def test_close_pair_exists(self):
        assert has_close_elements([1.0, 1.2, 3.0], 0.5) is True

    def test_exactly_at_threshold(self):
        assert has_close_elements([1.0, 1.5, 3.0], 0.5) is False

    def test_negative_threshold_returns_false(self):
        assert has_close_elements([1.0, 1.1], -0.1) is False

    def test_threshold_zero(self):
        assert has_close_elements([1.0, 1.0, 2.0], 0.0) is False

    def test_threshold_just_above_zero(self):
        assert has_close_elements([1.0, 1.0, 2.0], 0.0001) is True

    def test_invalid_threshold_bool(self):
        with pytest.raises(NumberValidationError, match="threshold must be numeric"):
            has_close_elements([1.0, 2.0], True)

    def test_invalid_threshold_string(self):
        with pytest.raises(NumberValidationError, match="threshold must be numeric"):
            has_close_elements([1.0, 2.0], "0.5")

    def test_example_from_docstring(self):
        assert has_close_elements([1.0, 2.0, 3.0], 0.5) is False
        assert has_close_elements([1.0, 2.8, 3.0, 4.0, 5.0, 2.0], 0.3) is True

    def test_normalization_applied(self):
        assert has_close_elements([1, "2.0", 3.0], 1.5) is True


class TestCountClosePairs:
    def test_empty(self):
        assert count_close_pairs([], 0.5) == 0

    def test_single(self):
        assert count_close_pairs([1.0], 0.5) == 0

    def test_no_close_pairs(self):
        assert count_close_pairs([1.0, 2.0, 3.0], 0.5) == 0

    def test_one_close_pair(self):
        assert count_close_pairs([1.0, 1.2, 3.0], 0.5) == 1

    def test_multiple_close_pairs(self):
        assert count_close_pairs([1.0, 1.1, 1.2, 2.0], 0.5) == 3

    def test_all_pairs_close(self):
        assert count_close_pairs([1.0, 1.1, 1.2], 0.5) == 3

    def test_exactly_at_threshold(self):
        assert count_close_pairs([1.0, 1.5, 2.0], 0.5) == 0

    def test_negative_threshold_returns_zero(self):
        assert count_close_pairs([1.0, 1.1], -0.1) == 0

    def test_threshold_zero(self):
        assert count_close_pairs([1.0, 1.0, 2.0], 0.0) == 0

    def test_duplicate_values(self):
        assert count_close_pairs([1.0, 1.0, 1.0], 0.0) == 0
        assert count_close_pairs([1.0, 1.0, 1.0], 0.1) == 3

    def test_normalization_applied(self):
        assert count_close_pairs([1, "2.0", 2.1], 0.5) == 1


class TestMinDistance:
    def test_empty(self):
        assert min_distance([]) is None

    def test_single(self):
        assert min_distance([1.0]) is None

    def test_two_elements(self):
        assert min_distance([5.0, 3.0]) == 2.0

    def test_three_elements(self):
        assert min_distance([1.0, 3.0, 2.0]) == 1.0

    def test_duplicate_values(self):
        assert min_distance([1.0, 1.0, 2.0]) == 0.0

    def test_normalization_applied(self):
        assert min_distance([1, "2.5", 3.0]) == 0.5