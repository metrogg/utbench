from collections import Counter
import math

def task_func(nested_dict):
    counter = Counter()
    for sub_dict in nested_dict.values():
        counter.update(sub_dict)
    counter.pop('ele', None)
    return {k: math.sin(v) for k, v in counter.items()}
