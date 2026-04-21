import binascii
import urllib.parse

def task_func(url):
    try:
        parsed_url = urllib.parse.urlparse(url)
        query = urllib.parse.parse_qs(parsed_url.query).get('q', [None])[0]
        return binascii.unhexlify(query).decode('utf-8') if query else None
    except (binascii.Error, UnicodeDecodeError):
        return None
