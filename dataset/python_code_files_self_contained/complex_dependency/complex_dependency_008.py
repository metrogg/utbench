import os
import random

def task_func(directory, n_files):
    if not os.path.exists(directory):
        os.makedirs(directory)
    for i in range(n_files):
        filename = os.path.join(directory, f'file_{i + 1}.txt')
        with open(filename, 'w') as file:
            file.write(str(random.randint(0, 9)))
            file.seek(0)
    return n_files
