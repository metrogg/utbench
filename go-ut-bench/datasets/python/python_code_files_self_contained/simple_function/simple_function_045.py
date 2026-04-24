import pickle
import os

def task_func(filename, data):
    try:
        directory = os.path.dirname(filename)
        if directory and (not os.path.exists(directory)):
            os.makedirs(directory)
        with open(filename, 'wb') as f:
            pickle.dump(data, f)
        return True
    except Exception as e:
        print(f'An error occurred: {e}')
        return False
