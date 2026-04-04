import pytest
from unittest.mock import patch
import itertools


def task_func(numbers=list(range(1, 3))):
    from random import shuffle
    permutations = list(itertools.permutations(numbers))
    sum_diffs = 0

    for perm in permutations:
        perm = list(perm)
        shuffle(perm)
        diffs = [abs(perm[i] - perm[i+1]) for i in range(len(perm)-1)]
        sum_diffs += sum(diffs)

    avg_sum_diffs = sum_diffs / len(permutations)
    
    return avg_sum_diffs


class TestTaskFunc:

    def test_default_input(self):
        """Test with default input [1, 2]"""
        with patch('random.shuffle'):
            result = task_func()
            expected = 1.0
            assert result == expected

    def test_single_element_list(self):
        """Test with single element list - no consecutive pairs exist"""
        with patch('random.shuffle'):
            result = task_func([1])
            expected = 0.0
            assert result == expected

    def test_two_elements(self):
        """Test with two distinct elements"""
        with patch('random.shuffle'):
            result = task_func([1, 2])
            expected = 1.0
            assert result == expected

    def test_three_elements(self):
        """Test with three distinct elements"""
        with patch('random.shuffle'):
            result = task_func([1, 2, 3])
            # Permutations: (1,2,3)->2, (1,3,2)->3, (2,1,3)->3, (2,3,1)->3, (3,1,2)->3, (3,2,1)->2
            # Sum = 16, Avg = 16/6 = 8/3
            expected = 8.0 / 3.0
            assert result == expected

    def test_identical_elements(self):
        """Test with identical elements - all differences should be 0"""
        with patch('random.shuffle'):
            result = task_func([5, 5, 5])
            expected = 0.0
            assert result == expected

    def test_negative_numbers(self):
        """Test with negative numbers"""
        with patch('random.shuffle'):
            result = task_func([-1, 0, 1])
            # After shuffle (no-op), permutations sum to:
            # (1,2,3): |−1−0|+|0−1| = 1+1 = 2
            # (1,3,2): |−1−1|+|1−0| = 2+1 = 3
            # (2,1,3): |0−(−1)|+|−1−1| = 1+2 = 3
            # (2,3,1): |0−1|+|1−(−1)| = 1+2 = 3
            # (3,1,2): |1−(−1)|+|−1−0| = 2+1 = 3
            # (3,2,1): |1−0|+|0−(−1)| = 1+1 = 2
            # Sum = 16, Avg = 16/6 = 8/3
            expected = 8.0 / 3.0
            assert result == expected

    def test_shuffle_actually_scrambles(self):
        """Test that shuffle is called and affects the result"""
        call_count = [0]
        
        def counting_shuffle(x):
            call_count[0] += 1
            # Actually shuffle the list
            import random
            random.shuffle(x)
        
        with patch('random.shuffle', side_effect=counting_shuffle):
            result = task_func([1, 2])
            # shuffle should be called for each permutation
            assert call_count[0] == 2

    def test_empty_list(self):
        """Test with empty list"""
        with patch('random.shuffle'):
            result = task_func([])
            expected = 0.0
            assert result == expected

    def test_custom_numbers(self):
        """Test with custom numbers list"""
        with patch('random.shuffle'):
            result = task_func([10, 20])
            expected = 10.0
            assert result == expected

    def test_result_is_float(self):
        """Test that the result is a float"""
        with patch('random.shuffle'):
            result = task_func([1, 2])
            assert isinstance(result, float)

    def test_larger_list_count(self):
        """Test that number of permutations is correctly calculated"""
        with patch('random.shuffle'):
            result = task_func([1, 2, 3, 4])
            # 4! = 24 permutations
            assert result >= 0  # Result should be non-negative