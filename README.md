# UT-Bench

UT-Bench 是一个用于横向评测大语言模型和编码 Agent 单元测试生成能力的基准工具。当前有效实现位于 [go-ut-bench](go-ut-bench/)，根目录的 `docs/` 主要保存调研、历史记录和外部工具资料；日常使用、开发和排障请以 [go-ut-bench/docs](go-ut-bench/docs/README.md) 为准。

## 项目背景

LLM 已经可以生成单元测试，但“能生成”不等于“生成得可靠”。在真实工程里，测试质量至少要回答这些问题：

- 生成的测试能否编译或被测试框架加载
- 测试本身是否能稳定通过
- 是否真正覆盖目标代码，而不是只跑空测试
- 是否能通过变异测试发现逻辑缺陷
- 同一模型在不同语言、不同场景下表现是否一致
- 编码 Agent 相比纯模型 API 是否真的有收益
- Agent 加入 Skill、工具说明或工作区上下文后提升多少

UT-Bench 的目标就是把这些问题变成可重复执行的评测流水线。它将模型、Agent、Skill 统一抽象为 `subject`，在同一数据集、同一提示词策略、同一评测工具链下生成测试、执行评测、汇总报告，并把结果落到文件产物和 SQLite 中，方便横向对比和历史追踪。

## 当前定位

当前 Go 版本不仅支持纯模型 API，也支持 CLI Agent 评测：

```text
subject = framework + model + optional skill
```

典型 subject：

- `model_api__deepseek-v4-flash__no_skill`：纯模型基线
- `opencode__deepseek-v4-flash__no_skill`：OpenCode Agent 基线
- `opencode__deepseek-v4-flash__unit_test_skill`：Agent + Skill

核心链路：

```text
sample -> subject generates tests -> compile/test/coverage/mutation -> report
```

## 能力概览

| 能力 | 说明 |
| --- | --- |
| 多语言评测 | Python、Go、Java、C++ |
| 多维指标 | 编译、测试通过率、覆盖率、变异测试、耗时、失败归因 |
| 统一 subject | 支持纯模型、CLI Agent、Agent + Skill 横向比较 |
| 自动化流水线 | `dataset -> generate -> evaluate -> report -> optional ingest` |
| 增量与复用 | checkpoint、生成资产复用、可选评测资产复用 |
| 报告输出 | JSON 汇总与 HTML 可视化报告 |
| SQLite 存储 | 支持跨运行查询、资产审计和复用解释 |
| Web 管理 UI | 新建任务、任务列表、报告查看、模型/资产管理 |
| Docker 工具链 | 评测镜像包含 Python/Go/Java/C++ 工具链和 mutation 工具 |
| Agent 沙箱 | 当前支持 Docker / local 方向，后续可扩展 gVisor、Firecracker、E2B 等后端 |

## 支持语言和工具

| 语言 | 测试框架 | 覆盖率 | 变异测试 |
| --- | --- | --- | --- |
| Python | pytest | coverage | mutmut |
| Go | go test | go test -cover | go-mutesting |
| Java | Maven/JUnit | JaCoCo | pitest |
| C++ | GoogleTest | gcov | mull |

## 仓库结构

```text
ut-bench/
  README.md                 # 当前总入口
  AGENTS.md                 # Codex/Agent 工作约定
  docs/                     # 父仓库层面的调研、外部资料、历史归档
  go-ut-bench/              # 当前有效的 Go CLI 实现
    cmd/utbench/            # CLI 入口
    internal/               # 核心模块
    configs/                # 模型、Agent、Skill、数据集配置
    datasets/               # 评测数据集
    docker/                 # Agent 沙箱等 Docker 资源
    docs/                   # 当前主文档
    artifacts/              # 运行产物
    storage/                # SQLite 数据库与缓存
```

根目录里的旧 Python/脚本式材料已经不再作为主实现入口。开发和运行命令默认进入 `go-ut-bench/` 执行。

## 快速上手

### 1. 进入 Go 实现目录

```bash
cd go-ut-bench
```

### 2. 配置 API Key

```bash
cp .env.example .env
```

按需要填写对应模型的密钥。常见变量包括：

| 环境变量 | 用途 |
| --- | --- |
| `DEEPSEEK_API_KEY` | DeepSeek |
| `DASHSCOPE_API_KEY` | Qwen / DashScope 兼容模型 |
| `MINIMAX_API_KEY` | MiniMax |
| `VOLCENGINE_API_KEY` | 豆包原始配置 |
| `ARK_API_KEY` | 火山 ARK 系列模型 |
| `ANTHROPIC_AUTH_TOKEN` | Codex 自定义模型端点认证 |

