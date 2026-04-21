import subprocess
import os
import signal
import time

def task_func(process_name: str) -> int:
    try:
        pids = subprocess.check_output(['pgrep', '-f', process_name]).decode().split('\n')[:-1]
    except subprocess.CalledProcessError:
        pids = []
    for pid in pids:
        os.kill(int(pid), signal.SIGTERM)
    time.sleep(1)
    return len(pids)
