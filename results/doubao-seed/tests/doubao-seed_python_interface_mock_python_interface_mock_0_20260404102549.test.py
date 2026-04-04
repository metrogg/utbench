import pytest
from unittest.mock import patch, call
import itertools
import random

# Include the source function for self-contained test execution
import itertools
from random import shuffle
def task_func(numbers=list(range(1, 3))):
    permutations = list(itertools.permutations(numbers))
    sum_diffs = 0
    for perm in permutations:
        perm = list(perm)
        shuffle(perm)
        diffs = [abs(perm[i] - perm[i+1]) for i in range(len(perm)-1)]
        sum_diffs += sum(diffs)
    avg_sum_diffs = sum_diffs / len(permutations)
    return avg_sum_diffs


def test_single_element_list():
    result = task_func([9])
    assert isinstance(result, float)
    assert result == 0.0


def test_empty_list_input():
    result = task_func([])
    assert result == 0.0


def test_all_duplicate_elements():
    numbers = [4, 4, 4, 4]
    with patch('random.shuffle'):
        result = task_func(numbers)
        assert result == 0.0


def test_default_input_mocked_shuffle_no_op():
    with patch('random.shuffle', side_effect=lambda x: x) as mock_shuffle:
        result = task_func()
        assert result == 1.0
        assert mock_shuffle.call_count == 2
        mock_shuffle.assert_has_calls([call([1, 2]), call([2, 1])], any_order=True)


def test_three_elements_mocked_shuffle_no_op():
    numbers = [1, 2, 3]
    expected_avg = 16 / 6
    with patch('random.shuffle', side_effect=lambda x: x):
        result = task_func(numbers)
        assert pytest.approx(result) == expected_avg


def test_shuffle_called_correct_number_of_times():
    numbers = [1, 2, 3, 4]
    expected_call_count = len(list(itertools.permutations(numbers)))
    with patch('random.shuffle', side_effect=lambda x: x) as mock_shuffle:
        task_func(numbers)
        assert mock_shuffle.call_count == expected_call_count


def test_non_numeric_elements_raises_type_error():
    invalid_numbers = [1, "string", 3.5]
    with pytest.raises(TypeError):
        task_func(invalid_numbers)


def test_mocked_shuffle_reverse_permutation():
    def reverse_list(lst):
        lst.reverse()
    numbers = [1, 2, 3]
    expected_avg = 16 / 6
    with patch('random.shuffle', side_effect=reverse_list):
        result = task_func(numbers)
        assert pytest.approx(result) == expected_avg