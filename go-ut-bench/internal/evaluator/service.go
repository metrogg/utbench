// evaluator 鍖呮彁渚涘崟鍏冩祴璇曡瘎娴嬪姛鑳?// 璐熻矗缂栬瘧銆佽繍琛屾祴璇曘€佹敹闆嗚鐩栫巼銆佹墽琛屽彉寮傛祴璇曞苟鐢熸垚璇勬祴鎶ュ憡
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

type Service struct {
	logger *obs.Logger
}

type Output struct {
	Result     contracts.EvaluationResultSet
	ResultPath string
}

// evalTask is passed between evaluation workers.
type evalTask struct {
	index int
	item  contracts.GeneratedCase
}

type evalResultItem struct {
	index int
	row   contracts.EvaluationResult
}

// NewService 鍒涘缓鏂扮殑璇勬祴鏈嶅姟瀹炰緥
// 鍙傛暟:
//   - logger: 鏃ュ織璁板綍鍣ㄥ疄渚?//
//
// 杩斿洖鍊?
//   - *Service: 鏂扮殑鏈嶅姟瀹炰緥
func NewService(logger *obs.Logger) *Service {
	SetMutationLogger(logger)
	return &Service{logger: logger}
}

// Evaluate 鎵ц瀹屾暣鐨勮瘎娴嬫祦绋?// 鍙傛暟:
//   - ctx: 涓婁笅鏂囷紝鐢ㄤ簬鍙栨秷鎿嶄綔
//   - spec: 杩愯瑙勬牸璇存槑
//   - manifestPath: 鐢熸垚鐨勬祴璇曟竻鍗曟枃浠惰矾寰?//
//
// 杩斿洖鍊?
//   - Output: 璇勬祴缁撴灉杈撳嚭
//   - error: 璇勬祴澶辫触鏃剁殑閿欒
//
// 鍔熻兘璇存槑:
//  1. 璇诲彇鐢熸垚鐨勬祴璇曟竻鍗?//  2. 浣跨敤worker姹犲苟琛岃瘎娴嬫瘡涓牱鏈?//  3. 瀵规瘡涓牱鏈墽琛岋細缂栬瘧 -> 娴嬭瘯 -> 瑕嗙洊鐜?-> 鍙樺紓娴嬭瘯
//  4. 姹囨€荤粨鏋滃苟鍐欏叆JSON鏂囦欢
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

	// 璁＄畻worker鏁伴噺
	workerCount := spec.Workers
	if workerCount <= 0 {
		workerCount = min(8, max(2, runtime.NumCPU()))
	}

	// 杈撳嚭璇勬祴閰嶇疆淇℃伅
	total := len(manifest.Cases)
	progress := obs.NewProgressReporter(total, "evaluate")
	progress.PrintStageStart("璇勬祴娴嬭瘯", fmt.Sprintf("鏍锋湰: %d | 鍙樺紓: %v | Workers: %d",
		total, spec.MutationEnabled, workerCount))

	// 鍒涘缓杈撳嚭鐩綍
	runRoot := filepath.Join(spec.OutputRoot, "runs", spec.RunID)
	evalRoot := filepath.Join(runRoot, "evaluation")
	if err := os.MkdirAll(evalRoot, 0o755); err != nil {
		return Output{}, err
	}
	tasks := make(chan evalTask, workerCount*2)
	results := make(chan evalResultItem, workerCount*2)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Fprintf(os.Stderr, "[worker panic] %v\n", r)
				}
			}()
			for t := range tasks {
				row := s.evaluateOne(ctx, spec, t.item)
				select {
				case <-ctx.Done():
					return
				case results <- evalResultItem{index: t.index, row: row}:
				}
			}
		}()
	}

	go func() {
		defer close(tasks)
		for i, item := range manifest.Cases {
			tasks <- evalTask{index: i, item: item}
		}
	}()

	// 娑堣垂鑰咃細鏀堕泦缁撴灉
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect and display first-pass evaluation results.
	resultItems := make([]evalResultItem, 0, len(manifest.Cases))
	completed := 0
	for result := range results {
		completed++
		resultItems = append(resultItems, result)
		row := result.row

		taskResult := obs.TaskResult{
			Model:         row.Model,
			Language:      row.Language,
			SampleID:      row.SampleID,
			Success:       !evaluationFailed(row),
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

		status := evaluationStatus(row)
		progress.PrintTaskLine(completed, total, row.Model, row.Language, row.SampleID, status, fmt.Sprintf("%dms", getRuntimeMS(row.RuntimeMS)))

		if completed%5 == 0 {
			progress.PrintStats()
		}
	}

	progress.PrintStats()
	progress.PrintStageDone("评测测试", obs.StageStats{
		Total:    total,
		Success:  countSuccessfulResults(resultItems),
		Duration: time.Since(progress.GetStartTime()),
	})

	if err := ctx.Err(); err != nil {
		return Output{}, err
	}

	sort.Slice(resultItems, func(i, j int) bool { return resultItems[i].index < resultItems[j].index })
	rows := make([]contracts.EvaluationResult, 0, len(resultItems))
	for _, result := range resultItems {
		rows = append(rows, result.row)
	}

	// 鎺掑簭缁撴灉锛氭寜妯″瀷 -> 璇█ -> 鏍锋湰ID
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Model == rows[j].Model {
			if rows[i].Language == rows[j].Language {
				return rows[i].SampleID < rows[j].SampleID
			}
			return rows[i].Language < rows[j].Language
		}
		return rows[i].Model < rows[j].Model
	})

	set := contracts.EvaluationResultSet{
		SchemaVersion:  contracts.SchemaVersion,
		RunID:          spec.RunID,
		EvaluatedAtUTC: time.Now().UTC(),
		ManifestPath:   manifestPath,
		Results:        rows,
	}
	resultPath := filepath.Join(evalRoot, "evaluation_result.json")
	if err := contracts.WriteJSON(resultPath, set); err != nil {
		return Output{}, err
	}

	// 濡傛灉鍚敤鍙樺紓娴嬭瘯涓旂瓥鐣ヤ负fail锛屾鏌ユ槸鍚︽湁閿欒
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

