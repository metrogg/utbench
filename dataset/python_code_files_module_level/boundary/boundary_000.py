import collections.abc
import inspect
import sys
import types
from dataclasses import _MISSING_TYPE, MISSING, Field, field, fields, make_dataclass
from functools import lru_cache, wraps
from typing import Annotated, Any, Callable, ForwardRef, Literal, Optional, Type, TypeVar, Union, get_args, get_origin, overload
try:
    from typing import NotRequired, Required
except ImportError:
    try:
        from typing_extensions import NotRequired, Required
    except ImportError:
        Required = type('Required', (), {})
        NotRequired = type('NotRequired', (), {})
from .errors import StrictDataclassClassValidationError, StrictDataclassDefinitionError, StrictDataclassFieldValidationError
Validator_T = Callable[[Any], None]
T = TypeVar('T')
TypedDictType = TypeVar('TypedDictType', bound=dict[str, Any])
_TYPED_DICT_DEFAULT_VALUE = object()

@overload
def strict(cls: Type[T]) -> Type[T]:
    ...

@overload
def strict(*, accept_kwargs: bool=False) -> Callable[[Type[T]], Type[T]]:
    ...

def strict(cls: Optional[Type[T]]=None, *, accept_kwargs: bool=False) -> Union[Type[T], Callable[[Type[T]], Type[T]]]:

    def wrap(cls: Type[T]) -> Type[T]:
        if not hasattr(cls, '__dataclass_fields__'):
            raise StrictDataclassDefinitionError(f"Class '{cls.__name__}' must be a dataclass before applying @strict.")
        field_validators: dict[str, list[Validator_T]] = {}
        for f in fields(cls):
            validators = []
            validators.append(_create_type_validator(f))
            custom_validator = f.metadata.get('validator')
            if custom_validator is not None:
                if not isinstance(custom_validator, list):
                    custom_validator = [custom_validator]
                for validator in custom_validator:
                    if not _is_validator(validator):
                        raise StrictDataclassDefinitionError(f"Invalid validator for field '{f.name}': {validator}. Must be a callable taking a single argument.")
                validators.extend(custom_validator)
            field_validators[f.name] = validators
        cls.__validators__ = field_validators
        original_setattr = cls.__setattr__

        def __strict_setattr__(self: Any, name: str, value: Any) -> None:
            for validator in self.__validators__.get(name, []):
                try:
                    validator(value)
                except (ValueError, TypeError) as e:
                    raise StrictDataclassFieldValidationError(field=name, cause=e) from e
            original_setattr(self, name, value)
        cls.__setattr__ = __strict_setattr__
        if accept_kwargs:
            original_init = cls.__init__

            @wraps(original_init)
            def __init__(self, *args, **kwargs: Any) -> None:
                dataclass_fields = {f.name for f in fields(cls)}
                standard_kwargs = {k: v for k, v in kwargs.items() if k in dataclass_fields}
                if len(args) > 0:
                    raise ValueError(f'When `accept_kwargs=True`, {cls.__name__} accepts only keyword arguments, but found `{len(args)}` positional args.')
                for f in fields(cls):
                    if f.name in standard_kwargs:
                        setattr(self, f.name, standard_kwargs[f.name])
                    elif f.default is not MISSING:
                        setattr(self, f.name, f.default)
                    elif f.default_factory is not MISSING:
                        setattr(self, f.name, f.default_factory())
                    else:
                        raise TypeError(f"Missing required field - '{f.name}'")
                additional_kwargs = {}
                for name, value in kwargs.items():
                    if name not in dataclass_fields:
                        additional_kwargs[name] = value
                self.__post_init__(**additional_kwargs)
            cls.__init__ = __init__
            if not hasattr(cls, '__post_init__'):

                def __post_init__(self, **kwargs: Any) -> None:
                    for name, value in kwargs.items():
                        setattr(self, name, value)
                cls.__post_init__ = __post_init__
            original_repr = cls.__repr__

            @wraps(original_repr)
            def __repr__(self) -> str:
                standard_repr = original_repr(self)
                additional_kwargs = [f'*{k}={v!r}' for k, v in self.__dict__.items() if k not in cls.__dataclass_fields__]
                additional_repr = ', '.join(additional_kwargs)
                return f'{standard_repr[:-1]}, {additional_repr})' if additional_kwargs else standard_repr
            if cls.__dataclass_params__.repr is True:
                cls.__repr__ = __repr__
        class_validators = []
        for name in dir(cls):
            if not name.startswith('validate_'):
                continue
            method = getattr(cls, name)
            if not callable(method):
                continue
            if len(inspect.signature(method).parameters) != 1:
                raise StrictDataclassDefinitionError(f"Class '{cls.__name__}' has a class validator '{name}' that takes more than one argument. Class validators must take only 'self' as an argument. Methods starting with 'validate_' are considered to be class validators.")
            class_validators.append(method)
        cls.__class_validators__ = class_validators

        def validate(self: T) -> None:
            for validator in cls.__class_validators__:
                try:
                    validator(self)
                except (ValueError, TypeError) as e:
                    raise StrictDataclassClassValidationError(validator=validator.__name__, cause=e) from e
        validate.__is_defined_by_strict_decorator__ = True
        if hasattr(cls, 'validate'):
            if not getattr(cls.validate, '__is_defined_by_strict_decorator__', False):
                raise StrictDataclassDefinitionError(f"Class '{cls.__name__}' already implements a method called 'validate'. This method name is reserved when using the @strict decorator on a dataclass. If you want to keep your own method, please rename it.")
        cls.validate = validate
        initial_init = cls.__init__

        @wraps(initial_init)
        def init_with_validate(self, *args, **kwargs) -> None:
            initial_init(self, *args, **kwargs)
            cls.validate(self)
        setattr(cls, '__init__', init_with_validate)
        return cls
    return wrap(cls) if cls is not None else wrap

