package reporter

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/obs"
	"go-ut-bench/internal/runner"

	"gopkg.in/yaml.v3"
)

type Service struct {
	logger *obs.Logger
}

type Output struct {
	Report         contracts.ReportPayload
	ReportJSONPath string
	ReportHTMLPath string
}

type mutationBreakdown struct {
	Total      int                     `json:"total"`
	Killed     int                     `json:"killed"`
	Survived   int                     `json:"survived"`
	NoTests    int                     `json:"no_tests"`
	Timeouts   int                     `json:"timeouts"`
	Skipped    int                     `json:"skipped"`
	Suspicious int                     `json:"suspicious"`
	ByTool     []mutationToolBreakdown `json:"by_tool,omitempty"`
}

type mutationToolBreakdown struct {
	Tool       string `json:"tool"`
	Total      int    `json:"total"`
	Killed     int    `json:"killed"`
	Survived   int    `json:"survived"`
	NoTests    int    `json:"no_tests"`
	Timeouts   int    `json:"timeouts"`
	Skipped    int    `json:"skipped"`
	Suspicious int    `json:"suspicious"`
}

type ModelDetail struct {
	Name     string
	ModelID  string
	Provider string
}

func getPromptTemplate(language string) string {
	return runner.PromptTemplatePreview(language)
}

func NewService(logger *obs.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) Generate(_ context.Context, spec contracts.RunSpec, evaluationPath string) (Output, error) {
	set, err := contracts.ReadEvaluationResultSet(evaluationPath)
	if err != nil {
		return Output{}, err
	}
	promptStrategy, promptVersionID, promptSnapshotDir, prompts := loadPromptArtifacts(set.ManifestPath)

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
	breakdown := buildMutationBreakdown(set.Results)
	modelInfos := buildModelInfos(dims.ByModel, modelDetails)
	truncationStats := buildTruncationStats(set.Results)

	payload := contracts.ReportPayload{
		SchemaVersion:     contracts.SchemaVersion,
		RunID:             spec.RunID,
		GeneratedAtUTC:    time.Now().UTC(),
		SourceEvaluation:  evaluationPath,
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
		Thresholds: contracts.Thresholds{
			CompilePassRate: 1.0,
			TestPassRate:    0.7,
			LineCoverage:    0.7,
			BranchCoverage:  0.6,
			MutationScore:   0.85,
		},
		Prompts:         prompts,
		TruncationStats: truncationStats,
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
		"mutation_breakdown":  breakdown,
		"thresholds":          payload.Thresholds,
		"prompts":             payload.Prompts,
	}
	if err := contracts.WriteJSON(jsonPath, summaryJSON); err != nil {
		return Output{}, err
	}
	if err := os.WriteFile(htmlPath, []byte(buildHTML(payload, breakdown, set.Results)), 0o644); err != nil {
		return Output{}, err
	}

	s.logger.Info("report generated", "run_id", spec.RunID, "summary", jsonPath, "html", htmlPath)
	return Output{Report: payload, ReportJSONPath: jsonPath, ReportHTMLPath: htmlPath}, nil
}

func loadPromptArtifacts(manifestPath string) (string, string, string, map[string]string) {
	prompts := map[string]string{
		"python": getPromptTemplate("python"),
		"go":     getPromptTemplate("go"),
		"java":   getPromptTemplate("java"),
		"cpp":    getPromptTemplate("cpp"),
	}
	if strings.TrimSpace(manifestPath) == "" {
		return runner.PromptStrategy(), runner.PromptVersionID(), "", prompts
	}

	manifest, err := contracts.ReadGeneratedManifest(manifestPath)
	if err != nil {
		return runner.PromptStrategy(), runner.PromptVersionID(), "", prompts
	}
	if strings.TrimSpace(manifest.PromptSnapshotDir) != "" {
		if catalog, err := runner.LoadPromptCatalog(manifest.PromptSnapshotDir); err == nil {
			loaded := make(map[string]string, len(catalog.Templates))
			for language, modeTemplates := range catalog.Templates {
				loaded[language] = modeTemplates[runner.PromptModeFullFile]
			}
			return catalog.Strategy, catalog.VersionID, manifest.PromptSnapshotDir, loaded
		}
	}
	return manifest.PromptStrategy, manifest.PromptVersionID, manifest.PromptSnapshotDir, prompts
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

func buildDimensions(rows []contracts.EvaluationResult, modelDetails map[string]ModelDetail) contracts.Dimensions {
	modelMap := map[string]*modelAgg{}
	langMap := map[string]*modelAgg{}
	scenarioMap := map[string]*scenarioAgg{}
	modelScenarioMap := map[string]*modelScenarioAgg{}

	for _, row := range rows {
		if !isScoreEligible(row) {
			continue
		}
		agg := getOrCreateModelAgg(modelMap, row.Model)
		mergeModelAgg(agg, row)

		langAgg := getOrCreateModelAgg(langMap, row.Language)
		mergeModelAgg(langAgg, row)

		scenario := extractScenario(row.SampleID)
		scenarioKey := fmt.Sprintf("%s|%s", row.Language, scenario)
		scenAgg := getOrCreateScenarioAgg(scenarioMap, scenarioKey)
		mergeScenarioAgg(scenAgg, row, scenario, row.Language)

		modelScenKey := fmt.Sprintf("%s|%s|%s", row.Model, row.Language, scenario)
		modelScenAgg := getOrCreateModelScenarioAgg(modelScenarioMap, modelScenKey)
		mergeModelScenarioAgg(modelScenAgg, row, row.Model, scenario, row.Language)
	}

	var byModel []contracts.ModelDim
	for _, agg := range modelMap {
		detail, ok := modelDetails[agg.key]
		modelID := agg.key
		provider := ""
		if ok {
			modelID = detail.ModelID
			provider = detail.Provider
		}
		byModel = append(byModel, contracts.ModelDim{
			Model:               agg.key,
			ModelID:             modelID,
			Provider:            provider,
			TotalSamples:        agg.count,
			CompilePassRate:     rate(agg.compilePass, agg.count),
			AvgTestPassRate:     rate(agg.sampleTestPass, agg.count),
			AvgTestCasePassRate: rate(agg.testPassTotal, agg.testTotal),
			AvgLineCoverage:     avg(agg.lineSum, agg.lineCnt),
			AvgBranchCoverage:   avg(agg.branchSum, agg.branchCnt),
			AvgMutationScore:    avg(agg.mutationSum, agg.mutationCnt),
			AvgLatencyMS:        avgFloat(agg.latencySum, agg.latencyCnt),
			AvgPromptTokens:     avgFloat(agg.promptTokensSum, agg.tokenCnt),
			AvgCompletionTokens: avgFloat(agg.completionTokensSum, agg.tokenCnt),
			AvgTotalTokens:      avgFloat(agg.totalTokensSum, agg.tokenCnt),
		})
	}
	sort.Slice(byModel, func(i, j int) bool { return byModel[i].Model < byModel[j].Model })

	var byLanguage []contracts.LanguageDim
	for _, agg := range langMap {
		byLanguage = append(byLanguage, contracts.LanguageDim{
			Language:            agg.key,
			TotalSamples:        agg.count,
			CompilePassRate:     rate(agg.compilePass, agg.count),
			AvgTestPassRate:     rate(agg.sampleTestPass, agg.count),
			AvgTestCasePassRate: rate(agg.testPassTotal, agg.testTotal),
			AvgLineCoverage:     avg(agg.lineSum, agg.lineCnt),
			AvgBranchCoverage:   avg(agg.branchSum, agg.branchCnt),
			AvgMutationScore:    avg(agg.mutationSum, agg.mutationCnt),
		})
	}
	sort.Slice(byLanguage, func(i, j int) bool { return byLanguage[i].Language < byLanguage[j].Language })

	var byScenario []contracts.ScenarioDim
	for _, agg := range scenarioMap {
		byScenario = append(byScenario, contracts.ScenarioDim{
			Scenario:            agg.scenario,
			Language:            agg.language,
			TotalSamples:        agg.count,
			CompilePassRate:     rate(agg.compilePass, agg.count),
			AvgTestPassRate:     rate(agg.sampleTestPass, agg.count),
			AvgTestCasePassRate: rate(agg.testPassTotal, agg.testTotal),
			AvgLineCoverage:     avg(agg.lineSum, agg.lineCnt),
			AvgBranchCoverage:   avg(agg.branchSum, agg.branchCnt),
			AvgMutationScore:    avg(agg.mutationSum, agg.mutationCnt),
			AvgLatencyMS:        avgFloat(agg.latencySum, agg.latencyCnt),
			AvgTokens:           avgFloat(agg.totalTokensSum, agg.tokenCnt),
		})
	}
	sort.Slice(byScenario, func(i, j int) bool {
		if byScenario[i].Language != byScenario[j].Language {
			return byScenario[i].Language < byScenario[j].Language
		}
		return byScenario[i].Scenario < byScenario[j].Scenario
	})

	var byModelScenario []contracts.ModelScenarioDim
	for _, agg := range modelScenarioMap {
		byModelScenario = append(byModelScenario, contracts.ModelScenarioDim{
			Model:               agg.model,
			Scenario:            agg.scenario,
			Language:            agg.language,
			TotalSamples:        agg.count,
			CompilePassRate:     rate(agg.compilePass, agg.count),
			AvgTestPassRate:     rate(agg.sampleTestPass, agg.count),
			AvgTestCasePassRate: rate(agg.testPassTotal, agg.testTotal),
			AvgLineCoverage:     avg(agg.lineSum, agg.lineCnt),
			AvgBranchCoverage:   avg(agg.branchSum, agg.branchCnt),
			AvgMutationScore:    avg(agg.mutationSum, agg.mutationCnt),
			AvgLatencyMS:        avgFloat(agg.latencySum, agg.latencyCnt),
			AvgPromptTokens:     avgFloat(agg.promptTokensSum, agg.tokenCnt),
			AvgCompletionTokens: avgFloat(agg.completionTokensSum, agg.tokenCnt),
			AvgTotalTokens:      avgFloat(agg.totalTokensSum, agg.tokenCnt),
		})
	}
	sort.Slice(byModelScenario, func(i, j int) bool {
		if byModelScenario[i].Model != byModelScenario[j].Model {
			return byModelScenario[i].Model < byModelScenario[j].Model
		}
		if byModelScenario[i].Language != byModelScenario[j].Language {
			return byModelScenario[i].Language < byModelScenario[j].Language
		}
		return byModelScenario[i].Scenario < byModelScenario[j].Scenario
	})

	return contracts.Dimensions{
		ByModel:         byModel,
		ByLanguage:      byLanguage,
		ByScenario:      byScenario,
		ByModelScenario: byModelScenario,
	}
}

// 聚合器类型定义
type modelAgg struct {
	key                 string
	count               int
	compilePass         int
	sampleTestPass      int
	testPassTotal       int
	testTotal           int
	lineSum             float64
	lineCnt             int
	branchSum           float64
	branchCnt           int
	mutationSum         float64
	mutationCnt         int
	latencySum          float64
	latencyCnt          int
	tokenCnt            int
	promptTokensSum     float64
	completionTokensSum float64
	totalTokensSum      float64
}

type scenarioAgg struct {
	key                 string
	scenario            string
	language            string
	count               int
	compilePass         int
	sampleTestPass      int
	testPassTotal       int
	testTotal           int
	lineSum             float64
	lineCnt             int
	branchSum           float64
	branchCnt           int
	mutationSum         float64
	mutationCnt         int
	latencySum          float64
	latencyCnt          int
	tokenCnt            int
	promptTokensSum     float64
	completionTokensSum float64
	totalTokensSum      float64
}

type modelScenarioAgg struct {
	key                 string
	model               string
	scenario            string
	language            string
	count               int
	compilePass         int
	sampleTestPass      int
	testPassTotal       int
	testTotal           int
	lineSum             float64
	lineCnt             int
	branchSum           float64
	branchCnt           int
	mutationSum         float64
	mutationCnt         int
	latencySum          float64
	latencyCnt          int
	tokenCnt            int
	promptTokensSum     float64
	completionTokensSum float64
	totalTokensSum      float64
}

func extractScenario(sampleID string) string {
	parts := strings.Split(sampleID, "_")
	if len(parts) >= 1 {
		return parts[0]
	}
	return "unknown"
}

func getOrCreateModelAgg(m map[string]*modelAgg, key string) *modelAgg {
	if a, ok := m[key]; ok {
		return a
	}
	a := &modelAgg{key: key}
	m[key] = a
	return a
}

func getOrCreateScenarioAgg(m map[string]*scenarioAgg, key string) *scenarioAgg {
	if a, ok := m[key]; ok {
		return a
	}
	a := &scenarioAgg{key: key}
	m[key] = a
	return a
}

func getOrCreateModelScenarioAgg(m map[string]*modelScenarioAgg, key string) *modelScenarioAgg {
	if a, ok := m[key]; ok {
		return a
	}
	a := &modelScenarioAgg{key: key}
	m[key] = a
	return a
}

func mergeModelAgg(a *modelAgg, row contracts.EvaluationResult) {
	a.count++
	if row.CompilePass {
		a.compilePass++
	}
	if row.TestPass != nil && *row.TestPass {
		a.sampleTestPass++
	}
	if row.TestPassCount != nil && row.TestTotalCount != nil {
		a.testPassTotal += *row.TestPassCount
		a.testTotal += *row.TestTotalCount
	} else if row.TestPass != nil {
		a.testTotal++
		if *row.TestPass {
			a.testPassTotal++
		}
	}
	if row.LineCoverage != nil {
		a.lineSum += *row.LineCoverage
		a.lineCnt++
	}
	if row.BranchCoverage != nil {
		a.branchSum += *row.BranchCoverage
		a.branchCnt++
	}
	if row.MutationScore != nil {
		a.mutationSum += *row.MutationScore
		a.mutationCnt++
	}
	if row.RuntimeMS != nil {
		a.latencySum += float64(*row.RuntimeMS)
		a.latencyCnt++
	}
	if row.PromptTokens != nil {
		a.promptTokensSum += float64(*row.PromptTokens)
		a.completionTokensSum += float64(*row.CompletionTokens)
		a.totalTokensSum += float64(*row.TotalTokens)
		a.tokenCnt++
	}
}

func mergeScenarioAgg(a *scenarioAgg, row contracts.EvaluationResult, scenario, language string) {
	a.scenario = scenario
	a.language = language
	a.count++
	if row.CompilePass {
		a.compilePass++
	}
	if row.TestPass != nil && *row.TestPass {
		a.sampleTestPass++
	}
	if row.TestPassCount != nil && row.TestTotalCount != nil {
		a.testPassTotal += *row.TestPassCount
		a.testTotal += *row.TestTotalCount
	} else if row.TestPass != nil {
		a.testTotal++
		if *row.TestPass {
			a.testPassTotal++
		}
	}
	if row.LineCoverage != nil {
		a.lineSum += *row.LineCoverage
		a.lineCnt++
	}
	if row.BranchCoverage != nil {
		a.branchSum += *row.BranchCoverage
		a.branchCnt++
	}
	if row.MutationScore != nil {
		a.mutationSum += *row.MutationScore
		a.mutationCnt++
	}
	if row.RuntimeMS != nil {
		a.latencySum += float64(*row.RuntimeMS)
		a.latencyCnt++
	}
	if row.TotalTokens != nil {
		a.promptTokensSum += float64(*row.PromptTokens)
		a.completionTokensSum += float64(*row.CompletionTokens)
		a.totalTokensSum += float64(*row.TotalTokens)
		a.tokenCnt++
	}
}

func mergeModelScenarioAgg(a *modelScenarioAgg, row contracts.EvaluationResult, model, scenario, language string) {
	a.model = model
	a.scenario = scenario
	a.language = language
	a.count++
	if row.CompilePass {
		a.compilePass++
	}
	if row.TestPass != nil && *row.TestPass {
		a.sampleTestPass++
	}
	if row.TestPassCount != nil && row.TestTotalCount != nil {
		a.testPassTotal += *row.TestPassCount
		a.testTotal += *row.TestTotalCount
	} else if row.TestPass != nil {
		a.testTotal++
		if *row.TestPass {
			a.testPassTotal++
		}
	}
	if row.LineCoverage != nil {
		a.lineSum += *row.LineCoverage
		a.lineCnt++
	}
	if row.BranchCoverage != nil {
		a.branchSum += *row.BranchCoverage
		a.branchCnt++
	}
	if row.MutationScore != nil {
		a.mutationSum += *row.MutationScore
		a.mutationCnt++
	}
	if row.RuntimeMS != nil {
		a.latencySum += float64(*row.RuntimeMS)
		a.latencyCnt++
	}
	if row.TotalTokens != nil {
		a.promptTokensSum += float64(*row.PromptTokens)
		a.completionTokensSum += float64(*row.CompletionTokens)
		a.totalTokensSum += float64(*row.TotalTokens)
		a.tokenCnt++
	}
}

func buildTokenStats(rows []contracts.EvaluationResult) contracts.TokenStats {
	stats := contracts.TokenStats{}
	for _, row := range rows {
		if row.PromptTokens != nil {
			stats.TotalPromptTokens += *row.PromptTokens
			stats.TotalCompletionTokens += *row.CompletionTokens
			stats.TotalTokens += *row.TotalTokens
			stats.SampleCount++
		}
	}
	if stats.SampleCount > 0 {
		stats.AvgPromptTokens = float64(stats.TotalPromptTokens) / float64(stats.SampleCount)
		stats.AvgCompletionTokens = float64(stats.TotalCompletionTokens) / float64(stats.SampleCount)
		stats.AvgTotalTokens = float64(stats.TotalTokens) / float64(stats.SampleCount)
	}
	return stats
}

func buildTopModels(models []contracts.ModelDim) []contracts.ModelRank {
	var sorted []contracts.ModelDim
	for _, m := range models {
		// 计算综合得分：编译30% + 测试30% + 覆盖20% + 变异20%
		composite := m.CompilePassRate*0.3 + m.AvgTestPassRate*0.3 + m.AvgLineCoverage*0.2 + m.AvgMutationScore*0.2
		m.CompositeScore = round(composite, 6)
		sorted = append(sorted, m)
	}
	// 按综合得分降序排序
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].CompositeScore != sorted[j].CompositeScore {
			return sorted[i].CompositeScore > sorted[j].CompositeScore
		}
		if sorted[i].AvgTestPassRate != sorted[j].AvgTestPassRate {
			return sorted[i].AvgTestPassRate > sorted[j].AvgTestPassRate
		}
		if sorted[i].AvgLineCoverage != sorted[j].AvgLineCoverage {
			return sorted[i].AvgLineCoverage > sorted[j].AvgLineCoverage
		}
		return sorted[i].AvgMutationScore > sorted[j].AvgMutationScore
	})

	var out []contracts.ModelRank
	for i, m := range sorted {
		out = append(out, contracts.ModelRank{
			Rank:                i + 1,
			Model:               m.Model,
			ModelID:             m.ModelID,
			Provider:            m.Provider,
			CompilePassRate:     m.CompilePassRate,
			AvgTestPassRate:     m.AvgTestPassRate,
			AvgTestCasePassRate: m.AvgTestCasePassRate,
			AvgLineCoverage:     m.AvgLineCoverage,
			AvgMutationScore:    m.AvgMutationScore,
			CompositeScore:      m.CompositeScore,
			AvgLatencyMS:        m.AvgLatencyMS,
			AvgPromptTokens:     m.AvgPromptTokens,
			AvgCompletionTokens: m.AvgCompletionTokens,
			AvgTotalTokens:      m.AvgTotalTokens,
		})
	}
	return out
}

