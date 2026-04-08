from __future__ import annotations

from typing import Any


def build_failure_breakdown(rows: list[dict[str, Any]], top_n: int) -> list[dict[str, Any]]:
    grouped: dict[tuple[str, str], dict[str, Any]] = {}

    for row in rows:
        stage, message = pick_error(row)
        if stage is None or message is None:
            continue
        error_type = classify_error(message)
        key = (stage, error_type)
        item = grouped.setdefault(
            key,
            {
                "stage": stage,
                "error_type": error_type,
                "count": 0,
                "example_sample": row.get("sample_id"),
                "example_model": row.get("model"),
                "example_message": short_message(message),
            },
        )
        item["count"] += 1

    ranked = sorted(grouped.values(), key=lambda item: (-int(item["count"]), item["stage"]))
    return ranked[: max(top_n, 0)]


def pick_error(row: dict[str, Any]) -> tuple[str | None, str | None]:
    for stage, field in (
        ("compile", "compile_error"),
        ("test", "test_error"),
        ("coverage", "coverage_error"),
        ("mutation", "mutation_error"),
    ):
        value = row.get(field)
        if isinstance(value, str) and value.strip():
            return stage, value
    return None, None


def classify_error(message: str) -> str:
    text = message.lower()
    if "no module named" in text or "modulenotfounderror" in text:
        return "module_not_found"
    if "nameerror" in text:
        return "name_error"
    if "assertionerror" in text:
        return "assertion_failure"
    if "importerror" in text:
        return "import_error"
    if "syntaxerror" in text:
        return "syntax_error"
    if "timeout" in text:
        return "timeout"
    if "stopiteration" in text:
        return "stop_iteration"
    if "recursionerror" in text:
        return "recursion_error"
    return "other"


def short_message(message: str, max_len: int = 200) -> str:
    one_line = " ".join(part.strip() for part in message.splitlines() if part.strip())
    if len(one_line) <= max_len:
        return one_line
    return one_line[: max_len - 3] + "..."
