import re
import os
import glob

def task_func(dir_path):
    if not os.path.isdir(dir_path):
        raise ValueError('Specified directory does not exist.')
    result = {}
    file_paths = glob.glob(f'{dir_path}/**/*.txt', recursive=True)
    for file_path in file_paths:
        with open(file_path, 'r') as file:
            content = file.read()
        matches = re.findall('\\berror\\b', content, re.IGNORECASE)
        result[os.path.relpath(file_path, dir_path)] = len(matches)
    return result
