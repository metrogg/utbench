import itertools
import json

def task_func(json_list, r):
    try:
        data = json.loads(json_list)
        number_list = data['number_list']
        return list(itertools.combinations(number_list, r))
    except Exception as e:
        raise e
