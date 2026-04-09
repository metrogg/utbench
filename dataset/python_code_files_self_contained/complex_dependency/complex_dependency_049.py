import re
import os
from pathlib import Path
import glob

def task_func(directory_path: str, regex_pattern: str='\\\\(.+?\\\\)|\\\\w') -> dict:
    FILE_PATTERN = '*.txt'
    match_dict = {}
    file_paths = glob.glob(os.path.join(directory_path, FILE_PATTERN))
    for file_path in file_paths:
        with open(file_path, 'r') as file:
            content = file.read()
            matches = re.findall(regex_pattern, content)
            match_dict[Path(file_path).name] = matches
    return match_dict
