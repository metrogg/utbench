from __future__ import annotations

import json
from pathlib import Path

from benchmark.evaluator.pipeline import Evaluator


def test_evaluator_pipeline_runs_and_returns_summary() -> None:
    evaluator = Evaluator(results_root="results")
    payload = evaluator.run(models=["deepseek"], languages=["python"])
    assert "summary" in payload
    assert "results" in payload
    assert payload["summary"]["total_samples"] >= 0
    assert payload["filters"]["latest_only"] is True


def test_evaluator_marks_non_self_contained_source_as_test_failure(tmp_path: Path) -> None:
    results_root = tmp_path / "results"
    model = "m1"
    tests_dir = results_root / model / "tests"
    reports_dir = results_root / model / "reports"
    tests_dir.mkdir(parents=True)
    reports_dir.mkdir(parents=True)

    source = tmp_path / "dataset" / "python" / "complex_dependency" / "complex_dependency_000.py"
    source.parent.mkdir(parents=True)
    source.write_text(
        "from ..pkg import x\n"
        "\n"
        "def f():\n"
        "    return x\n",
        encoding="utf-8",
    )

    test_file = tests_dir / f"{model}_python_complex_dependency_000_20260101010101.test.py"
    test_file.write_text("def test_placeholder():\n    assert True\n", encoding="utf-8")

    metadata_file = reports_dir / f"{model}_python_complex_dependency_000_20260101010101.metadata.json"
    metadata_file.write_text(
        json.dumps({"sample_path": str(source.resolve())}, ensure_ascii=False),
        encoding="utf-8",
    )

    evaluator = Evaluator(results_root=results_root)
    payload = evaluator.run(models=[model], languages=["python"], latest_only=True)

    assert payload["summary"]["total_samples"] == 1
    result = payload["results"][0]
    assert result["compile_pass"] is True
    assert result["test_pass"] is False
    assert isinstance(result["test_error"], str)
    assert "source precheck failed" in result["test_error"]
    assert result["sample_bucket"] == "non_self_contained"


def test_evaluator_only_self_contained_filter(tmp_path: Path) -> None:
    results_root = tmp_path / "results"
    model = "m1"
    tests_dir = results_root / model / "tests"
    reports_dir = results_root / model / "reports"
    tests_dir.mkdir(parents=True)
    reports_dir.mkdir(parents=True)

    # allowlisted sample
    src_sc = tmp_path / "dataset" / "python" / "boundary_python_boundary_0.py"
    src_sc.parent.mkdir(parents=True)
    src_sc.write_text("def f():\n    return 1\n", encoding="utf-8")
    test_sc = tests_dir / f"{model}_python_boundary_python_boundary_0_20260101010101.test.py"
    test_sc.write_text(
        "from boundary_python_boundary_0 import f\n\ndef test_f():\n    assert f()==1\n",
        encoding="utf-8",
    )
    meta_sc = reports_dir / f"{model}_python_boundary_python_boundary_0_20260101010101.metadata.json"
    meta_sc.write_text(json.dumps({"sample_path": str(src_sc.resolve())}), encoding="utf-8")

    # non-allowlisted sample
    src_nsc = tmp_path / "dataset" / "python" / "other" / "other_000.py"
    src_nsc.parent.mkdir(parents=True)
    src_nsc.write_text("def g():\n    return 2\n", encoding="utf-8")
    test_nsc = tests_dir / f"{model}_python_other_000_20260101010102.test.py"
    test_nsc.write_text("from other_000 import g\n\ndef test_g():\n    assert g()==2\n", encoding="utf-8")
    meta_nsc = reports_dir / f"{model}_python_other_000_20260101010102.metadata.json"
    meta_nsc.write_text(json.dumps({"sample_path": str(src_nsc.resolve())}), encoding="utf-8")

    evaluator = Evaluator(results_root=results_root)
    payload = evaluator.run(models=[model], languages=["python"], latest_only=True, only_self_contained=True)

    assert payload["summary"]["total_samples"] == 1
    assert payload["results"][0]["sample_id"] == "boundary_python_boundary_0"
