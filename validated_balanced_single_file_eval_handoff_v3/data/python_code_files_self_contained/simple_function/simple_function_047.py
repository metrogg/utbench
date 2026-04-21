import random
from collections import Counter
CARDS = ['2', '3', '4', '5', '6', '7', '8', '9', '10', 'J', 'Q', 'K', 'A']

def task_func(x=1):
    result = []
    card_counts = Counter()
    for i in range(x):
        drawn = random.sample(CARDS, 5)
        result.append(drawn)
        card_counts.update(drawn)
    return (result, card_counts)
