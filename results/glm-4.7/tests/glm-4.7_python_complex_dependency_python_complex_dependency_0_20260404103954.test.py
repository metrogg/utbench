import pytest
from unittest.mock import patch
import random
import statistics

# Source code included to ensure the test is runnable as a standalone script
def task_func(LETTERS):
    random_dict = {k: [random.randint(0, 100) for _ in range(random.randint(1, 10))] for k in LETTERS}
    sorted_dict = dict(sorted(random_dict.items(), key=lambda item: statistics.mean(item[1]), reverse=True))
    return sorted_dict

def test_task_func_normal_sorting():
    """
    Test that the dictionary is sorted correctly by the mean of values in descending order.
    """
    letters = ['a', 'b', 'c']
    
    # Mocking random.randint sequence:
    # 'a': length 2, values [10, 20] -> mean 15.0
    # 'b': length 1, values [30]    -> mean 30.0
    # 'c': length 3, values [1, 2, 3] -> mean 2.0
    # Expected order: 'b' (30.0), 'a' (15.0), 'c' (2.0)
    side_effects = [
        2, 10, 20,    # For 'a'
        1, 30,        # For 'b'
        3, 1, 2, 3    # For 'c'
    ]
    
    with patch('random.randint', side_effect=side_effects):
        result = task_func(letters)
        
        assert list(result.keys()) == ['b', 'a', 'c']
        assert result['a' ] == [10, 20]
        assert result['b'] == [30]
        assert result['c'] == [1, 2, 3]

def test_task_func_empty_input():
    """
    Test that an empty input list returns an empty dictionary.
    """
    with patch('random.randint') as mock_rand:
        result = task_func([])
        
        assert result == {}
        # Ensure no random calls were made since the loop doesn't execute
        mock_rand.assert_not_called()

def test_task_func_single_element():
    """
    Test behavior with a single letter input.
    """
    letters = ['x']
    
    # Mocking:
    # 'x': length 1, value [50] -> mean 50.0
    side_effects = [1, 50]
    
    with patch('random.randint', side_effect=side_effects):
        result = task_func(letters)
        
        assert result == {'x': [50]}

def test_task_func_equal_means():
    """
    Test sorting stability when means are equal.
    Python's sort is stable, so the original order should be preserved.
    """
    letters = ['a', 'b']
    
    # Mocking:
    # 'a': length 2, values [5, 5] -> mean 5.0
    # 'b': length 2, values [5, 5] -> mean 5.0
    side_effects = [2, 5, 5, 2, 5, 5]
    
    with patch('random.randint', side_effect=side_effects):
        result = task_func(letters)
        
        # Check that both are present and order is preserved (a then b)
        assert list(result.keys()) == ['a', 'b']
        assert result == {'a': [5, 5], 'b': [5, 5]}

def test_task_func_boundary_values():
    """
    Test with boundary values for random integers (0 and 100).
    """
    letters = ['min', 'max']
    
    # Mocking:
    # 'min': length 1, value [0] -> mean 0.0
    # 'max': length 1, value [100] -> mean 100.0
    side_effects = [1, 0, 1, 100]
    
    with patch('random.randint', side_effect=side_effects):
        result = task_func(letters)
        
        assert list(result.keys()) == ['max', 'min']
        assert result == {'max': [100], 'min': [0]}

def test_task_func_structure():
    """
    Verify the structure of the returned dictionary and values.
    """
    letters = ['k1', 'k2']
    
    # Mocking:
    # 'k1': length 3, values [10, 20, 30]
    # 'k2': length 2, values [5, 5]
    side_effects = [3, 10, 20, 30, 2, 5, 5]
    
    with patch('random.randint', side_effect=side_effects):
        result = task_func(letters)
        
        assert isinstance(result, dict)
        assert set(result.keys()) == set(letters)
        
        # Check types of values
        for key in result:
            assert isinstance(result[key], list)
            for val in result[key]:
                assert isinstance(val, int)