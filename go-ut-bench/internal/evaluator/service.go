// evaluator 包提供单元测试评测功能
// 负责编译、运行测试、收集覆盖率、执行变异测试并生成评测报告
package evaluator

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
	"go-ut-bench/internal/store"
)

// Service 评测服务结构
// 提供完整的评测流程管理
type Service struct {
	logger *obs.Logger // 日志记录器
}

// cleanupSemaphore 清理工作目录的并发限制信号量
// 限制同时清理的目录数量为2，避免系统资源占用过高
var cleanupSemaphore = make(chan struct{}, 2)

const evaluatorVersion = "utbench-evaluator.v1"
const defaultScorePolicyVersion = "default-v2"

// Output 评测操作的输出结果
// 包含评测结果集和结果文件路径
type Output struct {
	Result     contracts.EvaluationResultSet
	ResultPath string
}

// evalTask 评测任务结构
// 用于 worker 之间传递任务
type evalTask struct {
	index   int                     // 任务序号
	item    contracts.GeneratedCase // 待评测的生成结果
	reused  *store.ReusableEvaluationResult
	envHash string
}

// evalResultItem 评测结果项
// 包含序号和评测结果
type evalResultItem struct {
	index int                        // 任务序号
	row   contracts.EvaluationResult // 评测结果
}

type evaluationReusePlan struct {
	ByTaskKey    map[string]store.ReusableEvaluationResult
	ReusableHits int
}

