# UT-Bench 面试讲解手册

> 这份文档不是给用户看的快速上手文档，而是给项目作者自己做面试准备的“深度讲解稿”。
> 目标是：你不仅知道项目“做了什么”，还知道“为什么这样做”“核心难点在哪”“关键实现怎么落地”“如果面试官追问某个点，应该怎么展开回答”。

---

## 1. 我会先怎么用一句话介绍这个项目

### 1.1 最短版本

UT-Bench 是一个面向 **大模型和 Coding Agent 的单元测试生成评测平台**，支持多语言、多模型、多 Agent 框架的统一接入，并通过编译、测试通过率、覆盖率、变异测试、成本与轨迹等多个维度对生成能力做横向评测。

### 1.2 面试里更完整的版本

我做的这个项目本质上不是一个“调用 LLM 生成测试代码的小工具”，而是一个 **Agent Benchmark / Evaluation Platform**。  
它解决的问题是：如果我想公平地比较不同模型、不同 Agent 框架、不同 Skill 组合在“自动生成单元测试”这个任务上的真实效果，不能只看它有没有输出代码，而要看：

- 生成的测试能不能编译
- 能不能跑通
- 覆盖率如何
- 变异测试是否能真正杀死错误实现
- token 和成本是多少
- Agent 在过程中用了多少轮交互、多少工具、读写了多少文件
- 相同实验条件下能不能复用历史结果，避免重复花钱

所以 UT-Bench 的定位是：**把 LLM/Agent 测试生成做成一个可复现、可对比、可追踪、可复用的标准化评测系统**。

---

## 2. 项目为什么值得做

### 2.1 业务/研究上的真实问题

如果只靠人工体验去评估模型或 Agent 的测试生成能力，会有几个问题：

1. **不可比**
   - 每次问的问题不一样
   - 不同人给的 prompt 不一样
   - 不同语言、不同项目、不同样本难度不一致

2. **不完整**
   - 只看“生成出来了没有”意义很弱
   - 真正重要的是测试质量，而不是文本看起来像不像测试代码

3. **不可复现**
   - 没有固定数据集
   - 没有固定运行环境
   - 没有产物留档
   - 下次模型升级后无法严谨对比

4. **成本高**
   - 每次都重复跑一遍
   - 明明实验条件没变，还是重新花 token、重新跑评测

### 2.2 我的解决思路

我做这个项目时，核心思路不是“把一次评测跑通”，而是把整个事情拆成一套完整的工程体系：

- 数据集治理
- 被测对象抽象
- 生成阶段和评测阶段解耦
- Agent runtime 和 evaluator runtime 分离
- 结果产物结构化
- 历史资产可复用
- 报告和 Web 可视化

也就是说，我把这件事从“脚本自动化”推进成了“平台化评测”。

---

## 3. 项目整体架构怎么讲

### 3.1 先讲主线

整个系统的主线非常清晰，可以用一句话概括：

```text
样本发现 -> 测试生成 -> 编译/测试/覆盖率/变异评测 -> 报告聚合 -> 资产入库与复用
```

代码里这个主线是由 [internal/orchestrator/service.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/orchestrator/service.go) 统一编排的。

### 3.2 模块划分

我一般会把它拆成六层来讲：

1. **入口层**
   - CLI：`cmd/utbench/main.go`
   - Web：`internal/web/`

2. **编排层**
   - `internal/orchestrator/`
   - 负责组织 generate / evaluate / report 各阶段

3. **数据集层**
   - `internal/dataset/`
   - 负责样本发现、索引、校验、过滤

4. **生成层**
   - `internal/runner/`
   - 负责纯模型 API 调用和 CLI Agent 执行

5. **评测层**
   - `internal/evaluator/`
   - 负责编译、测试、覆盖率、变异测试、失败归因

6. **结果层**
   - `internal/reporter/`
   - `internal/store/`
   - 前者做聚合分析与报告，后者做 SQLite 资产存储与复用

