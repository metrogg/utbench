# UT-Bench Agent 评测升级 —— 需求整理🎯 项目定位UT-Bench 当前是一个评测 LLM 单元测试生成能力的基准工具，核心流程：代码样本 → LLM生成测试 → 编译/测试/覆盖率/变异评测 → 报告升级目标：从"评测模型"升级为"评测 Agent"，即在现有能力基础上，支持对编码 Agent 的评测，并能够横向对比 Agent 与纯模型 API 的效果差异。--需求一：支持多种类型的 Agent 接入需求描述：支持各种类型的编码 Agent 接入评测，包括但不限于：CLI 类 Agent：Claude Code、OpenCode、CodeBuddy 等框架类 Agent：SWE-Agent、Aider、AutoCodeRover 等纯模型 API（作为基线对比）核心挑战：不同 Agent 的交互方式各异（CLI 进程、HTTP API、IDE 插件等），需要统一的接入抽象层。关键设计点：定义统一的 AgentInterface 抽象层通过适配器模式适配不同 Agent 类型复用现有 Runner 模块的依赖注入设计需求二：支持 Agent 的框架 × 模型排列组合需求描述：Agent = 框架 + 底层模型，支持任意组合的排列评测，例如：Claude Code + DeepSeek V4CodeBuddy + GLM-5.1Aider + GPT-4.1SWE-Agent + DeepSeek V4纯模型 API（无框架基线）核心挑战：不是固定两种之间的搭配，而是支持 N×M 的笛卡尔积组合评测。关键设计点：AgentSpec 数据结构包含 Framework 和 ModelName 两个维度自动生成所有合法组合（排除不兼容的搭配）每种组合独立运行评测，结果可横向对比需求三：明确输入和输出输入：数据集样本（被测代码），目标是让 Agent 为该样本生成一份优秀的单元测试代码对于单文件样本：被测源码文件对于仓库级样本：仓库路径 + 任务描述（如"为该仓库的 xxx 模块生成单元测试"）输出：核心输出：生成的单元测试代码文件过程输出（Agent 特有）：操作轨迹（工具调用序列、文件读写记录、交互轮次、Token 消耗等）需求四：Agent 在安全的沙箱中执行需求描述：Agent 需要在一个安全隔离的沙箱环境中运行，因为 Agent 会读写文件、执行命令等操作。核心挑战：尚未确定沙箱方案，需要学习业界做法。关键调研方向：Docker 容器沙箱：OpenHands、SWE-Agent 的做法E2B Sandbox：专为 AI Agent 设计的云端沙箱Firecracker microVM：轻量级虚拟机方案gVisor：用户态内核额外安全层关键设计点：为每种语言（Python/Go/Java/C++）构建专用沙箱镜像沙箱内预装编译器、测试框架、覆盖率工具超时自动销毁，资源完全隔离需求五：评测样本的设计需求描述：需要仓库级样本来充分发挥 Agent 能力，同时要复用现有单文件样本。子需求 5a：仓库级样本来源从 SWE-bench、BigCodeBench 等现有基准中选取自建仓库样本（从 GitHub 热门项目提取）覆盖四种语言：Python、Go、Java、C++子需求 5b：单文件样本复用将现有单文件样本包装成"最小仓库"结构，让 Agent 也能在其上操作单文件样本对纯模型 API 和 Agent 都适用，便于对比子需求 5c：Agent vs 纯模型 API 对比在相同样本上，对比 Agent 和纯模型 API 的效果对比维度：编译通过率、测试通过率、覆盖率、变异得分、Token 效率、时间效率等量化"Agent 比纯模型 API 提升了多少"需求六：保留项目自身特色需求描述：在升级到 Agent 评测的同时，保留 UT-Bench 的独特优势：| 现有特色 | 升级策略 | |---------|---------| | 多语言覆盖（Python/Go/Java/C++） | 保留并扩展到 Agent 场景 | | 多维度评测（编译/测试/覆盖率/变异） | 保留 + 新增 Agent 过程维度 | | 变异测试（四种语言全覆盖） | 保留，Agent 场景更有意义 | | 环境指纹与可复现性 | 保留并强化（沙箱环境指纹） | | 横向对比排行榜 | 升级为多 Agent 排行榜 |📊 需求全景图graph TB subgraph "输入层" S1[单文件样本] S2[仓库级样本] end subgraph "被测对象层" A1[纯模型API - 基线] A2[Claude Code + 模型X] A3[CodeBuddy + 模型Y] A4[其他Agent框架 + 模型Z] A5[...N×M排列组合] end subgraph "执行层" E1[统一AgentInterface] E2[Docker沙箱环境] E3[操作轨迹采集] end subgraph "输出层" O1[生成的单元测试代码] O2[Agent操作轨迹] end subgraph "评测层" R1[编译通过率] R2[测试通过率] R3[覆盖率] R4[变异得分] R5[Agent过程效率] R6[工具使用精准度] R7[Agent vs 模型 对比] end subgraph "报告层" RP1[多Agent排行榜] RP2[操作轨迹可视化] RP3[提升量对比报告] end S1 --> A1 & A2 & A3 & A4 S2 --> A2 & A3 & A4 & A5 A1 --> E1 A2 --> E1 A3 --> E1 A4 --> E1 A5 --> E1 E1 --> E2 E2 --> E3 E3 --> O1 & O2 O1 --> R1 & R2 & R3 & R4 O2 --> R5 & R6 R1 & R2 & R3 & R4 & R5 & R6 --> R7 R7 --> RP1 & RP2 & RP3---这是我现在的一个新需求 ,就是希望从原来的评测模型到评测agent可以支持agent的横评或者是agent对比单模型api,看看我们要怎么做,怎么设计比较好,可以用到哪些开源的东西 我们可以用或者是可以参照学习,具体我哪里没说好或者不清楚的可以问我

> ## UT-Bench Agent 评测升级 —— 需求整理
>
> ### 🎯 项目定位
>
> **UT-Bench** 当前是一个**评测 LLM 单元测试生成能力**的基准工具，核心流程：
>
> ```
> 代码样本 → LLM生成测试 → 编译/测试/覆盖率/变异评测 → 报告
> ```
>
> **升级目标**：从"**评测模型**"升级为"**评测 Agent**"，即在现有能力基础上，支持对编码 Agent 的评测，并能够横向对比 Agent 与纯模型 API 的效果差异。
>
> --
>
> ### 需求一：支持多种类型的 Agent 接入
>
> **需求描述**：支持各种类型的编码 Agent 接入评测，包括但不限于：
> - **CLI 类 Agent**：Claude Code、OpenCode、CodeBuddy 等
> - **框架类 Agent**：SWE-Agent、Aider、AutoCodeRover 等
> - **纯模型 API**（作为基线对比）
>
> **核心挑战**：不同 Agent 的交互方式各异（CLI 进程、HTTP API、IDE 插件等），需要统一的接入抽象层。
>
> **关键设计点**：
> - 定义统一的 `AgentInterface` 抽象层
> - 通过适配器模式适配不同 Agent 类型
> - 复用现有 Runner 模块的依赖注入设计
>
> ---
>
> ### 需求二：支持 Agent 的框架 × 模型排列组合
>
> **需求描述**：Agent = 框架 + 底层模型，支持任意组合的排列评测，例如：
> - Claude Code + DeepSeek V4
> - CodeBuddy + GLM-5.1
> - Aider + GPT-4.1
> - SWE-Agent + DeepSeek V4
> - 纯模型 API（无框架基线）
>
> **核心挑战**：不是固定两种之间的搭配，而是支持 N×M 的笛卡尔积组合评测。
>
> **关键设计点**：
> - `AgentSpec` 数据结构包含 `Framework` 和 `ModelName` 两个维度
> - 自动生成所有合法组合（排除不兼容的搭配）
> - 每种组合独立运行评测，结果可横向对比
>
> ---
>
> ### 需求三：明确输入和输出
>
> **输入**：
> - 数据集样本（被测代码），目标是让 Agent **为该样本生成一份优秀的单元测试代码**
> - 对于单文件样本：被测源码文件
> - 对于仓库级样本：仓库路径 + 任务描述（如"为该仓库的 xxx 模块生成单元测试"）
>
> **输出**：
> - **核心输出**：生成的单元测试代码文件
> - **过程输出**（Agent 特有）：操作轨迹（工具调用序列、文件读写记录、交互轮次、Token 消耗等）
>
> ---
>
> ### 需求四：Agent 在安全的沙箱中执行
>
> **需求描述**：Agent 需要在一个安全隔离的沙箱环境中运行，因为 Agent 会读写文件、执行命令等操作。
>
> **核心挑战**：尚未确定沙箱方案，需要学习业界做法。
>
> **关键调研方向**：
> - **Docker 容器沙箱**：OpenHands、SWE-Agent 的做法
> - **E2B Sandbox**：专为 AI Agent 设计的云端沙箱
> - **Firecracker microVM**：轻量级虚拟机方案
> - **gVisor**：用户态内核额外安全层
>
> **关键设计点**：
> - 为每种语言（Python/Go/Java/C++）构建专用沙箱镜像
> - 沙箱内预装编译器、测试框架、覆盖率工具
> - 超时自动销毁，资源完全隔离
>
> ---
>
> ### 需求五：评测样本的设计
>
> **需求描述**：需要仓库级样本来充分发挥 Agent 能力，同时要复用现有单文件样本。
>
> **子需求 5a：仓库级样本来源**
> - 从 SWE-bench、BigCodeBench 等现有基准中选取
> - 自建仓库样本（从 GitHub 热门项目提取）
> - 覆盖四种语言：Python、Go、Java、C++
>
> **子需求 5b：单文件样本复用**
> - 将现有单文件样本包装成"最小仓库"结构，让 Agent 也能在其上操作
> - 单文件样本对纯模型 API 和 Agent 都适用，便于对比
>
> **子需求 5c：Agent vs 纯模型 API 对比**
> - 在相同样本上，对比 Agent 和纯模型 API 的效果
> - 对比维度：编译通过率、测试通过率、覆盖率、变异得分、Token 效率、时间效率等
> - 量化"Agent 比纯模型 API 提升了多少"
>
> ---
>
> ### 需求六：保留项目自身特色
>
> **需求描述**：在升级到 Agent 评测的同时，保留 UT-Bench 的独特优势：
>
> | 现有特色 | 升级策略 |
> |---------|---------|
> | **多语言覆盖**（Python/Go/Java/C++） | 保留并扩展到 Agent 场景 |
> | **多维度评测**（编译/测试/覆盖率/变异） | 保留 + 新增 Agent 过程维度 |
> | **变异测试**（四种语言全覆盖） | 保留，Agent 场景更有意义 |
> | **环境指纹与可复现性** | 保留并强化（沙箱环境指纹） |
> | **横向对比排行榜** | 升级为多 Agent 排行榜 |
>
> ---
>
> ### 📊 需求全景图
>
> ```mermaid
> graph TB
>     subgraph "输入层"
>         S1[单文件样本]
>         S2[仓库级样本]
>     end
>
>     subgraph "被测对象层"
>         A1[纯模型API - 基线]
>         A2[Claude Code + 模型X]
>         A3[CodeBuddy + 模型Y]
>         A4[其他Agent框架 + 模型Z]
>         A5[...N×M排列组合]
>     end
>
>     subgraph "执行层"
>         E1[统一AgentInterface]
>         E2[Docker沙箱环境]
>         E3[操作轨迹采集]
>     end
>
>     subgraph "输出层"
>         O1[生成的单元测试代码]
>         O2[Agent操作轨迹]
>     end
>
>     subgraph "评测层"
>         R1[编译通过率]
>         R2[测试通过率]
>         R3[覆盖率]
>         R4[变异得分]
>         R5[Agent过程效率]
>         R6[工具使用精准度]
>         R7[Agent vs 模型 对比]
>     end
>
>     subgraph "报告层"
>         RP1[多Agent排行榜]
>         RP2[操作轨迹可视化]
>         RP3[提升量对比报告]
>     end
>
>     S1 --> A1 & A2 & A3 & A4
>     S2 --> A2 & A3 & A4 & A5
>     A1 --> E1
>     A2 --> E1
>     A3 --> E1
>     A4 --> E1
>     A5 --> E1
>     E1 --> E2
>     E2 --> E3
>     E3 --> O1 & O2
>     O1 --> R1 & R2 & R3 & R4
>     O2 --> R5 & R6
>     R1 & R2 & R3 & R4 & R5 & R6 --> R7
>     R7 --> RP1 & RP2 & RP3
> ```
>
> ---这是我现在的一个新需求 ,就是希望从原来的评测模型到评测agent可以支持agent的横评或者是agent对比单模型api,看看我们要怎么做,怎么设计比较好,可以用到哪些开源的东西 我们可以用或者是可以参照学习,具体我哪里没说好或者不清楚的可以问我

<details><summary>15 previous messages</summary>

