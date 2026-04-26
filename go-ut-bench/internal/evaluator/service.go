// evaluator 包提供单元测试评测功能
// 负责编译、运行测试、收集覆盖率、执行变异测试并生成评测报告
package evaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/obs"
)

// Service 评测服务结构
// 提供完整的测试评测流程管理
type Service struct {
	logger *obs.Logger // 日志记录器
}

var cleanupSemaphore = make(chan struct{}, 2)

// Output 评测操作的输出结果
// 包含评测结果集和结果文件路径
type Output struct {
	Result     contracts.EvaluationResultSet // 评测结果集
	ResultPath string                        // 结果JSON文件路径
}

// evalTask 评测任务结构
// 用于在worker之间传递评测任务
type evalTask struct {
	item contracts.GeneratedCase // 要评测的生成案例
}

// NewService 创建新的评测服务实例
// 参数:
//   - logger: 日志记录器实例
//
// 返回值:
//   - *Service: 新的服务实例
func NewService(logger *obs.Logger) *Service {
	SetMutationLogger(logger)
	return &Service{logger: logger}
}

// Evaluate 执行完整的评测流程
// 参数:
//   - ctx: 上下文，用于取消操作
//   - spec: 运行规格说明
//   - manifestPath: 生成的测试清单文件路径
//
// 返回值:
//   - Output: 评测结果输出
//   - error: 评测失败时的错误
//
// 功能说明:
//  1. 读取生成的测试清单
//  2. 使用worker池并行评测每个样本
//  3. 对每个样本执行：编译 -> 测试 -> 覆盖率 -> 变异测试
//  4. 汇总结果并写入JSON文件
func (s *Service) Evaluate(ctx context.Context, spec contracts.RunSpec, manifestPath string) (Output, error) {
	s.logger.Debug(
		"evaluate options",
		"mutation_enabled", spec.MutationEnabled,
		"mutation_policy", spec.MutationPolicy,
		"mutation_timeout", spec.MutationTimeout,
	)
	manifest, err := contracts.ReadGeneratedManifest(manifestPath)
	if err != nil {
		return Output{}, err
	}

	// 采集评测环境指纹
	envFingerprint := CaptureEnvironmentFingerprint(ctx, false, "")
	s.logger.Debug("environment fingerprint captured", "fingerprint", envFingerprint.FingerprintHash())

	// 计算worker数量
	workerCount := spec.Workers
	if workerCount <= 0 {
		workerCount = min(8, max(2, runtime.NumCPU()))
	}

	// 输出评测配置信息
	total := len(manifest.Cases)
	progress := obs.NewProgressReporter(total, "evaluate")
	progress.PrintStageStart("评测测试", fmt.Sprintf("样本: %d | 变异: %v | Workers: %d",
		total, spec.MutationEnabled, workerCount))

	// 创建输出目录
	runRoot := filepath.Join(spec.OutputRoot, "runs", spec.RunID)
	evalRoot := filepath.Join(runRoot, "evaluation")
	if err := os.MkdirAll(evalRoot, 0o755); err != nil {
		return Output{}, err
	}
	tasks := make(chan evalTask, workerCount*2)
	results := make(chan contracts.EvaluationResult, workerCount*2)
	tracker := newActiveEvalTracker(s.logger)
	watchdogDone := make(chan struct{})
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-watchdogDone:
				return
			case <-ticker.C:
				tracker.printStalled(60 * time.Second)
			}
		}
	}()
	defer close(watchdogDone)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				// 捕获worker中的panic，防止整个程序崩溃
				if r := recover(); r != nil {
					fmt.Fprintf(os.Stderr, "[worker panic] %v\n", r)
				}
			}()
			for t := range tasks {
				key, setPhase, done := tracker.start(t.item)
				_ = key
				item := s.evaluateOne(ctx, spec, t.item, setPhase)
				done()
				select {
				case <-ctx.Done():
					return
				case results <- item:
				}
			}
		}()
	}

	// 生产者：发送所有评测任务
	go func() {
		defer close(tasks)
		for _, item := range manifest.Cases {
			tasks <- evalTask{item: item}
		}
	}()

	// 消费者：收集结果
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集并显示评测结果
	rows := make([]contracts.EvaluationResult, 0, len(manifest.Cases))
	completed := 0
	for row := range results {
		completed++
		rows = append(rows, row)

		taskResult := obs.TaskResult{
			Model:         row.Model,
			Language:      row.Language,
			SampleID:      row.SampleID,
			Success:       row.CompilePass,
			CompilePass:   row.CompilePass,
			TestPass:      row.TestPass != nil && *row.TestPass,
			LineCoverage:  getCoverageValue(row.LineCoverage),
			MutationScore: getCoverageValue(row.MutationScore),
		}
		if row.MutationTool != "" {
			taskResult.MutationTool = row.MutationTool
			if row.MutationTotal != nil {
				taskResult.MutationTotal = *row.MutationTotal
			}
			if row.MutationKilled != nil {
				taskResult.MutationKilled = *row.MutationKilled
			}
			if row.MutationSurvived != nil {
				taskResult.MutationSurvived = *row.MutationSurvived
			}
		}
		if !row.CompilePass {
			taskResult.Error = row.CompileError
		} else if row.TestPass != nil && !*row.TestPass {
			taskResult.Error = row.TestError
		}
		progress.OnTaskDone(taskResult)

		status := "PASS"
		if !row.CompilePass {
			status = "FAIL(compile)"
		} else if row.TestPass != nil && !*row.TestPass {
			status = "FAIL(test)"
		}
		progress.PrintTaskLine(completed, total, row.Model, row.Language, row.SampleID, status, fmt.Sprintf("%dms", getRuntimeMS(row.RuntimeMS)))

		if completed%5 == 0 {
			progress.PrintStats()
		}
	}

	progress.PrintStats()
	progress.PrintStageDone("评测测试", obs.StageStats{
		Total:    total,
		Success:  completed,
		Duration: time.Since(progress.GetStartTime()),
	})

	// 检查是否被取消
	if err := ctx.Err(); err != nil {
		return Output{}, err
	}

	// 排序结果：按模型 -> 语言 -> 样本ID
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Model == rows[j].Model {
			if rows[i].Language == rows[j].Language {
				return rows[i].SampleID < rows[j].SampleID
			}
			return rows[i].Language < rows[j].Language
		}
		return rows[i].Model < rows[j].Model
	})

	// 构建结果集
	set := contracts.EvaluationResultSet{
		SchemaVersion:          contracts.SchemaVersion,
		RunID:                  spec.RunID,
		EvaluatedAtUTC:         time.Now().UTC(),
		ManifestPath:           manifestPath,
		Results:                rows,
		EnvironmentFingerprint: envFingerprint.FingerprintHash(),
		EnvironmentJSON:        envFingerprint.ToJSON(),
	}
	resultPath := filepath.Join(evalRoot, "evaluation_result.json")
	if err := contracts.WriteJSON(resultPath, set); err != nil {
		return Output{}, err
	}

	// 如果启用变异测试且策略为fail，检查是否有错误
	if spec.MutationEnabled && strings.EqualFold(strings.TrimSpace(spec.MutationPolicy), "fail") {
		mutationErrCount := 0
		for _, row := range rows {
			if row.MutationError != "" {
				mutationErrCount++
			}
		}
		if mutationErrCount > 0 {
			return Output{Result: set, ResultPath: resultPath}, fmt.Errorf("mutation stage failed on %d samples", mutationErrCount)
		}
	}

	s.logger.Info("evaluation finished", "run_id", spec.RunID, "total_results", len(rows), "path", resultPath)
	return Output{Result: set, ResultPath: resultPath}, nil
}

