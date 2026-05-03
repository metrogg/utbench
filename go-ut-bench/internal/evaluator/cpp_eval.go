// cpp_eval.go 提供 C++ 语言单元测试评测功能
// 使用 CMake/Make 进行编译
// 使用 GoogleTest 进行测试执行
// 使用 gcov 进行覆盖率收集
// 使用 Mull 进行变异测试
package evaluator

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// cppCMakeTemplate CMake 构建模板
// 配置 GTest 和覆盖率选项
const cppCMakeTemplate = `cmake_minimum_required(VERSION 3.10)
project(utbench_eval)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(GTest REQUIRED)

enable_testing()

add_executable(test_runner
    %s
)

target_link_libraries(test_runner GTest::gtest_main gcov)

target_compile_options(test_runner PRIVATE --coverage -fprofile-arcs -ftest-coverage)

add_test(NAME AllTests COMMAND test_runner)
`

// cppMullConfigTemplate 使用 Mull 默认配置（所有变异器）
// 不指定 mutators，让 Mull 使用默认的全部变异器
const cppMullConfigTemplate = `excludePaths:
  - ".*\\.h$"
  - ".*\\.hpp$"
  - "^/usr/.*"
  - ".*googletest.*"

timeout: 30000
`

const cppMullCMakeTemplate = `cmake_minimum_required(VERSION 3.10)
project(utbench_mull)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

set(CMAKE_C_COMPILER clang-18)
set(CMAKE_CXX_COMPILER clang++-18)

find_package(GTest REQUIRED)

enable_testing()

add_executable(test_runner_mull
    %s
)

target_link_libraries(test_runner_mull GTest::gtest_main -lpthread -ldl)

target_compile_options(test_runner_mull PRIVATE
    -fpass-plugin=%s
    -g -grecord-command-line
    -fPIC
    -O0
)

set_target_properties(test_runner_mull PROPERTIES
    BUILD_RPATH "/usr/lib/llvm-18/lib;/usr/lib/x86_64-linux-gnu;/lib/x86_64-linux-gnu;/usr/lib64"
    INSTALL_RPATH "/usr/lib/llvm-18/lib;/usr/lib/x86_64-linux-gnu;/lib/x86_64-linux-gnu;/usr/lib64"
)

add_test(NAME AllTests COMMAND test_runner_mull)
`

// prepareCppWorkspace 准备 C++ 评测工作区
// 创建 CMake 项目结构，复制源码和测试文件
//
// 参数:
//   - testPath: 生成的测试文件路径
//   - samplePath: 源码文件路径
//
// 返回值:
//   - string: 工作目录路径（失败时为空）
//   - string: 测试文件名（失败时为错误信息）
//   - string: 源码文件名
//   - string: 源码文件名（去掉扩展名）
//   - string: 错误信息（成功时为空）
func prepareCppWorkspace(testPath, samplePath string) (string, string, string, string, string) {
	testSource, err := os.ReadFile(testPath)
	if err != nil {
		return "", "", "", "", fmt.Sprintf("failed to read generated test: %s", err)
	}

	sourceBase := filepath.Base(samplePath)
	sourceStem := strings.TrimSuffix(sourceBase, filepath.Ext(sourceBase))
	testFileName := sourceStem + "_test.cpp"

	sourceData, err := os.ReadFile(samplePath)
	if err != nil {
		return "", "", "", "", fmt.Sprintf("failed to read source: %s", err)
	}

	testSource = stripRedundantDefinitions(testSource, sourceData)

	workdir, err := os.MkdirTemp("", "utbench_cpp_eval_")
	if err != nil {
		return "", "", "", "", fmt.Sprintf("failed to create temp dir: %s", err)
	}

	buildDir := filepath.Join(workdir, "build")
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", "", fmt.Sprintf("failed to create build dir: %s", err)
	}

	if err := os.WriteFile(filepath.Join(workdir, sourceBase), sourceData, 0644); err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", "", fmt.Sprintf("failed to write source: %s", err)
	}

	headerPattern := regexp.MustCompile(`#include\s+"([^"]+)"`)
	headerMatches := headerPattern.FindAllStringSubmatch(string(testSource), -1)
	for _, match := range headerMatches {
		if len(match) < 2 {
			continue
		}
		headerName := match[1]
		if isSystemProvidedCppHeader(headerName) {
			continue
		}
		headerPath := filepath.Join(workdir, headerName)
		if _, err := os.Stat(headerPath); os.IsNotExist(err) {
			if strings.HasSuffix(headerName, ".h") || strings.HasSuffix(headerName, ".hpp") {
				if err := os.MkdirAll(filepath.Dir(headerPath), 0o755); err != nil {
					_ = os.RemoveAll(workdir)
					return "", "", "", "", fmt.Sprintf("failed to create header dir for %s: %s", headerName, err)
				}
				declHeader := generatePlaceholderHeader(headerName)
				if err := os.WriteFile(headerPath, []byte(declHeader), 0644); err != nil {
					_ = os.RemoveAll(workdir)
					return "", "", "", "", fmt.Sprintf("failed to write header %s: %s", headerName, err)
				}
			}
		}
	}

	hasSourceInclude := bytes.Contains(testSource, []byte("#include \""+sourceBase+"\"")) ||
		bytes.Contains(testSource, []byte("#include <"+sourceBase+">")) ||
		bytes.Contains(testSource, []byte("#include \"source.cpp\"")) ||
		bytes.Contains(testSource, []byte("#include <source.cpp>"))

	modifiedTestSource := testSource
	if !hasSourceInclude {
		sourceInclude := []byte("#include \"" + sourceBase + "\"\n")
		modifiedTestSource = append(sourceInclude, testSource...)
	} else {
		modifiedTestSource = bytes.ReplaceAll(modifiedTestSource,
			[]byte("#include \"source.cpp\""),
			[]byte("#include \""+sourceBase+"\""))
		modifiedTestSource = bytes.ReplaceAll(modifiedTestSource,
			[]byte("#include <source.cpp>"),
			[]byte("#include \""+sourceBase+"\""))
	}

	if err := os.WriteFile(filepath.Join(workdir, testFileName), modifiedTestSource, 0644); err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", "", fmt.Sprintf("failed to write test: %s", err)
	}

	cmakeContent := fmt.Sprintf(cppCMakeTemplate, testFileName)
	if err := os.WriteFile(filepath.Join(workdir, "CMakeLists.txt"), []byte(cmakeContent), 0644); err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", "", fmt.Sprintf("failed to write CMakeLists.txt: %s", err)
	}

	return workdir, testFileName, sourceBase, sourceStem, ""
}

