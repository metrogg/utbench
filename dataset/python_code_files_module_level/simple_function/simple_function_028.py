from __future__ import annotations
from collections.abc import Generator, Sequence
import textwrap
from typing import Any, NamedTuple, TypeVar, overload
from .token import Token

class _NesterTokens(NamedTuple):
    opening: Token
    closing: Token
_NodeType = TypeVar('_NodeType', bound='SyntaxTreeNode')

class SyntaxTreeNode:

    def __init__(self, tokens: Sequence[Token]=(), *, create_root: bool=True) -> None:
        self.token: Token | None = None
        self.nester_tokens: _NesterTokens | None = None
        self._parent: Any = None
        self._children: list[Any] = []
        if create_root:
            self._set_children_from_tokens(tokens)
            return
        if not tokens:
            raise ValueError('Can only create root from empty token sequence. Set `create_root=True`.')
        elif len(tokens) == 1:
            inline_token = tokens[0]
            if inline_token.nesting:
                raise ValueError('Unequal nesting level at the start and end of token stream.')
            self.token = inline_token
            if inline_token.children:
                self._set_children_from_tokens(inline_token.children)
        else:
            self.nester_tokens = _NesterTokens(tokens[0], tokens[-1])
            self._set_children_from_tokens(tokens[1:-1])

    def __repr__(self) -> str:
        return f'{type(self).__name__}({self.type})'

    @overload
    def __getitem__(self: _NodeType, item: int) -> _NodeType:
        ...

    @overload
    def __getitem__(self: _NodeType, item: slice) -> list[_NodeType]:
        ...

    def __getitem__(self: _NodeType, item: int | slice) -> _NodeType | list[_NodeType]:
        return self.children[item]

    def to_tokens(self: _NodeType) -> list[Token]:

        def recursive_collect_tokens(node: _NodeType, token_list: list[Token]) -> None:
            if node.type == 'root':
                for child in node.children:
                    recursive_collect_tokens(child, token_list)
            elif node.token:
                token_list.append(node.token)
            else:
                assert node.nester_tokens
                token_list.append(node.nester_tokens.opening)
                for child in node.children:
                    recursive_collect_tokens(child, token_list)
                token_list.append(node.nester_tokens.closing)
        tokens: list[Token] = []
        recursive_collect_tokens(self, tokens)
        return tokens

    @property
    def children(self: _NodeType) -> list[_NodeType]:
        return self._children

    @children.setter
    def children(self: _NodeType, value: list[_NodeType]) -> None:
        self._children = value

    @property
    def parent(self: _NodeType) -> _NodeType | None:
        return self._parent

    @parent.setter
    def parent(self: _NodeType, value: _NodeType | None) -> None:
        self._parent = value

    @property
    def is_root(self) -> bool:
        return not (self.token or self.nester_tokens)

    @property
    def is_nested(self) -> bool:
        return bool(self.nester_tokens)

    @property
    def siblings(self: _NodeType) -> Sequence[_NodeType]:
        if not self.parent:
            return [self]
        return self.parent.children

    @property
    def type(self) -> str:
        if self.is_root:
            return 'root'
        if self.token:
            return self.token.type
        assert self.nester_tokens
        return self.nester_tokens.opening.type.removesuffix('_open')

    @property
    def next_sibling(self: _NodeType) -> _NodeType | None:
        self_index = self.siblings.index(self)
        if self_index + 1 < len(self.siblings):
            return self.siblings[self_index + 1]
        return None

    @property
    def previous_sibling(self: _NodeType) -> _NodeType | None:
        self_index = self.siblings.index(self)
        if self_index - 1 >= 0:
            return self.siblings[self_index - 1]
        return None

    def _add_child(self, tokens: Sequence[Token]) -> None:
        child = type(self)(tokens, create_root=False)
        child.parent = self
        self.children.append(child)

    def _set_children_from_tokens(self, tokens: Sequence[Token]) -> None:
        reversed_tokens = list(reversed(tokens))
        while reversed_tokens:
            token = reversed_tokens.pop()
            if not token.nesting:
                self._add_child([token])
                continue
            if token.nesting != 1:
                raise ValueError('Invalid token nesting')
            nested_tokens = [token]
            nesting = 1
            while reversed_tokens and nesting:
                token = reversed_tokens.pop()
                nested_tokens.append(token)
                nesting += token.nesting
            if nesting:
                raise ValueError(f'unclosed tokens starting {nested_tokens[0]}')
            self._add_child(nested_tokens)

    def pretty(self, *, indent: int=2, show_text: bool=False, _current: int=0) -> str:
        prefix = ' ' * _current
        text = prefix + f'<{self.type}'
        if not self.is_root and self.attrs:
            text += ' ' + ' '.join((f'{k}={v!r}' for k, v in self.attrs.items()))
        text += '>'
        if show_text and (not self.is_root) and (self.type in ('text', 'text_special')) and self.content:
            text += '\n' + textwrap.indent(self.content, prefix + ' ' * indent)
        for child in self.children:
            text += '\n' + child.pretty(indent=indent, show_text=show_text, _current=_current + indent)
        return text

    def walk(self: _NodeType, *, include_self: bool=True) -> Generator[_NodeType, None, None]:
        if include_self:
            yield self
        for child in self.children:
            yield from child.walk(include_self=True)

    def _attribute_token(self) -> Token:
        if self.token:
            return self.token
        if self.nester_tokens:
            return self.nester_tokens.opening
        raise AttributeError('Root node does not have the accessed attribute')

    @property
    def tag(self) -> str:
        return self._attribute_token().tag

    @property
    def attrs(self) -> dict[str, str | int | float]:
        return self._attribute_token().attrs

    def attrGet(self, name: str) -> None | str | int | float:
        return self._attribute_token().attrGet(name)

    @property
    def map(self) -> tuple[int, int] | None:
        map_ = self._attribute_token().map
        if map_:
            return tuple(map_)
        return None

    @property
    def level(self) -> int:
        return self._attribute_token().level

    @property
    def content(self) -> str:
        return self._attribute_token().content

    @property
    def markup(self) -> str:
        return self._attribute_token().markup

    @property
    def info(self) -> str:
        return self._attribute_token().info

    @property
    def meta(self) -> dict[Any, Any]:
        return self._attribute_token().meta

    @property
    def block(self) -> bool:
        return self._attribute_token().block

    @property
    def hidden(self) -> bool:
        return self._attribute_token().hidden
