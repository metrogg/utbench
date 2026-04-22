package runner

import (
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
		"You are an expert unit testing engineer.",
		"你是一名资深单元测试工程师。",
		"MUST import target symbols from local module `boundary_001`",
		"if x <= 0:",
		"Error Prevention Checklist（错误预防清单，仅内部执行）",
		"No placeholder tests like `assert True`.",
		"Dependencies detected（检测到依赖）: math",
		"Return raw test code only (no Markdown fences).",
	}
	for _, item := range checks {
		if !strings.Contains(prompt, item) {
			t.Fatalf("prompt missing expected content: %q", item)
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
