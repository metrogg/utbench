import os
import re
import shutil
SOURCE_DIR = '/source/dir'
TARGET_DIR = '/target/dir'
FILE_PATTERN = re.compile('^(.*?)-\\d+\\.json$')

def task_func():
    SOURCE_DIR = '/source/dir'
    TARGET_DIR = '/target/dir'
    FILE_PATTERN = re.compile('^(.*?)-\\d+\\.json$')
    for filename in os.listdir(SOURCE_DIR):
        match = FILE_PATTERN.match(filename)
        if match is not None:
            prefix = match.group(1)
            new_filename = f'{prefix}.json'
            shutil.move(os.path.join(SOURCE_DIR, filename), os.path.join(TARGET_DIR, new_filename))
