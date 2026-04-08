from __future__ import annotations

from dataclasses import dataclass
from random import Random
from statistics import mean, pstdev
from typing import Iterable


class DatasetBuilderError(ValueError):
    pass


@dataclass(frozen=True)
class LetterStats:
    letter: str
    values: tuple[int, ...]
    avg: float
    spread: float


def _normalize_letters(letters: Iterable[object]) -> list[str]:
    normalized: list[str] = []
    seen: set[str] = set()
    for idx, raw in enumerate(letters):
        if not isinstance(raw, str) or len(raw) != 1 or not raw.isalpha():
            raise DatasetBuilderError(
                f"letters[{idx}] must be single alphabetic character"
            )
        upper = raw.upper()
        if upper in seen:
            raise DatasetBuilderError(f"duplicated letter: {upper}")
        seen.add(upper)
        normalized.append(upper)
    if not normalized:
        raise DatasetBuilderError("letters must not be empty")
    if len(normalized) > 26:
        raise DatasetBuilderError("letters too many")
    return normalized


def build_letter_values(
    letters: Iterable[object],
    *,
    seed: int = 11,
    min_size: int = 2,
    max_size: int = 6,
) -> dict[str, list[int]]:
    if min_size < 1 or max_size < min_size:
        raise DatasetBuilderError("invalid list-size bounds")
    normalized = _normalize_letters(letters)
    rng = Random(seed)

    result: dict[str, list[int]] = {}
    for letter in normalized:
        size = rng.randint(min_size, max_size)
        result[letter] = [rng.randint(0, 100) for _ in range(size)]
    return result


def summarize_letter_values(data: dict[str, list[int]]) -> list[LetterStats]:
    items: list[LetterStats] = []
    for letter, values in data.items():
        if not values:
            raise DatasetBuilderError(f"empty list for letter {letter}")
        avg = float(mean(values))
        spread = float(pstdev(values)) if len(values) > 1 else 0.0
        items.append(
            LetterStats(
                letter=letter,
                values=tuple(values),
                avg=avg,
                spread=spread,
            )
        )
    items.sort(key=lambda item: (item.avg, -item.spread, item.letter), reverse=True)
    return items


def task_func(letters: Iterable[str]) -> dict[str, list[int]]:
    """Return a dict sorted by descending mean value."""

    raw = build_letter_values(letters)
    stats = summarize_letter_values(raw)
    ordered: dict[str, list[int]] = {}
    for item in stats:
        ordered[item.letter] = list(item.values)
    return ordered


def top_letter(letters: Iterable[str]) -> tuple[str, float]:
    ordered = task_func(letters)
    first_letter = next(iter(ordered))
    return first_letter, float(mean(ordered[first_letter]))
