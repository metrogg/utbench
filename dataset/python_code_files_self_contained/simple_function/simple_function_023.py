import os
import random
import json

def task_func(directory, n):
    if not os.path.exists(directory):
        os.makedirs(directory)
    for i in range(n):
        filename = str(i) + '.json'
        filepath = os.path.join(directory, filename)
        with open(filepath, 'w') as file:
            json.dump({'number': random.randint(1, 100)}, file)
            file.seek(0)
    return directory
