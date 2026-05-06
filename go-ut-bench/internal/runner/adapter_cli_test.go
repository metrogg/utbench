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
		SandboxRunRequest{Workspace: workRoot, DockerImage: "utbench-agent-base:latest"},
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
		SandboxRunRequest{Workspace: workRoot, DockerImage: "utbench-agent-base:latest"},
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

func TestParseCodeBuddyJSONOutput(t *testing.T) {
	trace := &AgentTrace{Framework: "codebuddy"}
	stdout := `{"result":"Test file created","session_id":"sess_abc123","usage":{"input_tokens":5000,"output_tokens":1200}}`
	parseCodeBuddyJSONOutput(trace, stdout)

	if trace.SessionID != "sess_abc123" {
		t.Fatalf("session_id = %q, want %q", trace.SessionID, "sess_abc123")
	}
	if trace.PromptTokens == nil || *trace.PromptTokens != 5000 {
		t.Fatalf("prompt_tokens = %v, want 5000", trace.PromptTokens)
	}
	if trace.CompletionTokens == nil || *trace.CompletionTokens != 1200 {
		t.Fatalf("completion_tokens = %v, want 1200", trace.CompletionTokens)
	}
	if trace.TokenSource != "actual" {
		t.Fatalf("token_source = %q, want %q", trace.TokenSource, "actual")
	}
	if trace.UsageSourceDetail != "codebuddy_json_output" {
		t.Fatalf("usage_source_detail = %q, want %q", trace.UsageSourceDetail, "codebuddy_json_output")
	}
}

func TestParseCodeBuddyJSONOutputWithMixedOutput(t *testing.T) {
	trace := &AgentTrace{Framework: "codebuddy"}
	stdout := "Some log lines\nTool call: write_file\nAnother log line\n" +
		`{"result":"Done","session_id":"sess_mixed","usage":{"prompt_tokens":3000,"completion_tokens":800,"total_tokens":3800}}`
	parseCodeBuddyJSONOutput(trace, stdout)

	if trace.SessionID != "sess_mixed" {
		t.Fatalf("session_id = %q, want %q", trace.SessionID, "sess_mixed")
	}
	if trace.PromptTokens == nil || *trace.PromptTokens != 3000 {
		t.Fatalf("prompt_tokens = %v, want 3000", trace.PromptTokens)
	}
	if trace.CompletionTokens == nil || *trace.CompletionTokens != 800 {
		t.Fatalf("completion_tokens = %v, want 800", trace.CompletionTokens)
	}
}

func TestSummarizeAgentCommandErrorPrefersActionableTail(t *testing.T) {
	stderr := strings.Join([]string{
		"Reading additional input from stdin...",
		"OpenAI Codex v0.43.0",
		"2026-05-06T04:59:11.123Z ERROR codex_api::endpoint::responses_websocket: failed to connect to websocket: HTTP error: 401 Unauthorized, url: wss://api.openai.com/v1/responses",
		"ERROR: Reconnecting... 5/5",
	}, "\n")

	got := summarizeAgentCommandError(stderr, "exit status 1", 300)
	if strings.Contains(got, "Reading additional input") {
		t.Fatalf("summary kept startup prefix: %q", got)
	}
	if !strings.Contains(got, "401 Unauthorized") {
		t.Fatalf("summary = %q, want actionable 401 detail", got)
	}
}

func TestExtractFinalJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "single JSON object",
			input: `{"key":"value"}`,
			want:  `{"key":"value"}`,
		},
		{
			name:  "JSON at end of mixed output",
			input: "log line 1\nlog line 2\n{\"result\":\"ok\"}",
			want:  `{"result":"ok"}`,
		},
		{
			name:  "no JSON",
			input: "just plain text",
			want:  "",
		},
		{
			name:  "nested braces in JSON",
			input: `{"outer":{"inner":"val"}}`,
			want:  `{"outer":{"inner":"val"}}`,
		},
		{
			name:  "JSON array",
			input: `[{"type":"message","role":"user"},{"type":"message","role":"assistant"}]`,
			want:  `[{"type":"message","role":"user"},{"type":"message","role":"assistant"}]`,
		},
		{
			name:  "JSON array with trailing text",
			input: "some log\n[{\"type\":\"msg\"}]",
			want:  `[{"type":"msg"}]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractFinalJSON(tt.input)
			if got != tt.want {
				t.Fatalf("extractFinalJSON() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseUsageAndSessionWithCodeBuddy(t *testing.T) {
	trace := &AgentTrace{Framework: "codebuddy"}
	stdout := `{"result":"Done","session_id":"sess_cb","usage":{"input_tokens":2000,"output_tokens":500}}`
	parseUsageAndSession(trace, stdout, "")

	if trace.SessionID != "sess_cb" {
		t.Fatalf("session_id = %q, want %q", trace.SessionID, "sess_cb")
	}
	if trace.PromptTokens == nil || *trace.PromptTokens != 2000 {
		t.Fatalf("prompt_tokens = %v, want 2000", trace.PromptTokens)
	}
	if trace.CompletionTokens == nil || *trace.CompletionTokens != 500 {
		t.Fatalf("completion_tokens = %v, want 500", trace.CompletionTokens)
	}
	if got, want := trace.UsageSourceDetail, "codebuddy_json_output"; got != want {
		t.Fatalf("usage_source_detail = %q, want %q", got, want)
	}
}

func TestParseClaudeCodeJSONOutput(t *testing.T) {
	trace := &AgentTrace{Framework: "claudecode"}
	stdout := `{"session_id":"sess_claude","usage":{"input_tokens":4100,"output_tokens":900,"total_tokens":5000}}`
	parseClaudeCodeJSONOutput(trace, stdout)

	if trace.SessionID != "sess_claude" {
		t.Fatalf("session_id = %q, want %q", trace.SessionID, "sess_claude")
	}
	if trace.PromptTokens == nil || *trace.PromptTokens != 4100 {
		t.Fatalf("prompt_tokens = %v, want 4100", trace.PromptTokens)
	}
	if trace.CompletionTokens == nil || *trace.CompletionTokens != 900 {
		t.Fatalf("completion_tokens = %v, want 900", trace.CompletionTokens)
	}
	if trace.TotalTokens == nil || *trace.TotalTokens != 5000 {
		t.Fatalf("total_tokens = %v, want 5000", trace.TotalTokens)
	}
	if got, want := trace.UsageSourceDetail, "claudecode_json_output"; got != want {
		t.Fatalf("usage_source_detail = %q, want %q", got, want)
	}
}

func TestParseUsageAndSessionWithClaudeCodeJSONL(t *testing.T) {
	trace := &AgentTrace{Framework: "claudecode"}
	stdout := strings.Join([]string{
		`{"type":"init","session_id":"sess_stream"}`,
		`{"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"bash","input":{"command":"go test ./..."}}]}}`,
		`{"type":"result","usage":{"input_tokens":1200,"output_tokens":300,"total_tokens":1500}}`,
	}, "\n")
	parseUsageAndSession(trace, stdout, "")

	if trace.SessionID != "sess_stream" {
		t.Fatalf("session_id = %q, want %q", trace.SessionID, "sess_stream")
	}
	if trace.PromptTokens == nil || *trace.PromptTokens != 1200 {
		t.Fatalf("prompt_tokens = %v, want 1200", trace.PromptTokens)
	}
	if trace.CompletionTokens == nil || *trace.CompletionTokens != 300 {
		t.Fatalf("completion_tokens = %v, want 300", trace.CompletionTokens)
	}
	if got, want := trace.UsageSourceDetail, "claudecode_jsonl_lines"; got != want {
		t.Fatalf("usage_source_detail = %q, want %q", got, want)
	}
	tools := parseClaudeCodeToolCalls(stdout, "")
	if len(tools) != 1 || tools[0].Tool != "bash" {
		t.Fatalf("unexpected tools: %+v", tools)
	}
}
