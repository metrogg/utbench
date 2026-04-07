from __future__ import annotations

from dataclasses import dataclass
from math import isfinite
from typing import Iterable, Sequence


class NumberValidationError(ValueError):
    pass


@dataclass(frozen=True)
class PairDistance:
    left_index: int
    right_index: int
    left_value: float
    right_value: float
    distance: float


def _coerce_number(value: object, *, index: int) -> float:
    if isinstance(value, bool):
        raise NumberValidationError(f"numbers[{index}] must be numeric, got bool")
    if not isinstance(value, (int, float)):
        raise NumberValidationError(
            f"numbers[{index}] must be numeric, got {type(value).__name__}"
        )
    number = float(value)
    if not isfinite(number):
        raise NumberValidationError(f"numbers[{index}] must be finite")
    return number


def _normalize_numbers(numbers: Iterable[object]) -> list[float]:
    normalized: list[float] = []
    for idx, raw in enumerate(numbers):
        normalized.append(_coerce_number(raw, index=idx))
    return normalized


def nearest_pair(numbers: Sequence[float]) -> PairDistance | None:
    if len(numbers) < 2:
        return None
    indexed = sorted(enumerate(numbers), key=lambda item: item[1])
    best: PairDistance | None = None

    for (left_i, left_v), (right_i, right_v) in zip(indexed, indexed[1:]):
        gap = abs(right_v - left_v)
        candidate = PairDistance(
            left_index=min(left_i, right_i),
            right_index=max(left_i, right_i),
            left_value=left_v,
            right_value=right_v,
            distance=gap,
        )
        if best is None or candidate.distance < best.distance:
            best = candidate

    return best


def has_close_elements(numbers: Sequence[float], threshold: float) -> bool:
    """Check whether any pair of numbers is closer than threshold.

    >>> has_close_elements([1.0, 2.0, 3.0], 0.5)
    False
    >>> has_close_elements([1.0, 2.8, 3.0, 4.0, 5.0, 2.0], 0.3)
    True
    """

    if not isinstance(threshold, (int, float)) or isinstance(threshold, bool):
        raise NumberValidationError("threshold must be numeric")
    threshold_value = float(threshold)
    if threshold_value < 0:
        return False

    normalized = _normalize_numbers(numbers)
    pair = nearest_pair(normalized)
    if pair is None:
        return False
    return pair.distance < threshold_value


def count_close_pairs(numbers: Sequence[float], threshold: float) -> int:
    if threshold < 0:
        return 0
    normalized = _normalize_numbers(numbers)
    sorted_values = sorted(normalized)
    count = 0
    left = 0
    for right in range(len(sorted_values)):
        while left < right and sorted_values[right] - sorted_values[left] >= threshold:
            left += 1
        count += right - left
    return count


def min_distance(numbers: Sequence[float]) -> float | None:
    normalized = _normalize_numbers(numbers)
    pair = nearest_pair(normalized)
    if pair is None:
        return None
    return pair.distance
