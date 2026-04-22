from collections import Counter
import itertools
import string

def task_func(word: str) -> dict:
    ALPHABETS = string.ascii_lowercase
    permutations = [''.join(x) for x in itertools.permutations(ALPHABETS, 2)]
    combinations = permutations + [x * 2 for x in ALPHABETS]
    word_combinations = [''.join(x) for x in zip(word, word[1:])]
    word_counter = Counter(word_combinations)
    return {key: word_counter.get(key, 0) for key in combinations}
