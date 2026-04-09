import shutil
from pathlib import Path
from typing import List

def task_func(source_dir: str, target_dir: str, extensions: List[str]) -> int:
    if Path(source_dir).is_dir() == False:
        raise ValueError('source_dir does not exist.')
    if Path(target_dir).is_dir() == False:
        raise ValueError('target_dir does not exist.')
    count = 0
    for extension in extensions:
        for file_name in Path(source_dir).glob(f'*{extension}'):
            shutil.move(str(file_name), target_dir)
            count += 1
    return count
