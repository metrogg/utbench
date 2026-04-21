import json
import csv
import os
import base64

def task_func(raw_string, filename, output_dir):
    decoded_string = base64.b64decode(raw_string).decode('utf-8')
    data = json.loads(decoded_string)
    os.makedirs(output_dir, exist_ok=True)
    file_path = os.path.join(output_dir, f'{filename}.csv')
    with open(file_path, 'w', newline='') as f:
        writer = csv.writer(f)
        for key, value in data.items():
            writer.writerow([key, value])
    return file_path
