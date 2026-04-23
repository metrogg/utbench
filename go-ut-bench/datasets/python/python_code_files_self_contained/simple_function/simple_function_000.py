import itertools
import math

def task_func(x):
    min_length = math.inf
    min_subseq = []
    for r in range(1, len(x) + 1):
        for subseq in itertools.combinations(x.items(), r):
            length = sum((length for letter, length in subseq))
            if length < min_length:
                min_length = length
                min_subseq = [letter for letter, length in subseq]
    return min_subseq
