# UTBench Web 管理后台

基于 Go 内嵌 HTTP Server + 前端 SPA 的 Benchmark 运行管理界面。

---

## 快速启动

### 开发模式（直接运行）

```bash
cd go-ut-bench
go run ./cmd/utbench web -config ./configs/models.yaml
```

### 编译后运行

```bash
cd go-ut-bench
go build -o utbench.exe ./cmd/utbench

# 启动 Web 管理后台（默认端口 8080）
./utbench.exe web

# 指定端口和配置
./utbench.exe web -addr :9090 -config ./configs/models.yaml -dataset-root ./datasets

# 指定 agent 配置与 Docker 运行参数
./utbench.exe web -config ./configs/models.yaml -agents-config ./configs/agents.example.yaml -docker-image utbench:latest -project-root .
```

> **注意**：`--agents-config` 会自动从 `--config` 所在目录探测 `agents.yaml` 或 `agents.example.yaml`，通常无需手动指定。

打开浏览器访问 `http://localhost:8080` 即可进入管理界面。

---

## 命令参数

```
utbench web [flags]

  -addr                  HTTP 监听地址，默认 :8080
  -config                模型配置文件路径，默认 ./configs/models.yaml
  -dataset-root          数据集根目录，默认 ./datasets
  -output-root           产物根目录，默认 ./artifacts
  -db-path               SQLite 数据库路径，默认 ./storage/utbench.db
  -agents-config         Agent/skill 配置路径，默认自动从 --config 目录探测
  -docker-image          评测镜像名，默认 utbench:latest
  -docker-control-image  控制面镜像名，默认 utbench-control:latest
  -docker-eval-image     评测镜像名（优先于 --docker-image）
  -project-root          Docker 挂载的项目根目录，默认 .
  -env-file              Docker 运行时传入的 .env 文件，默认 ./.env
```

`--agents-config` 探测顺序：`<config-dir>/agents.yaml` → `<config-dir>/agents.example.yaml`，找到第一个即停止。若都不存在，Web UI 仅支持 model_api（纯 API）模式，新建任务页会显示提示。

---

## 功能概览

### 总览（Dashboard）

- 总任务数 / 运行中 / 已完成 / 失败 统计卡片
- 最近 10 次任务速览，点击跳转到详情页
- 顶部状态栏实时显示 Docker 环境状态（Docker 就绪 / 镜像就绪）

### 新建任务

交互式表单，完整覆盖 CLI 的主要运行参数。被测对象（Subject）支持**自由组合**：每行可独立选择 Framework × Model × Skill，可添加多行。

| 参数 | 控件 | 说明 |
|------|------|------|
| 任务 ID | 文本输入 | 留空自动生成 |
| 模式 | 下拉选择 | `full`（全量）/ `incremental`（增量） |
| Subject（自由组合） | 动态下拉行 | 每行三个下拉框：Framework（含 `model_api`）、Model、Skill（含 `no_skill`），可添加/删除多行，右侧实时显示生成的 Subject ID |
| 模型 Baseline | 多选卡片 | baseline / 兼容回退入口；未选 subject 时仍可按纯模型运行 |
| 语言 | 多选卡片 | python / go / java / cpp |
| 类型 | 下拉选择 | `self_contained` / `repo_level` |
| 场景 | 下拉选择 | `boundary` / `simple_function` / `complex_dependency` / `interface_mock` |
| 等级 | 文本输入 | 可选，如 `l1`、`l2` |
| 最大样本数 | 数字输入 | 每场景样本上限，0 = 不限 |
| 并发数 | 数字输入 | 并发 worker 数，0 = 自动 |
| 空跑模式 | 复选框 | 跳过真实 API 调用，验证流程 |
| 变异测试 | 复选框 | 开启 mutmut/go-mutesting/pitest/mull |
| 变异超时 | 数字输入 | 秒数，默认 1800 |
| 入库保存 | 复选框 | 完成后自动写入 v2 SQLite 数据库 |
| 复用历史生成 | 复选框 | 同 `subject_version + sample_uid + prompt` 时复用数据库里的 generated test，避免重复调用模型 |
| 复用历史评测 | 复选框 | 同 `generated_test + evaluation_env_fingerprint + mutation_config` 时复用评测结果 |
| Docker 执行 | 复选框 | 在 `utbench:latest` 容器中运行（Windows 下 mutmut/mull 必需） |

