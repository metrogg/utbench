import pytest
import random
import statistics
from unittest.mock import patch

# Assume task_func is imported from the module where it's defined, adjust import as needed
from your_module import task_func


def test_task_func_empty_letters():
    result = task_func([])
    assert isinstance(result, dict)
    assert result == {}


def test_task_func_single_letter_input():
    with patch("random.randint") as mock_randint:
        # Side effect sequence: list length=2, elements 50,70
        mock_randint.side_effect = [2, 50, 70]
        result = task_func(["x"])
        
        assert list(result.keys()) == ["x"]
        assert result["x"] == [50, 70]
        assert statistics.mean(result["x"]) == 60
        assert mock_randint.call_count == 3
        mock_randint.assert_any_call(1, 10)
        mock_randint.assert_any_call(0, 100)


def test_task_func_multiple_letters_sorted_by_mean_desc():
    test_letters = ["a", "b", "c"]
    # Side effect sequence:
    # a: len=2, values [10,30] → mean=20
    # b: len=3, values [40,50,60] → mean=50
    # c: len=1, value [10] → mean=10
    side_effects = [2, 10, 30, 3, 40, 50, 60, 1, 10]
    
    with patch("random.randint") as mock_randint:
        mock_randint.side_effect = side_effects
        result = task_func(test_letters)
        
        assert list(result.keys()) == ["b", "a", "c"]
        assert result["b"] == [40, 50, 60]
        assert result["a"] == [10, 30]
        assert result["c"] == [10]
        means = [statistics.mean(v) for v in result.values()]
        assert means == sorted(means, reverse=True)


def test_task_func_same_mean_preserves_input_order():
    test_letters = ["p", "q", "r"]
    # All lists have mean=20
    side_effects = [2, 20, 20, 1, 20, 2, 10, 30]
    
    with patch("random.randint") as mock_randint:
        mock_randint.side_effect = side_effects
        result = task_func(test_letters)
        
        assert list(result.keys()) == test_letters
        assert all(statistics.mean(v) == 20 for v in result.values())


def test_task_func_output_value_properties():
    test_letters = ["d", "e", "f", "g"]
    result = task_func(test_letters)
    
    assert set(result.keys()) == set(test_letters)
    for val in result.values():
        assert isinstance(val, list)
        assert 1 <= len(val) <= 10
        for num in val:
            assert isinstance(num, int)
            assert 0 <= num <= 100
    # Verify descending sort of means
    means = [statistics.mean(v) for v in result.values()]
    assert means == sorted(means, reverse=True)


def test_task_func_non_hashable_letter_raises_type_error():
    with pytest.raises(TypeError):
        task_func(["a", [1, 2], "c"])