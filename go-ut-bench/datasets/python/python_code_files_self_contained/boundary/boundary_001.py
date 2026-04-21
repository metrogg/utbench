def safe_divide(numerator: float, denominator: float) -> float:
    if denominator == 0:
        raise ZeroDivisionError("Cannot divide by zero")
    if numerator == 0:
        return 0.0
    result = numerator / denominator
    if result > 1e308:
        return float('inf')
    if result < -1e308:
        return float('-inf')
    return result