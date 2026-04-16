import random
import string

def task_func(max_length, n_samples, seed=None):
    if max_length < 1:
        raise ValueError('max_length must be larger than or equal to 1.')
    LETTERS = string.ascii_lowercase
    if seed is not None:
        random.seed(seed)
    all_combinations = []
    for i in range(n_samples):
        random_length = random.randint(1, max_length)
        combination = ''.join(random.choices(LETTERS, k=random_length))
        all_combinations.append(combination)
    return all_combinations
