import zipfile
import os
import re
import shutil

def task_func(source_dir: str, target_dir: str, archive_name: str='archive.zip') -> str:
    os.makedirs(source_dir, exist_ok=True)
    os.makedirs(target_dir, exist_ok=True)
    archive_path = os.path.join(target_dir, archive_name)
    with zipfile.ZipFile(archive_path, 'w') as archive:
        for file in os.listdir(source_dir):
            if re.search('_processed$', os.path.splitext(file)[0]):
                archive.write(os.path.join(source_dir, file), arcname=file)
                shutil.move(os.path.join(source_dir, file), target_dir)
    return archive_path
