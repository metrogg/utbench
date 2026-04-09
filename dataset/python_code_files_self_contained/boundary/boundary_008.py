from multiprocessing import Pool
import math

def calculate_factorial(number: int) -> tuple:
    return (number, math.factorial(number))

def task_func(numbers: list) -> dict:
    if not all((isinstance(n, int) and n >= 0 for n in numbers)):
        raise ValueError('All elements in the list must be integers')
    with Pool() as pool:
        factorial_dict = dict(pool.starmap(calculate_factorial, [(i,) for i in numbers]))
    return factorial_dict
