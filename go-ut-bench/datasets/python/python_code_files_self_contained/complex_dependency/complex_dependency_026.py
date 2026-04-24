import os
from pathlib import Path
import glob
import shutil

def task_func(source_directory: str, target_directory: str):
    moved_files = 0
    if not os.path.exists(source_directory):
        raise FileNotFoundError('source_directory must exist.')
    if not os.path.exists(target_directory):
        os.makedirs(target_directory)
    for extension in ['.txt', '.docx', '.xlsx', '.csv']:
        filepaths = glob.glob(os.path.join(source_directory, '**', '*' + extension), recursive=True)
        for filepath in filepaths:
            filename = Path(filepath).name
            stem = Path(filepath).stem
            target_filepath = os.path.join(target_directory, filename)
            count = 1
            while os.path.exists(target_filepath):
                new_filename = f'{stem}-{count}{extension}'
                target_filepath = os.path.join(target_directory, new_filename)
                count += 1
            shutil.move(filepath, target_filepath)
            moved_files += 1
    return moved_files
