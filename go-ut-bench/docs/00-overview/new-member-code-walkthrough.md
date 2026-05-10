# 新成员代码导读

本文面向第一次接触 `go-ut-bench/` 的开发者。目标不是覆盖所有细节，而是先建立一套和当前代码一致的心智模型，避免被历史文档或父仓库材料带偏。

## 先记住这三件事

1. 当前有效实现都在 `go-ut-bench/` 下，父仓库根目录不是主代码入口。
2. 这个项目现在评测的不是“模型名”本身，而是统一抽象后的 `subject = framework + model + optional skill`。
3. 主链路是：`dataset -> runner -> evaluator -> reporter -> store/web`，由 `internal/orchestrator/service.go` 串起来。

## 你应该从哪里开始读

推荐按下面顺序读代码：

1. `cmd/utbench/main.go`
2. `internal/orchestrator/service.go`
3. `internal/contracts/spec.go`
4. `internal/contracts/results.go`
5. `internal/dataset/service.go`
6. `internal/runner/service.go`
7. `internal/runner/subjects.go`
8. `internal/evaluator/service.go`
9. `internal/reporter/service.go`
10. `internal/store/sqlite.go`
11. `internal/web/server.go`
12. `internal/web/run_manager.go`

如果只想先跑通主线，前 8 个文件已经够用。

## 项目真实定位

从当前代码看，UT-Bench 已经不是“调用多个 LLM 生成单测并打分”的小工具，而是一个统一评测平台：

- 既支持纯模型 API baseline，也支持 CLI Agent。
- 既支持 prompt 级评测，也支持 skill 注入和 sandbox 隔离。
- 既支持一次性运行，也支持基于 SQLite 的资产沉淀、复用、对比和 Web 管理。

入口代码见 `cmd/utbench/main.go`。当前命令面实际包括：

- `run`
- `generate`
- `evaluate`
- `report`
- `db`
- `assets`
- `dataset`
- `doctor`
- `web`

其中 `assets`、`doctor`、`web` 往往不会出现在旧文档的第一层介绍里，但它们都已经是正式功能。

## 统一抽象：subject

当前系统最重要的抽象不是 model，而是 `subject`。相关加载逻辑在 `internal/agentconfig/config.go` 和 `internal/runner/subjects.go`。

`subject` 统一表示一个被测对象：

```text
subject = framework + model + optional skill
```

典型场景有三类：

- `model_api + model + no_skill`：纯模型基线
- `cli_agent + framework + model + no_skill`：Agent 基线
- `cli_agent + framework + model + skill`：Agent + skill

这件事会影响你读几乎所有模块：

- `runner` 调度的是 `subject × sample`
- `reporter` 比较的是 subject 维度结果
- `store` 复用时会同时看 subject、subject version、prompt/version、环境指纹

如果没有显式提供 `agents.yaml`，系统也会自动生成 `model_api__<model>__no_skill` 这类默认 subject。

## 数据是怎么流的

### 1. CLI / Web 组装 RunSpec

CLI 在 `cmd/utbench/main.go` 解析参数，Web 在 `internal/web/` 组装运行请求。两者最终都汇总到 `contracts.RunSpec`。

`RunSpec` 是贯穿全流程的主配置对象，定义在 `internal/contracts/spec.go`。

### 2. Dataset 发现样本

`internal/dataset/service.go` 负责：

- 扫描目录或读取 manifest
- 根据路径推断语言、类别、场景
- 计算源码 MD5
- 按语言、场景、数量做过滤

当前代码里支持的类别是：

- `self_contained`
- `repo_level`

注意：旧文档里常见的 `module_level` 已不是当前代码契约。

### 3. Runner 生成测试

`internal/runner/service.go` 是生成阶段主入口。

它会做这些事：

- 读取 `models.yaml`
- 读取 `agents.yaml`
- 解析出可运行的 subject 列表
- 为每个 `subject × sample` 生成任务
- 生成 prompt 快照
- 检查是否可以复用历史 generation
- 使用 worker pool 并发执行
- 写出 `generated_manifest.json`

纯模型 API 和 CLI Agent 在这里统一编排，但底层通过 adapter 分开实现：

- `adapter_model.go`
- `adapter_cli.go`

### 4. Evaluator 执行评测

`internal/evaluator/service.go` 负责读取 manifest，然后并发执行：

- compile
- sample tests
- coverage
- mutation

它还有几块比较关键的逻辑：

- 评测环境指纹采集
- 评测结果复用
- 失败归因
- Python repo-level workspace 处理
- mutation policy 和超时控制

