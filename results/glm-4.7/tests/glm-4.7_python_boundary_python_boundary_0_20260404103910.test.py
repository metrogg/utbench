import pytest
from typing import List

# Source code to be tested (included for runnability)
def below_zero(operations: List[int]) -> bool:
    balance = 0
    for op in operations:
        balance += op
        if balance < 0:
            return True
    return False


@pytest.mark.parametrize("operations, expected", [
    # --- Normal Paths ---
    # All positive numbers
    ([1, 2, 3], False),
    # Mixed numbers, balance never drops below zero
    ([10, -5, -4], False),
    # Mixed numbers, balance hits exactly zero
    ([5, -5], False),
    # Balance recovers but never goes negative
    ([10, -5, 1, -5, 1], False),

    # --- Boundary Conditions ---
    # Empty list
    ([], False),
    # List with a single zero
    ([0], False),
    # List with a single positive integer
    ([1], False),
    # List with a single negative integer (immediate failure)
    ([-1], True),
    # Balance hits zero multiple times
    ([1, -1, 1, -1], False),

    # --- Negative Paths (Balance drops below zero) ---
    # Drops below zero at the end
    ([10, -5, -6], True),
    # Drops below zero in the middle
    ([1, 2, -4, 5], True),
    # Drops below zero immediately
    ([-5, 10, 20], True),
    # Large numbers
    ([10**9, -10**9 - 1], True),
])
def test_below_zero_scenarios(operations: List[int], expected: bool) -> None:
    assert below_zero(operations) == expected