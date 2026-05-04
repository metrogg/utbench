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
| `--db-path` | SQLite 资产库路径，默认 `./storage/utbench.db` |
| `--reuse-generated` | 复用相同资产 key 的历史生成结果，默认开启 |
| `--reuse-evaluation` | 复用相同 evaluation key 的历史评测结果，默认关闭 |
| `--max-samples` | 样本上限 |
| `--mode` | `full` 或 `incremental` |
| `--reset-checkpoint` | 重置 checkpoint |
| `--dry-run` | 不调用真实模型 |
| `--mutation-enabled` | 是否启用变异测试 |
| `--mutation-timeout` | 变异超时秒数 |
| `--test-timeout` | 测试执行超时秒数 |
| `--workers` | 并发 worker 数 |
| `--output-root` | 输出根目录 |
| `--run-id` | 指定 run ID |

## 5. 资产复用

UT-Bench 现在会把生成结果和评测结果索引到本地 SQLite。生成复用默认开启：

```bash
./utbench run \
  --models deepseek-v4-flash \
  --langs cpp \
  --config ./configs/models.yaml \
  --agents-config ./configs/agents.example.yaml \
  --subjects model_api__deepseek-v4-flash__no_skill,opencode__deepseek-v4-flash__no_skill \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --db-path ./storage/utbench.db \
  --reuse-generated=true \
  --reuse-evaluation=false
```

生成复用按 `subject_version_id + sample_uid + prompt + dependency + sandbox/env fingerprint` 判断，不会跨 subject、skill、framework 或环境复用。命中复用时，当前 run 仍会写出完整 manifest，并标注来源 run/case。

当前 generation 和 evaluation 都会在阶段开始前先批量扫描可复用资产，进度头会显示 `预判可复用: N`，然后再把剩余任务下发给 worker。

评测复用目前需要显式开启：

```bash
./utbench run ... --reuse-evaluation=true
```

建议在 evaluator 环境和 mutation 配置稳定后再打开。

`evaluation_key` 现在已经包含 `evaluation_env_fingerprint`。这个 fingerprint 会记录本地工具链版本或 Docker digest，因此不同评测环境不会直接复用同一条评测结果。

### 查询资产

```bash
./utbench assets subjects --db-path ./storage/utbench.db
./utbench assets generations --subject opencode__deepseek-v4-flash__no_skill --lang go
./utbench assets evaluations --subject opencode__deepseek-v4-flash__no_skill --lang go --sample boundary_000
./utbench assets explain-reuse --subject opencode__deepseek-v4-flash__no_skill --lang go --sample boundary_000
```

`explain-reuse` 会输出当前数据库里最新的可复用生成候选、`generation_key`、来源 run 和 artifact 路径。
如果要手工比对当前环境，可以额外传 `--generation-key`、`--dependency-fingerprint`、`--generation-env-fingerprint`。

### Web 资产管理

启动 Web UI 后：

- “新建任务”页支持**自由组合** Framework × Model × Skill，每行可独立选择，可添加多行，无需预定义 subject。
- `--agents-config` 会自动从 `--config` 目录探测，通常无需手动指定。
- “任务列表”和”任务详情”都以 `subject` 为主展示，不再只看模型名。
- “注册表”页会同时展示 `models / frameworks / skills / subjects`。
- 进入”资产管理 → 资产审计”，可以直接查看：
  - `subject` 资产列表
  - `subject_versions`
  - 生成资产历史
  - 评测资产历史
  - `subject -> language -> sample` 树形历史视图
  - 单个 `subject + language + sample` 的复用解释

### 分步命令

- `utbench generate --manifest` 不存在，`generate` 会直接产出 `generated_manifest.json`
- `utbench evaluate --manifest <generated_manifest.json>`
- `utbench report --evaluation <evaluation_result.json>`

## 6. `agents.yaml` 结构

示例见 [configs/agents.example.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.example.yaml)。

如果使用当前仓库提供的 `OpenCode` 示例，还要先构建对应的 Agent 沙箱镜像（详见 [DOCKER_GUIDE.md](./DOCKER_GUIDE.md)）：

```bash
docker build -t utbench-agent-opencode-python:latest -f docker/agents/opencode/python.Dockerfile docker/agents/opencode/
docker build -t utbench-agent-opencode-go:latest -f docker/agents/opencode/go.Dockerfile docker/agents/opencode/
docker build -t utbench-agent-opencode-java:latest -f docker/agents/opencode/java.Dockerfile docker/agents/opencode/
docker build -t utbench-agent-opencode-cpp:latest -f docker/agents/opencode/cpp.Dockerfile docker/agents/opencode/
```

