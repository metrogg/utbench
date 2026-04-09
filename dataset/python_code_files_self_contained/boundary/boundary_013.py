import re
import socket

def task_func(ip_addresses: list) -> dict:
    IP_REGEX = '[0-9]+(?:\\.[0-9]+){3}'
    hostnames = {}
    for ip in ip_addresses:
        if re.match(IP_REGEX, ip):
            try:
                hostname = socket.gethostbyaddr(ip)[0]
                hostnames[ip] = hostname
            except (socket.herror, socket.gaierror):
                hostnames[ip] = None
    return hostnames
