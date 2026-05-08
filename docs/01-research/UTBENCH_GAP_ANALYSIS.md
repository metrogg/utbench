# UT-Bench 差距分析：当前实现 vs 行业顶尖项目

> 分析时间：2026-05-03
> 对照项目：SWE-bench/Pro、SWE-agent、OpenHands、Aider、DeepEval、Langfuse、mini-swe-agent

---

## 总览评分

| 维度 | 当前水平 | 行业顶尖 | 差距等级 |
|:---|:---|:---|:---|
| 评测指标体系 | ★★★★☆ | ★★★★★ | 🟡 中 |
| 数据集设计 | ★★☆☆☆ | ★★★★★ | 🔴 大 |
| Agent 沙箱与隔离 | ★★★☆☆ | ★★★★★ | 🔴 大 |
| 评测公平性/可复现 | ★★★★☆ | ★★★★★ | 🟡 中 |
| 可观测性/Trace | ★★★☆☆ | ★★★★★ | 🔴 大 |
| 成本追踪 | ★★★☆☆ | ★★★★☆ | 🟡 中 |
| 评测流程健壮性 | ★★★☆☆ | ★★★★★ | 🔴 大 |
| Prompt 工程 | ★★★★☆ | ★★★★★ | 🟡 中 |
| 报告/排行榜 | ★★★★☆ | ★★★★★ | 🟢 小 |

---

## 一、数据集设计（🔴 差距最大）

### 行业标准（SWE-bench Pro）

```
每个 task 实例包含:
├── instance_id          唯一标识
├── repo                 仓库地址
├── base_commit          精确的 commit hash
├── problem_statement    问题描述
├── hints_text           讨论上下文
├── patch                黄金修复 patch
├── test_patch           测试 patch
├── FAIL_TO_PASS         必须通过的测试列表
├── PASS_TO_PASS         必须继续通过的测试列表
├── version              版本
└── environment_setup_commit  环境锁定 commit
```

**关键设计**：
- **FAIL_TO_PASS / PASS_TO_PASS** — 双重验证契约，防止"破窗"
- **base_commit 锁定** — 完全可复现的仓库快照
- **environment_setup_commit** — 环境版本锁定

### 当前 UT-Bench

```go
// self_contained: 单文件，无上下文
type SampleRef struct {
    ID, Language, Category, Scenario, Path, SourceMD5
}

// repo_level: 极简 meta
type RepoLevelMeta struct {
    SampleID, ModuleImport, PackageName, TargetFile, WorkspaceRoot, Requirements
}
```

### 差距清单

| # | 差距 | 影响 | 优先级 |
|:---|:---|:---|:---|
| D1 | **无 FAIL_TO_PASS / PASS_TO_PASS 契约** | Agent 生成的测试可能破坏已有测试而不被发现 | P0 |
| D2 | **无 base_commit 锁定** | 仓库级样本不可复现，依赖的包版本漂移 | P0 |
| D3 | **无环境锁定（environment_setup_commit）** | 同一样本在不同时间跑可能结果不同 | P0 |
| D4 | **repo_level meta 信息太薄** | 缺少 relevant_files（上下文文件）、dependency_level、difficulty 等关键字段 | P1 |
| D5 | **样本量太少** | self_contained 800个，repo_level 仅 1 个雏形（dateutil） | P1 |
| D6 | **无难度分级标注** | SWE-bench Pro 按修改文件数/行数分级，UT-Bench 无法按难度分析 | P1 |
| D7 | **无污染防护** | 数据集全部来自公开代码，存在训练数据污染风险 | P2 |

### 优化方案

