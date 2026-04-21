from collections import Counter
from random import choice, seed
POSSIBLE_ITEMS = ['apple', 'banana', 'cherry', 'date', 'elderberry']

def task_func(list_of_lists):
    seed(42)
    baskets = []
    for list_ in list_of_lists:
        basket = Counter()
        for _ in list_:
            basket[choice(POSSIBLE_ITEMS)] += 1
        baskets.append(basket)
    return baskets
