# 项目文件职责说明

本文用于说明 `go-ut-bench/` 模块内哪些文件属于源码和可复用资产，哪些属于本地运行产物，以及各目录在系统中的职责边界。父仓库根目录下的旧 `README.md`、`docker-compose.yml`、`docker.sh` 仅作为历史材料参考，当前开发以 `go-ut-bench/` 为准。

## 提交边界

| 路径 | 职责 | Git 建议 |
| --- | --- | --- |
| `go.mod`, `go.sum` | Go 模块定义和依赖锁定。 | 提交。 |
| `cmd/utbench/` | CLI 入口，解析子命令并连接 orchestrator、runner、evaluator、reporter、store、web 等服务。 | 提交。 |
| `internal/` | 核心业务代码，按包边界拆分为数据集、生成、评测、报告、持久化、Web 管理等模块。 | 提交。 |
| `configs/*.example.yaml` | 可公开的配置模板。 | 提交。 |
| `configs/models.yaml`, `configs/agents.yaml` | 项目级模型和 Agent 配置；密钥只通过环境变量引用，不写入文件。 | 团队共享配置可提交，个人私有配置不要提交。 |
| `configs/skills/` | 注入到 Agent 工作区的 skill 提示词、参考资料和辅助脚本。 | 提交。 |
| `datasets/` | benchmark 输入样本，按语言和场景组织。 | 提交经过整理的数据集。 |
| `schemas/` | 生成清单、评测结果等 JSON Schema。 | 提交。 |
| `migrations/` | SQLite 表结构迁移脚本。 | 提交。 |
| `docs/` | 当前有效的架构、设计、用户指南、运维和集成文档。 | 提交。 |
| `docker/`, `Dockerfile*` | 应用镜像、评测镜像和 Agent 镜像构建定义。 | 提交。 |
| `build.ps1`, `build.sh` | 本地或 CI 构建 Docker 镜像的辅助脚本。 | 提交。 |
| `run_bench.ps1`, `run_bench.sh` | 使用预构建镜像运行 benchmark 的便捷脚本。 | 提交。 |
| `.env.example` | 环境变量模板，说明需要配置哪些密钥。 | 提交。 |
| `.env`, `.env.*` | 本地密钥和机器相关环境配置。 | 不提交。 |
| `artifacts/` | 运行目录、报告、manifest、Agent 工作区、日志、临时构建文件。 | 只提交 `artifacts/.gitkeep`。 |
| `storage/` | 本地 SQLite 数据库、WAL 文件和运行时存储。 | 只提交 `storage/.gitkeep`。 |
| `.claude/`, `.codebuddy/`, `.opencode/`, `.windsurf/`, `.workbuddy/` | 本地助手或 IDE 状态。 | 不提交。 |
| `.m2-cache/`, `.gocache/`, `.gomodcache/`, `node_modules/`, `target/`, `.gradle/` | 工具链或依赖缓存。 | 不提交。 |

## 顶层目录职责

| 路径 | 具体职责 |
| --- | --- |
| `cmd/utbench/main.go` | 注册 CLI 子命令，包括 `run`、`generate`、`evaluate`、`report`、`db`、`dataset`、`doctor`、`web` 等入口。 |
| `configs/` | 保存模型、Agent、应用配置模板，以及 Agent skill 资产。 |
| `datasets/` | 保存 Python、Go、Java、C++ 的自包含或模块级样本；评测流程只应读取这里的输入样本，不应写入运行结果。 |
| `docker/` | 保存额外 Docker 构建上下文，尤其是 Agent 运行环境相关镜像。 |
| `docs/` | 保存 `go-ut-bench/` 模块自己的项目文档，和父仓库 `docs/` 的调研归档区分开。 |
| `internal/` | Go 内部包，承载 CLI 背后的所有实现。 |
| `migrations/` | 数据库版本化结构定义，当前 SQLite 初始化由这里驱动。 |
| `schemas/` | 对外 JSON 结果文件的结构约束，便于校验和工具集成。 |
| `artifacts/` | 运行时输出根目录，包括生成结果、评测结果、报告、检查点和临时工作区。 |
| `storage/` | 本地持久化数据目录，通常保存 `utbench.db`。 |

## `internal` 模块职责

