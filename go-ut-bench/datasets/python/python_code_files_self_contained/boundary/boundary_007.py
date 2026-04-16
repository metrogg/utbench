import random
import string
from collections import Counter

def task_func(num_strings, string_length):
    strings = [''.join(random.choices(string.ascii_lowercase, k=string_length)) for _ in range(num_strings)]
    characters = ''.join(strings)
    character_counter = Counter(characters)
    most_common_characters = character_counter.most_common()
    return most_common_characters
