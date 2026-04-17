import os
import csv
from collections import Counter

def task_func(file_path):
    if not os.path.isfile(file_path):
        return None
    word_counter = Counter()
    with open(file_path, 'r') as f:
        csv_reader = csv.reader(f, delimiter=',', skipinitialspace=True)
        for row in csv_reader:
            for word in row:
                word_counter[word.strip()] += 1
    if not word_counter:
        return None
    most_common_word, frequency = word_counter.most_common(1)[0]
    return (most_common_word, frequency)
