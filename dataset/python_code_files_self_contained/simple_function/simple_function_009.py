import random
import string
POSSIBLE_LETTERS = ['a', 'b', 'c']

def task_func(word):
    if not all((char in string.ascii_letters for char in word)):
        raise ValueError('Input must only contain letters.')
    if len(word) < 2:
        return ['' for _ in range(len(POSSIBLE_LETTERS))]
    pairs = [''.join(x) for x in zip(word, word[1:])]
    random_pairs = [random.choice(pairs) for _ in range(len(POSSIBLE_LETTERS))]
    return random_pairs
