import collections
import random
import string

def task_func(length=100):
    if length < 0:
        raise ValueError
    random_string = ''.join(random.choices(string.ascii_uppercase + string.ascii_lowercase, k=length))
    char_counts = collections.Counter(random_string)
    return dict(char_counts)
