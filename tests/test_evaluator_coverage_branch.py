from __future__ import annotations

from pathlib import Path

from benchmark.evaluator.coverage import collect_coverage


def test_collect_coverage_python_includes_branch_when_present(tmp_path: Path) -> None:
    src = tmp_path / "branchy.py"
    src.write_text(
        """
def sign(v):
    if v > 0:
        return 1
    return -1
""".strip()
        + "\n",
        encoding="utf-8",
    )

    test_file = tmp_path / "test_branchy.py"
    test_file.write_text(
        """
from branchy import sign

def test_sign_positive_only():
    assert sign(2) == 1
""".strip()
        + "\n",
        encoding="utf-8",
    )

    line_cov, branch_cov, _func_cov, err = collect_coverage(
        language="python",
        generated_test_path=test_file,
        source_path=src,
    )

    assert err is None
    assert line_cov is not None
    assert 0.0 <= line_cov <= 1.0
    # 存在 if 分支时应给出分支覆盖率。
    assert branch_cov is not None
    assert 0.0 <= branch_cov <= 1.0
