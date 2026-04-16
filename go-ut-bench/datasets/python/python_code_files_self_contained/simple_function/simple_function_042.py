import binascii
import hashlib

def task_func(input_string, verify_hash=None):
    if not isinstance(input_string, str):
        raise TypeError('Input must be a string')
    if verify_hash is not None and (not isinstance(verify_hash, str)):
        raise TypeError('verify_hash must be a string or None')
    hashed_bytes = hashlib.sha256(input_string.encode()).digest()
    hex_encoded_hash = binascii.hexlify(hashed_bytes).decode()
    if verify_hash is not None:
        return hex_encoded_hash == verify_hash
    return hex_encoded_hash
