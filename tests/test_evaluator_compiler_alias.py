from __future__ import annotations

from pathlib import Path

from benchmark.evaluator.compiler import python_source_alias_candidates


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
