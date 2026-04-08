from __future__ import annotations

from typing import Any

from .contracts import EvaluationResult


def build_summary(results: list[EvaluationResult]) -> dict[str, Any]:
    """阶段5：聚合指标并输出汇总。"""
    total = len(results)
    compile_pass_count = sum(1 for row in results if row.compile_pass)
    executable_rows = [row for row in results if row.test_pass is not None]
    test_pass_count = sum(1 for row in results if row.test_pass is True)
    test_pass_rate_values = [row.test_pass_rate for row in results if row.test_pass_rate is not None]

    line_cov_values = [row.line_coverage for row in results if row.line_coverage is not None]
    branch_cov_values = [row.branch_coverage for row in results if row.branch_coverage is not None]
    func_cov_values = [row.function_coverage for row in results if row.function_coverage is not None]
    passed_line_cov_values = [row.passed_line_coverage for row in results if row.passed_line_coverage is not None]
    passed_branch_cov_values = [row.passed_branch_coverage for row in results if row.passed_branch_coverage is not None]
    passed_func_cov_values = [row.passed_function_coverage for row in results if row.passed_function_coverage is not None]
    mutation_values = [row.mutation_score for row in results if row.mutation_score is not None]
    mutation_stats = [row for row in results if row.mutation_total is not None]
    mutation_executed = [row for row in results if _mutation_processed(row) > 0]
    mutation_total_mutants = sum(row.mutation_total or 0 for row in results)
    mutation_killed_mutants = sum(row.mutation_killed or 0 for row in results)
    mutation_survived_mutants = sum(row.mutation_survived or 0 for row in results)
    mutation_no_tests_mutants = sum(row.mutation_no_tests or 0 for row in results)
    mutation_not_checked_mutants = sum(row.mutation_not_checked or 0 for row in results)
    mutation_timeout_mutants = sum(row.mutation_timeout or 0 for row in results)
    mutation_skipped_mutants = sum(row.mutation_skipped or 0 for row in results)
    mutation_suspicious_mutants = sum(row.mutation_suspicious or 0 for row in results)
    mutation_effective_mutants = mutation_killed_mutants + mutation_survived_mutants
    self_contained_samples = [row for row in results if row.sample_bucket == "self_contained"]
    non_self_contained_samples = [row for row in results if row.sample_bucket == "non_self_contained"]
    unknown_bucket_samples = [
        row
        for row in results
        if row.sample_bucket not in {"self_contained", "non_self_contained"}
    ]

    sc_compile_pass_count = sum(1 for row in self_contained_samples if row.compile_pass)
    sc_test_pass_count = sum(1 for row in self_contained_samples if row.test_pass is True)
    nsc_compile_pass_count = sum(1 for row in non_self_contained_samples if row.compile_pass)
    nsc_test_pass_count = sum(1 for row in non_self_contained_samples if row.test_pass is True)

    return {
        "total_samples": total,
        "compile_pass_count": compile_pass_count,
        "compile_pass_rate": _rate(compile_pass_count, total),
        "test_pass_count": test_pass_count,
        "test_pass_rate": _rate(test_pass_count, compile_pass_count),
        "avg_test_pass_rate": _avg(test_pass_rate_values),
        "avg_line_coverage": _avg(line_cov_values),
        "avg_branch_coverage": _avg(branch_cov_values),
        "avg_function_coverage": _avg(func_cov_values),
        "avg_passed_line_coverage": _avg(passed_line_cov_values),
        "avg_passed_branch_coverage": _avg(passed_branch_cov_values),
        "avg_passed_function_coverage": _avg(passed_func_cov_values),
        "avg_mutation_score": _avg(mutation_values),
        "stage2_executable_samples": len(executable_rows),
        "mutation_stats_samples": len(mutation_stats),
        "mutation_executed_samples": len(mutation_executed),
        "mutation_total_mutants": mutation_total_mutants,
        "mutation_killed_mutants": mutation_killed_mutants,
        "mutation_survived_mutants": mutation_survived_mutants,
        "mutation_no_tests_mutants": mutation_no_tests_mutants,
        "mutation_not_checked_mutants": mutation_not_checked_mutants,
        "mutation_timeout_mutants": mutation_timeout_mutants,
        "mutation_skipped_mutants": mutation_skipped_mutants,
        "mutation_suspicious_mutants": mutation_suspicious_mutants,
        "mutation_kill_rate": _rate(mutation_killed_mutants, mutation_total_mutants),
        "mutation_effective_mutants": mutation_effective_mutants,
        "mutation_effective_kill_rate": _rate(mutation_killed_mutants, mutation_effective_mutants),
        "sample_bucket_counts": {
            "self_contained": len(self_contained_samples),
            "non_self_contained": len(non_self_contained_samples),
            "unknown": len(unknown_bucket_samples),
        },
        "self_contained_compile_pass_rate": _rate(sc_compile_pass_count, len(self_contained_samples)),
        "self_contained_test_pass_rate": _rate(sc_test_pass_count, sc_compile_pass_count),
        "non_self_contained_compile_pass_rate": _rate(
            nsc_compile_pass_count,
            len(non_self_contained_samples),
        ),
        "non_self_contained_test_pass_rate": _rate(nsc_test_pass_count, nsc_compile_pass_count),
    }


def _mutation_processed(row: EvaluationResult) -> int:
    fields = (
        row.mutation_killed,
        row.mutation_survived,
        row.mutation_no_tests,
        row.mutation_timeout,
        row.mutation_skipped,
        row.mutation_suspicious,
    )
    return sum(value or 0 for value in fields)


def _avg(values: list[float]) -> float | None:
    if not values:
        return None
    return round(sum(values) / len(values), 6)


def _rate(numerator: int, denominator: int) -> float:
    if denominator <= 0:
        return 0.0
    return round(numerator / denominator, 6)
