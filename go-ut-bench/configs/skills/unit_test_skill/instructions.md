# Unit Test Generation Skill

Generate a strong but compact test file. This skill should improve precision without taking over the agent's workflow.

## Core priorities

1. Prefer **correct, runnable tests** over broad suites.
2. Cover **happy path + key boundaries + one or two meaningful failure paths**.
3. Use only **public behavior visible in the source**.
4. If local compile/test feedback contradicts an assumption, **trust the runtime and revise**.
5. Avoid oversized parameter matrices, speculative property tests, and unnecessary helpers.

## Minimal workflow

Before writing the final file, check only:

1. What functions/types are actually callable?
2. Which branches or edge conditions are obvious from the implementation?
3. Which dependencies truly need a fake or stub?
4. What is the smallest set of tests that would catch real regressions?

## Design rules

- Use behavioral assertions, not placeholders.
- Add a test only when it checks a distinct behavior.
- Keep fixtures local and small.
- Avoid asserting implementation trivia that will make tests brittle.

## Mocking rules

- Mock external behavior, not the unit under test.
- Prefer seams already exposed by the source.
- Do not patch standard-library internals unless the source explicitly supports it.

## Language-specific guardrails

### Python
- Prefer plain pytest functions.
- Use `pytest.raises` for failure paths.
- Import every helper module used by the test file explicitly.

### Go
- Prefer small table-driven tests when they clarify the file.
- Keep package names and symbol usage consistent with the source.

### Java
- Use JUnit 5 Jupiter only.
- Match package declaration, class names, and static/instance usage exactly.

### C++
- Use straightforward GoogleTest `TEST`, `EXPECT_*`, and `ASSERT_*`.
- Do not use operators or methods that the source does not define.
- Include standard headers required by the generated test code.
- Prefer a smaller suite that compiles cleanly over a broad suite built on guesses.

## Final checklist

- The file is runnable as written.
- Assertions are specific.
- The suite is no larger than needed.
- No invented APIs, operators, globals, or hidden seams.
- If the agent did local verification, the final file matches those results.
