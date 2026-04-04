# Create a dictionary in which keys are random letters and values are lists of random integers. The dictionary is then sorted by the mean of the values in descending order, demonstrating the use of the statistics library.
# The function should output with:
#     dict: The sorted dictionary with letters as keys and lists of integers as values, sorted by their mean values.
# You should write self-contained code starting with:

import random
import statistics
def task_func(LETTERS):

    random_dict = {k: [random.randint(0, 100) for _ in range(random.randint(1, 10))] for k in LETTERS}
    sorted_dict = dict(sorted(random_dict.items(), key=lambda item: statistics.mean(item[1]), reverse=True))
    return sorted_dict