func (s *Service) evaluateOne(ctx context.Context, spec contracts.RunSpec, item contracts.GeneratedCase, setPhase func(string)) (result contracts.EvaluationResult) {
	start := time.Now()
	row := contracts.EvaluationResult{
		Model:             item.Model,
		Language:          item.Language,
		SampleID:          item.SampleID,
		GeneratedTestPath: item.GeneratedTestPath,
		SourcePath:        item.SamplePath,
		PromptTokens:      item.PromptTokens,
		CompletionTokens:  item.CompletionTokens,
		TotalTokens:       item.TotalTokens,
		Truncated:         item.Truncated,
	}
	defer func() {
		finalizeEvaluationResult(&row, start)
		result = row
	}()

	if !item.Success {
		row.CompilePass = false
		if item.Error != nil {
			row.CompileError = "generation failed: " + item.Error.Message
		} else {
			row.CompileError = "generation failed"
		}
		return row
	}

	s.logger.Debug("evaluating", "model", item.Model, "lang", item.Language, "sample", item.SampleID)

	if strings.EqualFold(item.Language, "python") {
		setPhase("python.prepare")
		var workdir, testName, sourceBase, sourceStem, packageName, targetFile string
		var prepErr string
		isModuleLevel := isModuleLevelSample(item.SamplePath)
		if isModuleLevel {
			workdir, testName, packageName, targetFile, prepErr = preparePythonModuleLevelWorkspace(item.GeneratedTestPath, item.SamplePath)
		} else {
			workdir, testName, sourceBase, sourceStem, prepErr = preparePythonWorkspace(item.GeneratedTestPath, item.SamplePath)
		}
		if prepErr != "" {
			row.CompilePass = false
			row.CompileError = prepErr
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}
		if !isModuleLevel {
			defer cleanupWorkspaceAsync(workdir, item.Model, item.Language, item.SampleID, s.logger)
		}

		setPhase("python.compile")
		compilePass, compileErr := pythonCompileCheck(filepath.Join(workdir, testName))
		row.CompilePass = compilePass
		if !compilePass {
			row.CompileError = compileErr
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}

		testTimeout := spec.TestTimeout
		if testTimeout <= 0 {
			testTimeout = defaultTestTimeoutSeconds
		}

		var pass bool
		var testErr string
		var runtimeMs int
		setPhase("python.test")
		if isModuleLevel {
			pass, testErr, runtimeMs = executePythonTestsInWorkspace(workdir, testName, packageName, testTimeout)
		} else {
			pass, testErr, runtimeMs = executePythonTests(workdir, testName, testTimeout)
		}
		row.TestPass = &pass
		if !pass && testErr != "" {
			row.TestError = testErr
		}
		if runtimeMs > 0 {
			row.RuntimeMS = &runtimeMs
		}
		passCnt, totalCnt := parsePytestCounts(testErr)
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
		ensureTestCountsFromPass(&row)

		assertCnt, testCnt, density := estimateAssertionDensity(filepath.Join(workdir, testName), item.Language)
		row.AssertionCount = &assertCnt
		row.TestCaseCount = &testCnt
		row.AssertionDensity = &density

		var targets []string
		if isModuleLevel {
			if targetFile != "" {
				targets = []string{targetFile}
			}
		} else {
			targets = inferMutationTargets(workdir, testName, sourceBase)
		}

		if isModuleLevel {
			if packageName != "" {
				setPhase("python.coverage")
				lineCov, branchCov, covErr := collectPythonCoverageInWorkspace(workdir, testName, packageName, targetFile, testTimeout)
				if covErr != "" && pass {
					row.CoverageError = covErr
				} else if covErr == "" {
					row.LineCoverage = &lineCov
					row.BranchCoverage = &branchCov
				}
			}
		} else if sourceBase != "" {
			setPhase("python.coverage")
			lineCov, branchCov, covErr := collectPythonCoverage(workdir, testName, sourceBase, sourceStem, targets, testTimeout)
			if covErr != "" && pass {
				row.CoverageError = covErr
			} else if covErr == "" {
				row.LineCoverage = &lineCov
				row.BranchCoverage = &branchCov
			}
		}

		if spec.MutationEnabled && !strings.EqualFold(strings.TrimSpace(spec.MutationPolicy), "skip") {
			mutationTargets := targets
			if len(mutationTargets) == 0 && sourceBase != "" {
				mutationTargets = []string{sourceBase}
			}
			testPassed := 0
			if row.TestPassCount != nil {
				testPassed = *row.TestPassCount
			}
			testTotal := 0
			if row.TestTotalCount != nil {
				testTotal = *row.TestTotalCount
			}
			mutationSkipReason := ""
			if !shouldRunMutationAfterSampleTests(row) {
				mutationSkipReason = "mutmut: baseline tests failed, skipping mutation"
			} else if !isModuleLevel && !pythonTestImportsAnyMutationTarget(workdir, testName, mutationTargets) {
				mutationSkipReason = "mutmut: generated tests do not import mutation target, skipping mutation"
			}
			if mutationSkipReason != "" {
				zero := 0.0
				row.MutationScore = &zero
				row.MutationError = mutationSkipReason
				row.MutationTool = "mutmut"
			} else {
				checkResult := CheckTestPassRate(testPassed, testTotal, "mutmut", GetMinPassRateForTool("mutmut"))
				if !checkResult.ShouldRun {
					zero := 0.0
					row.MutationScore = &zero
					row.MutationError = checkResult.Message
					row.MutationTool = "mutmut"
				} else {
					setPhase("python.mutation")
					mutationStart := time.Now()
					mutationScore, mutationStats, mutationErr := collectPythonMutation(ctx, workdir, testName, mutationTargets, spec.MutationTimeout, testErr)
					mutationElapsed := time.Since(mutationStart)
					s.logger.ToFile("evaluator").Trace("mutation_result",
						"model", item.Model,
						"language", item.Language,
						"sample_id", item.SampleID,
						"tool", "mutmut",
						"score", mutationScore,
						"total", mutationStats.Total,
						"killed", mutationStats.Killed,
						"survived", mutationStats.Survived,
						"elapsed_seconds", int(mutationElapsed.Seconds()),
						"error", mutationErr,
					)
					if mutationErr != "" {
						row.MutationError = mutationErr
					} else {
						row.MutationScore = &mutationScore
					}
					if mutationStats.Total > 0 {
						total := mutationStats.Total
						killed := mutationStats.Killed
						survived := mutationStats.Survived
						noTests := mutationStats.NoTests
						timeouts := mutationStats.Timeout
						skipped := mutationStats.Skipped
						suspicious := mutationStats.Suspicious
						row.MutationTotal = &total
						row.MutationKilled = &killed
						row.MutationSurvived = &survived
						row.MutationNoTests = &noTests
						row.MutationTimeouts = &timeouts
						row.MutationSkipped = &skipped
						row.MutationSuspicious = &suspicious
					}
					row.MutationTool = "mutmut"
				}
			}
		}
	} else if strings.EqualFold(item.Language, "go") {
		setPhase("go.prepare")
		workdir, testName, sourceBase, _ := prepareGoWorkspace(item.GeneratedTestPath, item.SamplePath)
		if workdir == "" {
			row.CompilePass = false
			row.CompileError = testName // prepErr stored in testName
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}
		defer cleanupWorkspaceAsync(workdir, item.Model, item.Language, item.SampleID, s.logger)

		setPhase("go.compile")
		compilePass, compileErr := goCompileCheck(workdir, testName)
		row.CompilePass = compilePass
		if !compilePass {
			row.CompileError = compileErr
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}

		setPhase("go.test")
		pass, testErr, runtimeMs := executeGoTests(workdir, testName, sourceBase)
		row.TestPass = &pass
		if !pass && testErr != "" {
			row.TestError = testErr
		}
		if runtimeMs > 0 {
			row.RuntimeMS = &runtimeMs
		}

		passCnt, totalCnt := parseGoTestCounts(testErr)
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
		ensureTestCountsFromPass(&row)

		assertCnt, testCnt, density := estimateAssertionDensity(filepath.Join(workdir, testName), item.Language)
		row.AssertionCount = &assertCnt
		row.TestCaseCount = &testCnt
		row.AssertionDensity = &density

		if sourceBase != "" {
			setPhase("go.coverage")
			lineCov, branchCov, covErr := collectGoCoverage(workdir, testName, sourceBase)
			if covErr != "" && pass {
				row.CoverageError = covErr
			} else if covErr == "" {
				row.LineCoverage = &lineCov
				row.BranchCoverage = &branchCov
			}
		}

		if spec.MutationEnabled && !strings.EqualFold(strings.TrimSpace(spec.MutationPolicy), "skip") {
			setPhase("go.mutation")
			mutationStart := time.Now()
			mutationScore := 0.0
			mutationStats := mutationStats{}
			mutationErr := ""
			if !shouldRunMutationAfterSampleTests(row) {
				mutationErr = "go-mutesting: baseline tests failed, skipping mutation"
			} else {
				testPassed := 0
				if row.TestPassCount != nil {
					testPassed = *row.TestPassCount
				}
				testTotal := 0
				if row.TestTotalCount != nil {
					testTotal = *row.TestTotalCount
				}
				mutationScore, mutationStats, mutationErr = collectGoMutation(ctx, workdir, testName, sourceBase, spec.MutationTimeout, row.TestPassRate, testPassed, testTotal)
			}
			mutationElapsed := time.Since(mutationStart)
			s.logger.ToFile("evaluator").Trace("mutation_result",
				"model", item.Model,
				"language", item.Language,
				"sample_id", item.SampleID,
				"tool", "go-mutesting",
				"score", mutationScore,
				"total", mutationStats.Total,
				"killed", mutationStats.Killed,
				"survived", mutationStats.Survived,
				"elapsed_seconds", int(mutationElapsed.Seconds()),
				"error", mutationErr,
			)
			row.MutationScore = &mutationScore
			row.MutationTotal = &mutationStats.Total
			row.MutationKilled = &mutationStats.Killed
			row.MutationSurvived = &mutationStats.Survived
			row.MutationNoTests = &mutationStats.NoTests
			row.MutationTimeouts = &mutationStats.Timeout
			row.MutationSkipped = &mutationStats.Skipped
			row.MutationSuspicious = &mutationStats.Suspicious
			if mutationErr != "" {
				row.MutationError = mutationErr
			}
			row.MutationTool = "go-mutesting"
		}

	} else if strings.EqualFold(item.Language, "java") {
		setPhase("java.prepare")
		workdir, testName, _, className := prepareJavaWorkspace(item.GeneratedTestPath, item.SamplePath)
		if workdir == "" {
			row.CompilePass = false
			row.CompileError = testName // prepErr stored in testName
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}
		defer cleanupWorkspaceAsync(workdir, item.Model, item.Language, item.SampleID, s.logger)

		setPhase("java.compile")
		compilePass, compileErr := javaCompileCheck(workdir)
		row.CompilePass = compilePass
		if !compilePass {
			row.CompileError = compileErr
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}

		testTimeout := spec.TestTimeout
		if testTimeout <= 0 {
			testTimeout = 120
		}
		setPhase("java.test")
		pass, testErr, runtimeMs := executeJavaTestsWithTimeout(workdir, testTimeout)
		row.TestPass = &pass
		if !pass && testErr != "" {
			row.TestError = testErr
		}
		if runtimeMs > 0 {
			row.RuntimeMS = &runtimeMs
		}

		passCnt, totalCnt := parseJavaTestCounts(testErr)
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
		ensureTestCountsFromPass(&row)

		assertCnt, testCnt, density := estimateAssertionDensity(filepath.Join(workdir, "src", "test", "java", testName), item.Language)
		row.AssertionCount = &assertCnt
		row.TestCaseCount = &testCnt
		row.AssertionDensity = &density

		if className != "" {
			setPhase("java.coverage")
			lineCov, branchCov, covErr := collectJavaCoverage(workdir, className)
			if covErr != "" && pass {
				row.CoverageError = covErr
			} else if covErr == "" {
				row.LineCoverage = &lineCov
				row.BranchCoverage = &branchCov
			}
		}

		if spec.MutationEnabled && !strings.EqualFold(strings.TrimSpace(spec.MutationPolicy), "skip") {
			setPhase("java.mutation")
			mutationStart := time.Now()
			testPassed := 0
			if row.TestPassCount != nil {
				testPassed = *row.TestPassCount
			}
			testTotal := 0
			if row.TestTotalCount != nil {
				testTotal = *row.TestTotalCount
			}
			mutationScore := 0.0
			mutationStats := mutationStats{}
			mutationErr := ""
			if !shouldRunMutationAfterSampleTests(row) {
				mutationErr = "PITest: baseline tests failed, skipping mutation"
			} else {
				mutationScore, mutationStats, mutationErr = collectJavaMutation(ctx, workdir, className, spec.MutationTimeout, row.TestPassRate, testPassed, testTotal)
			}
			mutationElapsed := time.Since(mutationStart)
			s.logger.ToFile("evaluator").Trace("mutation_result",
				"model", item.Model,
				"language", item.Language,
				"sample_id", item.SampleID,
				"tool", "pitest",
				"score", mutationScore,
				"total", mutationStats.Total,
				"killed", mutationStats.Killed,
				"survived", mutationStats.Survived,
				"elapsed_seconds", int(mutationElapsed.Seconds()),
				"error", mutationErr,
			)
			row.MutationScore = &mutationScore
			row.MutationTotal = &mutationStats.Total
			row.MutationKilled = &mutationStats.Killed
			row.MutationSurvived = &mutationStats.Survived
			row.MutationNoTests = &mutationStats.NoTests
			row.MutationTimeouts = &mutationStats.Timeout
			row.MutationSkipped = &mutationStats.Skipped
			row.MutationSuspicious = &mutationStats.Suspicious
			if mutationErr != "" {
				row.MutationError = mutationErr
			}
			row.MutationTool = "pitest"
		}

	} else if strings.EqualFold(item.Language, "cpp") {
		setPhase("cpp.prepare")
		workdir, testName, sourceBase, _, prepErr := prepareCppWorkspace(item.GeneratedTestPath, item.SamplePath)
		if workdir == "" {
			row.CompilePass = false
			row.CompileError = prepErr
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}
		defer cleanupWorkspaceAsync(workdir, item.Model, item.Language, item.SampleID, s.logger)

		setPhase("cpp.compile")
		compilePass, compileErr := cppCompileCheck(workdir)
		row.CompilePass = compilePass
		if !compilePass {
			row.CompileError = compileErr
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}

		setPhase("cpp.test")
		pass, testErr, runtimeMs := executeCppTests(workdir)
		row.TestPass = &pass
		if !pass && testErr != "" {
			row.TestError = testErr
		}
		if runtimeMs > 0 {
			row.RuntimeMS = &runtimeMs
		}

		passCnt, totalCnt := parseCppTestCounts(testErr)
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
		ensureTestCountsFromPass(&row)

		assertCnt, testCnt, density := estimateCppAssertionDensity(filepath.Join(workdir, testName))
		row.AssertionCount = &assertCnt
		row.TestCaseCount = &testCnt
		row.AssertionDensity = &density

		if testName != "" {
			setPhase("cpp.coverage")
			lineCov, branchCov, covErr := collectCppCoverage(workdir, testName)
			if covErr != "" && pass {
				row.CoverageError = covErr
			} else if covErr == "" {
				row.LineCoverage = &lineCov
				row.BranchCoverage = &branchCov
			}
		}

		if spec.MutationEnabled && !strings.EqualFold(strings.TrimSpace(spec.MutationPolicy), "skip") {
			setPhase("cpp.mutation")
			mutationStart := time.Now()
			testPassed := 0
			if row.TestPassCount != nil {
				testPassed = *row.TestPassCount
			}
			testTotal := 0
			if row.TestTotalCount != nil {
				testTotal = *row.TestTotalCount
			}
			mutationScore := 0.0
			mutationStats := mutationStats{}
			mutationErr := ""
			if !shouldRunMutationAfterSampleTests(row) {
				mutationErr = "Mull: baseline tests failed, skipping mutation"
			} else {
				mutationScore, mutationStats, mutationErr = collectCppMutation(ctx, workdir, sourceBase, spec.MutationTimeout, row.TestPassRate, testPassed, testTotal)
			}
			mutationElapsed := time.Since(mutationStart)
			s.logger.ToFile("evaluator").Trace("mutation_result",
				"model", item.Model,
				"language", item.Language,
				"sample_id", item.SampleID,
				"tool", "mull",
				"score", mutationScore,
				"total", mutationStats.Total,
				"killed", mutationStats.Killed,
				"survived", mutationStats.Survived,
				"elapsed_seconds", int(mutationElapsed.Seconds()),
				"error", mutationErr,
			)
			row.MutationScore = &mutationScore
			row.MutationTotal = &mutationStats.Total
			row.MutationKilled = &mutationStats.Killed
			row.MutationSurvived = &mutationStats.Survived
			row.MutationNoTests = &mutationStats.NoTests
			row.MutationTimeouts = &mutationStats.Timeout
			row.MutationSkipped = &mutationStats.Skipped
			row.MutationSuspicious = &mutationStats.Suspicious
			if mutationErr != "" {
				row.MutationError = mutationErr
			}
			row.MutationTool = "mull"
		}

	} else {
		compilePass := true
		pass := true
		lineCov := 0.0
		branchCov := 0.0
		mutation := 0.0
		assertCnt, testCnt, density := estimateAssertionDensity(item.GeneratedTestPath, item.Language)

		row.CompilePass = compilePass
		row.TestPass = &pass
		row.LineCoverage = &lineCov
		row.BranchCoverage = &branchCov
		row.MutationScore = &mutation
		row.AssertionCount = &assertCnt
		row.TestCaseCount = &testCnt
		row.AssertionDensity = &density
	}

	return row
}

