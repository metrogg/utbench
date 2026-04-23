package evaluator

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const defaultTestTimeoutSeconds = 120

const javaPomTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>

    <groupId>utbench</groupId>
    <artifactId>utbench-eval</artifactId>
    <version>1.0.0</version>
    <packaging>jar</packaging>

    <properties>
        <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
        <maven.compiler.source>17</maven.compiler.source>
        <maven.compiler.target>17</maven.compiler.target>
    </properties>

    <dependencies>
        <dependency>
            <groupId>org.junit.jupiter</groupId>
            <artifactId>junit-jupiter</artifactId>
            <version>5.10.2</version>
            <scope>test</scope>
        </dependency>
        <dependency>
            <groupId>org.junit.vintage</groupId>
            <artifactId>junit-vintage-engine</artifactId>
            <version>5.10.2</version>
            <scope>test</scope>
        </dependency>
    </dependencies>

    <build>
        <plugins>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-compiler-plugin</artifactId>
                <version>3.11.0</version>
                <configuration>
                    <source>17</source>
                    <target>17</target>
                </configuration>
            </plugin>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-surefire-plugin</artifactId>
                <version>3.1.2</version>
                <configuration>
                    <!-- 即使测试失败也继续运行，不中断构建 -->
                    <testFailureIgnore>true</testFailureIgnore>
                    <!-- 跳过无法运行的测试类 -->
                    <skipAfterFailureCount>0</skipAfterFailureCount>
                </configuration>
            </plugin>
            <plugin>
                <groupId>org.jacoco</groupId>
                <artifactId>jacoco-maven-plugin</artifactId>
                <version>0.8.11</version>
                <executions>
                    <execution>
                        <id>prepare-agent</id>
                        <goals>
                            <goal>prepare-agent</goal>
                        </goals>
                    </execution>
                    <execution>
                        <id>report</id>
                        <phase>test</phase>
                        <goals>
                            <goal>report</goal>
                        </goals>
                    </execution>
                </executions>
            </plugin>
            <plugin>
                <groupId>org.pitest</groupId>
                <artifactId>pitest-maven</artifactId>
                <version>1.19.6</version>
                <configuration>
                    <targetClasses>%s</targetClasses>
                    <targetTests>%s</targetTests>
                    <outputFormats>XML,CSV</outputFormats>
                    <mutators>
                        <mutator>CONDITIONALS_BOUNDARY</mutator>
                        <mutator>INCREMENTS</mutator>
                        <mutator>INVERT_NEGS</mutator>
                        <mutator>MATH</mutator>
                        <mutator>NEGATE_CONDITIONALS</mutator>
                        <mutator>EMPTY_RETURNS</mutator>
                        <mutator>FALSE_RETURNS</mutator>
                        <mutator>TRUE_RETURNS</mutator>
                        <mutator>NULL_RETURNS</mutator>
                        <mutator>PRIMITIVE_RETURNS</mutator>
                        <mutator>VOID_METHOD_CALLS</mutator>
                    </mutators>
                    <timeoutFactor>1.5</timeoutFactor>
                    <timeoutConstant>%d</timeoutConstant>
                    <threads>1</threads>
                    <!-- 关键配置：即使基线测试有失败也继续运行变异 -->
                    <failWhenNoMutations>false</failWhenNoMutations>
                    <skipFailingTests>true</skipFailingTests>
                    <excludedTestMethods>
                        <!-- 自动排除失败的测试方法 -->
                        <excludedTestMethod>*#*fail*</excludedTestMethod>
                        <excludedTestMethod>*#*error*</excludedTestMethod>
                    </excludedTestMethods>
                    <!-- 变异测试运行失败也不中断，继续生成报告 -->
                    <failOnError>false</failOnError>
                    <failWhenNoCoverage>false</failWhenNoCoverage>
                </configuration>
            </plugin>
        </plugins>
    </build>
