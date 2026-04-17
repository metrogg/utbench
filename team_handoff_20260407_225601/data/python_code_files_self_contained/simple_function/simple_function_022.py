import os
import glob

def task_func(directory_path):
    new_filenames = []
    for filename in glob.glob(os.path.join(directory_path, '*')):
        base_name = os.path.basename(filename)
        new_base_name = '.'.join(base_name.split('.')[::-1])
        os.rename(filename, os.path.join(directory_path, new_base_name))
        new_filenames.append(new_base_name)
    return new_filenames
