import json
import random
from datetime import datetime, timedelta
USERS = ['Alice', 'Bob', 'Charlie', 'Dave', 'Eve']

def task_func(file_path, num_entries, seed=None):
    if seed is not None:
        random.seed(seed)
    log_entries = []
    current_time = datetime.now()
    for _ in range(num_entries):
        user = random.choice(USERS)
        action = random.choice(['login', 'logout', 'view_page', 'edit_profile', 'post_message'])
        timestamp = current_time.strftime('%Y-%m-%dT%H:%M:%S')
        log_entries.append({'user': user, 'action': action, 'timestamp': timestamp})
        current_time -= timedelta(minutes=random.randint(1, 60))
    with open(file_path, 'w') as json_file:
        json.dump(log_entries, json_file, indent=4)
    return file_path
