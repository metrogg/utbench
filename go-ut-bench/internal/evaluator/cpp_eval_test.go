package evaluator

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestCppMutationIntegration(t *testing.T) {
	if _, err := exec.LookPath("mull-runner-15"); err != nil {
		t.Skip("mull-runner-15 not installed, skipping mutation test")
	}

	workdir := filepath.Join(os.Getenv("HOME"), "cpp_mutation_workspace")
	if err := os.RemoveAll(workdir); err != nil && !os.IsNotExist(err) {
		t.Fatalf("failed to clean workdir: %v", err)
	}

	sourceFile := filepath.Join(workdir, "sample.cpp")
	testFile := filepath.Join(workdir, "sample_test.cpp")

	os.MkdirAll(workdir, 0755)

	os.WriteFile(sourceFile, []byte("int add(int a, int b) { return a + b; }\nint multiply(int a, int b) { return a * b; }\n"), 0644)

	os.WriteFile(testFile, []byte("#include <gtest/gtest.h>\nint add(int a, int b);\nint multiply(int a, int b);\n\nTEST(AddTest, Basic) {\n    EXPECT_EQ(add(2, 3), 5);\n}\n\nTEST(MultiplyTest, Basic) {\n    EXPECT_EQ(multiply(2, 3), 6);\n}\n"), 0644)

	prepWorkdir, testName, sourceBase, _, prepErr := prepareCppWorkspace(testFile, sourceFile)
	if prepErr != "" {
		t.Fatalf("prepareCppWorkspace failed: %s", prepErr)
	}
	defer os.RemoveAll(prepWorkdir)

	t.Logf("workdir: %s, test: %s, source: %s", prepWorkdir, testName, sourceBase)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	passRate := 1.0
	score, stats, mutErr := collectCppMutation(ctx, prepWorkdir, sourceBase, 60, &passRate, 2, 2)
	if mutErr != "" {
		t.Logf("mutation error: %s", mutErr)
	}
	t.Logf("mutation score: %.2f%%, total: %d, killed: %d, survived: %d", score*100, stats.Total, stats.Killed, stats.Survived)

	if stats.Total > 0 {
		t.Logf("Mutation testing SUCCESS: found %d mutants", stats.Total)
		if stats.Killed > 0 {
			t.Logf("Killed %d mutants (%.2f%%)", stats.Killed, score*100)
		}
	}

	os.RemoveAll(workdir)
}