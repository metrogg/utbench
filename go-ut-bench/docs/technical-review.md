# go-ut-bench 技术审查报告

审查日期: 2026-04-19
审查范围: `go-ut-bench/internal/` 所有模块

---

## 一、严重问题 (P0 - 必须立即修复)

### 1. 资源泄漏：临时目录未清理

**位置**: `evaluator/service.go:199`

**严重程度**: 高

**问题描述**:
Python workspace 使用 `defer cleanupWorkspace(workdir)` 清理，但只在非 module_level 时执行。module_level 样本的 workspace 不清理，可能导致磁盘空间累积。

**代码片段**:
```go
if !isModuleLevel {
    defer cleanupWorkspace(workdir)
}
```

**改进建议**:
- 所有路径都应在评测完成后清理临时文件
- 或显式记录 module_level workspace 生命周期管理策略
- 建议添加 `--keep-workspace` 参数用于调试场景

---

### 2. 潜在的 panic 未恢复

**位置**: `evaluator/service.go:74-77`

**严重程度**: 高

**问题描述**:
worker goroutine 使用 recover 捕获 panic，但仅打印日志后继续，panic 内容丢失，无法追踪根本原因。

**代码片段**:
```go
defer func() {
    if r := recover(); r != nil {
        fmt.Fprintf(os.Stderr, "[worker panic] %v\n", r)
    }
}()
```

**改进建议**:
- 记录完整 panic 信息到结果结构
- 或返回包含 panic 详情的结果
- 建议格式: `EvaluationResult.CompileError = "panic: " + r + " [stack trace]"`

---

### 3. 并发竞态：checkpoint 写入

**位置**: `runner/service.go:136-141`

**严重程度**: 高

**问题描述**:
多 worker 并发写入 checkpoint 文件，虽有 mutex 保护 map，但 `saveCheckpoint` 文件写入无锁保护，可能导致文件内容交错。

**代码片段**:
```go
ckptMu.Lock()
completed[key] = struct{}{}
_ = saveCheckpoint(checkpointPath, completed)
ckptMu.Unlock()
```

**改进建议**:
- 使用文件锁（`os.OpenFile` + `syscall.Flock`）
- 或改用单 goroutine 负责所有 checkpoint 写入（channel 传递写入请求）
- 或每次写入前先读取现有内容合并后再写入

---

### 4. exec.CommandContext 超时后的进程残留

**位置**: 
- `evaluator/cpp_eval.go:350-359`
- `evaluator/java_eval.go:339-344`
- `evaluator/go_eval.go:215-220`
- `evaluator/mutation.go:56-61`

**严重程度**: 高

**问题描述**:
`exec.CommandContext` 超时后，Go 只取消等待，**不会杀死子进程**。mutmut/maven/go-mutesting/mull 超时后可能继续运行，消耗资源。

**代码片段**:
```go
runCtx, cancelRun := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
defer cancelRun()
cmd := exec.CommandContext(runCtx, mutestingPath, "./...")
runOut, runErr := cmd.CombinedOutput()
```

**改进建议**:
- 超时后显式调用 `cmd.Process.Kill()`
- 或使用以下模式:

```go
done := make(chan error, 1)
go func() {
    done <- cmd.Wait()
}()
select {
case <-ctx.Done():
    cmd.Process.Kill()
    return ctx.Err()
case err := <-done:
    return err
}
```

---

## 二、重要问题 (P1 - 近期修复)

### 5. HTTP Client 超时过长

**位置**: `runner/api.go:30`

**严重程度**: 中

**问题描述**:
HTTP client 超时设为 120 秒，与某些模型 API 的现实响应时间不匹配。若模型返回慢，会阻塞 worker。

**代码片段**:
```go
client: &http.Client{Timeout: 120 * time.Second}
```

**改进建议**:
- 支持配置化超时参数（从 models.yaml 的 `benchmark.timeouts.api_call` 读取）
- 或分层超时（连接超时 10s vs 总超时 120s）

