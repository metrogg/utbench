import os
import sys
import importlib
from pkgutil import iter_modules

def task_func(package_name):
    added_modules = []
    try:
        package = importlib.import_module(package_name)
    except ImportError:
        raise ImportError(f"The package '{package_name}' is not installed! Please install the package first using 'pip install {package_name}'")
    for _, module_name, _ in iter_modules(package.__path__):
        module_path = os.path.join(package.__path__[0], module_name)
        if module_path not in sys.path:
            sys.path.append(module_path)
            added_modules.append(module_name)
    return added_modules
