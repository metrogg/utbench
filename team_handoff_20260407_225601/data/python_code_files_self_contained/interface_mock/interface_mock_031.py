import os
from datetime import datetime

def task_func(filepath: str) -> dict:
    try:
        size = os.path.getsize(filepath)
        mtime = os.path.getmtime(filepath)
        mtime = datetime.fromtimestamp(mtime).strftime('%Y-%m-%d %H:%M:%S')
    except OSError as e:
        raise Exception(f'Error: {e}')
    return {'size': f'{size} bytes', 'last_modified': mtime}
