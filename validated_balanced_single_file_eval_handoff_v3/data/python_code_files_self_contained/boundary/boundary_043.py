import math
import itertools
from functools import reduce

def task_func(numbers):
    sum_log_products = 0
    for r in range(1, len(numbers) + 1):
        combinations = itertools.combinations(numbers, r)
        for combination in combinations:
            product = reduce(lambda x, y: x * y, combination)
            sum_log_products += math.log(product)
    return sum_log_products
