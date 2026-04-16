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
    {"id":"complex_dependency_go_complex_dependency_0","language":"go","category":"complex_dependency","path":"go/complex_dependency_go_complex_dependency_0.go"}
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
		DatasetClass:    contracts.DatasetClassComplexDependency,
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
	if samples[0].Category != contracts.DatasetClassComplexDependency {
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
		DatasetClass:    contracts.DatasetClassModuleLevel,
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

	spec.DatasetClass = contracts.DatasetClassComplexDependency
	samples2, err := svc.DiscoverSamples(spec)
	if err != nil {
		t.Fatal(err)
	}
	if len(samples2) != 1 {
		t.Fatalf("expected backward-compatible class filter result, got %d", len(samples2))
	}
}
