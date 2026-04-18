package evaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

type mutationStats struct {
	Total      int
	Killed     int
	Survived   int
	NoTests    int
	NotChecked int
	Timeout    int
	Skipped    int
	Suspicious int
}

func collectPythonMutation(ctx context.Context, workdir, testName string, mutationTargets []string, timeoutSeconds int, testOutput string) (float64, mutationStats, string) {
	if len(mutationTargets) == 0 {
		return 0, mutationStats{}, "missing mutation targets"
	}
	if timeoutSeconds <= 0 {
		timeoutSeconds = 120
	}
	failingTests := collectFailingTestsByRerun(ctx, workdir, testName)
	if len(failingTests) == 0 {
		failingTests = collectFailingTestsFromPytestOutput(testOutput)
	}

	pyprojectPath := filepath.Join(workdir, "pyproject.toml")
	if err := os.WriteFile(pyprojectPath, []byte(buildMutmutPyproject(mutationTargets, testName, failingTests)), 0o644); err != nil {
		return 0, mutationStats{}, err.Error()
	}
	env, envErr := buildMutmutEnv(workdir)
	if envErr != nil {
		return 0, mutationStats{}, envErr.Error()
	}

	py := pythonExecutable()
	mutantsDir := filepath.Join(workdir, "mutants")
	_ = os.RemoveAll(mutantsDir)

	if err := preCreateMutantsDirectory(workdir, mutantsDir, mutationTargets, testName); err != nil {
		return 0, mutationStats{}, "failed to pre-create mutants directory: " + err.Error()
	}

	runCtx, cancelRun := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancelRun()
	runCmd := exec.CommandContext(runCtx, py, "-m", "mutmut", "run")
	runCmd.Dir = workdir
	runCmd.Env = env
	runOut, runErr := runCmd.CombinedOutput()

	exportCtx, cancelExport := context.WithTimeout(ctx, 30*time.Second)
	defer cancelExport()
	exportCmd := exec.CommandContext(exportCtx, py, "-m", "mutmut", "export-cicd-stats")
	exportCmd.Dir = workdir
	exportCmd.Env = env
	exportOut, exportErr := exportCmd.CombinedOutput()

	statsFile := filepath.Join(workdir, "mutants", "mutmut-cicd-stats.json")
	raw, err := os.ReadFile(statsFile)
	if err != nil {
		return 0, mutationStats{}, formatMutationError("mutmut stats file not found", runErr, runOut, exportErr, exportOut)
	}

	stats, parseErr := parseMutationStats(raw)
	if parseErr != "" {
		return 0, mutationStats{}, parseErr
	}
	if metaStats, ok := extractMetaMutationStats(workdir, mutationTargets); ok {
		stats = metaStats
	}
	if stats.Total <= 0 {
		return 0, stats, "mutmut produced zero mutants"
	}

	processed := stats.Killed + stats.Survived + stats.NoTests + stats.Timeout + stats.Skipped + stats.Suspicious
	if processed <= 0 {
		return 0, stats, formatMutationError("mutmut did not execute any mutants", runErr, runOut, exportErr, exportOut)
	}
	if stats.NotChecked > 0 {
		return 0, stats, formatMutationError(
			fmt.Sprintf("mutmut run incomplete (%d/%d)", processed, stats.Total),
			runErr,
			runOut,
			exportErr,
			exportOut,
		)
	}

	score := round(float64(stats.Killed)/float64(stats.Total), 6)
	return score, stats, ""
}

