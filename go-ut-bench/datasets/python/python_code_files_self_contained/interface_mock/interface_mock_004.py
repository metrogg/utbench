import os
import shutil

def task_func(path, delimiter='/'):
    if not path or not isinstance(path, str):
        raise ValueError('Path must be a non-empty string')
    if not os.path.exists(path):
        raise FileNotFoundError(f"Path '{path}' does not exist")
    path_components = path.strip(delimiter).split(delimiter)
    if not all(path_components):
        raise ValueError('Path contains invalid components')
    results = []
    for index, component in enumerate(path_components):
        sub_path = delimiter.join(path_components[:index + 1])
        if not sub_path.startswith(delimiter):
            sub_path = delimiter + sub_path
        usage = shutil.disk_usage(sub_path)
        results.append((component, {'total': usage.total, 'used': usage.used, 'free': usage.free}))
    return results
