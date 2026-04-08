from __future__ import annotations

from typing import Any

from .common import avg, numeric_values, rate, to_int


def build_metrics(rows: list[dict[str, Any]]) -> dict[str, Any]:
    total = len(rows)
    compile_pass_count = sum(1 for row in rows if row.get("compile_pass") is True)
    executable = sum(1 for row in rows if row.get("test_pass") is not None)
    test_pass_count = sum(1 for row in rows if row.get("test_pass") is True)
    mutation_total_mutants = sum(to_int(row.get("mutation_total")) or 0 for row in rows)
    mutation_killed_mutants = sum(to_int(row.get("mutation_killed")) or 0 for row in rows)
    mutation_survived_mutants = sum(to_int(row.get("mutation_survived")) or 0 for row in rows)
    mutation_no_tests_mutants = sum(to_int(row.get("mutation_no_tests")) or 0 for row in rows)
    mutation_not_checked_mutants = sum(to_int(row.get("mutation_not_checked")) or 0 for row in rows)
    mutation_timeout_mutants = sum(to_int(row.get("mutation_timeout")) or 0 for row in rows)
    mutation_skipped_mutants = sum(to_int(row.get("mutation_skipped")) or 0 for row in rows)
    mutation_suspicious_mutants = sum(to_int(row.get("mutation_suspicious")) or 0 for row in rows)
    mutation_effective_mutants = mutation_killed_mutants + mutation_survived_mutants

    return {
        "total_samples": total,
        "compile_pass_count": compile_pass_count,
        "compile_pass_rate": rate(compile_pass_count, total),
        "test_pass_count": test_pass_count,
        "test_pass_rate": rate(test_pass_count, compile_pass_count),
        "stage2_executable_samples": executable,
        "avg_line_coverage": avg(numeric_values(rows, "line_coverage")),
        "avg_branch_coverage": avg(numeric_values(rows, "branch_coverage")),
        "avg_function_coverage": avg(numeric_values(rows, "function_coverage")),
        "avg_mutation_score": avg(numeric_values(rows, "mutation_score")),
        "avg_eval_runtime_ms": avg(numeric_values(rows, "eval_runtime_ms")),
        "avg_generation_latency_ms": avg(numeric_values(rows, "generation_latency_ms")),
        "avg_prompt_tokens": avg(numeric_values(rows, "prompt_tokens")),
        "avg_completion_tokens": avg(numeric_values(rows, "completion_tokens")),
        "avg_total_tokens": avg(numeric_values(rows, "total_tokens")),
        "mutation_total_mutants": mutation_total_mutants,
        "mutation_killed_mutants": mutation_killed_mutants,
        "mutation_survived_mutants": mutation_survived_mutants,
        "mutation_no_tests_mutants": mutation_no_tests_mutants,
        "mutation_not_checked_mutants": mutation_not_checked_mutants,
        "mutation_timeout_mutants": mutation_timeout_mutants,
        "mutation_skipped_mutants": mutation_skipped_mutants,
        "mutation_suspicious_mutants": mutation_suspicious_mutants,
        "mutation_kill_rate": rate(mutation_killed_mutants, mutation_total_mutants),
        "mutation_effective_mutants": mutation_effective_mutants,
        "mutation_effective_kill_rate": rate(mutation_killed_mutants, mutation_effective_mutants),
    }


def group_metrics(rows: list[dict[str, Any]], group_keys: list[str]) -> list[dict[str, Any]]:
    grouped: dict[tuple[str, ...], list[dict[str, Any]]] = {}
    for row in rows:
        key = tuple(str(row.get(field) or "unknown") for field in group_keys)
        grouped.setdefault(key, []).append(row)

    output: list[dict[str, Any]] = []
    for key, bucket in grouped.items():
        record = {group_keys[idx]: key[idx] for idx in range(len(group_keys))}
        record.update(build_metrics(bucket))
        output.append(record)

    output.sort(key=lambda row: tuple(str(row.get(field) or "") for field in group_keys))
    return output
