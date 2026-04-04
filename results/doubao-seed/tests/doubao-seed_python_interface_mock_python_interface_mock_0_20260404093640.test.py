import itertools
from random import shuffle
import pytest
from unittest.mock import patch, call


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


def test_single_element_list_returns_zero():
    assert task_func([5]) == 0.0


def test_empty_list_returns_zero():
    assert task_func([]) == 0.0


def test_all_duplicate_elements_returns_zero():
    assert task_func([2, 2, 2, 2]) == 0.0


@patch(f"{__name__}.shuffle")
def test_default_input_with_no_op_shuffle_returns_expected(mock_shuffle):
    mock_shuffle.side_effect = lambda x: x
    result = task_func()
    assert result == 1.0
    assert isinstance(result, float)
    assert mock_shuffle.call_count == 2
    mock_shuffle.assert_has_calls([call([1, 2]), call([2, 1])], any_order=True)


@patch(f"{__name__}.shuffle")
def test_three_elements_with_sorted_shuffle_returns_expected(mock_shuffle):
    def mock_sort(lst):
        lst.sort()
    mock_shuffle.side_effect = mock_sort
    result = task_func([1, 2, 3])
    assert result == 2.0
    assert mock_shuffle.call_count == 6


@patch(f"{__name__}.shuffle")
def test_negative_numbers_returns_correct_average(mock_shuffle):
    mock_shuffle.side_effect = lambda x: x
    result = task_func([-1, 1])
    assert result == 2.0


@patch(f"{__name__}.shuffle")
def test_float_numbers_returns_correct_average(mock_shuffle):
    def mock_sort(lst):
        lst.sort()
    mock_shuffle.side_effect = mock_sort
    result = task_func([1.5, 2.5, 3.5])
    assert result == 2.0
    assert isinstance(result, float)


@patch(f"{__name__}.shuffle")
def test_shuffle_called_correct_number_of_times(mock_shuffle):
    mock_shuffle.side_effect = lambda x: x
    task_func([1, 2, 3, 4])
    assert mock_shuffle.call_count == 24  # 4! = 24


def test_non_numeric_elements_raises_type_error():
    with pytest.raises(TypeError):
        task_func([1, "a", 3])


def test_non_iterable_input_raises_type_error():
    with pytest.raises(TypeError):
        task_func(None)