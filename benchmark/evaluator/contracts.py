from __future__ import annotations

from dataclasses import asdict, dataclass
from typing import Any


@dataclass
class EvaluationSample:
    model: str
    language: str
    sample_id: str
    generated_test_path: str
    source_path: str | None
    timestamp: str


@dataclass
class EvaluationResult:
    model: str
    language: str
    sample_id: str
    generated_test_path: str
    source_path: str | None
    compile_pass: bool
    test_pass: bool | None
    line_coverage: float | None
    branch_coverage: float | None
    function_coverage: float | None
    mutation_score: float | None
    compile_error: str | None
    test_error: str | None
    coverage_error: str | None
    mutation_error: str | None
    runtime_ms: int | None
    sample_bucket: str | None = None
    sample_bucket_reason: str | None = None
    mutation_total: int | None = None
    mutation_killed: int | None = None
    mutation_survived: int | None = None
    mutation_no_tests: int | None = None
    mutation_not_checked: int | None = None
    mutation_timeout: int | None = None
    mutation_skipped: int | None = None
    mutation_suspicious: int | None = None
    test_pass_count: int | None = None
    test_total_count: int | None = None
    passed_line_coverage: float | None = None
    passed_branch_coverage: float | None = None
    passed_function_coverage: float | None = None

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)
