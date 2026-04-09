import os
import time
OUTPUT_DIR = './output'

def task_func(dataset, filename, output_dir=OUTPUT_DIR):
    start_time = time.time()
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)
    filepath = os.path.join(output_dir, filename)
    with open(filepath, 'w', newline='') as f:
        for i, df in enumerate(dataset):
            if i > 0:
                f.write('------\n')
            df.to_csv(f, index=False, header=True, mode='a')
            if i < len(dataset) - 1:
                f.write('\n')
    end_time = time.time()
    cost = f'Operation completed in {end_time - start_time} seconds.'
