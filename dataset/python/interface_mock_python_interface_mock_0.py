from __future__ import annotations

from dataclasses import dataclass
from itertools import permutations
from random import Random
from typing import Callable, Iterable, Sequence


class PermutationValidationError(ValueError):
    pass


@dataclass(frozen=True)
class PermutationMetric:
    permutation: tuple[int, ...]
    shuffled: tuple[int, ...]
    distance_sum: int


def _normalize_numbers(numbers: Iterable[object]) -> list[int]:
    normalized: list[int] = []
    for idx, raw in enumerate(numbers):
        if isinstance(raw, bool) or not isinstance(raw, int):
            raise PermutationValidationError(
                f"numbers[{idx}] must be int, got {type(raw).__name__}"
            )
        normalized.append(raw)
    if len(normalized) < 2:
        raise PermutationValidationError("numbers must contain at least two values")
    if len(normalized) > 8:
        raise PermutationValidationError("numbers must contain at most eight values")
    return normalized


def _distance_sum(values: Sequence[int]) -> int:
    total = 0
    for left, right in zip(values, values[1:]):
        total += abs(left - right)
    return total


def permutation_metrics(
    numbers: Iterable[object],
    *,
    rng: Random | None = None,
    shuffle_fn: Callable[[list[int]], None] | None = None,
) -> list[PermutationMetric]:
    values = _normalize_numbers(numbers)
    local_rng = rng or Random(0)
    metrics: list[PermutationMetric] = []

    for perm in permutations(values):
        shuffled_list = list(perm)
        if shuffle_fn is not None:
            shuffle_fn(shuffled_list)
        else:
            local_rng.shuffle(shuffled_list)
        metric = PermutationMetric(
            permutation=perm,
            shuffled=tuple(shuffled_list),
            distance_sum=_distance_sum(shuffled_list),
        )
        metrics.append(metric)

    return metrics


def task_func(numbers: Sequence[int] | None = None) -> float:
    """Average absolute-distance sum over shuffled permutations.

    The default input keeps the runtime small and deterministic.
    """

    input_numbers = list(range(1, 4)) if numbers is None else list(numbers)
    metrics = permutation_metrics(input_numbers, rng=Random(42))
    total = sum(metric.distance_sum for metric in metrics)
    return total / len(metrics)


def max_distance_permutation(numbers: Sequence[int]) -> PermutationMetric:
    metrics = permutation_metrics(numbers, rng=Random(7))
    return max(metrics, key=lambda item: item.distance_sum)
