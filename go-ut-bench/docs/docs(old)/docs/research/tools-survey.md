# 相关工具调研

## 1. 调研目标

为完整评测流水线选择工具链，覆盖：

- 生成阶段（Runner）
- 编译/执行阶段（Evaluator）
- 覆盖率与变异测试
- 报告聚合与可视化

## 2. 当前已落地工具

### 2.1 Runner 层

- 语言：Python
- 配置解析：`pyyaml`
- HTTP 调用：`urllib.request`
- 并发：`concurrent.futures.ThreadPoolExecutor`
- 续跑：checkpoint JSON

优势：

- 依赖轻、部署简单
- 便于跨 provider 统一接入

限制：

- `urllib` 连接复用能力有限，后续可评估 `httpx`

## 3. Evaluator 候选工具（规划）

### 3.1 Java

- 编译：`javac`
- 测试：JUnit（Maven/Gradle）
- 覆盖率：JaCoCo
- 变异测试：PIT

### 3.2 Python

- 语法/静态检查：`python -m py_compile` / `ruff` / `pylint`
- 测试：pytest
- 覆盖率：coverage.py
- 变异测试：mutmut / cosmic-ray

### 3.3 Go

- 编译：`go build`
- 测试：`go test`
- 覆盖率：`go test -cover` / `go tool cover`
- 变异测试：go-mutesting（可行性需验证）

### 3.4 C++

- 编译：g++ / clang++
- 测试：GoogleTest
- 覆盖率：gcov / llvm-cov
- 变异测试：Mull（可行性需验证）

## 4. Reporter 候选工具（规划）

- 数据格式：JSON + CSV
- 可视化：Python（matplotlib / plotly）或前端静态页面
- 报告输出：Markdown + HTML 双格式

## 5. 工具选型原则

- 跨平台优先
- 社区成熟度高
- 命令行可自动化
- 输出可结构化解析

## 6. 集成策略建议

- 先打通最小闭环：编译 + 测试执行 + 覆盖率
- 再增量接入变异测试
- 最后统一报告层做多模型对比

## 7. 风险与规避

- 工具版本差异导致结果不可比
  - 规避：锁定工具版本，记录执行环境
- 样本依赖导致执行不稳定
  - 规避：优先使用自包含样本并隔离外部依赖
- 多语言工具链成本高
  - 规避：按语言分阶段上线，先覆盖主战场
