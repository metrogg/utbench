from __future__ import annotations

import json
from pathlib import Path
from typing import Any

from .common import to_float, to_int
from .contracts import _LANGUAGE_TOKENS


def read_json(path: Path) -> dict[str, Any]:
    payload = json.loads(path.read_text(encoding="utf-8"))
    if not isinstance(payload, dict):
        raise ValueError(f"Invalid JSON payload in {path}")
    return payload


def enrich_result_row(
    row: dict[str, Any],
    results_root: Path,
    runner_metrics: dict[str, dict[str, int | None]],
) -> dict[str, Any]:
    sample_id = str(row.get("sample_id") or "")
    language = str(row.get("language") or "")
    complexity, scenario = parse_sample_id(sample_id, language)
    generation_metrics = extract_generation_metrics(row, results_root, runner_metrics)

    return {
        "model": row.get("model"),
        "language": language or None,
        "sample_id": sample_id,
        "complexity": complexity,
        "scenario": scenario,
        "generated_test_path": row.get("generated_test_path"),
        "source_path": row.get("source_path"),
        "compile_pass": row.get("compile_pass"),
        "test_pass": row.get("test_pass"),
        "line_coverage": to_float(row.get("line_coverage")),
        "branch_coverage": to_float(row.get("branch_coverage")),
        "function_coverage": to_float(row.get("function_coverage")),
        "mutation_score": to_float(row.get("mutation_score")),
        "mutation_total": to_int(row.get("mutation_total")),
        "mutation_killed": to_int(row.get("mutation_killed")),
        "mutation_survived": to_int(row.get("mutation_survived")),
        "mutation_no_tests": to_int(row.get("mutation_no_tests")),
        "mutation_not_checked": to_int(row.get("mutation_not_checked")),
        "mutation_timeout": to_int(row.get("mutation_timeout")),
        "mutation_skipped": to_int(row.get("mutation_skipped")),
        "mutation_suspicious": to_int(row.get("mutation_suspicious")),
        "eval_runtime_ms": to_int(row.get("runtime_ms")),
        "generation_latency_ms": generation_metrics.get("generation_latency_ms"),
        "prompt_tokens": generation_metrics.get("prompt_tokens"),
        "completion_tokens": generation_metrics.get("completion_tokens"),
        "total_tokens": generation_metrics.get("total_tokens"),
        "generation_metrics_source": generation_metrics.get("generation_metrics_source"),
        "generation_metrics_reason": generation_metrics.get("generation_metrics_reason"),
        "compile_error": row.get("compile_error"),
        "test_error": row.get("test_error"),
        "coverage_error": row.get("coverage_error"),
        "mutation_error": row.get("mutation_error"),
    }


def parse_sample_id(sample_id: str, language: str) -> tuple[str | None, str | None]:
    if not sample_id:
        return None, None

    parts = sample_id.split("_")
    if len(parts) < 3:
        return None, None

    lang_idx = next((idx for idx, part in enumerate(parts) if part in _LANGUAGE_TOKENS), None)
    if lang_idx is None and language in _LANGUAGE_TOKENS and language in parts:
        lang_idx = parts.index(language)
    if lang_idx is None:
        return None, None

    complexity = "_".join(parts[:lang_idx]) or None
    scenario_parts = parts[lang_idx + 1 :]
    if scenario_parts and scenario_parts[-1].isdigit():
        scenario_parts = scenario_parts[:-1]
    scenario = "_".join(scenario_parts) or None
    return complexity, scenario


def extract_generation_metrics(
    row: dict[str, Any],
    results_root: Path,
    runner_metrics: dict[str, dict[str, int | None]],
) -> dict[str, int | None | str]:
    generated_test_path = str(row.get("generated_test_path") or "")
    metadata_file = guess_metadata_path(generated_test_path, results_root)
    if metadata_file and metadata_file.exists():
        payload = read_json(metadata_file)
        tokens = payload.get("tokens", {})
        if not isinstance(tokens, dict):
            tokens = {}
        return {
            "generation_latency_ms": to_int(payload.get("latency_ms")),
            "prompt_tokens": to_int(tokens.get("prompt_tokens")),
            "completion_tokens": to_int(tokens.get("completion_tokens")),
            "total_tokens": to_int(tokens.get("total_tokens")),
            "generation_metrics_source": "metadata",
            "generation_metrics_reason": None,
        }

    fallback = runner_metrics.get(Path(generated_test_path).name)
    if fallback:
        return {
            "generation_latency_ms": fallback.get("generation_latency_ms"),
            "prompt_tokens": fallback.get("prompt_tokens"),
            "completion_tokens": fallback.get("completion_tokens"),
            "total_tokens": fallback.get("total_tokens"),
            "generation_metrics_source": "runner_summary",
            "generation_metrics_reason": "metadata_missing_fallback_runner_summary",
        }
    return {
        "generation_latency_ms": None,
        "prompt_tokens": None,
        "completion_tokens": None,
        "total_tokens": None,
        "generation_metrics_source": "missing",
        "generation_metrics_reason": "metadata_and_runner_summary_missing",
    }


def guess_metadata_path(generated_test_path: str, results_root: Path) -> Path | None:
    if not generated_test_path:
        return None

    candidates = [Path(generated_test_path)]
    if not Path(generated_test_path).is_absolute():
        candidates.append((Path.cwd() / generated_test_path).resolve())
        candidates.append((results_root.parent / generated_test_path).resolve())

    for test_path in candidates:
        if ".test." not in test_path.name:
            continue
        stem = test_path.name.split(".test.", maxsplit=1)[0]
        reports_dir = test_path.parent.parent / "reports"
        metadata_path = reports_dir / f"{stem}.metadata.json"
        if metadata_path.exists():
            return metadata_path

    test_path = candidates[-1]
    if ".test." not in test_path.name:
        return None
    stem = test_path.name.split(".test.", maxsplit=1)[0]
    return test_path.parent.parent / "reports" / f"{stem}.metadata.json"


def load_runner_metrics(results_root: Path) -> dict[str, dict[str, int | None]]:
    mapping: dict[str, dict[str, int | None]] = {}
    for summary_file in sorted(results_root.glob("runner_summary_*.json")):
        payload = read_json(summary_file)
        rows = payload.get("results", [])
        if not isinstance(rows, list):
            continue
        for row in rows:
            if not isinstance(row, dict):
                continue
            generated_test_path = row.get("generated_test_path")
            if not generated_test_path:
                continue
            name = Path(str(generated_test_path)).name
            mapping[name] = {
                "generation_latency_ms": to_int(row.get("latency_ms")),
                "prompt_tokens": to_int(row.get("prompt_tokens")),
                "completion_tokens": to_int(row.get("completion_tokens")),
                "total_tokens": to_int(row.get("total_tokens")),
            }
    return mapping
