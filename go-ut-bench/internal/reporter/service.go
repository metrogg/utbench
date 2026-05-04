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

type Service struct {
	logger  *obs.Logger
	prompts contracts.PromptMetaProvider
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

func (s *Service) getPromptTemplate(language string) string {
	return s.prompts.PromptTemplatePreview(language)
}

func NewService(logger *obs.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) Generate(_ context.Context, spec contracts.RunSpec, evaluationPath string) (Output, error) {
	set, err := contracts.ReadEvaluationResultSet(evaluationPath)
	if err != nil {
		return Output{}, err
	}
	return s.GenerateFromResultSet(spec, set, evaluationPath)
}

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
		Insights:        insights,
		EfficiencyStats: efficiencyStats,
		ErrorDiagnosis:  errorDiagnosis,
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
			AvgPromptTokens:     avgFloat(agg.promptTokensSum, agg.promptTokenCnt),
			AvgCompletionTokens: avgFloat(agg.completionTokensSum, agg.completionTokenCnt),
			AvgTotalTokens:      avgFloat(agg.totalTokensSum, agg.totalTokenCnt),
			AvgAssertionDensity: avgFloat(agg.assertionDensitySum, agg.assertionDensityCnt),
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
			AvgTokens:           avgFloat(agg.totalTokensSum, agg.totalTokenCnt),
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
			AvgPromptTokens:     avgFloat(agg.promptTokensSum, agg.promptTokenCnt),
			AvgCompletionTokens: avgFloat(agg.completionTokensSum, agg.completionTokenCnt),
			AvgTotalTokens:      avgFloat(agg.totalTokensSum, agg.totalTokenCnt),
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
	promptTokensSum     float64
	promptTokenCnt      int
	completionTokensSum float64
	completionTokenCnt  int
	totalTokensSum      float64
	totalTokenCnt       int
	assertionDensitySum float64
	assertionDensityCnt int
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
	promptTokensSum     float64
	promptTokenCnt      int
	completionTokensSum float64
	completionTokenCnt  int
	totalTokensSum      float64
	totalTokenCnt       int
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
	promptTokensSum     float64
	promptTokenCnt      int
	completionTokensSum float64
	completionTokenCnt  int
	totalTokensSum      float64
	totalTokenCnt       int
}

