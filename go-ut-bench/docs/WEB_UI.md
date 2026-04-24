# UTBench Web 管理后台

基于 Go 内嵌 HTTP Server + 前端 SPA 的 Benchmark 运行管理界面。

---

## 快速启动

```bash
# 在项目根目录构建 CLI
go build -o utbench ./cmd/utbench/

# 启动 Web 管理后台（默认端口 8080）
./utbench web

# 指定端口和配置
./utbench web --addr :9090 --config ./configs/models.yaml --dataset-root ./datasets
```

打开浏览器访问 `http://localhost:8080` 即可进入管理界面。

---

## 命令参数

```
utbench web [flags]

  --addr        HTTP 监听地址，默认 :8080
  --config      模型配置文件路径，默认 ./configs/models.yaml
  --dataset-root 数据集根目录，默认 ./datasets
  --output-root  产物根目录，默认 ./artifacts
  --db-path      SQLite 数据库路径，默认 ./storage/utbench.db
```

---

## 功能概览

### Dashboard

- 总运行数 / 进行中 / 已完成 / 失败 统计卡片
- 最近 10 次运行速览，点击跳转到详情页

### 新建运行

交互式表单，完整覆盖 CLI 的所有运行参数：

| 参数 | 控件 | 说明 |
|------|------|------|
| Run ID | 文本输入 | 留空自动生成 |
| 运行模式 | 下拉选择 | `full`（全量）/ `incremental`（断点续跑） |
| 模型 | 多选卡片 | 从 `configs/models.yaml` 读取，仅显示 `enabled: true` 的模型 |
| 语言 | 多选卡片 | python / go / java / cpp |
| Class | 下拉选择 | `self_contained` / `module_level` |
| Scenario | 下拉选择 | `boundary` / `simple_function` / `complex_dependency` / `interface_mock` |
| Level | 文本输入 | 可选，如 `l1`、`l2` |
| Max Samples | 数字输入 | 每场景样本上限，0 = 不限 |
| Workers | 数字输入 | 并发 worker 数，0 = 自动 |
| Dry Run | 复选框 | 跳过真实 API 调用，验证流程 |
| 变异测试 | 复选框 | 开启 mutmut/gremlins/pitest/mull |
| 变异超时 | 数字输入 | 秒数，默认 1800 |
| 写入数据库 | 复选框 | 完成后自动 ingest 到 SQLite |

提交后自动跳转到运行详情页，开始实时跟踪日志。

### 所有运行

- 完整运行历史列表（合并内存中的活跃运行 + `artifacts/runs/*/run_summary.json` 扫描）
- 支持按 **run_id / 模型名 / 语言** 文本过滤
- 支持按 **状态** 过滤：pending / running / completed / failed
- 显示模型、语言、场景、样本上限、开始时间、耗时

### 运行详情

#### 实时日志（SSE 推送）

- 后台 goroutine 运行评测，通过 Server-Sent Events 向前端推送结构化日志
- 日志区域自动滚动到底部
- 运行结束后 SSE 自动关闭

#### 评测报告

运行完成后（Status = `completed`）可切换到「评测报告」标签：

- **汇总卡片**：总样本数、编译通过率、测试通过率、平均行覆盖率、平均变异得分
- **模型横向对比柱状图**（Chart.js）：编译通过率 / 测试通过率 / 行覆盖率 / 变异得分
- **模型排名表**：按综合得分排序，含编译/测试/覆盖/变异/延迟/Token 数据
- **失败样本摘要**：按错误类型聚合，含示例信息

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
│  /api/config   GET  →  models.yaml 解析      │
│  /api/runs     GET  →  运行列表聚合          │
│  /api/runs     POST →  提交新运行            │
│  /api/runs/{id}      GET  →  运行详情        │
│  /api/runs/{id}/events  GET  →  SSE 日志流   │
│  /api/runs/{id}/report  GET  →  评测报告     │
│  /                    →  静态 SPA 文件       │
└─────────────────────────────────────────────┘
                      ↑
┌─────────────────────────────────────────────┐
│         RunManager (internal/web)            │
│  ─────────────────────────────────────────  │
│  内存运行状态管理 + goroutine 异步执行        │
│  lineWriter → 日志捕获 → SSE 广播            │
│  复用 orchestrator.Run 完成完整流水线        │
└─────────────────────────────────────────────┘
```

### 关键组件

| 文件 | 职责 |
|------|------|
| `internal/web/server.go` | HTTP Server、REST API 路由、CORS、静态文件服务 |
| `internal/web/run_manager.go` | 异步运行管理器：goroutine 执行、日志捕获（`lineWriter`）、SSE 订阅广播 |
| `internal/web/static/index.html` | 完整 SPA：Alpine.js 状态管理 + Tailwind 样式 + Chart.js 可视化 |

---

## API 参考

### GET /api/config

返回当前系统配置，用于前端动态渲染模型/语言/场景选项。

```json
{
  "models": [
    { "name": "deepseek", "provider": "deepseek", "model_id": "deepseek-chat", "enabled": true }
  ],
  "languages": ["python", "go", "java", "cpp"],
  "scenarios": ["boundary", "simple_function", "complex_dependency", "interface_mock"],
  "classes": ["self_contained", "module_level"],
  "dataset_root": "./datasets",
  "config_path": "./configs/models.yaml"
}
```

### GET /api/runs

返回所有运行记录（内存 + 磁盘扫描）。

```json
[
  {
    "run_id": "20240424T123456.123456789Z",
    "status": "running",
    "started_at": "2024-04-24T12:34:56Z",
    "ended_at": null,
    "error": "",
    "spec": { "models": ["deepseek"], "languages": ["python"], ... }
  }
]
```

### POST /api/runs

提交新运行。

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
  "ingest": true
}
```

**响应：**

```json
{ "run_id": "my_run_001", "status": "pending" }
```

### GET /api/runs/{run_id}

获取运行详情（含当前日志快照）。

```json
{
  "run_id": "my_run_001",
  "status": "running",
  "started_at": "2024-04-24T12:34:56Z",
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

这意味着 Web 提交的运行和 CLI 提交的运行产物完全兼容，可以互相查看。

---

## 常见问题

### Web 启动后浏览器无法访问

检查防火墙或绑定地址。默认绑定 `:8080`（所有接口），若只想本地访问：

```bash
./utbench web --addr 127.0.0.1:8080
```

### 运行页面日志不更新

SSE 连接在运行结束后会自动关闭。如果运行状态卡在 `pending` 或 `running`，检查后端日志是否有报错。

### 报告页面显示「报告将在完成后可用」

只有 `status === 'completed'` 时报告才会生成。若运行失败，请检查「实时日志」中的错误信息。

### 如何查看历史运行的报告

运行列表会自动扫描 `artifacts/runs/*/run_summary.json`，即使重启 Web Server 也能看到历史记录（但无法获取已完成运行的实时日志，因为进程已退出）。

---

## 文件清单

```
cmd/utbench/main.go              # web 子命令注册
internal/web/server.go           # HTTP Server + REST API + SSE
internal/web/run_manager.go      # 异步运行管理 + 日志捕获
internal/web/static/index.html   # 前端 SPA (Alpine.js + Tailwind + Chart.js)
```
