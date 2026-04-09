import heapq
import random

def task_func(obj_list, attr, top_n=5, seed=None):
    random.seed(seed)
    attr_values = [getattr(obj, attr) for obj in obj_list]
    if len(attr_values) == 0:
        return ([], None)
    top_values = heapq.nlargest(top_n, attr_values)
    random_value = random.choice(attr_values)
    return (top_values, random_value)
