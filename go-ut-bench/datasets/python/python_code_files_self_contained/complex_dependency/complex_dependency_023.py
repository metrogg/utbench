import csv
import random
DATA = ['Temperature', 'Humidity', 'Pressure']
RANGE = {'Temperature': (-50, 50), 'Humidity': (0, 100), 'Pressure': (980, 1040)}

def task_func(file_name='data.csv'):
    with open(file_name, 'w', newline='') as file:
        writer = csv.writer(file)
        writer.writerow(['Time'] + DATA)
        for hour in range(24):
            row = [f'{hour}:00']
            for data_type in DATA:
                min_val, max_val = RANGE[data_type]
                row.append(random.uniform(min_val, max_val))
            writer.writerow(row)
    return file_name