---

### 6. 随机数生成器未 Seed

**位置**: `runner/api.go:80`

**严重程度**: 中

**问题描述**:
使用 `rand.Int63n` 做 jitter 但未调用 `rand.Seed`，每次程序启动 jitter 序列相同。

**代码片段**:
```go
sleep += float64(time.Duration(rand.Int63n(int64(200 * time.Millisecond))))
```

**改进建议**:
- 使用 `rand.New(rand.NewSource(time.Now().UnixNano()))` 创建独立随机源
- 或 Go 1.20+ 直接使用 `rand.Int64()` (自动 seed)

---

### 7. 覆盖率解析逻辑过于简化

**位置**: `evaluator/go_eval.go:183-198`

**严重程度**: 中

**问题描述**:
`estimateGoStmtCount` 用逗号/点数量估算语句数，不准确。对复杂 range 格式可能产生错误结果。

**代码片段**:
```go
func estimateGoStmtCount(rangeStr string) int {
    count := 1
    for _, ch := range rangeStr {
        if ch == ',' || ch == '.' {
            count++
        }
    }
    if count > 1 {
        count = count / 2
    }
    ...
}
```

**改进建议**:
- 使用正规 Go coverage 解析方法（按行数而非 range 拆分）
- 或至少解析 range 格式提取真实行范围

---

### 8. 未处理 io.ReadAll 错误

**位置**: `runner/api.go:116`

**严重程度**: 中

**问题描述**:
`rawBytes, _ := io.ReadAll(resp.Body)` 忽略读取错误，可能导致空响应被误判为成功。

**代码片段**:
```go
rawBytes, _ := io.ReadAll(resp.Body)
rawText := string(rawBytes)
```

**改进建议**:
- 检查 ReadAll 错误并返回 `read_error` 类型

```go
rawBytes, err := io.ReadAll(resp.Body)
if err != nil {
    return "", nil, nil, nil, nil, &contracts.ErrorInfo{
        Kind: "read_error", 
        Message: err.Error(), 
        Retryable: true
    }
}
```

---

### 9. 测试计数解析覆盖不全

**位置**: `evaluator/service.go:959-1010`

**严重程度**: 中

**问题描述**:
`parsePytestCounts` 依赖特定格式，对 pytest 新版输出（如 `=== short test summary info ===`）可能解析失败。

**改进建议**:
- 增加更多 pytest 输出格式匹配模式
- 或使用 pytest-json-report 插件获取结构化输出
- 或使用 `pytest --report-json` 参数

---

### 10. SQLite Schema 升级不安全

**位置**: `store/sqlite.go:89-103`

**严重程度**: 中

**问题描述**:
ALTER TABLE 用 `_, _ = s.db.ExecContext(ctx, stmt)` 忽略所有错误，若列已存在可能返回错误但被忽略，不保证 schema 状态一致。

**代码片段**:
```go
alterDDLs := []string{
    `ALTER TABLE sample_results ADD COLUMN test_pass_count INTEGER;`,
    ...
}
for _, stmt := range alterDDLs {
    _, _ = s.db.ExecContext(ctx, stmt)
}
```

**改进建议**:
- 检查 sqlite_master 确认列是否存在后再 ALTER
- 或使用版本化 migration（如 `migrations/0002_add_columns.sql`）

---

### 11. 缺少 Context 传播验证

**位置**: `evaluator/service.go`

**严重程度**: 中

**问题描述**:
`evaluateOne` 内部创建了新的 context.WithTimeout，但未检查父 context 是否已取消，可能导致取消信号丢失。

**改进建议**:
- 在创建子 context 前检查 `ctx.Err()`

```go
if ctx.Err() != nil {
    return contracts.EvaluationResult{
        CompilePass: false,
        CompileError: "parent context cancelled",
    }
}
runCtx, cancel := context.WithTimeout(ctx, timeout)
```

