import os
import re
import shutil

def task_func(source_directory: str, target_directory: str, pattern: str='\\d{4}') -> int:
    moved_files_count = 0
    if not os.path.exists(source_directory) or not os.path.isdir(source_directory):
        return 0
    if not os.path.exists(target_directory):
        os.makedirs(target_directory)
    for root, _, files in os.walk(source_directory):
        for file in files:
            if re.search(pattern, file):
                shutil.move(os.path.join(root, file), os.path.join(target_directory, file))
                moved_files_count += 1
    return moved_files_count
