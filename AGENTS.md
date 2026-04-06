# AGENTS.md

## 1) Purpose

This file is the operating guide for agentic coding tools working in `ut-bench`.
Follow these rules for implementation, validation, and reporting.

Project focus (current branch state):
- Runner is implemented and actively used.
- Evaluator exists and is being extended.
- Reporter is not fully implemented.

Primary objective:
- Evaluate multi-model unit-test generation quality across languages.

## 2) Repository Snapshot

Important paths:
- `benchmark/runner/` runner pipeline and CLI
- `benchmark/evaluator/` evaluator pipeline and CLI
- `benchmark/config/models.yaml` model and benchmark configuration
- `dataset/` benchmark source samples by language/complexity
- `results/` generated tests, artifacts, summaries, checkpoints

Current CLIs:
- `python -m benchmark.runner`
- `python -m benchmark.evaluator`

## 3) Rule Files Scan (Cursor/Copilot)

Scanned locations and status:
- `.cursorrules`: not found
- `.cursor/rules/`: not found
- `.github/copilot-instructions.md`: not found

If these files are added later:
- Read them before coding.
- Merge constraints into this document.
- If conflicts exist, follow the stricter rule.

## 4) Build / Lint / Test Commands

Run commands from repository root.

### 4.1 Environment

```bash
python --version
python -m pip install pyyaml
```

If evaluator coverage is used:

```bash
python -m pip install pytest coverage
```

### 4.2 Runner commands

Full run on enabled models:

```bash
python -m benchmark.runner
```

Model/language scoped run:

```bash
python -m benchmark.runner --model deepseek --lang python --max-samples 2
```

Dry run (no real API call):

```bash
python -m benchmark.runner --model deepseek --lang python --max-samples 1 --dry-run
```

### 4.3 Evaluator commands

Evaluate all discovered outputs:

```bash
python -m benchmark.evaluator --results-root results
```

Model/language scoped evaluation:

```bash
python -m benchmark.evaluator --results-root results --model deepseek --lang python
```

### 4.4 Test commands

Run all tests:

```bash
python -m pytest -q
```

Run a single test file:

```bash
python -m pytest tests/test_evaluator_pipeline.py -q
```

Run a single test case (preferred during iteration):

```bash
python -m pytest tests/test_evaluator_pipeline.py::test_evaluator_pipeline_runs_and_returns_summary -q
```

Keyword filter:

```bash
python -m pytest -k coverage -q
```

## 5) Configuration Rules

Model config is dict-based, keyed by model name in `models.yaml`.

Required model fields:
- `enabled`
- `provider`
- `config.api_endpoint`
- `config.model`
- `config.api_key_env`

Secrets handling:
- Store API keys in environment variables only.
- Do not commit plaintext keys.

PowerShell example:

```powershell
$env:DEEPSEEK_API_KEY="..."
$env:DASHSCOPE_API_KEY="..."
$env:MINIMAX_API_KEY="..."
$env:VOLCENGINE_API_KEY="..."
```

## 6) Code Style Guidelines

### 6.1 Imports
- Group imports: stdlib, third-party, local.
- Avoid wildcard imports.
- Remove unused imports.

### 6.2 Formatting
- Follow PEP 8.
- Use 4-space indentation.
- Keep lines near <= 100 chars when practical.
- Prefer small single-purpose functions.

### 6.3 Types
- Add type hints for new/changed public functions.
- Use explicit return types for non-trivial functions.
- Prefer modern generics (`list[str]`, `dict[str, Any]`).

### 6.4 Naming
- `snake_case` for functions/variables/modules.
- `PascalCase` for classes.
- `UPPER_SNAKE_CASE` for constants.
- Test files should start with `test_`.

### 6.5 Error handling
- Fail fast on invalid inputs.
- Raise specific exceptions.
- Never use bare `except:`.
- Preserve actionable context in error messages.

### 6.6 Filesystem and subprocess
- Use `pathlib.Path` for paths.
- Use `encoding="utf-8"` for text I/O.
- Use explicit subprocess timeouts.
- Capture stdout/stderr for diagnostics.

### 6.7 Logging
- Keep logs concise and structured.
- Use `logging` for pipeline modules.
- Avoid noisy debug output in committed code.

### 6.8 Comments and docs
- Add short comments/docstrings for non-obvious logic.
- Keep README and user-facing docs in Chinese unless requested otherwise.

## 7) Runner/Evaluator Contract Notes

Runner outputs are stored under `results/<model>/...`.
Evaluator should consume generated tests from `results/<model>/tests/`.

When extending evaluator:
- Keep fields stable in summary JSON.
- Preserve `None` for unavailable metrics instead of fake zeros.
- Keep mutation testing optional unless environment is ready.

## 8) Git Hygiene

- Do not revert unrelated user changes.
- Avoid destructive git commands unless explicitly requested.
- Keep changes scoped (runner vs evaluator vs docs).
- Do not commit secrets or local credentials.

## 9) Agent Execution Checklist

Before coding:
- Read target module and related config/docs.
- Confirm command(s) to validate changes.

During coding:
- Make minimal reversible edits.
- Keep behavior and docs aligned.

After coding:
- Run smallest meaningful tests first, then broader suite.
- Report changed files, validation commands, and residual TODOs.
