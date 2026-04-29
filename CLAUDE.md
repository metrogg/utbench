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

### 快速演练（不调用 API）

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

# 导入 SQLite
./utbench ingest --evaluation ./artifacts/runs/<run-id>/evaluation/evaluation_result.json --db-path ./storage/utbench.db
```

### Docker（推荐用于运行评估工具链）

```bash
cd go-ut-bench
docker build -t utbench:latest .

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

### 数据集管理

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

# 验证数据集就绪状态
./utbench dataset validate --dataset-root ./datasets
```

### 工具链诊断

```bash
# 检查所有工具链（Go、Python、Java、C++、变异测试工具）
./utbench doctor
```

### 数据库查询

```bash
# 所有已导入数据概览
./utbench db overview --db-path ./storage/utbench.db

# 列出运行记录和结果
./utbench db list-runs --db-path ./storage/utbench.db
./utbench db list-results --run-id <run-id> --db-path ./storage/utbench.db
```

### Web UI

```bash
./utbench web --db-path ./storage/utbench.db
```

## 架构

五阶段流水线，由 `internal/orchestrator/service.go` 编排：

1. **数据集**（`internal/dataset/`）— 发现样本；从目录结构推断语言/类别/场景；计算 MD5 哈希；应用过滤器
2. **运行器**（`internal/runner/`）— 工作池并发调用 LLM API；基于检查点的增量执行；指数退避重试；截断检测和自动续写（当 `finish_reason: "length"` 或代码块不完整时最多重试 3 次）
3. **评估器**（`internal/evaluator/`）— 针对特定语言的编译→测试→覆盖率→变异测试，在隔离临时目录中执行；Python 的 `self_contained` 和 `module_level` 均支持
4. **报告器**（`internal/reporter/`）— 多维度聚合（按模型、语言、场景）；通过 Chart.js CDN 生成 HTML 报告；变异测试分解（总数/已杀死/存活/无测试/超时/跳过/可疑）；截断统计及调优建议
5. **存储**（`internal/store/`）— 基于 `(run_id, model, language, sample_id)` 的 SQLite 更新插入；所有表基于 SHA256 去重

### 并发模型

- 可配置 `--workers` 的工作池（默认 4）
- 任务以**轮询方式分配到各模型**，防止单个模型独占工作线程
- 每模型速率限制：API 调用间隔最低 200ms + 100ms 抖动
- 评估器中的看门狗在任务停滞 5 分钟后终止

### 综合评分

`0.3×编译 + 0.3×测试通过率 + 0.2×覆盖率 + 0.2×变异测试` — 按模型对所有合格样本计算。

### 失败来源分类

评估器对每个失败来源进行分类以确定评分资格：
- `model` — 生成的代码有缺陷（有评分资格，计入模型扣分）
- `environment` — 缺少工具链或系统依赖（排除评分）
- `dataset` — 样本本身存在问题（排除评分）
- `tool` — 评估器基础设施故障（排除评分）
- `none` — 无失败

只有 `model` 类型的失败影响综合评分。其他来源被排除并单独报告。

### 各语言评估工具链

| 语言 | 编译 | 测试 | 覆盖率 | 变异测试 |
|------|------|------|--------|----------|
| Python | `py_compile` | `pytest` | `coverage json` | `mutmut` |
| Go | `go build ./...` | `go test -v -coverprofile` | coverprofile 解析 | `go-mutesting` |
| Java | `mvn compile` | `mvn test`（JUnit5） | JaCoCo | PITest |
| C++ | `cmake --build` | `ctest` | `gcov` | Mull（基于 LLVM） |

### 数据契约

`internal/contracts/` 是所有跨阶段类型的唯一真相来源。关键文件：
- `spec.go` — `RunSpec`、`SampleRef`、`ModuleLevelMeta`
- `constants.go` — `SchemaVersion = "v0.1.0"`、`DatasetClass`、`RunMode`
- `results.go` — `GeneratedManifest`、`EvaluationResultSet`、`ReportPayload`

JSON 读写辅助函数位于 `internal/contracts/`。

### 检查点机制

运行器在增量模式下对 `(dataset_root + classes + level + manifest + max_samples + models)` 进行哈希 → SHA1 → `artifacts/checkpoints/runner_<hash>.checkpoint.json`。每个已完成任务的键：`<model>|<language>|<sample_id>`。作用域内参数的任何更改都会使检查点失效。

### 模型配置

`configs/models.yaml` — 每个模型的 provider/endpoint/api_key_env。默认代码路径为 `../benchmark/config/models.yaml`（相对于工作目录）；本地运行时始终传递 `--config ./configs/models.yaml`，Docker 中传递 `--config /app/configs/models.yaml`。

### 扩展系统

**添加新语言**：在 `internal/contracts/constants.go` 的 `SupportedLanguages` 中添加语言常量，然后在 `internal/evaluator/service.go` 中按照现有的每语言模式实现评估流程（编译/测试/覆盖率/变异测试）。

**添加新模型**：在 `configs/models.yaml` 中添加条目。如果提供商使用非标准 API 格式（非 OpenAI 兼容），需在 `internal/runner/api.go` 中添加特定提供商的响应解析。

## 常见陷阱

- **默认 `--class` 为 `self_contained`**：当前 Python 和 Go 数据集也是 `self_contained`，因此默认值正确。仅当使用需要工作区上下文的实际模块级样本时才使用 `--class module_level`。
- **模块级样本**需要包含 `workspace_root` 和 `module_import` 的 `meta.json`；评估器不会清理其工作区（原地重用）。
- **检查点失效**：更改 models、langs、class、level、manifest、max-samples 或 dataset-root 中的任何一个都会改变哈希值并开始新的运行。
- **变异测试工具**：Python 使用 `mutmut`，Go 使用 `go-mutesting`，Java 使用 PITest，C++ 使用 Mull。变异测试需要 Linux；Windows 上请使用 Docker。
- **Windows 上的行尾**：使用 `git add --renormalize .` 进行规范化。

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
| `DASHSCOPE_API_KEY` | Qwen、GLM（Dashscope） |
| `MINIMAX_API_KEY` | MiniMax M2.7 |
| `MINIMAX2.5_API_KEY` | MiniMax M2.5 |
| `VOLCENGINE_API_KEY` | doubao-seed |
| `MIMO_V2.5_API_KEY` | MiMo V2.5 |
| `MIMO_V2.5_PRO_API_KEY` | MiMo V2.5 Pro |

## 提示词系统

运行器使用三种提示词模式（`internal/runner/prompt.go`）：
- `full_file` — 自包含样本的默认模式；完整源码 + 指令
- `completion` — 截断后续写模式
- `module_level` — 需要工作区上下文和模块导入的样本

提示词策略：`structured-v1`。系统消息强调仅生成可运行的测试，无解释或占位符。使用 `BuildPromptCatalog()` 检查模板。

## 首先阅读的关键文件

- `go-ut-bench/cmd/utbench/main.go` — CLI 入口点、命令路由
- `go-ut-bench/internal/orchestrator/service.go` — 流水线编排逻辑
- `go-ut-bench/internal/contracts/spec.go` — 核心数据结构（RunSpec、SampleRef）
- `go-ut-bench/internal/contracts/constants.go` — SchemaVersion、DatasetClass、RunMode
- `go-ut-bench/internal/runner/prompt.go` — 提示词构造
- `go-ut-bench/internal/evaluator/service.go` — 各语言评估流水线
