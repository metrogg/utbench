# AGENTS.md

## 1) Project Overview

`go-ut-bench` is a Go CLI tool that benchmarks multi-model unit-test generation quality across multiple programming languages. It is the Go reimplementation of the legacy Python `ut-bench` workflow.

Supported languages and evaluation toolchains:

| Language | Test Framework | Coverage | Mutation Testing |
|----------|----------------|----------|------------------|
| Python   | pytest         | coverage | mutmut           |
| Go       | go test        | go test -cover | gremlins    |
| Java     | Maven / JUnit 5 | JaCoCo  | pitest           |
| C++      | GoogleTest     | gcov     | mull             |

Pipeline stages:

1. `generate` — call model APIs and write generated tests plus manifest.
2. `evaluate` — compile, execute, collect coverage and mutation metrics.
3. `report` — aggregate results into JSON and HTML reports.
4. `ingest` — persist evaluation results into SQLite.
5. `run` — orchestrate `generate -> evaluate -> report -> optional ingest`.

Current functional state:

- **Python**: compile / test / coverage / mutation work, including `self_contained` and `module_level`.
- **Go / Java / C++**: compile / test / coverage work; mutation is wired and depends on local toolchain availability.

## 2) Technology Stack

- **Go 1.21+**
- **Dependencies**
  - `gopkg.in/yaml.v3`
  - `modernc.org/sqlite`
- **External toolchains**
  - Python: `pytest`, `coverage`, `mutmut`
  - Go: `gremlins`
  - Java: JDK 17+, Maven 3+, pitest Maven plugin
  - C++: CMake, GoogleTest, gcov, mull

## 3) Build and Run Commands

### Build

Linux / macOS:

```bash
go build -o utbench ./cmd/utbench/
```

Windows PowerShell:

```powershell
go build -o utbench.exe ./cmd/utbench/
```

### Quick start (dry-run, verified against current CLI)

```bash
./utbench run \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --scenario simple_function \
  --max-samples 1 \
  --dry-run
```

### Full pipeline

```bash
./utbench run \
  --models deepseek,qwen,minimax \
  --langs python,go,java,cpp \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --scenario boundary,simple_function \
  --max-samples 5 \
  --mutation-enabled \
  --mutation-timeout 1800
```

Use model keys exactly as defined in `configs/models.yaml`. Current examples there include `deepseek`, `qwen`, `minimax`, `doubao-seed`, and `glm-4.7`.

### Step-by-step

If you want all artifacts to stay under one run directory, reuse the same `--run-id`:

```bash
RUN_ID=demo_step_001

# 1. Generate tests only
./utbench generate \
  --run-id "$RUN_ID" \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --max-samples 5 \
  --dry-run

# 2. Evaluate an existing manifest
./utbench evaluate \
  --run-id "$RUN_ID" \
  --manifest ./artifacts/runs/$RUN_ID/generated/generated_manifest.json \
  --output-root ./artifacts \
  --mutation-enabled

# 3. Generate report
./utbench report \
  --run-id "$RUN_ID" \
  --input ./artifacts/runs/$RUN_ID/evaluation/evaluation_result.json \
  --output-root ./artifacts

# 4. Ingest into SQLite
./utbench ingest \
  --input ./artifacts/runs/$RUN_ID/evaluation/evaluation_result.json \
  --db-path ./storage/utbench.db
```

Important:

- `evaluate` uses `--manifest`, not `--input`
- `ingest` uses `--db-path`, not `--db`
- report artifacts live under `report/`, not `reports/`

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

# Validate dataset layout and print dependency scan output
./utbench dataset stats --dataset-root ./datasets
```

`dataset stats` currently invokes `python3` for the dependency scan. On Windows, use WSL / Docker or make sure `python3` is resolvable on `PATH`.

## 4) Code Organization

```text
cmd/utbench/           # CLI entry point
internal/
  config/              # Defaults and validation helpers
  contracts/           # Shared cross-stage data structures
  dataset/             # Discovery, index, manifest, layout checks
  evaluator/           # Language-specific evaluation logic
  obs/                 # slog-based logger wrapper
  orchestrator/        # Full-pipeline orchestration
  reporter/            # JSON + HTML report generation
  runner/              # Model client, prompts, checkpoints
  store/               # SQLite init and ingestion
configs/               # Example configs and manifests
datasets/              # Benchmark samples
artifacts/             # Run outputs
storage/               # SQLite database file
migrations/            # SQL DDL snapshots
schemas/               # JSON schemas
```

## 5) Configuration

### Model config

CLI default path: `../benchmark/config/models.yaml`.

Practical note: when running from this repository root, prefer `--config ./configs/models.yaml`.

Example structure:

```yaml
models:
  deepseek:
    enabled: true
    provider: deepseek
    config:
      api_endpoint: "https://api.deepseek.com/v1"
      model: "deepseek-chat"
      api_key_env: "DEEPSEEK_API_KEY"
```

Required per model:

- `enabled`
- `provider`
- `config.api_endpoint`
- `config.model`
- `config.api_key_env`

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

Examples currently referenced by `configs/models.yaml`:

```bash
export DEEPSEEK_API_KEY="..."
export DASHSCOPE_API_KEY="..."
export MINIMAX_API_KEY="..."
export VOLCENGINE_API_KEY="..."
export ARK_API_KEY="..."
```

## 6) Output Structure

```text
artifacts/
  checkpoints/
    runner_<hash>.checkpoint.json
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
```

- `generated_manifest.json` is the required input for `evaluate`
- `evaluation_result.json` is the required input for `report` and `ingest`
- `run_summary.json` is written by `run`

## 7) Dataset Layout

```text
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

- Language is inferred from the top-level dataset directory
- Classes: `self_contained`, `module_level`
- Scenarios: `boundary`, `simple_function`, `complex_dependency`, `interface_mock`
- Module-level samples use `meta.json`

## 8) Code Style Guidelines

- Use standard Go formatting (`gofmt`)
- Keep package names short and lowercase
- Keep internal packages under `internal/`
- Prefer explicit `if err != nil` handling
- Use `internal/obs.Logger`
- Use `0o755` for directories and `0o644` for files
- Use `filepath.Join`
- Keep missing metrics absent (`nil`) rather than fabricating zeros

## 9) Testing Instructions

Run unit tests:

```bash
go test ./internal/...
```

Current test packages include:

- `internal/dataset`
- `internal/evaluator`
- `internal/reporter`
- `internal/runner`
- `internal/store`

## 10) Security Considerations

- API keys come from environment variables only
- Reports and logs must not leak secrets
- Outputs stay under the configured `--output-root`
- SQLite ingestion uses transactions
- Evaluator uses isolated temp directories per sample

## 11) Common Pitfalls

- **Config path**: CLI default is `../benchmark/config/models.yaml`; in this repo, pass `--config ./configs/models.yaml`
- **Step-by-step runs**: reuse the same `--run-id` if you want one artifact tree
- **Evaluate input**: use `--manifest`
- **Ingest DB flag**: use `--db-path`
- **Report directory**: outputs go to `report/`
- **Dataset validation**: there is no `dataset validate` subcommand; use `dataset stats`
- **Windows dataset stats**: dependency scan currently expects `python3`
- **Dry-run**: it validates the pipeline shape, not real model quality
