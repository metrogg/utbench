// runner 包提供测试生成功能
// 负责调用LLM API生成单元测试，支持多模型并行和checkpoint恢复
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

// Service 测试生成服务结构
// 提供完整的测试生成流程管理
type Service struct {
	logger *obs.Logger // 日志记录器
}

// Output 生成操作的输出结果
// 包含生成的测试清单和文件路径
type Output struct {
	Manifest     contracts.GeneratedManifest // 生成清单
	ManifestPath string                      // 清单文件路径
}

// task 生成任务结构
// 定义一个模型对一个样本的生成任务
type task struct {
	model  modelConfig         // 模型配置
	sample contracts.SampleRef // 样本引用
}

// NewService 创建新的生成服务实例
// 参数:
//   - logger: 日志记录器实例
//
// 返回值:
//   - *Service: 新的服务实例
func NewService(logger *obs.Logger) *Service {
	return &Service{logger: logger}
}

// Generate 执行测试生成流程
// 参数:
//   - ctx: 上下文，用于取消操作
//   - spec: 运行规格说明
//   - samples: 要处理的样本列表
//
// 返回值:
//   - Output: 生成结果输出
//   - error: 生成失败时的错误
//
// 功能说明:
//  1. 加载模型配置
//  2. 创建输出目录结构
//  3. 检查checkpoint（增量模式）
//  4. 使用worker池并行调用LLM API生成测试
//  5. 保存测试文件和元数据
//  6. 生成清单文件
func (s *Service) Generate(ctx context.Context, spec contracts.RunSpec, samples []contracts.SampleRef) (Output, error) {
	modelConfigs, err := loadModelConfigs(spec.ConfigPath, spec.Models)
	if err != nil {
		return Output{}, err
	}

	// 创建输出目录结构
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

	// 处理checkpoint
	checkpointPath := buildCheckpointPath(spec, modelConfigs)
	completed := map[string]struct{}{}
	if spec.ResetCheckpoint {
		_ = os.Remove(checkpointPath)
	}
	if spec.Mode == contracts.RunModeIncremental {
		completed, _ = loadCheckpoint(checkpointPath)
	}

	// 输出配置信息
	totalTasks := len(modelConfigs) * len(samples)
	workerCount := spec.Workers
	if workerCount <= 0 {
		workerCount = min(16, max(2, runtime.NumCPU()))
	}
	progress := obs.NewProgressReporter(totalTasks, "generate")
	progress.PrintStageStart("生成测试", fmt.Sprintf("模型: %s | 样本: %d | Workers: %d",
		strings.Join(getModelNames(modelConfigs), ", "), len(samples), workerCount))

	// 创建worker池
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

	// 发送任务，处理checkpoint过滤
	// 使用轮询方式分配任务，确保并发时每个worker处理不同模型的任务
	skippedByCheckpoint := 0
	go func() {
		defer close(tasks)

		// 为每个模型创建一个样本迭代器
		type modelIterator struct {
			model   modelConfig
			samples []contracts.SampleRef
			index   int
		}
		iterators := make([]modelIterator, len(modelConfigs))
		for i, model := range modelConfigs {
			iterators[i] = modelIterator{model: model, samples: samples, index: 0}
		}

		// 轮询分配任务
		activeModels := len(iterators)
		for activeModels > 0 {
			activeModels = 0
			for i := range iterators {
				it := &iterators[i]
				// 跳过已完成的模型
				for it.index < len(it.samples) {
					sample := it.samples[it.index]
					it.index++

					// 检查checkpoint
					if spec.Mode == contracts.RunModeIncremental {
						key := taskKey(it.model.Name, sample.Language, sample.ID)
						if _, ok := completed[key]; ok {
							skippedByCheckpoint++
							continue
						}
					}

					// 发送任务
					select {
					case <-ctx.Done():
						return
					case tasks <- task{model: it.model, sample: sample}:
					}
					activeModels++
					break // 每个模型每次只发送一个任务
				}
				if it.index < len(it.samples) {
					activeModels++
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
	completedCount := 0
	for item := range results {
		completedCount++
		cases = append(cases, item)

		var tokens int
		if item.TotalTokens != nil {
			tokens = *item.TotalTokens
		}
		taskResult := obs.TaskResult{
			Model:     item.Model,
			Language:  item.Language,
			SampleID:  item.SampleID,
			Success:   item.Success,
			Truncated: item.Truncated,
			LatencyMS: item.LatencyMS,
			Tokens:    tokens,
		}
		if !item.Success && item.Error != nil {
			taskResult.Error = item.Error.Message
		}
		progress.OnTaskDone(taskResult)

		status := "OK"
		if item.Truncated {
			status = "TRUNC"
		}
		if !item.Success {
			status = "FAIL"
			if item.Error != nil {
				status = fmt.Sprintf("FAIL(%s)", trimErrorMsg(item.Error.Message, 30))
			}
			if item.Truncated {
				status = "FAIL(truncated)"
			}
		}

		extra := ""
		if tokens > 0 {
			extra = fmt.Sprintf("| %d tokens", tokens)
		}
		progress.PrintTaskLine(completedCount, totalTasks-skippedByCheckpoint, item.Model, item.Language, item.SampleID, status, extra)

		if completedCount%5 == 0 {
			progress.PrintStats()
		}

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

	if skippedByCheckpoint > 0 {
		fmt.Fprintf(os.Stderr, "\n[Runner] Skipped %d tasks (already completed in checkpoint)\n", skippedByCheckpoint)
	}

	progress.PrintStats()
	progress.PrintStageDone("生成测试", obs.StageStats{
		Total:    totalTasks - skippedByCheckpoint,
		Success:  completedCount - skippedByCheckpoint,
		Duration: time.Since(progress.GetStartTime()),
	})

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
	var truncated bool
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
		generated, response, latency, pTok, cTok, tTok, isTruncated, genErr := client.generateTest(
			ctx,
			modelCfg,
			sample.Language,
			sample.Path,
			string(sourceCode),
		)
		truncated = isTruncated
		if genErr != nil {
			_ = contracts.WriteJSON(respPath, map[string]any{"error": genErr, "truncated": truncated})
			return contracts.GeneratedCase{
				Model:             model,
				Language:          sample.Language,
				SampleID:          sample.ID,
				SamplePath:        sample.Path,
				GeneratedTestPath: testPath,
				ResponsePath:      respPath,
				GeneratedAtUTC:    time.Now().UTC(),
				Success:           false,
				Truncated:         truncated,
				Error:             genErr,
			}
		}
		content = generated
		rawResponse = response
		promptTokens = pTok
		completionTokens = cTok
		totalTokens = tTok
		latencyMS = latency

		s.logger.LogAPIRequest(model, sample.Language, sample.ID, 0, latencyMS)
		s.logger.LogAPIResponse(model, sample.Language, sample.ID, genErr == nil, truncated, errorMsgSafe(genErr))
		s.logger.ToFile("runner").Trace("generate_response",
			"model", model,
			"language", sample.Language,
			"sample_id", sample.ID,
			"prompt_tokens", pTok,
			"completion_tokens", cTok,
			"total_tokens", tTok,
			"latency_ms", latencyMS,
			"truncated", truncated,
			"success", genErr == nil,
		)
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
		"truncated":      truncated,
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
		Truncated:         truncated,
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
		strings.Join(spec.DatasetClasses, ","),
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

func getModelNames(configs []modelConfig) []string {
	names := make([]string, 0, len(configs))
	for _, c := range configs {
		names = append(names, c.Name)
	}
	return names
}

func getLanguagesFromSamples(samples []contracts.SampleRef) string {
	langs := make(map[string]int)
	for _, s := range samples {
		langs[s.Language]++
	}
	var parts []string
	for _, l := range []string{"python", "go", "java", "cpp"} {
		if langs[l] > 0 {
			parts = append(parts, fmt.Sprintf("%s:%d", l, langs[l]))
		}
	}
	return strings.Join(parts, ", ")
}

func trimErrorMsg(msg string, max int) string {
	msg = strings.TrimSpace(msg)
	if len(msg) <= max {
		return msg
	}
	return msg[:max] + "..."
}

func errorMsgSafe(err *contracts.ErrorInfo) string {
	if err == nil {
		return ""
	}
	return trimErrorMsg(err.Message, 100)
}
