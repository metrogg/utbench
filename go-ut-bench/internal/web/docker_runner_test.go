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
	cfg := DockerConfig{ImageName: "utbench:latest", ProjectRoot: "/repo"}

	args := buildDockerRunArgs(spec, orchestrator.Options{}, cfg)
	joined := strings.Join(args, " ")

	mustContain(t, joined, "-v /repo:/workspace")
	mustContain(t, joined, "-w /workspace")
	mustContain(t, joined, "--entrypoint /bin/sh")
	mustContain(t, joined, "utbench:latest -lc 'go' 'run' './cmd/utbench' 'run'")
	mustContain(t, joined, "--output-root' '/app/artifacts")
}

func TestBuildDockerEvaluateArgsUsesMountedSource(t *testing.T) {
	spec := contracts.RunSpec{
		RunID:           "run-1",
		MutationEnabled: true,
		MutationTimeout: 120,
		MutationPolicy:  "warn",
	}
	cfg := DockerConfig{ImageName: "utbench:latest", ProjectRoot: "/repo"}

	args := buildDockerSourceArgs(cfg, []string{
		"go", "run", "./cmd/utbench", "evaluate",
		"--run-id", spec.RunID,
		"--manifest", "/app/artifacts/runs/" + spec.RunID + "/generated/generated_manifest.json",
	})
	joined := strings.Join(args, " ")

	mustContain(t, joined, "-v /repo:/workspace")
	mustContain(t, joined, "utbench:latest -lc 'go' 'run' './cmd/utbench' 'evaluate'")
	mustContain(t, joined, "--manifest' '/app/artifacts/runs/run-1/generated/generated_manifest.json")
}

func mustContain(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected %q to contain %q", haystack, needle)
	}
}
