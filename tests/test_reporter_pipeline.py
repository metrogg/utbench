from __future__ import annotations

import json
from pathlib import Path

from benchmark.reporter.pipeline import Reporter


def _write_json(path: Path, payload: dict) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")


def test_reporter_generates_json_csv_html(tmp_path: Path) -> None:
    results_root = tmp_path / "results"
    evaluator_summary = results_root / "evaluator_summary_20260101T000000000000Z.json"

    generated_test_path = (
        results_root
        / "deepseek"
        / "tests"
        / "deepseek_python_simple_function_python_simple_function_0_20260101010101.test.py"
    )
    generated_test_path.parent.mkdir(parents=True, exist_ok=True)
    generated_test_path.write_text("def test_placeholder():\n    assert True\n", encoding="utf-8")

    metadata_path = (
        results_root
        / "deepseek"
        / "reports"
        / "deepseek_python_simple_function_python_simple_function_0_20260101010101.metadata.json"
    )
    _write_json(
        metadata_path,
        {
            "latency_ms": 1234,
            "tokens": {
                "prompt_tokens": 11,
                "completion_tokens": 22,
                "total_tokens": 33,
            },
            "success": True,
        },
    )

    _write_json(
        evaluator_summary,
        {
            "filters": {"models": ["deepseek"], "languages": ["python"]},
            "summary": {"total_samples": 1},
            "results": [
                {
                    "model": "deepseek",
                    "language": "python",
                    "sample_id": "simple_function_python_simple_function_0",
                    "generated_test_path": str(generated_test_path),
                    "source_path": "dataset/python/simple_function_python_simple_function_0.py",
                    "compile_pass": True,
                    "test_pass": True,
                    "line_coverage": 0.8,
                    "branch_coverage": 0.6,
                    "function_coverage": None,
                    "mutation_score": 0.7,
                    "mutation_total": 10,
                    "mutation_killed": 7,
                    "mutation_survived": 3,
                    "mutation_no_tests": 0,
                    "mutation_not_checked": 0,
                    "mutation_timeout": 0,
                    "mutation_skipped": 0,
                    "mutation_suspicious": 0,
                    "mutation_error": None,
                    "compile_error": None,
                    "test_error": None,
                    "coverage_error": None,
                    "runtime_ms": 321,
                }
            ],
        },
    )

    reporter = Reporter(results_root=results_root)
    payload = reporter.run(evaluator_summary=evaluator_summary)

    assert payload["summary"]["total_samples"] == 1
    assert payload["summary"]["compile_pass_rate"] == 1.0
    assert payload["summary"]["test_pass_rate"] == 1.0
    assert payload["summary"]["avg_generation_latency_ms"] == 1234.0
    assert payload["summary"]["avg_total_tokens"] == 33.0
    assert payload["summary"]["mutation_total_mutants"] == 10
    assert payload["summary"]["mutation_killed_mutants"] == 7
    assert payload["summary"]["mutation_survived_mutants"] == 3
    assert payload["summary"]["mutation_not_checked_mutants"] == 0
    assert payload["summary"]["mutation_effective_mutants"] == 10
    assert payload["summary"]["mutation_effective_kill_rate"] == 0.7
    sample = payload["report"]["sample_rows"][0]
    assert sample["mutation_total"] == 10
    assert sample["mutation_killed"] == 7
    assert sample["mutation_survived"] == 3
    assert sample["mutation_no_tests"] == 0
    assert sample["mutation_not_checked"] == 0
    assert sample["mutation_timeout"] == 0
    assert sample["generation_metrics_source"] == "metadata"
    assert sample["generation_metrics_reason"] is None

    artifacts = payload["artifacts"]
    assert Path(artifacts["json"]).exists()
    assert Path(artifacts["by_model_csv"]).exists()
    assert Path(artifacts["by_language_csv"]).exists()
    assert Path(artifacts["sample_csv"]).exists()
    assert Path(artifacts["html"]).exists()


