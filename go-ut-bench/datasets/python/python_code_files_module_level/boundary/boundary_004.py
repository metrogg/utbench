from __future__ import annotations
from collections.abc import Iterable
from collections.abc import Sequence
from dataclasses import dataclass
from dataclasses import KW_ONLY
import os
from pathlib import Path
import sys
from typing import Literal
from typing import TypeAlias
import iniconfig
from .exceptions import UsageError
from _pytest.outcomes import fail
from _pytest.pathlib import absolutepath
from _pytest.pathlib import commonpath
from _pytest.pathlib import safe_exists

@dataclass(frozen=True)
class ConfigValue:
    value: object
    _: KW_ONLY
    origin: Literal['file', 'override']
    mode: Literal['ini', 'toml']
ConfigDict: TypeAlias = dict[str, ConfigValue]

def _parse_ini_config(path: Path) -> iniconfig.IniConfig:
    try:
        return iniconfig.IniConfig(str(path))
    except iniconfig.ParseError as exc:
        raise UsageError(str(exc)) from exc

def load_config_dict_from_file(filepath: Path) -> ConfigDict | None:
    if filepath.suffix == '.ini':
        iniconfig = _parse_ini_config(filepath)
        if 'pytest' in iniconfig:
            return {k: ConfigValue(v, origin='file', mode='ini') for k, v in iniconfig['pytest'].items()}
        elif filepath.name in {'pytest.ini', '.pytest.ini'}:
            return {}
    elif filepath.suffix == '.cfg':
        iniconfig = _parse_ini_config(filepath)
        if 'tool:pytest' in iniconfig.sections:
            return {k: ConfigValue(v, origin='file', mode='ini') for k, v in iniconfig['tool:pytest'].items()}
        elif 'pytest' in iniconfig.sections:
            fail(CFG_PYTEST_SECTION.format(filename='setup.cfg'), pytrace=False)
    elif filepath.suffix == '.toml':
        if sys.version_info >= (3, 11):
            import tomllib
        else:
            import tomli as tomllib
        toml_text = filepath.read_text(encoding='utf-8')
        try:
            config = tomllib.loads(toml_text)
        except tomllib.TOMLDecodeError as exc:
            raise UsageError(f'{filepath}: {exc}') from exc
        if filepath.name in ('pytest.toml', '.pytest.toml'):
            pytest_config = config.get('pytest', {})
            if pytest_config:
                return {k: ConfigValue(v, origin='file', mode='toml') for k, v in pytest_config.items()}
            return {}
        else:
            tool_pytest = config.get('tool', {}).get('pytest', {})
            toml_config = {k: v for k, v in tool_pytest.items() if k != 'ini_options'}
            ini_config = tool_pytest.get('ini_options', None)
            if toml_config and ini_config:
                raise UsageError(f'{filepath}: Cannot use both [tool.pytest] (native TOML types) and [tool.pytest.ini_options] (string-based INI format) simultaneously. Please use [tool.pytest] with native TOML types (recommended) or [tool.pytest.ini_options] for backwards compatibility.')
            if toml_config:
                return {k: ConfigValue(v, origin='file', mode='toml') for k, v in toml_config.items()}
            elif ini_config is not None:

                def make_scalar(v: object) -> str | list[str]:
                    return v if isinstance(v, list) else str(v)
                return {k: ConfigValue(make_scalar(v), origin='file', mode='ini') for k, v in ini_config.items()}
    return None

def locate_config(invocation_dir: Path, args: Iterable[Path]) -> tuple[Path | None, Path | None, ConfigDict, Sequence[str]]:
    config_names = ['pytest.toml', '.pytest.toml', 'pytest.ini', '.pytest.ini', 'pyproject.toml', 'tox.ini', 'setup.cfg']
    args = [x for x in args if not str(x).startswith('-')]
    if not args:
        args = [invocation_dir]
    found_pyproject_toml: Path | None = None
    ignored_config_files: list[str] = []
    for arg in args:
        argpath = absolutepath(arg)
        for base in (argpath, *argpath.parents):
            for config_name in config_names:
                p = base / config_name
                if p.is_file():
                    if p.name == 'pyproject.toml' and found_pyproject_toml is None:
                        found_pyproject_toml = p
                    ini_config = load_config_dict_from_file(p)
                    if ini_config is not None:
                        index = config_names.index(config_name)
                        for remainder in config_names[index + 1:]:
                            p2 = base / remainder
                            if p2.is_file() and load_config_dict_from_file(p2) is not None:
                                ignored_config_files.append(remainder)
                        return (base, p, ini_config, ignored_config_files)
    if found_pyproject_toml is not None:
        return (found_pyproject_toml.parent, found_pyproject_toml, {}, [])
    return (None, None, {}, [])