// cppCompileCheck 检查 C++ 测试代码是否能编译通过
// 使用 CMake 和 Make 进行编译检查
//
// 参数:
//   - workdir: 工作目录（包含 CMakeLists.txt）
//
// 返回值:
//   - bool: 编译是否通过
//   - string: 编译错误信息（成功时为空）
func cppCompileCheck(workdir string) (bool, string) {
	buildDir := filepath.Join(workdir, "build")

	cmakeCtx, cancelCMake := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelCMake()
	cmakeOut, cmakeErr := runCommandWithProcessGroupKill(cmakeCtx, "cmake", []string{".."}, buildDir, nil)
	if cmakeCtx.Err() != nil {
		return false, "cmake timed out after 120s"
	}
	if cmakeErr != nil {
		return false, trimErr(string(cmakeOut), 2000)
	}

	makeCtx, cancelMake := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancelMake()
	makeOut, makeErr := runCommandWithProcessGroupKill(makeCtx, "make", []string{"-j2"}, buildDir, nil)
	if makeCtx.Err() != nil {
		return false, "make timed out after 120s"
	}
	if makeErr != nil {
		return false, trimErr(string(makeOut), 2000)
	}

	return true, ""
}

// executeCppTests 执行 C++ 测试
// 运行 build 目录中的 test_runner 可执行文件
//
// 参数:
//   - workdir: 工作目录
//
// 返回值:
//   - bool: 测试是否通过
//   - string: 测试输出或错误信息
//   - int: 执行耗时（毫秒）
func executeCppTests(workdir string) (bool, string, int) {
	buildDir := filepath.Join(workdir, "build")

	runCtx, cancel := context.WithTimeout(context.Background(), defaultTestTimeoutSeconds*time.Second)
	defer cancel()
	started := time.Now()
	output, err := runCommandWithProcessGroupKill(runCtx, "./test_runner", nil, buildDir, nil)
	latency := int(time.Since(started).Milliseconds())
	if runCtx.Err() != nil {
		return false, fmt.Sprintf("cpp test timed out after %ds", defaultTestTimeoutSeconds), latency
	}
	if err == nil {
		return true, string(output), latency
	}
	return false, trimErr(string(output), 4000), latency
}

func parseCppTestCounts(output string) (*int, *int) {
	passedPattern := regexp.MustCompile(`\[(\d+)\/(\d+)\] PASSED`)
	failedPattern := regexp.MustCompile(`\[(\d+)\/(\d+)\] FAILED`)

	passedMatches := passedPattern.FindAllStringSubmatch(output, -1)
	failedMatches := failedPattern.FindAllStringSubmatch(output, -1)

	passed := len(passedMatches)
	failed := len(failedMatches)

	if passed > 0 || failed > 0 {
		total := passed + failed
		return &passed, &total
	}

	passedPattern2 := regexp.MustCompile(`(\d+)\s+tests?\s+from`)
	if match := passedPattern2.FindStringSubmatch(output); match != nil && len(match) > 1 {
		total := parseIntOrZero(match[1])
		if total > 0 {
			failedCount := 0
			failedPattern2 := regexp.MustCompile(`\[  FAILED  \]\s*(\d+)\s*tests?`)
			if match2 := failedPattern2.FindStringSubmatch(output); match2 != nil && len(match2) > 1 {
				failedCount = parseIntOrZero(match2[1])
			} else {
				failedPattern3 := regexp.MustCompile(`(\d+)\s+FAILED`)
				if match3 := failedPattern3.FindStringSubmatch(output); match3 != nil && len(match3) > 1 {
					failedCount = parseIntOrZero(match3[1])
				}
			}
			passedCount := total - failedCount
			if passedCount < 0 {
				passedCount = 0
			}
			return &passedCount, &total
		}
	}

	if strings.Contains(output, "[ PASSED  ]") && !strings.Contains(output, "[ FAILED  ]") {
		p := 1
		t := 1
		return &p, &t
	}

	return nil, nil
}

