import ast
import re

def task_func(text_file: str) -> list:
    with open(text_file, 'r') as file:
        text = file.read()
    pattern = re.compile('\\{[^{}]*\\{[^{}]*\\}[^{}]*\\}|\\{[^{}]*\\}')
    matches = pattern.findall(text)
    results = [ast.literal_eval(match) for match in matches]
    return results
