package dataset

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildIndexAndManifest(t *testing.T) {
	root := t.TempDir()
	datasetRoot := filepath.Join(root, "datasets")
	paths := []string{
		filepath.Join(datasetRoot, "python", "python_code_files_self_contained", "boundary"),
		filepath.Join(datasetRoot, "python", "python_code_files_module_level", "boundary"),
	}
	for _, p := range paths {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	_ = os.WriteFile(filepath.Join(paths[0], "boundary_000.py"), []byte("def f():\n    return 1\n"), 0o644)
	_ = os.WriteFile(filepath.Join(paths[0], "boundary_001.py"), []byte("def g():\n    return 2\n"), 0o644)
	_ = os.WriteFile(filepath.Join(paths[1], "boundary_000.py"), []byte("def h():\n    return 3\n"), 0o644)

	svc := NewService()
	indexPath := filepath.Join(root, "dataset_index.json")
	idxSummary, err := svc.BuildIndex(datasetRoot, indexPath)
	if err != nil {
		t.Fatal(err)
	}
	if idxSummary.Total != 3 {
		t.Fatalf("expected 3 indexed samples, got %d", idxSummary.Total)
	}

	manifestPath := filepath.Join(root, "dataset_l1.json")
	mSummary, err := svc.BuildManifest(ManifestBuildOptions{
		IndexPath:        indexPath,
		Level:            "l1",
		OutputPath:       manifestPath,
		Languages:        []string{"python"},
		ClassFilter:      "self_contained",
		ScenarioFilter:   "boundary",
		LimitPerScenario: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if mSummary.Total != 1 {
		t.Fatalf("expected 1 sample in manifest, got %d", mSummary.Total)
	}

	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		t.Fatal(err)
	}
	samples, ok := payload["samples"].([]any)
	if !ok || len(samples) != 1 {
		t.Fatalf("manifest samples invalid: %v", payload["samples"])
	}
}
