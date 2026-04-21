import math
import random
import statistics
RADIUS = 5

def task_func(n):
    distances = []
    for _ in range(n):
        theta = 2 * math.pi * random.random()
        r = RADIUS * math.sqrt(random.random())
        x = r * math.cos(theta)
        y = r * math.sin(theta)
        distance = math.sqrt(x ** 2 + y ** 2)
        distances.append(distance)
    return round(statistics.mean(distances), 4)