| 路径 | 具体职责 |
| --- | --- |
| `internal/agentconfig` | 加载和校验 CLI Agent、framework、skill 注入等配置。 |
| `internal/config` | 提供通用应用配置结构和配置读取辅助逻辑。 |
| `internal/contracts` | 跨阶段数据契约的唯一来源，包括 `RunSpec`、`SampleRef`、manifest、评测结果、报告结构、Schema 版本、prompt 元数据和公共错误类型。 |
| `internal/ctrl` | 提供暂停/恢复控制门，用于 Web 触发的运行控制。 |
| `internal/dataset` | 发现样本、构建索引、生成 manifest、校验数据集结构，并处理跨平台进程辅助逻辑。 |
| `internal/evaluator` | 执行各语言的编译、测试、覆盖率和变异测试；包含 Python、Go、Java、C++ 的语言后端和环境指纹采集。 |
| `internal/obs` | 封装结构化日志和终端进度展示，避免业务模块直接依赖全局日志状态。 |
| `internal/orchestrator` | 编排完整 benchmark 流程：数据集选择、生成、评测、报告和可选入库。 |
| `internal/reporter` | 聚合评测结果，计算排名和综合分，生成洞察、运行摘要和内嵌静态资源的 HTML 报告。 |
| `internal/runner` | 调用模型或 Agent 生成单测，负责 prompt 构造、模型配置、API 适配、自动续写、检查点、沙箱和生成复用。 |
| `internal/store` | 管理 SQLite schema、结果入库、查询、列表、运行记录和自动化任务存储。 |
| `internal/web` | 提供 Web 管理 UI、API handler、环境检查、Docker/local 运行、构建管理、自动化调度和静态资源服务。 |

## 重要文件职责

| 路径 | 具体职责 |
| --- | --- |
| `internal/contracts/spec.go` | 定义运行规格、样本引用和仓库级样本元数据。 |
| `internal/contracts/constants.go` | 定义 Schema 版本、数据集类别、运行模式等常量。 |
| `internal/contracts/results.go` | 定义生成清单、评测结果集和报告载荷等跨阶段结果结构。 |
| `internal/runner/prompt.go` | 构建 `full_file`、`completion`、`repo_level` 三类 prompt。 |
| `internal/runner/api.go` | 实现 LLM API 调用、重试、截断检测和自动续写。 |
| `internal/runner/models.go` | 从 YAML 加载模型 provider、endpoint、api key 环境变量等配置。 |
| `internal/runner/service.go` | 组织生成任务并行执行、检查点写入和生成结果落盘。 |
| `internal/evaluator/service.go` | 组织评测流水线，调度语言后端并汇总指标。 |
| `internal/evaluator/*_eval*.go` | 分语言实现编译、测试、覆盖率和结果解析。 |
| `internal/evaluator/mutation*.go` | 变异测试公共结构和平台差异处理。 |
| `internal/reporter/service.go` | 报告生成主流程。 |
| `internal/reporter/aggregation.go` | 按模型、语言、场景聚合指标。 |
| `internal/reporter/ranking.go` | 计算排名和综合得分。 |
| `internal/reporter/static/` | HTML 报告内嵌使用的 CSS 和 JavaScript。 |
| `internal/store/sqlite.go` | SQLite 连接、schema 初始化和核心写入逻辑。 |
| `internal/store/ingest.go` | 将运行产物导入数据库。 |
| `internal/web/server.go` | Web 服务路由注册和静态资源挂载。 |
| `internal/web/*_handlers.go` | Web API handler，按模型、环境、密钥、自动化等功能拆分。 |
| `internal/web/static/` | Web 管理 UI 的 HTML、CSS、JavaScript 和页面片段。 |

## 数据集目录约定

自包含 benchmark 样本放在：

```text
datasets/<language>/<language>_code_files_self_contained/<scenario>/<sample_id>.<ext>
```

当前主要语言包括 `python`、`go`、`java`、`cpp`。常见场景包括 `simple_function`、`boundary`、`complex_dependency`、`interface_mock`。生成出的测试、复制出的工作区、包缓存、覆盖率文件、变异测试输出和 Agent trace 都是运行产物，应进入 `artifacts/` 或临时工作区，不能混入 `datasets/`。

## Git 清理建议

应提交：源码、测试、经过整理的 benchmark 输入数据、可复用配置模板、文档、schema、数据库迁移、Docker 文件和脚本。

不应提交：本地密钥、个人助手状态、编译后二进制、SQLite 数据库、日志、覆盖率输出、生成的 manifest/result、包缓存、benchmark 运行工作区。

如果某个文件已经被 Git 跟踪，但之后希望它只保留在本地，可以在确认不需要保留历史版本后执行：

```bash
git rm --cached <path>
```