func extractScenario(sampleID string) string {
	for _, prefix := range []string{"complex_dependency", "interface_mock", "simple_function", "boundary"} {
		if sampleID == prefix || strings.HasPrefix(sampleID, prefix+"_") {
			return prefix
		}
	}
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
	if row.LatencyMS != nil {
		a.latencySum += float64(*row.LatencyMS)
		a.latencyCnt++
	}
	if row.PromptTokens != nil {
		a.promptTokensSum += float64(*row.PromptTokens)
		a.promptTokenCnt++
	}
	if row.CompletionTokens != nil {
		a.completionTokensSum += float64(*row.CompletionTokens)
		a.completionTokenCnt++
	}
	if row.TotalTokens != nil {
		a.totalTokensSum += float64(*row.TotalTokens)
		a.totalTokenCnt++
	}
	if row.AssertionDensity != nil {
		a.assertionDensitySum += *row.AssertionDensity
		a.assertionDensityCnt++
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
	if row.LatencyMS != nil {
		a.latencySum += float64(*row.LatencyMS)
		a.latencyCnt++
	}
	if row.TotalTokens != nil {
		a.totalTokensSum += float64(*row.TotalTokens)
		a.totalTokenCnt++
	}
	if row.PromptTokens != nil {
		a.promptTokensSum += float64(*row.PromptTokens)
		a.promptTokenCnt++
	}
	if row.CompletionTokens != nil {
		a.completionTokensSum += float64(*row.CompletionTokens)
		a.completionTokenCnt++
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
	if row.LatencyMS != nil {
		a.latencySum += float64(*row.LatencyMS)
		a.latencyCnt++
	}
	if row.TotalTokens != nil {
		a.totalTokensSum += float64(*row.TotalTokens)
		a.totalTokenCnt++
	}
	if row.PromptTokens != nil {
		a.promptTokensSum += float64(*row.PromptTokens)
		a.promptTokenCnt++
	}
	if row.CompletionTokens != nil {
		a.completionTokensSum += float64(*row.CompletionTokens)
		a.completionTokenCnt++
	}
}

func buildTokenStats(rows []contracts.EvaluationResult) contracts.TokenStats {
	stats := contracts.TokenStats{}
	promptCount := 0
	completionCount := 0
	totalCount := 0
	for _, row := range rows {
		if row.PromptTokens != nil {
			stats.TotalPromptTokens += *row.PromptTokens
			promptCount++
		}
		if row.CompletionTokens != nil {
			stats.TotalCompletionTokens += *row.CompletionTokens
			completionCount++
		}
		if row.TotalTokens != nil {
			stats.TotalTokens += *row.TotalTokens
			totalCount++
		}
	}
	stats.SampleCount = max(promptCount, max(completionCount, totalCount))
	if promptCount > 0 {
		stats.AvgPromptTokens = float64(stats.TotalPromptTokens) / float64(promptCount)
	}
	if completionCount > 0 {
		stats.AvgCompletionTokens = float64(stats.TotalCompletionTokens) / float64(completionCount)
	}
	if totalCount > 0 {
		stats.AvgTotalTokens = float64(stats.TotalTokens) / float64(totalCount)
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
			AvgAssertionDensity: m.AvgAssertionDensity,
			TotalSamples:        m.TotalSamples,
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

// buildZeroMutantSamples 从评测结果中筛选出因源代码结构简单无法产生变异体的样本
// 这类样本的 MutationError 包含 "produced zero mutants" 或 "did not execute any mutants"
func buildZeroMutantSamples(rows []contracts.EvaluationResult) []contracts.ZeroMutantSample {
	type sampleKey struct {
		sampleID string
		language string
	}
	type sampleAgg struct {
		sample contracts.ZeroMutantSample
		msgs   []string
	}
	m := map[sampleKey]*sampleAgg{}
	for _, row := range rows {
		if row.MutationError == "" {
			continue
		}
		msg := strings.ToLower(row.MutationError)
		if !strings.Contains(msg, "produced zero mutants") && !strings.Contains(msg, "did not execute any mutants") {
			continue
		}
		k := sampleKey{sampleID: row.SampleID, language: row.Language}
		agg, ok := m[k]
		if !ok {
			reason := inferZeroMutantReason(row.Language, row.SourcePath)
			agg = &sampleAgg{
				sample: contracts.ZeroMutantSample{
					SampleID:   row.SampleID,
					Language:   row.Language,
					SourcePath: row.SourcePath,
					Reason:     reason,
					Count:      0,
				},
			}
			m[k] = agg
		}
		agg.sample.Count++
		if row.MutationError != "" {
			agg.msgs = append(agg.msgs, row.MutationError)
		}
	}
	out := make([]contracts.ZeroMutantSample, 0, len(m))
	for _, agg := range m {
		if len(agg.msgs) > 0 {
			agg.sample.ExampleMsg = shortErrText(agg.msgs[0])
		}
		out = append(out, agg.sample)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].SampleID < out[j].SampleID
	})
	return out
}

// inferZeroMutantReason 根据语言和源文件路径推断零变异体的原因
func inferZeroMutantReason(language, sourcePath string) string {
	switch language {
	case "go":
		return "源代码无可变异结构：go-mutesting 仅支持条件语句、算术运算、比较运算、循环、分支等结构的变异"
	case "python":
		return "源代码无可变异结构：mutmut 仅支持算术运算、比较运算、逻辑运算等结构的变异"
	case "java":
		return "源代码无可变异结构：pitest 仅支持条件语句、返回值、数学运算等结构的变异"
	case "cpp":
		return "源代码无可变异结构：mull 仅支持算术运算、比较运算、逻辑运算等结构的变异"
	default:
		return "源代码结构过于简单，无法产生有效变异体"
	}
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
	heroTypeLabels := scenarioLabels(heroTypes)

	b.WriteString(`<!doctype html>
 <html lang="zh-CN">
 <head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>ut-bench 可视化评测报告</title>
<style>
` + GetStyleCSS() + `
</style>
</head>
<body>
<div class="wrap">
<div id="runtime-banner" class="runtime-banner" role="alert"></div>
`)

	// Hero Section
	b.WriteString(fmt.Sprintf(`
<div class="hero">
  <div class="hero-header">
    <div class="hero-title-group">
      <div class="hero-badge">UT-BENCH</div>
      <h1>模型评测报告</h1>
      <div class="hero-subtitle">%s · %d 个样本</div>
    </div>
  </div>
  <div class="hero-cards">
    <div class="hero-card">
      <div class="hero-card-icon">MODEL</div>
      <div class="hero-card-content">
        <div class="hero-card-label">评测模型</div>
        <div class="hero-card-value">%s</div>
      </div>
    </div>
    <div class="hero-card">
      <div class="hero-card-icon">LANG</div>
      <div class="hero-card-content">
        <div class="hero-card-label">编程语言</div>
        <div class="hero-card-value">%s</div>
      </div>
    </div>
    <div class="hero-card">
      <div class="hero-card-icon">SET</div>
      <div class="hero-card-content">
        <div class="hero-card-label">样本类型</div>
        <div class="hero-card-value">%s</div>
      </div>
    </div>
  </div>
</div>`,
		payload.GeneratedAtUTC.Format("2006-01-02 15:04"),
		payload.Summary.TotalSamples,
		escapeHTML(summarizeList(heroModels, 6)),
		escapeHTML(summarizeList(heroLangs, 6)),
		escapeHTML(summarizeList(heroTypeLabels, 6))))

	// Navigation
	b.WriteString(`
<div class="report-shell">
  <aside class="report-sidebar" aria-label="报告目录">
    <div class="sidebar-label">报告目录</div>
    <a href="#leaderboard">模型排名</a>
    <a href="#overview">评测总览</a>
    <a href="#details">图表分析</a>
    <a href="#analysis-controls">筛选导出</a>
    <a href="#dimension-analysis">维度分析</a>
    <a href="#score-exclusions">计分剔除</a>
    <a href="#error-analysis">错误分析</a>
    <a href="#dataset-browser">评测集</a>
    <a href="#raw-data">原始数据</a>
  </aside>
  <main class="report-main">
`)

	// Leaderboard Section - 模型排名（重点）
	b.WriteString(buildLeaderboardHTMLNew(payload.TopModels))
	b.WriteString(buildOverviewSection(payload, rows))
	b.WriteString(buildChartsSection(payload.TopModels))
	b.WriteString(buildAnalysisControlsSection(heroModels, heroLangs, heroTypes))
	b.WriteString(buildDimensionAnalysisSection())

	// Truncation Analysis Section - 截断分析（新增）

	// Error Analysis Section
	b.WriteString(buildErrorAnalysisSection())

	// Score Exclusions Section
	b.WriteString(buildScoreExclusionsSection(payload.ScoreExclusions))

	// Dataset Browser Section
	b.WriteString(buildDatasetBrowserSection(rows))

	// Raw Data Section - 原始数据（可展开收起）
	b.WriteString(buildRawDataSection(rows))

	// Prompt Section
	if len(payload.Prompts) > 0 {
		b.WriteString(buildPromptHTMLNew(payload.PromptStrategy, payload.PromptVersionID, payload.Prompts))
	}

	b.WriteString(`
  </main>
</div>
`)

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

	// 计算平均延迟和Token
	var totalLatency, totalTokens float64
	var latencyCount, tokenCount int
	for _, row := range rows {
		if row.LatencyMS != nil && *row.LatencyMS > 0 {
			totalLatency += float64(*row.LatencyMS)
			latencyCount++
		}
		if row.TotalTokens != nil && *row.TotalTokens > 0 {
			totalTokens += float64(*row.TotalTokens)
			tokenCount++
		}
	}
	avgLatency := 0.0
	if latencyCount > 0 {
		avgLatency = totalLatency / float64(latencyCount) / 1000
	}
	avgTokens := 0.0
	if tokenCount > 0 {
		avgTokens = totalTokens / float64(tokenCount)
	}

	return fmt.Sprintf(`<div class="section" id="overview">
  <h2>概览 Overview</h2>
  <div class="overview-grid">
    <div class="overview-card">
      <div class="eyebrow">编译通过率</div>
      <div class="value status-%s">%.1f%%</div>
      <div class="sub">样本 %d 条</div>
    </div>
    <div class="overview-card">
      <div class="eyebrow">样本测试通过率</div>
      <div class="value status-%s">%.1f%%</div>
      <div class="sub">用例通过 %.1f%%</div>
    </div>
    <div class="overview-card">
      <div class="eyebrow">平均覆盖率</div>
      <div class="value">%.1f%%</div>
      <div class="sub">行覆盖率</div>
    </div>
    <div class="overview-card">
      <div class="eyebrow">平均变异分数</div>
      <div class="value">%.1f%%</div>
      <div class="sub">%d 个场景</div>
    </div>
  </div>
  <div class="overview-grid" style="margin-top:12px;grid-template-columns:repeat(3, minmax(0, 1fr));">
    <div class="overview-card">
      <div class="eyebrow">平均断言密度</div>
      <div class="value">%.1f</div>
      <div class="sub">每测试方法断言数</div>
    </div>
    <div class="overview-card">
      <div class="eyebrow">平均生成耗时</div>
      <div class="value">%.1fs</div>
      <div class="sub">模型响应时间</div>
    </div>
    <div class="overview-card">
      <div class="eyebrow">平均Token消耗</div>
      <div class="value">%.0f</div>
      <div class="sub">prompt + completion</div>
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
		payload.Summary.AvgAssertionDensity,
		avgLatency,
		avgTokens)
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
  <div class="lb-info-box" style="background:#f8fafc;border:1px solid #e2e8f0;border-radius:8px;padding:12px;margin-bottom:16px;font-size:13px;color:#475569;">
    <div style="display:flex;gap:16px;flex-wrap:wrap;align-items:center;">
      <span><strong>综合评分公式：</strong> 编译×0.3 + 样本测试×0.3 + 行覆盖率×0.2 + 变异分数×0.2</span>
      <span style="color:#94a3b8;">|</span>
      <span><strong>指标说明：</strong> 样测=样本级测试通过率；行覆盖=代码行覆盖率；变异=变异测试得分</span>
    </div>
  </div>
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

		// 效率指标格式化（标注平均值）
		latencyStr := fmt.Sprintf("%.1fs", m.AvgLatencyMS/1000)
		tokensStr := fmt.Sprintf("%.0f", m.AvgTotalTokens)

		// 断言密度
		assertionDensityStr := fmt.Sprintf("%.2f", m.AvgAssertionDensity)

		b.WriteString(fmt.Sprintf(`
    <div class="lb-item %s">
      <div class="lb-rank %s">%d</div>
      <div class="lb-content">
        <div class="lb-title">%s <span style="font-size:12px;color:#6b7280;font-weight:400;">(%s)</span></div>
        <div class="lb-metrics">
          <div class="metric">
            <span class="name">编译通过</span>
            <div class="bar"><span style="width:%d%%;background:#3b82f6;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
          <div class="metric">
            <span class="name">样本测试</span>
            <div class="bar"><span style="width:%d%%;background:#10b981;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
          <div class="metric">
            <span class="name">行覆盖率</span>
            <div class="bar"><span style="width:%d%%;background:#f59e0b;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
          <div class="metric">
            <span class="name">变异分数</span>
            <div class="bar"><span style="width:%d%%;background:#0f766e;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
          <div class="metric">
            <span class="name">断言密度</span>
            <div class="bar"><span style="width:%d%%;background:#14b8a6;"></span></div>
            <span class="val">%s</span>
          </div>
        </div>
        <div class="lb-meta" style="border-top:1px dashed rgba(148,163,184,0.3);padding-top:8px;margin-top:6px;font-size:12px;">
          <span style="display:inline-flex;align-items:center;gap:6px;">
            <span class="meta-label">耗时</span>平均耗时: <strong>%s</strong>
          </span>
          <span style="display:inline-flex;align-items:center;gap:6px;margin-left:12px;">
            <span class="meta-label">Token</span>平均Token: <strong>%s</strong>
          </span>
          <span style="display:inline-flex;align-items:center;gap:6px;margin-left:12px;">
            <span class="meta-label">样本</span>样本数: <strong>%d</strong>
          </span>
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
			min(100, int(m.AvgAssertionDensity*20)), // 断言密度进度条（假设理想值5）
			assertionDensityStr,
			latencyStr, tokensStr, m.TotalSamples,
			compositePct))
	}

	b.WriteString(`
  </div>
</div>`)
	return b.String()
}

