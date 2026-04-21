import json
import os
PREFIXES = ['is_', 'has_', 'can_', 'should_']

def task_func(directory):
    stats = {prefix: 0 for prefix in PREFIXES}
    for filename in os.listdir(directory):
        if filename.endswith('.json'):
            with open(f'{directory}/{filename}', 'r') as f:
                data = json.load(f)
            for key in data.keys():
                for prefix in PREFIXES:
                    if key.startswith(prefix):
                        stats[prefix] += 1
    return stats
