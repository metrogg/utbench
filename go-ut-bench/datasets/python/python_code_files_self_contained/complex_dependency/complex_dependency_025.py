import csv
import os
from datetime import datetime
from random import randint
SENSORS = ['Temperature', 'Humidity', 'Pressure']
OUTPUT_DIR = './output'

def task_func(hours, output_dir=OUTPUT_DIR):
    FILE_PATH = os.path.join(output_dir, 'sensor_data.csv')
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)
    data = [['Time'] + SENSORS]
    for i in range(hours):
        row = [datetime.now().strftime('%Y-%m-%d %H:%M:%S.%f')] + [randint(0, 100) for _ in SENSORS]
        data.append(row)
    with open(FILE_PATH, 'w', newline='') as f:
        writer = csv.writer(f)
        writer.writerows(data)
    return FILE_PATH
