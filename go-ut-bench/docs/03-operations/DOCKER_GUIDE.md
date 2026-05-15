# Docker 使用指南

UT-Bench 的 Docker 体系分为两种镜像：

| 镜像 | Dockerfile | 用途 |
|:---|:---|:---|
| `utbench:latest` | `Dockerfile` | **评测镜像**：包含全部语言工具链（Python/Go/Java/C++）、coverage、mutation 工具。用于 `run`/`evaluate`/`doctor` 等需要编译和测试的命令 |
| `utbench-agent-opencode:latest` | `docker/agents/opencode/Dockerfile` | **Agent 沙箱镜像**：统一镜像，包含 Python/Go/Java/C++ 运行时 + OpenCode CLI。用于 `cli_agent` 类型的 subject 执行 |

控制面（Web/CLI）直接在宿主机运行，不需要单独的 Docker 镜像。

```
宿主机
├── ./utbench web / ./utbench run    ← 控制面：Web / CLI / 编排 / SQLite / 报告
│   └── generation 在宿主机执行      ← Agent 沙箱容器由宿主机直接启动
│
├── utbench:latest 容器              ← 评测面：compile / test / coverage / mutation
│   └── 不需要 docker.sock
│
└── utbench-agent-opencode 容器      ← Agent 沙箱：OpenCode 等 Agent 的隔离执行环境
```

Web 模式下流水线自动拆分：generation 和 report 在宿主机 in-process 执行（Agent 沙箱可用宿主机 Docker），
evaluation 在 eval 容器内执行（语言工具链齐全）。eval 容器不需要 docker.sock。

> 如果你想深入了解为什么这样拆分，以及后续演进方向（gVisor / E2B / Firecracker），看这里：
> [runtime-sandbox-redesign.md](./runtime-sandbox-redesign.md)

---

## 快速开始

### 1. 构建镜像

使用统一构建脚本（推荐）：

```bash
# Linux / WSL
./build.sh all          # 构建全部镜像
./build.sh all --cn     # 国内网络加速

# Windows PowerShell
.\build.ps1 all         # 构建全部镜像
.\build.ps1 all -Cn     # 国内网络加速
```

按需构建：

```bash
./build.sh eval             # 仅评测镜像
./build.sh sandbox          # 统一 Agent 沙箱（全语言）
./build.sh verify           # 验证镜像状态
```

手动构建（不使用脚本）：

```bash
# 评测镜像（跑评测必须）
docker build -t utbench:latest .

# Agent 沙箱镜像（跑 Agent subject 必须）
docker build -t utbench-agent-opencode:latest -f docker/agents/opencode/Dockerfile docker/agents/opencode/
```

注意：

- `FROM node:...`、`FROM python:...` 等基础镜像是否走国内镜像，由你的 Docker daemon 配置决定
- 当前只替换 Debian 主仓 `deb.debian.org/debian`，默认保留 `security.debian.org`
- 构建脚本的 `--cn` / `-Cn` 参数会自动配置 Ubuntu/Debian 镜像源、Go 下载源、npm/pip 镜像源

### 2. 配置 API 密钥

```bash
cp .env.example .env
# 编辑 .env 文件，填入你的 API 密钥
```

`.env` 文件内容：

```bash
DEEPSEEK_API_KEY=your-deepseek-api-key        # DeepSeek
DASHSCOPE_API_KEY=xxx          # 通义千问 (Qwen)
MINIMAX_API_KEY=xxx            # Minimax
VOLCENGINE_API_KEY=xxx         # 豆包 (Doubao)
ARK_API_KEY=xxx                # 豆包 ARK 版本
```

---

## 使用场景

### 场景 A：宿主机运行（推荐）

Web / CLI 直接在宿主机运行，评测通过 Docker 容器执行。

```bash
# 1. 构建镜像
./build.sh all --cn

# 2. 配置 API key
cp .env.example .env

# 3. 启动 Web UI
./utbench web --addr :8080 --config ./configs/models.yaml

# 或 CLI 运行
./utbench run --models deepseek --langs python --max-samples 2
```

### 场景 B：Docker 运行评测

评测在 Docker 容器内执行，宿主机直接调用。

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
    --max-samples 2
```

---

## 架构说明

UT-Bench 采用三层平面架构：

- **控制面**（control plane）：CLI / Web / SQLite / 报告
- **评测执行面**（evaluation plane）：compile / test / coverage / mutation
- **Agent 沙箱执行面**（agent sandbox plane）：Agent 隔离执行、workspace 读写、trace 采集

Web 模式下，Agent 沙箱在 eval 容器内以 local 模式运行，不依赖 docker.sock。

源码保护机制：
- Docker 模式：源文件通过 `:ro` 只读挂载到容器内
- Local 模式：源文件 `chmod 0444` 设为只读
- 两种模式均提供纵深防御

---

## 挂载目录说明

| 容器路径 | 说明 |
|:---|:---|
| `/app/datasets` | 数据集目录 |
| `/app/artifacts` | 输出结果目录 |
| `/app/configs` | 配置文件目录 |
| `/app/storage` | SQLite 数据库（可选） |

---

## 常用命令示例

### 评测环境自检

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

### 数据集校验

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

### 跑纯模型 baseline

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

### 跑 Agent subject（需要 Agent 沙箱镜像）

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek-v4-flash \
    --langs python \
    --config /app/configs/models.yaml \
    --agents-config /app/configs/agents.example.yaml \
    --subjects model_api__deepseek-v4-flash__no_skill,opencode__deepseek-v4-flash__no_skill \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --class self_contained \
    --max-samples 1
```

### Dry-run（不调用 API）

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

### 进入容器调试

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

## 注意事项

### 数据集类别

当前仓库内置数据集为 `self_contained`，使用 `--class self_contained`。`repo_level` 属于预留，除非你明确恢复 repo-level 数据集，否则不要用于正式横评。

### 路径问题

容器内运行时使用容器路径：配置 `/app/configs/models.yaml`，数据集 `/app/datasets`，输出 `/app/artifacts`。

### 源码保护

Docker 模式下，源文件通过 `:ro` 只读挂载到容器内，防止 Agent 意外修改被测源码。Local 模式下使用 `chmod 0444`。后续计划支持更强隔离（gVisor、E2B、Firecracker）。

### Windows

PowerShell 使用 `${PWD}`，Bash 使用 `\` 续行。WSL/Linux Bash 不要使用 PowerShell 的反引号续行。

---

## 本地运行（不使用 Docker）

### 前置要求

**Python:**
```bash
pip install pytest coverage mutmut
```

**Go:**
```bash
go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
```

**Java:**
- JDK 17+，Maven 3+

**C++:**
```bash
sudo apt-get install cmake clang-19 libgtest-dev g++ mull-19
```

### 运行命令

```bash
go build -o utbench ./cmd/utbench/
export DEEPSEEK_API_KEY="your-deepseek-api-key"
./utbench run --models deepseek --langs python --max-samples 5
```

---

## 查看帮助

```bash
docker run --rm utbench:latest --help
docker run --rm utbench:latest run --help
```
