import string
import random

def task_func(text, seed=None):
    if seed is not None:
        random.seed(seed)
    if not text:
        raise ValueError('text cannot be empty.')
    password = ''
    for char in text:
        random_lowercase = random.choice(string.ascii_lowercase)
        random_digit = random.choice(string.digits)
        if char.isalpha():
            password += random_lowercase
        elif char.isdigit():
            password += random_digit
        elif char == ' ':
            if random.random() < 0.5:
                password += random_lowercase
            else:
                password += random_digit
        else:
            password += char
    return password
