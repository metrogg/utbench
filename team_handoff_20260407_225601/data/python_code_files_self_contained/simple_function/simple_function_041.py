import random
import string

def task_func(length: int, predicates: list, seed: int=None):
    if seed is not None:
        random.seed(seed)
    if length < 0:
        raise ValueError('Length must be non-negative.')
    predicate_functions = {'has_uppercase': lambda x: any((c.isupper() for c in x)), 'has_lowercase': lambda x: any((c.islower() for c in x)), 'has_special_chars': lambda x: any((c in string.punctuation for c in x)), 'has_numbers': lambda x: any((c.isdigit() for c in x))}
    predicates = list(set(predicates))
    if any((p not in predicate_functions for p in predicates)):
        raise KeyError(f'Invalid predicate provided.')
    characters = string.ascii_letters + string.digits + string.punctuation
    generated_string = ''.join(random.choices(characters, k=length))
    results = {predicate: predicate_functions[predicate](generated_string) for predicate in predicates}
    return (generated_string, results)
