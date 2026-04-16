import hashlib
import random
import struct
KEYS = ['470FC614', '4A0FC614', '4B9FC614', '4C8FC614', '4D7FC614']

def task_func(hex_keys=KEYS, seed=42):
    random.seed(seed)
    hex_key = random.choice(hex_keys)
    try:
        float_num = struct.unpack('!f', bytes.fromhex(hex_key))[0]
    except ValueError as e:
        raise ValueError('Invalid hexadecimal string in hex_keys.') from e
    hashed_float = hashlib.md5(str(float_num).encode()).hexdigest()
    return hashed_float