> Stream error
>
> Reconnecting... 2/5
> timeout waiting for child process to exit
>
> Stream error
>
> Reconnecting... 3/5
> timeout waiting for child process to exit
>
> Stream error
>
> Reconnecting... 4/5
> timeout waiting for child process to exit
>
> Stream error
>
> Reconnecting... 5/5
> timeout waiting for child process to exit
>
> 我会先读一下现有 Go 项目的结构、runner/evaluator/config/dataset 这些关键点，再把需求拆成可落地的架构方案。现在先做非破坏性的代码探索，不改文件。
>
> <details><summary>Ran 3 commands</summary>
>
> - Ran `Get-ChildItem -Force`
> - Ran `Get-ChildItem -Force`
> - Ran `rg --files`
>
> </details>
>
> 我已经确认主项目在 `go-ut-bench`，现有实现是模型 Runner + 评测器 + 数据集服务的分层结构。下一步我会看现有接口和结果结构，判断“AgentInterface”应该接在哪一层，避免把模型基线和 Agent 跑偏成两套系统。
>
> <details><summary>Ran 8 commands</summary>
>
> - Ran `Get-Content internal\runner\service.go`
> - Ran `Get-Content internal\runner\api.go`
> - Ran `Get-Content internal\contracts\spec.go`
> - Ran `Get-Content internal\contracts\results.go`
> - Ran `Get-Content internal\evaluator\service.go`
> - Ran `Get-Content internal\dataset\service.go`
> - Ran `Get-ChildItem cmd\utbench -Force`
> - Ran `rg "Generate\(|Evaluate\(|RunSpec|models|langs|class|report" cmd internal -n`
>
> </details>
>
> 现有系统的关键形态很清楚：生成阶段目前强绑定 `models.yaml` 和 HTTP Chat Completions；评测阶段只消费 `GeneratedManifest`，所以升级时最稳的是把“生成测试”抽象成统一 Subject/Generator，继续产出兼容 manifest，让编译、覆盖、变异、报告尽量复用。
>
> <details><summary>Ran 7 commands, searched web 3 times</summary>
>
> - Ran `rg "flag\.|StringVar|BoolVar|run" cmd\utbench\main.go -n`
>
> - Ran `Get-Content internal\runner\models.go`
>
> - Ran `Get-Content configs\models.yaml`
>
> - Ran `Get-Content internal\reporter\service.go -TotalCount 260`
>
> - Ran `Get-Content internal\orchestrator\service.go`
>
> - Ran `Get-Content cmd\utbench\main.go -TotalCount 360`
>
> - Ran `Get-Content migrations\0001_init.sql`
>
> <details><summary>Searched web 3 times</summary>
>
> - Searched web for `SWE-agent sandbox Docker documentation agent environments`
> - Searched web for `Firecracker microVMs security isolation open source documentation`
> - Searched web for `aider documentation model settings OpenAI compatible API command line`
>
> </details>
>
> </details>
>
> 我看完了主要链路，也查了几类现成方案：OpenHands/SWE-agent 都把“Agent 可执行环境”当成沙箱/容器层，SWE-agent/mini-SWE-agent 还天然产出 trajectory；Aider/OpenCode/Claude Code 更适合作为 CLI 适配器样例；E2B/gVisor/Firecracker 是安全等级逐步升高的执行底座。这里有几个会直接影响第一版设计边界的选择，需要你定一下。
>
> User input response
>
> - 第一版你希望优先落地哪种范围？
>   - 渐进升级(Recommended)
> - Agent执行沙箱第一版采用哪个默认路线？
>   - Docker本地(Recommended)
> - 首批适配的Agent对象怎么选？
>   - CLI优先(Recommended)
>
> <details><summary>Searched web 4 times</summary>
>
> <details><summary>Searched web 4 times</summary>
>
> - Searched web for `SWE-agent documentation trajectories Docker environment`
> - Searched web for `site:gvisor.dev/docs user guide runsc docker sandbox gVisor`
> - Searched web for `Firecracker microVM official documentation lightweight virtualization security isolation`
> - Searched web for `Firecracker microVM official documentation security isolation jailer`
>
> </details>
>
> </details>
>
> Plan
>
> # UT-Bench Agent 评测升级方案
>
> ## Summary
>
> 第一版按你选定的“渐进升级 + 本地 Docker 沙箱 + CLI Agent 优先”落地：保留现有 `generate -> evaluate -> report` 主链路，把“纯模型 API 生成测试”升级为统一的 `Subject` 生成层，其中 `Subject = model_api baseline` 或 `CLI Agent framework + model`。现有评测器继续消费生成后的测试文件，因此编译、测试、覆盖率、变异测试可以最大化复用。
>
> 参考实现优先学习：OpenHands 的 Docker sandbox/custom sandbox 设计、SWE-agent/mini-SWE-agent 的 Docker 环境与 trajectory 输出、Aider/OpenCode 的 CLI + model/provider 配置方式；E2B、gVisor、Firecracker 作为后续增强路线。
>
> ## Key Changes
>
> - 新增统一生成抽象：`SubjectInterface`，替代当前 runner 对 `models.yaml + HTTP API` 的硬绑定。
>   - `model_api` adapter：复用现有 HTTP API 逻辑，作为纯模型基线。
>   - `cli_agent` adapter：运行 Aider/OpenCode/Claude Code/CodeBuddy 等 CLI，并从沙箱工作区收集新增/修改的测试文件。
>   - 后续可扩展 `http_agent`、`swe_agent`、`openhands` adapter。
>
> - 新增配置文件 `configs/agents.yaml`：
>   - `models`: 复用或引用现有模型配置。
>   - `frameworks`: 定义 Agent 框架、命令模板、环境变量、输出文件匹配规则、兼容模型列表。
>   - `subjects`: 支持显式 subject，也支持 `frameworks × models` 笛卡尔积生成。
>   - 纯模型基线表示为 `framework: model_api`。
>
> - CLI 增加新参数，同时保持兼容：
>   - `--subjects aider+deepseek-v4-flash,model_api+deepseek-v4-flash`
>   - `--agents-config ./configs/agents.yaml`
>   - `--agent-mode cli`
>   - 旧的 `--models` 继续走 `model_api` baseline。
>
> - 数据集升级：
>   - 保留现有 `self_contained` 和 `module_level`。
>   - 新增内部准备步骤 `PreparedSample`：把单文件样本包装成最小 repo，再交给 Agent。
>   - 第一版不直接引入完整 SWE-bench，只预留 `repo_level` schema；等基础链路稳定后再接 SWE-bench/BigCodeBench 子集。
>
> - Docker 沙箱第一版：
>   - 每个 subject × sample 创建独立临时 workspace。
>   - 容器只挂载该 workspace 为读写，数据集和配置按需只读挂载。
>   - 设置 timeout、CPU/memory 限制、可选 `--network none`。
>   - 生成阶段和评测阶段都记录沙箱镜像 digest、命令、退出码、耗时、资源限制。
>
> - 输出结构升级：
>   - `GeneratedCase` 增加 `subject_id`, `subject_kind`, `agent_framework`, `agent_model`, `trace_path`, `workspace_diff_path`, `sandbox_fingerprint`。
>   - 继续填充 legacy `model` 字段为 `subject_id`，保证现有 evaluator/reporter 不立即断裂。
>   - Agent 过程输出保存为 JSONL：命令、stdout/stderr 摘要、文件变更、token/费用字段、交互轮次。
>
> - 报告升级：
>   - 排行榜维度从 `model` 扩展为 `subject`。
>   - 新增 Agent vs Model API 对比表：同一模型、同一样本下，对比 `agent_framework+model` 相对 `model_api+model` 的覆盖率、变异得分、耗时、token/成本增量。
>   - 新增过程效率指标：生成耗时、命令轮次、文件修改数、失败阶段、trace 可查看路径。
>
> ## Implementation Plan
>
> 1. 抽象 runner：
>    - 在 `internal/runner` 内拆出 `Generator` 接口。
>    - 当前 API 调用迁移为 `modelAPIAdapter`。
>    - 新增 `cliAgentAdapter`，通过命令模板启动 CLI Agent。
>
> 2. 增加 subject 配置解析：
>    - 新建 `internal/agentconfig` 或 `internal/subject`。
>    - 解析 `agents.yaml`，生成合法 subject 列表。
>    - 做兼容过滤：框架不支持的模型直接跳过并写入 manifest metadata。
>
> 3. 增加样本工作区准备：
>    - `self_contained` 单文件包装成最小 repo。
>    - `module_level` 复制现有 workspace。
>    - 输出固定任务说明：要求只新增/修改测试文件，不改被测源码。
>
> 4. 增加 Docker sandbox runner：
>    - 统一执行 CLI Agent 命令。
>    - 捕获 stdout/stderr、退出码、耗时。
>    - 执行前后扫描 workspace diff，按 `output_globs` 找测试文件。
>    - 超时后强制销毁容器。
>
> 5. 扩展 manifest/result/report：
>    - manifest 记录 subject 和 trace。
>    - evaluator 继续按生成测试文件评测。
>    - reporter 聚合维度切到 subject，同时保留旧字段展示。
>
> 6. 首批适配器：
>    - `model_api`: 现有模型基线。
>    - `aider`: 使用非交互命令/输入 prompt 模式。
>    - `opencode`: 使用 `--model` 或配置文件指定 provider/model。
>    - `claude_code`/`codebuddy`: 先按通用 CLI command-template 适配，具体命令通过配置注入。
>
> ## Test Plan
>
> - 单元测试：
>   - subject 笛卡尔积生成、兼容模型过滤、legacy `--models` 行为。
>   - `agents.yaml` 解析错误、缺失命令、缺失模型配置。
>   - workspace diff 能正确识别新增测试文件。
>
> - 集成测试：
>   - `model_api` dry-run 生成结果与旧 runner 兼容。
>   - fake CLI Agent 在 Docker 中写入测试文件，manifest 记录 trace。
>   - timeout/非零退出码/无测试文件输出都能变成明确失败原因。
>
> - 回归测试：
>   - 现有 `go test ./cmd/utbench ./internal/evaluator ./internal/runner ./internal/reporter`。
>   - 用 `--dry-run --models deepseek-v4-flash --class self_contained --max-samples 1` 验证旧链路不破。
>   - 用 fake agent 跑 full pipeline，确认 evaluator/report 不依赖旧模型名。
>
> ## Assumptions
>
> - 第一版只实现本地 Docker 沙箱，不接 E2B/Firecracker。
> - 第一版优先 CLI Agent；SWE-agent/OpenHands 作为第二阶段 adapter。
> - 纯模型 API 必须作为 `model_api+model` subject 进入同一排行榜。
> - Agent 只能在样本副本内写文件，不能修改原始 dataset。
> - 仓库级真实样本先预留 schema，不在第一版强行接入完整 SWE-bench。
>
> ## References
>
> - OpenHands Docker sandbox/custom sandbox: [docs.openhands.dev](https://docs.openhands.dev/openhands/usage/runtimes/docker)
> - SWE-agent Docker environments: [swe-agent.com](https://swe-agent.com/latest/config/environments/)
> - mini-SWE-agent trajectory output: [mini-swe-agent.com](https://mini-swe-agent.com/latest/usage/output_files/)
> - Aider CLI/model options: [aider.chat](https://aider.chat/docs/config/options.html)
> - OpenCode provider/model config: [opencode.ai](https://opencode.ai/docs/models/)
> - E2B sandbox: [e2b.dev](https://www.e2b.dev/docs)
> - gVisor Docker runtime: [gvisor.dev](https://gvisor.dev/docs/user_guide/quick_start/docker/)
> - Firecracker microVM: [firecracker-microvm.github.io](https://firecracker-microvm.github.io/)
>
> Plan implementation
>
> Status: completed
> # UT-Bench Agent 评测升级方案
>
> ## Summary
>
> 第一版按你选定的“渐进升级 + 本地 Docker 沙箱 + CLI Agent 优先”落地：保留现有 `generate -> evaluate -> report` 主链路，把“纯模型 API 生成测试”升级为统一的 `Subject` 生成层，其中 `Subject = model_api baseline` 或 `CLI Agent framework + model`。现有评测器继续消费生成后的测试文件，因此编译、测试、覆盖率、变异测试可以最大化复用。
>
> 参考实现优先学习：OpenHands 的 Docker sandbox/custom sandbox 设计、SWE-agent/mini-SWE-agent 的 Docker 环境与 trajectory 输出、Aider/OpenCode 的 CLI + model/provider 配置方式；E2B、gVisor、Firecracker 作为后续增强路线。
>
> ## Key Changes
>
> - 新增统一生成抽象：`SubjectInterface`，替代当前 runner 对 `models.yaml + HTTP API` 的硬绑定。
>   - `model_api` adapter：复用现有 HTTP API 逻辑，作为纯模型基线。
>   - `cli_agent` adapter：运行 Aider/OpenCode/Claude Code/CodeBuddy 等 CLI，并从沙箱工作区收集新增/修改的测试文件。
>   - 后续可扩展 `http_agent`、`swe_agent`、`openhands` adapter。
>
> - 新增配置文件 `configs/agents.yaml`：
>   - `models`: 复用或引用现有模型配置。
>   - `frameworks`: 定义 Agent 框架、命令模板、环境变量、输出文件匹配规则、兼容模型列表。
>   - `subjects`: 支持显式 subject，也支持 `frameworks × models` 笛卡尔积生成。
>   - 纯模型基线表示为 `framework: model_api`。
>
> - CLI 增加新参数，同时保持兼容：
>   - `--subjects aider+deepseek-v4-flash,model_api+deepseek-v4-flash`
>   - `--agents-config ./configs/agents.yaml`
>   - `--agent-mode cli`
>   - 旧的 `--models` 继续走 `model_api` baseline。
>
> - 数据集升级：
>   - 保留现有 `self_contained` 和 `module_level`。
>   - 新增内部准备步骤 `PreparedSample`：把单文件样本包装成最小 repo，再交给 Agent。
>   - 第一版不直接引入完整 SWE-bench，只预留 `repo_level` schema；等基础链路稳定后再接 SWE-bench/BigCodeBench 子集。
>
> - Docker 沙箱第一版：
>   - 每个 subject × sample 创建独立临时 workspace。
>   - 容器只挂载该 workspace 为读写，数据集和配置按需只读挂载。
>   - 设置 timeout、CPU/memory 限制、可选 `--network none`。
>   - 生成阶段和评测阶段都记录沙箱镜像 digest、命令、退出码、耗时、资源限制。
>
> - 输出结构升级：
>   - `GeneratedCase` 增加 `subject_id`, `subject_kind`, `agent_framework`, `agent_model`, `trace_path`, `workspace_diff_path`, `sandbox_fingerprint`。
>   - 继续填充 legacy `model` 字段为 `subject_id`，保证现有 evaluator/reporter 不立即断裂。
>   - Agent 过程输出保存为 JSONL：命令、stdout/stderr 摘要、文件变更、token/费用字段、交互轮次。
>
> - 报告升级：
>   - 排行榜维度从 `model` 扩展为 `subject`。
>   - 新增 Agent vs Model API 对比表：同一模型、同一样本下，对比 `agent_framework+model` 相对 `model_api+model` 的覆盖率、变异得分、耗时、token/成本增量。
>   - 新增过程效率指标：生成耗时、命令轮次、文件修改数、失败阶段、trace 可查看路径。
>
> ## Implementation Plan
>
> 1. 抽象 runner：
>    - 在 `internal/runner` 内拆出 `Generator` 接口。
>    - 当前 API 调用迁移为 `modelAPIAdapter`。
>    - 新增 `cliAgentAdapter`，通过命令模板启动 CLI Agent。
>
> 2. 增加 subject 配置解析：
>    - 新建 `internal/agentconfig` 或 `internal/subject`。
>    - 解析 `agents.yaml`，生成合法 subject 列表。
>    - 做兼容过滤：框架不支持的模型直接跳过并写入 manifest metadata。
>
> 3. 增加样本工作区准备：
>    - `self_contained` 单文件包装成最小 repo。
>    - `module_level` 复制现有 workspace。
>    - 输出固定任务说明：要求只新增/修改测试文件，不改被测源码。
>
> 4. 增加 Docker sandbox runner：
>    - 统一执行 CLI Agent 命令。
>    - 捕获 stdout/stderr、退出码、耗时。
>    - 执行前后扫描 workspace diff，按 `output_globs` 找测试文件。
>    - 超时后强制销毁容器。
>
> 5. 扩展 manifest/result/report：
>    - manifest 记录 subject 和 trace。
>    - evaluator 继续按生成测试文件评测。
>    - reporter 聚合维度切到 subject，同时保留旧字段展示。
>
> 6. 首批适配器：
>    - `model_api`: 现有模型基线。
>    - `aider`: 使用非交互命令/输入 prompt 模式。
>    - `opencode`: 使用 `--model` 或配置文件指定 provider/model。
>    - `claude_code`/`codebuddy`: 先按通用 CLI command-template 适配，具体命令通过配置注入。
>
> ## Test Plan
>
> - 单元测试：
>   - subject 笛卡尔积生成、兼容模型过滤、legacy `--models` 行为。
>   - `agents.yaml` 解析错误、缺失命令、缺失模型配置。
>   - workspace diff 能正确识别新增测试文件。
>
> - 集成测试：
>   - `model_api` dry-run 生成结果与旧 runner 兼容。
>   - fake CLI Agent 在 Docker 中写入测试文件，manifest 记录 trace。
>   - timeout/非零退出码/无测试文件输出都能变成明确失败原因。
>
> - 回归测试：
>   - 现有 `go test ./cmd/utbench ./internal/evaluator ./internal/runner ./internal/reporter`。
>   - 用 `--dry-run --models deepseek-v4-flash --class self_contained --max-samples 1` 验证旧链路不破。
>   - 用 fake agent 跑 full pipeline，确认 evaluator/report 不依赖旧模型名。
>
> ## Assumptions
>
> - 第一版只实现本地 Docker 沙箱，不接 E2B/Firecracker。
> - 第一版优先 CLI Agent；SWE-agent/OpenHands 作为第二阶段 adapter。
> - 纯模型 API 必须作为 `model_api+model` subject 进入同一排行榜。
> - Agent 只能在样本副本内写文件，不能修改原始 dataset。
> - 仓库级真实样本先预留 schema，不在第一版强行接入完整 SWE-bench。
>
> ## References
>
> - OpenHands Docker sandbox/custom sandbox: [docs.openhands.dev](https://docs.openhands.dev/openhands/usage/runtimes/docker)
> - SWE-agent Docker environments: [swe-agent.com](https://swe-agent.com/latest/config/environments/)
> - mini-SWE-agent trajectory output: [mini-swe-agent.com](https://mini-swe-agent.com/latest/usage/output_files/)
> - Aider CLI/model options: [aider.chat](https://aider.chat/docs/config/options.html)
> - OpenCode provider/model config: [opencode.ai](https://opencode.ai/docs/models/)
> - E2B sandbox: [e2b.dev](https://www.e2b.dev/docs)
> - gVisor Docker runtime: [gvisor.dev](https://gvisor.dev/docs/user_guide/quick_start/docker/)
> - Firecracker microVM: [firecracker-microvm.github.io](https://firecracker-microvm.github.io/)
</details>

> 这个agent啊比如Claude code +deepseekv4后面还可以加一个skill这样,我们还能评价skill的能力,当然我们可以选择要或者不要,但是预留这个

<details><summary>4 previous messages</summary>

> 这个补充很关键：被测对象不能只建模成 `framework + model`，应该预留成 `framework + model + skill/profile`，其中 skill 是可选第三维。这样以后既能横评 Agent，也能做“同一 Agent 加不加某个 skill 的增益评测”。我需要把 skill 的定义边界锁一下。
>
> User input response
>
> - 这里的 skill 第一版应该按哪种抽象预留？
>   - 能力包(Recommended)
> - skill 能力评测第一版按什么方式呈现？
>   - A/B增益(Recommended)
>
> Plan
>
> # UT-Bench Agent + Skill 评测升级方案
>
> ## Summary
>
> 被测对象升级为三元组：`Subject = framework + model + optional skill`。第一版仍按“渐进升级 + 本地 Docker 沙箱 + CLI Agent 优先”落地，同时预留并实现 skill 维度的基础数据结构和 A/B 增益报告。
>
> 核心对比方式：
>
> - 纯模型基线：`model_api + deepseek-v4 + no_skill`
> - Agent 基线：`claude_code + deepseek-v4 + no_skill`
> - Agent + Skill：`claude_code + deepseek-v4 + unit_test_skill`
>
> 这样可以回答两类问题：Agent 是否比纯模型 API 强；同一个 Agent 加 skill 后提升多少。
>
> ## Key Changes
>
> - 新增统一被测对象结构 `SubjectSpec`：
>   - `subject_id`: 稳定唯一名，如 `claude_code__deepseek-v4__unit_test_skill`
>   - `kind`: `model_api` / `cli_agent` / 后续 `http_agent`、`swe_agent`
>   - `framework`: `model_api`、`claude_code`、`aider`、`opencode`、`codebuddy`
>   - `model`: `deepseek-v4`、`glm-5` 等
>   - `skill`: 可选，缺省为 `no_skill`
>   - `labels/tags`: 用于报告筛选
>
> - 新增 `SkillSpec` 抽象，定义为“能力包”，而不是绑定某个 Agent 的原生 skill：
>   - `name`, `version`, `description`
>   - `instruction_path`: skill 说明文档或 prompt 片段
>   - `files`: 可挂载到沙箱的脚本、模板、示例
>   - `inject_mode`: `prompt_append` / `workspace_mount` / `agent_native`
>   - `compatible_frameworks`
>   - `compatible_languages`
>   - `enabled`
>
> - `agents.yaml` 扩展为三维组合：
>   - `frameworks × models × skills`
>   - `skills` 可为空，系统自动生成 `no_skill` baseline。
>   - 支持排除不兼容组合，例如某 skill 只适合 Python，或只适合 Claude Code/Aider。
>
> - 生成阶段统一走 `SubjectInterface`：
>   - `model_api` adapter 复用现有 HTTP 模型调用。
>   - `cli_agent` adapter 在 Docker 沙箱里运行命令模板。
>   - skill 注入由 adapter 统一处理，不让 evaluator 感知不同 Agent 的细节。
>
> - 输出结构扩展：
>   - `GeneratedCase` 增加 `subject_id`, `subject_kind`, `agent_framework`, `agent_model`, `skill_name`, `skill_version`, `trace_path`, `workspace_diff_path`, `sandbox_fingerprint`。
>   - 旧 `model` 字段继续填 `subject_id`，保证 evaluator/reporter 兼容。
>
> - 报告新增两组对比：
>   - Agent vs Model API：同模型、同样本比较 `framework+model+no_skill` 相对 `model_api+model+no_skill`。
>   - Skill A/B 增益：同 framework、model、sample 比较 `skill` 相对 `no_skill`。
>   - 指标包括编译通过率、测试通过率、覆盖率、变异得分、耗时、token/成本、操作轮次。
>
> ## Implementation Plan
>
> 1. 抽象 runner：
>    - 在 `internal/runner` 拆出 `SubjectGenerator` 接口。
>    - 当前模型 API 改造成 `modelAPIAdapter`。
>    - 新增 `cliAgentAdapter`，负责命令执行、trace 捕获、测试文件收集。
>
> 2. 增加 subject/skill 配置解析：
>    - 新建 `internal/subject` 或 `internal/agentconfig`。
>    - 解析 `agents.yaml` 中的 models、frameworks、skills。
>    - 生成合法三元组组合，并自动加入 `no_skill` baseline。
>    - checkpoint key 从 `model|language|sample` 升级为 `subject_id|language|sample`。
>
> 3. 增加 skill 注入机制：
>    - `prompt_append`: 把 skill instruction 附加到任务说明。
>    - `workspace_mount`: 把 skill 文件复制/挂载到样本 workspace。
>    - `agent_native`: 只记录配置，具体注入由对应 Agent adapter 实现。
>    - 第一版优先实现 `prompt_append` 和 `workspace_mount`。
>
> 4. 增加 Docker 沙箱执行：
>    - 每个 `subject × sample` 独立 workspace。
>    - 只允许写 workspace，不改原始 dataset。
>    - 记录容器镜像、命令、退出码、耗时、stdout/stderr 摘要和文件 diff。
>    - 超时强制销毁容器。
>
> 5. 报告与数据库升级：
>    - reporter 聚合维度从 `model` 扩展为 `subject`。
>    - 增加 skill 维度表和 A/B uplift 表。
>    - SQLite 结果表增加 subject/framework/model/skill 字段；旧 model 字段保留兼容。
>
> ## Test Plan
>
> - 单元测试：
>   - `framework × model × skill` 组合生成。
>   - skill 兼容性过滤。
>   - `no_skill` baseline 自动生成。
>   - checkpoint key 包含 skill，避免复用错误结果。
>
> - 集成测试：
>   - fake CLI Agent + no_skill 能生成测试文件并通过 manifest。
>   - fake CLI Agent + skill 能接收到注入文件或 prompt。
>   - 同一 sample 下生成 no_skill 与 skill 结果，report 能计算 uplift。
>
> - 回归测试：
>   - 旧 `--models` 仍等价于 `model_api + model + no_skill`。
>   - 现有 evaluator 对扩展 manifest 继续可用。
>   - 现有模型排行榜不丢失，只是展示名变为 subject。
>
> ## Assumptions
>
> - 第一版 skill 是通用“能力包”，不绑定 Claude Code/Codex 原生 skill。
> - Skill 评测第一版采用 A/B 增益：同一 framework+model 对比 `no_skill` 和指定 skill。
> - 第一版只实现本地 Docker 沙箱和 CLI Agent。
> - 仓库级真实样本先预留 schema，优先复用现有单文件/module_level 样本包装成最小 repo。
> - 纯模型 API 必须始终作为 baseline 参与同一套结果与报告。
>
> Plan implementation
>
> Status: completed
> # UT-Bench Agent + Skill 评测升级方案
>
> ## Summary
>
> 被测对象升级为三元组：`Subject = framework + model + optional skill`。第一版仍按“渐进升级 + 本地 Docker 沙箱 + CLI Agent 优先”落地，同时预留并实现 skill 维度的基础数据结构和 A/B 增益报告。
>
> 核心对比方式：
>
> - 纯模型基线：`model_api + deepseek-v4 + no_skill`
> - Agent 基线：`claude_code + deepseek-v4 + no_skill`
> - Agent + Skill：`claude_code + deepseek-v4 + unit_test_skill`
>
> 这样可以回答两类问题：Agent 是否比纯模型 API 强；同一个 Agent 加 skill 后提升多少。
>
> ## Key Changes
>
> - 新增统一被测对象结构 `SubjectSpec`：
>   - `subject_id`: 稳定唯一名，如 `claude_code__deepseek-v4__unit_test_skill`
>   - `kind`: `model_api` / `cli_agent` / 后续 `http_agent`、`swe_agent`
>   - `framework`: `model_api`、`claude_code`、`aider`、`opencode`、`codebuddy`
>   - `model`: `deepseek-v4`、`glm-5` 等
>   - `skill`: 可选，缺省为 `no_skill`
>   - `labels/tags`: 用于报告筛选
>
> - 新增 `SkillSpec` 抽象，定义为“能力包”，而不是绑定某个 Agent 的原生 skill：
>   - `name`, `version`, `description`
>   - `instruction_path`: skill 说明文档或 prompt 片段
>   - `files`: 可挂载到沙箱的脚本、模板、示例
>   - `inject_mode`: `prompt_append` / `workspace_mount` / `agent_native`
>   - `compatible_frameworks`
>   - `compatible_languages`
>   - `enabled`
>
> - `agents.yaml` 扩展为三维组合：
>   - `frameworks × models × skills`
>   - `skills` 可为空，系统自动生成 `no_skill` baseline。
>   - 支持排除不兼容组合，例如某 skill 只适合 Python，或只适合 Claude Code/Aider。
>
> - 生成阶段统一走 `SubjectInterface`：
>   - `model_api` adapter 复用现有 HTTP 模型调用。
>   - `cli_agent` adapter 在 Docker 沙箱里运行命令模板。
>   - skill 注入由 adapter 统一处理，不让 evaluator 感知不同 Agent 的细节。
>
> - 输出结构扩展：
>   - `GeneratedCase` 增加 `subject_id`, `subject_kind`, `agent_framework`, `agent_model`, `skill_name`, `skill_version`, `trace_path`, `workspace_diff_path`, `sandbox_fingerprint`。
>   - 旧 `model` 字段继续填 `subject_id`，保证 evaluator/reporter 兼容。
>
> - 报告新增两组对比：
>   - Agent vs Model API：同模型、同样本比较 `framework+model+no_skill` 相对 `model_api+model+no_skill`。
>   - Skill A/B 增益：同 framework、model、sample 比较 `skill` 相对 `no_skill`。
>   - 指标包括编译通过率、测试通过率、覆盖率、变异得分、耗时、token/成本、操作轮次。
>
> ## Implementation Plan
>
> 1. 抽象 runner：
>    - 在 `internal/runner` 拆出 `SubjectGenerator` 接口。
>    - 当前模型 API 改造成 `modelAPIAdapter`。
>    - 新增 `cliAgentAdapter`，负责命令执行、trace 捕获、测试文件收集。
>
> 2. 增加 subject/skill 配置解析：
>    - 新建 `internal/subject` 或 `internal/agentconfig`。
>    - 解析 `agents.yaml` 中的 models、frameworks、skills。
>    - 生成合法三元组组合，并自动加入 `no_skill` baseline。
>    - checkpoint key 从 `model|language|sample` 升级为 `subject_id|language|sample`。
>
> 3. 增加 skill 注入机制：
>    - `prompt_append`: 把 skill instruction 附加到任务说明。
>    - `workspace_mount`: 把 skill 文件复制/挂载到样本 workspace。
>    - `agent_native`: 只记录配置，具体注入由对应 Agent adapter 实现。
>    - 第一版优先实现 `prompt_append` 和 `workspace_mount`。
>
> 4. 增加 Docker 沙箱执行：
>    - 每个 `subject × sample` 独立 workspace。
>    - 只允许写 workspace，不改原始 dataset。
>    - 记录容器镜像、命令、退出码、耗时、stdout/stderr 摘要和文件 diff。
>    - 超时强制销毁容器。
>
> 5. 报告与数据库升级：
>    - reporter 聚合维度从 `model` 扩展为 `subject`。
>    - 增加 skill 维度表和 A/B uplift 表。
>    - SQLite 结果表增加 subject/framework/model/skill 字段；旧 model 字段保留兼容。
>
> ## Test Plan
>
> - 单元测试：
>   - `framework × model × skill` 组合生成。
>   - skill 兼容性过滤。
>   - `no_skill` baseline 自动生成。
>   - checkpoint key 包含 skill，避免复用错误结果。
>
> - 集成测试：
>   - fake CLI Agent + no_skill 能生成测试文件并通过 manifest。
>   - fake CLI Agent + skill 能接收到注入文件或 prompt。
>   - 同一 sample 下生成 no_skill 与 skill 结果，report 能计算 uplift。
>
> - 回归测试：
>   - 旧 `--models` 仍等价于 `model_api + model + no_skill`。
>   - 现有 evaluator 对扩展 manifest 继续可用。
>   - 现有模型排行榜不丢失，只是展示名变为 subject。
>
> ## Assumptions
>
> - 第一版 skill 是通用“能力包”，不绑定 Claude Code/Codex 原生 skill。
> - Skill 评测第一版采用 A/B 增益：同一 framework+model 对比 `no_skill` 和指定 skill。
> - 第一版只实现本地 Docker 沙箱和 CLI Agent。
> - 仓库级真实样本先预留 schema，优先复用现有单文件/module_level 样本包装成最小 repo。
> - 纯模型 API 必须始终作为 baseline 参与同一套结果与报告。
</details>

> PLEASE IMPLEMENT THIS PLAN:
> # UT-Bench Agent + Skill 评测升级方案
>
> ## Summary
>
> 被测对象升级为三元组：`Subject = framework + model + optional skill`。第一版仍按“渐进升级 + 本地 Docker 沙箱 + CLI Agent 优先”落地，同时预留并实现 skill 维度的基础数据结构和 A/B 增益报告。
>
> 核心对比方式：
>
> - 纯模型基线：`model_api + deepseek-v4 + no_skill`
> - Agent 基线：`claude_code + deepseek-v4 + no_skill`
> - Agent + Skill：`claude_code + deepseek-v4 + unit_test_skill`
>
> 这样可以回答两类问题：Agent 是否比纯模型 API 强；同一个 Agent 加 skill 后提升多少。
>
> ## Key Changes
>
> - 新增统一被测对象结构 `SubjectSpec`：
>   - `subject_id`: 稳定唯一名，如 `claude_code__deepseek-v4__unit_test_skill`
>   - `kind`: `model_api` / `cli_agent` / 后续 `http_agent`、`swe_agent`
>   - `framework`: `model_api`、`claude_code`、`aider`、`opencode`、`codebuddy`
>   - `model`: `deepseek-v4`、`glm-5` 等
>   - `skill`: 可选，缺省为 `no_skill`
>   - `labels/tags`: 用于报告筛选
>
> - 新增 `SkillSpec` 抽象，定义为“能力包”，而不是绑定某个 Agent 的原生 skill：
>   - `name`, `version`, `description`
>   - `instruction_path`: skill 说明文档或 prompt 片段
>   - `files`: 可挂载到沙箱的脚本、模板、示例
>   - `inject_mode`: `prompt_append` / `workspace_mount` / `agent_native`
>   - `compatible_frameworks`
>   - `compatible_languages`
>   - `enabled`
>
> - `agents.yaml` 扩展为三维组合：
>   - `frameworks × models × skills`
>   - `skills` 可为空，系统自动生成 `no_skill` baseline。
>   - 支持排除不兼容组合，例如某 skill 只适合 Python，或只适合 Claude Code/Aider。
>
> - 生成阶段统一走 `SubjectInterface`：
>   - `model_api` adapter 复用现有 HTTP 模型调用。
>   - `cli_agent` adapter 在 Docker 沙箱里运行命令模板。
>   - skill 注入由 adapter 统一处理，不让 evaluator 感知不同 Agent 的细节。
>
> - 输出结构扩展：
>   - `GeneratedCase` 增加 `subject_id`, `subject_kind`, `agent_framework`, `agent_model`, `skill_name`, `skill_version`, `trace_path`, `workspace_diff_path`, `sandbox_fingerprint`。
>   - 旧 `model` 字段继续填 `subject_id`，保证 evaluator/reporter 兼容。
>
> - 报告新增两组对比：
>   - Agent vs Model API：同模型、同样本比较 `framework+model+no_skill` 相对 `model_api+model+no_skill`。
>   - Skill A/B 增益：同 framework、model、sample 比较 `skill` 相对 `no_skill`。
>   - 指标包括编译通过率、测试通过率、覆盖率、变异得分、耗时、token/成本、操作轮次。
>
> ## Implementation Plan
>
> 1. 抽象 runner：
>    - 在 `internal/runner` 拆出 `SubjectGenerator` 接口。
>    - 当前模型 API 改造成 `modelAPIAdapter`。
>    - 新增 `cliAgentAdapter`，负责命令执行、trace 捕获、测试文件收集。
>
> 2. 增加 subject/skill 配置解析：
>    - 新建 `internal/subject` 或 `internal/agentconfig`。
>    - 解析 `agents.yaml` 中的 models、frameworks、skills。
>    - 生成合法三元组组合，并自动加入 `no_skill` baseline。
>    - checkpoint key 从 `model|language|sample` 升级为 `subject_id|language|sample`。
>
> 3. 增加 skill 注入机制：
>    - `prompt_append`: 把 skill instruction 附加到任务说明。
>    - `workspace_mount`: 把 skill 文件复制/挂载到样本 workspace。
>    - `agent_native`: 只记录配置，具体注入由对应 Agent adapter 实现。
>    - 第一版优先实现 `prompt_append` 和 `workspace_mount`。
>
> 4. 增加 Docker 沙箱执行：
>    - 每个 `subject × sample` 独立 workspace。
>    - 只允许写 workspace，不改原始 dataset。
>    - 记录容器镜像、命令、退出码、耗时、stdout/stderr 摘要和文件 diff。
>    - 超时强制销毁容器。
>
> 5. 报告与数据库升级：
>    - reporter 聚合维度从 `model` 扩展为 `subject`。
>    - 增加 skill 维度表和 A/B uplift 表。
>    - SQLite 结果表增加 subject/framework/model/skill 字段；旧 model 字段保留兼容。
>
> ## Test Plan
>
> - 单元测试：
>   - `framework × model × skill` 组合生成。
>   - skill 兼容性过滤。
>   - `no_skill` baseline 自动生成。
>   - checkpoint key 包含 skill，避免复用错误结果。
>
> - 集成测试：
>   - fake CLI Agent + no_skill 能生成测试文件并通过 manifest。
>   - fake CLI Agent + skill 能接收到注入文件或 prompt。
>   - 同一 sample 下生成 no_skill 与 skill 结果，report 能计算 uplift。
>
> - 回归测试：
>   - 旧 `--models` 仍等价于 `model_api + model + no_skill`。
>   - 现有 evaluator 对扩展 manifest 继续可用。
>   - 现有模型排行榜不丢失，只是展示名变为 subject。
>
> ## Assumptions
>
> - 第一版 skill 是通用“能力包”，不绑定 Claude Code/Codex 原生 skill。
> - Skill 评测第一版采用 A/B 增益：同一 framework+model 对比 `no_skill` 和指定 skill。
> - 第一版只实现本地 Docker 沙箱和 CLI Agent。
> - 仓库级真实样本先预留 schema，优先复用现有单文件/module_level 样本包装成最小 repo。
> - 纯模型 API 必须始终作为 baseline 参与同一套结果与报告。

<details><summary>32 previous messages</summary>

> 我会按第一版可落地范围实现：先把 subject/skill 配置、runner 生成路径、manifest/result/report 字段打通，再加 fake CLI agent 可测试的沙箱执行。开始前先确认当前工作树和相关存储/报告代码，避免覆盖已有改动。
>
> <details><summary>Ran 4 commands</summary>
>
> - Ran `git status --short`
> - Ran `rg "sample_results|GeneratedCase|EvaluationResult|Ingest" internal\store -n`
> - Ran `rg "buildTopModels|buildDimensions|ModelDim|ModelRank|ReportPayload|summaryJSON" internal\reporter internal\contracts -n`
> - Ran `rg "GeneratedCase\{" -n internal cmd`
>
> </details>
>
> 工作树里只有仓库外的未提交文件，`go-ut-bench` 本身是干净的。我会集中改 Go 项目的 contracts、runner、reporter、store 和 CLI 配置，不碰外层文件。
>
> <details><summary>Ran 8 commands</summary>
>
> - Ran `Get-Content internal\store\sqlite.go -TotalCount 360`
> - Ran `Get-Content internal\store\sqlite.go | Select-Object -Skip 1180 -First 260`
> - Ran `Get-Content internal\store\sqlite.go | Select-Object -Skip 360 -First 240`
> - Ran `Get-Content internal\reporter\service.go | Select-Object -Skip 420 -First 500`
> - Ran `Get-Content internal\runner\checkpoint_test.go`
> - Ran `Get-Content internal\runner\api_test.go`
> - Ran `Get-Content internal\reporter\service_test.go`
> - Ran `Get-Content internal\store\sqlite_test.go`
>
> </details>
>
> 我会把实现拆成三块提交到代码：配置/数据结构、runner 的 subject adapter、report/store 的 subject/skill 字段与增益统计。CLI Agent 的 Docker 执行会支持 `sandbox_mode: docker`，同时保留 `local` 仅用于 fake agent 测试，默认仍是 Docker。
>
> <details><summary>Ran 4 commands</summary>
>
> - Ran `Get-Content internal\runner\prompt.go | Select-Object -First 120`
> - Ran `rg "LoadPromptCatalog|PromptMode|buildPrompt" internal\runner\prompt.go -n`
> - Ran `Get-Content go.mod`
> - Ran `rg "runCommandWithProcessGroupKill|copyDir|copyFile|filepath.WalkDir" internal\evaluator internal\runner -n`
>
> </details>
>
> 先改公共契约结构：`RunSpec` 增加 subject/agent 配置入口，`GeneratedCase` 和 `EvaluationResult` 增加 subject/framework/model/skill/trace 字段，报告 payload 增加 Agent-vs-API 和 Skill uplift 两张结构化表。
>
> <details><summary>Edited 2 files, ran 1 command</summary>
>
> - Ran `Get-Content internal\contracts\spec.go | Select-Object -First 45`
> - Wrote `./go-ut-bench/internal/contracts/spec.go` (+27 -0)
> - Wrote `./go-ut-bench/internal/contracts/results.go` (+16 -7)
> - Wrote `./go-ut-bench/internal/contracts/results.go` (+15 -6)
> - Wrote `./go-ut-bench/internal/contracts/results.go` (+61 -12)
>
> </details>
>
> 接下来新增 `internal/agentconfig`：它只负责把 `agents.yaml` 解析成确定的 subject 列表，并保证 `no_skill` 和 `model_api` baseline 自动存在。runner 会消费这个解析结果。
>
> <details><summary>Created 1 file, edited 1 file, ran 4 commands</summary>
>
> - Wrote `./go-ut-bench/internal/agentconfig/config.go` (+442 -0)
> - Wrote `./go-ut-bench/internal/runner/service.go` (+15 -9)
> - Ran `Get-Content internal\runner\service.go | Select-Object -Skip 70 -First 120`
> - Ran `Get-Content internal\runner\service.go | Select-Object -Skip 190 -First 90`
> - Wrote `./go-ut-bench/internal/runner/service.go` (+13 -9)
> - Wrote `./go-ut-bench/internal/runner/service.go` (+43 -43)
> - Wrote `./go-ut-bench/internal/runner/service.go` (+94 -67)
> - Ran `Get-Content internal\runner\service.go | Select-Object -Skip 340 -First 280`
> - Ran `Get-Content internal\runner\service.go | Select-Object -Skip 620 -First 80`
> - Wrote `./go-ut-bench/internal/runner/service.go` (+38 -21)
> - Wrote `./go-ut-bench/internal/runner/service.go` (+16 -15)
>
> </details>
>
> 现在实现 runner 的核心适配层。`model_api` 会直接复用原 API 调用；`cli_agent` 会准备独立 workspace、注入 skill、写任务 prompt、执行 Docker 或 local 命令、收集 trace 和生成的测试文件。
>
> <details><summary>Created 1 file, edited 2 files</summary>
>
> - Wrote `./go-ut-bench/internal/runner/subjects.go` (+586 -0)
> - Wrote `./go-ut-bench/internal/runner/subjects.go` (+11 -1)
> - Wrote `./go-ut-bench/internal/runner/subjects.go` (+8 -0)
> - Wrote `./go-ut-bench/internal/evaluator/service.go` (+20 -11)
>
> </details>
>
> runner/evaluator 已经能携带 subject 元数据。下一步把 CLI 的 `--agents-config` 和 `--subjects` 接进去，并保留旧 `--models` 入口作为 `model_api + model + no_skill`。
>
> <details><summary>Edited 2 files</summary>
>
> - Wrote `./go-ut-bench/cmd/utbench/main.go` (+8 -0)
> - Wrote `./go-ut-bench/internal/dataset/service.go` (+1 -1)
>
> </details>
>
> 接着改报告聚合：排行榜仍沿用 by-model 的表名以兼容前端，但每个 row 会带 subject/framework/model/skill；同时新增两张 JSON 表用于 Agent-vs-API 和 Skill A/B 增益。
>
> <details><summary>Edited 1 file, ran 2 commands</summary>
>
> - Wrote `./go-ut-bench/internal/reporter/service.go` (+13 -7)
> - Wrote `./go-ut-bench/internal/reporter/service.go` (+28 -8)
> - Wrote `./go-ut-bench/internal/reporter/service.go` (+10 -4)
> - Wrote `./go-ut-bench/internal/reporter/service.go` (+174 -3)
> - Ran `rg "func avg|func avgFloat|func rate" internal\reporter\service.go -n`
> - Ran `Get-Content internal\reporter\service.go | Select-Object -Skip 1568 -First 35`
> - Wrote `./go-ut-bench/internal/reporter/service.go` (+58 -4)
>
> </details>
>
> 报告 JSON 已经会产出 subject 维度和两类 uplift 表。现在更新 SQLite 入库，让 subject/framework/model/skill 字段落库；同时用 `ALTER TABLE ... ADD COLUMN` 兼容已有数据库。
>
> <details><summary>Edited 1 file, ran 2 commands</summary>
>
> - Wrote `./go-ut-bench/internal/store/sqlite.go` (+32 -14)
> - Wrote `./go-ut-bench/internal/store/sqlite.go` (+26 -8)
> - Wrote `./go-ut-bench/internal/store/sqlite.go` (+35 -8)
> - Wrote `./go-ut-bench/internal/store/sqlite.go` (+39 -24)
> - Wrote `./go-ut-bench/internal/store/sqlite.go` (+44 -29)
> - Ran `Get-Content internal\store\sqlite.go | Select-Object -Skip 1000 -First 90`
> - Ran `Get-Content internal\store\sqlite.go | Select-Object -Skip 1070 -First 70`
> - Wrote `./go-ut-bench/internal/store/sqlite.go` (+23 -16)
> - Wrote `./go-ut-bench/internal/store/sqlite.go` (+11 -2)
>
> </details>
>
> 主要代码改完了，先跑 `gofmt` 再编译测试，下一步会根据编译错误补齐遗漏的字段/SQL 占位符。
>
> <details><summary>Ran 2 commands</summary>
>
> - Ran `gofmt -w internal\contracts\spec.go internal\contracts\results.go internal\agentconfig\config.go internal\runner\service.go internal\runner\subjects.go internal\evaluator\service.go internal\reporter\service.go internal\store\sqlite.go cmd\utbench\main.go`
> - Ran `go test ./internal/agentconfig ./internal/runner ./internal/reporter ./internal/store ./internal/evaluator ./cmd/utbench`
>
> </details>
>
> 测试第一次失败不是代码错误，是 Go 默认 cache 写到了用户目录且被沙箱拦住。我会把 `GOCACHE` 指到项目内 `.gocache` 后重跑。
>
> <details><summary>Edited 3 files, ran 7 commands</summary>
>
> - Ran `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/agentconfig ./internal/runner ./internal/reporter ./internal/store ./internal/evaluator ./cmd/utbench`
> - Wrote `./go-ut-bench/internal/runner/service.go` (+0 -1)
> - Ran `gofmt -w internal\runner\service.go`
> - Ran `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/agentconfig ./internal/runner ./internal/reporter ./internal/store ./internal/evaluator ./cmd/utbench`
> - Ran `rg "func firstNonEmpty" internal\reporter\service.go -n`
> - Ran `Get-Content internal\reporter\service.go | Select-Object -Skip 1368 -First 25; Get-Content internal\reporter\service.go | Select-Object -Skip 1598 -First 25`
> - Wrote `./go-ut-bench/internal/reporter/service.go` (+0 -9)
> - Wrote `./go-ut-bench/internal/runner/checkpoint_test.go` (+10 -9)
> - Ran `gofmt -w internal\reporter\service.go internal\runner\checkpoint_test.go`
> - Ran `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/agentconfig ./internal/runner ./internal/reporter ./internal/store ./internal/evaluator ./cmd/utbench`
>
> </details>
>
> 基础测试已经通过。为了覆盖新能力，我再补三类测试：subject 笛卡尔积/兼容过滤、fake CLI agent 本地生成测试文件、报告里的 Agent-vs-API 与 Skill uplift 计算。
>
> <details><summary>Created 2 files, edited 1 file, ran 2 commands</summary>
>
> - Wrote `./go-ut-bench/internal/agentconfig/config_test.go` (+61 -0)
> - Wrote `./go-ut-bench/internal/runner/subjects_test.go` (+81 -0)
> - Wrote `./go-ut-bench/internal/reporter/service_test.go` (+77 -0)
> - Ran `gofmt -w internal\agentconfig\config_test.go internal\runner\subjects_test.go internal\reporter\service_test.go`
> - Ran `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/agentconfig ./internal/runner ./internal/reporter ./internal/store ./internal/evaluator ./cmd/utbench`
>
> </details>
>
> 核心测试通过了。我再加一个 `configs/agents.example.yaml`，让新入口的配置格式可直接参考；随后跑一次相关测试确认没有格式/编译问题。
>
> <details><summary>Created 1 file, ran 3 commands</summary>
>
> - Wrote `./go-ut-bench/configs/agents.example.yaml` (+46 -0)
> - Ran `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/agentconfig ./internal/runner ./internal/reporter ./internal/store ./internal/evaluator ./cmd/utbench`
> - Ran `git diff --stat`
> - Ran `git status --short`
>
> </details>
>
> 我注意到配置里预留了 framework/skill 的语言兼容性，runner 目前还没按样本语言过滤。这里补上：任务总数和 checkpoint 都只针对合法的 subject × sample 组合计算。
>
> <details><summary>Edited 3 files, ran 4 commands</summary>
>
> - Wrote `./go-ut-bench/internal/runner/service.go` (+7 -1)
> - Wrote `./go-ut-bench/internal/runner/subjects.go` (+30 -0)
> - Ran `gofmt -w internal\runner\service.go internal\runner\subjects.go; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/agentconfig ./internal/runner ./internal/reporter ./internal/store ./internal/evaluator ./cmd/utbench`
> - Ran `git diff -- go-ut-bench/internal/agentconfig go-ut-bench/internal/runner go-ut-bench/internal/reporter go-ut-bench/internal/contracts go-ut-bench/internal/evaluator go-ut-bench/internal/store go-ut-bench/internal/dataset go-ut-bench/cmd/utbench go-ut-bench/configs/agents.example.yaml --stat`
> - Ran `git status --short go-ut-bench`
> - Wrote `./go-ut-bench/internal/reporter/service.go` (+3 -0)
> - Ran `gofmt -w internal\reporter\service.go; $env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/agentconfig ./internal/runner ./internal/reporter ./internal/store ./internal/evaluator ./cmd/utbench`
>
> </details>
>
> 相关包测试都通过了。我再跑一次全仓库 Go 测试，确认 web/store/orchestrator 这类间接编译点没有被新字段影响。
>
> <details><summary>Ran 1 command</summary>
>
> - Ran `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./...`
>
> </details>
>
> `go test ./...` 会把历史 `artifacts/runs/.../generated/tests/...` 下的模型生成 Go 文件也当成包编译，失败来自这些历史产物，不是本次改动。我改跑源码包范围 `./cmd/... ./internal/...`。
>
> <details><summary>Ran 2 commands</summary>
>
> - Ran `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./cmd/... ./internal/...`
> - Ran `git status --short`
>
> </details>
</details>

已实现第一版 Agent + Skill 评测升级，核心改动在 `go-ut-bench` 内：

- 新增 `Subject = framework + model + skill` 配置解析：`internal/agentconfig`
- Runner 支持 `model_api` baseline 和 `cli_agent` adapter，含 skill prompt 注入、workspace 文件注入、trace/diff/sandbox fingerprint 记录
- CLI 新增 `--agents-config` 和 `--subjects`
- `GeneratedCase` / `EvaluationResult` / SQLite 入库扩展 subject、agent、skill、trace 字段
- Reporter 新增 `agent_comparisons` 和 `skill_uplifts` 两类 A/B 增益结果
- 新增示例配置：`go-ut-bench/configs/agents.example.yaml`
- 新增覆盖测试：subject 矩阵生成、fake CLI agent、Agent-vs-API 和 Skill uplift

验证结果：

```powershell
$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./cmd/... ./internal/...
```

通过。

我也跑了 `go test ./...`，它会把历史 `artifacts/runs/.../generated/tests/...` 里的模型生成 Go 文件当作包编译，因此失败在历史产物，不是本次代码改动。

注意：仓库根部已有未关联改动 `CLAUDE.md`、`.codebuddy/`、`PROJECT_ANALYSIS_REPORT.md`，我没有处理这些。

Model changed

gpt-5.5 -> gpt-5.4

> 现在实现了什么,接下来要做什么,文档是不是要更新

**现状**

第一版主链路已经打通了，重点是把“被测对象”从单纯 `model` 扩成了 `subject = framework + model + optional skill`，并且没有把现有 evaluator/reporter 推翻重做。

已经实现的部分：

- `RunSpec`、`GeneratedCase`、`EvaluationResult`、`ReportPayload` 扩了 `subject/framework/model/skill/trace/sandbox_fingerprint` 字段。
- 新增了 `internal/agentconfig`，能解析 `agents.yaml` 风格配置，自动展开：
  - `framework × model × skill`
  - 自动补 `model_api__<model>__no_skill` baseline
  - 支持兼容性过滤和 `no_skill` baseline。
- runner 已经抽象成两类 subject：
  - `model_api`：复用原有 HTTP LLM 生成逻辑
  - `cli_agent`：支持独立 workspace、skill 注入、trace/diff 记录、Docker 或 local 执行
- CLI 已加入口：
  - `--agents-config`
  - `--subjects`
  - 旧 `--models` 仍可用，默认走 `model_api + no_skill`
- reporter 已支持新增两类聚合结果：
  - `agent_comparisons`
  - `skill_uplifts`
- SQLite 入库已兼容 subject/skill/trace 新字段，并做了老库 `ALTER TABLE` 兼容。
- 新增了示例配置 [agents.example.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.example.yaml)。

验证过的范围：

- `go test ./cmd/... ./internal/...` 通过
- `go test ./...` 不适合作为回归标准，因为会把历史 `artifacts/runs/.../generated/tests/...` 里的产物也当 Go 包编译

**还没做完的**

现在更像“架子和主链路可运行”，还没到“可直接大规模横评多 Agent”的完成态。

下一步最该做的是：

1. 真正落地至少一个真实 CLI Agent 适配器
现在测试里主要是 fake/local agent，配置结构和执行框架有了，但还没验证 `aider`、`claude code`、`opencode` 这类真实命令模板在 Docker 内的端到端行为。

2. 补强 Docker 沙箱细节
目前 runner 支持 Docker/local 两种执行方式，但还缺：
- 更严格的挂载边界
- 容器镜像规范
- 资源限制策略的验证
- 网络禁用策略的实际联调

3. 把单文件样本包装成更稳定的“最小 repo”
现在 `module_level` 复制 workspace、`self_contained` 复制单文件，这够第一版跑通，但还不够统一。后面最好抽一个明确的 `PreparedSample` 层。

4. 报告前端/UI 展示补齐
report JSON 已有 `agent_comparisons` 和 `skill_uplifts`，但现有 HTML/Web UI 还没系统展示这些新维度。

5. 数据库查询接口补齐
目前入库有了，但 DB 查询结果、facets、列表页筛选还基本沿用旧 `model` 视角，应该继续补 `subject/framework/skill` 维度。

**文档要不要更新**

要，而且现在就该补，不然代码和使用方式已经脱节了。

优先更新这几份：

- [readme.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/readme.md)
要加新的概念：`subject`、`skill`、`agents-config`、`subjects`

- [QUICKSTART.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/QUICKSTART.md)
要补最小可跑示例，尤其是：
  - 纯模型 baseline 跑法
  - Agent + no_skill 跑法
  - Agent + skill 跑法

- [docs/USER_GUIDE.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/USER_GUIDE.md)
要补配置文件格式和输出说明

- [docs/output-convention.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/output-convention.md)
要补 `trace_path`、`workspace_diff_path`、`subject_id` 这些新产物

如果你要我继续，我建议下一步直接做两件事：

1. 更新文档
2. 选一个真实 CLI Agent 先打通端到端，优先 `aider` 或你们内部最容易接的那个

> 把文档更新,然后cli agent话是要怎么调用,挑一个claudecode或者opencode看看哪个合适,docker沙箱我们现在只是把整个项目用了个docker吧,其实不算沙箱,我们最终要做到怎么样的沙箱,我的理解是不是要隔离到每一个agent执行的一个样本

<details><summary>15 previous messages</summary>

> 我会先做两件事：更新仓库里最直接相关的文档，然后查 `Claude Code` 和 `OpenCode` 的官方接入方式，判断哪个更适合作为第一个真实 CLI agent 适配目标。顺手我也会把“现在的 Docker 只是运行环境，不是最终沙箱”的边界写清楚。
>
> <details><summary>Searched web 3 times</summary>
>
> <details><summary>Searched web 3 times</summary>
>
> - Searched web for `Claude Code official docs CLI non-interactive command model provider docs`
> - Searched web for `https://docs.anthropic.com/en/docs/claude-code/cli-reference`
> - Searched web for `'run' in https://opencode.ai/docs/cli/`
>
> </details>
>
> </details>
>
> 我已经确认了官方能力边界：`OpenCode` 明确支持 `opencode run --model provider/model --agent ... --format json`，并且 provider 可以自定义 `baseURL`；`Claude Code` 则有稳定的非交互 `claude -p` 和系统提示追加能力，但模型切换主要面向 Anthropic 自身模型，不像 OpenCode 那样天然适合“任意框架 × 任意底层模型”。我先更新仓库文档，再把第一个真实 CLI agent 的建议和最终沙箱形态写进去。
>
> <details><summary>Ran 4 commands</summary>
>
> - Ran `Get-Content readme.md`
> - Ran `Get-Content QUICKSTART.md`
> - Ran `Get-Content docs\USER_GUIDE.md`
> - Ran `Get-Content docs\output-convention.md`
>
> </details>
>
> Context compaction
>
> Source: automatic
> Status: completed
>
> 我先把文档现状和实际实现对齐，再给你一个明确的 CLI Agent 和沙箱建议。重点会放在三件事：怎么用、现在做到哪、最终沙箱边界是什么。
>
> <details><summary>Ran 4 commands</summary>
>
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\readme.md"`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\QUICKSTART.md"`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\docs\USER_GUIDE.md"`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\docs\output-convention.md"`
>
> </details>
>
> 我先把文档里的命令、字段和实际代码对齐，避免写成“设计稿文档”。然后我会把 CLI Agent 的推荐和沙箱边界直接写进说明里。
>
> <details><summary>Ran 10 commands</summary>
>
> - Ran `rg -n "agent_traces|workspace_diff_path|sandbox_fingerprint|agents-config|subjects|subject_id|skill_uplifts|agent_comparisons|docker_image|execution_mode" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\cmd\utbench"`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\configs\agents.example.yaml"`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\subjects.go"`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\cmd\utbench\main.go"`
> - Ran `Get-Content -TotalCount 260 "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\agentconfig\config.go"`
> - Ran `rg -n "type SkillSpec|type SubjectSpec|RunSpec struct" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\contracts"`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\contracts\spec.go"`
> - Ran `rg -n "generated/metadata|generated/tests|agent_workspaces|generated_manifest.json|report_summary.json|evaluation_result.json" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\reporter" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\evaluator"`
> - Ran `rg -n "prompts|responses|metadata|testsRoot|promptPath|respPath|genRoot" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\service.go"`
> - Ran `Get-Content "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\service.go" | Select-Object -Skip 280 -First 120`
>
> </details>
>
> 我现在直接改文档，把“subject/skill/cli_agent/沙箱”都写进去，并把旧的 model-only 命令修正掉。文档会按当前实现来写，不会把还没做的部分写成已支持。
>
> <details><summary>Created 4 files</summary>
>
> - Wrote `./go-ut-bench/QUICKSTART.md` (+215 -0)
> - Wrote `./go-ut-bench/docs/USER_GUIDE.md` (+339 -0)
> - Wrote `./go-ut-bench/docs/output-convention.md` (+224 -0)
> - Wrote `./go-ut-bench/readme.md` (+206 -0)
>
> </details>
>
> 文档主体已经改完了。我再做一遍快速自检，主要看链接、命令参数名和目录路径有没有写偏。
>
> <details><summary>Ran 4 commands</summary>
>
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\readme.md"`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\QUICKSTART.md"`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\docs\USER_GUIDE.md"`
> - Ran `git -C "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench" status --short -- "go-ut-bench/readme.md" "go-ut-bench/QUICKSTART.md" "go-ut-bench/docs/USER_GUIDE.md" "go-ut-bench/docs/output-convention.md"`
>
> </details>
>
> 有两处小偏差我顺手收一下：一个是 README 的快速开始链接，一个是示例模型名我改成已在 Agent 示例里使用的 `deepseek-v4-flash`，避免文档自己打架。
>
> <details><summary>Edited 2 files, ran 1 command</summary>
>
> - Wrote `./go-ut-bench/docs/USER_GUIDE.md` (+1 -1)
> - Wrote `./go-ut-bench/readme.md` (+1 -1)
> - Ran `git -C "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench" diff -- "go-ut-bench/readme.md" "go-ut-bench/QUICKSTART.md" "go-ut-bench/docs/USER_GUIDE.md" "go-ut-bench/docs/output-convention.md"`
>
> </details>
</details>

文档已经更新，主要是把当前实现和用法写实了，没有再按旧的 model-only 口径写。

更新的文件：
- [readme.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/readme.md)
- [QUICKSTART.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/QUICKSTART.md)
- [USER_GUIDE.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/USER_GUIDE.md)
- [output-convention.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/output-convention.md)

这次文档补了三块关键内容：
- `subject = framework + model + optional skill`
- `cli_agent` 的调用方式、`agents.yaml` 结构、模板变量和执行契约
- Docker 运行环境 vs 真正 Agent 沙箱的区别，以及当前实现边界

关于 CLI Agent 选型，我建议先接 `OpenCode`，不是 `Claude Code`。原因很简单：你现在要的是 `framework × model` 的横评矩阵，`OpenCode` 更适合作为“Agent 壳 + 多模型后端”的统一入口；`Claude Code` 更适合后面做 Claude 自身工作流的专项评测。文档里也已经把这个建议写进去了。

关于你说的沙箱理解，你理解是对的。现在外层 `docker run utbench:latest ...` 只是把整个 benchmark 运行在 Docker 里，这不等于最终 Agent 沙箱。最终正确的隔离粒度应该是：

```text
一个 agent subject 执行一个 sample
= 一个独立 workspace
= 最好再对应一个独立容器
```

也就是按 `subject × sample` 隔离。当前第一版代码已经朝这个方向走了：
- 每个 `subject × sample` 都会建独立 workspace
- `sandbox_mode: docker` 时会起一个内层容器，只挂载这个 workspace
- 会记录 trace、workspace diff、sandbox fingerprint

但这还不是最终完成态。后面要补的是：
- 给每种语言准备专用 Agent 基础镜像
- 固定镜像 digest 和依赖版本
- 更严格的只读挂载边界
- CPU / memory / timeout / network 策略
- 必要时再上 gVisor / Firecracker / E2B 这类更强隔离层

这次是文档更新，没有改代码逻辑，所以我没再跑测试。下一步最值得做的是直接把 `OpenCode` 适配器打通，顺手把第一套“每样本一个内层容器”的镜像规范定下来。

> 那现在的意思是让opencode当基座?opencode+模型+skill吗,目前这样可以,不过后面我们想看claude code 和opencode 之间的评测又该怎么做

不是让 `OpenCode` 当唯一基座，而是把它当第一个比较适合落地的 `framework` 适配器。

也就是现在的抽象还是这个：

```text
subject = framework + model + optional skill
```

所以如果先接 `OpenCode`，被测对象会是：

```text
opencode + model + skill
```

比如：

- `opencode__deepseek-v4-flash__no_skill`
- `opencode__deepseek-v4-flash__unit_test_skill`

这只是第一步，不是说以后系统就围着 `OpenCode` 转。

你后面想看 `Claude Code` 和 `OpenCode` 之间的评测，设计上就是再加一个 `claude_code` framework：

- `opencode__deepseek-v4-flash__no_skill`
- `claude_code__deepseek-v4-flash__no_skill`

如果两边都能稳定接同一个底层模型，这就是最干净的“框架对框架”对比，因为只变 `framework`，不变 `model`。

但这里有一个很现实的点要提前讲清楚：`Claude Code` 未必适合你想要的那种“任意底层模型切换”。所以后面会有两类对比：

1. 同模型对比
   条件是两个 framework 都支持同一个 model。
   这时比较的是纯框架能力差异，最有说服力。

2. 各自原生能力对比
   比如：
   - `claude_code + Claude系模型`
   - `opencode + DeepSeek / GPT / GLM`

   这时比较的是“完整 Agent 方案”，不是纯框架差异。这个也有价值，但口径要单独标，不然会混淆。

所以我建议你们后面的评测表分两层：

- `Framework apples-to-apples`
  - 固定相同 model
  - 比较 `claude_code vs opencode`
- `End-to-end stack`
  - 比较完整组合
  - `framework + model + skill` 作为整体被测对象

换句话说：

- 现在先落 `OpenCode`，是因为它更容易把 `framework × model × skill` 这套机制跑通
- 后面加 `Claude Code`，不是替代 `OpenCode`
- 两者最终都只是 `framework` 维度里的候选项

你们最终的横评对象应该长这样：

```text
model_api__deepseek-v4-flash__no_skill
opencode__deepseek-v4-flash__no_skill
opencode__deepseek-v4-flash__unit_test_skill
claude_code__claude-sonnet__no_skill
claude_code__claude-sonnet__unit_test_skill
```

如果后面 `Claude Code` 也能稳定挂同一个非 Claude 模型，再补：

```text
claude_code__deepseek-v4-flash__no_skill
```

那时就能做真正严格的框架对比了。

下一步最合理的是：
1. 先把 `OpenCode` 接成第一个真实 framework
2. 配一套 `opencode + deepseek-v4-flash + no_skill / skill`
3. 再评估 `Claude Code` 的模型切换能力和非交互调用稳定性
4. 最后决定它是进入“同模型对比组”还是“原生方案对比组”

> 了解了,开始做吧

<details><summary>15 previous messages</summary>

> 我先把 `OpenCode` 作为第一个真实 framework 落地，但不会把架构绑死在它上面。先做两件最关键的：一是让配置能把宿主环境里的 API key 透传进每个 Agent 沙箱，二是补一套可直接参考的 `OpenCode` 配置和测试。
>
> <details><summary>Ran 3 commands, searched web 4 times</summary>
>
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\agentconfig\config.go"`
>
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\subjects_test.go"`
>
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\agentconfig\config_test.go"`
>
> <details><summary>Searched web 4 times</summary>
>
> - Searched web for `https://opencode.ai/docs/models/`
> - Searched web for `'opencode run' in https://opencode.ai/docs/cli/`
> - Searched web for `https://opencode.ai/docs/providers/`
> - Searched web for `'OPENROUTER_API_KEY' in https://opencode.ai/docs/providers/`
>
> </details>
>
> </details>
>
> 我已经把关键约束确认清楚了：`OpenCode` 现在确实支持非交互 `run`、`--model provider/model`，也支持用 `OPENCODE_CONFIG` 或 `OPENCODE_CONFIG_CONTENT` 注入配置。这样我们不需要把它写死进代码，可以继续走通用 framework 适配，只补它真正需要的环境和配置能力。
>
> <details><summary>Ran 7 commands, searched web 5 times</summary>
>
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\configs\models.yaml"`
>
> - Ran `rg -n "type modelConfig|Model string|provider" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner"`
>
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\models.go"`
>
> - Ran `rg -n "api_key_env" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner"`
>
> <details><summary>Searched web 2 times</summary>
>
> - Searched web for `site:opencode.ai/docs OPENCODE_CONFIG_CONTENT custom provider baseURL envKey openrouter openai-compatible`
> - Searched web for `site:opencode.ai/docs opencode run --model --format json`
>
> </details>
>
> - Ran `rg -n "type subjectTarget" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner"`
>
> - Ran `Get-Content "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\service.go" | Select-Object -First 80`
>
> - Ran `Get-ChildItem -Recurse "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\configs" | Select-Object FullName`
>
> <details><summary>Searched web 3 times</summary>
>
> - Searched web for `site:opencode.ai/docs/cli opencode run --agent --format json`
> - Searched web for `site:opencode.ai/docs OPENCODE_CONFIG_CONTENT providers baseURL apiKey`
> - Searched web for `site:opencode.ai/docs OPENCODE_CONFIG_CONTENT OPENCODE_CONFIG`
>
> </details>
>
> </details>
>
> 我现在补的是让真实 `OpenCode` 配置能工作，而不是再造一个专门的 `opencode` adapter。这样后面加 `Claude Code` 还是同一套框架，只是换一个 `framework` 条目。
>
> <details><summary>Edited 2 files, ran 1 command</summary>
>
> - Wrote `./go-ut-bench/internal/agentconfig/config.go` (+3 -0)
> - Wrote `./go-ut-bench/internal/runner/subjects.go` (+71 -10)
> - Ran `rg -n "renderCommand\(|runShellCommand\(|runAgentCommand\(" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner"`
>
> </details>
>
> 我在补两个测试点：一个验证 `agents.yaml` 能解析 `env_from_host`，另一个验证真实执行时，模型 API key 和渲染后的 `OPENCODE_CONFIG_CONTENT` 能进入 Agent 进程。这样后面接 OpenCode 至少不是“看起来能配”。
>
> <details><summary>Edited 2 files, searched web 2 times</summary>
>
> - Wrote `./go-ut-bench/internal/agentconfig/config_test.go` (+45 -0)
>
> - Wrote `./go-ut-bench/internal/runner/subjects_test.go` (+98 -0)
>
> - Wrote `./go-ut-bench/internal/runner/subjects_test.go` (+1 -0)
>
> <details><summary>Searched web 2 times</summary>
>
> - Searched web for `site:opencode.ai/docs/cli "--output-last-message" opencode`
> - Searched web for `site:opencode.ai/docs/providers "providers" "baseURL" "apiKey" "@ai-sdk/openai-compatible"`
>
> </details>
>
> </details>
>
> 我把示例配置切到 `OpenCode`，这样仓库里的“example”就和我们实际要接的第一个 framework 对上了。配置里会直接示范怎么把模型 endpoint、模型 ID 和 API key env 拼成 `OPENCODE_CONFIG_CONTENT`。
>
> <details><summary>Created 2 files, edited 3 files, ran 4 commands</summary>
>
> - Wrote `./go-ut-bench/configs/agents.example.yaml` (+83 -0)
> - Wrote `./go-ut-bench/configs/skills/unit_test_skill.md` (+9 -0)
> - Ran `rg -n "aider__|aider|agents.example.yaml|env_from_host|OpenCode|opencode__|opencode" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\readme.md" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\QUICKSTART.md" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\docs\USER_GUIDE.md"`
> - Ran `Get-Content "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\QUICKSTART.md" | Select-Object -Skip 70 -First 60`
> - Ran `Get-Content "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\docs\USER_GUIDE.md" | Select-Object -Skip 60 -First 130`
> - Ran `Get-Content "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\readme.md" | Select-Object -Skip 108 -First 20`
> - Wrote `./go-ut-bench/QUICKSTART.md` (+3 -1)
> - Wrote `./go-ut-bench/docs/USER_GUIDE.md` (+12 -6)
> - Wrote `./go-ut-bench/readme.md` (+1 -1)
>
> </details>
>
> 我先把 Go 代码格式化，再跑最小回归：`agentconfig`、`runner`、`reporter`、`cmd/utbench`。这轮主要验证的是新加的环境透传、模板渲染和示例配置不会把主链路打断。
>
> <details><summary>Ran 3 commands, searched web 1 time</summary>
>
> - Ran `gofmt -w "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\agentconfig\config.go" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\agentconfig\config_test.go" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\subjects.go" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\subjects_test.go"`
>
> - Ran `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/agentconfig ./internal/runner ./internal/reporter ./cmd/utbench`
>
> - Ran `git -C "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench" status --short -- "go-ut-bench"`
>
> <details><summary>Searched web 1 time</summary>
>
> - Searched web for `site:opencode.ai/docs install opencode npm -g`
>
> </details>
>
> </details>
>
> 我再补一个最小的 `OpenCode` Agent 镜像骨架，不然 `utbench-agent-opencode:latest` 只有配置名，没有对应的 build path。这个不会替代 benchmark 主镜像，只是内层每样本沙箱要用的 agent 镜像。
>
> <details><summary>Created 1 file, edited 2 files, ran 1 command, searched web 1 time</summary>
>
> - Wrote `./go-ut-bench/QUICKSTART.md` (+6 -0)
> - Wrote `./go-ut-bench/docker/agents/opencode/Dockerfile` (+9 -0)
> - Wrote `./go-ut-bench/docs/USER_GUIDE.md` (+6 -0)
>
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\configs\agents.example.yaml"`
>
> <details><summary>Searched web 1 time</summary>
>
> - Searched web for `site:opencode.ai/docs/provider config "provider" opencode json`
>
> </details>
>
> </details>
>
> <details><summary>Plan</summary>
>
> OpenCode-focused framework support is now practical: host env passthrough, templated agent env, example config, skill asset, and focused tests are in place.
>
> - [x] Inspect current agent config and runner env handling to identify gaps for a real OpenCode adapter
> - [x] Implement host-env passthrough and any supporting config/schema updates for cli_agent sandboxes
> - [x] Add OpenCode example config/skill assets and focused tests
> - [x] Update docs to show concrete OpenCode usage and verify the changed test targets
>
> </details>
</details>

我已经把第一步做成可落地的 `OpenCode` framework 路径了，不是把架构绑到 `OpenCode` 上，而是让它成为第一个真实可接的 `framework`。

这次实际补了四件事：

1. `cli_agent` 现在支持宿主环境变量透传
   代码里新增了 `env_from_host`，并且会自动把当前模型配置里的 `api_key_env` 一起带进每个 Agent 进程或容器。这样 `OpenCode` 这类需要在内层沙箱里拿 API key 的框架能真正跑起来。

2. `framework.env` 现在支持模板渲染
   所以你可以直接在 `agents.yaml` 里按当前 subject 的模型信息拼出 `OPENCODE_CONFIG_CONTENT`，里面可以用：
   - `{{.ModelID}}`
   - `{{.ModelEndpoint}}`
   - `{{.ModelAPIKeyEnv}}`
   - `{{.ModelProvider}}`

3. 仓库里的 example 已切到 `OpenCode`
   我把 [agents.example.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.example.yaml) 改成了 `opencode` 示例，并加了一个 skill 文件：
   - [unit_test_skill.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/skills/unit_test_skill.md)

4. 加了一个内层 Agent 镜像骨架
   现在有了 [docker/agents/opencode/Dockerfile](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docker/agents/opencode/Dockerfile)，对应 `utbench-agent-opencode:latest` 这个沙箱镜像名，不再只是文档里空写一个名字。

文档也同步更新了：
- [readme.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/readme.md)
- [QUICKSTART.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/QUICKSTART.md)
- [USER_GUIDE.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/USER_GUIDE.md)
- [output-convention.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/output-convention.md)

测试我跑了这组，都是过的：

```powershell
$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/agentconfig ./internal/runner ./internal/reporter ./cmd/utbench
```

当前状态要讲清楚：

- benchmark 侧已经能把 `OpenCode` 当成一个真实 `framework` 来驱动
- 现在还没做的是“真实安装好的 OpenCode 镜像 + 真实模型 key + 一次端到端样本跑通”
- 也就是说，框架接线和配置能力已经就位，但还差最后那一步 smoke run

你后面想比较 `Claude Code` 和 `OpenCode`，架构上已经不用再改。到时候只要再加一个 `claude_code` framework：
- 一个 `claude_code` 的 `agents.yaml` 配置
- 一个对应的 Agent 沙箱镜像
- 一次真实 smoke run 验证

下一步最该做的是两件事：

1. 构建 `utbench-agent-opencode:latest`
2. 用一个最小样本跑一次真实 `opencode__deepseek-v4-flash__no_skill`

这样我们就能从“配置和测试通过”进入“真实 benchmark 可跑”的阶段。

> 现在他去对某一个样本执行生成代码是怎么样的一个流程

现在对“某一个样本”执行生成，流程是这样的。

先定一个被测对象，也就是一个 `subject`：

```text
subject = framework + model + skill
```

比如：

```text
opencode__deepseek-v4-flash__unit_test_skill
```

然后对这个 `subject × sample`，runner 会单独走一遍生成流程。

**1. 选中样本和 subject**

UT-Bench 先从数据集里拿到一个 `SampleRef`，里面有：
- `language`
- `sample_id`
- `sample_path`
- `class`
- `scenario`

然后把它和一个 `subject` 配对，形成一个最小任务单元：

```text
one task = one subject × one sample
```

**2. 准备 prompt**

runner 会先按样本类型生成基础 prompt：

- `self_contained`：直接围绕这个源码文件生成测试
- `module_level`：会带上 workspace / module 上下文

如果这个 subject 带了 skill，而且 skill 是 `prompt_append`，就把 skill instruction 追加到 prompt 后面。

最后会生成一个完整 prompt 文件，写到：

```text
generated/prompts/rendered/<subject>/<lang>/<sample>.prompt.txt
```

**3. 准备 Agent 工作区**

如果是 `model_api` baseline，到这里就直接调模型 API 了，不需要工作区。

如果是 `cli_agent`，会先创建一个独立 workspace：

```text
artifacts/runs/<run-id>/agent_workspaces/<subject>/<lang>/<sample>/
```

然后把样本放进去：

- `self_contained`：复制这个源码文件进去
- `module_level`：复制整个最小 workspace / repo 进去

如果 skill 需要 `workspace_mount`，也会把 skill 文件复制进这个 workspace。

**4. 写入 Agent 执行契约**

runner 会在 workspace 里生成一个 `utbench_agent_prompt.md`，里面除了主 prompt，还会明确告诉 Agent：

- 只能在当前 workspace 里工作
- 不要改原始源码行为
- 生成一个完整测试文件
- 最终把测试文件写到指定输出路径
- 源文件在哪里
- skill 文件在哪里

也就是说，Agent 拿到的不是一句裸 prompt，而是一份带执行约束的任务说明。

**5. 渲染 CLI 命令**

接着用 `framework.command` 模板渲染真实命令。

比如 `OpenCode` 这种 framework，会把这些变量带进去：

- `ModelID`
- `ModelEndpoint`
- `ModelAPIKeyEnv`
- `ContainerPrompt`
- `ContainerOutput`
- `ContainerSkillDir`

同时渲染环境变量，比如：

- `OPENCODE_CONFIG_CONTENT`
- 以及宿主环境透传进来的 API key

**6. 启动执行**

如果 `sandbox_mode=local`：
- 直接在本机 shell 里执行这个命令，但工作目录仍然是该样本的独立 workspace

如果 `sandbox_mode=docker`：
- 启动一个内层容器
- 只挂载这个 workspace 到 `/workspace`
- 在容器里执行命令
- 可选禁网、限 CPU、限内存、限时

这一层才是 Agent 真正的样本级执行沙箱。

**7. 记录执行轨迹**

执行前后会做两件事：

- 对 workspace 做文件快照
- 把 stdout / stderr / exit code / duration 写入 trace

产物包括：

```text
generated/metadata/agent_traces/<subject>/<lang>/<sample>.trace.jsonl
generated/metadata/agent_traces/<subject>/<lang>/<sample>.diff.json
```

所以后面你可以回看：
- Agent 跑了多久
- 命令是否成功
- 改了哪些文件

**8. 找生成出的测试文件**

命令跑完后，runner 会按这个顺序找测试文件：

1. 先看约定输出路径有没有文件
2. 再按 `output_globs` 搜
3. 再看 workspace diff 里新增/修改的测试文件

找到后把内容读出来，作为这次生成结果。

如果没找到，就标记为：

```text
agent_output_error
```

**9. 结果落盘**

成功的话，会把最终测试代码写到：

```text
generated/tests/<subject>/<lang>/<sample>.test<ext>
```

同时写：
- response json
- metadata json
- manifest case 记录

这个 case 里会带：
- `subject_id`
- `agent_framework`
- `agent_model`
- `skill_name`
- `trace_path`
- `workspace_diff_path`
- `sandbox_fingerprint`

**10. 后续进入评测**

生成阶段结束后，这个样本的产物就和普通 generated case 一样，后面 evaluator 会继续做：

- 编译
- 测试
- 覆盖率
- 变异测试

所以从 evaluator 视角看，它并不关心测试是纯模型生成的，还是 Agent 在沙箱里生成的。

如果你要一句话概括：

```text
当前流程就是：
对每个 subject × sample，
先准备独立 prompt 和独立 workspace，
再在独立执行环境里跑 Agent，
收集测试文件、trace 和 diff，
最后把结果送进原有 evaluator。
```

如果你愿意，下一步我可以直接给你画一张“单样本执行时序图”，把 `runner -> workspace -> sandbox -> trace -> manifest` 这条链画出来。

> 那这里的docker沙箱是怎么样的

现在这个 Docker 沙箱，准确说是“每个 `subject × sample` 一个内层容器”的执行沙箱，但还不是最终加强版安全沙箱。

它当前是这样工作的。

**1. 外层和内层是两层 Docker**

先分清楚：

- 外层 Docker
  你平时跑的：
  ```bash
  docker run utbench:latest ...
  ```
  这是 UT-Bench 本身的运行环境。

- 内层 Docker
  当某个 framework 配了：
  ```yaml
  sandbox_mode: docker
  ```
  runner 会再起一个容器，专门给这个样本上的 Agent 用。

所以真正的 Agent 沙箱，是内层这个容器，不是外层 benchmark 容器。

**2. 沙箱的粒度**

当前粒度是：

```text
一个 subject 执行一个 sample
= 一个独立 workspace
= 一个独立 docker 容器
```

比如：

```text
opencode__deepseek-v4-flash__unit_test_skill
× python/simple_function_001
```

就会有自己单独的 workspace 和单独的容器。

这点很重要，因为它保证了：
- 样本之间不会互相污染
- skill 之间不会互相污染
- 不同 framework 的运行痕迹不会混在一起
- 同一样本上的 baseline 和 agent 结果可公平比较

**3. 容器里实际挂了什么**

当前实现里，内层容器只挂一个目录：

```text
-v <workspace>:/workspace
-w /workspace
```

也就是说，Agent 在容器里看到的主要就是这个样本的工作区。

这个 workspace 里通常包括：
- 被测源码文件，或最小 repo 副本
- `utbench_agent_prompt.md`
- `.utbench/skills/...` 里的 skill 文件
- Agent 最终要写出的测试文件位置

所以现在的核心隔离手段是：

```text
只把当前样本工作区暴露给当前 Agent 容器
```

**4. 当前容器启动时带的限制**

现在代码里已经支持这些限制：

- `--rm`
  - 跑完自动销毁容器
- `--network none`
  - 如果配置了 `network_disabled: true`
- `--cpus <n>`
- `--memory <size>`
- 超时控制
  - runner 用 context timeout 杀掉这次执行

也就是 framework 配置里这些字段已经进沙箱执行了：

```yaml
sandbox_mode: docker
docker_image: utbench-agent-opencode:latest
timeout_seconds: 600
network_disabled: true
cpu: "2"
memory: 2g
```

**5. 容器里怎么执行命令**

当前逻辑等价于：

```bash
docker run --rm \
  --network none \
  --cpus 2 \
  --memory 2g \
  -e ... \
  -v <workspace>:/workspace \
  -w /workspace \
  utbench-agent-opencode:latest \
  /bin/sh -lc "<rendered command>"
```

也就是说：
- UT-Bench 先渲染出 framework command
- 再把这条命令放进容器里执行
- 容器退出后收集 stdout/stderr 和工作区改动

**6. 环境变量怎么进沙箱**

当前有两类：

- 静态 `env`
  - 直接在 `agents.yaml` 里写
  - 支持模板渲染
- `env_from_host`
  - 从宿主机透传进容器
  - 典型是 API key

比如 `OpenCode` 这次就是这么用的：
- `OPENCODE_CONFIG_CONTENT` 走静态模板渲染
- `DEEPSEEK_API_KEY` / `ARK_API_KEY` 这类走 `env_from_host`

所以沙箱不是完全无环境，而是“受控注入环境”。

**7. 它现在能防什么**

现在这版已经能防住不少最直接的问题：

- 一个样本改坏另一个样本目录
- 一个 Agent 执行污染另一个 Agent 的输出
- 不同 subject 共享同一个运行现场
- 无限制访问整个 benchmark 工作目录
- 没有限制地长期占用资源

从“评测公平性”和“基础隔离”角度，这已经比“整个项目一个大容器里跑所有 Agent”强很多。

**8. 它现在还防不住什么**

但如果你问的是严格意义上的高强度安全沙箱，那现在还不够。

当前还缺这些：

- 更严格的只读根文件系统
- 更明确的 mount 白名单
- 非 root 用户运行约束
- seccomp / capabilities / no-new-privileges
- syscall 级限制
- 镜像 digest 固定
- 更强隔离层，比如：
  - gVisor
  - Firecracker microVM
  - E2B 这类托管沙箱

换句话说，现在是：

```text
“按样本隔离的 Docker 执行沙箱”
```

但还不是：

```text
“高强度、对抗性安全沙箱”
```

**9. 你们最终应该做到什么程度**

如果目标是做正式 Agent benchmark，我建议最终标准至少是这样：

- 每个 `subject × sample` 独立 workspace
- 每个 `subject × sample` 独立容器
- 原始数据集不直接暴露，只复制副本
- 默认禁网，只有明确需要时才开放
- CPU / memory / timeout 强制限制
- 非 root 用户运行
- rootfs 尽量只读
- 只允许写 `/workspace`
- 全量记录：
  - command
  - stdout/stderr
  - file diff
  - exit code
  - image digest
  - env fingerprint
- 任务结束后销毁容器和临时状态

如果再往上提安全级别，再考虑：
- Docker + gVisor
- Firecracker microVM
- E2B

**10. 一句话总结**

你现在可以把它理解成：

```text
当前已经实现的是
“每个 Agent 在每个样本上的独立 Docker 执行单元”

还没完全实现的是
“高强度、不可逃逸、可审计完备的安全沙箱”
```

如果你要，我下一步可以直接把“当前内层 Docker 沙箱”和“最终目标沙箱”画成一个对比表，顺手列出该补的字段和实现项。

> 是dood吗

看运行方式。

如果 **UT-Bench 跑在宿主机上**，那不是 DOOD。只是 UT-Bench 这个进程直接调用宿主机的 `docker run`。

如果 **UT-Bench 自己也跑在 Docker 容器里**，而你又想让它在里面继续起每样本的 Agent 容器，那通常就是 **DOOD**：

- 容器里有 `docker` CLI
- 挂宿主机的 Docker socket
- 容器内执行的 `docker run` 实际在控制宿主机 daemon

也就是典型的：

```text
-v /var/run/docker.sock:/var/run/docker.sock
```

或 Windows 对应的 Docker daemon 接入方式。

关键点是：**我们现在代码层面是“调用 docker CLI 起内层容器”**。
所以：

- 宿主机运行 UT-Bench：不是 DOOD
- 容器内运行 UT-Bench 且继续起样本容器：就是 DOOD 风格

还有一个很实际的问题：
**你现在文档里的外层 `docker run utbench:latest ...` 并没有把 Docker socket/CLI 配进去。**
所以按当前文档，`sandbox_mode: docker` 的 CLI Agent 在外层容器里其实**不能直接工作**，除非你额外做 DOOD 接线。

所以当前更准确的结论是：

1. 代码支持“每样本起一个 Docker 沙箱”
2. 如果 benchmark 跑宿主机，这条路直接成立
3. 如果 benchmark 跑外层容器，要么做 DOOD，要么改成别的沙箱后端

我建议后面定一个明确策略：

- **开发期**：UT-Bench 跑宿主机，Agent 样本沙箱用 host Docker
- **正式版**：再决定是
  - DOOD
  - 远程 Docker daemon
  - E2B
  - Firecracker / gVisor

如果你要，我下一步可以直接把这三种部署形态画清楚：
- 宿主机跑 UT-Bench
- 外层容器 + DOOD
- 外层容器 + 远程沙箱后端

> 我们一般都是用docker来启用这个工具的

那你们现在这个前提下，**如果还要在评测过程中给每个样本再起一个 Agent 容器，本质上就会走 DOOD 路线**，或者至少是“容器内控制外部 Docker daemon”的路线。

也就是说，现状会变成两层：

```text
外层：utbench:latest
内层：每个 subject × sample 的 agent sandbox container
```

如果外层本身在 Docker 里跑，想让它再执行：

```bash
docker run ...
```

那必须满足至少一条：

1. 外层容器挂宿主机 Docker socket
   这是最典型的 DOOD

2. 外层容器连接远程 Docker daemon
   本质上还是“把容器内的 docker CLI 指向外部 daemon”

所以你说得对，按你们现在“工具通常通过 Docker 启动”的使用方式，后面要做每样本一个 Agent 沙箱，**不能只停留在代码里写 `docker run`**，必须把部署形态设计清楚。

我建议你们把方案分成两阶段。

**阶段 1：先把功能跑通**
直接接受 DOOD。

做法是：
- 外层 `utbench:latest` 里安装 `docker` CLI
- 挂 Docker socket
- 让外层容器能调用宿主机 Docker daemon
- 内层样本沙箱继续按 `subject × sample` 起容器

这样优点是最省改动，最容易把整条链跑通。
缺点也很明显：**安全边界比较弱**。挂 Docker socket 基本等于把宿主机 Docker 控制权给了外层容器。

**阶段 2：再把沙箱后端抽离**
等功能稳定后，把“样本沙箱执行器”从现在这个直接 `docker run`，抽成一个独立 backend，比如：

- `docker_local`
- `docker_remote`
- `e2b`
- `firecracker`
- `gvisor`

这样外层 benchmark 仍然是 Docker 启动，但内层不一定非得走 DOOD。

我建议你们最终目标是这个结构：

```text
UT-Bench container
  -> sandbox backend interface
     -> local docker daemon / remote docker / E2B / microVM
```

这样就不会把“benchmark 运行方式”和“agent 沙箱方式”绑死。

所以结论很明确：

- **短期**：如果你们现在普遍是 Docker 启动 UT-Bench，那要支持每样本一个 Agent 容器，基本就是 DOOD
- **中期**：应该把 DOOD 当过渡方案，不要当最终方案
- **长期**：把沙箱执行层抽象成独立 backend

从工程优先级看，我建议下一步不是先上更重的安全隔离，而是先做这两件事：

1. 明确外层 Docker 的 DOOD 启动方式
   也就是把 `utbench:latest` 怎么拿到 Docker daemon 控制权写清楚

2. 在代码里把当前 `docker run` 执行路径再包一层 `SandboxRunner` 抽象
   这样后面替换后端不会动 runner 主流程

如果你同意，我下一步可以直接继续做：
- 把当前内层 Docker 执行抽成 `SandboxRunner`
- 补一份“外层 Docker + DOOD” 的运行文档和启动命令
- 顺手把 `OpenCode` 的真实 smoke run 路径也整理出来

> 开始搞吧

<details><summary>9 previous messages</summary>

> 我先把内层容器执行从 runner 里剥成一个 `SandboxRunner`，这样 DOOD 只是一个 backend，不会和生成逻辑绑死。然后补外层 Docker 的 DOOD 用法文档，明确现在怎么跑、风险是什么。
>
> <details><summary>Ran 3 commands</summary>
>
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\Dockerfile"`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\docs\DOCKER_GUIDE.md"`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\subjects.go"`
>
> </details>
>
> 我会顺手把外层镜像也补上 `docker` CLI，不然文档写 DOOD，镜像本身却起不了内层样本容器，链路还是断的。安全上这只是过渡方案，我会在文档里把风险和边界写明。
>
> <details><summary>Created 1 file, edited 3 files, ran 5 commands</summary>
>
> - Wrote `./go-ut-bench/internal/runner/service.go` (+3 -2)
> - Wrote `./go-ut-bench/internal/runner/sandbox.go` (+172 -0)
> - Wrote `./go-ut-bench/internal/runner/sandbox.go` (+4 -2)
> - Ran `Get-Content "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\subjects.go" | Select-Object -First 340`
> - Ran `Get-Content "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\subjects.go" | Select-Object -Skip 340 -First 260`
> - Ran `Get-Content "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\subjects.go" | Select-Object -Skip 600 -First 80`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\subjects.go" | Select-String -Pattern "func sandboxFingerprint|func safePathName|func uniqueSortedStrings" -Context 0,20`
> - Wrote `./go-ut-bench/internal/runner/subjects.go` (+20 -112)
> - Ran `rg -n "type commandOutput|commandOutput" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner"`
> - Wrote `./go-ut-bench/internal/runner/subjects.go` (+0 -5)
>
> </details>
>
> 我现在补的是“外层容器怎么控制每样本内层容器”这条链。代码里已经抽出来了，文档里会明确写成 DOOD 过渡方案，并且要求显式挂 Docker socket；不挂的话，`sandbox_mode: docker` 的 framework 就不能用。
>
> <details><summary>Edited 3 files</summary>
>
> - Wrote `./go-ut-bench/Dockerfile` (+8 -7)
> - Wrote `./go-ut-bench/docs/DOCKER_GUIDE.md` (+121 -17)
> - Wrote `./go-ut-bench/readme.md` (+2 -0)
>
> </details>
>
> 我把抽象层和 DOOD 文档都补好了，接下来先做一轮编译级回归，确保 `runner` 和 CLI 没被这次重构打断。
>
> <details><summary>Ran 2 commands</summary>
>
> - Ran `gofmt -w "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\service.go" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\subjects.go" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\sandbox.go"`
> - Ran `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/runner ./internal/agentconfig ./internal/reporter ./cmd/utbench`
>
> </details>
>
> <details><summary>Plan</summary>
>
> Sandbox execution is now abstracted behind SandboxRunner, the outer Docker image can drive per-sample agent containers via DOOD, and the docs now explain that deployment path explicitly.
>
> - [x] Inspect current Docker docs/image setup and runner execution path for per-sample sandboxes
> - [x] Extract sandbox execution behind a SandboxRunner abstraction and keep local/docker backends behavior-compatible
> - [x] Document outer Docker + DOOD usage, constraints, and commands for Agent sandboxes
> - [x] Run focused regression tests for runner/cmd/docs-adjacent changes
>
> </details>
</details>

已经开工并落地了两部分：代码抽象和 DOOD 文档。

这次代码上做了这些：

- 把每样本 Agent 执行从 runner 主逻辑里抽成了 `SandboxRunner`
  - 新文件：[sandbox.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/sandbox.go)
  - `Service` 现在持有 `sandboxRunner`，`generateWithCLIAgent` 不再直接写死 `docker run`
  - 当前后端还是两种：
    - `local`
    - `docker`
- `docker` 后端仍保持现有行为
  - 每个 `subject × sample` 一个独立 workspace
  - 每个 `subject × sample` 一个独立容器
  - 环境变量渲染和 `env_from_host` 透传保持可用
- 沙箱指纹现在按实际 `SandboxRunRequest` 计算，不再只看 framework 配置

相关改动文件：
- [service.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/service.go)
- [subjects.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/subjects.go)
- [sandbox.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/sandbox.go)

外层 Docker / DOOD 这块也补了：

- 外层 benchmark 镜像现在安装了 `docker.io`
  - [Dockerfile](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/Dockerfile)
- 文档明确写了当前 `sandbox_mode: docker` 的实际部署方式就是 DOOD 过渡方案
  - [DOCKER_GUIDE.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/DOCKER_GUIDE.md)
  - [readme.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/readme.md)

现在文档已经明确了：
- 外层 `utbench:latest` 只是 orchestration 容器
- 真正的 Agent 样本级沙箱是内层容器
- 如果外层也跑在 Docker 里，要跑 `sandbox_mode: docker`，就必须挂 Docker socket
- 这是过渡方案，不是最终安全方案

另外之前做的 `OpenCode` 路径还在：
- [agents.example.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.example.yaml)
- [opencode Dockerfile](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docker/agents/opencode/Dockerfile)

我跑过的回归：

```powershell
$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/runner ./internal/agentconfig ./internal/reporter ./cmd/utbench
```

通过。

还没做的是实际构建镜像并跑一次真实 smoke run。我这轮没去跑 Docker build，因为那会依赖镜像下载和外部网络。下一步最值当的是直接做一次端到端验证：

1. 构建外层镜像 `utbench:latest`
2. 构建内层镜像 `utbench-agent-opencode:latest`
3. 挂 Docker socket 跑一个最小样本
4. 确认 `agent_workspaces/`、`trace.jsonl`、`generated/tests/` 都按预期落盘

如果你同意，我下一步就直接做这个真实 smoke run 路径。

> 目前的完成情况是怎么样,对了UT-Bench Agent 评测升级方案写一个文档,然后把已完成的或者哪里改了写上去,端到端验证我来我要怎么做

<details><summary>4 previous messages</summary>

> 我先补一份正式文档，把方案、已完成项、改动文件和待办都落下来。然后我给你一套可直接照着跑的端到端验证步骤，包括要构建哪些镜像、挂什么、看哪些产物。
>
> <details><summary>Ran 2 commands</summary>
>
> - Ran `Get-ChildItem "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\docs" | Select-Object Name`
> - Ran `git -C "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench" status --short -- "go-ut-bench"`
>
> </details>
>
> 我先把一个关键点修正掉：当前 `OpenCode` 示例里如果 `network_disabled: true`，真实端到端是跑不通的，因为内层 Agent 容器自己要去调模型 API。这个我会一并写进状态文档和验证步骤里。
>
> <details><summary>Created 1 file, edited 2 files</summary>
>
> - Wrote `./go-ut-bench/configs/agents.example.yaml` (+1 -1)
> - Wrote `./go-ut-bench/docs/AGENT_UPGRADE.md` (+409 -0)
> - Wrote `./go-ut-bench/readme.md` (+1 -0)
>
> </details>
</details>

我已经把正式文档写好了：

- [AGENT_UPGRADE.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/AGENT_UPGRADE.md)

我也把入口加到了：
- [readme.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/readme.md)

这份文档里已经写了：

- 升级目标和当前架构
- 当前完成情况
- 已完成内容
- 新增/修改了哪些文件
- 还没完成什么
- 你要怎么做端到端验证
- 成功后该检查哪些产物
- 失败时优先看哪里

另外我顺手修正了一个会影响真实跑通的问题：

- [agents.example.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.example.yaml)

这里的 `OpenCode` 示例现在改成了：

```yaml
network_disabled: false
```

原因很直接：`OpenCode` 在内层 Agent 容器里要自己访问模型 API，禁网的话真实端到端会直接失败。

**目前完成情况**

可以把当前状态理解成：

- 第一版架构改造：已完成
- `subject = framework + model + skill`：已完成
- 纯模型 baseline 兼容：已完成
- 通用 CLI Agent 执行链路：已完成
- Skill 注入：已完成第一版
- 每样本 Docker 沙箱：已完成第一版
- `SandboxRunner` 抽象：已完成
- `OpenCode` 作为首个真实 framework 的接入入口：已完成
- 外层 Docker + DOOD 文档和接线：已完成
- 真实 smoke run：还没做，需要你本地验证
- `Claude Code` framework：还没开始
- 更强安全沙箱后端：还没开始

**你要怎么做端到端验证**

最小验证流程已经写在文档第 6 节，我这里再给你一版精简步骤。

1. 准备 `.env`
至少有：

```dotenv
DEEPSEEK_API_KEY=...
```

2. 构建外层 benchmark 镜像

```bash
docker build -t utbench:latest .
```

3. 构建内层 OpenCode Agent 镜像

```bash
docker build -t utbench-agent-opencode:latest -f ./docker/agents/opencode/Dockerfile .
```

4. 运行最小 smoke run
建议在 WSL / Linux / macOS bash 下跑，因为 DOOD 的 Docker socket 挂载更直接：

```bash
docker run --rm --env-file .env \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek-v4-flash \
    --langs python \
    --config /app/configs/models.yaml \
    --agents-config /app/configs/agents.example.yaml \
    --subjects model_api__deepseek-v4-flash__no_skill,opencode__deepseek-v4-flash__no_skill,opencode__deepseek-v4-flash__unit_test_skill \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --class self_contained \
    --max-samples 1 \
    --run-id agent_smoke_001
```

5. 成功后检查这些产物：

- `artifacts/runs/agent_smoke_001/generated/generated_manifest.json`
- `artifacts/runs/agent_smoke_001/generated/tests/`
- `artifacts/runs/agent_smoke_001/generated/metadata/agent_traces/`
- `artifacts/runs/agent_smoke_001/evaluation/evaluation_result.json`
- `artifacts/runs/agent_smoke_001/report/report_summary.json`

重点看：

- `generated_manifest.json` 里有没有 3 个 subject
- `generated/tests/` 里有没有真正生成的测试
- `agent_traces/*.trace.jsonl` 里有没有 OpenCode 的 stderr/stdout
- `report_summary.json` 里有没有 `agent_comparisons` 和 `skill_uplifts`

**如果失败，优先看**

- `generated/metadata/*.response.json`
- `generated/metadata/agent_traces/*.trace.jsonl`
- `logs/`

最常见的问题会是：

- 外层没挂 Docker socket
- `utbench-agent-opencode:latest` 没构建
- API key 没透传
- 内层镜像里 `opencode` CLI 实际安装名或行为不一致
- 网络权限不对

我这轮新增的是文档和示例配置修正，没有再改 Go 逻辑；上一轮代码改动后的回归测试已经通过。你先按上面的 smoke run 跑一遍，把结果或报错贴回来，我再继续收敛 `OpenCode` 的真实适配。

> 439.9 E: Failed to fetch http://mirrors.aliyun.com/ubuntu/pool/main/libx/libxcomposite/libxcomposite1_0.4.5-1build3_amd64.deb  502  Bad Gateway [IP: 183.60.240.239 80]
> 439.9 E: Sub-process /usr/bin/dpkg returned an error code (1)
> ------
> Dockerfile:33
> --------------------
>   32 |     # System dependencies - Install LLVM 19 for Mull 19 compatibility
>   33 | >>> RUN apt-get install -y --fix-missing --no-install-recommends \
>   34 | >>>     build-essential \
>   35 | >>>     gcc \
>   36 | >>>     g++ \
>   37 | >>>     cmake \
>   38 | >>>     docker.io \
>   39 | >>>     git \
>   40 | >>>     unzip \
>   41 | >>>     python3 \
>   42 | >>>     python3-pip \
>   43 | >>>     python3-dev \
>   44 | >>>     python3-pytest \
>   45 | >>>     python3-six \
>   46 | >>>     python3-dateutil \
>   47 | >>>     python3-venv \
>   48 | >>>     openjdk-21-jdk \
>   49 | >>>     maven \
>   50 | >>>     libgtest-dev \
>   51 | >>>     gdb \
>   52 | >>>     sudo \
>   53 | >>>     clang-19 \
>   54 | >>>     llvm-19 \
>   55 | >>>     llvm-19-dev \
>   56 | >>>     libclang-19-dev \
>   57 | >>>     lld-19 \
>   58 | >>>     libstdc++6 \
>   59 | >>>     libgcc-s1 \
>   60 | >>>     && cd /usr/src/gtest && cmake . && make && cp lib/*.a /usr/lib/ \
>   61 | >>>     && rm -rf /var/lib/apt/lists/*
>   62 |
> --------------------
> ERROR: failed to build: failed to solve: process "/bin/sh -c apt-get install -y --fix-missing --no-install-recommends     build-essential     gcc     g++     cmake     docker.io     git     unzip     python3     python3-pip     python3-dev     python3-pytest     python3-six     python3-dateutil     python3-venv     openjdk-21-jdk     maven     libgtest-dev     gdb     sudo     clang-19     llvm-19     llvm-19-dev     libclang-19-dev     lld-19     libstdc++6     libgcc-s1     && cd /usr/src/gtest && cmake . && make && cp lib/*.a /usr/lib/     && rm -rf /var/lib/apt/lists/*" did not complete successfully: exit code: 100

<details><summary>2 previous messages</summary>

> 这是镜像源可用性问题，不是依赖本身有冲突。当前 Dockerfile 把 Ubuntu 源强行改成阿里云镜像，一旦镜像节点返回 `502`，整层就直接炸。我要把这段改成“优先官方源，按需允许自定义镜像”，避免把构建稳定性绑死在单一镜像站上。
>
> <details><summary>Edited 1 file</summary>
>
> - Wrote `./go-ut-bench/Dockerfile` (+22 -15)
>
> </details>
</details>

这个错误我已经处理了，问题根因是 Dockerfile 把 Ubuntu 源硬改成了阿里云，镜像站一旦 `502`，构建必挂。

我改了 [Dockerfile](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/Dockerfile:1) 的策略：

- 默认不再强制使用阿里云镜像
- 默认走 Ubuntu 官方源
- 如果你们确实需要国内镜像，现在改成了可选参数 `UBUNTU_MIRROR`

也就是说现在有两种构建方式。

直接用官方源构建：

```bash
docker build -t utbench:latest .
```

如果你想手动指定镜像源，再这样传：

```bash
docker build \
  --build-arg UBUNTU_MIRROR=http://mirrors.aliyun.com/ubuntu \
  -t utbench:latest .
```

如果阿里云继续不稳定，可以换成别的镜像，比如中科大或清华，形式一样：

```bash
docker build \
  --build-arg UBUNTU_MIRROR=https://mirrors.ustc.edu.cn/ubuntu \
  -t utbench:latest .
```

你现在先直接重跑一次不带 `UBUNTU_MIRROR` 的构建。
如果还卡在网络层，我再继续把 `apt` 这段做成更强的重试和 fallback。

> ========================================
> [STAGE] 生成测试
> ========================================
> 4 样本 × 3 被测对象 = 12 任务 | Workers: 16
>
> [1/12]  opencode__deepseek-v4-flash__no_skill | python | interface_mock_000 | FAIL(agent command failed: cat: /wo...) |
> [2/12]  opencode__deepseek-v4-flash__no_skill | python | boundary_000 | FAIL(agent command failed: cat: /wo...) |
> [3/12]  opencode__deepseek-v4-flash__no_skill | python | simple_function_000 | FAIL(agent command failed: cat: /wo...) |
> [4/12]  opencode__deepseek-v4-flash__unit_test_skill | python | boundary_000 | FAIL(agent command failed: cat: /wo...) |
> [5/12]  opencode__deepseek-v4-flash__no_skill | python | complex_dependency_000 | FAIL(agent command failed: cat: /wo...) |
> ========================================
> [STAT] 实时统计
>    已处理: 5/12 (41%)
>    阶段已耗时: 00:04
>    预计剩余: ~00:05
> ========================================
> [6/12]  opencode__deepseek-v4-flash__unit_test_skill | python | simple_function_000 | FAIL(agent command failed: cat: /wo...) |
> [7/12]  opencode__deepseek-v4-flash__unit_test_skill | python | complex_dependency_000 | FAIL(agent command failed: cat: /wo...) |
> [8/12]  opencode__deepseek-v4-flash__unit_test_skill | python | interface_mock_000 | FAIL(agent command failed: cat: /wo...) |
> [9/12]  model_api__deepseek-v4-flash__no_skill | python | boundary_000 | OK | | 2249 tokens
> [10/12]  model_api__deepseek-v4-flash__no_skill | python | simple_function_000 | OK | | 2409 tokens
> ========================================
> [STAT] 实时统计
>    已处理: 10/12 (83%)
>    阶段已耗时: 00:22
>    预计剩余: ~00:04
> ========================================
> [11/12]  model_api__deepseek-v4-flash__no_skill | python | interface_mock_000 | OK | | 2813 tokens
> [12/12]  model_api__deepseek-v4-flash__no_skill | python | complex_dependency_000 | OK | | 7532 tokens
> ========================================
> [STAT] 实时统计
>    已处理: 12/12 (100%)
>    阶段已耗时: 01:20
>    预计剩余: ~00:00
> ========================================
>
> ========================================
> [DONE] 生成测试完成
>    总计: 12 | 成功: 12 | 失败: 0
>    耗时: 01:20
> ========================================失败了,肯定有问题吧,记得做好日志和错误定位,方便每次出问题的时候能够很好的定位问题,看看现在是什么问题

<details><summary>12 previous messages</summary>

> 我先不猜，直接从这次 run 的 trace、response 和工作区产物里定位。你这条日志里最有价值的信息是 `cat: /wo...`，大概率是内层容器命令在读取 prompt 文件时路径不对或者 workspace 挂载没对上。
>
> <details><summary>Ran 3 commands</summary>
>
> - Ran `Get-ChildItem -Path "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\artifacts\runs" | Sort-Object LastWriteTime -Descending | Select-Object -First 3 FullName,LastWriteTime`
> - Ran `rg -n "cat: /wo|agent command failed|opencode__deepseek-v4-flash" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\artifacts\runs"`
> - Ran `Get-ChildItem -Recurse -File "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\artifacts\runs" | Where-Object { $_.Name -like '*.trace.jsonl' -or $_.Name -like '*.response.json' } | Sort-Object LastWriteTime -Descending | Select-Object -First 20 FullName,LastWriteTime`
>
> </details>
>
> 我定位到了，不是 OpenCode 本身先坏，而是 **DOOD 下内层容器挂载的 workspace 路径错了**。外层容器里看到的是 `/app/artifacts/...`，但内层 `docker run -v /app/artifacts/...:/workspace` 是交给宿主机 daemon 解析的，宿主机上并没有这个路径，所以内层容器拿到的是空目录，`/workspace/utbench_agent_prompt.md` 自然不存在。
> 另外还有一个次要问题：生成阶段最后那句 `成功: 12 | 失败: 0` 统计口径也有 bug，和前面的 `FAIL(...)` 不一致，这个我一并修。
>
> <details><summary>Ran 3 commands</summary>
>
> - Ran `rg -n "成功:|失败:|PrintStageDone|completedCount|successCount|progress" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\service.go"`
> - Ran `Get-Content "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\service.go" | Select-Object -Skip 120 -First 160`
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\obs\progress.go"`
>
> </details>
>
> 我直接把这两个问题一起修掉：
> 1. DOOD 下把外层容器路径映射到宿主机路径，不再把 `/app/artifacts/...` 直接拿去给宿主机 daemon 用。
> 2. 生成阶段的成功/失败统计按真实 `item.Success` 计算，不再把“任务跑完”误当成“生成成功”。
>
> <details><summary>Edited 3 files, ran 1 command</summary>
>
> - Ran `Get-Content -Raw "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\sandbox.go"`
> - Wrote `./go-ut-bench/internal/runner/sandbox.go` (+38 -2)
> - Wrote `./go-ut-bench/internal/runner/subjects.go` (+12 -11)
> - Wrote `./go-ut-bench/internal/runner/subjects.go` (+1 -1)
> - Wrote `./go-ut-bench/internal/runner/service.go` (+11 -2)
>
> </details>
>
> 我再补两个回归点：一个测 DOOD 路径映射，一个把文档里的运行命令补上 `UTBENCH_SANDBOX_HOST_OUTPUT_ROOT`，否则你下次还会踩同一个坑。
>
> <details><summary>Created 1 file, edited 2 files</summary>
>
> - Wrote `./go-ut-bench/internal/runner/sandbox_test.go` (+30 -0)
> - Wrote `./go-ut-bench/docs/AGENT_UPGRADE.md` (+6 -1)
> - Wrote `./go-ut-bench/docs/DOCKER_GUIDE.md` (+4 -0)
>
> </details>
>
> 我顺手把错误定位口径也收紧了：现在没有 host path 映射时，它会直接报一个明确的 DOOD 配置错误，不会再绕一圈变成 `cat /workspace/utbench_agent_prompt.md not found` 这种二次症状。
>
> <details><summary>Ran 2 commands</summary>
>
> - Ran `gofmt -w "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\sandbox.go" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\sandbox_test.go" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\subjects.go" "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\service.go"`
> - Ran `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/runner ./internal/agentconfig ./internal/reporter ./cmd/utbench`
>
> </details>
>
> 这是测试平台差异，不是逻辑问题。Windows 下 `filepath.Join` 会产出反斜杠路径，我把断言改成平台无关。
>
> <details><summary>Edited 1 file, ran 2 commands</summary>
>
> - Wrote `./go-ut-bench/internal/runner/sandbox_test.go` (+5 -2)
> - Ran `gofmt -w "C:\Users\wzd\Desktop\速通ing\腾讯mini(多模型单元测试生成效果横向评测)\ut-bench\go-ut-bench\internal\runner\sandbox_test.go"`
> - Ran `$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/runner ./internal/agentconfig ./internal/reporter ./cmd/utbench`
>
> </details>
</details>

问题已经定位清楚了，而且不是猜测，是从你这次 run 的产物里直接对出来的。

**这次实际问题**

核心不是 OpenCode 先挂了，而是 **DOOD 下内层容器挂载了错误的 workspace 路径**。

你现在是：

- 外层：`utbench:latest` 容器
- 内层：每个 `subject × sample` 的 Agent 容器
- 外层通过 Docker socket 控制宿主机 daemon

但之前代码里内层挂载用的是外层容器看到的路径：

```text
/app/artifacts/runs/.../agent_workspaces/...
```

这个路径对外层容器是对的，**但对宿主机 Docker daemon 不是宿主机真实路径**。于是内层容器拿到的是空目录，才会出现：

```text
cat: /workspace/utbench_agent_prompt.md: No such file or directory
Error: You must provide a message or a command
```

这个错误链我已经从 `evaluation_result.json` 和 `report_summary.json` 里确认了。

**我已经修了什么**

1. 修了 DOOD 路径映射
   新增了 `UTBENCH_SANDBOX_HOST_OUTPUT_ROOT` 支持，外层容器会把：
   - 容器内 `/app/artifacts/...`
   映射成：
   - 宿主机真实 `<host artifacts root>/...`

   相关代码：
   - [sandbox.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/sandbox.go)

2. 修了错误提示
   现在如果你走 DOOD，但是没传 host path 映射，不会再绕成 `cat /workspace/... not found`，而会直接报更明确的错误，大意是：
   - 这个 workspace 看起来是容器内路径
   - 需要设置 `UTBENCH_SANDBOX_HOST_OUTPUT_ROOT`

3. 修了生成阶段统计 bug
   之前最后那句：
   - `成功: 12 | 失败: 0`
   是错的，它把“任务完成”当成“生成成功”了。现在按 `item.Success` 真正统计。

   相关代码：
   - [service.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/service.go)

4. 文档也补了
   - [DOCKER_GUIDE.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/DOCKER_GUIDE.md)
   - [AGENT_UPGRADE.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/AGENT_UPGRADE.md)

**日志和定位这块现在怎么样**

这次其实已经能定位，只是错误提示不够一针见血。当前已经有这些定位材料：

- `generated/metadata/*.response.json`
- `generated/metadata/agent_traces/*.trace.jsonl`
- `generated/metadata/agent_traces/*.diff.json`
- `evaluation/evaluation_result.json`
- `report/report_summary.json`

这次就是靠这些把问题锁到 DOOD 挂载路径错了。现在我把“前置配置错误”也前移成显式报错了，后面定位会快很多。

**你现在怎么重跑**

先重新构建外层镜像：

```bash
docker build -t utbench:latest .
```

然后按这个方式跑，关键是 **多了 `UTBENCH_SANDBOX_HOST_OUTPUT_ROOT`**：

```bash
docker run --rm --env-file .env \
  -e UTBENCH_SANDBOX_HOST_OUTPUT_ROOT="$(pwd)/artifacts" \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -w /app \
  utbench:latest run \
    --models deepseek-v4-flash \
    --langs python \
    --config /app/configs/models.yaml \
    --agents-config /app/configs/agents.example.yaml \
    --subjects model_api__deepseek-v4-flash__no_skill,opencode__deepseek-v4-flash__no_skill,opencode__deepseek-v4-flash__unit_test_skill \
    --dataset-root /app/datasets \
    --output-root /app/artifacts \
    --class self_contained \
    --max-samples 1 \
    --run-id agent_smoke_002
```

如果你在 Windows PowerShell 原生环境跑，不要直接照搬这条 Bash 命令；优先用 WSL 跑这条。DOOD 这块在 WSL/Linux 下更稳，也更少路径翻译问题。

**我这边已经做过的回归**

```powershell
$env:GOCACHE=(Join-Path (Get-Location) '.gocache'); go test ./internal/runner ./internal/agentconfig ./internal/reporter ./cmd/utbench
```

通过。

你现在先跑一遍 `agent_smoke_002`。如果还失败，把新的：
- `evaluation_result.json`
- 对应 `trace.jsonl`
- 外层终端报错

贴回来，我继续收敛下一层问题。
