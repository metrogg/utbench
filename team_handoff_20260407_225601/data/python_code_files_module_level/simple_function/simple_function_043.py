import functools
import os
import platform
import sys
import textwrap
import types
import warnings
import numpy as np
from numpy._core import ndarray
from numpy._utils import set_module
__all__ = ['get_include', 'info', 'show_runtime']

@set_module('numpy')
def show_runtime():
    from pprint import pprint
    from numpy._core._multiarray_umath import __cpu_baseline__, __cpu_dispatch__, __cpu_features__
    config_found = [{'numpy_version': np.__version__, 'python': sys.version, 'uname': platform.uname()}]
    features_found, features_not_found = ([], [])
    for feature in __cpu_dispatch__:
        if __cpu_features__[feature]:
            features_found.append(feature)
        else:
            features_not_found.append(feature)
    config_found.append({'simd_extensions': {'baseline': __cpu_baseline__, 'found': features_found, 'not_found': features_not_found}})
    config_found.append({'ignore_floating_point_errors_in_matmul': not np._core._multiarray_umath._blas_supports_fpe(None)})
    try:
        from threadpoolctl import threadpool_info
        config_found.extend(threadpool_info())
    except ImportError:
        print('WARNING: `threadpoolctl` not found in system! Install it by `pip install threadpoolctl`. Once installed, try `np.show_runtime` again for more detailed build information')
    pprint(config_found)

@set_module('numpy')
def get_include():
    import numpy
    if numpy.show_config is None:
        d = os.path.join(os.path.dirname(numpy.__file__), '_core', 'include')
    else:
        import numpy._core as _core
        d = os.path.join(os.path.dirname(_core.__file__), 'include')
    return d

class _Deprecate:

    def __init__(self, old_name=None, new_name=None, message=None):
        self.old_name = old_name
        self.new_name = new_name
        self.message = message

    def __call__(self, func, *args, **kwargs):
        old_name = self.old_name
        new_name = self.new_name
        message = self.message
        if old_name is None:
            old_name = func.__name__
        if new_name is None:
            depdoc = f'`{old_name}` is deprecated!'
        else:
            depdoc = f'`{old_name}` is deprecated, use `{new_name}` instead!'
        if message is not None:
            depdoc += '\n' + message

        @functools.wraps(func)
        def newfunc(*args, **kwds):
            warnings.warn(depdoc, DeprecationWarning, stacklevel=2)
            return func(*args, **kwds)
        newfunc.__name__ = old_name
        doc = func.__doc__
        if doc is None:
            doc = depdoc
        else:
            lines = doc.expandtabs().split('\n')
            indent = _get_indent(lines[1:])
            if lines[0].lstrip():
                doc = indent * ' ' + doc
            else:
                skip = len(lines[0]) + 1
                for line in lines[1:]:
                    if len(line) > indent:
                        break
                    skip += len(line) + 1
                doc = doc[skip:]
            depdoc = textwrap.indent(depdoc, ' ' * indent)
            doc = f'{depdoc}\n\n{doc}'
        newfunc.__doc__ = doc
        return newfunc

def _get_indent(lines):
    indent = sys.maxsize
    for line in lines:
        content = len(line.lstrip())
        if content:
            indent = min(indent, len(line) - content)
    if indent == sys.maxsize:
        indent = 0
    return indent

def deprecate(*args, **kwargs):
    warnings.warn('`deprecate` is deprecated, use `warn` with `DeprecationWarning` instead. (deprecated in NumPy 2.0)', DeprecationWarning, stacklevel=2)
    if args:
        fn = args[0]
        args = args[1:]
        return _Deprecate(*args, **kwargs)(fn)
    else:
        return _Deprecate(*args, **kwargs)

def deprecate_with_doc(msg):
    warnings.warn('`deprecate` is deprecated, use `warn` with `DeprecationWarning` instead. (deprecated in NumPy 2.0)', DeprecationWarning, stacklevel=2)
    return _Deprecate(message=msg)

def _split_line(name, arguments, width):
    firstwidth = len(name)
    k = firstwidth
    newstr = name
    sepstr = ', '
    arglist = arguments.split(sepstr)
    for argument in arglist:
        if k == firstwidth:
            addstr = ''
        else:
            addstr = sepstr
        k = k + len(argument) + len(addstr)
        if k > width:
            k = firstwidth + 1 + len(argument)
            newstr = newstr + ',\n' + ' ' * (firstwidth + 2) + argument
        else:
            newstr = newstr + addstr + argument
    return newstr
_namedict = None
_dictlist = None

def _makenamedict(module='numpy'):
    module = __import__(module, globals(), locals(), [])
    thedict = {module.__name__: module.__dict__}
    dictlist = [module.__name__]
    totraverse = [module.__dict__]
    while True:
        if len(totraverse) == 0:
            break
        thisdict = totraverse.pop(0)
        for x in thisdict.keys():
            if isinstance(thisdict[x], types.ModuleType):
                modname = thisdict[x].__name__
                if modname not in dictlist:
                    moddict = thisdict[x].__dict__
                    dictlist.append(modname)
                    totraverse.append(moddict)
                    thedict[modname] = moddict
    return (thedict, dictlist)

