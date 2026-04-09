import subprocess
import os
import threading

def task_func(script_path: str, timeout: int=60) -> str:

    def target():
        subprocess.call(['python', script_path])
    thread = threading.Thread(target=target)
    thread.start()
    thread.join(timeout)
    if thread.is_alive():
        os.system(f'pkill -f "{script_path}"')
        thread.join()
        return 'Terminating process due to timeout.'
    else:
        return 'Script executed successfully.'
