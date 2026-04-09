import re
import os
import glob
import mimetypes

def task_func(directory, file_pattern, suffix):
    os.chdir(directory)
    files = glob.glob(file_pattern)
    file_types = {}
    for file in files:
        if re.search(suffix, file):
            file_type = mimetypes.guess_type(file)[0]
            file_types[file] = file_type
    return file_types
