import subprocess
import time
import json
import platform
LOGFILE_PATH = 'logfile.log'

def task_func(interval, duration):
    if interval <= 0 or duration <= 0:
        raise ValueError('Interval and duration must be greater than zero.')
    start_time = time.time()
    try:
        with open(LOGFILE_PATH, 'w', encoding='utf-8') as logfile:
            while time.time() - start_time <= duration:
                operation_start_time = time.time()
                if platform.system() == 'Windows':
                    command = ['typeperf', '\\Processor(_Total)\\% Processor Time', '-sc', '1']
                else:
                    command = ['top', '-b', '-n1']
                output = subprocess.check_output(command)
                cpu_usage_line = output.decode('utf-8').split('\n')[2] if platform.system() == 'Windows' else output.decode('utf-8').split('\n')[2]
                cpu_usage = cpu_usage_line.split(',')[-1].strip().replace('"', '') if platform.system() == 'Windows' else cpu_usage_line.split(':')[1].split(',')[0].strip()
                log_data = {'timestamp': time.time(), 'cpu_usage': cpu_usage}
                json.dump(log_data, logfile)
                logfile.write('\n')
                sleep_time = max(0, interval - (time.time() - operation_start_time))
                time.sleep(sleep_time)
    except IOError as e:
        print(f'Error writing to file {LOGFILE_PATH}: {e}')
        return None
    return LOGFILE_PATH
