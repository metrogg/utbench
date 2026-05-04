# UT-Bench Architecture Optimization: Insights from Anthropic's "Demystifying Evals for AI Agents"

> 分析时间：2026-05-04
> 来源：https://www.anthropic.com/engineering/demystifying-evals-for-ai-agents
> 对照文档：UTBENCH_GAP_ANALYSIS.md + research_agent_benchmark_architecture.md + MEMORY.md 3D Architecture

---

## 一、Anthropic 文章核心洞察摘要

### 1.1 评估的基本结构

| 概念 | 定义 | UT-Bench对应 |
|:---|:---|:---|
| **Task** | 单个测试，有明确输入和成功标准 | Sample (self_contained / repo_level) |
| **Trial** | 对Task的一次尝试；因非确定性需多次trial | 当前仅单次运行 |
| **Grader** | 对agent输出某方面的评分逻辑；一个Task可有多个Grader | 当前仅Code Grader(测试执行) |
| **Transcript** | Trial的完整记录——输出、工具调用、推理、中间结果 | AgentTrace (部分实现) |
| **Outcome** | Trial结束时环境中的最终状态 | EvaluationResult |
| **Evaluation harness** | 端到端运行eval的基础设施 | UT-Bench runner |
| **Evaluation suite** | 一组设计用于衡量特定能力的Task集合 | Dataset (per class/level) |

### 1.2 三种Grader类型

| 类型 | 方法 | 优势 | 劣势 |
|:---|:---|:---|:---|
| **Code-Based** | 字符串匹配、二值测试(FAIL_TO_PASS/PASS_TO_PASS)、静态分析、Outcome验证、工具调用验证、Transcript分析 | 快、客观、可复现 | 对主观任务不足 |
| **Model-Based** | Rubric评分、自然语言断言、成对比较、参考解评估、多judge共识 | 灵活、可扩展、捕获细微差别 | 非确定性、需校准 |
| **Human** | 领域专家评审、众包判断、抽检校准、A/B测试 | 黄金标准 | 昂贵、慢 |

### 1.3 Anthropic 十大关键方法论

1. **Grade outcomes, not paths** — 评估结果，不惩罚创造性路径
2. **Build in partial credit** — 部分成功比完全失败有意义
3. **Balanced problem sets** — 正/负样本平衡，避免单侧优化
4. **Trial isolation** — 每次trial从干净环境开始
5. **pass@k / pass^k** — 区分"至少一次成功"和"每次都成功"
6. **Multi-grader composition** — Code+LLM+Human多grader组合
7. **Eval saturation monitoring** — 100%通过率=无改进信号
8. **Reference solution** — 每个task必须有已知正确解
9. **Swiss cheese model** — 多种评估方法叠加互补
10. **Capability vs Regression evals** — 能力评估vs回归守卫，task可"毕业"

---

## 二、UT-Bench 现有架构 vs Anthropic 最佳实践 差距分析

### 2.1 差距总览

| Anthropic方法论 | UT-Bench现状 | 差距等级 | 影响 |
|:---|:---|:---|:---|
| Grade outcomes, not paths | 编译+测试是outcome，但缺transcript质量评估 | 🟡 部分 | 无法区分"好的过程+差的结果"和"差的过程+差的结果" |
| Build in partial credit | CompositeScore是加权求和，无partial credit语义 | 🔴 缺失 | 一个测试编译通过但运行失败的样本=0分，但实际比完全不编译的强 |
| Balanced problem sets | 无难度分级，无正负样本平衡 | 🔴 缺失 | 可能导致某个模型在某类样本上过拟合 |
| Trial isolation | Docker沙箱隔离，但无自动回滚和状态校验 | 🟡 部分 | 可能出现共享状态导致的相关性失败 |
| pass@k / pass^k | 仅单次pass_rate，无多次trial统计 | 🔴 缺失 | 无法区分"偶尔能过"和"稳定能过" |
| Multi-grader composition | 仅有Code grader(测试执行)，缺LLM+Human | 🔴 缺失 | 无法评估测试质量、代码风格等主观维度 |
| Eval saturation monitoring | 无饱和度监控，无"毕业"机制 | 🔴 缺失 | 评测集可能逐渐失去区分度 |
| Reference solution | self_contained有源码，但无gold test | 🟡 部分 | 无法验证eval基础设施本身是否正确 |
| Swiss cheese model | 仅有自动化eval | 🔴 缺失 | 无法发现自动化eval遗漏的问题 |
| Capability vs Regression evals | 无区分，所有task用同一套 | 🔴 缺失 | 高通过率的task不提供改进信号 |

### 2.2 最关键差距排序

| 优先级 | 差距 | 根因 | 优化方向 |
|:---|:---|:---|:---|
| **P0** | Multi-Grader Composition缺失 | 原架构仅设计了Metric Interface，未分层Grader类型 | 新增Layer 4: Grader层 |
| **P0** | pass@k / pass^k缺失 | 当前设计仅单次运行 | Layer 3 EventStream增加trial维度 |
| **P0** | Partial Credit缺失 | CompositeScore是binary加权 | Grader层增加partial credit语义 |
| **P1** | Capability vs Regression未区分 | 所有task一视同仁 | 新增Layer 6: Lifecycle层 |
| **P1** | Balanced Dataset缺失 | 数据集设计阶段未考虑 | Layer 2 Harness增加难度分级 |
| **P2** | Swiss Cheese缺失 | 超出自动化eval范围 | 架构预留Human-in-the-loop接口 |

---

## 三、优化方案：5层 → 6层架构

### 3.1 原架构 (5层)

```
Layer 1: Pluggable Sandbox
Layer 2: Configurable Harness
Layer 3: Event-Sourced Engine
Layer 4: Composable Metrics
Layer 5: 3D Dimension Aggregation
```

### 3.2 优化后架构 (6层)

```
Layer 1: Pluggable Sandbox (增强: Trial隔离校验)
Layer 2: Configurable Harness (增强: Reference Solution + Balanced Dataset)
Layer 3: Event-Sourced Engine (增强: Outcome校验 + trial维度)
Layer 4: Multi-Grader Composition (NEW: 从原Metrics层拆分独立)
Layer 5: 3D Dimension Aggregation (增强: pass@k/pass^k)
Layer 6: Eval Lifecycle & Saturation (NEW)
```

### 3.3 各层详细优化

---

#### Layer 1: Pluggable Sandbox (增强)

**Anthropic洞察**：
> "Each trial should be isolated from a clean environment. Shared state between runs can cause correlated failures or artificially inflate performance."

**优化点**：