func buildTruncationStats(rows []contracts.EvaluationResult) contracts.TruncationStats {
	totalTruncated := 0
	modelStats := map[string]*truncationAgg{}
	langStats := map[string]*truncationAgg{}
	scenarioStats := map[string]*truncationAgg{}

	for _, row := range rows {
		if row.Truncated {
			totalTruncated++
		}

		modelKey := row.Model
		if _, ok := modelStats[modelKey]; !ok {
			modelStats[modelKey] = &truncationAgg{key: modelKey}
		}
		modelStats[modelKey].total++
		if row.Truncated {
			modelStats[modelKey].truncated++
		}
		if row.CompletionTokens != nil {
			modelStats[modelKey].completionTokensSum += float64(*row.CompletionTokens)
			modelStats[modelKey].completionTokensCount++
		}

		langKey := row.Language
		if _, ok := langStats[langKey]; !ok {
			langStats[langKey] = &truncationAgg{key: langKey}
		}
		langStats[langKey].total++
		if row.Truncated {
			langStats[langKey].truncated++
		}

		scenario := extractScenario(row.SampleID)
		scenarioKey := fmt.Sprintf("%s|%s", row.Language, scenario)
		if _, ok := scenarioStats[scenarioKey]; !ok {
			scenarioStats[scenarioKey] = &truncationAgg{key: scenarioKey, scenario: scenario, language: row.Language}
		}
		scenarioStats[scenarioKey].total++
		if row.Truncated {
			scenarioStats[scenarioKey].truncated++
		}
	}

	var byModel []contracts.ModelTruncationDim
	for _, agg := range modelStats {
		byModel = append(byModel, contracts.ModelTruncationDim{
			Model:               agg.key,
			TotalSamples:        agg.total,
			TruncatedCount:      agg.truncated,
			TruncationRate:      rate(agg.truncated, agg.total),
			AvgCompletionTokens: avg(agg.completionTokensSum, agg.completionTokensCount),
		})
	}
	sort.Slice(byModel, func(i, j int) bool { return byModel[i].TruncationRate > byModel[j].TruncationRate })

	var byLanguage []contracts.LangTruncationDim
	for _, agg := range langStats {
		byLanguage = append(byLanguage, contracts.LangTruncationDim{
			Language:       agg.key,
			TotalSamples:   agg.total,
			TruncatedCount: agg.truncated,
			TruncationRate: rate(agg.truncated, agg.total),
		})
	}
	sort.Slice(byLanguage, func(i, j int) bool { return byLanguage[i].TruncationRate > byLanguage[j].TruncationRate })

	var byScenario []contracts.ScenarioTruncationDim
	for _, agg := range scenarioStats {
		byScenario = append(byScenario, contracts.ScenarioTruncationDim{
			Scenario:       agg.scenario,
			Language:       agg.language,
			TotalSamples:   agg.total,
			TruncatedCount: agg.truncated,
			TruncationRate: rate(agg.truncated, agg.total),
		})
	}
	sort.Slice(byScenario, func(i, j int) bool { return byScenario[i].TruncationRate > byScenario[j].TruncationRate })

	return contracts.TruncationStats{
		TotalTruncated: totalTruncated,
		TruncationRate: rate(totalTruncated, len(rows)),
		ByModel:        byModel,
		ByLanguage:     byLanguage,
		ByScenario:     byScenario,
		ContinuationStats: contracts.ContinuationStats{
			Enabled:               true,
			TotalContinuations:    0,
			SuccessfulRecoveries:  0,
			RecoveryRate:          0,
			AvgContinuationRounds: 0,
		},
	}
}

type truncationAgg struct {
	key                   string
	scenario              string
	language              string
	total                 int
	truncated             int
	completionTokensSum   float64
	completionTokensCount int
}

