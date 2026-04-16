import os
import ctypes
import sys
import subprocess

def task_func(filepath):
    if not isinstance(filepath, str):
        raise TypeError('Invalid filepath type')
    elif filepath == '' or not os.path.exists(filepath):
        raise OSError('Invalid filepath')
    else:
        lib = ctypes.CDLL(filepath)
    uname = os.uname()
    print(f'System: {uname.sysname}')
    print(f'Node Name: {uname.nodename}')
    print(f'Release: {uname.release}')
    print(f'Version: {uname.version}')
    print(f'Machine: {uname.machine}')
    python_version = sys.version
    print(f'Python Version: {python_version}')
    pip_version = subprocess.check_output(['pip', '--version'])
    print(f"PIP Version: {pip_version.decode('utf-8')}")
    return lib._name
