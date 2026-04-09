import re
import string
import random

def task_func(text: str, seed=None) -> str:
    if seed is not None:
        random.seed(seed)
    text = re.sub('[%s]' % re.escape(string.punctuation), '', text)
    REPLACEMENTS = {' ': '_', '\t': '__', '\n': '___'}
    for k, v in REPLACEMENTS.items():
        text = text.replace(k, v)
    text = ''.join((random.choice([k.upper(), k]) for k in text))
    return text