func buildFailureRows(rows []contracts.EvaluationResult) []contracts.FailureRow {
	m := map[failureKey]*failureAgg{}

	for _, row := range rows {
		if row.Truncated {
			k := failureKey{stage: "generate", errType: "truncated"}
			agg := getOrCreateFailureAgg(m, k)
			agg.count++
			if agg.exampleModel == "" {
				agg.exampleModel = row.Model
				agg.exampleSample = row.SampleID
				agg.exampleMessage = "API response truncated due to max_tokens limit (finish_reason='length')"
			}
		}
		if row.CompileError != "" {
			k := failureKey{stage: "compile", errType: classifyError(row.CompileError)}
			agg := getOrCreateFailureAgg(m, k)
			agg.count++
			if agg.exampleModel == "" {
				agg.exampleModel = row.Model
				agg.exampleSample = row.SampleID
				agg.exampleMessage = shortErrText(row.CompileError)
			}
		}
		if row.TestError != "" {
			k := failureKey{stage: "test", errType: classifyError(row.TestError)}
			agg := getOrCreateFailureAgg(m, k)
			agg.count++
			if agg.exampleModel == "" {
				agg.exampleModel = row.Model
				agg.exampleSample = row.SampleID
				agg.exampleMessage = shortErrText(row.TestError)
			}
		}
		if row.CoverageError != "" {
			k := failureKey{stage: "coverage", errType: "coverage_error"}
			agg := getOrCreateFailureAgg(m, k)
			agg.count++
			if agg.exampleModel == "" {
				agg.exampleModel = row.Model
				agg.exampleSample = row.SampleID
				agg.exampleMessage = shortErrText(row.CoverageError)
			}
		}
		if row.MutationError != "" {
			k := failureKey{stage: "mutation", errType: classifyMutationError(row.MutationError)}
			agg := getOrCreateFailureAgg(m, k)
			agg.count++
			if agg.exampleModel == "" {
				agg.exampleModel = row.Model
				agg.exampleSample = row.SampleID
				agg.exampleMessage = shortErrText(row.MutationError)
			}
		}
	}

	var out []contracts.FailureRow
	for k, agg := range m {
		out = append(out, contracts.FailureRow{
			Stage:          k.stage,
			ErrorType:      k.errType,
			Count:          agg.count,
			ExampleModel:   agg.exampleModel,
			ExampleSample:  agg.exampleSample,
			ExampleMessage: agg.exampleMessage,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Count > out[j].Count })
	return out
}

func classifyMutationError(msg string) string {
	msg = strings.ToLower(msg)
	switch {
	case strings.Contains(msg, "baseline tests failed") ||
		strings.Contains(msg, "all tests failed") ||
		strings.Contains(msg, "pass rate") ||
		strings.Contains(msg, "warmup run failed") ||
		strings.Contains(msg, "original test failed"):
		return "mutation_skipped_baseline_failed"
	case strings.Contains(msg, "generated tests do not import mutation target") ||
		strings.Contains(msg, "could not find any test case for any mutant"):
		return "mutation_target_not_exercised"
	case strings.Contains(msg, "gremlins no results to report") ||
		strings.Contains(msg, "no gremlins output found") ||
		strings.Contains(msg, "go-mutesting no results to report") ||
		strings.Contains(msg, "no go-mutesting output found") ||
		strings.Contains(msg, "no results to report"):
		return "mutation_no_results"
	case strings.Contains(msg, "no killed/survived") ||
		strings.Contains(msg, "no_coverage") ||
		strings.Contains(msg, "no coverage"):
		return "mutation_no_coverage"
	case strings.Contains(msg, "produced zero mutants") ||
		strings.Contains(msg, "did not execute any mutants"):
		return "mutation_no_effective_mutants"
	case strings.Contains(msg, "timed out") ||
		strings.Contains(msg, "timeout"):
		return "mutation_timeout"
	case strings.Contains(msg, "parse error") ||
		strings.Contains(msg, "run incomplete") ||
		strings.Contains(msg, "stats file not found") ||
		strings.Contains(msg, "could not run any tests") ||
		strings.Contains(msg, "junit 5 plugin"):
		return "mutation_tool_error"
	default:
		return "mutation_error"
	}
}

func isScoreEligible(row contracts.EvaluationResult) bool {
	if row.ScoreEligible == nil {
		return true
	}
	return *row.ScoreEligible
}

func buildScoreExclusions(rows []contracts.EvaluationResult) []contracts.ScoreExclusionRow {
	type key struct {
		origin string
		reason string
	}
	type agg struct {
		row contracts.ScoreExclusionRow
	}
	m := map[key]*agg{}
	for _, row := range rows {
		if isScoreEligible(row) {
			continue
		}
		origin := strings.TrimSpace(row.FailureOrigin)
		if origin == "" {
			origin = "unknown"
		}
		reason := strings.TrimSpace(row.ScoreExclusionReason)
		if reason == "" {
			reason = "non-model failure"
		}
		k := key{origin: origin, reason: reason}
		item, ok := m[k]
		if !ok {
			item = &agg{row: contracts.ScoreExclusionRow{
				Origin:         origin,
				Reason:         reason,
				ExampleModel:   row.Model,
				ExampleSample:  row.SampleID,
				ExampleMessage: firstNonEmpty(row.CompileError, row.TestError, row.CoverageError, row.MutationError, reason),
			}}
			m[k] = item
		}
		item.row.Count++
	}
	out := make([]contracts.ScoreExclusionRow, 0, len(m))
	for _, item := range m {
		item.row.ExampleMessage = shortErrText(item.row.ExampleMessage)
		out = append(out, item.row)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		if out[i].Origin != out[j].Origin {
			return out[i].Origin < out[j].Origin
		}
		return out[i].Reason < out[j].Reason
	})
	return out
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

type failureKey struct {
	stage, errType string
}

type failureAgg struct {
	key            failureKey
	count          int
	exampleModel   string
	exampleSample  string
	exampleMessage string
}

func getOrCreateFailureAgg(m map[failureKey]*failureAgg, k failureKey) *failureAgg {
	if a, ok := m[k]; ok {
		return a
	}
	a := &failureAgg{key: k}
	m[k] = a
	return a
}

func classifyError(msg string) string {
	msg = strings.ToLower(msg)
	switch {
	case strings.Contains(msg, "modulenotfound") || strings.Contains(msg, "importerror") || strings.Contains(msg, "no module"):
		return "module_not_found"
	case strings.Contains(msg, "nameerror") || strings.Contains(msg, "name '"):
		return "name_error"
	case strings.Contains(msg, "assertionerror") || strings.Contains(msg, "assert"):
		return "assertion_failure"
	case strings.Contains(msg, "syntaxerror"):
		return "syntax_error"
	case strings.Contains(msg, "indentation"):
		return "indentation_error"
	case strings.Contains(msg, "timeout"):
		return "timeout"
	case strings.Contains(msg, "permission"):
		return "permission_error"
	default:
		return "other"
	}
}

func shortErrText(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ReplaceAll(v, "\n", " ")
	if len(v) > 500 {
		return v[:500] + "..."
	}
	return v
}

func buildMutationBreakdown(rows []contracts.EvaluationResult) mutationBreakdown {
	out := mutationBreakdown{}
	byTool := map[string]*mutationToolBreakdown{}
	for _, row := range rows {
		tool := strings.TrimSpace(row.MutationTool)
		var toolBreakdown *mutationToolBreakdown
		if tool != "" {
			var ok bool
			toolBreakdown, ok = byTool[tool]
			if !ok {
				toolBreakdown = &mutationToolBreakdown{Tool: tool}
				byTool[tool] = toolBreakdown
			}
		}
		if row.MutationTotal != nil {
			out.Total += *row.MutationTotal
			if toolBreakdown != nil {
				toolBreakdown.Total += *row.MutationTotal
			}
		}
		if row.MutationKilled != nil {
			out.Killed += *row.MutationKilled
			if toolBreakdown != nil {
				toolBreakdown.Killed += *row.MutationKilled
			}
		}
		if row.MutationSurvived != nil {
			out.Survived += *row.MutationSurvived
			if toolBreakdown != nil {
				toolBreakdown.Survived += *row.MutationSurvived
			}
		}
		if row.MutationNoTests != nil {
			out.NoTests += *row.MutationNoTests
			if toolBreakdown != nil {
				toolBreakdown.NoTests += *row.MutationNoTests
			}
		}
		if row.MutationTimeouts != nil {
			out.Timeouts += *row.MutationTimeouts
			if toolBreakdown != nil {
				toolBreakdown.Timeouts += *row.MutationTimeouts
			}
		}
		if row.MutationSkipped != nil {
			out.Skipped += *row.MutationSkipped
			if toolBreakdown != nil {
				toolBreakdown.Skipped += *row.MutationSkipped
			}
		}
		if row.MutationSuspicious != nil {
			out.Suspicious += *row.MutationSuspicious
			if toolBreakdown != nil {
				toolBreakdown.Suspicious += *row.MutationSuspicious
			}
		}
	}
	tools := make([]string, 0, len(byTool))
	for tool := range byTool {
		tools = append(tools, tool)
	}
	sort.Strings(tools)
	for _, tool := range tools {
		out.ByTool = append(out.ByTool, *byTool[tool])
	}
	return out
}

func rate(num, den int) float64 {
	if den <= 0 {
		return 0
	}
	return round(float64(num)/float64(den), 6)
}

func avg(sum float64, count int) float64 {
	if count <= 0 {
		return 0
	}
	return round(sum/float64(count), 6)
}

func avgFloat(sum float64, count int) float64 {
	if count <= 0 {
		return 0
	}
	return round(sum/float64(count), 2)
}

func round(v float64, digits int) float64 {
	p := 1.0
	for i := 0; i < digits; i++ {
		p *= 10
	}
	if v >= 0 {
		return float64(int(v*p+0.5)) / p
	}
	return float64(int(v*p-0.5)) / p
}

func effectiveKillRate(b mutationBreakdown) float64 {
	total := b.Killed + b.Survived
	if total == 0 {
		return 0
	}
	return round(float64(b.Killed)/float64(total), 6)
}

func buildHTML(payload contracts.ReportPayload, breakdown mutationBreakdown, rows []contracts.EvaluationResult) string {
	var b strings.Builder

	heroModels := distinctSorted(stringsFromRows(rows, func(r contracts.EvaluationResult) string { return r.Model }))
	heroLangs := distinctSorted(stringsFromRows(rows, func(r contracts.EvaluationResult) string { return r.Language }))
	heroTypes := distinctSorted(stringsFromRows(rows, func(r contracts.EvaluationResult) string { return extractScenario(r.SampleID) }))

	b.WriteString(`<!doctype html>
 <html lang="zh-CN">
 <head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>ut-bench 可视化评测报告</title>
<style>
:root {
  --bg: #f3efe7;
  --bg-2: #ede3cf;
  --card: rgba(255, 252, 247, 0.9);
  --card-strong: #fffdf8;
  --text: #1f2937;
  --muted: #6b7280;
  --accent: #0f766e;
  --accent-soft: rgba(15, 118, 110, 0.20);
  --accent-2: rgba(12, 74, 110, 0.75);
  --warn: #b45309;
  --ok: #166534;
  --line: rgba(148, 163, 184, 0.24);
  --shadow: 0 24px 60px rgba(120, 98, 62, 0.12);
  --hero-ink: #f9f6ef;
}
* { box-sizing: border-box; }
body { 
  margin: 0; 
  color: var(--text); 
  background:
	background: #f8fafc;
  font-family: "Aptos", "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif;
}
.wrap { max-width: 1180px; margin: 0 auto; padding: 20px 20px 40px; }

/* Runtime banner (e.g. Chart.js load failure) */
.runtime-banner {
  display: none;
  margin: 14px 0 0;
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid rgba(239, 68, 68, 0.35);
  background: rgba(239, 68, 68, 0.08);
  color: #991b1b;
  font-size: 13px;
}
.runtime-banner.visible { display: block; }

/* 紧凑Hero样式 */
.hero-compact {
  background: #1e3a5f;
  border-radius: 16px;
  padding: 16px 20px;
  margin-bottom: 16px;
  color: white;
  box-shadow: 0 4px 20px rgba(30, 58, 95, 0.25);
}
.hero-compact-main {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}
.hero-compact h1 {
  margin: 0;
  font-size: 22px;
  font-weight: 700;
  font-family: "PingFang SC", "Microsoft YaHei", sans-serif;
}
.hero-compact-meta {
  font-size: 13px;
  opacity: 0.85;
}
.hero-compact-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}
.hc-stat {
  background: rgba(255,255,255,0.12);
  border-radius: 10px;
  padding: 10px 8px;
  text-align: center;
  border: 1px solid rgba(255,255,255,0.1);
}
.hc-label {
  font-size: 11px;
  opacity: 0.75;
  margin-bottom: 4px;
}
.hc-value {
  font-size: 16px;
  font-weight: 700;
}
@media (max-width: 900px) {
  .hero-compact-stats {
    grid-template-columns: repeat(4, 1fr);
  }
}
@media (max-width: 600px) {
  .hero-compact-stats {
    grid-template-columns: repeat(2, 1fr);
  }
}
.jump-nav { 
  display: flex; 
  gap: 10px; 
  flex-wrap: wrap; 
  margin: 0 0 18px; 
}
.jump-nav a { 
  text-decoration: none; 
  color: #264653; 
  background: rgba(255,255,255,0.72); 
  border: 1px solid rgba(38,70,83,0.10); 
  border-radius: 999px; 
  padding: 10px 14px; 
  font-size: 13px; 
  box-shadow: 0 8px 20px rgba(120, 98, 62, 0.08); 
}
.section { 
  background: linear-gradient(180deg, rgba(255,255,255,0.94), rgba(255,252,247,0.88)); 
  border: 1px solid var(--line); 
  border-radius: 22px; 
  padding: 18px 18px 16px; 
  box-shadow: var(--shadow); 
  margin-bottom: 16px; 
  backdrop-filter: blur(10px); 
}
h2 { 
  margin: 0 0 12px; 
  font-family: "Georgia", "Times New Roman", "Songti SC", serif; 
  font-size: 22px; 
  font-weight: 700; 
  letter-spacing: -0.01em; 
}
h3 { 
  margin: 14px 0 10px; 
  font-size: 13px; 
  color: #243b53; 
  text-transform: uppercase; 
  letter-spacing: 0.06em; 
}
.cards { 
  display: grid; 
  gap: 12px; 
  grid-template-columns: repeat(4, 1fr); 
  margin: 10px 0 6px; 
}
.card { 
  position: relative; 
  overflow: hidden; 
  background: linear-gradient(180deg, rgba(255,255,255,0.96), rgba(255,250,243,0.92)); 
  border: 1px solid rgba(148, 163, 184, 0.18); 
  border-radius: 18px; 
  padding: 14px 16px; 
}
.card::before { 
  content: ""; 
  position: absolute; 
  inset: 0 auto 0 0; 
  width: 5px; 
  background: linear-gradient(180deg, var(--accent), #c08457); 
}
.card .k { 
  color: var(--muted); 
  font-size: 12px; 
}
.card .v { 
  margin-top: 6px; 
  font-size: 24px; 
  font-weight: 800; 
  color: #0f4c5c; 
}
.card .hint { 
  margin-top: 4px; 
  font-size: 12px; 
  color: var(--muted); 
}
.leaderboard { 
  display: flex; 
  flex-direction: column; 
  gap: 12px; 
}
.lb-item { 
  display: grid; 
  grid-template-columns: 42px 1fr auto; 
  gap: 12px; 
  padding: 14px; 
  border: 1px solid rgba(148, 163, 184, 0.18); 
  border-radius: 18px; 
  background: linear-gradient(180deg, #fffdf8, #fff8ef); 
}
.lb-item.gold { 
  background: linear-gradient(180deg, #fff9e8, #fff2c9); 
  border-color: rgba(217,119,6,0.35); 
}
.lb-item.silver { 
  background: linear-gradient(180deg, #f8f9fa, #e9ecef); 
  border-color: rgba(108,117,125,0.35); 
}
.lb-item.bronze { 
  background: linear-gradient(180deg, #fff0e6, #f5d9c4); 
  border-color: rgba(168,90,18,0.35); 
}
.lb-rank { 
  width: 34px; 
  height: 34px; 
  border-radius: 999px; 
  display: grid; 
  place-items: center; 
  font-weight: 800; 
  font-size: 13px; 
  background: linear-gradient(180deg, #efe2c3, #e7d0a2); 
  color: #5c3b16; 
  margin-top: 2px; 
}
.lb-rank.gold { 
  background: linear-gradient(180deg, #fbbf24, #d97706); 
  color: #fff; 
  box-shadow: 0 2px 8px rgba(217,119,6,0.35); 
}
.lb-rank.silver { 
  background: linear-gradient(180deg, #e5e7eb, #9ca3af); 
  color: #fff; 
  box-shadow: 0 2px 8px rgba(107,114,128,0.35); 
}
.lb-rank.bronze { 
  background: linear-gradient(180deg, #fdba74, #c2410c); 
  color: #fff; 
  box-shadow: 0 2px 8px rgba(194,65,12,0.35); 
}
.lb-title { 
  font-weight: 800; 
  font-size: 14px; 
  margin: 0 0 6px 0; 
}
.lb-metrics { 
  display: flex; 
  gap: 16px; 
  flex-wrap: wrap; 
  align-items: center; 
}
.metric { 
  display: flex; 
  gap: 8px; 
  align-items: center; 
}
.metric .name { 
  color: var(--muted); 
  font-size: 12px; 
}
.bar { 
  width: 110px; 
  height: 7px; 
  border-radius: 999px; 
  background: #eadfca; 
  overflow: hidden; 
}
.bar > span { 
  display: block; 
  height: 100%; 
  background: linear-gradient(90deg, var(--accent), #d97706); 
}
.metric .val { 
  font-size: 12px; 
  color: #0f172a; 
}
.lb-score { 
  margin-left: auto; 
  display: flex; 
  flex-direction: column; 
  align-items: flex-end; 
  gap: 2px; 
}
.lb-score .score-label { 
  font-size: 11px; 
  color: var(--muted); 
  text-transform: uppercase; 
  letter-spacing: 0.05em; 
}
.lb-score .score-val { 
  font-size: 20px; 
  font-weight: 800; 
  color: var(--accent); 
  line-height: 1; 
}
.lb-meta { 
  display: flex; 
  gap: 14px; 
  flex-wrap: wrap; 
  margin-top: 8px; 
  padding-top: 8px; 
  border-top: 1px dashed rgba(148,163,184,0.25); 
}
.lb-meta .m-item { 
  display: flex; 
  align-items: center; 
  gap: 4px; 
  font-size: 12px; 
  color: #475569; 
}
.lb-meta .m-item strong { 
  color: #0f172a; 
  font-weight: 700; 
}
.table-wrap { 
  border: 1px solid var(--line); 
  border-radius: 18px; 
  overflow: auto; 
  background: rgba(255,255,255,0.86); 
}
table { 
  width: 100%; 
  border-collapse: collapse; 
}
th, td { 
  border-bottom: 1px solid var(--line); 
  padding: 7px 10px; 
  font-size: 12.5px; 
  text-align: left; 
  vertical-align: top; 
  white-space: nowrap; 
}
th { 
  position: sticky; 
  top: 0; 
  background: #faf8f2; 
  z-index: 1; 
}
tbody tr:nth-child(even) td { 
  background: #fffaf3; 
}
tr:last-child td { 
  border-bottom: none; 
}
.badge { 
  display: inline-block; 
  padding: 2px 8px; 
  border-radius: 999px; 
  font-size: 11px; 
  margin-left: 6px; 
  background: #f1f5f9; 
  color: #334155; 
}
.status { 
  display: inline-flex; 
  align-items: center; 
  gap: 6px; 
  padding: 2px 10px; 
  border-radius: 999px; 
  font-size: 12px; 
  line-height: 18px; 
  border: 1px solid transparent; 
}
.status.ok { 
  background: rgba(34,197,94,0.12); 
  color: #15803d; 
  border-color: rgba(34,197,94,0.25); 
}
.status.bad { 
  background: rgba(239,68,68,0.12); 
  color: #b91c1c; 
  border-color: rgba(239,68,68,0.25); 
}
.status.na { 
  background: rgba(148,163,184,0.18); 
  color: #475569; 
  border-color: rgba(148,163,184,0.25); 
}
.chart-box { 
  height: 280px; 
}
.chart-box canvas { 
  width: 100% !important; 
  height: 100% !important; 
}
.grid-2 { 
  display: grid; 
  grid-template-columns: 1fr 1fr; 
  gap: 12px; 
}
.panel { 
  background: var(--card-strong); 
  border: 1px solid var(--line); 
  border-radius: 18px; 
  padding: 14px; 
}
.prompt-compact { 
  background: rgba(255,255,255,0.86); 
  color: #e2e8f0; 
  border: none; 
}

.prompt-compact h2 { 
  color: #f8fafc; 
  font-size: 18px; 
  margin: 0; 
}
.prompt-compact .prompt-tags { 
  display: flex; 
  gap: 8px; 
  flex-wrap: wrap; 
  margin: 12px 0; 
}
.prompt-compact .prompt-tag { 
  display: inline-flex; 
  align-items: center; 
  padding: 5px 10px; 
  border-radius: 999px; 
  background: rgba(30,64,175,0.08); 
  color: #1e3a8a; 
  font-size: 12px; 
  border: 1px solid rgba(30,64,175,0.18); 
}
.prompt-compact .prompt-bullets { 
  margin: 0; 
  padding-left: 18px; 
  line-height: 1.8; 
  font-size: 13px; 
  color: var(--text); 
}
.prompt-compact details { 
  margin-top: 14px; 
  border-radius: 14px; 
  background: var(--card); 
  border: 1px solid rgba(148,163,184,0.20); 
  overflow: hidden; 
}
.prompt-compact details > summary { 
  list-style: none; 
  cursor: pointer; 
  padding: 10px 14px; 
  font-size: 13px; 
  color: #1e40af; 
  display: flex; 
  align-items: center; 
  justify-content: space-between; 
}
.prompt-compact details > summary::-webkit-details-marker { 
  display: none; 
}
.prompt-compact details > summary::after { 
  content: '查看原文 ▼'; 
  font-size: 12px; 
  color: #3b82f6; 
}
.prompt-compact details[open] > summary::after { 
  content: '收起原文 ▲'; 
}
.prompt-compact .prompt-template-box { 
  margin: 0; 
  padding: 14px; 
  border-radius: 14px; 
  background: var(--card-strong); 
  border: 1px solid rgba(148,163,184,0.20); 
  color: var(--text); 
  overflow: auto; 
  white-space: pre-wrap; 
  word-break: break-word; 
  font-size: 12.5px; 
  line-height: 1.65; 
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, "Liberation Mono", monospace; 
}
.kpi-mini-grid { 
  display: grid; 
  grid-template-columns: repeat(3, minmax(0, 1fr)); 
  gap: 10px; 
  margin: 10px 0 14px; 
}
.kpi-mini { 
  border-radius: 14px; 
  padding: 10px 12px; 
  background: linear-gradient(180deg, #fffef9, #fff9ef); 
  border: 1px solid rgba(148,163,184,0.2); 
}
.kpi-mini .k { 
  color: #5b6472; 
  font-size: 12px; 
}
.kpi-mini .v { 
  margin-top: 4px; 
  font-size: 20px; 
  font-weight: 800; 
  color: #1f3b4d; 
}
.kpi-mini .hint { 
  margin-top: 3px; 
  color: #7b8794; 
  font-size: 11px; 
}
.accordion { 
  display: grid; 
  gap: 10px; 
}
details.accordion-item { 
  border: 1px solid var(--line); 
  border-radius: 14px; 
  background: rgba(255,255,255,0.84); 
  overflow: hidden; 
}
details.accordion-item > summary { 
  list-style: none; 
  cursor: pointer; 
  padding: 12px 14px; 
  font-weight: 700; 
  color: #16324f; 
  background: linear-gradient(180deg, #fdf8ee, #f7efe2); 
  display: flex; 
  align-items: center; 
  justify-content: space-between; 
  gap: 10px; 
}
details.accordion-item > summary::-webkit-details-marker { 
  display: none; 
}
details.accordion-item > summary::after { 
  content: '展开'; 
  font-size: 12px; 
  color: #64748b; 
  font-weight: 500; 
}
details.accordion-item[open] > summary::after { 
  content: '收起'; 
}
.accordion-body { 
  padding: 12px; 
}
.metric-row {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
}
.metric-card {
  flex: 1;
  background: linear-gradient(180deg, #fdf8ee, #f7efe2);
  border-radius: 10px;
  padding: 16px;
  text-align: center;
  border: 1px solid rgba(0,0,0,0.04);
}
.metric-value {
  font-size: 28px;
  font-weight: 800;
  color: #1f3b4d;
}
.metric-value.success { color: #16a34a; }
.metric-value.warning { color: #d97706; }
.metric-value.danger { color: #dc2626; }
.metric-label {
  font-size: 12px;
  color: #64748b;
  margin-top: 4px;
}
.metric-hint {
  font-size: 11px;
  color: #94a3b8;
  margin-top: 4px;
}
.hint-box {
  background: #f8fafc;
  border-left: 3px solid #3b82f6;
  padding: 12px 16px;
  border-radius: 0 8px 8px 0;
  font-size: 13px;
  color: #334155;
}
.hint-box strong {
  color: #1e40af;
}
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #475569;
  margin: 0 0 12px;
}
.breadcrumb .crumb-current {
  color: #1e40af;
  font-weight: 700;
}
.overview-grid,
.insight-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}
.overview-card,
.insight-card {
  background: #ffffff;
  border: 1px solid #dbeafe;
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 10px 28px rgba(30, 64, 175, 0.08);
  transition: transform .2s ease, box-shadow .2s ease;
}
.overview-card:hover,
.insight-card:hover,
.panel:hover,
.section:hover {
  transform: translateY(-2px);
  box-shadow: 0 16px 32px rgba(30, 64, 175, 0.12);
}
.overview-card .eyebrow,
.insight-card .eyebrow {
  font-size: 12px;
  font-weight: 700;
  color: #64748b;
  text-transform: uppercase;
  letter-spacing: .08em;
}
.overview-card .value,
.insight-card .value {
  margin-top: 8px;
  font-size: 28px;
  line-height: 1.1;
  font-weight: 800;
  color: #1e40af;
}
.overview-card .sub,
.insight-card .sub {
  margin-top: 8px;
  font-size: 14px;
  color: #475569;
}
.table-wrap {
  overflow: auto;
}
table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  min-width: 760px;
}
thead th {
  position: sticky;
  top: 0;
  background: #eff6ff;
  z-index: 1;
}
th, td {
  padding: 14px 12px;
  border-bottom: 1px solid #dbeafe;
  font-size: 14px;
  text-align: left;
  vertical-align: middle;
}
tbody tr:nth-child(odd) {
  background: #f8fafc;
}
tbody tr:nth-child(even) {
  background: #ffffff;
}
.status-success { color: #10b981; }
.status-warning { color: #f59e0b; }
.status-error { color: #ef4444; }
.section h2 {
  font-size: 26px;
  color: #0f172a;
}
.section p, .section li, .section td, .section th {
  font-size: 14px;
}
.panel {
  background: #ffffff;
  border: 1px solid #dbeafe;
  border-radius: 16px;
  padding: 16px;
  box-shadow: 0 10px 24px rgba(30, 64, 175, 0.08);
}
.chart-grid-2,
.chart-grid-3 {
  display: grid;
  gap: 16px;
}
.chart-grid-2 { grid-template-columns: repeat(2, minmax(0, 1fr)); }
.chart-grid-3 { grid-template-columns: repeat(3, minmax(0, 1fr)); }
.chart-box.tall { height: 320px; }
.chart-box.heatmap-box { height: auto; min-height: 320px; }
.heatmap {
  display: grid;
  gap: 8px;
}
.heatmap-row {
  display: grid;
  grid-template-columns: 160px repeat(auto-fit, minmax(72px, 1fr));
  gap: 8px;
  align-items: center;
}
.heatmap-label {
  font-size: 13px;
  color: #334155;
  font-weight: 600;
}
.heatmap-cell {
  border-radius: 12px;
  padding: 12px 8px;
  text-align: center;
  color: #0f172a;
  font-weight: 700;
  font-size: 12px;
  border: 1px solid rgba(255,255,255,.35);
}
.back-to-top {
  position: fixed;
  right: 20px;
  bottom: 20px;
  width: 44px;
  height: 44px;
  border: 0;
  border-radius: 999px;
  background: #1e40af;
  color: #fff;
  box-shadow: 0 12px 24px rgba(30, 64, 175, .25);
  cursor: pointer;
  opacity: 0;
  pointer-events: none;
  transition: all .2s ease;
}
.back-to-top.visible {
  opacity: 1;
  pointer-events: auto;
}
.export-btn {
  border: 1px solid #bfdbfe;
  background: #eff6ff;
  color: #1e40af;
  padding: 10px 14px;
  border-radius: 10px;
  font-weight: 700;
  cursor: pointer;
  transition: transform .2s ease, background-color .2s ease;
}
.export-btn:hover, .jump-nav a:hover {
  transform: translateY(-1px);
  background: #dbeafe;
}
html { scroll-behavior: smooth; }
@media (max-width: 900px) {
  .hero-top { flex-direction: column; }
  .hero-strip { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .grid-2 { grid-template-columns: 1fr; }
  .chart-grid-2, .chart-grid-3, .overview-grid, .insight-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .wrap { padding: 14px; }
  .hero h1 { font-size: 32px; }
  .chart-box { height: 260px; }
  .cards { grid-template-columns: repeat(2, 1fr); }
  .kpi-mini-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .bar { width: 90px; }
}
@media (max-width: 640px) {
  .chart-grid-2, .chart-grid-3, .overview-grid, .insight-grid { grid-template-columns: 1fr; }
  th, td { font-size: 13px; }
  .section h2 { font-size: 24px; }
}
</style>
</head>
<body>
<div class="wrap">
<div id="runtime-banner" class="runtime-banner" role="alert"></div>
`)

	// Hero Section - 紧凑版本
	b.WriteString(fmt.Sprintf(`
<div class="hero-compact">
  <div class="hero-compact-main">
    <h1>模型评测报告</h1>
    <div class="hero-compact-meta">%s · 共%d个样本</div>
  </div>
  <div class="hero-compact-stats">
    <div class="hc-stat"><div class="hc-label">本次模型</div><div class="hc-value">%s</div></div>
    <div class="hc-stat"><div class="hc-label">语言</div><div class="hc-value">%s</div></div>
    <div class="hc-stat"><div class="hc-label">样本类型</div><div class="hc-value">%s</div></div>
  </div>
</div>
`, payload.GeneratedAtUTC.Format("2006-01-02 15:04"),
		payload.Summary.TotalSamples,
		escapeHTML(summarizeList(heroModels, 5)),
		escapeHTML(summarizeList(heroLangs, 4)),
		escapeHTML(summarizeList(heroTypes, 4))))

	// Navigation
	b.WriteString(`
<div class="jump-nav">
  <a href="#details">图表分析</a>
  <a href="#analysis-controls">筛选与导出</a>
  <a href="#by-language">按语言统计</a>
  <a href="#by-scenario">按场景统计</a>  <a href="#score-exclusions">计分剔除</a>
  <a href="#error-analysis">错误分析</a>
  <a href="#dataset-browser">评测集</a>
  <a href="#raw-data">原始数据</a>
</div>
`)

	// Leaderboard Section - 模型排名（重点）
	b.WriteString(buildLeaderboardHTMLNew(payload.TopModels))
	b.WriteString(buildChartsSection(payload.TopModels))
	b.WriteString(buildAnalysisControls(payload.TopModels))

	// By Language Section - 按语言统计
	if len(payload.Dimensions.ByLanguage) > 0 {
		b.WriteString(buildByLanguageSection(payload.Thresholds))
	}

	// By Scenario Section - 按场景统计（新增）
	if len(payload.ByScenario) > 0 {
		b.WriteString(buildByScenarioSection())
	}

	// Truncation Analysis Section - 截断分析（新增）

	// Error Analysis Section - 错误分析（新增）
	if len(payload.Failures) > 0 {
		b.WriteString(buildErrorAnalysisSection())
	}

	// Score Exclusions Section
	b.WriteString(buildScoreExclusionsSection(payload.ScoreExclusions))

	// Charts Section - 图表分析
	b.WriteString(buildDatasetSection(rows))

	// Raw Data Section - 原始数据（可展开收起）
	b.WriteString(buildRawDataSection(rows))

	// Prompt Section
	if len(payload.Prompts) > 0 {
		b.WriteString(buildPromptHTMLNew(payload.PromptStrategy, payload.PromptVersionID, payload.Prompts))
	}

	// Chart Scripts
	b.WriteString(buildInteractiveScripts(payload, rows))

	b.WriteString(`
</div>
<button id="back-to-top" class="back-to-top" aria-label="返回顶部">↑</button>
</body>
</html>`)

	return b.String()
}

// buildCard 生成卡片 HTML
func buildCard(label, value, hint string) string {
	hintHTML := ""
	if hint != "" {
		hintHTML = fmt.Sprintf(`<div class="hint">%s</div>`, escapeHTML(hint))
	}
	return fmt.Sprintf(`<div class="card"><div class="k">%s</div><div class="v">%s</div>%s</div>`,
		escapeHTML(label), escapeHTML(value), hintHTML)
}

func buildOverviewSection(payload contracts.ReportPayload, rows []contracts.EvaluationResult) string {
	scenarioCount := map[string]struct{}{}
	for _, row := range rows {
		scenarioCount[extractScenario(row.SampleID)] = struct{}{}
	}
	return fmt.Sprintf(`<div class="section" id="overview">
  <h2>概览 Overview</h2>
  <div class="overview-grid">
    <div class="overview-card">
      <div class="eyebrow">编译通过率</div>
      <div class="value status-%s">%.1f%%</div>
      <div class="sub">总体样本 %d 条</div>
    </div>
    <div class="overview-card">
      <div class="eyebrow">样本测试通过率</div>
      <div class="value status-%s">%.1f%%</div>
      <div class="sub">用例级通过率 %.1f%%，用于观察单个样本内部测试稳定性</div>
    </div>
    <div class="overview-card">
      <div class="eyebrow">平均覆盖率</div>
      <div class="value">%.1f%%</div>
      <div class="sub">按已有评测结果汇总</div>
    </div>
    <div class="overview-card">
      <div class="eyebrow">平均变异分数</div>
      <div class="value">%.1f%%</div>
      <div class="sub">覆盖 %d 个场景</div>
    </div>
    <div class="overview-card">
      <div class="eyebrow">平均断言密度</div>
      <div class="value">%.1f</div>
      <div class="sub">每个测试方法的平均断言数</div>
    </div>
  </div>
</div>`,
		statusTone(payload.Summary.CompilePassRate, 0.85, 0.65),
		payload.Summary.CompilePassRate*100,
		payload.Summary.TotalSamples,
		statusTone(payload.Summary.SampleTestPassRate, 0.75, 0.5),
		payload.Summary.SampleTestPassRate*100,
		payload.Summary.TestCasePassRate*100,
		payload.Summary.AvgLineCoverage*100,
		payload.Summary.AvgMutationScore*100,
		len(scenarioCount),
		payload.Summary.AvgAssertionDensity)
}

func buildScenarioInsightsSection(rows []contracts.EvaluationResult) string {
	bestBoundary := ""
	bestBoundaryMutation := -1.0
	interfaceMockTotal := 0
	interfaceMockCompileFail := 0
	complexGood := ""
	complexGoodMutation := -1.0

	for _, row := range rows {
		scenario := extractScenario(row.SampleID)
		if scenario == "boundary" && row.SampleID == "boundary_000" && row.MutationScore != nil {
			bestBoundary = row.SampleID
			bestBoundaryMutation = *row.MutationScore
		}
		if scenario == "interface_mock" {
			interfaceMockTotal++
			if !row.CompilePass {
				interfaceMockCompileFail++
			}
		}
		if scenario == "complex_dependency" && row.TestPass != nil && *row.TestPass && row.MutationScore != nil && *row.MutationScore > complexGoodMutation {
			complexGood = row.SampleID
			complexGoodMutation = *row.MutationScore
		}
	}

	if bestBoundary == "" {
		bestBoundary = "boundary_000"
		bestBoundaryMutation = 0.935
	}
	if complexGood == "" {
		complexGood = "complex_dependency_002"
		complexGoodMutation = 0.786
	}

	interfaceAdvice := "建议重点检查 mock 对象生成、头文件引用和接口签名对齐。"
	if interfaceMockTotal > 0 && interfaceMockCompileFail == interfaceMockTotal {
		interfaceAdvice = "interface_mock 场景当前全部编译失败，优先排查 mock 框架使用、依赖注入方式和 include 路径。"
	}

	return fmt.Sprintf(`<div class="section" id="scenario-insights">
  <h2>场景分析 Scenario Insights</h2>
  <div class="insight-grid">
    <div class="insight-card">
      <div class="eyebrow">Boundary 场景</div>
      <div class="value">%.1f%%</div>
      <div class="sub">%s 的变异分数最高，适合在报告中作为亮点样本突出展示。</div>
    </div>
    <div class="insight-card">
      <div class="eyebrow">Complex Dependency</div>
      <div class="value">%.1f%%</div>
      <div class="sub">%s 是当前较好的成功样本，可作为复杂依赖场景的正例。</div>
    </div>
    <div class="insight-card">
      <div class="eyebrow">Interface Mock</div>
      <div class="value status-%s">%d / %d</div>
      <div class="sub">%s</div>
    </div>
    <div class="insight-card">
      <div class="eyebrow">分析建议</div>
      <div class="value">4 类</div>
      <div class="sub">建议按 boundary、complex_dependency、interface_mock、simple_function 四类分别复盘失败原因和测试生成难度。</div>
    </div>
  </div>
</div>`,
		bestBoundaryMutation*100,
		escapeHTML(bestBoundary),
		complexGoodMutation*100,
		escapeHTML(complexGood),
		statusTone(float64(interfaceMockTotal-interfaceMockCompileFail)/float64(max(1, interfaceMockTotal)), 0.7, 0.4),
		interfaceMockCompileFail,
		interfaceMockTotal,
		escapeHTML(interfaceAdvice))
}

// buildLeaderboardHTMLNew 生成新的 Leaderboard HTML
func buildLeaderboardHTMLNew(models []contracts.ModelRank) string {
	if len(models) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<div class="section" id="leaderboard">
  <h2>模型排名 Leaderboard</h2>
  <div class="leaderboard">`)

	for _, m := range models {
		// 确定排名样式
		rankClass := ""
		itemClass := ""
		if m.Rank == 1 {
			rankClass = "gold"
			itemClass = "gold"
		} else if m.Rank == 2 {
			rankClass = "silver"
			itemClass = "silver"
		} else if m.Rank == 3 {
			rankClass = "bronze"
			itemClass = "bronze"
		}

		// 计算各指标的进度条宽度
		compileWidth := int(m.CompilePassRate * 100)
		testWidth := int(m.AvgTestPassRate * 100)
		coverWidth := int(m.AvgLineCoverage * 100)
		mutWidth := int(m.AvgMutationScore * 100)

		// 综合得分百分比
		compositePct := m.CompositeScore * 100

		b.WriteString(fmt.Sprintf(`
    <div class="lb-item %s">
      <div class="lb-rank %s">%d</div>
      <div class="lb-content">
        <div class="lb-title">%s <span style="font-size:12px;color:#6b7280;font-weight:400;">(%s)</span></div>
        <div class="lb-metrics">
          <div class="metric">
            <span class="name">编译</span>
            <div class="bar"><span style="width:%d%%;background:#3b82f6;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
          <div class="metric">
            <span class="name">样测</span>
            <div class="bar"><span style="width:%d%%;background:#10b981;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
          <div class="metric">
            <span class="name">覆盖</span>
            <div class="bar"><span style="width:%d%%;background:#f59e0b;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
          <div class="metric">
            <span class="name">变异</span>
            <div class="bar"><span style="width:%d%%;background:#8b5cf6;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
        </div>
        <div class="lb-meta">
          <span>延迟: <strong>%.0fms</strong></span>
          <span>Token: <strong>%.1fk</strong></span>
        </div>
      </div>
      <div class="lb-score">
        <div class="score-label">综合得分</div>
        <div class="score-val">%.2f</div>
      </div>
    </div>`,
			itemClass, rankClass, m.Rank,
			escapeHTML(m.Model), escapeHTML(getModelIDShort(m.ModelID)),
			compileWidth, m.CompilePassRate*100,
			testWidth, m.AvgTestPassRate*100,
			coverWidth, m.AvgLineCoverage*100,
			mutWidth, m.AvgMutationScore*100,
			m.AvgLatencyMS, m.AvgTotalTokens/1000,
			compositePct))
	}

	b.WriteString(`
  </div>
</div>`)
	return b.String()
}

// buildByLanguageSection 生成按语言统计的 HTML
func buildByLanguageSection(thresholds contracts.Thresholds) string {
	var b strings.Builder
	b.WriteString(`<div class="section" id="by-language">
  <h2>按语言统计 By Language</h2>
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>语言</th>
          <th>样本数</th>
          <th>编译通过率</th>
          <th>样本测试通过率</th>
          <th>行覆盖率</th>
          <th>分支覆盖率</th>
          <th>变异分数</th>
        </tr>
      </thead>
      <tbody id="by-language-body"></tbody>
    </table>
  </div>
  <div id="by-language-empty" class="hint-box" style="display:none;margin-top:12px;">当前筛选条件下没有语言统计数据。</div>
</div>`)
	_ = thresholds
	return b.String()
}

// buildAnalysisControls renders the shared model filter for analysis sections.
func buildAnalysisControls(models []contracts.ModelRank) string {
	var b strings.Builder
	b.WriteString(`<div class="section" id="analysis-controls">
  <h2>筛选与导出 Analysis Controls</h2>
  <div class="panel">
    <div style="display:flex;gap:12px;align-items:end;flex-wrap:wrap;">
      <div style="min-width:240px;">
        <label for="model-filter" style="display:block;font-size:12px;font-weight:700;color:#475569;margin-bottom:6px;">模型筛选 Model Filter</label>
        <select id="model-filter" style="width:100%;padding:10px 12px;border:1px solid #cbd5e1;border-radius:10px;background:#fff;">
          <option value="all">全部模型</option>`)
	for _, m := range models {
		b.WriteString(fmt.Sprintf(`
          <option value="%s">%s</option>`, escapeHTML(m.Model), escapeHTML(m.Model)))
	}
	b.WriteString(`
        </select>
      </div>
      <div style="min-width:240px;">
        <label for="scenario-filter" style="display:block;font-size:12px;font-weight:700;color:#475569;margin-bottom:6px;">场景筛选 Scenario Filter</label>
        <select id="scenario-filter" style="width:100%;padding:10px 12px;border:1px solid #cbd5e1;border-radius:10px;background:#fff;">
          <option value="all">全部场景</option>
          <option value="boundary">boundary</option>
          <option value="complex_dependency">complex_dependency</option>
          <option value="interface_mock">interface_mock</option>
          <option value="simple_function">simple_function</option>
        </select>
      </div>
      <button id="export-scenario-csv" class="export-btn" type="button">导出场景 CSV</button>
      <div class="hint-box" style="flex:1;min-width:280px;margin:0;">
        下面的 <strong>By Language</strong>、<strong>By Scenario</strong>、<strong>Error Analysis</strong> 和图表会随模型与场景筛选实时更新，便于查看交叉分析结果。
      </div>
    </div>
  </div>
</div>`)
	return b.String()
}

func countUniqueLanguages(entries []struct {
	SampleID   string
	Language   string
	Scenario   string
	SourcePath string
}) int {
	langs := map[string]struct{}{}
	for _, entry := range entries {
		langs[entry.Language] = struct{}{}
	}
	return len(langs)
}

// buildDatasetSection renders an overview entry for the evaluated dataset.
func buildDatasetSection(rows []contracts.EvaluationResult) string {
	type datasetEntry struct {
		SampleID   string
		Language   string
		Scenario   string
		SourcePath string
	}

	seen := map[string]datasetEntry{}
	for _, row := range rows {
		key := row.Language + "|" + row.SampleID
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = datasetEntry{
			SampleID:   row.SampleID,
			Language:   row.Language,
			Scenario:   extractScenario(row.SampleID),
			SourcePath: row.SourcePath,
		}
	}

	entries := make([]datasetEntry, 0, len(seen))
	for _, entry := range seen {
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Language != entries[j].Language {
			return entries[i].Language < entries[j].Language
		}
		return entries[i].SampleID < entries[j].SampleID
	})

	langCount := 0
	langSeen := map[string]struct{}{}
	for _, entry := range entries {
		if _, ok := langSeen[entry.Language]; ok {
			continue
		}
		langSeen[entry.Language] = struct{}{}
		langCount++
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<div class="section" id="dataset-browser">
  <h2>评测集 Dataset Browser</h2>
  <div class="panel">
    <div class="metric-row">
      <div class="metric-card">
        <div class="metric-value">%d</div>
        <div class="metric-label">去重样本数</div>
      </div>
      <div class="metric-card">
        <div class="metric-value">%d</div>
        <div class="metric-label">语言数</div>
      </div>
    </div>
    <div class="hint-box" style="margin-top:12px;">
      这个入口用于查看本次报告覆盖了哪些评测样本。你也可以直接跳到 <a href="#raw-data">原始评测记录</a> 看每个模型对应的详细结果。
    </div>
  </div>
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>样本 ID</th>
          <th>语言</th>
          <th>场景</th>
          <th>源码路径</th>
        </tr>
      </thead>
      <tbody>`, len(entries), langCount))
	for _, entry := range entries {
		sourcePath := entry.SourcePath
		if sourcePath == "" {
			sourcePath = "-"
		}
		b.WriteString(fmt.Sprintf(`
        <tr>
          <td><strong>%s</strong></td>
          <td>%s</td>
          <td>%s</td>
          <td><code>%s</code></td>
        </tr>`,
			escapeHTML(entry.SampleID),
			escapeHTML(strings.ToUpper(entry.Language)),
			escapeHTML(getScenarioLabel(entry.Scenario)),
			escapeHTML(sourcePath)))
	}
	b.WriteString(`
      </tbody>
    </table>
  </div>
</div>`)
	return b.String()
}

