package runner

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeSandboxRunner struct {
	t             *testing.T
	workspace     string
	exportPayload string
	commands      []string
}

func (f *fakeSandboxRunner) Run(_ context.Context, req SandboxRunRequest) (SandboxRunResult, error) {
	f.commands = append(f.commands, req.Command)
	if !strings.Contains(req.Command, "opencode export") {
		f.t.Fatalf("unexpected sandbox command: %s", req.Command)
	}
	exportPath := filepath.Join(f.workspace, ".utbench", "opencode", "session_export.json")
	if err := os.MkdirAll(filepath.Dir(exportPath), 0o755); err != nil {
		f.t.Fatalf("mkdir export path: %v", err)
	}
	if err := os.WriteFile(exportPath, []byte(f.exportPayload), 0o644); err != nil {
		f.t.Fatalf("write export payload: %v", err)
	}
	return SandboxRunResult{ExitCode: 0}, nil
}

func TestCollectOpenCodeSessionExportUsesExportedUsage(t *testing.T) {
	workRoot := t.TempDir()
	traceDir := filepath.Join(t.TempDir(), "trace")
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		t.Fatalf("mkdir trace dir: %v", err)
	}

	runner := &fakeSandboxRunner{
		t:         t,
		workspace: workRoot,
		exportPayload: `{
			"session": {
				"id": "ses_test",
				"messages": [
					{"role":"assistant","usage":{"input_tokens":1200,"output_tokens":300,"total_tokens":1500}},
					{"role":"assistant","usage":{"input_tokens":800,"output_tokens":200,"total_tokens":1000}}
				]
			}
		}`,
	}

	trace := AgentTrace{
		Framework: "opencode",
		SessionID: "ses_test",
	}

	collectOpenCodeSessionExport(
		context.Background(),
		runner,
		SandboxRunRequest{Workspace: workRoot, DockerImage: "utbench-agent-opencode-go:latest"},
		workRoot,
		traceDir,
		"boundary_000",
		&trace,
	)

	if trace.SessionExportError != "" {
		t.Fatalf("unexpected session export error: %s", trace.SessionExportError)
	}
	if trace.SessionExportPath == "" {
		t.Fatalf("expected session export path to be recorded")
	}
	if got, want := trace.UsageSourceDetail, "opencode_session_export"; got != want {
		t.Fatalf("usage source detail = %q, want %q", got, want)
	}
	if got, want := trace.TokenSource, "actual"; got != want {
		t.Fatalf("token source = %q, want %q", got, want)
	}
	if trace.PromptTokens == nil || *trace.PromptTokens != 2000 {
		t.Fatalf("prompt tokens = %v, want 2000", trace.PromptTokens)
	}
	if trace.CompletionTokens == nil || *trace.CompletionTokens != 500 {
		t.Fatalf("completion tokens = %v, want 500", trace.CompletionTokens)
	}
	if trace.TotalTokens == nil || *trace.TotalTokens != 2500 {
		t.Fatalf("total tokens = %v, want 2500", trace.TotalTokens)
	}
	if len(runner.commands) != 1 {
		t.Fatalf("expected 1 export command, got %d", len(runner.commands))
	}
	if _, err := os.Stat(trace.SessionExportPath); err != nil {
		t.Fatalf("session export file not persisted: %v", err)
	}
}

func TestCollectOpenCodeSessionExportUsesOpenCodeMessageTokensSchema(t *testing.T) {
	workRoot := t.TempDir()
	traceDir := filepath.Join(t.TempDir(), "trace")
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		t.Fatalf("mkdir trace dir: %v", err)
	}

	runner := &fakeSandboxRunner{
		t:         t,
		workspace: workRoot,
		exportPayload: `{
			"info": {"id":"ses_test"},
			"messages": [
				{"info":{"tokens":{"total":11840,"input":105,"output":52}}},
				{"info":{"tokens":{"total":14584,"input":651,"output":1310}}},
				{"info":{"tokens":{"total":14718,"input":136,"output":92}}},
				{"info":{"tokens":{"total":14861,"input":159,"output":96}}},
				{"info":{"tokens":{"total":16301,"input":1323,"output":32}}}
			]
		}`,
	}

	trace := AgentTrace{
		Framework: "opencode",
		SessionID: "ses_test",
	}

	collectOpenCodeSessionExport(
		context.Background(),
		runner,
		SandboxRunRequest{Workspace: workRoot, DockerImage: "utbench-agent-opencode-go:latest"},
		workRoot,
		traceDir,
		"boundary_000",
		&trace,
	)

	if trace.SessionExportError != "" {
		t.Fatalf("unexpected session export error: %s", trace.SessionExportError)
	}
	if got, want := trace.TokenSource, "actual"; got != want {
		t.Fatalf("token source = %q, want %q", got, want)
	}
	if got, want := trace.UsageSourceDetail, "opencode_session_export"; got != want {
		t.Fatalf("usage source detail = %q, want %q", got, want)
	}
	if trace.PromptTokens == nil || *trace.PromptTokens != 2374 {
		t.Fatalf("prompt tokens = %v, want 2374", trace.PromptTokens)
	}
	if trace.CompletionTokens == nil || *trace.CompletionTokens != 1582 {
		t.Fatalf("completion tokens = %v, want 1582", trace.CompletionTokens)
	}
	if trace.TotalTokens == nil || *trace.TotalTokens != 72304 {
		t.Fatalf("total tokens = %v, want 72304", trace.TotalTokens)
	}
}

func TestParseFileWritesFiltersOpenCodeWorkspaceNoise(t *testing.T) {
	got := parseFileWrites(
		[]string{
			"go mod init workspace",
			"go test ./workspace/generated_test.go",
			"cat /workspace/utbench_agent_prompt.md",
		},
		[]string{
			".utbench/xdg-config/opencode/config.json",
			".utbench/xdg-data/opencode/opencode.db",
			"go.mod",
			"generated_test.go",
			"test_runner",
		},
	)

	want := []string{"./workspace/generated_test.go", "generated_test.go"}
	if len(got) != len(want) {
		t.Fatalf("files_written len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("files_written[%d] = %q, want %q (all=%v)", i, got[i], want[i], got)
		}
	}
}
