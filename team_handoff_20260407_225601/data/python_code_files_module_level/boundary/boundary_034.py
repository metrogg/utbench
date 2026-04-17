from __future__ import annotations
from collections import abc, defaultdict
import copy
from typing import TYPE_CHECKING, Any, DefaultDict, overload
import numpy as np
from pandas._libs.writers import convert_json_to_lines
from pandas.util._decorators import set_module
import pandas as pd
from pandas import DataFrame, Series
if TYPE_CHECKING:
    from collections.abc import Iterable
    from pandas._typing import IgnoreRaise, Scalar

def convert_to_line_delimits(s: str) -> str:
    if not s[0] == '[' and s[-1] == ']':
        return s
    s = s[1:-1]
    return convert_json_to_lines(s)

@overload
def nested_to_record(ds: dict, prefix: str=..., sep: str=..., level: int=..., max_level: int | None=...) -> dict[str, Any]:
    ...

@overload
def nested_to_record(ds: list[dict], prefix: str=..., sep: str=..., level: int=..., max_level: int | None=...) -> list[dict[str, Any]]:
    ...

def nested_to_record(ds: dict | list[dict], prefix: str='', sep: str='.', level: int=0, max_level: int | None=None) -> dict[str, Any] | list[dict[str, Any]]:
    singleton = False
    if isinstance(ds, dict):
        ds = [ds]
        singleton = True
    new_ds = []
    for d in ds:
        new_d = copy.deepcopy(d)
        for k, v in d.items():
            if not isinstance(k, str):
                k = str(k)
            if level == 0:
                newkey = k
            else:
                newkey = prefix + sep + k
            if not isinstance(v, dict) or (max_level is not None and level >= max_level):
                if level != 0:
                    v = new_d.pop(k)
                    new_d[newkey] = v
                continue
            v = new_d.pop(k)
            new_d.update(nested_to_record(v, newkey, sep, level + 1, max_level))
        new_ds.append(new_d)
    if singleton:
        return new_ds[0]
    return new_ds

def _normalize_json(data: Any, key_string: str, normalized_dict: dict[str, Any], separator: str) -> dict[str, Any]:
    if isinstance(data, dict):
        for key, value in data.items():
            new_key = f'{key_string}{separator}{key}'
            if not key_string:
                new_key = new_key.removeprefix(separator)
            _normalize_json(data=value, key_string=new_key, normalized_dict=normalized_dict, separator=separator)
    else:
        normalized_dict[key_string] = data
    return normalized_dict

def _normalize_json_ordered(data: dict[str, Any], separator: str) -> dict[str, Any]:
    top_dict_ = {k: v for k, v in data.items() if not isinstance(v, dict)}
    nested_dict_ = _normalize_json(data={k: v for k, v in data.items() if isinstance(v, dict)}, key_string='', normalized_dict={}, separator=separator)
    return {**top_dict_, **nested_dict_}

def _simple_json_normalize(ds: dict | list[dict], sep: str='.') -> dict | list[dict] | Any:
    normalized_json_object = {}
    if isinstance(ds, dict):
        normalized_json_object = _normalize_json_ordered(data=ds, separator=sep)
    elif isinstance(ds, list):
        normalized_json_list = [_simple_json_normalize(row, sep=sep) for row in ds]
        return normalized_json_list
    return normalized_json_object

def _validate_meta(meta: str | list[str | list[str]] | None) -> None:
    if meta is None:
        return
    if isinstance(meta, str):
        return
    for item in meta:
        if isinstance(item, list):
            for subitem in item:
                if not isinstance(subitem, str):
                    raise TypeError(f'All elements in nested meta paths must be strings. Found {type(subitem).__name__}: {subitem!r}')
        elif not isinstance(item, str):
            raise TypeError(f"All elements in 'meta' must be strings or lists of strings. Found {type(item).__name__}: {item!r}")

