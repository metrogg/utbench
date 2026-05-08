# UT-Bench 项目详细介绍

## 一、项目概述

### 1.1 项目定位

UT-Bench（go-ut-bench）是一个基于 Go 语言开发的 **LLM 单元测试生成能力横向评测基准工具**。该工具的核心目标是：

- **横向对比**：在同一基准数据集上，对多个大语言模型（LLM）生成单元测试的能力进行公平、客观的横向对比评测
- **多维评估**：从编译通过率、测试通过率、代码覆盖率、变异测试得分等多个维度全面评估测试生成质量
- **自动化流程**：提供端到端的自动化评测流水线，从测试生成到报告生成全程自动化

### 1.2 核心价值

1. **科学评测**：采用业界标准的评测指标和方法，确保评测结果的科学性和可信度
2. **公平对比**：所有模型使用相同的提示词策略、相同的数据集样本，确保横向对比的公平性
3. **可视化呈现**：生成包含图表、排名、洞察分析的 HTML 报告，便于结果分析和决策
4. **历史追踪**：通过 SQLite 数据库持久化评测数据，支持跨运行对比和历史趋势分析

### 1.3 适用场景

- **模型选型**：帮助企业或研究团队选择最适合代码测试生成任务的 LLM
- **能力评估**：评估 LLM 在软件测试领域的实际应用能力
- **基准研究**：为学术研究提供标准的评测基准和数据
- **质量监控**：持续监控模型版本升级后的测试生成能力变化

---

## 二、整体架构

### 2.1 架构图

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CLI / Web UI                                    │
│                     (cmd/utbench/main.go + internal/web)                     │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Orchestrator (流水线编排)                             │
│                     (internal/orchestrator/service.go)                       │
│                                                                              │
│   Phase: full → generate → evaluate → report                                │
└─────────────────────────────────────────────────────────────────────────────┘
          │                    │                    │                    │
          ▼                    ▼                    ▼                    ▼
    ┌──────────┐         ┌──────────┐         ┌──────────┐         ┌──────────┐
    │  Dataset │         │  Runner  │         │ Evaluator│         │ Reporter │
    │ Service  │         │ Service  │         │ Service  │         │ Service  │
    └──────────┘         └──────────┘         └──────────┘         └──────────┘
          │                    │                    │                    │
          │                    │                    │                    │
          ▼                    ▼                    ▼                    ▼
    ┌──────────┐         ┌──────────┐         ┌──────────┐         ┌──────────┐
    │ 数据集   │         │ LLM API  │         │ 工具链   │         │ HTML生成 │
    │ 发现与   │         │ 调用与   │         │ 执行：   │         │ Chart.js │
    │ 过滤     │         │ Worker   │         │ 编译/测试│         │ 可视化   │
    │          │         │ Pool     │         │ /覆盖率  │         │          │
    │          │         │          │         │ /变异    │         │          │
    └──────────┘         └──────────┘         └──────────┘         └──────────┘
          │                    │                    │                    │
          └────────────────────┴────────────────────┴────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         SQLite Store (数据持久化)                            │
│                     (internal/store/sqlite.go)                               │
│                                                                              │
│   Tables: experiments, generation_runs, evaluation_runs,                    │
│           generated_cases, evaluation_results, artifacts, ...               │
└─────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 五阶段流水线详解

项目采用 **五阶段流水线架构**，由 `internal/orchestrator/service.go` 统一编排：

| 阶段 | 服务模块 | 核心功能 | 输入 | 输出 |
|------|----------|----------|------|------|
| 1. Dataset | `internal/dataset/` | 数据集样本发现与过滤 | RunSpec（运行配置） | []SampleRef（样本引用列表） |
| 2. Runner | `internal/runner/` | LLM API 调用，生成单元测试代码 | RunSpec + Samples | GeneratedManifest（生成清单） |
| 3. Evaluator | `internal/evaluator/` | 编译 → 测试 → 覆盖率 → 变异测试 | Manifest + 测试代码 | EvaluationResultSet（评测结果） |
| 4. Reporter | `internal/reporter/` | 多维度聚合分析与 HTML 报告生成 | EvaluationResultSet | ReportPayload + report.html |
| 5. Store | `internal/store/` | SQLite 持久化，支持历史查询与复用 | Run 目录产物 | DB Tables |

---

## 三、核心模块详解

### 3.1 数据契约层 (internal/contracts/)

`internal/contracts/` 是整个项目的 **数据结构定义中心**，定义了跨模块传递的所有数据类型。

#### 3.1.1 核心数据结构

**RunSpec（运行规格说明）**

定义于 `spec.go`，包含所有影响评测运行的配置参数：

```go
type RunSpec struct {
    RunID           string    // 唯一运行ID
    Models          []string  // 要评测的模型列表
    Languages       []string  // 编程语言列表
    DatasetClasses  []string  // 数据集类别（self_contained/repo_level）
    DatasetScenario string    // 场景过滤
    DatasetLevel    string    // 数据集级别
    DatasetRoot     string    // 数据集根目录
    ConfigPath      string    // 模型配置文件路径
    Mode            RunMode   // 运行模式：full/incremental
    DryRun          bool      // 试运行模式
    ReuseGenerated  bool      // 是否复用历史生成结果
    DBPath          string    // SQLite 数据库路径
    MutationEnabled bool      // 是否启用变异测试
    MutationTimeout int       // 变异测试超时（秒）
    TestTimeout     int       // 测试执行超时（秒）
    MaxSamples      int       // 最大样本数限制
    Workers         int       // 并发 Worker 数量
    OutputRoot      string    // 输出根目录
}
```

