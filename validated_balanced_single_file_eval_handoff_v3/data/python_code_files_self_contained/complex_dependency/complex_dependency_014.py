import sqlite3
from random import choice, seed
import os

def task_func(db_name, table_name, num_entries, random_seed=None):
    NAMES = ['John', 'Jane', 'Steve', 'Emma', 'Liam', 'Olivia']
    AGES = range(18, 65)
    HEIGHTS = range(150, 200)
    if random_seed:
        seed(random_seed)
    if num_entries < 0:
        raise ValueError('num_entries must not be negative')
    conn = sqlite3.connect(db_name)
    cur = conn.cursor()
    cur.execute(f'CREATE TABLE {table_name} (name TEXT, age INTEGER, height INTEGER)')
    for _ in range(num_entries):
        name = choice(NAMES)
        age = choice(AGES)
        height = choice(HEIGHTS)
        cur.execute(f'INSERT INTO {table_name} VALUES (?, ?, ?)', (name, age, height))
    conn.commit()
    return os.path.abspath(db_name)
