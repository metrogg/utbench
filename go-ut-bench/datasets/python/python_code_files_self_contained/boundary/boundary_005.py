import json
import random
import hashlib
from datetime import datetime

def task_func(utc_datetime, salt='salt', password_length=10, seed=0):
    random.seed(seed)
    if not isinstance(utc_datetime, datetime):
        raise ValueError('Input should be a datetime object')
    if not isinstance(salt, str):
        raise ValueError('Salt should be a string')
    utc_time_str = utc_datetime.strftime('%Y-%m-%d %H:%M:%S')
    salted_string = utc_time_str + salt
    password = ''.join((random.choice('abcdefghijklmnopqrstuvwxyz0123456789') for _ in range(password_length)))
    hashed_password = hashlib.sha256((password + salted_string).encode('utf-8')).hexdigest()
    password_json_str = json.dumps(hashed_password)
    return password_json_str
