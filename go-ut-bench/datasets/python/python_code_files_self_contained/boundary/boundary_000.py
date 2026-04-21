def clamp(value: int, min_val: int, max_val: int) -> int:
    if min_val > max_val:
        raise ValueError("min_val must be less than or equal to max_val")
    if value < min_val:
        return min_val
    if value > max_val:
        return max_val
    return value