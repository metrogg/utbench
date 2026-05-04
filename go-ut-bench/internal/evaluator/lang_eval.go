package evaluator

import (
	"context"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"
)

// WorkspaceContext 工作空间上下文
// 封装各语言 prepare*Workspace 返回的不同元组为统一结构
type WorkspaceContext struct {
	// Workdir 临时工作目录路径
	Workdir string
	// TestPath 测试文件名（相对 workdir）
	TestPath string
	// SourceBase 源文件名（不含扩展名），用于覆盖率收集和变异测试
	SourceBase string
	// SourceStem 源文件名（含扩展名），Python/C++ 特有
	SourceStem string
	// Extra 语言特有数据，如 Java 的 className、Python 的 packageName 等
	Extra map[string]string
	// ShouldCleanup 是否需要在评测结束后清理工作目录
	// Python repo_level 模式下为 false（复用 in-place）
	ShouldCleanup bool
}

// LanguageEvaluator 语言评测器接口
// 每种语言实现此接口，提供该语言的编译、测试、覆盖率、变异测试能力
//
// 实现要求：
//   - PrepareWorkspace: 创建临时工作目录，复制源码和测试文件，生成构建配置
//   - CompileCheck: 运行语言特定的编译检查，返回是否通过及错误信息
//   - ExecuteTests: 执行测试，返回是否通过、原始输出、运行耗时(ms)
//   - ParseTestCounts: 从测试输出中解析通过/总数
//   - EstimateAssertionDensity: 估算断言密度
//   - CollectCoverage: 收集行覆盖率和分支覆盖率
//   - CollectMutation: 执行变异测试，返回变异分数和统计信息
//   - MutationTool: 返回变异测试工具名称（如 "mutmut"、"go-mutesting"）
type LanguageEvaluator interface {
	// PrepareWorkspace 准备评测工作空间
	// 返回 WorkspaceContext 和错误。错误非空时评测终止。
	PrepareWorkspace(item contracts.GeneratedCase) (*WorkspaceContext, error)

	// CompileCheck 编译检查
	// 返回 (pass bool, compileError string)
	CompileCheck(ws *WorkspaceContext) (bool, string)

	// ExecuteTests 执行测试
	// 返回 (pass bool, testOutput string, runtimeMS int)
	// testOutput 用于后续 ParseTestCounts 和变异测试
	ExecuteTests(ws *WorkspaceContext, timeoutSeconds int) (bool, string, int)

	// ParseTestCounts 从测试输出解析通过数和总数
	// 返回 (*passed, *total)，任一为 nil 表示无法解析
	ParseTestCounts(testOutput string) (*int, *int)

	// EstimateAssertionDensity 估算断言密度
	// 参数为工作目录和测试文件名，由实现自行构建完整路径
	// 返回 (assertionCount, testCaseCount, density)
	EstimateAssertionDensity(workdir, testPath string) (int, int, float64)

	// CollectCoverage 收集代码覆盖率
	// 返回 (lineCoverage, branchCoverage, error)
	CollectCoverage(ws *WorkspaceContext, testPassed bool, timeoutSeconds int) (float64, float64, string)

	// CollectMutation 执行变异测试
	// 返回 (mutationScore, mutationStats, error)
	CollectMutation(ctx context.Context, ws *WorkspaceContext, input MutationInput) (float64, mutationStats, string)

	// MutationTool 返回变异测试工具名称
	MutationTool() string
}

// MutationInput 变异测试输入参数
// 统一各语言变异测试函数的输入
type MutationInput struct {
	// TimeoutSeconds 变异测试超时时间
	TimeoutSeconds int
	// TestPassRate 测试通过率
	TestPassRate *float64
	// TestPassed 已通过的测试数
	TestPassed int
	// TestTotal 总测试数
	TestTotal int
	// TestOutput 原始测试输出，Python mutmut 需要
	TestOutput string
	// Extra 语言特有参数
	Extra map[string]string
}

// evalRegistry 语言评测器注册表
var evalRegistry = map[string]LanguageEvaluator{}

// RegisterLanguageEvaluator 注册语言评测器
// 语言名称统一使用小写（python, go, java, cpp）
func RegisterLanguageEvaluator(lang string, eval LanguageEvaluator) {
	evalRegistry[strings.ToLower(lang)] = eval
}

// getLanguageEvaluator 获取指定语言的评测器
// 返回 nil 表示该语言未注册
func getLanguageEvaluator(lang string) LanguageEvaluator {
	return evalRegistry[strings.ToLower(lang)]
}

// computeTestPassRate 从测试计数计算通过率
// 同时更新 row 的 TestPassCount、TestTotalCount、TestPassRate 字段
// 这是所有语言共用的逻辑，消除 evaluateOne 中的重复代码
func computeTestPassRate(row *contracts.EvaluationResult, passCnt, totalCnt *int) {
	if passCnt != nil {
		row.TestPassCount = passCnt
	}
	if totalCnt != nil {
		row.TestTotalCount = totalCnt
	}
	if passCnt != nil && totalCnt != nil && *totalCnt > 0 {
		rate := round(float64(*passCnt)/float64(*totalCnt), 6)
		row.TestPassRate = &rate
	}
	ensureTestCountsFromPass(row)
}