func finalizeEvaluationResult(row *contracts.EvaluationResult, start time.Time) {
	totalRuntimeMS := int(time.Since(start).Milliseconds())
	row.RuntimeMS = &totalRuntimeMS
	origin, reason := classifyFailureOrigin(*row)
	if origin == "" {
		origin = "none"
	}
	row.FailureOrigin = origin
	eligible := origin == "none" || origin == "model"
	row.ScoreEligible = &eligible
	if !eligible {
		row.ScoreExclusionReason = reason
	}
}

func ensureTestCountsFromPass(row *contracts.EvaluationResult) {
	// 如果已经有用例级数据，不需要补充
	if row.TestPassCount != nil || row.TestTotalCount != nil {
		return
	}
	// 如果没有 TestPass 信息，无法推断，保持 nil 表示"未知"
	if row.TestPass == nil {
		return
	}
	// 注意：这里设置的是样本级数据（样本整体是否通过）
	// 用例级数据应该通过解析测试框架输出获得
	// 如果解析失败，保持 nil 是正确的做法，不应该强制设置默认值
	// 因为这会混淆样本级和用例级的概念
}

func shouldRunMutationAfterSampleTests(row contracts.EvaluationResult) bool {
	if row.TestPass != nil && !*row.TestPass {
		return false
	}
	return true
}

