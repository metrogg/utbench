import json
import urllib.parse
import hmac
import hashlib

def task_func(req_data, secret_key):
    if not isinstance(req_data, dict):
        raise TypeError('req_data must be a dictionary')
    json_req_data = json.dumps(req_data)
    hmac_obj = hmac.new(secret_key.encode(), json_req_data.encode(), hashlib.sha256)
    hmac_signature = hmac_obj.hexdigest()
    url_encoded_signature = urllib.parse.quote_plus(hmac_signature)
    return url_encoded_signature