def get_common_ancestor(invocation_dir: Path, paths: Iterable[Path]) -> Path:
    common_ancestor: Path | None = None
    for path in paths:
        if not path.exists():
            continue
        if common_ancestor is None:
            common_ancestor = path
        elif common_ancestor in path.parents or path == common_ancestor:
            continue
        elif path in common_ancestor.parents:
            common_ancestor = path
        else:
            shared = commonpath(path, common_ancestor)
            if shared is not None:
                common_ancestor = shared
    if common_ancestor is None:
        common_ancestor = invocation_dir
    elif common_ancestor.is_file():
        common_ancestor = common_ancestor.parent
    return common_ancestor

def get_dirs_from_args(args: Iterable[str]) -> list[Path]:

    def is_option(x: str) -> bool:
        return x.startswith('-')

    def get_file_part_from_node_id(x: str) -> str:
        return x.split('::')[0]

    def get_dir_from_path(path: Path) -> Path:
        if path.is_dir():
            return path
        return path.parent
    possible_paths = (absolutepath(get_file_part_from_node_id(arg)) for arg in args if not is_option(arg))
    return [get_dir_from_path(path) for path in possible_paths if safe_exists(path)]

def parse_override_ini(override_ini: Sequence[str] | None) -> ConfigDict:
    overrides = {}
    for ini_config in override_ini or ():
        try:
            key, user_ini_value = ini_config.split('=', 1)
        except ValueError as e:
            raise UsageError(f'-o/--override-ini expects option=value style (got: {ini_config!r}).') from e
        else:
            overrides[key] = ConfigValue(user_ini_value, origin='override', mode='ini')
    return overrides
CFG_PYTEST_SECTION = '[pytest] section in {filename} files is no longer supported, change to [tool:pytest] instead.'

def determine_setup(*, inifile: str | None, override_ini: Sequence[str] | None, args: Sequence[str], rootdir_cmd_arg: str | None, invocation_dir: Path) -> tuple[Path, Path | None, ConfigDict, Sequence[str]]:
    rootdir = None
    dirs = get_dirs_from_args(args)
    ignored_config_files: Sequence[str] = []
    if inifile:
        inipath_ = absolutepath(inifile)
        inipath: Path | None = inipath_
        inicfg = load_config_dict_from_file(inipath_) or {}
        if rootdir_cmd_arg is None:
            rootdir = inipath_.parent
    else:
        ancestor = get_common_ancestor(invocation_dir, dirs)
        rootdir, inipath, inicfg, ignored_config_files = locate_config(invocation_dir, [ancestor])
        if rootdir is None and rootdir_cmd_arg is None:
            for possible_rootdir in (ancestor, *ancestor.parents):
                if (possible_rootdir / 'setup.py').is_file():
                    rootdir = possible_rootdir
                    break
            else:
                if dirs != [ancestor]:
                    rootdir, inipath, inicfg, _ = locate_config(invocation_dir, dirs)
                if rootdir is None:
                    rootdir = get_common_ancestor(invocation_dir, [invocation_dir, ancestor])
                    if is_fs_root(rootdir):
                        rootdir = ancestor
    if rootdir_cmd_arg:
        rootdir = absolutepath(os.path.expandvars(rootdir_cmd_arg))
        if not rootdir.is_dir():
            raise UsageError(f"Directory '{rootdir}' not found. Check your '--rootdir' option.")
    ini_overrides = parse_override_ini(override_ini)
    inicfg.update(ini_overrides)
    assert rootdir is not None
    return (rootdir, inipath, inicfg, ignored_config_files)

def is_fs_root(p: Path) -> bool:
    return os.path.splitdrive(str(p))[1] == os.sep
