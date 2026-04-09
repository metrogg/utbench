import csv
from collections import Counter
import operator

def task_func(csv_file, csv_delimiter):
    words = []
    with open(csv_file, 'r') as f:
        reader = csv.reader(f, delimiter=csv_delimiter)
        for row in reader:
            words.extend(row)
    word_counter = Counter(words)
    most_common_words = sorted(word_counter.items(), key=operator.itemgetter(1), reverse=True)
    return most_common_words
