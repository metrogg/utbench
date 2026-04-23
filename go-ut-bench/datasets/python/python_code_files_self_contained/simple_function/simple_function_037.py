import inspect
import types
import json

def task_func(f):
    spec = inspect.getfullargspec(f)
    annotations = {k: v.__name__ if isinstance(v, type) else str(v) for k, v in spec.annotations.items()}
    info = {'function_name': f.__name__, 'args': spec.args, 'defaults': spec.defaults, 'annotations': annotations, 'is_lambda': isinstance(f, types.LambdaType)}
    return json.dumps(info)
