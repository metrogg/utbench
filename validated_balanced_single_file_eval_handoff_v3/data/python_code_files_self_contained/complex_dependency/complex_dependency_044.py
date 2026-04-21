import os
import glob
import shutil

def task_func(directory, archive_dir='archive'):
    if not os.path.exists(archive_dir):
        os.makedirs(archive_dir)
    json_files = glob.glob(os.path.join(directory, '*.json'))
    error_messages = []
    for json_file in json_files:
        try:
            shutil.move(json_file, archive_dir)
        except Exception as e:
            error_message = f'Unable to move {json_file} due to {str(e)}'
            error_messages.append(error_message)
    return (len(error_messages) == 0, error_messages)