```go
type SandboxRunRequest struct {
    // ... 原有字段 ...
    
    // 新增: Trial隔离
    TrialIndex     int       // 当前trial编号(从1开始)
    TotalTrials    int       // 总trial数(默认1)
    CleanSnapshot  string    // 干净环境快照ID,每次trial前恢复
    
    // 新增: Outcome校验点
    PreRunHash     string    // 运行前workspace文件哈希(验证隔离)
}

// 新增: Trial隔离验证
func (r *SandboxRunner) VerifyIsolation(req SandboxRunRequest) error {
    // 验证workspace状态是否恢复到CleanSnapshot
    // 如果不是,说明上一个trial残留了状态
}
```

**实现要点**：
- 每个trial前执行 `git checkout . && git clean -fd` 确保干净状态
- 记录运行前workspace哈希，运行后对比，验证隔离性
- 源码只读挂载 `-v {source}:/workspace/src:ro`

---

#### Layer 2: Configurable Harness (增强)

**Anthropic洞察**：
> "Create a reference solution: a known working output that passes all graders, proving the task is solvable."
> "Build balanced problem sets — test both cases where a behavior should occur and where it shouldn't."

**优化点**：

```yaml
# 新增: Task定义增强
task:
  id: "python_sort_001"
  difficulty: "medium"  # easy/medium/hard/expert
  category: "algorithm"
  
  # NEW: Reference solution
  reference_solution: "tests/test_sort_gold.py"
  reference_passes: true  # gold test是否通过(验证eval本身)
  
  # NEW: Balanced design
  expected_pass_rate: 0.4  # 预期通过率(用于饱和度监控)
  negative_cases:  # 不应该出现的行为
    - "should_not_modify_source_code"
    - "should_not_use_time_sleep_in_test"
  
  # NEW: 双重验证契约(借鉴SWE-bench)
  must_pass_existing: true
  existing_tests_pattern: "tests/test_original.py"
```

**实现要点**：
- 每个self_contained样本增加 `reference_solution` 字段
- Harness启动时先跑reference solution验证eval基础设施
- 数据集按难度分级: easy(单函数) / medium(单类) / hard(跨模块) / expert(仓库级)
- 正负样本平衡: 50%生成测试应通过 + 50%生成测试不应通过的场景

---

#### Layer 3: Event-Sourced Engine (增强)

**Anthropic洞察**：
> "Grade what the agent produced, not the path it took."
> "Outcome = the final state in the environment at the end of a trial."

**优化点**：

```go
// 增强的TaggedEvent
type TaggedEvent struct {
    ID        string
    Timestamp time.Time
    Type      string  // "action" | "observation" | "outcome"
    Cause     *string
    
    // 3D维度标签
    AgentFramework string
    Model          string
    Skill          string
    
    // NEW: Trial维度
    TrialIndex     int
    
    // NEW: Outcome类型事件
    Outcome        *OutcomeEvent  // nil for action/observation
    Data           map[string]any
}

// NEW: Outcome事件(区分过程和结果)
type OutcomeEvent struct {
    Type           string   // "compile" | "test_run" | "coverage" | "mutation"
    Success        bool     // 是否成功
    PartialCredit  float64  // 0.0~1.0 部分得分
    Details        map[string]any
}

// NEW: Transcript质量事件(不参与主评分,但记录在案)
type TranscriptEvent struct {
    StepCount      int
    RedundantSteps int     // 冗余步骤数
    ToolUsageScore float64 // 工具使用效率
    StrategyScore  float64 // 策略合理性(LLM评)
}
```

**关键设计决策**：
- **Outcome vs Process分离**：主评分只看Outcome(test pass/coverage/mutation)，Transcript质量作为附加维度
- **Partial Credit在Outcome层**：编译通过=0.3, 测试部分通过=0.5, 全部通过=1.0
- **Trial维度**：EventStream支持同一(Agent,Model,Skill)三元组的多次trial

---

#### Layer 4: Multi-Grader Composition (NEW — 从原Metrics层拆分)

**Anthropic洞察**：
> "Choose deterministic graders where possible, LLM graders where necessary, human graders judiciously."
> "Use structured rubrics grading each dimension with an isolated LLM judge."

**这是最大的架构改动**。原来的Layer 4"Composable Metrics"被拆分为两层：
- **Layer 4: Grader层** — 定义"谁来评"和"怎么评"
- **Layer 5: Aggregation层** — 定义"怎么聚合"(3D交叉表)

```go
// Grader Interface (统一评分契约)
type Grader interface {
    Name() string
    Type() GraderType  // Code / LLM / Human
    Grade(ctx context.Context, input GraderInput) GraderResult
}

type GraderType string
const (
    GraderCode   GraderType = "code"    // 确定性
    GraderLLM    GraderType = "llm"     // 模型评分
    GraderHuman  GraderType = "human"   // 人工评分
)

type GraderInput struct {
    SampleID       string
    GeneratedCode  string
    Transcript     []TaggedEvent
    Outcome        OutcomeEvent
    ReferenceSolution string  // 可选
}

type GraderResult struct {
    Score         float64  // 0.0~1.0
    PartialCredit float64  // 0.0~1.0 (部分得分)
    Reason        string
    GraderType    GraderType
    Assertions    []AssertionResult  // 细粒度断言结果
}

type AssertionResult struct {
    Name    string
    Passed  bool
    Weight  float64
    Detail  string
}
```

**三种Grader的具体实现**：

```go
// 1. Code Grader (确定性)
type CodeGrader struct {
    name       string
    assertions []CodeAssertion
}

type CodeAssertion struct {
    Name     string          // "compile_pass" | "test_pass" | "coverage_70" | "must_pass_existing"
    Type     string          // "binary" | "threshold"
    Command  string          // 执行命令
    Weight   float64         // 权重
    Partial  PartialCreditFn // 部分得分函数
}

// 例: 测试通过率的partial credit
// 5个测试通过3个 → partial credit = 0.6
func TestPassPartialCredit(passed, total int) float64 {
    return float64(passed) / float64(total)
}

// 2. LLM Grader (模型评分)
type LLMGrader struct {
    name       string
    rubric     string          // 评分准则
    model      string          // 评分用的模型
    dimensions []RubricDimension
}

type RubricDimension struct {
    Name        string  // "test_quality" | "code_style" | "assertion_meaningfulness"
    Description string
    Weight      float64
    MaxScore    float64
}

// 3. Human Grader (人工评分 — 校准用)
type HumanGrader struct {
    name       string
    dimensions []RubricDimension
    reviewers  []Reviewer
}

// Human Grader主要用于校准LLM Grader
// 定期抽样 → 人工评分 → 与LLM Grader对比 → 更新rubric
```

