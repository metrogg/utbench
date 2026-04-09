import json
import os
import glob
KEY = 'mynewkey'
VALUE = 'mynewvalue'

def task_func(directory):
    files = glob.glob(os.path.join(directory, '*.json'))
    updated_files = 0
    for file in files:
        with open(file, 'r+') as f:
            data = json.load(f)
            if KEY not in data:
                data[KEY] = VALUE
                f.seek(0)
                f.truncate()
                json.dump(data, f)
                updated_files += 1
    return updated_files
