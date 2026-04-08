from collections import Counter
import re

def task_func(word: str) -> list:
    clean_word = re.sub('[^a-z]', '', word.lower())
    if len(clean_word) < 2:
        return []
    pairs = [clean_word[i:i + 2] for i in range(len(clean_word) - 1)]
    pair_counter = Counter(pairs)
    most_common = pair_counter.most_common(1)
    return most_common
