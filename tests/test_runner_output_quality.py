from __future__ import annotations

from benchmark.runner.runner import ModelConfig, Runner


def test_sanitize_removes_think_block() -> None:
    runner = object.__new__(Runner)
    text = "<think>reasoning</think>\nimport pytest\n\n\ndef test_ok():\n    assert True\n"
    out = Runner._sanitize_model_output(runner, text, "python")
    assert "<think>" not in out
    assert out.startswith("import pytest")


def test_validate_generated_test_detects_python_syntax_error() -> None:
    runner = object.__new__(Runner)
    model = ModelConfig(
        model_name="m",
        provider="x",
        endpoint="https://x",
        model="m",
        api_key_env="X",
        parameters={"max_tokens": 4096},
        timeout_ms=1000,
    )

    reason = Runner._validate_generated_test(
        runner,
        code="def bad(:\n    pass\n",
        language="python",
        completion_tokens=4096,
        model_config=model,
    )
    assert reason is not None
    assert "syntax error" in reason


def test_validate_generated_test_rejects_think_tags() -> None:
    runner = object.__new__(Runner)
    model = ModelConfig(
        model_name="m",
        provider="x",
        endpoint="https://x",
        model="m",
        api_key_env="X",
        parameters={"max_tokens": 4096},
        timeout_ms=1000,
    )

    reason = Runner._validate_generated_test(
        runner,
        code="<think>...</think>",
        language="python",
        completion_tokens=100,
        model_config=model,
    )
    assert reason == "contains leaked reasoning tags"
