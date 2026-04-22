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

	// 输出评测配置信息
	total := len(manifest.Cases)
	fmt.Fprintf(os.Stderr, "\n[Evaluator] Starting evaluation of %d samples\n", total)
	fmt.Fprintf(os.Stderr, "[Evaluator] Languages: %s | Mutation: %v\n",
		getLanguagesSummary(manifest.Cases), spec.MutationEnabled)
	fmt.Fprintf(os.Stderr, "[Evaluator] Workers: %d\n\n", min(16, max(2, runtime.NumCPU())))

	// 创建输出目录
	runRoot := filepath.Join(spec.OutputRoot, "runs", spec.RunID)
	evalRoot := filepath.Join(runRoot, "evaluation")
	if err := os.MkdirAll(evalRoot, 0o755); err != nil {
		return Output{}, err
	}

	// 创建worker池处理评测任务
	workerCount := min(16, max(2, runtime.NumCPU()))
	tasks := make(chan evalTask, workerCount*2)
	results := make(chan contracts.EvaluationResult, workerCount*2)

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
				item := s.evaluateOne(ctx, spec, t.item)
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
		status := "PASS"
		if !row.CompilePass {
			status = "FAIL(compile)"
		} else if row.TestPass != nil && !*row.TestPass {
			status = "FAIL(test)"
		}
		fmt.Fprintf(os.Stderr, "[%d/%d] %s | %s | %s | %s | %dms\n",
			completed, total, row.Model, row.Language, row.SampleID, status,
			getRuntimeMS(row.RuntimeMS))
		rows = append(rows, row)
	}

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