**Grader组合策略**（借鉴Anthropic）：

```yaml
# Task级别配置Grader组合
task:
  graders:
    - type: code
      required: true
      assertions:
        - name: compile_pass
          weight: 0.15
        - name: test_pass
          weight: 0.25
          partial_credit: true  # 启用partial credit
        - name: must_pass_existing
          weight: 0.10
          veto: true  # 一票否决: 破坏已有测试直接0分
        - name: coverage_threshold
          weight: 0.15
          threshold: 0.7
          partial_credit: true
        - name: mutation_score
          weight: 0.15
          threshold: 0.5
          partial_credit: true
    
    - type: llm
      required: false  # 可选,不参与主评分
      dimensions:
        - name: test_quality
          rubric: "prompts/test_quality_rubric.md"
          weight: 0.10
        - name: code_style
          rubric: "prompts/code_style_rubric.md"
          weight: 0.05
        - name: assertion_meaningfulness
          rubric: "prompts/assertion_rubric.md"
          weight: 0.05
    
    - type: human
      required: false  # 校准用
      sample_rate: 0.05  # 抽检5%
      dimensions:
        - name: overall_quality
          weight: 1.0
  
  # 组合策略
  scoring:
    strategy: "weighted"  # weighted | binary | hybrid
    veto_graders: ["must_pass_existing"]  # 一票否决grader
    partial_credit: true
```

---

#### Layer 5: 3D Dimension Aggregation (增强)

**Anthropic洞察**：
> "pass@k: likelihood of at least one correct solution in k attempts — for tools where one success matters"
> "pass^k: probability that all k trials succeed — for customer-facing agents where users expect reliable behavior"

**优化点**：

```go
// 增强的3D聚合，加入pass@k / pass^k
type Dimension3D struct {
    // 单维度聚合
    ByModel   map[string]AggregatedScore
    ByAgent   map[string]AggregatedScore
    BySkill   map[string]AggregatedScore
    
    // 双维度交叉
    ByAgentModel map[AgentModelKey]AggregatedScore
    ByAgentSkill map[AgentSkillKey]AggregatedScore
    ByModelSkill map[ModelSkillKey]AggregatedScore
    
    // 三维度全交叉
    ByAgentModelSkill map[AgentModelSkillKey]AggregatedScore
    
    // Delta分析
    AgentDeltas map[string]DeltaAnalysis
    ModelDeltas map[string]DeltaAnalysis
    SkillDeltas map[string]DeltaAnalysis
    InteractionEffects []InteractionEffect
    
    // NEW: pass@k / pass^k 统计
    PassAtK    map[AgentModelSkillKey]PassKStats
    PassAllK   map[AgentModelSkillKey]PassKStats
}

// NEW: 多次trial统计
type PassKStats struct {
    K              int       // trial次数
    PerTrialRate   float64   // 单次trial成功率
    PassAtK        float64   // 至少1次成功的概率
    PassAllK       float64   // 全部成功的概率 = PerTrialRate^K
    Trials         []TrialResult
}

type TrialResult struct {
    TrialIndex  int
    Passed      bool
    Score       float64
    PartialCredit float64
    Duration    time.Duration
    TokenUsage  TokenUsage
}
```

**对UT-Bench场景的映射**：

| 统计量 | 含义 | UT-Bench解读 |
|:---|:---|:---|
| pass@1 | 单次成功率 | Agent一次生成就对的概率 |
| pass@3 | 3次中至少1次成功 | Agent看到失败后能否修正 |
| pass^3 | 3次都成功 | Agent的稳定性/可靠性 |
| pass@1之差 | Agent vs 纯API | **Agent增益** |
| pass^3之差 | opencode vs claudecode | **Agent稳定性差异** |
| pass@k(python) vs pass@k(go) | 跨语言对比 | **语言维度** |

---

#### Layer 6: Eval Lifecycle & Saturation (NEW)

**Anthropic洞察**：
> "An eval at 100% tracks regressions but provides no signal for improvement."
> "Tasks that once measured 'Can we do this at all?' then measure 'Can we still do this reliably?'"

**这是完全由Anthropic方法论驱动的新增层**。

```go
// Eval生命周期管理
type EvalLifecycle struct {
    SuiteID     string
    Tasks       []TaskLifecycle
    LastUpdated time.Time
}

type TaskLifecycle struct {
    TaskID          string
    Phase           TaskPhase  // capability | transitioning | regression
    CurrentPassRate float64
    History         []PassRateSnapshot
    
    // 饱和度检测
    SaturationLevel float64   // 0.0~1.0, 越高越饱和
    LastSignalDate  time.Time  // 最后一次提供改进信号的日期
}

type TaskPhase string
const (
    PhaseCapability    TaskPhase = "capability"     // 低通过率,驱动改进
    PhaseTransitioning TaskPhase = "transitioning"  // 正在"毕业"
    PhaseRegression    TaskPhase = "regression"      // ≈100%,回归守卫
)

// "毕业"逻辑
func (l *EvalLifecycle) CheckGraduation(taskID string) GraduationDecision {
    task := l.getTask(taskID)
    
    // 条件: 连续3次运行pass_rate > 95%
    if task.CurrentPassRate > 0.95 && task.consecutiveHighRuns >= 3 {
        return GraduationDecision{
            ShouldGraduate: true,
            Reason:         "pass_rate > 95% for 3 consecutive runs",
            NewPhase:       PhaseRegression,
        }
    }
    
    return GraduationDecision{ShouldGraduate: false}
}

// 饱和度监控
func (l *EvalLifecycle) ComputeSaturation() float64 {
    // 饱和度 = 平均通过率的加权平均
    // 高通过率的task贡献更多饱和度
    totalWeight := 0.0
    totalSaturation := 0.0
    for _, task := range l.Tasks {
        weight := 1.0
        totalSaturation += task.CurrentPassRate * weight
        totalWeight += weight
    }
    return totalSaturation / totalWeight
}
```

**Capability vs Regression双轨运行**：

```
┌──────────────────────────────────────────────────────────────┐
│                    Evaluation Suite                           │
├────────────────────────┬─────────────────────────────────────┤
│  Capability Suite      │  Regression Suite                   │
│  (驱动改进)            │  (守卫不退步)                        │
│                        │                                     │
│  - 新增的hard task     │  - 已"毕业"的task                   │
│  - 预期pass率 < 50%    │  - 预期pass率 ≈ 100%               │
│  - pass@1最重要        │  - pass^k最重要                     │
│  - 发现模型差距        │  - 发现回归                          │
│  - 每次全量运行        │  - CI每次运行,快速反馈              │
│                        │                                     │
│  Saturation < 60%      │  Saturation ≈ 95%                   │
│  → 有改进信号          │  → 无改进信号,纯守卫                │
└────────────────────────┴─────────────────────────────────────┘
```