### 3.3 如果面试官问“为什么这样分层”

你可以这样回答：

> 因为这个系统有两个天然解耦点。第一是“生成”和“评测”本身是不同阶段，生成阶段的输入是 prompt 和样本，输出是测试代码；评测阶段的输入是测试代码和样本，输出是 compile/test/coverage/mutation 等指标。第二是“业务控制面”和“执行运行时”也是不同东西，Web/CLI/资产管理是控制面，真正跑编译测试和 Agent 的是执行面。如果不先把这两组边界拆开，后面一接入多个 Agent、多语言、多运行环境，复杂度会指数增长。  

---

## 4. 核心抽象：为什么我定义了 subject

### 4.1 这个抽象解决什么问题

最开始如果只评测模型 API，维度很简单，就是“模型名”。  
但接入 Agent 后，比较对象不再只是模型，而是：

- 纯模型：`model_api + model`
- Agent 基线：`framework + model`
- Agent + skill：`framework + model + skill`

如果没有一个统一的抽象，后面的配置、调度、报告、数据库都会非常乱。

### 4.2 subject 的定义

我把统一被测对象定义成：

```text
subject = framework + model + optional skill
```

典型形式：

```text
model_api__deepseek-v4-flash__no_skill
opencode__deepseek-v4-flash__no_skill
opencode__deepseek-v4-flash__unit_test_skill
codebuddy__deepseek-v4-flash__no_skill
claudecode__deepseek-v4-flash__unit_test_skill
```

### 4.3 这个抽象的价值

这个抽象带来几个直接好处：

1. **调度统一**
   - 无论是纯模型还是 Agent，最终都统一成 `subject × sample`

2. **报告统一**
   - 比较维度不再是“某个模型”而是“一个完整实验对象”

3. **资产复用统一**
   - 复用是按 subject 及其精确版本来判断，而不是模糊按模型名判断

4. **未来扩展更自然**
   - 后面接入 `claudecode`、更多框架、更多 skill，不需要重写整体数据模型

### 4.4 subject 和 subject_version 的区别

这是很容易被问到的点。

#### subject

`subject` 是用户视角的身份，例如：

```text
opencode__deepseek-v4-flash__unit_test_skill
```

#### subject_version

`subject_version` 是复用与实验严谨性视角下的精确版本。  
它不仅看 `subject_id`，还会考虑：

- 模型配置
- framework 配置
- skill 内容
- agent 命令模板
- sandbox 指纹
- docker image / digest
- 环境契约指纹

这些字段已经进入资产层设计和结果结构里，见：

- [docs/asset-management-design.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/asset-management-design.md)
- [internal/contracts/results.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/contracts/results.go)

如果面试官问“为什么不直接用 subject_id 做缓存”，你可以答：

> 因为 `opencode + deepseek + unit_test_skill` 这个名字一样，不代表实验条件一样。只要 skill 文件内容、agent 命令模板、镜像版本、sandbox 策略、环境契约有变化，结果都不应该被当成同一个实验条件复用。所以我专门做了 `subject_version` 和各种 fingerprint。  

---

## 5. 样本和数据集是怎么设计的

### 5.1 当前支持两类样本

当前代码和文档中使用两个数据集类别：

- `self_contained`
- `repo_level`

不要再说旧名字 `module_level`，现在已经收敛到 `repo_level`。

### 5.2 self_contained 的含义

它面向单文件或轻量上下文的任务。  
一个样本通常是一段源代码，要求模型或 Agent 为它生成测试。

这类样本比较适合：

- 快速 benchmark
- 多语言横向对比
- 初始能力验证

### 5.3 repo_level 的含义

repo-level 是为了未来的真实工程任务准备的。  
它强调：

- workspace 级上下文
- target files
- build/test command
- dependency/lock files
- 更真实的仓库结构

这也是后续往 SWE-bench 那种更强任务契约演进的基础。

### 5.4 sample_uid 为什么重要

