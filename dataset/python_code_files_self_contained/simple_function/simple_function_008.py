import shutil
import os
import fnmatch
import itertools

def task_func(src_dir, dst_dir):
    FILE_PATTERNS = ['*.txt', '*.docx']
    matching_files = list(itertools.chain.from_iterable((fnmatch.filter(os.listdir(src_dir), pattern) for pattern in FILE_PATTERNS)))
    for filename in matching_files:
        shutil.copy2(os.path.join(src_dir, filename), dst_dir)
    return dst_dir