核心结构（新版 sandbox 配置块）：

```yaml
models:
  - deepseek-v4-flash

frameworks:
  opencode:
    kind: cli_agent
    sandbox:
      provider: docker
      mode: docker
      images:
        python: utbench-agent-opencode-python:latest
        go: utbench-agent-opencode-go:latest
        java: utbench-agent-opencode-java:latest
        cpp: utbench-agent-opencode-cpp:latest
      timeout_seconds: 600
      network_disabled: false
      cpu: "2"
      memory: 2g
    preflight:
      python:
        - python3 --version
        - pytest --version
      go:
        - go version
    forbidden_command_patterns:
      - apt-get install
      - pip install
    env:
      OPENCODE_CONFIG_CONTENT: |
        {"provider":{"utbench":{"options":{"baseURL":"{{.ModelEndpoint}}","apiKey":"{env:{{.ModelAPIKeyEnv}}}"}}}}
    env_from_host:
      - DEEPSEEK_API_KEY
    command: >
      mkdir -p "$XDG_CONFIG_HOME/opencode" "$XDG_DATA_HOME" "{{.ContainerWorkdir}}/.utbench/opencode" &&
      printf '%s' "$OPENCODE_CONFIG_CONTENT" > "$XDG_CONFIG_HOME/opencode/config.json" &&
      PROMPT="$(cat {{.ContainerPrompt}})" &&
      opencode run --print-logs --dangerously-skip-permissions --model utbench/{{.ModelID}} "$PROMPT"
    compatible_models:
      - deepseek-v4-flash
    compatible_languages:
      - python
      - go
      - java
      - cpp
    output_globs:
      - "test_*.py"
      - "*_test.go"
      - "*Test.java"
      - "*test*.cpp"

skills:
  unit_test_skill:
    version: "2"
    inject_mode: prompt_append
    instruction_path: ./skills/unit_test_skill/instructions.md
    compatible_frameworks:
      - opencode
    compatible_languages:
      - python
      - go
      - java
      - cpp
```

旧写法（`sandbox_mode` + `docker_images`）仍然兼容，但建议迁移到 `sandbox.*` 块。

### framework 字段

| 字段 | 说明 |
|------|------|
| `kind` | `model_api`（纯 API）或 `cli_agent`（Agent 框架） |
| `command` | 命令模板（仅 `cli_agent`） |
| `sandbox.provider` | `local` 或 `docker`（后续支持 `e2b`、`firecracker`） |
| `sandbox.mode` | `local` 或 `docker` |
| `sandbox.images` | 按语言选择的 Agent 沙箱镜像 |
| `sandbox.timeout_seconds` | 单样本 Agent 超时 |
| `sandbox.cpu` / `sandbox.memory` | 资源限制 |
| `sandbox.network_disabled` | 是否禁网 |
| `preflight` | 按语言定义的执行前环境检查命令 |
| `forbidden_command_patterns` | 识别环境漂移的命令模式，如 `apt-get install` |
| `output_globs` | 测试文件发现规则 |
| `env` | 传给 Agent 的环境变量 |
| `env_from_host` | 从宿主环境透传进 Agent 进程或容器的变量名 |
| `compatible_models` | 允许的模型列表 |
| `compatible_languages` | 允许的语言列表 |

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

### inject_mode 说明

| 模式 | 行为 |
|------|------|
| `prompt_append` | 将 `instruction_path` 内容追加到 prompt 末尾，辅助文件复制到 `.utbench/skills/` |
| `workspace_mount` | 仅将 `files` 中的文件复制到 `.utbench/skills/`，不修改 prompt |
| `agent_native` | 转换为 Agent 框架原生 skill 格式写入 workspace：CodeBuddy → `.codebuddy/skills/<name>/SKILL.md`，OpenCode → `.opencode/skills/<name>.md`。Agent 启动后自动加载原生 skill，无需手动注入 prompt |

## 7. CLI Agent 的调用约定

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

### 评测环境约束

UT-Bench 现在会在每个 `subject × sample` 的 Agent 执行前做 preflight。典型检查包括：

- Python: `python3 --version`, `pytest --version`
- Go: `go version`
- Java: `java -version`, `mvn -version`
- C++: `g++ --version`, `cmake --version`

如果 preflight 失败，任务会直接以 `sandbox_preflight_error` 终止，而不是让 Agent 在运行中自己尝试 `apt-get install`。

UT-Bench 也会检查环境漂移命令。默认会拦截这类操作：

