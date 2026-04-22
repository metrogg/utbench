package evaluator

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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

const cppMullConfigTemplate = `mutators:
  - cxx_arithmetic
  - cxx_comparison
  - cxx_boundary
  - cxx_bitwise
  - cxx_logical
  - cxx_increment_decrement
  - cxx_remove_void_call

excludePaths:
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

set(CMAKE_C_COMPILER clang-15)
set(CMAKE_CXX_COMPILER clang++-15)

find_package(GTest REQUIRED)

enable_testing()

add_executable(test_runner_mull
    %s
    %s
)

target_link_libraries(test_runner_mull GTest::gtest_main)

target_compile_options(test_runner_mull PRIVATE
    -fpass-plugin=%s
    -g -grecord-command-line
    -fPIC
    -O0
)

add_test(NAME AllTests COMMAND test_runner_mull)
`

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
		headerPath := filepath.Join(workdir, headerName)
		if _, err := os.Stat(headerPath); os.IsNotExist(err) {
			if strings.HasSuffix(headerName, ".h") || strings.HasSuffix(headerName, ".hpp") {
				declHeader := generateDeclarationsHeader(string(sourceData), strings.TrimSuffix(headerName, filepath.Ext(headerName)))
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

func cppCompileCheck(workdir string) (bool, string) {
	buildDir := filepath.Join(workdir, "build")

	cmakeCmd := exec.Command("cmake", "..")
	cmakeCmd.Dir = buildDir
	cmakeOut, cmakeErr := cmakeCmd.CombinedOutput()
	if cmakeErr != nil {
		return false, trimErr(string(cmakeOut), 2000)
	}

	makeCmd := exec.Command("make", "-j2")
	makeCmd.Dir = buildDir
	makeOut, makeErr := makeCmd.CombinedOutput()
	if makeErr != nil {
		return false, trimErr(string(makeOut), 2000)
	}

	return true, ""
}

func executeCppTests(workdir string) (bool, string, int) {
	buildDir := filepath.Join(workdir, "build")

	cmd := exec.Command("./test_runner")
	cmd.Dir = buildDir
	started := time.Now()
	output, err := cmd.CombinedOutput()
	latency := int(time.Since(started).Milliseconds())
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
			failedPattern2 := regexp.MustCompile(`(\d+)\s+FAILED`)
			failedMatch := failedPattern2.FindStringSubmatch(output)
			failedCount := 0
			if failedMatch != nil && len(failedMatch) > 1 {
				failedCount = parseIntOrZero(failedMatch[1])
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

	gcovCmd := exec.Command("gcov", "-b", testFileName)
	gcovCmd.Dir = gcovDir
	gcovOut, gcovErr := gcovCmd.CombinedOutput()
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

func collectCppMutation(ctx context.Context, workdir, sourceBase string, timeoutSeconds int) (float64, mutationStats, string) {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 120
	}

	mullRunner := findMullRunner()
	if mullRunner == "" {
		return 0, mutationStats{}, "Mull not installed. Install: sudo apt-get install -y llvm-15 clang-15 mull-15"
	}

	// Find mull-ir-frontend plugin path
	mullFrontend := findMullFrontend()
	if mullFrontend == "" {
		return 0, mutationStats{}, "mull-ir-frontend-15 not found"
	}

	mullBuildDir := filepath.Join(workdir, "build_mull")
	if err := os.MkdirAll(mullBuildDir, 0755); err != nil {
		return 0, mutationStats{}, "failed to create mull build dir: " + err.Error()
	}

	testFileName := strings.TrimSuffix(sourceBase, filepath.Ext(sourceBase)) + "_test.cpp"
	
	// Write mull.yml config
	mullConfigPath := filepath.Join(mullBuildDir, "mull.yml")
	if err := os.WriteFile(mullConfigPath, []byte(cppMullConfigTemplate), 0644); err != nil {
		return 0, mutationStats{}, "failed to write mull.yml: " + err.Error()
	}
	
	// Generate CMakeLists.txt with correct plugin path
	cmakeContent := fmt.Sprintf(cppMullCMakeTemplate, sourceBase, testFileName, mullFrontend)

	mullCmakePath := filepath.Join(mullBuildDir, "CMakeLists.txt")
	if err := os.WriteFile(mullCmakePath, []byte(cmakeContent), 0644); err != nil {
		return 0, mutationStats{}, "failed to write Mull CMakeLists: " + err.Error()
	}

	sourceCopy := filepath.Join(mullBuildDir, sourceBase)
	testCopy := filepath.Join(mullBuildDir, testFileName)

	if err := copyFile(filepath.Join(workdir, sourceBase), sourceCopy); err != nil {
		return 0, mutationStats{}, "failed to copy source: " + err.Error()
	}
	if err := copyFile(filepath.Join(workdir, testFileName), testCopy); err != nil {
		return 0, mutationStats{}, "failed to copy test: " + err.Error()
	}

	cmakeCmd := exec.Command("cmake", ".")
	cmakeCmd.Dir = mullBuildDir
	cmakeCmd.Env = append(os.Environ(), "CC=clang-15", "CXX=clang++-15")
	cmakeOut, cmakeErr := cmakeCmd.CombinedOutput()
	if cmakeErr != nil {
		return 0, mutationStats{}, trimErr("cmake failed: "+string(cmakeOut), 1000)
	}

	makeCmd := exec.Command("make", "-j2")
	makeCmd.Dir = mullBuildDir
	makeCmd.Env = append(os.Environ(), "CC=clang-15", "CXX=clang++-15")
	makeOut, makeErr := makeCmd.CombinedOutput()
	if makeErr != nil {
		return 0, mutationStats{}, trimErr("make failed: "+string(makeOut), 1000)
	}

	execPath := filepath.Join(mullBuildDir, "test_runner_mull")
	if _, err := os.Stat(execPath); err != nil {
		return 0, mutationStats{}, "mull executable not found after build"
	}

	runCtx, cancelRun := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancelRun()

	mullCmd := exec.CommandContext(runCtx, mullRunner, execPath)
	mullCmd.Dir = mullBuildDir
	mullCmd.Env = append(os.Environ(),
		"LD_LIBRARY_PATH=/usr/lib/llvm-15/lib:/usr/lib/x86_64-linux-gnu:/usr/lib64:"+os.Getenv("LD_LIBRARY_PATH"),
	)
	mullOut, mullErr := mullCmd.CombinedOutput()

	stats, parseErr := parseMullOutput(string(mullOut))
	if parseErr != "" {
		return 0, stats, formatMutationError("mull parse error", mullErr, mullOut, nil, nil)
	}

	if stats.Total <= 0 {
		return 0, stats, formatMutationError("mull produced zero mutants", mullErr, mullOut, nil, nil)
	}

	processed := stats.Killed + stats.Survived + stats.NoTests + stats.Timeout + stats.Skipped
	if processed <= 0 {
		return 0, stats, formatMutationError("mull did not execute any mutants", mullErr, mullOut, nil, nil)
	}

	if stats.Killed+stats.Survived <= 0 {
		return 0, stats, formatMutationError("mull no killed/survived results", mullErr, mullOut, nil, nil)
	}

	score := round(float64(stats.Killed)/float64(stats.Killed+stats.Survived), 6)
	return score, stats, ""
}

func findMullRunner() string {
	candidates := []string{"mull-runner-15", "mull-runner-14", "mull-runner"}
	for _, c := range candidates {
		if path, err := exec.LookPath(c); err == nil {
			return path
		}
	}
	return ""
}

func findMullFrontend() string {
	// Try common paths for mull-ir-frontend
	candidates := []string{
		"/usr/lib/mull-ir-frontend-15",
		"/usr/lib/llvm-15/lib/mull-ir-frontend-15.so",
		"/usr/lib/x86_64-linux-gnu/mull-ir-frontend-15.so",
		"/usr/local/lib/mull-ir-frontend-15.so",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

func parseMullOutput(output string) (mutationStats, string) {
	stats := mutationStats{}

	if strings.Contains(output, "Original test failed (warmup run)") {
		return stats, "mull warmup run failed: baseline test has failing assertions"
	}

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

	killedPattern := regexp.MustCompile(`Killed mutants\s*\((\d+)/(\d+)\)`)
	if match := killedPattern.FindStringSubmatch(output); match != nil && len(match) > 2 {
		stats.Killed = parseIntOrZero(match[1])
		stats.Total = parseIntOrZero(match[2])
	}

	survivedPattern := regexp.MustCompile(`Survived mutants\s*\((\d+)/(\d+)\)`)
	if match := survivedPattern.FindStringSubmatch(output); match != nil && len(match) > 1 {
		stats.Survived = parseIntOrZero(match[1])
	}

	if stats.Total == 0 {
		killedCount := strings.Count(output, "Killed:")
		survivedCount := strings.Count(output, "Survived:")
		if killedCount > 0 || survivedCount > 0 {
			stats.Killed = killedCount
			stats.Survived = survivedCount
			stats.Total = killedCount + survivedCount
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

func estimateCppAssertionDensity(testPath string) (int, int, float64) {
	raw, err := os.ReadFile(testPath)
	if err != nil {
		return 0, 0, 0
	}
	text := string(raw)

	assertCount := strings.Count(text, "EXPECT_") + strings.Count(text, "ASSERT_")

	testPattern := regexp.MustCompile(`TEST\s*\(\s*[^,]+\s*,\s*[^)]+\s*\)`)
	testMatches := testPattern.FindAllString(text, -1)
	testCount := len(testMatches)

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
