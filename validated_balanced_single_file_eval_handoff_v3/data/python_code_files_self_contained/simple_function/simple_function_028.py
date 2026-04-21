from collections import Counter
import random
HAND_RANKS = ['2', '3', '4', '5', '6', '7', '8', '9', '10', 'J', 'Q', 'K', 'A']
SUITS = ['H', 'D', 'C', 'S']

def task_func():
    hand = []
    for _ in range(5):
        rank = random.choice(HAND_RANKS)
        suit = random.choice(SUITS)
        card = f'{rank}{suit}'
        hand.append(card)
    rank_counts = Counter([card[:-1] for card in hand])
    return (hand, rank_counts)
