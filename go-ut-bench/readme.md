# go-ut-bench

UT-Bench 的 Go CLI。当前版本既能评测纯模型 API，也能评测编码 Agent，并支持 `framework + model + optional skill` 的横向对比。

核心链路：

```text
sample -> subject generates tests -> compile/test/coverage/mutation -> report
```

其中 `subject` 是统一被测对象：

- `model_api + model + no_skill`：纯模型基线
- `cli_agent + framework + model + no_skill`：Agent 基线
- `cli_agent + framework + model + skill`：Agent + skill

## 当前能力

- 多语言评测：Python、Go、Java、C++
- 多维指标：编译、测试、覆盖率、变异测试
- 统一 subject 抽象：`framework__model__skill`
- 纯模型 API baseline：自动生成 `model_api__<model>__no_skill`
- 通用 CLI Agent 适配器：通过 `agents.yaml` 的命令模板驱动
- skill 注入：
  - `prompt_append`
  - `workspace_mount`
  - `agent_native` 预留
- Agent 过程产物：
  - trace
  - workspace diff
  - sandbox fingerprint
- 报告新增：
  - `agent_comparisons`
  - `skill_uplifts`

## 核心概念

### Subject

被测对象统一表示为：

```text
subject = framework + model + optional skill
```

示例：

- `model_api__deepseek-v4-flash__no_skill`
- `opencode__deepseek-v4-flash__no_skill`
- `opencode__deepseek-v4-flash__unit_test_skill`

### Skill

Skill 是通用能力包，不绑定某个 Agent 的原生 skill 机制。第一版支持两种注入：

- `prompt_append`：把 skill 说明追加到 prompt
- `workspace_mount`：把 skill 文件复制到 Agent 工作区

### CLI Agent

第一版 `cli_agent` 通过 `agents.yaml` 中的命令模板调用。UT-Bench 会为每个 `subject × sample` 准备独立工作区，并把模板渲染成 shell 命令执行。

最关键的模板变量：

- `{{.Model}}`：UT-Bench 模型名
- `{{.ModelID}}`：底层模型 ID
- `{{.PromptFile}}` / `{{.ContainerPrompt}}`
- `{{.OutputFile}}` / `{{.ContainerOutput}}`
- `{{.SourceFile}}`
- `{{.SkillDir}}` / `{{.ContainerSkillDir}}`
- `{{.Workspace}}`
- `{{.Language}}`
- `{{.SampleID}}`

## 快速开始

### 1. 构建

```bash
go build -o utbench ./cmd/utbench/
```

或：

```bash
docker build -t utbench:latest .
```

### 2. 纯模型 baseline

```bash
./utbench run \
  --models deepseek-v4-flash \
  --langs python \
  --config ./configs/models.yaml \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --class self_contained \
  --max-samples 2
```

这等价于运行 subject：

```text
model_api__deepseek-v4-flash__no_skill
```

### 3. Agent + Skill

先准备 Agent 配置文件，例如 [agents.example.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.example.yaml)。

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

如果不传 `--subjects`，系统会自动展开所有合法的 `framework × model × skill` 组合，并始终保留纯模型 baseline。

## CLI 命令

| 命令 | 说明 |
|------|------|
| `utbench run` | 完整流程：generate -> evaluate -> report |
| `utbench generate` | 仅生成测试 |
| `utbench evaluate` | 仅评测，输入 `generated_manifest.json` |
| `utbench report` | 仅报告，输入 `evaluation_result.json` |
| `utbench db` | SQLite 入库与查询 |
| `utbench dataset` | 数据集索引、manifest、validate |
| `utbench doctor` | 检查工具链并运行 canary |
| `utbench web` | 启动 Web 管理界面 |

## Docker 与沙箱

要分清两层：

- 外层 `docker run utbench:latest ...`：只是运行 UT-Bench 本身的执行环境
- 内层 `cli_agent` 的 `sandbox_mode: docker`：才是按 `subject × sample` 启动的 Agent 执行沙箱

当前实现里，CLI Agent 在 `sandbox_mode: docker` 时会：

- 为每个 `subject × sample` 创建独立工作区
- 启动一个独立容器
- 按语言选择固定的内层 Agent 镜像
- 平台托管样本依赖准备，如 `requirements.txt` / `go.mod` / `pom.xml`
- 执行前先跑 preflight，自检语言工具链
- 只挂载该工作区到 `/workspace`
- 可选 `--network none`
- 拦截环境漂移命令，如 `apt-get install`
- 记录 trace、diff、sandbox fingerprint
- 超时后结束该容器

如果外层 UT-Bench 本身也是用 Docker 启动的，那么这条链当前走的是 DOOD 方式：外层容器内的 `docker` CLI 控制宿主机 Docker daemon。具体挂载方式见 [docs/DOCKER_GUIDE.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/DOCKER_GUIDE.md)。

这已经比“整项目一个 Docker 容器”更接近真正沙箱，但还不是最终形态。后续还需要补：

- 更严格的镜像和依赖管理
- 更强的只读挂载边界
- 更细的资源治理
- 可选 gVisor / Firecracker / E2B 一类隔离增强

## 当前建议的首个真实 CLI Agent

如果要先打通一个真实 CLI Agent，优先建议 `OpenCode`，不是 `Claude Code`。

原因很直接：

- UT-Bench 的目标是 `framework × model` 的 N×M 组合
- `OpenCode` 更适合作为“Agent 外壳 + 任意底层模型”的统一入口
- `Claude Code` 更适合评测 Claude 自身工作流，不是最自然的多模型矩阵入口

所以顺序建议是：

1. 先做 `OpenCode`
2. 再补 `Claude Code`

## 输出结构

```text
artifacts/runs/<run-id>/
  generated/
    generated_manifest.json
    tests/
    prompts/
    metadata/
  evaluation/
    evaluation_result.json
  report/
    report_summary.json
    report.html
  agent_workspaces/
  logs/
  run_summary.json
```

更细的说明见 [docs/output-convention.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/output-convention.md)。

## 文档

- [快速开始](QUICKSTART.md)
- [Agent 升级状态](docs/AGENT_UPGRADE.md)
- [用户指南](docs/USER_GUIDE.md)
- [输出规范](docs/output-convention.md)
- [CLI 参数](docs/cli-spec.md)
- [Docker 使用](docs/DOCKER_GUIDE.md)
- [架构设计](docs/architecture-mvp.md)
