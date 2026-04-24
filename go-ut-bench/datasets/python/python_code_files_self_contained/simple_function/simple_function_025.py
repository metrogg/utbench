from collections import Counter
from operator import itemgetter
import itertools
ANIMAL = ['cat', 'camel', 'cow', 'dog', 'elephant', 'fox', 'giraffe', 'hippo', 'iguana', 'jaguar']

def task_func(animal_dict):
    animal_dict_copy = {}
    for i in animal_dict:
        if i in ANIMAL:
            animal_dict_copy[i] = animal_dict[i]
    letters = list(itertools.chain.from_iterable(animal_dict_copy.keys()))
    count_dict = dict(Counter(letters))
    sorted_dict = dict(sorted(count_dict.items(), key=itemgetter(1), reverse=True))
    return sorted_dict
