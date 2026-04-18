package evaluator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func goCompileCheck(workdir, testFile string) (bool, string) {
	cmd := exec.Command("go", "build", "-o", "/dev/null", testFile)
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

func executeGoTests(workdir, testFile string) (bool, string, int) {
	cmd := exec.Command("go", "test", "-v", filepath.Base(testFile))
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
	return nil, nil
}

func collectGoCoverage(workdir, testFile, sourceBase string) (float64, float64, string) {
	coverFile := filepath.Join(workdir, "cover.out")
	cmd := exec.Command("go", "test", "-coverprofile="+filepath.Base(coverFile), filepath.Base(testFile))
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

		filePath := parts[0]
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