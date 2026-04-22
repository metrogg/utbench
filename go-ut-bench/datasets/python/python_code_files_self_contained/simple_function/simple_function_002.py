from string import ascii_lowercase
import re
from collections import Counter
LETTERS_PATTERN = re.compile('^(.*?)-[a-z]$')
LETTERS = ascii_lowercase

def task_func(string):
    match = re.search('^(.*)-', string)
    if match:
        prefix = match.group(1)
    else:
        prefix = string if string.isalpha() else ''
    letter_counts = Counter(prefix)
    result = {letter: 0 for letter in ascii_lowercase}
    result.update({letter: letter_counts.get(letter, 0) for letter in letter_counts if letter in result})
    return result