---

## 四、对原有Gap Analysis的影响

### 4.1 Gap优先级调整

原Gap Analysis的优先级排序需要根据Anthropic方法论重新调整：

| 排名 | 原优先级 | 优化项 | 新优先级 | 变化 | 原因 |
|:---|:---|:---|:---|:---|:---|
| 1 | P0 | 源码只读挂载 | P0 | 不变 | 评测公平性基础 |
| 2 | P0 | **Multi-Grader架构** | **P0** | **新增** | Anthropic最强调的方法论 |
| 3 | P0 | **pass@k / pass^k** | **P0** | **新增** | 区分"偶尔能过"vs"稳定能过" |
| 4 | P0 | **Partial Credit** | **P0** | **新增** | 评估精度大幅提升 |
| 5 | P0 | must_pass_existing | P0 | 不变 | 双重验证契约 |
| 6 | P1 | **Capability/Regression双轨** | **P1** | **新增** | 评测长期可持续性 |
| 7 | P1 | 交互轮次上限 | P1 | 不变 | 资源控制 |
| 8 | P1 | **Balanced Dataset** | **P1** | **升级** | Anthropic强调单侧优化的危害 |
| 9 | P1 | 结构化trajectory | P1 | 不变 | 可观测性基础 |
| 10 | P1 | **Reference Solution** | **P1** | **升级** | Anthropic强调验证eval本身 |

### 4.2 工作量重估

| 阶段 | 内容 | 工作量 |
|:---|:---|:---|
| Phase 1: 核心Grader架构 | Multi-Grader Interface + CodeGrader重构 + Partial Credit + must_pass_existing | 5天 |
| Phase 2: Trial与统计 | pass@k/pass^k + Trial隔离 + Reference Solution验证 + 评测重试 | 4天 |
| Phase 3: Lifecycle与平衡 | Capability/Regression双轨 + 饱和度监控 + 难度分级 + Balanced Dataset | 3天 |
| Phase 4: LLM Grader | LLMGrader实现 + Rubric设计 + 人工校准流程 | 3天 |
| **总计** | | **15天** (原10.5天+4.5天新增) |

---

## 五、UT-Bench的独特优势(不变+增强)

Anthropic文章验证了UT-Bench的几个独特优势：

### 5.1 被Anthropic验证的优势

1. **Coverage + Mutation Testing** — Anthropic的Coding Agent eval建议用unit tests + LLM rubric + static analysis，我们天然有前两项
2. **3D Matrix** — Anthropic文章没有涉及多维度交叉评测，这是我们的蓝海
3. **Skill Evaluation** — Anthropic文章完全没有提到skill评测概念

### 5.2 Anthropic文章新增的差异化方向

4. **Partial Credit + pass@k 三维切片** — pass@k按(Agent, Model, Skill)切片是全新视角
5. **Capability/Regression双轨 + 饱和度监控** — 所有竞品都没有这个
6. **Multi-Grader with Calibration** — Code+LLM+Human三层Grader + 人工校准

---

## 六、YAML评估配置完整示例

借鉴Anthropic文章中Coding Agent的YAML eval配置，设计UT-Bench的task定义：

```yaml
task:
  id: "python_flask_route_001"
  desc: "Generate unit tests for Flask route handlers including error cases"
  difficulty: "medium"
  category: "web_framework"
  language: "python"
  
  source:
    file: "samples/python/flask_routes.py"
    target_functions: ["handle_get_users", "handle_create_user", "handle_delete_user"]
    
  # Reference solution (验证eval本身)
  reference_solution: "samples/python/flask_routes_gold_test.py"
  reference_passes: true
  
  # Balanced design
  expected_pass_rate: 0.35
  negative_cases:
    - "should_not_modify_route_definitions"
    - "should_not_use_flask_test_client_incorrectly"
  
  graders:
    # Grader 1: Code (确定性, 必须)
    - type: code
      required: true
      assertions:
        - name: compile_pass
          weight: 0.15
          partial_credit: false  # binary
          
        - name: test_pass
          weight: 0.25
          partial_credit: true   # 5/7 tests pass = 0.71
          
        - name: must_pass_existing
          weight: 0.10
          veto: true             # 破坏已有测试→0分
          
        - name: line_coverage
          weight: 0.15
          threshold: 0.70
          partial_credit: true   # 60% coverage = 0.86
          
        - name: branch_coverage
          weight: 0.10
          threshold: 0.60
          partial_credit: true
          
        - name: mutation_score
          weight: 0.15
          threshold: 0.50
          partial_credit: true
    
    # Grader 2: LLM (模型评分, 可选)
    - type: llm
      required: false
      model: "deepseek-v4"  # 评分模型
      dimensions:
        - name: test_quality
          rubric: |
            Evaluate the generated unit tests on:
            1. Assertion quality: Are assertions specific and meaningful?
            2. Edge case coverage: Are boundary conditions tested?
            3. Test isolation: Are tests independent?
            4. Naming: Do test names clearly describe what's being tested?
          weight: 0.06
          
        - name: code_style
          rubric: |
            Evaluate code style on:
            1. PEP 8 compliance
            2. Proper imports
            3. No hardcoded values
            4. Appropriate use of fixtures
          weight: 0.04
          
    # Grader 3: Human (校准用)
    - type: human
      required: false
      sample_rate: 0.05
      dimensions:
        - name: overall_quality
          weight: 1.0
  
  # Trial配置
  trials:
    count: 3                    # 每个三元组跑3次
    isolation: "snapshot"       # 每次trial前恢复快照
    metrics: ["pass@1", "pass@3", "pass^3"]
  
  # 双重验证(借鉴SWE-bench)
  dual_verification:
    must_pass_existing: true
    existing_tests_pattern: "tests/test_original.py"
  
  # 成本追踪
  tracked_metrics:
    - type: cost
      metrics: [total_cost_usd, cost_per_coverage_point, tokens_per_test_case]
    - type: latency
      metrics: [time_to_first_token, total_duration]
    - type: transcript
      metrics: [n_turns, n_tool_calls, redundant_steps]
  
  # Lifecycle
  lifecycle:
    phase: "capability"  # 初始为capability eval
    graduation_threshold: 0.95
    graduation_consecutive_runs: 3
```

