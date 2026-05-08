# Agent 评测生态开源项目调研

> 调研时间：2026-05-02
> 调研范围：AI 编码 Agent 评测 Benchmark、Agent 框架、评测工具/平台、单元测试生成专项

---

## 一、生态总览

AI 编码 Agent 评测生态在 2025-2026 年经历了**爆发式增长**。根据 Awesome-Repo-Level-Code-Generation 的统计，仅仓库级代码任务相关论文就超过 **198 篇**（其中 Issue 修复 ~70 篇，Benchmark 65+ 个），且 2025-2026 一年半就贡献了 100+ 篇。

研究可分为 **7 大方向**：

| 方向 | 论文数 | 活跃度 | 代表项目 |
|:---|:---|:---|:---|
| Issue 修复（Bug Fixing） | ~70 | 极高 | SWE-agent, OpenHands, AutoCodeRover |
| 数据集与 Benchmark | ~65 | 极高 | SWE-bench, LiveCodeBench, Terminal-Bench |
| 仓库级代码补全 | ~32 | 高 | RepoCoder, CodeRAG, RepoBench |
| 代码问答 | ~9 | 中 | SWE-QA, RepoChat |
| 代码翻译 | ~8 | 增长中 | RustRepoTrans, C2SaferRust |
| **单元测试生成** | **~7** | **增长中** | **UT-Bench 定位** |
| 任务合成（训练数据生成） | ~7 | 中 | SWE-Smith, SWE-Synth, SWE-Gym |

---

## 二、核心 Benchmark 深度分析

### 2.1 SWE-bench 家族（行业事实标准）

SWE-bench 是目前影响力最大的编码 Agent Benchmark，已形成完整的生态系统。

| 变体 | 年份 | 规模 | 特点 | 状态 |
|:---|:---|:---|:---|:---|
| **SWE-bench** | 2024 | 2,294 tasks | 原始版，12 个 Python 仓库的 GitHub Issues | 基准 |
| **SWE-bench Verified** | 2024 | 500 tasks | 人工验证子集，质量更高 | **已饱和/失效** |
| **SWE-bench+** | 2024 | - | 引入 Human Alignment Score（人工代码审查加权） | 改进 |
| **SWE-bench Multimodal** | 2025 | - | 扩展到视觉界面、GUI 交互评测 | 新方向 |
| **SWE-bench Pro** | 2025 | 1,865 tasks | Scale AI 出品，多语言，GPL 仓库防污染 | **当前推荐** |
| **SWE-bench Live** | 2025 | 动态 | 微软出品，持续从 GitHub 拉取新任务 | 动态评测 |
| **SWE-rebench** | 2025 | 21,000+ tasks | 自动化任务收集管道，适合 RL 训练 | 训练用 |
| **SWE-bench++** | 2025 | - | 扩展版本 | 新 |
| **Multi-SWE-bench** | 2025 | - | 多语言扩展版本 | 多语言 |
| **SWE-Compass** | 2025 | - | 导航式评测框架 | 工具 |

#### SWE-bench Verified 为什么失效了？

2026 年 4 月 OpenAI 正式承认 **SWE-bench Verified 已无法区分前沿编码模型**，三个根因：

1. **数据污染**：真实 GitHub Issues 出现在模型训练语料中，高分反映的是"记忆"而非"推理"。OpenAI 控制污染变量后发现，顶级模型分数显著趋同。
2. **任务粒度失配**：单点代码修改 vs 真实开发的高度迭代性（跨文件追踪、反复调试）。
3. **指标过于简单**：二元 pass/fail（通过测试 = 成功），忽略可读性、性能、副作用。

#### SWE-bench Pro 为什么是当前推荐？

| 特性 | SWE-bench Verified | SWE-bench Pro |
|:---|:---|:---|
| 语言 | Python only | **Python, JS, TS, Go, Rust** |
| 任务复杂度 | 平均 1-2 文件 | **平均 4.1 文件, 107 行** |
| 污染风险 | 高（公开 GitHub） | **低（GPL 仓库 + 私有代码）** |
| 模型区分度 | Top 6 差距仅 1.3 分 | **30-40 分的区分度** |
| 标准化脚手架 | 无 | **250 轮交互上限** |
| 开源 | Yes | **Yes**（Scale 出品） |