```json
// 提议的 repo_level meta.json
{
  "sample_id": "django_orm_001",
  "source": {
    "repo_url": "https://github.com/django/django",
    "base_commit": "419b7d...",
    "license": "BSD-3"
  },
  "target": {
    "target_files": ["django/db/models/query.py"],
    "target_classes": ["QuerySet"],
    "target_functions": ["annotate", "filter"]
  },
  "context": {
    "context_scope": "module",
    "dependency_level": 3,
    "relevant_files": ["django/db/models/sql/query.py", "django/db/models/manager.py"]
  },
  "environment": {
    "language": "python",
    "python_version": "3.11",
    "setup_commands": ["pip install -e ."],
    "test_runner": "pytest",
    "pinned_deps": {"django": "5.0.1", "pytest": "8.0.0"}
  },
  "evaluation": {
    "run_tests_command": "pytest tests/query/ --tb=short",
    "must_pass_existing": true,
    "existing_tests_pattern": "tests/query/test_query.py",
    "timeout_seconds": 120
  },
  "difficulty": {
    "dependency_level": 3,
    "scenario": "complex_dependency",
    "estimated_loc": 150,
    "tags": ["orm", "database", "query-builder"]
  }
}
```

---

## 二、Agent 沙箱与隔离（🔴 差距大）

### 行业标准

| 项目 | 沙箱能力 |
|:---|:---|
| **SWE-agent** | SWE-ReX 抽象层：`run_command()`, `read_file()`, `write_file()`, `close()`；Docker + 沙箱策略 |
| **OpenHands** | 4 种 Runtime 后端：Docker / Local / Remote / E2B；EventStream 事件流；容器资源监控 |
| **Aider** | Docker 容器 + git worktree 隔离；自动回滚 |
| **Terminal-Bench** | Docker 内真实终端；限制交互轮次 |

### 当前 UT-Bench

```go
type SandboxRunRequest struct {
    Mode, Workspace, Command, DockerImage, NetworkDisabled, CPU, Memory, TimeoutSeconds
}
// 仅 2 种模式: docker / local
// 无 seccomp/AppArmor
// 无资源监控
// 无自动回滚
// 无交互轮次限制
```

### 差距清单

| # | 差距 | 影响 | 优先级 |
|:---|:---|:---|:---|
| D8 | **源码文件可被 Agent 修改** | Agent 可能改源码让测试通过，而非生成正确的测试 | P0 |
| D9 | **无交互轮次限制** | Agent 可能无限循环消耗资源 | P0 |
| D10 | **无容器资源监控** | OOM/高 CPU 不被发现，静默失败 | P1 |
| D11 | **无自动回滚机制** | Agent 修改了不该改的文件无法恢复 | P1 |
| D12 | **仅 Docker/Local 两种后端** | 无法适配 E2B/腾讯云 Cube 等托管沙箱 | P2 |
| D13 | **孤儿容器无自动清理** | 取消任务后容器可能残留 | P2 |

### 优化方案

1. **源码只读挂载**：`-v {source}:/workspace/src:ro`，Agent 只能在 `/workspace/tests/` 写入
2. **交互轮次上限**：参考 SWE-bench Pro 的 250 轮，UT-Bench 可设 50 轮
3. **资源监控**：`docker stats` 采样，记录 peak_memory / peak_cpu
4. **自动回滚**：执行前 `git init + git add . + git commit`，失败后 `git checkout .`
5. **SandboxRunner 接口扩展**：增加 `Stats()` 方法返回资源使用数据

---

## 三、可观测性/Trace（🔴 差距大）

### 行业标准

| 项目 | Trace 能力 |
|:---|:---|
| **SWE-agent** | 完整 trajectory：每个 action（command）+ observation（输出）结构化记录 |
| **OpenHands** | EventStream：Action/Observation 事件流，支持多模态 |
| **Langfuse** | Trace → Span → Generation 三级归因；cost per span；prompt 版本追踪 |
| **DeepEval** | LLM-as-Judge 评估 trace 质量 |

### 当前 UT-Bench

```
AgentTrace 已收集:
├── InteractionCount, ToolCallCount
├── FilesRead, FilesWritten, CommandsExecuted
├── Stdout/Stderr（原始文本）
├── Workspace diff（JSON）
└── Session export（OpenCode 专属）

缺失:
├── 无结构化 trajectory（action-by-action）
├── 无 cost per action 归因
├── 无 LLM-as-Judge 评估 trace 质量
└── 无 Trace 可视化
```

### 差距清单

