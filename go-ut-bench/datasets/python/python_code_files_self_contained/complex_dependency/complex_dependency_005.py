import os
import os.path
import csv
import collections
FILE_NAME = 'file_sizes.csv'

def task_func(my_path):
    file_sizes = collections.defaultdict(int)
    for dirpath, dirnames, filenames in os.walk(my_path):
        for f in filenames:
            fp = os.path.join(dirpath, f)
            file_sizes[f] += os.path.getsize(fp)
    with open(os.path.join(my_path, FILE_NAME), 'w', newline='') as csvfile:
        writer = csv.writer(csvfile)
        writer.writerow(['File Name', 'Size'])
        for row in file_sizes.items():
            writer.writerow(row)
    return os.path.join(my_path, FILE_NAME)