func preCreateMutantsDirectory(workdir, mutantsDir string, mutationTargets []string, testName string) error {
	absWorkdir, err := filepath.Abs(workdir)
	if err != nil {
		return fmt.Errorf("failed to get absolute workdir: %w", err)
	}
	absMutantsDir, err := filepath.Abs(mutantsDir)
	if err != nil {
		return fmt.Errorf("failed to get absolute mutantsDir: %w", err)
	}

	if err := os.MkdirAll(absMutantsDir, 0o755); err != nil {
		return err
	}

	if testName != "" {
		testsDestDir := filepath.Join(absMutantsDir, "tests")
		if err := os.MkdirAll(testsDestDir, 0o755); err != nil {
			return err
		}

		testFileName := filepath.Base(testName)
		srcPath := filepath.Join(absWorkdir, testName)
		dstPath := filepath.Join(testsDestDir, testFileName)
		data, err := os.ReadFile(srcPath)
		if err != nil {
			return fmt.Errorf("failed to read test file %s: %w", srcPath, err)
		}
		if err := os.WriteFile(dstPath, data, 0o644); err != nil {
			return fmt.Errorf("failed to write test file %s: %w", dstPath, err)
		}
	}

	validTargets := 0
	for _, target := range mutationTargets {
		absTargetPath := filepath.Join(absWorkdir, target)
		if _, err := os.Stat(absTargetPath); err != nil {
			continue
		}
		validTargets++

		absPackageDir := filepath.Dir(absTargetPath)
		if absPackageDir == absWorkdir {
			continue
		}

		relPackageDir, err := filepath.Rel(absWorkdir, absPackageDir)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", absPackageDir, err)
		}
		mutantsPackageDir := filepath.Join(absMutantsDir, relPackageDir)
		if err := copyPackageDependencies(absPackageDir, mutantsPackageDir, filepath.Base(target)); err != nil {
			return fmt.Errorf("failed to copy dependencies for %s: %w", target, err)
		}
	}

	if validTargets == 0 {
		return fmt.Errorf("no valid mutation targets found in %s", absWorkdir)
	}
	return nil
}

func copyPackageDependencies(srcPackageDir, dstPackageDir, targetFileName string) error {
	if err := os.MkdirAll(dstPackageDir, 0o755); err != nil {
		return err
	}

	entries, err := os.ReadDir(srcPackageDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		name := entry.Name()
		if name == targetFileName {
			continue
		}

		srcPath := filepath.Join(srcPackageDir, name)
		dstPath := filepath.Join(dstPackageDir, name)

		if entry.IsDir() {
			if name == "__pycache__" {
				continue
			}
			if err := copyPackageDependencies(srcPath, dstPath, ""); err != nil {
				return err
			}
			continue
		}

		if strings.HasSuffix(name, ".py") {
			data, err := os.ReadFile(srcPath)
			if err != nil {
				return err
			}
			if err := os.WriteFile(dstPath, data, 0o644); err != nil {
				return err
			}
		}
	}
	return nil
}

func buildMutmutPyproject(sourceNames []string, testName string, failingTests []string) string {
	sources, _ := json.Marshal(sourceNames)
	args := []string{"-q", "--tb=no", "--maxfail=9999"}
	if len(failingTests) > 0 {
		excludes := make([]string, 0, len(failingTests))
		for _, item := range failingTests {
			excludes = append(excludes, "not "+item)
		}
		args = append(args, "-k", strings.Join(excludes, " and "))
	}
	argsJSON, _ := json.Marshal(args)
	return "[tool.mutmut]\n" +
		"paths_to_mutate = " + string(sources) + "\n" +
		"pytest_add_cli_args = " + string(argsJSON) + "\n"
}

