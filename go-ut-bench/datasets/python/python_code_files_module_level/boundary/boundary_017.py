from __future__ import annotations
from collections.abc import Callable, Generator, Iterable, Mapping, MutableMapping
from contextlib import contextmanager
from typing import Any, Literal, overload
from . import helpers, presets
from .common import normalize_url, utils
from .parser_block import ParserBlock
from .parser_core import ParserCore
from .parser_inline import ParserInline
from .renderer import RendererHTML, RendererProtocol
from .rules_core.state_core import StateCore
from .token import Token
from .utils import EnvType, OptionsDict, OptionsType, PresetType
try:
    import linkify_it
except ModuleNotFoundError:
    linkify_it = None
_PRESETS: dict[str, PresetType] = {'default': presets.default.make(), 'js-default': presets.js_default.make(), 'zero': presets.zero.make(), 'commonmark': presets.commonmark.make(), 'gfm-like': presets.gfm_like.make()}

class MarkdownIt:

    def __init__(self, config: str | PresetType='commonmark', options_update: Mapping[str, Any] | None=None, *, renderer_cls: Callable[[MarkdownIt], RendererProtocol]=RendererHTML):
        self.utils = utils
        self.helpers = helpers
        self.inline = ParserInline()
        self.block = ParserBlock()
        self.core = ParserCore()
        self.renderer = renderer_cls(self)
        self.linkify = linkify_it.LinkifyIt() if linkify_it else None
        if options_update and (not isinstance(options_update, Mapping)):
            raise TypeError(f'options_update should be a mapping: {options_update}\n(Perhaps you intended this to be the renderer_cls?)')
        self.configure(config, options_update=options_update)

    def __repr__(self) -> str:
        return f'{self.__class__.__module__}.{self.__class__.__name__}()'

    @overload
    def __getitem__(self, name: Literal['inline']) -> ParserInline:
        ...

    @overload
    def __getitem__(self, name: Literal['block']) -> ParserBlock:
        ...

    @overload
    def __getitem__(self, name: Literal['core']) -> ParserCore:
        ...

    @overload
    def __getitem__(self, name: Literal['renderer']) -> RendererProtocol:
        ...

    @overload
    def __getitem__(self, name: str) -> Any:
        ...

    def __getitem__(self, name: str) -> Any:
        return {'inline': self.inline, 'block': self.block, 'core': self.core, 'renderer': self.renderer}[name]

    def set(self, options: OptionsType) -> None:
        self.options = OptionsDict(options)

    def configure(self, presets: str | PresetType, options_update: Mapping[str, Any] | None=None) -> MarkdownIt:
        if isinstance(presets, str):
            if presets not in _PRESETS:
                raise KeyError(f"Wrong `markdown-it` preset '{presets}', check name")
            config = _PRESETS[presets]
        else:
            config = presets
        if not config:
            raise ValueError("Wrong `markdown-it` config, can't be empty")
        options = config.get('options', {}) or {}
        if options_update:
            options = {**options, **options_update}
        self.set(options)
        if 'components' in config:
            for name, component in config['components'].items():
                rules = component.get('rules', None)
                if rules:
                    self[name].ruler.enableOnly(rules)
                rules2 = component.get('rules2', None)
                if rules2:
                    self[name].ruler2.enableOnly(rules2)
        return self

    def get_all_rules(self) -> dict[str, list[str]]:
        rules = {chain: self[chain].ruler.get_all_rules() for chain in ['core', 'block', 'inline']}
        rules['inline2'] = self.inline.ruler2.get_all_rules()
        return rules

    def get_active_rules(self) -> dict[str, list[str]]:
        rules = {chain: self[chain].ruler.get_active_rules() for chain in ['core', 'block', 'inline']}
        rules['inline2'] = self.inline.ruler2.get_active_rules()
        return rules

    def enable(self, names: str | Iterable[str], ignoreInvalid: bool=False) -> MarkdownIt:
        result = []
        if isinstance(names, str):
            names = [names]
        for chain in ['core', 'block', 'inline']:
            result.extend(self[chain].ruler.enable(names, True))
        result.extend(self.inline.ruler2.enable(names, True))
        missed = [name for name in names if name not in result]
        if missed and (not ignoreInvalid):
            raise ValueError(f'MarkdownIt. Failed to enable unknown rule(s): {missed}')
        return self

    def disable(self, names: str | Iterable[str], ignoreInvalid: bool=False) -> MarkdownIt:
        result = []
        if isinstance(names, str):
            names = [names]
        for chain in ['core', 'block', 'inline']:
            result.extend(self[chain].ruler.disable(names, True))
        result.extend(self.inline.ruler2.disable(names, True))
        missed = [name for name in names if name not in result]
        if missed and (not ignoreInvalid):
            raise ValueError(f'MarkdownIt. Failed to disable unknown rule(s): {missed}')
        return self

    @contextmanager
    def reset_rules(self) -> Generator[None, None, None]:
        chain_rules = self.get_active_rules()
        yield
        for chain, rules in chain_rules.items():
            if chain != 'inline2':
                self[chain].ruler.enableOnly(rules)
        self.inline.ruler2.enableOnly(chain_rules['inline2'])

    def add_render_rule(self, name: str, function: Callable[..., Any], fmt: str='html') -> None:
        if self.renderer.__output__ == fmt:
            self.renderer.rules[name] = function.__get__(self.renderer)

    def use(self, plugin: Callable[..., None], *params: Any, **options: Any) -> MarkdownIt:
        plugin(self, *params, **options)
        return self

    def parse(self, src: str, env: EnvType | None=None) -> list[Token]:
        env = {} if env is None else env
        if not isinstance(env, MutableMapping):
            raise TypeError(f'Input data should be a MutableMapping, not {type(env)}')
        if not isinstance(src, str):
            raise TypeError(f'Input data should be a string, not {type(src)}')
        state = StateCore(src, self, env)
        self.core.process(state)
        return state.tokens

    def render(self, src: str, env: EnvType | None=None) -> Any:
        env = {} if env is None else env
        return self.renderer.render(self.parse(src, env), self.options, env)

    def parseInline(self, src: str, env: EnvType | None=None) -> list[Token]:
        env = {} if env is None else env
        if not isinstance(env, MutableMapping):
            raise TypeError(f'Input data should be an MutableMapping, not {type(env)}')
        if not isinstance(src, str):
            raise TypeError(f'Input data should be a string, not {type(src)}')
        state = StateCore(src, self, env)
        state.inlineMode = True
        self.core.process(state)
        return state.tokens

    def renderInline(self, src: str, env: EnvType | None=None) -> Any:
        env = {} if env is None else env
        return self.renderer.render(self.parseInline(src, env), self.options, env)

    def validateLink(self, url: str) -> bool:
        return normalize_url.validateLink(url)

    def normalizeLink(self, url: str) -> str:
        return normalize_url.normalizeLink(url)

    def normalizeLinkText(self, link: str) -> str:
        return normalize_url.normalizeLinkText(link)
