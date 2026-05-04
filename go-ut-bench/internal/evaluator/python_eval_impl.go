package evaluator

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"go-ut-bench/internal/contracts"
)

func init() {
	RegisterLanguageEvaluator("python", &PythonEvaluator{})
}

// PythonEvaluator Python 语言评测器
// 内部处理 repo_level 和 self_contained 两种模式的差异
type PythonEvaluator struct{}

func (e *PythonEvaluator) PrepareWorkspace(item contracts.GeneratedCase) (*WorkspaceContext, error) {
	isRepoLevel := isRepoLevelSample(item.SamplePath)
	extra := map[string]string{}

	var workdir, testName, sourceBase, sourceStem string
	var prepErr string

	if isRepoLevel {
		var packageName, targetFile string
		workdir, testName, packageName, targetFile, prepErr = preparePythonRepoLevelWorkspace(item.GeneratedTestPath, item.SamplePath)
		extra["packageName"] = packageName
		extra["targetFile"] = targetFile
		extra["isRepoLevel"] = "true"
	} else {
		workdir, testName, sourceBase, sourceStem, prepErr = preparePythonWorkspace(item.GeneratedTestPath, item.SamplePath)
	}

	if prepErr != "" {
		return nil, errors.New(prepErr)
	}

	return &WorkspaceContext{
		Workdir:       workdir,
		TestPath:      testName,
		SourceBase:    sourceBase,
		SourceStem:    sourceStem,
		Extra:         extra,
		ShouldCleanup: !isRepoLevel, // repo_level 复用 in-place，不清理
	}, nil
}

func (e *PythonEvaluator) CompileCheck(ws *WorkspaceContext) (bool, string) {
	return pythonCompileCheck(filepath.Join(ws.Workdir, ws.TestPath))
}

func (e *PythonEvaluator) ExecuteTests(ws *WorkspaceContext, timeoutSeconds int) (bool, string, int) {
	if ws.Extra["isRepoLevel"] == "true" {
		return executePythonTestsInWorkspace(ws.Workdir, ws.TestPath, ws.Extra["packageName"], timeoutSeconds)
	}
	return executePythonTests(ws.Workdir, ws.TestPath, timeoutSeconds)
}

func (e *PythonEvaluator) ParseTestCounts(testOutput string) (*int, *int) {
	return parsePytestCounts(testOutput)
}

func (e *PythonEvaluator) EstimateAssertionDensity(workdir, testPath string) (int, int, float64) {
	raw, err := os.ReadFile(filepath.Join(workdir, testPath))
	if err != nil {
		return 0, 0, 0
	}
	return estimatePythonAssertionDensity(string(raw))
}

func (e *PythonEvaluator) CollectCoverage(ws *WorkspaceContext, testPassed bool, timeoutSeconds int) (float64, float64, string) {
	if ws.Extra["isRepoLevel"] == "true" {
		packageName := ws.Extra["packageName"]
		if packageName == "" {
			return 0, 0, ""
		}
		return collectPythonCoverageInWorkspace(ws.Workdir, ws.TestPath, packageName, ws.Extra["targetFile"], timeoutSeconds)
	}
	if ws.SourceBase == "" {
		return 0, 0, ""
	}
	targets := inferPythonMutationTargets(ws.Workdir, ws.TestPath, ws.SourceBase)
	return collectPythonCoverage(ws.Workdir, ws.TestPath, ws.SourceBase, ws.SourceStem, targets, timeoutSeconds)
}

func (e *PythonEvaluator) CollectMutation(ctx context.Context, ws *WorkspaceContext, input MutationInput) (float64, mutationStats, string) {
	// 确定变异目标
	mutationTargets := e.resolveMutationTargets(ws)
	if len(mutationTargets) == 0 && ws.SourceBase != "" {
		mutationTargets = []string{ws.SourceBase}
	}

	// Python 特有检查：测试文件是否导入了变异目标
	if ws.Extra["isRepoLevel"] != "true" && !pythonTestImportsAnyMutationTarget(ws.Workdir, ws.TestPath, mutationTargets) {
		return 0, mutationStats{}, "mutmut: generated tests do not import mutation target, skipping mutation"
	}

	// Python 使用 testOutput（原始测试输出）而非 testPassRate
	return collectPythonMutation(ctx, ws.Workdir, ws.TestPath, mutationTargets, input.TimeoutSeconds, input.TestOutput)
}

func (e *PythonEvaluator) MutationTool() string {
	return "mutmut"
}

// resolveMutationTargets 解析 Python 变异测试的目标文件列表
func (e *PythonEvaluator) resolveMutationTargets(ws *WorkspaceContext) []string {
	if ws.Extra["isRepoLevel"] == "true" {
		targetFile := ws.Extra["targetFile"]
		if targetFile != "" {
			return []string{targetFile}
		}
		return nil
	}
	return inferPythonMutationTargets(ws.Workdir, ws.TestPath, ws.SourceBase)
}

// inferPythonMutationTargets 推断 Python 单文件模式的变异目标
func inferPythonMutationTargets(workdir, testName, sourceBase string) []string {
	return inferMutationTargets(workdir, testName, sourceBase)
}