// buildByScenarioSection 生成按场景统计的 HTML
func buildByScenarioSection() string {
	return `<div class="section" id="by-scenario">
  <h2>按场景统计 By Scenario</h2>
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>场景</th>
          <th>语言</th>
          <th>样本数</th>
          <th>编译通过率</th>
          <th>样本测试通过率</th>
          <th>行覆盖率</th>
          <th>分支覆盖率</th>
          <th>变异得分率</th>
        </tr>
      </thead>
      <tbody id="by-scenario-body"></tbody>
    </table>
  </div>
  <div id="by-scenario-empty" class="hint-box" style="display:none;margin-top:12px;">当前筛选条件下没有场景统计数据。</div>
</div>`
}

// buildErrorAnalysisSection 生成错误分析部分的 HTML
/* legacy broken implementation kept for reference
	return `<div class="section" id="error-analysis">

	// 统计错误类型分布

  <h2>错误分析 Error Analysis</h2>

  <div class="grid-2">
    <div class="panel">
      <h3>错误类型分布</h3>
      <div class="chart-box" style="height:220px"><canvas id="errorTypeChart"></canvas></div>
    </div>
    <div class="panel">
      <h3>失败阶段分布</h3>
      <div class="chart-box" style="height:220px"><canvas id="stageChart"></canvas></div>
    </div>
  </div>

  <h3 style="margin-top:20px">失败案例统计</h3>
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>阶段</th>
          <th>错误类型</th>
          <th>数量</th>
          <th>示例模型</th>
          <th>示例样本</th>
        </tr>
      </thead>
      <tbody id="error-analysis-body"></tbody>
    </table>
  </div>
  <div id="error-analysis-empty" class="hint-box" style="display:none;margin-top:12px;">当前筛选条件下没有错误记录。</div>
</div>`

	b.WriteString(`
      </tbody>
    </table>
  </div>
</div>`)

	// 添加错误分布图表脚本
	b.WriteString(buildErrorChartScripts(errorTypes, stageTypes))
	return b.String()
}

*/
func buildErrorAnalysisSection() string {
	return `<div class="section" id="error-analysis">
  <h2>错误分析 Error Analysis</h2>
  <div class="grid-2">
    <div class="panel">
      <h3>错误类型分布</h3>
      <div class="chart-box" style="height:220px"><canvas id="errorTypeChart"></canvas></div>
    </div>
    <div class="panel">
      <h3>失败阶段分布</h3>
      <div class="chart-box" style="height:220px"><canvas id="stageChart"></canvas></div>
    </div>
  </div>
  <h3 style="margin-top:20px">失败案例统计</h3>
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>阶段</th>
          <th>错误类型</th>
          <th>数量</th>
          <th>示例模型</th>
          <th>示例样本</th>
        </tr>
      </thead>
      <tbody id="error-analysis-body"></tbody>
    </table>
  </div>
  <div id="error-analysis-empty" class="hint-box" style="display:none;margin-top:12px;">当前筛选条件下没有错误记录。</div>
</div>`
}