---

## 三、安全问题 (P1)

### 12. API Key 在日志中可能泄露

**位置**: `runner/api.go:43-50`

**严重程度**: 中

**问题描述**:
错误消息包含 `model.APIKeyEnv`，若日志被公开分享可能透露使用哪些 API key。

**代码片段**:
```go
Message: fmt.Sprintf("missing API key env var: %s (model=%s)", model.APIKeyEnv, model.Name)
```

**改进建议**:
- 仅记录 model name，不记录 key env 名称
- 或在敏感输出前检查 verbose 模式

---

### 13. 路径注入风险

**位置**: `evaluator/service.go` 多处

**严重程度**: 中

**问题描述**:
从外部输入读取 `samplePath`、`testPath` 后直接使用 `filepath.Join` 和 `os.ReadFile`，未校验路径是否在预期目录内。恶意样本可能访问意外文件。

**改进建议**:
- 使用 `filepath.Rel` 校验路径是否在 dataset_root/workdir 范围内
- 禁止 `../` 跳出预期目录

```go
func validatePathWithinRoot(path, root string) error {
    rel, err := filepath.Rel(root, path)
    if err != nil {
        return err
    }
    if strings.HasPrefix(rel, "..") {
        return fmt.Errorf("path escapes root directory")
    }
    return nil
}
```

---

### 14. 生成的测试代码直接执行

**位置**: `evaluator/service.go`

**严重程度**: 中

**问题描述**:
模型生成的测试代码直接写入文件并执行（pytest/go test），若模型返回恶意代码（如 `import os; os.system('rm -rf /')`），存在执行风险。

**改进建议**:
- 增加代码静态分析检查，禁止危险调用（如 `os.system`, `subprocess`, `exec`）
- 或限制执行环境（容器化，如 Docker/沙箱）
- 或在 prompt 中明确禁止此类代码

---

## 四、架构设计问题 (P1/P2)

### 15. 缺乏接口抽象，难以 Mock

**严重程度**: 中

**问题描述**:
`runner.Service`, `evaluator.Service` 等均为具体结构体，测试无法 mock HTTP Client 或 subprocess 调用。当前测试依赖真实文件系统和外部工具。

**改进建议**:
- 定义接口（如 `ModelAPIClient`, `Executor`）
- 使用依赖注入，便于测试和生产切换实现

```go
type ModelAPIClient interface {
    GenerateTest(ctx context.Context, model, language, source string) (string, error)
}

type Executor interface {
    Run(ctx context.Context, cmd string, args ...string) (string, error)
}
```

---

### 16. 硬编码工具路径/版本

**严重程度**: 中

**问题描述**:
C++ eval 硬编码 `clang-15`, `mull-runner-15`；Java eval 硬编码 Maven；Go eval 硬编码 `go-mutesting`。不同环境版本不匹配会失败。

**代码片段**:
```go
// cpp_eval.go
cmakeCmd.Env = append(os.Environ(), "CC=clang-15", "CXX=clang++-15")

// go_eval.go
candidates := []string{"go-mutesting", ...}
```

**改进建议**:
- 支持配置化工具路径
- 或提供 `--toolchain-config` 参数
- 或自动检测可用版本

---

### 17. Reporter HTML 过大

**位置**: `reporter/service.go:425-691`

**严重程度**: 低

**问题描述**:
`buildHTML` 函数生成 800+ 行内联 HTML/JS/CSS，不利于维护和定制。用户无法自定义报告样式。

**改进建议**:
- 分离 HTML 模板为独立文件（如 `templates/report.html`）
- 支持模板覆盖（`--template-dir` 参数）
- 提取 JS 为单独 CDN 或本地文件

---

### 18. 缺少结构化日志级别

**位置**: `obs/logger.go`

**严重程度**: 低

**问题描述**:
`obs.Logger` 仅区分 verbose/non-verbose，无 INFO/WARN/ERROR 分级。大量评测时日志难以过滤关键信息。

