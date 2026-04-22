import string
import random

def task_func(text, seed=None):

    def replace_with_random_char(c):
        if c.isalpha():
            if c.islower():
                return random.choice(string.ascii_lowercase)
            else:
                return random.choice(string.ascii_uppercase)
        return c
    if seed is not None:
        random.seed(seed)
    return ''.join((replace_with_random_char(c) for c in text))
