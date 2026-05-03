# UT-Bench Long-term Memory

## Project Overview
- UT-Bench: 从"评测LLM模型生成单元测试能力"升级为"评测编码Agent生成单元测试能力"
- 项目路径: `go-ut-bench/` 主代码
- 支持的subject三元组: framework + model + optional skill
- 已接入的framework: opencode (首个真实framework)
- Agent评测端到端已跑通: agent_smoke_002 (12/12成功)

## Dataset Architecture
- 三层DatasetClass: self_contained(800个) → module_level(1个) → repository_level(规划中)
- 仓库级样本meta.json设计为5维度: source/target/context/environment/evaluation
- 现有module_level样本dateutil是仓库级的雏形，但不具备完整仓库构建上下文

## Key Technical Decisions
- DOOD沙箱方案: 外层benchmark容器通过Docker socket控制宿主机daemon启动内层Agent沙箱
- SandboxRunner抽象: local/docker后端，后续可接gvisor/firecracker/cube
- Skill注入模式: prompt_append / workspace_mount / agent_native(预留)

## User Preferences
- 语言: 中文
- 用户关注仓库级数据集设计

## Agent Benchmark Landscape (2026-05 调研)
- SWE-bench Verified已失效(污染), SWE-bench Pro是当前推荐(多语言,1865 tasks,GPL防污染)
- 单元测试生成方向仅~7篇论文, UT-Bench处于蓝海定位
- 极简主义趋势: mini-swe-agent(100行Python) > 复杂框架
- 成本和延迟是行业盲区, UT-Bench可做差异化
- 必关注: SWE-agent(trajectory), OpenHands(Runtime), DeepEval(50+指标), Aider(错误分类), Langfuse(成本追踪)
