// reporter 包提供评测报告生成功能
// 负责汇总评测结果、生成 JSON/HTML 报告、构建多维度分析数据
// 支持按模型、语言、场景等维度进行聚合分析
package reporter

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/obs"

	"gopkg.in/yaml.v3"
)

// Service 报告生成服务
// 持有日志记录器实例和提示词元数据提供者
type Service struct {
	logger  *obs.Logger
	prompts contracts.PromptMetaProvider
}

// Output 报告生成输出
// 包含报告数据和文件路径
type Output struct {
	Report         contracts.ReportPayload // 报告数据结构
	ReportJSONPath string                  // JSON 报告文件路径
	ReportHTMLPath string                  // HTML 报告文件路径
}

// ModelDetail 模型详细信息
// 包含模型标识、具体型号、提供商信息
type ModelDetail struct {
	Name     string // 模型标识名（如 deepseek）
	ModelID  string // 具体型号（如 deepseek-chat）
	Provider string // 提供商（如 deepseek、dashscope）
}

func (s *Service) getPromptTemplate(language string) string {
	return s.prompts.PromptTemplatePreview(language)
}

// NewService 创建报告生成服务实例
// 参数:
//   - logger: 日志记录器
//   - prompts: 提示词元数据提供者
//
// 返回值:
//   - *Service: 报告服务实例
func NewService(logger *obs.Logger, prompts contracts.PromptMetaProvider) *Service {
	return &Service{logger: logger, prompts: prompts}
}

// Generate 从评测结果文件生成报告
// 读取 evaluation_result.json，构建多维度分析数据，生成 JSON 和 HTML 报告
//
// 参数:
//   - _ctx: 上下文（当前未使用）
//   - spec: 运行规格说明
//   - evaluationPath: 评测结果文件路径
//
// 返回值:
//   - Output: 报告输出（包含报告数据和文件路径）
//   - error: 生成过程中的错误
func (s *Service) Generate(_ context.Context, spec contracts.RunSpec, evaluationPath string) (Output, error) {
	set, err := contracts.ReadEvaluationResultSet(evaluationPath)
	if err != nil {
		return Output{}, err
	}
	return s.GenerateFromResultSet(spec, set, evaluationPath)
}

// GenerateFromResultSet 从评测结果集直接生成报告
// 不需要读取文件，直接使用传入的结果集数据
//
// 参数:
//   - spec: 运行规格说明
//   - set: 评测结果集
//   - sourceEvaluation: 源评测文件路径（用于记录）
//
// 返回值:
//   - Output: 报告输出
//   - error: 生成过程中的错误
func (s *Service) GenerateFromResultSet(spec contracts.RunSpec, set contracts.EvaluationResultSet, sourceEvaluation string) (Output, error) {
	promptStrategy, promptVersionID, promptSnapshotDir, prompts := s.loadPromptArtifacts(set.ManifestPath, sourceEvaluation)
	manifest := s.loadGeneratedManifest(set.ManifestPath, sourceEvaluation)

	reportRoot := filepath.Join(spec.OutputRoot, "runs", spec.RunID, "report")
	if err := os.MkdirAll(reportRoot, 0o755); err != nil {
		return Output{}, err
	}

	modelDetails := loadModelDetails(spec.ConfigPath)
	summary := buildSummary(set.Results)
	dims := buildDimensions(set.Results, modelDetails)
	tokenStats := buildTokenStats(set.Results)
	topModels := buildTopModels(dims.ByModel)
	failures := buildFailureRows(set.Results)
	scoreExclusions := buildScoreExclusions(set.Results)
	zeroMutantSamples := buildZeroMutantSamples(set.Results)
	breakdown := buildMutationBreakdown(set.Results)
	modelInfos := buildModelInfos(dims.ByModel, modelDetails)
	truncationStats := buildTruncationStats(set.Results)
	// 新增：洞察、效率、错误诊断
	insights := buildInsights(topModels, dims, summary, failures)
	efficiencyStats := buildEfficiencyStats(topModels, set.Results)
	errorDiagnosis := buildErrorDiagnosis(set.Results)
	agentComparisons := buildAgentComparisons(set.Results)
	skillUplifts := buildSkillUplifts(set.Results)
	comparisonViews := buildComparisonViews(set.Results)

	payload := contracts.ReportPayload{
		SchemaVersion:     contracts.SchemaVersion,
		RunID:             spec.RunID,
		GeneratedAtUTC:    time.Now().UTC(),
		SourceEvaluation:  sourceEvaluation,
		PromptStrategy:    promptStrategy,
		PromptVersionID:   promptVersionID,
		PromptSnapshotDir: promptSnapshotDir,
		Summary:           summary,
		Dimensions:        dims,
		TopModels:         topModels,
		ModelInfos:        modelInfos,
		ByScenario:        dims.ByScenario,
		ByModelScenario:   dims.ByModelScenario,
		TokenStats:        tokenStats,
		Failures:          failures,
		ScoreExclusions:   scoreExclusions,
		ZeroMutantSamples: zeroMutantSamples,
		Thresholds: contracts.Thresholds{
			CompilePassRate: 1.0,
			TestPassRate:    0.7,
			LineCoverage:    0.7,
			BranchCoverage:  0.6,
			MutationScore:   0.85,
		},
		Prompts:         prompts,
		TruncationStats: truncationStats,
		// 新增字段
		Insights:         insights,
		EfficiencyStats:  efficiencyStats,
		ErrorDiagnosis:   errorDiagnosis,
		RuntimeSummary:   buildRuntimeSummary(manifest, set),
		AgentComparisons: agentComparisons,
		SkillUplifts:     skillUplifts,
		ComparisonViews:  comparisonViews,
	}

	jsonPath := filepath.Join(reportRoot, "report_summary.json")
	htmlPath := filepath.Join(reportRoot, "report.html")
	summaryJSON := map[string]any{
		"schema_version":      payload.SchemaVersion,
		"run_id":              payload.RunID,
		"generated_at_utc":    payload.GeneratedAtUTC,
		"source_evaluation":   payload.SourceEvaluation,
		"prompt_strategy":     payload.PromptStrategy,
		"prompt_version_id":   payload.PromptVersionID,
		"prompt_snapshot_dir": payload.PromptSnapshotDir,
		"summary":             payload.Summary,
		"dimensions":          payload.Dimensions,
		"top_models":          payload.TopModels,
		"model_infos":         payload.ModelInfos,
		"by_scenario":         payload.ByScenario,
		"by_model_scenario":   payload.ByModelScenario,
		"token_stats":         payload.TokenStats,
		"failures":            payload.Failures,
		"score_exclusions":    payload.ScoreExclusions,
		"zero_mutant_samples": payload.ZeroMutantSamples,
		"mutation_breakdown":  breakdown,
		"thresholds":          payload.Thresholds,
		"prompts":             payload.Prompts,
		"truncation_stats":    payload.TruncationStats,
		"insights":            payload.Insights,
		"efficiency_stats":    payload.EfficiencyStats,
		"error_diagnosis":     payload.ErrorDiagnosis,
		"runtime_summary":     payload.RuntimeSummary,
		"agent_comparisons":   payload.AgentComparisons,
		"skill_uplifts":       payload.SkillUplifts,
		"comparison_views":    payload.ComparisonViews,
	}
	if err := contracts.WriteJSON(jsonPath, summaryJSON); err != nil {
		return Output{}, err
	}
	if err := os.WriteFile(htmlPath, []byte(buildHTML(payload, breakdown, set.Results, s.prompts)), 0o644); err != nil {
		return Output{}, err
	}

	s.logger.Info("report generated", "run_id", spec.RunID, "summary", jsonPath, "html", htmlPath)
	return Output{Report: payload, ReportJSONPath: jsonPath, ReportHTMLPath: htmlPath}, nil
}

