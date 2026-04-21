import json
from datetime import datetime
from decimal import Decimal

def task_func(my_obj):

    class DateTimeEncoder(json.JSONEncoder):

        def default(self, obj):
            if isinstance(obj, datetime):
                return obj.isoformat()
            if isinstance(obj, Decimal):
                return str(obj)
            return json.JSONEncoder.default(self, obj)
    return json.dumps(my_obj, cls=DateTimeEncoder)
