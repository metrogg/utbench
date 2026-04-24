import random
import statistics
AGE_RANGE = (22, 60)

def task_func(dict1):
    emp_ages = []
    for prefix, num_employees in dict1.items():
        if not prefix.startswith('EMP$$'):
            continue
        for _ in range(num_employees):
            age = random.randint(*AGE_RANGE)
            emp_ages.append(age)
    if not emp_ages:
        return (0, 0, [])
    mean_age = statistics.mean(emp_ages)
    median_age = statistics.median(emp_ages)
    mode_age = statistics.multimode(emp_ages)
    return (mean_age, median_age, mode_age)
