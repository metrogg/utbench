import os
import errno
import shutil

def task_func(filename, dest_dir):
    try:
        os.makedirs(dest_dir, exist_ok=True)
    except OSError as e:
        if e.errno != errno.EEXIST:
            raise
    dest = shutil.copy(filename, dest_dir)
    with open(filename, 'w') as original_file:
        original_file.truncate(0)
    return os.path.abspath(dest)
