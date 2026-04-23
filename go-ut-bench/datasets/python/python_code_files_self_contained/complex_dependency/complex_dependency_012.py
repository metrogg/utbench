import sqlite3
import random

def task_func(db_path, num_entries, users=['Alice', 'Bob', 'Charlie', 'Dave', 'Eve'], countries=['USA', 'UK', 'Canada', 'Australia', 'India'], random_seed=None):
    random.seed(random_seed)
    conn = sqlite3.connect(db_path)
    c = conn.cursor()
    c.execute('\n        CREATE TABLE users\n        (id INTEGER PRIMARY KEY, name TEXT, age INTEGER, country TEXT)\n    ')
    for _ in range(num_entries):
        user = random.choice(users)
        age = random.randint(20, 60)
        country = random.choice(countries)
        c.execute('INSERT INTO users (name, age, country) VALUES (?, ?, ?)', (user, age, country))
    conn.commit()
    conn.close()
    return db_path
