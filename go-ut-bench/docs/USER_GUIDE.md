# UT-Bench 用户指南

## 1. 定位

UT-Bench 现在评测的不是单一模型，而是统一被测对象 `subject`：

```text
subject = framework + model + optional skill
```

典型对比方式：

- `model_api + model + no_skill`
- `cli_agent + framework + model + no_skill`
- `cli_agent + framework + model + skill`

这允许你同时回答两类问题：

- Agent 相比纯模型 API 提升了多少
- 同一 Agent 加 skill 之后提升了多少

## 2. 支持范围

### 语言

| 语言 | 编译/测试 | 覆盖率 | 变异测试 |
|------|-----------|--------|----------|
| Python | pytest | coverage | mutmut |
| Go | go test | go test -cover | go-mutesting |
| Java | Maven/JUnit | JaCoCo | pitest |
| C++ | GoogleTest | gcov | mull |

### Subject 类型

| kind | 说明 |
|------|------|
| `model_api` | 纯模型 API baseline |
| `cli_agent` | 外部 CLI Agent，通过命令模板调用 |

后续可扩展 `http_agent`、`swe_agent` 一类适配器，但当前代码里还没有实现。

## 3. 运行模式

### 纯模型 baseline

```bash
./utbench run \
  --models deepseek-v4-flash \
  --langs python,go \
  --config ./configs/models.yaml \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --class self_contained
```

当只传 `--models` 时，系统会自动生成：

```text
model_api__<model>__no_skill
```

### Agent / Skill 模式

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

如果只传 `--agents-config` 不传 `--subjects`，系统会自动展开所有合法的：

```text
framework × model × skill
```

并自动保留 `model_api__<model>__no_skill`。

## 4. 关键参数

### `utbench run`

| 参数 | 说明 |
|------|------|
| `--config` | 模型配置文件 |
| `--models` | 模型列表 |
| `--agents-config` | Agent/Skill 配置文件 |
| `--subjects` | 显式选择 subject，格式 `framework__model__skill` |
| `--langs` | 语言列表 |
| `--dataset-root` | 数据集目录 |
| `--dataset-manifest` | 数据集 manifest |
| `--class` | 数据集类别，支持逗号分隔 |
| `--scenario` | 场景过滤 |
| `--level` | 数据集级别 |
| `--max-samples` | 样本上限 |
| `--mode` | `full` 或 `incremental` |
| `--reset-checkpoint` | 重置 checkpoint |
| `--dry-run` | 不调用真实模型 |
| `--reuse-generated` | 从 SQLite 复用历史 generated case |
| `--mutation-enabled` | 是否启用变异测试 |
| `--mutation-timeout` | 变异超时秒数 |
| `--test-timeout` | 测试执行超时秒数 |
| `--workers` | 并发 worker 数 |
| `--output-root` | 输出根目录 |
| `--run-id` | 指定 run ID |

### 分步命令

- `utbench generate --manifest` 不存在，`generate` 会直接产出 `generated_manifest.json`
- `utbench evaluate --manifest <generated_manifest.json>`
- `utbench report --evaluation <evaluation_result.json>`

## 5. `agents.yaml` 结构

示例见 [configs/agents.example.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.example.yaml)。

如果使用当前仓库提供的 `OpenCode` 示例，还要先构建对应的内层 Agent 镜像：

```bash
docker build -t utbench-agent-opencode:latest -f ./docker/agents/opencode/Dockerfile .
```

核心结构：

```yaml
models:
  - deepseek-v4-flash

frameworks:
  opencode:
    kind: cli_agent
    sandbox_mode: docker
    docker_image: utbench-agent-opencode:latest
    timeout_seconds: 600
    network_disabled: true
    cpu: "2"
    memory: 2g
    env:
      OPENCODE_CONFIG_CONTENT: |
        {"provider":{"utbench":{"options":{"baseURL":"{{.ModelEndpoint}}","apiKey":"{env:{{.ModelAPIKeyEnv}}}"}}}}
    env_from_host:
      - DEEPSEEK_API_KEY
    command: >
      opencode run --print-logs --dangerously-skip-permissions --model utbench/{{.ModelID}} "$(cat {{.ContainerPrompt}})"
    compatible_models:
      - deepseek-v4-flash
    compatible_languages:
      - python
      - go
    output_globs:
      - "test_*.py"
      - "*_test.go"

skills:
  unit_test_skill:
    version: "1"
    inject_mode: prompt_append
    instruction_path: ./skills/unit_test_skill.md
    compatible_frameworks:
      - opencode

subjects:
  - framework: opencode
    model: deepseek-v4-flash
    skill: unit_test_skill
```

### framework 字段

| 字段 | 说明 |
|------|------|
| `kind` | 目前主要是 `cli_agent` |
| `command` | 命令模板 |
| `sandbox_mode` | `docker` 或 `local` |
| `docker_image` | 内层 Agent 沙箱镜像 |
| `timeout_seconds` | 单样本 Agent 超时 |
| `output_globs` | 测试文件发现规则 |
| `env` | 传给 Agent 的环境变量 |
| `env_from_host` | 从宿主环境透传进 Agent 进程或容器的变量名 |
| `compatible_models` | 允许的模型列表 |
| `compatible_languages` | 允许的语言列表 |
| `network_disabled` | 是否禁网 |
| `cpu` / `memory` | 资源限制 |

