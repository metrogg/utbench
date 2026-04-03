# AGENTS.md - Working Guide for Agentic Contributors

## Collaboration Rules (必须遵守)

**当信息不全或不确定时，不要强行判断执行，必须先向用户询问确认。**

以下情况必须暂停并询问：
- 需求有多种理解方式，无法确定用户意图
- 需要修改的文件/位置不明确
- 存在多个可行方案，需要用户选择
- 发现用户的设计可能存在明显问题
- 缺少关键上下文（如具体的错误信息、期望的输出格式等）

询问时应：
1. 说明你的理解
2. 列出不确定的地方
3. 提供可选方案及各自的利弊
4. 请用户确认后再执行

## 1) Purpose and Scope

This file defines practical rules for coding agents working in `ut-bench`.
Use it as the default execution guide for planning, editing, testing, and docs.

Current repository state (master branch) is scaffold-first:
- Core scripts exist but are placeholders.
- Most design/research docs are TODO stubs.
- Benchmark runtime modules are not implemented yet.

Because of this, prefer accurate status reporting over pretending features exist.

## 2) Collaboration Rule (Critical)

When information is incomplete or ambiguous, do not force a decision.
Ask the user for confirmation before executing irreversible or assumption-heavy changes.

Ask when:
- Requirements can be interpreted in multiple valid ways.
- Target files/paths are unclear.
- Multiple implementation choices have different tradeoffs.
- A command may delete or overwrite user work.

When asking, include:
1. Your understanding.
2. What is uncertain.
3. Recommended default and alternatives.
4. What changes based on the choice.

## 3) Source of Truth and Rules

- Primary project intent: `README.md`.
- Design docs: `docs/design/*.md`.
- Script behavior: `scripts/*.sh`.
- Model config shape: `benchmark/config/models.yaml`.

Rule-file scan result:
- No `.cursorrules` found.
- No `.cursor/rules/` found.
- No `.github/copilot-instructions.md` found.

If these files are added later, merge their rules into this guide.

## 4) Repository Map (Current)

- `benchmark/`
  - `README.md` (module-level description)
  - `config/models.yaml` (placeholder config)
- `dataset/`
  - `README.md` (dataset description placeholder)
- `docs/`
  - `design/` (benchmark, metrics, dataset docs; currently TODO-heavy)
  - `research/` (model/tools survey TODO stubs)
- `scripts/`
  - `setup.sh` (placeholder)
  - `run_benchmark.sh` (placeholder)
  - `gen_report.sh` (placeholder)

## 5) Build / Lint / Test Commands

Use these commands from repo root.

### Project Scripts (Current, Placeholder)

```bash
bash scripts/setup.sh
bash scripts/run_benchmark.sh
bash scripts/gen_report.sh
```

These scripts currently print TODO messages and exit successfully.
Treat them as scaffolding hooks, not production entry points.

### Shell Script Validation

```bash
bash -n scripts/setup.sh
bash -n scripts/run_benchmark.sh
bash -n scripts/gen_report.sh
```

### Python Syntax Check (per file)

```bash
python -m py_compile path/to/file.py
```

### YAML Sanity Check

```bash
python -c "import pathlib,yaml; yaml.safe_load(pathlib.Path('benchmark/config/models.yaml').read_text(encoding='utf-8'))"
```

### Pytest Commands (Use When Tests Exist)

There is no stable root test suite on this branch yet.
When tests are added, use:

```bash
pytest
pytest path/to/test_file.py
pytest path/to/test_file.py::TestClass::test_case
pytest -k "keyword" -q
```

Single-test execution standard for agents:
- Prefer node id form: `pytest file.py::Class::test_name`.
- Use `-q` for concise logs in iterative debugging.

## 6) Coding Style Guidelines

These rules apply to all new code added to this repo.

### Imports

- Order imports: standard library, third-party, local modules.
- Keep imports explicit and minimal.
- Prefer absolute imports within project packages.
- Avoid `from x import *`.

### Formatting

- Follow PEP 8 for Python.
- Use 4 spaces for indentation.
- Keep lines reasonably short (target <= 100 chars).
- Keep functions small and single-purpose.

### Types

- Add type hints for all new public functions and methods.
- Add return type annotations consistently.
- Use concrete built-ins (`list[str]`, `dict[str, Any]`) on modern Python.

### Naming

- Functions/variables: `snake_case`.
- Classes: `PascalCase`.
- Constants: `UPPER_SNAKE_CASE`.
- Test files: `test_<module>.py`.

### Error Handling

- Raise specific exceptions with actionable messages.
- Do not swallow errors with bare `except:`.
- Use `except Exception as exc` only when re-raising with context.
- Validate inputs early and fail fast.

### File and Path Handling

- Use `pathlib.Path` over raw string path concatenation.
- Always specify file encoding for text I/O (`utf-8`).
- Keep generated outputs under `results/` unless user asks otherwise.

### Subprocess and External Commands

- Use explicit timeouts.
- Capture stdout/stderr for diagnostics.
- Check exit codes and surface failures clearly.

### Logging and Output

- Prefer structured, concise logs.
- Use `print` for CLI scripts only.
- Include enough context to debug failed samples quickly.

## 7) Test and Evaluation Conventions

- Keep deterministic tests (no random outcomes without fixed seeds).
- Isolate filesystem/network side effects with temp dirs or mocks.
- For benchmark metrics, store raw measurement fields first.
- Do not compute or report weighted total score unless explicitly required.

## 8) Documentation Conventions

- Update docs when behavior or interfaces change.
- Keep design docs aligned with implementation status.
- Mark placeholders explicitly as TODO instead of inventing details.

## 9) Git Hygiene

- Never revert user-owned unrelated changes.
- Keep commits focused by concern (docs vs scripts vs runtime code).
- Avoid destructive git operations unless user explicitly requests.
- Before major edits, inspect `git status` and report risky states.

## 10) Agent Execution Checklist

Before changes:
- Confirm task scope and target files.
- Read relevant docs/scripts first.

During changes:
- Apply minimal, reversible edits.
- Keep behavior and docs in sync.

After changes:
- Run the smallest meaningful validation command.
- Report what was changed, what was verified, and what remains TODO.
