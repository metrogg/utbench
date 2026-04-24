import random
from string import ascii_uppercase

def task_func(dict1):
    employee_ids = []
    for prefix, num_employees in dict1.items():
        for _ in range(num_employees):
            random_str = ''.join((random.choice(ascii_uppercase) for _ in range(5)))
            employee_ids.append(f'{prefix}{random_str}')
    return employee_ids
