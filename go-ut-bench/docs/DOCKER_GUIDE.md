# Docker 使用指南

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

---

## 挂载目录说明

| 容器路径 | 说明 |
|---------|------|
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
- `module_level` 属于预留/旧数据说明；除非你明确恢复 module-level 数据集，否则不要用于正式横评

### 路径问题

在容器内运行时，所有路径使用容器路径：
- 配置文件：`/app/configs/models.yaml`
- 数据集：`/app/datasets`
- 输出：`/app/artifacts`

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