---

## 七、Swiss Cheese Model 在 UT-Bench 的应用

Anthropic强调"没有单一评估方法能捕获所有问题"。UT-Bench可以构建自己的Swiss Cheese：

```
┌────────────────────────────────────────────────────────────────┐
│                  UT-Bench Swiss Cheese Model                   │
├──────────────────┬──────────────────┬──────────────────────────┤
│                  │                  │                          │
│  Layer 1:        │  Layer 2:        │  Layer 3:               │
│  Automated Evals │  Transcript      │  Human Review            │
│  (每次CI)        │  Analysis        │  (定期校准)              │
│                  │  (每次运行)       │                          │
│  - compile_pass  │  - 策略合理性    │  - 测试质量抽检          │
│  - test_pass     │  - 工具效率      │  - LLM Grader校准       │
│  - coverage      │  - 冗余步骤      │  - 评测集有效性验证      │
│  - mutation      │  - 成本归因      │  - 难度分级校准          │
│  - must_pass     │                  │                          │
│                  │                  │                          │
│  捕获:           │  捕获:           │  捕获:                   │
│  功能正确性      │  过程效率问题    │  主观质量问题            │
│                  │                  │                          │
└──────────────────┴──────────────────┴──────────────────────────┘
   L1漏掉的 ←──── L2可能捕获 ←──── L3可能捕获
   (如: 编译通过  (如: 10轮冗余   (如: 测试断言
    但测试无用)    但最终过了)      无意义)
```

---

## 八、总结：Anthropic方法论对UT-Bench的核心价值

| 维度 | Anthropic洞察 | UT-Bench改动 | 预期收益 |
|:---|:---|:---|:---|
| **评估精度** | Partial Credit | Outcome层增加partial credit | 0分vs 0.6分的区分度 |
| **评估可靠性** | pass@k / pass^k | 多次trial + 统计 | 区分"偶尔能过"vs"稳定能过" |
| **评估深度** | Multi-Grader | 三层Grader架构 | Code+LLM+Human互补 |
| **评估可持续** | Lifecycle + Saturation | Capability/Regression双轨 | 评测集长期有效 |
| **评估公平** | Balanced + Reference | 难度分级 + gold solution | 避免单侧优化,验证eval本身 |
| **评估全面** | Swiss Cheese | 三层互补 | 捕获单一方法遗漏的问题 |
| **3D独特性** | (Anthropic未涉及) | pass@k三维切片 | 全行业独有的多维度统计 |

**一句话总结**：Anthropic文章验证了UT-Bench在Coverage+Mutation+Skill方面的独特优势，同时在**评估精度(partial credit)、评估可靠性(pass@k)、评估深度(Multi-Grader)、评估可持续(Lifecycle)**四个维度给出了我们之前完全没想到的关键优化方向。这些优化将UT-Bench从"能跑"提升到"专业级评测框架"。

---

## 九、现有代码到新架构的迁移映射

> 本节对照 go-ut-bench 现有代码，精确标注每个优化项涉及的文件、结构体、函数及改动方式。

### 9.1 三维标识传播现状（无需改动，已完备）

| 阶段 | 文件 | 结构体 | 三维字段 | 状态 |
|:---|:---|:---|:---|:---|
| 配置 | `agentconfig/config.go` | `SubjectSpec` | `Framework, Model, Skill` | ✅ 完备 |
| 配置 | `agentconfig/config.go` | `SubjectID()` 函数 | `{framework}__{model}__{skill}` | ✅ 完备 |
| 生成 | `contracts/results.go` | `GeneratedCase` | `AgentFramework, AgentModel, SkillName, SkillVersion` | ✅ 完备 |
| 评测 | `contracts/results.go` | `EvaluationResult` | `AgentFramework, AgentModel, SkillName, SkillVersion` | ✅ 完备 |
| 报告 | `reporter/aggregation.go` | `ModelDim` | `SubjectID, AgentFramework, AgentModel, SkillName` | ✅ 完备 |
| 报告 | `contracts/results.go` | `AgentComparisonRow` | `Framework, Model, Skill` | ✅ 完备 |
| 报告 | `contracts/results.go` | `SkillUpliftRow` | `Framework, Model, Skill, SkillVersion` | ✅ 完备 |

### 9.2 需要改动的文件和函数

#### 9.2.1 Layer 1: Sandbox 改动

| 文件 | 函数/结构体 | 改动类型 | 具体改动 |
|:---|:---|:---|:---|
| `runner/sandbox.go` | `SandboxRunRequest` | **新增字段** | `TrialIndex int`, `TotalTrials int`, `CleanSnapshot string`, `PreRunHash string` |
| `runner/sandbox.go` | `SandboxRunner` 接口 | **新增方法** | `VerifyIsolation(req) error` |
| `runner/sandbox.go` | `defaultSandboxRunner.Run()` | **修改逻辑** | Docker 模式下增加 `--read-only` 源码挂载逻辑；每次 trial 前执行 `git checkout .` |
| `runner/adapter_cli.go` | `generateCLIAgent()` | **修改逻辑** | 步骤1前增加 trial 循环；步骤6快照后增加 `PreRunHash` 计算；步骤12后增加隔离校验 |

#### 9.2.2 Layer 2: Harness 改动

| 文件 | 函数/结构体 | 改动类型 | 具体改动 |
|:---|:---|:---|:---|
| `contracts/spec.go` | `SampleRef` | **新增字段** | `Difficulty string`, `ReferenceSolution string`, `ExpectedPassRate float64`, `NegativeCases []string` |
| `contracts/spec.go` | `RepoLevelMeta` | **新增字段** | `MustPassExisting bool`, `ExistingTestsPattern string`, `ReferenceSolution string` |
| **新文件** | `harness/reference.go` | **新增** | `VerifyReferenceSolution() error` — 运行 gold test 验证 eval 基础设施 |

#### 9.2.3 Layer 3: EventStream 改动

| 文件 | 函数/结构体 | 改动类型 | 具体改动 |
|:---|:---|:---|:---|
| `runner/adapter.go` | `AgentTrace` | **新增字段** | `TrialIndex int`, `OutcomeEvents []OutcomeEvent`, `TranscriptQuality *TranscriptQuality` |
| **新文件** | `contracts/events.go` | **新增** | `TaggedEvent`, `OutcomeEvent`, `TranscriptQuality` 结构体定义 |
| `runner/adapter_cli.go` | `parseAgentOutput()` | **修改逻辑** | 增加 Outcome 事件的解析（从 stdout/stderr 提取 compile/test/coverage/mutation 结果） |

