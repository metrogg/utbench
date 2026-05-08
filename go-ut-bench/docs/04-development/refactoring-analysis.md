# go-ut-bench 重构分析报告

> 分析日期: 2026-05-03
> 分析范围: `go-ut-bench/internal/` 全部 13 个包 + `cmd/utbench/main.go`
> 目标: 为遵循 SOLID 原则的系统性重构提供依据

---

## 一、项目现状概览

### 1.1 规模数据

| 指标 | 数值 |
|------|------|
| Go 源文件数 | 62 |
| 总代码行数 | ~30,000 |
| internal 包数 | 13 |
| 最大文件 | `reporter/html_sections.go` (1,549 行) |
| 最长函数 | `evaluator/service.go:evaluateOne` (584 行) |
| CLI 入口 | `cmd/utbench/main.go` (1,610 行) |

### 1.2 各包规模

| 包名 | 文件数 | 总行数 | 最大文件 | 评价 |
|------|--------|--------|----------|------|
| `contracts` | 4 | 792 | constants.go (218) | 合理 |
| `config` | 1 | 81 | config.go (81) | 合理 |
| `agentconfig` | 2 | 652 | config.go (568) | 合理 |
| `ctrl` | 1 | 106 | gate.go (106) | 合理 |
| `obs` | 2 | 640 | logger.go (620) | 合理 |
| `dataset` | 6 | 1,647 | service.go (711) | 合理 |
| `orchestrator` | 1 | 247 | service.go (247) | 合理 |
| `store` | 5 | 3,137 | ingest.go (878) | 合理（已拆分） |
| `evaluator` | 16 | 6,561 | service.go (1,672) | 合理（Phase 2 重构后） |
| `runner` | 16 | 6,510 | service.go (1,225) | 偏大 |
| `reporter` | 9 | 4,890 | html_sections.go (1,549) | 合理（已拆分） |
| `web` | 10 | 4,480 | server.go (1,827) | 偏大 |

### 1.3 五阶段流水线架构

```
                    ┌─────────────┐
                    │  contracts  │  Layer 0: 数据结构、常量
                    └──────┬──────┘
                           │
          ┌────────────────┼────────────────┐
          │                │                │
    ┌─────┴─────┐   ┌─────┴─────┐   ┌─────┴─────┐
    │  dataset   │   │   ctrl    │   │    obs    │  Layer 1: 基础设施
    └─────┬─────┘   └───────────┘   └───────────┘
          │
    ┌─────┴─────┐
    │  runner    │◄─── 生成单元测试（LLM API + Agent Sandbox）
    └─────┬─────┘
          │
    ┌─────┴─────┐
    │ evaluator  │◄─── 评测（编译→测试→覆盖率→变异）
    └─────┬─────┘
          │
    ┌─────┴─────┐
    │ reporter   │◄─── 报告生成（JSON + HTML）
    └─────┬─────┘
          │
    ┌─────┴─────┐
    │orchestrator│◄─── 编排 generate → evaluate → report
    └─────┬─────┘
          │
    ┌─────┴─────┐
    │    web     │◄─── HTTP 管理 UI
    └─────┬─────┘
          │
    ┌─────┴─────┐
    │cmd/utbench │◄─── CLI 入口
    └───────────┘
```

---

## 二、包间依赖关系

### 2.1 依赖图

```
contracts          (0 内部依赖 — 最底层)
ctrl               (0 内部依赖 — 独立原语)
obs                (0 内部依赖 — 独立日志)
config             → contracts
agentconfig        → contracts
dataset            → contracts
store              → contracts
evaluator          → contracts, obs, store
runner             → contracts, agentconfig, ctrl, obs  ← 已解耦 store
reporter           → contracts, obs                   ← 已解耦 runner
orchestrator       → contracts, dataset, evaluator, reporter, runner, store
web                → contracts, agentconfig, ctrl, dataset, evaluator,
                     obs, orchestrator, reporter, runner, store
```

### 2.2 不合理的依赖方向

#### 问题 1: reporter → runner（违反分层原则）

- **位置**: `reporter/service.go` 第 71, 224, 230, 233 行
- **现状**: 报告生成层 import runner 包，调用 `runner.PromptTemplatePreview()`、`runner.PromptStrategy()`、`runner.PromptVersionID()`、`runner.LoadPromptCatalog()`
- **影响**: reporter 无法脱离 runner 编译和测试；runner 的任何改动都可能影响 reporter
- **根因**: prompt 元数据（strategy, version, catalog）定义在 runner 包中，但 reporter 也需要展示这些信息
- **修复**: 将 prompt 元数据提取到 `contracts` 包或独立的 `promptmeta` 包

#### 问题 2: runner → store（跨层依赖）

- **位置**: `runner/service.go` 第 91-93 行; `runner/reuse_plan.go` 第 19, 34 行
- **现状**: runner 直接 import store 包查询 SQLite 复用记录
- **影响**: runner 无法脱离 SQLite 进行单元测试
- **修复**: 定义 `ReuseStore` 接口，runner 依赖接口而非具体实现