Framework 下拉框会根据 `agents.yaml` 中的 `compatible_models` 自动过滤可用 Model；Skill 下拉框会根据 `compatible_frameworks` 自动过滤。未加载 `agents.yaml` 时仅显示 `model_api`，页面会提示原因。

提交前页面会显示当前选中的 `subject / framework / skill` 规模和预计任务量。提交后自动跳转到任务详情页，开始实时跟踪日志。

### 任务列表

- 完整任务历史列表（合并内存中的活跃任务 + `artifacts/runs/*/run_summary.json` 扫描）
- 支持按 **任务 ID / subject / 模型名 / 语言** 文本过滤
- 支持按 **状态** 过滤：等待中 / 运行中 / 已完成 / 失败
- 显示被测对象、语言、场景、样本上限、开始时间、耗时

### 数据管理

面向可复现评测和横向对比的数据管理入口：

- 查看当前 SQLite 数据库路径和核心计数：生成运行、评测运行、生成样本、评测结果、artifact、报告
- 以 `subject -> language -> sample` 树形方式浏览历史资产
- 列出数据库中的运行记录，点击某个 run 可联动筛选评测结果和 artifact
- 查看评测结果表：模型、语言、样本、编译、测试、覆盖率、变异分
- 查看 artifact 索引：`generated_manifest`、`generated_test`、`model_response`、`generation_metadata`、`evaluation_result`、`report_summary`、`report_html` 等
- 输入已有 `run_id` 可补录 `artifacts/runs/<run_id>` 中的 manifest、evaluation、report 和关联 artifact
- 在“跨运行对比报告”中多选历史运行、模型和语言，生成新的数据库报告，用于把多个历史运行中的模型结果放到同一份报告中比较
- 对所选运行展示环境一致性提示：同一 `eval_env_id` 下的排名可作为正式横向比较；跨环境或环境指纹缺失时报告仅供参考

### 注册表

- 在“注册表”页统一查看 `models / frameworks / skills / subjects`
- `models.yaml` 仍负责模型连接配置
- `agents.yaml` 负责 framework、skill 和 subject 的组合注册
- 这个页面的目标是让“模型”和“agent 被测对象”不再混在同一层语义里

数据库 schema 见 [database-design.md](database-design.md)。Web 仍会保留基于 `artifacts/runs` 的任务列表，数据库页负责长期保存和跨运行查询。

### 任务详情

#### 实时日志（SSE 推送）

- 后台 goroutine 运行评测，通过 Server-Sent Events 向前端推送结构化日志
- 日志区域自动滚动到底部
- 运行结束后 SSE 自动关闭

#### 评测报告

运行完成后（Status = `completed`）可切换到「报告」标签：

- **汇总卡片**：总样本数、编译通过率、测试通过率、平均行覆盖率、平均变异得分
- **模型横向对比柱状图**（Chart.js）：编译通过率 / 测试通过率 / 行覆盖率 / 变异得分
- **排行榜**：按综合得分排序，含编译/测试/覆盖/变异/延迟/Token 数据
- **失败项摘要**：按错误类型聚合，含示例信息
- **打开 HTML 报告**：点击「打开 HTML 报告」按钮可直接在新标签页查看原始 HTML 报告文件

HTML 报告页采用可视化评测报告结构：

- 顶部紧凑概览展示模型、语言、样本类型和关键指标
- 模型排名区按综合得分展示各模型的编译、测试、覆盖、变异、延迟和 Token
- 图表分析区支持综合、语言、场景和模型详情切换
- 语言/场景统计、失败分析、计分剔除、零变异体、原始数据和 prompt 快照作为下钻区域

### 运行时拓扑与镜像构建

页面顶部状态栏显示 Docker 环境和三种镜像的状态：

- **Docker 就绪**：Docker daemon 可用（含版本号）
- **拓扑模式**：`host_control` / `container_control` / `container_control+nested_docker` 等
- **控制镜像**：`utbench-control:latest` 是否就绪
- **评测镜像**：`utbench:latest` 是否就绪

