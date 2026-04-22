import json
from glob import glob

def task_func(directory, string):
    json_files = glob(f'{directory}/**/*.json', recursive=True)
    found_files = []
    for file in json_files:
        try:
            with open(file, 'r') as f:
                data = json.load(f)
                if string in data:
                    found_files.append(str(file))
        except (IOError, json.JSONDecodeError):
            continue
    return found_files
