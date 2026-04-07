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
    def test_coerce_int(self):
        result = _coerce_number(42, index=0)
        assert result == 42.0
        assert isinstance(result, float)

    def test_coerce_float(self):
        result = _coerce_number(3.14, index=1)
        assert result == 3.14

    def test_coerce_float_from_int(self):
        result = _coerce_number(7, index=2)
        assert result == 7.0
        assert isinstance(result, float)

    def test_coerce_bool_raises(self):
        with pytest.raises(NumberValidationError) as exc_info:
            _coerce_number(True, index=0)
        assert "got bool" in str(exc_info.value)

    def test_coerce_string_raises(self):
        with pytest.raises(NumberValidationError) as exc_info:
            _coerce_number("abc", index=1)
        assert "got str" in str(exc_info.value)

    def test_coerce_none_raises(self):
        with pytest.raises(NumberValidationError) as exc_info:
            _coerce_number(None, index=2)
        assert "got NoneType" in str(exc_info.value)

    def test_coerce_inf_raises(self):
        with pytest.raises(NumberValidationError) as exc_info:
            _coerce_number(float("inf"), index=3)
        assert "must be finite" in str(exc_info.value)

    def test_coerce_nan_raises(self):
        with pytest.raises(NumberValidationError) as exc_info:
            _coerce_number(float("nan"), index=4)
        assert "must be finite" in str(exc_info.value)


class TestNormalizeNumbers:
    def test_normalize_empty(self):
        result = _normalize_numbers([])
        assert result == []

    def test_normalize_ints(self):
        result = _normalize_numbers([1, 2, 3])
        assert result == [1.0, 2.0, 3.0]

    def test_normalize_floats(self):
        result = _normalize_numbers([1.5, 2.5, 3.5])
        assert result == [1.5, 2.5, 3.5]

    def test_normalize_mixed(self):
        result = _normalize_numbers([1, 2.5, 3])
        assert result == [1.0, 2.5, 3.0]

    def test_normalize_invalid_raises(self):
        with pytest.raises(NumberValidationError):
            _normalize_numbers([1, "invalid", 3])


class TestNearestPair:
    def test_nearest_pair_empty(self):
        result = nearest_pair([])
        assert result is None

    def test_nearest_pair_single(self):
        result = nearest_pair([1.0])
        assert result is None

    def test_nearest_pair_two_ascending(self):
        result = nearest_pair([1.0, 2.0])
        assert result == PairDistance(
            left_index=0, right_index=1, left_value=1.0, right_value=2.0, distance=1.0
        )

    def test_nearest_pair_two_descending(self):
        result = nearest_pair([2.0, 1.0])
        assert result == PairDistance(
            left_index=0, right_index=1, left_value=1.0, right_value=2.0, distance=1.0
        )

    def test_nearest_pair_three_sorted(self):
        result = nearest_pair([1.0, 2.0, 3.0])
        assert result == PairDistance(
            left_index=0, right_index=1, left_value=1.0, right_value=2.0, distance=1.0
        )

    def test_nearest_pair_three_unsorted(self):
        result = nearest_pair([3.0, 1.0, 2.0])
        assert result == PairDistance(
            left_index=1, right_index=2, left_value=1.0, right_value=2.0, distance=1.0
        )

    def test_nearest_pair_equal_values(self):
        result = nearest_pair([5.0, 5.0, 5.0])
        assert result == PairDistance(
            left_index=0, right_index=1, left_value=5.0, right_value=5.0, distance=0.0
        )

    def test_nearest_pair_negative_values(self):
        result = nearest_pair([-5.0, -3.0, -1.0])
        assert result == PairDistance(
            left_index=0,
            right_index=1,
            left_value=-5.0,
            right_value=-3.0,
            distance=2.0,
        )

    def test_nearest_pair_mixed_sign(self):
        result = nearest_pair([-2.0, 0.0, 2.0])
        assert result == PairDistance(
            left_index=0, right_index=1, left_value=-2.0, right_value=0.0, distance=2.0
        )

    def test_nearest_pair_tie_distance_picks_smaller_indices(self):
        result = nearest_pair([0.0, 5.0, 10.0, 15.0])
        assert result == PairDistance(
            left_index=0, right_index=1, left_value=0.0, right_value=5.0, distance=5.0
        )


