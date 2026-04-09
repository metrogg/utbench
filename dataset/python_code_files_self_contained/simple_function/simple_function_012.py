import csv
import json
import os

def task_func(file_name):
    if not os.path.exists(file_name):
        raise FileNotFoundError('File does not exist.')
    data = []
    with open(file_name, 'r') as f:
        csv_reader = csv.DictReader(f)
        for row in csv_reader:
            data.append(row)
    json_file_name = file_name.split('.')[0] + '.json'
    with open(json_file_name, 'w') as f:
        json.dump(data, f)
    return json_file_name
