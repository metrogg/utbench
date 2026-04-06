from pyarrow._compute import Function, FunctionOptions, FunctionRegistry, HashAggregateFunction, HashAggregateKernel, Kernel, ScalarAggregateFunction, ScalarAggregateKernel, ScalarFunction, ScalarKernel, VectorFunction, VectorKernel, ArraySortOptions, AssumeTimezoneOptions, CastOptions, CountOptions, CumulativeOptions, CumulativeSumOptions, DayOfWeekOptions, DictionaryEncodeOptions, RunEndEncodeOptions, ElementWiseAggregateOptions, ExtractRegexOptions, ExtractRegexSpanOptions, FilterOptions, IndexOptions, InversePermutationOptions, JoinOptions, ListSliceOptions, ListFlattenOptions, MakeStructOptions, MapLookupOptions, MatchSubstringOptions, ModeOptions, NullOptions, PadOptions, PairwiseOptions, PartitionNthOptions, PivotWiderOptions, QuantileOptions, RandomOptions, RankOptions, RankQuantileOptions, ReplaceSliceOptions, ReplaceSubstringOptions, RoundBinaryOptions, RoundOptions, RoundTemporalOptions, RoundToMultipleOptions, ScalarAggregateOptions, ScatterOptions, SelectKOptions, SetLookupOptions, SkewOptions, SliceOptions, SortOptions, SplitOptions, SplitPatternOptions, StrftimeOptions, StrptimeOptions, StructFieldOptions, TakeOptions, TDigestOptions, TrimOptions, Utf8NormalizeOptions, VarianceOptions, WeekOptions, WinsorizeOptions, ZeroFillOptions, call_function, function_registry, get_function, list_functions, call_tabular_function, register_scalar_function, register_tabular_function, register_aggregate_function, register_vector_function, UdfContext, Expression
from collections import namedtuple
import inspect
from textwrap import dedent
import warnings
import pyarrow as pa
from pyarrow import _compute_docstrings
from pyarrow.vendored import docscrape

def _get_arg_names(func):
    return func._doc.arg_names
_OptionsClassDoc = namedtuple('_OptionsClassDoc', ('params',))

def _scrape_options_class_doc(options_class):
    if not options_class.__doc__:
        return None
    doc = docscrape.NumpyDocString(options_class.__doc__)
    return _OptionsClassDoc(doc['Parameters'])

def _decorate_compute_function(wrapper, exposed_name, func, options_class):
    cpp_doc = func._doc
    wrapper.__arrow_compute_function__ = dict(name=func.name, arity=func.arity, options_class=cpp_doc.options_class, options_required=cpp_doc.options_required)
    wrapper.__name__ = exposed_name
    wrapper.__qualname__ = exposed_name
    doc_pieces = []
    summary = cpp_doc.summary
    if not summary:
        arg_str = 'arguments' if func.arity > 1 else 'argument'
        summary = f'Call compute function {func.name!r} with the given {arg_str}'
    doc_pieces.append(f'{summary}.\n\n')
    description = cpp_doc.description
    if description:
        doc_pieces.append(f'{description}\n\n')
    doc_addition = _compute_docstrings.function_doc_additions.get(func.name)
    doc_pieces.append(dedent('        Parameters\n        ----------\n        '))
    arg_names = _get_arg_names(func)
    for arg_name in arg_names:
        if func.kind in ('vector', 'scalar_aggregate'):
            arg_type = 'Array-like'
        else:
            arg_type = 'Array-like or scalar-like'
        doc_pieces.append(f'{arg_name} : {arg_type}\n')
        doc_pieces.append('    Argument to compute function.\n')
    if options_class is not None:
        options_class_doc = _scrape_options_class_doc(options_class)
        if options_class_doc:
            for p in options_class_doc.params:
                doc_pieces.append(f'{p.name} : {p.type}\n')
                for s in p.desc:
                    doc_pieces.append(f'    {s}\n')
        else:
            warnings.warn(f'Options class {options_class.__name__} does not have a docstring', RuntimeWarning)
            options_sig = inspect.signature(options_class)
            for p in options_sig.parameters.values():
                doc_pieces.append(dedent(f'                {p.name} : optional\n                    Parameter for {options_class.__name__} constructor. Either `options`\n                    or `{p.name}` can be passed, but not both at the same time.\n                '))
        doc_pieces.append(dedent(f'            options : pyarrow.compute.{options_class.__name__}, optional\n                Alternative way of passing options.\n            '))
    doc_pieces.append(dedent('        memory_pool : pyarrow.MemoryPool, optional\n            If not passed, will allocate memory from the default memory pool.\n        '))
    if doc_addition is not None:
        stripped = dedent(doc_addition).strip('\n')
        doc_pieces.append(f'\n{stripped}\n')
    wrapper.__doc__ = ''.join(doc_pieces)
    return wrapper

