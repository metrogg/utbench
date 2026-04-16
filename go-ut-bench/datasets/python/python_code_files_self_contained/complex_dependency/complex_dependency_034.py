import os
import shutil

def task_func(src_folder, backup_dir):
    if not os.path.isdir(src_folder):
        raise ValueError(f"Source folder '{src_folder}' does not exist.")
    backup_folder = os.path.join(backup_dir, os.path.basename(src_folder))
    shutil.copytree(src_folder, backup_folder)
    try:
        shutil.rmtree(src_folder)
        return True
    except Exception as e:
        print(f'Error while deleting source folder: {e}')
        return False