点击「构建镜像」打开构建弹窗，可选择构建目标：

- **评测镜像**（`Dockerfile`）：含全部语言工具链，用于 run / evaluate / doctor
- **控制镜像**（`Dockerfile.control`）：仅含 utbench 二进制 + docker CLI，用于 web 编排

环境检查页（`/environment`）显示完整的运行时拓扑卡片，包括拓扑模式、Docker 状态、镜像状态、嵌套沙箱就绪状态、Host Output Root 等。

任务详情页头部在运行完成后会显示本次运行的沙箱 provider、Agent 框架和评测环境指纹。

---

## 性能策略

Web UI 采用以下策略避免系统卡顿：

- **按页面轮询**：`loadRuns`（4s）仅在总览/任务列表/任务详情页生效；`loadDatabase`（12s）仅在资产管理页生效；`loadEnv`（30s）仅在新建任务/注册表页生效。切换页面时自动跳过无关轮询。
- **环境检测缓存**：`/api/env` 后端对 `DetectEnv` 结果缓存 15 秒，避免每次轮询都 spawn Docker 子进程。
- **任务列表缓存**：`/api/runs` 磁盘扫描缓存 10 秒，减少 `filepath.Glob` + JSON 解析频率。
- **日志非响应式**：运行日志通过 rAF 批处理直接写 DOM，计数器不触发 Alpine.js 响应式更新。
- **日志上限**：服务端 `RunEntry.logs` 上限 10000 行，前端 DOM 上限 5000 行。
- **资产页懒加载**：资产管理页首次进入时加载核心数据（subjects/generations/evaluations），次要数据延迟加载不阻塞 UI。

---

## 技术架构

```
┌─────────────────────────────────────────────┐
│  Browser (Alpine.js + Tailwind + Chart.js)  │
│  ─────────────────────────────────────────  │
│  Dashboard  NewRun  RunsList  RunDetail     │
│    │          │        │         │          │
│    └──────────┴────────┴─────────┘          │
│              REST API + SSE                   │
└─────────────────────────────────────────────┘
                      ↑
              :8080 HTTP Server
                      ↑
┌─────────────────────────────────────────────┐
│         Go HTTP Server (internal/web)        │
│  ─────────────────────────────────────────  │
│  /api/config          GET  →  models.yaml 解析     │
│  /api/env             GET  →  环境检测（Docker/工具）│
│  /api/env/build-image GET  →  构建任务状态         │
│  /api/env/build-image POST →  启动镜像构建        │
│  /api/env/build-image/:id/events SSE → 构建日志流 │
│  /api/runs            GET  →  任务列表聚合          │
│  /api/runs            POST →  提交新任务            │
│  /api/runs/{id}       GET  →  任务详情             │
│  /api/runs/{id}/events    SSE → 实时日志流        │
│  /api/runs/{id}/report    GET →  评测报告 JSON     │
│  /api/runs/{id}/report-html GET → 原始 HTML 报告   │
│  /                        →  静态 SPA 文件         │
└─────────────────────────────────────────────┘
                      ↑
┌─────────────────────────────────────────────┐
│         RunManager (internal/web)            │
│  ─────────────────────────────────────────  │
│  内存任务状态管理 + goroutine 异步执行        │
│  本地进程执行 / Docker 容器执行（可选）         │
│  lineWriter → 日志捕获 → SSE 广播            │
│  复用 orchestrator.Run 完成完整流水线        │
└─────────────────────────────────────────────┘
```

### 关键组件

| 文件 | 职责 |
|------|------|
| `internal/web/server.go` | HTTP Server、REST API 路由、CORS、静态文件服务 |
| `internal/web/run_manager.go` | 异步任务管理器：goroutine 执行、日志捕获（`lineWriter`）、SSE 订阅广播 |
| `internal/web/env.go` | 环境检测：Docker daemon、镜像存在性、本地变异工具 |
| `internal/web/build_manager.go` | Docker 镜像构建任务管理，支持 SSE 日志流 |
| `internal/web/docker_runner.go` | Docker 容器执行后端，自动挂载产物目录 |
| `internal/web/env_handlers.go` | `/api/env` 和 `/api/env/build-image` API 端点 |
| `internal/web/static/index.html` | 完整 SPA：Alpine.js 状态管理 + Tailwind 样式 + Chart.js 可视化 |
| `internal/store/sqlite.go` | v2 SQLite schema、manifest/evaluation/report 入库和数据库查询接口 |