def _info(obj, output=None):
    extra = ''
    tic = ''
    bp = lambda x: x
    cls = getattr(obj, '__class__', type(obj))
    nm = getattr(cls, '__name__', cls)
    strides = obj.strides
    endian = obj.dtype.byteorder
    if output is None:
        output = sys.stdout
    print('class: ', nm, file=output)
    print('shape: ', obj.shape, file=output)
    print('strides: ', strides, file=output)
    print('itemsize: ', obj.itemsize, file=output)
    print('aligned: ', bp(obj.flags.aligned), file=output)
    print('contiguous: ', bp(obj.flags.contiguous), file=output)
    print('fortran: ', obj.flags.fortran, file=output)
    print(f'data pointer: {hex(obj.ctypes._as_parameter_.value)}{extra}', file=output)
    print('byteorder: ', end=' ', file=output)
    if endian in ['|', '=']:
        print(f'{tic}{sys.byteorder}{tic}', file=output)
        byteswap = False
    elif endian == '>':
        print(f'{tic}big{tic}', file=output)
        byteswap = sys.byteorder != 'big'
    else:
        print(f'{tic}little{tic}', file=output)
        byteswap = sys.byteorder != 'little'
    print('byteswap: ', bp(byteswap), file=output)
    print(f'type: {obj.dtype}', file=output)

@set_module('numpy')
def info(object=None, maxwidth=76, output=None, toplevel='numpy'):
    global _namedict, _dictlist
    import inspect
    import pydoc
    if hasattr(object, '_ppimport_importer') or hasattr(object, '_ppimport_module'):
        object = object._ppimport_module
    elif hasattr(object, '_ppimport_attr'):
        object = object._ppimport_attr
    if output is None:
        output = sys.stdout
    if object is None:
        info(info)
    elif isinstance(object, ndarray):
        _info(object, output=output)
    elif isinstance(object, str):
        if _namedict is None:
            _namedict, _dictlist = _makenamedict(toplevel)
        numfound = 0
        objlist = []
        for namestr in _dictlist:
            try:
                obj = _namedict[namestr][object]
                if id(obj) in objlist:
                    print(f'\n     *** Repeat reference found in {namestr} *** ', file=output)
                else:
                    objlist.append(id(obj))
                    print(f'     *** Found in {namestr} ***', file=output)
                    info(obj)
                    print('-' * maxwidth, file=output)
                numfound += 1
            except KeyError:
                pass
        if numfound == 0:
            print(f'Help for {object} not found.', file=output)
        else:
            print('\n     *** Total of %d references found. ***' % numfound, file=output)
    elif inspect.isfunction(object) or inspect.ismethod(object):
        name = object.__name__
        try:
            arguments = str(inspect.signature(object))
        except Exception:
            arguments = '()'
        if len(name + arguments) > maxwidth:
            argstr = _split_line(name, arguments, maxwidth)
        else:
            argstr = name + arguments
        print(' ' + argstr + '\n', file=output)
        print(inspect.getdoc(object), file=output)
    elif inspect.isclass(object):
        name = object.__name__
        try:
            arguments = str(inspect.signature(object))
        except Exception:
            arguments = '()'
        if len(name + arguments) > maxwidth:
            argstr = _split_line(name, arguments, maxwidth)
        else:
            argstr = name + arguments
        print(' ' + argstr + '\n', file=output)
        doc1 = inspect.getdoc(object)
        if doc1 is None:
            if hasattr(object, '__init__'):
                print(inspect.getdoc(object.__init__), file=output)
        else:
            print(inspect.getdoc(object), file=output)
        methods = pydoc.allmethods(object)
        public_methods = [meth for meth in methods if meth[0] != '_']
        if public_methods:
            print('\n\nMethods:\n', file=output)
            for meth in public_methods:
                thisobj = getattr(object, meth, None)
                if thisobj is not None:
                    methstr, other = pydoc.splitdoc(inspect.getdoc(thisobj) or 'None')
                print(f'  {meth}  --  {methstr}', file=output)
    elif hasattr(object, '__doc__'):
        print(inspect.getdoc(object), file=output)

def safe_eval(source):
    warnings.warn('`safe_eval` is deprecated. Use `ast.literal_eval` instead. Be aware of security implications, such as memory exhaustion based attacks (deprecated in NumPy 2.0)', DeprecationWarning, stacklevel=2)
    import ast
    return ast.literal_eval(source)

def _median_nancheck(data, result, axis):
    if data.size == 0:
        return result
    potential_nans = data.take(-1, axis=axis)
    n = np.isnan(potential_nans)
    if np.ma.isMaskedArray(n):
        n = n.filled(False)
    if not n.any():
        return result
    if isinstance(result, np.generic):
        return potential_nans
    np.copyto(result, potential_nans, where=n)
    return result

def _opt_info():
    from numpy._core._multiarray_umath import __cpu_baseline__, __cpu_dispatch__, __cpu_features__
    if len(__cpu_baseline__) == 0 and len(__cpu_dispatch__) == 0:
        return ''
    enabled_features = ' '.join(__cpu_baseline__)
    for feature in __cpu_dispatch__:
        if __cpu_features__[feature]:
            enabled_features += f' {feature}*'
        else:
            enabled_features += f' {feature}?'
    return enabled_features

def drop_metadata(dtype, /):
    if dtype.fields is not None:
        found_metadata = dtype.metadata is not None
        names = []
        formats = []
        offsets = []
        titles = []
        for name, field in dtype.fields.items():
            field_dt = drop_metadata(field[0])
            if field_dt is not field[0]:
                found_metadata = True
            names.append(name)
            formats.append(field_dt)
            offsets.append(field[1])
            titles.append(None if len(field) < 3 else field[2])
        if not found_metadata:
            return dtype
        structure = {'names': names, 'formats': formats, 'offsets': offsets, 'titles': titles, 'itemsize': dtype.itemsize}
        return np.dtype(structure, align=dtype.isalignedstruct)
    elif dtype.subdtype is not None:
        subdtype, shape = dtype.subdtype
        new_subdtype = drop_metadata(subdtype)
        if dtype.metadata is None and new_subdtype is subdtype:
            return dtype
        return np.dtype((new_subdtype, shape))
    else:
        if dtype.metadata is None:
            return dtype
        return np.dtype(dtype.str)