func pythonTestImportsAnyMutationTarget(workdir, testName string, mutationTargets []string) bool {
	if len(mutationTargets) == 0 {
		return false
	}
	raw, err := os.ReadFile(filepath.Join(workdir, testName))
	if err != nil {
		return true
	}
	text := string(raw)
	for _, target := range mutationTargets {
		module := strings.TrimSuffix(filepath.ToSlash(target), ".py")
		module = strings.Trim(module, "/")
		module = strings.ReplaceAll(module, "/", ".")
		if module == "" {
			continue
		}
		quoted := regexp.QuoteMeta(module)
		fromRe := regexp.MustCompile(`(?m)^\s*from\s+` + quoted + `\s+import\b`)
		importRe := regexp.MustCompile(`(?m)^\s*import\s+(?:[a-zA-Z_][a-zA-Z0-9_]*\s*,\s*)*` + quoted + `(?:\s+as\s+[a-zA-Z_][a-zA-Z0-9_]*)?(?:\s*(?:,|$))`)
		if fromRe.MatchString(text) || importRe.MatchString(text) {
			return true
		}
	}
	return false
}

func classifyFailureOrigin(row contracts.EvaluationResult) (string, string) {
	if row.CompileError == "" && row.TestError == "" && row.CoverageError == "" && row.MutationError == "" && !row.Truncated {
		return "none", ""
	}
	for _, msg := range []string{row.CompileError, row.TestError, row.CoverageError, row.MutationError} {
		if msg == "" {
			continue
		}
		if isDatasetFailureMessage(msg) {
			return "dataset", shortFailureReason(msg)
		}
	}
	if generatedTestDidNotPass(row) {
		return "model", ""
	}
	for _, msg := range []string{row.CompileError, row.TestError, row.CoverageError, row.MutationError} {
		if msg == "" {
			continue
		}
		if isEnvironmentFailureMessage(msg) {
			return "environment", shortFailureReason(msg)
		}
		if isToolFailureMessage(msg) {
			return "tool", shortFailureReason(msg)
		}
	}
	return "model", ""
}

func generatedTestDidNotPass(row contracts.EvaluationResult) bool {
	if row.CompilePass == false && row.CompileError != "" {
		return true
	}
	if row.TestPass != nil && !*row.TestPass {
		return true
	}
	if row.TestPassRate != nil && *row.TestPassRate < 1.0 {
		return true
	}
	return false
}

func isDatasetFailureMessage(msg string) bool {
	msg = strings.ToLower(msg)
	return strings.Contains(msg, "dataset root") ||
		strings.Contains(msg, "source file not found") ||
		strings.Contains(msg, "target file not found") ||
		strings.Contains(msg, "module_level sample missing metadata") ||
		strings.Contains(msg, "module_level workspace not found") ||
		strings.Contains(msg, "module_level workspace_root not set") ||
		strings.Contains(msg, "failed to read source")
}

func isEnvironmentFailureMessage(msg string) bool {
	msg = strings.ToLower(msg)
	if strings.Contains(msg, "pitest") || strings.Contains(msg, "junit 5 plugin") {
		return false
	}
	return strings.Contains(msg, "permission denied") ||
		strings.Contains(msg, "access is denied") ||
		strings.Contains(msg, "executable file not found") ||
		strings.Contains(msg, "not installed") ||
		strings.Contains(msg, "command not found")
}