// collectCppCoverage 收集 C++ 测试覆盖率数据
// 使用 gcov 工具从 .gcno/.gcda 文件解析覆盖率
//
// 参数:
//   - workdir: 工作目录
//   - testFileName: 测试文件名
//
// 返回值:
//   - float64: 行覆盖率（0-1）
//   - float64: 分支覆盖率（0-1）
//   - string: 错误信息（成功时为空）
func collectCppCoverage(workdir, testFileName string) (float64, float64, string) {
	buildDir := filepath.Join(workdir, "build")
	gcovDir := filepath.Join(buildDir, "CMakeFiles", "test_runner.dir")

	if _, err := os.Stat(gcovDir); err != nil {
		return 0, 0, "gcov directory not found: " + gcovDir
	}

	gcnoFile := filepath.Join(gcovDir, testFileName+".gcno")
	gcdaFile := filepath.Join(gcovDir, testFileName+".gcda")
	if _, err := os.Stat(gcnoFile); err != nil {
		return 0, 0, "gcno file not found: " + gcnoFile
	}

	gcnoLink := filepath.Join(gcovDir, strings.TrimSuffix(testFileName, filepath.Ext(testFileName))+".gcno")
	gcdaLink := filepath.Join(gcovDir, strings.TrimSuffix(testFileName, filepath.Ext(testFileName))+".gcda")

	if _, err := os.Stat(gcnoLink); err != nil {
		if err := copyFile(gcnoFile, gcnoLink); err != nil {
			return 0, 0, "failed to copy gcno file: " + err.Error()
		}
	}

	if _, err := os.Stat(gcdaFile); err == nil {
		if _, err := os.Stat(gcdaLink); err != nil {
			if err := copyFile(gcdaFile, gcdaLink); err != nil {
				return 0, 0, "failed to copy gcda file: " + err.Error()
			}
		}
	}

	gcovCtx, cancelGCov := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelGCov()
	gcovOut, gcovErr := runCommandWithProcessGroupKill(gcovCtx, "gcov", []string{"-b", testFileName}, gcovDir, nil)
	if gcovCtx.Err() != nil {
		return 0, 0, "gcov timed out after 30s"
	}
	if gcovErr != nil {
		return 0, 0, trimErr(string(gcovOut), 2000)
	}

	gcovFile := filepath.Join(gcovDir, testFileName+".gcov")
	raw, err := os.ReadFile(gcovFile)
	if err != nil {
		return 0, 0, fmt.Sprintf("failed to read gcov file: %s", err)
	}

	return parseGcovOutput(string(raw))
}

func parseGcovOutput(content string) (float64, float64, string) {
	lines := strings.Split(content, "\n")

	totalLines := 0
	coveredLines := 0
	totalBranches := 0
	coveredBranches := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Skip metadata lines (Source:, Graph:, Data:, Runs:)
		if strings.HasPrefix(line, "Source:") || strings.HasPrefix(line, "Graph:") ||
			strings.HasPrefix(line, "Data:") || strings.HasPrefix(line, "Runs:") {
			continue
		}

		// Skip function/call summary lines
		if strings.HasPrefix(line, "function ") || strings.HasPrefix(line, "call ") {
			continue
		}

		// Parse branch coverage: "branch  0 taken 100% (fallthrough)" or "branch  1 taken 0%"
		// or "branch  2 taken never"
		if strings.HasPrefix(line, "branch") {
			totalBranches++
			pct, ok := parseBranchPercent(line)
			if ok && pct > 0 {
				coveredBranches++
			}
			continue
		}

		parts := strings.SplitN(line, ":", 3)
		if len(parts) < 3 {
			continue
		}

		execCountStr := strings.TrimSpace(parts[0])
		lineNumStr := strings.TrimSpace(parts[1])
		lineCode := strings.TrimSpace(parts[2])

		// Skip if lineNum is 0 (metadata) or not a valid line number
		if lineNumStr == "0" {
			continue
		}
		if _, err := strconv.Atoi(lineNumStr); err != nil {
			continue
		}

		// Skip empty code lines
		if lineCode == "" {
			continue
		}

		// Only count lines that are in code blocks (not "-:" prefix)
		// Lines with "-" count are not in any code block
		if execCountStr == "-" {
			continue
		}

		totalLines++

		if execCountStr != "" && execCountStr != "######" && execCountStr != "======" {
			if count, err := strconv.Atoi(execCountStr); err == nil && count > 0 {
				coveredLines++
			}
		}
	}

	if totalLines == 0 {
		return 0, 0, "no coverage data"
	}

	lineCov := round(float64(coveredLines)/float64(totalLines), 6)

	var branchCov float64
	if totalBranches > 0 {
		branchCov = round(float64(coveredBranches)/float64(totalBranches), 6)
	}

	return lineCov, branchCov, ""
}

func parseBranchPercent(line string) (float64, bool) {
	if !strings.Contains(line, "taken") {
		return 0, false
	}
	idx := strings.Index(line, "taken")
	afterTaken := strings.TrimSpace(line[idx+5:])
	if strings.HasPrefix(afterTaken, "never") {
		return 0, false
	}
	parts := strings.Fields(afterTaken)
	if len(parts) == 0 {
		return 0, false
	}
	pctStr := strings.TrimSuffix(parts[0], "%")
	pct, err := strconv.ParseFloat(pctStr, 64)
	if err != nil {
		return 0, false
	}
	return pct, true
}

