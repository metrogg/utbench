import os
import re

def task_func(log_file_path: str, keywords: list):
    if not os.path.exists(log_file_path):
        raise FileNotFoundError(f'Log file {log_file_path} does not exist.')
    formatted_lines = []
    with open(log_file_path, 'r') as log:
        for line in log:
            for keyword in keywords:
                if keyword in line:
                    parts = re.split('\\s+', line.strip(), maxsplit=2)
                    if len(parts) == 3:
                        formatted_line = f'{keyword:>{20}} : {parts[1]:>{20}} : {parts[2]:>{20}}'
                        formatted_lines.append(formatted_line)
                    else:
                        formatted_lines.append(f'Line format unexpected: {line.strip()}')
    return formatted_lines
