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
	subject := agentconfig.ResolvedSubject{
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
			Preflight: map[string][]string{
				"python": {"go version"},
			},
		},
		Skill: contracts.SkillSpec{Name: agentconfig.NoSkill, Enabled: true},
	}
	adapter := newCLIAgentAdapter(NewSandboxRunner())
	result := adapter.Generate(context.Background(), AgentGenerateRequest{
		Subject:    subject,
		Model:      modelConfig{Name: "deepseek", Model: "deepseek-chat"},
		Sample:     contracts.SampleRef{ID: "sample_000", Language: "python", Path: samplePath},
		Prompt:     "Generate tests",
		TestPath:   filepath.Join(tmp, "out.py"),
		MetaRoot:   filepath.Join(tmp, "metadata"),
		OutputRoot: tmp,
		RunID:      "run_cli_agent",
	})
	if result.Error != nil {
		t.Fatalf("adapter.Generate returned error: %+v raw=%+v", result.Error, result.RawResponse)
	}
	if result.Truncated {
		t.Fatalf("did not expect truncation")
	}
	if result.LatencyMS <= 0 {
		t.Fatalf("expected latency to be recorded")
	}
	if !strings.Contains(result.Code, "def test_generated") {
		t.Fatalf("unexpected generated code: %s", result.Code)
	}
	if result.Trace.TracePath == "" || result.Trace.WorkspaceDiffPath == "" || result.Trace.SandboxFingerprint == "" {
		t.Fatalf("expected trace artifacts, got %+v", result.Trace)
	}
	// 验证新增的 trace 字段
	if result.Trace.InteractionCount < 1 {
		t.Fatalf("expected interaction_count >= 1, got %d", result.Trace.InteractionCount)
	}
	if result.Trace.StartedAt.IsZero() || result.Trace.FinishedAt.IsZero() {
		t.Fatalf("expected started_at and finished_at to be set")
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
	subject := agentconfig.ResolvedSubject{
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
			Preflight: map[string][]string{
				"python": {"go version"},
			},
			Env: map[string]string{
				"OPENCODE_CONFIG_CONTENT": `{"endpoint":"{{.ModelEndpoint}}","api_key_env":"{{.ModelAPIKeyEnv}}","model":"{{.ModelID}}"}`,
			},
			EnvFromHost: []string{"TEST_AGENT_API_KEY"},
		},
		Skill: contracts.SkillSpec{Name: "unit_test_skill", Enabled: true},
	}
	adapter := newCLIAgentAdapter(NewSandboxRunner())
	result := adapter.Generate(context.Background(), AgentGenerateRequest{
		Subject: subject,
		Model: modelConfig{
			Name:      "deepseek",
			Model:     "deepseek-v4-flash",
			Provider:  "deepseek",
			Endpoint:  "https://api.deepseek.com",
			APIKeyEnv: "TEST_AGENT_API_KEY",
		},
		Sample:     contracts.SampleRef{ID: "sample_001", Language: "python", Path: samplePath},
		Prompt:     "Generate tests",
		TestPath:   filepath.Join(tmp, "out_env.py"),
		MetaRoot:   filepath.Join(tmp, "metadata"),
		OutputRoot: tmp,
		RunID:      "run_cli_agent_env",
	})
	if result.Error != nil {
		t.Fatalf("adapter.Generate returned error: %+v raw=%+v", result.Error, result.RawResponse)
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

func TestFrameworkDockerImageAndPreflightCommands(t *testing.T) {
	// 新优先级: Sandbox.Image → DockerImage → DockerImages[lang] → images["default"]
	fw := agentconfig.FrameworkSpec{
		DockerImage: "fallback:latest",
		DockerImages: map[string]string{
			"python":  "python:latest",
			"default": "default:latest",
		},
		Preflight: map[string][]string{
			"python": {"python3 --version", "pytest --version"},
		},
	}
	// DockerImage 优先于 DockerImages
	if got := frameworkDockerImage(fw, "python"); got != "fallback:latest" {
		t.Fatalf("unexpected python image: %q (expected DockerImage to take priority)", got)
	}
	if got := frameworkDockerImage(fw, "go"); got != "fallback:latest" {
		t.Fatalf("unexpected fallback image: %q (expected DockerImage to take priority)", got)
	}
	// 仅当 DockerImage 为空时才回退到 DockerImages
	fw2 := agentconfig.FrameworkSpec{
		DockerImages: map[string]string{
			"python":  "python:latest",
			"default": "default:latest",
		},
	}
	if got := frameworkDockerImage(fw2, "python"); got != "python:latest" {
		t.Fatalf("unexpected python image without DockerImage: %q", got)
	}
	if got := frameworkDockerImage(fw2, "go"); got != "default:latest" {
		t.Fatalf("unexpected fallback image without DockerImage: %q", got)
	}
	commands := frameworkPreflightCommands(fw, "python")
	if len(commands) != 2 || commands[0] != "pytest --version" || commands[1] != "python3 --version" {
		t.Fatalf("unexpected preflight commands: %+v", commands)
	}
}

func TestBuildSandboxRunRequestUsesSandboxBlock(t *testing.T) {
	// 新优先级: Sandbox.Image 统一镜像 > DockerImage > Sandbox.Images[lang]
	fw := agentconfig.FrameworkSpec{
		Sandbox: agentconfig.SandboxSpec{
			Provider:        "docker",
			Mode:            "docker",
			Image:           "unified:latest",
			Images:          map[string]string{"python": "python:v2"},
			TimeoutSeconds:  777,
			NetworkDisabled: false,
			CPU:             "3",
			Memory:          "3g",
		},
		DockerImage:     "legacy:latest",
		SandboxMode:     "local",
		TimeoutSeconds:  100,
		NetworkDisabled: true,
		CPU:             "1",
		Memory:          "1g",
	}
	req := buildSandboxRunRequest("C:\\out", fw, "python", "C:\\ws", "echo ok", nil, nil)
	if req.Provider != "docker" {
		t.Fatalf("unexpected provider: %q", req.Provider)
	}
	// Sandbox.Image 优先
	if req.DockerImage != "unified:latest" {
		t.Fatalf("unexpected image: %q (expected Sandbox.Image to take priority)", req.DockerImage)
	}
	if req.TimeoutSeconds != 777 {
		t.Fatalf("unexpected timeout: %d", req.TimeoutSeconds)
	}
	if req.NetworkDisabled {
		t.Fatalf("expected nested sandbox config to override network_disabled")
	}
	if req.CPU != "3" || req.Memory != "3g" {
		t.Fatalf("unexpected resources: cpu=%q memory=%q", req.CPU, req.Memory)
	}

	// 无 Sandbox.Image 时回退到 DockerImage
	fw2 := agentconfig.FrameworkSpec{
		Sandbox: agentconfig.SandboxSpec{
			Provider:       "docker",
			Mode:           "docker",
			Images:         map[string]string{"python": "python:v2"},
			TimeoutSeconds: 500,
		},
		DockerImage: "legacy:latest",
	}
	req2 := buildSandboxRunRequest("C:\\out", fw2, "python", "C:\\ws", "echo ok", nil, nil)
	if req2.DockerImage != "legacy:latest" {
		t.Fatalf("unexpected image fallback: %q (expected DockerImage)", req2.DockerImage)
	}

	// 无 Sandbox.Image 也无 DockerImage 时回退到 Sandbox.Images[lang]
	fw3 := agentconfig.FrameworkSpec{
		Sandbox: agentconfig.SandboxSpec{
			Provider:       "docker",
			Mode:           "docker",
			Images:         map[string]string{"python": "python:v2"},
			TimeoutSeconds: 500,
		},
	}
	req3 := buildSandboxRunRequest("C:\\out", fw3, "python", "C:\\ws", "echo ok", nil, nil)
	if req3.DockerImage != "python:v2" {
		t.Fatalf("unexpected per-language fallback: %q (expected Sandbox.Images[lang])", req3.DockerImage)
	}
}

func TestDetectSandboxPolicyViolation(t *testing.T) {
	commands := []string{
		"python3 --version",
		"apt-get install -y python3 python3-pip",
	}
	violation := detectSandboxPolicyViolation(commands, []string{"apt-get install"})
	if violation == "" {
		t.Fatalf("expected sandbox policy violation")
	}
}

func TestDetectSandboxPolicyViolationIgnoresOutputNoise(t *testing.T) {
	commands := []string{
		"g++ --version",
		"ctest --output-on-failure",
	}
	violation := detectSandboxPolicyViolation(commands, []string{"npm install", "pip install"})
	if violation != "" {
		t.Fatalf("expected no sandbox policy violation from unrelated commands, got %q", violation)
	}
}

func TestBuildSampleEnvironmentSetupCommandsPythonRepoLevel(t *testing.T) {
	tmp := t.TempDir()
	sampleDir := filepath.Join(tmp, "simple_function_000")
	if err := os.MkdirAll(sampleDir, 0o755); err != nil {
		t.Fatal(err)
	}
	samplePath := filepath.Join(sampleDir, "entry.py")
	if err := os.WriteFile(samplePath, []byte("def f():\n    return 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	meta := `{"sample_id":"simple_function_000","module_import":"pkg.mod","workspace_root":"./workspace","requirements":["six","requests==2.31.0"]}`
	if err := os.WriteFile(filepath.Join(sampleDir, "meta.json"), []byte(meta), 0o644); err != nil {
		t.Fatal(err)
	}
	commands := buildSampleEnvironmentSetupCommands(contracts.SampleRef{
		ID:       "simple_function_000",
		Language: "python",
		Path:     samplePath,
	}, sampleDir)
	if len(commands) != 1 {
		t.Fatalf("expected one command, got %+v", commands)
	}
	if !strings.Contains(commands[0], "python3 -m pip install") || !strings.Contains(commands[0], "six") || !strings.Contains(commands[0], "requests==2.31.0") {
		t.Fatalf("unexpected python dependency setup command: %q", commands[0])
	}
}

func TestBuildSampleEnvironmentSetupCommandsWorkspaceFiles(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module demo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "pom.xml"), []byte("<project/>"), 0o644); err != nil {
		t.Fatal(err)
	}
	goCommands := buildSampleEnvironmentSetupCommands(contracts.SampleRef{Language: "go"}, tmp)
	if len(goCommands) != 1 || goCommands[0] != "go mod download" {
		t.Fatalf("unexpected go commands: %+v", goCommands)
	}
	javaCommands := buildSampleEnvironmentSetupCommands(contracts.SampleRef{Language: "java"}, tmp)
	if len(javaCommands) != 1 || javaCommands[0] != "mvn -q -DskipTests dependency:go-offline" {
		t.Fatalf("unexpected java commands: %+v", javaCommands)
	}
}

func TestInjectAgentNativeSkillForClaudeCodeRenamesInstructionToSkillMD(t *testing.T) {
	tmp := t.TempDir()
	skillSrc := filepath.Join(tmp, "instructions.md")
	refDir := filepath.Join(tmp, "references")
	if err := os.WriteFile(skillSrc, []byte("# unit test skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(refDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(refDir, "checklist.md"), []byte("- check\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	workRoot := filepath.Join(tmp, "workspace")
	if err := os.MkdirAll(workRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	dest, err := injectAgentNativeSkill(workRoot, "claudecode", contracts.SkillSpec{
		Name:            "unit_test_skill",
		Enabled:         true,
		InstructionPath: skillSrc,
		Files:           []string{refDir},
	})
	if err != nil {
		t.Fatalf("injectAgentNativeSkill returned error: %v", err)
	}

	wantDir := filepath.Join(workRoot, ".claude", "skills", "unit_test_skill")
	if dest != wantDir {
		t.Fatalf("dest = %q, want %q", dest, wantDir)
	}
	if _, err := os.Stat(filepath.Join(dest, "SKILL.md")); err != nil {
		t.Fatalf("expected SKILL.md to exist: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "references", "checklist.md")); err != nil {
		t.Fatalf("expected copied references dir to exist: %v", err)
	}
}

func TestInjectAgentNativeSkillForCodexAddsFrontmatter(t *testing.T) {
	tmp := t.TempDir()
	skillSrc := filepath.Join(tmp, "instructions.md")
	if err := os.WriteFile(skillSrc, []byte("# unit test skill\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	workRoot := filepath.Join(tmp, "workspace")
	if err := os.MkdirAll(workRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	dest, err := injectAgentNativeSkill(workRoot, "codex", contracts.SkillSpec{
		Name:            "unit_test_skill",
		Enabled:         true,
		Description:     "Generate compact unit tests.",
		InstructionPath: skillSrc,
	})
	if err != nil {
		t.Fatalf("injectAgentNativeSkill returned error: %v", err)
	}

	raw, err := os.ReadFile(filepath.Join(dest, "SKILL.md"))
	if err != nil {
		t.Fatalf("expected SKILL.md to exist: %v", err)
	}
	text := string(raw)
	if !strings.HasPrefix(text, "---\nname: unit_test_skill\n") {
		t.Fatalf("expected Codex skill frontmatter, got %q", text)
	}
	if !strings.Contains(text, "# unit test skill") {
		t.Fatalf("expected original skill content to be preserved, got %q", text)
	}
}

func TestGenerateWithCLIAgentFailsOnSandboxPreflight(t *testing.T) {
	tmp := t.TempDir()
	samplePath := filepath.Join(tmp, "sample.py")
	if err := os.WriteFile(samplePath, []byte("def add(a, b):\n    return a + b\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	subject := agentconfig.ResolvedSubject{
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
			Command:     `go version`,
			Preflight: map[string][]string{
				"python": {"definitely_missing_utbench_command"},
			},
		},
		Skill: contracts.SkillSpec{Name: agentconfig.NoSkill, Enabled: true},
	}
	adapter := newCLIAgentAdapter(NewSandboxRunner())
	result := adapter.Generate(context.Background(), AgentGenerateRequest{
		Subject:    subject,
		Model:      modelConfig{Name: "deepseek", Model: "deepseek-chat"},
		Sample:     contracts.SampleRef{ID: "sample_000", Language: "python", Path: samplePath},
		Prompt:     "Generate tests",
		TestPath:   filepath.Join(tmp, "out.py"),
		MetaRoot:   filepath.Join(tmp, "metadata"),
		OutputRoot: tmp,
		RunID:      "run_cli_agent_preflight_fail",
	})
	if result.Error == nil || result.Error.Kind != "sandbox_preflight_error" {
		t.Fatalf("expected sandbox_preflight_error, got %+v", result.Error)
	}
	if len(result.Trace.PreflightChecks) != 1 || result.Trace.PreflightChecks[0].Passed {
		t.Fatalf("expected failed preflight checks, got %+v", result.Trace.PreflightChecks)
	}
}