#### 9.2.4 Layer 4: Multi-Grader 改动（最大改动）

| 文件 | 函数/结构体 | 改动类型 | 具体改动 |
|:---|:---|:---|:---|
| **新文件** | `grader/interface.go` | **新增** | `Grader` 接口, `GraderInput`, `GraderResult`, `AssertionResult`, `GraderType` 常量 |
| **新文件** | `grader/code_grader.go` | **新增** | `CodeGrader` 实现 — 重构现有 evaluator 逻辑为 Code Grader |
| **新文件** | `grader/llm_grader.go` | **新增** | `LLMGrader` 实现 — rubric 评分 |
| **新文件** | `grader/human_grader.go` | **新增** | `HumanGrader` 实现 — 抽检 + 校准 |
| **新文件** | `grader/composer.go` | **新增** | `GraderComposer` — 按 YAML 配置组合多个 Grader，处理 veto 和 partial credit |
| `contracts/results.go` | `EvaluationResult` | **新增字段** | `GraderResults []GraderResultJSON`, `PartialCredit float64`, `MustPassExisting *bool` |
| `contracts/constants.go` | `ScoreWeights` | **修改** | 增加字段 `TestQuality float64`, `Efficiency float64`；增加 `VetoGraders []string` |

**现有 evaluator 的迁移路径**：

```
当前流程:
  runner.Generate() → evaluator.Evaluate() → reporter.Generate()

迁移后:
  runner.Generate() → grader.Compose().Grade() → reporter.Generate()
                            ↓
                     ┌──────┼──────┐
                  CodeGrader  LLMGrader  HumanGrader
                  (从现有      (新增)      (新增)
                  evaluator
                  重构而来)
```

**关键迁移步骤**：
1. 将 `evaluator/` 中的编译、测试、覆盖率、变异测试逻辑封装为 `CodeAssertion`
2. `CodeGrader` 的 `Grade()` 方法 = 现有 `evaluator.Evaluate()` 的核心逻辑
3. 每个断言增加 `PartialCreditFn` 字段（之前是 binary pass/fail）
4. `GraderComposer` 按权重汇总，处理 veto 逻辑

#### 9.2.5 Layer 5: 3D Aggregation 改动

| 文件 | 函数/结构体 | 改动类型 | 具体改动 |
|:---|:---|:---|:---|
| `reporter/aggregation.go` | `buildDimensions()` | **重构** | 当前 `ByModel` 实际是 `BySubject`（key=row.Model），需拆为独立 `ByAgent`、`ByModel`、`BySkill` |
| `contracts/results.go` | `Dimensions` | **新增字段** | `ByAgent []AgentDim`, `BySkill []SkillDim`, `ByAgentModel []AgentModelDim`, `ByAgentSkill []AgentSkillDim`, `ByModelSkill []ModelSkillDim`, `ByAgentModelSkill []AgentModelSkillDim` |
| `contracts/results.go` | **新增结构体** | **新增** | `PassKStats`, `TrialResult` |
| `reporter/ranking.go` | `buildAgentComparisons()` | **增强** | 使用 pass@k/pass^k 代替单次 pass_rate 做 delta 分析 |

**`buildDimensions()` 核心改动详解**：

```go
// 当前: 所有 subject 按同一个 row.Model key 聚合，三维字段取首条
// 问题: opencode__deepseek__no_skill 和 aider__deepseek__no_skill 被合并到同一行
func buildDimensions(rows []EvaluationResult) Dimensions {
    byModel := make(map[string]*modelAgg)
    for _, row := range rows {
        key := row.Model  // BUG: 同 model 不同 framework/skill 被合并
        agg, ok := byModel[key]
        if !ok {
            agg = &modelAgg{}
            // 首条记录的三维字段直接覆盖
            agg.SubjectID = row.SubjectID
            agg.AgentFramework = row.AgentFramework
            // ...
        }
        agg.add(row)
        byModel[key] = agg
    }
    // ...
}

// 优化后: 按5个独立维度分别聚合
func buildDimensionsV2(rows []EvaluationResult) Dimensions {
    byModel   := aggregateBy(rows, func(r EvaluationResult) string { return r.AgentModel })
    byAgent   := aggregateBy(rows, func(r EvaluationResult) string { return r.AgentFramework })
    bySkill   := aggregateBy(rows, func(r EvaluationResult) string { return r.SkillName })
    byAgentModel := aggregateBy2D(rows, func(r EvaluationResult) string {
        return r.AgentFramework + "|" + r.AgentModel
    })
    // ... 以及 pass@k 统计
}
```

#### 9.2.6 Layer 6: Lifecycle 改动

| 文件 | 函数/结构体 | 改动类型 | 具体改动 |
|:---|:---|:---|:---|
| **新文件** | `lifecycle/manager.go` | **新增** | `EvalLifecycle`, `TaskLifecycle`, `TaskPhase` 结构体和毕业逻辑 |
| **新文件** | `lifecycle/saturation.go` | **新增** | `ComputeSaturation()` 饱和度计算 |
| `contracts/spec.go` | `SampleRef` | **新增字段** | `Phase string` (capability/regression), `GraduationThreshold float64` |

---

## 十、pass@k 数学公式

Anthropic 文章给出了精确的数学定义，UT-Bench 需要严格遵循：

### 10.1 定义

设 n = 总 trial 数，c = 成功 trial 数，k = 考察的 trial 数。

**pass@k** = 在 k 次尝试中至少成功 1 次的概率：

```
pass@k = 1 - C(n-c, k) / C(n, k)
```

其中 C(n, k) 是组合数。当 n = k 时退化为 `pass@k = c/n`。

**pass^k** = k 次尝试全部成功的概率：

```
pass^k = (c/n)^k
```

### 10.2 实现注意事项

```go
// pass@k 计算需注意数值稳定性
// 直接计算组合数容易溢出,用对数空间
func PassAtK(n int, c int, k int) float64 {
    if n-c < k {
        return 1.0 // 成功数已经足够,必定至少1次成功
    }
    // log(C(n-c, k)) - log(C(n, k))
    // = sum_{i=0}^{k-1} [log(n-c-i) - log(n-i)]
    logResult := 0.0
    for i := 0; i < k; i++ {
        logResult += math.Log(float64(n-c-i)) - math.Log(float64(n-i))
    }
    return 1.0 - math.Exp(logResult)
}

func PassAllK(perTrialRate float64, k int) float64 {
    return math.Pow(perTrialRate, float64(k))
}
```