**2026 年 5 月 Agent 系统排行榜（SWE-bench Pro 自定义脚手架）**：

| Agent | 基础模型 | 得分 |
|:---|:---|:---|
| GPT-5.3-Codex (CLI) | GPT-5.3-Codex | 57.0% |
| Claude Code | Opus 4.5 | 55.4% |
| Auggie (Augment Code) | Opus 4.5 | 51.8% |
| Cursor | Opus 4.5 | 50.2% |

#### SWE-bench 的数据格式

每个 task 实例包含：

```json
{
  "instance_id": "django__django-16379",
  "repo": "django/django",
  "base_commit": "419b7d...abc",
  "problem_statement": "When using QuerySet.annotate()...",
  "hints_text": "Issue discussion...",
  "created_at": "2023-01-15T10:00:00Z",
  "patch": "diff --git a/django/db/models/query.py...",
  "test_patch": "diff --git a/tests/...",
  "version": "1.1",
  "FAIL_TO_PASS": ["test_annotate.test_regression_16379"],
  "PASS_TO_PASS": ["test_annotate.test_basic", ...],
  "environment_setup_commit": "5f3a2c...def"
}
```

### 2.2 LiveCodeBench（最干净的竞赛编程 benchmark）

| 属性 | 说明 |
|:---|:---|
| **来源** | LeetCode, Codeforces, AtCoder 最新题目 |
| **核心优势** | 题目在模型训练截止后才发布，**零污染** |
| **语言** | 多语言竞赛编程 |
| **局限** | 竞赛编程 ≠ 软件工程，不测多文件编辑和代码库导航 |
| **开源** | Yes |
| **推荐场景** | 测"纯推理能力"，作为模型选择参考 |

### 2.3 Terminal-Bench（最接近开发者真实工作流）

| 属性 | 说明 |
|:---|:---|
| **评测内容** | 文件编辑、git 操作、测试运行、多步调试、环境搭建 |
| **运行环境** | Docker 容器内真实终端 |
| **核心优势** | 测 Agent 在终端中"读代码→跑测试→解读输出→修复→提交"的完整循环 |
| **创建者** | Paul Gauthier（Aider 作者）参与 |
| **2026 排行榜** | GPT-5.3 Codex ~77.3%, Claude Code ~72%, Aider ~67% |
| **开源** | Yes |

### 2.4 Aider Polyglot Benchmark（纯模型代码编辑能力）

| 属性 | 说明 |
|:---|:---|
| **任务** | 133 个 Exercism 编程练习，纯代码编辑 |
| **语言** | Python, JS, TS, Go, Rust, C++, Java |
| **评测对象** | 模型本身（非 Agent 编排能力） |
| **指标** | pass_rate、syntax_errors、context_window_exhausted、$/task |
| **核心优势** | **隔离 Agent 编排，纯粹测模型代码能力** |
| **开源** | Yes |
| **推荐场景** | 选择模型时参考（同模型 Agent A vs Agent B 的公平基线） |

### 2.5 EvalPlus（HumanEval+/MBPP+）（测试扩展基准）

| 属性 | 说明 |
|:---|:---|
| **核心创新** | 将 HumanEval 的每个问题从 1 个测试扩展到 **80 个**测试用例 |
| **目的** | 消除"窄测试通过但实际逻辑有错"的假阳性 |
| **语言** | Python |
| **状态** | 已饱和（前沿模型 >95%），严重污染 |
| **对 UT-Bench 的启发** | 测试用例扩展思路可用于验证 Agent 生成的测试质量 |

### 2.6 HumanEval / MBPP（历史基准，已饱和）

| Benchmark | 样本量 | 语言 | 2026 顶级分数 | 状态 |
|:---|:---|:---|:---|:---|
| HumanEval | 164 | Python | ~96% (DeepSeek R1) | **已饱和，仅作 sanity check** |
| MBPP | 974 | Python | ~95% | **已饱和** |

### 2.7 BigCodeBench（真实库使用代码生成）

| 属性 | 说明 |
|:---|:---|
| **任务** | 1,140 个涉及真实库使用的编程任务 |
| **语言** | Python, JS, Java, Go |
| **核心优势** | 测的是"真实场景下的库 API 调用能力"，不是纯算法 |
| **发表** | ICLR 2025 Oral |
| **开源** | Yes (github.com/bigcode-project/bigcodebench) |

