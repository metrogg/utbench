import json
from collections import Counter
import random

def task_func(my_dict, keys):
    if len(set(keys)) != 10:
        raise ValueError('keys parameter must contain exactly 10 unique elements')
    for key in keys:
        my_dict[key] = random.randint(1, 100)
    json_filename = 'updated_dictionary.json'
    txt_filename = 'key_frequencies.txt'
    with open(json_filename, 'w') as json_file:
        json.dump(my_dict, json_file, indent=4)
    key_counts = Counter(my_dict.keys())
    with open(txt_filename, 'w') as txt_file:
        for key, count in key_counts.items():
            txt_file.write(f'{key}: {count}\n')
    return (my_dict, json_filename, txt_filename)