class TestHasCloseElements:
    def test_has_close_elements_empty(self):
        result = has_close_elements([], 0.5)
        assert result is False

    def test_has_close_elements_single(self):
        result = has_close_elements([1.0], 0.5)
        assert result is False

    def test_has_close_elements_no_close(self):
        result = has_close_elements([1.0, 2.0, 3.0], 0.5)
        assert result is False

    def test_has_close_elements_close_exists(self):
        result = has_close_elements([1.0, 1.2, 3.0], 0.5)
        assert result is True

    def test_has_close_elements_exact_threshold_not_close(self):
        result = has_close_elements([1.0, 2.0, 3.0], 1.0)
        assert result is False

    def test_has_close_elements_below_threshold_close(self):
        result = has_close_elements([1.0, 1.9, 3.0], 1.0)
        assert result is True

    def test_has_close_elements_negative_threshold_returns_false(self):
        result = has_close_elements([1.0, 1.1, 2.0], -0.1)
        assert result is False

    def test_has_close_elements_zero_threshold_with_equal_values(self):
        result = has_close_elements([1.0, 1.0, 2.0], 0.0)
        assert result is False

    def test_has_close_elements_zero_threshold_with_different_values(self):
        result = has_close_elements([1.0, 1.1, 2.0], 0.0)
        assert result is False

    def test_has_close_elements_bool_threshold_raises(self):
        with pytest.raises(NumberValidationError) as exc_info:
            has_close_elements([1.0, 2.0], True)
        assert "threshold must be numeric" in str(exc_info.value)

    def test_has_close_elements_string_threshold_raises(self):
        with pytest.raises(NumberValidationError) as exc_info:
            has_close_elements([1.0, 2.0], "0.5")
        assert "threshold must be numeric" in str(exc_info.value)

    def test_has_close_elements_none_threshold_raises(self):
        with pytest.raises(NumberValidationError) as exc_info:
            has_close_elements([1.0, 2.0], None)
        assert "threshold must be numeric" in str(exc_info.value)

    def test_has_close_elements_invalid_number_in_list_raises(self):
        with pytest.raises(NumberValidationError):
            has_close_elements([1.0, "invalid", 3.0], 0.5)


class TestCountClosePairs:
    def test_count_close_pairs_empty(self):
        result = count_close_pairs([], 0.5)
        assert result == 0

    def test_count_close_pairs_single(self):
        result = count_close_pairs([1.0], 0.5)
        assert result == 0

    def test_count_close_pairs_two_close(self):
        result = count_close_pairs([1.0, 1.2], 0.5)
        assert result == 1

    def test_count_close_pairs_two_not_close(self):
        result = count_close_pairs([1.0, 2.0], 0.5)
        assert result == 0

    def test_count_close_pairs_three_all_close(self):
        result = count_close_pairs([1.0, 1.1, 1.2], 0.5)
        assert result == 3

    def test_count_close_pairs_three_some_close(self):
        result = count_close_pairs([1.0, 1.4, 2.0], 0.5)
        assert result == 1

    def test_count_close_pairs_equal_values_zero_threshold(self):
        result = count_close_pairs([5.0, 5.0, 5.0], 0.0)
        assert result == 0

    def test_count_close_pairs_equal_values_positive_threshold(self):
        result = count_close_pairs([5.0, 5.0, 5.0], 0.1)
        assert result == 3

    def test_count_close_pairs_negative_threshold_returns_zero(self):
        result = count_close_pairs([1.0, 1.1, 2.0], -0.1)
        assert result == 0

    def test_count_close_pairs_large_threshold(self):
        result = count_close_pairs([1.0, 2.0, 3.0], 10.0)
        assert result == 3

    def test_count_close_pairs_exact_threshold_not_counted(self):
        result = count_close_pairs([1.0, 2.0, 3.0], 1.0)
        assert result == 1

    def test_count_close_pairs_unsorted_input(self):
        result = count_close_pairs([3.0, 1.0, 2.0], 1.0)
        assert result == 1

    def test_count_close_pairs_invalid_number_raises(self):
        with pytest.raises(NumberValidationError):
            count_close_pairs([1.0, "invalid", 3.0], 0.5)


class TestMinDistance:
    def test_min_distance_empty(self):
        result = min_distance([])
        assert result is None

    def test_min_distance_single(self):
        result = min_distance([1.0])
        assert result is None

    def test_min_distance_two(self):
        result = min_distance([1.0, 2.0])
        assert result == 1.0

    def test_min_distance_three(self):
        result = min_distance([1.0, 3.0, 2.0])
        assert result == 1.0

    def test_min_distance_equal_values(self):
        result = min_distance([5.0, 5.0, 5.0])
        assert result == 0.0

    def test_min_distance_negative_values(self):
        result = min_distance([-5.0, -3.0, -1.0])
        assert result == 2.0

    def test_min_distance_invalid_number_raises(self):
        with pytest.raises(NumberValidationError):
            min_distance([1.0, "invalid", 3.0])