---

## 三、Agent 评测框架与工具

### 3.1 评测框架对比

| 框架 | 类型 | 开源 | 核心能力 | 对 UT-Bench 的价值 |
|:---|:---|:---|:---|:---|
| **SWE-agent** | Agent框架+评测 | Yes | SWE-bench 专用 Agent，内置 Docker 沙箱，原生产出 trajectory | trajectory 格式、沙箱设计参考 |
| **OpenHands** | Agent平台 | Yes | 最完整的 Agent 框架，多后端 Runtime (Docker/E2B/Local/Remote)，EventStream 事件流 | Runtime 多后端架构、EventStream 模型 |
| **mini-swe-agent** | 极简Agent | Yes | ~100 行 Python，仅 bash shell，无工具接口，score >74% on Verified | **极简主义参考：Agent 编排复杂度不是关键** |
| **DeepEval** | LLM评测框架 | Yes (Apache 2.0) | 50+ 预构建指标，Pytest 风格，支持幻觉检测、相关性、RAGAS | **指标体系参考** |
| **Langfuse** | 可观测性平台 | Yes (全开源) | Trace追踪、cost归因、prompt管理、自定义指标 | **成本追踪参考** |
| **LangSmith** | 调试平台 | No | LangChain 深度集成，trace调试，评估套件 | 仅 LangChain 生态适用 |
| **Braintrust** | 评估平台 | No | 生产trace转测试用例，CI/CD集成，GitHub Action | CI/CD 评估流程参考 |
| **Galileo** | 评测平台 | No | Luna-2 评估模型（<200ms），自动化根因分析，运行时保护 | 企业级 Agent 评测参考 |

### 3.2 重点框架详细分析

#### SWE-agent

```
Agent ←→ ACI (Agent-Computer Interface) ←→ Docker Sandbox
                                                  ├─ bash shell
                                                  ├─ 文件系统
                                                  └─ 测试运行
```

- **SWE-ReX** 沙箱抽象：`run_command()`, `read_file()`, `write_file()`, `close()`
- **Trajectory 产出**：完整的工具调用序列（action + observation），用于分析 Agent 行为
- **mini-swe-agent** 的启示：仅 ~100 行代码就能达到 74% Verified，说明 **模型推理能力 > Agent 编排复杂度**

#### OpenHands

```
Agent ←→ EventStream ←→ ActionExecutor (REST API)
                         └─ Container
                             ├─ Bash Shell
                             ├─ Browser (Playwright)
                             ├─ Jupyter
                             ├─ VSCode
                             └─ Agent Skills (插件)
```

- **Runtime 抽象**最完善：Docker / Local / Remote / E2B 四个后端
- **EventStream** 模型：Action/Observation 事件流，支持多模态
- **插件系统**：Jupyter、VSCode、Agent Skills

#### DeepEval

```python
from deepeval import evaluate
from deepeval.metrics import AnswerRelevancyMetric, HallucinationMetric

# 类 Pytest 风格的评估
def test_agent_response():
    response = my_agent.run("为这个函数生成单元测试")
    assert AnswerRelevancyMetric().measure(response) > 0.8
    assert HallucinationMetric().measure(response) < 0.1
```

- 50+ 预构建指标：幻觉检测、答案相关性、RAGAS、上下文利用率等
- 支持 LLM-as-Judge 自动化评估
- **对 UT-Bench 的启发**：可以用 LLM-as-Judge 评估 Agent 生成的测试质量

---

## 四、单元测试生成专项 Benchmark

这是 UT-Bench 的核心定位，但目前该方向的独立 Benchmark 非常少（仅 ~7 篇论文）。

### 4.1 相关论文

| 论文 | 日期 | 核心思路 |
|:---|:---|:---|
| **Rethinking Agent-Generated Tests for LLM-based SWE Agents** | 2026-02 | 重新思考 Agent 生成测试的价值，发现自动生成的测试可能误导 SWE Agent |
| **Execution-Feedback Driven Test Generation from SWE Issues** | 2025-08 | 从 SWE Issues 通过执行反馈驱动测试生成 |
| **AssertFlip: Reproducing Bugs via Inversion of LLM-Generated Tests** | 2025-07 | 反转 LLM 生成的通过测试来复现 Bug |
| **Issue2Test: Generating Reproducing Test Cases from Issue Reports** | 2025-03 | 从 Issue 报告生成复现测试用例 |
| **Agentic Bug Reproduction at Google** | 2025-02 | Google 内部的 Agent Bug 复现实践 |
| **LLMs as Continuous Learners** | 2024-11 | LLM 作为持续学习者，改进缺陷代码的复现 |

