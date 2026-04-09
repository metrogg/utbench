import re
import os
import shutil

def task_func(directory):
    pattern = re.compile('(like|what)', re.IGNORECASE)
    interesting_files = [file for file in os.listdir(directory) if pattern.search(file)]
    if not os.path.exists(os.path.join(directory, 'Interesting Files')):
        os.mkdir(os.path.join(directory, 'Interesting Files'))
    for file in interesting_files:
        shutil.move(os.path.join(directory, file), os.path.join(directory, 'Interesting Files'))
    return interesting_files
