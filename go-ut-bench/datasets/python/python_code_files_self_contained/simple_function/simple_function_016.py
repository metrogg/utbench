from itertools import combinations
import math

def task_func(x, w):
    max_weight = -math.inf
    max_substr = ''
    for start, end in combinations(range(len(x) + 1), 2):
        substr = x[start:end]
        weight = sum((w.get(c, 0) for c in substr))
        if weight > max_weight:
            max_weight = weight
            max_substr = substr
    return max_substr
