import ast
import json
from collections import Counter

def task_func(file_pointer):
    data = json.load(file_pointer)
    key_frequency_counter = Counter()
    for item in data:
        if isinstance(item, str):
            try:
                item = ast.literal_eval(item)
            except ValueError:
                continue
        if isinstance(item, dict):
            key_frequency_counter.update(item.keys())
    return key_frequency_counter
