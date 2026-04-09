import random
from itertools import combinations
import math

def task_func(n):
    if n < 2:
        return None
    points = [(random.random(), random.random()) for i in range(n)]
    closest_pair = min(combinations(points, 2), key=lambda pair: math.hypot(pair[0][0] - pair[1][0], pair[0][1] - pair[1][1]))
    return closest_pair
