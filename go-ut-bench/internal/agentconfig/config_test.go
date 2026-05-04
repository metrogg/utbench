package agentconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadExpandsFrameworkModelSkillMatrix(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "agents.yaml")
	raw := `
models: [deepseek, qwen]
frameworks:
  aider:
    kind: cli_agent
    command: "aider {{.PromptFile}}"
    sandbox_mode: local
    compatible_models: [deepseek]
skills:
  unit_test_skill:
    version: "1"
    inject_mode: prompt_append
    compatible_frameworks: [aider]
`
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	subjects, err := Load(path, []string{"deepseek", "qwen"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, subject := range subjects {
		got[subject.Spec.ID] = true
	}
	want := []string{
		"model_api__deepseek__no_skill",
		"model_api__qwen__no_skill",
		"aider__deepseek__no_skill",
		"aider__deepseek__unit_test_skill",
	}
	for _, id := range want {
		if !got[id] {
			t.Fatalf("expected subject %s in %+v", id, got)
		}
	}
	if got["aider__qwen__no_skill"] {
		t.Fatalf("incompatible model should have been filtered")
	}
}

func TestLoadFiltersSelectedSubjects(t *testing.T) {
	subjects, err := Load("", []string{"deepseek"}, []string{"model_api__deepseek__no_skill"})
	if err != nil {
		t.Fatal(err)
	}
	if len(subjects) != 1 || subjects[0].Spec.ID != "model_api__deepseek__no_skill" {
		t.Fatalf("unexpected selected subjects: %+v", subjects)
	}
}

func TestLoadParsesEnvFromHost(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "agents.yaml")
	raw := `
models: [deepseek]
frameworks:
  opencode:
    kind: cli_agent
    command: "opencode run {{.ContainerPrompt}}"
    sandbox_mode: local
    env:
      OPENCODE_CONFIG_CONTENT: |
        {"provider":"{{.ModelProvider}}","endpoint":"{{.ModelEndpoint}}","api_key_env":"{{.ModelAPIKeyEnv}}"}
    env_from_host:
      - OPENAI_API_KEY
      - ARK_API_KEY
`
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	subjects, err := Load(path, []string{"deepseek"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(subjects) != 2 {
		t.Fatalf("expected baseline + opencode, got %d", len(subjects))
	}
	var opencode ResolvedSubject
	for _, subject := range subjects {
		if subject.Spec.Framework == "opencode" {
			opencode = subject
			break
		}
	}
	if opencode.Spec.ID == "" {
		t.Fatalf("expected opencode subject in %+v", subjects)
	}
	if got, want := opencode.Framework.EnvFromHost, []string{"ARK_API_KEY", "OPENAI_API_KEY"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("unexpected env_from_host: got=%v want=%v", got, want)
	}
	if _, ok := opencode.Framework.Env["OPENCODE_CONFIG_CONTENT"]; !ok {
		t.Fatalf("expected templated env to be preserved")
	}
}

func TestLoadParsesLanguageAwareSandboxFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "agents.yaml")
	raw := `
models: [deepseek]
frameworks:
  opencode:
    kind: cli_agent
    command: "opencode run {{.ContainerPrompt}}"
    sandbox_mode: docker
    docker_images:
      python: utbench-agent-opencode-python:latest
      go: utbench-agent-opencode-go:latest
    preflight:
      python:
        - python3 --version
        - pytest --version
      go:
        - go version
    forbidden_command_patterns:
      - apt-get install
      - pip install
`
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	subjects, err := Load(path, []string{"deepseek"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var opencode ResolvedSubject
	for _, subject := range subjects {
		if subject.Spec.Framework == "opencode" {
			opencode = subject
			break
		}
	}
	if opencode.Spec.ID == "" {
		t.Fatalf("expected opencode subject in %+v", subjects)
	}
	if got := opencode.Framework.DockerImages["python"]; got != "utbench-agent-opencode-python:latest" {
		t.Fatalf("unexpected python docker image: %q", got)
	}
	if got := opencode.Framework.Preflight["go"]; len(got) != 1 || got[0] != "go version" {
		t.Fatalf("unexpected go preflight commands: %+v", got)
	}
	if got := opencode.Framework.ForbiddenCommandPatterns; len(got) != 2 || got[0] != "apt-get install" || got[1] != "pip install" {
		t.Fatalf("unexpected forbidden command patterns: %+v", got)
	}
}

func TestLoadPrefersNestedSandboxBlock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "agents.yaml")
	raw := `
models: [deepseek]
frameworks:
  opencode:
    kind: cli_agent
    command: "opencode run {{.ContainerPrompt}}"
    sandbox:
      provider: docker
      mode: docker
      images:
        python: utbench-agent-opencode-python:v2
      timeout_seconds: 900
      network_disabled: false
      cpu: "3"
      memory: 3g
    docker_images:
      python: legacy-ignored:latest
    timeout_seconds: 120
    network_disabled: true
`
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	subjects, err := Load(path, []string{"deepseek"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	var opencode ResolvedSubject
	for _, subject := range subjects {
		if subject.Spec.Framework == "opencode" {
			opencode = subject
			break
		}
	}
	if opencode.Spec.ID == "" {
		t.Fatalf("expected opencode subject in %+v", subjects)
	}
	if got := opencode.Framework.Sandbox.Provider; got != "docker" {
		t.Fatalf("unexpected sandbox provider: %q", got)
	}
	if got := opencode.Framework.Sandbox.Images["python"]; got != "utbench-agent-opencode-python:v2" {
		t.Fatalf("unexpected sandbox image: %q", got)
	}
	if got := opencode.Framework.TimeoutSeconds; got != 900 {
		t.Fatalf("unexpected timeout seconds: %d", got)
	}
	if got := opencode.Framework.NetworkDisabled; got {
		t.Fatalf("expected nested sandbox config to override network_disabled")
	}
	if got := opencode.Framework.CPU; got != "3" {
		t.Fatalf("unexpected cpu: %q", got)
	}
	if got := opencode.Framework.Memory; got != "3g" {
		t.Fatalf("unexpected memory: %q", got)
	}
}
