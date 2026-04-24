import os
import re
import json
import glob

def task_func(directory_path: str) -> list:
    if not os.path.exists(directory_path):
        raise FileNotFoundError(f'Directory {directory_path} not found.')
    json_files = glob.glob(directory_path + '/*.json')
    processed_files = []
    for json_file in json_files:
        with open(json_file, 'r') as file:
            data = json.load(file)
        escaped_data = json.dumps(data, ensure_ascii=False)
        escaped_data = re.sub('(?<!\\\\)"', '\\\\\\"', escaped_data)
        with open(json_file, 'w') as file:
            file.write(escaped_data)
        processed_files.append(json_file)
    return processed_files
