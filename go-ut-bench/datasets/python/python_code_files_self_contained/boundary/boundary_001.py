import os
import re
import hashlib

def task_func(path, delimiter):
    path_components = re.split(f'({delimiter})', path)
    hashes = []
    for component in path_components:
        if not component:
            continue
        if component != delimiter and os.path.isfile(component):
            with open(component, 'rb') as f:
                hashes.append(hashlib.sha256(f.read()).hexdigest())
        else:
            hashes.append(None)
    return list(zip(path_components, hashes))
