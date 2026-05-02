# UT-Bench Agent 评测升级方案与实施状态

## 1. 目标

UT-Bench 原本评测的是“模型生成单元测试”的能力：

```text
sample -> model generates tests -> evaluate -> report
```

升级后的目标是评测统一被测对象 `subject`：

```text
subject = framework + model + optional skill
```

典型对象：

- `model_api__deepseek-v4-flash__no_skill`
- `opencode__deepseek-v4-flash__no_skill`
- `opencode__deepseek-v4-flash__unit_test_skill`

这样可以同时回答两类问题：

1. Agent 相比纯模型 API 提升了多少
2. 同一 Agent 加 skill 之后提升了多少

## 2. 当前总体完成情况

当前可以认为已经完成了 **第一版架构和主链路改造**，但还没有完成“真实线上可批量横评”的最后落地。

完成度拆分：

- 核心数据结构与配置：已完成
- 纯模型 baseline 兼容：已完成
- AgentAdapter 统一抽象层：已完成
- 通用 CLI Agent 主链路：已完成
- Skill v2 专业版：已完成
- 扩展 Agent Trace：已完成
- GeneratedCase Agent 追踪摘要：已完成
- 运行时显示改进（Agent 感知）：已完成
- 每样本 Docker 沙箱执行：已完成第一版
- 按语言选择内层 Agent 镜像：已完成第一版
- 沙箱 preflight 与环境漂移拦截：已完成第一版
- 报告与数据库扩展：已完成第一版
- OpenCode 作为首个真实 framework 的配置入口：已完成
- 外层 Docker + DOOD 路径：已完成文档和代码接线
- 真实端到端 smoke run：已完成（agent_smoke_002）
- Claude Code framework：未开始
- 更强安全沙箱后端：未开始

## 3. 已完成内容

### 3.1 被测对象从 model 升级为 subject

已经新增：

- `SubjectSpec`
- `SkillSpec`
- `RunSpec.Subjects`
- `RunSpec.AgentsConfigPath`

现在系统支持：

- 旧模式：只传 `--models`
- 新模式：传 `--agents-config` + `--subjects`

只传 `--models` 时，系统会自动生成：

```text
model_api__<model>__no_skill
```

### 3.2 Agent/Skill 配置解析

已经新增配置解析模块：

- `internal/agentconfig`

支持：

- `framework × model × skill` 展开
- 自动加入 `model_api` baseline
- `compatible_models`
- `compatible_languages`
- `compatible_frameworks`
- `env`
- `env_from_host`

### 3.3 AgentAdapter 统一抽象层

已经新增：

- `internal/runner/adapter.go` — `AgentAdapter` 接口定义
- `internal/runner/adapter_cli.go` — CLI Agent 适配器实现
- `internal/runner/adapter_model.go` — 纯模型 API 适配器实现

`AgentAdapter` 接口是所有 Agent 生成器的统一入口：

```go
type AgentAdapter interface {
    Generate(ctx context.Context, req AgentGenerateRequest) AgentGenerateResult
}
```

当前实现：

- `modelAPIAdapter` — 复用原有 HTTP LLM API 逻辑
- `cliAgentAdapter` — 在沙箱中执行 CLI Agent，收集丰富化 trace

新增 framework 只需实现 `AgentAdapter` 接口，不用改 runner 主逻辑。

### 3.4 通用 CLI Agent 执行链路

CLI Agent 单样本执行流程：

1. 准备样本 prompt
2. 准备独立 workspace
3. 注入 skill prompt / skill files
4. 渲染 framework command
5. 渲染 framework env
6. 执行 Agent（通过 SandboxRunner）
7. 实时收集 trace（工具调用、文件读写、命令执行）
8. 收集输出测试文件
9. 落盘 trace、diff、metadata、manifest

### 3.5 Skill v2 专业版

已支持注入模式：

- `prompt_append` — 将 skill 指令追加到 prompt
- `workspace_mount` — 将 skill 文件复制到 workspace
- `agent_native` — 预留，由 Agent 框架自行处理

skill v2 目录结构：

```
configs/skills/unit_test_skill/
├── instructions.md       # 主指令：分析阶段、测试结构、命名规范、断言质量
├── checklist.md          # 预检清单：结构、断言、Mock、隔离、完整性
└── examples/
    └── python_example.py # 示例：带 Mock 的 pytest 测试
```

skill v2 能力：

- 结构化的测试生成方法论（分析→结构→命名→断言→Mock）
- 语言特定的模板（Python/Go/Java/C++）
- 预检清单确保质量
- 示例代码供 Agent 参考

### 3.6 扩展 Agent Trace

trace 文件现在是完整的 `AgentTrace` 结构化 JSONL，包含：

