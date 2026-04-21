# AGENTS.md

## 1) Project Overview

`go-ut-bench` is a Go CLI tool that benchmarks multi-model unit-test generation quality across programming languages. It is the Go reimplementation of the legacy Python `ut-bench` toolchain.

Supported languages and evaluation toolchains:

| Language | Test Framework | Coverage | Mutation Testing |
|----------|---------------|----------|------------------|
| Python   | pytest        | coverage | mutmut           |
| Go       | go test       | go test -cover | avito-tech/go-mutesting |
| Java     | Maven / JUnit 4 | JaCoCo | pitest           |
| C++      | GoogleTest    | gcov     | mull-15          |

The pipeline stages are:

1. `generate` — call LLM APIs to generate unit tests for dataset samples.
2. `evaluate` — compile, execute, collect coverage and mutation scores.
3. `report` — aggregate results into JSON and HTML reports.
4. `ingest` — persist evaluation results into SQLite.
5. `run` — orchestrates the full pipeline (`dataset -> generate -> evaluate -> report -> optional ingest`).

Current functional state:
- **Python**: full evaluation chain (compile → test → coverage → mutation) works, including `self_contained` and `module_level` samples.
- **Go / Java / C++**: compile, test, and coverage are implemented; mutation testing is wired but may require local toolchain setup.

## 2) Technology Stack

- **Go 1.21+**
- **Dependencies** (see `go.mod`):
  - `gopkg.in/yaml.v3` — model config parsing
  - `modernc.org/sqlite` — pure-Go SQLite driver (CGO-free)
- **External per-language toolchains** (must be installed separately):
  - Python: `pytest`, `coverage`, `mutmut`
  - Go: `go-mutesting`
  - Java: Maven 3+, JaCoCo, pitest Maven plugin
  - C++: CMake, GoogleTest, gcov, mull

## 3) Build and Run Commands

### Build

```bash
go build -o utbench ./cmd/utbench/
```

### Quick start (dry-run, no API calls)

```bash
./utbench run \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --max-samples 2 \
  --dry-run
```

### Full pipeline

```bash
./utbench run \
  --models deepseek,minimax,doubao \
  --langs python,go,java,cpp \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --class self_contained \
  --scenario boundary,simple_function \
  --max-samples 5 \
  --mutation-enabled \
  --mutation-timeout 1800
```

### Step-by-step

```bash
# 1. Generate tests only
./utbench generate --models deepseek --langs python --max-samples 5

# 2. Evaluate an existing manifest
./utbench evaluate \
  --input ./artifacts/runs/<run-id>/generated/generated_manifest.json \
  --mutation-enabled

# 3. Generate report
./utbench report \
  --input ./artifacts/runs/<run-id>/evaluation/evaluation_result.json

# 4. Ingest into SQLite
./utbench ingest \
  --input ./artifacts/runs/<run-id>/evaluation/evaluation_result.json \
  --db ./storage/utbench.db
```

### Dataset management

```bash
# Index all samples
./utbench dataset index --dataset-root ./datasets --output ./configs/dataset_index.json

# Build filtered manifest
./utbench dataset manifest \
  --index ./configs/dataset_index.json \
  --langs python,go \
  --class self_contained \
  --level l1 \
  --limit-per-scenario 20 \
  --output ./configs/dataset_l1.json

# Validate dataset layout and Python dependencies
./utbench dataset validate --dataset-root ./datasets --langs python
```

## 4) Code Organization

```
cmd/utbench/           # CLI entry point (main.go)
internal/
  config/              # Default values and validation helpers
  contracts/           # Cross-stage data structures and schema version
  dataset/             # Sample discovery, index/manifest builders, layout validation
  evaluator/           # Language-specific evaluation (compile, test, coverage, mutation)
  obs/                 # Structured logger wrapper around slog
  orchestrator/        # High-level pipeline orchestration
  reporter/            # JSON + HTML report generation
  runner/              # Model API client, prompt builder, checkpoint logic
  store/               # SQLite schema init and upsert ingestion
configs/               # Example configs and dataset manifests
datasets/              # Benchmark samples: <lang>/<class>/<scenario>/<files>
artifacts/             # Default output root (run artifacts)
storage/               # SQLite database file
migrations/            # SQL DDL snapshots
schemas/               # JSON Schema for evaluation_result and generated_manifest
```

### Key modules

- **`internal/contracts`** — single source of truth for:
  - `RunSpec`, `SampleRef`, `GeneratedManifest`, `EvaluationResultSet`, `ReportPayload`
  - `SchemaVersion = "v0.1.0"`
  - JSON read/write helpers (`WriteJSON`, `ReadGeneratedManifest`, etc.)

- **`internal/runner`** — `Service.Generate()` spawns a worker pool (default `min(16, max(2, CPU))`) to call model APIs concurrently. Supports:
  - Retry with exponential backoff for transient HTTP errors
  - Checkpointing in `incremental` mode (SHA1-scoped checkpoint files under `artifacts/checkpoints/`)
  - Dry-run placeholder generation
  - Prompt construction with bilingual (Chinese/English) instructions

- **`internal/evaluator`** — `Service.Evaluate()` spawns a worker pool to execute per-language toolchains in isolated temp directories. Supports Python self-contained and module-level samples.

