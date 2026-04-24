import urllib.request
import os
import zipfile
TARGET_DIR = 'downloaded_files'
TARGET_ZIP_FILE = 'downloaded_files.zip'

def task_func(url):
    os.makedirs(TARGET_DIR, exist_ok=True)
    urllib.request.urlretrieve(url, TARGET_ZIP_FILE)
    with zipfile.ZipFile(TARGET_ZIP_FILE, 'r') as zip_ref:
        zip_ref.extractall(TARGET_DIR)
    if os.path.exists(TARGET_ZIP_FILE):
        os.remove(TARGET_ZIP_FILE)
    return TARGET_DIR
