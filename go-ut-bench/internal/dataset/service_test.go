package dataset

import (
	"os"
	"path/filepath"
	"testing"

	"go-ut-bench/internal/contracts"
)

func TestDiscoverSamplesFromManifest(t *testing.T) {
	root := t.TempDir()
	datasetRoot := filepath.Join(root, "dataset")
	if err := os.MkdirAll(filepath.Join(datasetRoot, "python"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(datasetRoot, "go"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(datasetRoot, "python", "boundary_000.py"), []byte("def x():\n    return 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(datasetRoot, "go", "complex_dependency_go_complex_dependency_0.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	manifestPath := filepath.Join(root, "dataset_l1.json")
	manifest := `{
  "level": "l1",
  "samples": [
    {"id":"boundary_000","language":"python","category":"self_contained","path":"python/boundary_000.py"},
    {"id":"complex_dependency_go_complex_dependency_0","language":"go","category":"module_level","path":"go/complex_dependency_go_complex_dependency_0.go"}
  ]
}`
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	spec := contracts.RunSpec{
		RunID:           "r1",
		DatasetRoot:     datasetRoot,
		OutputRoot:      root,
		ConfigPath:      "dummy",
		DatasetManifest: manifestPath,
		DatasetClasses:  []string{"module_level"},
		Languages:       []string{"go", "python"},
	}

	samples, err := svc.DiscoverSamples(spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(samples) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(samples))
	}
	if samples[0].Language != "go" {
		t.Fatalf("expected go sample, got %s", samples[0].Language)
	}
	if samples[0].Category != contracts.DatasetClassModuleLevel {
		t.Fatalf("unexpected category: %s", samples[0].Category)
	}
}

func TestDiscoverSamplesScenarioAndModuleLevelClass(t *testing.T) {
	root := t.TempDir()
	datasetRoot := filepath.Join(root, "datasets")
	path := filepath.Join(datasetRoot, "python", "python_code_files_module_level", "boundary")
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(path, "boundary_000.py"), []byte("def f():\n    return 1\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	spec := contracts.RunSpec{
		RunID:           "r2",
		DatasetRoot:     datasetRoot,
		OutputRoot:      root,
		ConfigPath:      "dummy",
		Languages:       []string{"python"},
		DatasetClasses:  []string{"module_level"},
		DatasetScenario: "boundary",
		MaxSamples:      10,
	}

	samples, err := svc.DiscoverSamples(spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(samples) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(samples))
	}
	if samples[0].Category != contracts.DatasetClassModuleLevel {
		t.Fatalf("expected class module_level, got %s", samples[0].Category)
	}
	if samples[0].Scenario != "boundary" {
		t.Fatalf("expected scenario boundary, got %s", samples[0].Scenario)
	}

	spec.DatasetClasses = []string{"self_contained"}
	samples2, err := svc.DiscoverSamples(spec)
	if err == nil {
		t.Fatalf("expected no samples for self_contained filter, got %d", len(samples2))
	}
}

func TestDiscoverSamplesManifestMaxSamplesPerLanguageScenario(t *testing.T) {
	root := t.TempDir()
	datasetRoot := filepath.Join(root, "dataset")
	if err := os.MkdirAll(filepath.Join(datasetRoot, "python"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(datasetRoot, "go"), 0o755); err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		"python/boundary_000.py":           "def f():\n    return 1\n",
		"python/boundary_001.py":           "def f():\n    return 2\n",
		"python/complex_dependency_000.py": "def f():\n    return 3\n",
		"python/complex_dependency_001.py": "def f():\n    return 4\n",
		"go/boundary_000.go":               "package main\n",
		"go/boundary_001.go":               "package main\n",
		"go/complex_dependency_000.go":     "package main\n",
		"go/complex_dependency_001.go":     "package main\n",
	}
	for relPath, content := range files {
		abs := filepath.Join(datasetRoot, relPath)
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	manifestPath := filepath.Join(root, "dataset_l1.json")
	manifest := `{
  "level": "l1",
  "samples": [
    {"id":"boundary_000","language":"python","category":"self_contained","path":"python/boundary_000.py","scenario":"boundary"},
    {"id":"boundary_001","language":"python","category":"self_contained","path":"python/boundary_001.py","scenario":"boundary"},
    {"id":"complex_dependency_000","language":"python","category":"self_contained","path":"python/complex_dependency_000.py","scenario":"complex_dependency"},
    {"id":"complex_dependency_001","language":"python","category":"self_contained","path":"python/complex_dependency_001.py","scenario":"complex_dependency"},
    {"id":"boundary_000","language":"go","category":"self_contained","path":"go/boundary_000.go","scenario":"boundary"},
    {"id":"boundary_001","language":"go","category":"self_contained","path":"go/boundary_001.go","scenario":"boundary"},
    {"id":"complex_dependency_000","language":"go","category":"self_contained","path":"go/complex_dependency_000.go","scenario":"complex_dependency"},
    {"id":"complex_dependency_001","language":"go","category":"self_contained","path":"go/complex_dependency_001.go","scenario":"complex_dependency"}
  ]
}`
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	svc := NewService()
	spec := contracts.RunSpec{
		RunID:           "r3",
		DatasetRoot:     datasetRoot,
		OutputRoot:      root,
		ConfigPath:      "dummy",
		DatasetManifest: manifestPath,
		Languages:       []string{"python", "go"},
		MaxSamples:      1,
	}

	samples, err := svc.DiscoverSamples(spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(samples) != 4 {
		t.Fatalf("expected 4 samples (2 langs x 2 scenarios x 1 each), got %d", len(samples))
	}

	type key struct {
		lang     string
		scenario string
	}
	counts := map[key]int{}
	for _, sample := range samples {
		counts[key{lang: sample.Language, scenario: sample.Scenario}]++
	}

	expected := []key{
		{lang: "python", scenario: "boundary"},
		{lang: "python", scenario: "complex_dependency"},
		{lang: "go", scenario: "boundary"},
		{lang: "go", scenario: "complex_dependency"},
	}
	for _, k := range expected {
		if counts[k] != 1 {
			t.Fatalf("expected 1 sample for %s/%s, got %d", k.lang, k.scenario, counts[k])
		}
	}
}
