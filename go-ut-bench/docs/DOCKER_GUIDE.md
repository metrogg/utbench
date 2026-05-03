# Docker 使用指南

这份文档分两类场景：

- 外层容器只运行 UT-Bench 本身
- 外层容器还要为每个 `subject × sample` 再起一个 Agent 沙箱容器

第二种场景就是当前 Agent 评测要走的路径。因为 UT-Bench 通常本身就跑在 Docker 里，这里采用的是过渡性的 **DOOD** 方式：外层容器内的 `docker` CLI 去控制宿主机 Docker daemon。

## DOOD 结论

如果某个 framework 配置了：

```yaml
sandbox_mode: docker
```

那你运行外层 `utbench:latest` 时，必须至少满足两点：

1. 外层镜像里有 `docker` CLI
2. 外层容器能访问宿主机 Docker daemon

当前仓库的 `Dockerfile` 已经安装了 `docker.io`，所以还差第二点：挂 Docker socket。

不挂 Docker socket 的结果很简单：

- 纯模型 baseline 可以跑
- `sandbox_mode: local` 的 CLI Agent 可以跑
- `sandbox_mode: docker` 的 CLI Agent 不能跑

## 快速开始

### 1. 构建 Docker 镜像

```bash
# Linux/macOS
docker build -t utbench:latest .

# Windows PowerShell
docker build -t utbench:latest .
```

### 2. 配置 API 密钥

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，填入你的 API 密钥
```

`.env` 文件内容：

```bash
DEEPSEEK_API_KEY=sk-xxx        # DeepSeek
DASHSCOPE_API_KEY=xxx          # 通义千问 (Qwen)
MINIMAX_API_KEY=xxx            # Minimax
VOLCENGINE_API_KEY=xxx         # 豆包 (Doubao)
ARK_API_KEY=xxx                # 豆包 ARK 版本
```

### 3. 先做评测环境自检

正式评测前先运行 `doctor`，确认四种语言的 compile/test/coverage/mutation 工具链都能正常工作。

**Linux/macOS/WSL Bash 使用反斜杠 `\` 续行：**

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest doctor \
    --langs python,go,java,cpp \
    --mutation-enabled \
    --mutation-timeout 120
```

**Windows PowerShell 使用反引号 `` ` `` 续行：**

```powershell
docker run --rm --env-file .env `
  -v "${PWD}/datasets:/app/datasets" `
  -v "${PWD}/artifacts:/app/artifacts" `
  -v "${PWD}/configs:/app/configs" `
  utbench:latest doctor `
    --langs python,go,java,cpp `
    --mutation-enabled `
    --mutation-timeout 120
```

期望输出包含：

```text
Doctor: OK
[canary] cpp compile=true test=true coverage=true mutation=true
[canary] go compile=true test=true coverage=true mutation=true
[canary] java compile=true test=true coverage=true mutation=true
[canary] python compile=true test=true coverage=true mutation=true
```

### 4. 再做数据集校验

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest dataset validate \
    --dataset-root /app/datasets \
    --langs python,go,java,cpp \
    --class self_contained \
    --strict
```

如果需要保存 JSON 报告：

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest dataset validate \
    --dataset-root /app/datasets \
    --langs python,go,java,cpp \
    --class self_contained \
    --strict \
    --json /app/artifacts/dataset_validate_report.json
```

期望 `Dataset validation: OK` 且 `errors: 0`。warnings 是风险提示，默认不阻断正式评测。

### 5. 运行评测

#### 5a. 只跑模型 API 或本地 CLI Agent

如果你没有用 `sandbox_mode: docker` 的 framework，普通运行方式就够了。

评测阶段如果长时间停在最后几个样本，优先看实时日志里的 `[EVAL-WARN]` 行。它会输出仍在运行的 `model/language/sample_id/phase`，用于区分是 compile/test/coverage/mutation 外部工具超时，还是临时工作区清理耗时。当前 evaluator 对 Go、Java、C++、Python 的外部命令都设置了超时和进程组清理，临时目录清理会在后台进行，不再阻塞结果落盘。

