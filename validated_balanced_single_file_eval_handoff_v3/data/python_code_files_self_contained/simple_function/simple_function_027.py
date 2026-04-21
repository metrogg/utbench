import random
import string
import hashlib
import time

def task_func(data_dict: dict, seed=0) -> dict:
    random.seed(seed)
    SALT_LENGTH = 5
    data_dict.update(dict(a=1))
    salt = ''.join((random.choice(string.ascii_lowercase) for _ in range(SALT_LENGTH)))
    for key in data_dict.keys():
        data_dict[key] = hashlib.sha256((str(data_dict[key]) + salt).encode()).hexdigest()
    data_dict['timestamp'] = time.time()
    return data_dict
