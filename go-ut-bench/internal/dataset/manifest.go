package dataset

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"go-ut-bench/internal/contracts"
)

type datasetManifest struct {
	Level   string `json:"level"`
	Samples []struct {
		ID       string                 `json:"id"`
		Language string                 `json:"language"`
		Category contracts.DatasetClass `json:"category"`
		Scenario string                 `json:"scenario"`
		Path     string                 `json:"path"`
	} `json:"samples"`
}

func loadManifest(path string, datasetRoot string) (datasetManifest, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return datasetManifest{}, err
	}
	var mf datasetManifest
	if err := json.Unmarshal(raw, &mf); err != nil {
		return datasetManifest{}, err
	}
	for i := range mf.Samples {
		sample := &mf.Samples[i]
		if sample.Path == "" {
			sample.Path = filepath.Join(sample.Language, sample.ID+languageExtByName(sample.Language))
		}
		if !filepath.IsAbs(sample.Path) {
			candidate := filepath.Clean(sample.Path)
			if _, err := os.Stat(candidate); err == nil {
				sample.Path = candidate
			} else {
				sample.Path = filepath.Join(datasetRoot, sample.Path)
			}
		}
		sample.Path = filepath.Clean(sample.Path)
		sample.Language = strings.ToLower(strings.TrimSpace(sample.Language))
		sample.Scenario = strings.TrimSpace(sample.Scenario)
	}
	return mf, nil
}