// populateMutationResult 将变异测试结果填充到评测结果行
// 所有语言共用的变异结果写入逻辑
func populateMutationResult(row *contracts.EvaluationResult, score float64, stats mutationStats, mutationErr, tool string) {
	if mutationErr != "" {
		row.MutationError = mutationErr
	} else {
		row.MutationScore = &score
	}
	if stats.Total > 0 {
		total := stats.Total
		killed := stats.Killed
		survived := stats.Survived
		noTests := stats.NoTests
		timeouts := stats.Timeout
		skipped := stats.Skipped
		suspicious := stats.Suspicious
		row.MutationTotal = &total
		row.MutationKilled = &killed
		row.MutationSurvived = &survived
		row.MutationNoTests = &noTests
		row.MutationTimeouts = &timeouts
		row.MutationSkipped = &skipped
		row.MutationSuspicious = &suspicious
	}
	row.MutationTool = tool
}

// populateCoverageResult 将覆盖率结果填充到评测结果行
// 所有语言共用的覆盖率写入逻辑
func populateCoverageResult(row *contracts.EvaluationResult, lineCov, branchCov float64, covErr string, testPassed bool) {
	if covErr != "" && testPassed {
		row.CoverageError = covErr
	} else if covErr == "" {
		row.LineCoverage = &lineCov
		row.BranchCoverage = &branchCov
	}
}

// evalWithLanguageEvaluator 使用 LanguageEvaluator 接口执行评测
// 这是重构后的 evaluateOne 核心语言分发逻辑
// RuntimeMS 由 evaluateOne 的 defer finalizeEvaluationResult 统一处理
func (s *Service) evalWithLanguageEvaluator(
	ctx context.Context,
	spec contracts.RunSpec,
	item contracts.GeneratedCase,
	row *contracts.EvaluationResult,
	start time.Time,
	setPhase func(string),
	langEval LanguageEvaluator,
) {
	lang := strings.ToLower(item.Language)

	// 1. PrepareWorkspace
	setPhase(lang + ".prepare")
	ws, prepErr := langEval.PrepareWorkspace(item)
	if prepErr != nil {
		row.CompilePass = false
		row.CompileError = prepErr.Error()
		return
	}
	if ws.ShouldCleanup {
		defer cleanupWorkspaceAsync(ws.Workdir, item.Model, item.Language, item.SampleID, s.logger)
	}

	// 2. CompileCheck
	setPhase(lang + ".compile")
	compilePass, compileErr := langEval.CompileCheck(ws)
	row.CompilePass = compilePass
	if !compilePass {
		row.CompileError = compileErr
		return
	}

	// 3. ExecuteTests
	testTimeout := spec.TestTimeout
	if testTimeout <= 0 {
		testTimeout = defaultTestTimeoutSeconds
	}
	setPhase(lang + ".test")
	pass, testOutput, runtimeMs := langEval.ExecuteTests(ws, testTimeout)
	row.TestPass = &pass
	if !pass && testOutput != "" {
		row.TestError = testOutput
	}
	if runtimeMs > 0 {
		row.RuntimeMS = &runtimeMs
	}

	// 4. ParseTestCounts + computeTestPassRate (共用逻辑)
	passCnt, totalCnt := langEval.ParseTestCounts(testOutput)
	computeTestPassRate(row, passCnt, totalCnt)

	// 5. EstimateAssertionDensity
	assertCnt, testCnt, density := langEval.EstimateAssertionDensity(ws.Workdir, ws.TestPath)
	row.AssertionCount = &assertCnt
	row.TestCaseCount = &testCnt
	row.AssertionDensity = &density

	// 6. CollectCoverage
	setPhase(lang + ".coverage")
	lineCov, branchCov, covErr := langEval.CollectCoverage(ws, pass, testTimeout)
	populateCoverageResult(row, lineCov, branchCov, covErr, pass)

	// 7. CollectMutation (条件执行)
	if spec.MutationEnabled && !strings.EqualFold(strings.TrimSpace(spec.MutationPolicy), "skip") {
		if shouldRunMutationAfterSampleTests(*row) {
			testPassed := 0
			if row.TestPassCount != nil {
				testPassed = *row.TestPassCount
			}
			testTotal := 0
			if row.TestTotalCount != nil {
				testTotal = *row.TestTotalCount
			}
			input := MutationInput{
				TimeoutSeconds: spec.MutationTimeout,
				TestPassRate:   row.TestPassRate,
				TestPassed:     testPassed,
				TestTotal:      testTotal,
				TestOutput:     testOutput,
			}
			setPhase(lang + ".mutation")
			mutationScore, mutationStats, mutationErr := langEval.CollectMutation(ctx, ws, input)
			populateMutationResult(row, mutationScore, mutationStats, mutationErr, langEval.MutationTool())
		} else {
			zero := 0.0
			row.MutationScore = &zero
			row.MutationError = langEval.MutationTool() + ": baseline tests failed, skipping mutation"
			row.MutationTool = langEval.MutationTool()
		}
	}
}
