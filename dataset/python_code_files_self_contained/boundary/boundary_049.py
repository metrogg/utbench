import os
import hashlib
import base64

def task_func(password, PREFIX='ME', SALT_LENGTH=16):
    if SALT_LENGTH < 0:
        raise ValueError
    salt = os.urandom(SALT_LENGTH)
    salted_password = PREFIX + password + salt.hex()
    hashed_password = hashlib.sha256(salted_password.encode()).digest()
    return base64.b64encode(hashed_password).decode()
