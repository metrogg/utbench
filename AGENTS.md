# AGENTS.md

## 1) Repository Structure

This is a workspace containing a Go CLI benchmark tool and datasets:

- `go-ut-bench/` — main Go CLI tool for multi-model unit-test generation evaluation (has its own AGENTS.md)
- `validated_balanced_single_file_eval_handoff_v3/` — validated dataset (200 samples per language: Python, JavaScript, Java, C++, Go)
- `docker-compose.yml`, `docker.sh` — Docker orchestration for running go-ut-bench
- `多模型单测生成效果横向评测.md` — Chinese project documentation

## 2) Primary Usage (Docker)

The benchmark runs via Docker. Build and run from `go-ut-bench/`:

```bash
cd go-ut-bench
docker build -t utbench:latest .
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/../benchmark/config:/app/config" \
  utbench:latest run --models deepseek --langs python --max-samples 10 --dry-run
```

Or use root-level helper: `./docker.sh build && ./docker.sh run-python`

## 3) Environment Setup

Copy `.env.example` to `.env` in `go-ut-bench/` and fill in API keys:

```
DEEPSEEK_API_KEY=...
DASHSCOPE_API_KEY=...     # qwen
MINIMAX_API_KEY=...
VOLCENGINE_API_KEY=...    # doubao
```

## 4) Model Config

Location: `go-ut-bench/configs/models.yaml`

Enabled models: `qwen`, `deepseek`, `minimax`, `doubao-seed`

Required per model: `enabled`, `provider`, `config.api_endpoint`, `config.model`, `config.api_key_env`

## 5) For Development Work in go-ut-bench

See `go-ut-bench/AGENTS.md` for full details on:
- Go build commands (`go build -o utbench ./cmd/utbench/`)
- CLI commands (`run`, `generate`, `evaluate`, `report`, `ingest`, `dataset`)
- Dataset layout and code organization
- Testing (`go test ./internal/...`)
- Output structure

## 6) Pre-built Docker Image

`go-ut-bench/utbench.tar` contains a pre-built Docker image (~3GB). Load it:

```bash
docker load -i go-ut-bench/utbench.tar
```

## 7) Datasets

- `go-ut-bench/datasets/` — dataset for benchmark runs
- `validated_balanced_single_file_eval_handoff_v3/data/` — validated 5-language dataset (self-contained samples)

## 8) Gotchas

- Root AGENTS.md does NOT describe a Python benchmark (that project was replaced by go-ut-bench)
- When running Docker, model config path needs explicit override: `--config /app/config/models.yaml`
- `mutmut` (Python mutation testing) may have issues on Windows; Docker image uses Linux