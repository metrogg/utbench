import pytest
import itertools
import sys
from random import shuffle

# Source code under test
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

# Tests
def test_task_func_default_args(monkeypatch):
    """
    Test with default arguments [1, 2].
    Mock shuffle to be a no-op (identity function).
    Permutations: (1, 2) and (2, 1).
    Diffs: |1-2| = 1, |2-1| = 1.
    Sum = 2, Count = 2, Avg = 1.0.
    """
    monkeypatch.setattr(sys.modules[__name__], 'shuffle', lambda x: None)
    assert task_func() == 1.0

def test_task_func_single_element(monkeypatch):
    """
    Test with a single element list [1].
    Mock shuffle to be a no-op.
    Permutations: (1,).
    Diffs: [] (empty list).
    Sum = 0, Count = 1, Avg = 0.0.
    """
    monkeypatch.setattr(sys.modules[__name__], 'shuffle', lambda x: None)
    assert task_func([1]) == 0.0

def test_task_func_custom_list_sort_mock(monkeypatch):
    """
    Test with [1, 2, 3].
    Mock shuffle to sort the list ascending.
    All permutations become [1, 2, 3] after mock.
    Diffs: |1-2| + |2-3| = 1 + 1 = 2.
    Total permutations = 6.
    Sum = 2 * 6 = 12.
    Avg = 12 / 6 = 2.0.
    """
    monkeypatch.setattr(sys.modules[__name__], 'shuffle', lambda x: x.sort())
    assert task_func([1, 2, 3]) == 2.0

def test_task_func_floats(monkeypatch):
    """
    Test with floating point numbers [1.5, 2.5].
    Mock shuffle to be a no-op.
    Permutations: (1.5, 2.5) and (2.5, 1.5).
    Diffs: |1.5-2.5| = 1.0, |2.5-1.5| = 1.0.
    Sum = 2.0, Count = 2, Avg = 1.0.
    """
    monkeypatch.setattr(sys.modules[__name__], 'shuffle', lambda x: None)
    assert task_func([1.5, 2.5]) == 1.0

def test_task_func_reverse_mock(monkeypatch):
    """
    Test with [1, 2, 3].
    Mock shuffle to reverse the list.
    Perm (1, 2, 3) -> [3, 2, 1] -> |3-2|+|2-1| = 2
    Perm (1, 3, 2) -> [2, 3, 1] -> |2-3|+|3-1| = 3
    Perm (2, 1, 3) -> [3, 1, 2] -> |3-1|+|1-2| = 3
    Perm (2, 3, 1) -> [1, 3, 2] -> |1-3|+|3-2| = 3
    Perm (3, 1, 2) -> [2, 1, 3] -> |2-1|+|1-3| = 3
    Perm (3, 2, 1) -> [1, 2, 3] -> |1-2|+|2-3| = 2
    Sum = 2+3+3+3+3+2 = 16
    Avg = 16 / 6.
    """
    def reverse_shuffle(lst):
        lst.reverse()
        
    monkeypatch.setattr(sys.modules[__name__], 'shuffle', reverse_shuffle)
    assert task_func([1, 2, 3]) == 16 / 6