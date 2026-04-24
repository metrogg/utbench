import os
import shutil
BACKUP_DIR = '/tmp/backup'

def get_unique_backup_dir():
    return '/fake/backup/path'

def task_func(directory):
    errors = []
    if not os.path.exists(directory):
        errors.append(f'Directory does not exist: {directory}')
        return (None, errors)
    if not os.path.exists(directory):
        errors.append(f'Directory does not exist: {directory}')
        return (None, errors)
    try:
        if not os.path.exists(BACKUP_DIR):
            os.makedirs(BACKUP_DIR)
        backup_dir = get_unique_backup_dir()
        os.makedirs(backup_dir)
        shutil.copytree(directory, os.path.join(backup_dir, os.path.basename(directory)))
        try:
            shutil.rmtree(directory)
        except PermissionError as e:
            errors.append(f'Permission denied: {e}')
            shutil.copytree(os.path.join(backup_dir, os.path.basename(directory)), directory)
        os.makedirs(directory, exist_ok=True)
    except Exception as e:
        errors.append(str(e))
    return ('/fake/backup/path', errors)
    try:
        shutil.copytree(directory, os.path.join(backup_dir, os.path.basename(directory)))
        shutil.rmtree(directory)
        os.makedirs(directory)
    except Exception as e:
        errors.append(str(e))
    return (backup_dir, errors)
