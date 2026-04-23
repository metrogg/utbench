from collections import deque
import math

def task_func(l):
    if not l:
        return deque()
    dq = deque(l)
    dq.rotate(3)
    numeric_sum = sum((item for item in dq if isinstance(item, (int, float))))
    if numeric_sum > 0:
        print(f'The square root of the sum of numeric elements: {math.sqrt(numeric_sum)}')
    return dq
