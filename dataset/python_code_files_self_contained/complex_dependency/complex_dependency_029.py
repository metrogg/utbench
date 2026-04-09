import subprocess
import os
import shutil
from glob import glob

def task_func(src_folder, dst_folder):
    if not os.path.isdir(src_folder):
        raise ValueError(f"Source folder '{src_folder}' does not exist.")
    if not os.path.isdir(dst_folder):
        raise ValueError(f"Destination folder '{dst_folder}' does not exist.")
    processes = []
    failed_files = []
    for file in glob(os.path.join(src_folder, '*')):
        process = subprocess.Popen(['gzip', file])
        processes.append((process, file))
    for process, file in processes:
        retcode = process.wait()
        if retcode != 0:
            failed_files.append(os.path.basename(file))
    for file in glob(os.path.join(src_folder, '*.gz')):
        try:
            shutil.move(file, dst_folder)
        except Exception as e:
            failed_files.append(os.path.basename(file))
    if failed_files:
        return {'success': False, 'message': 'Some files failed to compress or move.', 'failed_files': failed_files}
    else:
        return {'success': True, 'message': 'All files compressed and moved successfully.', 'failed_files': []}
