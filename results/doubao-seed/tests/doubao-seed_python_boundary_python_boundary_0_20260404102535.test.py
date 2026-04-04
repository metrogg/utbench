import pytest
from typing import List
# Assume below_zero is imported from the module where it's defined
from your_module import below_zero


@pytest.mark.parametrize(
    "operations, expected",
    [
        # Normal path cases
        ([1, 2, 3], False),
        ([1, 2, -4, 5], True),
        # Boundary cases
        ([], False),
        ([-1], True),
        ([3, 2, -6], True),
        ([2, -1, -1], False),
        ([5, -3, -2, 4], False),
        ([0, 0, 0], False),
        ([-5, 10, 20], True),
        ([10], False),
        ([1_000_000, -2_000_000], True),
        ([2_000_000, -1_000_000, -500_000], False),
        ([0, -1], True),
        ([5, -5, -1], True),
        ([5, -4, -1], False),
        ([2, -3, 100], True),
    ]
)
def test_below_zero_valid_inputs(operations: List[int], expected: bool) -> None:
    assert below_zero(operations) == expected


def test_below_zero_none_input() -> None:
    with pytest.raises(TypeError):
        below_zero(None)


def test_below_zero_non_iterable_input() -> None:
    with pytest.raises(TypeError):
        below_zero(12345)


def test_below_zero_non_numeric_elements() -> None:
    with pytest.raises(TypeError):
        below_zero([1, "2", -3])