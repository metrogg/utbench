import pickle
import os
import random
import string

def task_func(strings, filename=None):
    if filename is None:
        filename = ''.join((random.choice(string.ascii_uppercase + string.digits) for _ in range(10))) + '.pkl'
    with open(filename, 'wb') as file:
        pickle.dump(strings, file)
    with open(filename, 'rb') as file:
        loaded_strings = pickle.load(file)
    os.remove(filename)
    return loaded_strings
