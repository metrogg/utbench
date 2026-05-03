// runner 包提供测试生成功能
// 负责调用LLM API生成单元测试，支持多模型并行和checkpoint恢复
package runner

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"go-ut-bench/internal/agentconfig"
	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/ctrl"
	"go-ut-bench/internal/obs"
	"go-ut-bench/internal/store"
)

// Service 测试生成服务结构
// 提供完整的测试生成流程管理
type Service struct {
	logger        *obs.Logger   // 日志记录器
	sandboxRunner SandboxRunner // Agent 样本级沙箱执行器
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
	subject subjectTarget       // 被测对象配置
	sample  contracts.SampleRef // 样本引用
}

type subjectTarget struct {
	subject agentconfig.ResolvedSubject
	model   modelConfig
}

// NewService 创建新的生成服务实例
// 参数:
//   - logger: 日志记录器实例
//
// 返回值:
//   - *Service: 新的服务实例
func NewService(logger *obs.Logger) *Service {
	return &Service{logger: logger, sandboxRunner: NewSandboxRunner()}
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
	subjects, err := loadSubjectTargets(spec, modelConfigs)
	if err != nil {
		return Output{}, err
	}
	var reuseStore *store.SQLiteStore
	if spec.ReuseGenerated && strings.TrimSpace(spec.DBPath) != "" && !spec.DryRun {
		if db, openErr := store.OpenSQLite(spec.DBPath); openErr == nil {
			if initErr := db.Init(ctx); initErr == nil {
				reuseStore = db
				defer reuseStore.Close()
			} else {
				_ = db.Close()
				s.logger.Warn("reuse generated disabled", "reason", initErr.Error())
			}
		} else {
			s.logger.Warn("reuse generated disabled", "reason", openErr.Error())
		}
	}

	// 创建输出目录结构
	runRoot := filepath.Join(spec.OutputRoot, "runs", spec.RunID)
	genRoot := filepath.Join(runRoot, "generated")
	testRoot := filepath.Join(genRoot, "tests")
	metaRoot := filepath.Join(genRoot, "metadata")
	promptRoot := filepath.Join(genRoot, "prompts")
	if err := os.MkdirAll(testRoot, 0o755); err != nil {
		return Output{}, err
	}
	if err := os.MkdirAll(metaRoot, 0o755); err != nil {
		return Output{}, err
	}
	promptCatalog, err := WritePromptCatalog(promptRoot)
	if err != nil {
		return Output{}, err
	}

	// 处理checkpoint
	checkpointPath := buildCheckpointPath(spec, subjects)
	completed := map[string]struct{}{}
	if spec.ResetCheckpoint {
		_ = os.Remove(checkpointPath)
	}
	if spec.Mode == contracts.RunModeIncremental {
		completed, _ = loadCheckpoint(checkpointPath)
	}

	// 输出配置信息
	totalTasks := countEligibleSubjectTasks(subjects, samples)
	if totalTasks == 0 {
		return Output{}, fmt.Errorf("no eligible subject/sample tasks after compatibility filtering")
	}
	workerCount := spec.Workers
	if workerCount <= 0 {
		workerCount = min(16, max(2, runtime.NumCPU()))
	}
	progress := obs.NewProgressReporter(totalTasks, "generate")
	// 构建 subject 列表详情
	subjectLines := ""
	for _, sub := range subjects {
		kind := sub.subject.Spec.Kind
		if kind == "" {
			kind = "model_api"
		}
		skill := sub.subject.Spec.Skill
		if skill == "" {
			skill = "no_skill"
		}
		subjectLines += fmt.Sprintf("   [%s] %s | model=%s | skill=%s\n", kind, sub.subject.Spec.ID, sub.subject.Spec.Model, skill)
	}
	progress.PrintStageStart("生成测试", fmt.Sprintf("%d 样本 × %d 被测对象 = %d 任务 | Workers: %d\n%s",
		len(samples), len(subjects), totalTasks, workerCount, subjectLines))

	// 创建worker池
	tasks := make(chan task)
	results := make(chan contracts.GeneratedCase)

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range tasks {
				// Web-triggered pause gate; it blocks before starting the next API call.
				if err := ctrl.Wait(ctx); err != nil {
					return
				}
				item := s.generateOne(ctx, spec, testRoot, metaRoot, promptRoot, promptCatalog.VersionID, reuseStore, t.subject, t.sample)
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

		// 为每个被测对象创建一个样本迭代器
		type subjectIterator struct {
			subject subjectTarget
			samples []contracts.SampleRef
			index   int
		}
		iterators := make([]subjectIterator, len(subjects))
		for i, subject := range subjects {
			iterators[i] = subjectIterator{subject: subject, samples: samples, index: 0}
		}

		// 轮询分配任务
		activeSubjects := len(iterators)
		for activeSubjects > 0 {
			activeSubjects = 0
			for i := range iterators {
				it := &iterators[i]
				// 跳过已完成的被测对象
				for it.index < len(it.samples) {
					sample := it.samples[it.index]
					it.index++
					if !subjectSupportsLanguage(it.subject, sample.Language) {
						continue
					}

					// 检查checkpoint
					if spec.Mode == contracts.RunModeIncremental {
						key := taskKey(it.subject.subject.Spec.ID, sample.Language, sample.ID)
						if _, ok := completed[key]; ok {
							skippedByCheckpoint++
							continue
						}
					}

					// 发送任务
					select {
					case <-ctx.Done():
						return
					case tasks <- task{subject: it.subject, sample: sample}:
					}
					activeSubjects++
					break // 每个模型每次只发送一个任务
				}
				if it.index < len(it.samples) {
					activeSubjects++
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
	successCount := 0
	failedCount := 0
	for item := range results {
		completedCount++
		cases = append(cases, item)
		if item.Success {
			successCount++
		} else {
			failedCount++
		}

		var tokens int
		if item.TotalTokens != nil {
			tokens = *item.TotalTokens
		}
		taskResult := obs.TaskResult{
			Model:            item.Model,
			Language:         item.Language,
			SampleID:         item.SampleID,
			Success:          item.Success,
			Truncated:        item.Truncated,
			LatencyMS:        item.LatencyMS,
			Tokens:           tokens,
			SubjectID:        item.SubjectID,
			SubjectKind:      item.SubjectKind,
			AgentFramework:   item.AgentFramework,
			SkillName:        item.SkillName,
			InteractionCount: item.InteractionCount,
			ToolCallCount:    item.ToolCallCount,
			FilesRead:        item.FilesReadCount,
			FilesWritten:     item.FilesWriteCount,
			CommandsExecuted: item.CommandCount,
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

		// 使用 Agent 感知的任务行显示
		progress.PrintAgentTaskLine(completedCount, totalTasks-skippedByCheckpoint, taskResult, status)

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
		Total:    completedCount,
		Success:  successCount,
		Failed:   failedCount,
		Skipped:  skippedByCheckpoint,
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
		SchemaVersion:     contracts.SchemaVersion,
		RunID:             spec.RunID,
		CreatedAtUTC:      time.Now().UTC(),
		Spec:              spec,
		PromptStrategy:    promptCatalog.Strategy,
		PromptVersionID:   promptCatalog.VersionID,
		PromptSnapshotDir: promptRoot,
		Cases:             cases,
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

// generateOne 执行单个样本的测试生成
// 处理文件准备、prompt 构建、API 调用、结果保存等完整流程
//
// 参数:
//   - ctx: 上下文
//   - spec: 运行规格
//   - testRoot: 测试文件输出目录
//   - metaRoot: 元数据输出目录
//   - promptRoot: prompt 输出目录
//   - promptVersionID: prompt 版本标识
//   - reuseStore: 复用数据库（可选）
//   - modelCfg: 模型配置
//   - sample: 样本引用
//
// 返回值:
//   - GeneratedCase: 生成结果
func (s *Service) generateOne(ctx context.Context, spec contracts.RunSpec, testRoot, metaRoot, promptRoot, promptVersionID string, reuseStore *store.SQLiteStore, target subjectTarget, sample contracts.SampleRef) contracts.GeneratedCase {
	subject := target.subject.Spec
	model := subject.ID
	started := time.Now()
	skillVersion := target.subject.Skill.Version
	ext := languageExt(sample.Language)
	testRel := filepath.Join(model, sample.Language, fmt.Sprintf("%s.test%s", sample.ID, ext))
	testPath := filepath.Join(testRoot, testRel)
	respPath := filepath.Join(metaRoot, fmt.Sprintf("%s_%s_%s.response.json", model, sample.Language, sample.ID))
	promptPath := ""
	promptPathCandidate := filepath.Join(promptRoot, "rendered", model, sample.Language, fmt.Sprintf("%s.prompt.txt", sample.ID))
	promptMode := string(PromptModeFullFile)
	if loadRepoLevelMetaForRunner(sample.Path) != nil {
		promptMode = string(PromptModeRepoLevel)
	}

	if spec.Mode == contracts.RunModeIncremental {
		if _, err := os.Stat(testPath); err == nil {
			latency := int(time.Since(started).Milliseconds())
			if _, err := os.Stat(respPath); err != nil {
				respPath = ""
			}
			return contracts.GeneratedCase{
				Model:             model,
				SubjectID:         subject.ID,
				SubjectKind:       subject.Kind,
				AgentFramework:    subject.Framework,
				AgentModel:        subject.Model,
				SkillName:         subject.Skill,
				SkillVersion:      skillVersion,
				Language:          sample.Language,
				SampleID:          sample.ID,
				SamplePath:        sample.Path,
				PromptVersionID:   promptVersionID,
				PromptMode:        promptMode,
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
			SubjectID:         subject.ID,
			SubjectKind:       subject.Kind,
			AgentFramework:    subject.Framework,
			AgentModel:        subject.Model,
			SkillName:         subject.Skill,
			SkillVersion:      skillVersion,
			Language:          sample.Language,
			SampleID:          sample.ID,
			SamplePath:        sample.Path,
			PromptVersionID:   promptVersionID,
			PromptMode:        promptMode,
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
	renderedPrompt := ""
	trace := subjectTrace{}
	var agentSummary agentTraceSummary
	identity := generationIdentity{}

	if spec.DryRun {
		content = buildPlaceholderTest(sample.Language, sample.ID)
	} else {
		sourceCode, readErr := os.ReadFile(sample.Path)
		if readErr != nil {
			return contracts.GeneratedCase{
				Model:             model,
				SubjectID:         subject.ID,
				SubjectKind:       subject.Kind,
				AgentFramework:    subject.Framework,
				AgentModel:        subject.Model,
				SkillName:         subject.Skill,
				SkillVersion:      skillVersion,
				Language:          sample.Language,
				SampleID:          sample.ID,
				SamplePath:        sample.Path,
				PromptVersionID:   promptVersionID,
				PromptMode:        promptMode,
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
		renderedPrompt = buildPrompt(sample.Language, sample.Path, string(sourceCode))
		if err := os.MkdirAll(filepath.Dir(promptPathCandidate), 0o755); err == nil {
			if err := os.WriteFile(promptPathCandidate, []byte(renderedPrompt), 0o644); err == nil {
				promptPath = promptPathCandidate
			}
		}

		sourceSHA := sha256Bytes(sourceCode)
		identity = buildGenerationIdentity(target, target.model, sample, sourceCode, renderedPrompt, promptVersionID)
		if reuseStore != nil {
			if reused, ok, reuseErr := reuseStore.FindReusableGeneratedAsset(ctx, identity.GenerationKey); reuseErr == nil && ok {
				if copyErr := copyFile(reused.GeneratedTestPath, testPath); copyErr == nil {
					metadataPath := filepath.Join(metaRoot, fmt.Sprintf("%s_%s_%s.metadata.json", model, sample.Language, sample.ID))
					metadata := map[string]any{
						"model":                        model,
						"subject_id":                   subject.ID,
						"subject_kind":                 subject.Kind,
						"agent_framework":              subject.Framework,
						"agent_model":                  subject.Model,
						"skill_name":                   subject.Skill,
						"skill_version":                skillVersion,
						"language":                     sample.Language,
						"sample_id":                    sample.ID,
						"sample_uid":                   identity.SampleUID,
						"sample_path":                  sample.Path,
						"prompt_strategy":              PromptStrategy(),
						"prompt_version_id":            promptVersionID,
						"prompt_mode":                  promptMode,
						"prompt_path":                  promptPath,
						"scenario":                     sample.Scenario,
						"generated_test_path":          testPath,
						"dataset_class":                sample.Category,
						"source_md5":                   sample.SourceMD5,
						"source_sha256":                sourceSHA,
						"subject_version_id":           identity.SubjectVersionID,
						"generation_key":               identity.GenerationKey,
						"dependency_fingerprint":       identity.DependencyFingerprint,
						"generation_env_fingerprint":   identity.GenerationEnvFingerprint,
						"reused":                       true,
						"reuse_stage":                  "generation",
						"reuse_key":                    identity.GenerationKey,
						"reuse_reason":                 "generation_key_match",
						"reused_from_run_id":           reused.RunID,
						"reused_from_case_id":          reused.GeneratedCaseID,
						"reused_generated_test_sha256": reused.GeneratedTestSHA256,
						"created_at_utc":               time.Now().UTC(),
						"success":                      true,
					}
					_ = contracts.WriteJSON(metadataPath, metadata)
					s.logger.Info("reuse generated test", "subject", model, "language", sample.Language, "sample_id", sample.ID, "from_run", reused.RunID)
					return contracts.GeneratedCase{
						Model:                    model,
						SubjectID:                subject.ID,
						SubjectKind:              subject.Kind,
						AgentFramework:           subject.Framework,
						AgentModel:               subject.Model,
						SkillName:                subject.Skill,
						SkillVersion:             skillVersion,
						Language:                 sample.Language,
						SampleID:                 sample.ID,
						SampleUID:                identity.SampleUID,
						SamplePath:               sample.Path,
						PromptVersionID:          promptVersionID,
						PromptMode:               promptMode,
						PromptPath:               promptPath,
						GeneratedTestPath:        testPath,
						ResponsePath:             reused.ResponsePath,
						MetadataPath:             metadataPath,
						TracePath:                reused.TracePath,
						WorkspaceDiffPath:        reused.WorkspaceDiffPath,
						SandboxFingerprint:       reused.SandboxFingerprint,
						SubjectVersionID:         identity.SubjectVersionID,
						GenerationKey:            identity.GenerationKey,
						DependencyFingerprint:    identity.DependencyFingerprint,
						GenerationEnvFingerprint: identity.GenerationEnvFingerprint,
						Reused:                   true,
						ReuseStage:               "generation",
						ReuseKey:                 identity.GenerationKey,
						ReuseReason:              "generation_key_match",
						ReusedFromRunID:          reused.RunID,
						ReusedFromCaseID:         reused.GeneratedCaseID,
						LatencyMS:                int(time.Since(started).Milliseconds()),
						PromptTokens:             reused.PromptTokens,
						CompletionTokens:         reused.CompletionTokens,
						TotalTokens:              reused.TotalTokens,
						TokenSource:              reused.TokenSource,
						EstimatedCostUSD:         reused.EstimatedCostUSD,
						CostSource:               reused.CostSource,
						GeneratedAtUTC:           time.Now().UTC(),
						Success:                  true,
					}
				}
			} else if reuseErr != nil {
				s.logger.Warn("reuse lookup failed", "subject", model, "language", sample.Language, "sample_id", sample.ID, "generation_key", identity.GenerationKey, "error", reuseErr.Error())
			}
		}

		generated, response, subjectTrace, latency, pTok, cTok, tTok, isTruncated, genErr, agentSmry := s.generateWithSubject(ctx, spec, target, sample, renderedPrompt, testPath, metaRoot)
		trace = subjectTrace
		agentSummary = agentSmry
		truncated = isTruncated
		if genErr != nil {
			_ = contracts.WriteJSON(respPath, map[string]any{"error": genErr, "truncated": truncated, "trace_path": trace.TracePath, "workspace_diff_path": trace.WorkspaceDiffPath})
			return contracts.GeneratedCase{
				Model:                    model,
				SubjectID:                subject.ID,
				SubjectKind:              subject.Kind,
				AgentFramework:           subject.Framework,
				AgentModel:               subject.Model,
				SkillName:                subject.Skill,
				SkillVersion:             skillVersion,
				Language:                 sample.Language,
				SampleID:                 sample.ID,
				SampleUID:                identity.SampleUID,
				SamplePath:               sample.Path,
				PromptVersionID:          promptVersionID,
				PromptMode:               promptMode,
				PromptPath:               promptPath,
				GeneratedTestPath:        testPath,
				ResponsePath:             respPath,
				LatencyMS:                latency,
				PromptTokens:             promptTokens,
				CompletionTokens:         completionTokens,
				TotalTokens:              totalTokens,
				TokenSource:              trace.TokenSource,
				EstimatedCostUSD:         trace.EstimatedCostUSD,
				CostSource:               trace.CostSource,
				TracePath:                trace.TracePath,
				WorkspaceDiffPath:        trace.WorkspaceDiffPath,
				SandboxFingerprint:       trace.SandboxFingerprint,
				SubjectVersionID:         identity.SubjectVersionID,
				GenerationKey:            identity.GenerationKey,
				DependencyFingerprint:    identity.DependencyFingerprint,
				GenerationEnvFingerprint: identity.GenerationEnvFingerprint,
				GeneratedAtUTC:           time.Now().UTC(),
				Success:                  false,
				Truncated:                truncated,
				Error:                    genErr,
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
			Model:              model,
			SubjectID:          subject.ID,
			SubjectKind:        subject.Kind,
			AgentFramework:     subject.Framework,
			AgentModel:         subject.Model,
			SkillName:          subject.Skill,
			SkillVersion:       skillVersion,
			Language:           sample.Language,
			SampleID:           sample.ID,
			SamplePath:         sample.Path,
			GeneratedTestPath:  testPath,
			TracePath:          trace.TracePath,
			WorkspaceDiffPath:  trace.WorkspaceDiffPath,
			SandboxFingerprint: trace.SandboxFingerprint,
			GeneratedAtUTC:     time.Now().UTC(),
			Success:            false,
			Error:              &contracts.ErrorInfo{Kind: "write_error", Message: err.Error(), Retryable: false},
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
		"model":                      model,
		"subject_id":                 subject.ID,
		"subject_kind":               subject.Kind,
		"agent_framework":            subject.Framework,
		"agent_model":                subject.Model,
		"skill_name":                 subject.Skill,
		"skill_version":              skillVersion,
		"language":                   sample.Language,
		"sample_id":                  sample.ID,
		"sample_uid":                 identity.SampleUID,
		"sample_path":                sample.Path,
		"prompt_strategy":            PromptStrategy(),
		"prompt_version_id":          promptVersionID,
		"prompt_mode":                promptMode,
		"prompt_path":                promptPath,
		"scenario":                   sample.Scenario,
		"generated_test_path":        testPath,
		"response_path":              respPath,
		"trace_path":                 trace.TracePath,
		"workspace_diff_path":        trace.WorkspaceDiffPath,
		"sandbox_fingerprint":        trace.SandboxFingerprint,
		"dataset_class":              sample.Category,
		"source_md5":                 sample.SourceMD5,
		"subject_version_id":         identity.SubjectVersionID,
		"generation_key":             identity.GenerationKey,
		"dependency_fingerprint":     identity.DependencyFingerprint,
		"generation_env_fingerprint": identity.GenerationEnvFingerprint,
		"latency_ms":                 latencyForMeta,
		"tokens": map[string]any{
			"prompt_tokens":      promptTokens,
			"completion_tokens":  completionTokens,
			"total_tokens":       totalTokens,
			"token_source":       trace.TokenSource,
			"estimated_cost_usd": trace.EstimatedCostUSD,
			"cost_source":        trace.CostSource,
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
		Model:                    model,
		SubjectID:                subject.ID,
		SubjectKind:              subject.Kind,
		AgentFramework:           subject.Framework,
		AgentModel:               subject.Model,
		SkillName:                subject.Skill,
		SkillVersion:             skillVersion,
		Language:                 sample.Language,
		SampleID:                 sample.ID,
		SampleUID:                identity.SampleUID,
		SamplePath:               sample.Path,
		PromptVersionID:          promptVersionID,
		PromptMode:               promptMode,
		PromptPath:               promptPath,
		GeneratedTestPath:        testPath,
		ResponsePath:             respPath,
		MetadataPath:             metadataPath,
		LatencyMS:                latency,
		PromptTokens:             promptTokens,
		CompletionTokens:         completionTokens,
		TotalTokens:              totalTokens,
		TokenSource:              trace.TokenSource,
		EstimatedCostUSD:         trace.EstimatedCostUSD,
		CostSource:               trace.CostSource,
		TracePath:                trace.TracePath,
		WorkspaceDiffPath:        trace.WorkspaceDiffPath,
		SandboxFingerprint:       trace.SandboxFingerprint,
		SubjectVersionID:         identity.SubjectVersionID,
		GenerationKey:            identity.GenerationKey,
		DependencyFingerprint:    identity.DependencyFingerprint,
		GenerationEnvFingerprint: identity.GenerationEnvFingerprint,
		GeneratedAtUTC:           time.Now().UTC(),
		Success:                  true,
		Truncated:                truncated,
		InteractionCount:         agentSummary.InteractionCount,
		ToolCallCount:            agentSummary.ToolCallCount,
		FilesReadCount:           agentSummary.FilesRead,
		FilesWriteCount:          agentSummary.FilesWritten,
		CommandCount:             agentSummary.CommandsExecuted,
	}
}

// languageExt 返回编程语言的文件扩展名
//
// 参数:
//   - language: 编程语言名称
//
// 返回值:
//   - string: 文件扩展名（如 ".py"、"java"）
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

// buildPlaceholderTest 构建占位符测试代码
// 用于 dry-run 模式，生成简单但不执行真实 API 调用的测试
//
// 参数:
//   - language: 编程语言
//   - sampleID: 样本 ID
//
// 返回值:
//   - string: 占位符测试代码
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

// sha256Bytes 计算数据的 SHA256 哈希值
//
// 参数:
//   - data: 输入数据
//
// 返回值:
//   - string: 十六进制哈希字符串
func sha256Bytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// copyFile 复制文件到目标路径
//
// 参数:
//   - src: 源文件路径
//   - dst: 目标文件路径
//
// 返回值:
//   - error: 复制失败时的错误
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

// sanitizeIdentifier 将字符串规范化为合法标识符
// 移除特殊字符，转换为小写
//
// 参数:
//   - raw: 原始字符串
//
// 返回值:
//   - string: 规范化后的标识符
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

// taskKey 构建任务的唯一标识键
// 格式为 "model|language|sampleID"
//
// 参数:
//   - model: 模型名称
//   - language: 编程语言
//   - sampleID: 样本 ID
//
// 返回值:
//   - string: 任务键
func taskKey(model, language, sampleID string) string {
	return model + "|" + language + "|" + sampleID
}

// buildCheckpointPath 构建 checkpoint 文件路径
// 根据运行规格参数生成唯一哈希，确保参数变化时 checkpoint 失效
//
// 参数:
//   - spec: 运行规格
//   - models: 模型配置列表
//
// 返回值:
//   - string: checkpoint 文件路径
func buildCheckpointPath(spec contracts.RunSpec, subjects []subjectTarget) string {
	subjectIDs := make([]string, 0, len(subjects))
	for _, item := range subjects {
		subjectIDs = append(subjectIDs, item.subject.Spec.ID)
	}
	sort.Strings(subjectIDs)

	langs := append([]string{}, spec.Languages...)
	for i := range langs {
		langs[i] = strings.ToLower(strings.TrimSpace(langs[i]))
	}
	sort.Strings(langs)

	scope := fmt.Sprintf(
		"subjects=%s;langs=%s;class=%s;level=%s;manifest=%s;max=%d;dataset=%s;agents=%s",
		strings.Join(subjectIDs, ","),
		strings.Join(langs, ","),
		strings.Join(spec.DatasetClasses, ","),
		spec.DatasetLevel,
		spec.DatasetManifest,
		spec.MaxSamples,
		spec.DatasetRoot,
		spec.AgentsConfigPath,
	)
	h := sha1.Sum([]byte(scope))
	hash := hex.EncodeToString(h[:])[:12]
	return filepath.Join(spec.OutputRoot, "checkpoints", "runner_"+hash+".checkpoint.json")
}

// loadCheckpoint 加载 checkpoint 文件
// 解析已完成的任务列表，用于增量运行
//
// 参数:
//   - path: checkpoint 文件路径
//
// 返回值:
//   - map[string]struct{}: 已完成任务键集合
//   - error: 加载错误（文件不存在时返回空集合）
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

// saveCheckpoint 保存 checkpoint 文件
// 将已完成任务列表写入 JSON 文件
//
// 参数:
//   - path: checkpoint 文件路径
//   - completed: 已完成任务键集合
//
// 返回值:
//   - error: 保存错误
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

// getModelNames 从模型配置列表提取模型名称
//
// 参数:
//   - configs: 模型配置列表
//
// 返回值:
//   - []string: 模型名称列表
func getModelNames(configs []modelConfig) []string {
	names := make([]string, 0, len(configs))
	for _, c := range configs {
		names = append(names, c.Name)
	}
	return names
}

// getLanguagesFromSamples 从样本列表统计语言分布
//
// 参数:
//   - samples: 样本列表
//
// 返回值:
//   - string: 语言分布统计字符串（如 "python:10, go:5"）
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

// trimErrorMsg 截断错误消息到指定长度
//
// 参数:
//   - msg: 原始消息
//   - max: 最大长度
//
// 返回值:
//   - string: 截断后的消息
func trimErrorMsg(msg string, max int) string {
	msg = strings.TrimSpace(msg)
	if len(msg) <= max {
		return msg
	}
	return msg[:max] + "..."
}

// errorMsgSafe 安全提取错误消息
// 从 ErrorInfo 结构中获取截断后的消息
//
// 参数:
//   - err: 错误信息结构
//
// 返回值:
//   - string: 截断后的错误消息（无错误时返回空字符串）
func errorMsgSafe(err *contracts.ErrorInfo) string {
	if err == nil {
		return ""
	}
	return trimErrorMsg(err.Message, 100)
}