// collectCppMutation 执行 C++ 变异测试
// 使用 Mull 工具对源码进行变异并计算变异得分
//
// 参数:
//   - ctx: 上下文
//   - workdir: 工作目录
//   - sourceBase: 源码文件名
//   - timeoutSeconds: 超时时间（秒），最大限制300秒
//   - testPassRate: 测试通过率（用于判断是否运行变异测试）
//   - testPassed: 通过的测试数
//   - testTotal: 总测试数
//
// 返回值:
//   - float64: 变异得分（0-1）
//   - mutationStats: 变异统计数据
//   - string: 错误信息（成功时为空）
func collectCppMutation(ctx context.Context, workdir, sourceBase string, timeoutSeconds int, testPassRate *float64, testPassed, testTotal int) (float64, mutationStats, string) {
	// 限制最大超时时间为300秒（与 Python 版本一致），防止Mull无限卡住
	if timeoutSeconds <= 0 || timeoutSeconds > 300 {
		timeoutSeconds = 300
	}

	minPassRate := GetMinPassRateForTool("mull")
	passed := 0
	total := 0
	if testPassed > 0 || testTotal > 0 {
		passed = testPassed
		total = testTotal
	} else if testPassRate != nil {
		total = 100
		passed = int(math.Round(*testPassRate * float64(total)))
		if passed == 0 && *testPassRate > 0 {
			passed = 1
		}
		if passed > total {
			passed = total
		}
	}

	checkResult := CheckTestPassRate(passed, total, "Mull", minPassRate)
	if !checkResult.ShouldRun {
		return 0, mutationStats{}, checkResult.Message
	}

	mullRunner := findMullRunner()
	if mullRunner == "" {
		return 0, mutationStats{}, "Mull not installed. Install: curl -1sLf 'https://dl.cloudsmith.io/public/mull-project/mull-stable/setup.deb.sh' | bash && apt-get install -y mull-19"
	}

	// Find mull-ir-frontend plugin path
	mullFrontend := findMullFrontend()
	if mullFrontend == "" {
		return 0, mutationStats{}, "mull-ir-frontend-19 not found"
	}

	testFileName := strings.TrimSuffix(sourceBase, filepath.Ext(sourceBase)) + "_test.cpp"

	// Write mull.yml config
	mullConfigPath := filepath.Join(workdir, "mull.yml")
	if err := os.WriteFile(mullConfigPath, []byte(cppMullConfigTemplate), 0644); err != nil {
		return 0, mutationStats{}, "failed to write mull.yml: " + err.Error()
	}

	// 手动编译（与 Python 版本一致，不用 CMake）
	// Python: common_flags = ["-std=c++17", "-O0", "-g", "-fno-inline", "-fno-omit-frame-pointer"]
	// Mull 需要 -grecord-command-line 来获取编译标志
	commonFlags := []string{"-std=c++17", "-O0", "-g", "-grecord-command-line", "-fno-inline", "-fno-omit-frame-pointer"}
	if mullFrontend != "" {
		commonFlags = append(commonFlags, "-fpass-plugin="+mullFrontend)
	}

	// include dirs（与 Python 版本一致）
	includeDirs := []string{workdir}

	// 编译测试文件（测试文件通过 #include 包含了源文件，所以只编译测试文件即可）
	// 如果同时编译源文件，会导致函数多重定义链接错误
	testObj := "generated_test_mull.o"

	compileArgs := append([]string{}, commonFlags...)
	for _, inc := range includeDirs {
		compileArgs = append(compileArgs, "-I", inc)
	}
	compileArgs = append(compileArgs, "-c", testFileName, "-o", testObj)

	compileCtx, cancelCompile := context.WithTimeout(ctx, 120*time.Second)
	out, err := runCommandWithProcessGroupKill(compileCtx, "clang++-19", compileArgs, workdir, nil)
	compileCtxErr := compileCtx.Err()
	cancelCompile()
	if compileCtxErr == context.DeadlineExceeded {
		return 0, mutationStats{}, "compile for Mull timed out after 120s"
	}
	if compileCtxErr != nil {
		return 0, mutationStats{}, "compile for Mull canceled: " + compileCtxErr.Error()
	}
	if err != nil {
		return 0, mutationStats{}, fmt.Sprintf("compile for Mull failed\n%s", trimErr(string(out), 4000))
	}

	// 链接生成可执行文件（与 Python 版本一致，120 秒超时）
	// 需要链接 gtest 和 gtest_main（gtest_main 提供默认的 main 函数）
	linkArgs := append([]string{}, commonFlags...)
	linkArgs = append(linkArgs, testObj, "-o", "run_tests_mull")
	// 显式指定 gtest 库路径，确保能找到 gtest_main
	linkArgs = append(linkArgs, "/usr/lib/libgtest.a", "/usr/lib/libgtest_main.a")
	linkArgs = append(linkArgs, "-lpthread", "-ldl")

	linkCtx, cancelLink := context.WithTimeout(ctx, 120*time.Second)
	linkOut, linkErr := runCommandWithProcessGroupKill(linkCtx, "clang++-19", linkArgs, workdir, nil)
	linkCtxErr := linkCtx.Err()
	cancelLink()
	if linkCtxErr == context.DeadlineExceeded {
		return 0, mutationStats{}, "link for Mull timed out after 120s"
	}
	if linkCtxErr != nil {
		return 0, mutationStats{}, "link for Mull canceled: " + linkCtxErr.Error()
	}
	if linkErr != nil {
		return 0, mutationStats{}, fmt.Sprintf("link for Mull failed\n%s", trimErr(string(linkOut), 4000))
	}

	binaryPath := filepath.Join(workdir, "run_tests_mull")
	if _, err := os.Stat(binaryPath); err != nil {
		return 0, mutationStats{}, "mull executable not found after build: " + binaryPath
	}

	// 运行 mull-runner（与 Python 版本一致）
	runCtx, cancelRun := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancelRun()

	mullEnv := append(os.Environ(),
		"LD_LIBRARY_PATH=/usr/lib/llvm-19/lib:/usr/lib/x86_64-linux-gnu:/lib/x86_64-linux-gnu:/usr/lib64:"+os.Getenv("LD_LIBRARY_PATH"),
	)
	mullOut, mullErr := runCommandWithProcessGroupKill(runCtx, mullRunner, []string{binaryPath}, workdir, mullEnv)

	stats, parseErr := parseMullOutput(string(mullOut))
	if parseErr != "" {
		return 0, stats, formatMullError("mull error", mullErr, mullOut)
	}

	if stats.Total <= 0 {
		return 0, stats, formatMullError("mull produced zero mutants", mullErr, mullOut)
	}

	// 分数计算：使用 killed / total（与 Python 版本一致）
	score := round(float64(stats.Killed)/float64(stats.Total), 6)
	return score, stats, ""
}

func formatMullError(prefix string, runErr error, runOut []byte) string {
	runMsg := ""
	if runErr != nil {
		runMsg = runErr.Error()
	}
	return fmt.Sprintf(
		"%s; mull-runner err=%q; mull_out=%q",
		prefix,
		runMsg,
		trimErr(string(runOut), 1500),
	)
}

func findMullRunner() string {
	candidates := []string{"mull-runner-19", "mull-runner-18", "mull-runner"}
	for _, c := range candidates {
		if path, err := exec.LookPath(c); err == nil {
			return path
		}
	}
	return ""
}

