import hashlib
import io
import os

def task_func(file_path1, file_path2):
    if not os.path.exists(file_path1) or not os.path.exists(file_path2):
        raise FileNotFoundError('File not found! Please specify a valid filepath')
    with io.open(file_path1, 'rb') as file1, io.open(file_path2, 'rb') as file2:
        file1_hash = hashlib.md5(file1.read()).hexdigest()
        file2_hash = hashlib.md5(file2.read()).hexdigest()
    return file1_hash == file2_hash
