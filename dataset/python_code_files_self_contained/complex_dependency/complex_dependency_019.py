import urllib.request
import zipfile
import os
import urllib.error

def task_func(url: str, save_path: str='downloaded_file.zip', extract_path: str='extracted_files') -> str:
    try:
        if os.path.exists(save_path):
            os.remove(save_path)
        urllib.request.urlretrieve(url, save_path)
        if not os.path.exists(extract_path):
            os.makedirs(extract_path)
        with zipfile.ZipFile(save_path, 'r') as zip_ref:
            zip_ref.extractall(extract_path)
        os.remove(save_path)
        return extract_path
    except urllib.error.URLError as e:
        return f'URL Error: {e.reason}'