为了支持样本复用和稳定身份，我做了 `sample_uid` 概念。

它不是简单用文件名，而是结合：

- language
- dataset_class
- scenario
- sample_id
- 源码或 workspace 指纹
- repo manifest 指纹

这样做的好处是：

- 同名不同内容的样本不会混淆
- repo-level 未来可以自然扩展
- 资产复用有稳定样本身份

---

## 6. 生成阶段怎么实现

生成阶段的核心在 [internal/runner/](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner)。

### 6.1 为什么我把生成层单独抽出来

因为“生成测试代码”本身就有两条完全不同的执行路径：

1. **model_api**
   - 直接调模型接口

2. **cli_agent**
   - 启动 OpenCode / CodeBuddy / Claude Code 这样的 Agent CLI

它们表面上都产出测试代码，但运行方式、可观测信息、失败模式完全不一样。

### 6.2 统一的输出

不管是 model_api 还是 cli_agent，最终都要归一化成 `GeneratedCase`，也就是：

- 生成是否成功
- 生成的测试文件路径
- prompt 和 response 快照
- token / cost
- latency
- trace / workspace diff
- sandbox 信息
- 复用信息

这也是为什么 [internal/contracts/results.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/contracts/results.go) 很重要，它是整个系统的数据契约中心。

### 6.3 model_api 路径

`model_api` 是纯模型基线，主要价值是：

- 作为不带 Agent 工作流的 baseline
- 对比“Agent 外壳到底带来了什么增益或损失”

它通常直接读取模型配置：

- `api_endpoint`
- `anthropic_endpoint`
- `api_key_env`
- `model`
- pricing / parameters

来自 [configs/models.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/models.yaml)。

### 6.4 cli_agent 路径

这是项目的核心复杂度之一。

CLI Agent 通过配置驱动，而不是把每个 Agent 写死在代码里。  
当前已经接入：

- `opencode`
- `codebuddy`
- `claudecode`

配置集中在：

- [configs/agents.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.yaml)

每个 framework 都有：

- kind
- sandbox 配置
- preflight
- forbidden command patterns
- env / env_from_host
- command template
- compatible languages
- output globs

### 6.5 为什么用“命令模板 + 适配器”而不是每个 Agent 单独写一套 runner

因为很多差异其实只在三层：

1. 命令模板怎么写
2. skill 怎么注入
3. 输出 JSON / stream-json 怎么解析

如果每接一个 Agent 都单独从头写一套 runner，维护成本会很高。  
所以我做的是：

- 公共 runner 负责工作区准备、prompt 渲染、执行、超时、产物落盘
- framework config 负责描述“怎么跑”
- adapter 负责解析“跑出来的结构化输出”

### 6.6 skill 注入怎么做

项目支持三类 skill 注入思路：

- `prompt_append`
- `workspace_mount`
- `agent_native`

这是一个很值得讲的设计点。  
因为不同 Agent 的“原生 skill 机制”并不统一，如果一开始把 skill 和某个 Agent 绑死，后面扩展会很难。

所以我的做法是：

1. skill 在 UT-Bench 里先抽象成“能力包”
2. 再根据不同 Agent 的特性选择注入方式
3. 对支持原生 skill 的 Agent，再单独适配它的目录约定

例如：

- Claude Code 适配时，原生 skill 被映射到 `.claude/skills/<name>/SKILL.md`
- 这部分逻辑在 [internal/runner/subjects.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/subjects.go)

这样做的价值是：**skill 是平台能力，不是某个 Agent 的私有配置**。

### 6.7 agent 输出为什么难

Agent 和纯 API 最大的不同点在于：  
你要的不只是“最后生成的文本”，还要尽量还原过程。

当前系统已经在做这些采集：

- token / cost
- session id
- tool calls
- interaction count
- files read / write count
- command count
- trace path
- workspace diff path

例如 Claude Code 接入时，我就专门做了：

- `stream-json` 输出解析
- usage 解析
- session id 补全
- tool call 提取