func buildMutmutEnv(workdir string) ([]string, error) {
	env := os.Environ()
	absWorkdir, err := filepath.Abs(workdir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}
	shimDir := filepath.Join(absWorkdir, ".mutmut_shim")
	mutantsDir := filepath.Join(absWorkdir, "mutants")
	if err := os.MkdirAll(shimDir, 0o755); err != nil {
		return nil, err
	}
	sitecustomize := filepath.Join(shimDir, "sitecustomize.py")
	shimCode := "import multiprocessing as _mp\n" +
		"import multiprocessing.context as _mpc\n" +
		"_orig = _mpc._default_context.set_start_method\n" +
		"def _safe(method, force=False):\n" +
		"    try:\n" +
		"        return _orig(method, force=force)\n" +
		"    except RuntimeError as e:\n" +
		"        if 'context has already been set' in str(e):\n" +
		"            return None\n" +
		"        raise\n" +
		"_mpc._default_context.set_start_method = _safe\n" +
		"_mp.set_start_method = _safe\n"
	if err := os.WriteFile(sitecustomize, []byte(shimCode), 0o644); err != nil {
		return nil, err
	}

	pyPath := ""
	for _, entry := range env {
		if strings.HasPrefix(entry, "PYTHONPATH=") {
			pyPath = strings.TrimPrefix(entry, "PYTHONPATH=")
			break
		}
	}
	merged := shimDir + string(os.PathListSeparator) + mutantsDir + string(os.PathListSeparator) + absWorkdir
	if strings.TrimSpace(pyPath) != "" {
		merged = merged + string(os.PathListSeparator) + pyPath
	}


	out := make([]string, 0, len(env)+1)
	set := false
	for _, entry := range env {
		if strings.HasPrefix(entry, "PYTHONPATH=") {
			out = append(out, "PYTHONPATH="+merged)
			set = true
			continue
		}
		out = append(out, entry)
	}
	if !set {
		out = append(out, "PYTHONPATH="+merged)
	}
	return out, nil
}

func parseMutationStats(raw []byte) (mutationStats, string) {
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return mutationStats{}, err.Error()
	}

	total, okTotal := toFloat(payload["total"])
	killed, okKilled := toFloat(payload["killed"])
	if !okTotal || !okKilled {
		return mutationStats{}, "invalid mutmut stats payload"
	}

	stats := mutationStats{
		Total:      int(total),
		Killed:     int(killed),
		Survived:   toIntDefault(payload["survived"]),
		NoTests:    toIntDefault(payload["no_tests"]),
		NotChecked: toIntDefault(payload["not_checked"]) + toIntDefault(payload["check_was_interrupted_by_user"]),
		Timeout:    toIntDefault(payload["timeout"]),
		Skipped:    toIntDefault(payload["skipped"]),
		Suspicious: toIntDefault(payload["suspicious"]) + toIntDefault(payload["segfault"]),
	}
	return stats, ""
}

func toIntDefault(v any) int {
	f, ok := toFloat(v)
	if !ok {
		return 0
	}
	return int(f)
}

func formatMutationError(prefix string, runErr error, runOut []byte, exportErr error, exportOut []byte) string {
	runMsg := ""
	if runErr != nil {
		runMsg = runErr.Error()
	}
	exportMsg := ""
	if exportErr != nil {
		exportMsg = exportErr.Error()
	}
	return fmt.Sprintf(
		"%s; mutmut run err=%q; mutmut export err=%q; run_out=%q; export_out=%q",
		prefix,
		runMsg,
		exportMsg,
		trimErr(string(runOut), 1200),
		trimErr(string(exportOut), 1200),
	)
}

func inferMutationTargets(workdir, testName, sourceBase string) []string {
	testPath := filepath.Join(workdir, testName)
	raw, err := os.ReadFile(testPath)
	if err != nil {
		if sourceBase == "" {
			return nil
		}
		return []string{sourceBase}
	}

	text := string(raw)
	re := regexp.MustCompile(`(?m)^\s*(?:from|import)\s+([a-zA-Z_][a-zA-Z0-9_]*)`)
	matches := re.FindAllStringSubmatch(text, -1)
	uniq := map[string]struct{}{}
	out := make([]string, 0)
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		name := strings.TrimSpace(match[1])
		if name == "pytest" || name == "unittest" {
			continue
		}
		fileName := name + ".py"
		if _, err := os.Stat(filepath.Join(workdir, fileName)); err == nil {
			if _, ok := uniq[fileName]; !ok {
				uniq[fileName] = struct{}{}
				out = append(out, fileName)
			}
		}
	}
	if len(out) > 0 {
		sort.Strings(out)
		return out
	}
	if sourceBase != "" {
		return []string{sourceBase}
	}
	return nil
}

