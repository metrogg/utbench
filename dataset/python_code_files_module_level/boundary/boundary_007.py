from __future__ import annotations
import tokenize
from typing import TYPE_CHECKING, Any
import warnings
from pandas.util._decorators import set_module
from pandas.util._exceptions import find_stack_level
from pandas.util._validators import validate_bool_kwarg
from pandas.core.dtypes.common import is_extension_array_dtype, is_string_dtype
from pandas.core.computation.engines import ENGINES
from pandas.core.computation.expr import PARSERS, Expr
from pandas.core.computation.parsing import tokenize_string
from pandas.core.computation.scope import ensure_scope
from pandas.core.generic import NDFrame
from pandas.io.formats.printing import pprint_thing
if TYPE_CHECKING:
    from pandas.core.computation.ops import BinOp

def _check_engine(engine: str | None) -> str:
    from pandas.core.computation.check import NUMEXPR_INSTALLED
    from pandas.core.computation.expressions import USE_NUMEXPR
    if engine is None:
        engine = 'numexpr' if USE_NUMEXPR else 'python'
    if engine not in ENGINES:
        valid_engines = list(ENGINES.keys())
        raise KeyError(f"Invalid engine '{engine}' passed, valid engines are {valid_engines}")
    if engine == 'numexpr' and (not NUMEXPR_INSTALLED):
        raise ImportError("'numexpr' is not installed or an unsupported version. Cannot use engine='numexpr' for query/eval if 'numexpr' is not installed")
    return engine

def _check_parser(parser: str) -> None:
    if parser not in PARSERS:
        raise KeyError(f"Invalid parser '{parser}' passed, valid parsers are {PARSERS.keys()}")

def _check_resolvers(resolvers) -> None:
    if resolvers is not None:
        for resolver in resolvers:
            if not hasattr(resolver, '__getitem__'):
                name = type(resolver).__name__
                raise TypeError(f"Resolver of type '{name}' does not implement the __getitem__ method")

def _check_expression(expr) -> None:
    if not expr:
        raise ValueError('expr cannot be an empty string')

def _convert_expression(expr) -> str:
    s = pprint_thing(expr)
    _check_expression(s)
    return s

def _check_for_locals(expr: str, stack_level: int, parser: str) -> None:
    at_top_of_stack = stack_level == 0
    not_pandas_parser = parser != 'pandas'
    if not_pandas_parser:
        msg = "The '@' prefix is only supported by the pandas parser"
    elif at_top_of_stack:
        msg = "The '@' prefix is not allowed in top-level eval calls.\nplease refer to your variables by name without the '@' prefix."
    if at_top_of_stack or not_pandas_parser:
        for toknum, tokval in tokenize_string(expr):
            if toknum == tokenize.OP and tokval == '@':
                raise SyntaxError(msg)

@set_module('pandas')
def eval(expr: str | BinOp, parser: str='pandas', engine: str | None=None, local_dict=None, global_dict=None, resolvers=(), level: int=0, target=None, inplace: bool=False) -> Any:
    inplace = validate_bool_kwarg(inplace, 'inplace')
    exprs: list[str | BinOp]
    if isinstance(expr, str):
        _check_expression(expr)
        exprs = [e.strip() for e in expr.splitlines() if e.strip() != '']
    else:
        exprs = [expr]
    multi_line = len(exprs) > 1
    if multi_line and target is None:
        raise ValueError('multi-line expressions are only valid in the context of data, use DataFrame.eval')
    engine = _check_engine(engine)
    _check_parser(parser)
    _check_resolvers(resolvers)
    ret = None
    first_expr = True
    target_modified = False
    for expr in exprs:
        expr = _convert_expression(expr)
        _check_for_locals(expr, level, parser)
        env = ensure_scope(level + 1, global_dict=global_dict, local_dict=local_dict, resolvers=resolvers, target=target)
        parsed_expr = Expr(expr, engine=engine, parser=parser, env=env)
        if engine == 'numexpr' and (is_extension_array_dtype(parsed_expr.terms.return_type) and (not is_string_dtype(parsed_expr.terms.return_type)) or (getattr(parsed_expr.terms, 'operand_types', None) is not None and any((is_extension_array_dtype(elem) and (not is_string_dtype(elem)) for elem in parsed_expr.terms.operand_types)))):
            warnings.warn("Engine has switched to 'python' because numexpr does not support extension array dtypes. Please set your engine to python manually.", RuntimeWarning, stacklevel=find_stack_level())
            engine = 'python'
        eng = ENGINES[engine]
        eng_inst = eng(parsed_expr)
        ret = eng_inst.evaluate()
        if parsed_expr.assigner is None:
            if multi_line:
                raise ValueError('Multi-line expressions are only valid if all expressions contain an assignment')
            if inplace:
                raise ValueError('Cannot operate inplace if there is no assignment')
        assigner = parsed_expr.assigner
        if env.target is not None and assigner is not None:
            target_modified = True
            if not inplace and first_expr:
                try:
                    target = env.target
                    if isinstance(target, NDFrame):
                        target = target.copy(deep=False)
                    else:
                        target = target.copy()
                except AttributeError as err:
                    raise ValueError('Cannot return a copy of the target') from err
            else:
                target = env.target
            try:
                if inplace and isinstance(target, NDFrame):
                    target.loc[:, assigner] = ret
                else:
                    target[assigner] = ret
            except (TypeError, IndexError) as err:
                raise ValueError('Cannot assign expression output to target') from err
            if not resolvers:
                resolvers = ({assigner: ret},)
            else:
                for resolver in resolvers:
                    if assigner in resolver:
                        resolver[assigner] = ret
                        break
                else:
                    resolvers += ({assigner: ret},)
            ret = None
            first_expr = False
    if inplace is False:
        return target if target_modified else ret
