import re
from pathlib import Path
import tarfile
PATTERN = '(?<!Distillr)\\\\\\\\AcroTray\\.exe'
DIRECTORY = 'C:\\\\SomeDir\\\\'

def task_func(directory=DIRECTORY, file_pattern=PATTERN):
    tar_path = Path(directory) / 'archive.tar'
    with tarfile.open(tar_path, 'w') as tar:
        for path in Path(directory).rglob('*'):
            if re.match(file_pattern, path.name):
                try:
                    tar.add(path, arcname=path.relative_to(directory))
                except PermissionError as e:
                    print(f'Skipping {path} due to permission error: {e}')
    return str(tar_path)