@set_module('pandas')
def json_normalize(data: dict | list[dict] | Series, record_path: str | list | None=None, meta: str | list[str | list[str]] | None=None, meta_prefix: str | None=None, record_prefix: str | None=None, errors: IgnoreRaise='raise', sep: str='.', max_level: int | None=None) -> DataFrame:
    _validate_meta(meta)

    def _pull_field(js: dict[str, Any], spec: list | str, extract_record: bool=False) -> Scalar | Iterable:
        result = js
        try:
            if isinstance(spec, list):
                for field in spec:
                    if result is None:
                        raise KeyError(field)
                    result = result[field]
            else:
                result = result[spec]
        except KeyError as e:
            if extract_record:
                raise KeyError(f'Key {e} not found. If specifying a record_path, all elements of data should have the path.') from e
            if errors == 'ignore':
                return np.nan
            else:
                raise KeyError(f"Key {e} not found. To replace missing values of {e} with np.nan, pass in errors='ignore'") from e
        return result

    def _pull_records(js: dict[str, Any], spec: list | str) -> list:
        result = _pull_field(js, spec, extract_record=True)
        if not isinstance(result, list):
            if pd.isnull(result):
                result = []
            else:
                raise TypeError(f'Path must contain list or null, but got {type(result).__name__} at {spec!r}')
        return result
    if isinstance(data, Series):
        index = data.index
    else:
        index = None
    if isinstance(data, list) and (not data):
        return DataFrame()
    elif isinstance(data, dict):
        data = [data]
    elif isinstance(data, abc.Iterable) and (not isinstance(data, str)):
        data = list(data)
        for item in data:
            if not isinstance(item, dict):
                msg = f'All items in data must be of type dict, found {type(item).__name__}'
                raise TypeError(msg)
    else:
        raise NotImplementedError
    if record_path is None and meta is None and (meta_prefix is None) and (record_prefix is None) and (max_level is None):
        return DataFrame(_simple_json_normalize(data, sep=sep), index=index)
    if record_path is None:
        if any(([isinstance(x, dict) for x in y.values()] for y in data)):
            data = nested_to_record(data, sep=sep, max_level=max_level)
        result = DataFrame(data, index=index)
        if record_prefix is not None:
            result = result.rename(columns=lambda x: f'{record_prefix}{x}')
        return result
    elif not isinstance(record_path, list):
        record_path = [record_path]
    if meta is None:
        meta = []
    elif not isinstance(meta, list):
        meta = [meta]
    _meta = [m if isinstance(m, list) else [m] for m in meta]
    records: list = []
    lengths = []
    meta_vals: DefaultDict = defaultdict(list)
    meta_keys = [sep.join(val) for val in _meta]

    def _recursive_extract(data, path, seen_meta, level: int=0) -> None:
        if isinstance(data, dict):
            data = [data]
        if len(path) > 1:
            for obj in data:
                for val, key in zip(_meta, meta_keys, strict=True):
                    if level + 1 == len(val):
                        seen_meta[key] = _pull_field(obj, val[-1])
                _recursive_extract(obj[path[0]], path[1:], seen_meta, level=level + 1)
        else:
            for obj in data:
                recs = _pull_records(obj, path[0])
                recs = [nested_to_record(r, sep=sep, max_level=max_level) if isinstance(r, dict) else r for r in recs]
                lengths.append(len(recs))
                for val, key in zip(_meta, meta_keys, strict=True):
                    if level + 1 > len(val):
                        meta_val = seen_meta[key]
                    else:
                        meta_val = _pull_field(obj, val[level:])
                    meta_vals[key].append(meta_val)
                records.extend(recs)
    _recursive_extract(data, record_path, {}, level=0)
    result = DataFrame(records)
    if record_prefix is not None:
        result = result.rename(columns=lambda x: f'{record_prefix}{x}')
    for k, v in meta_vals.items():
        if meta_prefix is not None:
            k = meta_prefix + k
        if k in result:
            raise ValueError(f'Conflicting metadata name {k}, need distinguishing prefix ')
        values = np.array(v, dtype=object)
        if values.ndim > 1:
            values = np.empty((len(v),), dtype=object)
            for i, val in enumerate(v):
                values[i] = val
        result[k] = values.repeat(lengths)
    if index is not None:
        result.index = index.repeat(lengths)
    return result
