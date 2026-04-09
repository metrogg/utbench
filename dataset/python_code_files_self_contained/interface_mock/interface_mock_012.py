import subprocess
import os
import json
from datetime import datetime

def task_func(script_name='backup.sh', log_file='/home/user/backup_log.json'):
    log_data = {}
    if not os.path.isfile(script_name):
        raise FileNotFoundError(f'Script {script_name} does not exist.')
    start_time = datetime.now()
    log_data['start_time'] = start_time.strftime('%Y-%m-%d %H:%M:%S')
    try:
        exit_status = subprocess.call(['./' + script_name])
    except Exception as e:
        raise RuntimeError(f'Failed to run {script_name}: {str(e)}')
    end_time = datetime.now()
    log_data['end_time'] = end_time.strftime('%Y-%m-%d %H:%M:%S')
    log_data['exit_status'] = exit_status
    with open(log_file, 'w') as f:
        json.dump(log_data, f)
    return log_data