基础信息：
- `subject_id`, `framework`, `model`, `skill`, `sample_id`, `language`

执行信息：
- `command`, `exit_code`, `duration_ms`, `started_at`, `finished_at`

Token 消耗：
- `prompt_tokens`, `completion_tokens`, `total_tokens`, `estimated_cost`

Agent 交互过程：
- `interaction_count` — 交互轮次
- `tool_calls` — 工具调用序列（tool, input, output, duration_ms, success）
- `files_read` — Agent 读取的文件列表
- `files_written` — Agent 写入/修改的文件列表
- `commands_executed` — Agent 执行的 shell 命令

输出：
- `stdout`, `stderr`, `workspace_diff`

解析逻辑：从 Agent 的 stdout/stderr 中自动提取工具调用、文件操作和命令执行记录。

### 3.7 GeneratedCase Agent 追踪摘要

`GeneratedCase` 已新增 5 个 Agent 追踪摘要字段，用于运行时显示和报告聚合：

- `interaction_count` — Agent 交互轮次
- `tool_call_count` — 工具调用次数
- `files_read_count` — 读取的文件数
- `files_write_count` — 写入的文件数
- `command_count` — 执行的命令数

这些字段从 `AgentTrace` 中提取，通过 `agentTraceSummary` 结构从 `generateWithSubject` 传递到 `generateOne`，最终写入 `GeneratedCase`。

### 3.8 运行时显示改进

运行时终端输出已升级为 Agent 感知模式：

启动阶段展示 subject 详情：
```
[STAGE] 生成测试
4 样本 × 3 被测对象 = 12 任务 | Workers: 4
   [cli_agent] opencode__deepseek-v4-flash__no_skill | model=deepseek-v4-flash | skill=no_skill
   [model_api] model_api__deepseek-v4-flash__no_skill | model=deepseek-v4-flash | skill=no_skill
```

每行任务进度展示 framework、skill、耗时、Agent 追踪摘要：
```
[1/12]  opencode | deepseek-v4-flash | no_skill | python | boundary_000 | OK | 02:56 | 2249 tok | 3轮,5工具,2写入,1读取
[9/12]  model_api | deepseek-v4-flash | no_skill | python | boundary_000 | OK | 00:25 | 1834 tok
```

多 framework 时统计面板自动分组：
```
[opencode] 7/8 成功 | 截断 0 | 平均 120s | 平均 3200 tok
[model_api] 4/4 成功 | 截断 0 | 平均 25s | 平均 1800 tok
```

### 3.9 每样本沙箱执行

当前已经实现的隔离粒度是：

```text
一个 subject 执行一个 sample
= 一个独立 workspace
= 一个独立执行单元
```

对于 `sandbox_mode: docker`：

- 每个样本会起一个独立内层容器
- 只挂当前 workspace 到 `/workspace`
- 可设置 CPU / memory / timeout
- 可选择禁网

注意：

- 对 OpenCode 这类需要主动访问模型 API 的 framework，真实运行时必须允许网络
- 所以当前 `agents.example.yaml` 已修正为 `network_disabled: false`

### 3.10 SandboxRunner 抽象

已经新增：

- `internal/runner/sandbox.go`

当前后端：

- `local`
- `docker`

这意味着“样本沙箱执行”已经从 runner 主逻辑中抽离，后续可以继续接：

- `docker_remote`
- `gvisor`
- `firecracker`
- `e2b`

而不需要重新改写生成主流程。

### 3.11 OpenCode 首个真实 framework 入口

已经新增：

- `configs/agents.example.yaml`
- `docker/agents/opencode/Dockerfile`

这不是把 OpenCode 变成唯一基座，而是把它作为第一个真实 `framework` 落地。

### 3.12 评测环境完整性

目前已经开始把“环境完整性由平台负责”落实到代码和配置：

- framework 支持 `docker_images`，可以按语言选择内层沙箱镜像
- framework 支持 `preflight`，在 Agent 执行前检查语言运行时和基础工具
- framework 支持 `forbidden_command_patterns`，识别环境漂移行为

默认理念是：

```text
固定镜像 + 固定工具链 + 每样本独立 workspace + 每样本独立容器
```

而不是让 Agent 自己在运行中安装依赖。

当前 `OpenCode` 示例已经改成按语言镜像：

- `utbench-agent-opencode-python`
- `utbench-agent-opencode-go`
- `utbench-agent-opencode-java`
- `utbench-agent-opencode-cpp`

并提供了典型 preflight：

- Python: `python3 --version`, `pytest --version`
- Go: `go version`
- Java: `java -version`, `mvn -version`
- C++: `g++ --version`, `cmake --version`

