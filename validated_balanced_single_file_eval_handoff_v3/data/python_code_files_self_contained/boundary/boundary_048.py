import os
import hashlib
DIRECTORY = './hashed_files'

def task_func(input_string):
    if not os.path.exists(DIRECTORY):
        os.makedirs(DIRECTORY)
    file_paths = []
    lines = input_string.split('\n')
    for line in lines:
        if line:
            line_hash = hashlib.sha256(line.encode()).hexdigest()
            filename = line_hash[:10] + '.txt'
            filepath = os.path.join(DIRECTORY, filename)
            with open(filepath, 'w', encoding='utf-8') as file:
                file.write(line_hash)
            file_paths.append(filepath)
    return file_paths