- `apt-get install`
- `apt install`
- `apk add`
- `yum install`
- `pip install`
- `npm install`

这条规则的目的很明确：评测平台负责提供完整、固定、可复现的环境；Agent 只负责在这个环境里生成测试。

### 样本依赖准备

除了语言基础镜像，UT-Bench 现在还会在 Agent 执行前按样本准备常见依赖：

- Python `repo_level` 元数据里的 `requirements`
- workspace 根目录的 `requirements.txt` / `requirements-dev.txt`
- `go.mod` -> `go mod download`
- `pom.xml` -> `mvn -q -DskipTests dependency:go-offline`

这一步由平台执行，并记录到 trace 里的 `environment_setup`。如果依赖准备失败，任务会以 `sample_env_prepare_error` 结束，不会把“缺依赖”误算成 Agent 生成能力问题。

## 8. Docker 与沙箱的边界

UT-Bench 的 Docker 体系分为三种镜像（详见 [DOCKER_GUIDE.md](./DOCKER_GUIDE.md)）：

| 镜像 | 用途 |
|:---|:---|
| `utbench:latest` | **评测镜像**：compile / test / coverage / mutation |
| `utbench-control:latest` | **控制面镜像**：Web / CLI / 编排（不含评测工具链） |
| `utbench-agent-{fw}-{lang}:latest` | **Agent 沙箱镜像**：Agent CLI + 语言运行时 |

### 外层 Docker（评测镜像）

```text
docker run utbench:latest ...
```

这是 UT-Bench 的评测执行环境，包含全部语言工具链。

### 控制面 Docker

```text
docker run utbench-control:latest web ...
```

只包含 utbench 二进制 + docker CLI，用于 Web UI 编排场景。

### 内层 Agent 沙箱

当 framework 配置为 `sandbox.provider: docker` 时，UT-Bench 会对每个 `subject × sample`：

1. 创建独立工作区
2. 启动独立容器（使用 `sandbox.images.{lang}` 指定的镜像）
3. 只挂载该工作区到 `/workspace`
4. 可选禁网
5. 设置 CPU / memory / timeout
6. 回收容器

隔离粒度：一个 Agent 执行一个样本，对应一个独立工作区，对应一个独立容器。这既是安全边界，也是公平边界。

### 当前状态

当前代码已支持：

- 工作区按 `subject × sample` 拆分
- `sandbox.provider: docker` 时按样本起内层容器
- 记录 trace、diff、sandbox fingerprint
- 三种镜像独立版本化和构建

后续演进方向：

- 更强的隔离层（gVisor / Firecracker / E2B）
- Evaluator backend 独立抽象
- 远程 sandbox provider

## 9. 报告新增内容

除了原有的编译、测试、覆盖率、变异得分，当前报告还会额外产出：

- `agent_comparisons`
  - 同模型、同样本下，Agent 相对 `model_api` baseline 的提升
- `skill_uplifts`
  - 同 framework + model、同样本下，skill 相对 `no_skill` 的提升

常见比较：

- `opencode__deepseek-v4-flash__no_skill` vs `model_api__deepseek-v4-flash__no_skill`
- `opencode__deepseek-v4-flash__unit_test_skill` vs `opencode__deepseek-v4-flash__no_skill`

## 10. 数据集与样本准备

### `self_contained`

直接使用单文件样本。对纯模型 API 和 Agent 都适用。

### `repo_level`

如果样本旁边存在 repo-level metadata，runner 会复制整个 workspace，再把目标文件交给 Agent。

这意味着 Agent 视角里拿到的是一个最小可操作 repo，而不是孤立源码片段。

## 11. 现在优先接哪个真实 CLI Agent

建议先接 `OpenCode`。

理由不是品牌偏好，而是适配目标：

- 你要的是 `framework × model` 的横评能力
- `OpenCode` 更适合作为通用 Agent 框架入口
- `Claude Code` 更适合作为 Claude 自身工作流的专项评测

所以建议顺序：

1. `OpenCode`
2. `Claude Code`

## 12. 常见问题

### `--models` 和 `--subjects` 要不要同时传

建议传。

原因是 subject 展开和模型配置加载都需要模型集合。第一版里，`agents.yaml` 也会对模型名做交集过滤。

### 只传 `--models` 能不能跑

可以。那就是纯模型 baseline 模式。

### 只传 `--agents-config` 不传 `--subjects` 能不能跑

可以。系统会自动展开全部合法组合。

### `sandbox_mode: local` 有什么用

主要用于本地调试和测试。正式跑 Agent 横评时，应该优先用 `sandbox_mode: docker`。