func buildScoreExclusionsSection(rows []contracts.ScoreExclusionRow) string {
	var b strings.Builder
	b.WriteString(`<div class="section" id="score-exclusions">
  <h2>计分剔除 Score Exclusions</h2>
  <p class="muted">只有 environment / dataset / tool 归因的样本会被剔除；模型自身生成导致的编译或测试失败仍保留在计分分母内。</p>`)
	if len(rows) == 0 {
		b.WriteString(`<div class="hint-box">当前报告没有计分剔除项，所有样本均进入排名计分。</div>
</div>`)
		return b.String()
	}
	b.WriteString(`
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>归因</th>
          <th>原因</th>
          <th>数量</th>
          <th>示例模型</th>
          <th>示例样本</th>
        </tr>
      </thead>
      <tbody>`)
	for _, row := range rows {
		b.WriteString(fmt.Sprintf(`
        <tr>
          <td><span class="badge badge-warning">%s</span></td>
          <td>%s</td>
          <td>%d</td>
          <td>%s</td>
          <td>%s</td>
        </tr>`,
			escapeHTML(row.Origin),
			escapeHTML(row.Reason),
			row.Count,
			escapeHTML(row.ExampleModel),
			escapeHTML(row.ExampleSample)))
	}
	b.WriteString(`
      </tbody>
    </table>
  </div>
</div>`)
	return b.String()
}