func collectFailingTestsFromPytestOutput(output string) []string {
	lines := strings.Split(output, "\n")
	set := map[string]struct{}{}
	out := make([]string, 0)
	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if !strings.Contains(line, "::") || !strings.Contains(line, "FAILED") {
			continue
		}
		idx := strings.Index(line, "FAILED")
		if idx >= 0 {
			line = strings.TrimSpace(line[idx+len("FAILED"):])
		}
		line = strings.Split(line, "[")[0]
		parts := strings.Split(line, "::")
		if len(parts) < 2 {
			continue
		}
		selector := ""
		if len(parts) >= 3 {
			className := strings.TrimSpace(parts[len(parts)-2])
			testName := strings.TrimSpace(strings.Split(parts[len(parts)-1], " - ")[0])
			if className != "" && testName != "" {
				selector = fmt.Sprintf("(%s and %s)", className, testName)
			}
		} else {
			modulePart := strings.TrimSpace(parts[0])
			moduleName := strings.TrimSuffix(filepath.Base(modulePart), filepath.Ext(modulePart))
			testName := strings.TrimSpace(strings.Split(parts[1], " - ")[0])
			if moduleName != "" && testName != "" {
				selector = fmt.Sprintf("(%s and %s)", moduleName, testName)
			}
		}
		if selector == "" {
			continue
		}
		if _, ok := set[selector]; ok {
			continue
		}
		set[selector] = struct{}{}
		out = append(out, selector)
	}
	sort.Strings(out)
	return out
}

func collectFailingTestsByRerun(ctx context.Context, workdir, testName string) []string {
	py := pythonExecutable()
	runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(runCtx, py, "-m", "pytest", testName, "-q", "--tb=no", "--maxfail=9999")
	cmd.Dir = workdir
	out, _ := cmd.CombinedOutput()
	return collectFailingTestsFromPytestOutput(string(out))
}

func extractMetaMutationStats(workdir string, mutationTargets []string) (mutationStats, bool) {
	mutantsDir := filepath.Join(workdir, "mutants")
	if _, err := os.Stat(mutantsDir); err != nil {
		return mutationStats{}, false
	}
	stats := mutationStats{}
	hasAny := false
	for _, target := range mutationTargets {
		targetName := strings.TrimSuffix(target, filepath.Ext(target))
		metaPath := filepath.Join(mutantsDir, targetName+".py.meta")
		if _, err := os.Stat(metaPath); err != nil {
			fallback := filepath.Join(mutantsDir, target+".meta")
			if _, err2 := os.Stat(fallback); err2 != nil {
				continue
			}
			metaPath = fallback
		}
		hasAny = true
		raw, err := os.ReadFile(metaPath)
		if err != nil {
			continue
		}
		var payload map[string]any
		if err := json.Unmarshal(raw, &payload); err != nil {
			continue
		}
		exitByKey, _ := payload["exit_code_by_key"].(map[string]any)
		for _, code := range exitByKey {
			stats.Total++
			status := mutmutStatusFromExitCode(code)
			switch status {
			case "killed":
				stats.Killed++
			case "survived":
				stats.Survived++
			case "no_tests":
				stats.NoTests++
			case "timeout":
				stats.Timeout++
			case "skipped":
				stats.Skipped++
			case "not_checked":
				stats.NotChecked++
			default:
				stats.Suspicious++
			}
		}
	}
	if !hasAny || stats.Total <= 0 {
		return mutationStats{}, false
	}
	return stats, true
}

func mutmutStatusFromExitCode(code any) string {
	if code == nil {
		return "not_checked"
	}
	parsed, ok := toFloat(code)
	if !ok {
		return "suspicious"
	}
	v := int(parsed)
	switch v {
	case 1, 3:
		return "killed"
	case 0:
		return "survived"
	case 5, 33:
		return "no_tests"
	case 2:
		return "not_checked"
	case -24:
		return "timeout"
	case -11:
		return "suspicious"
	default:
		return "suspicious"
	}
}