def test_reporter_failure_breakdown_extracts_types(tmp_path: Path) -> None:
    results_root = tmp_path / "results"
    evaluator_summary = results_root / "evaluator_summary_20260101T000000000000Z.json"

    _write_json(
        evaluator_summary,
        {
            "filters": {},
            "summary": {},
            "results": [
                {
                    "model": "m1",
                    "language": "python",
                    "sample_id": "boundary_python_boundary_0",
                    "generated_test_path": "results/m1/tests/a.test.py",
                    "source_path": "dataset/python/boundary_python_boundary_0.py",
                    "compile_pass": True,
                    "test_pass": False,
                    "line_coverage": None,
                    "branch_coverage": None,
                    "function_coverage": None,
                    "mutation_score": None,
                    "compile_error": None,
                    "test_error": "ModuleNotFoundError: No module named 'solution'",
                    "coverage_error": None,
                    "mutation_error": None,
                    "runtime_ms": 100,
                    "mutation_total": None,
                    "mutation_killed": None,
                    "mutation_survived": None,
                    "mutation_no_tests": None,
                    "mutation_not_checked": None,
                    "mutation_timeout": None,
                    "mutation_skipped": None,
                    "mutation_suspicious": None,
                },
                {
                    "model": "m2",
                    "language": "python",
                    "sample_id": "boundary_python_boundary_1",
                    "generated_test_path": "results/m2/tests/b.test.py",
                    "source_path": "dataset/python/boundary_python_boundary_1.py",
                    "compile_pass": False,
                    "test_pass": None,
                    "line_coverage": None,
                    "branch_coverage": None,
                    "function_coverage": None,
                    "mutation_score": None,
                    "compile_error": "SyntaxError: invalid syntax",
                    "test_error": None,
                    "coverage_error": None,
                    "mutation_error": None,
                    "runtime_ms": None,
                    "mutation_total": None,
                    "mutation_killed": None,
                    "mutation_survived": None,
                    "mutation_no_tests": None,
                    "mutation_not_checked": None,
                    "mutation_timeout": None,
                    "mutation_skipped": None,
                    "mutation_suspicious": None,
                },
            ],
        },
    )

    reporter = Reporter(results_root=results_root)
    payload = reporter.run(evaluator_summary=evaluator_summary, top_n_errors=10)

    failure_rows = payload["report"]["failure_breakdown"]
    types = {(row["stage"], row["error_type"]) for row in failure_rows}
    assert ("test", "module_not_found") in types
    assert ("compile", "syntax_error") in types


def test_reporter_data_quality_panel_metrics_present(tmp_path: Path) -> None:
    results_root = tmp_path / "results"
    evaluator_summary = results_root / "evaluator_summary_20260101T000000000000Z.json"

    _write_json(
        evaluator_summary,
        {
            "filters": {},
            "summary": {},
            "results": [
                {
                    "model": "m1",
                    "language": "python",
                    "sample_id": "s1_python_boundary_0",
                    "generated_test_path": "results/m1/tests/a.test.py",
                    "source_path": "dataset/python/s1_python_boundary_0.py",
                    "compile_pass": True,
                    "test_pass": True,
                    "line_coverage": 0.9,
                    "branch_coverage": 0.7,
                    "function_coverage": None,
                    "mutation_score": 0.8,
                    "mutation_total": 5,
                    "mutation_killed": 4,
                    "mutation_survived": 1,
                    "mutation_no_tests": 0,
                    "mutation_not_checked": 0,
                    "mutation_timeout": 0,
                    "mutation_skipped": 0,
                    "mutation_suspicious": 0,
                    "compile_error": None,
                    "test_error": None,
                    "coverage_error": None,
                    "mutation_error": None,
                    "runtime_ms": 100,
                }
            ],
        },
    )

    reporter = Reporter(results_root=results_root)
    payload = reporter.run(evaluator_summary=evaluator_summary)
    html = Path(payload["artifacts"]["html"]).read_text(encoding="utf-8")
    assert "Metadata 完整样本" in html
    assert "覆盖率有效样本" in html
    assert "变异统计有效样本" in html


def test_reporter_respects_export_formats(tmp_path: Path) -> None:
    results_root = tmp_path / "results"
    evaluator_summary = results_root / "evaluator_summary_20260101T000000000000Z.json"

    _write_json(
        evaluator_summary,
        {
            "filters": {},
            "summary": {},
            "results": [
                {
                    "model": "m1",
                    "language": "python",
                    "sample_id": "simple_function_python_simple_function_0",
                    "generated_test_path": "results/m1/tests/a.test.py",
                    "source_path": "dataset/python/simple_function_python_simple_function_0.py",
                    "compile_pass": True,
                    "test_pass": True,
                    "line_coverage": 0.9,
                    "branch_coverage": 0.8,
                    "function_coverage": None,
                    "mutation_score": 0.7,
                    "compile_error": None,
                    "test_error": None,
                    "coverage_error": None,
                    "mutation_error": None,
                    "runtime_ms": 100,
                }
            ],
        },
    )

    reporter = Reporter(results_root=results_root)
    payload = reporter.run(
        evaluator_summary=evaluator_summary,
        formats=["json", "html"],
        chart_style="orange",
        thresholds={"test_pass_rate": 0.8},
    )

    artifacts = payload["artifacts"]
    assert "json" in artifacts
    assert "html" in artifacts
    assert "sample_csv" not in artifacts
    assert payload["report"]["chart_style"] == "orange"
    assert payload["report"]["thresholds"]["test_pass_rate"] == 0.8
