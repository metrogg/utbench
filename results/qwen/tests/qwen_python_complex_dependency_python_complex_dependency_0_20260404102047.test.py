import random
import statistics
import pytest

def test_empty_input_returns_empty_dict():
    random.seed(42)
    result = task_func([])
    assert result == {}

def test_single_key_generates_valid_list():
    random.seed(42)
    result = task_func(['A'])
    assert len(result) == 1
    assert 'A' in result
    values = result['A']
    assert isinstance(values, list)
    assert 1 <= len(values) <= 10
    assert all(0 <= v <= 100 for v in values)

def test_multiple_keys_sorted_by_mean_descending():
    random.seed(100)
    letters = ['A', 'B', 'C', 'D', 'E']
    result = task_func(letters)

    assert set(result.keys()) == set(letters)
    assert len(result) == len(letters)

    means = [statistics.mean(v) for v in result.values()]
    for i in range(len(means) - 1):
        assert means[i] >= means[i + 1]

    for v in result.values():
        assert isinstance(v, list)
        assert 1 <= len(v) <= 10
        assert all(isinstance(x, int) and 0 <= x <= 100 for x in v)

def test_deterministic_behavior_across_runs():
    letters = ['X', 'Y', 'Z']
    random.seed(99)
    res1 = task_func(letters)
    random.seed(99)
    res2 = task_func(letters)
    assert res1 == res2

def test_invalid_input_raises_type_error():
    with pytest.raises(TypeError):
        task_func(None)