import re
import json
from collections import Counter

def task_func(json_str, top_n=10):
    pattern = '(https?:\\/\\/(?:www\\.|(?!www))[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\\.[^\\s]{2,}|www\\.[a-zA-Z0-9][a-zA-Z0-9-]+[a-zA-Z0-9]\\.[^\\s]{2,}|https?:\\/\\/(?:www\\.|(?!www))[a-zA-Z0-9]+\\.[^\\s]{2,}|www\\.[a-zA-Z0-9]+\\.[^\\s]{2,})'
    data = json.loads(json_str)
    urls = []

    def extract(dictionary):
        for key, value in dictionary.items():
            if isinstance(value, dict):
                extract(value)
            elif isinstance(value, str) and re.match(pattern, value):
                urls.append(value)
    extract(data)
    if not urls:
        return {}
    elif len(urls) <= top_n:
        return dict(Counter(urls))
    return dict(Counter(urls).most_common(top_n))