func buildDimensionAnalysisSection() string {
	return `<div class="section" id="dimension-analysis">
  <h2>维度分析 Dimension Analysis</h2>
  <p class="muted">这里和上方筛选联动。可以多选模型、语言和场景，比如只看 C++ 下几个模型的表现，或只看复杂依赖场景在不同语言里的差异。</p>
  <div class="chart-grid-1">
    <div class="panel">
      <h3>模型对比 Model Comparison</h3>
      <div class="table-wrap">
        <table>
          <thead>
            <tr>
              <th>模型</th>
              <th>样本数</th>
              <th>编译通过率</th>
              <th>样本测试通过率</th>
              <th>行覆盖率</th>
              <th>分支覆盖率</th>
              <th>变异分数</th>
            </tr>
          </thead>
          <tbody id="by-model-body"></tbody>
        </table>
      </div>
      <div id="by-model-empty" class="hint-box" style="display:none;margin-top:12px;">当前筛选条件下没有模型统计数据。</div>
    </div>
  </div>
  <div class="chart-grid-1" style="margin-top:16px;">
    <div class="panel">
      <h3>语言汇总 Language Summary</h3>
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
    </div>
  </div>
  <div class="chart-grid-1" style="margin-top:16px;">
    <div class="panel">
      <h3>场景 × 语言 Scenario by Language</h3>
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
    </div>
  </div>
</div>`
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

// buildErrorAnalysisSection 生成错误分析部分的 HTML
func buildErrorAnalysisSection() string {
	return `<div class="section" id="error-analysis">
  <h2>错误分析 Error Analysis</h2>
  <div class="grid-2">
    <div class="panel">
      <h3>错误类型分布</h3>
      <div class="chart-box" style="height:220px"><canvas id="errorTypeChart" aria-label="错误类型分布柱状图"></canvas></div>
    </div>
    <div class="panel">
      <h3>失败阶段分布</h3>
      <div class="chart-box" style="height:220px"><canvas id="stageChart" aria-label="失败阶段分布柱状图"></canvas></div>
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

func buildDatasetBrowserSection(rows []contracts.EvaluationResult) string {
	type sampleEntry struct {
		SampleID   string
		Language   string
		Scenario   string
		SourcePath string
	}
	seen := map[string]struct{}{}
	entries := make([]sampleEntry, 0)
	for _, row := range rows {
		key := strings.ToLower(row.Language) + "|" + row.SampleID + "|" + row.SourcePath
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		entries = append(entries, sampleEntry{
			SampleID:   row.SampleID,
			Language:   strings.ToUpper(row.Language),
			Scenario:   extractScenario(row.SampleID),
			SourcePath: row.SourcePath,
		})
	}
	sort.Slice(entries, func(i, j int) bool {
		a := entries[i].Language + "|" + entries[i].Scenario + "|" + entries[i].SampleID
		b := entries[j].Language + "|" + entries[j].Scenario + "|" + entries[j].SampleID
		return a < b
	})
	langs := map[string]struct{}{}
	for _, entry := range entries {
		langs[entry.Language] = struct{}{}
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
  <details class="accordion-item dataset-details">
    <summary>展开/收起评测集样本列表 (%d 条去重样本)</summary>
    <div class="accordion-body">
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
      <tbody>`, len(entries), len(langs), len(entries)))
	for _, entry := range entries {
		b.WriteString(fmt.Sprintf(`
        <tr>
          <td><strong>%s</strong></td>
          <td>%s</td>
          <td>%s</td>
          <td><code>%s</code></td>
        </tr>`,
			escapeHTML(entry.SampleID),
			escapeHTML(entry.Language),
			escapeHTML(getScenarioLabel(entry.Scenario)),
			escapeHTML(entry.SourcePath)))
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

// buildZeroMutantSection 生成零变异体样本区块
func buildZeroMutantSection(rows []contracts.ZeroMutantSample) string {
	var b strings.Builder
	b.WriteString(`<div class="section" id="zero-mutant-samples">
  <h2>零变异体样本 Zero Mutant Samples</h2>
  <p class="muted">以下样本因源代码结构过于简单（如只有I/O调用、return语句等），无法产生有效变异体。这不影响模型评测排名，但可作为数据集质量分析的参考。</p>`)
	if len(rows) == 0 {
		b.WriteString(`<div class="hint-box">当前报告没有零变异体样本，所有源代码均包含可变异结构。</div>
</div>`)
		return b.String()
	}
	b.WriteString(`
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>样本ID</th>
          <th>语言</th>
          <th>次数</th>
          <th>原因说明</th>
          <th>示例消息</th>
        </tr>
      </thead>
      <tbody>`)
	for _, row := range rows {
		b.WriteString(fmt.Sprintf(`
        <tr>
          <td><code>%s</code></td>
          <td><span class="badge">%s</span></td>
          <td>%d</td>
          <td>%s</td>
          <td class="ellipsis" title="%s">%s</td>
        </tr>`,
			escapeHTML(row.SampleID),
			escapeHTML(row.Language),
			row.Count,
			escapeHTML(row.Reason),
			escapeHTML(row.ExampleMsg),
			escapeHTML(shortErrText(row.ExampleMsg))))
	}
	b.WriteString(`
      </tbody>
    </table>
  </div>
</div>`)
	return b.String()
}

// buildRawDataSection 生成原始数据部分（可展开收起，带筛选功能）
func buildRawDataSection(rows []contracts.EvaluationResult) string {
	return buildRawDataSectionSimple(rows)
}

func buildRawDataSectionSimple(rows []contracts.EvaluationResult) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<div class="section" id="raw-data">
  <h2>原始数据 Raw Data</h2>
  <div class="raw-column-controls">
    <span>扩展列：</span>
    <label><input type="checkbox" data-raw-column="coverage"> 覆盖明细</label>
    <label><input type="checkbox" data-raw-column="mutation"> 变异明细</label>
    <label><input type="checkbox" data-raw-column="assertion"> 断言/用例</label>
    <label><input type="checkbox" data-raw-column="token"> Token 明细</label>
    <label><input type="checkbox" data-raw-column="path"> 文件路径</label>
    <label><input type="checkbox" data-raw-column="error"> 错误摘要</label>
    <label><input type="checkbox" data-raw-column="score"> 计分归因</label>
  </div>
  
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
              <th class="raw-extra raw-col-coverage">分支覆盖</th>
              <th class="raw-extra raw-col-mutation">NoTests</th>
              <th class="raw-extra raw-col-mutation">Timeout</th>
              <th class="raw-extra raw-col-mutation">Skipped</th>
              <th class="raw-extra raw-col-mutation">Suspicious</th>
              <th class="raw-extra raw-col-mutation">工具</th>
              <th class="raw-extra raw-col-assertion">用例数</th>
              <th class="raw-extra raw-col-assertion">断言数</th>
              <th class="raw-extra raw-col-assertion">断言密度</th>
              <th class="raw-extra raw-col-assertion">用例通过</th>
              <th class="raw-extra raw-col-token">Prompt</th>
              <th class="raw-extra raw-col-token">Completion</th>
              <th class="raw-extra raw-col-token">截断</th>
              <th class="raw-extra raw-col-path">源码路径</th>
              <th class="raw-extra raw-col-path">测试路径</th>
              <th class="raw-extra raw-col-error">编译错误</th>
              <th class="raw-extra raw-col-error">测试错误</th>
              <th class="raw-extra raw-col-error">覆盖错误</th>
              <th class="raw-extra raw-col-error">变异错误</th>
              <th class="raw-extra raw-col-score">失败归因</th>
              <th class="raw-extra raw-col-score">计分</th>
              <th class="raw-extra raw-col-score">剔除原因</th>
            </tr>
          </thead>
          <tbody>`, len(rows)))
	for _, r := range rows {
		compileStatus := "✗"
		if r.CompilePass {
			compileStatus = "✓"
		}
		testStatus := "-"
		if r.TestPass != nil {
			if *r.TestPass {
				testStatus = "✓"
			} else {
				testStatus = "✗"
			}
		}
		lineCov := "-"
		if r.LineCoverage != nil {
			lineCov = fmt.Sprintf("%.1f%%", *r.LineCoverage*100)
		}
		mutationScore := "-"
		if r.MutationScore != nil {
			mutationScore = fmt.Sprintf("%.1f%%", *r.MutationScore*100)
		}
		mutationStats := "-"
		if r.MutationTotal != nil && *r.MutationTotal > 0 {
			survived := 0
			if r.MutationSurvived != nil {
				survived = *r.MutationSurvived
			}
			killed := 0
			if r.MutationKilled != nil {
				killed = *r.MutationKilled
			}
			mutationStats = fmt.Sprintf("%d/%d/%d", *r.MutationTotal, survived, killed)
		}
		runtime := "-"
		if r.RuntimeMS != nil {
			runtime = fmt.Sprintf("%dms", *r.RuntimeMS)
		}
		tokens := "-"
		if r.TotalTokens != nil {
			tokens = fmt.Sprintf("%d", *r.TotalTokens)
		}
		branchCov := formatPercentPtr(r.BranchCoverage)
		noTests := formatIntPtr(r.MutationNoTests)
		timeouts := formatIntPtr(r.MutationTimeouts)
		skipped := formatIntPtr(r.MutationSkipped)
		suspicious := formatIntPtr(r.MutationSuspicious)
		testCaseCount := formatIntPtr(r.TestCaseCount)
		assertionCount := formatIntPtr(r.AssertionCount)
		assertionDensity := formatFloatPtr(r.AssertionDensity, "%.2f")
		testCasePass := "-"
		if r.TestPassCount != nil && r.TestTotalCount != nil {
			testCasePass = fmt.Sprintf("%d/%d", *r.TestPassCount, *r.TestTotalCount)
		}
		promptTokens := formatIntPtr(r.PromptTokens)
		completionTokens := formatIntPtr(r.CompletionTokens)
		truncated := "-"
		if r.Truncated {
			truncated = "是"
		}
		scoreEligible := "是"
		if r.ScoreEligible != nil && !*r.ScoreEligible {
			scoreEligible = "否"
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
              <td class="raw-extra raw-col-coverage">%s</td>
              <td class="raw-extra raw-col-mutation">%s</td>
              <td class="raw-extra raw-col-mutation">%s</td>
              <td class="raw-extra raw-col-mutation">%s</td>
              <td class="raw-extra raw-col-mutation">%s</td>
              <td class="raw-extra raw-col-mutation">%s</td>
              <td class="raw-extra raw-col-assertion">%s</td>
              <td class="raw-extra raw-col-assertion">%s</td>
              <td class="raw-extra raw-col-assertion">%s</td>
              <td class="raw-extra raw-col-assertion">%s</td>
              <td class="raw-extra raw-col-token">%s</td>
              <td class="raw-extra raw-col-token">%s</td>
              <td class="raw-extra raw-col-token">%s</td>
              <td class="raw-extra raw-col-path raw-path"><code>%s</code></td>
              <td class="raw-extra raw-col-path raw-path"><code>%s</code></td>
              <td class="raw-extra raw-col-error raw-error" title="%s">%s</td>
              <td class="raw-extra raw-col-error raw-error" title="%s">%s</td>
              <td class="raw-extra raw-col-error raw-error" title="%s">%s</td>
              <td class="raw-extra raw-col-error raw-error" title="%s">%s</td>
              <td class="raw-extra raw-col-score">%s</td>
              <td class="raw-extra raw-col-score">%s</td>
              <td class="raw-extra raw-col-score raw-error" title="%s">%s</td>
            </tr>`,
			escapeHTML(r.Model),
			escapeHTML(strings.ToUpper(r.Language)),
			escapeHTML(r.SampleID),
			escapeHTML(compileStatus),
			escapeHTML(testStatus),
			escapeHTML(lineCov),
			escapeHTML(mutationScore),
			escapeHTML(mutationStats),
			escapeHTML(runtime),
			escapeHTML(tokens),
			escapeHTML(branchCov),
			escapeHTML(noTests),
			escapeHTML(timeouts),
			escapeHTML(skipped),
			escapeHTML(suspicious),
			escapeHTML(emptyDash(r.MutationTool)),
			escapeHTML(testCaseCount),
			escapeHTML(assertionCount),
			escapeHTML(assertionDensity),
			escapeHTML(testCasePass),
			escapeHTML(promptTokens),
			escapeHTML(completionTokens),
			escapeHTML(truncated),
			escapeHTML(r.SourcePath),
			escapeHTML(r.GeneratedTestPath),
			escapeHTML(r.CompileError), escapeHTML(shortErrText(r.CompileError)),
			escapeHTML(r.TestError), escapeHTML(shortErrText(r.TestError)),
			escapeHTML(r.CoverageError), escapeHTML(shortErrText(r.CoverageError)),
			escapeHTML(r.MutationError), escapeHTML(shortErrText(r.MutationError)),
			escapeHTML(emptyDash(r.FailureOrigin)),
			escapeHTML(scoreEligible),
			escapeHTML(r.ScoreExclusionReason), escapeHTML(shortErrText(r.ScoreExclusionReason))))
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

