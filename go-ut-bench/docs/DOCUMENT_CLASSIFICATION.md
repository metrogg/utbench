# 文档分类与去重索引

本页用于回答两个问题：哪些文档是当前入口，哪些只是补充或历史材料；哪些内容功能重合，后续可以合并或归档。

## 分类规则

| 等级 | 含义 | 维护要求 |
| --- | --- | --- |
| 主入口 | 日常使用、开发和排障优先阅读的文档 | 保持最新，入口链接必须可用 |
| 参考 | 仍有价值，但不是第一阅读入口 | 内容变动时按需更新 |
| 阶段记录 | 记录某次设计、审查、升级或治理过程 | 不要求持续同步，只保留上下文 |
| 外部资料 | 外部工具、Agent 或第三方资料 | 不作为 UT-Bench 行为依据 |
| 归档 | 已废弃、未采用或被新方案替代 | 默认不读取，除非追溯历史 |
| 生成产物 | 运行、缓存、模型输出或构建产物 | 不纳入人工文档维护 |

## 主入口

| 文档 | 分类 | 说明 |
| --- | --- | --- |
| [README.md](README.md) | 主入口 | 当前文档总导航和分类规则。 |
| [01-user-guides/USER_GUIDE.md](01-user-guides/USER_GUIDE.md) | 主入口 | 面向使用者的总指南，覆盖快速开始、运行模式、Agent/Skill 和常见问题。 |
| [01-user-guides/startup-guide.md](01-user-guides/startup-guide.md) | 主入口 | 完整启动和参数说明，适合执行前查命令。 |
| [03-operations/DOCKER_GUIDE.md](03-operations/DOCKER_GUIDE.md) | 主入口 | Docker 构建、运行、DOOD 和沙箱边界说明。 |
| [02-design/architecture-mvp.md](02-design/architecture-mvp.md) | 主入口 | 系统架构、模块职责和 MVP 设计。 |
| [02-design/output-convention.md](02-design/output-convention.md) | 主入口 | 运行产物、checkpoint、数据库和 Agent 过程产物落盘规范。 |

## 按目录分类

| 目录 | 文档 | 等级 | 处理意见 |
| --- | --- | --- | --- |
| `00-overview/` | `project-introduction.md` | 主入口 | 项目介绍主入口。 |
| `00-overview/` | `project-file-map.md` | 参考 | 判断文件职责和提交边界时使用。 |
| `00-overview/` | `feature-details.md` | 参考 | 与项目介绍有重合，保留为能力清单；后续可抽取差异内容后合并。 |
| `01-user-guides/` | `USER_GUIDE.md` | 主入口 | 用户文档优先入口。 |
| `01-user-guides/` | `startup-guide.md` | 主入口 | 启动和参数手册。 |
| `01-user-guides/` | `cli-spec.md` | 参考 | 命令参考，避免和用户指南重复扩写场景教程。 |
| `01-user-guides/` | `WEB_UI.md` | 参考 | Web 页面和 API 说明。 |
| `02-design/` | `architecture-mvp.md` | 主入口 | 架构设计入口。 |
| `02-design/` | `database-design.md` | 主入口 | SQLite schema 和入库查询设计。 |
| `02-design/` | `asset-management-design.md` | 主入口 | 资产复用和评测复用设计。 |
| `02-design/` | `output-convention.md` | 主入口 | 输出目录和产物命名规范。 |
| `02-design/` | `prompt-design.md` | 参考 | Prompt 构造和后处理设计。 |
| `02-design/` | `runtime-sandbox-redesign.md` | 参考 | 沙箱后续改造方案。 |
| `02-design/` | `automation-scheduler-notification-design.md` | 参考 | 当前采用方向的定时评测和通知设计。 |
| `03-operations/` | `DOCKER_GUIDE.md` | 主入口 | Docker 主文档。 |
| `03-operations/` | `quickstart-docker.md` | 参考 | 与 Docker 指南重合，保留为精简流程；后续可并入 Docker 指南。 |
| `04-development/` | `dataset-governance.md` | 参考 | 数据集治理口径。 |
| `04-development/` | `refactoring-analysis.md` | 阶段记录 | 重构分析和阶段计划，保留但不作为当前架构唯一依据。 |
| `04-development/` | `technical-review.md` | 阶段记录 | 技术审查问题清单。 |
| `04-development/` | `AGENT_UPGRADE.md` | 阶段记录 | Agent 评测升级方案与实施状态。 |
| `05-integrations/` | `codebuddy-guide.md` | 外部资料 | 外部工具资料，体量较大，不作为主阅读入口。 |
| `99-archive/` | `automation-scheduler-notification.md` | 归档 | 旧版自动化评测方案，已被当前设计替代。 |

## 功能重合处理

| 重合主题 | 主入口 | 保留文档 | 处理建议 |
| --- | --- | --- | --- |
| 快速开始、CLI 参数、日常运行 | [USER_GUIDE.md](01-user-guides/USER_GUIDE.md) | `startup-guide.md`、`cli-spec.md` | 用户指南讲路径，启动指南讲完整命令，CLI spec 只保留命令参考。 |
| Docker 快速开始与完整 Docker 说明 | [DOCKER_GUIDE.md](03-operations/DOCKER_GUIDE.md) | `quickstart-docker.md` | quickstart 只保留最短路径，重复解释放回 Docker 指南。 |
| 项目介绍、功能详情、架构说明 | [project-introduction.md](00-overview/project-introduction.md) | `feature-details.md`、`architecture-mvp.md` | 项目介绍讲是什么，架构设计讲怎么实现，功能详情只当清单。 |
| 自动化调度与通知 | [automation-scheduler-notification-design.md](02-design/automation-scheduler-notification-design.md) | `99-archive/automation-scheduler-notification.md` | 旧版只留历史背景，不再扩展。 |
| 技术审查、重构分析、升级记录 | [refactoring-analysis.md](04-development/refactoring-analysis.md) | `technical-review.md`、`AGENT_UPGRADE.md` | 都是阶段记录，新增结论应沉淀回设计或用户文档。 |
| CodeBuddy/Skill 资料 | [codebuddy-guide.md](05-integrations/codebuddy-guide.md) | `configs/skills/**` | `docs` 只放说明，实际 skill 文件继续留在 `configs/skills/` 供运行使用。 |

## 不纳入人工文档整理的内容

| 路径 | 原因 |
| --- | --- |
| `artifacts/**` | 运行产物、模型输出、prompt 渲染结果和临时工作区。 |
| `storage/**` | 数据库和构建缓存。 |
| `.m2-cache/**` | Maven 缓存。 |
| `utbench`、`utbench.exe` | 构建产物。 |

## 后续整理建议

- 将 `quickstart-docker.md` 的独有内容并入 `DOCKER_GUIDE.md` 后归档或删除。
- 将 `feature-details.md` 中仍有效的能力清单浓缩到 `project-introduction.md`，减少概览类重复。
- 阶段记录只追加结论，不再扩写教程；可执行操作应回写到用户指南、设计文档或运维文档。
