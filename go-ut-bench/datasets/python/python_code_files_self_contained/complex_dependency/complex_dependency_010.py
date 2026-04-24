import csv
import random

def task_func(file_path, num_rows, gender=['Male', 'Female', 'Non-Binary'], countries=['USA', 'UK', 'Canada', 'Australia', 'India'], seed=None):
    FIELDS = ['Name', 'Age', 'Gender', 'Country']
    random.seed(seed)
    with open(file_path, 'w', newline='') as csv_file:
        writer = csv.DictWriter(csv_file, fieldnames=FIELDS)
        writer.writeheader()
        for _ in range(num_rows):
            writer.writerow({'Name': ''.join(random.choices('ABCDEFGHIJKLMNOPQRSTUVWXYZ', k=5)), 'Age': random.randint(20, 60), 'Gender': random.choice(gender), 'Country': random.choice(countries)})
    return file_path
