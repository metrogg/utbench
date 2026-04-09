import os
import shutil
BACKUP_DIR = '/tmp/backup'

def task_func(directory):
    if not os.path.exists(BACKUP_DIR):
        return f'Backup directory {BACKUP_DIR} does not exist. Cannot rollback update.'
    backups = sorted(os.listdir(BACKUP_DIR))
    latest_backup = backups[-1] if backups else None
    if not latest_backup:
        return f'No backups found in {BACKUP_DIR}. Cannot rollback update.'
    if os.path.exists(directory):
        shutil.rmtree(directory)
    shutil.copytree(os.path.join(BACKUP_DIR, latest_backup), directory)
    return directory
