from __future__ import annotations

from benchmark.evaluator.pipeline import Evaluator


def test_evaluator_pipeline_runs_and_returns_summary() -> None:
    evaluator = Evaluator(results_root="results")
    payload = evaluator.run(models=["deepseek"], languages=["python"])
    assert "summary" in payload
    assert "results" in payload
    assert payload["summary"]["total_samples"] >= 0
    assert payload["filters"]["latest_only"] is True
