package evaluator

import (
	"context"
	"errors"
	"path/filepath"

	"go-ut-bench/internal/contracts"
)

func init() {
	RegisterLanguageEvaluator("cpp", &CppEvaluator{})
}

// CppEvaluator C++ 语言评测器
type CppEvaluator struct{}

func (e *CppEvaluator) PrepareWorkspace(item contracts.GeneratedCase) (*WorkspaceContext, error) {
	workdir, testFileName, sourceBase, sourceStem, prepErr := prepareCppWorkspace(item.GeneratedTestPath, item.SamplePath)
	if workdir == "" {
		return nil, errors.New(prepErr)
	}
	return &WorkspaceContext{
		Workdir:       workdir,
		TestPath:      testFileName,
		SourceBase:    sourceBase,
		SourceStem:    sourceStem,
		ShouldCleanup: true,
	}, nil
}

func (e *CppEvaluator) CompileCheck(ws *WorkspaceContext) (bool, string) {
	return cppCompileCheck(ws.Workdir)
}

func (e *CppEvaluator) ExecuteTests(ws *WorkspaceContext, timeoutSeconds int) (bool, string, int) {
	return executeCppTests(ws.Workdir)
}

func (e *CppEvaluator) ParseTestCounts(testOutput string) (*int, *int) {
	return parseCppTestCounts(testOutput)
}

func (e *CppEvaluator) EstimateAssertionDensity(workdir, testPath string) (int, int, float64) {
	// C++ 的 estimateCppAssertionDensity 接受完整路径
	return estimateCppAssertionDensity(filepath.Join(workdir, testPath))
}

func (e *CppEvaluator) CollectCoverage(ws *WorkspaceContext, testPassed bool, timeoutSeconds int) (float64, float64, string) {
	if ws.TestPath == "" {
		return 0, 0, ""
	}
	return collectCppCoverage(ws.Workdir, ws.TestPath)
}

func (e *CppEvaluator) CollectMutation(ctx context.Context, ws *WorkspaceContext, input MutationInput) (float64, mutationStats, string) {
	return collectCppMutation(ctx, ws.Workdir, ws.SourceBase, input.TimeoutSeconds, input.TestPassRate, input.TestPassed, input.TestTotal)
}

func (e *CppEvaluator) MutationTool() string {
	return "mull"
}
