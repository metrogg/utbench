import csv
import os
OUTPUT_DIR = './output'

def task_func(df, filename, output_dir=OUTPUT_DIR):
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)
    file_path = os.path.join(output_dir, filename)
    df.to_csv(file_path, index=False, quoting=csv.QUOTE_NONNUMERIC)
    return os.path.abspath(file_path)