这部分位于：

- [internal/runner/adapter_cli.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/runner/adapter_cli.go)

如果面试官问“为什么要解析 trace”，你可以答：

> 因为 Agent 的价值不只在最终答案，还在工作过程。对于纯模型，token/cost 可能够了；但对 Agent，很多时候需要知道它做了多少轮交互、用了多少工具、有没有乱读写、是不是靠大量试错才得到结果。这些过程指标会直接影响我们对 Agent 工程效率和可靠性的判断。  

---

## 7. 为什么我专门做了 Agent 沙箱

这是面试里非常值得讲的点，因为它能体现你不是在“调个命令行工具”，而是在做可靠评测基础设施。

### 7.1 为什么需要沙箱

如果让 Agent 直接在评测机目录里乱跑，会有几个严重问题：

1. 污染工作区
2. 偷装依赖导致实验不公平
3. 不同样本之间相互影响
4. 失败难归因
5. 安全边界太弱

### 7.2 当前沙箱设计

当前的设计是：

- 外层 UT-Bench 负责控制
- 对每个 `subject × sample` 建独立工作区
- Agent 在独立沙箱中执行
- 可以按 framework 配置 `sandbox.provider` / `sandbox.mode` / `image`

当前主流模式是 Docker 沙箱：

- 一个样本一个隔离工作区
- 一个 subject × sample 一个沙箱执行实例
- preflight 先检查语言工具链
- 可限制网络
- 可限制 CPU / memory
- 拦截一类环境漂移命令，例如 `apt-get install`、`pip install`、`npm install`

配置可见：

- [configs/agents.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/agents.yaml)

### 7.3 为什么要把 sandbox provider 做成配置块

这不是为了“配置好看”，而是为了支持未来 runtime 抽象。

当前已经不是只有一个胖 Docker 镜像在跑所有东西了，而是开始明确三层：

- control plane
- evaluation plane
- agent sandbox plane

这套设计和思考已经记录在：

- [docs/runtime-sandbox-redesign.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/runtime-sandbox-redesign.md)

如果面试官问“为什么不让评测 Docker 和 Agent Docker 合并”，你可以这样答：

> 因为评测环境和 Agent 环境的职责不同。评测环境强调 compile/test/coverage/mutation 的稳定性；Agent 环境强调 CLI agent 可运行、被隔离、可追踪。两者生命周期不同，依赖不同，安全边界也不同。如果混成一个大镜像，短期方便，长期会很难治理，也很难保证实验条件清晰。  

---

## 8. 评测阶段怎么实现

评测阶段的核心在 [internal/evaluator/](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/evaluator)。

### 8.1 为什么评测阶段是系统可信度的核心

很多类似项目只做到“让模型输出一段测试代码”，这其实不能证明什么。  
真正有价值的是：这段代码能不能在标准环境下通过真实验证。

所以评测阶段是整个系统最决定结论可信度的一层。

### 8.2 多维指标

当前主要指标包括：

1. **CompilePass**
   - 测试代码和目标代码能不能编译/构建通过

2. **TestPass**
   - 测试能不能执行成功

3. **Coverage**
   - line coverage
   - branch coverage（有语言支持时）

4. **Mutation Score**
   - 不是看代码跑通，而是看测试能不能真正发现错误实现

5. **结构性指标**
   - assertion count
   - test case count
   - assertion density

6. **效率指标**
   - runtime
   - token
   - cost

### 8.3 为什么我加了变异测试

这是很容易加分的点。

只看覆盖率有一个问题：  
覆盖到了，不代表测试真的有效。

变异测试更接近测试质量本身。  
如果把源码做一些小修改，测试仍然全过，说明测试可能只是“碰到了代码”，没有真正约束行为。

所以我在各语言上尽量接入 mutation tool：

- Python：`mutmut`
- Go：`go-mutesting`
- Java：`pitest`
- C++：`mull`