### 5. Reporter 聚合分析

`internal/reporter/service.go` 不只是“把结果转成 HTML”。

它真实做了这些聚合：

- summary
- dimensions
- top models
- failures
- score exclusions
- truncation stats
- efficiency stats
- error diagnosis
- agent comparisons
- skill uplifts
- comparison views

所以如果你后面要改评分口径或排名规则，落点一般不只在一个文件里。

### 6. Store 持久化和复用

`internal/store/sqlite.go` 已经不是简单的 run/result 两张表。

当前 schema 包含至少这些主题：

- dataset
- subjects
- subject_versions
- prompt_profiles
- prompt_renderings
- generation_runs
- generated_cases
- evaluation_envs
- evaluation_runs
- evaluation_results
- report_snapshots
- asset 查询面
- automation / notification

从设计意图上看，这一层已经在支持“实验资产管理”，不是单纯“结果落库”。

### 7. Web 是第二主入口

`internal/web/server.go` 和 `internal/web/run_manager.go` 说明 Web 已经是完整控制面，不是只展示报告的薄壳。

它负责：

- 创建和管理 benchmark run
- 本地/ Docker 执行切换
- 暂停/恢复/取消
- 环境检查
- agent/skill 管理
- DB 浏览和报表
- automation / notification API

如果你要改运行生命周期、前后端联动、或者新增管理能力，优先看 Web。

## 目录怎么理解

最值得记住的是这些目录边界：

| 目录 | 真实职责 |
| --- | --- |
| `cmd/utbench/` | CLI 入口和参数编排 |
| `internal/contracts/` | 跨阶段数据契约唯一来源 |
| `internal/dataset/` | 数据集发现、校验、索引、manifest |
| `internal/runner/` | 生成阶段，包括 prompt、model API、agent adapter、checkpoint、sandbox |
| `internal/evaluator/` | 编译、测试、覆盖率、变异、失败归因、环境指纹 |
| `internal/reporter/` | 聚合分析和 HTML/JSON 报告 |
| `internal/store/` | SQLite schema、入库、查询、复用 |
| `internal/web/` | Web API、run manager、环境检查、自动化 |
| `configs/` | 模型、Agent、skill 配置 |
| `datasets/` | benchmark 输入样本 |
| `artifacts/` | 运行产物 |
| `storage/` | SQLite 数据库及运行时存储 |

## 你第一次上手时最容易踩的坑

### 1. 误把父仓库根目录当主模块

真正的 Go module 在 `go-ut-bench/go.mod`。父仓库根目录还有一个占位 `go.mod`，不要在外层跑构建或测试。

### 2. 误以为只有 model baseline

现在代码主线已经把 Agent 评测和 skill 注入并入同一套抽象，不能再按“多模型 LLM benchmark”理解全部设计。

### 3. 误用 `module_level`

代码契约现在使用 `repo_level`。旧词主要出现在历史文档和零散帮助文案里。

### 4. 误以为 Web 只是看报告

当前 Web 已经覆盖运行控制、环境检查、DB 浏览、agent/skill 管理和自动化接口。

### 5. 误以为 SQLite 只是结果存档

当前 SQLite 还承载 generation/evaluation 复用判断和资产查询。

## 建议的阅读方法

如果你接下来要做的事情是：

- 跑通主流程：先读 `main.go`、`orchestrator/service.go`、`dataset/service.go`、`runner/service.go`
- 接模型或 Agent：读 `runner/models.go`、`runner/subjects.go`、`runner/adapter*.go`、`agentconfig/config.go`
- 改 prompt：读 `runner/prompt.go`、`contracts/prompt_*.go`
- 改评测规则：读 `evaluator/service.go` 和对应语言后端
- 改打分或报告：读 `reporter/service.go` 及聚合/排名相关文件
- 改复用策略：读 `runner/reuse_*.go`、`evaluator/service.go`、`store/sqlite.go`
- 改 Web：读 `web/server.go`、`web/run_manager.go`

## 关联文档

- 项目总介绍：`docs/00-overview/project-introduction.md`
- 文件职责：`docs/00-overview/project-file-map.md`
- 架构设计：`docs/02-design/architecture-mvp.md`
- 数据集治理：`docs/04-development/dataset-governance.md`
- 本次基于代码的偏差审计：`docs/04-development/code-audit-2026-05-10.md`

## 一句话总结

把这个项目当成“可管理、可复用、可对比的单测生成评测平台”去理解，后面的代码结构才是连贯的。
