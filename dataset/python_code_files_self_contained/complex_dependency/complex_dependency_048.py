import os
import shutil

def task_func(directory, backup_directory):
    copied_files = []
    if not os.path.exists(backup_directory):
        os.makedirs(backup_directory)
    for filename in os.listdir(directory):
        if filename.endswith('.json'):
            src = os.path.join(directory, filename)
            dst = os.path.join(backup_directory, filename)
            shutil.copy(src, dst)
            copied_files.append(dst)
    return copied_files
