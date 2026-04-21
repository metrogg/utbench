import json
from datetime import datetime

def task_func(json_data):
    try:
        data = json.loads(json_data)
        datetime_str = data['utc_datetime']
        utc_datetime = datetime.strptime(datetime_str, '%Y-%m-%dT%H:%M:%S')
        return utc_datetime.weekday() >= 5
    except Exception as e:
        raise e
