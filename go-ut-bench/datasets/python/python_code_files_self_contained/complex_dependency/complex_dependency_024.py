import re
import socket
import urllib.parse

def task_func(myString):
    urls = re.findall('https?://[^\\s,]+', myString)
    ip_addresses = {}
    for url in urls:
        domain = urllib.parse.urlparse(url).netloc
        try:
            ip_addresses[domain] = socket.gethostbyname(domain)
        except socket.gaierror:
            ip_addresses[domain] = None
    return ip_addresses
