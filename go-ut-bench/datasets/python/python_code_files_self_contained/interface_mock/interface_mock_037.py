import subprocess
import ftplib
import os

def task_func(ftp_server='ftp.dlptest.com', ftp_user='dlpuser', ftp_password='your_ftp_password_here', ftp_dir='/ftp/test'):
    try:
        ftp_obj = ftplib.FTP(ftp_server)
    except Exception as e:
        raise Exception(f'Failed to connect to FTP server {ftp_server}: {str(e)}')
    try:
        ftp_obj.login(ftp_user, ftp_password)
    except Exception as e:
        raise Exception(f'Failed to log into FTP server {ftp_server} with user {ftp_user}: {str(e)}')
    try:
        ftp_obj.cwd(ftp_dir)
    except Exception as e:
        raise Exception(f'Failed to change to directory {ftp_dir} on server {ftp_server}: {str(e)}')
    download_dir = 'downloaded_files'
    if not os.path.exists(download_dir):
        os.makedirs(download_dir)
    downloaded_files = []
    for filename in ftp_obj.nlst():
        command = f'wget ftp://{ftp_user}:{ftp_password}@{ftp_server}{ftp_dir}/{filename} -P {download_dir}'
        subprocess.call(command, shell=True)
        downloaded_files.append(filename)
    ftp_obj.quit()
    return downloaded_files
