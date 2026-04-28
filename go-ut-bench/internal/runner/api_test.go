package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripMarkdownFence_RemovesTrailingFenceOnly(t *testing.T) {
	in := "import unittest\n\nclass T(unittest.TestCase):\n    pass\n```"
	out := stripMarkdownFence(in)
	if out == in {
		t.Fatalf("expected trailing fence to be removed")
	}
	if out == "" {
		t.Fatalf("expected non-empty output")
	}
}

func TestBuildPrompt_PythonIncludesHardRequirements(t *testing.T) {
	source := "import math\n\ndef calc(x):\n    if x <= 0:\n        return 0\n    return x + 1\n"
	prompt := buildPrompt("python", "/tmp/boundary_001.py", source)

	checks := []string{
		"Task: Generate one complete test file for the provided source code.",
		"Mode: full_file",
		"Framework: pytest",
		"Use pytest function-based tests and import target symbols from `module_under_test`.",
		"Detected dependencies: math",
		"Keep test inputs small and representative; do not create stress tests or huge inputs.",
		"Never access real networks, real credentials, or real external services.",
		"- if x <= 0:",
		"Output raw code only.",
		"```python",
	}
	for _, item := range checks {
		if !strings.Contains(prompt, item) {
			t.Fatalf("prompt missing expected content: %q", item)
		}
	}

	forbidden := []string{
		"nearest_pair",
		"sliding-window or two-pointer counting algorithms",
		"Trace through the algorithm by hand",
		"你是一名资深单元测试工程师。",
		"Scenario:",
		"Complexity:",
	}
	for _, item := range forbidden {
		if strings.Contains(prompt, item) {
			t.Fatalf("prompt should not contain benchmark-specific or duplicated guidance: %q", item)
		}
	}
}

func TestExtractDependencies_GoImportBlock(t *testing.T) {
	source := "package demo\n\nimport (\n    \"fmt\"\n    \"net/http\"\n)\n"
	deps := extractDependencies(source, "go")
	joined := strings.Join(deps, ",")
	if !strings.Contains(joined, "fmt") || !strings.Contains(joined, "net/http") {
		t.Fatalf("unexpected go deps: %v", deps)
	}
}

func TestPromptTemplatePreview_UsesSharedStructure(t *testing.T) {
	preview := PromptTemplatePreview("go")
	checks := []string{
		"System",
		"Return only runnable test code.",
		"Language: go",
		"Framework: Go testing package",
		"Mode: full_file",
		"package preview",
	}
	for _, item := range checks {
		if !strings.Contains(preview, item) {
			t.Fatalf("preview missing expected content: %q", item)
		}
	}
}

func TestBuildPrompt_JavaIncludesPackageAndClassConstraints(t *testing.T) {
	source := "import java.util.*;\n\nclass Solution {\n    public List<Double> solve(List<Double> xs) { return xs; }\n}\n"
	prompt := buildPrompt("java", "/tmp/simple_function_001.java", source)

	checks := []string{
		"Declared package: default package",
		"Declared classes: Solution.",
		"If the source file has no `package` declaration, the test file must also have no `package` declaration.",
		"Instantiate the exact class declared in the source before calling instance methods; only call methods statically when the source declares them as `static`.",
	}
	for _, item := range checks {
		if !strings.Contains(prompt, item) {
			t.Fatalf("java prompt missing expected content: %q", item)
		}
	}
}

func TestBuildPrompt_PythonIncludesImportDiscipline(t *testing.T) {
	source := "import csv\n\ndef task_func(x):\n    return x\n"
	prompt := buildPrompt("python", "/tmp/complex_dependency_001.py", source)

	checks := []string{
		"If the test code uses a module such as `csv`, `json`, `os`, `tempfile`, or `xml.etree.ElementTree`, import it explicitly in the test file even if the source imports it.",
		"Do not introduce third-party modules that are not already present in the source context.",
	}
	for _, item := range checks {
		if !strings.Contains(prompt, item) {
			t.Fatalf("python prompt missing expected content: %q", item)
		}
	}
}

func TestPromptCatalog_WriteAndLoad(t *testing.T) {
	dir := t.TempDir()
	catalog, err := WritePromptCatalog(dir)
	if err != nil {
		t.Fatalf("WritePromptCatalog() error = %v", err)
	}
	if catalog.Strategy == "" || catalog.VersionID == "" {
		t.Fatalf("expected prompt catalog strategy and version")
	}

	loaded, err := LoadPromptCatalog(dir)
	if err != nil {
		t.Fatalf("LoadPromptCatalog() error = %v", err)
	}
	if loaded.VersionID != catalog.VersionID {
		t.Fatalf("version mismatch: got %s want %s", loaded.VersionID, catalog.VersionID)
	}
	if _, err := os.Stat(filepath.Join(dir, "prompt_catalog.json")); err != nil {
		t.Fatalf("expected prompt catalog file to exist: %v", err)
	}
	if loaded.Templates["python"][PromptModeFullFile] == "" {
		t.Fatalf("expected python full-file template to be present")
	}
}

func TestExtractFinishReason_OpenAIFormat(t *testing.T) {
	tests := []struct {
		name     string
		response map[string]any
		provider string
		want     bool
	}{
		{
			name: "finish_reason_length",
			response: map[string]any{
				"choices": []any{
					map[string]any{
						"message":       map[string]any{"content": "test"},
						"finish_reason": "length",
					},
				},
			},
			provider: "deepseek",
			want:     true,
		},
		{
			name: "finish_reason_stop",
			response: map[string]any{
				"choices": []any{
					map[string]any{
						"message":       map[string]any{"content": "test"},
						"finish_reason": "stop",
					},
				},
			},
			provider: "deepseek",
			want:     false,
		},
		{
			name: "no_finish_reason",
			response: map[string]any{
				"choices": []any{
					map[string]any{
						"message": map[string]any{"content": "test"},
					},
				},
			},
			provider: "deepseek",
			want:     false,
		},
		{
			name: "dashscope_format_length",
			response: map[string]any{
				"output": map[string]any{
					"choices": []any{
						map[string]any{
							"message":       map[string]any{"content": "test"},
							"finish_reason": "length",
						},
					},
				},
			},
			provider: "dashscope",
			want:     true,
		},
		{
			name: "dashscope_format_stop",
			response: map[string]any{
				"output": map[string]any{
					"choices": []any{
						map[string]any{
							"message":       map[string]any{"content": "test"},
							"finish_reason": "stop",
						},
					},
				},
			},
			provider: "dashscope",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractFinishReason(tt.response, tt.provider)
			if got != tt.want {
				t.Errorf("extractFinishReason() = %v, want %v", got, tt.want)
			}
		})
	}
}