**SampleRef（样本引用）**

定义于 `spec.go`，唯一标识数据集中的样本：

```go
type SampleRef struct {
    ID        string       // 样本ID，格式：<场景>_<编号>
    Language  string       // 编程语言
    Category  DatasetClass // 类别：self_contained 或 repo_level
    Scenario  string       // 场景名
    Path      string       // 源文件绝对路径
    SourceMD5 string       // 源码 MD5 哈希
}
```

**GeneratedManifest（生成清单）**

定义于 `results.go`，记录测试生成阶段的完整产出：

```go
type GeneratedManifest struct {
    SchemaVersion     string          // 数据结构版本
    RunID             string          // 关联的运行ID
    CreatedAtUTC      time.Time       // 创建时间
    Spec              RunSpec         // 运行配置快照
    PromptStrategy    string          // 提示词策略名
    PromptVersionID   string          // 提示词版本ID
    PromptSnapshotDir string          // 提示词快照目录
    Cases             []GeneratedCase // 所有生成案例列表
}
```

**EvaluationResult（评测结果）**

定义于 `results.go`，记录单个样本的完整评测数据：

```go
type EvaluationResult struct {
    Model             string   // 模型名称
    Language          string   // 编程语言
    SampleID          string   // 样本ID
    CompilePass       bool     // 编译是否通过
    TestPass          *bool    // 测试是否通过
    LineCoverage      *float64 // 行覆盖率（0-100）
    BranchCoverage    *float64 // 分支覆盖率（0-100）
    MutationScore     *float64 // 变异得分（0-100）
    MutationTotal     *int     // 变异体总数
    MutationKilled    *int     // 被杀死变异体数
    MutationSurvived  *int     // 存活变异体数
    TestCaseCount     *int     // 测试用例数
    AssertionCount    *int     // 断言数
    RuntimeMS         *int     // 评测耗时（毫秒）
    LatencyMS         *int     // API调用耗时（毫秒）
    FailureOrigin     string   // 失败归因
    ScoreEligible     *bool    // 是否计入排名
}
```

#### 3.1.2 常量定义

定义于 `constants.go`：

```go
const SchemaVersion = "v0.1.0"  // 数据结构版本

type DatasetClass string
const (
    DatasetClassSelfContained DatasetClass = "self_contained"  // 单文件自包含
    DatasetClassRepoLevel   DatasetClass = "repo_level"    // 多文件仓库级
)

type RunMode string
const (
    RunModeFull        RunMode = "full"        // 完整运行
    RunModeIncremental RunMode = "incremental" // 断点续跑
)

var SupportedLanguages = []string{"python", "java", "go", "cpp"}
```

### 3.2 数据集发现模块 (internal/dataset/)

#### 3.2.1 核心职责

- **样本发现**：扫描数据集目录，识别所有样本文件
- **自动分类**：从目录结构推断样本的语言、类别、场景
- **过滤应用**：根据 RunSpec 参数过滤样本（语言、场景、数量限制）
- **哈希计算**：计算源码 MD5，用于校验和复用判断

#### 3.2.2 目录结构推断规则

数据集采用标准目录结构，Dataset Service 从路径自动推断属性：

```
datasets/
  python/
    self_contained/
      boundary/
        boundary_001.py
        boundary_002.py
      simple_function/
        simple_001.py
    repo_level/
      dateutil/
        meta.json
        parser.py
  go/
    self_contained/
      boundary/
        boundary_001.go
```

推断逻辑：
- **Language**：从第一级目录名推断（python/go/java/cpp）
- **Category**：从第二级目录名推断（self_contained/repo_level）
- **Scenario**：从第三级目录名推断（boundary/simple_function/complex_dependency/interface_mock）
- **SampleID**：从文件名推断（去掉扩展名）

#### 3.2.3 核心方法

```go
// ValidateSpec 验证运行配置的有效性
func (s *Service) ValidateSpec(spec contracts.RunSpec) error

// DiscoverSamples 发现并过滤符合条件的样本
func (s *Service) DiscoverSamples(spec contracts.RunSpec) ([]contracts.SampleRef, error)

// BuildIndex 构建数据集索引文件
func (s *Service) BuildIndex(datasetRoot, outputPath string) (*IndexSummary, error)

// BuildManifest 根据索引构建清单文件
func (s *Service) BuildManifest(opts ManifestBuildOptions) (*ManifestSummary, error)

// ValidateReadiness 验证数据集准备状态
func (s *Service) ValidateReadiness(opts ValidateOptions) *ValidationReport
```

### 3.3 测试生成模块 (internal/runner/)

#### 3.3.1 核心职责

- **LLM API 调用**：并发调用多个 LLM API 生成测试代码
- **Worker Pool**：并发 Worker池管理，提高生成效率
- **Checkpoint机制**：支持断点续跑，避免重复生成
- **自动续写**：检测截断响应（finish_reason="length"），自动续写补全
- **提示词管理**：三种提示词模式，针对不同样本类型

#### 3.3.2 提示词系统

定义于 `prompt.go`，提供三种提示词模式：

| 模式 | 适用场景 | 特点 |
|------|----------|------|
| `full_file` | self_contained 类型样本 | 默认模式，提供完整源码 |
| `repo_level` | repo_level 类型样本 | 包含 workspace、module_import 等上下文 |
| `completion` | 截断后续写 | 基于已生成内容继续补全 |