---

## API 参考

### GET /api/config

返回当前系统配置，用于前端动态渲染模型/语言/场景选项。

```json
{
  "models": [
    { "name": "deepseek", "provider": "deepseek", "model_id": "deepseek-chat", "enabled": true }
  ],
  "frameworks": [
    { "name": "opencode", "kind": "cli_agent", "sandbox_mode": "docker", "compatible_models": ["deepseek-v4-flash"], "compatible_languages": ["python", "go", "java", "cpp"] }
  ],
  "skills": [
    { "name": "unit_test_skill", "version": "2", "inject_mode": "prompt_append", "compatible_frameworks": ["opencode", "claude_code"], "compatible_languages": ["python", "go", "java", "cpp"] }
  ],
  "subjects": [
    { "id": "opencode__deepseek-v4-flash__no_skill", "kind": "cli_agent", "framework": "opencode", "model": "deepseek-v4-flash", "skill": "no_skill" }
  ],
  "languages": ["python", "go", "java", "cpp"],
  "scenarios": ["boundary", "simple_function", "complex_dependency", "interface_mock"],
  "classes": ["self_contained", "repo_level"],
  "dataset_root": "./datasets",
  "config_path": "./configs/models.yaml",
  "agents_config_path": "./configs/agents.example.yaml",
  "agents_config_error": ""
}
```

`agents_config_error` 在 agents config 加载失败时返回错误信息（如文件不存在、模型交集为空等），前端会据此显示提示。`frameworks`、`skills`、`subjects` 仅在 agents config 加载成功时有值。

### GET /api/env

返回环境检测信息，用于前端显示 Docker/镜像状态和推荐提示。

```json
{
  "docker_available": true,
  "docker_version": "24.0.7",
  "image_present": true,
  "image_name": "utbench:latest",
  "project_root": "f:/Code/ut-bench/...",
  "os": "windows",
  "native_tools": { "mutmut": false, "go-mutesting": false, "pitest": false, "mull": false },
  "recommendation": ""
}
```

### POST /api/env/build-image

启动 Docker 镜像构建任务。

**响应：**

```json
{ "build_id": "build_20240424_123456", "status": "pending" }
```

### GET /api/env/build-image/:id/events

SSE 端点，推送镜像构建实时日志。

### GET /api/runs

返回所有任务记录（内存 + 磁盘扫描）。

```json
[
  {
    "run_id": "20240424T123456.123456789Z",
    "status": "running",
    "started_at": "2024-04-24T12:34:56Z",
    "ended_at": null,
    "error": "",
    "use_docker": false,
    "spec": { "models": ["deepseek"], "languages": ["python"], ... }
  }
]
```

### POST /api/runs

提交新任务。

**请求体：**

```json
{
  "run_id": "my_run_001",
  "models": ["deepseek", "qwen"],
  "languages": ["python", "go"],
  "class": "self_contained",
  "scenario": "boundary",
  "level": "",
  "max_samples": 10,
  "workers": 4,
  "mode": "full",
  "dry_run": false,
  "mutation_enabled": false,
  "mutation_timeout": 1800,
  "mutation_policy": "warn",
  "ingest": true,
  "use_docker": false
}
```

**响应：**

```json
{ "run_id": "my_run_001", "status": "pending" }
```

### GET /api/runs/{run_id}

获取任务详情（含当前日志快照）。

```json
{
  "run_id": "my_run_001",
  "status": "running",
  "started_at": "2024-04-24T12:34:56Z",
  "use_docker": false,
  "spec": { ... },
  "logs": ["[12:34:56.123] run started...", "..."]
}
```

### GET /api/runs/{run_id}/events

SSE 端点，推送实时日志。

```
Content-Type: text/event-stream

data: {"type":"log","payload":"[12:34:56.123] run started  id=..."}

data: {"type":"log","payload":"[12:34:57.456] generating tests..."}

data: {"type":"done","payload":"completed"}
```

### GET /api/runs/{run_id}/report

