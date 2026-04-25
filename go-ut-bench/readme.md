# go-ut-bench

<<<<<<< HEAD
多模型、多语言单元测试生成效果评测 CLI。本文档已按当前仓库里的 `cmd/utbench/main.go` 和实际命令行为核对。

## 当前命令总览

```text
utbench run
utbench generate
utbench evaluate
utbench report
utbench ingest
utbench dataset index
utbench dataset manifest
utbench dataset stats
```

## 支持语言

| 语言 | 测试框架 | 覆盖率 | 变异测试 |
|------|----------|--------|----------|
| Python | pytest | coverage | mutmut |
| Go | go test | go test -cover | gremlins |
| Java | Maven / JUnit 5 | JaCoCo | pitest |
| C++ | GoogleTest | gcov | mull |

## 构建

Linux / macOS:

```bash
go build -o utbench ./cmd/utbench/
```

Windows PowerShell:

```powershell
go build -o utbench.exe ./cmd/utbench/
```

提示：CLI 的代码默认从 `../benchmark/config/models.yaml` 读取模型配置；如果你直接在当前仓库根目录运行，通常应显式传入 `--config ./configs/models.yaml`。

## 快速开始

查看帮助：

```bash
./utbench --help
```

```powershell
.\utbench.exe --help
```

最小 dry-run（已按当前实现核对）：

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

运行完成后会生成：

```text
artifacts/runs/<run-id>/generated/generated_manifest.json
artifacts/runs/<run-id>/evaluation/evaluation_result.json
artifacts/runs/<run-id>/report/report_summary.json
artifacts/runs/<run-id>/report/report.html
artifacts/runs/<run-id>/run_summary.json
```

## 分步执行

如果希望 `generate / evaluate / report` 的产物落在同一个 `run-id` 目录下，必须显式复用同一个 `--run-id`。

```bash
RUN_ID=demo_local_001

./utbench generate \
  --run-id "$RUN_ID" \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --max-samples 2 \
  --dry-run

./utbench evaluate \
  --run-id "$RUN_ID" \
  --manifest ./artifacts/runs/$RUN_ID/generated/generated_manifest.json \
  --output-root ./artifacts

./utbench report \
  --run-id "$RUN_ID" \
  --input ./artifacts/runs/$RUN_ID/evaluation/evaluation_result.json \
  --output-root ./artifacts

./utbench ingest \
  --input ./artifacts/runs/$RUN_ID/evaluation/evaluation_result.json \
  --db-path ./storage/utbench.db
```

关键点：

- `evaluate` 的输入参数是 `--manifest`，不是 `--input`
- `report` 的输入参数是 `--input`
- `ingest` 使用 `--db-path`，不是 `--db`
- 报告目录是 `report/`，不是 `reports/`

## 数据集命令

```bash
./utbench dataset index \
  --dataset-root ./datasets \
  --output ./configs/dataset_index.json

./utbench dataset manifest \
  --index ./configs/dataset_index.json \
  --langs python,go \
  --class self_contained \
  --level l1 \
  --limit-per-scenario 20 \
  --output ./configs/dataset_l1.json

./utbench dataset stats --dataset-root ./datasets
```

说明：

- `dataset stats` 会先校验数据集布局，再做依赖扫描
- 当前 Windows 环境下该命令依赖 `python3` 在 `PATH` 中可用；如果本机只有 `python`，建议使用 WSL / Docker 或自行补齐命令别名

## 环境要求

最小要求：

- Go 1.21+
- dry-run 仅需能编译本项目

完整评测常见依赖：

```bash
pip install pytest coverage mutmut
go install github.com/go-gremlins/gremlins/cmd/gremlins@latest
```

还需要：

- Java: JDK 17+、Maven 3+
- C++: CMake、GoogleTest、gcov、mull

## 模型配置

当前仓库示例配置在 [configs/models.yaml](configs/models.yaml)。

示例里常用的模型键包括：

- `deepseek`
- `qwen`
- `minimax`
- `doubao-seed`
- `glm-4.7`

对应环境变量以 `configs/models.yaml` 为准，当前文件里涉及：

