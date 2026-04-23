# AGENTS.md

## 1) Repository Structure

- `go-ut-bench/` — main Go CLI tool for multi-model unit-test generation (has its own AGENTS.md)
- Root `README.md` — legacy Python benchmark docs; the Python tool no longer exists at root
- `docker-compose.yml`, `docker.sh` — Docker orchestration scripts (must run from `go-ut-bench/`)

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
VOLCENGINE_API_KEY=...    # doubao
ARK_API_KEY=...           # newer doubao/glm models
BIGMODEL_API_KEY=...      # glm
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

- **Docker context**: `docker-compose.yml` and `docker.sh` at repo root expect a `Dockerfile` at root, but the actual `Dockerfile` is in `go-ut-bench/`. Always run Docker commands from `go-ut-bench/`.
- **Default `--class` is `self_contained`**: Python and Go datasets have ALL samples as `module_level`, not `self_contained`. Use `--class module_level` for those languages, or `--class self_contained` for Java/C++ only.
- **Model config path**: default in code is `../benchmark/config/models.yaml`. When running in Docker, always use `--config /app/configs/models.yaml`.
- **Line endings on Windows**: files may have CRLF; normalize with `git add --renormalize .`
- **mutmut on Windows**: mutation testing is primarily validated on Linux; Docker recommended
- **Checkpoint scope**: changing models, languages, class, level, manifest, max-samples, or dataset root invalidates the checkpoint