系统消息（System Message）强调：

```
"You are a senior unit test generation model. 
Return only runnable test code. 
Do not include explanations, markdown fences, or placeholder tests. 
Use only the provided source and context. 
Do not invent APIs, imports, or behavior."
```

提示词策略版本：`structured-v1`，版本ID通过 SHA1 哈希生成，确保可追溯。

#### 3.3.3 Checkpoint机制

Runner 支持增量执行，避免重复调用 API：

```go
// Checkpoint路径计算
func buildCheckpointPath(spec contracts.RunSpec) string {
    hash := sha1.Sum([]byte(
        spec.DatasetRoot + 
        strings.Join(spec.DatasetClasses, ",") +
        spec.DatasetLevel +
        spec.DatasetManifest +
        strconv.Itoa(spec.MaxSamples) +
        strings.Join(spec.Models, ",")
    ))
    return filepath.Join(spec.OutputRoot, "checkpoints", 
        "runner_" + hex.EncodeToString(hash[:]) + ".checkpoint.json")
}
```

Checkpoint 文件格式：

```json
{
    "run_id": "run_123",
    "completed_tasks": ["deepseek|python|boundary_001", "qwen|python|boundary_001"],
    "updated_at": "2026-04-28T10:00:00Z"
}
```

任何改变参数的操作会生成新的 Hash，导致 Checkpoint 失效重新运行。

#### 3.3.4 Worker Pool 并发模型

```go
type Service struct {
    logger *obs.Logger
    // Worker池并发执行模型调用
}

func (s *Service) Generate(ctx context.Context, spec contracts.RunSpec, samples []contracts.SampleRef) (*Output, error) {
    // 创建任务队列
    tasks := buildTasks(spec, samples)
    
    // Worker池并发执行
    var wg sync.WaitGroup
    for i := 0; i < spec.Workers; i++ {
        wg.Add(1)
        go s.worker(ctx, &wg, tasks, results)
    }
    wg.Wait()
    
    return buildManifest(results)
}
```

#### 3.3.5 截断检测与自动续写

当 LLM 返回 `finish_reason: "length"` 时，表示输出被 max_tokens 限制截断。Runner 会自动检测并进行续写：

```go
func (s *Service) handleTruncation(ctx context.Context, req GenerateRequest, existingTest string) (string, error) {
    // 最多3轮续写
    for round := 0; round < 3; round++ {
        continuation := s.generateCompletion(ctx, req, existingTest)
        existingTest += continuation
        
        if !isTruncated(continuation) {
            return existingTest, nil  // 成功补全
        }
    }
    return existingTest, fmt.Errorf("max continuation rounds reached")
}
```

### 3.4 评测执行模块 (internal/evaluator/)

#### 3.4.1 核心职责

- **编译验证**：检查生成的测试代码能否编译通过
- **测试运行**：执行测试用例，收集通过率
- **覆盖率收集**：使用语言对应工具收集行覆盖率、分支覆盖率
- **变异测试**：运行变异测试工具，计算变异得分
- **失败归因**：自动分析失败原因（model/environment/dataset/tool）

#### 3.4.2 评测流水线（单样本）

每个样本经历四个阶段：

```
┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐
│ Compile │ → │ Test    │ → │ Coverage│ → │ Mutation│
│ Stage   │    │ Stage   │    │ Stage   │    │ Stage   │
└─────────┘    └─────────┘    └─────────┘    └─────────┘
     │              │              │              │
     ▼              ▼              ▼              ▼
 [编译结果]    [测试通过率]   [覆盖率数据]   [变异得分]
```

任何阶段失败都会中断后续阶段，并记录失败归因。

#### 3.4.3 语言对应工具链

| 语言 | 编译器 | 测试框架 | 覆盖率工具 | 变异测试工具 |
|------|--------|----------|------------|--------------|
| Python | pylint | pytest | coverage.py | mutmut |
| Go | go build | go test | go tool cover | go-mutesting |
| Java | javac | JUnit 5 | JaCoCo | PITest |
| C++ | g++ | GoogleTest | gcov | mull |

#### 3.4.4 隔离执行环境

评测在 **临时目录** 中隔离执行，避免相互干扰：

```go
func (s *Service) evaluateOne(ctx context.Context, spec contracts.RunSpec, case contracts.GeneratedCase) (*contracts.EvaluationResult, error) {
    // 创建隔离临时目录
    workDir, err := os.MkdirTemp("", "eval_"+case.SampleID+"_")
    if err != nil {
        return nil, err
    }
    defer os.RemoveAll(workDir)  // 评测完成后清理
    
    // 复制源码和测试文件到工作目录
    copyFiles(workDir, case)
    
    // 在工作目录中执行评测阶段
    compileResult := s.compile(ctx, workDir, case.Language)
    testResult := s.runTests(ctx, workDir, case.Language)
    coverageResult := s.collectCoverage(ctx, workDir, case.Language)
    mutationResult := s.runMutation(ctx, workDir, case.Language)
    
    return buildResult(compileResult, testResult, coverageResult, mutationResult)
}
```

对于 repo_level 类型样本，workspace 不会被清理，保留用于后续分析。

#### 3.4.5 失败归因机制

Evaluator 自动分析失败原因，便于定位问题：

