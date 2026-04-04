import pytest
import random
import statistics
# Adjust import path as needed to match your project structure
from your_module import task_func


def test_task_func_empty_input_returns_empty_dict():
    assert task_func([]) == {}


def test_task_func_single_letter_input_returns_valid_entry():
    random.seed(42)
    result = task_func(["a"])
    assert list(result.keys()) == ["a"]
    value = result["a"]
    assert isinstance(value, list)
    assert 1 <= len(value) <= 10
    assert all(isinstance(num, int) and 0 <= num <= 100 for num in value)


def test_task_func_multiple_entries_sorted_by_mean_descending():
    random.seed(42)
    test_letters = ["x", "y", "z"]
    result = task_func(test_letters)
    assert set(result.keys()) == set(test_letters)
    values = list(result.values())
    means = [statistics.mean(val) for val in values]
    assert means == sorted(means, reverse=True)
    for val in values:
        assert 1 <= len(val) <= 10
        assert all(isinstance(num, int) and 0 <= num <= 100 for num in val)


def test_task_func_duplicate_input_letters_return_unique_keys():
    random.seed(99)
    test_letters = ["a", "a", "b", "c", "b", "d"]
    result = task_func(test_letters)
    assert set(result.keys()) == {"a", "b", "c", "d"}


def test_task_func_non_string_hashable_keys_work_correctly():
    random.seed(123)
    test_keys = [1, 2, 3, 4]
    result = task_func(test_keys)
    assert set(result.keys()) == set(test_keys)
    means = [statistics.mean(v) for v in result.values()]
    assert means == sorted(means, reverse=True)


def test_task_func_unhashable_input_elements_raise_type_error():
    with pytest.raises(TypeError):
        task_func([1, [2, 3], "a", {"key": "val"}])