#### 问题 3: web 扇出过大（上帝包）

- **位置**: `web/run_manager.go` import 9 个内部包
- **现状**: web 包同时依赖 contracts, ctrl, dataset, evaluator, obs, orchestrator, reporter, runner, store
- **影响**: 任何包的改动都可能影响 web；无法独立测试 web 层
- **修复**: web 层应只依赖 orchestrator 接口，不应直接依赖 evaluator/runner/store

---

## 三、SOLID 原则违反分析

### 3.1 S — 单一职责原则 (SRP)

#### 3.1.1 evaluator/service.go — 最严重的 SRP 违反

**2,168 行，57 个函数**，承担了以下全部职责：

| 职责 | 行数范围 | 函数数 | 说明 |
|------|---------|--------|------|
| Worker 池调度 | L95-L315 | 2 | `Evaluate` + 结果收集循环 |
| 单样本评测主流程 | L316-L890 | 1 | `evaluateOne` 574 行 |
| Python 工作区准备 | L1230-L1400 | 3 | `preparePythonWorkspace` 等 |
| Python 编译/测试/覆盖率 | L1400-L1614 | 6 | 应在 `python_eval.go` 中 |
| Python 导入重写 | L1616-L1690 | 4 | `rewriteGeneratedTestImports` 等 |
| 测试结果解析 | L1883-L1960 | 1 | `parsePytestCounts` |
| 断言密度估算 | L1695-L1880 | 5 | 4 语言各一个函数 + 分发函数 |
| 失败分类 | L1084-L1210 | 5 | `classifyFailureOrigin` 等 |
| 结果判定/完成 | L900-L1050 | 7 | `evaluationFailed` 等 |
| 结果哈希/复用 | L2100-L2170 | 5 | `evaluationKeyForItem` 等 |
| 工具函数 | L1960-L2100 | 15 | `round`, `trimErr`, `toFloat` 等 |

**对比**: Go/Java/C++ 各有独立的 `go_eval.go`、`java_eval.go`、`cpp_eval.go`，但 Python 的 ~400 行评测函数散落在 `service.go` 中，与其他语言不一致。

#### 3.1.2 runner/adapter_cli.go — 5+ 职责

**1,050 行**，职责列表：

| 职责 | 行数范围 | 说明 |
|------|---------|------|
| CLI Agent 执行编排 | L19-L384 | `generateCLIAgent` 366 行 |
| Agent 输出解析 | L386-L521 | `parseAgentOutput` 等 |
| Usage/Session 解析 | L621-L875 | token 用量递归 JSON 解析 |
| OpenCode Session 导出 | L668-L764 | 框架特定逻辑 |
| Token 成本估算 | L877-L929 | `estimateCostUSD` |
| 路径/命令过滤工具 | L572-L1050 | 20+ 辅助函数 |

#### 3.1.3 reporter/service.go — 3,764 行单文件

包含：HTML 报告生成、原始数据表构建（312行+198行）、多维分析（230行）、截断分析（207行）、洞察面板、场景分析、评分计算。

#### 3.1.4 巨型函数

| 函数 | 行数 | 圈复杂度 | 问题 |
|------|------|---------|------|
| `evaluator/service.go:evaluateOne` | 584 | ~35-40 | 4 语言 × 5 阶段全在一个函数 |
| `runner/service.go:generateOne` | 539 | ~25 | checkpoint + 复用 + prompt + API + 文件 + metadata |
| `store/sqlite.go:Init` | 423 | ~15 | schema DDL + 迁移 + 索引全在一个函数 |
| `runner/adapter_cli.go:generateCLIAgent` | 366 | ~20 | workspace + skill + 命令 + 执行 + 结果 |

### 3.2 O — 开放封闭原则 (OCP)

**这是项目最严重的系统性问题。** 几乎所有扩展点都需要修改现有代码，而非通过新增代码实现。

#### 3.2.1 语言扩展：需要改 6+ 处

新增一门语言（如 JavaScript）需要修改的位置：

| # | 文件 | 位置 | 修改内容 |
|---|------|------|---------|
| 1 | `contracts/constants.go` | L38 | `SupportedLanguages` 添加 |
| 2 | `runner/prompt.go` | L38 | `promptLanguages` 添加 |
| 3 | `runner/prompt.go` | L404 | `promptLanguageRules` 添加 case |
| 4 | `runner/prompt.go` | L231-288 | `buildFullFilePrompt` 可能需要特殊处理 |
| 5 | `runner/api.go` | L809 | `languageFramework` 添加 case |
| 6 | `runner/api.go` | L678 | `extractCode` 添加语言别名 |
| 7 | `evaluator/service.go` | L372-887 | `evaluateOne` 添加 100+ 行分支 |
| 8 | `evaluator/service.go` | L1695 | `estimateAssertionDensity` 添加 case |
| 9 | `evaluator/service.go` | L2080 | `getLanguagesSummary` 添加 |
| 10 | `reporter/service.go` | L217 | `loadPromptArtifacts` 添加 |
| 11 | `reporter/service.go` | L3298 | 提示词展示面板添加 |
| 12 | `reporter/service.go` | L3688 | `getScenarioLabel` 可能需要 |
| 13 | `reporter/static/charts.js` | L68 | 场景解析添加 |
| 14 | `reporter/static/charts.js` | L99 | 错误标签添加 |

