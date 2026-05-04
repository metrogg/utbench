package web

import (
	"strings"
	"testing"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/orchestrator"
)

func TestBuildDockerRunArgsUsesMountedSource(t *testing.T) {
	spec := contracts.RunSpec{
		RunID:     "run-1",
		Models:    []string{"deepseek"},
		Languages: []string{"python"},
	}
	cfg := DockerConfig{EvalImageName: "utbench:latest", ProjectRoot: "/repo"}

	args := buildDockerRunArgs(spec, orchestrator.Options{}, cfg)
	joined := strings.Join(args, " ")

	mustContain(t, joined, "-v /repo/datasets:/app/datasets")
	mustContain(t, joined, "-v /repo/artifacts:/app/artifacts")
	mustContain(t, joined, "utbench:latest")
	mustContain(t, joined, "--output-root /app/artifacts")
	mustContain(t, joined, "--models deepseek")
}

func TestBuildDockerEvaluateArgsUsesSourceRunManifest(t *testing.T) {
	spec := contracts.RunSpec{
		RunID:           "run-1",
		MutationEnabled: true,
		MutationTimeout: 120,
		MutationPolicy:  "warn",
	}
	cfg := DockerConfig{EvalImageName: "utbench:latest", ProjectRoot: "/repo"}

	args := buildDockerRunArgs(spec, orchestrator.Options{
		Phase:       "evaluate",
		SourceRunID: "source-run",
	}, cfg)
	joined := strings.Join(args, " ")

	mustContain(t, joined, "--run-id run-1")
	mustContain(t, joined, "--manifest /app/artifacts/runs/source-run/generated/generated_manifest.json")
	mustContain(t, joined, "--mutation-enabled")
	mustContain(t, joined, "--mutation-timeout 120")
}

func mustContain(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected %q to contain %q", haystack, needle)
	}
}