这些信息也写在 [configs/models.yaml](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/configs/models.yaml) 的 language 配置中。

### 8.4 失败归因为什么重要

这是你项目成熟度的一个关键体现。

如果一个 case 失败了，不代表就是模型差。失败可能来自：

- 模型能力问题
- Agent 工作流问题
- 样本本身问题
- 环境问题
- 工具链问题

所以结果结构里我加入了：

- `failure_origin`
- `score_eligible`
- `score_exclusion_reason`

也就是说，我在尽量把“能力失败”和“环境失败”区分开，避免污染模型排名。

如果面试官问“为什么这么麻烦”，你可以答：

> 因为 benchmark 最容易犯的错误就是把环境问题误判成模型问题。只要评测环境不稳定、工具链有问题、样本有缺陷，最后统计出来的通过率就会失真。所以我在结果结构和评测逻辑上专门做了失败归因和计分资格控制。  

---

## 9. 报告层怎么设计

### 9.1 报告不是简单汇总

报告层的目标不是“把几个数字打印出来”，而是把一次 run 提炼成可决策的信息。

所以我在报告层做了几类聚合：

1. 基础汇总
   - compile/test/coverage/mutation 平均表现

2. 排名维度
   - 模型/subject 的综合对比

3. 效率维度
   - token efficiency
   - time efficiency
   - cost estimate

4. 错误诊断
   - compile errors
   - test errors
   - mutation errors

5. 自动洞察
   - best model
   - weak scenarios
   - language gaps
   - recommendations

6. Agent 维度分析
   - agent comparisons
   - skill uplifts
   - runtime summary

### 9.2 runtime summary 为什么要进报告

这是后来平台化过程中补的重要能力。

报告里不应该只看到“谁分高”，还应该看到“是在什么运行时条件下跑出来的”。

因此报告里加入了：

- evaluator environment fingerprint
- sandbox providers
- sandbox images
- agent frameworks
- docker/local backed subject counts

这部分来自：

- [internal/reporter/runtime_summary.go](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/reporter/runtime_summary.go)

这能帮助跨 run 做严谨对比。

---

## 10. 资产管理和结果复用是怎么做的

这是你项目里非常有“平台味”的一块，也是很适合在面试里展开的。

### 10.1 为什么要做复用

这个系统很贵，贵在两件事：

1. 生成阶段会花 token 和时间
2. 评测阶段会花编译/测试/变异测试成本

如果完全相同实验条件下每次都重跑，是明显浪费。

### 10.2 我的复用不是“按名字碰运气”

我没有做那种简单的“同模型同样本就直接复用”，而是设计了完整资产模型：

- `subject`
- `subject_version`
- `sample_uid`
- `generation_key`
- `evaluation_key`

文档在：

- [docs/asset-management-design.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/asset-management-design.md)

### 10.3 generation_key 关注什么

生成复用要保证这些条件一致：

- subject_version
- sample_uid
- prompt rendering
- prompt version
- language
- dataset class
- dependency fingerprint
- generation env fingerprint

也就是说：**只要真正影响生成结果的条件变了，就不应该复用。**

### 10.4 evaluation_key 关注什么

评测复用要关注：

- generated test 内容
- sample_uid
- evaluation environment fingerprint
- evaluator version
- mutation config
- dependency fingerprint

评测复用的默认策略更保守，因为评测环境比生成更容易受工具链变化影响。

### 10.5 为什么 SQLite + 文件系统混合存储

这是很典型的工程权衡。

#### 文件系统保存大产物

例如：

- generated test
- prompt snapshot
- model response
- agent trace
- workspace diff
- report html

这些适合落文件，方便查看和调试。

#### SQLite 保存索引和指纹

例如：

- subject
- subject_version
- generation_key
- evaluation_key
- artifact metadata
- run 和 asset 的关联

这样做的好处是：

- 查询快
- 结构化好管理
- 复用逻辑清晰
- 产物仍然可直接落地查看

---

## 11. Web 管理界面怎么讲