**总计: 至少 14 处修改，遗漏任何一处都会导致静默 bug。**

#### 3.2.2 Provider 扩展：需要改 5 个函数

新增一个 LLM Provider（如 Anthropic）需要修改 `runner/api.go` 中：

1. `resolveEndpoint` (L290) — 端点路径拼接
2. `buildPayload` (L313) — 请求体格式
3. `buildContinuationPayload` (L351) — 续写请求体
4. `extractResponseText` (L507) — 响应文本提取
5. `extractFinishReason` (L576) — 结束原因提取

#### 3.2.3 场景扩展：需要改 4+ 处

新增一个场景类型（如 `algorithm`）需要修改：

1. `reporter/service.go:667` — `extractScenario` 前缀列表
2. `reporter/service.go:3688` — `getScenarioLabel` 标签映射
3. `reporter/static/charts.js:68` — JS 场景解析
4. `reporter/static/charts.js:78` — JS 标签映射

#### 3.2.4 评分权重：4 处独立维护

| 位置 | 内容 | 用途 |
|------|------|------|
| `reporter/service.go:1097` | `m.CompilePassRate*0.3 + ...` | 实际计算 |
| `reporter/service.go:1998` | `编译×0.3 + 样本测试×0.3 + ...` | HTML 描述 |
| `reporter/service_html_new.go:277` | 同上 | HTML 描述 |
| `reporter/service_insights.go:98` | 同上 | HTML 描述 |
| `store/sqlite.go:1789` | `{"compile":0.30,...}` | DB 存储（**未被使用**） |

### 3.3 L — 里氏替换原则 (LSP)

**基本遵守。** 两个接口设计合理：

- `SandboxRunner` (`sandbox.go:39`): `Run(ctx, SandboxRunRequest) (SandboxRunResult, error)`
- `AgentAdapter` (`adapter.go`): `Generate(ctx, AgentGenerateRequest) (AgentGenerateResult, error)`

**问题**: `modelConfig` 是结构体而非接口，所有消费者直接依赖字段，无法替换配置来源。

### 3.4 I — 接口隔离原则 (ISP)

`AgentAdapter` 接口只有一个方法，设计干净。

**缺失**: 没有 `LanguageEvaluator` 接口。每种语言的评测函数是裸函数（`goCompileCheck`、`javaCompileCheck` 等），调用方直接依赖具体实现。

### 3.5 D — 依赖倒置原则 (DIP)

| 违反 | 高层模块 | 低层模块 | 影响 |
|------|---------|---------|------|
| runner → store | `runner.Service` | `store.SQLiteStore` | 无法脱离 DB 测试 runner |
| reporter → runner | `reporter.Service` | `runner.PromptCatalog` | 无法脱离 runner 测试 reporter |
| evaluator 内部 | `evaluateOne` | `goCompileCheck()` 等 | 无法 mock 语言执行器 |
| web → 一切 | `web.Server` | 9 个内部包 | 上帝包 |

---

## 四、重复代码分析

### 4.1 evaluateOne 中的 4 重复制

`evaluator/service.go` 的 `evaluateOne` 函数中，python/go/java/cpp 四个分支包含以下几乎逐字重复的代码块：

#### 重复块 1: 测试结果赋值（出现 4 次，~8 行/次）

```go
row.TestPass = &pass
if !pass && testErr != "" {
    row.TestError = testErr
}
if runtimeMs > 0 {
    row.RuntimeMS = &runtimeMs
}
```

#### 重复块 2: 测试计数解析与填充（出现 4 次，~15 行/次）

```go
passCnt, totalCnt := parseXxxTestCounts(testErr)
if passCnt != nil { row.TestPassCount = passCnt }
if totalCnt != nil { row.TestTotalCount = totalCnt }
if passCnt != nil && totalCnt != nil && *totalCnt > 0 {
    rate := round(float64(*passCnt)/float64(*totalCnt), 6)
    row.TestPassRate = &rate
}
ensureTestCountsFromPass(&row)
```

#### 重复块 3: 变异测试结果赋值（出现 4 次，~15 行/次）

```go
row.MutationScore = &mutationScore
row.MutationTotal = &mutationStats.Total
row.MutationKilled = &mutationStats.Killed
row.MutationSurvived = &mutationStats.Survived
row.MutationNoTests = &mutationStats.NoTests
row.MutationTimeouts = &mutationStats.Timeout
row.MutationSkipped = &mutationStats.Skipped
row.MutationSuspicious = &mutationStats.Suspicious
if mutationErr != "" { row.MutationError = mutationErr }
```

#### 重复块 4: 变异测试前置检查（出现 4 次，~10 行/次）

