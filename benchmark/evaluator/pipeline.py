from __future__ import annotations

import json
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

from .aggregate import build_summary
from .compiler import compile_check, python_source_precheck
from .contracts import EvaluationResult
from .coverage import collect_coverage, collect_mutation_score
from .executor import execute_tests
from .loader import load_generated_samples_with_mode
from .sample_policy import PythonSelfContainedPolicy


class Evaluator:
    """评测流水线：编译 -> 执行 -> 覆盖率 -> 变异 -> 聚合。"""

    def __init__(self, results_root: str | Path = "results") -> None:
        self.results_root = Path(results_root)
        self._python_policy = PythonSelfContainedPolicy()

    def run(
        self,
        models: list[str] | None = None,
        languages: list[str] | None = None,
        output_dir: str | Path | None = None,
        latest_only: bool = True,
        only_self_contained: bool = False,
    ) -> dict[str, Any]:
        all_samples = load_generated_samples_with_mode(
            results_root=self.results_root,
            latest_only=latest_only,
        )
        prefiltered_samples = [
            sample
            for sample in all_samples
            if (not models or sample.model in models)
            and (not languages or sample.language in languages)
        ]
        if not only_self_contained:
            samples = prefiltered_samples
        else:
            samples = []
            for sample in prefiltered_samples:
                if sample.language != "python":
                    continue
                source_path = Path(sample.source_path) if sample.source_path else None
                bucket, _ = self._python_policy.classify(source_path)
                if bucket == "self_contained":
                    samples.append(sample)

        eval_results: list[EvaluationResult] = []
        for sample in samples:
            generated_test_path = Path(sample.generated_test_path)
            source_path = Path(sample.source_path) if sample.source_path else None

            compile_pass, compile_error = compile_check(sample.language, generated_test_path)
            sample_bucket = None
            sample_bucket_reason = None
            if sample.language == "python":
                sample_bucket, sample_bucket_reason = self._python_policy.classify(source_path)

            test_pass: bool | None = None
            test_error: str | None = None
            runtime_ms: int | None = None
            source_precheck_error: str | None = None
            if compile_pass and sample.language == "python":
                source_precheck_error = python_source_precheck(source_path)

            if compile_pass:
                if source_precheck_error:
                    test_pass = False
                    test_error = f"source precheck failed: {source_precheck_error}"
                else:
                    test_pass, test_error, runtime_ms = execute_tests(
                        language=sample.language,
                        generated_test_path=generated_test_path,
                        source_path=source_path,
                    )

            line_cov = None
            branch_cov = None
            function_cov = None
            coverage_error = None
            mutation_score = None
            mutation_error = None
            mutation_total = None
            mutation_killed = None
            mutation_survived = None
            mutation_no_tests = None
            mutation_not_checked = None
            mutation_timeout = None
            mutation_skipped = None
            mutation_suspicious = None

            if test_pass:
                line_cov, branch_cov, function_cov, coverage_error = collect_coverage(
                    language=sample.language,
                    generated_test_path=generated_test_path,
                    source_path=source_path,
                )
                mutation_score, mutation_error, mutation_stats = collect_mutation_score(
                    language=sample.language,
                    generated_test_path=generated_test_path,
                    source_path=source_path,
                )
                if mutation_stats:
                    mutation_total = mutation_stats.get("total")
                    mutation_killed = mutation_stats.get("killed")
                    mutation_survived = mutation_stats.get("survived")
                    mutation_no_tests = mutation_stats.get("no_tests")
                    mutation_not_checked = mutation_stats.get("not_checked")
                    mutation_timeout = mutation_stats.get("timeout")
                    mutation_skipped = mutation_stats.get("skipped")
                    mutation_suspicious = mutation_stats.get("suspicious")

            eval_results.append(
                EvaluationResult(
                    model=sample.model,
                    language=sample.language,
                    sample_id=sample.sample_id,
                    generated_test_path=sample.generated_test_path,
                    source_path=sample.source_path,
                    compile_pass=compile_pass,
                    test_pass=test_pass,
                    line_coverage=line_cov,
                    branch_coverage=branch_cov,
                    function_coverage=function_cov,
                    mutation_score=mutation_score,
                    compile_error=compile_error,
                    test_error=test_error,
                    coverage_error=coverage_error,
                    mutation_error=mutation_error,
                    runtime_ms=runtime_ms,
                    sample_bucket=sample_bucket,
                    sample_bucket_reason=sample_bucket_reason,
                    mutation_total=mutation_total,
                    mutation_killed=mutation_killed,
                    mutation_survived=mutation_survived,
                    mutation_no_tests=mutation_no_tests,
                    mutation_not_checked=mutation_not_checked,
                    mutation_timeout=mutation_timeout,
                    mutation_skipped=mutation_skipped,
                    mutation_suspicious=mutation_suspicious,
                )
            )

        summary = build_summary(eval_results)
        payload = {
            "evaluated_at_utc": datetime.now(timezone.utc).isoformat(),
            "filters": {
                "models": models or [],
                "languages": languages or [],
                "latest_only": latest_only,
                "only_self_contained": only_self_contained,
            },
            "summary": summary,
            "results": [item.to_dict() for item in eval_results],
        }

        out_dir = Path(output_dir) if output_dir else self.results_root
        out_dir.mkdir(parents=True, exist_ok=True)
        run_id = datetime.now(timezone.utc).strftime("%Y%m%dT%H%M%S%fZ")
        out_file = out_dir / f"evaluator_summary_{run_id}.json"
        out_file.write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")
        return payload
