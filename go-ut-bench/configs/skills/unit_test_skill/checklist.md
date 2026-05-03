# Pre-flight Checklist

Before outputting the test file, verify each item:

## Structure
- [ ] One test file, self-contained
- [ ] All imports at the top
- [ ] No `if __name__ == "__main__"` block
- [ ] No print statements (use assert messages instead)

## Assertions
- [ ] Every test has at least one meaningful assertion
- [ ] No placeholder assertions (`assert True`, `assert result is not None`)
- [ ] Error paths assert specific exception types
- [ ] Numeric assertions use tolerance for floats: `assert abs(a - b) < 1e-9`

## Mocking
- [ ] External I/O is mocked (network, filesystem, subprocess)
- [ ] Time-dependent code is mocked
- [ ] Random-dependent code is seeded or mocked
- [ ] Mock assertions verify call arguments

## Isolation
- [ ] No test depends on another test's state
- [ ] No global side effects
- [ ] Temp files use `tempfile` and are cleaned up
- [ ] No hardcoded network ports or hostnames

## Completeness
- [ ] Happy path covered
- [ ] At least one boundary case
- [ ] At least one error/invalid input case
- [ ] All source code branches exercised
