import json
import os

def task_func(filename, data):
    try:
        with open(filename, 'w') as f:
            json.dump(data, f)
        file_exists = os.path.exists(filename)
        if not file_exists:
            return (False, None)
        with open(filename, 'r') as f:
            written_data = json.load(f)
            if written_data != data:
                return (False, None)
        return (True, written_data)
    except Exception as e:
        return (False, None)