```go
type FailureOrigin string

const (
    FailureOriginModel      FailureOrigin = "model"      // 模型生成代码问题
    FailureOriginEnvironment FailureOrigin = "environment" // 环境配置问题
    FailureOriginDataset    FailureOrigin = "dataset"    // 数据集样本问题
    FailureOriginTool       FailureOrigin = "tool"       // 评测工具问题
    FailureOriginNone       FailureOrigin = "none"       // 无失败
)
```

归因规则示例：
- 编译失败 + 语法错误 → `model`（模型生成错误代码）
- 编译失败 + 缺少依赖 → `environment`（环境配置问题）
- 测试失败 + 导入错误 → `model`（模型生成错误导入）
- 变异工具崩溃 → `tool`（工具链问题）

### 3.5 报告生成模块 (internal/reporter/)

#### 3.5.1 核心职责

- **多维度聚合**：按模型、语言、场景等维度聚合统计数据
- **模型排名**：计算综合得分，生成模型排名列表
- **洞察分析**：自动生成最佳模型、弱项场景、改进建议等洞察
- **HTML报告**：生成包含 Chart.js 图表的可视化 HTML 报告

#### 3.5.2 综合得分计算

采用加权计算模型综合得分：

```
CompositeScore = 0.30 × CompilePassRate + 
                 0.30 × TestCasePassRate + 
                 0.20 × AvgLineCoverage + 
                 0.20 × AvgMutationScore
```

权重含义：
- **编译通过率 30%**：基础能力，生成的代码必须能编译
- **测试通过率 30%**：核心能力，测试必须能正确运行
- **行覆盖率 20%**：测试充分性，覆盖更多代码路径
- **变异得分 20%**：测试有效性，真正发现潜在缺陷

#### 3.5.3 多维度分析

Reporter 从多个维度分析评测结果：

```go
type Dimensions struct {
    ByModel         []ModelDim         // 按模型维度
    ByLanguage      []LanguageDim      // 按语言维度
    ByScenario      []ScenarioDim      // 按场景维度
    ByModelScenario []ModelScenarioDim // 模型+场景交叉
}
```

每个维度统计的关键指标：
- 编译通过率
- 测试通过率（样本级 + 用例级）
- 平均行覆盖率、分支覆盖率
- 平均变异得分
- 平均延迟、Token消耗

#### 3.5.4 洞察分析生成

Reporter 自动生成洞察结论：

```go
type Insights struct {
    BestModel      InsightItem   // 最佳模型
    WeakScenarios  []InsightItem // 弱项场景
    StrongScenarios []InsightItem // 强项场景
    LanguageGaps   []InsightItem // 语言差异
    Recommendations []InsightItem // 改进建议
    BenchmarkNotes []InsightItem // 评测说明
}
```

洞察示例：
- **最佳模型**：`deepseek 综合得分最高（82.5），变异得分领先`
- **弱项场景**：`complex_dependency 场景覆盖率偏低，建议关注复杂依赖处理`
- **语言差异**：`Go 变异得分普遍低于 Python，可能因工具链差异`
- **改进建议**：`建议提高 max_tokens 参数以减少截断率`

#### 3.5.5 HTML报告结构

生成的 `report.html` 包含：

1. **头部概览**：运行配置、总样本数、关键指标汇总
2. **模型排名表**：排名、得分、各项指标对比
3. **多维度图表**：
   - 模型对比柱状图（编译率、覆盖率、变异得分）
   - 场景对比雷达图
   - 语言对比热力图（已删除）
4. **截断统计**：截断率、按模型/语言分布
5. **失败分析**：常见错误类型、示例错误信息
6. **效率统计**：Token效率排名、时间效率排名
7. **洞察结论**：自动生成的分析结论

### 3.6 数据持久化模块 (internal/store/)

#### 3.6.1 核心职责

- **数据入库**：将评测产物（manifest、evaluation、report）入库 SQLite
- **历史查询**：支持查询历史运行、评测结果
- **复用判断**：查找可复用的历史生成结果（相同源码+模型+prompt）
- **跨运行对比**：支持从数据库选择多运行数据生成对比报告

#### 3.6.2 SQLite V2 数据库架构

数据库 Schema Version：`store.v2`

核心表结构：

```
┌──────────────────────────────────────────────────────────────────┐
│                        核心业务表                                 │
├──────────────────────────────────────────────────────────────────┤
│ experiments        实验管理（实验ID、名称、描述）                 │
│ generation_runs    生成运行记录（RunID、Spec、Prompt信息）        │
│ generated_cases    生成案例详情（模型、语言、样本、Token消耗）    │
│ evaluation_runs    评测运行记录（环境指纹、ScorePolicy）          │
│ evaluation_results 评测结果详情（覆盖率、变异得分、失败归因）     │
│ report_snapshots   报告快照（JSON/HTML artifact关联）             │
├──────────────────────────────────────────────────────────────────┤
│                        数据集管理表                               │
├──────────────────────────────────────────────────────────────────┤
│ dataset_samples    样本定义（语言、类别、场景、SHA256）           │
│ dataset_snapshots  数据集快照（指纹、样本数量）                   │
│ dataset_snapshot_members 快照成员关联                            │
├──────────────────────────────────────────────────────────────────┤
│                        配置与元数据表                             │
├──────────────────────────────────────────────────────────────────┤
│ model_configs      模型配置（Provider、APIEndpoint、参数）        │
│ prompt_profiles    提示词配置（策略、版本ID）                     │
│ prompt_renderings  提示词渲染记录（实际使用的Prompt）             │
│ evaluation_envs    评测环境（环境指纹、详情JSON）                 │
│ score_policies     评分策略（权重配置）                           │
├──────────────────────────────────────────────────────────────────┤
│                        产物管理表                                 │
├──────────────────────────────────────────────────────────────────┤
│ artifacts          产物记录（类型、路径、SHA256、大小）           │
│ run_artifacts      Run与Artifact关联（角色：manifest/source等）   │
│ evaluation_stage_results 评测阶段详情（命令、耗时、状态）         │
│ log_events         日志事件记录                                  │
└──────────────────────────────────────────────────────────────────┘
```

