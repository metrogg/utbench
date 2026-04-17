import json
import os
import re

def task_func(file_path, attribute, INPUT_JSON={'type': 'object', 'properties': {'name': {'type': str}, 'age': {'type': int}, 'email': {'type': str}}, 'required': ['name', 'age', 'email']}, EMAIL_REGEX='^[a-z0-9]+[\\._]?[a-z0-9]+[@]\\w+[.]\\w{2,3}$'):
    if not os.path.isfile(file_path):
        raise ValueError(f'{file_path} does not exist.')
    with open(file_path, 'r') as f:
        data = json.load(f)
    for key in INPUT_JSON['required']:
        if key not in data:
            raise ValueError(f'{key} is missing from the JSON object.')
        if not isinstance(data[key], INPUT_JSON['properties'][key]['type']):
            raise ValueError(f"{key} is not of type {INPUT_JSON['properties'][key]['type']}.")
    if 'email' in data and (not re.fullmatch(EMAIL_REGEX, data['email'])):
        raise ValueError('Email is not valid.')
    return data[attribute]
