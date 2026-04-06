from __future__ import annotations

from typing import Any

from .contracts import EvaluationResult


def build_summary(results: list[EvaluationResult]) -> dict[str, Any]:
    """阶段5：聚合指标并输出汇总。"""
    total = len(results)
    compile_pass_count = sum(1 for row in results if row.compile_pass)
    executable_rows = [row for row in results if row.test_pass is not None]
    test_pass_count = sum(1 for row in results if row.test_pass is True)

    line_cov_values = [row.line_coverage for row in results if row.line_coverage is not None]
    branch_cov_values = [row.branch_coverage for row in results if row.branch_coverage is not None]
    func_cov_values = [row.function_coverage for row in results if row.function_coverage is not None]
    mutation_values = [row.mutation_score for row in results if row.mutation_score is not None]
    mutation_executed = [row for row in results if row.mutation_total is not None]
    mutation_total_mutants = sum(row.mutation_total or 0 for row in results)
    mutation_killed_mutants = sum(row.mutation_killed or 0 for row in results)
    mutation_survived_mutants = sum(row.mutation_survived or 0 for row in results)
    mutation_no_tests_mutants = sum(row.mutation_no_tests or 0 for row in results)
    mutation_not_checked_mutants = sum(row.mutation_not_checked or 0 for row in results)
    mutation_timeout_mutants = sum(row.mutation_timeout or 0 for row in results)
    mutation_skipped_mutants = sum(row.mutation_skipped or 0 for row in results)
    mutation_suspicious_mutants = sum(row.mutation_suspicious or 0 for row in results)
    mutation_effective_mutants = mutation_killed_mutants + mutation_survived_mutants

    return {
        "total_samples": total,
        "compile_pass_count": compile_pass_count,
        "compile_pass_rate": _rate(compile_pass_count, total),
        "test_pass_count": test_pass_count,
        "test_pass_rate": _rate(test_pass_count, compile_pass_count),
        "avg_line_coverage": _avg(line_cov_values),
        "avg_branch_coverage": _avg(branch_cov_values),
        "avg_function_coverage": _avg(func_cov_values),
        "avg_mutation_score": _avg(mutation_values),
        "stage2_executable_samples": len(executable_rows),
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
    }


def _avg(values: list[float]) -> float | None:
    if not values:
        return None
    return round(sum(values) / len(values), 6)


def _rate(numerator: int, denominator: int) -> float:
    if denominator <= 0:
        return 0.0
    return round(numerator / denominator, 6)
