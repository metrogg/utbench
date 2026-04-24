import json
import base64
import unicodedata

def task_func(json_file: str) -> dict:
    ENCODING = 'utf-8'
    with open(json_file, 'r') as f:
        data = json.load(f)
    decoded_data = {k: unicodedata.normalize('NFC', base64.b64decode(v).decode(ENCODING)) for k, v in data.items()}
    return decoded_data
