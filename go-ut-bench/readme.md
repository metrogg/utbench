# go-ut-bench

多语言单元测试生成效果横向评测 CLI 工具。

## 支持语言

| 语言 | 测试框架 | 覆盖率工具 | 变异测试 |
|------|---------|-----------|---------|
| Python | pytest | coverage | mutmut |
| Go | go test | go test -cover | go-mutesting |
| Java | Maven/JUnit | JaCoCo | pitest |
| C++ | GoogleTest | gcov | mull |

## 支持模型

| 模型 | Provider | 说明 |
|------|----------|------|
| deepseek | deepseek | DeepSeek Chat |
| qwen | dashscope | 通义千问 3.6-plus |
| minimax | minimax | MiniMax M2.7 |
| doubao-seed | volcengine | 豆包 Seed 2.0 Pro |
| doubao-seed-2.0-lite | volcengine | 豆包 Seed 2.0 Lite |
| doubao-seed-1.6 | volcengine | 豆包 Seed 1.6 |
| doubao-seed-2.0-pro-v2 | volcengine | 豆包 Seed 2.0 Pro V2 |

## 快速开始

### Docker 运行（推荐）

```bash
# 构建镜像
docker build -t utbench:latest .

# 配置 API 密钥
cp .env.example .env

# 运行测试
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

### 本地运行

```bash
go build -o utbench ./cmd/utbench/
export DEEPSEEK_API_KEY="sk-xxx"
./utbench run --models deepseek --langs python --max-samples 5
```

## 命令

| 命令 | 说明 |
|------|------|
| `utbench run` | 完整流程 (generate → evaluate → report) |
| `utbench generate` | 仅生成单元测试 |
| `utbench evaluate` | 评测已生成的单元测试 |
| `utbench report` | 生成评测报告 |
| `utbench ingest` | 结果导入 SQLite |
| `utbench dataset` | 数据集管理 |

## 输出结构

```
artifacts/runs/<run-id>/
  generated/           # 生成的测试文件
  evaluation/          # 评测结果 JSON
  report/              # HTML 报告
  run.log              # 运行日志
  api.log              # API 调用日志
```

## 特性

- **截断自动续写**：检测到输出截断时自动发送续写请求
- **截断统计分析**：报告中显示截断率、续写统计、调优建议
- **增量运行**：支持 checkpoint 断点续跑
- **变异测试**：可选启用变异测试评估测试质量

## 文档

- [用户指南](docs/USER_GUIDE.md) - 完整使用文档
- [CLI 参数](docs/cli-spec.md) - 命令行参数详解
- [Docker 使用](docs/DOCKER_GUIDE.md) - Docker 运行指南
- [架构设计](docs/architecture-mvp.md) - 系统架构

## 环境要求

### Python
```bash
pip install pytest coverage mutmut
```

### Go
```bash
go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest
```

### Java
- JDK 17+
- Maven 3+

### C++
```bash
sudo apt-get install cmake clang-15 libgtest-dev g++-15 mull-15
```