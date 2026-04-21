import inspect
import types
import math

def task_func(f):
    spec = inspect.getfullargspec(f)
    info = {'function_name': f.__name__, 'sqrt_args': math.sqrt(len(spec.args))}
    if spec.defaults:
        info['lambda_in_defaults'] = sum((1 for d in spec.defaults if isinstance(d, types.LambdaType)))
    else:
        info['lambda_in_defaults'] = 0
    return info