func isToolFailureMessage(msg string) bool {
	msg = strings.ToLower(msg)
	if strings.Contains(msg, "all tests failed") ||
		strings.Contains(msg, "no tests found, skipping mutation") ||
		strings.Contains(msg, "generated tests do not import mutation target") ||
		strings.Contains(msg, "could not find any test case for any mutant") ||
		strings.Contains(msg, "pytest timed out") ||
		strings.Contains(msg, "coverage run timed out") {
		return false
	}
	return strings.Contains(msg, "coverage json failed") ||
		strings.Contains(msg, "coverage files empty") ||
		strings.Contains(msg, "stats file not found") ||
		strings.Contains(msg, "gremlins no results to report") ||
		strings.Contains(msg, "no gremlins output found") ||
		strings.Contains(msg, "go-mutesting no results to report") ||
		strings.Contains(msg, "no go-mutesting output found") ||
		strings.Contains(msg, "no results to report") ||
		strings.Contains(msg, "produced zero mutants") ||
		strings.Contains(msg, "did not execute any mutants") ||
		strings.Contains(msg, "run incomplete") ||
		strings.Contains(msg, "parse error") ||
		strings.Contains(msg, "pitest could not run any tests") ||
		strings.Contains(msg, "pitest no killed/survived results") ||
		strings.Contains(msg, "pitest requires junit 5 plugin") ||
		strings.Contains(msg, "pitest junit 5 plugin is not installed")
}

func shortFailureReason(msg string) string {
	msg = trimErr(msg, 240)
	if msg == "" {
		return "non-model failure"
	}
	return msg
}

func preparePythonWorkspace(testPath string, sourcePath string) (string, string, string, string, string) {
	workdir, err := os.MkdirTemp("", "utbench_eval_")
	if err != nil {
		return "", "", "", "", err.Error()
	}

	raw, err := os.ReadFile(testPath)
	if err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", "", err.Error()
	}
	testSource := string(raw)
	if sourcePath != "" {
		testSource = rewriteGeneratedTestImports(testSource, sourcePath)
	}

	sourceBase := ""
	sourceStem := ""
	if sourcePath != "" {
		sourceBase = filepath.Base(sourcePath)
		sourceStem = strings.TrimSuffix(sourceBase, filepath.Ext(sourceBase))
		srcRaw, err := os.ReadFile(sourcePath)
		if err == nil {
			target := filepath.Join(workdir, sourceBase)
			_ = os.WriteFile(target, srcRaw, 0o644)
			aliases := inferAliasModules(sourceStem)
			for _, alias := range aliases {
				if alias == sourceStem || !isValidModuleName(alias) {
					continue
				}
				targetAlias := filepath.Join(workdir, alias+".py")
				_ = os.WriteFile(targetAlias, srcRaw, 0o644)
			}
		}
	}

	testName := normalizedPytestFilename(filepath.Base(testPath))
	if err := os.WriteFile(filepath.Join(workdir, testName), []byte(testSource), 0o644); err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", "", err.Error()
	}

	return workdir, testName, sourceBase, sourceStem, ""
}

func cleanupWorkspace(workdir string) {
	_ = os.RemoveAll(workdir)
}

func cleanupWorkspaceAsync(workdir, model, language, sampleID string, logger *obs.Logger) {
	if strings.TrimSpace(workdir) == "" {
		return
	}
	go func() {
		cleanupSemaphore <- struct{}{}
		defer func() { <-cleanupSemaphore }()
		start := time.Now()
		fmt.Printf("        [CLEANUP] start | %s | %s | %s | %s\n", model, language, sampleID, workdir)
		err := os.RemoveAll(workdir)
		elapsed := time.Since(start)
		if err != nil {
			fmt.Printf("        [CLEANUP-WARN] failed | %s | %s | %s | elapsed=%s | err=%v\n", model, language, sampleID, elapsed.Round(time.Second), err)
			if logger != nil {
				logger.ToFile("evaluator").Trace("cleanup_failed",
					"model", model,
					"language", language,
					"sample_id", sampleID,
					"workdir", workdir,
					"elapsed_ms", elapsed.Milliseconds(),
					"error", err.Error(),
				)
			}
			return
		}
		if elapsed >= 2*time.Second {
			fmt.Printf("        [CLEANUP] done | %s | %s | %s | elapsed=%s\n", model, language, sampleID, elapsed.Round(time.Second))
		}
		if logger != nil {
			logger.ToFile("evaluator").Trace("cleanup_done",
				"model", model,
				"language", language,
				"sample_id", sampleID,
				"workdir", workdir,
				"elapsed_ms", elapsed.Milliseconds(),
			)
		}
	}()
}

