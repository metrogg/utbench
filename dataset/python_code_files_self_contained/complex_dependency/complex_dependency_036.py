import os
import shutil

def task_func(src_dir, dest_dir, extension):
    files_moved = 0
    for file_name in os.listdir(src_dir):
        if file_name.endswith(extension):
            shutil.move(os.path.join(src_dir, file_name), os.path.join(dest_dir, file_name))
            files_moved += 1
    return files_moved
