from __future__ import annotations

import ast
import importlib.util
import keyword
import re
import shutil
import sys
import tempfile
from pathlib import Path


_PLACEHOLDER_SOURCE_MODULES = {"solution", "your_module", "module_name", "src"}
_LANGUAGE_TOKENS = {"python", "java", "go", "cpp", "javascript"}
_ALLOWED_TEST_IMPORT_PREFIXES = {
    "pytest",
    "unittest",
    "mock",
    "typing",
    "collections",
    "pathlib",
    "tempfile",
    "json",
    "math",
    "re",
    "os",
    "sys",
    "datetime",
    "itertools",
    "functools",
    "dataclasses",
    "subprocess",
    "asyncio",
}
_STDLIB_MODULES = set(getattr(sys, "stdlib_module_names", set()))


def compile_check(language: str, generated_test_path: Path) -> tuple[bool, str | None]:
    """阶段1：编译/语法验证。"""
    if language == "python":
        return _compile_python(generated_test_path)
    if language in {"java", "go", "cpp", "javascript"}:
        return False, f"Language toolchain not implemented yet: {language}"
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
    """为 Python 测试执行创建隔离目录，完整模拟包结构以支持相对导入。"""
    workdir = Path(tempfile.mkdtemp(prefix="utbench_eval_"))

    generated_test_source = generated_test_path.read_text(encoding="utf-8")
    if source_path and source_path.exists():
        generated_test_source = rewrite_generated_test_imports(
            generated_test_source=generated_test_source,
            source_path=source_path,
        )

    source_dir = source_path.parent if source_path else None
    source_stem = source_path.stem if source_path else None

    if source_dir and source_stem:
        pkg_root = _find_package_root(source_dir)
        if pkg_root:
            _copy_package_to_workspace(pkg_root, workdir)
        else:
            target_src = workdir / f"{source_stem}.py"
            target_src.write_text(source_path.read_text(encoding="utf-8"), encoding="utf-8")

        alias_modules = python_source_alias_candidates(
            source_path=source_path,
            generated_test_source=generated_test_source,
        )
        for module_name in sorted(alias_modules):
            if module_name == source_stem:
                continue
            if not _is_valid_module_name(module_name):
                continue
            alias_file = workdir / f"{module_name}.py"
            if alias_file.exists():
                continue
            alias_file.write_text(source_path.read_text(encoding="utf-8"), encoding="utf-8")

    target_test = workdir / _normalized_pytest_filename(generated_test_path)
    target_test.write_text(generated_test_source, encoding="utf-8")

    return workdir, target_test


def _find_package_root(source_dir: Path) -> Path | None:
    """向上查找最近的包根目录（包含 __init__.py 或与父级同名目录结构）。"""
    if not source_dir.exists():
        return None
    if (source_dir / "__init__.py").exists():
        return source_dir
    parent = source_dir.parent
    if parent.name == "dataset":
        return None
    if source_dir.name == parent.name:
        return _find_package_root(parent)
    if (parent / "__init__.py").exists():
        return parent
    return None


def _copy_package_to_workspace(pkg_root: Path, workdir: Path) -> None:
    """将完整包结构复制到 workspace，确保相对导入可用。"""
    target_pkg = workdir / pkg_root.name
    target_pkg.mkdir(parents=True, exist_ok=True)
    for src_file in pkg_root.rglob("*.py"):
        rel = src_file.relative_to(pkg_root)
        dst_file = target_pkg / rel
        dst_file.parent.mkdir(parents=True, exist_ok=True)
        dst_file.write_text(src_file.read_text(encoding="utf-8"), encoding="utf-8")


def cleanup_execution_workspace(workdir: Path) -> None:
    """Best-effort cleanup for temporary evaluator workspace."""
    try:
        shutil.rmtree(workdir, ignore_errors=True)
    except Exception:
        pass


def python_source_precheck(source_path: Path | None) -> str | None:
    """预检查 Python 源码是否满足当前评测沙箱的基本可执行条件。"""
    if source_path is None:
        return "missing source file"
    if not source_path.exists() or not source_path.is_file():
        return f"source file not found: {source_path}"

    try:
        source = source_path.read_text(encoding="utf-8")
    except Exception as exc:  # pylint: disable=broad-except
        return f"failed to read source: {exc}"

    if re.search(r"^\s*from\s+\.+", source, flags=re.MULTILINE):
        return "source uses relative imports and is not self-contained"

    missing_deps = _find_missing_external_imports(source_path=source_path, source=source)
    if missing_deps:
        preview = ", ".join(sorted(missing_deps)[:5])
        return f"missing runtime dependencies: {preview}"

    return None


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


