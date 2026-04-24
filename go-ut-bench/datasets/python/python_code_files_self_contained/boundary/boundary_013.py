import collections
from itertools import zip_longest
from random import choices

def task_func(l1, l2, K=10):
    combined = [val for pair in zip_longest(l1, l2) for val in pair if val is not None]
    sample = choices(combined, k=K)
    freq = collections.Counter(sample)
    return freq