如果 preflight 失败，会直接报 `sandbox_preflight_error`。

如果 Agent 尝试执行：

- `apt-get install`
- `pip install`
- `npm install`

会报 `sandbox_policy_error`。

另外，平台现在会在 Agent 执行前准备样本依赖：

- Python `repo_level` 元数据里的 `requirements`
- workspace 根目录的 `requirements.txt` / `requirements-dev.txt`
- `go.mod`
- `pom.xml`

这一步属于平台责任，不属于 Agent 自主行为；如果失败，会报 `sample_env_prepare_error`。

当前推荐比较对象：

- `model_api__deepseek-v4-flash__no_skill`
- `opencode__deepseek-v4-flash__no_skill`
- `opencode__deepseek-v4-flash__unit_test_skill`

### 3.13 报告、结果和数据库扩展

已经扩展：

- `GeneratedCase`
- `EvaluationResult`
- `ReportPayload`
- SQLite schema

新增字段包括：

- `subject_id`
- `subject_kind`
- `agent_framework`
- `agent_model`
- `skill_name`
- `skill_version`
- `trace_path`
- `workspace_diff_path`
- `sandbox_fingerprint`

报告中已新增：

- `agent_comparisons`
- `skill_uplifts`

### 3.14 外层 Docker + DOOD 路径

由于你们通常用 Docker 启动 UT-Bench，而每个样本又希望起一个内层 Agent 容器，所以当前外层部署方案是 **DOOD 过渡方案**。

当前已经完成：

- 外层 benchmark 镜像安装 `docker.io`
- 文档明确要求挂 Docker socket
- 说明了内层容器才是真正的样本级 Agent 沙箱

## 4. 已修改 / 新增文件

### 4.1 新增文件

- `configs/agents.example.yaml`
- `configs/skills/unit_test_skill/` — skill v2 目录
  - `instructions.md`
  - `checklist.md`
  - `examples/python_example.py`
- `docker/agents/opencode/Dockerfile`
- `docs/AGENT_UPGRADE.md`
- `internal/agentconfig/config.go`
- `internal/agentconfig/config_test.go`
- `internal/runner/adapter.go` — AgentAdapter 接口定义
- `internal/runner/adapter_cli.go` — CLI Agent 适配器（含 trace 解析）
- `internal/runner/adapter_model.go` — 纯模型 API 适配器
- `internal/runner/sandbox.go`
- `internal/runner/subjects.go`
- `internal/runner/subjects_test.go`

### 4.2 已修改的核心文件

- `Dockerfile`
- `readme.md`
- `QUICKSTART.md`
- `docs/DOCKER_GUIDE.md`
- `docs/USER_GUIDE.md`
- `docs/output-convention.md`
- `cmd/utbench/main.go`
- `internal/contracts/spec.go`
- `internal/contracts/results.go`
- `internal/dataset/service.go`
- `internal/evaluator/service.go`
- `internal/reporter/service.go`
- `internal/reporter/service_test.go`
- `internal/runner/checkpoint_test.go`
- `internal/runner/service.go`
- `internal/obs/progress.go`
- `internal/store/sqlite.go`

## 5. 端到端 Smoke Run 结果

### 5.1 `agent_smoke_002` 首次成功

配置：3 个 subject × 4 个 Python self_contained 样本 = 12 个任务。

| 排名 | Subject | 编译 | 测试通过率 | 覆盖率 | 变异分 | 综合分 | 平均耗时 |
|------|---------|------|-----------|--------|--------|--------|---------|
| 1 | `opencode__deepseek-v4-flash__no_skill` | 100% | 50% | 98.8% | 46.7% | 0.737 | 175.8s |
| 2 | `model_api__deepseek-v4-flash__no_skill` | 100% | 50% | 72.9% | 45.7% | 0.687 | 25.6s |
| 3 | `opencode__deepseek-v4-flash__unit_test_skill` | 100% | 50% | 57.3% | 45.7% | 0.656 | 70.0s |

Agent vs Model API 对比（同模型、同 skill）：

- 覆盖率 delta: **+25.9%**（Agent 显著更好）
- 变异分 delta: +1.0%
- 耗时 delta: +150s（Agent 慢 ~7 倍，预期中）

产物全部按预期落盘：

- `generated/tests/` — 12 个生成的测试文件
- `agent_workspaces/` — 8 个 OpenCode 样本工作区
- `agent_traces/*.trace.jsonl` — 完整 Agent 操作轨迹
- `agent_traces/*.diff.json` — workspace 文件变更
- `evaluation/evaluation_result.json` — 全部评测结果
- `report/report_summary.json` — 含 `agent_comparisons` + `skill_uplifts`
- `report/report.html` — HTML 报告

