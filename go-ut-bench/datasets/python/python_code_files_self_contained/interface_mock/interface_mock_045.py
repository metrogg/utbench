from collections import Counter
import os
import csv
FILE_DIR = './yourdictfiles/'

def task_func(output_file, test_directory):
    total_words = 0
    try:
        word_counts = Counter()
        for file_name in os.listdir(test_directory):
            if not file_name.endswith('.txt'):
                continue
            with open(os.path.join(test_directory, file_name), 'r') as file:
                words = file.read().split()
                word_counts.update(words)
        with open(output_file, 'w') as file:
            writer = csv.writer(file)
            writer.writerow(['Word', 'Count'])
            writer.writerows(word_counts.items())
        for word in word_counts:
            total_words += word_counts[word]
    except Exception as e:
        print(e)
    return total_words