**改进建议**:
- 增加日志级别配置（DEBUG/INFO/WARN/ERROR）
- 支持按 module/级别过滤

---

## 五、性能问题 (P2)

### 19. Worker 数计算未考虑实际任务类型

**位置**: 
- `runner/service.go:72`
- `evaluator/service.go:65`

**严重程度**: 低

**问题描述**:
`min(16, max(2, runtime.NumCPU()))` 对 I/O 密集型任务（HTTP 调用）可能太少，对 CPU 密集型（变异测试）可能太多。

**改进建议**:
- Runner 和 Evaluator 分别配置不同 worker 数
- Runner 可用更高并发（如 32）- I/O 密集
- Evaluator 用较低并发（如 CPU 数）- CPU 密集

---

### 20. 内存重复分配

**位置**: `evaluator/service.go` 多处 trimErr

**严重程度**: 低

**问题描述**:
每次调用 `trimErr` 都创建新字符串，高并发时 GC 压力大。

**改进建议**:
- 使用固定大小的 buffer pool 或预分配
- 或直接截断而非创建新字符串

---

### 21. 覆盖率 JSON 重复写入

**位置**: `evaluator/service.go:795-815`

**严重程度**: 低

**问题描述**:
coverage run 和 json export 分两次调用 pytest，实际可合并为一次。

**改进建议**:
- 合并为单次 pytest 调用
- 使用 `coverage run -m pytest ... --cov-report=json`

---

## 六、可维护性问题 (P2)

### 22. 测试覆盖率严重不足

**位置**: `*_test.go` 文件

**严重程度**: 中

**问题描述**:
仅 12 个测试函数，覆盖率约 10-15%。关键逻辑如 `Evaluate()`, `Generate()`, `collectPythonMutation()` 无测试。

**当前测试统计**:
| 模块 | 测试文件 | 测试数 |
|------|----------|--------|
| dataset | `index_test.go`, `service_test.go` | 3 |
| evaluator | `cpp_eval_test.go`, `service_test.go` | 3 |
| reporter | `service_test.go` | 1 |
| runner | `api_test.go`, `checkpoint_test.go` | 5 |
| store | `sqlite_test.go` | 1 |

**改进建议**:
- 按模块补充单元测试，优先覆盖错误处理路径
- 建议目标覆盖率: 60%+

---

### 23. 缺少错误码规范

**位置**: `contracts/errors.go`

**严重程度**: 低

**问题描述**:
`ErrorInfo.Kind` 为自由字符串，无枚举或常量定义，难以统一处理。

**改进建议**:
- 定义 `ErrorKind` 常量枚举

```go
type ErrorKind string

const (
    ErrorKindAuthConfig    ErrorKind = "auth_config_error"
    ErrorKindNetwork       ErrorKind = "network_error"
    ErrorKindTimeout       ErrorKind = "timeout"
    ErrorKindQuality       ErrorKind = "quality_error"
    ErrorKindPayload       ErrorKind = "payload_error"
    ErrorKindRequestBuild  ErrorKind = "request_build_error"
    ErrorKindRead          ErrorKind = "read_error"
)
```

---

### 24. 重复代码：语言 eval 函数

**位置**: 
- `evaluator/go_eval.go`
- `evaluator/java_eval.go`
- `evaluator/cpp_eval.go`

**严重程度**: 低

**问题描述**:
各语言 eval 代码结构相似（workspace prep → compile → test → coverage → mutation），大量重复逻辑。

**改进建议**:
- 提取公共流程框架
- 语言差异通过 Strategy 模式实现

```go
type LanguageEvaluator interface {
    PrepareWorkspace(testPath, samplePath string) (workdir, testFile, sourceFile, error string)
    Compile(workdir, testFile string) (pass bool, error string)
    ExecuteTests(workdir, testFile string) (pass bool, output string, latency int)
    CollectCoverage(workdir, testFile, sourceFile string) (lineCov, branchCov float64, error string)
    CollectMutation(ctx context.Context, workdir, testFile string, timeout int) (score float64, stats mutationStats, error string)
}
```

