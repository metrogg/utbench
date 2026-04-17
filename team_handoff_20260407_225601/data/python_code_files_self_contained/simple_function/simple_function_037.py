import subprocess
import os
import glob
import time

def task_func(test_dir):
    execution_times = {}
    py_scripts = glob.glob(os.path.join(test_dir, '*.py'))
    for py_script in py_scripts:
        start_time = time.time()
        subprocess.call(['python', py_script])
        end_time = time.time()
        execution_times[os.path.basename(py_script)] = end_time - start_time
    return execution_times
