import re
import urllib.request
import json
IP_REGEX = '[0-9]+(?:\\.[0-9]+){3}'

def task_func(API_URL):
    try:
        response = urllib.request.urlopen(API_URL)
        data = json.loads(response.read())
        ip = data['ip']
        if re.match(IP_REGEX, ip):
            return ip
        else:
            return 'Invalid IP address received'
    except Exception as e:
        return str(e)
