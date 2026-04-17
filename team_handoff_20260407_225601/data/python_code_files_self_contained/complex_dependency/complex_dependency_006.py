from datetime import datetime
import os
from pathlib import Path
DATE_FORMAT = '%Y-%m-%d %H:%M:%S'

def task_func(file_path):
    if not Path(file_path).exists():
        raise FileNotFoundError(f"No such file or directory: '{file_path}'")
    creation_time = os.path.getctime(file_path)
    formatted_time = datetime.fromtimestamp(creation_time).strftime(DATE_FORMAT)
    return formatted_time