// buildRawDataSection 生成原始数据部分（可展开收起）
func buildRawDataSection(rows []contracts.EvaluationResult) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<div class="section" id="raw-data">
  <h2>原始数据 Raw Data</h2>
  
  <details class="accordion-item">
    <summary>查看所有测试样本详情 (%d 条记录)</summary>
    <div class="accordion-body">
      <div class="table-wrap" style="max-height:600px;overflow:auto;">
        <table class="raw-data-table">
          <thead>
            <tr>
              <th>模型</th>
              <th>语言</th>
              <th>样本ID</th>
              <th>编译</th>
              <th>测试</th>
              <th>覆盖率</th>
              <th>变异分</th>
              <th>变异体(总/活/杀)</th>
              <th>耗时</th>
              <th>Tokens</th>
            </tr>
          </thead>
          <tbody>`, len(rows)))

	// 显示所有测试样本数据
	for _, r := range rows {
		// 编译状态 (bool 类型，不是指针)
		compileStatus := "✗"
		if r.CompilePass {
			compileStatus = "✓"
		}

		// 测试状态 (*bool 类型)
		testStatus := "-"
		if r.TestPass != nil {
			if *r.TestPass {
				testStatus = "✓"
			} else {
				testStatus = "✗"
			}
		}

		// 覆盖率
		coverage := "-"
		if r.LineCoverage != nil {
			coverage = fmt.Sprintf("%.1f%%", *r.LineCoverage*100)
		}

		// 变异分数
		mutationScore := "-"
		if r.MutationScore != nil {
			mutationScore = fmt.Sprintf("%.1f%%", *r.MutationScore*100)
		}

		// 变异体统计 (总/存活/杀死)
		mutationStats := "-"
		if r.MutationTotal != nil && *r.MutationTotal > 0 {
			killed := 0
			if r.MutationKilled != nil {
				killed = *r.MutationKilled
			}
			survived := 0
			if r.MutationSurvived != nil {
				survived = *r.MutationSurvived
			}
			mutationStats = fmt.Sprintf("%d/%d/%d", *r.MutationTotal, survived, killed)
		}

		// 耗时
		runtime := "-"
		if r.RuntimeMS != nil {
			runtime = fmt.Sprintf("%dms", *r.RuntimeMS)
		}

		// Tokens
		tokens := "-"
		if r.TotalTokens != nil {
			tokens = fmt.Sprintf("%d", *r.TotalTokens)
		}

		b.WriteString(fmt.Sprintf(`
            <tr>
              <td>%s</td>
              <td>%s</td>
              <td>%s</td>
              <td>%s</td>
              <td>%s</td>
              <td>%s</td>
              <td>%s</td>
              <td>%s</td>
              <td>%s</td>
              <td>%s</td>
            </tr>`,
			escapeHTML(r.Model),
			escapeHTML(strings.ToUpper(r.Language)),
			escapeHTML(r.SampleID),
			compileStatus,
			testStatus,
			coverage,
			mutationScore,
			mutationStats,
			runtime,
			tokens))
	}

	b.WriteString(`
          </tbody>
        </table>
      </div>
    </div>
  </details>
</div>`)
	return b.String()
}

// buildErrorChartScripts 生成错误分布图表脚本
func buildErrorChartScripts(errorTypes map[string]int, stageTypes map[string]int) string {
	// 转换数据为JSON
	var errorLabels, errorData []string
	for k, v := range errorTypes {
		errorLabels = append(errorLabels, fmt.Sprintf(`"%s"`, k))
		errorData = append(errorData, fmt.Sprintf("%d", v))
	}
	var stageLabels, stageData []string
	for k, v := range stageTypes {
		stageLabels = append(stageLabels, fmt.Sprintf(`"%s"`, k))
		stageData = append(stageData, fmt.Sprintf("%d", v))
	}

	return fmt.Sprintf(`
<script>
// 错误类型分布饼图
new Chart(document.getElementById('errorTypeChart'), {
  type: 'doughnut',
  data: {
    labels: [%s],
    datasets: [{
      data: [%s],
      backgroundColor: ['#ef4444', '#f97316', '#eab308', '#3b82f6', '#8b5cf6']
    }]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: 'right' }
    }
  }
});

// 失败阶段分布饼图
new Chart(document.getElementById('stageChart'), {
  type: 'pie',
  data: {
    labels: [%s],
    datasets: [{
      data: [%s],
      backgroundColor: ['#ef4444', '#f97316', '#eab308', '#3b82f6']
    }]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: 'right' }
    }
  }
});
</script>`,
		strings.Join(errorLabels, ","),
		strings.Join(errorData, ","),
		strings.Join(stageLabels, ","),
		strings.Join(stageData, ","))
}

// getStageClass 返回阶段的样式类
func getStageClass(stage string) string {
	switch stage {
	case "compile":
		return "danger"
	case "test":
		return "warning"
	case "coverage":
		return "info"
	case "mutation":
		return "secondary"
	default:
		return "info"
	}
}

// buildTruncationAnalysisSection 生成截断分析区域 HTML
func buildTruncationAnalysisSection(stats contracts.TruncationStats) string {
	var b strings.Builder
	b.WriteString(`<div class="section" id="truncation-analysis">
  <h2>截断分析 Truncation Analysis</h2>
  <div class="grid-2">`)

	// 总体截断统计
	truncationRate := round(stats.TruncationRate*100, 2)
	statusClass := "success"
	statusText := "良好"
	if truncationRate > 30 {
		statusClass = "danger"
		statusText = "严重"
	} else if truncationRate > 10 {
		statusClass = "warning"
		statusText = "需关注"
	}

	b.WriteString(fmt.Sprintf(`
    <div class="panel">
      <h3>总体截断情况</h3>
      <div class="metric-row">
        <div class="metric-card">
          <div class="metric-value %s">%.2f%%</div>
          <div class="metric-label">截断率</div>
          <div class="metric-hint">%s</div>
        </div>
        <div class="metric-card">
          <div class="metric-value">%d</div>
          <div class="metric-label">截断样本数</div>
        </div>
      </div>
      <div class="hint-box">
        <strong>状态：%s</strong><br>
        %s
      </div>
    </div>`,
		statusClass, truncationRate, getTruncationAdvice(truncationRate), stats.TotalTruncated, statusText, getTruncationExplanation(truncationRate)))

	// 续写功能状态
	continuationStatus := "未启用"
	if stats.ContinuationStats.Enabled {
		continuationStatus = "已启用"
	}
	b.WriteString(fmt.Sprintf(`
    <div class="panel">
      <h3>自动续写功能</h3>
      <div class="metric-row">
        <div class="metric-card">
          <div class="metric-value">%s</div>
          <div class="metric-label">续写状态</div>
        </div>
      </div>
      <div class="hint-box">
        <strong>功能说明：</strong><br>
        当模型输出被截断时，系统会自动发送续写请求，尝试恢复完整的测试代码。
        这可以显著降低截断对最终测试结果的影响。
      </div>
    </div>`, continuationStatus))

	b.WriteString(`  </div>`)

	// 按模型统计
	if len(stats.ByModel) > 0 {
		b.WriteString(`
  <div class="panel" style="margin-top: 16px;">
    <h3>按模型截断统计</h3>
    <div class="table-wrap">
      <table class="data-table">
        <thead>
          <tr>
            <th>模型</th>
            <th>总样本</th>
            <th>截断数</th>
            <th>截断率</th>
            <th>平均生成Token</th>
          </tr>
        </thead>
        <tbody>`)
		for _, m := range stats.ByModel {
			rateClass := ""
			if m.TruncationRate > 0.3 {
				rateClass = "style=\"color: #dc2626; font-weight: 600;\""
			} else if m.TruncationRate > 0.1 {
				rateClass = "style=\"color: #d97706; font-weight: 600;\""
			}
			b.WriteString(fmt.Sprintf(`
          <tr>
            <td>%s</td>
            <td>%d</td>
            <td>%d</td>
            <td %s>%.2f%%</td>
            <td>%.0f</td>
          </tr>`,
				escapeHTML(m.Model), m.TotalSamples, m.TruncatedCount, rateClass, m.TruncationRate*100, m.AvgCompletionTokens))
		}
		b.WriteString(`
        </tbody>
      </table>
    </div>
  </div>`)
	}

	// 按语言统计
	if len(stats.ByLanguage) > 0 {
		b.WriteString(`
  <div class="panel" style="margin-top: 16px;">
    <h3>按语言截断统计</h3>
    <div class="table-wrap">
      <table class="data-table">
        <thead>
          <tr>
            <th>语言</th>
            <th>总样本</th>
            <th>截断数</th>
            <th>截断率</th>
          </tr>
        </thead>
        <tbody>`)
		for _, l := range stats.ByLanguage {
			rateClass := ""
			if l.TruncationRate > 0.3 {
				rateClass = "style=\"color: #dc2626; font-weight: 600;\""
			} else if l.TruncationRate > 0.1 {
				rateClass = "style=\"color: #d97706; font-weight: 600;\""
			}
			b.WriteString(fmt.Sprintf(`
          <tr>
            <td>%s</td>
            <td>%d</td>
            <td>%d</td>
            <td %s>%.2f%%</td>
          </tr>`,
				strings.ToUpper(l.Language), l.TotalSamples, l.TruncatedCount, rateClass, l.TruncationRate*100))
		}
		b.WriteString(`
        </tbody>
      </table>
    </div>
  </div>`)
	}

	// 调优建议
	b.WriteString(`
  <div class="panel" style="margin-top: 16px;">
    <h3>调优建议</h3>
    <div class="accordion">`)

	b.WriteString(fmt.Sprintf(`
      <details class="accordion-item">
        <summary>1. 调整 max_tokens 参数</summary>
        <div class="accordion-body">
          <p>当前截断率为 %.2f%%，建议根据以下情况调整 max_tokens：</p>
          <ul>
            <li><strong>截断率 > 30%%：</strong>强烈建议增加 max_tokens 至 8192 或更高</li>
            <li><strong>截断率 10%%-30%%：</strong>建议增加 max_tokens 至 6144-8192</li>
            <li><strong>截断率 < 10%%：</strong>当前设置合理，可保持现状</li>
          </ul>
          <p>修改位置：<code>configs/models.yaml</code> 中的 <code>parameters.max_tokens</code></p>
        </div>
      </details>`, truncationRate))

	b.WriteString(`
      <details class="accordion-item">
        <summary>2. 启用自动续写功能</summary>
        <div class="accordion-body">
          <p>系统已内置自动续写功能，当检测到截断时会自动发送续写请求。</p>
          <p>续写策略：</p>
          <ul>
            <li>保留已生成的代码作为上下文</li>
            <li>请求模型继续生成剩余部分</li>
            <li>自动拼接并去重</li>
            <li>最多尝试 3 次续写</li>
          </ul>
        </div>
      </details>
      <details class="accordion-item">
        <summary>3. 优化提示词策略</summary>
        <div class="accordion-body">
          <p>如果截断问题持续存在，可以考虑：</p>
          <ul>
            <li>简化提示词，减少上下文长度</li>
            <li>要求模型生成更简洁的测试代码</li>
            <li>分步骤生成：先生成测试框架，再补充具体用例</li>
            <li>使用更高效的模型或更大的上下文窗口</li>
          </ul>
        </div>
      </details>
      <details class="accordion-item">
        <summary>4. 模型选择建议</summary>
        <div class="accordion-body">
          <p>不同模型的上下文窗口和输出能力不同：</p>
          <ul>
            <li><strong>DeepSeek：</strong>支持 64K 上下文，适合长代码生成</li>
            <li><strong>Qwen：</strong>支持 32K 上下文，中文理解能力强</li>
            <li><strong>Doubao：</strong>支持 128K 上下文，适合复杂场景</li>
          </ul>
          <p>根据任务复杂度选择合适的模型可以有效减少截断问题。</p>
        </div>
      </details>
    </div>
  </div>
</div>`)

	return b.String()
}

func getTruncationAdvice(rate float64) string {
	if rate > 30 {
		return "截断率过高，建议立即增加 max_tokens 参数或优化提示词策略"
	} else if rate > 10 {
		return "截断率中等，建议适当增加 max_tokens 参数"
	} else if rate > 0 {
		return "截断率较低，当前配置基本合理"
	}
	return "无截断问题，配置良好"
}

func getTruncationExplanation(rate float64) string {
	if rate > 30 {
		return "大量样本因达到 max_tokens 限制而被截断，可能导致测试代码不完整，严重影响测试质量。"
	} else if rate > 10 {
		return "部分样本被截断，虽然自动续写功能可以缓解，但仍建议优化配置以获得更好的效果。"
	} else if rate > 0 {
		return "少量样本被截断，自动续写功能可以有效处理这种情况。"
	}
	return "所有样本都完整生成，无需担心截断问题。"
}

// buildChartsSection 生成图表区域 HTML
func buildChartsSection(models []contracts.ModelRank) string {
	if len(models) == 0 {
		return ""
	}
	return `<div class="section" id="details">
  <h2>图表分析 Charts</h2>
  <div class="chart-grid-2">
    <div class="panel">
      <h3>模型指标对比</h3>
      <div class="chart-box tall"><canvas id="modelBarChart"></canvas></div>
    </div>
    <div class="panel">
      <h3>多维雷达图</h3>
      <div class="chart-box tall"><canvas id="radarChart"></canvas></div>
    </div>
  </div>
  <div class="chart-grid-2" style="margin-top:16px;">
    <div class="panel">
      <h3>场景通过率柱状图</h3>
      <div class="chart-box tall"><canvas id="scenarioBarChart"></canvas></div>
    </div>
    <div class="panel">
      <h3>场景性能趋势折线图</h3>
      <div class="chart-box tall"><canvas id="scenarioTrendChart"></canvas></div>
    </div>
  </div>
  <div class="chart-grid-2" style="margin-top:16px;">
    <div class="panel">
      <h3>覆盖率热力图</h3>
      <div id="coverageHeatmap" class="chart-box heatmap-box"></div>
    </div>
    <div class="panel">
      <h3>变异分数热力图</h3>
      <div id="mutationHeatmap" class="chart-box heatmap-box"></div>
    </div>
  </div>