```go
testPassed := 0
if row.TestPassCount != nil { testPassed = *row.TestPassCount }
testTotal := 0
if row.TestTotalCount != nil { testTotal = *row.TestTotalCount }
checkResult := CheckTestPassRate(testPassed, testTotal, toolName, GetMinPassRateForTool(toolName))
```

**总重复代码量: ~192 行（4 × 48 行）**

### 4.2 generateOne 中的复制粘贴

`runner/service.go` 的 `generateOne` 中，复用逻辑有两段几乎相同的 80+ 行代码（L547-641 和 L642-738），分别处理不同的复用来源。

### 4.3 各语言评测函数的公共签名模式

每种语言都实现了相同的 6 个函数，但没有统一接口：

| 阶段 | Python | Go | Java | C++ |
|------|--------|-----|------|-----|
| 准备工作区 | `preparePythonWorkspace` | `prepareGoWorkspace` | `prepareJavaWorkspace` | `prepareCppWorkspace` |
| 编译检查 | `pythonCompileCheck` | `goCompileCheck` | `javaCompileCheck` | `cppCompileCheck` |
| 执行测试 | `executePythonTests` | `executeGoTests` | `executeJavaTestsWithTimeout` | `executeCppTests` |
| 解析计数 | `parsePytestCounts` | `parseGoTestCounts` | `parseJavaTestCounts` | `parseCppTestCounts` |
| 收集覆盖率 | `collectPythonCoverage` | `collectGoCoverage` | `collectJavaCoverage` | `collectCppCoverage` |
| 收集变异 | `collectPythonMutation` | `collectGoMutation` | `collectJavaMutation` | `collectCppMutation` |

---

## 五、硬编码清单

### 5.1 散落的列表型硬编码

| 概念 | 散落数量 | 典型位置 |
|------|---------|---------|
| 语言列表 | 6+ 处 | `contracts/constants.go:38`, `runner/prompt.go:38`, `reporter/service.go:217`, `reporter/service.go:3298`, `reporter/charts.js:68`, `reporter/service.go:3688` |
| 场景列表 | 4+ 处 | `reporter/service.go:667`, `reporter/service.go:1912`, `reporter/charts.js:68`, `reporter/service.go:3688` |
| 评分权重 | 4 处 | `reporter/service.go:1097`, `service.go:1998`, `service_html_new.go:277`, `service_insights.go:98` |
| Provider 名称 | 6+ 处 | `runner/api.go` 中 6 处 `model.Provider == "xxx"` |
| 变异工具阈值 | 1 处 | `evaluator/mutation_types.go:126` switch-case |
| 错误类型标签 | 2 处 | `reporter/service.go` + `reporter/charts.js:99` (Go 和 JS 各维护一份) |

### 5.2 工具版本硬编码

| 工具 | 硬编码版本 | 位置 |
|------|-----------|------|
| Java compiler | 17 | `java_eval.go:38` |
| JUnit | 5.10.2 | `java_eval.go:45` |
| maven-compiler-plugin | 3.11.0 | `java_eval.go:61` |
| maven-surefire-plugin | 3.1.2 | `java_eval.go:70` |
| JaCoCo | 0.8.11 | `java_eval.go:82` |
| PITest | 1.19.6 | `java_eval.go:101` |
| C++ standard | C++17 | `cpp_eval.go:28` |
| clang | clang-18/19 | `cpp_eval.go:62`, `cpp_eval.go:564` |
| Go version | 1.24 | `go_eval.go:88` |
| GoogleTest | 系统安装 | `cpp_eval.go:30` |
| Chart.js | 4.4.1 | `reporter/charts.js:22` |

### 5.3 路径硬编码

| 路径 | 位置 | 风险 |
|------|------|------|
| `/workspace` | `sandbox.go:99`, `adapter_cli.go:57` | 容器内工作区固定 |
| `/app/artifacts` | `sandbox.go:129` | 容器内输出路径默认值 |
| `/opt/venv/bin/python` | `evaluator/service.go:2030` | Docker 内 Python 路径 |
| `/usr/lib/libgtest.a` | `cpp_eval.go:582` | gtest 静态库路径 |
| `/usr/lib/llvm-19/lib` | `cpp_eval.go:609` | LLVM 库路径 |
| `/var/run/docker.sock` | `sandbox.go:200` | Docker socket 检测 |
| `//bin/sh` | `sandbox.go:210` | Windows MSYS 绕过 |

### 5.4 超时/魔法数字

| 值 | 用途 | 位置 |
|------|------|------|
| 180s | 默认测试超时 | `evaluator/java_eval.go:20` |
| 120s | cmake/make/Java 测试超时 | 多处 |
| 300s | Mull 最大超时 | `cpp_eval.go:497` |
| 30s | gcov/mutmut export 超时 | 多处 |
| 5000ms | PITest timeoutConstant | `java_eval.go:225` |
| 200ms | 全局 API 最小间隔 | `runner/api.go:38` |
| 300s | HTTP 客户端超时 | `runner/api.go:46` |
| 3 | 重试次数 | `runner/api.go:47` |
| 0.3/0.3/0.2/0.2 | 评分权重 | `reporter/service.go:1097` |
| 1.0/0.8/0.5 | 变异测试最小通过率 | `evaluator/mutation_types.go:126` |

