import collections
import operator
import os
import shutil

def task_func(data_dict, source_directory, backup_directory):
    data_dict.update({'a': 1})
    counter = collections.Counter(data_dict.values())
    sorted_dict = sorted(counter.items(), key=operator.itemgetter(1), reverse=True)
    backup_status = False
    if os.path.isdir(source_directory):
        shutil.copytree(source_directory, backup_directory, dirs_exist_ok=True)
        backup_status = True
    return (data_dict, sorted_dict, backup_status)