#### 3.6.3 入库流程

三种入库入口：

```go
// 1. 入库生成清单（含样本快照、模型配置、提示词记录）
func (s *SQLiteStore) IngestManifestFile(ctx context.Context, path string) (IngestSummary, error)

// 2. 入库评测结果（含环境指纹、ScorePolicy、阶段详情）
func (s *SQLiteStore) IngestEvaluationFile(ctx context.Context, path string) (IngestSummary, error)

// 3. 入库报告快照
func (s *SQLiteStore) IngestReportFile(ctx context.Context, path string) (IngestSummary, error)

// 4. 入库整个Run目录（自动检测并入库三种文件）
func (s *SQLiteStore) IngestRun(ctx context.Context, opts IngestRunOptions) (IngestSummary, error)
```

#### 3.6.4 生成结果复用机制

当 `ReuseGenerated: true` 时，Runner 会查询数据库寻找可复用的历史生成：

```go
func (s *SQLiteStore) FindReusableGeneratedCase(ctx context.Context, 
    model, language, sampleID, sourceSHA256, promptVersionID string) (ReusableGeneratedCase, bool, error) {
    // 查询条件：
    // 1. 相同模型 (model)
    // 2. 相同语言 (language)
    // 3. 相同样本ID (sample_id)
    // 4. 相同源码SHA256 (source_sha256)
    // 5. 相同Prompt版本 (prompt_version_id)
    // 6. 生成成功 (success=1)
    // 7. 测试文件未删除
    
    // 返回最早的成功生成记录
}
```

复用机制价值：
- **节省API调用成本**：避免重复调用相同模型
- **加速评测流程**：跳过生成阶段，直接进入评测
- **确保公平性**：所有复用结果来自相同Prompt版本

---

## 四、模型配置系统

### 4.1 配置文件结构

模型配置位于 `configs/models.yaml`，采用 YAML 格式：

```yaml
models:
  deepseek-v4-flash:
    enabled: true
    provider: deepseek
    config:
      api_endpoint: "https://api.deepseek.com"
      model: "deepseek-v4-flash"
      api_key_env: "DEEPSEEK_API_KEY"
      parameters:
        temperature: 0.7
        top_p: 0.9
        max_tokens: 4096

  minimax2.7:
    enabled: true
    provider: minimax
    config:
      api_endpoint: "https://api.minimax.chat/v1"
      model: "MiniMax-M2.7"
      api_key_env: "MINIMAX_API_KEY"
      parameters:
        temperature: 0.7
        top_p: 0.9
        max_tokens: 4096

benchmark:
  thresholds:
    compile_rate: 0.7
    execution_rate: 0.70
    line_coverage: 0.7
    branch_coverage: 0.6
    function_coverage: 0.8
    mutation_score: 0.70
  timeouts:
    api_call: 120000
    compilation: 60000
    test_execution: 30000
    mutation_test: 120000
  parallel:
    max_concurrent_models: 6
    max_concurrent_samples: 8

languages:
  python:
    compiler: "pylint"
    test_framework: "pytest"
    coverage_tool: "coverage"
    mutation_tool: "mutmut"
  go:
    compiler: "go build"
    test_framework: "go test"
    coverage_tool: "go tool cover"
    mutation_tool: "go-mutesting"
```

### 4.2 支持的模型

| 模型标识 | Provider | 实际型号 | API Key 环境变量 |
|----------|----------|----------|------------------|
| deepseek-v4-flash | deepseek | deepseek-v4-flash | DEEPSEEK_API_KEY |
| minimax2.7 | minimax | MiniMax-M2.7 | MINIMAX_API_KEY |
| minimax2.5 | minimax | minimax-M2.5 | MINIMAX2.5_API_KEY |
| doubao-seed | volcengine | doubao-seed-2.0-pro | VOLCENGINE_API_KEY |
| glm-5 | dashscope | glm-5 | DASHSCOPE_API_KEY |
| qwen3.6-flash | dashscope | qwen3.6-flash | DASHSCOPE_API_KEY |
| deepseek-v3.2 | volcengine | ep-20260425131055-4q4ck | ARK_API_KEY |
| doubao-seed-1.6-flash | volcengine | ep-20260425130953-qr2sj | ARK_API_KEY |
| doubao-seed-1.8 | volcengine | ep-20260425130558-gs24c | ARK_API_KEY |
| doubao-seed-2.0-code | volcengine | ep-20260425130248-lqxth | ARK_API_KEY |

### 4.3 Provider 类型

| Provider | API 格式 | 特点 |
|----------|----------|------|
| deepseek | OpenAI兼容 | 直接调用 DeepSeek API |
| minimax | OpenAI兼容 | MiniMax 平台 API |
| volcengine | OpenAI兼容 | 火山引擎平台，需要 endpoint ID |
| dashscope | OpenAI兼容 | 阿里云百炼平台 |

