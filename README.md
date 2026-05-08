# UT-Bench

LLM 单测生成能力横向评测基准工具。

## 快速开始

### 第一步：构建 Docker 镜像（推荐）

镜像预装了所有评测工具链（pytest/coverage/mutmut、go-mutesting、JDK/Maven/pitest、clang/gcov/mull）：

```bash
cd go-ut-bench
docker build -t utbench:latest .
```

### 第二步：配置 API 密钥

创建 `go-ut-bench/.env` 文件：

```dotenv
DEEPSEEK_API_KEY=sk-xxx
DASHSCOPE_API_KEY=sk-xxx      # 阿里云百炼 (qwen)
MINIMAX_API_KEY=xxx
ARK_API_KEY=xxx               # 火山引擎 (doubao-seed / glm-4.7)
```

### 第三步：选择使用方式

#### 方式 A：Web 管理界面（可视化操作）

```bash
cd go-ut-bench
go run ./cmd/utbench web
```

打开浏览器访问 `http://localhost:8080`：

- **总览**：任务统计、Docker/镜像状态
- **新建任务**：交互式表单，勾选模型、语言、场景等参数
- **任务列表**：历史任务、状态过滤
- **任务详情**：实时日志、评测报告、横向对比图表
- **数据管理**：跨运行对比、历史数据入库

#### 方式 B：命令行（CLI）

```bash
# Linux/macOS
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest run \
    --models deepseek,qwen \
    --langs python,go \
    --config /app/configs/models.yaml \
    --max-samples 10

# Windows PowerShell
docker run --rm --env-file .env `
  -v "${PWD}/datasets:/app/datasets" `
  -v "${PWD}/artifacts:/app/artifacts" `
  -v "${PWD}/configs:/app/configs" `
  utbench:latest run `
    --models deepseek `
    --langs python `
    --config /app/configs/models.yaml `
    --max-samples 5
```

---

## 详细文档

| 文档 | 说明 |
|------|------|
| [go-ut-bench 文档导航](go-ut-bench/docs/README.md) | 当前 Go CLI 的主文档入口 |
| [用户指南](go-ut-bench/docs/01-user-guides/USER_GUIDE.md) | 快速开始、运行模式、Agent/Skill 和常见问题 |
| [完整启动指南](go-ut-bench/docs/01-user-guides/startup-guide.md) | 参数详解和常见场景 |
| [Docker 使用指南](go-ut-bench/docs/03-operations/DOCKER_GUIDE.md) | Docker 构建、运行和沙箱边界 |
| [Web UI 使用说明](go-ut-bench/docs/01-user-guides/WEB_UI.md) | Web 管理界面使用说明 |
| [文档分类与去重索引](go-ut-bench/docs/DOCUMENT_CLASSIFICATION.md) | 文档保留等级、重合关系和归档建议 |

---

## 功能概览

### 评测流水线

```
数据集 → 生成测试 (LLM) → 编译 → 运行测试 → 覆盖率 → 变异测试 → 报告
```

### 支持的语言和评测工具

| 语言 | 测试框架 | 覆盖率 | 变异测试 |
|------|----------|--------|----------|
| Python | pytest | coverage | mutmut |
| Go | go test | go tool cover | go-mutesting |
| Java | JUnit 5 | JaCoCo | pitest |
| C++ | GoogleTest | gcov | mull |

### 支持的模型

| 模型 | 提供商 |
|------|--------|
| deepseek | DeepSeek |
| qwen | 阿里云百炼 |
| minimax | MiniMax |
| doubao-seed / doubao-seed-2.0-lite / doubao-seed-1.6 / doubao-seed-2.0-pro-v2 | 火山引擎 |
| glm-4.7 | 火山引擎 ARK |

### 数据集场景

- `boundary` — 边界值测试
- `simple_function` — 简单函数
- `complex_dependency` — 复杂依赖
- `interface_mock` — 接口 mock

---

## 目录结构

```
go-ut-bench/
  cmd/utbench/           # CLI 入口
  internal/
    orchestrator/        # 流水线编排
    runner/              # LLM API 调用、断点续跑
    evaluator/           # 编译/测试/覆盖率/变异
    reporter/            # HTML 报告生成
    store/               # SQLite 持久化
    web/                 # Web 管理界面
  configs/models.yaml    # 模型配置
  datasets/              # 数据集目录
  artifacts/             # 运行产物（报告、日志）
  storage/               # SQLite 数据库
```

---

## 常用场景速查

### 验证环境（空跑）

```bash
docker run --rm -v "$(pwd)/datasets:/app/datasets" -v "$(pwd)/artifacts:/app/artifacts" -v "$(pwd)/configs:/app/configs" \
  utbench:latest run --models deepseek --langs python --config /app/configs/models.yaml --max-samples 1 --dry-run
```

### 多模型横向对比

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" -v "$(pwd)/artifacts:/app/artifacts" -v "$(pwd)/configs:/app/configs" -v "$(pwd)/storage:/app/storage" \
  utbench:latest run \
    --run-id compare_001 \
    --models deepseek,qwen,minimax,glm-4.7 \
    --langs python,go \
    --config /app/configs/models.yaml \
    --max-samples 20 \
    --ingest --db-path /app/storage/utbench.db
```

### 带变异测试的完整评测

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" -v "$(pwd)/artifacts:/app/artifacts" -v "$(pwd)/configs:/app/configs" \
  utbench:latest run \
    --models deepseek --langs python \
    --config /app/configs/models.yaml \
    --max-samples 5 \
    --mutation-enabled
```

---

## 报告查看

每次运行产物位于 `artifacts/runs/<run-id>/`：

```
artifacts/runs/<run-id>/
  generated/               # 生成的测试文件
  evaluation/              # 评测结果 JSON
  report/
    report_summary.json    # 汇总 JSON
    report.html            # 可视化 HTML 报告 ⭐
```

用浏览器直接打开 `report.html` 查看模型横向对比图表、覆盖率分布、变异得分等。
