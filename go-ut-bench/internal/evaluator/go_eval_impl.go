package evaluator

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"go-ut-bench/internal/contracts"
)

func init() {
	RegisterLanguageEvaluator("go", &GoEvaluator{})
}

// GoEvaluator Go 语言评测器
type GoEvaluator struct{}

func (e *GoEvaluator) PrepareWorkspace(item contracts.GeneratedCase) (*WorkspaceContext, error) {
	workdir, testFile, sourceBase, prepErr := prepareGoWorkspace(item.GeneratedTestPath, item.SamplePath)
	if workdir == "" {
		return nil, errors.New(prepErr)
	}
	return &WorkspaceContext{
		Workdir:       workdir,
		TestPath:      testFile,
		SourceBase:    sourceBase,
		ShouldCleanup: true,
	}, nil
}

func (e *GoEvaluator) CompileCheck(ws *WorkspaceContext) (bool, string) {
	return goCompileCheck(ws.Workdir, ws.TestPath)
}

func (e *GoEvaluator) ExecuteTests(ws *WorkspaceContext, timeoutSeconds int) (bool, string, int) {
	return executeGoTests(ws.Workdir, ws.TestPath, ws.SourceBase)
}

func (e *GoEvaluator) ParseTestCounts(testOutput string) (*int, *int) {
	return parseGoTestCounts(testOutput)
}

func (e *GoEvaluator) EstimateAssertionDensity(workdir, testPath string) (int, int, float64) {
	raw, err := os.ReadFile(filepath.Join(workdir, testPath))
	if err != nil {
		return 0, 0, 0
	}
	return estimateGoAssertionDensity(string(raw))
}

func (e *GoEvaluator) CollectCoverage(ws *WorkspaceContext, testPassed bool, timeoutSeconds int) (float64, float64, string) {
	if ws.SourceBase == "" {
		return 0, 0, ""
	}
	return collectGoCoverage(ws.Workdir, ws.TestPath, ws.SourceBase)
}

func (e *GoEvaluator) CollectMutation(ctx context.Context, ws *WorkspaceContext, input MutationInput) (float64, mutationStats, string) {
	return collectGoMutation(ctx, ws.Workdir, ws.TestPath, ws.SourceBase, input.TimeoutSeconds, input.TestPassRate, input.TestPassed, input.TestTotal)
}

func (e *GoEvaluator) MutationTool() string {
	return "go-mutesting"
}
