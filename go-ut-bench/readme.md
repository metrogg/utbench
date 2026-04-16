## go-ut-bench

Go 版评测工具工作区（MVP 起步版本）。

目标：保留原有 `runner / evaluator / reporter / dataset` 认知模型，做成可持续迭代的 CLI 工具。

### 当前约定

- 语言范围（MVP）：`python`、`java`、`go`、`cpp`
- 数据集分层（MVP）：
  - 大类：`self_contained` / `module_level`
  - 子类：`boundary` / `simple_function` / `complex_dependency` / `interface_mock`
- 工具运行模式：支持分步执行，也支持一键编排
- 产物目录：`./artifacts`，数据库目录：`./storage`

### 命令设计（草案）

- `utbench generate`：生成单测
- `utbench evaluate`：评测已生成单测
- `utbench report`：生成报告
- `utbench ingest`：结果入库 SQLite
- `utbench run`：编排命令（默认串联 generate -> evaluate -> report，可选 ingest）
- `utbench dataset`：数据集索引/校验/查看

建议的数据集治理流程：

1. `utbench dataset index --dataset-root ./datasets --output ./configs/dataset_index.json`
2. `utbench dataset manifest --index ./configs/dataset_index.json --level l1 --output ./configs/dataset_l1.json --limit-per-scenario 20`
3. `utbench run --dataset-manifest ./configs/dataset_l1.json ...`

关键参数（MVP）：

- `--config`：模型配置文件（默认 `../benchmark/config/models.yaml`）
- `--models`：模型多选（逗号分隔）
- `--langs`：语言多选（逗号分隔）
- `--class`：数据集大类（`self_contained` / `module_level` / `complex_dependency`）
- `--scenario`：数据集子类（`boundary` / `simple_function` / `complex_dependency` / `interface_mock`）
- `--level`：评测集分级（如 `l1`）
- `--dataset-manifest`：显式样本清单（如 `./configs/dataset_l1.json`）
- `--mode`：`full` / `incremental`
- `--reset-checkpoint`：重置当前作用域 checkpoint
- `--mutation-enabled`：开启/关闭变异阶段
- `--mutation-timeout`：变异阶段超时（秒）
- `--mutation-policy`：变异异常策略（`warn` / `fail`）
- `--max-samples`：样本数量限制
- `--dry-run`：跳过模型 API，生成占位测试

### 目录结构

```text
go-ut-bench/
  cmd/utbench/                 # CLI 入口
  internal/
    contracts/                 # 阶段间统一数据契约
    config/                    # 配置加载与校验
    dataset/                   # 数据集管理
    runner/                    # 单测生成
    evaluator/                 # 评测执行
    reporter/                  # 报告生成
    orchestrator/              # 流程编排
    store/                     # 结果持久化（SQLite）
    obs/                       # 日志/观测（预留）
  docs/                        # 架构与规范文档
  migrations/                  # SQLite 初始化与升级脚本
  schemas/                     # JSON 结果协议
  configs/                     # 示例配置
  go.mod
```

### 快速开始（当前脚手架）

```bash
go run ./cmd/utbench help
go run ./cmd/utbench run --dry-run --dataset-root ./datasets --config ../benchmark/config/models.yaml
```

> 注：当前代码为架构脚手架，先定边界和契约，再逐步填充真实执行逻辑。

每次 `run` 会在 `artifacts/runs/<run-id>/run_summary.json` 写出本次执行快照与产物路径。
