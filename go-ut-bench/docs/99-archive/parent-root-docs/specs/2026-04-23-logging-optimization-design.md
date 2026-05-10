# 日志系统优化设计方案

**日期**: 2026-04-23
**状态**: 已确认，待实施
**作者**: OpenCode Agent

---

## 1. 背景与问题

当前 `go-ut-bench` 的日志系统存在以下问题：

1. **输出混杂**: 大量使用 `fmt.Fprintf(os.Stderr, ...)` 直接输出，`obs.Logger` 几乎未被使用
2. **进度信息匮乏**: 只有 `[d/total] model | lang | sample | status` 的简单输出，用户不知道当前在哪个阶段
3. **错误信息被截断**: 错误消息被截断到30字符，难以定位问题根因
4. **没有分层日志**: Debug/Info/Warn/Error 没有明确分工，verbose 模式也没有增加有用的调试信息
5. **长时间任务无进度**: 变异测试等耗时操作只有开始/结束标记，中间完全黑盒

---

## 2. 设计目标

1. **清晰的分层日志**: Trace/Debug/Info/Warn/Error 各司其职
2. **中文面向用户**: 终端输出使用中文，文件日志使用结构化英文
3. **实时进度可见**: 每5个样本刷新统计面板，显示阶段计时和预计剩余时间
4. **完整错误保留**: 文件日志保留完整错误详情，终端只显示状态
5. **多目标输出**: 终端（简洁）+ 文件（详细JSON）同时输出

---

## 3. 架构设计

### 3.1 整体架构

终端输出（中文、简洁、实时）：
- 阶段标题（如"阶段 1/3: 生成测试"）
- 进度行（含状态、耗时、Token等信息）
- 统计面板（每5个样本刷新）
- 状态提示（成功/失败/截断，不展开错误详情）

文件日志（JSON Lines，详细、可分析）：
- artifacts/runs/<run_id>/logs/
  - run.log: 整体流程日志
  - runner.log: 测试生成阶段（含API请求详情）
  - evaluator.log: 评测阶段（含编译/测试/覆盖率/变异测试详情）
  - errors.log: 错误汇总（便于快速排查问题）

设计原则:
- 终端面向用户: 中文、简洁、有进度感
- 文件面向调试: 完整、结构化、可检索
- 不修改业务逻辑: 只在关键节点插入日志，不影响现有执行流程

### 3.2 核心组件

#### 3.2.1 增强 Logger (internal/obs/logger.go)

新增 Trace 级别（比 Debug 更细）:
- Trace: API请求/响应详情、HTTP headers、Token消耗
- Debug: 内部状态变化、临时文件路径、命令执行详情
- Info: 阶段开始/结束、关键进度、统计摘要
- Warn: 非致命错误（如覆盖率收集失败但测试通过）、降级处理
- Error: 致命错误、编译失败、测试失败

关键功能:
- 同时输出到终端（中文、简洁）和文件（JSON、完整）
- 支持上下文字段自动附加（如 run_id, model, sample）
- 分类日志（写入不同文件：runner, evaluator, api, error）

接口设计:
```go
type Logger struct {
    terminalHandler slog.Handler              // 终端输出：中文、简洁
    fileHandlers    map[string]slog.Handler   // 文件输出：JSON、详细
    context         map[string]any            // 固定上下文（如 run_id, model）
    level           Level
}

func (l *Logger) Trace(msg string, attrs ...any)
func (l *Logger) Debug(msg string, attrs ...any)
func (l *Logger) Info(msg string, attrs ...any)
func (l *Logger) Warn(msg string, attrs ...any)
func (l *Logger) Error(msg string, attrs ...any)
func (l *Logger) With(attrs ...any) *Logger
func (l *Logger) WithContext(ctx context.Context) *Logger
func (l *Logger) ToFile(category string) *Logger  // "runner", "evaluator", "api", "error"
```

#### 3.2.2 进度报告器 (internal/obs/progress.go)

职责: 在终端显示实时进度和统计面板，与 Logger 解耦。

```go
type ProgressReporter struct {
    total          int
    completed      int
    successCount   int
    failCount      int
    truncatedCount int
    coverageSum    float64
    coverageCount  int
    mutationSum    float64
    mutationCount  int
    startTime      time.Time
    stage          string
    mu             sync.Mutex
}

func (pr *ProgressReporter) OnTaskStart(model, lang, sampleID string)
func (pr *ProgressReporter) OnTaskDone(result TaskResult)
func (pr *ProgressReporter) PrintStats()
func (pr *ProgressReporter) PrintStageStart(name string)
func (pr *ProgressReporter) PrintStageDone(name string, stats StageStats)
```

