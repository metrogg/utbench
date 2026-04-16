import collections
import json
import os

def task_func(directory_path: str) -> dict:
    key_counts = collections.defaultdict(int)
    for filename in os.listdir(directory_path):
        if filename.endswith('.json'):
            file_path = os.path.join(directory_path, filename)
            with open(file_path, 'r') as json_file:
                data = json.load(json_file)
                for key in data.keys():
                    key_counts[key] += 1
    return dict(key_counts)
