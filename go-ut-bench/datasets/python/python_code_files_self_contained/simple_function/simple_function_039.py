import random
import string
import collections
VALID_CHARACTERS = string.ascii_letters + string.digits

def task_func(n_strings, string_length):
    strings = [''.join((random.choice(VALID_CHARACTERS) for _ in range(string_length))) for _ in range(n_strings)]
    character_counts = collections.Counter(''.join(strings))
    return dict(character_counts)