### skill 字段

| 字段 | 说明 |
|------|------|
| `version` | skill 版本号 |
| `description` | skill 描述 |
| `instruction_path` | prompt 说明文件 |
| `files` | 复制进工作区的文件或目录 |
| `inject_mode` | `prompt_append` / `workspace_mount` / `agent_native` |
| `compatible_frameworks` | 兼容的 framework |
| `compatible_languages` | 兼容的语言 |

## 6. CLI Agent 的调用约定

当前实现不是为某个 Agent 写死命令，而是统一模板渲染。

### 模板变量

常用变量：

- `{{.Model}}`
- `{{.ModelID}}`
- `{{.Framework}}`
- `{{.SubjectID}}`
- `{{.Skill}}`
- `{{.Language}}`
- `{{.SampleID}}`
- `{{.Workspace}}`
- `{{.PromptFile}}`
- `{{.OutputFile}}`
- `{{.SourceFile}}`
- `{{.SkillDir}}`
- `{{.ContainerWorkdir}}`
- `{{.ContainerPrompt}}`
- `{{.ContainerOutput}}`
- `{{.ContainerSkillDir}}`

### 执行契约

UT-Bench 会给 Agent 一个明确约束：

- 只能在当前工作区内工作
- 不要修改原始源码行为
- 生成一个完整的单元测试文件
- 最终测试文件必须写到指定路径

如果 Agent 没在预期路径输出，UT-Bench 会按 `output_globs` 和工作区 diff 再查一次。

## 7. Docker 与沙箱的边界

很多人会把这两件事混在一起：

### 外层 Docker

```text
docker run utbench:latest ...
```

这是 UT-Bench 自己的运行环境，不等于 Agent 沙箱。

### 内层 Agent 沙箱

当 framework 配置为：

```yaml
sandbox_mode: docker
```

UT-Bench 会对每个 `subject × sample`：

1. 创建独立工作区
2. 启动独立容器
3. 只挂载该工作区到 `/workspace`
4. 可选禁网
5. 设置 CPU / memory / timeout
6. 回收容器

所以最终正确的隔离粒度，确实就是：

```text
一个 Agent 执行一个样本，对应一个独立工作区，最好再对应一个独立容器
```

这既是安全边界，也是公平边界。

### 当前状态

当前代码已经支持这个方向的第一版：

- 工作区按 `subject × sample` 拆分
- `sandbox_mode: docker` 时按样本起内层容器
- 记录 trace、diff、sandbox fingerprint

但还不能说“沙箱体系已经做完”，因为还缺：

- 更严格的镜像供应和固定 digest
- 更明确的只读挂载策略
- 更系统的网络、syscall、权限收敛
- 语言专用基础镜像
- 更强的隔离层，例如 gVisor / Firecracker / E2B

## 8. 报告新增内容

除了原有的编译、测试、覆盖率、变异得分，当前报告还会额外产出：

- `agent_comparisons`
  - 同模型、同样本下，Agent 相对 `model_api` baseline 的提升
- `skill_uplifts`
  - 同 framework + model、同样本下，skill 相对 `no_skill` 的提升

常见比较：

- `opencode__deepseek-v4-flash__no_skill` vs `model_api__deepseek-v4-flash__no_skill`
- `opencode__deepseek-v4-flash__unit_test_skill` vs `opencode__deepseek-v4-flash__no_skill`

## 9. 数据集与样本准备

### `self_contained`

直接使用单文件样本。对纯模型 API 和 Agent 都适用。

### `module_level`

如果样本旁边存在 module-level metadata，runner 会复制整个 workspace，再把目标文件交给 Agent。

这意味着 Agent 视角里拿到的是一个最小可操作 repo，而不是孤立源码片段。

## 10. 现在优先接哪个真实 CLI Agent

建议先接 `OpenCode`。

理由不是品牌偏好，而是适配目标：

- 你要的是 `framework × model` 的横评能力
- `OpenCode` 更适合作为通用 Agent 框架入口
- `Claude Code` 更适合作为 Claude 自身工作流的专项评测

所以建议顺序：

1. `OpenCode`
2. `Claude Code`

## 11. 常见问题

### `--models` 和 `--subjects` 要不要同时传

建议传。

原因是 subject 展开和模型配置加载都需要模型集合。第一版里，`agents.yaml` 也会对模型名做交集过滤。

### 只传 `--models` 能不能跑

可以。那就是纯模型 baseline 模式。

### 只传 `--agents-config` 不传 `--subjects` 能不能跑

可以。系统会自动展开全部合法组合。

### `sandbox_mode: local` 有什么用

主要用于本地调试和测试。正式跑 Agent 横评时，应该优先用 `sandbox_mode: docker`。
