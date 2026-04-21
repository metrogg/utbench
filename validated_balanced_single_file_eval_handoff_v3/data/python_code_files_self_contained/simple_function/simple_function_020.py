from collections import defaultdict
from random import randint

def task_func(dict1):
    employee_data = defaultdict(list)
    for prefix, num_employees in dict1.items():
        if not prefix.startswith('EMP$$'):
            continue
        salaries = [randint(1, 100) for _ in range(num_employees)]
        employee_data[prefix].extend(salaries)
    return dict(employee_data)