---

## 五、CLI 命令系统

### 5.1 命令概览

CLI 入口定义于 `cmd/utbench/main.go`：

```
utbench - unified test-bench CLI

Usage:
  utbench run          Run full pipeline (generate → evaluate → report)
  utbench generate     Generate unit tests only
  utbench evaluate     Evaluate existing generated tests
  utbench report       Generate reports from evaluation results
  utbench db           Manage SQLite benchmark database
  utbench dataset      Dataset management (index, manifest, stats)
  utbench doctor       Check evaluator toolchains with canary tests
  utbench web          Launch Web management UI
  utbench help         Show this help
```

### 5.2 run 命令详解

完整流水线执行命令：

```bash
utbench run \
  --models deepseek,qwen \
  --langs python,go \
  --config ./configs/models.yaml \
  --dataset-root ./datasets \
  --output-root ./artifacts \
  --max-samples 20 \
  --workers 16 \
  --mutation-enabled \
  --mutation-timeout 600 \
  --test-timeout 180 \
  --mode full \
  --ingest \
  --db-path ./storage/utbench.db
```

关键参数说明：

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--models` | 逗号分隔的模型列表 | 无（必须指定） |
| `--langs` | 逗号分隔的语言列表 | 无（必须指定） |
| `--config` | 模型配置文件路径 | `../benchmark/config/models.yaml` |
| `--dataset-root` | 数据集根目录 | `./datasets` |
| `--output-root` | 输出根目录 | `./artifacts` |
| `--max-samples` | 最大样本数限制 | 0（不限制） |
| `--workers` | 并发 Worker 数量 | 16 |
| `--mutation-enabled` | 是否启用变异测试 | true |
| `--mutation-timeout` | 变异测试超时（秒） | 600 |
| `--test-timeout` | 测试执行超时（秒） | 180 |
| `--mode` | 运行模式 | `full` |
| `--dry-run` | 试运行（不调用API） | false |
| `--reuse-generated` | 复用历史生成 | false |
| `--ingest` | 入库 SQLite | false |
| `--db-path` | SQLite 数据库路径 | `./storage/utbench.db` |

### 5.3 db 命令详解

数据库管理命令集：

```bash
# 初始化数据库
utbench db init --db-path ./storage/utbench.db

# 入库生成清单
utbench db ingest-manifest --manifest ./artifacts/runs/<run-id>/generated/generated_manifest.json --db-path ./storage/utbench.db

# 入库评测结果
utbench db ingest-evaluation --evaluation ./artifacts/runs/<run-id>/evaluation/evaluation_result.json --db-path ./storage/utbench.db

# 入库整个Run目录
utbench db ingest-run --run-id <run-id> --output-root ./artifacts --db-path ./storage/utbench.db

# 查看数据库概览
utbench db overview --db-path ./storage/utbench.db --limit 10

# 列出运行记录
utbench db list-runs --db-path ./storage/utbench.db --limit 50

# 从数据库生成对比报告
utbench db report \
  --run-ids run_001,run_002 \
  --models deepseek,qwen \
  --langs python \
  --score-eligible-only \
  --db-path ./storage/utbench.db \
  --output-root ./artifacts
```

### 5.4 doctor 命令详解

环境诊断命令，检查工具链和运行 Canary 测试：

```bash
utbench doctor \
  --langs python,go,java,cpp \
  --mutation-enabled \
  --mutation-timeout 120 \
  --test-timeout 60
```

检查内容：
- **工具链检查**：检查编译器、测试框架、覆盖率工具、变异工具是否安装
- **Canary 测试**：运行简单的示例测试，验证完整流水线

输出示例：
```
Doctor: OK
  [tool ok] python: Python 3.11.4
  [tool ok] pytest: pytest 7.4.0
  [tool ok] coverage: coverage 7.3.0
  [tool ok] mutmut: mutmut 2.4.3
  [canary] python compile=true test=true coverage=true mutation=true
```

---

## 六、Web UI 系统

### 6.1 Web UI 功能

Web 管理界面提供可视化操作入口：

- **总览页**：任务统计、Docker/镜像状态
- **新建任务**：交互式表单，勾选模型、语言、场景等参数
- **任务列表**：历史任务、状态过滤
- **任务详情**：实时日志、评测报告、横向对比图表
- **数据管理**：跨运行对比、历史数据入库

### 6.2 启动 Web UI

```bash
cd go-ut-bench
go run ./cmd/utbench web --addr :8080
```

打开浏览器访问 `http://localhost:8080`。

### 6.3 Web API 端点

定义于 `internal/web/server.go`：

| 端点 | 方法 | 功能 |
|------|------|------|
| `/api/config` | GET | 获取模型配置 |
| `/api/env` | GET | 获取环境变量状态 |
| `/api/runs` | GET | 获取运行列表 |
| `/api/runs/:id` | GET | 获取运行详情 |
| `/api/runs/:id/logs` | GET (SSE) | 实时日志流 |
| `/api/db/overview` | GET | 数据库概览 |
| `/api/db/runs` | GET | 数据库运行列表 |
| `/api/db/results` | GET | 数据库评测结果 |
| `/api/start` | POST | 启动新运行 |
| `/api/build` | POST | 构建 Docker 镜像 |

---

## 七、产物目录结构