---

## 六、测试输出解析脆弱性

评测器的核心逻辑是"调外部工具 → 解析文本输出"，这层耦合极其脆弱：

### 6.1 各语言的解析依赖

| 语言 | 工具 | 解析方式 | 脆弱性 |
|------|------|---------|--------|
| Python | pytest | 正则 `(\d+)\s+passed` | pytest 8.x 输出格式变化 |
| Go | go test | 正则 `--- (PASS\|FAIL):` | go test -json 格式不同 |
| Java | Maven Surefire | 正则 `Tests run:\s*(\d+)` | 国际化可能导致不同语言 |
| C++ | GoogleTest | 正则 `\[(\d+)/(\d+)\] PASSED` | GTest 1.14+ 格式变化 |
| Python | coverage.py | JSON 字段名 `covered_lines` | 版本更新可能改字段名 |
| Java | JaCoCo | XML 解析 | 大版本间结构变化 |
| C++ | gcov | 冒号分隔文本 `exec:line:code` | gcov 12+ 格式变化 |
| Python | mutmut | JSON + .meta 文件 | mutmut 版本变化 |
| Java | PITest | XML/HTML/CSV | 多格式解析 |
| C++ | Mull | 10+ 正则模式 | 极其脆弱 |
| Go | go-mutesting | 正则 `mutation score is` | 项目不活跃 |

### 6.2 兜底逻辑问题

当正则解析失败时，多处使用 `passed=1, total=1` 兜底：

| 位置 | 条件 | 兜底值 |
|------|------|--------|
| `java_eval.go:415` | 包含 "BUILD SUCCESS" | 1/1 |
| `cpp_eval.go:301` | 包含 "[ PASSED ]" | 1/1 |
| `go_eval.go:155` | 包含 "PASS" 且不含 "FAIL" | 1/1 |

这会掩盖真实的测试问题：如果一个样本有 10 个测试但只有 1 个通过，解析失败后会报告 1/1（100%通过率）。

---

## 七、可测试性评估

### 7.1 容易测试的函数（纯逻辑，~40 个）

纯输入输出、无 IO 依赖的函数，如 `classifyFailureOrigin`、`parsePytestCounts`、`CheckTestPassRate`、`round`、`trimErr` 等。这些已有较好的测试覆盖。

### 7.2 难以测试的函数（~20 个）

| 函数 | 难点 |
|------|------|
| `evaluateOne` (584行) | 4 语言分支 + 外部工具进程，无法 mock |
| `generateOne` (539行) | 文件系统 + API 调用 + checkpoint |
| `Evaluate` / `Generate` | 并发 worker + 文件系统 + SQLite |
| 各语言 `compileCheck` / `executeTests` | 依赖外部命令（go/mvn/cmake/pytest） |
| 各语言 `collectMutation` | 依赖变异测试工具二进制 |

### 7.3 可测试性改进方向

1. **引入 `LanguageEvaluator` 接口** — 可以 mock 语言执行器来测试调度逻辑
2. **引入 `LLMProvider` 接口** — 可以 mock API 调用来测试续写/重试逻辑
3. **引入 `ReuseStore` 接口** — 可以 mock 存储层来测试复用逻辑
4. **拆分巨型函数** — 将 584 行的 `evaluateOne` 拆为 5 个可独立测试的子函数

---

## 八、CLI 入口分析

### 8.1 main.go 结构

`cmd/utbench/main.go` 1,610 行，使用手工 switch-case 路由：

```
utbench <command> [subcommand] [flags]
```

一级命令 10 个：`run`, `generate`, `evaluate`, `report`, `db`, `assets`, `dataset`, `doctor`, `web`, `help`

二级子命令：`db` 有 9 个, `assets` 有 4 个, `dataset` 有 4 个

### 8.2 问题

1. **没有使用 cobra/urfave-cli** — 手工 flag 解析导致大量样板代码
2. **业务逻辑下沉不足** — `doctor` 的金丝雀测试生成、`db` 的查询逻辑等本应在 internal 包中
3. **flag 定义散落** — 每个命令函数内部独立定义 flag，无法复用

---

## 九、重构方案

### 9.1 重构原则

1. **渐进式** — 每一步独立可验证，不破坏现有功能
2. **测试先行** — 重构前先补充关键路径的测试覆盖
3. **接口驱动** — 通过接口解耦，而非通过包重组
4. **单一来源** — 每个概念只在一个地方定义

### 9.2 分阶段计划

#### Phase 0: 准备（无代码改动）

- [x] 输出本分析报告
- [ ] 补充 `evaluateOne` 关键路径的集成测试（使用真实工具）
- [ ] 补充 `generateOne` 的单元测试（mock API）
- [ ] 确认 Docker 构建和 CI 流水线正常

#### Phase 1: 统一常量（低成本高收益，1-2 天） ✅ 已完成

**目标**: 消除散落的硬编码列表，建立单一来源。

**完成内容**:

