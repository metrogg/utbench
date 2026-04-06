from __future__ import annotations

import ast
import importlib.util
import keyword
import re
import tempfile
from pathlib import Path


_PLACEHOLDER_SOURCE_MODULES = {"solution", "your_module", "module_name", "src"}
_LANGUAGE_TOKENS = {"python", "java", "go", "cpp", "javascript"}


def compile_check(language: str, generated_test_path: Path) -> tuple[bool, str | None]:
    """阶段1：编译/语法验证。"""
    if language == "python":
        return _compile_python(generated_test_path)
    if language in {"java", "go", "cpp", "javascript"}:
        # 其他语言先返回占位通过，后续接入真实工具链。
        return True, None
    return False, f"Unsupported language: {language}"


def _compile_python(generated_test_path: Path) -> tuple[bool, str | None]:
    try:
        source = generated_test_path.read_text(encoding="utf-8")
        ast.parse(source)
        return True, None
    except Exception as exc:  # pylint: disable=broad-except
        return False, str(exc)


def prepare_python_execution_workspace(
    generated_test_path: Path,
    source_path: Path | None,
) -> tuple[Path, Path]:
    """为 Python 测试执行创建隔离目录。"""
    workdir = Path(tempfile.mkdtemp(prefix="utbench_eval_"))

    generated_test_source = generated_test_path.read_text(encoding="utf-8")
    target_test = workdir / _normalized_pytest_filename(generated_test_path)
    target_test.write_text(generated_test_source, encoding="utf-8")

    if source_path and source_path.exists():
        source_text = source_path.read_text(encoding="utf-8")
        target_src = workdir / source_path.name
        target_src.write_text(source_text, encoding="utf-8")

        alias_modules = python_source_alias_candidates(
            source_path=source_path,
            generated_test_source=generated_test_source,
        )
        for module_name in sorted(alias_modules):
            if module_name == source_path.stem:
                continue
            if not _is_valid_module_name(module_name):
                continue
            alias_file = workdir / f"{module_name}.py"
            if alias_file.exists():
                continue
            alias_file.write_text(source_text, encoding="utf-8")

    return workdir, target_test


def python_source_alias_candidates(
    source_path: Path,
    generated_test_source: str | None = None,
    generated_test_path: Path | None = None,
) -> set[str]:
    """推断可用于导入被测源码的模块别名。"""
    aliases: set[str] = {source_path.stem}
    aliases.update(_PLACEHOLDER_SOURCE_MODULES)
    aliases.update(_aliases_from_sample_id(source_path.stem))
    aliases.update(_aliases_from_source_symbols(source_path))

    if generated_test_source is None and generated_test_path and generated_test_path.exists():
        generated_test_source = generated_test_path.read_text(encoding="utf-8")

    if generated_test_source:
        aliases.update(_aliases_from_generated_test_imports(generated_test_source))

    return {name for name in aliases if _is_valid_module_name(name)}


def _normalized_pytest_filename(generated_test_path: Path) -> str:
    """将生成文件名转换为 pytest 可稳定导入的测试文件名。"""
    stem = re.sub(r"[^0-9a-zA-Z_]", "_", generated_test_path.stem)
    if not stem:
        stem = "generated"
    if stem[0].isdigit():
        stem = f"_{stem}"
    if not stem.startswith("test_"):
        stem = f"test_{stem}"
    return f"{stem}.py"


def _aliases_from_sample_id(sample_id: str) -> set[str]:
    aliases: set[str] = set()
    parts = sample_id.split("_")
    if len(parts) < 3:
        return aliases

    lang_idx = next((idx for idx, part in enumerate(parts) if part in _LANGUAGE_TOKENS), None)
    if lang_idx is None:
        return aliases

    scenario_parts = parts[lang_idx + 1 :]
    if scenario_parts and scenario_parts[-1].isdigit():
        scenario_parts = scenario_parts[:-1]
    scenario = "_".join(scenario_parts)

    complexity_parts = parts[:lang_idx]
    complexity = "_".join(complexity_parts)

    if scenario:
        aliases.add(scenario)
    if complexity:
        aliases.add(complexity)
    return aliases


def _aliases_from_source_symbols(source_path: Path) -> set[str]:
    aliases: set[str] = set()
    try:
        tree = ast.parse(source_path.read_text(encoding="utf-8"))
    except Exception:  # pylint: disable=broad-except
        return aliases

    for node in tree.body:
        if isinstance(node, (ast.FunctionDef, ast.ClassDef)):
            aliases.add(node.name)
    return aliases


def _aliases_from_generated_test_imports(source: str) -> set[str]:
    aliases: set[str] = set()
    try:
        tree = ast.parse(source)
    except Exception:  # pylint: disable=broad-except
        return aliases

    for node in ast.walk(tree):
        if isinstance(node, ast.Import):
            for alias in node.names:
                top = alias.name.split(".", maxsplit=1)[0]
                if _should_use_alias(top):
                    aliases.add(top)
        elif isinstance(node, ast.ImportFrom):
            if node.level and node.level > 0:
                continue
            if node.module:
                module_name = node.module.split(".", maxsplit=1)[0]
                if _should_use_alias(module_name):
                    aliases.add(module_name)
    return aliases


def _should_use_alias(module_name: str) -> bool:
    if not _is_valid_module_name(module_name):
        return False
    if module_name in _PLACEHOLDER_SOURCE_MODULES:
        return True
    return importlib.util.find_spec(module_name) is None


def _is_valid_module_name(name: str) -> bool:
    if not name:
        return False
    if keyword.iskeyword(name):
        return False
    return re.match(r"^[a-zA-Z_][a-zA-Z0-9_]*$", name) is not None
