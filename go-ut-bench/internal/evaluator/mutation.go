// evaluator/mutation.go 提供变异测试功能
// 使用 mutmut（Python）、go-mutesting（Go）、PITest（Java）、Mull（C++）执行变异测试
// 计算变异得分，统计 killed/survived/no_tests 等状态
package evaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"go-ut-bench/internal/obs"
)

// mutationLogger 变异测试专用日志记录器
var mutationLogger *obs.Logger

// SetMutationLogger 设置变异测试日志记录器
//
// 参数:
//   - logger: 日志记录器实例
func SetMutationLogger(logger *obs.Logger) {
	mutationLogger = logger
}

// logMutation 记录变异测试日志
//
// 参数:
//   - level: 日志级别
//   - msg: 日志消息
//   - fields: 日志字段
func logMutation(level, msg string, fields ...any) {
	if mutationLogger != nil {
		mutationLogger.ToFile("evaluator").Trace(msg, fields...)
	}
}

// collectPythonMutation 执行 Python 变异测试
// 使用 mutmut 工具对源码进行变异并计算变异得分
//
// 参数:
//   - ctx: 上下文
//   - workdir: 工作目录
//   - testName: 测试文件名
//   - mutationTargets: 变异目标文件列表
//   - timeoutSeconds: 超时时间（秒）
//   - testOutput: 测试输出（用于提取失败测试）
//
// 返回值:
//   - float64: 变异得分（0-1）
//   - mutationStats: 变异统计数据
//   - string: 错误信息（成功时为空）
func collectPythonMutation(ctx context.Context, workdir, testName string, mutationTargets []string, timeoutSeconds int, testOutput string) (float64, mutationStats, string) {
	if len(mutationTargets) == 0 {
		return 0, mutationStats{}, "missing mutation targets"
	}
	if timeoutSeconds <= 0 {
		timeoutSeconds = MutationTimeoutSeconds
	}

	fmt.Printf("        [MUTATION] Python mutmut 开始 | 目标: %v | 超时: %ds | 测试文件: %s\n", mutationTargets, timeoutSeconds, testName)
	logMutation("DEBUG-1", "mutation_start", "workdir", workdir, "test_name", testName, "targets", mutationTargets, "timeout_seconds", timeoutSeconds)

	// Step 1: 收集失败测试
	fmt.Printf("        [MUTATION] 步骤1: 收集失败测试...\n")
	logMutation("DEBUG-2", "mutation_step", "step", "collect_failing_tests", "workdir", workdir)
	failingTestsStart := time.Now()
	failingTests, failingErr := collectFailingTestsByRerun(ctx, workdir, testName)
	if len(failingTests) == 0 && failingErr == "" {
		failingTests = collectFailingTestsFromPytestOutput(testOutput)
	}
	logMutation("DEBUG-2", "mutation_step_done", "step", "collect_failing_tests", "elapsed_ms", time.Since(failingTestsStart).Milliseconds(), "failing_count", len(failingTests), "error", failingErr)
	fmt.Printf("        [MUTATION] 步骤1完成 | 发现失败测试: %d | 耗时: %dms\n", len(failingTests), time.Since(failingTestsStart).Milliseconds())

	if failingErr != "" && len(failingTests) == 0 {
		return 0, mutationStats{}, "pytest rerun before mutation failed: " + failingErr
	}

	if len(failingTests) > 0 {
		logMutation("DEBUG-3", "failing_tests", "tests", failingTests)
	}

	// Step 2: 创建配置文件
	fmt.Printf("        [MUTATION] 步骤2: 创建配置文件...\n")
	logMutation("DEBUG-2", "mutation_step", "step", "create_pyproject")
	pyprojectPath := filepath.Join(workdir, "pyproject.toml")
	if err := os.WriteFile(pyprojectPath, []byte(buildMutmutPyproject(mutationTargets, testName, failingTests)), 0o644); err != nil {
		logMutation("ERROR", "mutation_step_error", "step", "create_pyproject", "error", err.Error())
		return 0, mutationStats{}, err.Error()
	}
	logMutation("DEBUG-2", "mutation_step_done", "step", "create_pyproject")
	fmt.Printf("        [MUTATION] 步骤2完成\n")

	// Step 3: 构建环境变量
	fmt.Printf("        [MUTATION] 步骤3: 构建环境变量...\n")
	logMutation("DEBUG-2", "mutation_step", "step", "build_env")
	env, envErr := buildMutmutEnv(workdir)
	if envErr != nil {
		logMutation("ERROR", "mutation_step_error", "step", "build_env", "error", envErr.Error())
		return 0, mutationStats{}, envErr.Error()
	}
	logMutation("DEBUG-2", "mutation_step_done", "step", "build_env")
	fmt.Printf("        [MUTATION] 步骤3完成\n")

	py := pythonExecutable()
	mutantsDir := filepath.Join(workdir, "mutants")
	_ = os.RemoveAll(mutantsDir)

	// Step 4: 创建 mutants 目录
	fmt.Printf("        [MUTATION] 步骤4: 创建 mutants 目录结构...\n")
	logMutation("DEBUG-2", "mutation_step", "step", "pre_create_mutants_dir", "targets", mutationTargets)
	preCreateStart := time.Now()
	if err := preCreateMutantsDirectory(workdir, mutantsDir, mutationTargets, testName); err != nil {
		logMutation("ERROR", "mutation_step_error", "step", "pre_create_mutants_dir", "error", err.Error())
		return 0, mutationStats{}, "failed to pre-create mutants directory: " + err.Error()
	}
	logMutation("DEBUG-2", "mutation_step_done", "step", "pre_create_mutants_dir", "elapsed_ms", time.Since(preCreateStart).Milliseconds())
	fmt.Printf("        [MUTATION] 步骤4完成 | 耗时: %dms\n", time.Since(preCreateStart).Milliseconds())

	// Step 5: 运行 mutmut run
	fmt.Printf("        [MUTATION] 步骤5: 运行 mutmut run (超时=%ds)...\n", timeoutSeconds)
	logMutation("DEBUG-1", "mutation_step", "step", "mutmut_run", "timeout_seconds", timeoutSeconds, "python", py)
	mutmutRunStart := time.Now()
	runCtx, cancelRun := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancelRun()
	runOut, runErr := runCommandWithProcessGroupKill(runCtx, py, []string{"-m", "mutmut", "run"}, workdir, env)
	mutmutRunElapsed := time.Since(mutmutRunStart)
	logMutation("DEBUG-1", "mutation_step_done", "step", "mutmut_run", "elapsed_ms", mutmutRunElapsed.Milliseconds(), "run_err", runErr)

	if runErr != nil {
		fmt.Printf("        [MUTATION] 步骤5完成(有错误) | 耗时: %dms | 错误: %v\n", mutmutRunElapsed.Milliseconds(), runErr)
		logMutation("DEBUG-2", "mutmut_run_output", "output", string(runOut))
	} else {
		fmt.Printf("        [MUTATION] 步骤5完成 | 耗时: %dms\n", mutmutRunElapsed.Milliseconds())
	}

	// Step 6: 导出统计信息
	fmt.Printf("        [MUTATION] 步骤6: 导出统计信息...\n")
	logMutation("DEBUG-2", "mutation_step", "step", "mutmut_export")
	exportCtx, cancelExport := context.WithTimeout(ctx, 30*time.Second)
	defer cancelExport()
	exportOut, exportErr := runCommandWithProcessGroupKill(exportCtx, py, []string{"-m", "mutmut", "export-cicd-stats"}, workdir, env)
	logMutation("DEBUG-2", "mutation_step_done", "step", "mutmut_export", "export_err", exportErr)
	fmt.Printf("        [MUTATION] 步骤6完成\n")

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

// preCreateMutantsDirectory 预创建 mutants 目录结构
// 复制测试文件和依赖到 mutants 目录
//
// 参数:
//   - workdir: 工作目录
//   - mutantsDir: mutants 目录路径
//   - mutationTargets: 变异目标列表
//   - testName: 测试文件名
//
// 返回值:
//   - error: 创建错误
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

// copyPackageDependencies 复制包依赖到目标目录
// 递归复制 Python 文件，排除 __pycache__ 和目标文件
//
// 参数:
//   - srcPackageDir: 源包目录
//   - dstPackageDir: 目标包目录
//   - targetFileName: 目标文件名（不复制）
//
// 返回值:
//   - error: 复制错误
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

// buildMutmutPyproject 构建 mutmut 配置文件内容
// 设置 paths_to_mutate 和 pytest 参数
//
// 参数:
//   - sourceNames: 源文件列表
//   - testName: 测试文件名
//   - failingTests: 需排除的失败测试列表
//
// 返回值:
//   - string: pyproject.toml 内容
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

// buildMutmutEnv 构建 mutmut 运行环境变量
// 设置 PYTHONPATH 和 multiprocessing shim
//
// 参数:
//   - workdir: 工作目录
//
// 返回值:
//   - []string: 环境变量列表
//   - error: 构建错误
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

// parseMutationStats 解析变异测试统计 JSON
// 从 mutmut-cicd-stats.json 提取 killed/survived 等计数
//
// 参数:
//   - raw: JSON 原始数据
//
// 返回值:
//   - mutationStats: 统计数据
//   - string: 解析错误（成功时为空）
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
	return formatMutationToolError("mutmut", prefix, runErr, runOut, exportErr, exportOut)
}

