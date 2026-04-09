import random
from collections import Counter
from statistics import mode

def task_func(list_length=1000, range_start=1, range_end=10, random_seed=None):
    random.seed(random_seed)
    random_list = [random.randint(range_start, range_end) for _ in range(list_length)]
    counter = Counter(random_list)
    numbers = ((number, count) for number, count in counter.items())
    return (mode(random_list), numbers)