**Linux/macOS:**

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek \
    --langs python,go,java,cpp \
    --config /app/configs/models.yaml \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --class self_contained \
    --max-samples 2 \
    --mutation-enabled \
    --mutation-timeout 360
```

**Windows PowerShell:**

```powershell
docker run --rm --env-file .env `
  -v "${PWD}/datasets:/app/datasets" `
  -v "${PWD}/artifacts:/app/artifacts" `
  -v "${PWD}/configs:/app/configs" `
  -w /app `
  utbench:latest run `
    --models deepseek `
    --langs python,go,java,cpp `
    --config /app/configs/models.yaml `
    --dataset-root /app/datasets `
    --output-root /app/artifacts `
    --class self_contained `
    --max-samples 2 `
    --mutation-enabled `
    --mutation-timeout 360
```

#### 5b. 跑每样本 Docker Agent 沙箱

如果你要跑 `OpenCode` 这种 `sandbox_mode: docker` 的 framework，外层容器必须挂 Docker socket。

先构建外层 benchmark 镜像和按语言划分的内层 Agent 镜像：

```bash
docker build -t utbench:latest .
docker build -t utbench-agent-opencode-python:latest -f ./docker/agents/opencode/python.Dockerfile .
docker build -t utbench-agent-opencode-go:latest -f ./docker/agents/opencode/go.Dockerfile .
docker build -t utbench-agent-opencode-java:latest -f ./docker/agents/opencode/java.Dockerfile .
docker build -t utbench-agent-opencode-cpp:latest -f ./docker/agents/opencode/cpp.Dockerfile .
```

如果你希望优先使用国内镜像源，可以给内层 Agent 镜像显式传构建参数：

```bash
docker build \
  --build-arg DEBIAN_MIRROR=http://mirrors.ustc.edu.cn/debian \
  --build-arg NPM_REGISTRY=https://registry.npmmirror.com \
  --build-arg PIP_INDEX_URL=https://pypi.tuna.tsinghua.edu.cn/simple \
  -t utbench-agent-opencode-python:latest \
  -f ./docker/agents/opencode/python.Dockerfile .
```

Go / Java / C++ 镜像也支持 `DEBIAN_MIRROR` 和 `NPM_REGISTRY`。  
注意：

- 这些参数只影响 Dockerfile 构建阶段里的 apt / npm / pip 下载；`FROM node:...`、`FROM python:...` 这类基础镜像本身是否走国内镜像，仍由你的 Docker daemon 配置决定。
- 当前实现只替换 Debian 主仓 `deb.debian.org/debian`，默认保留官方 `security.debian.org`。这是刻意的：很多国内镜像在 `debian-security` 上同步延迟更明显，容易触发 `File has unexpected size`。

Linux/macOS/WSL：

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
    --max-samples 1
```

这个运行方式的含义是：

- 外层 `utbench:latest` 容器只负责 orchestration
- 每个 `subject × sample` 的 Agent 执行由外层容器内的 `docker` CLI 再起一个内层容器
- 内层容器只挂当前样本工作区到 `/workspace`

`UTBENCH_SANDBOX_HOST_OUTPUT_ROOT` 的作用是把外层容器里的 `/app/artifacts/...` 映射回宿主机真实路径。没有它时，宿主机 Docker daemon 会把 `/app/artifacts/...` 当成宿主机路径解析，结果就是内层容器拿到空目录。

也就是说，真正的 Agent 样本级沙箱是内层容器，不是外层 benchmark 容器

### 5c. 为什么要按语言拆内层镜像

现在的原则是：

- benchmark 平台负责提供完整工具链
- Agent 不能在运行中自行安装依赖

所以内层镜像会按语言拆分：

- `utbench-agent-opencode-python`
- `utbench-agent-opencode-go`
- `utbench-agent-opencode-java`
- `utbench-agent-opencode-cpp`

