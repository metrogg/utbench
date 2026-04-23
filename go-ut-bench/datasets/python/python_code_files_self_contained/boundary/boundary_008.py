import json
import re
from collections import Counter
REPLACE_NONE = 'None'

def task_func(json_str):
    data = json.loads(json_str)
    processed_data = {}
    for key, value in data.items():
        if value is None:
            continue
        if isinstance(value, str) and re.match('[^@]+@[^@]+\\.[^@]+', value):
            value = REPLACE_NONE
        processed_data[key] = value
    value_counts = Counter(processed_data.values())
    return {'data': processed_data, 'value_counts': value_counts}
