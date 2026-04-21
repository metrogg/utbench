import random
from collections import Counter

def task_func(strings: list) -> dict:
    if not strings:
        return Counter()
    pattern = '}'
    random_choices = random.choices(strings, k=10)
    pattern_counts = Counter([string.count(pattern) for string in random_choices])
    return pattern_counts
