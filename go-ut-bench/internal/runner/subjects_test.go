package runner

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-ut-bench/internal/agentconfig"
	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/obs"
)

func TestGenerateWithCLIAgentLocalFake(t *testing.T) {
	tmp := t.TempDir()
	samplePath := filepath.Join(tmp, "sample.py")
	if err := os.WriteFile(samplePath, []byte("def add(a, b):\n    return a + b\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fakeAgent := filepath.Join(tmp, "fake_agent.go")
	fakeAgentSrc := `package main
import (
	"os"
)
func main() {
	if len(os.Args) < 2 { panic("missing output") }
	content := "import pytest\n\nfrom module_under_test import add\n\nfunc test_placeholder():\n    pass\n"
	content = "import pytest\n\n\ndef test_generated():\n    assert True\n"
	if err := os.WriteFile(os.Args[1], []byte(content), 0644); err != nil { panic(err) }
}`
	if err := os.WriteFile(fakeAgent, []byte(fakeAgentSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	target := subjectTarget{
		subject: agentconfig.ResolvedSubject{
			Spec: contracts.SubjectSpec{
				ID:        "fake_agent__deepseek__no_skill",
				Kind:      agentconfig.KindCLIAgent,
				Framework: "fake_agent",
				Model:     "deepseek",
				Skill:     agentconfig.NoSkill,
			},
			Framework: agentconfig.FrameworkSpec{
				Name:        "fake_agent",
				Kind:        agentconfig.KindCLIAgent,
				Enabled:     true,
				SandboxMode: "local",
				Command:     `go run "` + fakeAgent + `" "{{.OutputFile}}"`,
			},
			Skill: contracts.SkillSpec{Name: agentconfig.NoSkill, Enabled: true},
		},
		model: modelConfig{Name: "deepseek", Model: "deepseek-chat"},
	}
	spec := contracts.RunSpec{OutputRoot: tmp, RunID: "run_cli_agent"}
	svc := NewService(obs.NewLogger(false, ""))
	code, raw, trace, latency, _, _, _, truncated, errInfo := svc.generateWithCLIAgent(
		context.Background(),
		spec,
		target,
		contracts.SampleRef{ID: "sample_000", Language: "python", Path: samplePath},
		"Generate tests",
		filepath.Join(tmp, "out.py"),
		filepath.Join(tmp, "metadata"),
	)
	if errInfo != nil {
		t.Fatalf("generateWithCLIAgent returned error: %+v raw=%+v", errInfo, raw)
	}
	if truncated {
		t.Fatalf("did not expect truncation")
	}
	if latency <= 0 {
		t.Fatalf("expected latency to be recorded")
	}
	if !strings.Contains(code, "def test_generated") {
		t.Fatalf("unexpected generated code: %s", code)
	}
	if trace.TracePath == "" || trace.WorkspaceDiffPath == "" || trace.SandboxFingerprint == "" {
		t.Fatalf("expected trace artifacts, got %+v", trace)
	}
}

func TestGenerateWithCLIAgentRendersEnvAndPassesThroughHostSecrets(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("TEST_AGENT_API_KEY", "secret-token")
	samplePath := filepath.Join(tmp, "sample.py")
	if err := os.WriteFile(samplePath, []byte("def add(a, b):\n    return a + b\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	fakeAgent := filepath.Join(tmp, "fake_agent_env.go")
	fakeAgentSrc := `package main
import (
	"encoding/json"
	"os"
)
func main() {
	if len(os.Args) < 3 { panic("missing args") }
	outPath := os.Args[1]
	envPath := os.Args[2]
	content := "import pytest\n\n\ndef test_generated():\n    assert True\n"
	if err := os.WriteFile(outPath, []byte(content), 0644); err != nil { panic(err) }
	payload := map[string]string{
		"api_key": os.Getenv("TEST_AGENT_API_KEY"),
		"config": os.Getenv("OPENCODE_CONFIG_CONTENT"),
	}
	raw, err := json.Marshal(payload)
	if err != nil { panic(err) }
	if err := os.WriteFile(envPath, raw, 0644); err != nil { panic(err) }
}`
	if err := os.WriteFile(fakeAgent, []byte(fakeAgentSrc), 0o644); err != nil {
		t.Fatal(err)
	}
	envCapture := filepath.Join(tmp, "env_capture.json")
	target := subjectTarget{
		subject: agentconfig.ResolvedSubject{
			Spec: contracts.SubjectSpec{
				ID:        "opencode__deepseek__unit_test_skill",
				Kind:      agentconfig.KindCLIAgent,
				Framework: "opencode",
				Model:     "deepseek",
				Skill:     "unit_test_skill",
			},
			Framework: agentconfig.FrameworkSpec{
				Name:        "opencode",
				Kind:        agentconfig.KindCLIAgent,
				Enabled:     true,
				SandboxMode: "local",
				Command:     `go run "` + fakeAgent + `" "{{.OutputFile}}" "` + envCapture + `"`,
				Env: map[string]string{
					"OPENCODE_CONFIG_CONTENT": `{"endpoint":"{{.ModelEndpoint}}","api_key_env":"{{.ModelAPIKeyEnv}}","model":"{{.ModelID}}"}`,
				},
				EnvFromHost: []string{"TEST_AGENT_API_KEY"},
			},
			Skill: contracts.SkillSpec{Name: "unit_test_skill", Enabled: true},
		},
		model: modelConfig{
			Name:      "deepseek",
			Model:     "deepseek-v4-flash",
			Provider:  "deepseek",
			Endpoint:  "https://api.deepseek.com",
			APIKeyEnv: "TEST_AGENT_API_KEY",
		},
	}
	spec := contracts.RunSpec{OutputRoot: tmp, RunID: "run_cli_agent_env"}
	svc := NewService(obs.NewLogger(false, ""))
	_, raw, _, _, _, _, _, _, errInfo := svc.generateWithCLIAgent(
		context.Background(),
		spec,
		target,
		contracts.SampleRef{ID: "sample_001", Language: "python", Path: samplePath},
		"Generate tests",
		filepath.Join(tmp, "out_env.py"),
		filepath.Join(tmp, "metadata"),
	)
	if errInfo != nil {
		t.Fatalf("generateWithCLIAgent returned error: %+v raw=%+v", errInfo, raw)
	}
	payloadRaw, err := os.ReadFile(envCapture)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]string
	if err := json.Unmarshal(payloadRaw, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["api_key"] != "secret-token" {
		t.Fatalf("expected host env passthrough, got %q", payload["api_key"])
	}
	wantConfigParts := []string{
		`"endpoint":"https://api.deepseek.com"`,
		`"api_key_env":"TEST_AGENT_API_KEY"`,
		`"model":"deepseek-v4-flash"`,
	}
	for _, part := range wantConfigParts {
		if !strings.Contains(payload["config"], part) {
			t.Fatalf("expected config to contain %s, got %s", part, payload["config"])
		}
	}
}
