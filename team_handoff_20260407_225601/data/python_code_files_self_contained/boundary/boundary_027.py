import collections
import random
import itertools
ANIMALS = ['Cat', 'Dog', 'Elephant', 'Lion', 'Tiger', 'Bear', 'Giraffe', 'Horse', 'Rabbit', 'Snake', 'Zebra']

def task_func(animal_dict, max_count=10, seed=0):
    if max_count < 1:
        raise ValueError('max_count must be a positive integer')
    random.seed(seed)
    reversed_dict = {v: [] for v in animal_dict.values() if isinstance(v, str) and v in ANIMALS}
    for k, v in animal_dict.items():
        if isinstance(v, str) and v in ANIMALS:
            reversed_dict[v].append(k)
    animal_counter = collections.Counter(itertools.chain.from_iterable([[v] * random.randint(1, max_count) for v in ANIMALS]))
    return (reversed_dict, animal_counter)
