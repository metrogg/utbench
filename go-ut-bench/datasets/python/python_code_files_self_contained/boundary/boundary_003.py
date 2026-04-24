import random
import string
LETTERS = string.ascii_letters

def task_func(num_words, word_length):
    if num_words < 0 or word_length < 0:
        raise ValueError('num_words and word_length must be non-negative')
    random.seed(42)
    words = [''.join((random.choice(LETTERS) for _ in range(word_length))) for _ in range(num_words)]
    return words
