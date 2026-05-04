// reporter 包 - 多维度聚合分析
// 包含 buildDimensions、模型/场景聚合器、Token 统计、截断统计等
package reporter

import (
	"fmt"
	"sort"
	"strings"

	"go-ut-bench/internal/contracts"
)

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
		if !ok && agg.agentModel != "" {
			detail, ok = modelDetails[agg.agentModel]
		}
		modelID := agg.key
		provider := ""
		if ok {
			modelID = detail.ModelID
			provider = detail.Provider
		}
		byModel = append(byModel, contracts.ModelDim{
			Model:                 agg.key,
			SubjectID:             agg.subjectID,
			SubjectKind:           agg.subjectKind,
			AgentFramework:        agg.agentFramework,
			AgentModel:            agg.agentModel,
			SkillName:             agg.skillName,
			SkillVersion:          agg.skillVersion,
			ModelID:               modelID,
			Provider:              provider,
			TotalSamples:          agg.count,
			CompilePassRate:       rate(agg.compilePass, agg.count),
			AvgTestPassRate:       rate(agg.sampleTestPass, agg.count),
			AvgTestCasePassRate:   rate(agg.testPassTotal, agg.testTotal),
			AvgLineCoverage:       avg(agg.lineSum, agg.lineCnt),
			AvgBranchCoverage:     avg(agg.branchSum, agg.branchCnt),
			AvgMutationScore:      avg(agg.mutationSum, agg.mutationCnt),
			AvgLatencyMS:          avgFloat(agg.latencySum, agg.latencyCnt),
			AvgPromptTokens:       avgFloat(agg.promptTokensSum, agg.promptTokenCnt),
			AvgCompletionTokens:   avgFloat(agg.completionTokensSum, agg.completionTokenCnt),
			AvgTotalTokens:        avgFloat(agg.totalTokensSum, agg.totalTokenCnt),
			AvgAssertionDensity:   avgFloat(agg.assertionDensitySum, agg.assertionDensityCnt),
			ActualTokenSamples:    agg.actualTokenSamples,
			EstimatedTokenSamples: agg.estimatedTokenSamples,
			PartialTokenSamples:   agg.partialTokenSamples,
			MissingTokenSamples:   agg.missingTokenSamples,
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
	key                   string
	subjectID             string
	subjectKind           string
	agentFramework        string
	agentModel            string
	skillName             string
	skillVersion          string
	count                 int
	compilePass           int
	sampleTestPass        int
	testPassTotal         int
	testTotal             int
	lineSum               float64
	lineCnt               int
	branchSum             float64
	branchCnt             int
	mutationSum           float64
	mutationCnt           int
	latencySum            float64
	latencyCnt            int
	promptTokensSum       float64
	promptTokenCnt        int
	completionTokensSum   float64
	completionTokenCnt    int
	totalTokensSum        float64
	totalTokenCnt         int
	assertionDensitySum   float64
	assertionDensityCnt   int
	actualTokenSamples    int
	estimatedTokenSamples int
	partialTokenSamples   int
	missingTokenSamples   int
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
	for _, prefix := range contracts.SupportedScenarios {
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
	if a.subjectID == "" {
		a.subjectID = firstNonEmpty(row.SubjectID, row.Model)
		a.subjectKind = row.SubjectKind
		a.agentFramework = row.AgentFramework
		a.agentModel = row.AgentModel
		a.skillName = row.SkillName
		a.skillVersion = row.SkillVersion
	}
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
	switch strings.TrimSpace(row.TokenSource) {
	case "actual":
		a.actualTokenSamples++
	case "estimated":
		a.estimatedTokenSamples++
	case "partial":
		a.partialTokenSamples++
	default:
		a.missingTokenSamples++
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
		switch strings.ToLower(strings.TrimSpace(row.TokenSource)) {
		case "actual":
			stats.ActualSampleCount++
		case "estimated":
			stats.EstimatedSampleCount++
		case "partial":
			stats.PartialSampleCount++
		case "missing":
			stats.MissingSampleCount++
		default:
			if row.TotalTokens != nil || row.PromptTokens != nil || row.CompletionTokens != nil {
				stats.PartialSampleCount++
			} else {
				stats.MissingSampleCount++
			}
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
