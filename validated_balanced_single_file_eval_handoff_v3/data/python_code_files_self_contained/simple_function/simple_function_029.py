from collections import Counter
import itertools
import random
ALPHABET = 'abcdefghijklmnopqrstuvwxyz'

def task_func(list_of_lists, seed=0):
    random.seed(seed)
    flattened_list = list(itertools.chain(*list_of_lists))
    for list_item in list_of_lists:
        if list_item == []:
            flattened_list += random.sample(ALPHABET, 10)
    counter = Counter(flattened_list)
    return counter
