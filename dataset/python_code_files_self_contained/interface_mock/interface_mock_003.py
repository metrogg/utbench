import re
import urllib.parse
import ssl
import socket

def task_func(myString):
    urls = re.findall('https://[^\\s,]+', myString)
    ssl_expiry_dates = {}
    for url in urls:
        try:
            domain = urllib.parse.urlparse(url).netloc
            context = ssl.create_default_context()
            with socket.create_connection((domain, 443)) as sock:
                with context.wrap_socket(sock, server_hostname=domain) as ssock:
                    ssl_expiry_dates[domain] = ssock.getpeercert()['notAfter']
        except ssl.SSLError:
            continue
    return ssl_expiry_dates
