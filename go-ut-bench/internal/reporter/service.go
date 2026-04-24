package reporter

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"
	"go-ut-bench/internal/obs"

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
	Total      int    `json:"total"`
	Killed     int    `json:"killed"`
	Survived   int    `json:"survived"`
	NoTests    int    `json:"no_tests"`
	Timeouts   int    `json:"timeouts"`
	Skipped    int    `json:"skipped"`
	Suspicious int    `json:"suspicious"`
	Tool       string `json:"tool,omitempty"`
}

type ModelDetail struct {
	Name     string
	ModelID  string
	Provider string
}

// getPromptTemplate 返回通用的中英双语提示词模板
func getPromptTemplate() string {
	return `You are an expert unit testing engineer.
你是一名资深单元测试工程师。
Generate high-quality unit tests based on the following specification.
请基于以下规范生成高质量单元测试。

## Role & Objective（角色与目标）
- Goal: produce executable tests that match source behavior exactly.
- 目标：生成可执行且与源码行为严格一致的测试。

## Language & Framework（语言与框架）
- Language（语言）: {language}
- Test Framework（测试框架）: {framework}

## Step-by-Step Workflow（分步流程）
1) Identify callable symbols and input/output contracts from source.
2) Build a test matrix: normal, boundary, and exception paths.
3) Derive expected values only from implementation semantics.
4) Write deterministic, runnable tests with clear assertions.
5) Self-check syntax/imports/assertions before final output.

## Test Requirements（测试要求）
- Cover normal paths, boundary conditions, and error/exception behavior.
- 覆盖正常路径、边界条件和异常行为。
- Keep tests deterministic and runnable.
- 保持测试可重复、可执行（避免随机性）。
- Use clear assertions with meaningful expected values.
- 使用清晰断言和有意义的期望值。
- Coverage targets（覆盖率目标，供参考）: {coverage_targets}
- Mock requirements（Mock 要求）: {mock_requirement}

## Semantic Alignment Hard Rules（语义对齐硬约束）
- Derive expected values strictly from the given source code behavior.
- 期望值必须严格依据给定源码行为推导，不要按题型常识脑补。
- Respect exact comparison semantics in code (<, <=, >, >=, ==).
- 必须严格遵守源码比较符号语义（尤其阈值边界等于时）。
- Include explicit boundary-equality assertions when threshold/limit checks exist.
- 当存在阈值/边界判断时，必须包含"等于边界"的断言样例。
- If implementation looks counter-intuitive, still assert implementation behavior.
- 若实现与常识不一致，也必须以源码实现为准。
- If docstring/comment conflicts with implementation, trust implementation.
- 若注释/文档示例与实现冲突，以实现为准。

## Output Format（输出格式）
- Return raw test code only (no Markdown fences).
- 仅输出原始测试代码，不要 Markdown 代码块。
- Do not include explanations.
- 不要输出解释文字。

## Source Code Under Test（被测源码）
{source_code}`
}