### 7.1 输出目录结构

每次运行产物位于 `artifacts/runs/<run-id>/`：

```
artifacts/
  runs/
    <run-id>/
      generated/
        generated_manifest.json     # 生成清单
        <model>/
          <language>/
            <sample_id>_test.py     # 生成的测试文件
            <sample_id>_response.json  # API响应（调试用）
            <sample_id>_metadata.json  # 元数据
      evaluation/
        evaluation_result.json      # 评测结果
      report/
        report_summary.json         # 报告汇总JSON
        report.html                 # 可视化HTML报告
      logs/
        orchestrator.log            # 编排日志
        runner.log                  # 生成日志
        evaluator.log               # 评测日志
      run_summary.json              # 运行汇总
  checkpoints/
    runner_<hash>.checkpoint.json   # Runner断点文件
```

### 7.2 Prompt 快照目录

当运行时，Prompt 模板会被快照保存：

```
artifacts/
  runs/
    <run-id>/
      prompts/
        system.txt                  # 系统消息
        python_full_file.prompt.txt # Python完整文件模式模板
        python_completion.prompt.txt
        python_repo_level.prompt.txt
        go_full_file.prompt.txt
        ...
        prompt_catalog.json         # Prompt目录索引
```

---

## 八、Docker 支持

### 8.1 Docker 镜像构建

推荐使用 Docker 运行，确保工具链环境一致：

```bash
cd go-ut-bench
docker build -t utbench:latest .
```

Dockerfile 包含：
- Go 编译环境
- Python 3 + pytest + coverage + mutmut
- JDK 17 + Maven + PITest
- g++ + GoogleTest + gcov + mull

### 8.2 Docker 运行示例

```bash
# Linux/macOS
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -v "$(pwd)/storage:/app/storage" \
  utbench:latest run \
    --models deepseek,qwen \
    --langs python,go \
    --config /app/configs/models.yaml \
    --max-samples 20 \
    --ingest --db-path /app/storage/utbench.db

# Windows PowerShell
docker run --rm --env-file .env `
  -v "${PWD}/datasets:/app/datasets" `
  -v "${PWD}/artifacts:/app/artifacts" `
  -v "${PWD}/configs:/app/configs" `
  -v "${PWD}/storage:/app/storage" `
  utbench:latest run `
    --models deepseek `
    --langs python `
    --config /app/configs/models.yaml `
    --max-samples 5
```

---

## 九、评测指标详解

### 9.1 编译通过率

**定义**：生成的测试代码能否通过编译/语法检查

**语言对应检查方式**：
- Python：pylint 语法检查
- Go：go build
- Java：javac
- C++：g++

**意义**：衡量模型生成语法正确代码的基础能力

### 9.2 测试通过率

**两层定义**：

1. **样本级通过率**：测试文件是否全部通过（无失败用例）
2. **用例级通过率**：通过的测试用例数 / 总测试用例数

**意义**：衡量模型生成测试的正确性，能否真正验证源码行为

### 9.3 覆盖率

**行覆盖率**：被执行的代码行数 / 总代码行数

**分支覆盖率**：被执行的分支数 / 总分支数

**意义**：衡量测试的充分性，覆盖更多代码路径意味着更全面的测试

### 9.4 变异测试得分

**定义**：被杀死的变异体数 / 总变异体数

**变异体**：对源码进行微小语义修改（如 `+` 改 `-`，`>` 改 `>=`）

**杀死条件**：测试用例在变异后产生不同的输出（失败）

**意义**：衡量测试的有效性，真正能够发现潜在缺陷的测试才是有价值的测试

**变异工具输出详解**：

| 状态 | 说明 |
|------|------|
| killed | 测试检测到变异 |
| survived | 测试未检测到变异（测试无效） |
| no_tests | 该变异体无对应测试 |
| timeouts | 变异测试超时 |
| suspicious | 变异后测试结果可疑 |
| skipped | 被跳过的变异体 |

---

## 十、数据集设计

### 10.1 数据集类别

| 类别 | 特点 | 示例 |
|------|------|------|
| `self_contained` | 单文件自包含，无外部依赖 | 简单函数、边界值测试 |
| `repo_level` | 多文件仓库级，需要 workspace 上下文 | python-dateutil、第三方库测试 |

### 10.2 数据集场景

| 场景 | 难度 | 特点 |
|------|------|------|
| `boundary` | 中等 | 边界值测试，需要处理极限情况 |
| `simple_function` | 低 | 简单函数，基础测试能力 |
| `complex_dependency` | 高 | 复杂依赖，需要 mock/stub |
| `interface_mock` | 高 | 接口 mock，需要理解抽象依赖 |

### 10.3 Module Level 样本元数据

repo_level 类型样本需要 `meta.json` 提供额外上下文：

```json
{
    "sample_id": "dateutil_parser",
    "module_import": "dateutil.parser",
    "package_name": "dateutil",
    "target_file": "parser.py",
    "workspace_root": "/path/to/dateutil_workspace",
    "requirements": ["python-dateutil"]
}
```

---

## 十一、环境变量配置

创建 `go-ut-bench/.env` 文件配置 API 密钥：

```dotenv
# DeepSeek
DEEPSEEK_API_KEY=sk-xxx

# 阿里云百炼（Qwen、GLM-5）
DASHSCOPE_API_KEY=sk-xxx

# MiniMax
MINIMAX_API_KEY=xxx
MINIMAX2.5_API_KEY=xxx