### 4.2 UT-Bench 的差异化优势

目前没有一个 Benchmark 专门且系统地评测 **"Agent 生成单元测试的能力"**，UT-Bench 的独特定位：

| 维度 | SWE-bench / SWE-agent | Aider Benchmark | UT-Bench |
|:---|:---|:---|:---|
| 评测任务 | Issue 修复 | 代码编辑（通过测试） | **单元测试生成** |
| 评测对象 | Agent 系统 | 模型本身 | **Agent vs 模型 API** |
| 评测指标 | pass/fail | pass_rate + 错误分类 | **覆盖率 + 变异测试 + 多维度** |
| 多语言 | Python（多） | 多语言 | **Python/Go/Java/C++** |
| Skill 评测 | 无 | 无 | **有（Skill 增益对比）** |
| 成本追踪 | 无 | $/task | **规划中** |
| 仓库级 | 有（完整仓库） | 无（单文件） | **规划中** |

### 4.3 可借鉴的测试生成评估指标

来自相关论文和 Benchmark：

1. **AssertFlip 指标**：生成的测试能否通过反转变异来检测 Bug
2. **Issue2Test 指标**：生成的测试能否复现 Issue 描述的问题
3. **Test Execution Rate**：生成的测试中能成功执行的比例
4. **Bug Detection Rate**：生成的测试能发现多少被注入的 Bug（变异测试）
5. **Flaky Test Rate**：生成的测试是否稳定（非 Flaky）

---

## 五、关键趋势与行业信号

### 5.1 趋势一：极简主义正在获胜

mini-swe-agent（~100 行 Python）在 SWE-bench Verified 上达到 74%，匹配甚至超过复杂框架。这说明：

> **瓶颈不是 Agent 编排复杂度，而是底层模型的推理能力和 shell 操作能力。**

对 UT-Bench 的启示：
- 不要过度设计 Agent 执行逻辑
- 重点应放在：prompt 质量、模型选择、评测指标

### 5.2 趋势二：Benchmark 污染已成核心问题

- HumanEval/MBPP/EvalPlus：严重污染，已饱和
- SWE-bench Verified：OpenAI 确认失效
- LiveCodeBench：目前最干净（动态更新）
- SWE-bench Pro：GPL 仓库防污染

对 UT-Bench 的启示：
- **自建数据集是长期壁垒**——开源 benchmark 迟早会被污染
- 自建数据集 + 私有评测 = 竞争优势
- 可以定期用 LiveCodeBench 做模型选择参考

### 5.3 趋势三：成本和延迟是盲区

**目前没有一个主流 benchmark 评测 $/task 和延迟**。但实际生产中：
- GPT-5.3 Codex 在 SWE-bench Pro 上得分 57%，但单 task 成本可能是 DeepSeek V3.2 的 10x+
- Agent 系统延迟差异巨大（5 秒 vs 5 分钟）

对 UT-Bench 的启示：
- **成本效率指标是差异化机会**：`cost_per_coverage_point`、`tokens_per_test_case`
- 在排行榜中增加"性价比"维度

### 5.4 趋势四：从评估到训练

SWE-Smith（合成训练数据）、SWE-RL（强化学习）、SWE-Gym（训练环境）——社区正在从"评测"走向"训练+评测"。

对 UT-Bench 的启示：
- 评测数据可以反向用于训练（测试生成数据 → 训练更好的测试生成模型）
- 评测框架设计应考虑数据复用

### 5.5 趋势五：Scaffolding（脚手架）影响 > 模型差异

SWE-bench Pro 上，SEAL 标准脚手架（~43%）vs 最佳 Agent 系统（~57%）差距 **14 分**，超过模型之间的差距。

> **同一个模型在不同脚手架上的得分可能差 10+ 分。**

对 UT-Bench 的启示：
- UT-Bench 的 Agent 执行环境（prompt 格式、skill 注入、workspace 准备）是评测公平性的关键
- 报告中应明确标注脚手架版本

---

## 六、对 UT-Bench 的可操作建议