### 10.3 在 3D 中的语义

| 场景 | 公式 | UT-Bench 含义 |
|:---|:---|:---|
| `pass@1(opencode, deepseek, tdd_skill)` | 单次成功率 | TDD skill 在 opencode+deepseek 上的即时成功率 |
| `pass@3(opencode, deepseek, tdd_skill)` | 3次至少1次 | TDD skill 在 opencode+deepseek 上的"可修正率" |
| `pass^3(opencode, deepseek, tdd_skill)` | 3次全过 | TDD skill 在 opencode+deepseek 上的稳定性 |
| `pass@1(opencode, deepseek, no_skill)` | baseline | 无 skill 的基准成功率 |
| **Agent 增益** | `pass@1(with_agent) - pass@1(model_api)` | Agent 框架带来的提升 |
| **Skill 增益** | `pass@1(with_skill) - pass@1(no_skill)` | Skill 带来的提升 |
| **交互效应** | `pass@1(opencode+skill) - pass@1(opencode) - pass@1(model_api+skill) + pass@1(model_api)` | Agent 和 Skill 的协同增益 |

---

## 十一、LLM Grader 成本估算

### 11.1 调用量估算

| 参数 | 值 | 说明 |
|:---|:---|:---|
| 样本数 | 800 | self_contained 全量 |
| 3D 三元组数 | ~15 | 3 frameworks × 5 models × 1 skill(初始) |
| Trial 数 | 3 | pass@3 需要 |
| LLM 评分维度 | 3 | test_quality + code_style + assertion_meaningfulness |
| **总调用量** | 800 × 15 × 3 × 3 = **108,000** | |

### 11.2 单次调用成本估算

| 评分模型 | Input tokens | Output tokens | 单次成本(USD) |
|:---|:---|:---|:---|
| deepseek-v4 | ~800 (rubric+code) | ~200 (评分+reason) | ~$0.001 |
| gpt-4o-mini | ~800 | ~200 | ~$0.0005 |
| claude-3.5-haiku | ~800 | ~200 | ~$0.001 |

### 11.3 总成本估算

| 方案 | 模型 | 总调用量 | 总成本(USD) | 说明 |
|:---|:---|:---|:---|:---|
| 最低成本 | gpt-4o-mini | 108,000 | **~$54** | 质量可能不足 |
| 推荐方案 | deepseek-v4 | 108,000 | **~$108** | 性价比最优 |
| 最高质量 | claude-3.5-haiku | 108,000 | **~$108** | 评分一致性好 |

**优化措施**：
1. **缓存**：同一 (code, rubric) 组合只评一次，不重复调用 → 预计减少 30%
2. **只对 Code Grader 通过的样本调用 LLM Grader** → 预计减少 40-60%
3. **抽检模式**：初期只对 10% 样本跑 LLM Grader → 成本降至 ~$10

---

## 十二、LLM Grader Rubric 设计

### 12.1 test_quality_rubric.md

```markdown
# Unit Test Quality Rubric

You are evaluating the quality of a generated unit test file. Score each dimension from 0.0 to 1.0.

## Input
- Source code file (the code being tested)
- Generated test file (the test to evaluate)
- Test execution results (pass/fail per test case)

## Dimensions

### 1. Assertion Quality (weight: 0.4)
Score based on:
- 1.0: All assertions are specific, meaningful, and test distinct behaviors
- 0.7: Most assertions are specific, some are trivial (e.g., assert result is not None)
- 0.4: Many assertions are weak or redundant
- 0.1: Almost no meaningful assertions (only smoke tests like "no exception")
- 0.0: No assertions at all

### 2. Edge Case Coverage (weight: 0.3)
Score based on:
- 1.0: Boundary values, empty inputs, null/None, type errors all covered
- 0.7: Some edge cases covered, missing obvious ones
- 0.4: Only happy path tested
- 0.1: Minimal edge case testing
- 0.0: No edge cases tested

### 3. Test Isolation & Independence (weight: 0.2)
Score based on:
- 1.0: Each test is fully independent, proper setup/teardown, no shared mutable state
- 0.7: Mostly independent, minor coupling
- 0.4: Some tests depend on execution order or shared state
- 0.1: Tests are heavily coupled
- 0.0: Tests cannot run independently

### 4. Naming & Documentation (weight: 0.1)
Score based on:
- 1.0: Test names clearly describe what is being tested; docstrings present
- 0.7: Test names are mostly descriptive
- 0.4: Test names are generic (test_function_1, test_function_2)
- 0.1: No meaningful names
- 0.0: Tests have no names

## Output Format

```json
{
  "assertion_quality": { "score": 0.0-1.0, "reason": "..." },
  "edge_case_coverage": { "score": 0.0-1.0, "reason": "..." },
  "test_isolation": { "score": 0.0-1.0, "reason": "..." },
  "naming": { "score": 0.0-1.0, "reason": "..." },
  "overall": 0.0-1.0,
  "overall_reason": "..."
}
```

If you don't have enough information to evaluate a dimension, output "unknown" for that dimension.
```

### 12.2 code_style_rubric.md

```markdown
# Code Style Rubric

Evaluate the code style of the generated unit test file.

## Dimensions

### 1. Standard Compliance (weight: 0.3)
- Python: PEP 8 compliance
- Java: Google Java Style / Sun convention
- Go: gofmt compliance
- C++: LLVM/Google style

### 2. Import & Dependency (weight: 0.2)
- Proper imports, no unused imports, correct mock/patch usage

### 3. Hardcoded Values (weight: 0.3)
- 1.0: No hardcoded test data; uses fixtures/factories/parameterization
- 0.5: Some hardcoded values that should be parameterized
- 0.0: All test data hardcoded, brittle to changes

### 4. Fixture Usage (weight: 0.2)
- 1.0: Appropriate use of setup/teardown, fixtures, test doubles
- 0.5: Some fixture usage, could be improved
- 0.0: No fixtures, all setup inline

## Output Format (same as test_quality_rubric)
```

### 12.3 assertion_rubric.md

```markdown
# Assertion Meaningfulness Rubric

Evaluate whether each assertion in the test is actually testing something meaningful,
or if it's a "fake" assertion that would pass regardless of code correctness.

## Red flags (score 0.0-0.3)
- `assert True` or equivalent
- `assert result is not None` without further checks
- `assert len(result) > 0` without checking content
- Assertions that only verify function didn't raise exception
- Assertions that compare against the function's own output

## Good patterns (score 0.7-1.0)
- Assert specific expected values
- Assert type + content
- Assert error messages or exception types
- Parameterized assertions testing multiple inputs
- Assertions that would fail if implementation changes

## Output Format (same as test_quality_rubric)
```