| # | 差距 | 影响 | 优先级 |
|:---|:---|:---|:---|
| D14 | **无结构化 trajectory** | 无法分析 Agent "为什么失败"——只看到最终结果 | P0 |
| D15 | **无 cost per action 归因** | 无法知道哪个工具调用/哪轮对话最耗 token | P1 |
| D16 | **无 LLM-as-Judge 评估** | 无法自动评估 Agent 过程质量 | P1 |
| D17 | **Trace 不可视化** | 难以直观理解 Agent 行为模式 | P2 |

### 优化方案

1. **统一 Trajectory 格式**（兼容 SWE-agent 格式）：
```json
{
  "trajectory": [
    {"step": 1, "action": "command", "content": "cat src/main.py", "observation": "...", "duration_ms": 120},
    {"step": 2, "action": "write_file", "content": "test_main.py", "observation": "File written", "duration_ms": 50}
  ]
}
```

2. **Cost 归因**：每步记录 input_tokens / output_tokens / estimated_cost
3. **LLM-as-Judge 评分**：用小模型对 trajectory 评分（效率/冗余/策略合理性）

---

## 四、评测公平性/可复现（🟡 中等差距）

### 行业标准

| 实践 | 来源 |
|:---|:---|
| **Scaffolding 版本标注** | SWE-bench Pro：同一个模型在不同脚手架上差 14 分 |
| **PASS_TO_PASS 验证** | SWE-bench：生成的测试不能破坏已有测试 |
| **pass_rate_n 多轮通过率** | Aider：给 Agent 多次机会 |
| **Pinned dependency** | SWE-bench：environment_setup_commit 锁定 |
| **公开 Leaderboard** | benchlm.ai / SWE-bench Leaderboard |

### 当前 UT-Bench 已有

- ✅ Environment fingerprint（SHA256）
- ✅ Generation/evaluation key 复用
- ✅ Score eligibility 排除环境/数据/工具故障
- ✅ Prompt catalog + version ID

### 差距清单

| # | 差距 | 影响 | 优先级 |
|:---|:---|:---|:---|
| D18 | **无 must_pass_existing 验证** | Agent 生成的新测试可能破坏旧测试 | P0 |
| D19 | **无 pass_rate_n 多轮通过率** | 只看单次结果，不反映 Agent 迭代能力 | P1 |
| D20 | **无 Scaffolding 版本标注** | 报告中未明确标注 Agent 框架版本 | P1 |
| D21 | **依赖版本未 pin** | pytest/coverage.py 版本变化影响结果 | P1 |
| D22 | **无公开 Leaderboard** | 无法与社区对比 | P2 |

### 优化方案

1. **must_pass_existing**：评测阶段先跑已有测试，全通过后才评新生成测试
2. **pass_rate_1/2/3**：给 Agent 最多 3 次机会（看到失败输出后重试）
3. **Scaffolding 标注**：报告中固定输出 `framework + version + prompt_strategy + prompt_version`
4. **Dockerfile 中 pin 依赖版本**：`pip install pytest==8.0.0 coverage==7.4.0 mutmut==2.5.0`

---

## 五、评测指标体系（🟡 中等差距）

### 行业标准

| 指标 | 来源 | UT-Bench 状态 |
|:---|:---|:---|
| Pass/Fail | SWE-bench | ✅ CompilePass + TestPass |
| 行/分支覆盖率 | 行业通用 | ✅ LineCoverage + BranchCoverage |
| 变异测试 | 行业通用 | ✅ MutationScore |
| AssertFlip | AssertFlip 论文 | ❌ 缺失 |
| Flaky Test Rate | Google | ❌ 缺失 |
| Test Execution Rate | 行业通用 | ✅ TestPassRate |
| $/task | Aider | ✅ EstimatedCostUSD |
| tokens_per_coverage_point | 自创 | ❌ 缺失 |
| syntax_errors / error 分类 | Aider | ❌ 缺失（仅有 FailureOrigin） |
| LLM-as-Judge | DeepEval | ❌ 缺失 |
| 代码风格/可维护性 | DeepEval | ❌ 缺失 |

### 当前综合评分

```go
CompositeScore = Compile×0.3 + Test×0.3 + Coverage×0.2 + Mutation×0.2
```

**问题**：
- 权重硬编码，不可配置
- 无 Test Quality 维度（代码风格、断言质量、flaky 风险）
- 无效率维度（cost_per_coverage_point）