def rewrite_generated_test_imports(
    generated_test_source: str,
    source_path: Path,
) -> str:
    """将明显错误的导入改写到本地被测源码模块，提升评测可执行性。"""
    source_module = source_path.stem
    source_exports = _source_exported_names(source_path)
    if not source_exports:
        return generated_test_source

    try:
        tree = ast.parse(generated_test_source)
    except Exception:  # pylint: disable=broad-except
        return generated_test_source

    rewriter = _GeneratedTestImportRewriter(
        source_module=source_module,
        source_exports=source_exports,
    )
    rewritten = rewriter.visit(tree)
    if not rewriter.changed:
        return generated_test_source

    ast.fix_missing_locations(rewritten)
    try:
        return ast.unparse(rewritten)
    except Exception:  # pylint: disable=broad-except
        return generated_test_source


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


def _source_exported_names(source_path: Path) -> set[str]:
    try:
        tree = ast.parse(source_path.read_text(encoding="utf-8"))
    except Exception:  # pylint: disable=broad-except
        return set()

    exported: set[str] = set()
    for node in tree.body:
        if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef, ast.ClassDef)):
            exported.add(node.name)
            continue
        if isinstance(node, ast.Assign):
            for target in node.targets:
                if isinstance(target, ast.Name):
                    exported.add(target.id)
            if _is_dunder_all_assign(node):
                exported.update(_extract_dunder_all(node.value))
            continue
        if isinstance(node, ast.AnnAssign) and isinstance(node.target, ast.Name):
            exported.add(node.target.id)

    return {name for name in exported if _is_valid_module_name(name)}


def _is_dunder_all_assign(node: ast.Assign) -> bool:
    return any(isinstance(target, ast.Name) and target.id == "__all__" for target in node.targets)


def _extract_dunder_all(value: ast.AST) -> set[str]:
    names: set[str] = set()
    if isinstance(value, (ast.List, ast.Tuple, ast.Set)):
        for item in value.elts:
            if isinstance(item, ast.Constant) and isinstance(item.value, str):
                names.add(item.value)
    return names


def _find_missing_external_imports(source_path: Path, source: str) -> set[str]:
    try:
        tree = ast.parse(source)
    except Exception:  # pylint: disable=broad-except
        return set()

    imported: set[str] = set()
    for node in ast.walk(tree):
        if isinstance(node, ast.Import):
            for alias in node.names:
                top = alias.name.split(".", maxsplit=1)[0]
                if top:
                    imported.add(top)
        elif isinstance(node, ast.ImportFrom):
            if node.level and node.level > 0:
                continue
            if node.module:
                top = node.module.split(".", maxsplit=1)[0]
                if top:
                    imported.add(top)

    missing: set[str] = set()
    for module in imported:
        if module == source_path.stem:
            continue
        if _is_stdlib_module(module):
            continue
        if _is_local_module(source_path.parent, module):
            continue
        if importlib.util.find_spec(module) is None:
            missing.add(module)

    return missing


def _is_stdlib_module(module_name: str) -> bool:
    if module_name in _STDLIB_MODULES:
        return True
    return module_name in {"__future__"}


def _is_local_module(base_dir: Path, module_name: str) -> bool:
    file_candidate = base_dir / f"{module_name}.py"
    package_candidate = base_dir / module_name
    if file_candidate.exists() and file_candidate.is_file():
        return True
    return package_candidate.exists() and package_candidate.is_dir()


class _GeneratedTestImportRewriter(ast.NodeTransformer):
    def __init__(self, source_module: str, source_exports: set[str]) -> None:
        self.source_module = source_module
        self.source_exports = source_exports
        self.changed = False

    def visit_ImportFrom(self, node: ast.ImportFrom) -> ast.AST:  # noqa: N802
        module_name = node.module or ""
        imported_names = [alias.name for alias in node.names if alias.name != "*"]
        should_rewrite = False

        if node.level and node.level > 0:
            should_rewrite = True
        elif module_name in _PLACEHOLDER_SOURCE_MODULES:
            should_rewrite = True
        elif module_name and module_name != self.source_module:
            top = module_name.split(".", maxsplit=1)[0]
            if top not in _ALLOWED_TEST_IMPORT_PREFIXES and imported_names:
                if all(name in self.source_exports for name in imported_names):
                    should_rewrite = True

        if should_rewrite:
            self.changed = True
            return ast.copy_location(
                ast.ImportFrom(module=self.source_module, names=node.names, level=0),
                node,
            )

        return node

    def visit_Import(self, node: ast.Import) -> ast.AST:  # noqa: N802
        names: list[ast.alias] = []
        changed = False
        for alias in node.names:
            top = alias.name.split(".", maxsplit=1)[0]
            if top in _PLACEHOLDER_SOURCE_MODULES and top != self.source_module:
                changed = True
                names.append(
                    ast.alias(
                        name=self.source_module,
                        asname=alias.asname or top,
                    )
                )
            else:
                names.append(alias)

        if changed:
            self.changed = True
            return ast.copy_location(ast.Import(names=names), node)

        return node