func isModuleLevelSample(samplePath string) bool {
	dir := filepath.Dir(samplePath)
	base := filepath.Base(samplePath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	if name == "entry" {
		metaPath := filepath.Join(dir, "meta.json")
		if _, err := os.Stat(metaPath); err == nil {
			return true
		}
	}
	metaPath := filepath.Join(dir, name+".meta.json")
	if _, err := os.Stat(metaPath); err == nil {
		var meta contracts.ModuleLevelMeta
		if raw, err := os.ReadFile(metaPath); err == nil {
			if err := json.Unmarshal(raw, &meta); err == nil && meta.ModuleImport != "" {
				return true
			}
		}
	}
	return false
}

func loadModuleLevelMeta(samplePath string) *contracts.ModuleLevelMeta {
	sampleDir := filepath.Dir(samplePath)
	entryBase := filepath.Base(samplePath)
	entryExt := filepath.Ext(entryBase)
	entryName := entryBase[:len(entryBase)-len(entryExt)]
	if entryName == "entry" {
		metaPath := filepath.Join(sampleDir, "meta.json")
		if raw, err := os.ReadFile(metaPath); err == nil {
			var meta contracts.ModuleLevelMeta
			if err := json.Unmarshal(raw, &meta); err == nil {
				if strings.HasPrefix(meta.WorkspaceRoot, ".") {
					meta.WorkspaceRoot = filepath.Join(sampleDir, meta.WorkspaceRoot)
				}
				return &meta
			}
		}
		return nil
	}
	dir := filepath.Dir(samplePath)
	base := filepath.Base(samplePath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	metaPath := filepath.Join(dir, name+".meta.json")
	var meta contracts.ModuleLevelMeta
	if raw, err := os.ReadFile(metaPath); err == nil {
		if err := json.Unmarshal(raw, &meta); err == nil {
			if strings.HasPrefix(meta.WorkspaceRoot, ".") {
				meta.WorkspaceRoot = filepath.Join(dir, meta.WorkspaceRoot)
			}
			return &meta
		}
	}
	return nil
}

func preparePythonModuleLevelWorkspace(testPath string, samplePath string) (string, string, string, string, string) {
	meta := loadModuleLevelMeta(samplePath)
	if meta == nil {
		return "", "", "", "", "module_level sample missing metadata"
	}
	workspaceRoot := meta.WorkspaceRoot
	if workspaceRoot == "" {
		return "", "", "", "", "module_level workspace_root not set in metadata"
	}
	if _, err := os.Stat(workspaceRoot); err != nil {
		return "", "", "", "", "module_level workspace not found: " + workspaceRoot
	}
	testFileName := normalizedPytestFilename(filepath.Base(testPath))
	generatedSrc, err := os.ReadFile(testPath)
	if err != nil {
		return "", "", "", "", "failed to read generated test: " + err.Error()
	}
	testsDir := filepath.Join(workspaceRoot, "tests")
	if err := os.MkdirAll(testsDir, 0o755); err != nil {
		return "", "", "", "", "failed to create tests dir: " + err.Error()
	}
	testDest := filepath.Join(testsDir, testFileName)
	if err := os.WriteFile(testDest, generatedSrc, 0o644); err != nil {
		return "", "", "", "", "failed to write test file: " + err.Error()
	}
	testName := filepath.Join("tests", testFileName)
	return workspaceRoot, testName, meta.PackageName, meta.TargetFile, ""
}

func executePythonTestsInWorkspace(workdir, testName, packageName string, timeoutSeconds int) (bool, string, int) {
	py := pythonExecutable()
	env := os.Environ()
	env = append(env, "PYTHONPATH="+workdir)
	runCtx, cancel := context.WithTimeout(context.Background(), normalizedTimeout(timeoutSeconds))
	defer cancel()
	started := time.Now()
	output, err := runCommandWithProcessGroupKill(runCtx, py, []string{"-m", "pytest", testName, "-q", "--maxfail=9999"}, workdir, env)
	latency := int(time.Since(started).Milliseconds())
	if runCtx.Err() != nil {
		return false, fmt.Sprintf("pytest timed out after %ds", timeoutSeconds), latency
	}
	if err == nil {
		return true, string(output), latency
	}
	return false, trimErr(string(output), 4000), latency
}

func collectPythonCoverageInWorkspace(workdir, testName, packageName, targetFile string, timeoutSeconds int) (float64, float64, string) {
	if packageName == "" {
		return 0, 0, "missing package name for module_level coverage"
	}
	py := pythonExecutable()
	absWorkdir, err := filepath.Abs(workdir)
	if err != nil {
		return 0, 0, "failed to get absolute path: " + err.Error()
	}
	jsonPath := filepath.Join(absWorkdir, ".coverage.utbench.json")
	env := os.Environ()
	env = append(env, "PYTHONPATH="+absWorkdir)
	env = append(env, "COVERAGE_FILE="+filepath.Join(absWorkdir, ".coverage.utbench"))
	runCtx, cancelRun := context.WithTimeout(context.Background(), normalizedTimeout(timeoutSeconds))
	defer cancelRun()
	runOut, runErr := runCommandWithProcessGroupKill(runCtx, py, []string{"-m", "coverage", "run", "--branch", "--source", packageName, "-m", "pytest", testName, "-q", "--maxfail=9999"}, absWorkdir, env)
	if runCtx.Err() != nil {
		return 0, 0, fmt.Sprintf("coverage run timed out after %ds", timeoutSeconds)
	}
	if runErr != nil {
		return 0, 0, "coverage run failed: " + trimErr(string(runOut), 800)
	}
	jsonCtx, cancelJSON := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelJSON()
	if out, err := runCommandWithProcessGroupKill(jsonCtx, py, []string{"-m", "coverage", "json", "-o", jsonPath}, absWorkdir, env); err != nil {
		return 0, 0, "coverage json failed: " + trimErr(string(out), 800)
	}
	raw, err := os.ReadFile(jsonPath)
	if err != nil {
		return 0, 0, err.Error()
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return 0, 0, err.Error()
	}
	files, _ := payload["files"].(map[string]any)
	if targetFile != "" && files != nil {
		for filePath, anyDetail := range files {
			if strings.Contains(filePath, targetFile) || filepath.Base(filePath) == filepath.Base(targetFile) {
				detail, _ := anyDetail.(map[string]any)
				summary, _ := detail["summary"].(map[string]any)
				line := extractLineCoverage(summary)
				branch := extractBranchCoverage(summary)
				return line, branch, ""
			}
		}
	}
	if totals, ok := payload["totals"].(map[string]any); ok {
		line := extractLineCoverage(totals)
		branch := extractBranchCoverage(totals)
		return line, branch, ""
	}
	return 0, 0, "coverage files empty"
}

func pythonCompileCheck(path string) (bool, string) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return false, err.Error()
	}
	text := string(raw)
	if strings.TrimSpace(text) == "" {
		return false, "empty test file"
	}
	py := pythonExecutable()
	runCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	output, err := runCommandWithProcessGroupKill(runCtx, py, []string{"-m", "py_compile", path}, "", nil)
	if runCtx.Err() != nil {
		return false, "python compile timed out after 30s"
	}
	if err != nil {
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			msg = err.Error()
		}
		return false, msg
	}
	if !strings.Contains(text, "def test_") && !strings.Contains(text, "import pytest") {
		return false, "invalid python test structure"
	}
	return true, ""
}

func executePythonTests(workdir, filename string, timeoutSeconds int) (bool, string, int) {
	py := pythonExecutable()
	runCtx, cancel := context.WithTimeout(context.Background(), normalizedTimeout(timeoutSeconds))
	defer cancel()
	started := time.Now()
	output, err := runCommandWithProcessGroupKill(runCtx, py, []string{"-m", "pytest", filename, "-q", "--maxfail=9999"}, workdir, nil)
	latency := int(time.Since(started).Milliseconds())
	if runCtx.Err() != nil {
		return false, fmt.Sprintf("pytest timed out after %ds", timeoutSeconds), latency
	}
	if err == nil {
		return true, string(output), latency
	}
	return false, trimErr(string(output), 4000), latency
}

func collectPythonCoverage(workdir, filename, sourceBase, sourceStem string, aliases []string, timeoutSeconds int) (float64, float64, string) {
	if sourceBase == "" {
		return 0, 0, "missing source path"
	}
	py := pythonExecutable()
	jsonPath := filepath.Join(workdir, ".coverage.utbench.json")
	aliasSet := map[string]struct{}{}
	for _, alias := range aliases {
		name := strings.TrimSpace(alias)
		if name == "" {
			continue
		}
		base := filepath.Base(name)
		stem := strings.TrimSuffix(base, filepath.Ext(base))
		aliasSet[base] = struct{}{}
		aliasSet[stem] = struct{}{}
	}

	runCtx, cancelRun := context.WithTimeout(context.Background(), normalizedTimeout(timeoutSeconds))
	defer cancelRun()
	_, _ = runCommandWithProcessGroupKill(runCtx, py, []string{"-m", "coverage", "run", "--branch", "-m", "pytest", filename, "-q", "--maxfail=9999"}, workdir, nil)
	if runCtx.Err() != nil {
		return 0, 0, fmt.Sprintf("coverage run timed out after %ds", timeoutSeconds)
	}

	jsonCtx, cancelJSON := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelJSON()
	if out, err := runCommandWithProcessGroupKill(jsonCtx, py, []string{"-m", "coverage", "json", "-o", jsonPath}, workdir, nil); err != nil {
		return 0, 0, "coverage json failed: " + trimErr(string(out), 800)
	}

	raw, err := os.ReadFile(jsonPath)
	if err != nil {
		return 0, 0, err.Error()
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return 0, 0, err.Error()
	}

	files, _ := payload["files"].(map[string]any)
	if len(files) == 0 {
		return 0, 0, "coverage files empty"
	}

	for filePath, anyDetail := range files {
		base := filepath.Base(filePath)
		stem := strings.TrimSuffix(base, filepath.Ext(base))
		if _, ok := aliasSet[base]; !ok && base != sourceBase && stem != sourceStem {
			if _, ok2 := aliasSet[stem]; !ok2 {
				continue
			}
		}
		detail, _ := anyDetail.(map[string]any)
		summary, _ := detail["summary"].(map[string]any)
		line := extractLineCoverage(summary)
		branch := extractBranchCoverage(summary)
		return line, branch, ""
	}

	if totals, ok := payload["totals"].(map[string]any); ok {
		line := extractLineCoverage(totals)
		branch := extractBranchCoverage(totals)
		if line > 0 || branch > 0 {
			return line, branch, ""
		}
	}

	return 0, 0, "source file not found in coverage report"
}

