import collections
import random
import json
PREFICES = ['EMP$$', 'MAN$$', 'DEV$$', 'HR$$']
LEVELS = ['Junior', 'Mid', 'Senior']

def task_func(department_data):
    level_data = collections.defaultdict(list)
    for prefix, num_employees in department_data.items():
        if prefix not in PREFICES:
            continue
        for _ in range(num_employees):
            level = random.choice(LEVELS)
            level_data[prefix].append(level)
    return json.dumps(level_data)