def validate_typed_dict(schema: type[TypedDictType], data: dict) -> None:
    strict_cls = _build_strict_cls_from_typed_dict(schema)
    strict_cls(**data)

@lru_cache
def _build_strict_cls_from_typed_dict(schema: type[TypedDictType]) -> Type:
    type_hints = _get_typed_dict_annotations(schema)
    if not getattr(schema, '__total__', True):
        for key, value in type_hints.items():
            origin = get_origin(value)
            if origin is Annotated:
                base, *meta = get_args(value)
                if not _is_required_or_notrequired(base):
                    base = NotRequired[base]
                type_hints[key] = Annotated[tuple([base] + list(meta))]
            elif not _is_required_or_notrequired(value):
                type_hints[key] = NotRequired[value]
    fields = []
    for key, value in type_hints.items():
        if get_origin(value) is Annotated:
            base, *meta = get_args(value)
            fields.append((key, base, field(default=_TYPED_DICT_DEFAULT_VALUE, metadata={'validator': meta[0]})))
        else:
            fields.append((key, value, field(default=_TYPED_DICT_DEFAULT_VALUE)))
    return strict(make_dataclass(schema.__name__, fields))

def _get_typed_dict_annotations(schema: type[TypedDictType]) -> dict[str, Any]:
    try:
        import annotationlib
        return annotationlib.get_annotations(schema)
    except ImportError:
        return {name: value if value is not None else type(None) for name, value in schema.__dict__.get('__annotations__', {}).items()}

def validated_field(validator: Union[list[Validator_T], Validator_T], default: Union[Any, _MISSING_TYPE]=MISSING, default_factory: Union[Callable[[], Any], _MISSING_TYPE]=MISSING, init: bool=True, repr: bool=True, hash: Optional[bool]=None, compare: bool=True, metadata: Optional[dict]=None, **kwargs: Any) -> Any:
    if not isinstance(validator, list):
        validator = [validator]
    if metadata is None:
        metadata = {}
    metadata['validator'] = validator
    return field(default=default, default_factory=default_factory, init=init, repr=repr, hash=hash, compare=compare, metadata=metadata, **kwargs)

def as_validated_field(validator: Validator_T):

    def _inner(default: Union[Any, _MISSING_TYPE]=MISSING, default_factory: Union[Callable[[], Any], _MISSING_TYPE]=MISSING, init: bool=True, repr: bool=True, hash: Optional[bool]=None, compare: bool=True, metadata: Optional[dict]=None, **kwargs: Any):
        return validated_field(validator, default=default, default_factory=default_factory, init=init, repr=repr, hash=hash, compare=compare, metadata=metadata, **kwargs)
    return _inner

def type_validator(name: str, value: Any, expected_type: Any) -> None:
    origin = get_origin(expected_type)
    args = get_args(expected_type)
    if expected_type is Any:
        return
    elif (validator := _BASIC_TYPE_VALIDATORS.get(origin)):
        validator(name, value, args)
    elif isinstance(expected_type, type):
        _validate_simple_type(name, value, expected_type)
    elif isinstance(expected_type, ForwardRef) or isinstance(expected_type, str):
        return
    elif origin is Required:
        if value is _TYPED_DICT_DEFAULT_VALUE:
            raise TypeError(f"Field '{name}' is required but missing.")
        type_validator(name, value, args[0])
    elif origin is NotRequired:
        if value is _TYPED_DICT_DEFAULT_VALUE:
            return
        type_validator(name, value, args[0])
    else:
        raise TypeError(f"Unsupported type for field '{name}': {expected_type}")

def _validate_union(name: str, value: Any, args: tuple[Any, ...]) -> None:
    errors = []
    for t in args:
        try:
            type_validator(name, value, t)
            return
        except TypeError as e:
            errors.append(str(e))
    raise TypeError(f"Field '{name}' with value {repr(value)} doesn't match any type in {args}. Errors: {'; '.join(errors)}")

