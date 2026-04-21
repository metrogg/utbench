import random
import bisect
from array import array

def task_func(n=10, total=100):
    nums = []
    while sum(nums) != total:
        nums = [random.randint(0, total) for _ in range(n)]
    nums.sort()
    nums = array('i', nums)
    new_num = random.randint(0, total)
    pos = bisect.bisect(nums, new_num)
    return (nums, pos)
