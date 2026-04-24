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

### 3. 运行评测

**Linux/macOS:**

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek \
    --langs python,java \
    --config /app/configs/models.yaml \
    --dataset-root /app/datasets \
    --max-samples 5
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
    --langs python,java `
    --config /app/configs/models.yaml `
    --dataset-root /app/datasets `
    --max-samples 5
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
    --class module_level
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

- **Python/Go**: 使用 `--class module_level`
- **Java/C++**: 使用 `--class self_contained`

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