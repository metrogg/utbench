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

## Gap Analysis (2026-05-03)
- 9维评估: 数据集(🔴)、沙箱(🔴)、Trace(🔴)、健壮性(🔴)为最大短板; 指标/公平性/成本/Prompt(🟡); 报告(🟢)
- Top 5紧急优化: 源码只读挂载(D8) > must_pass_existing(D1/D18) > 交互轮次上限(D9) > 评测重试(D23) > meta.json扩展(D4)
- UT-Bench独特优势(所有竞品没有): 覆盖率+变异测试+Skill评测
- 详细分析: UTBENCH_GAP_ANALYSIS.md

## 3D Evaluation Architecture (2026-05-04)
- 用户核心诉求: Agent×Model×Skill三维评测，控制变量分离各自贡献
- 例: opencode+deepseek vs deepseek纯API → Agent增益; 同Agent不同Model → 模型差距; 同Model+Agent加Skill → Skill增益
- 5层架构设计: Sandbox(可插拔) → Harness(YAML配置) → Engine(事件源) → Metrics(可组合) → Dimension(3D聚合)
- 开源借鉴: SWE-bench(分层Docker缓存+双重验证), SWE-agent(YAML配置+ACI+Guardrail), OpenHands(EventStream+多后端Runtime), DeepEval(Strategy Pattern+统一0-1评分), mini-swe-agent(极简基线)
- 现有代码基础好: SubjectID=framework__model__skill已编码三维, GeneratedCase/EvaluationResult已有AgentFramework/SkillName字段
- 主要改动: 报告聚合层需3D交叉表, EventStream维度标签, Metric Interface统一化, 新增ByAgent/BySkill/ByAgentModelSkill聚合

## Anthropic Eval方法论优化 (2026-05-04)
- 来源: Anthropic "Demystifying Evals for AI Agents" (2026-01-09)
- 关键洞察导致架构从5层→6层,核心新增2层:
  - Layer 4: Multi-Grader Composition (Code+LLM+Human三层Grader,取代原Metrics层)
  - Layer 6: Eval Lifecycle & Saturation (Capability vs Regression双轨,饱和度监控,"毕业"机制)
- 原有4层增强:
  - Layer 1 Sandbox: +Trial隔离校验(pre-run hash验证)
  - Layer 2 Harness: +Reference Solution验证+Balanced Dataset(难度分级+正负样本)
  - Layer 3 Engine: +Outcome/Process分离+Partial Credit+trial维度
  - Layer 5 Aggregation: +pass@k/pass^k三维切片
- 7大关键差距(Anthropic视角): Multi-Grader🔴、pass@k🔴、Partial Credit🔴、Lifecycle🔴、Balanced Dataset🔴、Reference Solution🟡、Swiss Cheese🔴
- 工作量重估: 15天(原10.5天+4.5天新增),详细: ANTHROPIC_EVAL_OPTIMIZATION.md
- Anthropic验证了UT-Bench独特优势: Coverage+Mutation+Skill+3D Matrix,且pass@k三维切片是全行业独有

## Repo-Level评测调研 (2026-05-04)
- 调研了13个行业项目,报告: REPO_LEVEL_EVALUATION_RESEARCH.md
- 最接近UT-Bench的项目: TestGenEval(同为测试生成,有@pass指标和变异得分)
- 环境构建方案: SWE-Bench++的模板+LLM修正(5轮),构建失败后降级到RepoST沙盒模式
- 数据集构建5阶段Pipeline: 仓库采集验证→目标函数提取→环境构建→黄金测试收集→质量保证
- meta.json扩展5维度: source(base_commit/license/stars) + target(files/functions/classes) + context(relevant_files/cross_deps) + environment(pinned_deps/dockerfile_template) + evaluation(existing_tests_expected_pass/must_pass_existing)
- 关键借鉴: SWE-bench FAIL/PASS契约, TestGenEval @pass指标, FEA-Bench仓库选择策略, CrossCodeEval跨文件依赖提取, SWE-Bench++ 4层AutoQA
- 实施路线: Phase1(1-2周,5仓库)→Phase2(2-4周,30-50仓库+模板)→Phase3(1-2月,100+仓库+4语言)→Phase4(2-3月,Leaderboard+论文)
- UT-Bench vs 竞品差异: 任务不是"修Bug"而是"为覆盖不足的函数生成测试",评测指标是覆盖率+变异得分提升

## 首次完整Run数据 (2026-05-04)
- 160样本: 5Subjects × 4语言(C++/Go/Java/Python) × 4场景
- Subjects: model_api(deepseek-v4-flash, minimax2.7), opencode(+deepseek, +minimax), opencode+skill(+deepseek)
- **Agent模式质变**: 测试通过率30-38%(API)→74-94%(Agent), 变异得分28-42%→56-76%
- **行覆盖率是假象**: API已达85-92%, Agent仅+1-7%; 变异得分差距30-50pp
- **最佳区分语言**: C++最难(75%编译/45%通过), Java最简单(API即82.5%通过)
- **Skill增益微弱**: +3.2pp通过率, +2.4pp变异(31样本可能不够统计显著)
- **成本**: Agent模式token消耗~50x(API 4k vs Agent 195k)
- **异常**: MiniMax Java变异得分负增益(-12.8pp), Agent循环引入错误

## 新Agent接入决策 (2026-05-05)
- **Claude Code胜出, Codex CLI出局**: Codex v0.80.0+强制Responses API,第三方模型全不兼容
- **Claude Code多模型方案**: 各厂商均提供Anthropic兼容端点(`/anthropic`),只需3个环境变量(BASE_URL+AUTH_TOKEN+MODEL)
- **厂商兼容接口**: DeepSeek(api.deepseek.com/anthropic)、百炼(dashscope.aliyuncs.com/apps/anthropic)、智谱(open.bigmodel.cn/api/anthropic)、MiniMax(api.minimaxi.com/anthropic 国内/api.minimax.io 国际)、豆包(ark.cn-beijing.volces.com/api/coding Coding Plan或/api/compatible)
- **Claude Code headless**: `claude -p --bare --output-format stream-json --verbose --max-turns 50 --permission-mode bypassPermissions`
- **优势vs CodeBuddy**: 零配置文件(纯env注入)、NDJSON丰富事件流、原生total_cost_usd/token追踪
- **stream-json陷阱**: assistant事件同一message.id多次发射需去重; tool_result在user事件中不在assistant中; result字段是total_cost_usd不是cost_usd
- **6层适配**: agents.yaml + Dockerfile + injectAgentNativeSkill(.claude/skills/) + parseAgentOutput(stream-json) + sessionExport(按需) + smoke test
- **完整参考文档**: docs/CLAUDE_CODE_REFERENCE.md

## CodeBuddy Agent 接入 (2026-05-04)
- 已完成代码层面接入: docker/agents/codebuddy/Dockerfile + configs/agents.yaml + Go解析代码
- CodeBuddy模型配置: models.json({id,url(/chat/completions结尾),apiKey:${ENV_VAR}}), 非简单env var
- headless命令: `codebuddy -p -y --output-format json --model <id> --max-turns 50 "prompt"`
- 模型注入方式: heredoc写models.json + sed替换(避免Go template/shell变量冲突)
- Go代码新增: parseCodeBuddyJSONOutput(), extractFinalJSONObject(), 4个测试全PASS
- 待做: Linux+Docker环境构建镜像并smoke test
