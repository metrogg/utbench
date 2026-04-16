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
