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
| `utbench doctor` | 检查评测工具链并运行 canary 自检 |

## 输出结构

```
artifacts/runs/<run-id>/
  generated/           # 生成的测试文件
  evaluation/          # 评测结果 JSON
  report/              # HTML 报告
  logs/                # run/evaluator/api/errors 结构化日志
```

## 特性

- **截断自动续写**：检测到输出截断时自动发送续写请求
- **截断统计分析**：报告中显示截断率、续写统计、调优建议
- **增量运行**：支持 checkpoint 断点续跑
- **变异测试**：可选启用变异测试评估测试质量
- **评测自检**：`utbench doctor` 检查工具版本并运行临时 canary 样本
- **数据集审计**：`utbench dataset validate` 统计样本并标记外部 I/O、非确定性和复杂度风险

## 数据集说明

当前仓库内置数据集实际为 `self_contained`：Python、Go、Java、C++ 各 4 个场景，每个场景 50 个样本。正式运行请显式使用 `--class self_contained`，避免旧文档中的 module-level 说明造成样本集合不一致。

## 报告口径

- `compile_pass_rate`、`sample_test_pass_rate`、`avg_test_pass_rate` 都是样本级口径。
- `test_case_pass_rate`、`avg_test_case_pass_rate` 是测试用例级口径，用来补充说明单个样本内部测试函数通过情况。
- 排名和综合分默认使用样本级测试通过率，避免样本内测试函数数量差异放大分数。
- `avg_latency_ms` 现在表示单样本完整评测耗时，不再是某个子阶段的局部时间。
- 所有语言都必须样本级测试通过后才运行变异测试；基线测试失败的样本变异分记 0，并归入模型问题。
- 变异测试失败会按原因细分展示：`mutation_skipped_baseline_failed`、`mutation_target_not_exercised`、`mutation_no_results`、`mutation_no_coverage`、`mutation_no_effective_mutants`、`mutation_timeout`、`mutation_tool_error`、`mutation_error`。
- 工具/环境/数据集问题不进入模型排名分母；模型生成代码导致的编译或测试失败仍进入排名。

```bash
./utbench doctor --langs python,go,java,cpp --mutation-enabled --mutation-timeout 120
./utbench dataset validate --dataset-root ./datasets --langs python,go,java,cpp --class self_contained --strict
```

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
