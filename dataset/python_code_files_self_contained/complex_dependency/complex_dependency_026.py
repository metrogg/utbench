import binascii
import hashlib
import re
OUTPUT_DIR = './output'

def task_func(directory: str, pattern: str='(?<!Distillr)\\\\AcroTray\\.exe') -> dict:
    hashes = {}
    for root, _, files in os.walk(directory):
        for file in files:
            if re.search(pattern, file):
                path = os.path.join(root, file)
                with open(path, 'rb') as f:
                    data = f.read()
                    hash_digest = hashlib.sha256(data).digest()
                    hashes[path] = binascii.hexlify(hash_digest).decode()
    return hashes
