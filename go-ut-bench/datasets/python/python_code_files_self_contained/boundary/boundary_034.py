import json
import os

def task_func(file_path):
    if not os.path.exists(file_path):
        return False
    with open(file_path, 'r') as file:
        try:
            data = json.load(file)
        except json.JSONDecodeError:
            return False
    return isinstance(data, list) and all((isinstance(item, dict) for item in data))