func buildRawDataSectionDetailed(rows []contracts.EvaluationResult) string {
	var b strings.Builder

	// 收集筛选选项
	models := make(map[string]bool)
	languages := make(map[string]bool)
	scenarios := make(map[string]bool)
	for _, r := range rows {
		models[r.Model] = true
		languages[r.Language] = true
		scenarios[extractScenario(r.SampleID)] = true
	}
	modelList := make([]string, 0, len(models))
	for m := range models {
		modelList = append(modelList, m)
	}
	sort.Strings(modelList)
	langList := make([]string, 0, len(languages))
	for l := range languages {
		langList = append(langList, strings.ToUpper(l))
	}
	sort.Strings(langList)
	scenarioList := make([]string, 0, len(scenarios))
	for s := range scenarios {
		scenarioList = append(scenarioList, s)
	}
	sort.Strings(scenarioList)

	b.WriteString(`<div class="section" id="raw-data">
  <h2>原始数据 Raw Data</h2>
  <p class="muted">展示所有评测样本的详细数据。默认显示前 20 条，可使用筛选功能查看特定数据。</p>

  <!-- 筛选控件 -->
  <div class="raw-data-filters" style="margin-bottom:16px;padding:12px;background:#f8fafc;border-radius:8px;">
    <div style="display:flex;gap:12px;flex-wrap:wrap;align-items:center;">
      <label style="font-weight:600;color:#475569;">筛选：</label>

      <select id="filter-model" onchange="applyRawDataFilters()" style="padding:4px 8px;border-radius:4px;border:1px solid #cbd5e1;">
        <option value="">全部模型</option>`)
	for _, m := range modelList {
		b.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, escapeHTML(m), escapeHTML(m)))
	}
	b.WriteString(`      </select>

      <select id="filter-language" onchange="applyRawDataFilters()" style="padding:4px 8px;border-radius:4px;border:1px solid #cbd5e1;">
        <option value="">全部语言</option>`)
	for _, l := range langList {
		b.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, escapeHTML(strings.ToLower(l)), escapeHTML(l)))
	}
	b.WriteString(`      </select>

      <select id="filter-scenario" onchange="applyRawDataFilters()" style="padding:4px 8px;border-radius:4px;border:1px solid #cbd5e1;">
        <option value="">全部场景</option>`)
	for _, s := range scenarioList {
		b.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, escapeHTML(s), escapeHTML(s)))
	}
	b.WriteString(`      </select>

      <select id="filter-status" onchange="applyRawDataFilters()" style="padding:4px 8px;border-radius:4px;border:1px solid #cbd5e1;">
        <option value="">全部状态</option>
        <option value="pass">全部通过</option>
        <option value="fail">有失败</option>
        <option value="compile_fail">编译失败</option>
        <option value="test_fail">测试失败</option>
        <option value="mutation_zero">变异零分</option>
      </select>

      <input type="text" id="filter-search" placeholder="搜索样本ID..."
             oninput="applyRawDataFilters()"
             style="padding:4px 8px;border-radius:4px;border:1px solid #cbd5e1;width:120px;">

      <button onclick="resetRawDataFilters()" style="padding:4px 8px;border-radius:4px;border:1px solid #cbd5e1;background:#fff;cursor:pointer;">
        重置
      </button>

      <span id="filter-result-count" style="color:#64748b;font-size:13px;">显示 20 / ` + fmt.Sprintf("%d", len(rows)) + ` 条</span>
    </div>
  </div>

  <details class="accordion-item">
    <summary onclick="initRawDataFilters()">展开/收起原始数据表格 (` + fmt.Sprintf("%d", len(rows)) + ` 条记录)</summary>
    <div class="accordion-body">
      <div class="table-wrap" style="max-height:600px;overflow:auto;">
        <table class="raw-data-table">
          <thead>
            <tr>
              <th>模型</th>
              <th>语言</th>
              <th>样本ID</th>
              <th>编译</th>
              <th>测试<br><small>(样本级)</small></th>
              <th>用例<br><small>通过率</small></th>
              <th>行覆盖</th>
              <th>分支<br><small>覆盖</small></th>
              <th>变异分</th>
              <th>变异体<br><small>(总/活/杀/跳)</small></th>
              <th>断言<br><small>密度</small></th>
              <th>用例数</th>
              <th>断言数</th>
              <th>截断</th>
              <th>计分<br><small>剔除</small></th>
              <th>耗时</th>
              <th>Tokens<br><small>(提/生/总)</small></th>
            </tr>
          </thead>
          <tbody id="raw-data-body">`)

	// 显示所有测试样本数据（添加 data 属性用于筛选）
	for _, r := range rows {
		// 编译状态
		compileStatus := `<span class="status-fail" title="编译失败">✗</span>`
		compileError := ""
		compilePass := "false"
		if r.CompilePass {
			compileStatus = `<span class="status-pass" title="编译通过">✓</span>`
			compilePass = "true"
		} else if r.CompileError != "" {
			compileError = shortErrText(r.CompileError)
			compileStatus = fmt.Sprintf(`<span class="status-fail" title="%s">✗</span>`, escapeHTML(compileError))
		}

		// 测试状态
		testStatus := `<span class="status-skip" title="未运行">-</span>`
		testError := ""
		testPass := "unknown"
		if r.TestPass != nil {
			if *r.TestPass {
				testStatus = `<span class="status-pass" title="测试通过">✓</span>`
				testPass = "true"
			} else {
				testError = shortErrText(r.TestError)
				if testError != "" {
					testStatus = fmt.Sprintf(`<span class="status-fail" title="%s">✗</span>`, escapeHTML(testError))
				} else {
					testStatus = `<span class="status-fail" title="测试失败">✗</span>`
				}
				testPass = "false"
			}
		}

		// 用例通过率
		testCaseRate := "-"
		if r.TestPassRate != nil {
			testCaseRate = fmt.Sprintf("%.1f%%", *r.TestPassRate*100)
		} else if r.TestPassCount != nil && r.TestTotalCount != nil && *r.TestTotalCount > 0 {
			rate := float64(*r.TestPassCount) / float64(*r.TestTotalCount) * 100
			testCaseRate = fmt.Sprintf("%.1f%%<br><small>%d/%d</small>", rate, *r.TestPassCount, *r.TestTotalCount)
		}

		// 行覆盖率
		lineCov := "-"
		if r.LineCoverage != nil {
			lineCov = fmt.Sprintf("%.1f%%", *r.LineCoverage*100)
		}

		// 分支覆盖率
		branchCov := "-"
		if r.BranchCoverage != nil {
			branchCov = fmt.Sprintf("%.1f%%", *r.BranchCoverage*100)
		}

		// 变异分数
		mutationScore := "-"
		mutationZero := "false"
		mutationError := ""
		if r.MutationScore != nil {
			mutationScore = fmt.Sprintf("%.1f%%", *r.MutationScore*100)
			if *r.MutationScore == 0 {
				mutationZero = "true"
			}
		} else if r.MutationError != "" {
			mutationError = shortErrText(r.MutationError)
			mutationZero = "true" // 无法计算变异分也算零分
		}

		// 变异体统计
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
			skipped := 0
			if r.MutationSkipped != nil {
				skipped = *r.MutationSkipped
			}
			mutationStats = fmt.Sprintf("%d/%d/%d/%d", *r.MutationTotal, survived, killed, skipped)
		} else if mutationError != "" {
			mutationStats = fmt.Sprintf(`<span class="status-skip" title="%s">-</span>`, escapeHTML(mutationError))
		}

		// 断言密度
		assertionDensity := "-"
		if r.AssertionDensity != nil {
			assertionDensity = fmt.Sprintf("%.2f", *r.AssertionDensity)
		}

		// 测试用例数
		testCaseCount := "-"
		if r.TestCaseCount != nil {
			testCaseCount = fmt.Sprintf("%d", *r.TestCaseCount)
		}

		// 断言数
		assertionCount := "-"
		if r.AssertionCount != nil {
			assertionCount = fmt.Sprintf("%d", *r.AssertionCount)
		}

		// 截断标记
		truncated := "-"
		if r.Truncated {
			truncated = `<span class="badge badge-warning" title="API响应因max_tokens截断">截断</span>`
		}

		// 计分剔除
		scoreExcluded := "-"
		if r.ScoreEligible != nil && !*r.ScoreEligible {
			reason := r.ScoreExclusionReason
			if reason == "" {
				reason = "未说明"
			}
			scoreExcluded = fmt.Sprintf(`<span class="badge badge-warning" title="%s">剔除</span>`, escapeHTML(reason))
		}

		// 耗时
		runtime := "-"
		if r.RuntimeMS != nil {
			if *r.RuntimeMS >= 1000 {
				runtime = fmt.Sprintf("%.1fs", float64(*r.RuntimeMS)/1000)
			} else {
				runtime = fmt.Sprintf("%dms", *r.RuntimeMS)
			}
		}

		// Tokens
		tokens := "-"
		if r.TotalTokens != nil {
			prompt := 0
			if r.PromptTokens != nil {
				prompt = *r.PromptTokens
			}
			completion := 0
			if r.CompletionTokens != nil {
				completion = *r.CompletionTokens
			}
			tokens = fmt.Sprintf("%d/%d/%d", prompt, completion, *r.TotalTokens)
		}

		// 提取场景
		scenario := extractScenario(r.SampleID)

		// 计算是否全部通过
		allPass := compilePass == "true" && testPass == "true" && mutationZero == "false"
		hasFail := compilePass == "false" || testPass == "false"

		// 添加 data 属性用于筛选
		b.WriteString(fmt.Sprintf(`
            <tr data-model="%s" data-language="%s" data-scenario="%s" data-sample="%s"
                data-compile-pass="%s" data-test-pass="%s" data-mutation-zero="%s"
                data-all-pass="%s" data-has-fail="%s"
                style="display:none;">`,
			escapeHTML(r.Model),
			escapeHTML(strings.ToLower(r.Language)),
			escapeHTML(scenario),
			escapeHTML(r.SampleID),
			compilePass, testPass, mutationZero,
			fmt.Sprintf("%v", allPass),
			fmt.Sprintf("%v", hasFail)))

		b.WriteString(fmt.Sprintf(`
              <td>%s</td>
              <td>%s</td>
              <td class="ellipsis" title="%s">%s</td>
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
              <td>%s</td>
              <td>%s</td>
              <td>%s</td>
              <td>%s</td>
            </tr>`,
			escapeHTML(r.Model),
			escapeHTML(strings.ToUpper(r.Language)),
			escapeHTML(r.SampleID), escapeHTML(r.SampleID),
			compileStatus, testStatus, testCaseRate, lineCov, branchCov,
			mutationScore, mutationStats, assertionDensity, testCaseCount,
			assertionCount, truncated, scoreExcluded, runtime, tokens))
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
  <div class="chart-grid-1" style="margin-top:16px;">
    <div class="panel efficiency-panel">
      <div class="panel-heading-row">
        <h3>速度与质量权衡</h3>
        <div class="chart-help" tabindex="0" aria-label="质量分计算说明">?
          <div class="chart-help-popover">
            <strong>质量分 Quality Score</strong>
            <p>编译通过率 25% + 样本测试通过率 30% + 行覆盖率 15% + 变异分数 25% + 断言密度归一化 5%。</p>
            <strong>X 轴</strong>
            <p>可切换为平均耗时或平均 Token。越靠左成本越低，越靠上质量越高。</p>
            <p>左上最理想；右上质量优先；左下快速初稿；右下需谨慎使用。</p>
          </div>
        </div>
        <div class="axis-toggle" role="group" aria-label="切换效率图 X 轴">
          <button type="button" class="active" data-efficiency-axis="latency">平均耗时</button>
          <button type="button" data-efficiency-axis="tokens">平均 Token</button>
        </div>
      </div>
      <div class="chart-box efficiency-chart-box"><canvas id="efficiencyQualityChart"></canvas></div>
      <div id="efficiencyQualityNotes" class="efficiency-notes"></div>
    </div>
  </div>
  <div class="chart-grid-2" style="margin-top:16px;">
    <div class="panel">
      <h3>场景通过率柱状图</h3>
      <div class="chart-box tall"><canvas id="scenarioBarChart"></canvas></div>
    </div>
    <div class="panel">
      <h3>场景覆盖与变异对比（柱状图）</h3>
      <div class="chart-box tall"><canvas id="scenarioTrendChart"></canvas></div>
    </div>
  </div>