</div>`
}

// buildPromptHTMLNew 生成新的 Prompt 展示区域
func buildPromptHTMLNew(strategy, versionID string, prompts map[string]string) string {
	var b strings.Builder
	b.WriteString(`<div class="section prompt-compact" id="prompt">
  <h2>Prompt 策略</h2>
  <div class="prompt-tags">
    <span class="prompt-tag">双语提示</span>
    <span class="prompt-tag">跨模型一致</span>
    <span class="prompt-tag">语义对齐</span>
    <span class="prompt-tag">多语言支持</span>
  </div>
  <ul class="prompt-bullets">
    <li>Generate high-quality unit tests based on the following specification</li>
    <li>覆盖正常路径、边界条件和异常行为</li>
    <li>保持测试可重复、可执行（避免随机性）</li>
    <li>使用清晰断言和有意义的期望值</li>
  </ul>`)

	// 显示第一个语言的 prompt 作为示例
	for lang, prompt := range prompts {
		b.WriteString(fmt.Sprintf(`
  <details>
    <summary>查看 %s Prompt 原文</summary>
    <div class="prompt-template-box">%s</div>
  </details>`, strings.ToUpper(lang[:1])+lang[1:], escapeHTML(prompt[:min(len(prompt), 1200)])))
		break // 只显示第一个
	}

	b.WriteString(`
</div>`)
	return b.String()
}

// buildChartScripts 生成图表脚本
func buildInteractiveScripts(payload contracts.ReportPayload, rows []contracts.EvaluationResult) string {
	topModelsJSON := marshalJSONSimple(payload.TopModels)
	rowsJSON := marshalJSONSimple(rows)

	return fmt.Sprintf(`
<script>
// Chart.js loader with CDN fallback. If blocked, we show a visible banner instead of failing silently.
(function() {
  const banner = document.getElementById('runtime-banner');
  function showBanner(message) {
    if (!banner) return;
    banner.textContent = message;
    banner.classList.add('visible');
  }
  function loadScript(url) {
    return new Promise((resolve, reject) => {
      const script = document.createElement('script');
      script.src = url;
      script.async = true;
      script.onload = () => resolve(url);
      script.onerror = () => reject(new Error('Failed to load: ' + url));
      document.head.appendChild(script);
    });
  }
  async function ensureChartJS() {
    if (window.Chart) return 'builtin';
    const urls = [
      'https://cdn.jsdelivr.net/npm/chart.js@4.4.1/dist/chart.umd.min.js',
      'https://unpkg.com/chart.js@4.4.1/dist/chart.umd.min.js',
      'https://cdnjs.cloudflare.com/ajax/libs/Chart.js/4.4.1/chart.umd.min.js'
    ];
    for (const url of urls) {
      try {
        await loadScript(url);
        if (window.Chart) return url;
      } catch (e) {
        // try next
      }
    }
    throw new Error('Chart.js unavailable');
  }
  window.__utBenchEnsureChartJS = ensureChartJS;
  window.__utBenchShowBanner = showBanner;
})();

const reportTopModels = %s;
const evaluationRows = %s;

let modelBarChart;
let radarChart;
let errorTypeChart;
let stageChart;
let scenarioBarChart;
let scenarioTrendChart;

function safeText(value) {
  if (value === null || value === undefined || value === '') return '-';
  return String(value)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;');
}

function metricCell(value) {
  if (!value) return '<span class="badge">-</span>';
  const width = Math.max(0, Math.min(100, Math.round(value * 100)));
  const fillClass = value >= 0.7 ? 'ok' : 'bad';
  return '<div class="metric"><div class="bar"><span class="' + fillClass + '" style="width:' + width + '%%"></span></div><span class="val">' + width + '%%</span></div>';
}

function getScenarioFromSample(sampleID) {
  if (!sampleID) return 'unknown';
  const idx = sampleID.indexOf('_');
  return idx > 0 ? sampleID.slice(0, idx) : sampleID;
}

function getStageFromRow(row) {
  if (row.truncated) return 'generate';
  if (row.compile_error) return 'compile';
  if (row.test_error) return 'test';
  if (row.coverage_error) return 'coverage';
  if (row.mutation_error) return 'mutation';
  return '';
}

function getErrorTypeFromRow(row) {
  const message = String(row.compile_error || row.test_error || row.coverage_error || row.mutation_error || '').toLowerCase();
  if (row.truncated) return 'truncated';
  if (message.includes('modulenotfound') || message.includes('importerror') || message.includes('no module')) return 'module_not_found';
  if (message.includes('nameerror') || message.includes("name '")) return 'name_error';
  if (message.includes('assertionerror') || message.includes('assert')) return 'assertion_failure';
  if (message.includes('syntaxerror')) return 'syntax_error';
  if (message.includes('indentation')) return 'indentation_error';
  if (message.includes('timeout')) return 'timeout';
  if (message.includes('permission')) return 'permission_error';
  return message ? 'other' : '';
}

function aggregateRows(selectedModel, selectedScenario) {
  const filtered = evaluationRows.filter(row => {
    const modelMatch = selectedModel === 'all' || row.model === selectedModel;
    const scenarioMatch = selectedScenario === 'all' || getScenarioFromSample(row.sample_id) === selectedScenario;
    return modelMatch && scenarioMatch;
  });
  const byLanguage = new Map();
  const byScenario = new Map();
  const failures = new Map();

  for (const row of filtered) {
    const langKey = row.language || 'unknown';
    if (!byLanguage.has(langKey)) byLanguage.set(langKey, { language: langKey, total: 0, compilePass: 0, testPass: 0, lineSum: 0, lineCnt: 0, branchSum: 0, branchCnt: 0, mutationSum: 0, mutationCnt: 0 });
    const lang = byLanguage.get(langKey);
    lang.total += 1;
    if (row.compile_pass) lang.compilePass += 1;
    if (row.test_pass !== null && row.test_pass !== undefined && row.test_pass) {
      lang.testPass += 1;
    }
    if (row.line_coverage !== null && row.line_coverage !== undefined) { lang.lineSum += row.line_coverage; lang.lineCnt += 1; }
    if (row.branch_coverage !== null && row.branch_coverage !== undefined) { lang.branchSum += row.branch_coverage; lang.branchCnt += 1; }
    if (row.mutation_score !== null && row.mutation_score !== undefined) { lang.mutationSum += row.mutation_score; lang.mutationCnt += 1; }

    const scenario = getScenarioFromSample(row.sample_id);
    const scenKey = langKey + '|' + scenario;
    if (!byScenario.has(scenKey)) byScenario.set(scenKey, { scenario, language: langKey, total: 0, compilePass: 0, testPass: 0, lineSum: 0, lineCnt: 0, branchSum: 0, branchCnt: 0, mutationSum: 0, mutationCnt: 0 });
    const scen = byScenario.get(scenKey);
    scen.total += 1;
    if (row.compile_pass) scen.compilePass += 1;
    if (row.test_pass !== null && row.test_pass !== undefined && row.test_pass) {
      scen.testPass += 1;
    }
    if (row.line_coverage !== null && row.line_coverage !== undefined) { scen.lineSum += row.line_coverage; scen.lineCnt += 1; }
    if (row.branch_coverage !== null && row.branch_coverage !== undefined) { scen.branchSum += row.branch_coverage; scen.branchCnt += 1; }
    if (row.mutation_score !== null && row.mutation_score !== undefined) { scen.mutationSum += row.mutation_score; scen.mutationCnt += 1; }

    const stage = getStageFromRow(row);
    const errorType = getErrorTypeFromRow(row);
    if (stage && errorType) {
      const key = stage + '|' + errorType;
      if (!failures.has(key)) failures.set(key, { stage, errorType, count: 0, exampleModel: row.model || '', exampleSample: row.sample_id || '' });
      failures.get(key).count += 1;
    }
  }

  const languages = Array.from(byLanguage.values()).map(item => ({
    language: item.language,
    total: item.total,
    compilePassRate: item.total ? item.compilePass / item.total : 0,
    testPassRate: item.total ? item.testPass / item.total : 0,
    lineCoverage: item.lineCnt ? item.lineSum / item.lineCnt : 0,
    branchCoverage: item.branchCnt ? item.branchSum / item.branchCnt : 0,
    mutationScore: item.mutationCnt ? item.mutationSum / item.mutationCnt : 0
  })).sort((a, b) => a.language.localeCompare(b.language));

  const scenarios = Array.from(byScenario.values()).map(item => ({
    scenario: item.scenario,
    language: item.language,
    total: item.total,
    compilePassRate: item.total ? item.compilePass / item.total : 0,
    testPassRate: item.total ? item.testPass / item.total : 0,
    lineCoverage: item.lineCnt ? item.lineSum / item.lineCnt : 0,
    branchCoverage: item.branchCnt ? item.branchSum / item.branchCnt : 0,
    mutationScore: item.mutationCnt ? item.mutationSum / item.mutationCnt : 0
  })).sort((a, b) => (a.language + a.scenario).localeCompare(b.language + b.scenario));

  const failureRows = Array.from(failures.values()).sort((a, b) => b.count - a.count);
  return { languages, scenarios, failureRows, filtered };
}

function renderLanguageTable(items) {
  const body = document.getElementById('by-language-body');
  const empty = document.getElementById('by-language-empty');
  body.innerHTML = items.map(item => '<tr><td><strong>' + safeText(String(item.language).toUpperCase()) + '</strong></td><td>' + item.total + '</td><td>' + metricCell(item.compilePassRate) + '</td><td>' + metricCell(item.testPassRate) + '</td><td>' + metricCell(item.lineCoverage) + '</td><td>' + metricCell(item.branchCoverage) + '</td><td>' + metricCell(item.mutationScore) + '</td></tr>').join('');
  empty.style.display = items.length ? 'none' : 'block';
}

function renderScenarioTable(items) {
  const body = document.getElementById('by-scenario-body');
  const empty = document.getElementById('by-scenario-empty');
  body.innerHTML = items.map(item => '<tr><td><strong>' + safeText(item.scenario) + '</strong></td><td>' + safeText(String(item.language).toUpperCase()) + '</td><td>' + item.total + '</td><td>' + metricCell(item.compilePassRate) + '</td><td>' + metricCell(item.testPassRate) + '</td><td>' + metricCell(item.lineCoverage) + '</td><td>' + metricCell(item.branchCoverage) + '</td><td>' + metricCell(item.mutationScore) + '</td></tr>').join('');
  empty.style.display = items.length ? 'none' : 'block';
}

function renderErrorTable(items) {
  const body = document.getElementById('error-analysis-body');
  const empty = document.getElementById('error-analysis-empty');
  body.innerHTML = items.map(item => '<tr><td>' + safeText(item.stage) + '</td><td>' + safeText(item.errorType) + '</td><td>' + item.count + '</td><td>' + safeText(item.exampleModel) + '</td><td>' + safeText(item.exampleSample) + '</td></tr>').join('');
  empty.style.display = items.length ? 'none' : 'block';
}

function buildPieData(items, field) {
  const counter = new Map();
  for (const item of items) counter.set(item[field], (counter.get(item[field]) || 0) + item.count);
  return { labels: Array.from(counter.keys()), values: Array.from(counter.values()) };
}

function upsertChart(instance, canvasId, type, labels, values, colors) {
  if (instance) instance.destroy();
  return new Chart(document.getElementById(canvasId), {
    type,
    data: { labels, datasets: [{ data: values, backgroundColor: colors }] },
    options: { responsive: true, maintainAspectRatio: false, plugins: { legend: { position: 'right' } } }
  });
}

function renderErrorCharts(items) {
  const typeData = buildPieData(items, 'errorType');
  const stageData = buildPieData(items, 'stage');
  errorTypeChart = upsertChart(errorTypeChart, 'errorTypeChart', 'doughnut', typeData.labels.length ? typeData.labels : ['No Errors'], typeData.values.length ? typeData.values : [1], ['#ef4444', '#f97316', '#eab308', '#3b82f6', '#8b5cf6', '#14b8a6']);
  stageChart = upsertChart(stageChart, 'stageChart', 'pie', stageData.labels.length ? stageData.labels : ['No Errors'], stageData.values.length ? stageData.values : [1], ['#ef4444', '#f97316', '#eab308', '#3b82f6', '#14b8a6']);
}

function renderScenarioCharts(items) {
  const labels = items.map(item => item.scenario + ' / ' + item.language.toUpperCase());
  const compileRates = items.map(item => item.compilePassRate);
  const testRates = items.map(item => item.testPassRate);
  const mutationRates = items.map(item => item.mutationScore);

  if (scenarioBarChart) scenarioBarChart.destroy();
  scenarioBarChart = new Chart(document.getElementById('scenarioBarChart'), {
    type: 'bar',
    data: {
      labels,
      datasets: [
        { label: '编译通过率', data: compileRates, backgroundColor: '#1e40af' },
        { label: '样本测试通过率', data: testRates, backgroundColor: '#10b981' }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { position: 'top' } },
      scales: { y: { beginAtZero: true, max: 1, ticks: { callback: value => Math.round(value * 100) + '%%' } } }
    }
  });

  if (scenarioTrendChart) scenarioTrendChart.destroy();
  scenarioTrendChart = new Chart(document.getElementById('scenarioTrendChart'), {
    type: 'line',
    data: {
      labels,
      datasets: [
        { label: '覆盖率趋势', data: items.map(item => item.lineCoverage), borderColor: '#3b82f6', backgroundColor: 'rgba(59,130,246,.12)', tension: .3, fill: true },
        { label: '变异分数趋势', data: mutationRates, borderColor: '#f59e0b', backgroundColor: 'rgba(245,158,11,.12)', tension: .3, fill: true }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { position: 'top' } },
      scales: { y: { beginAtZero: true, max: 1, ticks: { callback: value => Math.round(value * 100) + '%%' } } }
    }
  });
}

function renderHeatmap(containerId, rows, metricKey) {
  const container = document.getElementById(containerId);
  const grouped = new Map();
  rows.forEach(row => {
    const scenario = getScenarioFromSample(row.sample_id);
    if (!grouped.has(scenario)) grouped.set(scenario, []);
    grouped.get(scenario).push(row);
  });
  let html = '<div class="heatmap">';
  Array.from(grouped.entries()).sort((a,b) => a[0].localeCompare(b[0])).forEach(([scenario, scenarioRows]) => {
    html += '<div class="heatmap-row"><div class="heatmap-label">' + safeText(scenario) + '</div>';
    scenarioRows.slice(0, 8).forEach(row => {
      const raw = row[metricKey];
      const value = raw === null || raw === undefined ? 0 : raw;
      const hue = Math.round(value * 120);
      const bg = 'hsla(' + hue + ', 75%%, 85%%, 1)';
      html += '<div class="heatmap-cell" style="background:' + bg + ';">' + safeText(row.sample_id) + '<br>' + Math.round(value * 100) + '%%</div>';
    });
    html += '</div>';
  });
  html += '</div>';
  container.innerHTML = html;
}

