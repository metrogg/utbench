package evaluator

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func goCompileCheck(workdir, testFile string) (bool, string) {
	_ = testFile
	// Compile check for Go tests must use `go test -c`; `go build` rejects *_test.go files.
	tempOutput := filepath.Join(workdir, "compile_check_output.test")
	if runtime.GOOS == "windows" {
		tempOutput = filepath.Join(workdir, "compile_check_output.test.exe")
	}
	defer os.Remove(tempOutput)

	cmd := exec.Command("go", "test", "-c", "-o", tempOutput, ".")
	cmd.Dir = workdir
	output, err := cmd.CombinedOutput()
	if err == nil {
		return true, ""
	}
	return false, trimErr(string(output), 2000)
}

func prepareGoWorkspace(testPath, samplePath string) (string, string, string, string) {
	testSource, err := os.ReadFile(testPath)
	if err != nil {
		return "", "", "", fmt.Sprintf("failed to read generated test: %s", err)
	}

	sourceBase := filepath.Base(samplePath)
	sourceStem := strings.TrimSuffix(sourceBase, filepath.Ext(sourceBase))
	testFileName := sourceStem + "_test.go"

	workdir, err := os.MkdirTemp("", "utbench_go_eval_")
	if err != nil {
		return "", "", "", fmt.Sprintf("failed to create temp dir: %s", err)
	}

	sourceData, err := os.ReadFile(samplePath)
	if err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", fmt.Sprintf("failed to read source: %s", err)
	}

	if err := os.WriteFile(filepath.Join(workdir, sourceBase), sourceData, 0o644); err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", fmt.Sprintf("failed to write source: %s", err)
	}

	goModContent := "module utbench_eval\n\ngo 1.22\n"
	if err := os.WriteFile(filepath.Join(workdir, "go.mod"), []byte(goModContent), 0o644); err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", fmt.Sprintf("failed to write go.mod: %s", err)
	}

	if err := os.WriteFile(filepath.Join(workdir, testFileName), testSource, 0o644); err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", fmt.Sprintf("failed to write test: %s", err)
	}

	return workdir, testFileName, sourceBase, sourceStem
}

func executeGoTests(workdir, testFile, sourceFile string) (bool, string, int) {
	cmd := exec.Command("go", "test", "-v", filepath.Base(testFile), filepath.Base(sourceFile))
	cmd.Dir = workdir
	started := time.Now()
	output, err := cmd.CombinedOutput()
	latency := int(time.Since(started).Milliseconds())
	if err == nil {
		return true, string(output), latency
	}
	return false, trimErr(string(output), 4000), latency
}

func parseGoTestCounts(output string) (*int, *int) {
	passed := 0
	failed := 0
	for _, match := range regexp.MustCompile(`--- (PASS|FAIL):`).FindAllStringSubmatch(output, -1) {
		if len(match) >= 2 {
			if match[1] == "PASS" {
				passed++
			} else if match[1] == "FAIL" {
				failed++
			}
		}
	}
	if passed > 0 || failed > 0 {
		total := passed + failed
		return &passed, &total
	}
	if strings.Contains(output, "PASS") && !strings.Contains(output, "FAIL") {
		passed = 1
		total := 1
		return &passed, &total
	}
	return nil, nil
}

func collectGoCoverage(workdir, testFile, sourceBase string) (float64, float64, string) {
	coverFile := filepath.Join(workdir, "cover.out")
	cmd := exec.Command("go", "test", "-coverprofile="+filepath.Base(coverFile), filepath.Base(testFile), filepath.Base(sourceBase))
	cmd.Dir = workdir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, 0, trimErr(string(output), 2000)
	}

	raw, err := os.ReadFile(coverFile)
	if err != nil {
		return 0, 0, fmt.Sprintf("failed to read coverage file: %s", err)
	}

	return parseGoCoverageOutput(string(raw), sourceBase)
}

func parseGoCoverageOutput(content, sourceBase string) (float64, float64, string) {
	lines := strings.Split(content, "\n")
	if len(lines) < 2 {
		return 0, 0, "coverage file too short"
	}

	modeLine := lines[0]
	if !strings.HasPrefix(modeLine, "mode:") {
		return 0, 0, "invalid coverage format"
	}

	totalStmts := 0
	coveredStmts := 0
	totalBranches := 0
	coveredBranches := 0

	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		if len(parts) < 3 {
			continue
		}

		filePathWithRange := parts[0]
		colonIdx := strings.Index(filePathWithRange, ":")
		if colonIdx == -1 {
			continue
		}
		filePath := filePathWithRange[:colonIdx]
		if !strings.HasSuffix(filePath, sourceBase) {
			continue
		}

		countStr := parts[len(parts)-1]
		var count int
		if countStr == "1" {
			count = 1
		} else {
			count = 0
		}

		rangeStr := parts[1]
		stmtCount := estimateGoStmtCount(rangeStr)
		totalStmts += stmtCount
		if count > 0 {
			coveredStmts += stmtCount
		}

		totalBranches += 2
		if count > 0 {
			coveredBranches += 2
		}
	}

	if totalStmts == 0 {
		return 0, 0, "no coverage data for target file"
	}

	lineCov := round(float64(coveredStmts)/float64(totalStmts), 6)
	branchCov := round(float64(coveredBranches)/float64(totalBranches), 6)
	return lineCov, branchCov, ""
}