### 11.1 Web 不只是“好看”

Web 的价值不是把 CLI 换个壳，而是把这个平台真正变成可操作的管理系统：

- 发起任务
- 查看环境
- 管理 API Key
- 浏览历史 run
- 看报告
- 看资产
- 查复用命中情况

### 11.2 Web 当前做了哪些能力

结合你当前的代码状态，Web 侧已经覆盖这些：

- 运行任务
- 查看任务详情和日志
- 构建镜像
- 环境检查
- subject / registry / assets 管理视图
- API key 管理
- 运行时拓扑状态展示

核心目录：

- [internal/web/](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/internal/web)

### 11.3 Web 和 Docker/runtime 的关系

这是当前项目很重要的一个工程点。

我不是把 Web 当静态页面，而是让它参与运行时管理：

- 感知 Docker 是否可用
- 区分 control image / eval image
- 支持镜像构建和环境检查

这背后其实体现的是：**项目已经从工具演进到平台**。

---

## 12. Docker 和运行时重构这块怎么讲

### 12.1 一开始的问题

项目最早只有一个“大而全”的评测镜像思路，但随着 Agent 接入，这个思路越来越不够用。

因为系统里其实有三种完全不同的职责：

1. control plane
2. evaluation plane
3. agent sandbox plane

### 12.2 我现在的方向

我已经在往这三个运行时层次拆：

- `utbench-control`
- `utbench-eval`
- `utbench-agent-*`

这套设计的文档基础是：

- [docs/runtime-sandbox-redesign.md](/C:/Users/wzd/Desktop/速通ing/腾讯mini(多模型单元测试生成效果横向评测)/ut-bench/go-ut-bench/docs/runtime-sandbox-redesign.md)

### 12.3 为什么这件事重要

因为如果控制面、评测面、Agent 沙箱面不分开，会有这些问题：

- 镜像过胖
- 重建慢
- 依赖边界不清
- 环境污染难排查
- 复用判断不严谨
- 将来接入更强隔离方案很难

如果面试官问“你为什么觉得自己项目难点在 runtime，而不是 prompt”，你可以答：

> prompt 很重要，但当系统从单模型脚本发展成多 Agent benchmark 平台后，决定可信度和可扩展性的其实是 runtime contract。只有把 control、evaluation、sandbox 的边界钉死，后面模型、Agent、skill 的接入才不会把系统拖垮。  

---

## 13. 你接入 OpenCode、CodeBuddy、Claude Code 时到底做了什么

这是很适合讲“工程抽象能力”的地方。

### 13.1 OpenCode

OpenCode 更适合做“多模型外壳”的 Agent 入口，因为它更容易对接 OpenAI-compatible provider。  
我接入它时，重点做了：

- command template
- config 注入
- sandbox preflight
- trace / usage / tool parsing
- skill 注入

### 13.2 CodeBuddy

CodeBuddy 的难点在于它对自定义模型配置和 stream-json 输出有自己的约定。  
我做的是把它纳入同一个 CLI agent 框架，而不是单独绕一套逻辑。

### 13.3 Claude Code

Claude Code 的难点不只是“把命令跑起来”，而是它的世界观和 OpenCode/CodeBuddy 不完全一样：

- 有自己的原生 skill 目录约定
- 有自己的 stream-json 输出格式
- 有 Anthropic / Bedrock / Vertex / Gateway 等 provider 模式

所以我做 Claude Code 适配时，主要处理了：

1. framework 配置
2. `.claude/skills/<name>/SKILL.md` 原生 skill 注入
3. `stream-json` 输出解析
4. usage / tool call / session id 采集
5. 统一 agent 沙箱镜像安装 Claude Code CLI

如果面试官问“为什么不把 Claude Code 当成普通命令直接跑”，你可以答：

> 因为真正的工作不是把命令跑起来，而是把它纳入统一 benchmark contract：能调度、能隔离、能采 usage、能解析 tool calls、能注入 skill、能写入统一结果结构。只有这样它才是平台里的一个 framework，而不是一个特例脚本。  