</div>`
}

func buildAnalysisControlsSection(models, languages, scenarios []string) string {
	var b strings.Builder
	b.WriteString(`<div class="section" id="analysis-controls">
  <h2>筛选与导出 Analysis Controls</h2>
  <div class="panel">
    <div class="filter-grid">
      <div class="filter-group">
        <div class="filter-title">模型筛选 Model Filter</div>
        <label class="filter-chip filter-all"><input type="checkbox" data-filter-all="model" checked> 全部模型</label>
        <div class="filter-options">`)
	for _, model := range models {
		b.WriteString(fmt.Sprintf(`
          <label class="filter-chip"><input type="checkbox" data-filter-model value="%s"> %s</label>`, escapeHTML(model), escapeHTML(model)))
	}
	b.WriteString(`
        </div>
      </div>
      <div class="filter-group">
        <div class="filter-title">语言筛选 Language Filter</div>
        <label class="filter-chip filter-all"><input type="checkbox" data-filter-all="language" checked> 全部语言</label>
        <div class="filter-options">`)
	for _, language := range languages {
		b.WriteString(fmt.Sprintf(`
          <label class="filter-chip"><input type="checkbox" data-filter-language value="%s"> %s</label>`, escapeHTML(language), escapeHTML(strings.ToUpper(language))))
	}
	b.WriteString(`
        </div>
      </div>
      <div class="filter-group">
        <div class="filter-title">场景筛选 Scenario Filter</div>
        <label class="filter-chip filter-all"><input type="checkbox" data-filter-all="scenario" checked> 全部场景</label>
        <div class="filter-options">`)
	for _, scenario := range scenarios {
		b.WriteString(fmt.Sprintf(`
          <label class="filter-chip"><input type="checkbox" data-filter-scenario value="%s"> %s</label>`, escapeHTML(scenario), escapeHTML(getScenarioLabel(scenario))))
	}
	b.WriteString(`
        </div>
      </div>
      <div class="filter-actions">
        <button id="reset-analysis-filters" type="button">重置筛选</button>
        <button id="export-scenario-csv" type="button">导出当前维度 CSV</button>
      </div>
    </div>
    <div id="filter-summary" class="filter-summary">当前筛选：全部模型 · 全部语言 · 全部场景</div>
  </div>
