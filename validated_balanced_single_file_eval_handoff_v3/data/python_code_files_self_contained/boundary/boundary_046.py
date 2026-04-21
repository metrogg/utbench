import sys
from configparser import ConfigParser
PATH_TO_APPEND = '/path/to/whatever'
CONFIG_FILE = '/path/to/config.ini'

def task_func(path_to_append=PATH_TO_APPEND, config_file=CONFIG_FILE):
    if isinstance(path_to_append, list):
        for path in path_to_append:
            sys.path.append(path)
    else:
        sys.path.append(path_to_append)
    config = ConfigParser()
    if not os.path.exists(config_file):
        open(config_file, 'a').close()
    config.read(config_file)
    path_str = ','.join(path_to_append) if isinstance(path_to_append, list) else path_to_append
    config.set('DEFAULT', 'path_to_append', path_str)
    with open(config_file, 'w') as file:
        config.write(file)
    return (config, config_file)