func NewService(logger *obs.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) Generate(_ context.Context, spec contracts.RunSpec, evaluationPath string) (Output, error) {
	set, err := contracts.ReadEvaluationResultSet(evaluationPath)
	if err != nil {
		return Output{}, err
	}

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
	breakdown := buildMutationBreakdown(set.Results)
	modelInfos := buildModelInfos(dims.ByModel, modelDetails)

	payload := contracts.ReportPayload{
		SchemaVersion:    contracts.SchemaVersion,
		RunID:            spec.RunID,
		GeneratedAtUTC:   time.Now().UTC(),
		SourceEvaluation: evaluationPath,
		Summary:          summary,
		Dimensions:       dims,
		TopModels:        topModels,
		ModelInfos:       modelInfos,
		ByScenario:       dims.ByScenario,
		ByModelScenario:  dims.ByModelScenario,
		TokenStats:       tokenStats,
		Failures:         failures,
		Thresholds: contracts.Thresholds{
			CompilePassRate: 1.0,
			TestPassRate:    0.7,
			LineCoverage:    0.7,
			BranchCoverage:  0.6,
			MutationScore:   0.85,
		},
		Prompts: map[string]string{
			"python": getPromptTemplate(),
			"go":     getPromptTemplate(),
			"java":   getPromptTemplate(),
			"cpp":    getPromptTemplate(),
		},
	}

	jsonPath := filepath.Join(reportRoot, "report_summary.json")
	htmlPath := filepath.Join(reportRoot, "report.html")
	summaryJSON := map[string]any{
		"schema_version":     payload.SchemaVersion,
		"run_id":             payload.RunID,
		"generated_at_utc":   payload.GeneratedAtUTC,
		"source_evaluation":  payload.SourceEvaluation,
		"summary":            payload.Summary,
		"dimensions":         payload.Dimensions,
		"top_models":         payload.TopModels,
		"model_infos":        payload.ModelInfos,
		"by_scenario":        payload.ByScenario,
		"by_model_scenario":  payload.ByModelScenario,
		"token_stats":        payload.TokenStats,
		"failures":           payload.Failures,
		"mutation_breakdown": breakdown,
		"thresholds":         payload.Thresholds,
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
	compilePass := 0
	testPassTotal := 0
	testTotal := 0
	fallbackPass := 0
	fallbackTotal := 0
	lineSum := 0.0
	lineCnt := 0
	mutationSum := 0.0
	mutationCnt := 0
	assertDensitySum := 0.0
	assertDensityCnt := 0

	for _, row := range rows {
		if row.CompilePass {
			compilePass++
		}
		if row.TestPassCount != nil && row.TestTotalCount != nil {
			testPassTotal += *row.TestPassCount
			testTotal += *row.TestTotalCount
		} else if row.TestPass != nil {
			fallbackTotal++
			if *row.TestPass {
				fallbackPass++
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

	if fallbackTotal > 0 {
		testPassTotal += fallbackPass
		testTotal += fallbackTotal
	}

	return contracts.ReportSummary{
		TotalSamples:        total,
		CompilePassCount:    compilePass,
		CompilePassRate:     rate(compilePass, total),
		TestPassCount:       testPassTotal,
		TestPassRate:        rate(testPassTotal, testTotal),
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
		dim := contracts.ModelDim{
			Model:               agg.key,
			ModelID:             modelID,
			Provider:            provider,
			TotalSamples:        agg.count,
			CompilePassRate:     rate(agg.compilePass, agg.count),
			AvgTestPassRate:     rate(agg.testPassTotal, agg.testTotal),
			AvgLineCoverage:     avg(agg.lineSum, agg.lineCnt),
			AvgBranchCoverage:   avg(agg.branchSum, agg.branchCnt),
			AvgMutationScore:    avg(agg.mutationSum, agg.mutationCnt),
			AvgLatencyMS:        avgFloat(agg.latencySum, agg.latencyCnt),
			AvgPromptTokens:     avgFloat(agg.promptTokensSum, agg.tokenCnt),
			AvgCompletionTokens: avgFloat(agg.completionTokensSum, agg.tokenCnt),
			AvgTotalTokens:      avgFloat(agg.totalTokensSum, agg.tokenCnt),
			PassCount:           agg.passCount,
			TokensPerPass:       round(avgFloat(agg.completionTokensSum, agg.passCount), 2),
			MsPerPass:           round(avgFloat(agg.latencySum, agg.passCount), 2),
		}
		dim.CompositeScore = round(compositeScore(dim), 6)
		byModel = append(byModel, dim)
	}
	applyEfficiencyBonus(byModel)
	sort.Slice(byModel, func(i, j int) bool { return byModel[i].Model < byModel[j].Model })

	var byLanguage []contracts.LanguageDim
	for _, agg := range langMap {
		byLanguage = append(byLanguage, contracts.LanguageDim{
			Language:          agg.key,
			TotalSamples:      agg.count,
			CompilePassRate:   rate(agg.compilePass, agg.count),
			AvgTestPassRate:   rate(agg.testPassTotal, agg.testTotal),
			AvgLineCoverage:   avg(agg.lineSum, agg.lineCnt),
			AvgBranchCoverage: avg(agg.branchSum, agg.branchCnt),
			AvgMutationScore:  avg(agg.mutationSum, agg.mutationCnt),
		})
	}
	sort.Slice(byLanguage, func(i, j int) bool { return byLanguage[i].Language < byLanguage[j].Language })

	var byScenario []contracts.ScenarioDim
	for _, agg := range scenarioMap {
		byScenario = append(byScenario, contracts.ScenarioDim{
			Scenario:          agg.scenario,
			Language:          agg.language,
			TotalSamples:      agg.count,
			CompilePassRate:   rate(agg.compilePass, agg.count),
			AvgTestPassRate:   rate(agg.testPassTotal, agg.testTotal),
			AvgLineCoverage:   avg(agg.lineSum, agg.lineCnt),
			AvgBranchCoverage: avg(agg.branchSum, agg.branchCnt),
			AvgMutationScore:  avg(agg.mutationSum, agg.mutationCnt),
			AvgLatencyMS:      avgFloat(agg.latencySum, agg.latencyCnt),
			AvgTokens:         avgFloat(agg.totalTokensSum, agg.tokenCnt),
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
			AvgTestPassRate:     rate(agg.testPassTotal, agg.testTotal),
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
	passCount           int // 编译通过且至少一条测试通过的样本数，用于效率分母
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
	samplePassed := false
	if row.TestPassCount != nil && row.TestTotalCount != nil {
		a.testPassTotal += *row.TestPassCount
		a.testTotal += *row.TestTotalCount
		if row.CompilePass && *row.TestPassCount > 0 {
			samplePassed = true
		}
	} else if row.TestPass != nil {
		a.testTotal++
		if *row.TestPass {
			a.testPassTotal++
			if row.CompilePass {
				samplePassed = true
			}
		}
	}
	if samplePassed {
		a.passCount++
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

// compositeScore 计算模型综合得分。
//
// 权重设计：
//   - 编译通过率 20%：基础门槛
//   - 测试通过率 20%：按编译通过率折减（未编译视为 0 贡献）
//   - 行覆盖率 15%：同样按编译通过率折减
//   - 变异得分 35%：最能区分模型"测试是否真的有效"的指标，主导排名
//   - 效率 bonus 10%：基于效率相对排名的加分项（无参照时返回 0）
//
// 所有输入均为 [0,1] 小数；返回值亦为 [0,1]。
func compositeScore(m contracts.ModelDim) float64 {
	cp := m.CompilePassRate
	// 编译失败的样本视为测试/覆盖/变异全 0，折减后的"有效指标"
	effTest := m.AvgTestPassRate * cp
	effLine := m.AvgLineCoverage * cp
	effMut := m.AvgMutationScore * cp
	return cp*0.20 + effTest*0.20 + effLine*0.15 + effMut*0.35
}

// applyEfficiencyBonus 在一组模型间根据 tokens_per_pass / ms_per_pass 的相对排名
// 为 composite_score 叠加最多 0.10 的 bonus（最快最便宜的模型得满额，最差得 0）。
// 需要在所有 ModelDim 都填充完 composite_score 后调用。
func applyEfficiencyBonus(models []contracts.ModelDim) {
	if len(models) <= 1 {
		return
	}
	minTok, maxTok := math.Inf(1), math.Inf(-1)
	minMs, maxMs := math.Inf(1), math.Inf(-1)
	for _, m := range models {
		if m.TokensPerPass > 0 {
			if m.TokensPerPass < minTok {
				minTok = m.TokensPerPass
			}
			if m.TokensPerPass > maxTok {
				maxTok = m.TokensPerPass
			}
		}
		if m.MsPerPass > 0 {
			if m.MsPerPass < minMs {
				minMs = m.MsPerPass
			}
			if m.MsPerPass > maxMs {
				maxMs = m.MsPerPass
			}
		}
	}
	tokSpan := maxTok - minTok
	msSpan := maxMs - minMs
	for i := range models {
		var tokScore, msScore float64
		if tokSpan > 0 && models[i].TokensPerPass > 0 {
			tokScore = 1 - (models[i].TokensPerPass-minTok)/tokSpan
		}
		if msSpan > 0 && models[i].MsPerPass > 0 {
			msScore = 1 - (models[i].MsPerPass-minMs)/msSpan
		}
		bonus := (tokScore + msScore) / 2 * 0.10
		models[i].CompositeScore = round(models[i].CompositeScore+bonus, 6)
	}
}

// buildTopModels 基于 buildDimensions 已计算好的 CompositeScore 生成排名。
// 不再重复计算，确保与 by_model 中的 composite_score 一致（含 efficiency bonus）。
func buildTopModels(models []contracts.ModelDim) []contracts.ModelRank {
	sorted := append([]contracts.ModelDim(nil), models...)
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
			AvgLineCoverage:     m.AvgLineCoverage,
			AvgMutationScore:    m.AvgMutationScore,
			CompositeScore:      m.CompositeScore,
			AvgLatencyMS:        m.AvgLatencyMS,
			AvgPromptTokens:     m.AvgPromptTokens,
			AvgCompletionTokens: m.AvgCompletionTokens,
			AvgTotalTokens:      m.AvgTotalTokens,
			TokensPerPass:       m.TokensPerPass,
			MsPerPass:           m.MsPerPass,
			PassCount:           m.PassCount,
		})
	}
	return out
}

func buildFailureRows(rows []contracts.EvaluationResult) []contracts.FailureRow {
	m := map[failureKey]*failureAgg{}

	for _, row := range rows {
		if row.Truncated {
			k := failureKey{stage: "generate", errType: "truncated"}
			agg := getOrCreateFailureAgg(m, k)
			agg.count++
			agg.byModel[row.Model]++
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
			agg.byModel[row.Model]++
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
			agg.byModel[row.Model]++
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
			agg.byModel[row.Model]++
			if agg.exampleModel == "" {
				agg.exampleModel = row.Model
				agg.exampleSample = row.SampleID
				agg.exampleMessage = shortErrText(row.CoverageError)
			}
		}
		if row.MutationError != "" {
			k := failureKey{stage: "mutation", errType: "mutation_error"}
			agg := getOrCreateFailureAgg(m, k)
			agg.count++
			agg.byModel[row.Model]++
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
			ByModel:        agg.byModel,
			ExampleModel:   agg.exampleModel,
			ExampleSample:  agg.exampleSample,
			ExampleMessage: agg.exampleMessage,
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Count > out[j].Count })
	return out
}

type failureKey struct {
	stage, errType string
}

type failureAgg struct {
	key            failureKey
	count          int
	byModel        map[string]int
	exampleModel   string
	exampleSample  string
	exampleMessage string
}

func getOrCreateFailureAgg(m map[failureKey]*failureAgg, k failureKey) *failureAgg {
	if a, ok := m[k]; ok {
		return a
	}
	a := &failureAgg{key: k, byModel: map[string]int{}}
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
	for _, row := range rows {
		if row.MutationTotal != nil {
			out.Total += *row.MutationTotal
		}
		if row.MutationKilled != nil {
			out.Killed += *row.MutationKilled
		}
		if row.MutationSurvived != nil {
			out.Survived += *row.MutationSurvived
		}
		if row.MutationNoTests != nil {
			out.NoTests += *row.MutationNoTests
		}
		if row.MutationTimeouts != nil {
			out.Timeouts += *row.MutationTimeouts
		}
		if row.MutationSkipped != nil {
			out.Skipped += *row.MutationSkipped
		}
		if row.MutationSuspicious != nil {
			out.Suspicious += *row.MutationSuspicious
		}
		if row.MutationTool != "" && out.Tool == "" {
			out.Tool = row.MutationTool
		}
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

	b.WriteString(`<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>ut-bench 可视化评测报告</title>
<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.1/dist/chart.umd.min.js"></script>
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
    radial-gradient(circle at top left, rgba(15,118,110,0.08), transparent 24%),
    radial-gradient(circle at top right, rgba(217,119,6,0.10), transparent 26%),
    linear-gradient(180deg, var(--bg-2) 0%, var(--bg) 18%, #f8f5ee 100%);
  font-family: "Aptos", "Segoe UI", "PingFang SC", "Microsoft YaHei", sans-serif;
}
.wrap { max-width: 1180px; margin: 0 auto; padding: 20px 20px 40px; }

/* 紧凑Hero样式 */
.hero-compact {
  background: linear-gradient(135deg, #1e3a5f 0%, #2d5a87 50%, #3d7ab5 100%);
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
  grid-template-columns: repeat(7, 1fr);
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
  background: linear-gradient(135deg, #0f172a 0%, #102a43 100%); 
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
  background: rgba(255,255,255,0.10); 
  color: #f8fafc; 
  font-size: 12px; 
  border: 1px solid rgba(255,255,255,0.12); 
}
.prompt-compact .prompt-bullets { 
  margin: 0; 
  padding-left: 18px; 
  line-height: 1.8; 
  font-size: 13px; 
  color: rgba(241,245,249,0.92); 
}
.prompt-compact details { 
  margin-top: 14px; 
  border-radius: 14px; 
  background: rgba(2,6,23,0.35); 
  border: 1px solid rgba(148,163,184,0.18); 
  overflow: hidden; 
}
.prompt-compact details > summary { 
  list-style: none; 
  cursor: pointer; 
  padding: 10px 14px; 
  font-size: 13px; 
  color: #bae6fd; 
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
  color: #7dd3fc; 
}
.prompt-compact details[open] > summary::after { 
  content: '收起原文 ▲'; 
}
.prompt-compact .prompt-template-box { 
  margin: 0; 
  padding: 14px; 
  border-radius: 0; 
  background: rgba(2,6,23,0.52); 
  border: none; 
  border-top: 1px solid rgba(148,163,184,0.22); 
  color: #e2e8f0; 
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
@media (max-width: 900px) {
  .hero-top { flex-direction: column; }
  .hero-strip { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .grid-2 { grid-template-columns: 1fr; }
  .wrap { padding: 14px; }
  .hero h1 { font-size: 32px; }
  .chart-box { height: 260px; }
  .cards { grid-template-columns: repeat(2, 1fr); }
  .kpi-mini-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .bar { width: 90px; }
}
</style>
</head>
<body>
<div class="wrap">
`)

	// Hero Section - 紧凑版本
	b.WriteString(fmt.Sprintf(`
<div class="hero-compact">
  <div class="hero-compact-main">
    <h1>模型评测报告</h1>
    <div class="hero-compact-meta">%s · 共%d个样本</div>
  </div>
  <div class="hero-compact-stats">
    <div class="hc-stat"><div class="hc-label">编译通过率</div><div class="hc-value">%.1f%%</div></div>
    <div class="hc-stat"><div class="hc-label">测试通过率</div><div class="hc-value">%.1f%%</div></div>
    <div class="hc-stat"><div class="hc-label">覆盖率</div><div class="hc-value">%.1f%%</div></div>
    <div class="hc-stat"><div class="hc-label">变异得分率</div><div class="hc-value">%.1f%%</div></div>
    <div class="hc-stat"><div class="hc-label">平均耗时</div><div class="hc-value">%dms</div></div>
    <div class="hc-stat"><div class="hc-label">Token消耗</div><div class="hc-value">%.1fk</div></div>
    <div class="hc-stat"><div class="hc-label">断言密度</div><div class="hc-value">%.2f</div></div>
  </div>
</div>
`, payload.GeneratedAtUTC.Format("2006-01-02 15:04"),
		payload.Summary.TotalSamples,
		payload.Summary.CompilePassRate*100,
		payload.Summary.TestPassRate*100,
		payload.Summary.AvgLineCoverage*100,
		payload.Summary.AvgMutationScore*100,
		getAvgLatency(payload.TopModels),
		getAvgTokens(payload.TopModels)/1000,
		payload.Summary.AvgAssertionDensity))

	// Navigation
	b.WriteString(`
<div class="jump-nav">
  <a href="#leaderboard">模型排名</a>
  <a href="#by-language">按语言统计</a>
  <a href="#by-scenario">按场景统计</a>
  <a href="#error-analysis">错误分析</a>
  <a href="#details">图表分析</a>
  <a href="#raw-data">原始数据</a>
</div>
`)

	// Leaderboard Section - 模型排名（重点）
	b.WriteString(buildLeaderboardHTMLNew(payload.TopModels))

	// By Language Section - 按语言统计
	if len(payload.Dimensions.ByLanguage) > 0 {
		b.WriteString(buildByLanguageSection(payload.Dimensions.ByLanguage, payload.Thresholds))
	}

	// By Scenario Section - 按场景统计（新增）
	if len(payload.ByScenario) > 0 {
		b.WriteString(buildByScenarioSection(payload.ByScenario))
	}

	// Error Analysis Section - 错误分析（新增）
	if len(payload.Failures) > 0 {
		b.WriteString(buildErrorAnalysisSection(payload.Failures))
	}

	// Charts Section - 图表分析
	b.WriteString(buildChartsSection(payload.TopModels))

	// Raw Data Section - 原始数据（可展开收起）
	b.WriteString(buildRawDataSection(rows))

	// Prompt Section
	if len(payload.Prompts) > 0 {
		b.WriteString(buildPromptHTMLNew(payload.Prompts))
	}

	// Chart Scripts
	b.WriteString(buildChartScripts(payload.TopModels))

	b.WriteString(`
</div>
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

// buildByLanguageSection 生成按语言统计的 HTML
func buildByLanguageSection(languages []contracts.LanguageDim, thresholds contracts.Thresholds) string {
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
          <th>测试通过率</th>
          <th>行覆盖率</th>
          <th>分支覆盖率</th>
          <th>变异分数</th>
        </tr>
      </thead>
      <tbody>`)

	for _, l := range languages {
		b.WriteString(fmt.Sprintf(`
        <tr>
          <td><strong>%s</strong></td>
          <td>%d</td>
          <td>%s</td>
          <td>%s</td>
          <td>%s</td>
          <td>%s</td>
          <td>%s</td>
        </tr>`,
			strings.ToUpper(l.Language[:1])+l.Language[1:],
			l.TotalSamples,
			progressBarNew(l.CompilePassRate, thresholds.CompilePassRate),
			progressBarNew(l.AvgTestPassRate, thresholds.TestPassRate),
			progressBarNew(l.AvgLineCoverage, thresholds.LineCoverage),
			progressBarNew(l.AvgBranchCoverage, thresholds.BranchCoverage),
			progressBarNew(l.AvgMutationScore, thresholds.MutationScore)))
	}
	b.WriteString(`
      </tbody>
    </table>
  </div>
</div>`)
	return b.String()
}

// buildByScenarioSection 生成按场景统计的 HTML
func buildByScenarioSection(scenarios []contracts.ScenarioDim) string {
	if len(scenarios) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString(`<div class="section" id="by-scenario">
  <h2>按场景统计 By Scenario</h2>
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>场景</th>
          <th>语言</th>
          <th>样本数</th>
          <th>编译通过率</th>
          <th>测试通过率</th>
          <th>行覆盖率</th>
          <th>分支覆盖率</th>
          <th>变异得分率</th>
        </tr>
      </thead>
      <tbody>`)

	for _, s := range scenarios {
		b.WriteString(fmt.Sprintf(`
        <tr>
          <td><strong>%s</strong></td>
          <td>%s</td>
          <td>%d</td>
          <td>%s</td>
          <td>%s</td>
          <td>%s</td>
          <td>%s</td>
          <td>%s</td>
        </tr>`,
			escapeHTML(getScenarioLabel(s.Scenario)),
			escapeHTML(strings.ToUpper(s.Language)),
			s.TotalSamples,
			progressBarNew(s.CompilePassRate, 0.7),
			progressBarNew(s.AvgTestPassRate, 0.7),
			progressBarNew(s.AvgLineCoverage, 0.7),
			progressBarNew(s.AvgBranchCoverage, 0.6),
			progressBarNew(s.AvgMutationScore, 0.5)))
	}
	b.WriteString(`
      </tbody>
    </table>
  </div>
</div>`)
	return b.String()
}

// buildErrorAnalysisSection 生成错误分析部分的 HTML
func buildErrorAnalysisSection(failures []contracts.FailureRow) string {
	if len(failures) == 0 {
		return ""
	}

	// 统计错误类型分布
	errorTypes := make(map[string]int)
	stageTypes := make(map[string]int)
	for _, f := range failures {
		errorTypes[f.ErrorType] += f.Count
		stageTypes[f.Stage] += f.Count
	}

	var b strings.Builder
	b.WriteString(`<div class="section" id="error-analysis">
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
      <tbody>`)

	for _, f := range failures {
		b.WriteString(fmt.Sprintf(`
        <tr>
          <td><span class="badge badge-%s">%s</span></td>
          <td>%s</td>
          <td>%d</td>
          <td>%s</td>
          <td>%s</td>
        </tr>`,
			getStageClass(f.Stage),
			escapeHTML(f.Stage),
			escapeHTML(f.ErrorType),
			f.Count,
			escapeHTML(f.ExampleModel),
			escapeHTML(f.ExampleSample)))
	}
	b.WriteString(`
      </tbody>
    </table>
  </div>
</div>`)

	// 添加错误分布图表脚本
	b.WriteString(buildErrorChartScripts(errorTypes, stageTypes))
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

// buildChartsSection 生成图表区域 HTML
func buildChartsSection(models []contracts.ModelRank) string {
	if len(models) == 0 {
		return ""
	}
	return `<div class="section" id="details">
  <h2>图表分析 Charts</h2>
  <div class="grid-2">
    <div class="panel">
      <h3>模型指标对比</h3>
      <div class="chart-box"><canvas id="modelBarChart"></canvas></div>
    </div>
    <div class="panel">
      <h3>多维雷达图</h3>
      <div class="chart-box"><canvas id="radarChart"></canvas></div>
    </div>
  </div>
</div>`
}

// buildPromptHTMLNew 生成新的 Prompt 展示区域
func buildPromptHTMLNew(prompts map[string]string) string {
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
func buildChartScripts(models []contracts.ModelRank) string {
	if len(models) == 0 {
		return ""
	}

	modelNames, compileRates, testRates, lineCovs, mutScores := extractChartDataSimple(models)

	return fmt.Sprintf(`
<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.1/dist/chart.umd.min.js"></script>
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
