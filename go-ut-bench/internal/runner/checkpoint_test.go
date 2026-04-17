package runner

import (
	"path/filepath"
	"testing"

	"go-ut-bench/internal/contracts"
)

func TestCheckpointPathIncludesScopeFields(t *testing.T) {
	spec1 := contracts.RunSpec{
		OutputRoot:      filepath.Join(t.TempDir(), "out"),
		DatasetRoot:     "../dataset",
		DatasetClass:    contracts.DatasetClassSelfContained,
		DatasetLevel:    "l1",
		DatasetManifest: "./configs/dataset_index.json",
		MaxSamples:      10,
		Languages:       []string{"python"},
	}
	spec2 := spec1
	spec2.DatasetLevel = "l2"

	models := []modelConfig{{Name: "deepseek"}}
	p1 := buildCheckpointPath(spec1, models)
	p2 := buildCheckpointPath(spec2, models)
	if p1 == p2 {
		t.Fatalf("expected different checkpoint paths for different levels")
	}
}

func TestCheckpointRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "runner.checkpoint.json")
	completed := map[string]struct{}{
		"m|python|s1": {},
		"m|python|s2": {},
	}
	if err := saveCheckpoint(path, completed); err != nil {
		t.Fatal(err)
	}
	loaded, err := loadCheckpoint(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 checkpoint entries, got %d", len(loaded))
	}
}
