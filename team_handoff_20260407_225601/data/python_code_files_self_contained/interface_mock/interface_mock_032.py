from collections import defaultdict
import itertools
import json
import random

def task_func(LETTERS, n):
    combinations = list(itertools.combinations(LETTERS, n))
    letter_counts = defaultdict(int)
    for combination in combinations:
        for letter in combination:
            letter_counts[letter] += 1
    filename = f'letter_combinations_{random.randint(1, 100)}.json'
    with open(filename, 'w') as f:
        json.dump(letter_counts, f)
    return filename
