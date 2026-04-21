import re
import os
import subprocess

def task_func(dir_path, exe_pattern, execute_files=True):
    results = []
    for dirpath, dirnames, filenames in os.walk(os.path.normpath(dir_path)):
        for filename in filenames:
            if re.search(exe_pattern, filename):
                file_path = os.path.join(dirpath, filename)
                if execute_files:
                    result = subprocess.run([file_path], stdout=subprocess.PIPE)
                    results.append(result.stdout.decode('utf-8'))
                else:
                    results.append(file_path)
    return results
