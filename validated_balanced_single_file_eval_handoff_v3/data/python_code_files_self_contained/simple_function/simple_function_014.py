from collections import Counter
import random
from itertools import cycle
ELEMENTS = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J']

def task_func(l):
    if not l:
        return Counter()
    random.shuffle(l)
    l_cycled = cycle(l)
    counter = Counter((next(l_cycled) for _ in range(30)))
    keys = list(counter.keys())
    counter = Counter({k: counter[k] for k in keys[3:] + keys[:3]})
    return counter
