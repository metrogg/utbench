import subprocess
import os
import time
import glob

def task_func(r_script_path: str, output_path: str, duration: int) -> (bool, str):
    command = f'/usr/bin/Rscript --vanilla {r_script_path}'
    subprocess.call(command, shell=True)
    start_time = time.time()
    search_pattern = os.path.join(output_path, '*.csv')
    while time.time() - start_time < duration:
        if glob.glob(search_pattern):
            return (True, 'File generated successfully within the specified duration.')
        time.sleep(0.1)
    return (False, 'File not generated within the specified duration.')
