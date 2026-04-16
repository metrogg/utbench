import re
import os
from collections import Counter

def task_func(folder_path: str) -> dict:
    IP_REGEX = re.compile('\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}')
    counter = Counter()
    for filename in os.listdir(folder_path):
        if filename.endswith('.log'):
            with open(os.path.join(folder_path, filename)) as file:
                content = file.read()
                ips = re.findall(IP_REGEX, content)
                counter.update(ips)
    return dict(counter)
