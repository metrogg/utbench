import subprocess
import os
import random

def task_func(directory, file_list):
    if not file_list:
        return None
    file = random.choice(file_list)
    file_path = os.path.join(directory, file)
    try:
        process = subprocess.Popen(file_path)
        process.wait()
        return process.returncode
    except Exception as e:
        return None
