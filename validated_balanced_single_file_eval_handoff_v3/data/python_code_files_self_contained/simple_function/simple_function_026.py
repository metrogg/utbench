import string
import random
import re

def task_func(elements, pattern, seed=100):
    random.seed(seed)
    replaced_elements = []
    for element in elements:
        replaced = ''.join([random.choice(string.ascii_letters) for _ in element])
        formatted = '%{}%'.format(replaced)
        replaced_elements.append(formatted)
    concatenated_elements = ''.join(replaced_elements)
    search_result = re.search(pattern, concatenated_elements)
    return (replaced_elements, bool(search_result))
