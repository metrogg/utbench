package evaluator

import (
	"fmt"
	"go-ut-bench/internal/contracts"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParsePytestCounts(t *testing.T) {
	out := "1 failed, 14 passed in 0.09s"
	passed, total := parsePytestCounts(out)
	if passed == nil || total == nil {
		t.Fatalf("expected counts parsed")
	}
	if *passed != 14 || *total != 15 {
		t.Fatalf("unexpected parsed counts: passed=%d total=%d", *passed, *total)
	}
}

func TestParsePytestCountsPassedOnly(t *testing.T) {
	out := "13 passed in 0.05s"
	passed, total := parsePytestCounts(out)
	if passed == nil || total == nil {
		t.Fatalf("expected counts parsed")
	}
	if *passed != 13 || *total != 13 {
		t.Fatalf("unexpected parsed counts: passed=%d total=%d", *passed, *total)
	}
}

func TestParseGoTestCounts(t *testing.T) {
	out := "=== RUN   TestA\n--- PASS: TestA (0.00s)\n=== RUN   TestB\n--- FAIL: TestB (0.00s)\nFAIL\n"
	passed, total := parseGoTestCounts(out)
	if passed == nil || total == nil {
		t.Fatalf("expected counts parsed")
	}
	if *passed != 1 || *total != 2 {
		t.Fatalf("unexpected parsed counts: passed=%d total=%d", *passed, *total)
	}
}

func TestExtractAllClassNamesFromSourceIncludesInterfacesAndEnums(t *testing.T) {
	source := `
interface DataSource {
    String read();
}

class DataProcessor {
}

enum Mode {
    FAST, SLOW
}
`

	names := extractAllClassNamesFromSource(source)
	if len(names) != 3 {
		t.Fatalf("expected 3 types, got %d: %v", len(names), names)
	}
	set := map[string]bool{}
	for _, name := range names {
		set[name] = true
	}
	for _, expected := range []string{"DataSource", "DataProcessor", "Mode"} {
		if !set[expected] {
			t.Fatalf("missing type %s in %v", expected, names)
		}
	}
}

func TestSplitJavaSourceByClassesPreservesHeader(t *testing.T) {
	source := `
import java.util.*;

interface DataSource {
    String read();
}

class DataProcessor {
}
`

	parts := splitJavaSourceByClasses(source)
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}

	if !containsAll(parts["DataSource"], []string{"import java.util.*;", "interface DataSource"}) {
		t.Fatalf("DataSource part missing expected content: %q", parts["DataSource"])
	}
	if !containsAll(parts["DataProcessor"], []string{"import java.util.*;", "class DataProcessor"}) {
		t.Fatalf("DataProcessor part missing expected content: %q", parts["DataProcessor"])
	}
}