func findMullFrontend() string {
	// Try common paths for mull-ir-frontend (Mull 19 需要 LLVM 19)
	candidates := []string{
		"/usr/lib/mull-ir-frontend-19",
		"/usr/lib/llvm-19/lib/mull-ir-frontend-19.so",
		"/usr/lib/x86_64-linux-gnu/mull-ir-frontend-19.so",
		"/usr/local/lib/mull-ir-frontend-19.so",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

// parseMullOutput 解析 Mull 输出，支持 8 个统计维度
// 与 Python 版本保持一致：total, killed, survived, timeout, no_tests, not_checked, skipped, suspicious
func parseMullOutput(output string) (mutationStats, string) {
	stats := mutationStats{}

	if strings.Contains(output, "Original test failed (warmup run)") {
		return stats, "mull warmup run failed: baseline test has failing assertions"
	}

	// 情况1: 所有变异体都被杀死
	if strings.Contains(output, "All mutations have been killed") {
		re := regexp.MustCompile(`(\d+)/(\d+)\s*\.?\s*Finished`)
		matches := re.FindAllStringSubmatch(output, -1)
		if len(matches) > 0 && len(matches[len(matches)-1]) > 2 {
			lastMatch := matches[len(matches)-1]
			stats.Total = parseIntOrZero(lastMatch[2])
			stats.Killed = stats.Total
			stats.Survived = 0
			return stats, ""
		}
	}

	// 解析各个维度的统计信息（多种命名格式兼容，与 Python 版本一致）
	// Killed mutants
	killedPatterns := []string{
		`Killed mutants\s*\((\d+)/(\d+)\)`,
		`Killed\s*(?:mutants|mutations)?\s*[:=]\s*(\d+)`,
	}
	for _, pattern := range killedPatterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindStringSubmatch(output); match != nil && len(match) > 1 {
			stats.Killed = parseIntOrZero(match[1])
			if len(match) > 2 && stats.Total == 0 {
				stats.Total = parseIntOrZero(match[2])
			}
			break
		}
	}

	// Survived mutants
	survivedPatterns := []string{
		`Survived mutants\s*\((\d+)/(\d+)\)`,
		`Survived\s*(?:mutants|mutations)?\s*[:=]\s*(\d+)`,
		`Surviving\s*(?:mutants|mutations)?\s*[:=]\s*(\d+)`,
	}
	for _, pattern := range survivedPatterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindStringSubmatch(output); match != nil && len(match) > 1 {
			stats.Survived = parseIntOrZero(match[1])
			break
		}
	}

	// Timeout mutants
	timeoutPatterns := []string{
		`(?:Timed out|Timeout)\s*(?:mutants|mutations)?\s*[:=]\s*(\d+)`,
		`Timeout\s*[:=]\s*(\d+)`,
	}
	for _, pattern := range timeoutPatterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindStringSubmatch(output); match != nil && len(match) > 1 {
			stats.Timeout = parseIntOrZero(match[1])
			break
		}
	}

	// No tests / Not covered
	noTestsPatterns := []string{
		`(?:No tests|Not covered)\s*(?:mutants|mutations)?\s*[:=]\s*(\d+)`,
		`NoCoverage\s*[:=]\s*(\d+)`,
		`NotCovered\s*[:=]\s*(\d+)`,
	}
	for _, pattern := range noTestsPatterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindStringSubmatch(output); match != nil && len(match) > 1 {
			stats.NoTests = parseIntOrZero(match[1])
			break
		}
	}

	// Not checked
	notCheckedPatterns := []string{
		`Not checked\s*(?:mutants|mutations)?\s*[:=]\s*(\d+)`,
		`NotChecked\s*[:=]\s*(\d+)`,
	}
	for _, pattern := range notCheckedPatterns {
		re := regexp.MustCompile(pattern)
		if match := re.FindStringSubmatch(output); match != nil && len(match) > 1 {
			stats.NotChecked = parseIntOrZero(match[1])
			break
		}
	}

	// Skipped
	skippedPattern := regexp.MustCompile(`Skipped\s*(?:mutants|mutations)?\s*[:=]\s*(\d+)`)
	if match := skippedPattern.FindStringSubmatch(output); match != nil && len(match) > 1 {
		stats.Skipped = parseIntOrZero(match[1])
	}

	// Suspicious
	suspiciousPattern := regexp.MustCompile(`Suspicious\s*(?:mutants|mutations)?\s*[:=]\s*(\d+)`)
	if match := suspiciousPattern.FindStringSubmatch(output); match != nil && len(match) > 1 {
		stats.Suspicious = parseIntOrZero(match[1])
	}

	// 从 X/Y 格式提取总数（Mull 0.33.x 进度显示）
	if stats.Total == 0 {
		fractions := regexp.MustCompile(`(\d+)\s*/\s*(\d+)`)
		matches := fractions.FindAllStringSubmatch(output, -1)
		if len(matches) > 0 {
			maxTotal := 0
			for _, match := range matches {
				if len(match) > 2 {
					total := parseIntOrZero(match[2])
					if total > maxTotal {
						maxTotal = total
					}
				}
			}
			if maxTotal > 0 {
				stats.Total = maxTotal
			}
		}
	}

	// 备用：从 Killed:/Survived: 计数
	if stats.Total == 0 {
		killedCount := strings.Count(output, "Killed:")
		survivedCount := strings.Count(output, "Survived:")
		if killedCount > 0 || survivedCount > 0 {
			stats.Killed = killedCount
			stats.Survived = survivedCount
			stats.Total = killedCount + survivedCount
		}
	}

	// 从 Mutation score 反推 killed
	if stats.Killed == 0 && stats.Total > 0 {
		scorePattern := regexp.MustCompile(`[Mm]utation\s*[Ss]core\s*[:=]\s*(\d+(?:\.\d+)?)\s*%`)
		if match := scorePattern.FindStringSubmatch(output); match != nil && len(match) > 1 {
			scorePct := parseFloatOrZero(match[1])
			if scorePct > 0 {
				stats.Killed = int(float64(stats.Total) * scorePct / 100.0)
			}
		}
	}

	// 如果没有解析到 Total，尝试从各个维度计算
	if stats.Total == 0 {
		derivedTotal := stats.Killed + stats.Survived + stats.Timeout + stats.NoTests + stats.NotChecked + stats.Skipped + stats.Suspicious
		if derivedTotal > 0 {
			stats.Total = derivedTotal
		}
	}

	if stats.Total == 0 {
		return stats, "no mutation stats found in mull output"
	}

	return stats, ""
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func isSystemProvidedCppHeader(headerName string) bool {
	return strings.HasPrefix(headerName, "gtest/") || strings.HasPrefix(headerName, "gmock/")
}

