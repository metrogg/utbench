import re
import os
import shutil

def task_func(directory):
    for filename in os.listdir(directory):
        match = re.search('\\.(.*?)$', filename)
        if match:
            ext_dir = os.path.join(directory, match.group(1))
            if not os.path.exists(ext_dir):
                os.mkdir(ext_dir)
            shutil.move(os.path.join(directory, filename), ext_dir)
