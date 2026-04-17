import re
import collections

def task_func(string, patterns=['nnn', 'aaa', 'sss', 'ddd', 'fff']):
    if not isinstance(string, str):
        raise TypeError('Input string should be of type string.')
    if not isinstance(patterns, list):
        raise TypeError('patterns should be a list of strings.')
    if not all((isinstance(s, str) for s in patterns)):
        raise TypeError('patterns should be a list of strings.')
    pattern_counts = collections.defaultdict(int)
    for pattern in patterns:
        pattern_counts[pattern] = len(re.findall(pattern, string))
    return dict(pattern_counts)