## 6. 目前还没完成的部分

### 6.1 Claude Code framework

还没有接 `claude_code`。

未来要做的是：

- 新增 `claude_code` framework 配置
- 新增对应镜像
- 验证它是否支持你们要的模型切换能力

### 6.2 更强安全沙箱

当前还是：

```text
outer benchmark container + DOOD + per-sample inner docker container
```

这能满足第一版工程验证，但还不是最终高强度安全方案。

后续候选：

- remote Docker daemon
- gVisor
- Firecracker microVM
- E2B

## 7. 你要怎么做端到端验证

建议先做一个 **最小 smoke run**，只跑：

- 1 个语言：`python`
- 1 个样本
- 3 个 subject：
  - `model_api__deepseek-v4-flash__no_skill`
  - `opencode__deepseek-v4-flash__no_skill`
  - `opencode__deepseek-v4-flash__unit_test_skill`

### 7.1 准备 `.env`

至少保证有：

```dotenv
DEEPSEEK_API_KEY=...
```

如果你后面换别的模型，也补对应 key。

### 7.2 构建外层 benchmark 镜像

在 [go-ut-bench](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench) 目录执行：

```bash
docker build -t utbench:latest .
```

### 7.3 构建内层 OpenCode Agent 镜像

```bash
docker build -t utbench-agent-opencode-python:latest -f ./docker/agents/opencode/python.Dockerfile .
docker build -t utbench-agent-opencode-go:latest -f ./docker/agents/opencode/go.Dockerfile .
docker build -t utbench-agent-opencode-java:latest -f ./docker/agents/opencode/java.Dockerfile .
docker build -t utbench-agent-opencode-cpp:latest -f ./docker/agents/opencode/cpp.Dockerfile .
```

### 7.4 运行最小 smoke run

建议在 Linux/macOS/WSL Bash 下执行，因为 DOOD 的 Docker socket 挂载更直接：

```bash
docker run --rm --env-file .env \
  -e UTBENCH_SANDBOX_HOST_OUTPUT_ROOT="$(pwd)/artifacts" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek-v4-flash \
    --langs python \
    --config /app/configs/models.yaml \
    --agents-config /app/configs/agents.example.yaml \
    --subjects model_api__deepseek-v4-flash__no_skill,opencode__deepseek-v4-flash__no_skill,opencode__deepseek-v4-flash__unit_test_skill \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --class self_contained \
    --max-samples 1 \
    --run-id agent_smoke_001
```

### 7.5 运行成功后你要检查什么

看这些目录和文件：

```text
artifacts/runs/agent_smoke_001/
  generated/generated_manifest.json
  generated/tests/
  generated/metadata/
  generated/metadata/agent_traces/
  evaluation/evaluation_result.json
  report/report_summary.json
  agent_workspaces/
```

重点检查：

1. `generated_manifest.json`
   - 是否出现 3 个 subject 的 case

2. `generated/tests/`
   - 是否真的生成了测试文件

3. `generated/metadata/agent_traces/...trace.jsonl`
   - 是否记录了 OpenCode 执行轨迹

4. `generated/metadata/agent_traces/...diff.json`
   - 是否记录了 workspace 改动

5. `evaluation/evaluation_result.json`
   - 是否有 `subject_id`、`agent_framework`、`skill_name`

6. `report/report_summary.json`
   - 是否有 `agent_comparisons`
   - 是否有 `skill_uplifts`

### 7.6 如果失败，优先看哪里

先看：

- `artifacts/runs/<run-id>/generated/metadata/*.response.json`
- `artifacts/runs/<run-id>/generated/metadata/agent_traces/*.trace.jsonl`
- `artifacts/runs/<run-id>/logs/`

常见失败点：

1. 外层没挂 Docker socket
   - 内层容器根本起不来

2. `utbench-agent-opencode:latest` 没构建
   - 找不到内层镜像

3. API key 没透传
   - OpenCode 能启动，但调模型失败

4. 内层网络被禁
   - Agent 无法访问模型 API

5. 没传 `UTBENCH_SANDBOX_HOST_OUTPUT_ROOT`
   - DOOD 下内层容器会挂到空 workspace
   - 典型报错就是 `/workspace/utbench_agent_prompt.md: No such file or directory`

6. OpenCode CLI 包名或行为和镜像内实际安装结果不一致
   - 这时看 trace 里的 stderr

## 7. 当前建议的下一步

优先顺序建议：

1. 你先按 6 节做一次真实 smoke run
2. 把失败点或运行结果贴回来
3. 我继续补：
   - OpenCode 真实适配修正
   - Claude Code framework
   - 更严格的沙箱 backend

