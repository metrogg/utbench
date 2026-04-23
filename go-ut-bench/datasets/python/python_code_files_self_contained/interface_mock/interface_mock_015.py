import os
from datetime import datetime
LOG_DIR = './logs'

def task_func(metrics, filename, log_dir=LOG_DIR):
    if not isinstance(metrics, dict):
        raise ValueError('Metrics must be a dictionary')
    if not isinstance(filename, str):
        raise ValueError('Filename must be a string')
    try:
        with open(os.path.join(log_dir, filename), 'a') as f:
            f.write(f'{datetime.now()}\n')
            for key, value in metrics.items():
                f.write(f'{key}: {value}\n')
            f.write('\n')
        return True
    except Exception as e:
        print(f'An error occurred: {e}')
        return False
