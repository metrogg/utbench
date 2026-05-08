# UT-Bench 文档导航

本目录用于收纳父仓库层面的项目资料、调研报告、外部工具文档和历史记录。根目录只保留 `README.md`、`CLAUDE.md`、`go.mod` 这类入口或工具约定文件。

当前可执行项目文档以 [go-ut-bench/docs](../go-ut-bench/docs/README.md) 为准；父仓库文档的保留等级见 [父仓库文档分类与归档索引](DOCUMENT_CLASSIFICATION.md)。

## 推荐阅读顺序

1. 快速了解项目：先读 `README.md`，再读 [项目分析报告](00-overview/PROJECT_ANALYSIS_REPORT.md)。
2. 查研究背景：读 [研究与差距分析](01-research/README.md)，其中 `RESEARCH_REPORT.md` 和 `UTBENCH_GAP_ANALYSIS.md` 是主入口。
3. 查开发过程：读 [开发记录](02-development/README.md)。
4. 查外部 Agent/CLI 工具：读 [外部集成资料](03-integrations/README.md)。
5. 查历史上下文：读 [历史归档](99-archive/README.md)，默认不作为当前实现依据。

## 分类目录

| 目录 | 内容定位 |
| --- | --- |
| [00-overview](00-overview/) | 项目整体分析、当前能力和结构概览。 |
| [01-research](01-research/) | 行业调研、差距分析、架构优化建议、repo-level 评测研究。 |
| [02-development](02-development/) | 阶段性开发记录和升级状态。 |
| [03-integrations](03-integrations/) | Claude Code、OpenCode、CodeBuddy 等外部工具资料。 |
| [04-specs](04-specs/) | 由规格/设计流程沉淀的专题方案。 |
| [99-archive](99-archive/) | 体量大、历史性强、当前不建议作为主入口的材料。 |
| [DOCUMENT_CLASSIFICATION.md](DOCUMENT_CLASSIFICATION.md) | 父仓库文档保留等级、重合关系和归档说明。 |

## 功能重合处理

| 重合主题 | 主入口 | 其他资料处理 |
| --- | --- | --- |
| 项目整体分析 | [PROJECT_ANALYSIS_REPORT.md](00-overview/PROJECT_ANALYSIS_REPORT.md) | 根目录 `README.md` 保留为快速入口。 |
| Agent 评测生态调研 | [RESEARCH_REPORT.md](01-research/RESEARCH_REPORT.md) | `AGENT_BENCHMARK_SURVEY.md` 更偏 benchmark 生态清单，作为补充。 |
| UT-Bench 差距和优化路线 | [UTBENCH_GAP_ANALYSIS.md](01-research/UTBENCH_GAP_ANALYSIS.md) | `ANTHROPIC_EVAL_OPTIMIZATION.md` 从 Anthropic eval 方法论补充优化视角。 |
| 仓库级评测 | [REPO_LEVEL_EVALUATION_RESEARCH.md](01-research/REPO_LEVEL_EVALUATION_RESEARCH.md) | `research_agent_benchmark_architecture.md` 更偏英文架构模式总结，作为参考。 |
| Claude Code 接入 | [claudecode官方文档.md](03-integrations/claudecode官方文档.md) | `CLAUDE_CODE_REFERENCE.md` 是更长的官方资料摘录，作为参考文档。 |
| 项目历史上下文 | [history.md](99-archive/history.md) | 文件很大，默认只在追溯历史时读取。 |

## 新增文档规则

- 当前有效的项目说明放入 `00-overview/`。
- 调研、竞品分析、差距分析、架构建议放入 `01-research/`。
- 阶段性实施记录放入 `02-development/`。
- 外部工具官方资料、接入说明放入 `03-integrations/`。
- 规格流程产生的专题设计放入 `04-specs/`。
- 过大、过旧、重复或仅用于追溯的材料放入 `99-archive/`。