- **`internal/reporter`** — Generates `report_summary.json` and `report.html` (uses Chart.js CDN for visualizations).

- **`internal/store`** — SQLite tables: `runs`, `sample_results`. Upsert on `(run_id, model, language, sample_id)`.

## 5) Configuration

### Model config

Default path: `../benchmark/config/models.yaml` (relative to working dir).

Example structure:

```yaml
models:
  deepseek:
    enabled: true
    provider: deepseek
    config:
      api_endpoint: "https://api.deepseek.com/v1"
      model: "deepseek-coder"
      api_key_env: "DEEPSEEK_API_KEY"
      parameters:
        temperature: 0.7
        max_tokens: 4096
```

Required per model: `enabled`, `provider`, `config.api_endpoint`, `config.model`, `config.api_key_env`.

### App defaults

Defined in `internal/config/config.go`:
- Dataset root: `./datasets`
- Output root: `./artifacts`
- DB path: `./storage/utbench.db`
- Default models: `["deepseek"]`
- Default languages: `["python"]`
- Default class: `self_contained`
- Default mode: `full`

### Environment variables for API keys

```bash
export DEEPSEEK_API_KEY="..."
export MINIMAX_API_KEY="..."
export VOLCENGINE_API_KEY="..."      # doubao
export DASHSCOPE_API_KEY="..."       # qwen
export TEST_API_KEY="..."            # for test configs
```

## 6) Output Structure

```
artifacts/
  runs/
    <run-id>/
      generated/
        generated_manifest.json
        tests/
          <model>/
            <lang>/
              *.test.<ext>
        metadata/
          <model>_<lang>_<sample_id>.response.json
          <model>_<lang>_<sample_id>.metadata.json
      evaluation/
        evaluation_result.json
      report/
        report_summary.json
        report.html
      run_summary.json
  checkpoints/
    runner_<hash>.checkpoint.json
```

- `generated_manifest.json` is the required input for `evaluate`.
- `evaluation_result.json` is the required input for `report` and `ingest`.
- Response JSONs retain raw API responses for debugging.

## 7) Dataset Layout

```
datasets/
  python/
    python_code_files_self_contained/
      boundary/
      simple_function/
      complex_dependency/
      interface_mock/
    python_code_files_module_level/
      <sample-id>/
        entry.py
        meta.json
        workspace/
          ...
  go/
    go_code_files_self_contained/
      ...
  java/
    java_code_files_self_contained/
      ...
  cpp/
    cpp_code_files_self_contained/
      ...
```

- Language is inferred from the top-level directory name (`python`, `go`, `java`, `cpp`).
- Classes: `self_contained`, `module_level`.
- Scenarios: `boundary`, `simple_function`, `complex_dependency`, `interface_mock`.
- Module-level samples contain a `meta.json` with `module_import`, `package_name`, `target_file`, `workspace_root`.

## 8) Code Style Guidelines

- Go standard formatting (`gofmt`).
- Package names are short and lowercase (`runner`, `evaluator`, `contracts`).
- Internal packages live under `internal/`.
- Error handling: explicit `if err != nil` returns; no panic recovery except in evaluator worker goroutines.
- Logging: use `internal/obs.Logger` (slog wrapper). Pass `logger` into services.
- File permissions: `0o755` for directories, `0o644` for files.
- Use `filepath.Join` for all path construction.
- Keep `nil` / pointer types for unavailable metrics in JSON output (do not fake zeros in summary JSON where missing data should be represented as absent).

## 9) Testing Instructions

Run Go unit tests:

```bash
go test ./internal/...
```

Test files present:
- `internal/dataset/index_test.go`
- `internal/dataset/service_test.go`
- `internal/evaluator/cpp_eval_test.go`
- `internal/evaluator/service_test.go`
- `internal/reporter/service_test.go`
- `internal/runner/api_test.go`
- `internal/runner/checkpoint_test.go`
- `internal/store/sqlite_test.go`

Note: the project does not have exhaustive unit tests for the benchmark code itself; many evaluation paths are integration tests against external toolchains (pytest, Maven, go test, etc.).

## 10) Security Considerations

- API keys are read **only** from environment variables (`config.api_key_env`). Never commit keys.
- Reports and logs must not emit API keys or request payloads containing secrets.
- All artifact output is constrained under the configured `--output-root`.
- SQLite ingestion uses transactions to avoid partial writes.
- Evaluator runs each sample in isolated temp directories to reduce cross-sample pollution.

## 11) Common Pitfalls

- **Model config path**: default is `../benchmark/config/models.yaml`. If you run from a different working directory, override with `--config`.
- **Checkpoint scope**: changing models, languages, class, level, manifest, max-samples, or dataset root changes the checkpoint hash and starts fresh.
- **Windows mutation testing**: `mutmut` (Python) is primarily validated on Linux; Windows behavior may vary.
- **Module-level samples**: require `meta.json` with valid `workspace_root` and `module_import`. The evaluator does not clean up module-level workspaces (they are reused in-place).
- **Dry-run**: produces placeholder tests that are structurally valid but do not exercise real model behavior.
- **Language inference**: based on dataset directory name, not file extension. Legacy scenario folder names (`boundary`, `simple_function`, etc.) are mapped to `python` in some compatibility paths.