统计面板频率: 每 **5** 个样本刷新一次。

---

## 4. 终端输出规范

### 4.1 阶段 1: 生成测试

```
========================================
阶段 1/3: 生成测试
========================================
模型: deepseek, qwen | 语言: python, go | 样本数: 50 | 工作线程: 8

[1/50]  deepseek | python | boundary_01    | 成功 | 1,234 tokens | 1.2s
[2/50]  deepseek | python | boundary_02    | 成功 | 892 tokens  | 0.9s
[3/50]  qwen     | python | simple_01      | 成功 | 2,156 tokens | 2.1s
[4/50]  qwen     | go     | simple_02      | 截断 | 4,096/4,096 tokens (已达上限) | 0.8s
[5/50]  deepseek | go     | boundary_03    | 失败 | -- tokens   | 0.5s

========================================
实时统计
   已处理: 5/50 (10%)
   成功: 4 | 失败: 1 | 截断: 1
   阶段已耗时: 0:00:45
   预计剩余: ~2分30秒
========================================
```

字段说明:
- "1,234 tokens": 本次请求消耗的 Token 数
- "4,096/4,096 tokens (已达上限)": 截断时显示已达上限
- "阶段已耗时": 从开始计时，格式 H:MM:SS
- "预计剩余": 基于当前速度估算

### 4.2 阶段 2: 评测测试

评测流程：编译 -> 测试 -> 行覆盖 -> 变异

```
========================================
阶段 2/3: 评测测试
========================================
语言: python(25), go(25) | 变异测试: 已启用 | 工作线程: 8

[1/50]  deepseek | python | boundary_01    
        编译通过 -> 测试通过(3/3) -> 行覆盖: 85% -> 变异: mutmut | 20/20 | 存活: 5 | 杀死: 15 | 分数: 75% | 5.2s
        
[2/50]  deepseek | python | boundary_02    
        编译通过 -> 测试通过(5/5) -> 行覆盖: 92% -> 变异: mutmut | 15/15 | 存活: 3 | 杀死: 12 | 分数: 80% | 4.8s
        
[3/50]  qwen     | python | simple_01      
        编译通过 -> 测试失败(2/4) -> 行覆盖: 60% -> 跳过变异(测试未通过) | 2.1s
        
[4/50]  qwen     | go     | simple_02      
        编译失败 -> 跳过(编译未通过) | 0.8s
        
[5/50]  deepseek | go     | boundary_03    
        编译通过 -> 测试通过(3/3) -> 行覆盖: 45% -> 变异: gremlins | 10/10 | 存活: 8 | 杀死: 2 | 分数: 20% | 180.0s
        
[6/50]  deepseek | python | complex_01     
        编译通过 -> 测试通过(2/2) -> 行覆盖: 78% -> 变异: mutmut | 超时(30s) | 30.0s

========================================
实时统计
   已处理: 6/50 (12%)
   编译通过: 5/6 (83%)
   测试通过: 4/6 (67%)
   平均行覆盖: 72.0%
   平均变异分数: 58.3%
   阶段已耗时: 0:03:20
   预计剩余: ~8分15秒
========================================
```

变异测试显示规范:
- "变异: mutmut | 20/20 | 存活: 5 | 杀死: 15 | 分数: 75%"
  - "mutmut" = 使用的工具名（mutmut/gremlins/pitest/mull）
  - "20/20" = 生成20个变异体，全部执行
  - "存活: 5" = 5个变异体存活
  - "杀死: 15" = 15个变异体被杀死
  - "分数: 75%" = 变异分数

失败原因标记:
- "跳过变异(测试未通过)" = 测试没通过，跳过变异
- "跳过(编译未通过)" = 编译失败，跳过后续
- "变异: mutmut | 超时(30s)" = 变异测试超时
- "变异: mutmut | 失败: 配置错误" = 变异测试执行失败

### 4.3 阶段 3: 生成报告

```
========================================
阶段 3/3: 生成报告
========================================

报告生成完成 | 0:00:02

========================================
全部完成！
   运行ID: 20240115-103015
   总耗时: 12分30秒
   
   生成阶段: 50/50 (100%) | 耗时: 3分15秒
   评测阶段: 45/50 (90%) | 耗时: 8分45秒
   报告阶段: 1/1 | 耗时: 2秒
   
   编译通过率: 90% (45/50)
   测试通过率: 80% (40/50)
   平均行覆盖: 78.5%
   平均变异分数: 65.2%
========================================
```

---

## 5. 文件日志规范

### 5.1 输出路径

