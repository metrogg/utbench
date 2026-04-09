import os
import glob
from collections import Counter

def task_func(directory, extensions=['.txt', '.docx', '.xlsx', '.csv'], keep_zero=True):
    if not os.path.exists(directory):
        raise OSError('directory must exist.')
    counter = Counter()
    for suffix in extensions:
        count = len(glob.glob(os.path.join(directory, '**', '*' + suffix), recursive=True))
        if count:
            counter[suffix] += count
        elif keep_zero:
            counter[suffix] += count
    return counter
