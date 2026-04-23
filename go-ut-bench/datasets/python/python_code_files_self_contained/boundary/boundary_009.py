import os
import json

def task_func(config_path: str) -> dict:
    if not os.path.isfile(config_path):
        raise FileNotFoundError(f'The configuration file {config_path} does not exist.')
    with open(config_path) as f:
        config = json.load(f)
    return config
