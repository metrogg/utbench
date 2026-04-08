from __future__ import annotations

from typing import Any

from .contracts import (
    _COLUMN_LABELS_ZH,
    _ERROR_TYPE_LABELS_ZH,
    _GEN_METRICS_REASON_LABELS_ZH,
    _GEN_METRICS_SOURCE_LABELS_ZH,
    _STAGE_LABELS_ZH,
)


def to_int(value: Any) -> int | None:
    if isinstance(value, bool):
        return None
    if isinstance(value, int):
        return value
    if isinstance(value, float):
        return int(value)
    if isinstance(value, str):
        text = value.strip()
        if text and text.lstrip("-").isdigit():
            return int(text)
    return None


def to_float(value: Any) -> float | None:
    if isinstance(value, bool):
        return None
    if isinstance(value, (int, float)):
        return float(value)
    if isinstance(value, str):
        text = value.strip()
        if not text:
            return None
        try:
            return float(text)
        except ValueError:
            return None
    return None


def numeric_values(rows: list[dict[str, Any]], key: str) -> list[float]:
    values: list[float] = []
    for row in rows:
        value = row.get(key)
        if isinstance(value, bool):
            continue
        if isinstance(value, (int, float)):
            values.append(float(value))
    return values


def avg(values: list[float]) -> float | None:
    if not values:
        return None
    return round(sum(values) / len(values), 6)


def rate(numerator: int, denominator: int) -> float:
    if denominator <= 0:
        return 0.0
    return round(numerator / denominator, 6)


def column_title(column: str) -> str:
    return _COLUMN_LABELS_ZH.get(column, column)


def format_pct(value: Any) -> str:
    number = to_float(value)
    if number is None:
        return "N/A"
    return f"{number * 100:.2f}%"


def format_num(value: Any) -> str:
    number = to_float(value)
    if number is None:
        return "N/A"
    if abs(number - int(number)) < 1e-9:
        return str(int(number))
    return f"{number:.6f}".rstrip("0").rstrip(".")


def display_value(key: str, value: Any) -> str:
    if value is None:
        return "N/A"
    if isinstance(value, bool):
        return "鏄?" if value else "鍚?"
    if key == "stage":
        return _STAGE_LABELS_ZH.get(str(value), str(value))
    if key == "error_type":
        return _ERROR_TYPE_LABELS_ZH.get(str(value), str(value))
    if key == "generation_metrics_source":
        return _GEN_METRICS_SOURCE_LABELS_ZH.get(str(value), str(value))
    if key == "generation_metrics_reason":
        return _GEN_METRICS_REASON_LABELS_ZH.get(str(value), str(value))
    if isinstance(value, (int, float)):
        if key.endswith("_rate") or "coverage" in key or key.endswith("_score"):
            return format_pct(float(value))
        return format_num(float(value))
    return str(value)
