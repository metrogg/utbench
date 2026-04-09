import configparser
import os
import shutil

def task_func(config_file_path, archieve_dir='/home/user/archive'):
    config = configparser.ConfigParser()
    config.read(config_file_path)
    project_dir = config.get('Project', 'directory')
    if not os.path.isdir(project_dir):
        raise FileNotFoundError(f'Directory {project_dir} does not exist.')
    archive_file = f'{archieve_dir}/{os.path.basename(project_dir)}.zip'
    shutil.make_archive(base_name=os.path.splitext(archive_file)[0], format='zip', root_dir=project_dir)
    if not os.path.isfile(archive_file):
        raise Exception(f'Failed to create archive {archive_file}')
    return True
