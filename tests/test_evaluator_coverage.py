from __future__ import annotations

from pathlib import Path

from benchmark.evaluator.coverage import collect_coverage


def test_collect_coverage_python_smoke(tmp_path: Path) -> None:
    src = tmp_path / "calculator.py"
    src.write_text(
        """
def add(a, b):
    return a + b
""".strip()
        + "\n",
        encoding="utf-8",
    )

    test_file = tmp_path / "test_generated.py"
    test_file.write_text(
        """
from calculator import add

def test_add():
    assert add(1, 2) == 3
""".strip()
        + "\n",
        encoding="utf-8",
    )

    line_cov, branch_cov, func_cov, err = collect_coverage(
        language="python",
        generated_test_path=test_file,
        source_path=src,
    )

    assert err is None
    assert line_cov is not None
    assert 0.0 <= line_cov <= 1.0
    # 由于 add 函数无分支，branch_cov 可能为 None；若有值应在合法区间。
    if branch_cov is not None:
        assert 0.0 <= branch_cov <= 1.0
    assert func_cov is None


def test_collect_coverage_python_requires_source_in_report(tmp_path: Path) -> None:
    src = tmp_path / "calculator.py"
    src.write_text(
        """
def add(a, b):
    return a + b
""".strip()
        + "\n",
        encoding="utf-8",
    )

    test_file = tmp_path / "test_generated.py"
    test_file.write_text(
        """
def test_ok():
    assert True
""".strip()
        + "\n",
        encoding="utf-8",
    )

    line_cov, branch_cov, func_cov, err = collect_coverage(
        language="python",
        generated_test_path=test_file,
        source_path=src,
    )

    assert line_cov is None
    assert branch_cov is None
    assert func_cov is None
    assert err is not None
    assert "source file not found in coverage report" in err
