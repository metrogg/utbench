import difflib
import gzip

def task_func(file_path1, file_path2):
    with gzip.open(file_path1, 'rt') as file1, gzip.open(file_path2, 'rt') as file2:
        file1_content = file1.readlines()
        file2_content = file2.readlines()
        diff = difflib.ndiff(file1_content, file2_content)
        diff = [line for line in diff if line.startswith('+ ') or line.startswith('- ')]
    return ''.join(diff)
