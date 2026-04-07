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
        with pytest.raises(NumberValidationError) as exc:
            _coerce_number(True, index=2)
        assert "got bool" in str(exc.value)

    def test_non_numeric_raises(self):
        with pytest.raises(NumberValidationError) as exc:
            _coerce_number("string", index=3)
        assert "got str" in str(exc.value)

    def test_infinite_raises(self):
        with pytest.raises(NumberValidationError) as exc:
            _coerce_number(float("inf"), index=4)
        assert "must be finite" in str(exc.value)

    def test_nan_raises(self):
        with pytest.raises(NumberValidationError) as exc:
            _coerce_number(float("nan"), index=5)
        assert "must be finite" in str(exc.value)


class TestNormalizeNumbers:
    def test_empty(self):
        assert _normalize_numbers([]) == []

    def test_mixed_valid(self):
        result = _normalize_numbers([1, 2.5, -3])
        assert result == [1.0, 2.5, -3.0]

    def test_invalid_in_middle(self):
        with pytest.raises(NumberValidationError):
            _normalize_numbers([1, False, 3])


class TestNearestPair:
    def test_empty(self):
        assert nearest_pair([]) is None

    def test_single(self):
        assert nearest_pair([5.0]) is None

    def test_two_ascending(self):
        result = nearest_pair([1.0, 2.0])
        expected = PairDistance(
            left_index=0,
            right_index=1,
            left_value=1.0,
            right_value=2.0,
            distance=1.0,
        )
        assert result == expected

    def test_two_descending(self):
        result = nearest_pair([3.0, 1.0])
        expected = PairDistance(
            left_index=1,
            right_index=0,
            left_value=1.0,
            right_value=3.0,
            distance=2.0,
        )
        assert result == expected

    def test_three_with_tie_distance(self):
        # Values: [1.0, 3.0, 2.0] -> sorted by value: (0,1.0), (2,2.0), (1,3.0)
        # Gaps: 1.0 and 1.0. First candidate wins.
        result = nearest_pair([1.0, 3.0, 2.0])
        expected = PairDistance(
            left_index=0,
            right_index=2,
            left_value=1.0,
            right_value=2.0,
            distance=1.0,
        )
        assert result == expected

    def test_duplicate_values(self):
        result = nearest_pair([5.0, 5.0, 5.0])
        expected = PairDistance(
            left_index=0,
            right_index=1,
            left_value=5.0,
            right_value=5.0,
            distance=0.0,
        )
        assert result == expected

    def test_negative_numbers(self):
        result = nearest_pair([-10.0, -5.0, 0.0])
        expected = PairDistance(
            left_index=1,
            right_index=2,
            left_value=-5.0,
            right_value=0.0,
            distance=5.0,
        )
        assert result == expected


class TestHasCloseElements:
    def test_empty(self):
        assert has_close_elements([], 0.5) is False

    def test_single(self):
        assert has_close_elements([1.0], 0.5) is False

    def test_no_close_pair(self):
        assert has_close_elements([1.0, 2.0, 3.0], 0.5) is False

    def test_close_pair_exists(self):
        assert has_close_elements([1.0, 2.8, 3.0, 4.0, 5.0, 2.0], 0.3) is True

    def test_exact_threshold_not_close(self):
        # distance = 1.0, threshold = 1.0, condition is distance < threshold
        assert has_close_elements([1.0, 2.0], 1.0) is False

    def test_distance_less_than_threshold(self):
        assert has_close_elements([1.0, 1.5], 1.0) is True

    def test_negative_threshold_returns_false(self):
        assert has_close_elements([1.0, 2.0], -0.1) is False

    def test_zero_threshold_with_duplicates(self):
        # duplicate distance = 0, 0 < 0 is False
        assert has_close_elements([5.0, 5.0], 0.0) is False

    def test_zero_threshold_without_duplicates(self):
        assert has_close_elements([1.0, 2.0], 0.0) is False

    def test_invalid_threshold_type(self):
        with pytest.raises(NumberValidationError):
            has_close_elements([1.0, 2.0], "string")

    def test_bool_threshold_raises(self):
        with pytest.raises(NumberValidationError):
            has_close_elements([1.0, 2.0], True)

    def test_threshold_coercion(self):
        # int threshold is coerced to float
        assert has_close_elements([1.0, 1.9], 1) is True  # 0.9 < 1.0

    def test_numbers_validation_propagates(self):
        with pytest.raises(NumberValidationError):
            has_close_elements([1.0, "invalid"], 0.5)


class TestCountClosePairs:
    def test_empty(self):
        assert count_close_pairs([], 0.5) == 0

    def test_single(self):
        assert count_close_pairs([1.0], 0.5) == 0

    def test_negative_threshold(self):
        assert count_close_pairs([1.0, 2.0], -0.1) == 0

    def test_zero_threshold_no_duplicates(self):
        # sorted: [1,2,3], threshold=0
        # right=0: left=0, count+=0
        # right=1: while left<1 and 2-1 >=0 -> left=1, count+=0
        # right=2: while left<2 and 3-2 >=0 -> left=2, count+=0
        assert count_close_pairs([1.0, 2.0, 3.0], 0.0) == 0

    def test_zero_threshold_with_duplicates(self):
        # sorted: [5,5,5], threshold=0
        # right=0: left=0, count+=0
        # right=1: while left<1 and 5-5 >=0 -> left=1, count+=0
        # right=2: while left<2 and 5-5 >=0 -> left=2, count+=0
        assert count_close_pairs([5.0, 5.0, 5.0], 0.0) == 0

    def test_small_threshold(self):
        # sorted: [1,2,3], threshold=1.5
        # right=0: left=0, count+=0
        # right=1: while left<1 and 2-1 >=1.5? 1>=1.5? False, count+=1-0=1
        # right=2: while left<2 and 3-1 >=1.5? 2>=1.5 True -> left=1
        #          while left<2 and 3-2 >=1.5? 1>=1.5? False, count+=2-1=1
        # total = 2
        assert count_close_pairs([1.0, 2.0, 3.0], 1.5) == 2

    def test_large_threshold(self):
        # sorted: [1,2,3], threshold=10.0
        # All pairs satisfy distance < threshold
        # Number of pairs = n*(n-1)/2 = 3*2/2 = 3
        assert count_close_pairs([1.0, 2.0, 3.0], 10.0) == 3

    def test_exact_threshold_equals_distance(self):
        # sorted: [1,2,3], threshold=1.0
        # right=0: left=0, count+=0
        # right=1: while left<1 and 2-1 >=1.0? 1>=1.0 True -> left=1, count+=0
        # right=2: while left<2 and 3-1 >=1.0? 2>=1.0 True -> left=2, count+=0
        assert count_close_pairs([1.0, 2.0, 3.0], 1.0) == 0

    def test_unsorted_input(self):
        assert count_close_pairs([3.0, 1.0, 2.0], 1.5) == 2

    def test_numbers_validation_propagates(self):
        with pytest.raises(NumberValidationError):
            count_close_pairs([1.0, "invalid"], 0.5)


class TestMinDistance:
    def test_empty(self):
        assert min_distance([]) is None

    def test_single(self):
        assert min_distance([5.0]) is None

    def test_two(self):
        assert min_distance([1.0, 3.0]) == 2.0

    def test_duplicates(self):
        assert min_distance([5.0, 5.0]) == 0.0

    def test_three(self):
        assert min_distance([1.0, 3.0, 2.0]) == 1.0

    def test_numbers_validation_propagates(self):
        with pytest.raises(NumberValidationError):
            min_distance([1.0, "invalid"])