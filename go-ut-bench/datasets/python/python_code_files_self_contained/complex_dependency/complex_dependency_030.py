import os
import hashlib
import json
from pathlib import Path

def task_func(directory: str) -> str:
    hash_dict = {}
    for root, dirs, files in os.walk(directory):
        for file in files:
            file_path = Path(root) / file
            with open(file_path, 'rb') as f:
                bytes = f.read()
                readable_hash = hashlib.sha256(bytes).hexdigest()
                hash_dict[str(file_path)] = readable_hash
    json_file = Path(directory) / 'hashes.json'
    with open(json_file, 'w') as f:
        json.dump(hash_dict, f)
    return str(json_file)
