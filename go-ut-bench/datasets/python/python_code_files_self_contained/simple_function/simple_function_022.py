import random
import re

def task_func(text, seed=None):
    if seed is not None:
        random.seed(seed)

    def scramble_word(match):
        word = match.group(0)
        if len(word) > 3:
            middle = list(word[1:-1])
            random.shuffle(middle)
            return word[0] + ''.join(middle) + word[-1]
        else:
            return word
    pattern = '\\b\\w+\\b'
    scrambled_text = re.sub(pattern, scramble_word, text)
    return scrambled_text
