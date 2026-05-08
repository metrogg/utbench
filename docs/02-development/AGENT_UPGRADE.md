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
- 通用 CLI Agent 主链路：已完成
- Skill 注入机制：已完成第一版
- 每样本 Docker 沙箱执行：已完成第一版
- 报告与数据库扩展：已完成第一版
- OpenCode 作为首个真实 framework 的配置入口：已完成
- 外层 Docker + DOOD 路径：已完成文档和代码接线
- 真实端到端 smoke run：未完成，需你本地执行验证
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

### 3.3 通用 CLI Agent 执行链路

已经新增：

- `internal/runner/subjects.go`

CLI Agent 单样本执行流程现在是：

1. 准备样本 prompt
2. 准备独立 workspace
3. 注入 skill prompt / skill files
4. 渲染 framework command
5. 渲染 framework env
6. 执行 Agent
7. 收集输出测试文件
8. 落盘 trace、diff、metadata、manifest

### 3.4 Skill 第一版机制

已支持：

- `prompt_append`
- `workspace_mount`
- `agent_native` 预留

并新增示例 skill：

- `configs/skills/unit_test_skill.md`

### 3.5 每样本沙箱执行

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

### 3.6 SandboxRunner 抽象

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

### 3.7 OpenCode 首个真实 framework 入口

已经新增：

- `configs/agents.example.yaml`
- `docker/agents/opencode/Dockerfile`

这不是把 OpenCode 变成唯一基座，而是把它作为第一个真实 `framework` 落地。

当前推荐比较对象：

- `model_api__deepseek-v4-flash__no_skill`
- `opencode__deepseek-v4-flash__no_skill`
- `opencode__deepseek-v4-flash__unit_test_skill`

### 3.8 报告、结果和数据库扩展

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

### 3.9 外层 Docker + DOOD 路径

由于你们通常用 Docker 启动 UT-Bench，而每个样本又希望起一个内层 Agent 容器，所以当前外层部署方案是 **DOOD 过渡方案**。

当前已经完成：

- 外层 benchmark 镜像安装 `docker.io`
- 文档明确要求挂 Docker socket
- 说明了内层容器才是真正的样本级 Agent 沙箱

## 4. 已修改 / 新增文件

### 4.1 新增文件

- `configs/agents.example.yaml`
- `configs/skills/unit_test_skill.md`
- `docker/agents/opencode/Dockerfile`
- `go-ut-bench/docs/04-development/AGENT_UPGRADE.md`
- `internal/agentconfig/config.go`
- `internal/agentconfig/config_test.go`
- `internal/runner/sandbox.go`
- `internal/runner/subjects.go`
- `internal/runner/subjects_test.go`

### 4.2 已修改的核心文件

- `Dockerfile`
- `readme.md`
- `go-ut-bench/docs/01-user-guides/USER_GUIDE.md`
- `go-ut-bench/docs/03-operations/DOCKER_GUIDE.md`
- `go-ut-bench/docs/02-design/output-convention.md`
- `cmd/utbench/main.go`
- `internal/contracts/spec.go`
- `internal/contracts/results.go`
- `internal/dataset/service.go`
- `internal/evaluator/service.go`
- `internal/reporter/service.go`
- `internal/reporter/service_test.go`
- `internal/runner/checkpoint_test.go`
- `internal/runner/service.go`
- `internal/store/sqlite.go`

## 5. 目前还没完成的部分

### 5.1 真实端到端 smoke run

代码和配置已经准备好了，但还没有在你们真实环境里完成一次：

- 外层 `utbench:latest`
- 内层 `utbench-agent-opencode:latest`
- 真实 API key
- 真实样本
- 真实产物落盘检查

这一步需要你本地执行。

### 5.2 Claude Code framework

还没有接 `claude_code`。

未来要做的是：

- 新增 `claude_code` framework 配置
- 新增对应镜像
- 验证它是否支持你们要的模型切换能力

### 5.3 更强安全沙箱

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

## 6. 你要怎么做端到端验证

建议先做一个 **最小 smoke run**，只跑：

- 1 个语言：`python`
- 1 个样本
- 3 个 subject：
  - `model_api__deepseek-v4-flash__no_skill`
  - `opencode__deepseek-v4-flash__no_skill`
  - `opencode__deepseek-v4-flash__unit_test_skill`

### 6.1 准备 `.env`

至少保证有：

```dotenv
DEEPSEEK_API_KEY=...
```

如果你后面换别的模型，也补对应 key。

### 6.2 构建外层 benchmark 镜像

在 [go-ut-bench](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench) 目录执行：

```bash
docker build -t utbench:latest .
```

### 6.3 构建内层 OpenCode Agent 镜像

```bash
docker build -t utbench-agent-opencode:latest -f ./docker/agents/opencode/Dockerfile .
```

### 6.4 运行最小 smoke run

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

### 6.5 运行成功后你要检查什么

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

### 6.6 如果失败，优先看哪里

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