实际模型名、endpoint 和密钥变量以 [configs/models.yaml](go-ut-bench/configs/models.yaml) 为准。

### 3. 构建 CLI （可选）

```bash
go build -o utbench ./cmd/utbench/
```

Windows PowerShell 下也可以输出为：

```powershell
go build -o utbench.exe ./cmd/utbench/
```

### 4. 环境自检

```bash
./utbench doctor --langs python --mutation-enabled
```

如果要检查全部语言工具链：

```bash
./utbench doctor --langs python,go,java,cpp --mutation-enabled
```

### 5. 先跑 dry-run

dry-run 不调用真实 API，适合验证数据集、配置和输出目录是否通畅：

```bash
./utbench run \
  --models deepseek \
  --langs python \
  --max-samples 2 \
  --config ./configs/models.yaml \
  --dry-run
```

### 6. 跑一次最小正式评测

```bash
./utbench run \
  --models deepseek \
  --langs python \
  --max-samples 5 \
  --config ./configs/models.yaml
```

运行完成后查看：

```text
artifacts/runs/<run-id>/
  generated/generated_manifest.json
  evaluation/evaluation_result.json
  report/report_summary.json
  report/report.html
  run_summary.json
```

浏览器打开 `report/report.html` 可以查看可视化报告。

## 推荐使用方式 （推荐）

### Web UI

Web UI 适合日常任务管理、结果查看和资产审计：

```bash
./utbench web --addr :8080 --config ./configs/models.yaml
```

打开 `http://localhost:8080`，可以使用：

- 新建任务：选择模型、语言、数据集和 subject
- 任务列表：查看历史 run 和状态
- 任务详情：查看日志、报告和产物
- 注册表：查看 models / frameworks / skills / subjects
- 资产管理：查看生成和评测资产，解释复用命中情况

详细说明见 [Web UI 使用说明](go-ut-bench/docs/01-user-guides/WEB_UI.md)。

### CLI

常用命令：

| 命令 | 说明 |
| --- | --- |
| `utbench run` | 一键执行 generate -> evaluate -> report |
| `utbench generate` | 仅生成测试 |
| `utbench evaluate` | 读取 `generated_manifest.json` 执行评测 |
| `utbench report` | 读取 `evaluation_result.json` 生成报告 |
| `utbench doctor` | 检查语言工具链和评测环境 |
| `utbench dataset` | 数据集索引、manifest、validate |
| `utbench db` | SQLite 初始化、入库和查询 |
| `utbench assets` | 查询 subject、生成资产和评测资产 |
| `utbench web` | 启动 Web 管理 UI |

完整参数见 [CLI 命令文档](go-ut-bench/docs/01-user-guides/cli-spec.md) 和 [用户指南](go-ut-bench/docs/01-user-guides/USER_GUIDE.md)。

### Docker

推荐用构建脚本准备评测镜像和 Agent 沙箱镜像：

```bash
# Linux / WSL
./build.sh all
./build.sh all --cn

# Windows PowerShell
.\build.ps1 all
.\build.ps1 all -Cn
```

镜像角色：

| 镜像 | 用途 |
| --- | --- |
| `utbench:latest` | 评测镜像，包含多语言编译、测试、覆盖率和 mutation 工具 |
| `utbench-agent-opencode:latest` | Agent 沙箱镜像，用于 OpenCode CLI Agent subject |

更多挂载、DOOD、沙箱边界和排障说明见 [Docker 使用指南](go-ut-bench/docs/03-operations/DOCKER_GUIDE.md)。

## 典型评测场景

### 纯模型 baseline

```bash
./utbench run \
  --models deepseek \
  --langs python,go \
  --config ./configs/models.yaml \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --class self_contained \
  --max-samples 10
```

只传 `--models` 时，系统自动生成：

```text
model_api__<model>__no_skill
```

### Agent / Skill 横向比较

```bash
./utbench run \
  --models deepseek-v4-flash \
  --langs python \
  --config ./configs/models.yaml \
  --agents-config ./configs/agents.example.yaml \
  --subjects model_api__deepseek-v4-flash__no_skill,opencode__deepseek-v4-flash__no_skill,opencode__deepseek-v4-flash__unit_test_skill \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --class self_contained \
  --max-samples 1
```

如果只传 `--agents-config` 不传 `--subjects`，系统会自动展开合法的 `framework × model × skill` 组合，并保留纯模型 baseline。

### 生成、评测、报告分步执行

```bash
./utbench generate --models deepseek --langs python --max-samples 5 --config ./configs/models.yaml

./utbench evaluate --manifest ./artifacts/runs/<run-id>/generated/generated_manifest.json

./utbench report --evaluation ./artifacts/runs/<run-id>/evaluation/evaluation_result.json
```

### 入库和历史查询

