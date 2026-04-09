import dis
from inspect import ismethod, isfunction, istraceback, isframe, iscode
from .pointers import parent, reference, at, parents, children
from .logger import trace
__all__ = ['baditems', 'badobjects', 'badtypes', 'code', 'errors', 'freevars', 'getmodule', 'globalvars', 'nestedcode', 'nestedglobals', 'outermost', 'referredglobals', 'referrednested', 'trace', 'varnames']

def getmodule(object, _filename=None, force=False):
    from inspect import getmodule as getmod
    module = getmod(object, _filename)
    if module or not force:
        return module
    import builtins
    from .source import getname
    name = getname(object, force=True)
    return builtins if name in vars(builtins).keys() else None

def outermost(func):
    if ismethod(func):
        _globals = func.__func__.__globals__ or {}
    elif isfunction(func):
        _globals = func.__globals__ or {}
    else:
        return
    _globals = _globals.items()
    from .source import getsourcelines
    try:
        lines, lnum = getsourcelines(func, enclosing=True)
    except Exception:
        lines, lnum = ([], None)
    code = ''.join(lines)
    _locals = ((name, obj) for name, obj in _globals if name in code)
    for name, obj in _locals:
        try:
            if getsourcelines(obj) == (lines, lnum):
                return obj
        except Exception:
            pass
    return

def nestedcode(func, recurse=True):
    func = code(func)
    if not iscode(func):
        return []
    nested = set()
    for co in func.co_consts:
        if co is None:
            continue
        co = code(co)
        if co:
            nested.add(co)
            if recurse:
                nested |= set(nestedcode(co, recurse=True))
    return list(nested)

def code(func):
    if ismethod(func):
        func = func.__func__
    if isfunction(func):
        func = func.__code__
    if istraceback(func):
        func = func.tb_frame
    if isframe(func):
        func = func.f_code
    if iscode(func):
        return func
    return

def referrednested(func, recurse=True):
    import gc
    funcs = set()
    for co in nestedcode(func, recurse):
        for obj in gc.get_referrers(co):
            _ = getattr(obj, '__func__', None)
            if getattr(_, '__code__', None) is co:
                funcs.add(obj)
            elif getattr(obj, '__code__', None) is co:
                funcs.add(obj)
            elif getattr(obj, 'f_code', None) is co:
                funcs.add(obj)
            elif hasattr(obj, 'co_code') and obj is co:
                funcs.add(obj)
    return list(funcs)

def freevars(func):
    if ismethod(func):
        func = func.__func__
    if isfunction(func):
        closures = func.__closure__ or ()
        func = func.__code__.co_freevars
    else:
        return {}

    def get_cell_contents():
        for name, c in zip(func, closures):
            try:
                cell_contents = c.cell_contents
            except ValueError:
                continue
            yield (name, c.cell_contents)
    return dict(get_cell_contents())

def nestedglobals(func, recurse=True):
    func = code(func)
    if func is None:
        return list()
    import sys
    from .temp import capture
    CAN_NULL = sys.hexversion >= 51052711
    names = set()
    with capture('stdout') as out:
        try:
            dis.dis(func)
        except IndexError:
            pass
    for line in out.getvalue().splitlines():
        if '_GLOBAL' in line:
            name = line.split('(')[-1].split(')')[0]
            if CAN_NULL:
                names.add(name.replace('NULL + ', '').replace(' + NULL', ''))
            else:
                names.add(name)
    for co in getattr(func, 'co_consts', tuple()):
        if co and recurse and iscode(co):
            names.update(nestedglobals(co, recurse=True))
    return list(names)

def referredglobals(func, recurse=True, builtin=False):
    return globalvars(func, recurse, builtin).keys()

def globalvars(func, recurse=True, builtin=False):
    if ismethod(func):
        func = func.__func__
    if isfunction(func):
        globs = vars(getmodule(sum)).copy() if builtin else {}
        orig_func, func = (func, set())
        for obj in orig_func.__closure__ or {}:
            try:
                cell_contents = obj.cell_contents
            except ValueError:
                pass
            else:
                _vars = globalvars(cell_contents, recurse, builtin) or {}
                func.update(_vars)
                globs.update(_vars)
        globs.update(orig_func.__globals__ or {})
        if not recurse:
            func.update(orig_func.__code__.co_names)
        else:
            func.update(nestedglobals(orig_func.__code__))
            for key in func.copy():
                nested_func = globs.get(key)
                if nested_func is orig_func:
                    continue
                func.update(globalvars(nested_func, True, builtin))
    elif iscode(func):
        globs = vars(getmodule(sum)).copy() if builtin else {}
        if not recurse:
            func = func.co_names
        else:
            orig_func = func.co_name
            func = set(nestedglobals(func))
            for key in func.copy():
                if key is orig_func:
                    continue
                nested_func = globs.get(key)
                func.update(globalvars(nested_func, True, builtin))
    else:
        return {}
    return dict(((name, globs[name]) for name in func if name in globs))

def varnames(func):
    func = code(func)
    if not iscode(func):
        return ()
    return (func.co_varnames, func.co_cellvars)

def baditems(obj, exact=False, safe=False):
    if not hasattr(obj, '__iter__'):
        return [j for j in (badobjects(obj, 0, exact, safe),) if j is not None]
    obj = obj.values() if getattr(obj, 'values', None) else obj
    _obj = []
    [_obj.append(badobjects(i, 0, exact, safe)) for i in obj if i not in _obj]
    return [j for j in _obj if j is not None]

def badobjects(obj, depth=0, exact=False, safe=False):
    from dill import pickles
    if not depth:
        if pickles(obj, exact, safe):
            return None
        return obj
    return dict(((attr, badobjects(getattr(obj, attr), depth - 1, exact, safe)) for attr in dir(obj) if not pickles(getattr(obj, attr), exact, safe)))

def badtypes(obj, depth=0, exact=False, safe=False):
    from dill import pickles
    if not depth:
        if pickles(obj, exact, safe):
            return None
        return type(obj)
    return dict(((attr, badtypes(getattr(obj, attr), depth - 1, exact, safe)) for attr in dir(obj) if not pickles(getattr(obj, attr), exact, safe)))

def errors(obj, depth=0, exact=False, safe=False):
    from dill import pickles, copy
    if not depth:
        try:
            pik = copy(obj)
            if exact:
                assert pik == obj, 'Unpickling produces %s instead of %s' % (pik, obj)
            assert type(pik) == type(obj), 'Unpickling produces %s instead of %s' % (type(pik), type(obj))
            return None
        except Exception:
            import sys
            return sys.exc_info()[1]
    _dict = {}
    for attr in dir(obj):
        try:
            _attr = getattr(obj, attr)
        except Exception:
            import sys
            _dict[attr] = sys.exc_info()[1]
            continue
        if not pickles(_attr, exact, safe):
            _dict[attr] = errors(_attr, depth - 1, exact, safe)
    return _dict