1. 在 `contracts/constants.go` 中定义统一常量：
   - `SupportedLanguages` — 替代 6+ 处硬编码语言列表
   - `SupportedScenarios` — 替代 4+ 处硬编码场景列表
   - `ScenarioLabels` — 场景中英文标签映射
   - `ScoreWeights` 类型 + `DefaultWeights` — 替代 4 处硬编码权重
   - `ScoreWeights.String()` — 公式文本单一来源

2. 全局替换所有硬编码引用：
   - `runner/prompt.go` — promptLanguages → contracts.SupportedLanguages
   - `reporter/service.go` — 6 处替换（场景、权重、公式、语言）
   - `reporter/service_html_new.go` — 公式文本
   - `reporter/service_insights.go` — 公式文本
   - `reporter/embed.go` — BuildChartsJS 注入场景数据
   - `reporter/static/charts.js` — 硬编码场景数组 → 占位符注入
   - `evaluator/service.go` — 语言列表
   - `runner/service.go` — 语言列表
   - `web/server.go` — 语言和场景列表
   - `dataset/service.go` — 3 处场景引用
   - `dataset/service_test.go` — 语言列表
   - `store/sqlite.go` — 权重 JSON

3. 前端 `charts.js` 的场景列表通过 Go 模板注入，不再独立维护

**影响范围**: 纯重命名/引用替换，无逻辑变更
**风险**: 极低
**收益**: 加新语言/场景只需改一处

#### Phase 2: 引入 LanguageEvaluator 接口（核心 OCP 改造，3-5 天） ✅ 已完成

**目标**: 消除 `evaluateOne` 中 500 行 if/else 链。

**完成内容**:

1. 定义接口 (`evaluator/lang_eval.go`)：
   ```go
   type LanguageEvaluator interface {
       PrepareWorkspace(item contracts.GeneratedCase) (*WorkspaceContext, error)
       CompileCheck(ws *WorkspaceContext) (bool, string)
       ExecuteTests(ws *WorkspaceContext, timeoutSeconds int) (bool, string, int)
       ParseTestCounts(testOutput string) (*int, *int)
       EstimateAssertionDensity(workdir, testPath string) (int, int, float64)
       CollectCoverage(ws *WorkspaceContext, testPassed bool, timeoutSeconds int) (float64, float64, string)
       CollectMutation(ctx context.Context, ws *WorkspaceContext, input MutationInput) (float64, mutationStats, string)
       MutationTool() string
   }
   ```

2. 注册表机制 (`init()` + `RegisterLanguageEvaluator`)：
   - `go_eval_impl.go` — 包装 go_eval.go 函数
   - `java_eval_impl.go` — 包装 java_eval.go 函数
   - `cpp_eval_impl.go` — 包装 cpp_eval.go 函数
   - `python_eval_impl.go` — 处理 repo_level/self_contained 双模式

3. 提取共用逻辑为公共函数：
   - `computeTestPassRate(row, passCnt, totalCnt)` — 消除 4 处重复的测试计数处理
   - `populateMutationResult(row, score, stats, err, tool)` — 消除 4 处重复的变异结果写入
   - `populateCoverageResult(row, lineCov, branchCov, covErr, testPassed)` — 消除 4 处重复的覆盖率写入

4. `evaluateOne` 从 575 行缩减到 ~75 行：
   - 原 515 行 if/else-if 链 → 3 行注册表查询 + 分发
   - `service.go` 从 2169 行降到 1672 行（-497 行）

**影响范围**: evaluator 包内部重构
**风险**: 中等（已通过全部测试验证）
**收益**: 加新语言只需实现接口 + 1 行 init() 注册；evaluateOne 可测试

#### Phase 3: 引入 LLMProvider 接口（OCP 改造，2-3 天） ✅ 已完成

**目标**: 消除 api.go 中 6 处 provider if/switch。

**完成内容**:

1. 定义接口 (`runner/provider.go`)：
   ```go
   type LLMProvider interface {
       ResolveEndpoint(model modelConfig) string
       BuildPayload(model modelConfig, prompt string) map[string]any
       BuildContinuationPayload(model modelConfig, originalPrompt, generatedSoFar string) map[string]any
       ExtractResponseText(response map[string]any) (string, error)
       ExtractFinishReason(response map[string]any) bool
   }
   ```

2. 注册表机制 (`RegisterProvider` + `resolveProvider`)：
   - `deepseek`、`minimax`、`volcengine` → `OpenAICompatibleProvider`
   - `dashscope` (非 compatible-mode) → `DashscopeNativeProvider`
   - 未知 provider → 默认 `OpenAICompatibleProvider`

3. 两种实现 (`runner/provider_impl.go`)：
   - `OpenAICompatibleProvider` — 标准 OpenAI Chat Completions 格式
   - `DashscopeNativeProvider` — Dashscope 专有 input/parameters 格式