func TestGoCompileCheckWithTestFile(t *testing.T) {
	workdir := t.TempDir()
	goMod := `module utbench_eval

go 1.22
`
	if err := os.WriteFile(filepath.Join(workdir, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	source := `package main

func Add(a, b int) int { return a + b }
`
	if err := os.WriteFile(filepath.Join(workdir, "sample.go"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	testCode := `package main

import "testing"

func TestAdd(t *testing.T) {
	if Add(1, 2) != 3 {
		t.Fatalf("unexpected")
	}
}
`
	if err := os.WriteFile(filepath.Join(workdir, "sample_test.go"), []byte(testCode), 0o644); err != nil {
		t.Fatal(err)
	}

	ok, compileErr := goCompileCheck(workdir, "sample_test.go")
	if !ok {
		t.Fatalf("expected compile check pass, got error: %s", compileErr)
	}
}

func TestRewriteGeneratedTestImportsNeutralPythonModules(t *testing.T) {
	sourcePath := filepath.Join("tmp", "actual_module.py")
	input := strings.Join([]string{
		"from module_under_test import task_func",
		"from target_module import helper",
		"from solution import other",
		"import module_under_test",
		"import target_module as target",
	}, "\n")

	got := rewriteGeneratedTestImports(input, sourcePath)
	for _, forbidden := range []string{"module_under_test", "target_module", "solution"} {
		if strings.Contains(got, forbidden) {
			t.Fatalf("expected neutral import %q to be rewritten, got:\n%s", forbidden, got)
		}
	}
	if !strings.Contains(got, "from actual_module import task_func") ||
		!strings.Contains(got, "from actual_module import helper") ||
		!strings.Contains(got, "import actual_module") {
		t.Fatalf("unexpected rewritten imports:\n%s", got)
	}
}

func TestClassifyFailureOriginKeepsMutationToolErrorAsModelWhenTestsFailed(t *testing.T) {
	testPass := false
	rate := 0.5
	row := contracts.EvaluationResult{
		CompilePass:   true,
		TestPass:      &testPass,
		TestPassRate:  &rate,
		MutationError: "go-mutesting parse error: failed to parse mutation output",
	}

	origin, reason := classifyFailureOrigin(row)
	if origin != "model" {
		t.Fatalf("expected model origin for mutation error caused by failing tests, got origin=%q reason=%q", origin, reason)
	}
}

func TestClassifyFailureOriginKeepsPureMutationToolErrorExcluded(t *testing.T) {
	testPass := true
	rate := 1.0
	row := contracts.EvaluationResult{
		CompilePass:   true,
		TestPass:      &testPass,
		TestPassRate:  &rate,
		MutationError: "go-mutesting parse error: failed to parse mutation output",
	}

	origin, reason := classifyFailureOrigin(row)
	if origin != "tool" || reason == "" {
		t.Fatalf("expected pure mutation parser issue to remain tool origin, got origin=%q reason=%q", origin, reason)
	}
}

func TestClassifyFailureOriginTreatsGoMutestingNoResultsAsTool(t *testing.T) {
	testPass := true
	rate := 1.0
	row := contracts.EvaluationResult{
		CompilePass:   true,
		TestPass:      &testPass,
		TestPassRate:  &rate,
		MutationError: "go-mutesting: go-mutesting no results to report",
	}

	origin, reason := classifyFailureOrigin(row)
	if origin != "tool" || reason == "" {
		t.Fatalf("expected go-mutesting no-results to be tool origin, got origin=%q reason=%q", origin, reason)
	}
}

func TestShouldRunMutationAfterSampleTestsRequiresSamplePass(t *testing.T) {
	testPass := false
	rate := 1.0
	row := contracts.EvaluationResult{
		CompilePass:  true,
		TestPass:     &testPass,
		TestPassRate: &rate,
	}

	if shouldRunMutationAfterSampleTests(row) {
		t.Fatal("expected mutation to be skipped when sample-level tests failed even if parsed pass rate is 1.0")
	}
}

func TestPythonTestImportsAnyMutationTarget(t *testing.T) {
	dir := t.TempDir()
	testName := "test_sample.py"
	if err := os.WriteFile(filepath.Join(dir, testName), []byte("from sample import task_func\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !pythonTestImportsAnyMutationTarget(dir, testName, []string{"sample.py"}) {
		t.Fatal("expected direct target import to be detected")
	}
	if pythonTestImportsAnyMutationTarget(dir, testName, []string{"other.py"}) {
		t.Fatal("did not expect unrelated import to be treated as target import")
	}
}

func TestClassifyFailureOriginTreatsMissingPythonTargetImportAsModel(t *testing.T) {
	testPass := true
	rate := 1.0
	row := contracts.EvaluationResult{
		CompilePass:   true,
		TestPass:      &testPass,
		TestPassRate:  &rate,
		MutationError: "mutmut: generated tests do not import mutation target, skipping mutation",
	}

	origin, reason := classifyFailureOrigin(row)
	if origin != "model" || reason != "" {
		t.Fatalf("expected missing target import to be model origin, got origin=%q reason=%q", origin, reason)
	}
}

func TestClassifyFailureOriginTreatsMissingGeneratedCppIncludeAsModel(t *testing.T) {
	row := contracts.EvaluationResult{
		CompilePass:  false,
		CompileError: "fatal error: sorted_list_sum.cpp: No such file or directory",
	}

	origin, reason := classifyFailureOrigin(row)
	if origin != "model" || reason != "" {
		t.Fatalf("expected generated missing include to remain model origin, got origin=%q reason=%q", origin, reason)
	}
}

func TestClassifyFailureOriginMarksPitestJUnit5PluginAsTool(t *testing.T) {
	testPass := true
	rate := 1.0
	row := contracts.EvaluationResult{
		CompilePass:   true,
		TestPass:      &testPass,
		TestPassRate:  &rate,
		MutationError: "PITest JUnit 5 plugin is not installed",
	}

	origin, reason := classifyFailureOrigin(row)
	if origin != "tool" || reason == "" {
		t.Fatalf("expected missing pitest junit 5 plugin to be tool origin, got origin=%q reason=%q", origin, reason)
	}
}

func TestClassifyFailureOriginMarksPitestNoEffectiveResultsAsTool(t *testing.T) {
	testPass := true
	rate := 1.0
	row := contracts.EvaluationResult{
		CompilePass:   true,
		TestPass:      &testPass,
		TestPassRate:  &rate,
		MutationError: "pitest: pitest no killed/survived results",
	}

	origin, reason := classifyFailureOrigin(row)
	if origin != "tool" || reason == "" {
		t.Fatalf("expected pitest no-results to be tool origin, got origin=%q reason=%q", origin, reason)
	}
}

func TestJavaPomTemplateIncludesPitestJUnit5Plugin(t *testing.T) {
	pom := fmt.Sprintf(javaPomTemplate, "Calculator*", "CalculatorTest", 5000)
	for _, want := range []string{
		"<artifactId>pitest-maven</artifactId>",
		"<artifactId>pitest-junit5-plugin</artifactId>",
		"<version>1.2.1</version>",
	} {
		if !strings.Contains(pom, want) {
			t.Fatalf("expected generated pom to contain %q", want)
		}
	}
}

func TestParseGoMutestingOutputNoResultsToReport(t *testing.T) {
	stats, parseErr := parseGoMutestingOutput("No mutants generated\n")
	if parseErr != "go-mutesting no results to report" {
		t.Fatalf("expected explicit no-results error, got stats=%+v err=%q", stats, parseErr)
	}
	if stats.Total != 0 {
		t.Fatalf("expected zero stats for no-results output, got %+v", stats)
	}
}

func TestParseGoMutestingOutputCounts(t *testing.T) {
	output := strings.Join([]string{
		"The mutation score is 0.600000 (3 passed, 2 failed, 4 duplicated, 1 skipped, total is 6)",
	}, "\n")
	stats, parseErr := parseGoMutestingOutput(output)
	if parseErr != "" {
		t.Fatalf("unexpected parse error: %s", parseErr)
	}
	if stats.Killed != 3 || stats.Survived != 2 || stats.Duplicated != 4 || stats.Skipped != 1 || stats.Total != 6 {
		t.Fatalf("unexpected go-mutesting stats: %+v", stats)
	}
}

func TestParseGoMutestingOutputZeroMutantsSummary(t *testing.T) {
	output := "The mutation score is 0.000000 (0 passed, 0 failed, 0 duplicated, 0 skipped, total is 0)"
	stats, parseErr := parseGoMutestingOutput(output)
	if parseErr != "" {
		t.Fatalf("expected zero-mutant summary to parse successfully, got stats=%+v err=%q", stats, parseErr)
	}
	if stats.Total != 0 || stats.Killed != 0 || stats.Survived != 0 {
		t.Fatalf("unexpected zero-mutant stats: %+v", stats)
	}
}

func TestFinalizeEvaluationResultUsesEndToEndRuntime(t *testing.T) {
	initialRuntime := 1
	row := contracts.EvaluationResult{
		CompilePass: true,
		RuntimeMS:   &initialRuntime,
	}

	start := time.Now().Add(-25 * time.Millisecond)
	finalizeEvaluationResult(&row, start)

	if row.RuntimeMS == nil || *row.RuntimeMS < 20 {
		t.Fatalf("expected end-to-end runtime to overwrite stage runtime, got %v", row.RuntimeMS)
	}
}

func containsAll(s string, subs []string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}
