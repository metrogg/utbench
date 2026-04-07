from __future__ import annotations

from pathlib import Path

from benchmark.evaluator.sample_policy import PythonSelfContainedPolicy


def test_policy_classifies_allowlisted_sample(tmp_path: Path) -> None:
    allowlist = tmp_path / "allowlist.txt"
    allowlist.write_text("boundary/boundary_000.py\n", encoding="utf-8")

    source = tmp_path / "dataset" / "python" / "boundary" / "boundary_000.py"
    source.parent.mkdir(parents=True)
    source.write_text("def f():\n    return 1\n", encoding="utf-8")

    policy = PythonSelfContainedPolicy(allowlist_path=allowlist)
    bucket, reason = policy.classify(source)

    assert bucket == "self_contained"
    assert reason is None


def test_policy_classifies_non_allowlisted_sample(tmp_path: Path) -> None:
    allowlist = tmp_path / "allowlist.txt"
    allowlist.write_text("simple_function/simple_function_000.py\n", encoding="utf-8")

    source = tmp_path / "dataset" / "python" / "boundary" / "boundary_999.py"
    source.parent.mkdir(parents=True)
    source.write_text("def f():\n    return 1\n", encoding="utf-8")

    policy = PythonSelfContainedPolicy(allowlist_path=allowlist)
    bucket, reason = policy.classify(source)

    assert bucket == "non_self_contained"
    assert isinstance(reason, str)
    assert "not allowlisted" in reason
