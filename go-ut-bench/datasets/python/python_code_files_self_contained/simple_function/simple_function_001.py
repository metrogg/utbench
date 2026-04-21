def reverse_string(s: str) -> str:
    if s is None:
        raise ValueError("Input string cannot be None")
    return s[::-1]