import os
import shutil
import glob

def task_func(source_dir, dest_dir, extension):
    files = glob.glob(os.path.join(source_dir, f'*.{extension}'))
    for file in files:
        shutil.move(file, dest_dir)
    result = len(files)
    return result
