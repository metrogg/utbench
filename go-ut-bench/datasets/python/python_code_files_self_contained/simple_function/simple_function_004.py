from itertools import combinations
import math

def task_func(seq, letter_weight_dict):
    max_weight = -math.inf
    max_subseq = ''
    for r in range(1, len(seq) + 1):
        for subseq in combinations(seq, r):
            weight = sum((letter_weight_dict[c] for c in subseq))
            if weight > max_weight:
                max_weight = weight
                max_subseq = ''.join(subseq)
    return max_subseq
