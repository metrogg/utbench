# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Layout

The primary codebase is `go-ut-bench/` (Go CLI). Root `README.md`, `docker-compose.yml`, and `docker.sh` are **obsolete** — ignore them and work from `go-ut-bench/`.

**Important**: there is an unintegrated duplicate of `internal/runner/` at the repository root (`runner/api.go`, `runner/prompt.go`, etc.). This was added in a recent commit but is **outside the Go module** and not imported by the build. Always edit files under `go-ut-bench/internal/runner/` instead.

## Build and Run

All commands run from `go-ut-bench/`:

### Build

```bash
go build -o utbench ./cmd/utbench/
```

### Tests

```bash
go test ./internal/...
# Run a specific test
go test ./internal/runner/... -run TestCheckpoint
go test -v ./internal/evaluator/... -run TestPythonEval
```

### Quick dry-run (no API calls)

```bash
./utbench run --models deepseek --langs python --max-samples 2 --dry-run
```

### Full pipeline

```bash
./utbench run --models deepseek,qwen --langs python,go --max-samples 5 --mutation-enabled
```

### Step-by-step

```bash
# Generate tests only
./utbench generate --models deepseek --langs python --max-samples 5

# Evaluate an existing manifest
./utbench evaluate --manifest ./artifacts/runs/<run-id>/generated/generated_manifest.json

# Generate report
./utbench report --evaluation ./artifacts/runs/<run-id>/evaluation/evaluation_result.json

# Ingest into SQLite
./utbench ingest --evaluation ./artifacts/runs/<run-id>/evaluation/evaluation_result.json --db-path ./storage/utbench.db
```

### Docker (recommended when running evaluation toolchains)

```bash
cd go-ut-bench
docker build -t utbench:latest .

# Linux/macOS
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest run \
    --models deepseek \
    --langs python \
    --config /app/configs/models.yaml \
    --max-samples 2

# Windows PowerShell
docker run --rm --env-file .env `
  -v "${PWD}/datasets:/app/datasets" `
  -v "${PWD}/artifacts:/app/artifacts" `
  -v "${PWD}/configs:/app/configs" `
  utbench:latest run --models deepseek --langs python --max-samples 2
```

### Helper scripts

Pre-built convenience wrappers for Docker runs:

```bash
# Linux / WSL
./run_bench.sh [models] [langs] [max-samples] [mutation]

# Windows PowerShell
.\run_bench.ps1 [models] [langs] [max-samples] [mutation]
```

### Dataset management

```bash
# Index all samples
./utbench dataset index --dataset-root ./datasets --output ./configs/dataset_index.json

# Build filtered manifest
./utbench dataset manifest \
  --index ./configs/dataset_index.json \
  --langs python,go \
  --class module_level \
  --level l1 \
  --limit-per-scenario 20 \
  --output ./configs/dataset_l1.json
```

## Architecture

Five-stage pipeline orchestrated by `internal/orchestrator/service.go`:

1. **dataset** (`internal/dataset/`) — discovers samples; infers language/class/scenario from directory structure; computes MD5 hashes; applies filters
2. **runner** (`internal/runner/`) — worker pool calling LLM APIs concurrently; checkpoint-based incremental execution; retry with exponential backoff; truncation detection and auto-continuation (max 3 retries when `finish_reason: "length"` or incomplete code blocks)
3. **evaluator** (`internal/evaluator/`) — language-specific compile → test → coverage → mutation in isolated temp dirs; Python `self_contained` and `module_level` both supported
4. **reporter** (`internal/reporter/`) — multi-dimensional aggregation (by model, language, scenario); composite score = 0.3×compile + 0.3×pass_rate + 0.2×coverage + 0.2×mutation; HTML report via Chart.js CDN; mutation breakdown (total/killed/survived/no_tests/timeouts/skipped/suspicious); truncation statistics with tuning recommendations
5. **store** (`internal/store/`) — SQLite upsert on `(run_id, model, language, sample_id)`

### Data contracts

`internal/contracts/` is the single source of truth for all cross-stage types. Key files:
- `spec.go` — `RunSpec`, `SampleRef`, `ModuleLevelMeta`
- `constants.go` — `SchemaVersion = "v0.1.0"`, `DatasetClass`, `RunMode`
- `results.go` — `GeneratedManifest`, `EvaluationResultSet`, `ReportPayload`

JSON read/write helpers live in `internal/contracts/`.

### Checkpoint mechanism

Runner in incremental mode hashes `(dataset_root + classes + level + manifest + max_samples + models)` → SHA1 → `artifacts/checkpoints/runner_<hash>.checkpoint.json`. Key per completed task: `<model>|<language>|<sample_id>`. Any change to the scoped parameters invalidates the checkpoint.

### Model config

`configs/models.yaml` — provider/endpoint/api_key_env per model. Default code path is `../benchmark/config/models.yaml` (relative to working dir); always pass `--config ./configs/models.yaml` locally or `--config /app/configs/models.yaml` in Docker.

## Common Pitfalls

- **Default `--class` is `self_contained`**: Python and Go datasets are entirely `module_level`. Use `--class module_level` for those languages.
- **Module-level samples** require `meta.json` with `workspace_root` and `module_import`; the evaluator does not clean up their workspaces (reused in-place).
- **Checkpoint invalidation**: changing any of models, langs, class, level, manifest, max-samples, or dataset-root changes the hash and starts a fresh run.
- **Mutation testing tools**: Python uses `mutmut`, Go uses `gremlins` (`go install github.com/go-gremlins/gremlins/cmd/gremlins@latest`). Windows mutation testing for Python is validated on Linux only; use Docker on Windows.
- **Line endings on Windows**: normalize with `git add --renormalize .`

## Code Style

- `gofmt` formatting; package names short lowercase
- Explicit `if err != nil` returns; panic recovery only in evaluator/runner worker goroutines
- Use `internal/obs.Logger` (slog wrapper); pass `logger` into services rather than using a global
- Use `nil` pointers for missing metrics in JSON output — do not substitute zeros for absent data
- Use `filepath.Join` for all path construction
- File permissions: `0o755` for directories, `0o644` for files

## Environment Setup

Copy `.env.example` (at repo root) or create `go-ut-bench/.env` and fill in:

| Variable | Provider |
|---|---|
| `DEEPSEEK_API_KEY` | DeepSeek |
| `DASHSCOPE_API_KEY` | Qwen (Dashscope) |
| `MINIMAX_API_KEY` | MiniMax |
| `VOLCENGINE_API_KEY` | doubao-seed (original) |
| `ARK_API_KEY` | doubao-seed-2.0-lite/1.6/2.0-pro-v2, glm-4.7 |

# Prompt System

The runner uses three prompt modes (`internal/runner/prompt.go`):
- `full_file` — default for self-contained samples; full source + instructions
- `completion` — for continuation after truncation
- `module_level` — for samples with workspace context and module imports

Prompt strategy: `structured-v1`. System message emphasizes runnable tests only, no explanations or placeholders. Use `BuildPromptCatalog()` to inspect templates.

# Key Files to Read First

- `go-ut-bench/cmd/utbench/main.go` — CLI entry point, command routing
- `go-ut-bench/internal/orchestrator/service.go` — pipeline orchestration logic
- `go-ut-bench/internal/contracts/spec.go` — core data structures (RunSpec, SampleRef)
- `go-ut-bench/internal/contracts/constants.go` — SchemaVersion, DatasetClass, RunMode
- `go-ut-bench/internal/runner/prompt.go` — prompt construction
- `go-ut-bench/internal/evaluator/service.go` — evaluation pipeline per language