func (s *Service) loadGeneratedManifest(manifestPath, sourceEvaluation string) *contracts.GeneratedManifest {
	resolvedManifest := resolveArtifactPath(manifestPath, sourceEvaluation)
	if strings.TrimSpace(resolvedManifest) == "" {
		return nil
	}
	manifest, err := contracts.ReadGeneratedManifest(resolvedManifest)
	if err != nil {
		return nil
	}
	return &manifest
}

func (s *Service) loadPromptArtifacts(manifestPath, sourceEvaluation string) (string, string, string, map[string]string) {
	promptTemplates := make(map[string]string, len(contracts.SupportedLanguages))
	for _, lang := range contracts.SupportedLanguages {
		promptTemplates[lang] = s.getPromptTemplate(lang)
	}
	if strings.TrimSpace(manifestPath) == "" {
		return s.prompts.PromptStrategy(), s.prompts.PromptVersionID(), "", promptTemplates
	}

	resolvedManifest := resolveArtifactPath(manifestPath, sourceEvaluation)
	manifest, err := contracts.ReadGeneratedManifest(resolvedManifest)
	if err != nil {
		return s.prompts.PromptStrategy(), s.prompts.PromptVersionID(), "", promptTemplates
	}
	if strings.TrimSpace(manifest.PromptSnapshotDir) != "" {
		if catalog, err := contracts.LoadPromptCatalog(resolveArtifactPath(manifest.PromptSnapshotDir, sourceEvaluation)); err == nil {
			for language, modeTemplates := range catalog.Templates {
				if prompt := strings.TrimSpace(modeTemplates[contracts.PromptModeFullFile]); prompt != "" {
					promptTemplates[language] = prompt
				}
			}
		}
	}

	seenActual := map[string]struct{}{}
	for _, item := range manifest.Cases {
		lang := strings.ToLower(strings.TrimSpace(item.Language))
		if lang == "" {
			continue
		}
		if _, ok := seenActual[lang]; ok {
			continue
		}
		if strings.TrimSpace(item.PromptPath) == "" {
			continue
		}
		raw, err := os.ReadFile(resolveArtifactPath(item.PromptPath, sourceEvaluation))
		if err != nil || len(raw) == 0 {
			continue
		}
		promptTemplates[lang] = string(raw)
		seenActual[lang] = struct{}{}
	}

	return manifest.PromptStrategy, manifest.PromptVersionID, manifest.PromptSnapshotDir, promptTemplates
}

