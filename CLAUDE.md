# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 在此仓库中工作时提供指导。

## 语言规则

- 内部推理过程（thinking/推理/思考）必须使用简体中文
- 所有对用户的回复使用简体中文
- 代码注释、文档说明、问题分析均使用中文
- 变量名、函数名、类名等代码标识符本身保持英文不变
- 无论用户用什么语言提问，以上规则不可违反

## 仓库布局

主要代码库位于 `go-ut-bench/`（Go CLI）。根目录的 `README.md`、`docker-compose.yml` 和 `docker.sh` 已**过时** — 请忽略它们，以 `go-ut-bench/` 为准。

**重要**：仓库根目录存在一份未集成的 `internal/runner/` 副本（`runner/api.go`、`runner/prompt.go` 等）。这是最近提交中添加的，但在 **Go 模块之外**，不被构建导入。请始终编辑 `go-ut-bench/internal/runner/` 下的文件。

## 构建和运行

所有命令在 `go-ut-bench/` 目录下执行：

### 构建

```bash
go build -o utbench ./cmd/utbench/
```

### 测试

```bash
go test ./internal/...
# 运行特定测试
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

### 完整流水线

```bash
./utbench run --models deepseek,qwen --langs python,go --max-samples 5 --mutation-enabled
```

### 分步执行

```bash
# 仅生成测试
./utbench generate --models deepseek --langs python --max-samples 5

# 评估已有清单
./utbench evaluate --manifest ./artifacts/runs/<run-id>/generated/generated_manifest.json

# 生成报告
./utbench report --evaluation ./artifacts/runs/<run-id>/evaluation/evaluation_result.json

# Ingest into SQLite (replaces old `utbench ingest`)
./utbench db ingest-evaluation --evaluation ./artifacts/runs/<run-id>/evaluation/evaluation_result.json --db-path ./storage/utbench.db
```

### Docker（推荐用于运行评估工具链）

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

### 辅助脚本

预构建的 Docker 运行便捷脚本：

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
# 索引所有样本
./utbench dataset index --dataset-root ./datasets --output ./configs/dataset_index.json

# 构建过滤清单
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

`internal/contracts/` 是所有跨阶段类型的唯一真相来源。关键文件：
- `spec.go` — `RunSpec`、`SampleRef`、`ModuleLevelMeta`
- `constants.go` — `SchemaVersion = "v0.1.0"`、`DatasetClass`、`RunMode`
- `results.go` — `GeneratedManifest`、`EvaluationResultSet`、`ReportPayload`

JSON 读写辅助函数位于 `internal/contracts/`。

### 检查点机制

运行器在增量模式下对 `(dataset_root + classes + level + manifest + max_samples + models)` 进行哈希 → SHA1 → `artifacts/checkpoints/runner_<hash>.checkpoint.json`。每个已完成任务的键：`<model>|<language>|<sample_id>`。作用域内参数的任何更改都会使检查点失效。

### 模型配置

`configs/models.yaml` — provider/endpoint/api_key_env per model. Default code path is `../benchmark/config/models.yaml` (relative to working dir); always pass `--config ./configs/models.yaml` locally or `--config /app/configs/models.yaml` in Docker.

## Common Pitfalls

- **Default `--class` is `self_contained`**: current Python and Go datasets are also `self_contained`, so the default works correctly. Use `--class module_level` only when using actual module-level samples that require workspace context.
- **Module-level samples** require `meta.json` with `workspace_root` and `module_import`; the evaluator does not clean up their workspaces (reused in-place).
- **Checkpoint invalidation**: changing any of models, langs, class, level, manifest, max-samples, or dataset-root changes the hash and starts a fresh run.
- **Mutation testing tools**: Python uses `mutmut`, Go uses `go-mutesting` (`go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest`), Java uses `pitest` (Maven plugin), C++ uses `mull`. Windows mutation testing for Python is validated on Linux only; use Docker on Windows.
- **Line endings on Windows**: normalize with `git add --renormalize .`
- **`utbench ingest` is deprecated**: replaced by `utbench db ingest-evaluation`, `utbench db ingest-manifest`, `utbench db ingest-report`, and `utbench db ingest-run`.
- **Run with `--ingest`**: `./utbench run --ingest --db-path ./storage/utbench.db` auto-ingests the run directory into SQLite after completion.

## 代码风格

- `gofmt` 格式化；包名使用简短小写字母
- 显式 `if err != nil` 返回；仅在评估器/运行器的工作协程中进行 panic 恢复
- 使用 `internal/obs.Logger`（slog 包装器）；将 `logger` 传入服务而非使用全局变量
- 对 JSON 输出中缺失的指标使用 `nil` 指针 — 不要用零值替代缺失数据
- 所有路径构造使用 `filepath.Join`
- 文件权限：目录 `0o755`，文件 `0o644`

## 环境设置

复制 `.env.example`（仓库根目录）或创建 `go-ut-bench/.env` 并填写：

| 变量 | 提供商 |
|------|--------|
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
