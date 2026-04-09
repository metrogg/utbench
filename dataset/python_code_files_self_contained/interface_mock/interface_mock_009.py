import subprocess
import os
import shutil
import sys
DIRECTORY = 'c:\\Program Files\\VMware\\VMware Server'
BACKUP_DIRECTORY = 'c:\\Program Files\\VMware\\VMware Server\\Backup'

def task_func(filename):
    file_path = os.path.join(DIRECTORY, filename)
    backup_path = os.path.join(BACKUP_DIRECTORY, filename)
    try:
        shutil.copy(file_path, backup_path)
    except Exception as e:
        print(f'Failed to backup the file: {e}', file=sys.stderr)
        return -1
    try:
        process = subprocess.Popen(file_path)
        return process.poll()
    except Exception as e:
        print(f'Failed to execute the file: {e}', file=sys.stderr)
        return -1
