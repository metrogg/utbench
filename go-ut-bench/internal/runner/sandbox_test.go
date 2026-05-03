package runner

import (
	"path/filepath"
	"testing"
)

func TestResolveDockerWorkspaceMountUsesHostOutputRootOverride(t *testing.T) {
	t.Setenv("UTBENCH_SANDBOX_HOST_OUTPUT_ROOT", "/host/artifacts")
	req := SandboxRunRequest{
		Workspace:           "/app/artifacts/runs/run_1/agent_workspaces/opencode/python/sample_001",
		ContainerOutputRoot: "/app/artifacts",
	}
	got, err := resolveDockerWorkspaceMount(req)
	if err != nil {
		t.Fatalf("resolveDockerWorkspaceMount returned error: %v", err)
	}
	want := filepath.Join("/host/artifacts", "runs", "run_1", "agent_workspaces", "opencode", "python", "sample_001")
	if got != want {
		t.Fatalf("unexpected mount source: got=%s want=%s", got, want)
	}
}

func TestResolveDockerWorkspaceMountFailsFastForContainerLocalPathWithoutHostOverride(t *testing.T) {
	req := SandboxRunRequest{
		Workspace:           "/app/artifacts/runs/run_1/agent_workspaces/opencode/python/sample_001",
		ContainerOutputRoot: "/app/artifacts",
	}
	_, err := resolveDockerWorkspaceMount(req)
	if err == nil {
		t.Fatalf("expected error for container-local workspace without host override")
	}
}
