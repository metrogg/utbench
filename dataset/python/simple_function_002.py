import random
import string
LETTERS = string.ascii_letters
DIGITS = string.digits

def task_func(length, num_digits):
    random.seed(42)
    if length <= 0:
        raise ValueError('Length must be a positive integer.')
    if not 0 <= num_digits <= length:
        raise ValueError('num_digits must be a non-negative integer and less than or equal to length.')
    password = []
    for _ in range(length - num_digits):
        password.append(random.choice(LETTERS))
    for _ in range(num_digits):
        password.append(random.choice(DIGITS))
    random.shuffle(password)
    return ''.join(password)