---

## 14. 这个项目里你觉得最难的几个点

你可以明确说三个。

### 14.1 难点一：统一抽象

纯模型 API 和 CLI Agent 看起来都是“生成测试”，但它们的数据、运行时、失败模式差异很大。  
最难的是做出统一抽象，又不把重要细节丢掉。

### 14.2 难点二：运行时隔离和环境治理

只要接入真实 Agent，就必须处理：

- 独立工作区
- 镜像选择
- preflight
- 环境漂移命令
- Docker 嵌套
- 宿主机路径映射
- 环境 fingerprint

这部分复杂度远高于“写个 prompt 调 API”。

### 14.3 难点三：结果可信性

benchmark 最大的问题不是功能实现，而是结论是否可信。  
所以我花了很多力气在：

- 失败归因
- score_eligible
- environment fingerprint
- asset reuse keys
- runtime summary

这部分工作对外看不显眼，但它决定项目质量。

---

## 15. 如果面试官问“你的技术亮点是什么”

可以从四个角度说。

### 15.1 平台化抽象

我没有把它做成一个“评测脚本集合”，而是做成了统一的 benchmark 平台：

- subject
- skill
- runner
- evaluator
- reporter
- store

### 15.2 评测深度

不只评估“是否生成了测试”，而是覆盖：

- compile
- test
- coverage
- mutation
- cost
- trace

### 15.3 运行时工程

我把 Agent 评测里最麻烦的一层，也就是 sandbox/runtime 设计，真正往平台能力推进了。

### 15.4 资产复用

我引入 `subject_version / generation_key / evaluation_key`，让系统能在保证实验严谨性的前提下复用历史结果。

---

## 16. 项目的不足和下一步怎么讲

这个问题面试官非常可能问。

### 16.1 不要说“没什么问题”

正确说法是：我知道目前系统已经到了一个必须继续收架构账的阶段。

### 16.2 当前不足

可以坦诚说这些：

1. **repo-level 样本能力还在增强**
   - 当前单文件 benchmark 更成熟
   - 仓库级任务治理还会继续补强

2. **runtime 拆分还在进行中**
   - control / eval / agent sandbox 已经有设计和部分落地
   - 但 evaluator backend 还可以继续抽象

3. **Web 还是在持续重构**
   - 目前已具备平台能力
   - 但前端结构还可以继续模块化

4. **多 provider / 多 framework 的适配还在扩展**
   - Claude Code 已进入适配路径
   - 更完整的多 provider 验证还需继续做

### 16.3 下一步方向

你可以这样说：

> 下一步我会继续把 repo-level benchmark 做硬，把 evaluator 和 sandbox runtime contract 进一步收敛，同时完善更细的 trace 观测和复用规划，这样这个系统就会从“能用的评测平台”进一步变成“可信的 benchmark 基础设施”。  

---

## 17. 一段 2 分钟面试讲解稿

如果面试官让你快速介绍项目，你可以直接按这个说：

> 我做了一个叫 UT-Bench 的项目，它是一个面向大模型和 Coding Agent 的单元测试生成评测平台。它不是简单地调模型生成测试代码，而是把整个评测过程做成了标准化流水线：先从数据集中发现样本，然后统一调度被测对象生成测试，再做编译、测试执行、覆盖率和变异测试，最后生成报告并把结果入库。  
>
> 这个项目里我做的关键设计有三个。第一，我定义了 `subject = framework + model + skill` 这个统一抽象，把纯模型 API、OpenCode、CodeBuddy、Claude Code 这类不同被测对象统一到一个调度和报告体系里。第二，我把运行时分成控制面、评测执行面和 Agent 沙箱面，专门处理了 Docker 沙箱、preflight、环境指纹、失败归因这些工程问题，保证 benchmark 结果更可信。第三，我设计了 SQLite + 文件系统的资产管理层，引入 `subject_version`、`generation_key`、`evaluation_key` 来做历史结果复用，避免相同实验条件下重复花 token 和重复评测。  
>
> 我觉得这个项目最有价值的地方不是“能生成测试”，而是把 LLM/Agent 测试生成做成了一套可复现、可比较、可追踪、可扩展的评测基础设施。  