### 6.1 短期（1-2 周）

1. **接入 LiveCodeBench 或 Aider Polyglot 做模型选择参考**——在选模型时先用这些公开 benchmark 跑一遍，避免选到"看起来强但实际弱"的模型
2. **学习 Aider 的错误分类指标**——引入 `syntax_errors`, `exhausted_context_windows` 等细粒度失败分类
3. **引入 `pass_rate_n` 多轮通过率**——给 Agent 最多 3 次尝试机会，报告 `pass_rate_1/2/3`

### 6.2 中期（1-2 月）

4. **仓库级数据集建设**——参考 SWE-bench 的 meta.json 格式，但任务改为"单元测试生成"
5. **成本追踪**——采集 `agent_input_tokens`/`agent_output_tokens`，计算 `$` 和 `tokens_per_coverage_point`
6. **接入 Aider 或 Claude Code 作为第二个 framework**——实现真正的框架横评

### 6.3 长期（3-6 月）

7. **自建私有数据集**——从内部代码仓库提取，建立 UT-Bench 独有的评测壁垒
8. **Trajectory 评估**——用 LLM-as-Judge 评估 Agent 的执行过程质量
9. **Leaderboard 公开**——对标 benchlm.ai，建立公开的 UT-Bench 排行榜

---

## 七、关键开源项目清单

### 7.1 必须关注（直接影响 UT-Bench 设计）

| 项目 | GitHub | 关注点 |
|:---|:---|:---|
| **SWE-bench** | github.com/SWE-bench/SWE-bench | 数据集格式、评测流程 |
| **SWE-agent** | github.com/SWE-agent/SWE-agent | Trajectory 格式、沙箱抽象 |
| **mini-swe-agent** | github.com/SWE-agent/mini-swe-agent | 极简 Agent 设计参考 |
| **OpenHands** | github.com/All-Hands-AI/OpenHands | Runtime 多后端、EventStream |
| **Aider** | github.com/paul-gauthier/aider | Benchmark 指标体系、代码编辑评测 |
| **DeepEval** | github.com/confident-ai/deepeval | 50+ 评测指标、LLM-as-Judge |
| **Langfuse** | github.com/langfuse/langfuse | 成本追踪、Trace 归因 |
| **BigCodeBench** | github.com/bigcode-project/bigcodebench | 真实库使用的代码生成 |
| **EvalPlus** | github.com/evalplus/evalplus | HumanEval+/MBPP+ 测试扩展 |

### 7.2 仓库级代码生成论文索引

| 项目 | GitHub | 关注点 |
|:---|:---|:---|
| **Awesome-Repo-Level-Code-Generation** | github.com/YerbaPage/Awesome-Repo-Level-Code-Generation | **完整的论文索引**，持续更新 |
| **RepoBench** | github.com/Leolty/repobench | 仓库级代码补全 Benchmark |
| **CodeRAG-Bench** | github.com/code-rag-bench/code-rag-bench | RAG 增强代码生成 |
| **SWE-QA-Bench** | github.com/peng-weihan/SWE-QA-Bench | 仓库级代码问答 |
| **HumanEvo** | github.com/DeepSoftwareAnalytics/HumanEvo | 时间感知仓库级代码生成 |

### 7.3 可直接集成的工具

| 项目 | 用途 | 集成方式 |
|:---|:---|:---|
| **gVisor** | 沙箱安全升级 | Docker runtime 替换 |
| **E2B SDK** | 托管沙箱后端 | SandboxRunner 后端实现 |
| **腾讯云 Cube Sandbox** | 国内部署首选沙箱 | 兼容 E2B 接口 |

---

## 八、总结

UT-Bench 处于一个**蓝海定位**——目前没有专门的 Benchmark 系统性地评测"编码 Agent 生成单元测试的能力"。但有三个关键挑战：

1. **数据集壁垒**：自建高质量仓库级数据集是长期竞争力，也是最大工作量
2. **评测公平性**：脚手架设计（prompt、workspace、skill 注入）直接影响评测结果
3. **成本意识**：业界正在觉醒，$/task 和延迟评测是差异化机会

行业正在从"能不能做"走向"做得好不好、贵不贵、快不快"，UT-Bench 的多维度评测指标体系（覆盖率 + 变异测试 + 多语言）恰好站在这个趋势上。
