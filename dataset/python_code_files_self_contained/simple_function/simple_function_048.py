import csv
import re
from collections import Counter

def task_func(file_path, regex_pattern='\\(.+?\\)|\\w+|[\\W_]+'):
    with open(file_path, 'r') as file:
        reader = csv.reader(file)
        text = ' '.join((row[0] for row in reader))
        matches = re.findall(regex_pattern, text)
    counts = Counter(matches)
    return dict(counts)