// NewService 创建新的评测服务实例
//
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
//
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

	// 规范化路径分隔符：manifest 可能在 Windows 上生成（反斜杠），
	// 但评测可能在 Linux 容器内执行（需要正斜杠）。
	// 注意：filepath.ToSlash 在 Linux 上是 no-op（不会转换反斜杠），
	// 必须用 strings.ReplaceAll 确保跨平台一致。
	for i := range manifest.Cases {
		c := &manifest.Cases[i]
		c.SamplePath = toSlashCrossPlatform(c.SamplePath)
		c.PromptPath = toSlashCrossPlatform(c.PromptPath)
		c.GeneratedTestPath = toSlashCrossPlatform(c.GeneratedTestPath)
		c.ResponsePath = toSlashCrossPlatform(c.ResponsePath)
		c.MetadataPath = toSlashCrossPlatform(c.MetadataPath)
		c.TracePath = toSlashCrossPlatform(c.TracePath)
		c.WorkspaceDiffPath = toSlashCrossPlatform(c.WorkspaceDiffPath)
	}

	// 采集评测环境指纹
	envFingerprint := CaptureEnvironmentFingerprint(ctx, false, "")
	envFingerprintHash := envFingerprint.FingerprintHash()
	s.logger.Debug("environment fingerprint captured", "fingerprint", envFingerprintHash)

	var reuseStore *store.SQLiteStore
	if spec.ReuseEvaluation && strings.TrimSpace(spec.DBPath) != "" {
		if db, openErr := store.OpenSQLite(spec.DBPath); openErr == nil {
			if initErr := db.Init(ctx); initErr == nil {
				reuseStore = db
				defer reuseStore.Close()
			} else {
				_ = db.Close()
				s.logger.Warn("reuse evaluation disabled", "reason", initErr.Error())
			}
		} else {
			s.logger.Warn("reuse evaluation disabled", "reason", openErr.Error())
		}
	}
	reusePlan := prepareEvaluationReusePlan(ctx, spec, manifest.Cases, envFingerprintHash, reuseStore)

	// 计算worker数量
	workerCount := spec.Workers
	if workerCount <= 0 {
		workerCount = min(8, max(2, runtime.NumCPU()))
	}

	// 输出评测配置信息
	total := len(manifest.Cases)
	progress := obs.NewProgressReporterWithWriter(total, "evaluate", s.logger.Writer())
	stageHeader := fmt.Sprintf("样本: %d | 变异: %v | Workers: %d", total, spec.MutationEnabled, workerCount)
	if reusePlan.ReusableHits > 0 {
		stageHeader += fmt.Sprintf(" | 预判可复用: %d", reusePlan.ReusableHits)
	}
	progress.PrintStageStart("评测测试", stageHeader)

	// 鍒涘缓杈撳嚭鐩綍
	runRoot := filepath.Join(spec.OutputRoot, "runs", spec.RunID)
	evalRoot := filepath.Join(runRoot, "evaluation")
	if err := os.MkdirAll(evalRoot, 0o755); err != nil {
		return Output{}, err
	}
	tasks := make(chan evalTask, workerCount*2)
	results := make(chan evalResultItem, workerCount*2)
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
				if r := recover(); r != nil {
					fmt.Fprintf(os.Stderr, "[worker panic] %v\n", r)
				}
			}()
			for t := range tasks {
				key, setPhase, done := tracker.start(t.item)
				_ = key
				res := s.evaluateOne(ctx, spec, t.item, t.envHash, t.reused, setPhase)
				done()
				select {
				case <-ctx.Done():
					return
				case results <- evalResultItem{index: t.index, row: res}:
				}
			}
		}()
	}

	go func() {
		defer close(tasks)
		for i, item := range manifest.Cases {
			taskKey := evaluationTaskKey(item)
			var reused *store.ReusableEvaluationResult
			if row, ok := reusePlan.ByTaskKey[taskKey]; ok {
				copyRow := row
				reused = &copyRow
			}
			tasks <- evalTask{index: i, item: item, reused: reused, envHash: envFingerprintHash}
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

	// 鎺掑簭缁撴灉锛氭寜妯″瀷 -> 璇█ -> 样本ID
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
		SchemaVersion:          contracts.SchemaVersion,
		RunID:                  spec.RunID,
		EvaluatedAtUTC:         time.Now().UTC(),
		ManifestPath:           manifestPath,
		Results:                rows,
		EnvironmentFingerprint: envFingerprintHash,
		EnvironmentJSON:        envFingerprint.ToJSON(),
	}
	resultPath := filepath.Join(evalRoot, "evaluation_result.json")
	if err := contracts.WriteJSON(resultPath, set); err != nil {
		return Output{}, err
	}

	// 濡傛灉鍚敤变异娴嬭瘯涓旂瓥鐣ヤ负fail锛屾鏌ユ槸鍚︽湁閿欒
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

func (s *Service) evaluateOne(ctx context.Context, spec contracts.RunSpec, item contracts.GeneratedCase, evaluationEnvFingerprint string, reused *store.ReusableEvaluationResult, setPhase func(string)) (result contracts.EvaluationResult) {
	start := time.Now()
	row := contracts.EvaluationResult{
		Model:                    item.Model,
		SubjectID:                item.SubjectID,
		SubjectKind:              item.SubjectKind,
		AgentFramework:           item.AgentFramework,
		AgentModel:               item.AgentModel,
		SkillName:                item.SkillName,
		SkillVersion:             item.SkillVersion,
		Language:                 item.Language,
		SampleID:                 item.SampleID,
		SampleUID:                item.SampleUID,
		GeneratedTestPath:        item.GeneratedTestPath,
		SourcePath:               item.SamplePath,
		PromptTokens:             item.PromptTokens,
		CompletionTokens:         item.CompletionTokens,
		TotalTokens:              item.TotalTokens,
		TokenSource:              item.TokenSource,
		EstimatedCostUSD:         item.EstimatedCostUSD,
		CostSource:               item.CostSource,
		Truncated:                item.Truncated,
		TracePath:                item.TracePath,
		WorkspaceDiffPath:        item.WorkspaceDiffPath,
		SandboxProvider:          item.SandboxProvider,
		SandboxFingerprint:       item.SandboxFingerprint,
		EvaluationEnvFingerprint: evaluationEnvFingerprint,
		EvaluatorVersion:         evaluatorVersion,
		MutationConfigSHA256:     mutationConfigSHA256(spec),
	}
	row.EvaluationKey = evaluationKeyForItem(spec, item, row.MutationConfigSHA256, evaluationEnvFingerprint)
	if item.LatencyMS > 0 {
		row.LatencyMS = &item.LatencyMS
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

	if reused != nil {
		setPhase("evaluation.reuse_lookup")
		applyReusableEvaluation(&row, *reused)
		return row
	}

	s.logger.Debug("evaluating", "model", item.Model, "lang", item.Language, "sample", item.SampleID)

	langEval := getLanguageEvaluator(item.Language)
	if langEval != nil {
		s.evalWithLanguageEvaluator(ctx, spec, item, &row, start, setPhase, langEval)
	} else {
		// 未知语言：设置默认值
		pass := true
		lineCov := 0.0
		branchCov := 0.0
		mutation := 0.0
		assertCnt, testCnt, density := estimateAssertionDensity(item.GeneratedTestPath, item.Language)

		row.CompilePass = true
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

// evaluationFailed 检查评测是否失败
// 编译失败或测试失败都视为失败
//
// 参数:
//   - row: 评测结果
//
// 返回值:
//   - bool: 是否失败
func evaluationFailed(row contracts.EvaluationResult) bool {
	return !row.CompilePass || (row.TestPass != nil && !*row.TestPass)
}

// evaluationStatus 获取评测状态字符串
// 用于进度显示
//
// 参数:
//   - row: 评测结果
//
// 返回值:
//   - string: 状态标识（PASS、COMPILE_ERR、TEST_FAIL）
func evaluationStatus(row contracts.EvaluationResult) string {
	if !row.CompilePass {
		return "COMPILE_ERR"
	}
	if row.TestPass != nil && !*row.TestPass {
		return "TEST_FAIL"
	}
	return "PASS"
}

// countSuccessfulResults 统计成功的结果数量
// 编译通过且测试通过的视为成功
//
// 参数:
//   - items: 结果列表
//
// 返回值:
//   - int: 成功数量
func countSuccessfulResults(items []evalResultItem) int {
	count := 0
	for _, item := range items {
		if item.row.CompilePass && (item.row.TestPass == nil || *item.row.TestPass) {
			count++
		}
	}
	return count
}

func applyReusableEvaluation(row *contracts.EvaluationResult, reused store.ReusableEvaluationResult) {
	current := *row
	reusedRow := reused.Result
	reusedRow.Model = current.Model
	reusedRow.SubjectID = current.SubjectID
	reusedRow.SubjectKind = current.SubjectKind
	reusedRow.AgentFramework = current.AgentFramework
	reusedRow.AgentModel = current.AgentModel
	reusedRow.SkillName = current.SkillName
	reusedRow.SkillVersion = current.SkillVersion
	reusedRow.Language = current.Language
	reusedRow.SampleID = current.SampleID
	reusedRow.SampleUID = current.SampleUID
	reusedRow.GeneratedTestPath = current.GeneratedTestPath
	reusedRow.SourcePath = current.SourcePath
	reusedRow.PromptTokens = current.PromptTokens
	reusedRow.CompletionTokens = current.CompletionTokens
	reusedRow.TotalTokens = current.TotalTokens
	reusedRow.TokenSource = current.TokenSource
	reusedRow.EstimatedCostUSD = current.EstimatedCostUSD
	reusedRow.CostSource = current.CostSource
	reusedRow.LatencyMS = current.LatencyMS
	reusedRow.EvaluationEnvFingerprint = current.EvaluationEnvFingerprint
	reusedRow.EvaluationKey = current.EvaluationKey
	reusedRow.EvaluatorVersion = current.EvaluatorVersion
	reusedRow.MutationConfigSHA256 = current.MutationConfigSHA256
	reusedRow.Reused = true
	reusedRow.ReuseStage = "evaluation"
	reusedRow.ReuseKey = current.EvaluationKey
	reusedRow.ReuseReason = "evaluation_key_match"
	reusedRow.ReusedFromRunID = reused.RunID
	reusedRow.ReusedFromResultID = reused.EvaluationResultID
	*row = reusedRow
}

// finalizeEvaluationResult 完成评测结果的最终处理
// 设置总耗时、失败来源分类、评分资格等
//
// 参数:
//   - row: 评测结果指针
//   - start: 开始时间
func finalizeEvaluationResult(row *contracts.EvaluationResult, start time.Time) {
	if !(row.Reused && row.RuntimeMS != nil) {
		totalRuntimeMS := int(time.Since(start).Milliseconds())
		row.RuntimeMS = &totalRuntimeMS
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

// shouldRunMutationAfterSampleTests 检查是否应该在样本测试后运行变异测试
// 基线测试失败时跳过变异测试。
// 优先使用 test_pass_rate（解析自测试框架输出）判断，因为进程退出码
// 可能被覆盖率工具（如 gcov）等非测试因素干扰。
//
// 参数:
//   - row: 评测结果
//
// 返回值:
//   - bool: 是否应该运行变异测试
func shouldRunMutationAfterSampleTests(row contracts.EvaluationResult) bool {
	// 优先使用 test_pass_rate：如果解析到 100% 通过率，说明所有测试用例都通过了
	if row.TestPassRate != nil && *row.TestPassRate >= 1.0 {
		return true
	}
	// 有通过率但未达 100%，跳过变异
	if row.TestPassRate != nil && *row.TestPassRate < 1.0 {
		return false
	}
	// 没有 pass_rate 数据时，回退到退出码判断
	if row.TestPass != nil && !*row.TestPass {
		return false
	}
	return true
}

// pythonTestImportsAnyMutationTarget 检查生成的测试是否导入了变异目标模块
// 用于判断变异测试是否有效
//
// 参数:
//   - workdir: 工作目录
//   - testName: 测试文件名
//   - mutationTargets: 变异目标列表
//
// 返回值:
//   - bool: 是否导入任意目标
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

// classifyFailureOrigin 分类失败来源
// 区分 dataset、model、environment、tool 等来源
//
// 参数:
//   - row: 评测结果
//
// 返回值:
//   - string: 失败来源（none、dataset、model、environment、tool）
//   - string: 简短失败原因
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

// generatedTestDidNotPass 检查生成的测试是否未通过
// 编译失败、测试失败或测试通过率低于100%都视为未通过
//
// 参数:
//   - row: 评测结果
//
// 返回值:
//   - bool: 是否未通过
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

// isDatasetFailureMessage 判断消息是否为数据集相关失败
// 检查文件缺失、元数据缺失等数据集问题
//
// 参数:
//   - msg: 错误消息
//
// 返回值:
//   - bool: 是否为数据集失败
func isDatasetFailureMessage(msg string) bool {
	msg = strings.ToLower(msg)
	return strings.Contains(msg, "dataset root") ||
		strings.Contains(msg, "source file not found") ||
		strings.Contains(msg, "target file not found") ||
		strings.Contains(msg, "repo_level sample missing metadata") ||
		strings.Contains(msg, "repo_level workspace not found") ||
		strings.Contains(msg, "repo_level workspace_root not set") ||
		strings.Contains(msg, "failed to read source")
}

// isEnvironmentFailureMessage 判断消息是否为环境相关失败
// 检查权限、工具未安装等环境问题
//
// 参数:
//   - msg: 错误消息
//
// 返回值:
//   - bool: 是否为环境失败
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

// isToolFailureMessage 判断消息是否为工具相关失败
// 检查覆盖率工具、变异测试工具的问题
//
// 参数:
//   - msg: 错误消息
//
// 返回值:
//   - bool: 是否为工具失败
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

// preparePythonWorkspace 准备 Python 评测工作区
// 创建临时目录，复制源码和测试文件，处理导入别名
//
// 参数:
//   - testPath: 生成的测试文件路径
//   - sourcePath: 源码文件路径
//
// 返回值:
//   - string: 工作目录路径（失败时为空）
//   - string: 测试文件名（失败时为错误信息）
//   - string: 源码文件名
//   - string: 源码文件名（去掉扩展名）
//   - string: 错误信息（成功时为空）
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

// cleanupWorkspace 清理工作目录
//
// 参数:
//   - workdir: 工作目录路径
func cleanupWorkspace(workdir string) {
	_ = os.RemoveAll(workdir)
}

// cleanupWorkspaceAsync 异步清理工作目录
// 使用信号量限制并发清理数量
//
// 参数:
//   - workdir: 工作目录路径
//   - model: 模型名称
//   - language: 编程语言
//   - sampleID: 样本 ID
//   - logger: 日志记录器
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

// isRepoLevelSample 判断样本是否为仓库级别样本
// 检查是否存在 meta.json 或 {name}.meta.json 文件
//
// 参数:
//   - samplePath: 样本文件路径
//
// 返回值:
//   - bool: 是否为仓库级别样本
func isRepoLevelSample(samplePath string) bool {
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
		var meta contracts.RepoLevelMeta
		if raw, err := os.ReadFile(metaPath); err == nil {
			if err := json.Unmarshal(raw, &meta); err == nil && meta.ModuleImport != "" {
				return true
			}
		}
	}
	return false
}

func loadRepoLevelMeta(samplePath string) *contracts.RepoLevelMeta {
	sampleDir := filepath.Dir(samplePath)
	entryBase := filepath.Base(samplePath)
	entryExt := filepath.Ext(entryBase)
	entryName := entryBase[:len(entryBase)-len(entryExt)]
	if entryName == "entry" {
		metaPath := filepath.Join(sampleDir, "meta.json")
		if raw, err := os.ReadFile(metaPath); err == nil {
			var meta contracts.RepoLevelMeta
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
	var meta contracts.RepoLevelMeta
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

func preparePythonRepoLevelWorkspace(testPath string, samplePath string) (string, string, string, string, string) {
	meta := loadRepoLevelMeta(samplePath)
	if meta == nil {
		return "", "", "", "", "repo_level sample missing metadata"
	}
	workspaceRoot := meta.WorkspaceRoot
	if workspaceRoot == "" {
		return "", "", "", "", "repo_level workspace_root not set in metadata"
	}
	if _, err := os.Stat(workspaceRoot); err != nil {
		return "", "", "", "", "repo_level workspace not found: " + workspaceRoot
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
		return 0, 0, "missing package name for repo_level coverage"
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
	runCtx, cancel := context.WithTimeout(context.Background(), PythonCompileTimeoutSeconds*time.Second)
	defer cancel()
	output, err := runCommandWithProcessGroupKill(runCtx, py, []string{"-m", "py_compile", path}, "", nil)
	if runCtx.Err() != nil {
		return false, fmt.Sprintf("python compile timed out after %ds", PythonCompileTimeoutSeconds)
	}
	if err != nil {
		msg := strings.TrimSpace(string(output))
		if msg == "" {
			msg = err.Error()
		}
		return false, msg
	}
	if !strings.Contains(text, "def test_") && !strings.Contains(text, "import pytest") && !strings.Contains(text, "unittest.TestCase") {
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
	// 保留头尾各一半，中间用省略号连接，确保错误信息（通常在末尾）不被丢弃
	half := max/2 - 20
	if half < 100 {
		half = 100
	}
	return v[:half] + " ... [truncated] ... " + v[len(v)-half:]
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
	for _, l := range contracts.SupportedLanguages {
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

func evaluationKeyForItem(spec contracts.RunSpec, item contracts.GeneratedCase, mutationConfigSHA, evaluationEnvFingerprint string) string {
	generatedTestSHA := hashExistingFile(item.GeneratedTestPath)
	sampleUID := item.SampleUID
	if sampleUID == "" {
		sampleUID = strings.Join([]string{item.Language, item.SampleID, item.SamplePath}, "\x00")
	}
	payload := strings.Join([]string{
		generatedTestSHA,
		sampleUID,
		item.Language,
		evaluationEnvFingerprint,
		evaluatorVersion,
		defaultScorePolicyVersion,
		mutationConfigSHA,
		fmt.Sprintf("test-timeout=%d", spec.TestTimeout),
		item.DependencyFingerprint,
	}, "\x00")
	return "evaluation_" + sha256String(payload)[:16]
}

func evaluationTaskKey(item contracts.GeneratedCase) string {
	return strings.Join([]string{item.Model, item.Language, item.SampleID}, "\x00")
}

func prepareEvaluationReusePlan(ctx context.Context, spec contracts.RunSpec, items []contracts.GeneratedCase, evaluationEnvFingerprint string, reuseStore *store.SQLiteStore) evaluationReusePlan {
	out := evaluationReusePlan{ByTaskKey: make(map[string]store.ReusableEvaluationResult)}
	if reuseStore == nil || !spec.ReuseEvaluation {
		return out
	}
	mutationSHA := mutationConfigSHA256(spec)
	for _, item := range items {
		evaluationKey := evaluationKeyForItem(spec, item, mutationSHA, evaluationEnvFingerprint)
		reused, ok, err := reuseStore.FindReusableEvaluationAsset(ctx, evaluationKey)
		if err != nil {
			continue
		}
		if !ok {
			continue
		}
		out.ByTaskKey[evaluationTaskKey(item)] = reused
		out.ReusableHits++
	}
	return out
}

func mutationConfigSHA256(spec contracts.RunSpec) string {
	raw, _ := json.Marshal(map[string]any{
		"mutation_enabled": spec.MutationEnabled,
		"mutation_timeout": spec.MutationTimeout,
		"mutation_policy":  spec.MutationPolicy,
	})
	return sha256String(string(raw))
}

func hashExistingFile(path string) string {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func sha256String(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

// toSlashCrossPlatform 将路径中的反斜杠统一替换为正斜杠。
// filepath.ToSlash 在 Linux 上是 no-op（不会转换 Windows 反斜杠），
// 因此需要显式替换以确保跨平台一致性。
func toSlashCrossPlatform(path string) string {
	return strings.ReplaceAll(path, `\`, "/")
}
