# AGENTS.md

## 1) Repository Structure

- `go-ut-bench/` — main Go CLI tool for multi-model unit-test generation (has its own AGENTS.md)
- Root `README.md` — **obsolete** legacy Python docs; ignore and use `go-ut-bench/readme.md` instead

## 2) Primary Usage (Docker)

Build and run from `go-ut-bench/`:

```bash
cd go-ut-bench
docker build -t utbench:latest .
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest run \
    --models deepseek \
    --langs python \
    --config /app/configs/models.yaml \
    --max-samples 2
```

Or use helper scripts:
```bash
cd go-ut-bench
./run_bench.sh [models] [langs] [max-samples] [mutation]
```

## 3) Environment Setup

Copy `.env.example` to `.env` and fill in API keys:

```
DEEPSEEK_API_KEY=...
DASHSCOPE_API_KEY=...     # qwen
MINIMAX_API_KEY=...
VOLCENGINE_API_KEY=...    # doubao-seed (original)
ARK_API_KEY=...           # doubao-seed-2.0-lite, doubao-seed-1.6, doubao-seed-2.0-pro-v2, glm-4.7
BIGMODEL_API_KEY=...      # glm (legacy, if needed)
```

## 4) Model Config

Location: `go-ut-bench/configs/models.yaml`

Enabled models: `qwen`, `deepseek`, `minimax`, `doubao-seed`, `doubao-seed-2.0-lite`, `doubao-seed-1.6`, `doubao-seed-2.0-pro-v2`, `glm-4.7`

## 5) For Development in go-ut-bench

See `go-ut-bench/AGENTS.md` — covers build, CLI commands, dataset layout, testing, and output structure.

Quick build:
```bash
cd go-ut-bench
go build -o utbench ./cmd/utbench/
```

## 6) Pre-built Binary

`go-ut-bench/utbench` (Linux) and `go-ut-bench/utbench.exe` (Windows) may be present as pre-built artifacts.

## 7) Gotchas

- **Docker context**: run Docker commands from `go-ut-bench/`; the active Dockerfile is `go-ut-bench/Dockerfile`.
- **Dataset classes**: current code and docs use `self_contained` and `repo_level`. Do not use the old `module_level` name.
- **Model config path**: default in code is `../benchmark/config/models.yaml`. When running in Docker, always use `--config /app/configs/models.yaml`.
- **Line endings on Windows**: files may have CRLF; normalize with `git add --renormalize .`
- **mutmut on Windows**: mutation testing is primarily validated on Linux; Docker recommended
- **Checkpoint scope**: changing models, languages, class, level, manifest, max-samples, or dataset root invalidates the checkpoint
