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

### Environment self-check

```bash
# Check that all evaluation toolchains are working
./utbench doctor --langs python,go,java,cpp --mutation-enabled

# Validate dataset readiness
./utbench dataset validate --dataset-root ./datasets --langs python,go,java,cpp --class self_contained --strict
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

# Ingest into SQLite (replaces old `utbench ingest`)
./utbench db ingest-evaluation --evaluation ./artifacts/runs/<run-id>/evaluation/evaluation_result.json --db-path ./storage/utbench.db
```

### Docker (recommended when running evaluation toolchains)

```bash
cd go-ut-bench

# Build all images (one command)
./build.sh all          # Linux/WSL
.\build.ps1 all         # Windows PowerShell
./build.sh all --cn     # 国内网络加速

# Or build just the eval image
./build.sh eval

# Verify image status
./build.sh verify
```

Run evaluation in Docker:

```bash
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

### Web management UI

```bash
# Launch on localhost:8080 (supports Docker and local runs)
./utbench web --addr :8080 --config ./configs/models.yaml

# Custom database and Docker image
./utbench web --addr :8080 --db-path ./storage/utbench.db --docker-image utbench:latest
```

### Database management

```bash
./utbench db init --db-path ./storage/utbench.db
./utbench db ingest-run --run-id <run-id> --output-root ./artifacts --db-path ./storage/utbench.db
./utbench db overview --db-path ./storage/utbench.db
./utbench db list-results --run-id <run-id> --db-path ./storage/utbench.db
./utbench db report --run-ids <run-a>,<run-b> --models deepseek,qwen --langs python,go --db-path ./storage/utbench.db
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
3. **evaluator** (`internal/evaluator/`) — language-specific compile → test → coverage → mutation in isolated temp dirs; Python `self_contained` and `module_level` both supported; environment fingerprinting for cross-run comparison
4. **reporter** (`internal/reporter/`) — multi-dimensional aggregation (by model, language, scenario); composite score = 0.3×compile + 0.3×pass_rate + 0.2×coverage + 0.2×mutation; HTML report via Chart.js CDN; mutation breakdown (total/killed/survived/no_tests/timeouts/skipped/suspicious); truncation statistics with tuning recommendations; auto-generated insights and efficiency stats
5. **store** (`internal/store/`) — SQLite v2 schema; upsert on `(run_id, model, language, sample_id)`; artifact indexing; supports `--reuse-generated` for skipping API calls when same prompt+source+model already exists
6. **web** (`internal/web/`) — HTTP management UI (`utbench web`); supports Docker and local runs; model config display; database browsing; embedded static assets
7. **obs** (`internal/obs/`) — structured logging (slog wrapper); progress reporting with per-task status lines
8. **ctrl** (`internal/ctrl/`) — pause/resume gate for web-triggered pause during generation

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

- **Default `--class` is `self_contained`**: current Python and Go datasets are also `self_contained`, so the default works correctly. Use `--class module_level` only when using actual module-level samples that require workspace context.
- **Module-level samples** require `meta.json` with `workspace_root` and `module_import`; the evaluator does not clean up their workspaces (reused in-place).
- **Checkpoint invalidation**: changing any of models, langs, class, level, manifest, max-samples, or dataset-root changes the hash and starts a fresh run.
- **Mutation testing tools**: Python uses `mutmut`, Go uses `go-mutesting` (`go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest`), Java uses `pitest` (Maven plugin), C++ uses `mull`. Windows mutation testing for Python is validated on Linux only; use Docker on Windows.
- **Line endings on Windows**: normalize with `git add --renormalize .`
- **`utbench ingest` is deprecated**: replaced by `utbench db ingest-evaluation`, `utbench db ingest-manifest`, `utbench db ingest-report`, and `utbench db ingest-run`.
- **Run with `--ingest`**: `./utbench run --ingest --db-path ./storage/utbench.db` auto-ingests the run directory into SQLite after completion.

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
| `ARK_API_KEY` | doubao-seed-2.0-lite/1.6/2.0-pro-v2, glm-4.7, deepseek-v3.2 |
| `ANTHROPIC_AUTH_TOKEN` | Claude Code 自定义模型端点认证（设为对应 provider 的 API key） |

# Supported Languages and Tools

| Language | Test Framework | Coverage | Mutation Tool |
|----------|---------------|----------|--------------|
| Python   | pytest        | coverage | mutmut       |
| Go       | go test       | go test -cover | go-mutesting |
| Java     | JUnit 5 (Maven) | JaCoCo | pitest      |
| C++      | GoogleTest    | gcov     | mull         |

# Prompt System

The runner uses three prompt modes (`internal/runner/prompt.go`):
- `full_file` — default for self-contained samples; full source + instructions
- `completion` — for continuation after truncation
- `module_level` — for samples with workspace context and module imports

Prompt strategy: `structured-v1`. System message emphasizes runnable tests only, no explanations or placeholders. Use `BuildPromptCatalog()` to inspect templates. Version ID is a SHA1 hash of the entire catalog content for traceability.

# Key Files to Read First

- `go-ut-bench/cmd/utbench/main.go` — CLI entry point, command routing
- `go-ut-bench/internal/orchestrator/service.go` — pipeline orchestration logic
- `go-ut-bench/internal/contracts/spec.go` — core data structures (RunSpec, SampleRef)
- `go-ut-bench/internal/contracts/constants.go` — SchemaVersion, DatasetClass, RunMode
- `go-ut-bench/internal/contracts/results.go` — all result/report data structures
- `go-ut-bench/internal/runner/prompt.go` — prompt construction (3 modes)
- `go-ut-bench/internal/runner/api.go` — LLM API client with auto-continuation
- `go-ut-bench/internal/runner/models.go` — model config loading from YAML
- `go-ut-bench/internal/evaluator/service.go` — evaluation pipeline per language
- `go-ut-bench/internal/reporter/service.go` — report aggregation and generation
- `go-ut-bench/internal/store/sqlite.go` — SQLite v2 schema and queries
- `go-ut-bench/internal/web/server.go` — Web management UI server
- `go-ut-bench/configs/models.yaml` — model provider/endpoints/API key env vars
