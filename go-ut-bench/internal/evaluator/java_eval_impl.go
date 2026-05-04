package evaluator

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"go-ut-bench/internal/contracts"
)

func init() {
	RegisterLanguageEvaluator("java", &JavaEvaluator{})
}

// JavaEvaluator Java 语言评测器
type JavaEvaluator struct{}

func (e *JavaEvaluator) PrepareWorkspace(item contracts.GeneratedCase) (*WorkspaceContext, error) {
	workdir, testFileName, _, className := prepareJavaWorkspace(item.GeneratedTestPath, item.SamplePath)
	if workdir == "" {
		// prepareJavaWorkspace 在失败时将错误信息存在 testFileName 中
		return nil, errors.New(testFileName)
	}
	return &WorkspaceContext{
		Workdir:       workdir,
		TestPath:      testFileName,
		Extra:         map[string]string{"className": className},
		ShouldCleanup: true,
	}, nil
}

func (e *JavaEvaluator) CompileCheck(ws *WorkspaceContext) (bool, string) {
	return javaCompileCheck(ws.Workdir)
}

func (e *JavaEvaluator) ExecuteTests(ws *WorkspaceContext, timeoutSeconds int) (bool, string, int) {
	return executeJavaTestsWithTimeout(ws.Workdir, timeoutSeconds)
}

func (e *JavaEvaluator) ParseTestCounts(testOutput string) (*int, *int) {
	return parseJavaTestCounts(testOutput)
}

func (e *JavaEvaluator) EstimateAssertionDensity(workdir, testPath string) (int, int, float64) {
	fullPath := filepath.Join(workdir, "src", "test", "java", testPath)
	raw, err := os.ReadFile(fullPath)
	if err != nil {
		return 0, 0, 0
	}
	return estimateJavaAssertionDensity(string(raw))
}

func (e *JavaEvaluator) CollectCoverage(ws *WorkspaceContext, testPassed bool, timeoutSeconds int) (float64, float64, string) {
	className := ws.Extra["className"]
	if className == "" {
		return 0, 0, ""
	}
	return collectJavaCoverage(ws.Workdir, className)
}

func (e *JavaEvaluator) CollectMutation(ctx context.Context, ws *WorkspaceContext, input MutationInput) (float64, mutationStats, string) {
	className := ws.Extra["className"]
	return collectJavaMutation(ctx, ws.Workdir, className, input.TimeoutSeconds, input.TestPassRate, input.TestPassed, input.TestTotal)
}

func (e *JavaEvaluator) MutationTool() string {
	return "pitest"
}
