import os
import shutil
import string
INVALID_CHARACTERS = string.punctuation + string.whitespace

def task_func(directory_path):
    summary = {}
    for filename in os.listdir(directory_path):
        if any((char in INVALID_CHARACTERS for char in filename)):
            if not os.path.exists(os.path.join(directory_path, 'Invalid')):
                os.mkdir(os.path.join(directory_path, 'Invalid'))
            shutil.move(os.path.join(directory_path, filename), os.path.join(directory_path, 'Invalid'))
            summary['Invalid'] = summary.get('Invalid', 0) + 1
        else:
            extension = os.path.splitext(filename)[-1].strip('.')
            if not os.path.exists(os.path.join(directory_path, extension)):
                os.mkdir(os.path.join(directory_path, extension))
            shutil.move(os.path.join(directory_path, filename), os.path.join(directory_path, extension))
            summary[extension] = summary.get(extension, 0) + 1
    return summary