</project>
`

func prepareJavaWorkspace(testPath, samplePath string) (string, string, string, string) {
	testSource, err := os.ReadFile(testPath)
	if err != nil {
		return "", "", "", fmt.Sprintf("failed to read generated test: %s", err)
	}

	sourceBase := filepath.Base(samplePath)
	sourceStem := strings.TrimSuffix(sourceBase, filepath.Ext(sourceBase))

	sourceData, err := os.ReadFile(samplePath)
	if err != nil {
		return "", "", "", fmt.Sprintf("failed to read source: %s", err)
	}

	classNames := extractAllClassNamesFromSource(string(sourceData))
	if len(classNames) == 0 {
		classNames = []string{sourceStem}
	}

	primaryClassName := classNames[0]

	testClassName := extractTestClassNameFromTest(string(testSource))
	if testClassName == "" {
		testClassName = primaryClassName + "Test"
	}
	testFileName := testClassName + ".java"

	workdir, err := os.MkdirTemp("", "utbench_java_eval_")
	if err != nil {
		return "", "", "", fmt.Sprintf("failed to create temp dir: %s", err)
	}

	srcMainJava := filepath.Join(workdir, "src", "main", "java")
	srcTestJava := filepath.Join(workdir, "src", "test", "java")
	if err := os.MkdirAll(srcMainJava, 0o755); err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", fmt.Sprintf("failed to create src/main/java: %s", err)
	}
	if err := os.MkdirAll(srcTestJava, 0o755); err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", fmt.Sprintf("failed to create src/test/java: %s", err)
	}

	if len(classNames) == 1 {
		sourceDestPath := filepath.Join(srcMainJava, primaryClassName+".java")
		if err := os.WriteFile(sourceDestPath, sourceData, 0o644); err != nil {
			_ = os.RemoveAll(workdir)
			return "", "", "", fmt.Sprintf("failed to write source: %s", err)
		}
	} else {
		perClassSources := splitJavaSourceByClasses(string(sourceData))
		for className, classSource := range perClassSources {
			classFileName := className + ".java"
			if err := os.WriteFile(filepath.Join(srcMainJava, classFileName), []byte(classSource), 0o644); err != nil {
				_ = os.RemoveAll(workdir)
				return "", "", "", fmt.Sprintf("failed to write class %s: %s", classFileName, err)
			}
		}
	}

	testDestPath := filepath.Join(srcTestJava, testFileName)
	if err := os.WriteFile(testDestPath, testSource, 0o644); err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", fmt.Sprintf("failed to write test: %s", err)
	}

	pomContent := fmt.Sprintf(javaPomTemplate, primaryClassName+"*", classNames[0]+"Test", 5000)
	pomPath := filepath.Join(workdir, "pom.xml")
	if err := os.WriteFile(pomPath, []byte(pomContent), 0o644); err != nil {
		_ = os.RemoveAll(workdir)
		return "", "", "", fmt.Sprintf("failed to write pom.xml: %s", err)
	}

	return workdir, testFileName, sourceBase, primaryClassName
}

func extractJavaClassName(filename string) string {
	stem := strings.TrimSuffix(filename, filepath.Ext(filename))
	return stem
}

func extractClassNameFromSource(source string) string {
	classPattern := regexp.MustCompile(`(?:public\s+|private\s+|protected\s+)?class\s+([A-Za-z_][A-Za-z0-9_]*)`)
	match := classPattern.FindStringSubmatch(source)
	if match != nil && len(match) > 1 {
		return match[1]
	}
	return ""
}

func extractTestClassNameFromTest(test string) string {
	classPattern := regexp.MustCompile(`(?:public\s+)?class\s+([A-Za-z_][A-Za-z0-9_]*Test)`)
	match := classPattern.FindStringSubmatch(test)
	if match != nil && len(match) > 1 {
		return match[1]
	}
	return ""
}

func extractAllClassNamesFromSource(source string) []string {
	typePattern := regexp.MustCompile(`(?m)(?:^|\n)\s*(?:public\s+|private\s+|protected\s+)?(?:abstract\s+|final\s+|static\s+)*\b(?:class|interface|enum)\s+([A-Za-z_][A-Za-z0-9_]*)\b`)
	seen := make(map[string]bool)
	out := make([]string, 0)
	for _, match := range typePattern.FindAllStringSubmatch(source, -1) {
		if len(match) < 2 {
			continue
		}
		name := match[1]
		if seen[name] {
			continue
		}
		seen[name] = true
		out = append(out, name)
	}
	return out
}

func splitJavaSourceByClasses(source string) map[string]string {
	typePattern := regexp.MustCompile(`(?m)(?:^|\n)\s*(?:public\s+|private\s+|protected\s+)?(?:abstract\s+|final\s+|static\s+)*\b(?:class|interface|enum)\s+([A-Za-z_][A-Za-z0-9_]*)\b`)
	decls := typePattern.FindAllStringSubmatchIndex(source, -1)
	if len(decls) == 0 {
		return map[string]string{}
	}

	firstStart := decls[0][0]
	header := strings.TrimSpace(source[:firstStart])

	result := make(map[string]string)
	for i, match := range decls {
		if len(match) < 4 {
			continue
		}
		className := source[match[2]:match[3]]
		start := match[0]
		for start < len(source) && (source[start] == '\n' || source[start] == '\r') {
			start++
		}
		end := len(source)
		if i+1 < len(decls) {
			end = decls[i+1][0]
		}
		classSource := strings.TrimSpace(source[start:end])
		if header != "" {
			classSource = header + "\n\n" + classSource
		}
		result[className] = classSource
	}
	return result
}

func javaCompileCheck(workdir string) (bool, string) {
	cmd := exec.Command("mvn", "test-compile", "-q")
	cmd.Dir = workdir
	output, err := cmd.CombinedOutput()
	if err == nil {
		return true, ""
	}
	return false, trimErr(string(output), 2000)
}

func executeJavaTests(workdir string) (bool, string, int) {
	return executeJavaTestsWithTimeout(workdir, defaultTestTimeoutSeconds)
}

func executeJavaTestsWithTimeout(workdir string, timeoutSeconds int) (bool, string, int) {
	if timeoutSeconds <= 0 {
		timeoutSeconds = defaultTestTimeoutSeconds
	}
	cmd := exec.Command("mvn", "test", "-q")
	cmd.Dir = workdir
	started := time.Now()
	done := make(chan error, 1)
	var output []byte
	var err error
	go func() {
		output, err = cmd.CombinedOutput()
		done <- nil
	}()
	select {
	case <-done:
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		cmd.Process.Kill()
		return false, "test execution timed out", int(time.Since(started).Milliseconds())
	}
	latency := int(time.Since(started).Milliseconds())
	if err == nil {
		return true, string(output), latency
	}
	return false, trimErr(string(output), 4000), latency
}

func stripANSICodes(s string) string {
	return regexp.MustCompile(`\x1b\[[0-9;]*m`).ReplaceAllString(s, "")
}

func parseJavaTestCounts(output string) (*int, *int) {
	clean := stripANSICodes(output)

	// Match: Tests run: X, Failures: Y, Errors: Z
	// Note: Maven may output either "Failures" or "Errors" or both
	passedPattern := regexp.MustCompile(`Tests run:\s*(\d+),\s*Failures:\s*(\d+)(?:,\s*Errors:\s*(\d+))?`)
	match := passedPattern.FindStringSubmatch(clean)
	if match != nil && len(match) >= 3 {
		totalRuns := parseIntOrZero(match[1])
		failures := parseIntOrZero(match[2])
		errors := 0
		if len(match) >= 4 && match[3] != "" {
			errors = parseIntOrZero(match[3])
		}
		totalFailures := failures + errors
		passed := totalRuns - totalFailures
		if passed < 0 {
			passed = 0
		}
		return &passed, &totalRuns
	}

	if strings.Contains(clean, "BUILD SUCCESS") {
		passed := 1
		total := 1
		return &passed, &total
	}

	return nil, nil
}

func parseIntOrZero(s string) int {
	var result int
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			result = result*10 + int(ch-'0')
		}
	}
	return result
}

func collectJavaCoverage(workdir, className string) (float64, float64, string) {
	jacocoXML := filepath.Join(workdir, "target", "site", "jacoco", "jacoco.xml")
	if _, err := os.Stat(jacocoXML); err != nil {
		return 0, 0, "jacoco.xml not found"
	}

	raw, err := os.ReadFile(jacocoXML)
	if err != nil {
		return 0, 0, fmt.Sprintf("failed to read jacoco.xml: %s", err)
	}

	return parseJacocoXML(string(raw), className)
}

type JacocoReport struct {
	XMLName xml.Name        `xml:"report"`
	Package []JacocoPackage `xml:"package"`
}

type JacocoPackage struct {
	Name  string        `xml:"name,attr"`
	Class []JacocoClass `xml:"class"`
}

type JacocoClass struct {
	Name    string          `xml:"name,attr"`
	Counter []JacocoCounter `xml:"counter"`
}

type JacocoCounter struct {
	Type    string `xml:"type,attr"`
	Missed  int    `xml:"missed,attr"`
	Covered int    `xml:"covered,attr"`
}

func parseJacocoXML(content, className string) (float64, float64, string) {
	var report JacocoReport
	if err := xml.Unmarshal([]byte(content), &report); err != nil {
		return 0, 0, fmt.Sprintf("failed to parse jacoco.xml: %s", err)
	}

	for _, pkg := range report.Package {
		for _, cls := range pkg.Class {
			if strings.Contains(cls.Name, className) || cls.Name == className {
				lineMissed := 0
				lineCovered := 0
				branchMissed := 0
				branchCovered := 0

				for _, cnt := range cls.Counter {
					switch cnt.Type {
					case "LINE":
						lineMissed = cnt.Missed
						lineCovered = cnt.Covered
					case "BRANCH":
						branchMissed = cnt.Missed
						branchCovered = cnt.Covered
					}
				}

				totalLines := lineCovered + lineMissed
				if totalLines == 0 {
					return 0, 0, "no line coverage data"
				}

				lineCov := round(float64(lineCovered)/float64(totalLines), 6)

				var branchCov float64
				totalBranches := branchCovered + branchMissed
				if totalBranches > 0 {
					branchCov = round(float64(branchCovered)/float64(totalBranches), 6)
				}

				return lineCov, branchCov, ""
			}
		}
	}

	return 0, 0, "class not found in coverage report"
}

func collectJavaMutation(ctx context.Context, workdir, className string, timeoutSeconds int, testPassRate *float64, testPassed, testTotal int) (float64, mutationStats, string) {
	if timeoutSeconds <= 0 {
		timeoutSeconds = 120
	}

	minPassRate := GetMinPassRateForTool("pitest")
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

	checkResult := CheckTestPassRate(passed, total, "PITest", minPassRate)
	if !checkResult.ShouldRun {
		return 0, mutationStats{}, checkResult.Message
	}

	runCtx, cancelRun := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
	defer cancelRun()

	runOut, runErr := runCommandWithProcessGroupKill(runCtx, "mvn", []string{"org.pitest:pitest-maven:mutationCoverage", "-q"}, workdir, nil)

	stats, parseErr := parsePitXML(workdir)
	if parseErr != "" {
		return 0, stats, formatMutationError("pitest parse error", runErr, runOut, nil, nil)
	}

	if stats.Total <= 0 {
		return 0, stats, formatMutationError("pitest produced zero mutants", runErr, runOut, nil, nil)
	}

	processed := stats.Killed + stats.Survived + stats.NoTests + stats.Timeout + stats.Skipped + stats.Suspicious
	if processed <= 0 {
		return 0, stats, formatMutationError("pitest did not execute any mutants", runErr, runOut, nil, nil)
	}

	if stats.Killed+stats.Survived <= 0 {
		return 0, stats, formatMutationError("pitest no killed/survived results", runErr, runOut, nil, nil)
	}

	score := round(float64(stats.Killed)/float64(stats.Killed+stats.Survived), 6)
	return score, stats, ""
}

type PitMutationSummary struct {
	XMLName             xml.Name      `xml:"mutations"`
	Mutations           []PitMutation `xml:"mutation"`
	MutationsTotal      int           `xml:"mutationsTotal,attr"`
	MutationsKilled     int           `xml:"mutationsKilled,attr"`
	MutationsSurvived   int           `xml:"mutationsSurvived,attr"`
	MutationsNoCoverage int           `xml:"mutationsNoCoverage,attr"`
	MutationsTimedOut   int           `xml:"mutationsTimedOut,attr"`
	MutationsSkipped    int           `xml:"mutationsSkipped,attr"`
}

type PitMutation struct {
	Status string `xml:"status,attr"`
}

func parsePitXML(workdir string) (mutationStats, string) {
	pitXML := filepath.Join(workdir, "target", "pit-reports", "mutations.xml")
	if _, err := os.Stat(pitXML); err != nil {
		pitXML = filepath.Join(workdir, "target", "pit-reports", "index.html")
		if _, err := os.Stat(pitXML); err != nil {
			csvPath := findPitCSV(workdir)
			if csvPath != "" {
				return parsePitCSV(csvPath)
			}
			return mutationStats{}, "pitest report not found"
		}
		return parsePitHTMLSummary(pitXML)
	}

	raw, err := os.ReadFile(pitXML)
	if err != nil {
		return mutationStats{}, fmt.Sprintf("failed to read pitest report: %s", err)
	}

	var summary PitMutationSummary
	if err := xml.Unmarshal([]byte(raw), &summary); err != nil {
		return mutationStats{}, fmt.Sprintf("failed to parse pitest report: %s", err)
	}

	stats := mutationStats{
		Total:    summary.MutationsTotal,
		Killed:   summary.MutationsKilled,
		Survived: summary.MutationsSurvived,
		NoTests:  summary.MutationsNoCoverage,
		Timeout:  summary.MutationsTimedOut,
		Skipped:  summary.MutationsSkipped,
	}

	if stats.Total == 0 {
		stats.Total = len(summary.Mutations)
		for _, m := range summary.Mutations {
			switch strings.ToUpper(m.Status) {
			case "KILLED":
				stats.Killed++
			case "SURVIVED":
				stats.Survived++
			case "NO_COVERAGE":
				stats.NoTests++
			case "TIMED_OUT":
				stats.Timeout++
			case "SKIPPED":
				stats.Skipped++
			}
		}
	}

	return stats, ""
}

func parsePitHTMLSummary(path string) (mutationStats, string) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return mutationStats{}, fmt.Sprintf("failed to read pitest html: %s", err)
	}

	content := string(raw)

	stats := mutationStats{}
	stats.Total = extractFirstIntOrZero(content, `(\d+)\s*mutations`)
	stats.Killed = extractFirstIntOrZero(content, `(\d+)\s*killed`)
	stats.Survived = extractFirstIntOrZero(content, `(\d+)\s*survived`)
	stats.NoTests = extractFirstIntOrZero(content, `(\d+)\s*no coverage`)

	if stats.Total == 0 {
		return stats, "no mutation stats found in html"
	}

	return stats, ""
}

func findPitCSV(workdir string) string {
	pitDir := filepath.Join(workdir, "target", "pit-reports")
	files, err := os.ReadDir(pitDir)
	if err != nil {
		return ""
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".csv") {
			return filepath.Join(pitDir, f.Name())
		}
	}
	return ""
}

func parsePitCSV(path string) (mutationStats, string) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return mutationStats{}, fmt.Sprintf("failed to read pitest csv: %s", err)
	}

	lines := strings.Split(string(raw), "\n")
	stats := mutationStats{}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(strings.TrimSpace(line), "mutation") {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) < 2 {
			continue
		}

		stats.Total++
		status := strings.TrimSpace(strings.ToUpper(fields[1]))
		switch status {
		case "KILLED":
			stats.Killed++
		case "SURVIVED":
			stats.Survived++
		case "NO_COVERAGE":
			stats.NoTests++
		case "TIMED_OUT":
			stats.Timeout++
		case "SKIPPED":
			stats.Skipped++
		}
	}

	if stats.Total == 0 {
		return stats, "no mutations in csv"
	}

	return stats, ""
}
