import ctypes
import os
import shutil
import glob

def task_func(filepath, destination_dir):
    lib = ctypes.CDLL(filepath)
    dll_dir = os.path.dirname(filepath)
    dll_files = glob.glob(os.path.join(dll_dir, '*.dll'))
    for dll_file in dll_files:
        shutil.move(dll_file, destination_dir)
    return lib._name
