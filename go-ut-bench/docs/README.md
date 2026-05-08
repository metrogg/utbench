# UT-Bench 文档导航

本目录按阅读目的重新分类。优先从本页进入，不建议直接在 `docs/` 根目录里平铺新增文档。

如果需要判断某份文档是否还应该维护、是否与其他文档重合，先看 [文档分类与去重索引](DOCUMENT_CLASSIFICATION.md)。

## 推荐阅读顺序

1. 想快速使用：先读 [用户指南](01-user-guides/USER_GUIDE.md)，再按需要查看 [完整启动指南](01-user-guides/startup-guide.md) 或 [Docker 使用指南](03-operations/DOCKER_GUIDE.md)。
2. 想理解项目：先读 [项目详细介绍](00-overview/project-introduction.md)，再看 [项目文件职责说明](00-overview/project-file-map.md)。
3. 想开发或改造：先读 [架构设计](02-design/architecture-mvp.md)，再按模块查看数据库、Prompt、资产管理、沙箱等设计文档。
4. 想接入 Agent：先读 [用户指南](01-user-guides/USER_GUIDE.md) 中的 Agent/Skill 部分，再看 [Agent 评测升级记录](04-development/AGENT_UPGRADE.md) 和外部工具资料。

## 目录分类

| 目录 | 内容定位 |
| --- | --- |
| [00-overview](00-overview/) | 项目定位、整体能力、文件职责和功能总览。 |
| [01-user-guides](01-user-guides/) | 面向使用者的上手、CLI、Web UI 和日常运行说明。 |
| [02-design](02-design/) | 当前或仍有参考价值的系统设计文档。 |
| [03-operations](03-operations/) | Docker、部署、运行环境和运维相关说明。 |
| [04-development](04-development/) | 数据治理、重构分析、技术审查和迭代记录。 |
| [05-integrations](05-integrations/) | 外部工具、外部 Agent 或第三方集成资料。 |
| [99-archive](99-archive/) | 已废弃、未采用或被新文档替代的历史方案。 |
| [DOCUMENT_CLASSIFICATION.md](DOCUMENT_CLASSIFICATION.md) | 文档保留等级、重合关系和归档建议。 |

## 重复内容处理

目前没有直接删除旧文档，避免丢掉历史上下文；重复关系先在这里明确：

| 重复/重叠主题 | 主入口 | 其他文档处理方式 |
| --- | --- | --- |
| 快速开始、CLI 参数、常见运行方式 | [用户指南](01-user-guides/USER_GUIDE.md) | [完整启动指南](01-user-guides/startup-guide.md) 保留为详细参数手册；[CLI 命令文档](01-user-guides/cli-spec.md) 保留为命令参考。 |
| Docker 快速开始和 Docker 完整说明 | [Docker 使用指南](03-operations/DOCKER_GUIDE.md) | [Docker 快速开始](03-operations/quickstart-docker.md) 保留为精简版流程，后续可并入 Docker 使用指南。 |
| 项目介绍、功能详情、架构说明 | [项目详细介绍](00-overview/project-introduction.md) | [功能详细文档](00-overview/feature-details.md) 保留为功能清单；[架构设计](02-design/architecture-mvp.md) 保留为设计视角。 |
| 自动化调度与通知 | [定时评测与推送服务设计方案](02-design/automation-scheduler-notification-design.md) | 旧版 [自动化评测与通知方案](99-archive/automation-scheduler-notification.md) 已归档，文内也注明未实际采用。 |
| 技术审查、重构分析、升级记录 | [重构分析报告](04-development/refactoring-analysis.md) | [技术审查报告](04-development/technical-review.md) 和 [Agent 升级记录](04-development/AGENT_UPGRADE.md) 保留为阶段性历史记录。 |
| CodeBuddy 使用资料 | [CodeBuddy 指南](05-integrations/codebuddy-guide.md) | 这是外部工具资料，体量较大，不作为项目主文档入口。 |

## 新增文档规则

- 面向用户操作的文档放入 `01-user-guides/`。
- 面向系统设计或数据结构决策的文档放入 `02-design/`。
- 面向部署、Docker、环境排障的文档放入 `03-operations/`。
- 阶段性审查、重构计划、治理说明放入 `04-development/`。
- 外部工具资料放入 `05-integrations/`。
- 被替代、未采用或只保留历史上下文的方案放入 `99-archive/`。
- 新文档请在本页补一行入口说明，避免再次变成平铺文档堆。