4. 重构 `api.go`：
   - `generateTest` 使用 `resolveProvider(model)` 获取 provider 实例
   - `doOnce` 接收 `LLMProvider` 替代 `provider string`
   - 删除 5 个旧函数：`resolveEndpoint`、`buildPayload`、`buildContinuationPayload`、`extractResponseText`、`extractFinishReason`
   - `api.go` 从 1197 行降到 992 行（-205 行）
   - 消除 volcengine 冗余分支（与 default 完全相同）

**影响范围**: runner 包内部重构
**风险**: 低（已通过全部测试验证）
**收益**: 加新 provider 只需实现接口 + 注册；消除所有 provider if/switch

#### Phase 4: 拆分巨型文件 ✅ 已完成

**目标**: 将超过 2000 行的巨型文件拆分为职责单一的小文件。

1. `store/sqlite.go` (3092行) → 拆分为 4 个文件：
   - `store/sqlite.go` (842行) — 类型定义、连接、Schema DDL
   - `store/ingest.go` (878行) — 写入路径 ingest 逻辑
   - `store/query.go` (607行) — 读取路径查询/报告
   - `store/list.go` (810行) — 18 个 List* 方法 + 22 个工具函数

2. `reporter/service.go` (3764行) → 拆分为 5 个文件：
   - `reporter/service.go` (404行) — Service/Output/ModelDetail 类型 + 主编排
   - `reporter/aggregation.go` (602行) — 多维聚合类型 + merge 函数
   - `reporter/ranking.go` (1010行) — 评分、排名、比较、失败分析、工具函数
   - `reporter/html_sections.go` (1549行) — buildHTML + 各 HTML section 构建
   - `reporter/html_prompt_charts.go` (235行) — prompt HTML + Chart.js 脚本

**影响范围**: 文件级重组，不改接口
**风险**: 低（纯移动代码）
**收益**: 可读性、可维护性大幅提升
**验证**: `go build ./internal/... ./cmd/...` + `go test ./internal/...` 全部通过

> 注: evaluator/service.go 的巨型函数问题已在 Phase 2 通过 LanguageEvaluator 接口解决（evaluateOne 从 575 行降至 ~75 行），Phase 4 无需再拆。runner/adapter_cli.go 暂未拆分（优先级较低，可后续处理）。

#### Phase 5: 修复依赖方向 ✅ 已完成

1. **reporter → runner 解耦** ✅：
   - 定义 `contracts.PromptMetaProvider` 接口（PromptStrategy/PromptVersionID/PromptTemplatePreview）
   - `PromptCatalog`、`PromptMode` 类型和 `LoadPromptCatalog` 移到 `contracts` 包
   - runner 提供 `DefaultPromptMetaProvider` 实现
   - reporter 通过接口获取 prompt 元数据，不再 import runner

2. **runner → store 解耦** ✅：
   - 定义 `runner.GenerationReuseStore` 接口（FindReusableGeneratedAsset）
   - `ReusableGeneratedCase` 类型移到 `contracts` 包，store 使用类型别名
   - runner 依赖接口，store 实现接口
   - store 打开逻辑移到 orchestrator 和 cmd 层

**影响范围**: 接口级变更
**风险**: 低
**收益**: reporter 和 runner 包可独立编译和测试
**验证**: `go build ./internal/... ./cmd/...` + `go test ./internal/...` 全部通过

> 注: web 包瘦身（Docker runner 拆分）优先级较低，暂未处理。

#### Phase 6: 改进测试输出解析健壮性 ✅ 已完成

1. **消除 1/1 兜底** ✅：Go parser 不再对 bare "PASS" 输出伪造 passed=1/total=1，返回 nil
2. **Mutation 总数膨胀修复** ✅：去掉 mutation collectors 中 `passed=1` 下限，避免 0% 通过率被抬高
3. **Python 测试验证放宽** ✅：支持 `unittest.TestCase` 风格测试文件（service.go + api.go）
4. **统一超时管理** ✅：新建 `evaluator/timeout.go`，集中定义 `DefaultTestTimeoutSeconds`(180)、`MutationTimeoutSeconds`(120)、`CppMutationTimeoutSeconds`(300)、`PythonCompileTimeoutSeconds`(30)

**影响范围**: 评测器内部逻辑
**风险**: 低
**收益**: 评测结果准确性提升，超时配置可维护性改善
**验证**: `go build ./internal/... ./cmd/...` + `go test ./internal/...` 全部通过

### 9.3 时间估算

| Phase | 内容 | 预估工时 | 风险 | 依赖 | 状态 |
|-------|------|---------|------|------|------|
| 0 | 测试补充 | 2-3 天 | 低 | 无 | 待定 |
| 1 | 统一常量 | 1-2 天 | 极低 | 无 | ✅ 已完成 |
| 2 | LanguageEvaluator 接口 | 3-5 天 | 中 | Phase 1 | ✅ 已完成 |
| 3 | LLMProvider 接口 | 2-3 天 | 低 | 无 | ✅ 已完成 |
| 4 | 拆分巨型文件 | 2-3 天 | 低 | Phase 2, 3 | ✅ 已完成 |
| 5 | 修复依赖方向 | 2-3 天 | 中 | Phase 2, 3 | ✅ 已完成 |
| 6 | 解析健壮性 | 1-2 天 | 低 | Phase 2 | ✅ 已完成 |
| **合计** | | **13-21 天** | | |