func generatePlaceholderHeader(headerName string) string {
	guard := strings.ToUpper(strings.NewReplacer("/", "_", ".", "_", "-", "_").Replace(headerName))
	return fmt.Sprintf("#ifndef %s\n#define %s\n\n#endif\n", guard, guard)
}

// parseFloatOrZero 解析字符串为 float64，失败返回 0
func parseFloatOrZero(s string) float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0
	}
	return f
}

func estimateCppAssertionDensity(testPath string) (int, int, float64) {
	raw, err := os.ReadFile(testPath)
	if err != nil {
		return 0, 0, 0
	}
	text := string(raw)

	// 统计各类断言（概念上都是验证点）
	assertCount := 0

	// 1. GoogleTest 断言宏
	// EXPECT_*: 验证失败继续执行（软断言）
	// ASSERT_*: 验证失败终止当前测试（硬断言）
	assertCount += strings.Count(text, "EXPECT_")
	assertCount += strings.Count(text, "ASSERT_")

	// 2. gMock 验证（验证调用行为）
	// EXPECT_CALL(mock, method()): 设置期望调用，测试结束时验证
	// 这是 Mock 验证，概念上也是断言
	assertCount += strings.Count(text, "EXPECT_CALL(")

	// 3. Catch2 断言（另一个 C++ 测试框架）
	// REQUIRE_*: 硬断言（失败终止）
	// CHECK_*: 软断言（失败继续）
	assertCount += strings.Count(text, "REQUIRE(")
	assertCount += strings.Count(text, "CHECK(")
	assertCount += strings.Count(text, "REQUIRE_")
	assertCount += strings.Count(text, "CHECK_")

	// 统计测试用例：TEST, TEST_F, TEST_P 宏
	// GoogleTest 支持三种测试宏：
	//   TEST(TestSuite, TestName) - 普通测试
	//   TEST_F(TestFixture, TestName) - 固定测试（使用 fixture）
	//   TEST_P(TestFixture, TestName) - 参数化测试
	testPattern := regexp.MustCompile(`TEST(?:_F|_P)?\s*\(\s*[^,]+\s*,\s*[^)]+\s*\)`)
	testMatches := testPattern.FindAllString(text, -1)
	testCount := len(testMatches)

	// 统计 Catch2 测试用例
	// TEST_CASE("name") 或 TEST_CASE_METHOD(Fixture, "name")
	catchPattern := regexp.MustCompile(`TEST_CASE(?:_METHOD)?\s*\(\s*`)
	catchMatches := catchPattern.FindAllString(text, -1)
	testCount += len(catchMatches)

	if testCount <= 0 {
		return assertCount, 0, 0
	}

	return assertCount, testCount, round(float64(assertCount)/float64(testCount), 6)
}