def _validate_literal(name: str, value: Any, args: tuple[Any, ...]) -> None:
    if value not in args:
        raise TypeError(f"Field '{name}' expected one of {args}, got {value}")

def _validate_list(name: str, value: Any, args: tuple[Any, ...]) -> None:
    if not isinstance(value, list):
        raise TypeError(f"Field '{name}' expected a list, got {type(value).__name__}")
    item_type = args[0]
    for i, item in enumerate(value):
        try:
            type_validator(f'{name}[{i}]', item, item_type)
        except TypeError as e:
            raise TypeError(f"Invalid item at index {i} in list '{name}'") from e

def _validate_dict(name: str, value: Any, args: tuple[Any, ...]) -> None:
    if not isinstance(value, dict):
        raise TypeError(f"Field '{name}' expected a dict, got {type(value).__name__}")
    key_type, value_type = args
    for k, v in value.items():
        try:
            type_validator(f'{name}.key', k, key_type)
            type_validator(f'{name}[{k!r}]', v, value_type)
        except TypeError as e:
            raise TypeError(f"Invalid key or value in dict '{name}'") from e

def _validate_tuple(name: str, value: Any, args: tuple[Any, ...]) -> None:
    if not isinstance(value, tuple):
        raise TypeError(f"Field '{name}' expected a tuple, got {type(value).__name__}")
    if len(args) == 2 and args[1] is Ellipsis:
        for i, item in enumerate(value):
            try:
                type_validator(f'{name}[{i}]', item, args[0])
            except TypeError as e:
                raise TypeError(f"Invalid item at index {i} in tuple '{name}'") from e
    elif len(args) != len(value):
        raise TypeError(f"Field '{name}' expected a tuple of length {len(args)}, got {len(value)}")
    else:
        for i, (item, expected) in enumerate(zip(value, args)):
            try:
                type_validator(f'{name}[{i}]', item, expected)
            except TypeError as e:
                raise TypeError(f"Invalid item at index {i} in tuple '{name}'") from e

def _validate_set(name: str, value: Any, args: tuple[Any, ...]) -> None:
    if not isinstance(value, set):
        raise TypeError(f"Field '{name}' expected a set, got {type(value).__name__}")
    item_type = args[0]
    for i, item in enumerate(value):
        try:
            type_validator(f'{name} item', item, item_type)
        except TypeError as e:
            raise TypeError(f"Invalid item in set '{name}'") from e

def _validate_sequence(name: str, value: Any, args: tuple[Any, ...]) -> None:
    if not isinstance(value, collections.abc.Sequence):
        raise TypeError(f"Field '{name}' expected a Sequence, got {type(value).__name__}")
    if not args:
        return
    item_type = args[0]
    for i, item in enumerate(value):
        try:
            type_validator(f'{name}[{i}]', item, item_type)
        except TypeError as e:
            raise TypeError(f"Invalid item at index {i} in sequence '{name}'") from e

def _validate_simple_type(name: str, value: Any, expected_type: type) -> None:
    if not isinstance(value, expected_type):
        raise TypeError(f"Field '{name}' expected {expected_type.__name__}, got {type(value).__name__} (value: {repr(value)})")

def _create_type_validator(field: Field) -> Validator_T:

    def validator(value: Any) -> None:
        type_validator(field.name, value, field.type)
    return validator

def _is_validator(validator: Any) -> bool:
    if not callable(validator):
        return False
    signature = inspect.signature(validator)
    parameters = list(signature.parameters.values())
    if len(parameters) == 0:
        return False
    if parameters[0].kind not in (inspect.Parameter.POSITIONAL_OR_KEYWORD, inspect.Parameter.POSITIONAL_ONLY, inspect.Parameter.VAR_POSITIONAL):
        return False
    for parameter in parameters[1:]:
        if parameter.default == inspect.Parameter.empty:
            return False
    return True

def _is_required_or_notrequired(type_hint: Any) -> bool:
    return type_hint in (Required, NotRequired) or get_origin(type_hint) in (Required, NotRequired)
_BASIC_TYPE_VALIDATORS: dict[Any, Callable[[str, Any, tuple[Any, ...]], None]] = {Union: _validate_union, Literal: _validate_literal, list: _validate_list, dict: _validate_dict, tuple: _validate_tuple, set: _validate_set, collections.abc.Sequence: _validate_sequence}
if sys.version_info >= (3, 10):
    _BASIC_TYPE_VALIDATORS[types.UnionType] = _validate_union
__all__ = ['strict', 'validate_typed_dict', 'validated_field', 'Validator_T', 'StrictDataclassClassValidationError', 'StrictDataclassDefinitionError', 'StrictDataclassFieldValidationError']