### 优化方案

1. **新增指标**：
   - `ExistingTestsStillPass: bool` — must_pass_existing 验证
   - `FlakyTestRate: float64` — 同一测试跑 3 次，不一致率
   - `SyntaxErrorRate: float64` — 编译失败中语法错误占比
   - `CostPerCoveragePoint: float64` — `$ / line_coverage`
   - `TokensPerTestCase: float64` — `total_tokens / test_count`

2. **可配置权重**：
```yaml
scoring:
  weights:
    compile: 0.2
    test: 0.25
    coverage: 0.2
    mutation: 0.2
    test_quality: 0.1
    efficiency: 0.05
```

---

## 六、成本追踪（🟡 中等差距）

### 行业标准（Langfuse）

```
Trace → Span → Generation
  ├── input_tokens / output_tokens per span
  ├── estimated_cost per span
  ├── model 信息
  └── latency per span
```

### 当前 UT-Bench

```go
// 已有
PromptTokens, CompletionTokens, TotalTokens
EstimatedCostUSD, CostSource, TokenSource

// 缺失
// - 无 cost per action 归因（仅 total）
// - 无 Provider 级别的定价自动更新
// - 无 token 预算控制
// - 无成本效率指标
```

### 优化方案

1. **Per-action cost**：CLI Agent 每轮工具调用都记录 token 使用
2. **Cost budget**：`--max-cost-per-sample 0.5`，超出则终止
3. **效率指标**：`cost_per_coverage_point`、`tokens_per_test_case`

---

## 七、评测流程健壮性（🔴 差距大）

### 行业标准

| 能力 | 来源 | UT-Bench |
|:---|:---|:---|
| API 重试 + 退避 | SWE-agent | ✅ 3 次指数退避 |
| 评测阶段重试 | 行业通用 | ❌ 无 |
| 断路器 | 行业通用 | ❌ 无 |
| 并发 per-provider 限流 | 行业通用 | ❌ 仅 per-model 间隔 200ms |
| Docker 清理 | OpenHands | ❌ 无自动清理 |
| 磁盘空间监控 | 行业通用 | ❌ 无 |
| 热重载/继续运行 | SWE-agent | ✅ checkpoint + incremental |
| 进度 Web UI | OpenHands | ❌ 仅 CLI |

### 差距清单

| # | 差距 | 影响 | 优先级 |
|:---|:---|:---|:---|
| D23 | **评测阶段无重试** | 一次 coverage 工具超时就丢失数据 | P0 |
| D24 | **无 per-provider 并发限流** | 多模型并行可能触发 rate limit | P1 |
| D25 | **无磁盘空间监控** | artifacts 撑满磁盘后静默失败 | P1 |
| D26 | **孤儿容器无清理** | Ctrl+C 后容器残留 | P2 |

---

## 八、Prompt 工程（🟡 中等差距）

### 当前 UT-Bench 的 Prompt 优点

- ✅ 三种模式（full_file / completion / repo_level）
- ✅ 语言特化规则（Python/Go/Java/C++ 各自不同）
- ✅ 自动提取 dependencies / critical conditions / mock requirements
- ✅ Prompt catalog + SHA1 version

### 差距

| # | 差距 | 行业做法 | 优先级 |
|:---|:---|:---|:---|
| D27 | **无 execution feedback loop** | SWE-agent：Agent 可跑测试看输出再修改 | P1 |
| D28 | **无 few-shot examples** | Aider：提供 passing test 样例 | P1 |
| D29 | **Coverage target 硬编码** | `line>=70%, branch>=60%, function>=80%` 写死在代码中 | P2 |
| D30 | **无 A/B prompt 实验** | 无法对比不同 prompt 策略的效果 | P2 |

### 优化方案

1. **多轮反馈**：CLI Agent 天然支持（可跑测试再修改），model_api 增加 `--feedback-rounds 2`
2. **Few-shot**：在 prompt 中加入 1-2 个 passing test 示例
3. **可配置 coverage target**：
```yaml
prompt:
  coverage_targets:
    line: 70
    branch: 60
    function: 80
```

---

