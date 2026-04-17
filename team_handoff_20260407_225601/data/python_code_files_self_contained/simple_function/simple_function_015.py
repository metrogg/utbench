import pathlib
import os

def task_func(path: str, delimiter: str=os.path.sep) -> list:
    if not path:
        return []
    path = path.replace('\\', '/')
    path_obj = pathlib.Path(path)
    invalid_chars = set('<>:"|?*')
    if any((set(str(component)).intersection(invalid_chars) for component in path_obj.parts)):
        return []
    return [component for component in path_obj.parts if component and component != delimiter]
