from __future__ import annotations

from pathlib import Path

from benchmark.evaluator.compiler import (
    compile_check,
    python_source_alias_candidates,
    python_source_precheck,
    rewrite_generated_test_imports,
)


def test_python_source_alias_candidates_include_placeholders_and_symbols(tmp_path: Path) -> None:
    src = tmp_path / "simple_function_python_simple_function_0.py"
    src.write_text(
        """
def has_close_elements(numbers, threshold):
    return False
""".strip()
        + "\n",
        encoding="utf-8",
    )

    generated_test = tmp_path / "generated.test.py"
    generated_test.write_text(
        """
from solution import has_close_elements
from your_module import helper
""".strip()
        + "\n",
        encoding="utf-8",
    )

    aliases = python_source_alias_candidates(
        source_path=src,
        generated_test_path=generated_test,
    )

    assert "simple_function_python_simple_function_0" in aliases
    assert "simple_function" in aliases
    assert "solution" in aliases
    assert "your_module" in aliases
    assert "has_close_elements" in aliases


def test_rewrite_generated_test_imports_replaces_placeholder_imports(tmp_path: Path) -> None:
    src = tmp_path / "boundary_python_boundary_0.py"
    src.write_text(
        "def below_zero(operations):\n"
        "    return False\n",
        encoding="utf-8",
    )

    generated = (
        "from solution import below_zero\n"
        "from your_module import below_zero as alias_below_zero\n"
        "\n"
        "def test_case():\n"
        "    assert below_zero([1, -2]) is True\n"
    )

    rewritten = rewrite_generated_test_imports(
        generated_test_source=generated,
        source_path=src,
    )

    assert "from boundary_python_boundary_0 import below_zero" in rewritten
    assert "from solution import" not in rewritten
    assert "from your_module import" not in rewritten


def test_python_source_precheck_detects_relative_import_source(tmp_path: Path) -> None:
    src = tmp_path / "complex_dependency_000.py"
    src.write_text(
        "from .._models import Request\n"
        "\n"
        "def f():\n"
        "    return 1\n",
        encoding="utf-8",
    )

    err = python_source_precheck(src)
    assert err is not None
    assert "relative imports" in err


def test_compile_check_non_python_marks_not_implemented(tmp_path: Path) -> None:
    java_file = tmp_path / "generated.test.java"
    java_file.write_text("public class GeneratedTest {}\n", encoding="utf-8")

    ok, err = compile_check("java", java_file)
    assert ok is False
    assert err == "Language toolchain not implemented yet: java"


def test_rewrite_generated_test_imports_keeps_known_third_party_imports(tmp_path: Path) -> None:
    src = tmp_path / "sample.py"
    src.write_text(
        "def target():\n"
        "    return 1\n",
        encoding="utf-8",
    )

    generated = (
        "import pytest\n"
        "from sample import target\n"
        "\n"
        "def test_ok():\n"
        "    assert target() == 1\n"
    )

    rewritten = rewrite_generated_test_imports(
        generated_test_source=generated,
        source_path=src,
    )

    assert "import pytest" in rewritten