---

## 十三、Balanced Dataset 设计

### 13.1 "正样本" vs "负样本" 在 UT 生成场景的定义

在 Anthropic 的语境中，"balanced" 指的是：
- **正样本**：Agent 应该成功生成的测试（验证能力上限）
- **负样本**：Agent 不应该生成的测试（验证不会做错事）

在 UT 生成的独特场景下，"负样本"的定义不同于传统分类任务：

| 样本类型 | 含义 | 占比 | 示例 |
|:---|:---|:---|:---|
| **Easy positive** | 单函数，简单逻辑，预期高 pass 率 | 25% | `max()`, `is_even()` |
| **Medium positive** | 单类多方法，需要 mock | 25% | `UserRepository.save()`, `FlaskRoute.handle()` |
| **Hard positive** | 跨模块依赖，需要复杂 fixture | 20% | `Django ORM QuerySet`, `grpc Service` |
| **Negative: 不可测** | 源码本身无法被单元测试（纯 UI、纯 I/O） | 10% | GUI 渲染函数、裸 `main()` |
| **Negative: 陷阱** | 容易生成"假测试"的场景 | 10% | 时间相关函数、随机数、全局状态 |
| **Negative: 破坏性** | 容易让 Agent 修改源码来通过测试 | 10% | 有 bug 的源码（Agent 可能修 bug 而非写测试） |

### 13.2 难度分级标准

| 难度 | 依赖层级 | 目标函数数 | 需要 mock | 预期 pass 率 |
|:---|:---|:---|:---|:---|
| easy | 0（纯函数） | 1 | 否 | >80% |
| medium | 1-2（类内依赖） | 2-5 | 简单 mock | 40-60% |
| hard | 3+（跨模块） | 5+ | 复杂 fixture | 10-30% |
| expert | 仓库级 | 全类 | 深度 mock + 数据库 | <10% |

### 13.3 每个 self_contained 样本的扩展字段

```json
{
  "sample_id": "python_sort_001",
  "language": "python",
  "scenario": "simple_function",
  
  "difficulty": "easy",
  "dependency_level": 0,
  "target_function_count": 1,
  "requires_mock": false,
  "expected_pass_rate": 0.85,
  
  "sample_type": "easy_positive",
  "negative_signals": [],
  
  "reference_solution": "gold/test_sort_gold.py"
}
```

---

## 十四、Grader → Reporter 数据流

### 14.1 完整数据流图

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                              Generation Phase                                 │
│                                                                              │
│  agentconfig.Load()                                                          │
│       │                                                                      │
│       ▼                                                                      │
│  ResolvedSubject {framework, model, skill}  ──────────┐                     │
│       │                                                │                     │
│       ▼                                                │ SubjectID           │
│  AgentAdapter.Generate()                               │ 传播                 │
│       │                                                │                     │
│       ▼                                                ▼                     │
│  AgentGenerateResult ────────────────────► GeneratedCase                     │
│  (AgentTrace 已含 TrialIndex,             (三维字段已填)                      │
│   OutcomeEvents)                                                              │
└──────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                              Grading Phase (NEW)                              │
│                                                                              │
│  GraderComposer.Grade(input)                                                 │
│       │                                                                      │
│       ├──► CodeGrader.Grade() ──────── GraderResult {Score, PartialCredit,   │
│       │    (compile, test, coverage,     Assertions, Reason}                  │
│       │     mutation, must_pass_existing)                                     │
│       │                                                                      │
│       ├──► LLMGrader.Grade() ──────── GraderResult {Score, Reason,           │
│       │    (rubric评分)                  Dimensions}                          │
│       │                                                                      │
│       └──► HumanGrader.Grade() ────── GraderResult (仅抽样校准)              │
│                                                                              │
│  Composer 汇总:                                                               │
│  - 处理 veto (must_pass_existing 失败 → 总分 0)                               │
│  - 加权汇总 partial credit                                                    │
│  - 输出 CompositeGraderResult                                                │
└──────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                           EvaluationResult (增强)                             │
│                                                                              │
│  原有字段: CompilePass, TestPass, LineCoverage, BranchCoverage,              │
│           MutationScore, AssertionCount, TestCaseCount, ...                  │
│                                                                              │
│  新增字段:                                                                    │
│  - MustPassExisting *bool          // 破坏已有测试?                           │
│  - PartialCredit float64           // 部分得分 (0.0~1.0)                     │
│  - GraderResults []GraderResultJSON  // 每个 Grader 的详细结果               │
│  - TrialIndex int                  // 第几次 trial                           │
│  - TrialID string                  // trial 唯一 ID                          │
└──────────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│                           Reporter Phase (增强)                               │
│                                                                              │
│  buildDimensionsV2():                                                        │
│  ┌─────────────────────────────────────────────────────────┐                 │
│  │ ByModel   (key: AgentModel)                             │                 │
│  │ ByAgent   (key: AgentFramework)    ← NEW                │                 │
│  │ BySkill   (key: SkillName)         ← NEW                │                 │
│  │ ByAgentModel (key: framework|model) ← NEW               │                 │
│  │ ByAgentSkill (key: framework|skill) ← NEW               │                 │
│  │ ByModelSkill (key: model|skill)     ← NEW               │                 │
│  │ ByAgentModelSkill (key: f|m|s)      ← NEW               │                 │
│  └─────────────────────────────────────────────────────────┘                 │
│                                                                              │
│  buildPassKStats():  ← NEW                                                   │
│  ┌─────────────────────────────────────────────────────────┐                 │
│  │ Per (Agent, Model, Skill) 三元组:                        │                 │
│  │   pass@1, pass@3, pass^3                                 │                 │
│  │   PartialCredit 均值/中位数                               │                 │
│  │   Cost per trial                                         │                 │
│  └─────────────────────────────────────────────────────────┘                 │
│                                                                              │
│  buildDeltaAnalysis(): (增强)                                                 │
│  ┌─────────────────────────────────────────────────────────┐                 │
│  │ Agent Delta: pass@1(with_agent) - pass@1(model_api)     │                 │
│  │ Model Delta: pass@1(model_a) - pass@1(model_b)          │                 │
│  │ Skill Delta: pass@1(with_skill) - pass@1(no_skill)      │                 │
│  │ Interaction: Agent×Skill 协同效应                         │                 │
│  └─────────────────────────────────────────────────────────┘                 │
└──────────────────────────────────────────────────────────────────────────────┘
```