## 九、优先级排序：Top 10 立即可做的优化

| 排名 | 优化项 | 差距编号 | 预期收益 | 工作量 |
|:---|:---|:---|:---|:---|
| **1** | 源码只读挂载，防止 Agent 篡改 | D8 | 评测公平性核心 | 0.5d |
| **2** | must_pass_existing 验证 | D1,D18 | 防止"破窗"测试 | 1d |
| **3** | 交互轮次上限（50 轮） | D9 | 防止无限循环 | 0.5d |
| **4** | 评测阶段增加重试（1 次） | D23 | 减少偶发失败 | 0.5d |
| **5** | repo_level meta.json 扩展（5 维度） | D4 | 仓库级数据集基础 | 1d |
| **6** | 依赖版本 pin（Dockerfile） | D21,D3 | 可复现性 | 0.5d |
| **7** | 结构化 trajectory 输出 | D14 | Agent 行为分析基础 | 2d |
| **8** | Flaky Test Rate 指标 | — | 测试质量评估 | 1d |
| **9** | Cost per action 归因 | D15 | 成本优化决策 | 1d |
| **10** | pass_rate_n 多轮通过率 | D19 | Agent 迭代能力评估 | 1d |

**总计约 8.5 人天**，完成后 UT-Bench 的核心短板将大幅改善。

---

## 十、中长期路线图

### Phase 1：补齐核心短板（1-2 周）
- [ ] 源码只读 + 交互上限 + must_pass_existing
- [ ] 评测重试 + 依赖 pin
- [ ] repo_level meta.json 5 维度 schema
- [ ] Flaky Test Rate 指标

### Phase 2：提升可观测性（2-4 周）
- [ ] 结构化 trajectory（兼容 SWE-agent 格式）
- [ ] Cost per action 归因
- [ ] 容器资源监控（peak CPU/Memory）
- [ ] pass_rate_n 多轮通过率

### Phase 3：构建差异化壁垒（1-3 月）
- [ ] LLM-as-Judge 评估 trajectory 质量
- [ ] AssertFlip 指标（反转测试检测 Bug）
- [ ] 第二个 framework 接入（Aider / Claude Code）
- [ ] 公开 Leaderboard
- [ ] 私有数据集 + 防污染机制

---

## 附录：与各项目的功能对比矩阵

| 功能 | SWE-bench/Pro | SWE-agent | OpenHands | Aider | DeepEval | **UT-Bench** |
|:---|:---|:---|:---|:---|:---|:---|
| 评测任务类型 | Issue 修复 | Issue 修复 | 通用 | 代码编辑 | LLM 输出 | **UT 生成** |
| 多语言 | Python+5 | Python | Python | 7 语言 | Python | **4 语言** |
| 覆盖率指标 | ❌ | ❌ | ❌ | ❌ | ❌ | **✅** |
| 变异测试 | ❌ | ❌ | ❌ | ❌ | ❌ | **✅** |
| Skill 评测 | ❌ | ❌ | ❌ | ❌ | ❌ | **✅** |
| FAIL/PASS 契约 | ✅ | ✅ | ❌ | ❌ | ❌ | ❌→✅ |
| 交互轮次限制 | ✅(250) | ✅ | ✅ | ✅ | N/A | ❌→✅ |
| 结构化 trajectory | ✅ | ✅ | ✅ | ❌ | ❌ | ❌→✅ |
| Cost 追踪 | ❌ | ❌ | ❌ | ✅ | ❌ | **✅**(total) |
| LLM-as-Judge | ❌ | ❌ | ❌ | ❌ | ✅ | ❌ |
| Docker 沙箱 | ✅ | ✅ | ✅ | ✅ | ❌ | ✅ |
| 源码只读 | ✅ | ✅ | ❌ | ✅ | N/A | ❌→✅ |
| 依赖版本锁定 | ✅ | ✅ | ✅ | ❌ | ❌ | ❌→✅ |
| Leaderboard | ✅ | ✅ | ❌ | ✅ | ❌ | ❌ |

**UT-Bench 的独特优势**：覆盖率 + 变异测试 + Skill 评测——这三项在所有竞品中都没有。补齐短板后，差异化价值将非常突出。
