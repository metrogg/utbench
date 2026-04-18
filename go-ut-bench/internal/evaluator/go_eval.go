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