def _get_options_class(func):
    class_name = func._doc.options_class
    if not class_name:
        return None
    try:
        return globals()[class_name]
    except KeyError:
        warnings.warn(f'Python binding for {class_name} not exposed', RuntimeWarning)
        return None

def _handle_options(name, options_class, options, args, kwargs):
    if args or kwargs:
        if options is not None:
            raise TypeError(f"Function {name!r} called with both an 'options' argument and additional arguments")
        return options_class(*args, **kwargs)
    if options is not None:
        if isinstance(options, dict):
            return options_class(**options)
        elif isinstance(options, options_class):
            return options
        raise TypeError(f'Function {name!r} expected a {options_class} parameter, got {type(options)}')
    return None

def _make_generic_wrapper(func_name, func, options_class, arity):
    if options_class is None:

        def wrapper(*args, memory_pool=None):
            if arity is not Ellipsis and len(args) != arity:
                raise TypeError(f'{func_name} takes {arity} positional argument(s), but {len(args)} were given')
            if args and isinstance(args[0], Expression):
                return Expression._call(func_name, list(args))
            return func.call(args, None, memory_pool)
    else:

        def wrapper(*args, memory_pool=None, options=None, **kwargs):
            if arity is not Ellipsis:
                if len(args) < arity:
                    raise TypeError(f'{func_name} takes {arity} positional argument(s), but {len(args)} were given')
                option_args = args[arity:]
                args = args[:arity]
            else:
                option_args = ()
            options = _handle_options(func_name, options_class, options, option_args, kwargs)
            if args and isinstance(args[0], Expression):
                return Expression._call(func_name, list(args), options)
            return func.call(args, options, memory_pool)
    return wrapper

def _make_signature(arg_names, var_arg_names, options_class):
    from inspect import Parameter
    params = []
    for name in arg_names:
        params.append(Parameter(name, Parameter.POSITIONAL_ONLY))
    for name in var_arg_names:
        params.append(Parameter(name, Parameter.VAR_POSITIONAL))
    if options_class is not None:
        options_sig = inspect.signature(options_class)
        for p in options_sig.parameters.values():
            assert p.kind in (Parameter.POSITIONAL_OR_KEYWORD, Parameter.KEYWORD_ONLY)
            if var_arg_names:
                p = p.replace(kind=Parameter.KEYWORD_ONLY)
            params.append(p)
        params.append(Parameter('options', Parameter.KEYWORD_ONLY, default=None))
    params.append(Parameter('memory_pool', Parameter.KEYWORD_ONLY, default=None))
    return inspect.Signature(params)

def _wrap_function(name, func):
    options_class = _get_options_class(func)
    arg_names = _get_arg_names(func)
    has_vararg = arg_names and arg_names[-1].startswith('*')
    if has_vararg:
        var_arg_names = [arg_names.pop().lstrip('*')]
    else:
        var_arg_names = []
    wrapper = _make_generic_wrapper(name, func, options_class, arity=func.arity)
    wrapper.__signature__ = _make_signature(arg_names, var_arg_names, options_class)
    return _decorate_compute_function(wrapper, name, func, options_class)

