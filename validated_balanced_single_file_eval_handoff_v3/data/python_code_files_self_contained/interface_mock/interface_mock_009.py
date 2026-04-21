import os
import logging

def task_func(directory):
    logging.basicConfig(filename='jquery_removal.log', level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
    if not os.path.exists(directory):
        raise FileNotFoundError(f"Directory '{directory}' does not exist.")
    files = os.listdir(directory)
    removed_files = 0
    removed_file_names = []
    for file in files:
        if 'jquery' in file and file.endswith('.js'):
            try:
                os.remove(os.path.join(directory, file))
                removed_files += 1
                removed_file_names.append(file)
                logging.info(f'Removed jQuery file: {file}')
            except Exception as e:
                logging.error(f'Error while removing file {file}: {e}')
    return (removed_files, removed_file_names)
