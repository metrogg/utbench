import re
import json
from collections import defaultdict
import string

def task_func(json_string):
    try:
        data = json.loads(json_string)
        text = data.get('text', '')
    except json.JSONDecodeError:
        return {}
    text = re.sub('[^\\sa-zA-Z0-9]', '', text).lower().strip()
    text = text.translate({ord(c): None for c in string.punctuation})
    word_counts = defaultdict(int)
    for word in text.split():
        word_counts[word] += 1
    return dict(word_counts)
