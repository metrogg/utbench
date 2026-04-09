import os
import glob
import csv

def task_func(directory_path, file_extension='.csv'):
    data = {}
    for file in glob.glob(os.path.join(directory_path, '*' + file_extension)):
        filename = os.path.splitext(os.path.basename(file))[0]
        with open(file, 'r') as f:
            reader = csv.reader(f)
            data[filename] = list(reader)
    return data
