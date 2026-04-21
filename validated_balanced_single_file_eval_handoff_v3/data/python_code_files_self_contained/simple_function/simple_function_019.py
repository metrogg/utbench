import re
from urllib import request
import json
IP_REGEX = '[0-9]+(?:\\.[0-9]+){3}'

def task_func(ip_address):
    try:
        response = ip_address
        data = json.loads(response)
        ip = data['ip']
        if re.match(IP_REGEX, ip):
            return ip
        else:
            return 'Invalid IP address received'
    except Exception as e:
        return str(e)
