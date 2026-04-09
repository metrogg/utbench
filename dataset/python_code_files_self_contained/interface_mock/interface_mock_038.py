import subprocess
from ipaddress import IPv4Network

def task_func(ip_range):
    active_ips = {}
    for ip in IPv4Network(ip_range):
        try:
            subprocess.check_output(f'ping -c 1 {ip}', shell=True)
            active_ips[str(ip)] = True
        except subprocess.CalledProcessError:
            active_ips[str(ip)] = False
    return active_ips