func estimateGoStmtCount(rangeStr string) int {
	count := 1
	for _, ch := range rangeStr {
		if ch == ',' || ch == '.' {
			count++
		}
	}
	if count > 1 {
		count = count / 2
	}
	if count < 1 {
		count = 1
	}
	return count
}

func collectGoMutation(ctx context.Context, workdir, testFile, sourceBase string, timeoutSeconds int, testPassRate *float64, testPassed, testTotal int) (float64, mutationStats, string) {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 120
	}

	minPassRate := GetMinPassRateForTool("go-mutesting")
	passed := 0
	total := 0
	if testPassed > 0 || testTotal > 0 {
		passed = testPassed
		total = testTotal
	} else if testPassRate != nil {
		total = 1
		passed = int(*testPassRate * float64(total))
		if passed == 0 && *testPassRate > 0 {
			passed = 1
		}
	}

	checkResult := CheckTestPassRate(passed, total, "go-mutesting", minPassRate)
	if !checkResult.ShouldRun {
		return 0, mutationStats{}, checkResult.Message
	}

	targetPath := filepath.Join(workdir, sourceBase)
	if _, err := os.Stat(targetPath); err != nil {
		return 0, mutationStats{}, fmt.Sprintf("target file not found: %s", sourceBase)
	}

	mutestingPath := findGoMutesting()
	if mutestingPath == "" {
		return 0, mutationStats{}, "go-mutesting not installed. Install: go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@latest"
	}

	runCtx, cancelRun := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancelRun()

	runOut, runErr := runCommandWithProcessGroupKill(runCtx, mutestingPath, []string{"./..."}, workdir, nil)

	stats, parseErr := parseGoMutestingOutput(string(runOut))
	if parseErr != "" {
		return 0, stats, formatMutationError("go-mutesting parse error", runErr, runOut, nil, nil)
	}

	if stats.Total <= 0 {
		return 0, stats, formatMutationError("go-mutesting produced zero mutants", runErr, runOut, nil, nil)
	}

	processed := stats.Killed + stats.Survived + stats.NoTests + stats.Timeout + stats.Skipped + stats.Suspicious
	if processed <= 0 {
		return 0, stats, formatMutationError("go-mutesting did not execute any mutants", runErr, runOut, nil, nil)
	}

	if stats.Killed+stats.Survived <= 0 {
		return 0, stats, formatMutationError("go-mutesting no killed/survived results", runErr, runOut, nil, nil)
	}

	// 更科学的得分计算：测试能杀死的变异 / 所有生成的变异
	// 包含NoTests部分，能真实反映测试覆盖度
	effectiveTotal := stats.Killed + stats.Survived + stats.NoTests
	if effectiveTotal == 0 {
		return 0, stats, "no effective mutants found"
	}
	score := round(float64(stats.Killed)/float64(effectiveTotal), 6)
	return score, stats, ""
}

func findGoMutesting() string {
	candidates := []string{
		"go-mutesting",
		filepath.Join(os.Getenv("HOME"), "go", "bin", "go-mutesting"),
		"/usr/local/go/bin/go-mutesting",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
		if path, err := exec.LookPath("go-mutesting"); err == nil {
			return path
		}
	}
	return ""
}

func parseGoMutestingOutput(output string) (mutationStats, string) {
	stats := mutationStats{}

	stats.Total = extractFirstIntOrZero(output, `total is (\d+)`)
	passed := extractFirstIntOrZero(output, `(\d+) passed`)
	failed := extractFirstIntOrZero(output, `(\d+) failed`)
	skipped := extractFirstIntOrZero(output, `(\d+) skipped`)
	duplicated := extractFirstIntOrZero(output, `(\d+) duplicated`)

	stats.Killed = passed
	stats.Survived = failed
	stats.Skipped = skipped + duplicated

	if stats.Total == 0 {
		return stats, "no mutation stats found in output"
	}
	return stats, ""
}

func extractFirstIntOrZero(s, pattern string) int {
	if match := extractFirstInt(s, pattern); match != nil {
		return *match
	}
	return 0
}

func inferGoMutationTargets(workdir, testFile, sourceBase string) []string {
	testContent, err := os.ReadFile(filepath.Join(workdir, testFile))
	if err != nil {
		return []string{sourceBase}
	}

	testText := string(testContent)
	if strings.Contains(testText, filepath.Base(strings.TrimSuffix(sourceBase, ".go"))) {
		return []string{strings.TrimSuffix(sourceBase, ".go")}
	}
	return []string{sourceBase}
}