def _make_global_functions():
    g = globals()
    reg = function_registry()
    rewrites = {'and': 'and_', 'or': 'or_'}
    for cpp_name in reg.list_functions():
        name = rewrites.get(cpp_name, cpp_name)
        func = reg.get_function(cpp_name)
        if func.kind == 'hash_aggregate':
            continue
        if func.kind == 'scalar_aggregate' and func.arity == 0:
            continue
        assert name not in g, name
        g[cpp_name] = g[name] = _wrap_function(name, func)
_make_global_functions()
utf8_zfill = utf8_zero_fill = globals()['utf8_zero_fill']

def cast(arr, target_type=None, safe=None, options=None, memory_pool=None):
    safe_vars_passed = safe is not None or target_type is not None
    if safe_vars_passed and options is not None:
        raise ValueError("Must either pass values for 'target_type' and 'safe' or pass a value for 'options'")
    if options is None:
        target_type = pa.types.lib.ensure_type(target_type)
        if safe is False:
            options = CastOptions.unsafe(target_type)
        else:
            options = CastOptions.safe(target_type)
    return call_function('cast', [arr], options, memory_pool)

def index(data, value, start=None, end=None, *, memory_pool=None):
    if start is not None:
        if end is not None:
            data = data.slice(start, end - start)
        else:
            data = data.slice(start)
    elif end is not None:
        data = data.slice(0, end)
    if not isinstance(value, pa.Scalar):
        value = pa.scalar(value, type=data.type)
    elif data.type != value.type:
        value = pa.scalar(value.as_py(), type=data.type)
    options = IndexOptions(value=value)
    result = call_function('index', [data], options, memory_pool)
    if start is not None and result.as_py() >= 0:
        result = pa.scalar(result.as_py() + start, type=pa.int64())
    return result

def take(data, indices, *, boundscheck=True, memory_pool=None):
    options = TakeOptions(boundscheck=boundscheck)
    return call_function('take', [data, indices], options, memory_pool)

def fill_null(values, fill_value):
    if not isinstance(fill_value, (pa.Array, pa.ChunkedArray, pa.Scalar)):
        fill_value = pa.scalar(fill_value, type=values.type)
    elif values.type != fill_value.type:
        fill_value = pa.scalar(fill_value.as_py(), type=values.type)
    return call_function('coalesce', [values, fill_value])

def top_k_unstable(values, k, sort_keys=None, *, memory_pool=None):
    if sort_keys is None:
        sort_keys = []
    if isinstance(values, (pa.Array, pa.ChunkedArray)):
        sort_keys.append(('dummy', 'descending'))
    else:
        sort_keys = map(lambda key_name: (key_name, 'descending'), sort_keys)
    options = SelectKOptions(k, sort_keys)
    return call_function('select_k_unstable', [values], options, memory_pool)

def bottom_k_unstable(values, k, sort_keys=None, *, memory_pool=None):
    if sort_keys is None:
        sort_keys = []
    if isinstance(values, (pa.Array, pa.ChunkedArray)):
        sort_keys.append(('dummy', 'ascending'))
    else:
        sort_keys = map(lambda key_name: (key_name, 'ascending'), sort_keys)
    options = SelectKOptions(k, sort_keys)
    return call_function('select_k_unstable', [values], options, memory_pool)

def random(n, *, initializer='system', options=None, memory_pool=None):
    options = RandomOptions(initializer=initializer)
    return call_function('random', [], options, memory_pool, length=n)

def field(*name_or_index):
    n = len(name_or_index)
    if n == 1:
        if isinstance(name_or_index[0], (str, int)):
            return Expression._field(name_or_index[0])
        elif isinstance(name_or_index[0], tuple):
            return Expression._nested_field(name_or_index[0])
        else:
            raise TypeError(f'field reference should be str, multiple str, tuple or integer, got {type(name_or_index[0])}')
    else:
        return Expression._nested_field(name_or_index)

def scalar(value):
    return Expression._scalar(value)