func (s *Service) evaluateOne(ctx context.Context, spec contracts.RunSpec, item contracts.GeneratedCase) contracts.EvaluationResult {
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

	if !item.Success {
		row.CompilePass = false
		if item.Error != nil {
			row.CompileError = "generation failed: " + item.Error.Message
		} else {
			row.CompileError = "generation failed"
		}
		rt := int(time.Since(start).Milliseconds())
		row.RuntimeMS = &rt
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

		var pass bool
		var testErr string
		var runtimeMs int
		if isModuleLevel {
			pass, testErr, runtimeMs = executePythonTestsInWorkspace(workdir, testName, packageName)
		} else {
			pass, testErr, runtimeMs = executePythonTests(workdir, testName)
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
				lineCov, branchCov, covErr := collectPythonCoverageInWorkspace(workdir, testName, packageName, targetFile)
				if covErr != "" && pass {
					row.CoverageError = covErr
				} else if covErr == "" {
					row.LineCoverage = &lineCov
					row.BranchCoverage = &branchCov
				}
			}
		} else if sourceBase != "" {
			lineCov, branchCov, covErr := collectPythonCoverage(workdir, testName, sourceBase, sourceStem, targets)
			if covErr != "" && pass {
				row.CoverageError = covErr
			} else if covErr == "" {
				row.LineCoverage = &lineCov
				row.BranchCoverage = &branchCov
			}
		}

		if spec.MutationEnabled {
			mutationTargets := targets
			if len(mutationTargets) == 0 && sourceBase != "" {
				mutationTargets = []string{sourceBase}
			}
			mutationStart := time.Now()
			fmt.Fprintf(os.Stderr, "  [mutation] %s | %s | %s | starting...\n", item.Model, item.Language, item.SampleID)
			mutationScore, mutationStats, mutationErr := collectPythonMutation(ctx, workdir, testName, mutationTargets, spec.MutationTimeout, testErr)
			mutationElapsed := int(time.Since(mutationStart).Seconds())
			if mutationErr != "" {
				fmt.Fprintf(os.Stderr, "  [mutation] %s | %s | %s | ERROR after %ds\n", item.Model, item.Language, item.SampleID, mutationElapsed)
				if strings.EqualFold(spec.MutationPolicy, "warn") {
					row.MutationError = mutationErr
				} else {
					row.MutationError = mutationErr
				}
			} else {
				fmt.Fprintf(os.Stderr, "  [mutation] %s | %s | %s | done in %ds, score=%.2f\n", item.Model, item.Language, item.SampleID, mutationElapsed, mutationScore)
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

		if spec.MutationEnabled && spec.MutationPolicy != "skip" {
			mutationStart := time.Now()
			fmt.Fprintf(os.Stderr, "  [mutation] %s | %s | %s | starting...\n", item.Model, item.Language, item.SampleID)
			mutationScore, mutationStats, mutationErr := collectGoMutation(ctx, workdir, testName, sourceBase, spec.MutationTimeout, row.TestPassRate, 0, 0)
			mutationElapsed := int(time.Since(mutationStart).Seconds())
			row.MutationScore = &mutationScore
			row.MutationTotal = &mutationStats.Total
			row.MutationKilled = &mutationStats.Killed
			row.MutationSurvived = &mutationStats.Survived
			row.MutationNoTests = &mutationStats.NoTests
			row.MutationTimeouts = &mutationStats.Timeout
			row.MutationSkipped = &mutationStats.Skipped
			row.MutationSuspicious = &mutationStats.Suspicious
			if mutationErr != "" {
				fmt.Fprintf(os.Stderr, "  [mutation] %s | %s | %s | ERROR after %ds\n", item.Model, item.Language, item.SampleID, mutationElapsed)
				row.MutationError = mutationErr
			} else {
				fmt.Fprintf(os.Stderr, "  [mutation] %s | %s | %s | done in %ds, score=%.2f\n", item.Model, item.Language, item.SampleID, mutationElapsed, mutationScore)
			}
			row.MutationTool = "go-mutesting"
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

		if spec.MutationEnabled && spec.MutationPolicy != "skip" {
			mutationStart := time.Now()
			fmt.Fprintf(os.Stderr, "  [mutation] %s | %s | %s | starting...\n", item.Model, item.Language, item.SampleID)
			mutationScore, mutationStats, mutationErr := collectJavaMutation(ctx, workdir, className, spec.MutationTimeout, row.TestPassRate, 0, 0)
			mutationElapsed := int(time.Since(mutationStart).Seconds())
			row.MutationScore = &mutationScore
			row.MutationTotal = &mutationStats.Total
			row.MutationKilled = &mutationStats.Killed
			row.MutationSurvived = &mutationStats.Survived
			row.MutationNoTests = &mutationStats.NoTests
			row.MutationTimeouts = &mutationStats.Timeout
			row.MutationSkipped = &mutationStats.Skipped
			row.MutationSuspicious = &mutationStats.Suspicious
			if mutationErr != "" {
				fmt.Fprintf(os.Stderr, "  [mutation] %s | %s | %s | ERROR after %ds\n", item.Model, item.Language, item.SampleID, mutationElapsed)
				row.MutationError = mutationErr
			} else {
				fmt.Fprintf(os.Stderr, "  [mutation] %s | %s | %s | done in %ds, score=%.2f\n", item.Model, item.Language, item.SampleID, mutationElapsed, mutationScore)
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

		if spec.MutationEnabled && spec.MutationPolicy != "skip" {
			mutationStart := time.Now()
			fmt.Fprintf(os.Stderr, "  [mutation] %s | %s | %s | starting...\n", item.Model, item.Language, item.SampleID)
			mutationScore, mutationStats, mutationErr := collectCppMutation(ctx, workdir, sourceBase, spec.MutationTimeout, row.TestPassRate, 0, 0)
			mutationElapsed := int(time.Since(mutationStart).Seconds())
			row.MutationScore = &mutationScore
			row.MutationTotal = &mutationStats.Total
			row.MutationKilled = &mutationStats.Killed
			row.MutationSurvived = &mutationStats.Survived
			row.MutationNoTests = &mutationStats.NoTests
			row.MutationTimeouts = &mutationStats.Timeout
			row.MutationSkipped = &mutationStats.Skipped
			row.MutationSuspicious = &mutationStats.Suspicious
			if mutationErr != "" {
				fmt.Fprintf(os.Stderr, "  [mutation] %s | %s | %s | ERROR after %ds\n", item.Model, item.Language, item.SampleID, mutationElapsed)
				row.MutationError = mutationErr
			} else {
				fmt.Fprintf(os.Stderr, "  [mutation] %s | %s | %s | done in %ds, score=%.2f\n", item.Model, item.Language, item.SampleID, mutationElapsed, mutationScore)
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

	if row.RuntimeMS == nil {
		rt := int(time.Since(start).Milliseconds())
		row.RuntimeMS = &rt
	}
	return row
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

func executePythonTestsInWorkspace(workdir, testName, packageName string) (bool, string, int) {
	py := pythonExecutable()
	cmd := exec.Command(py, "-m", "pytest", testName, "-q", "--maxfail=9999")
	cmd.Dir = workdir
	env := os.Environ()
	env = append(env, "PYTHONPATH="+workdir)
	cmd.Env = env
	started := time.Now()
	output, err := cmd.CombinedOutput()
	latency := int(time.Since(started).Milliseconds())
	if err == nil {
		return true, string(output), latency
	}
	return false, trimErr(string(output), 4000), latency
}

func collectPythonCoverageInWorkspace(workdir, testName, packageName, targetFile string) (float64, float64, string) {
	if packageName == "" {
		return 0, 0, "missing package name for module_level coverage"
	}
	py := pythonExecutable()
	absWorkdir, err := filepath.Abs(workdir)
	if err != nil {
		return 0, 0, "failed to get absolute path: " + err.Error()
	}
	jsonPath := filepath.Join(absWorkdir, ".coverage.utbench.json")
	runCmd := exec.Command(py, "-m", "coverage", "run", "--branch", "--source", packageName, "-m", "pytest", testName, "-q", "--maxfail=9999")
	runCmd.Dir = absWorkdir
	env := os.Environ()
	env = append(env, "PYTHONPATH="+absWorkdir)
	env = append(env, "COVERAGE_FILE="+filepath.Join(absWorkdir, ".coverage.utbench"))
	runCmd.Env = env
	runOut, runErr := runCmd.CombinedOutput()
	if runErr != nil {
		return 0, 0, "coverage run failed: " + trimErr(string(runOut), 800)
	}
	jsonCmd := exec.Command(py, "-m", "coverage", "json", "-o", jsonPath)
	jsonCmd.Dir = absWorkdir
	jsonCmd.Env = append(env, "COVERAGE_FILE="+filepath.Join(absWorkdir, ".coverage.utbench"))
	if out, err := jsonCmd.CombinedOutput(); err != nil {
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

func executePythonTests(workdir, filename string) (bool, string, int) {
	py := pythonExecutable()
	cmd := exec.Command(py, "-m", "pytest", filename, "-q", "--maxfail=9999")
	cmd.Dir = workdir
	started := time.Now()
	output, err := cmd.CombinedOutput()
	latency := int(time.Since(started).Milliseconds())
	if err == nil {
		return true, string(output), latency
	}
	return false, trimErr(string(output), 4000), latency
}

func collectPythonCoverage(workdir, filename, sourceBase, sourceStem string, aliases []string) (float64, float64, string) {
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

	runCmd := exec.Command(py, "-m", "coverage", "run", "--branch", "-m", "pytest", filename, "-q", "--maxfail=9999")
	runCmd.Dir = workdir
	_, _ = runCmd.CombinedOutput()

	jsonCmd := exec.Command(py, "-m", "coverage", "json", "-o", jsonPath)
	jsonCmd.Dir = workdir
	if out, err := jsonCmd.CombinedOutput(); err != nil {
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
		if strings.HasPrefix(trimmed, "from solution import") || strings.HasPrefix(trimmed, "from your_module import") || strings.HasPrefix(trimmed, "from module_name import") || strings.HasPrefix(trimmed, "from src import") {
			lines[i] = strings.Replace(line, strings.Fields(trimmed)[1], sourceStem, 1)
		}
	}
	return strings.Join(lines, "\n")
}

func inferAliasModules(sourceStem string) []string {
	base := []string{sourceStem, "solution", "your_module", "module_name", "src"}
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
	if _, err := exec.LookPath("python"); err == nil {
		return "python"
	}
	return "python3"
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