---

## 18. 面试常见追问和推荐回答

### Q1：为什么用 Go 做这个项目？

推荐答法：

> 这个项目本质上是一个多阶段编排系统，需要比较强的工程组织能力、并发控制能力、二进制部署能力和跨平台 CLI 能力。Go 在这些点上很合适：一方面可以很自然地做 worker pool、文件系统管理、子进程编排、HTTP/Web 服务；另一方面单二进制部署和 Docker 化也比较方便。  

### Q2：为什么不用 Python？评测工具很多都是 Python 的。

推荐答法：

> 评测工具链里确实有不少 Python 组件，比如 mutmut，但平台本身和工具链本身不是一回事。我希望核心编排层和资产管理层是一个稳定、容易部署、并发和工程结构更清晰的系统，所以平台用 Go，而不是让平台层跟工具脚本层耦合在一起。  

### Q3：为什么不用只看覆盖率？

推荐答法：

> 覆盖率只能说明“跑到了”，不能说明“测得好”。所以我引入变异测试作为更强的质量指标。变异测试能更好反映测试是否真正约束了行为，这是我认为 benchmark 里很关键的一点。  

### Q4：为什么要做 SQLite 资产库？

推荐答法：

> 因为这个项目的成本很高，尤其是生成阶段的 token 成本和评测阶段的运行成本。只靠文件目录做历史管理会很难支持复用、筛选和 explain-reuse，所以我做了 SQLite 存索引和指纹、文件系统存大产物的混合方案。  

### Q5：你怎么保证实验公平？

推荐答法：

> 我主要从四个方面控制：统一数据集和样本过滤、统一 prompt/skill 策略、统一评测环境和沙箱策略、以及通过 subject_version 和环境指纹严格区分实验条件。这样尽量避免把环境漂移、依赖变化或配置差异误当成模型能力差异。  

### Q6：你觉得这个项目最体现你能力的地方是什么？

推荐答法：

> 我觉得不是某一个功能点，而是我把一个容易做成脚本集合的事情，推进成了一套有统一抽象、运行时治理、资产复用和结果可视化的评测平台。这里面既有架构设计，也有运行时细节，也有指标和可信性设计。  

---

## 19. 你自己准备面试时应该怎么用这份文档

建议按下面顺序准备，不要死背。

### 第一轮：讲主线

先把这四句话讲顺：

1. 项目是干什么的
2. 为什么值得做
3. 整体流水线是什么
4. 你做了哪些核心设计

### 第二轮：讲三个难点

你至少要能顺畅讲出：

1. subject 抽象
2. sandbox/runtime 设计
3. asset reuse 和结果可信性

### 第三轮：讲一个具体实现细节

建议挑你最熟的一块，比如：

- Claude Code 适配
- 变异测试接入
- generation/evaluation key 设计
- Web + Docker runtime 拆分

这类细节能体现你不是泛泛而谈。

### 第四轮：讲项目不足和下一步

面试里如果你能清楚讲出：

- 哪些是现在已经做好的
- 哪些是明确还没做完的
- 你为什么这样排优先级

会比只讲功能更有说服力。

---

## 20. 最后给你一个“答题原则”

讲这个项目时，不要把它说成：

- 一个调用模型生成测试代码的工具
- 一个 Web 页面
- 一个自动化脚本集合

更准确的说法应该是：

> 这是一个面向 LLM 和 Coding Agent 的测试生成 benchmark 平台，我重点解决的是统一抽象、运行时治理、结果可信性和历史资产复用问题。  

这句话一旦立住，后面的所有细节才会有逻辑支撑。

