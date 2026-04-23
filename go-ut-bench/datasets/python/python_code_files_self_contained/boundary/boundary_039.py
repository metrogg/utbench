from collections import Counter
import itertools
from random import randint

def task_func(T1, RANGE=100):
    int_list = [list(map(int, x)) for x in T1]
    flattened_list = list(itertools.chain(*int_list))
    total_nums = sum(flattened_list)
    random_nums = [randint(0, RANGE) for _ in range(total_nums)]
    counts = Counter(random_nums)
    return counts
