import subprocess
import shutil
import os

def task_func(script_path: str, temp_dir: str) -> str:
    try:
        shutil.copy(script_path, temp_dir)
        temp_script_path = os.path.join(temp_dir, os.path.basename(script_path))
        result = subprocess.call(['python', temp_script_path])
        print(result)
        if result == 0:
            return 'Script executed successfully!'
        else:
            return 'Script execution failed!'
    except Exception as e:
        return 'Script execution failed!'