返回 `report_summary.json` 的完整内容。

### GET /api/runs/{run_id}/report-html

直接返回原始 HTML 报告文件（`artifacts/runs/{run_id}/report/report.html`），Content-Type 为 `text/html; charset=utf-8`。可直接在浏览器中打开查看。

---

## Docker 执行模式

### 何时使用 Docker 执行

| 场景 | 推荐后端 |
|------|---------|
| Windows 本地运行，需要 Python mutmut 变异测试 | Docker（Windows 下 mutmut 依赖 Unix 工具链） |
| Windows 本地运行，需要 C++ mull 变异测试 | Docker（mull 依赖 LLVM/Clang 环境） |
| Linux/macOS，且已安装所有原生工具 | 本地进程 |
| 希望完全隔离依赖环境 | Docker |

### 使用流程

1. 首次使用：点击顶部状态栏「构建镜像」按钮，等待 `utbench:latest` 构建完成
2. 新建任务时勾选「Docker 执行」复选框
3. 提交任务后，后台将在容器中运行完整的 generate → evaluate → report 流水线
4. 容器产物自动挂载到宿主机的 `artifacts/runs/<run_id>/` 目录，与本地模式完全兼容

---

## 日志捕获机制

Web 后台通过自定义 `io.Writer`（`lineWriter`）将 Go `slog` 日志同时输出到 stderr 和内存缓冲区：

1. `NewLoggerWithWriter(verbose, multiWriter)` 创建双输出 Logger
2. `lineWriter` 按行缓冲，将完整日志行追加到 `RunEntry.logs`
3. `RunEntry.appendLog()` 广播给所有 SSE 订阅者
4. 前端通过 `EventSource` 接收实时日志，自动渲染到终端面板

---

## 与 CLI 的关系

Web 后台复用了完整的 CLI 执行引擎：

- `orchestrator.Service.Run()` — 同一套 generate → evaluate → report 流水线
- `contracts.RunSpec` — 同一套规格定义
- `artifacts/runs/<run_id>/` — 同一套产物目录结构

这意味着 Web 提交的任务和 CLI 提交的任务产物完全兼容，可以互相查看。

---

## 常见问题

### Web 启动后浏览器无法访问

检查防火墙或绑定地址。默认绑定 `:8080`（所有接口），若只想本地访问：

```bash
./utbench.exe web -addr 127.0.0.1:8080
```

### 运行页面日志不更新

SSE 连接在运行结束后会自动关闭。如果运行状态卡在 `等待中` 或 `运行中`，检查后端日志是否有报错。

### 报告页面显示「报告将在完成后可用」

只有 `status === 'completed'` 时报告才会生成。若运行失败，请检查「日志」标签中的错误信息。

### 如何查看历史运行的报告

任务列表会自动扫描 `artifacts/runs/*/run_summary.json`，即使重启 Web Server 也能看到历史记录（但无法获取已完成任务的实时日志，因为进程已退出）。

### 如何直接打开 HTML 报告

在任务详情页，当状态为「已完成」时，顶部按钮组会出现「打开 HTML 报告」按钮，点击即可在新标签页查看原始 HTML 报告文件。也可直接访问 `/api/runs/{run_id}/report-html`。

### Docker 镜像构建失败

检查 Dockerfile 路径和 Docker daemon 是否正常启动。确保 `docker build` 命令在项目根目录可正常执行。

### Windows 下 mutmut 无法运行

Windows 原生不支持 mutmut 所需的 Unix 工具链。请在「新建任务」时勾选「Docker 执行」，并确保「镜像就绪」状态为绿色。

---

## 文件清单

```
cmd/utbench/main.go                   # web 子命令注册
internal/web/server.go                # HTTP Server + REST API + SSE
internal/web/run_manager.go           # 异步任务管理 + 日志捕获
internal/web/env.go                   # 环境检测（Docker/工具链）
internal/web/build_manager.go         # Docker 镜像构建任务管理
internal/web/docker_runner.go         # Docker 容器执行后端
internal/web/env_handlers.go          # /api/env 和 /api/env/build-image API
internal/web/static/index.html        # 前端 SPA (Alpine.js + Tailwind + Chart.js)
```

