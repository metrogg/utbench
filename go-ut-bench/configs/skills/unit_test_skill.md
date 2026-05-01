Generate a production-grade unit test file for the target sample.

Requirements:
- Prefer deterministic tests over snapshot-style assertions.
- Cover normal path, edge cases, and at least one failure or invalid-input path when the target behavior supports it.
- Keep the test file self-contained inside the prepared workspace.
- Do not modify the source implementation unless the prompt explicitly asks for it.
- Use the language's standard test conventions and file naming.
- Avoid placeholder assertions such as `assert True` unless they are part of a setup smoke test.
