from __future__ import annotations

from pathlib import Path


def test_expected_shell_scripts_exist() -> None:
    repo_root = Path(__file__).resolve().parents[1]
    assert (repo_root / "scripts" / "setup.sh").exists()
    assert (repo_root / "scripts" / "run_benchmark.sh").exists()
    assert (repo_root / "scripts" / "gen_report.sh").exists()
