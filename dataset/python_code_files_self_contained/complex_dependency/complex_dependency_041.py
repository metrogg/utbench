import os
import re

def task_func(directory):
    BRACKET_PATTERN = '[(){}\\[\\]]'
    file_list = []
    for root, dirs, files in os.walk(directory):
        for file in files:
            if re.search(BRACKET_PATTERN, file):
                file_list.append(os.path.join(root, file))
    return file_list