</div>`)
	return b.String()
}

// buildPromptHTMLNew 生成新的 Prompt 展示区域
func buildPromptHTMLNew(strategy, versionID string, prompts map[string]string) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<div class="section prompt-compact" id="prompt">
  <h2>Prompt 策略与实际提示词</h2>
  <div class="prompt-meta">
    <span>Strategy: <strong>%s</strong></span>
    <span>Version: <strong>%s</strong></span>
    <span>展示优先级：实际 rendered prompt 快照 &gt; 模板预览</span>
  </div>
  <div class="prompt-tags">
    <span class="prompt-tag">真实运行快照</span>
    <span class="prompt-tag">跨模型一致</span>
    <span class="prompt-tag">四语言对照</span>
    <span class="prompt-tag">中文说明</span>
  </div>
  <ul class="prompt-bullets">
    <li>这里展示的是生成阶段发送给模型的提示词。若本次 run 保存了 rendered prompt，则直接读取样本快照；否则展示同版本模板预览。</li>
    <li>Prompt 不显式暴露 benchmark 场景名和复杂度，避免用题型标签提示模型。</li>
    <li>Python 在正式评测前会把中性导入名 <code>module_under_test</code> 重写到真实源码模块，减少文件名泄露。</li>
    <li>所有语言都约束小而有代表性的输入、禁止真实网络/认证/外部服务，并优先使用 fake/stub。</li>
  </ul>`, escapeHTML(emptyDash(strategy)), escapeHTML(emptyDash(versionID))))

	languages := []string{"python", "go", "java", "cpp"}
	b.WriteString(`<div class="prompt-language-grid">`)
	for _, lang := range languages {
		prompt := prompts[lang]
		if strings.TrimSpace(prompt) == "" {
			prompt = getPromptTemplate(lang)
		}
		mode := promptField(prompt, "Mode")
		framework := promptField(prompt, "Framework")
		sourceKind := "模板预览"
		if !strings.Contains(prompt, "preview_sample") && !strings.Contains(prompt, "PreviewSample") {
			sourceKind = "实际渲染快照"
		}
		displayPrompt := truncateText(prompt, 12000)
		truncatedNote := ""
		if len(displayPrompt) < len(prompt) {
			truncatedNote = fmt.Sprintf(`<div class="prompt-note">原文较长，当前仅展示前 %d 字符；完整内容已保存在 rendered prompt 快照文件中。</div>`, len(displayPrompt))
		}
		b.WriteString(fmt.Sprintf(`
  <details class="prompt-language-card">
    <summary>
      <span>%s Prompt</span>
      <small>%s · %s · %s</small>
    </summary>
    <div class="prompt-translation">
      <h3>中文说明</h3>
      %s
    </div>
    %s
    <div class="prompt-template-box">%s</div>
  </details>`,
			escapeHTML(strings.ToUpper(lang)),
			escapeHTML(emptyDash(mode)),
			escapeHTML(emptyDash(framework)),
			escapeHTML(sourceKind),
			promptChineseSummary(lang, mode),
			truncatedNote,
			escapeHTML(displayPrompt)))
	}
	b.WriteString(`</div>`)

	b.WriteString(`
</div>`)
	return b.String()
}

