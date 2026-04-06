from __future__ import annotations

from benchmark.evaluator.aggregate import build_summary
from benchmark.evaluator.contracts import EvaluationResult


def test_build_summary_basic_rates() -> None:
    rows = [
        EvaluationResult(
            model="m1",
            language="python",
            sample_id="s1",
            generated_test_path="a",
            source_path=None,
            compile_pass=True,
            test_pass=True,
            line_coverage=0.8,
            branch_coverage=None,
            function_coverage=None,
            mutation_score=None,
            compile_error=None,
            test_error=None,
            coverage_error=None,
            mutation_error=None,
            runtime_ms=100,
        ),
        EvaluationResult(
            model="m1",
            language="python",
            sample_id="s2",
            generated_test_path="b",
            source_path=None,
            compile_pass=False,
            test_pass=None,
            line_coverage=None,
            branch_coverage=None,
            function_coverage=None,
            mutation_score=None,
            compile_error="syntax error",
            test_error=None,
            coverage_error=None,
            mutation_error=None,
            runtime_ms=None,
        ),
    ]

    summary = build_summary(rows)
    assert summary["total_samples"] == 2
    assert summary["compile_pass_rate"] == 0.5
    assert summary["test_pass_rate"] == 1.0
    assert summary["mutation_executed_samples"] == 0
    assert summary["mutation_total_mutants"] == 0
    assert summary["mutation_killed_mutants"] == 0
    assert summary["mutation_survived_mutants"] == 0
    assert summary["mutation_no_tests_mutants"] == 0
    assert summary["mutation_not_checked_mutants"] == 0
    assert summary["mutation_timeout_mutants"] == 0
    assert summary["mutation_skipped_mutants"] == 0
    assert summary["mutation_suspicious_mutants"] == 0
    assert summary["mutation_kill_rate"] == 0.0
    assert summary["mutation_effective_mutants"] == 0
    assert summary["mutation_effective_kill_rate"] == 0.0
