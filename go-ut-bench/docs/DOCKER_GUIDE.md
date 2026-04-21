# Docker 使用指南

## 快速开始

### 1. 构建 Docker 镜像

```bash
# Linux/macOS
./docker.sh build

# Windows PowerShell
docker build -t utbench:latest .
```

### 2. 配置 API 密钥

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑 .env 文件，填入你的 API 密钥
```

### 3. 运行评测

```bash
# 快速测试（Python，dry-run）
./docker.sh run-python

# 评测所有语言
./docker.sh run-all

# 自定义参数
./docker.sh run run --models deepseek --langs python --max-samples 10
```

### 4. 进入容器调试

```bash
./docker.sh shell
```

## Docker Compose 方式

```bash
# 构建并运行
docker-compose up --build

# 运行特定服务
docker-compose run run-python

# 查看日志
docker-compose logs -f
```

## 挂载目录说明

| 容器路径 | 说明 |
|---------|------|
| `/app/datasets` | 数据集目录 |
| `/app/artifacts` | 输出结果目录 |
| `/app/storage` | SQLite 数据库 |
| `/app/configs` | 配置文件 |

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
- JaCoCo、pitest 插件（Maven 会自动下载）

**C++:**
- CMake 3.10+
- GoogleTest
- gcov
- mull（可选，用于 mutation testing）

### 运行命令

```bash
# 构建
go build -o utbench ./cmd/utbench/

# 运行
./utbench run --models deepseek --langs python --max-samples 5
```

## 跨平台修复说明

代码已修复以下 Windows 兼容性问题：

1. **go_eval.go**: `/dev/null` → 使用临时文件
2. **cpp_eval.go**: `os.Symlink` → 使用文件复制
3. **java_eval.go**: 正则匹配增加 nil 检查
4. **service.go**: 移除死代码，修复断言密度计算

## 常见问题

### Q: Windows 上 Docker 很慢？
A: 可以使用 WSL2 + Docker Desktop，或直接本地运行（需安装工具链）。

### Q: 某些语言评测失败？
A: 检查对应工具链是否正确安装：
```bash
# Python
pytest --version
coverage --version

# Go
go version
go-mutesting --version

# Java
java -version
mvn -version

# C++
cmake --version
gcov --version
```

### Q: 如何只评测特定语言？
A: 使用 `--langs` 参数：
```bash
./utbench run --langs python,go
```