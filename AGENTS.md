# AGENTS.md

## 1) Repository Structure

- `go-ut-bench/` — main Go CLI tool for multi-model unit-test generation (has its own AGENTS.md)
- `validated_balanced_single_file_eval_handoff_v3/` — validated dataset (200 samples × 5 languages)
- `docker-compose.yml`, `docker.sh` — Docker orchestration for go-ut-bench

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

## 3) Environment Setup

Copy `.env.example` to `.env` and fill in API keys:

```
DEEPSEEK_API_KEY=...
DASHSCOPE_API_KEY=...     # qwen
MINIMAX_API_KEY=...
VOLCENGINE_API_KEY=...    # doubao
```

## 4) Model Config

Location: `go-ut-bench/configs/models.yaml`

Enabled models: `qwen`, `deepseek`, `minimax`, `doubao-seed`

## 5) For Development in go-ut-bench

See `go-ut-bench/AGENTS.md` — covers build, CLI commands, dataset layout, testing, and output structure.

## 6) Pre-built Docker Image

`go-ut-bench/utbench.tar` (~3GB). Load it:

```bash
docker load -i go-ut-bench/utbench.tar
```

## 7) Gotchas

- **Default `--class` is `self_contained`**: Python and Go datasets have ALL samples as `module_level`, not `self_contained`. Use `--class module_level` for those languages, or `--class self_contained` for Java/C++ only.
- **Model config path**: always use `--config /app/configs/models.yaml` when running in Docker
- **Line endings on Windows**: files may have CRLF; normalize with `git add --renormalize .`
- **mutmut on Windows**: mutation testing is primarily validated on Linux; Docker recommended
- **Checkpoint scope**: changing models, languages, class, level, manifest, max-samples, or dataset root invalidates the checkpoint