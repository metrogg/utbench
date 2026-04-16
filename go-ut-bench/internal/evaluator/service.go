package evaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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

type evalTask struct {
	item contracts.GeneratedCase
}

func NewService(logger *obs.Logger) *Service {
	return &Service{logger: logger}
}

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

	runRoot := filepath.Join(spec.OutputRoot, "runs", spec.RunID)
	evalRoot := filepath.Join(runRoot, "evaluation")
	if err := os.MkdirAll(evalRoot, 0o755); err != nil {
		return Output{}, err
	}

	workerCount := min(16, max(2, runtime.NumCPU()))
	tasks := make(chan evalTask)
	results := make(chan contracts.EvaluationResult)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
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

	go func() {
		defer close(tasks)
		for _, item := range manifest.Cases {
			select {
			case <-ctx.Done():
				return
			case tasks <- evalTask{item: item}:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	rows := make([]contracts.EvaluationResult, 0, len(manifest.Cases))
	for row := range results {
		rows = append(rows, row)
	}
	if err := ctx.Err(); err != nil {
		return Output{}, err
	}

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

	if strings.EqualFold(item.Language, "python") {
		workdir, testName, sourceBase, sourceStem, prepErr := preparePythonWorkspace(item.GeneratedTestPath, item.SamplePath)
		if prepErr != "" {
			row.CompilePass = false
			row.CompileError = prepErr
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}
		defer cleanupWorkspace(workdir)

		compilePass, compileErr := pythonCompileCheck(filepath.Join(workdir, testName))
		row.CompilePass = compilePass
		if !compilePass {
			row.CompileError = compileErr
			rt := int(time.Since(start).Milliseconds())
			row.RuntimeMS = &rt
			return row
		}

		pass, testErr, runtimeMs := executePythonTests(workdir, testName)
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

		targets := inferMutationTargets(workdir, testName, sourceBase)

		if pass || (row.TestPassRate != nil && *row.TestPassRate >= 0.70) {
			lineCov, branchCov, covErr := collectPythonCoverage(workdir, testName, sourceBase, sourceStem, targets)
			if covErr != "" {
				row.CoverageError = covErr
			} else {
				row.LineCoverage = &lineCov
				row.BranchCoverage = &branchCov
			}

			if spec.MutationEnabled {
				mutationScore, mutationStats, mutationErr := collectPythonMutation(ctx, workdir, testName, targets, spec.MutationTimeout, testErr)
				if mutationErr != "" {
					if strings.EqualFold(spec.MutationPolicy, "warn") {
						row.MutationError = mutationErr
					} else {
						row.MutationError = mutationErr
					}
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
			}
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
	if out, err := runCmd.CombinedOutput(); err != nil {
		return 0, 0, "coverage run failed: " + trimErr(string(out), 800)
	}

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
	assertCount := strings.Count(text, "if ")
	testCount := strings.Count(text, "func Test")
	if testCount <= 0 {
		return assertCount, 0, 0
	}
	fset := token.NewFileSet()
	_, _ = parser.ParseFile(fset, "generated_test.go", text, parser.ParseComments)
	_ = ast.File{}
	return assertCount, testCount, round(float64(assertCount)/float64(testCount), 6)
}

func parsePytestCounts(output string) (*int, *int) {
	passed := extractFirstInt(output, `(\d+)\s+passed`)
	failed := extractFirstInt(output, `(\d+)\s+failed`)
	if passed == nil && failed == nil {
		return nil, nil
	}
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
