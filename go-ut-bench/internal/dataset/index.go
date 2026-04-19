package dataset

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"
)

type indexSample struct {
	ID        string                 `json:"id"`
	Language  string                 `json:"language"`
	Category  contracts.DatasetClass `json:"category"`
	Scenario  string                 `json:"scenario"`
	Path      string                 `json:"path"`
	SourceMD5 string                 `json:"source_md5"`
}

type datasetIndexFile struct {
	SchemaVersion string        `json:"schema_version"`
	GeneratedAt   time.Time     `json:"generated_at_utc"`
	DatasetRoot   string        `json:"dataset_root"`
	Samples       []indexSample `json:"samples"`
}

type BuildSummary struct {
	Path  string
	Total int
}

type ManifestBuildOptions struct {
	IndexPath        string
	Level            string
	OutputPath       string
	Languages        []string
	ClassFilter      string
	ScenarioFilter   string
	LimitPerScenario int
}

func (s *Service) BuildIndex(datasetRoot, outputPath string) (BuildSummary, error) {
	spec := contracts.RunSpec{
		RunID:       contracts.NewRunID(),
		DatasetRoot: strings.TrimSpace(datasetRoot),
		Languages:   append([]string{}, contracts.SupportedLanguages...),
		MaxSamples:  0,
	}
	samples, err := s.DiscoverSamples(spec)
	if err != nil {
		return BuildSummary{}, err
	}

	rows := make([]indexSample, 0, len(samples))
	for _, item := range samples {
		rows = append(rows, indexSample{
			ID:        item.ID,
			Language:  item.Language,
			Category:  item.Category,
			Scenario:  item.Scenario,
			Path:      item.Path,
			SourceMD5: item.SourceMD5,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Language == rows[j].Language {
			if rows[i].Category == rows[j].Category {
				if rows[i].Scenario == rows[j].Scenario {
					return rows[i].ID < rows[j].ID
				}
				return rows[i].Scenario < rows[j].Scenario
			}
			return rows[i].Category < rows[j].Category
		}
		return rows[i].Language < rows[j].Language
	})

	payload := datasetIndexFile{
		SchemaVersion: contracts.SchemaVersion,
		GeneratedAt:   time.Now().UTC(),
		DatasetRoot:   spec.DatasetRoot,
		Samples:       rows,
	}
	if err := contracts.WriteJSON(outputPath, payload); err != nil {
		return BuildSummary{}, err
	}
	return BuildSummary{Path: outputPath, Total: len(rows)}, nil
}

func (s *Service) BuildManifest(opts ManifestBuildOptions) (BuildSummary, error) {
	idx, err := readDatasetIndex(opts.IndexPath)
	if err != nil {
		return BuildSummary{}, err
	}

	langs := map[string]struct{}{}
	for _, lang := range opts.Languages {
		lang = strings.ToLower(strings.TrimSpace(lang))
		if lang != "" {
			langs[lang] = struct{}{}
		}
	}
	classFilter := strings.TrimSpace(opts.ClassFilter)
	scenarioFilter := normalizeScenario(opts.ScenarioFilter)
	limit := opts.LimitPerScenario
	if limit <= 0 {
		limit = 20
	}

	type groupKey struct {
		Lang     string
		Category contracts.DatasetClass
		Scenario string
	}
	grouped := map[groupKey][]indexSample{}

	for _, item := range idx.Samples {
		if len(langs) > 0 {
			if _, ok := langs[item.Language]; !ok {
				continue
			}
		}
		if classFilter != "" {
			classFilters := strings.Split(classFilter, ",")
			for i := range classFilters {
				classFilters[i] = strings.TrimSpace(classFilters[i])
			}
			if !matchDatasetClassFilter(classFilters, item.Category) {
				continue
			}
		}
		if scenarioFilter != "" && item.Scenario != scenarioFilter {
			continue
		}
		k := groupKey{Lang: item.Language, Category: item.Category, Scenario: item.Scenario}
		grouped[k] = append(grouped[k], item)
	}

	keys := make([]groupKey, 0, len(grouped))
	for k := range grouped {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Lang == keys[j].Lang {
			if keys[i].Category == keys[j].Category {
				return keys[i].Scenario < keys[j].Scenario
			}
			return keys[i].Category < keys[j].Category
		}
		return keys[i].Lang < keys[j].Lang
	})

	samples := make([]map[string]any, 0)
	for _, key := range keys {
		rows := grouped[key]
		sort.Slice(rows, func(i, j int) bool { return rows[i].ID < rows[j].ID })
		if len(rows) > limit {
			rows = rows[:limit]
		}
		for _, row := range rows {
			samples = append(samples, map[string]any{
				"id":       row.ID,
				"language": row.Language,
				"category": row.Category,
				"scenario": row.Scenario,
				"path":     row.Path,
			})
		}
	}

	if len(samples) == 0 {
		return BuildSummary{}, fmt.Errorf("no samples selected from index")
	}

	level := strings.TrimSpace(opts.Level)
	if level == "" {
		level = "l1"
	}
	payload := map[string]any{
		"level":   level,
		"samples": samples,
	}
	if err := contracts.WriteJSON(opts.OutputPath, payload); err != nil {
		return BuildSummary{}, err
	}
	return BuildSummary{Path: opts.OutputPath, Total: len(samples)}, nil
}

func readDatasetIndex(path string) (datasetIndexFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return datasetIndexFile{}, err
	}
	var out datasetIndexFile
	if err := json.Unmarshal(raw, &out); err != nil {
		return datasetIndexFile{}, err
	}
	for i := range out.Samples {
		if !filepath.IsAbs(out.Samples[i].Path) {
			out.Samples[i].Path = filepath.Clean(out.Samples[i].Path)
		}
	}
	return out, nil
}