---

### 25. 配置硬编码

**位置**: `cmd/utbench/main.go`

**严重程度**: 低

**问题描述**:
默认路径（`../benchmark/config/models.yaml`, `./datasets`, `./artifacts`）硬编码，不便于不同环境部署。

**改进建议**:
- 支持环境变量覆盖默认路径

```go
func getDefaultConfigPath() string {
    if p := os.Getenv("UTBENCH_CONFIG_PATH"); p != "" {
        return p
    }
    return "../benchmark/config/models.yaml"
}
```

---

## 七、文档问题 (P2)

### 26. module_level 样本元数据格式未说明

**严重程度**: 低

**问题描述**:
`readme.md` 未说明 module_level 样本的 `meta.json` 格式要求。

**改进建议**:
- 补充 module_level 样本元数据规范文档
- 示例 `meta.json` 结构:

```json
{
    "sample_id": "simple_function_000",
    "module_import": "mypackage",
    "package_name": "mypackage",
    "target_file": "main.py",
    "workspace_root": "./workspace",
    "requirements": ["pytest", "requests"]
}
```

---

### 27. 缺少 API 错误码和重试策略文档

**严重程度**: 低

**问题描述**:
无文档说明各错误类型及对应重试策略。

**改进建议**:
- 在 `docs/` 中补充错误处理指南
- 列出可重试错误类型（network_error, timeout, 429, 5xx）vs 不可重试（auth_error, quality_error）

---

### 28. architecture-mvp.md 未更新当前实现状态

**位置**: `docs/architecture-mvp.md`

**严重程度**: 低

**问题描述**:
文档未更新当前实现状态（已完成的部分）。

**改进建议**:
- 更新 v0.2 状态标记
- 补充各语言 evaluator 实现进度表

---

## 八、优先级排序

### 立即修复 (P0)
1. 资源泄漏：临时目录未清理 (#1)
2. 潜在的 panic 未恢复 (#2)
3. 并发竞态：checkpoint 写入 (#3)
4. exec.CommandContext 超时后的进程残留 (#4)

### 一周内修复 (P1)
5. HTTP Client 超时过长 (#5)
6. 随机数生成器未 Seed (#6)
7. 覆盖率解析逻辑过于简化 (#7)
8. 未处理 io.ReadAll 错误 (#8)
9. 测试计数解析覆盖不全 (#9)
10. SQLite Schema 升级不安全 (#10)
11. 缺少 Context 传播验证 (#11)
12. API Key 在日志中可能泄露 (#12)
13. 路径注入风险 (#13)
14. 生成的测试代码直接执行 (#14)
15. 缺乏接口抽象 (#15)
16. 硬编码工具路径/版本 (#16)

### 迭代中改进 (P2)
17. Reporter HTML 过大 (#17)
18. 缺少结构化日志级别 (#18)
19. Worker 数计算未考虑任务类型 (#19)
20. 内存重复分配 (#20)
21. 覆盖率 JSON 重复写入 (#21)
22. 测试覆盖率严重不足 (#22)
23. 缺少错误码规范 (#23)
24. 重复代码：语言 eval 函数 (#24)
25. 配置硬编码 (#25)
26-28. 文档问题

---

## 九、修复验证建议

每次修复后应执行:

```bash
# 1. 单元测试
go test ./internal/... -v -race

# 2. 构建
go build ./cmd/utbench

# 3. 端到端验证（dry-run）
./utbench run --dry-run --langs python --max-samples 3

# 4. 端到端验证（真实 API）
./utbench run --models deepseek --langs python --max-samples 2 --mutation-enabled
```

---

*报告生成工具: OpenCode*
*审查标准: Go 最佳实践、安全规范、并发模式*