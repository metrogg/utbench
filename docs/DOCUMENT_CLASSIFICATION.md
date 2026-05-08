# 父仓库文档分类与归档索引

父仓库层面的 `docs/` 主要保留调研、外部工具资料、历史记录和规格流程产物。当前可执行项目文档以 [go-ut-bench/docs](../go-ut-bench/docs/README.md) 为准。

## 分类规则

| 等级 | 含义 |
| --- | --- |
| 主入口 | 父仓库层面的当前导航或总览。 |
| 研究资料 | 调研、差距分析、方法论和路线建议。 |
| 外部资料 | Claude Code、OpenCode、CodeBuddy 等外部工具文档。 |
| 规格记录 | 由规格流程沉淀的专题设计。 |
| 归档 | 体量大、历史性强或当前不建议作为实现依据的材料。 |

## 文档索引

| 路径 | 等级 | 处理意见 |
| --- | --- | --- |
| [README.md](README.md) | 主入口 | 父仓库文档导航。 |
| [00-overview/PROJECT_ANALYSIS_REPORT.md](00-overview/PROJECT_ANALYSIS_REPORT.md) | 主入口 | 项目分析主入口；如与 Go CLI 当前实现冲突，以 `go-ut-bench/docs` 为准。 |
| [01-research/RESEARCH_REPORT.md](01-research/RESEARCH_REPORT.md) | 研究资料 | Agent 评测生态调研主入口。 |
| [01-research/AGENT_BENCHMARK_SURVEY.md](01-research/AGENT_BENCHMARK_SURVEY.md) | 研究资料 | benchmark 生态清单，作为补充。 |
| [01-research/UTBENCH_GAP_ANALYSIS.md](01-research/UTBENCH_GAP_ANALYSIS.md) | 研究资料 | 差距分析和路线图。 |
| [01-research/ANTHROPIC_EVAL_OPTIMIZATION.md](01-research/ANTHROPIC_EVAL_OPTIMIZATION.md) | 研究资料 | Anthropic eval 方法论参考。 |
| [01-research/REPO_LEVEL_EVALUATION_RESEARCH.md](01-research/REPO_LEVEL_EVALUATION_RESEARCH.md) | 研究资料 | 仓库级评测调研。 |
| [01-research/research_agent_benchmark_architecture.md](01-research/research_agent_benchmark_architecture.md) | 研究资料 | 英文架构模式研究。 |
| [02-development/AGENT_UPGRADE.md](02-development/AGENT_UPGRADE.md) | 规格记录 | 迁移前的 Agent 升级记录；当前版本参考 `go-ut-bench/docs/04-development/AGENT_UPGRADE.md`。 |
| [03-integrations/claudecode官方文档.md](03-integrations/claudecode官方文档.md) | 外部资料 | Claude Code 官方资料摘录。 |
| [03-integrations/CLAUDE_CODE_REFERENCE.md](03-integrations/CLAUDE_CODE_REFERENCE.md) | 外部资料 | Claude Code 长参考资料。 |
| [03-integrations/CODEBUDDY_OFFICIAL_DOCS.md](03-integrations/CODEBUDDY_OFFICIAL_DOCS.md) | 外部资料 | CodeBuddy 官方资料摘录。 |
| [03-integrations/OPENCODE_OFFICIAL_DOCS.md](03-integrations/OPENCODE_OFFICIAL_DOCS.md) | 外部资料 | OpenCode 官方资料摘录。 |
| [04-specs/superpowers/specs/2026-04-23-logging-optimization-design.md](04-specs/superpowers/specs/2026-04-23-logging-optimization-design.md) | 规格记录 | 日志优化专题设计。 |
| [04-specs/superpowers/specs/2026-04-26-test-plan-design.md](04-specs/superpowers/specs/2026-04-26-test-plan-design.md) | 规格记录 | 测试计划专题设计。 |
| [99-archive/history.md](99-archive/history.md) | 归档 | 历史上下文，默认不作为当前实现依据。 |

## 重合内容处理

| 重合主题 | 主入口 | 其他资料处理 |
| --- | --- | --- |
| 当前 Go CLI 使用和开发 | [go-ut-bench/docs](../go-ut-bench/docs/README.md) | 父仓库分析报告只用于背景追溯。 |
| Agent 评测研究 | [RESEARCH_REPORT.md](01-research/RESEARCH_REPORT.md) | survey、gap analysis 和 repo-level research 作为专题补充。 |
| 外部 Agent 工具资料 | [03-integrations/README.md](03-integrations/README.md) | 具体官方资料保留为参考，不直接决定实现。 |
| 历史对话和迁移过程 | [99-archive/history.md](99-archive/history.md) | 只在追溯决策背景时读取。 |
