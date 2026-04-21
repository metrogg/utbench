import csv
import sys

def task_func(filename):
    try:
        with open(filename, 'r+') as file:
            reader = csv.reader(file)
            rows = list(reader)
            file.seek(0)
            file.truncate()
            writer = csv.writer(file)
            writer.writerows(reversed(rows))
            file.seek(0)
    except Exception as e:
        print(f'An error occurred: {e}', file=sys.stderr)
    return filename