func evaluationFailed(row contracts.EvaluationResult) bool {
	if !row.CompilePass {
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

func evaluationStatus(row contracts.EvaluationResult) string {
	if !row.CompilePass {
		return "FAIL(compile)"
	}
	if row.TestPass != nil && !*row.TestPass {
		return "FAIL(test)"
	}
	if row.TestPassRate != nil && *row.TestPassRate < 1.0 {
		return "FAIL(test)"
	}
	return "PASS"
}

func countSuccessfulResults(items []evalResultItem) int {
	count := 0
	for _, item := range items {
		if !evaluationFailed(item.row) {
			count++
		}
	}
	return count
}

func (s *Service) evaluateOne(ctx context.Context, spec contracts.RunSpec, item contracts.GeneratedCase) (result contracts.EvaluationResult) {
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
			defer cleanupWorkspace(workdir)
		}

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
				lineCov, branchCov, covErr := collectPythonCoverageInWorkspace(workdir, testName, packageName, targetFile, testTimeout)
				if covErr != "" && pass {
					row.CoverageError = covErr
				} else if covErr == "" {
					row.LineCoverage = &lineCov
					row.BranchCoverage = &branchCov
				}
			}
		} else if sourceBase != "" {
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
			checkResult := CheckTestPassRate(testPassed, testTotal, "mutmut", GetMinPassRateForTool("mutmut"))
			if !checkResult.ShouldRun {
				row.MutationError = checkResult.Message
				row.MutationTool = "mutmut"
			} else {
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
	} else if strings.EqualFold(item.Language, "go") {
		workdir, testName, sourceBase, _ := prepareGoWorkspace(item.GeneratedTestPath, item.SamplePath)
		if workdir == "" {
			row.CompilePass = false
			row.CompileError = testName // prepErr stored in testName
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}
		defer cleanupWorkspace(workdir)

		compilePass, compileErr := goCompileCheck(workdir, testName)
		row.CompilePass = compilePass
		if !compilePass {
			row.CompileError = compileErr
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}

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
			lineCov, branchCov, covErr := collectGoCoverage(workdir, testName, sourceBase)
			if covErr != "" && pass {
				row.CoverageError = covErr
			} else if covErr == "" {
				row.LineCoverage = &lineCov
				row.BranchCoverage = &branchCov
			}
		}

		if spec.MutationEnabled && !strings.EqualFold(strings.TrimSpace(spec.MutationPolicy), "skip") {
			mutationStart := time.Now()
			mutationScore, mutationStats, mutationErr := collectGoMutation(ctx, workdir, testName, sourceBase, spec.MutationTimeout, row.TestPassRate, 0, 0)
			mutationElapsed := time.Since(mutationStart)
			s.logger.ToFile("evaluator").Trace("mutation_result",
				"model", item.Model,
				"language", item.Language,
				"sample_id", item.SampleID,
				"tool", "gremlins",
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
			row.MutationTool = "gremlins"
		}

	} else if strings.EqualFold(item.Language, "java") {
		workdir, testName, _, className := prepareJavaWorkspace(item.GeneratedTestPath, item.SamplePath)
		if workdir == "" {
			row.CompilePass = false
			row.CompileError = testName // prepErr stored in testName
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}
		defer cleanupWorkspace(workdir)

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
			lineCov, branchCov, covErr := collectJavaCoverage(workdir, className)
			if covErr != "" && pass {
				row.CoverageError = covErr
			} else if covErr == "" {
				row.LineCoverage = &lineCov
				row.BranchCoverage = &branchCov
			}
		}

		if spec.MutationEnabled && !strings.EqualFold(strings.TrimSpace(spec.MutationPolicy), "skip") {
			mutationStart := time.Now()
			testPassed := 0
			if row.TestPassCount != nil {
				testPassed = *row.TestPassCount
			}
			testTotal := 0
			if row.TestTotalCount != nil {
				testTotal = *row.TestTotalCount
			}
			mutationScore, mutationStats, mutationErr := collectJavaMutation(ctx, workdir, className, spec.MutationTimeout, row.TestPassRate, testPassed, testTotal)
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
		workdir, testName, sourceBase, _, prepErr := prepareCppWorkspace(item.GeneratedTestPath, item.SamplePath)
		if workdir == "" {
			row.CompilePass = false
			row.CompileError = prepErr
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}
		defer cleanupWorkspace(workdir)

		compilePass, compileErr := cppCompileCheck(workdir)
		row.CompilePass = compilePass
		if !compilePass {
			row.CompileError = compileErr
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}

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
			lineCov, branchCov, covErr := collectCppCoverage(workdir, testName)
			if covErr != "" && pass {
				row.CoverageError = covErr
			} else if covErr == "" {
				row.LineCoverage = &lineCov
				row.BranchCoverage = &branchCov
			}
		}

		if spec.MutationEnabled && !strings.EqualFold(strings.TrimSpace(spec.MutationPolicy), "skip") {
			mutationStart := time.Now()
			testPassed := 0
			if row.TestPassCount != nil {
				testPassed = *row.TestPassCount
			}
			testTotal := 0
			if row.TestTotalCount != nil {
				testTotal = *row.TestTotalCount
			}
			mutationScore, mutationStats, mutationErr := collectCppMutation(ctx, workdir, sourceBase, spec.MutationTimeout, row.TestPassRate, testPassed, testTotal)
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
	if row.RuntimeMS == nil {
		rt := int(time.Since(start).Milliseconds())
		row.RuntimeMS = &rt
	}
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
	if row.TestPass == nil || row.TestPassCount != nil || row.TestTotalCount != nil {
		return
	}
	total := 1
	passed := 0
	if *row.TestPass {
		passed = 1
	}
	rateValue := float64(passed)
	row.TestPassCount = &passed
	row.TestTotalCount = &total
	row.TestPassRate = &rateValue
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
		if isEnvironmentFailureMessage(msg) {
			return "environment", shortFailureReason(msg)
		}
	}
	if generatedTestDidNotPass(row) {
		return "model", ""
	}
	for _, msg := range []string{row.CompileError, row.TestError, row.CoverageError, row.MutationError} {
		if msg == "" {
			continue
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
	return strings.Contains(msg, "permission denied") ||
		strings.Contains(msg, "access is denied") ||
		strings.Contains(msg, "executable file not found") ||
		strings.Contains(msg, "not installed") ||
		strings.Contains(msg, "command not found") ||
		// Only for compile errors with missing headers (.h/.cpp files)
		// NOT for test errors - those with "no such file or directory" are model issues (wrong mock strategy)
		(strings.Contains(msg, "no such file or directory") &&
			(strings.Contains(msg, ".h\"") || strings.Contains(msg, ".h>") ||
				strings.Contains(msg, ".cpp\"") || strings.Contains(msg, ".cpp>")))
}

func isToolFailureMessage(msg string) bool {
	msg = strings.ToLower(msg)
	if strings.Contains(msg, "all tests failed") ||
		strings.Contains(msg, "no tests found, skipping mutation") ||
		strings.Contains(msg, "pytest timed out") ||
		strings.Contains(msg, "coverage run timed out") {
		return false
	}
	return strings.Contains(msg, "coverage json failed") ||
		strings.Contains(msg, "coverage files empty") ||
		strings.Contains(msg, "stats file not found") ||
		strings.Contains(msg, "produced zero mutants") ||
		strings.Contains(msg, "did not execute any mutants") ||
		strings.Contains(msg, "run incomplete") ||
		strings.Contains(msg, "parse error")
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
	cmd := exec.Command(py, "-m", "py_compile", path)
	output, err := cmd.CombinedOutput()
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
	raw, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, 0
	}
	text := string(raw)
	if strings.EqualFold(language, "go") {
		return estimateGoAssertionDensity(text)
	}

	assertCount := strings.Count(text, "assert ") + strings.Count(text, "pytest.raises")
	testCount := strings.Count(text, "def test_")
	if testCount <= 0 {
		return assertCount, 0, 0
	}
	return assertCount, testCount, round(float64(assertCount)/float64(testCount), 6)
}

func estimateGoAssertionDensity(text string) (int, int, float64) {
	assertCount := strings.Count(text, "if ") +
		strings.Count(text, "assert.") +
		strings.Count(text, "Expect(") +
		strings.Count(text, "require.")
	testCount := strings.Count(text, "func Test")
	if testCount <= 0 {
		return assertCount, 0, 0
	}
	return assertCount, testCount, round(float64(assertCount)/float64(testCount), 6)
}

func parsePytestCounts(output string) (*int, *int) {
	passed := extractFirstInt(output, `(\d+)\s+passed`)
	failed := extractFirstInt(output, `(\d+)\s+failed`)
	if passed != nil || failed != nil {
		if passed != nil && failed != nil {
			total := *passed + *failed
			return passed, &total
		}
		if passed != nil {
			total := *passed
			return passed, &total
		}
		return nil, nil
	}

	lines := strings.Split(output, "\n")
	var shortLine string
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		if strings.Contains(line, "[100%]") {
			shortLine = line
			break
		}
	}
	if shortLine == "" {
		return nil, nil
	}

	passedCount := 0
	failedCount := 0
	for _, ch := range shortLine {
		switch ch {
		case '.', 's', 'S':
			passedCount++
		case 'F', 'E', 'x', 'X', '!':
			failedCount++
		}
	}
	if passedCount == 0 && failedCount == 0 {
		return nil, nil
	}
	total := passedCount + failedCount
	if failedCount == 0 {
		return &passedCount, &total
	}
	if passedCount == 0 {
		return nil, &total
	}
	return &passedCount, &total
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