func generateDeclarationsHeader(sourceCode, sourceStem string) string {
	var decls []string
	decls = append(decls, "#ifndef "+strings.ToUpper(sourceStem)+"_DECL_H")
	decls = append(decls, "#define "+strings.ToUpper(sourceStem)+"_DECL_H")
	decls = append(decls, "")

	seenDecls := make(map[string]bool)

	enumPattern := regexp.MustCompile(`(?m)^\s*enum\s+(?:class\s+)?(\w+)(?:\s*:\s*\w+)?\s*\{[^}]+\};?`)
	enumMatches := enumPattern.FindAllString(sourceCode, -1)
	if len(enumMatches) > 0 {
		decls = append(decls, "// Enums")
		for _, e := range enumMatches {
			enumLine := strings.TrimSpace(e)
			if !seenDecls[enumLine] {
				seenDecls[enumLine] = true
				decls = append(decls, enumLine)
			}
		}
		decls = append(decls, "")
	}

	typedefPattern := regexp.MustCompile(`(?m)^\s*typedef\s+[^;]+;`)
	typedefMatches := typedefPattern.FindAllString(sourceCode, -1)
	if len(typedefMatches) > 0 {
		decls = append(decls, "// Typedefs")
		for _, td := range typedefMatches {
			t := strings.TrimSpace(td)
			if !seenDecls[t] {
				seenDecls[t] = true
				decls = append(decls, t)
			}
		}
		decls = append(decls, "")
	}

	usingPattern := regexp.MustCompile(`(?m)^\s*using\s+\w+\s*=\s*[^;]+;`)
	usingMatches := usingPattern.FindAllString(sourceCode, -1)
	if len(usingMatches) > 0 {
		decls = append(decls, "// Using declarations")
		for _, u := range usingMatches {
			t := strings.TrimSpace(u)
			if !seenDecls[t] {
				seenDecls[t] = true
				decls = append(decls, t)
			}
		}
		decls = append(decls, "")
	}

	classStructPattern := regexp.MustCompile(`(?m)^\s*(?:class|struct)\s+(\w+)(?:\s*:\s*(?:public|private|protected)\s+\w+)?\s*;`)
	classStructForwardMatches := classStructPattern.FindAllStringSubmatch(sourceCode, -1)
	classNames := make(map[string]bool)
	if len(classStructForwardMatches) > 0 {
		decls = append(decls, "// Forward declarations")
		for _, match := range classStructForwardMatches {
			if len(match) > 1 {
				name := match[1]
				if !classNames[name] {
					classNames[name] = true
					decls = append(decls, fmt.Sprintf("%s %s;", match[0][:len(match[0])-1], name))
				}
			}
		}
		decls = append(decls, "")
	}

	classDefPattern := regexp.MustCompile(`(?m)^\s*(class|struct)\s+(\w+)(?:\s*:\s*(?:public|private|protected)\s+\w+)?\s*\{`)
	classDefStarts := classDefPattern.FindAllStringIndex(sourceCode, -1)
	classEndPattern := regexp.MustCompile(`\}`)
	classDefEnds := classEndPattern.FindAllStringIndex(sourceCode, -1)

	classRanges := make([]struct{ start, end int }, 0)
	classNameMap := make(map[int]string)
	classIdx := 0
	for _, startIdx := range classDefStarts {
		if classIdx >= len(classDefEnds) {
			continue
		}
		className := classDefPattern.FindStringSubmatch(sourceCode[startIdx[0]:startIdx[1]])[2]
		braceCount := 0
		endIdx := -1
		for _, endIdxPair := range classDefEnds {
			if endIdxPair[0] < startIdx[1] {
				continue
			}
			endIdx = endIdxPair[0]
			for i := startIdx[1]; i < endIdx; i++ {
				if sourceCode[i] == '{' {
					braceCount++
				} else if sourceCode[i] == '}' {
					braceCount--
					if braceCount == 0 {
						break
					}
				}
			}
			if braceCount == 0 {
				break
			}
		}
		if endIdx > startIdx[0] {
			classRanges = append(classRanges, struct{ start, end int }{startIdx[0], endIdx + 1})
			classNameMap[startIdx[0]] = className
		}
		classIdx++
	}

	extractedClasses := make(map[string]bool)
	decls = append(decls, "// Classes and Structs")
	for _, cr := range classRanges {
		className := classNameMap[cr.start]
		if extractedClasses[className] {
			continue
		}
		extractedClasses[className] = true
		classBlock := sourceCode[cr.start:cr.end]
		firstLine := strings.Split(classBlock, "\n")[0]
		if strings.TrimSpace(firstLine) == "" {
			continue
		}
		if !strings.Contains(firstLine, "{") {
			continue
		}
		var classDecl string
		bracePos := strings.Index(classBlock, "{")
		if bracePos > 0 {
			firstPart := strings.TrimSpace(classBlock[:bracePos])
			if strings.HasSuffix(firstPart, "}") {
				classDecl = firstPart + ";"
			} else {
				classDecl = firstPart + " { /* ... */ };"
			}
		} else {
			classDecl = classBlock
		}
		classDecl = strings.TrimSpace(classDecl)
		if classDecl != "" && !seenDecls[classDecl] {
			seenDecls[classDecl] = true
			decls = append(decls, classDecl)
		}
	}
	decls = append(decls, "")

	funcPattern := regexp.MustCompile(`(?m)^\s*(?:inline\s+)?(?:static\s+)?(?:virtual\s+)?(?:const\s+)?(?:explicit\s+)?(?:unsigned\s+)?(?:signed\s+)?(?:void|int|char|short|long|float|double|bool|auto|\w+)(?:\s*\*|\s*&)*\s+(\w+)\s*\([^)]*\)\s*(?:const)?\s*(?:override)?\s*;`)
	funcMatches := funcPattern.FindAllStringSubmatch(sourceCode, -1)
	if len(funcMatches) > 0 {
		decls = append(decls, "// Function declarations")
		for _, match := range funcMatches {
			if len(match) > 1 {
				funcName := match[1]
				if funcName == "if" || funcName == "while" || funcName == "for" || funcName == "switch" {
					continue
				}
			}
			if len(match) > 0 {
				fn := strings.TrimSpace(match[0])
				if !seenDecls[fn] {
					seenDecls[fn] = true
					decls = append(decls, fn)
				}
			}
		}
		decls = append(decls, "")
	}

	decls = append(decls, "#endif // "+strings.ToUpper(sourceStem)+"_DECL_H")

	return strings.Join(decls, "\n")
}

func stripRedundantDefinitions(testSource, sourceData []byte) []byte {
	testCode := string(testSource)
	sourceCode := string(sourceData)

	classNames := extractClassStructNames(sourceCode)
	constNames := extractConstNames(sourceCode)
	funcNames := extractTopLevelFuncNames(sourceCode)
	structNames := extractStructNames(sourceCode)

	result := stripClassDefinitions(testCode, classNames)
	result = stripStructDefinitions(result, structNames)
	result = stripConstDefinitions(result, constNames)
	result = stripTopLevelFuncDefinitions(result, funcNames)

	return []byte(result)
}

func extractClassStructNames(sourceCode string) []string {
	pattern := regexp.MustCompile(`(?m)^\s*(?:class|struct)\s+(\w+)\s*(?:\{|:|\s*$)`)
	matches := pattern.FindAllStringSubmatch(sourceCode, -1)
	names := make(map[string]bool)
	for _, m := range matches {
		if len(m) > 1 {
			names[m[1]] = true
		}
	}
	result := make([]string, 0, len(names))
	for n := range names {
		result = append(result, n)
	}
	return result
}

func extractStructNames(sourceCode string) []string {
	pattern := regexp.MustCompile(`(?m)^\s*struct\s+(\w+)\s*\{`)
	matches := pattern.FindAllStringSubmatch(sourceCode, -1)
	names := make(map[string]bool)
	for _, m := range matches {
		if len(m) > 1 {
			names[m[1]] = true
		}
	}
	result := make([]string, 0, len(names))
	for n := range names {
		result = append(result, n)
	}
	return result
}