func formatMutationToolError(tool, prefix string, runErr error, runOut []byte, exportErr error, exportOut []byte) string {
	runMsg := ""
	if runErr != nil {
		runMsg = runErr.Error()
	}
	exportMsg := ""
	if exportErr != nil {
		exportMsg = exportErr.Error()
	}
	return fmt.Sprintf(
		"%s; %s run err=%q; %s export err=%q; run_out=%q; export_out=%q",
		prefix,
		tool,
		runMsg,
		tool,
		exportMsg,
		trimErr(string(runOut), 1200),
		trimErr(string(exportOut), 1200),
	)
}

// inferMutationTargets 从测试文件推断变异目标
// 分析 import 语句，查找可变异的源文件
//
// 参数:
//   - workdir: 工作目录
//   - testName: 测试文件名
//   - sourceBase: 默认源文件名
//
// 返回值:
//   - []string: 变异目标列表
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

func collectFailingTestsByRerun(ctx context.Context, workdir, testName string) ([]string, string) {
	py := pythonExecutable()
	runCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	out, _ := runCommandWithProcessGroupKill(runCtx, py, []string{"-m", "pytest", testName, "-q", "--tb=no", "--maxfail=9999"}, workdir, nil)
	if runCtx.Err() != nil {
		return nil, "pytest timed out after 30s"
	}
	return collectFailingTestsFromPytestOutput(string(out)), ""
}

// extractMetaMutationStats 从 .meta 文件提取变异统计
// 备用方案，当 cicd-stats.json 解析失败时使用
//
// 参数:
//   - workdir: 工作目录
//   - mutationTargets: 变异目标列表
//
// 返回值:
//   - mutationStats: 统计数据
//   - bool: 是否成功提取
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
