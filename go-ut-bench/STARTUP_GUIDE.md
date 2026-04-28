# go-ut-bench 完整启动指南

## 目录

1. [环境准备](#1-环境准备)
2. [构建 CLI](#2-构建-cli)
3. [配置 API 密钥](#3-配置-api-密钥)
4. [选择数据集](#4-选择数据集)
5. [选择模型](#5-选择模型)
6. [选择语言](#6-选择语言)
7. [控制样本数](#7-控制样本数)
8. [控制并发](#8-控制并发)
9. [变异测试](#9-变异测试)
10. [选择运行方式：本地 vs Docker](#10-选择运行方式本地-vs-docker)
11. [完整流水线：一键运行](#11-完整流水线一键运行)
12. [分步执行](#12-分步执行)
13. [查看结果](#13-查看结果)
14. [持久化到数据库](#14-持久化到数据库)
15. [常见场景速查](#15-常见场景速查)
16. [全量参数参考](#16-全量参数参考)

---

## 1. 环境准备

### 最低要求（仅 dry-run / 生成阶段）

- Go 1.21+

### 完整评测依赖（各语言评测工具）

| 语言 | 必须安装 | 变异测试|
|------|----------|-----------------|
| Python | `python3`、`pytest`、`coverage` | `mutmut` |
| Go | Go 工具链（已有即可） | `go-mutesting` |
| Java | JDK 17+、Maven 3+ | pitest（Maven 插件，自动拉取） |
| C++ | CMake、GoogleTest、`gcov`/`llvm-cov` | `mull` |

```bash
# Python 依赖
pip install pytest coverage mutmut

# Go 变异工具
go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
```

> **Windows 用户**：`dataset stats` 命令依赖 `python3`（注意是 `python3`，不是 `python`）。
> 建议使用 WSL、Docker 或在 PATH 中添加 `python3` 别名。

---

## 2. 构建 CLI

在项目根目录执行：

**Linux / macOS / WSL：**
```bash
go build -o utbench ./cmd/utbench/
```

**Windows PowerShell：**
```powershell
go build -o utbench.exe ./cmd/utbench/
```

验证构建成功：
```bash
./utbench --help
```

---

## 3. 配置 API 密钥

根据你要使用的模型，创建 `.env` 文件（或直接设置为系统环境变量）：

```dotenv
# DeepSeek
DEEPSEEK_API_KEY=sk-xxx

# 阿里云百炼 (qwen)
DASHSCOPE_API_KEY=sk-xxx

# MiniMax
MINIMAX_API_KEY=xxx

# 火山引擎 doubao-seed（旧端点）
VOLCENGINE_API_KEY=xxx

# 火山引擎 ARK（doubao-seed 新端点 / doubao-seed-2.0-lite / doubao-seed-1.6 / doubao-seed-2.0-pro-v2 / glm-4.7）
ARK_API_KEY=xxx
```

加载 `.env`（Linux/macOS）：
```bash
export $(grep -v '^#' .env | xargs)
```

加载 `.env`（Windows PowerShell）：
```powershell
Get-Content .env | ForEach-Object {
  if ($_ -match '^([^#][^=]*)=(.*)$') {
    [System.Environment]::SetEnvironmentVariable($Matches[1], $Matches[2], 'Process')
  }
}
```

> Docker 场景下直接用 `--env-file .env`，无需手动 export。

---

## 4. 选择数据集

### 4.1 数据集目录结构

```
datasets/
  python/
    python_code_files_self_contained/
      boundary/          # 边界值测试
      simple_function/   # 简单函数
      complex_dependency/ # 复杂依赖
      interface_mock/    # 接口 mock
    python_code_files_module_level/
      ...
  go/
    go_code_files_self_contained/
      boundary/
      simple_function/
      complex_dependency/
      interface_mock/
  java/   # 同上结构
  cpp/    # 同上结构
```

### 4.2 数据集选择参数

| 参数 | 可选值 | 说明 |
|------|--------|------|
| `--class` | `self_contained`、`module_level` | 数据集类别，逗号分隔可多选 |
| `--scenario` | `boundary`、`simple_function`、`complex_dependency`、`interface_mock` | 场景过滤，逗号分隔可多选，不填则全选 |
| `--level` | `l1`、`l2`、… | manifest level 过滤（与 `--dataset-manifest` 配合使用） |
| `--dataset-root` | 路径 | 数据集根目录，默认 `./datasets` |
| `--dataset-manifest` | JSON 路径 | 使用预生成的 manifest，与 `--class`/`--scenario` 二选一 |

### 4.3 使用预生成 manifest（精细控制）

先建索引，再生成 manifest：

```bash
# 第一步：扫描数据集，建索引
./utbench dataset index \
  --dataset-root ./datasets \
  --output ./configs/dataset_index.json

# 第二步：按条件筛选，生成 manifest
./utbench dataset manifest \
  --index ./configs/dataset_index.json \
  --langs python,go \
  --class self_contained \
  --scenario boundary,simple_function \
  --level l1 \
  --limit-per-scenario 20 \
  --output ./configs/my_dataset.json

# 第三步：运行时传入 manifest
./utbench run \
  --dataset-manifest ./configs/my_dataset.json \
  --models deepseek \
  --output-root ./artifacts \
  --config ./configs/models.yaml
```

### 4.4 查看数据集健康状态

```bash
./utbench dataset stats --dataset-root ./datasets
```

---

## 5. 选择模型

### 5.1 当前可用模型（`configs/models.yaml`）

| 模型键名 | 提供商 | 所需环境变量 |
|---------|--------|-------------|
| `deepseek` | DeepSeek | `DEEPSEEK_API_KEY` |
| `qwen` | 阿里云百炼 | `DASHSCOPE_API_KEY` |
| `minimax` | MiniMax | `MINIMAX_API_KEY` |
| `doubao-seed` | 火山引擎（旧端点） | `VOLCENGINE_API_KEY` |
| `doubao-seed-2.0-lite` | 火山引擎 ARK | `ARK_API_KEY` |
| `doubao-seed-1.6` | 火山引擎 ARK | `ARK_API_KEY` |
| `doubao-seed-2.0-pro-v2` | 火山引擎 ARK | `ARK_API_KEY` |
| `glm-4.7` | 火山引擎 ARK | `ARK_API_KEY` |

### 5.2 使用方式

```bash
# 单模型
--models deepseek

# 多模型并行评测（逗号分隔，无空格）
--models deepseek,qwen,minimax

# 全部模型
--models deepseek,qwen,minimax,doubao-seed,doubao-seed-2.0-lite,doubao-seed-1.6,doubao-seed-2.0-pro-v2,glm-4.7
```

### 5.3 添加或修改模型

编辑 `configs/models.yaml`，仿照已有格式添加：

```yaml
models:
  my-model:
    enabled: true
    provider: openai   # 兼容 OpenAI API 格式
    config:
      api_endpoint: "https://api.example.com/v1"
      model: "model-name"
      api_key_env: "MY_MODEL_API_KEY"
      parameters:
        temperature: 0.7
        top_p: 0.9
        max_tokens: 4096
```

---

## 6. 选择语言

| 语言键名 | 测试框架 | 覆盖率工具 |
|---------|----------|-----------|
| `python` | pytest | coverage |
| `go` | go test | go tool cover |
| `java` | JUnit 5 / Maven | JaCoCo |
| `cpp` | GoogleTest | gcov |

```bash
# 单语言
--langs python

# 多语言（逗号分隔）
--langs python,go,java,cpp
```

---

## 7. 控制样本数

```bash
# 每个语言/场景组合最多取 5 个样本
--max-samples 5

# 不限制（默认，跑全量数据集）
--max-samples 0
```

> 建议调试时先用 `--max-samples 1` 或 `--max-samples 2` 验证流程。

---

## 8. 控制并发

### 8.1 worker 数量

```bash
# 自动选择（默认）
--workers 0

# 指定 8 个 worker 并发
--workers 8
```

### 8.2 配置文件级别并发上限（`configs/models.yaml`）

```yaml
benchmark:
  parallel:
    max_concurrent_models: 6    # 同时跑的模型数上限
    max_concurrent_samples: 8   # 同时处理的样本数上限
```

### 8.3 运行模式

```bash
# 全量运行（默认）
--mode full

# 增量运行（跳过已有结果，断点续跑）
--mode incremental

# 重置断点，强制从头跑
--reset-checkpoint
```

---

## 9. 变异测试

变异测试耗时较长，默认关闭。

```bash
# 开启变异测试
--mutation-enabled

# 设置超时（秒，默认 1800）
--mutation-timeout 3600

# 变异失败策略：warn = 记录警告继续跑，fail = 直接报错退出
--mutation-policy warn
--mutation-policy fail
```

---

## 10. 选择运行方式：本地 vs Docker

### 方式 A：本地直接运行

适合：已安装好所有评测工具链、想快速迭代的场景。

```bash
# 确保 API 密钥已在环境变量中
./utbench run \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --max-samples 3
```

### 方式 B：Docker（推荐，评测工具链全预装）

Docker 镜像已预装：Python pytest/coverage/mutmut、Go go-mutesting、JDK 21 / Maven、C++ clang/gcov/GoogleTest/mull。

**第一步：构建镜像（只需一次）**

```bash
docker build -t utbench:latest .
```

**第二步：准备目录和密钥**

```bash
mkdir -p artifacts storage
# 确保 .env 文件已配置好 API 密钥
```

**第三步：运行**

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -v "$(pwd)/storage:/app/storage" \
  utbench:latest run \
    --models deepseek \
    --langs python \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --config /app/configs/models.yaml \
    --class self_contained \
    --max-samples 3
```

> **注意**：容器内路径以 `/app/` 开头，主机路径通过 `-v` 挂载映射。

**Windows PowerShell（Docker）：**

```powershell
docker run --rm --env-file .env `
  -v "${PWD}/datasets:/app/datasets" `
  -v "${PWD}/artifacts:/app/artifacts" `
  -v "${PWD}/configs:/app/configs" `
  -v "${PWD}/storage:/app/storage" `
  utbench:latest run `
    --models deepseek `
    --langs python `
    --dataset-root /app/datasets `
    --output-root /app/artifacts `
    --config /app/configs/models.yaml `
    --class self_contained `
    --max-samples 3
```

### 方式 C：使用仓库内置脚本

```bash
# Bash / WSL：参数顺序 models langs max_samples mutation
./run_bench.sh deepseek,minimax python,go 5 false
```

```powershell
# PowerShell
.\run_bench.ps1 -Models deepseek,minimax -Langs python,go -MaxSamples 5 -Mutation $false
```

---

## 11. 完整流水线：一键运行

`run` 命令自动完成 generate → evaluate → report 三个阶段。

### 最小 dry-run（验证流程，不调用真实 API）

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

### 标准运行（单模型、单语言）

```bash
./utbench run \
  --run-id my_run_001 \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --max-samples 10 \
  --workers 4
```

### 多模型横向对比

```bash
./utbench run \
  --run-id compare_001 \
  --models deepseek,qwen,minimax,glm-4.7 \
  --langs python,go \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --scenario boundary,simple_function \
  --max-samples 20 \
  --workers 8 \
  --ingest \
  --db-path ./storage/utbench.db
```

### 带变异测试的完整运行

```bash
./utbench run \
  --run-id full_mutation_001 \
  --models deepseek \
  --langs python,go \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --max-samples 5 \
  --mutation-enabled \
  --mutation-timeout 3600 \
  --mutation-policy warn \
  --ingest \
  --db-path ./storage/utbench.db
```

---

## 12. 分步执行

当需要单独控制每个阶段，或在不同机器上分阶段运行时使用。
**关键：三步必须使用同一个 `--run-id`，否则产物落入不同目录。**

```bash
RUN_ID=my_step_run_001

# 第一步：生成测试（调用 LLM API）
./utbench generate \
  --run-id "$RUN_ID" \
  --models deepseek,qwen \
  --langs python,go \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --max-samples 5

# 第二步：评测（编译 / 运行测试 / 统计覆盖率）
./utbench evaluate \
  --run-id "$RUN_ID" \
  --manifest ./artifacts/runs/$RUN_ID/generated/generated_manifest.json \
  --output-root ./artifacts \
  --mutation-enabled

# 第三步：生成报告
./utbench report \
  --run-id "$RUN_ID" \
  --input ./artifacts/runs/$RUN_ID/evaluation/evaluation_result.json \
  --output-root ./artifacts

# 第四步（可选）：写入数据库
./utbench ingest \
  --input ./artifacts/runs/$RUN_ID/evaluation/evaluation_result.json \
  --db-path ./storage/utbench.db
```

> **参数容易混淆的地方：**
> - `evaluate` 用 `--manifest`（不是 `--input`）
> - `report` 用 `--input`（不是 `--manifest`）
> - `ingest` 用 `--db-path`（不是 `--db`）

---

## 13. 查看结果

每次运行产物统一落在：

```
artifacts/runs/<run-id>/
  generated/
    generated_manifest.json       # 生成索引
    tests/                        # 生成的测试文件
    metadata/                     # 元数据
  evaluation/
    evaluation_result.json        # 评测明细
  report/
    report_summary.json           # 汇总 JSON
    report.html                   # 可视化 HTML 报告
  run_summary.json                # 整体运行摘要

artifacts/checkpoints/
  runner_<hash>.checkpoint.json   # 断点续跑用
```

用浏览器直接打开 `report.html` 查看可视化结果。

---

## 14. 持久化到数据库

```bash
# 运行时加 --ingest 一步到位
./utbench run ... --ingest --db-path ./storage/utbench.db

# 或事后单独写入
./utbench ingest \
  --input ./artifacts/runs/<run-id>/evaluation/evaluation_result.json \
  --db-path ./storage/utbench.db
```

数据库文件是标准 SQLite，可用任意 SQLite 客户端查询（如 `sqlite3`、DB Browser for SQLite）。

---

## 15. 常见场景速查

### 场景 A：第一次跑，验证环境是否通

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

### 场景 B：快速小批量真实评测

```bash
./utbench run \
  --models deepseek \
  --langs python \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --max-samples 3
```

### 场景 C：多模型 × 多语言全量横评（Docker）

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -v "$(pwd)/storage:/app/storage" \
  utbench:latest run \
    --run-id full_eval_001 \
    --models deepseek,qwen,minimax,glm-4.7 \
    --langs python,go,java,cpp \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --config /app/configs/models.yaml \
    --class self_contained \
    --max-samples 20 \
    --workers 8 \
    --ingest \
    --db-path /app/storage/utbench.db
```

### 场景 D：断点续跑（中途失败后继续）

```bash
./utbench run \
  --run-id my_run_001 \
  --mode incremental \
  --models deepseek,qwen \
  --langs python,go \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --max-samples 20
```

### 场景 E：只跑 Go 语言 boundary 场景

```bash
./utbench run \
  --models deepseek \
  --langs go \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --config ./configs/models.yaml \
  --class self_contained \
  --scenario boundary \
  --max-samples 10
```

### 场景 F：使用预筛选 manifest 精确控制样本

```bash
# 先生成 manifest
./utbench dataset manifest \
  --index ./configs/dataset_index.json \
  --langs python \
  --class self_contained \
  --scenario boundary \
  --limit-per-scenario 10 \
  --output ./configs/boundary_py.json

# 使用 manifest 运行
./utbench run \
  --dataset-manifest ./configs/boundary_py.json \
  --models deepseek,qwen \
  --output-root ./artifacts \
  --config ./configs/models.yaml
```

---

## 16. 全量参数参考

### `utbench run` / `utbench generate` 共享参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--config` | string | `../benchmark/config/models.yaml` | **必须显式指定** `./configs/models.yaml` |
| `--models` | string | 空 | 模型键名，逗号分隔 |
| `--langs` | string | 空 | 语言，逗号分隔：`python,go,java,cpp` |
| `--dataset-root` | string | `./datasets` | 数据集根目录 |
| `--dataset-manifest` | string | 空 | 预生成 manifest 路径，与 `--class`/`--scenario` 二选一 |
| `--class` | string | 空 | `self_contained` 或 `module_level`，逗号分隔 |
| `--scenario` | string | 空 | 场景过滤，逗号分隔，空=全选 |
| `--level` | string | 空 | manifest level 过滤 |
| `--max-samples` | int | `0` | 每个语言/场景最大样本数，`0`=不限 |
| `--output-root` | string | `./artifacts` | 产物根目录 |
| `--run-id` | string | 自动生成 | 本次运行 ID，分步执行时必须复用 |
| `--mode` | string | `full` | `full` 或 `incremental` |
| `--workers` | int | `0` | 并发 worker 数，`0`=自动 |
| `--mutation-enabled` | bool | `false` | 开启变异测试 |
| `--mutation-timeout` | int | `1800` | 变异测试超时（秒） |
| `--mutation-policy` | string | `warn` | 变异失败策略：`warn` 或 `fail` |
| `--dry-run` | bool | `false` | 跳过真实 API 调用 |
| `--reset-checkpoint` | bool | `false` | 清空断点，强制重跑 |
| `--ingest` | bool | `false` | 完成后自动写入 SQLite（仅 `run` 有效） |
| `--db-path` | string | `./storage/utbench.db` | SQLite 路径 |
| `-v` | bool | `false` | 详细日志 |

### `utbench evaluate`

| 参数 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `--manifest` | **是** | — | `generated_manifest.json` 路径 |
| `--run-id` | 否 | 自动生成 | 与 generate 同步复用 |
| `--output-root` | 否 | `./artifacts` | 产物根目录 |
| `--mutation-enabled` | 否 | `false` | 开启变异测试 |
| `--mutation-timeout` | 否 | `1800` | 超时秒数 |
| `--mutation-policy` | 否 | `warn` | `warn` 或 `fail` |
| `-v` | 否 | `false` | 详细日志 |

### `utbench report`

| 参数 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `--input` | **是** | — | `evaluation_result.json` 路径 |
| `--run-id` | 否 | 自动生成 | 与前步同步复用 |
| `--output-root` | 否 | `./artifacts` | 产物根目录 |
| `-v` | 否 | `false` | 详细日志 |

### `utbench ingest`

| 参数 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `--input` | **是** | — | `evaluation_result.json` 路径 |
| `--db-path` | 否 | `./storage/utbench.db` | SQLite 路径 |

### `utbench dataset manifest`

| 参数 | 必填 | 默认值 | 说明 |
|------|------|--------|------|
| `--index` | **是** | — | 索引文件路径（由 `dataset index` 生成） |
| `--output` | **是** | — | manifest 输出路径 |
| `--langs` | 否 | 空（全选） | 逗号分隔 |
| `--class` | 否 | 空（全选） | `self_contained` 或 `module_level` |
| `--scenario` | 否 | 空（全选） | 场景名，逗号分隔 |
| `--level` | 否 | 空（全选） | level 标签 |
| `--limit-per-scenario` | 否 | `20` | 每场景保留样本数 |