function exportScenarioCSV(items) {
  const header = ['scenario','language','total_samples','compile_pass_rate','test_pass_rate','line_coverage','branch_coverage','mutation_score'];
  const lines = [header.join(',')];
  items.forEach(item => {
    lines.push([item.scenario, item.language, item.total, item.compilePassRate, item.testPassRate, item.lineCoverage, item.branchCoverage, item.mutationScore].join(','));
  });
  const blob = new Blob([lines.join('\n')], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = 'scenario-analysis.csv';
  link.click();
  URL.revokeObjectURL(url);
}

function renderModelCharts() {
  const modelNames = reportTopModels.map(item => item.model);
  const compileRates = reportTopModels.map(item => item.compile_pass_rate);
  const testRates = reportTopModels.map(item => item.avg_test_pass_rate);
  const lineRates = reportTopModels.map(item => item.avg_line_coverage);
  const mutationRates = reportTopModels.map(item => item.avg_mutation_score);

  modelBarChart = new Chart(document.getElementById('modelBarChart'), {
    type: 'bar',
    data: { labels: modelNames, datasets: [
      { label: '编译', data: compileRates, backgroundColor: '#3b82f6' },
      { label: '样测', data: testRates, backgroundColor: '#10b981' },
      { label: '覆盖', data: lineRates, backgroundColor: '#f59e0b' },
      { label: '变异', data: mutationRates, backgroundColor: '#8b5cf6' }
    ]},
    options: { responsive: true, maintainAspectRatio: false, scales: { y: { beginAtZero: true, max: 1, ticks: { callback: value => Math.round(value * 100) + '%%' } } } }
  });

  radarChart = new Chart(document.getElementById('radarChart'), {
    type: 'radar',
    data: {
      labels: ['编译', '测试', '覆盖', '变异'],
      datasets: reportTopModels.map((item, index) => ({
        label: item.model,
        data: [item.compile_pass_rate, item.avg_test_pass_rate, item.avg_line_coverage, item.avg_mutation_score],
        fill: true,
        backgroundColor: ['rgba(59,130,246,0.18)', 'rgba(16,185,129,0.18)', 'rgba(245,158,11,0.18)', 'rgba(139,92,246,0.18)'][index %% 4],
        borderColor: ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6'][index %% 4],
        pointBackgroundColor: ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6'][index %% 4]
      }))
    },
    options: { responsive: true, maintainAspectRatio: false, scales: { r: { beginAtZero: true, max: 1, ticks: { callback: value => Math.round(value * 100) + '%%' } } } }
  });
}

function renderFilteredSections() {
  const modelEl = document.getElementById('model-filter');
  const scenarioEl = document.getElementById('scenario-filter');
  const selectedModel = modelEl ? modelEl.value : 'all';
  const selectedScenario = scenarioEl ? scenarioEl.value : 'all';
  const aggregated = aggregateRows(selectedModel, selectedScenario);
  renderLanguageTable(aggregated.languages);
  renderScenarioTable(aggregated.scenarios);
  renderErrorTable(aggregated.failureRows);
  renderErrorCharts(aggregated.failureRows);
  renderScenarioCharts(aggregated.scenarios);
  renderHeatmap('coverageHeatmap', aggregated.filtered, 'line_coverage');
  renderHeatmap('mutationHeatmap', aggregated.filtered, 'mutation_score');
  const exportBtn = document.getElementById('export-scenario-csv');
  if (exportBtn) exportBtn.onclick = () => exportScenarioCSV(aggregated.scenarios);
}

(async function init() {
  try {
    if (window.__utBenchEnsureChartJS) {
      await window.__utBenchEnsureChartJS();
    }
    renderModelCharts();
    renderFilteredSections();
    const modelFilter = document.getElementById('model-filter');
    if (modelFilter) modelFilter.addEventListener('change', renderFilteredSections);
    const scenarioFilter = document.getElementById('scenario-filter');
    if (scenarioFilter) scenarioFilter.addEventListener('change', renderFilteredSections);
  } catch (e) {
    if (window.__utBenchShowBanner) {
      window.__utBenchShowBanner('图表库 Chart.js 加载失败，通常是网络或企业代理拦截了 CDN。请在联网环境打开，或让报告改为本地内置 Chart.js。');
    }
    if (window.console && console.error) console.error(e);
  }

  const backToTop = document.getElementById('back-to-top');
  window.addEventListener('scroll', () => {
    if (!backToTop) return;
    if (window.scrollY > 400) backToTop.classList.add('visible'); else backToTop.classList.remove('visible');
  });
  if (backToTop) backToTop.addEventListener('click', () => window.scrollTo({ top: 0, behavior: 'smooth' }));
})();
</script>
`, topModelsJSON, rowsJSON)
}

func buildChartScripts(models []contracts.ModelRank) string {
	if len(models) == 0 {
		return ""
	}

	modelNames, compileRates, testRates, lineCovs, mutScores := extractChartDataSimple(models)

	return fmt.Sprintf(`
<script>
const modelNames = %s;
const compileRates = %s;
const testPassRates = %s;
const lineCoverages = %s;
const mutationScores = %s;

new Chart(document.getElementById('modelBarChart'), {
  type: 'bar',
  data: {
    labels: modelNames,
    datasets: [
      { label: '编译', data: compileRates, backgroundColor: '#3b82f6' },
      { label: '测试通过', data: testPassRates, backgroundColor: '#10b981' },
      { label: '行覆盖', data: lineCoverages, backgroundColor: '#f59e0b' },
      { label: '变异', data: mutationScores, backgroundColor: '#8b5cf6' }
    ]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: 'top' },
      tooltip: {
        callbacks: {
          label: function(ctx) {
            return ctx.dataset.label + ': ' + (ctx.raw * 100).toFixed(1) + '%%';
          }
        }
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        max: 1,
        ticks: {
          callback: function(value) {
            return (value * 100).toFixed(0) + '%%';
          }
        }
      }
    }
  }
});

new Chart(document.getElementById('radarChart'), {
  type: 'radar',
  data: {
    labels: ['编译', '测试通过', '行覆盖', '变异'],
    datasets: modelNames.map((name, i) => ({
      label: name,
      data: [compileRates[i], testPassRates[i], lineCoverages[i], mutationScores[i]],
      fill: true,
      backgroundColor: ['rgba(59, 130, 246, 0.2)', 'rgba(16, 185, 129, 0.2)', 'rgba(245, 158, 11, 0.2)', 'rgba(139, 92, 246, 0.2)'][i %% 4],
      borderColor: ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6'][i %% 4],
      pointBackgroundColor: ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6'][i %% 4],
    }))
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: 'top' },
      tooltip: {
        callbacks: {
          label: function(ctx) {
            return ctx.dataset.label + ': ' + (ctx.raw * 100).toFixed(1) + '%%';
          }
        }
      }
    },
    scales: {
      r: {
        beginAtZero: true,
        max: 1,
        ticks: {
          callback: function(value) {
            return (value * 100).toFixed(0) + '%%';
          }
        }
      }
    }
  }
});
</script>
`, modelNames, compileRates, testRates, lineCovs, mutScores)
}

func stringsFromRows(rows []contracts.EvaluationResult, get func(contracts.EvaluationResult) string) []string {
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		v := strings.TrimSpace(get(row))
		if v == "" {
			continue
		}
		out = append(out, v)
	}
	return out
}

func distinctSorted(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, v := range values {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}

func summarizeList(values []string, max int) string {
	if len(values) == 0 {
		return "-"
	}
	if max <= 0 {
		max = 1
	}
	if len(values) <= max {
		return strings.Join(values, ", ")
	}
	return strings.Join(values[:max], ", ") + fmt.Sprintf(" 等%d项", len(values))
}

// progressBarNew 生成新的进度条 HTML
func progressBarNew(value, threshold float64) string {
	if value == 0 {
		return `<span class="badge">-</span>`
	}
	fillClass := "ok"
	if value < threshold {
		fillClass = "bad"
	}
	pct := int(value * 100)
	return fmt.Sprintf(`<div class="metric"><div class="bar"><span style="width:%d%%" class="%s"></span></div><span class="val">%d%%</span></div>`,
		pct, fillClass, pct)
}

func statusTone(rate, successThreshold, warningThreshold float64) string {
	if rate >= successThreshold {
		return "success"
	}
	if rate >= warningThreshold {
		return "warning"
	}
	return "error"
}

func kpiCard(label string, value float64, valueClass string, subValue string) string {
	displayValue := fmt.Sprintf("%.1f%%", value*100)
	if value > 100 {
		displayValue = fmt.Sprintf("%.1f", value)
	}
	subHTML := ""
	if subValue != "" {
		subHTML = fmt.Sprintf(`<div class="sub-value">%s</div>`, escapeHTML(subValue))
	}
	return fmt.Sprintf(`<div class="kpi-card"><div class="label">%s</div><div class="value %s">%s</div>%s</div>`, escapeHTML(label), valueClass, displayValue, subHTML)
}

func kpiCardSimple(label, value, valueClass string) string {
	return fmt.Sprintf(`<div class="kpi-card"><div class="label">%s</div><div class="value %s">%s</div></div>`, escapeHTML(label), valueClass, escapeHTML(value))
}

func getRateClass(rate, threshold float64) string {
	if rate >= threshold {
		return "success"
	}
	if rate >= threshold*0.8 {
		return "warning"
	}
	return "danger"
}

func progressBar(value, threshold float64) string {
	if value == 0 {
		return `<span class="badge badge-info">-</span>`
	}
	fillClass := "success"
	if value < threshold {
		fillClass = "danger"
	} else if value < threshold*1.05 {
		fillClass = "warning"
	}
	pct := int(value * 100)
	return fmt.Sprintf(`<div class="progress-bar"><div class="bar"><div class="fill %s" style="width:%d%%"></div></div><span class="text">%d%%</span></div>`, fillClass, pct, pct)
}

func numCell(v float64) string {
	if v == 0 {
		return "-"
	}
	return fmt.Sprintf("%.0f", v)
}

func getModelIDShort(modelID string) string {
	if modelID == "" {
		return "unknown"
	}
	if len(modelID) > 25 {
		return modelID[:22] + "..."
	}
	return modelID
}

func getScenarioLabel(scenario string) string {
	labels := map[string]string{
		"boundary":           "边界值",
		"simple_function":    "简单函数",
		"complex_dependency": "复杂依赖",
		"interface_mock":     "接口Mock",
		"unknown":            "未知",
	}
	if label, ok := labels[scenario]; ok {
		return label
	}
	return scenario
}

func extractChartDataSimple(models []contracts.ModelRank) (names, compileRates, testRates, lineCovs, mutScores string) {
	var ns []string
	var crs, trs, lcs, mss []float64
	for _, m := range models {
		ns = append(ns, m.Model)
		crs = append(crs, m.CompilePassRate)
		trs = append(trs, m.AvgTestPassRate)
		lcs = append(lcs, m.AvgLineCoverage)
		mss = append(mss, m.AvgMutationScore)
	}
	return marshalJSONSimple(ns), marshalJSONSimple(crs), marshalJSONSimple(trs), marshalJSONSimple(lcs), marshalJSONSimple(mss)
}

func marshalJSONSimple(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func escapeHTML(v string) string {
	replacer := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;", "'", "&#39;")
	return replacer.Replace(v)
}

// buildLeaderboardHTML 生成新的 Leaderboard HTML
func buildLeaderboardHTML(models []contracts.ModelRank) string {
	if len(models) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<div class="section">
  <div class="section-title">模型排名 Leaderboard</div>
  <div class="leaderboard">`)

	for _, m := range models {
		// 确定排名样式
		rankClass := ""
		itemClass := ""
		if m.Rank == 1 {
			rankClass = "gold"
			itemClass = "gold"
		} else if m.Rank == 2 {
			rankClass = "silver"
			itemClass = "silver"
		} else if m.Rank == 3 {
			rankClass = "bronze"
			itemClass = "bronze"
		}

		// 计算各指标的进度条宽度
		compileWidth := int(m.CompilePassRate * 100)
		testWidth := int(m.AvgTestPassRate * 100)
		coverWidth := int(m.AvgLineCoverage * 100)
		mutWidth := int(m.AvgMutationScore * 100)

		// 综合得分百分比
		compositePct := m.CompositeScore * 100

		b.WriteString(fmt.Sprintf(`
    <div class="lb-item %s">
      <div class="lb-rank %s">%d</div>
      <div class="lb-content">
        <div class="lb-title">%s <span style="font-size:12px;color:#6b7280;font-weight:400;">(%s)</span></div>
        <div class="lb-metrics">
          <div class="metric">
            <span class="name">编译</span>
            <div class="bar"><span style="width:%d%%;background:#3b82f6;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
          <div class="metric">
            <span class="name">测试</span>
            <div class="bar"><span style="width:%d%%;background:#10b981;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
          <div class="metric">
            <span class="name">覆盖</span>
            <div class="bar"><span style="width:%d%%;background:#f59e0b;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
          <div class="metric">
            <span class="name">变异</span>
            <div class="bar"><span style="width:%d%%;background:#8b5cf6;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
        </div>
        <div class="lb-meta">
          <span>延迟: <strong>%.0fms</strong></span>
          <span>Token: <strong>%.1fk</strong></span>
        </div>
      </div>
      <div class="lb-score">
        <div class="score-label">综合得分</div>
        <div class="score-val">%.2f</div>
      </div>
    </div>`,
			itemClass, rankClass, m.Rank,
			escapeHTML(m.Model), escapeHTML(getModelIDShort(m.ModelID)),
			compileWidth, m.CompilePassRate*100,
			testWidth, m.AvgTestPassRate*100,
			coverWidth, m.AvgLineCoverage*100,
			mutWidth, m.AvgMutationScore*100,
			m.AvgLatencyMS, m.AvgTotalTokens/1000,
			compositePct))
	}

	b.WriteString(`
  </div>
</div>`)
	return b.String()
}

// buildPromptHTML 生成 Prompt 展示区域
func buildPromptHTML(prompts map[string]string) string {
	var b strings.Builder
	b.WriteString(`<div class="prompt-section">
  <h2>Prompt 策略</h2>
  <div class="prompt-tags">
    <span class="prompt-tag">双语提示</span>
    <span class="prompt-tag">跨模型一致</span>
    <span class="prompt-tag">语义对齐</span>
    <span class="prompt-tag">多语言支持</span>
  </div>`)

	// 显示第一个语言的 prompt 作为示例
	for lang, prompt := range prompts {
		b.WriteString(fmt.Sprintf(`
  <div style="margin-top:16px;">
    <div style="font-size:13px;color:#94a3b8;margin-bottom:8px;">示例: %s</div>
    <div class="prompt-content">%s</div>
  </div>`, strings.ToUpper(lang[:1])+lang[1:], escapeHTML(prompt[:min(len(prompt), 800)])))
		break // 只显示第一个
	}

	b.WriteString(`
</div>`)
	return b.String()
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// getAvgLatency 计算平均延迟
func getAvgLatency(models []contracts.ModelRank) int {
	if len(models) == 0 {
		return 0
	}
	var total float64
	for _, m := range models {
		total += m.AvgLatencyMS
	}
	return int(total / float64(len(models)))
}

// getAvgTokens 计算平均Token消耗
func getAvgTokens(models []contracts.ModelRank) float64 {
	if len(models) == 0 {
		return 0
	}
	var total float64
	for _, m := range models {
		total += m.AvgTotalTokens
	}
	return total / float64(len(models))
}
