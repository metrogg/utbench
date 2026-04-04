import pytest
from typing import List
# Adjust import path to match your project structure
from src import has_close_elements


@pytest.mark.parametrize("numbers, threshold, expected", [
    # Docstring provided examples
    ([1.0, 2.0, 3.0], 0.5, False),
    ([1.0, 2.8, 3.0, 4.0, 5.0, 2.0], 0.3, True),
    # Boundary: small list sizes
    ([], 0.5, False),
    ([9.9], 1.0, False),
    # Boundary: exact matches
    ([5.0, 5.0], 0.001, True),
    # Boundary: distance equal to threshold
    ([1.0, 3.0], 2.0, False),
    # Boundary: distance just below threshold
    ([2.0, 3.9], 2.0, True),
    # Negative number handling
    ([-2.0, -2.4], 0.5, True),
    ([-5.0, 0.0, 5.0], 0.9, False),
    # Threshold edge values
    ([1.0, 1.0], 0.0, False),
    ([1.0, 2.0], 0.0, False),
    ([-1.0, 1.0], -0.5, False),
    # Non-adjacent close elements (validates all pairs are checked)
    ([1.0, 3.0, 5.0, 1.1], 0.2, True),
    # Duplicate elements spread across list
    ([2.0, 4.0, 6.0, 4.0], 0.1, True),
    # Large value handling
    ([1000000.0, 1000000.15], 0.2, True),
    # Integer input compatibility
    ([1, 3, 5, 6], 0.9, True),
    # Floating point precision test
    ([0.1 + 0.2, 0.3], 1e-15, True),
    ([0.1 + 0.2, 0.3], 1e-16, False),
    # Very large threshold
    ([1.0, 10.0, 100.0], 200.0, True),
])
def test_has_close_elements_various_scenarios(numbers: List[float], threshold: float, expected: bool):
    assert has_close_elements(numbers, threshold) == expected


def test_has_close_elements_non_numeric_elements_raises_type_error():
    with pytest.raises(TypeError):
        has_close_elements(["a", 1.0, "b"], 0.5)


def test_has_close_elements_invalid_threshold_type_raises_type_error():
    with pytest.raises(TypeError):
        has_close_elements([1.0, 2.0], "0.5")