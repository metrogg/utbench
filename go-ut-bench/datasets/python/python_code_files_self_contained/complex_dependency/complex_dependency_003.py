import csv
import random

def task_func(csv_file='names.csv', latin_names=['Sopetón', 'Méndez', 'Gómez', 'Pérez', 'Muñoz'], names=['Smith', 'Johnson', 'Williams', 'Brown', 'Jones'], encoding='latin-1', rng_seed=None):
    if not isinstance(csv_file, str):
        raise TypeError('csv_file should be a string.')
    if not isinstance(names, list):
        raise TypeError('names should be a list.')
    if not isinstance(latin_names, list):
        raise TypeError('latin_names should be a list.')
    if rng_seed is not None:
        random.seed(rng_seed)
    with open(csv_file, 'w', newline='', encoding=encoding) as csvfile:
        fieldnames = ['Name', 'Age']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        for _ in range(50):
            if latin_names:
                writer.writerow({'Name': random.choice(latin_names), 'Age': random.randint(20, 50)})
            if names:
                writer.writerow({'Name': random.choice(names), 'Age': random.randint(20, 50)})
    return csv_file