### 9.4 重构前后对比（预期）

| 指标 | 重构前 | 重构后 |
|------|--------|--------|
| 加新语言需改文件数 | 14+ | 2（实现接口 + 注册） |
| 加新 Provider 需改文件数 | 5 | 1（实现接口 + 注册） |
| 最长函数 | 584 行 | ~100 行 |
| 最大文件 | 3,764 行 | 1,549 行（html_sections.go） |
| `evaluateOne` 圈复杂度 | 35-40 | ~8 |
| 重复代码 | ~192 行 | 0 |
| reporter → runner 依赖 | 有 | 无（通过 PromptMetaProvider 接口） |
| runner → store 依赖 | 有 | 无（通过 GenerationReuseStore 接口） |

---

## 十、风险与缓解

| 风险 | 概率 | 影响 | 缓解措施 |
|------|------|------|---------|
| 重构引入回归 bug | 中 | 高 | Phase 0 先补充测试覆盖；每步重构后跑完整 pipeline |
| 接口设计不合理 | 中 | 中 | 先实现最小接口，后续迭代扩展 |
| 重构范围膨胀 | 高 | 中 | 严格按 Phase 执行，不在此轮做功能变更 |
| 与功能开发冲突 | 中 | 低 | 重构分支定期从 master rebase |

---

## 附录 A: 关键文件索引

| 文件 | 行数 | 核心职责 | 重构优先级 |
|------|------|---------|-----------|
| `evaluator/service.go` | 2168 | 评测调度 + 4 语言评测逻辑 | Phase 2, 4 |
| `reporter/service.go` | 3764 | 报告生成 + 数据分析 | Phase 4 |
| `store/sqlite.go` | 3092 | SQLite 全部操作 | Phase 4 |
| `runner/service.go` | 1225 | 生成调度 + worker 池 | Phase 4 |
| `runner/api.go` | 1196 | LLM API 调用 + provider 分发 | Phase 3 |
| `runner/adapter_cli.go` | 1050 | CLI Agent 适配器 | Phase 4 |
| `evaluator/cpp_eval.go` | 1294 | C++ 评测 | Phase 2 |
| `web/server.go` | 1827 | HTTP 服务器 | Phase 5 |
| `cmd/utbench/main.go` | 1610 | CLI 入口 | 低优先级 |

## 附录 B: 各语言评测函数对照表

| 阶段 | Python | Go | Java | C++ |
|------|--------|-----|------|-----|
| 准备 | `preparePythonWorkspace` (service.go:1230) | `prepareGoWorkspace` (go_eval.go:62) | `prepareJavaWorkspace` (java_eval.go:150) | `prepareCppWorkspace` (cpp_eval.go:103) |
| 编译 | `pythonCompileCheck` (service.go:1500) | `goCompileCheck` (go_eval.go:29) | `javaCompileCheck` (java_eval.go:318) | `cppCompileCheck` (cpp_eval.go:211) |
| 测试 | `executePythonTests` (service.go:1529) | `executeGoTests` (go_eval.go:114) | `executeJavaTestsWithTimeout` (java_eval.go:348) | `executeCppTests` (cpp_eval.go:247) |
| 计数 | `parsePytestCounts` (service.go:1883) | `parseGoTestCounts` (go_eval.go:139) | `parseJavaTestCounts` (java_eval.go:379) | `parseCppTestCounts` (cpp_eval.go:264) |
| 覆盖 | `collectPythonCoverage` (service.go:1545) | `collectGoCoverage` (go_eval.go:175) | `collectJavaCoverage` (java_eval.go:445) | `collectCppCoverage` (cpp_eval.go:321) |
| 变异 | `collectPythonMutation` (mutation.go:58) | `collectGoMutation` (go_eval.go:309) | `collectJavaMutation` (java_eval.go:542) | `collectCppMutation` (cpp_eval.go:495) |

## 附录 C: 包间依赖矩阵

| 被依赖 → | contracts | ctrl | obs | config | agentconfig | store | dataset | runner | evaluator | reporter | orchestrator | web |
|----------|-----------|------|-----|--------|-------------|-------|---------|--------|-----------|----------|-------------|-----|
| contracts | - | | | | | | | | | | | |
| ctrl | | - | | | | | | | | | | |
| obs | | | - | | | | | | | | | |
| config | ✓ | | | - | | | | | | | | |
| agentconfig | ✓ | | | | - | | | | | | | |
| store | ✓ | | | | | - | | | | | | |
| dataset | ✓ | | | | | | - | | | | | |
| runner | ✓ | ✓ | ✓ | | ✓ | ✓ | | - | | | | |
| evaluator | ✓ | | ✓ | | | ✓ | | | - | | | |
| reporter | ✓ | | ✓ | | | | | ✓ | | - | | |
| orchestrator | ✓ | | | | | ✓ | ✓ | ✓ | ✓ | ✓ | - | |
| web | ✓ | ✓ | ✓ | | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ | - |
