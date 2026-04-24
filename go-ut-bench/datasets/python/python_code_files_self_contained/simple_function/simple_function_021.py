import re
import string
from random import choice

def task_func(n, pattern):
    while True:
        s = ''.join((choice(string.ascii_letters) for _ in range(n)))
        if re.match(pattern, s):
            return s
