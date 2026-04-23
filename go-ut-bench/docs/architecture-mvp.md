# go-ut-bench 架构设计（MVP）

## 1. 目标与背景

go-ut-bench 是 ut-bench 的 Go 版本评测工具工作区。

当前阶段目标：

- 保留原有认知模型：`dataset / runner / evaluator / reporter`
- 做成可编译的 CLI 工具，支持分步执行与一键编排
- 评测数据双落地：文件产物 + SQLite
- 首期聚焦 L1 能力闭环（先保证链路稳定，再逐步扩展）

## 2. MVP 范围

### 2.1 支持内容

- 子命令：`generate`、`evaluate`、`report`、`ingest`、`run`、`dataset`
- 语言范围（当前约定）：`python`、`java`、`go`、`cpp`
- 数据集分类（当前约定）：
  - 大类：`self_contained` / `module_level`
  - 子类：`boundary` / `simple_function` / `complex_dependency` / `interface_mock`
- 入库：SQLite（本地文件）

### 2.2 暂不支持

- 微服务化部署
- 分布式任务调度
- 全语言真实 toolchain 执行器（当前先留接口与脚手架）

## 3. 核心流程

### 3.1 分步流程

1. `generate`：选择模型 + 样本，生成单测并产出 `generated_manifest.json`
2. `evaluate`：读取 manifest，执行评测并产出 `evaluation_result.json`
3. `report`：读取评测结果，输出 JSON/HTML 报告
4. `ingest`：将评测结果写入 SQLite

### 3.2 一键流程

`run` 子命令内部按顺序编排：

`dataset -> generate -> evaluate -> report -> (optional ingest)`

数据集治理链路（推荐）：

`dataset index -> dataset manifest -> run`

## 4. 模块划分与职责

### 4.1 保留模块

- `internal/dataset`：样本发现、分类过滤、布局校验
- `internal/runner`：生成单测产物与 manifest
- `internal/evaluator`：执行评测并计算指标
- `internal/reporter`：聚合结果并生成报告

### 4.2 新增支撑模块

- `internal/orchestrator`：一键编排服务
- `internal/store`：SQLite 初始化与入库
- `internal/contracts`：跨阶段统一数据契约
- `internal/config`：默认配置与配置校验
- `internal/obs`：结构化日志能力

## 5. 数据契约

统一由 `internal/contracts` 管理：

- `RunSpec`：一次执行参数快照
- `GeneratedManifest`：生成阶段输出清单
- `EvaluationResultSet`：评测阶段结果集合
- `ReportPayload`：报告阶段聚合结果

设计原则：

- 每类产物都带 `schema_version`
- 支持后续字段增量，不破坏旧版本读取

## 6. 数据存储设计

### 6.1 文件产物

默认输入数据集路径：`./datasets`

默认输出路径：`<output-root>/runs/<run-id>/`

- `generated/generated_manifest.json`
- `generated/tests/<model>/<lang>/*.test.*`
- `generated/metadata/*.response.json`
- `evaluation/evaluation_result.json`
- `report/report_summary.json`
- `report/report.html`
- `run_summary.json`

报告 JSON 额外提供 `mutation_breakdown`（total/killed/survived/no_tests/timeouts/skipped/suspicious）。

其中：

- `generated_manifest.json` 是 evaluate 阶段唯一必需输入
- `response.json` 保留原始模型响应，用于排障与追溯

### 6.2 SQLite

当前最小表：

- `runs`
- `sample_results`

约束：

- `sample_results` 使用 `(run_id, model, language, sample_id)` 唯一约束，支持幂等 upsert

## 7. 并发与性能

- `runner` 与 `evaluator` 均使用 worker pool 并发处理样本
- worker 数默认基于 CPU 自动估算，并限制上界
- 任务失败不影响其他样本执行（样本级隔离）
- `runner` 在 `incremental` 模式下会写入 checkpoint（按 model/lang/sample 去重）

附加说明：

- runner 支持 API 重试与退避（仅可重试错误）
- evaluator 在 Python 场景使用隔离临时目录执行，降低导入污染风险

## 8. 可靠性与错误处理

- CLI 参数前置校验
- 当前数据集校验命令为 `dataset stats`，它会先校验布局，再尝试输出依赖扫描结果
- 样本级错误记录到结果结构，不中断全局 run
- 入库操作采用事务，避免部分写入

## 9. 安全与合规

- 认证信息仅允许环境变量传入（后续模型适配层实现）
- 报告与日志不写出密钥
- 输出路径统一在 `output-root` 下构造，避免不受控落盘

## 10. 可扩展性设计

- 新语言：在 evaluator 增加语言执行器
- 新指标：在 evaluator 增加指标采集器并扩展 contracts
- 新存储：在 store 增加新后端实现（保持 ingest 接口稳定）

## 11. 迭代计划建议

### v0.1（当前）

- 完成 CLI 骨架、目录结构、数据契约、最小入库能力

### v0.2

- 接入真实模型调用
- 接入 Python 真实评测链路（pytest/coverage/mutation）

当前进展：

- 已完成真实模型调用与错误分级重试（runner）
- 已完成 Python 编译/执行/覆盖率链路（evaluator）
- 已接入 mutmut 变异执行框架（结果按执行状态返回）

### v0.3

- 数据集 catalog/manifest 机制
- 报告维度与对比能力增强

## 12. 与旧版 Python 工具关系

- 保留旧版阶段概念，便于团队迁移
- 新版重点补齐：
  - 数据层（SQLite）
  - 结构化契约
  - 更清晰的模块边界