```
artifacts/runs/<run_id>/logs/
├── run.log           # 整体流程日志（Info 级别以上）
├── runner.log        # 测试生成阶段（Debug 级别以上）
├── evaluator.log     # 评测阶段（Debug 级别以上）
├── api_calls.log     # API 请求详情（Trace 级别）
└── errors.log        # 错误汇总（Error 级别）
```

### 5.2 格式规范

文件日志使用 JSON Lines 格式，便于后续分析：

```json
{
  "time": "2026-04-23T10:30:18Z",
  "level": "ERROR",
  "msg": "测试生成失败",
  "run_id": "20260423-103015",
  "model": "qwen",
  "language": "go",
  "sample_id": "simple_03",
  "stage": "generation",
  "error": {
    "kind": "compile_error",
    "message": "undefined: fmt.Println",
    "full_output": "...",
    "exit_code": 1
  },
  "latency_ms": 800,
  "request_id": "req_abc123"
}
```

### 5.3 日志分类规则

| 日志内容 | 目标文件 | 级别 |
|---------|---------|------|
| 阶段开始/结束 | run.log | Info |
| API 请求参数 | api_calls.log | Trace |
| API 响应内容 | api_calls.log | Trace |
| 编译命令执行 | evaluator.log | Debug |
| 测试执行输出 | evaluator.log | Debug |
| 覆盖率收集详情 | evaluator.log | Debug |
| 变异测试详情 | evaluator.log | Debug |
| 错误信息 | errors.log | Error |
| 统计摘要 | run.log | Info |

---

## 6. 错误处理规范

### 6.1 终端错误显示

原则: 终端只显示状态，不展开错误详情。

```
[14/50] qwen | go | simple_03 | 失败 | 编译错误 | 0.8s
```

用户可通过以下方式查看详情：
1. 查看文件日志: artifacts/runs/<run_id>/logs/errors.log
2. 查看响应文件: artifacts/runs/<run_id>/generated/metadata/...response.json

### 6.2 文件错误记录

原则: 文件日志保留完整错误信息，包括：
- 错误类型（kind）
- 错误消息（message）
- 完整输出（full_output，不截断）
- 退出码（exit_code）
- 上下文信息（run_id, model, sample 等）

---

## 7. 实施计划

### Phase 1: 基础设施改造
1. 改造 internal/obs/logger.go，增加 Trace 级别和文件输出能力
2. 新增 internal/obs/progress.go，实现 ProgressReporter
3. 新增 internal/obs/file_handler.go，实现分类文件日志

### Phase 2: Runner 阶段日志
1. 在 internal/runner/service.go 中插入阶段日志
2. 替换 fmt.Fprintf(os.Stderr, ...) 为 Logger + ProgressReporter
3. 增加 API 请求/响应的 Trace 日志

### Phase 3: Evaluator 阶段日志
1. 在 internal/evaluator/service.go 中插入阶段日志
2. 在每个评测步骤（编译、测试、覆盖率、变异）记录进度
3. 增加变异测试详细统计日志

### Phase 4: Orchestrator 阶段日志
1. 在 internal/orchestrator/service.go 中记录整体流程
2. 增加阶段切换时的汇总信息

### Phase 5: 测试与验证
1. 使用 --dry-run 验证日志输出格式
2. 使用小样本验证完整流程
3. 检查文件日志的 JSON 格式正确性

---

## 8. 兼容性说明

1. 命令行参数: 增加 --no-progress 参数，禁用终端进度条（纯日志模式）
2. 日志级别: --verbose 控制终端日志级别，不影响文件日志（文件始终记录 Trace 级别）
3. 向后兼容: 不改变现有 CLI 接口，只增加新的日志输出

---

## 9. 风险评估

| 风险 | 影响 | 缓解措施 |
|-----|------|---------|
| 日志写入影响性能 | 中 | 使用异步写入，文件日志不阻塞主流程 |
| 日志文件过大 | 低 | 按 run_id 分目录，不自动清理 |
| 终端输出格式错乱 | 低 | 使用互斥锁保证输出顺序 |
| 并发写入文件冲突 | 低 | 每个 worker 使用独立 Logger 实例 |

---

## 10. 验收标准

1. 终端显示中文进度，每5个样本刷新统计面板
2. 阶段计时器准确显示已耗时和预计剩余时间
3. 评测流程清晰展示：编译 -> 测试 -> 行覆盖 -> 变异
4. 变异测试显示：工具名、生成数、存活数、杀死数、分数
5. 终端不展开错误详情，文件日志保留完整错误信息
6. 文件日志按分类写入不同文件，格式为 JSON Lines
7. 支持 --no-progress 参数禁用终端进度条
8. 终端使用简洁状态（成功/失败/截断），文件保留完整英文日志
