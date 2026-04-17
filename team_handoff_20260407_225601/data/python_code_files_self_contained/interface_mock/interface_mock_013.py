import os
import shutil
import fnmatch

def task_func(source_directory, destination_directory, file_pattern):
    moved_files = []
    for path, dirs, files in os.walk(source_directory):
        for filename in fnmatch.filter(files, file_pattern):
            shutil.move(os.path.join(path, filename), os.path.join(destination_directory, filename))
            moved_files.append(filename)
    return moved_files
