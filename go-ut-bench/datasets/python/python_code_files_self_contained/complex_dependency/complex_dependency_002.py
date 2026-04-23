import subprocess
import csv
import glob
import random
import os

def task_func(file):
    if not os.path.exists(file):
        print('Provided file does not exist.')
        return []
    if not file.endswith('.csv'):
        print('Provided file is not a CSV.')
        return []
    try:
        subprocess.call(['split', '-n', '5', '-d', file, 'split_'])
        split_files = glob.glob('split_*')
        for split_file in split_files:
            with open(split_file, 'r') as f:
                reader = csv.reader(f)
                rows = list(reader)
            random.shuffle(rows)
            with open(split_file, 'w') as f:
                writer = csv.writer(f)
                writer.writerows(rows)
        return split_files
    except Exception as e:
        print(f'An error occurred: {e}')
        return []