```bash
./utbench db init --db-path ./storage/utbench.db

./utbench db ingest-run \
  --run-id <run-id> \
  --output-root ./artifacts \
  --db-path ./storage/utbench.db

./utbench db overview --db-path ./storage/utbench.db
```

也可以在运行时自动入库：

```bash
./utbench run \
  --models deepseek \
  --langs python \
  --max-samples 5 \
  --config ./configs/models.yaml \
  --ingest \
  --db-path ./storage/utbench.db
```

## 数据集与指标

当前默认数据集类别为 `self_contained`，适合单文件、自包含的函数或模块评测。数据集通常按语言、类别、场景组织：

```text
datasets/
  python/
    self_contained/
      boundary/
      simple_function/
      complex_dependency/
      interface_mock/
```

常见场景：

| 场景 | 说明 |
| --- | --- |
| `boundary` | 边界值测试 |
| `simple_function` | 简单函数 |
| `complex_dependency` | 复杂依赖 |
| `interface_mock` | 接口 mock |

评测指标包括：

- compile pass：测试是否能编译或加载
- test pass：测试是否能通过
- line / branch coverage：覆盖率
- mutation score：变异测试得分
- mutation breakdown：killed、survived、timeout、skipped 等明细
- latency / runtime：生成与评测耗时
- failure origin：失败归因

数据集治理说明见 [dataset-governance.md](go-ut-bench/docs/04-development/dataset-governance.md)。

## 架构速览

UT-Bench 采用分阶段流水线，由 `internal/orchestrator` 编排：

```text
dataset -> runner -> evaluator -> reporter -> store
```

| 阶段 | 模块 | 职责 |
| --- | --- | --- |
| dataset | `internal/dataset` | 发现样本、推断语言和场景、过滤数据集 |
| runner | `internal/runner` | 调模型或 Agent 生成测试，处理重试、续写、checkpoint |
| evaluator | `internal/evaluator` | 编译、测试、覆盖率和变异测试 |
| reporter | `internal/reporter` | 聚合指标，生成 JSON 和 HTML 报告 |
| store | `internal/store` | SQLite schema、入库、查询和资产索引 |
| web | `internal/web` | HTTP 管理 UI |
| contracts | `internal/contracts` | 跨阶段数据结构的唯一真相来源 |

更完整的设计见 [架构设计](go-ut-bench/docs/02-design/architecture-mvp.md)、[输出规范](go-ut-bench/docs/02-design/output-convention.md)、[数据库设计](go-ut-bench/docs/02-design/database-design.md)。

## 当前边界和注意事项

- 当前有效代码在 `go-ut-bench/`，不要把根目录旧材料当成主实现入口。
- 本仓库当前数据集主要是 `self_contained`；`repo_level` / `module_level` 需要对应元数据和工作区上下文。
- `artifacts/`、`storage/`、`.m2-cache/`、`utbench`、`utbench.exe` 是运行或构建产物，不属于人工文档整理范围。
- 纯模型评测和 Agent 评测都可能调用外部 API，请先确认 `.env`、`configs/models.yaml` 和网络环境。
- Agent 执行的是不可信代码生成流程。当前 Docker 沙箱是第一版方案，后续可向 gVisor、Firecracker、E2B 或云沙箱演进。
- Windows 上跑 Python mutation 建议使用 Docker；本地工具链差异会影响评测复用和结果可比性。

## 文档入口

| 目标 | 文档 |
| --- | --- |
| 快速使用 | [用户指南](go-ut-bench/docs/01-user-guides/USER_GUIDE.md) |
| 完整启动参数 | [完整启动指南](go-ut-bench/docs/01-user-guides/startup-guide.md) |
| CLI 命令参考 | [CLI 命令文档](go-ut-bench/docs/01-user-guides/cli-spec.md) |
| Web UI | [Web UI 使用说明](go-ut-bench/docs/01-user-guides/WEB_UI.md) |
| Docker 和沙箱 | [Docker 使用指南](go-ut-bench/docs/03-operations/DOCKER_GUIDE.md) |
| 项目介绍 | [项目详细介绍](go-ut-bench/docs/00-overview/project-introduction.md) |
| 架构设计 | [架构设计](go-ut-bench/docs/02-design/architecture-mvp.md) |
| 文档分类 | [文档分类与去重索引](go-ut-bench/docs/DOCUMENT_CLASSIFICATION.md) |
| 调研和历史资料 | [父仓库文档导航](docs/README.md) |

父仓库 `docs/` 中的研究资料包括 Agent 评测生态调研、UT-Bench 差距分析、repo-level 评测研究和外部工具资料。它们用于理解背景和后续路线；当前实现细节仍以 `go-ut-bench/docs/` 为准。
