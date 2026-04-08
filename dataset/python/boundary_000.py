from __future__ import annotations

from dataclasses import dataclass
from typing import Iterable, Sequence


class LedgerValidationError(ValueError):
    pass


@dataclass(frozen=True)
class LedgerSnapshot:
    index: int
    operation: int
    balance: int


def _coerce_operation(value: object, *, index: int) -> int:
    if isinstance(value, bool):
        raise LedgerValidationError(f"operation[{index}] must be int, got bool")
    if not isinstance(value, int):
        raise LedgerValidationError(
            f"operation[{index}] must be int, got {type(value).__name__}"
        )
    return value


def _normalize_operations(operations: Iterable[object]) -> list[int]:
    normalized: list[int] = []
    for idx, raw in enumerate(operations):
        op = _coerce_operation(raw, index=idx)
        normalized.append(op)
    if not normalized:
        return []
    if len(normalized) > 10_000:
        raise LedgerValidationError("too many operations")
    return normalized


def running_balance(operations: Iterable[object], *, initial_balance: int = 0) -> list[LedgerSnapshot]:
    if not isinstance(initial_balance, int) or isinstance(initial_balance, bool):
        raise LedgerValidationError("initial_balance must be int")
    normalized = _normalize_operations(operations)
    snapshots: list[LedgerSnapshot] = []
    balance = initial_balance
    for idx, op in enumerate(normalized):
        balance += op
        snapshots.append(LedgerSnapshot(index=idx, operation=op, balance=balance))
    return snapshots


def first_negative_index(
    operations: Iterable[object],
    *,
    initial_balance: int = 0,
) -> int | None:
    for snapshot in running_balance(operations, initial_balance=initial_balance):
        if snapshot.balance < 0:
            return snapshot.index
    return None


def below_zero(operations: Sequence[int]) -> bool:
    """Return True when balance goes below zero at any step.

    The account starts at zero. Positive numbers are deposits and negative numbers
    are withdrawals.

    >>> below_zero([1, 2, 3])
    False
    >>> below_zero([1, 2, -4, 5])
    True
    """

    return first_negative_index(operations, initial_balance=0) is not None


def minimum_balance(operations: Sequence[int], *, initial_balance: int = 0) -> int:
    min_balance = initial_balance
    current = initial_balance
    for op in _normalize_operations(operations):
        current += op
        if current < min_balance:
            min_balance = current
    return min_balance


def overdraft_count(operations: Sequence[int], *, initial_balance: int = 0) -> int:
    count = 0
    for snapshot in running_balance(operations, initial_balance=initial_balance):
        if snapshot.balance < 0:
            count += 1
    return count