# 火山引擎原始
VOLCENGINE_API_KEY=xxx

# 火山引擎 ARK（doubao-seed系列、deepseek-v3.2）
ARK_API_KEY=xxx
```

---

## 十二、常见使用场景

### 12.1 验证环境（试运行）

```bash
docker run --rm \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest run \
    --models deepseek \
    --langs python \
    --config /app/configs/models.yaml \
    --max-samples 1 \
    --dry-run
```

### 12.2 多模型横向对比

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  -v "$(pwd)/storage:/app/storage" \
  utbench:latest run \
    --run-id compare_001 \
    --models deepseek,qwen,minimax,glm-5 \
    --langs python,go \
    --config /app/configs/models.yaml \
    --max-samples 20 \
    --ingest --db-path /app/storage/utbench.db
```

### 12.3 完整评测（含变异测试）

```bash
docker run --rm --env-file .env \
  -v "$(pwd)/datasets:/app/datasets" \
  -v "$(pwd)/artifacts:/app/artifacts" \
  -v "$(pwd)/configs:/app/configs" \
  utbench:latest run \
    --models deepseek \
    --langs python \
    --config /app/configs/models.yaml \
    --max-samples 5 \
    --mutation-enabled \
    --mutation-timeout 600
```

### 12.4 断点续跑

```bash
# 第一次运行（中途停止）
docker run ... --run-id resume_test --max-samples 100

# 续跑（已完成的不重复执行）
docker run ... --run-id resume_test --mode incremental --max-samples 100
```

### 12.5 复用历史生成

```bash
# 复用数据库中相同源码+模型+Prompt的历史生成
docker run ... --reuse-generated --db-path /app/storage/utbench.db
```

---

## 十三、项目目录结构总览

```
go-ut-bench/
  cmd/utbench/                  # CLI 入口
    main.go                      # 命令路由和参数解析

  internal/                      # 内部模块（不对外暴露）
    contracts/                   # 数据契约（跨模块数据结构）
      spec.go                    # RunSpec, SampleRef, RepoLevelMeta
      results.go                 # GeneratedManifest, EvaluationResult, ReportPayload
      constants.go               # SchemaVersion, DatasetClass, RunMode

    orchestrator/                # 流水线编排
      service.go                 # Phase控制，协调各模块

    dataset/                     # 数据集发现与过滤
      service.go                 # ValidateSpec, DiscoverSamples, BuildIndex

    runner/                      # 测试生成
      service.go                 # Generate, WorkerPool, Checkpoint
      prompt.go                  # 三种Prompt模式，PromptCatalog
      api.go                     # LLM API调用适配

    evaluator/                   # 评测执行
      service.go                 # Evaluate, 四阶段流水线
      python.go                  # Python评测实现
      go.go                      # Go评测实现
      java.go                    # Java评测实现
      cpp.go                     # C++评测实现

    reporter/                    # 报告生成
      service.go                 # Generate, 多维度聚合，洞察分析
      html.go                    # HTML模板，Chart.js图表

    store/                       # 数据持久化
      sqlite.go                  # V2 Schema，入库，复用查找

    web/                         # Web UI
      server.go                  # HTTP路由，SSE日志流
      run_manager.go             # Docker运行管理

    obs/                         # 可观测性
      logger.go                  # slog封装，日志输出

  configs/                       # 配置文件
    models.yaml                  # 模型配置（Provider、APIEndpoint、参数）

  datasets/                      # 数据集目录
    python/
      self_contained/
        boundary/
        simple_function/
      repo_level/
    go/
    java/
    cpp/

  artifacts/                     # 运行产物
    runs/
      <run-id>/
        generated/
        evaluation/
        report/
        logs/
    checkpoints/

  storage/                       # SQLite数据库
    utbench.db

  .env                           # API密钥配置
  Dockerfile                     # Docker镜像构建
  docs/01-user-guides/USER_GUIDE.md        # 快速入门
  docs/01-user-guides/startup-guide.md     # 启动指南
```

---

## 十四、技术亮点总结

1. **端到端自动化流水线**：从数据集发现到报告生成，全程无需人工干预
2. **多维度评测指标**：编译、测试、覆盖率、变异测试四维度全面评估
3. **公平横向对比**：相同Prompt、相同样本、相同环境，确保公平性
4. **断点续跑机制**：Checkpoint支持中断后继续，适合大规模评测
5. **生成结果复用**：数据库查询相同条件历史生成，节省API成本
6. **自动失败归因**：智能分析失败原因（model/environment/dataset/tool）
7. **可视化洞察报告**：Chart.js图表 + 自动洞察分析，便于决策
8. **跨运行对比能力**：SQLite持久化支持历史数据查询和对比报告
9. **多语言支持**：Python、Go、Java、C++ 四种主流语言
10. **Docker环境隔离**：预装完整工具链，确保评测环境一致性

---

## 十五、版本与演进

- **数据结构版本**：`v0.1.0`（contracts.SchemaVersion）
- **数据库Schema版本**：`store.v2`
- **提示词策略**：`structured-v1`

项目持续演进，主要改进方向：
- 增加更多 LLM Provider 支持
- 扩展数据集场景覆盖
- 优化评测流水线效率
- 增强洞察分析深度

---

*本文档基于 go-ut-bench 源代码深入分析编写，涵盖项目的核心架构、数据流程、模块设计、配置系统和使用方法。*
