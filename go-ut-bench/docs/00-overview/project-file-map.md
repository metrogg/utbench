# 项目文件职责说明

这份文档用于说明项目中哪些文件属于源码资产，哪些属于本地运行产物，以及哪些文件通常不应该提交到 Git。

## 顶层文件和目录

| 路径 | 具体职责 | 提交建议 |
| --- | --- | --- |
| `go.mod`, `go.sum` | Go 模块定义和依赖校验锁定文件。 | 提交。 |
| `cmd/utbench/` | `utbench` 命令行入口，负责解析命令并串联各服务模块。 | 提交。 |
| `internal/` | 项目主要业务实现，按职责边界拆分为多个内部包。 | 提交。 |
| `configs/*.example.yaml` | 可公开的配置模板。 | 提交。 |
| `configs/models.yaml`, `configs/agents.yaml` | 项目级模型和 Agent 配置；密钥应放在环境变量里，不要写进这些文件。 | 如果配置对团队可复用，则提交。 |
| `configs/skills/` | 注入到 Agent 工作区的 skill 提示词、参考资料和辅助脚本。 | 提交。 |
| `datasets/` | 按语言和场景组织的 benchmark 源样本，是评测语料的一部分。 | 提交经过整理的数据集。 |
| `schemas/` | 生成清单和评测结果的 JSON Schema。 | 提交。 |
| `migrations/` | SQLite 数据库结构迁移脚本。 | 提交。 |
| `docs/` | 架构、用户指南、设计说明和项目约定文档。 | 提交。 |
| `docker/`, `Dockerfile*`, `.dockerignore` | 应用和 Agent 环境的容器构建定义。 | 提交。 |
| `build.ps1`, `build.sh`, `run_bench.ps1`, `run_bench.sh` | 本地构建和 benchmark 运行辅助脚本。 | 提交。 |
| `.env.example` | 环境变量模板，用于说明需要配置哪些变量。 | 提交。 |
| `.env` | 本地密钥和机器相关环境配置。 | 不提交。 |
| `artifacts/` | 运行报告、Agent 工作区、日志、manifest、临时构建文件等生成产物。 | 只保留并提交 `artifacts/.gitkeep`。 |
| `storage/` | 本地 SQLite 数据库、WAL 文件和构建缓存。 | 只保留并提交 `storage/.gitkeep`。 |
| `.claude/`, `.workbuddy/` | 本地助手或工具状态文件。 | 不提交。 |

## `internal` 模块职责

| 路径 | 具体职责 |
| --- | --- |
| `internal/agentconfig` | 加载并校验 CLI Agent 和 framework 配置。 |
| `internal/config` | 通用应用配置辅助逻辑。 |
| `internal/contracts` | 跨模块共享的数据契约，包括 spec、结果结构、prompt 元数据、常量和类型化错误。 |
| `internal/ctrl` | 轻量执行门控和控制原语。 |
| `internal/dataset` | 数据集索引、manifest 生成、校验，以及不同操作系统下的进程辅助逻辑。 |
| `internal/evaluator` | Python、Go、Java、C++ 的编译、测试、覆盖率和变异测试评测逻辑。 |
| `internal/obs` | 日志和进度展示工具。 |
| `internal/orchestrator` | 端到端 benchmark 流程编排。 |
| `internal/reporter` | 结果聚合、排名、洞察生成和内嵌 HTML 报告资源。 |
| `internal/runner` | 模型/Agent 调用、prompt 构造、沙箱、subject 展开和生成结果复用。 |
| `internal/store` | SQLite 持久化、结果入库、查询、列表和自动化任务存储。 |
| `internal/web` | Web 服务、API handler、自动化调度/通知、环境检查、Docker runner 和静态 UI。 |

## 数据集目录约定

整理后的 benchmark 样本应放在：

```text
datasets/<语言>/<语言>_code_files_self_contained/<场景>/<样本ID>.<扩展名>
```

当前预期语言包括 `python`、`go`、`java` 和 `cpp`。常见场景包括 `boundary`、`simple_function`、`complex_dependency` 和 `interface_mock`。

生成出来的测试、复制出来的工作区、包缓存、覆盖率文件、变异测试输出和 Agent trace 都属于运行产物，应放在 `artifacts/` 或临时工作区中，不应该混入 `datasets/`。

## Git 提交边界

应该提交：源码、经过整理的 benchmark 输入数据、可复用配置模板、文档、schema、数据库迁移脚本、Docker 文件和构建脚本。

不应该提交：本地密钥、本地助手状态、编译后的二进制文件、SQLite 数据库、日志、覆盖率输出、生成的 manifest/result、包缓存、benchmark 运行工作区。

如果某个文件已经被 Git 跟踪，但之后希望它只保留在本地，可以保留磁盘文件并从 Git 索引移除：

```bash
git rm --cached <path>
```
