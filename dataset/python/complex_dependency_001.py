import os
import shutil
import glob
import hashlib

def task_func(ROOT_DIR, DEST_DIR, SPECIFIC_HASH):
    files_moved = 0
    os.makedirs(DEST_DIR, exist_ok=True)
    for filename in glob.glob(os.path.join(ROOT_DIR, '*')):
        if not os.path.exists(filename) or os.path.isdir(filename):
            continue
        with open(filename, 'rb') as f:
            file_hash = hashlib.md5(f.read()).hexdigest()
        if file_hash == SPECIFIC_HASH:
            shutil.move(filename, DEST_DIR)
            files_moved += 1
    return files_moved
