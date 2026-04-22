import subprocess
from multiprocessing import Pool

def execute_command(command):
    process = subprocess.Popen(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE, shell=True)
    output, _ = process.communicate()
    return output

def task_func(commands):
    if not commands:
        return []
    with Pool(processes=len(commands)) as pool:
        outputs = pool.map(execute_command, commands)
    return outputs
