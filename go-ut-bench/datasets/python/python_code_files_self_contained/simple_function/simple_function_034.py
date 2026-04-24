import re
import random
import string

def task_func(n, pattern, seed=None):
    if seed is not None:
        random.seed(seed)
    rand_str = ''.join((random.choice(string.ascii_letters + string.digits) for _ in range(n)))
    matches = re.findall(pattern, rand_str)
    return matches