每个镜像都预装该语言做单元测试最基本的一套运行时和工具。UT-Bench 在执行前会运行 preflight，自检这些工具是否可用；如果缺工具，任务直接报 `sandbox_preflight_error`。

如果 Agent 在执行过程中尝试：

- `apt-get install`
- `pip install`
- `npm install`

UT-Bench 会把它记成 `sandbox_policy_error`。这是刻意的：环境完整性是平台责任，不应该交给被测 Agent。

---

## 挂载目录说明

| 容器路径 | 说明 |
|---------|------|
| `/var/run/docker.sock` | 宿主机 Docker daemon，只有 `sandbox_mode: docker` 时需要 |
| `/app/datasets` | 数据集目录 |
| `/app/artifacts` | 输出结果目录 |
| `/app/configs` | 配置文件目录 |
| `/app/storage` | SQLite 数据库（可选） |

---

## 常用命令示例

### 快速测试（2个样本）

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek \
    --langs python \
    --config /app/configs/models.yaml \
    --dataset-root /app/datasets \
    --max-samples 2 \
    --class self_contained
```

### 全量评测（多模型+多语言）

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek,qwen,minimax \
    --langs python,java,go,cpp \
    --config /app/configs/models.yaml \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --class self_contained \
    --mutation-enabled \
    --mutation-timeout 360
```

### 仅生成测试（调试用）

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest generate \
    --models deepseek \
    --langs cpp \
    --config /app/configs/models.yaml \
    --dataset-root /app/datasets \
    --max-samples 1 \
    --class self_contained
```

### 跑 OpenCode Agent 基线 + Skill

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
    --max-samples 1
```

### Dry-run（不调用API）

```bash
docker run --rm \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek \
    --langs python \
    --config /app/configs/models.yaml \
    --dataset-root /app/datasets \
    --max-samples 2 \
    --dry-run
```

---

## 进入容器调试

```bash
docker run --rm -it \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest \
  /bin/bash
```

---

## 查看帮助

```bash
docker run --rm utbench:latest --help
docker run --rm utbench:latest run --help
```

---

## 注意事项

### 数据集类别

当前仓库内置数据集实际为 `self_contained`：

- Python/Go/Java/C++ 都使用 `--class self_contained`
- `repo_level` 属于预留/旧数据说明；除非你明确恢复 repo-level 数据集，否则不要用于正式横评

### 路径问题

在容器内运行时，所有路径使用容器路径：
- 配置文件：`/app/configs/models.yaml`
- 数据集：`/app/datasets`
- 输出：`/app/artifacts`

### DOOD 风险

挂 Docker socket 意味着外层容器拥有较强的宿主机 Docker 控制能力。当前这是为了尽快把“每样本一个 Agent 容器”的链路跑通，不建议把它当最终安全方案。

当前建议把 DOOD 视为过渡方案：

1. 先把 Agent benchmark 跑通
2. 再把样本沙箱执行层抽成独立 backend
3. 后续再考虑 remote Docker、gVisor、Firecracker、E2B 一类替代方案

### Windows 路径

Windows PowerShell 使用 `${PWD}` 获取当前目录：
```powershell
-v "${PWD}/datasets:/app/datasets"
```

WSL/Linux Bash 不要使用 PowerShell 的反引号续行；Bash 里应使用 `\`。否则会出现 `-v: command not found` 或 `--langs: command not found`。

---

## 本地运行（不使用 Docker）

### 前置要求

安装对应语言的工具链：

**Python:**
```bash
pip install pytest coverage mutmut
```

**Go:**
```bash
go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
```

**Java:**
- JDK 17+
- Maven 3+

**C++:**
```bash
sudo apt-get install cmake clang-15 libgtest-dev g++-15 mull-15
```

### 运行命令

```bash
# 构建
go build -o utbench ./cmd/utbench/

# 设置环境变量
export DEEPSEEK_API_KEY="sk-xxx"

# 运行
./utbench run --models deepseek --langs python --max-samples 5
```