func rewriteGeneratedTestImports(source string, sourcePath string) string {
	sourceStem := strings.TrimSuffix(filepath.Base(sourcePath), filepath.Ext(sourcePath))
	lines := strings.Split(source, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "from solution import") ||
			strings.HasPrefix(trimmed, "from your_module import") ||
			strings.HasPrefix(trimmed, "from module_name import") ||
			strings.HasPrefix(trimmed, "from src import") ||
			strings.HasPrefix(trimmed, "from module_under_test import") ||
			strings.HasPrefix(trimmed, "from target_module import") {
			lines[i] = strings.Replace(line, strings.Fields(trimmed)[1], sourceStem, 1)
		} else if strings.HasPrefix(trimmed, "import module_under_test") ||
			strings.HasPrefix(trimmed, "import target_module") {
			lines[i] = strings.Replace(line, strings.Fields(trimmed)[1], sourceStem, 1)
		}
	}
	return strings.Join(lines, "\n")
}

func inferAliasModules(sourceStem string) []string {
	base := []string{sourceStem, "module_under_test", "target_module", "solution", "your_module", "module_name", "src"}
	uniq := map[string]struct{}{}
	out := make([]string, 0, len(base))
	for _, item := range base {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := uniq[item]; ok {
			continue
		}
		uniq[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func isValidModuleName(v string) bool {
	if v == "" {
		return false
	}
	for i, ch := range v {
		if i == 0 {
			if !(ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')) {
				return false
			}
			continue
		}
		if !(ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9')) {
			return false
		}
	}
	return true
}

func normalizedPytestFilename(origin string) string {
	stem := strings.TrimSuffix(origin, filepath.Ext(origin))
	var b strings.Builder
	for _, ch := range stem {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '_' {
			b.WriteRune(ch)
		} else {
			b.WriteRune('_')
		}
	}
	out := b.String()
	if out == "" {
		out = "generated"
	}
	if out[0] >= '0' && out[0] <= '9' {
		out = "_" + out
	}
	if !strings.HasPrefix(out, "test_") {
		out = "test_" + out
	}
	return out + ".py"
}

func estimateAssertionDensity(path, language string) (int, int, float64) {
	switch strings.ToLower(language) {
	case "go":
		raw, err := os.ReadFile(path)
		if err != nil {
			return 0, 0, 0
		}
		return estimateGoAssertionDensity(string(raw))
	case "python":
		raw, err := os.ReadFile(path)
		if err != nil {
			return 0, 0, 0
		}
		return estimatePythonAssertionDensity(string(raw))
	case "java":
		raw, err := os.ReadFile(path)
		if err != nil {
			return 0, 0, 0
		}
		return estimateJavaAssertionDensity(string(raw))
	case "cpp":
		return estimateCppAssertionDensity(path)
	}
	return 0, 0, 0
}

func estimatePythonAssertionDensity(text string) (int, int, float64) {
	// 统计各类断言（概念上都是验证点）
	assertCount := 0

	// 1. 基础 assert 关键字
	assertCount += strings.Count(text, "assert ")

	// 2. pytest 异常/警告检查（上下文管理器也是验证点）
	assertCount += strings.Count(text, "pytest.raises")
	assertCount += strings.Count(text, "pytest.warns")

	// 3. Mock 断言（验证调用行为）
	// 统计 assert_called 模式（覆盖 assert_called, assert_called_once, assert_called_with 等）
	assertCount += strings.Count(text, "assert_called")
	assertCount += strings.Count(text, "assert_not_called")

	// 4. unittest.TestCase 断言方法
	assertCount += strings.Count(text, "assertEqual")
	assertCount += strings.Count(text, "assertNotEqual")
	assertCount += strings.Count(text, "assertTrue")
	assertCount += strings.Count(text, "assertFalse")
	assertCount += strings.Count(text, "assertIs")
	assertCount += strings.Count(text, "assertIsNot")
	assertCount += strings.Count(text, "assertIsNone")
	assertCount += strings.Count(text, "assertIsNotNone")
	assertCount += strings.Count(text, "assertIn")
	assertCount += strings.Count(text, "assertNotIn")
	assertCount += strings.Count(text, "assertRaises")

	// 统计测试用例数（pytest 风格：def test_xxx）
	testCount := strings.Count(text, "def test_")

	// 统计 unittest 风格测试方法（以 test 开头的方法）
	// 但要排除 pytest 的 def test_
 unittestPattern := regexp.MustCompile(`(?m)^\s+def test_\w+\s*\(`)
 unittestMatches := unittestPattern.FindAllString(text, -1)
 unittestCount := len(unittestMatches)

	// 如果有 unittest 风格的测试方法，也计入
	// 注意：unittest 方法通常缩进在类内部，所以单独统计
	testCount += unittestCount

	if testCount <= 0 {
		return assertCount, 0, 0
	}
	return assertCount, testCount, round(float64(assertCount)/float64(testCount), 6)
}

func estimateGoAssertionDensity(text string) (int, int, float64) {
	// 统计各类断言（概念上都是验证点）
	assertCount := 0

	// 1. testify 断言库（最常用）
	// assert.Equal, assert.NotNil, require.Equal 等
	assertCount += strings.Count(text, "assert.")
	assertCount += strings.Count(text, "require.")

	// 2. gomega 断言库
	assertCount += strings.Count(text, "Expect(")
	assertCount += strings.Count(text, "ExpectWithOffset(")
	assertCount += strings.Count(text, "Eventually(")
	assertCount += strings.Count(text, "Consistently(")
	assertCount += strings.Count(text, "Ω(")      // Omega 别名
	assertCount += strings.Count(text, "Should(") // gomega 的 Should

	// 3. Go 原生 testing 包的失败标记
	// t.Error/t.Errorf - 标记失败但继续执行
	// t.Fatal/t.Fatalf - 标记失败并终止
	// 注意：这些是"验证点"，表示测试发现了问题
	assertCount += strings.Count(text, "t.Error(")
	assertCount += strings.Count(text, "t.Errorf(")
	assertCount += strings.Count(text, "t.Fatal(")
	assertCount += strings.Count(text, "t.Fatalf(")

	// 4. check 断言库（较少用）
	assertCount += strings.Count(text, "check.")

	// 注意：不再统计 "if "，因为它是控制流，不一定是断言
	// 正确的做法是统计 t.Error/t.Fatal 等失败标记

	// 统计测试用例
	testCount := strings.Count(text, "func Test")

	// 统计子测试：t.Run("name", func(t *testing.T) { ... })
	// 表格驱动测试中，每个 t.Run 是一个独立测试
	assertCount += strings.Count(text, "t.Run(")

	// 统计 Example 测试（可验证输出）
	testCount += strings.Count(text, "func Example")

	if testCount <= 0 {
		return assertCount, 0, 0
	}
	return assertCount, testCount, round(float64(assertCount)/float64(testCount), 6)
}

func estimateJavaAssertionDensity(text string) (int, int, float64) {
	// 统计各类断言（概念上都是验证点）
	assertCount := 0

	// 1. JUnit 5 Assertions.* methods
	assertCount += strings.Count(text, "Assertions.assertEquals")
	assertCount += strings.Count(text, "Assertions.assertTrue")
	assertCount += strings.Count(text, "Assertions.assertFalse")
	assertCount += strings.Count(text, "Assertions.assertNull")
	assertCount += strings.Count(text, "Assertions.assertNotNull")
	assertCount += strings.Count(text, "Assertions.assertThrows")
	assertCount += strings.Count(text, "Assertions.assertThat")
	assertCount += strings.Count(text, "Assertions.assertSame")
	assertCount += strings.Count(text, "Assertions.assertNotSame")
	assertCount += strings.Count(text, "Assertions.assertArrayEquals")
	assertCount += strings.Count(text, "Assertions.assertLinesMatch")
	assertCount += strings.Count(text, "Assertions.assertTimeout")
	assertCount += strings.Count(text, "Assertions.assertTimeoutPreemptively")
	assertCount += strings.Count(text, "Assertions.assertIterableEquals")
	assertCount += strings.Count(text, "Assertions.assertNotEquals")
	assertCount += strings.Count(text, "Assertions.assertDoesNotThrow")
	assertCount += strings.Count(text, "Assertions.fail")

	// 2. JUnit 4 style (static import, without Assertions prefix)
	assertCount += strings.Count(text, "assertEquals(")
	assertCount += strings.Count(text, "assertTrue(")
	assertCount += strings.Count(text, "assertFalse(")
	assertCount += strings.Count(text, "assertNull(")
	assertCount += strings.Count(text, "assertNotNull(")
	assertCount += strings.Count(text, "assertSame(")
	assertCount += strings.Count(text, "assertNotSame(")
	assertCount += strings.Count(text, "assertThrows(")
	assertCount += strings.Count(text, "assertThat(")
	assertCount += strings.Count(text, "assertArrayEquals(")
	assertCount += strings.Count(text, "assertDoesNotThrow(")
	assertCount += strings.Count(text, "expect(")
	assertCount += strings.Count(text, "fail(")

	// 3. Mockito 验证（验证调用行为）
	// verify(mock).method() 是验证点，确认 mock 被正确调用
	assertCount += strings.Count(text, "verify(")
	assertCount += strings.Count(text, "verifyNoMoreInteractions")
	assertCount += strings.Count(text, "verifyZeroInteractions")
	assertCount += strings.Count(text, "verifyNoInteractions")
	assertCount += strings.Count(text, "Mockito.verify")
	assertCount += strings.Count(text, "InOrder.verify")

	// 4. AssertJ 流式断言（现代 Java 测试常用）
	// assertThat(actual).isEqualTo(expected)
	assertCount += strings.Count(text, "assertThat(")
	assertCount += strings.Count(text, "Assertions.assertThat(") // 已在上面统计，但 AssertJ 也用这个

	// 5. Hamcrest matchers（虽然 assertThat 已统计，但 matcher 本身也是验证概念）
	// 注意：matcher 通常在 assertThat 内部，所以不重复统计

	// 统计测试用例：@Test 注解
	testCount := strings.Count(text, "@Test")

	// 注意：不统计 @Before/@After 等，它们不是测试方法

	if testCount <= 0 {
		return assertCount, 0, 0
	}
	return assertCount, testCount, round(float64(assertCount)/float64(testCount), 6)
}

func parsePytestCounts(output string) (*int, *int) {
	// 解析 pytest 摘要行，如 "2 passed, 1 failed, 1 skipped, 1 xfailed"
	// 注意：passed 和 failed 是实际执行并产生结果的测试
	// skipped/xfailed/xpassed/error 是特殊状态，不计入通过率分母
	passed := extractFirstInt(output, `(\d+)\s+passed`)
	failed := extractFirstInt(output, `(\d+)\s+failed`)
	skipped := extractFirstInt(output, `(\d+)\s+skipped`)
	xfailed := extractFirstInt(output, `(\d+)\s+xfailed`)
	xpassed := extractFirstInt(output, `(\d+)\s+xpassed`)
	errors := extractFirstInt(output, `(\d+)\s+error`)
	_ = skipped // 用于判断是否有特殊状态
	_ = xfailed
	_ = xpassed
	_ = errors

	// 如果有明确的 passed 或 failed 数字，优先使用
	if passed != nil || failed != nil {
		// 只统计真正执行的测试（passed + failed）
		// skipped/xfailed 等不计入分母，因为它们没有实际验证行为
		passCount := 0
		if passed != nil {
			passCount = *passed
		}
		totalCount := passCount
		if failed != nil {
			totalCount = passCount + *failed
		}
		if totalCount > 0 {
			return &passCount, &totalCount
		}
		return nil, nil
	}

	// 备用：解析进度条 [100%] 行
	// pytest 输出进度时，每个字符代表一个测试状态
	// . = passed, F = failed, E = error, s = skipped, x = xfailed, X = xpassed
	lines := strings.Split(output, "\n")
	var shortLine string
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		// pytest 7.0+ 使用不同的进度显示格式
		if strings.Contains(line, "[100%]") || strings.Contains(line, "passed") || strings.Contains(line, "failed") {
			shortLine = line
			break
		}
	}
	if shortLine == "" {
		return nil, nil
	}

	// 从进度条字符统计
	passedCount := 0
	failedCount := 0
	for _, ch := range shortLine {
		switch ch {
		case '.': // passed
			passedCount++
		case 'F', 'E', '!': // failed/error
			failedCount++
		// 's' = skipped, 'x' = xfailed, 'X' = xpassed - 不计入通过/失败分母
		}
	}

	if passedCount == 0 && failedCount == 0 {
		// 如果进度条没有字符，尝试从摘要行推断
		// 检查是否有特殊状态的测试但没有 passed/failed
		if skipped != nil || xfailed != nil || xpassed != nil || errors != nil {
			// 有特殊状态但没有 passed/failed，返回 nil
			return nil, nil
		}
		return nil, nil
	}

	total := passedCount + failedCount
	if total > 0 {
		return &passedCount, &total
	}
	return nil, nil
}

func extractFirstInt(text, pattern string) *int {
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(text)
	if len(match) < 2 {
		return nil
	}
	value := 0
	for _, ch := range match[1] {
		if ch < '0' || ch > '9' {
			return nil
		}
		value = value*10 + int(ch-'0')
	}
	return &value
}

func extractLineCoverage(summary map[string]any) float64 {
	covered, cOk := toFloat(summary["covered_lines"])
	total, tOk := toFloat(summary["num_statements"])
	if cOk && tOk && total > 0 {
		return round(covered/total, 6)
	}
	if pct, ok := toFloat(summary["percent_covered"]); ok {
		return round(pct/100.0, 6)
	}
	return 0
}

func extractBranchCoverage(summary map[string]any) float64 {
	covered, cOk := toFloat(summary["covered_branches"])
	total, tOk := toFloat(summary["num_branches"])
	if cOk && tOk && total > 0 {
		return round(covered/total, 6)
	}
	if pct, ok := toFloat(summary["percent_covered_branches"]); ok {
		return round(pct/100.0, 6)
	}
	return 0
}

func toFloat(v any) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

func trimErr(v string, max int) string {
	v = regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(v), " ")
	if len(v) <= max {
		return v
	}
	return v[:max]
}

func pythonExecutable() string {
	if path, err := exec.LookPath("/opt/venv/bin/python"); err == nil {
		return path
	}
	if path, err := exec.LookPath("/opt/venv/bin/python3"); err == nil {
		return path
	}
	if _, err := exec.LookPath("python"); err == nil {
		return "python"
	}
	return "python3"
}

func normalizedTimeout(seconds int) time.Duration {
	if seconds <= 0 {
		seconds = defaultTestTimeoutSeconds
	}
	return time.Duration(seconds) * time.Second
}

func round(v float64, digits int) float64 {
	p := 1.0
	for i := 0; i < digits; i++ {
		p *= 10
	}
	if v >= 0 {
		return float64(int(v*p+0.5)) / p
	}
	return float64(int(v*p-0.5)) / p
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func getLanguagesSummary(cases []contracts.GeneratedCase) string {
	langs := make(map[string]int)
	for _, c := range cases {
		langs[c.Language]++
	}
	var parts []string
	for _, l := range []string{"python", "go", "java", "cpp"} {
		if langs[l] > 0 {
			parts = append(parts, fmt.Sprintf("%s:%d", l, langs[l]))
		}
	}
	return strings.Join(parts, ", ")
}

func getRuntimeMS(ms *int) int {
	if ms == nil {
		return 0
	}
	return *ms
}

func getCoverageValue(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}
