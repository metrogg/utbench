package runner

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/obs"
)

type Service struct {
	logger *obs.Logger
}

type Output struct {
	Manifest     contracts.GeneratedManifest
	ManifestPath string
}

type task struct {
	model  modelConfig
	sample contracts.SampleRef
}

func NewService(logger *obs.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) Generate(ctx context.Context, spec contracts.RunSpec, samples []contracts.SampleRef) (Output, error) {
	modelConfigs, err := loadModelConfigs(spec.ConfigPath, spec.Models)
	if err != nil {
		return Output{}, err
	}

	runRoot := filepath.Join(spec.OutputRoot, "runs", spec.RunID)
	genRoot := filepath.Join(runRoot, "generated")
	testRoot := filepath.Join(genRoot, "tests")
	metaRoot := filepath.Join(genRoot, "metadata")
	if err := os.MkdirAll(testRoot, 0o755); err != nil {
		return Output{}, err
	}
	if err := os.MkdirAll(metaRoot, 0o755); err != nil {
		return Output{}, err
	}

	checkpointPath := buildCheckpointPath(spec, modelConfigs)
	completed := map[string]struct{}{}
	if spec.ResetCheckpoint {
		_ = os.Remove(checkpointPath)
	}
	if spec.Mode == contracts.RunModeIncremental {
		completed, _ = loadCheckpoint(checkpointPath)
	}

	workerCount := min(16, max(2, runtime.NumCPU()))
	tasks := make(chan task)
	results := make(chan contracts.GeneratedCase)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range tasks {
				item := s.generateOne(ctx, spec, testRoot, metaRoot, t.model, t.sample)
				select {
				case <-ctx.Done():
					return
				case results <- item:
				}
			}
		}()
	}

	totalTasks := len(modelConfigs) * len(samples)
	skippedByCheckpoint := 0
	go func() {
		defer close(tasks)
		for _, model := range modelConfigs {
			for _, sample := range samples {
				if spec.Mode == contracts.RunModeIncremental {
					key := taskKey(model.Name, sample.Language, sample.ID)
					if _, ok := completed[key]; ok {
						skippedByCheckpoint++
						continue
					}
				}
				select {
				case <-ctx.Done():
					return
				case tasks <- task{model: model, sample: sample}:
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	cases := make([]contracts.GeneratedCase, 0, totalTasks)
	var ckptMu sync.Mutex
	for item := range results {
		cases = append(cases, item)
		if spec.Mode == contracts.RunModeIncremental && item.Success {
			key := taskKey(item.Model, item.Language, item.SampleID)
			ckptMu.Lock()
			completed[key] = struct{}{}
			_ = saveCheckpoint(checkpointPath, completed)
			ckptMu.Unlock()
		}
	}
	if err := ctx.Err(); err != nil {
		return Output{}, err
	}

	sort.Slice(cases, func(i, j int) bool {
		if cases[i].Model == cases[j].Model {
			if cases[i].Language == cases[j].Language {
				return cases[i].SampleID < cases[j].SampleID
			}
			return cases[i].Language < cases[j].Language
		}
		return cases[i].Model < cases[j].Model
	})

	manifest := contracts.GeneratedManifest{
		SchemaVersion: contracts.SchemaVersion,
		RunID:         spec.RunID,
		CreatedAtUTC:  time.Now().UTC(),
		Spec:          spec,
		Cases:         cases,
	}
	manifestPath := filepath.Join(genRoot, "generated_manifest.json")
	if err := contracts.WriteJSON(manifestPath, manifest); err != nil {
		return Output{}, err
	}

	s.logger.Info(
		"generate finished",
		"run_id", spec.RunID,
		"total_cases", len(cases),
		"pending_tasks", len(cases),
		"skipped_by_checkpoint", skippedByCheckpoint,
		"checkpoint", checkpointPath,
		"manifest", manifestPath,
	)
	return Output{Manifest: manifest, ManifestPath: manifestPath}, nil
}

func (s *Service) generateOne(ctx context.Context, spec contracts.RunSpec, testRoot, metaRoot string, modelCfg modelConfig, sample contracts.SampleRef) contracts.GeneratedCase {
	model := modelCfg.Name
	started := time.Now()
	ext := languageExt(sample.Language)
	testRel := filepath.Join(model, sample.Language, fmt.Sprintf("%s.test%s", sample.ID, ext))
	testPath := filepath.Join(testRoot, testRel)
	respPath := filepath.Join(metaRoot, fmt.Sprintf("%s_%s_%s.response.json", model, sample.Language, sample.ID))

	if spec.Mode == contracts.RunModeIncremental {
		if _, err := os.Stat(testPath); err == nil {
			latency := int(time.Since(started).Milliseconds())
			if _, err := os.Stat(respPath); err != nil {
				respPath = ""
			}
			return contracts.GeneratedCase{
				Model:             model,
				Language:          sample.Language,
				SampleID:          sample.ID,
				SamplePath:        sample.Path,
				GeneratedTestPath: testPath,
				ResponsePath:      respPath,
				MetadataPath:      "",
				LatencyMS:         latency,
				GeneratedAtUTC:    time.Now().UTC(),
				Success:           true,
			}
		}
	}

	if err := os.MkdirAll(filepath.Dir(testPath), 0o755); err != nil {
		return contracts.GeneratedCase{
			Model:             model,
			Language:          sample.Language,
			SampleID:          sample.ID,
			SamplePath:        sample.Path,
			GeneratedTestPath: testPath,
			ResponsePath:      "",
			GeneratedAtUTC:    time.Now().UTC(),
			Success:           false,
			Error: &contracts.ErrorInfo{
				Kind:      "write_error",
				Message:   err.Error(),
				Retryable: false,
			},
		}
	}

	content := ""
	var rawResponse map[string]any
	var promptTokens *int
	var completionTokens *int
	var totalTokens *int
	latencyMS := 0

	if spec.DryRun {
		content = buildPlaceholderTest(sample.Language, sample.ID)
	} else {
		sourceCode, readErr := os.ReadFile(sample.Path)
		if readErr != nil {
			return contracts.GeneratedCase{
				Model:             model,
				Language:          sample.Language,
				SampleID:          sample.ID,
				SamplePath:        sample.Path,
				GeneratedTestPath: testPath,
				GeneratedAtUTC:    time.Now().UTC(),
				Success:           false,
				Error: &contracts.ErrorInfo{
					Kind:      "sample_read_error",
					Message:   readErr.Error(),
					Retryable: false,
				},
			}
		}

		client := newAPIClient()
		generated, response, latency, pTok, cTok, tTok, genErr := client.generateTest(
			ctx,
			modelCfg,
			sample.Language,
			sample.Path,
			string(sourceCode),
		)
		if genErr != nil {
			_ = contracts.WriteJSON(respPath, map[string]any{"error": genErr})
			return contracts.GeneratedCase{
				Model:             model,
				Language:          sample.Language,
				SampleID:          sample.ID,
				SamplePath:        sample.Path,
				GeneratedTestPath: testPath,
				ResponsePath:      respPath,
				GeneratedAtUTC:    time.Now().UTC(),
				Success:           false,
				Error:             genErr,
			}
		}
		content = generated
		rawResponse = response
		promptTokens = pTok
		completionTokens = cTok
		totalTokens = tTok
		latencyMS = latency
	}

	if err := os.WriteFile(testPath, []byte(content), 0o644); err != nil {
		return contracts.GeneratedCase{
			Model:             model,
			Language:          sample.Language,
			SampleID:          sample.ID,
			SamplePath:        sample.Path,
			GeneratedTestPath: testPath,
			GeneratedAtUTC:    time.Now().UTC(),
			Success:           false,
			Error:             &contracts.ErrorInfo{Kind: "write_error", Message: err.Error(), Retryable: false},
		}
	}

	if spec.DryRun {
		rawResponse = map[string]any{"dry_run": true}
	}
	_ = contracts.WriteJSON(respPath, rawResponse)
	latencyForMeta := latencyMS
	if latencyForMeta == 0 {
		latencyForMeta = int(time.Since(started).Milliseconds())
	}

	metadataPath := filepath.Join(metaRoot, fmt.Sprintf("%s_%s_%s.metadata.json", model, sample.Language, sample.ID))
	metadata := map[string]any{
		"model":               model,
		"language":            sample.Language,
		"sample_id":           sample.ID,
		"sample_path":         sample.Path,
		"scenario":            sample.Scenario,
		"generated_test_path": testPath,
		"response_path":       respPath,
		"dataset_class":       sample.Category,
		"source_md5":          sample.SourceMD5,
		"latency_ms":          latencyForMeta,
		"tokens": map[string]any{
			"prompt_tokens":     promptTokens,
			"completion_tokens": completionTokens,
			"total_tokens":      totalTokens,
		},
		"created_at_utc": time.Now().UTC(),
		"success":        true,
	}
	_ = contracts.WriteJSON(metadataPath, metadata)

	latency := latencyMS
	if latency == 0 {
		latency = int(time.Since(started).Milliseconds())
	}
	return contracts.GeneratedCase{
		Model:             model,
		Language:          sample.Language,
		SampleID:          sample.ID,
		SamplePath:        sample.Path,
		GeneratedTestPath: testPath,
		ResponsePath:      respPath,
		MetadataPath:      metadataPath,
		LatencyMS:         latency,
		PromptTokens:      promptTokens,
		CompletionTokens:  completionTokens,
		TotalTokens:       totalTokens,
		GeneratedAtUTC:    time.Now().UTC(),
		Success:           true,
	}
}

func languageExt(language string) string {
	switch strings.ToLower(language) {
	case "python":
		return ".py"
	case "java":
		return ".java"
	case "go":
		return ".go"
	case "cpp":
		return ".cpp"
	default:
		return ".txt"
	}
}

func buildPlaceholderTest(language, sampleID string) string {
	switch language {
	case "python":
		return fmt.Sprintf("import pytest\n\n\ndef test_placeholder_%s():\n    assert True\n", sanitizeIdentifier(sampleID))
	case "java":
		return "public class PlaceholderTest { public void testPlaceholder() { assert true; } }\n"
	case "go":
		return "package main\n\nimport \"testing\"\n\nfunc TestPlaceholder(t *testing.T) {}\n"
	case "cpp":
		return "#include <cassert>\nint main() { assert(true); return 0; }\n"
	default:
		return "placeholder test\n"
	}
}

func sanitizeIdentifier(raw string) string {
	raw = strings.ToLower(raw)
	b := strings.Builder{}
	for _, r := range raw {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
			continue
		}
		b.WriteRune('_')
	}
	if b.Len() == 0 {
		return "sample"
	}
	return b.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func taskKey(model, language, sampleID string) string {
	return model + "|" + language + "|" + sampleID
}

func buildCheckpointPath(spec contracts.RunSpec, models []modelConfig) string {
	modelNames := make([]string, 0, len(models))
	for _, item := range models {
		modelNames = append(modelNames, item.Name)
	}
	sort.Strings(modelNames)

	langs := append([]string{}, spec.Languages...)
	for i := range langs {
		langs[i] = strings.ToLower(strings.TrimSpace(langs[i]))
	}
	sort.Strings(langs)

	scope := fmt.Sprintf(
		"models=%s;langs=%s;class=%s;level=%s;manifest=%s;max=%d;dataset=%s",
		strings.Join(modelNames, ","),
		strings.Join(langs, ","),
		spec.DatasetClass,
		spec.DatasetLevel,
		spec.DatasetManifest,
		spec.MaxSamples,
		spec.DatasetRoot,
	)
	h := sha1.Sum([]byte(scope))
	hash := hex.EncodeToString(h[:])[:12]
	return filepath.Join(spec.OutputRoot, "checkpoints", "runner_"+hash+".checkpoint.json")
}

func loadCheckpoint(path string) (map[string]struct{}, error) {
	result := map[string]struct{}{}
	raw, err := os.ReadFile(path)
	if err != nil {
		return result, nil
	}
	var payload struct {
		Completed []string `json:"completed"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return result, nil
	}
	for _, item := range payload.Completed {
		if item != "" {
			result[item] = struct{}{}
		}
	}
	return result, nil
}

func saveCheckpoint(path string, completed map[string]struct{}) error {
	items := make([]string, 0, len(completed))
	for item := range completed {
		items = append(items, item)
	}
	sort.Strings(items)
	payload := map[string]any{
		"updated_at_utc": time.Now().UTC(),
		"completed":      items,
	}
	return contracts.WriteJSON(path, payload)
}