func resolveArtifactPath(pathValue, anchorPath string) string {
	pathValue = strings.TrimSpace(pathValue)
	if pathValue == "" {
		return pathValue
	}
	if _, err := os.Stat(pathValue); err == nil {
		return pathValue
	}

	normalized := filepath.ToSlash(pathValue)
	if strings.HasPrefix(normalized, "/app/artifacts/") {
		if artifactsRoot := findArtifactsRoot(anchorPath); artifactsRoot != "" {
			candidate := filepath.Join(artifactsRoot, filepath.FromSlash(strings.TrimPrefix(normalized, "/app/artifacts/")))
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
			return candidate
		}
	}
	if strings.HasPrefix(normalized, "artifacts/") {
		candidate := filepath.FromSlash(normalized)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return pathValue
}

func findArtifactsRoot(anchorPath string) string {
	if strings.TrimSpace(anchorPath) == "" {
		return ""
	}
	current := anchorPath
	if info, err := os.Stat(current); err == nil && !info.IsDir() {
		current = filepath.Dir(current)
	}
	for {
		if strings.EqualFold(filepath.Base(current), "artifacts") {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}
		current = parent
	}
}

func loadModelDetails(configPath string) map[string]ModelDetail {
	details := map[string]ModelDetail{}
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return details
	}
	var cfg struct {
		Models map[string]struct {
			Enabled  bool   `yaml:"enabled"`
			Provider string `yaml:"provider"`
			Config   struct {
				Model string `yaml:"model"`
			} `yaml:"config"`
		} `yaml:"models"`
	}
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return details
	}
	for name, item := range cfg.Models {
		details[name] = ModelDetail{
			Name:     name,
			ModelID:  item.Config.Model,
			Provider: item.Provider,
		}
	}
	return details
}

func buildModelInfos(models []contracts.ModelDim, details map[string]ModelDetail) []contracts.ModelInfo {
	var infos []contracts.ModelInfo
	for _, m := range models {
		detail, ok := details[m.Model]
		modelID := m.Model
		provider := ""
		if ok {
			modelID = detail.ModelID
			provider = detail.Provider
		}
		infos = append(infos, contracts.ModelInfo{
			Name:       m.Model,
			ModelID:    modelID,
			Provider:   provider,
			TotalCases: m.TotalSamples,
		})
	}
	return infos
}

func buildSummary(rows []contracts.EvaluationResult) contracts.ReportSummary {
	total := len(rows)
	eligibleTotal := 0
	excludedTotal := 0
	compilePass := 0
	sampleTestPass := 0
	testCasePassTotal := 0
	testCaseTotal := 0
	fallbackCasePass := 0
	fallbackCaseTotal := 0
	lineSum := 0.0
	lineCnt := 0
	mutationSum := 0.0
	mutationCnt := 0
	assertDensitySum := 0.0
	assertDensityCnt := 0

	for _, row := range rows {
		if !isScoreEligible(row) {
			excludedTotal++
			continue
		}
		eligibleTotal++
		if row.CompilePass {
			compilePass++
		}
		if row.TestPass != nil && *row.TestPass {
			sampleTestPass++
		}
		if row.TestPassCount != nil && row.TestTotalCount != nil {
			testCasePassTotal += *row.TestPassCount
			testCaseTotal += *row.TestTotalCount
		} else if row.TestPass != nil {
			fallbackCaseTotal++
			if *row.TestPass {
				fallbackCasePass++
			}
		}
		if row.LineCoverage != nil {
			lineSum += *row.LineCoverage
			lineCnt++
		}
		if row.MutationScore != nil {
			mutationSum += *row.MutationScore
			mutationCnt++
		}
		if row.AssertionDensity != nil {
			assertDensitySum += *row.AssertionDensity
			assertDensityCnt++
		}
	}

	if fallbackCaseTotal > 0 {
		testCasePassTotal += fallbackCasePass
		testCaseTotal += fallbackCaseTotal
	}

	return contracts.ReportSummary{
		TotalSamples:        total,
		EligibleSamples:     eligibleTotal,
		ExcludedSamples:     excludedTotal,
		CompilePassCount:    compilePass,
		CompilePassRate:     rate(compilePass, eligibleTotal),
		TestPassCount:       sampleTestPass,
		TestPassRate:        rate(sampleTestPass, eligibleTotal),
		SampleTestPassCount: sampleTestPass,
		SampleTestPassRate:  rate(sampleTestPass, eligibleTotal),
		TestCasePassCount:   testCasePassTotal,
		TestCasePassRate:    rate(testCasePassTotal, testCaseTotal),
		AvgLineCoverage:     avg(lineSum, lineCnt),
		AvgMutationScore:    avg(mutationSum, mutationCnt),
		AvgAssertionDensity: avg(assertDensitySum, assertDensityCnt),
	}
}