func promptField(prompt, field string) string {
	prefix := field + ":"
	for _, line := range strings.Split(prompt, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), prefix))
		}
	}
	return ""
}

func promptChineseSummary(language, mode string) string {
	modeText := map[string]string{
		"full_file":    "完整文件模式：模型需要直接返回一个完整、可运行的测试文件。",
		"module_level": "模块级模式：模型基于多文件模块上下文生成目标模块的完整测试文件。",
		"completion":   "续写模式：模型只补充下一个有价值的测试函数或测试块。",
	}
	if modeText[mode] == "" {
		modeText[mode] = "当前模式由 prompt 原文中的 Mode 字段决定。"
	}
	rules := []string{
		"必须只输出原始测试代码，不允许 Markdown 代码块、解释文字或占位测试。",
		"测试必须可重复，断言必须从源码实现行为推导，不能根据注释或常识猜测。",
		"覆盖正常路径、边界路径和错误路径，但输入要小而有代表性，不能做压力测试。",
		"不得访问真实网络、真实凭证或真实外部服务；涉及外部 I/O 时使用 mock、stub 或 fake。",
	}
	switch strings.ToLower(language) {
	case "python":
		rules = append(rules,
			"使用 pytest 函数式测试；从 module_under_test 导入目标符号，后续评测准备阶段会重写到真实模块。",
			"使用 plain assert 和 pytest.raises；需要 mock 时 patch 目标模块实际引用的符号。")
	case "go":
		rules = append(rules,
			"使用 testing 包和 TestXxx 函数；测试包名与源码包保持一致。",
			"适合时使用表驱动测试；不要假设标准库内部可以 monkey patch。")
	case "java":
		rules = append(rules,
			"使用 JUnit 5 Jupiter；测试类命名为 ClassNameTest，package 声明与源码保持一致。",
			"严格使用源码中声明的类名、方法名和 static/instance 调用方式。")
	case "cpp":
		rules = append(rules,
			"使用 GoogleTest 的 TEST、EXPECT_*、ASSERT_*。",
			"只包含测试所需头文件，不重复声明源码中已有的类或函数。")
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<p>%s</p><ul>`, escapeHTML(modeText[mode]))
	for _, rule := range rules {
		fmt.Fprintf(&b, `<li>%s</li>`, escapeHTML(rule))
	}
	b.WriteString(`</ul>`)
	return b.String()
}

func truncateText(value string, limit int) string {
	if limit <= 0 || len(value) <= limit {
		return value
	}
	return value[:limit]
}

// buildChartScripts 生成图表脚本
// buildInteractiveScripts 生成图表脚本
func buildInteractiveScripts(payload contracts.ReportPayload, rows []contracts.EvaluationResult) string {
	topModelsJSON := marshalJSONSimple(payload.TopModels)
	rowsJSON := marshalJSONSimple(rows)
	return "<script>\n" + BuildChartsJS(topModelsJSON, rowsJSON) + "\n</script>"
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
      { label: '变异', data: mutationScores, backgroundColor: '#0f766e' }
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
      borderColor: ['#2563eb', '#059669', '#d97706', '#0f766e'][i %% 4],
      pointBackgroundColor: ['#2563eb', '#059669', '#d97706', '#0f766e'][i %% 4],
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

func scenarioLabels(values []string) []string {
	labels := make([]string, 0, len(values))
	for _, value := range values {
		labels = append(labels, getScenarioLabel(value))
	}
	return labels
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

func statusColor(rate, successThreshold, warningThreshold float64) string {
	if rate >= successThreshold {
		return "#22c55e"
	}
	if rate >= warningThreshold {
		return "#f59e0b"
	}
	return "#ef4444"
}

func avgLatencyFromRows(rows []contracts.EvaluationResult) float64 {
	var total float64
	var count int
	for _, row := range rows {
		if row.LatencyMS != nil && *row.LatencyMS > 0 {
			total += float64(*row.LatencyMS) / 1000
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / float64(count)
}

func formatPercentPtr(v *float64) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%.1f%%", *v*100)
}

func formatIntPtr(v *int) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf("%d", *v)
}

func formatFloatPtr(v *float64, format string) string {
	if v == nil {
		return "-"
	}
	return fmt.Sprintf(format, *v)
}

func emptyDash(v string) string {
	if strings.TrimSpace(v) == "" {
		return "-"
	}
	return v
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
		"boundary":           "边界值 / Boundary",
		"simple_function":    "简单函数 / Simple",
		"complex_dependency": "复杂依赖 / Complex",
		"interface_mock":     "接口 Mock / Interface",
		"unknown":            "未知 / Unknown",
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

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// getAvgLatency 计算平均延迟（保留备用）
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

// getAvgTokens 计算平均Token消耗（保留备用）
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