- `DEEPSEEK_API_KEY`
- `DASHSCOPE_API_KEY`
- `MINIMAX_API_KEY`
- `VOLCENGINE_API_KEY`
- `ARK_API_KEY`

## 进一步文档

- [docs/quickstart.md](docs/quickstart.md)
- [docs/cli-spec.md](docs/cli-spec.md)
- [docs/DOCKER_GUIDE.md](docs/DOCKER_GUIDE.md)
=======
多语言单元测试生成效果横向评测 CLI 工具。

## 支持语言

| 语言 | 测试框架 | 覆盖率工具 | 变异测试 |
|------|---------|-----------|---------|
| Python | pytest | coverage | mutmut |
| Go | go test | go test -cover | go-mutesting |
| Java | Maven/JUnit | JaCoCo | pitest |
| C++ | GoogleTest | gcov | mull |

## 支持模型

| 模型 | Provider | 说明 |
|------|----------|------|
| deepseek | deepseek | DeepSeek Chat |
| qwen | dashscope | 通义千问 3.6-plus |
| minimax | minimax | MiniMax M2.7 |
| doubao-seed | volcengine | 豆包 Seed 2.0 Pro |
| doubao-seed-2.0-lite | volcengine | 豆包 Seed 2.0 Lite |
| doubao-seed-1.6 | volcengine | 豆包 Seed 1.6 |
| doubao-seed-2.0-pro-v2 | volcengine | 豆包 Seed 2.0 Pro V2 |

## 快速开始

### Docker 运行（推荐）

```bash
# 构建镜像
docker build -t utbench:latest .

# 配置 API 密钥
cp .env.example .env

# 运行测试
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek \
    --langs python,java \
    --config /app/configs/models.yaml \
    --dataset-root /app/datasets \
    --max-samples 5
```

### 本地运行

```bash
go build -o utbench ./cmd/utbench/
export DEEPSEEK_API_KEY="sk-xxx"
./utbench run --models deepseek --langs python --max-samples 5
```

## 命令

| 命令 | 说明 |
|------|------|
| `utbench run` | 完整流程 (generate → evaluate → report) |
| `utbench generate` | 仅生成单元测试 |
| `utbench evaluate` | 评测已生成的单元测试 |
| `utbench report` | 生成评测报告 |
| `utbench ingest` | 结果导入 SQLite |
| `utbench dataset` | 数据集管理 |
| `utbench doctor` | 检查评测工具链并运行 canary 自检 |

## 输出结构

```
artifacts/runs/<run-id>/
  generated/           # 生成的测试文件
  evaluation/          # 评测结果 JSON
  report/              # HTML 报告
  run.log              # 运行日志
  api.log              # API 调用日志
```

## 特性

- **截断自动续写**：检测到输出截断时自动发送续写请求
- **截断统计分析**：报告中显示截断率、续写统计、调优建议
- **增量运行**：支持 checkpoint 断点续跑
- **变异测试**：可选启用变异测试评估测试质量
- **评测自检**：`utbench doctor` 检查工具版本并运行临时 canary 样本
- **数据集审计**：`utbench dataset validate` 统计样本并标记外部 I/O、非确定性和复杂度风险

## 数据集说明

当前仓库内置数据集实际为 `self_contained`：Python、Go、Java、C++ 各 4 个场景，每个场景 50 个样本。正式运行请显式使用 `--class self_contained`，避免旧文档中的 module-level 说明造成样本集合不一致。

```bash
./utbench doctor --langs python,go,java,cpp --mutation-enabled --mutation-timeout 120
./utbench dataset validate --dataset-root ./datasets --langs python,go,java,cpp --class self_contained --strict
```

## 文档

- [用户指南](docs/USER_GUIDE.md) - 完整使用文档
- [CLI 参数](docs/cli-spec.md) - 命令行参数详解
- [Docker 使用](docs/DOCKER_GUIDE.md) - Docker 运行指南
- [架构设计](docs/architecture-mvp.md) - 系统架构

## 环境要求

### Python
```bash
pip install pytest coverage mutmut
```

### Go
```bash
go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
```

### Java
- JDK 17+
- Maven 3+

### C++
```bash
sudo apt-get install cmake clang-15 libgtest-dev g++-15 mull-15
```
>>>>>>> origin/feat/go