func extractConstNames(sourceCode string) []string {
	pattern := regexp.MustCompile(`(?m)^\s*(?:const|static\s+const)\s+\w+\s+(\w+)\s*=`)
	matches := pattern.FindAllStringSubmatch(sourceCode, -1)
	names := make(map[string]bool)
	for _, m := range matches {
		if len(m) > 1 {
			names[m[1]] = true
		}
	}
	result := make([]string, 0, len(names))
	for n := range names {
		result = append(result, n)
	}
	return result
}

func extractTopLevelFuncNames(sourceCode string) []string {
	pattern := regexp.MustCompile(`(?m)^\s*(?:inline\s+)?(?:static\s+)?(?:\w+(?:\s*\*|\s*&)?\s+)+(\w+)\s*\([^)]*\)\s*\{`)
	matches := pattern.FindAllStringSubmatch(sourceCode, -1)
	names := make(map[string]bool)
	for _, m := range matches {
		if len(m) > 1 {
			name := m[1]
			if name == "if" || name == "while" || name == "for" || name == "switch" || name == "main" {
				continue
			}
			names[name] = true
		}
	}
	result := make([]string, 0, len(names))
	for n := range names {
		result = append(result, n)
	}
	return result
}

func stripClassDefinitions(code string, classNames []string) string {
	for _, name := range classNames {
		pattern := regexp.MustCompile(`(?m)(?://[^\n]*\n)?\s*class\s+` + name + `\s*(?:\{|:\s*(?:public|private|protected)\s+\w+\s*\{)`)

		matches := pattern.FindAllStringIndex(code, -1)
		for _, match := range matches {
			startIdx := match[0]
			braceStart := strings.Index(code[startIdx:], "{")
			if braceStart == -1 {
				continue
			}
			braceStart += startIdx

			braceCount := 1
			endIdx := braceStart + 1
			for endIdx < len(code) && braceCount > 0 {
				if code[endIdx] == '{' {
					braceCount++
				} else if code[endIdx] == '}' {
					braceCount--
				}
				endIdx++
			}

			if braceCount == 0 {
				endMarker := endIdx
				for endMarker < len(code) && (code[endMarker] == ';' || code[endMarker] == '\n' || code[endMarker] == ' ' || code[endMarker] == '\t') {
					endMarker++
				}

				classBlock := code[startIdx:endMarker]
				replacement := "// class " + name + " definition removed (already in source)\n"
				code = strings.Replace(code, classBlock, replacement, 1)
			}
		}
	}
	return code
}

func stripStructDefinitions(code string, structNames []string) string {
	for _, name := range structNames {
		pattern := regexp.MustCompile(`(?m)(?://[^\n]*\n)?\s*struct\s+` + name + `\s*\{`)

		matches := pattern.FindAllStringIndex(code, -1)
		for _, match := range matches {
			startIdx := match[0]
			braceStart := strings.Index(code[startIdx:], "{")
			if braceStart == -1 {
				continue
			}
			braceStart += startIdx

			braceCount := 1
			endIdx := braceStart + 1
			for endIdx < len(code) && braceCount > 0 {
				if code[endIdx] == '{' {
					braceCount++
				} else if code[endIdx] == '}' {
					braceCount--
				}
				endIdx++
			}

			if braceCount == 0 {
				endMarker := endIdx
				for endMarker < len(code) && (code[endMarker] == ';' || code[endMarker] == '\n' || code[endMarker] == ' ' || code[endMarker] == '\t') {
					endMarker++
				}

				structBlock := code[startIdx:endMarker]
				replacement := "// struct " + name + " definition removed (already in source)\n"
				code = strings.Replace(code, structBlock, replacement, 1)
			}
		}
	}
	return code
}

func stripConstDefinitions(code string, constNames []string) string {
	for _, name := range constNames {
		constPattern := regexp.MustCompile(`(?m)^\s*(?:const|static\s+const)\s+\w+\s+` + name + `\s*=.*;`)
		code = constPattern.ReplaceAllString(code, "// const "+name+" definition removed (already in source)\n")
	}
	return code
}

func stripTopLevelFuncDefinitions(code string, funcNames []string) string {
	for _, name := range funcNames {
		pattern := regexp.MustCompile(`(?m)(?://[^\n]*\n)?\s*(?:inline\s+)?(?:static\s+)?(?:\w+(?:\s*\*|\s*&)?\s+)+` + name + `\s*\([^)]*\)\s*\{`)

		matches := pattern.FindAllStringIndex(code, -1)
		for _, match := range matches {
			startIdx := match[0]
			braceStart := strings.Index(code[startIdx:], "{")
			if braceStart == -1 {
				continue
			}
			braceStart += startIdx

			braceCount := 1
			endIdx := braceStart + 1
			for endIdx < len(code) && braceCount > 0 {
				if code[endIdx] == '{' {
					braceCount++
				} else if code[endIdx] == '}' {
					braceCount--
				}
				endIdx++
			}

			if braceCount == 0 {
				funcBlock := code[startIdx:endIdx]
				replacement := "// function " + name + " definition removed (already in source)\n"
				code = strings.Replace(code, funcBlock, replacement, 1)
			}
		}
	}
	return code
}

func removeSourceInclude(testContent []byte, sourceBase string) []byte {
	code := string(testContent)

	includePattern := regexp.MustCompile(`(?m)^\s*#include\s*"` + sourceBase + `"\s*\n?`)
	code = includePattern.ReplaceAllString(code, "")

	sourceCppPattern := regexp.MustCompile(`(?m)^\s*#include\s*"source\.cpp"\s*\n?`)
	code = sourceCppPattern.ReplaceAllString(code, "")

	anglePattern := regexp.MustCompile(`(?m)^\s*#include\s*<` + sourceBase + `>\s*\n?`)
	code = anglePattern.ReplaceAllString(code, "")

	return []byte(code)
}
