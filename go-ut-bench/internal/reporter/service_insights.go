// reporter/service_insights.go 提供洞察生成功能
// 自动分析评测结果，生成最佳模型、弱项场景、强项场景等洞察
package reporter

import (
	"fmt"
	"sort"
	"strings"

	"go-ut-bench/internal/contracts"
)

// buildInsights 生成自动洞察结论
func buildInsights(topModels []contracts.ModelRank, dims contracts.Dimensions, summary contracts.ReportSummary, failures []contracts.FailureRow) contracts.Insights {
	insights := contracts.Insights{}

	// 最佳模型洞察
	if len(topModels) > 0 {
		best := topModels[0]
		gapToSecond := 0.0
		if len(topModels) > 1 {
			gapToSecond = best.CompositeScore - topModels[1].CompositeScore
		}
		gapPct := gapToSecond * 100
		detail := fmt.Sprintf("综合得分 %.1f%%，领先第二名 %.1f%%。编译通过率 %.1f%%，测试通过率 %.1f%%，行覆盖率 %.1f%%，变异分数 %.1f%%。",
			best.CompositeScore*100, gapPct,
			best.CompilePassRate*100, best.AvgTestPassRate*100, best.AvgLineCoverage*100, best.AvgMutationScore*100)
		insights.BestModel = contracts.InsightItem{
			Category: "best_model",
			Title:    fmt.Sprintf("%s 表现最佳", best.Model),
			Detail:   detail,
			Icon:     "winner",
			Priority: 1,
		}
	}

	// 弱项场景洞察
	scenarioWeakness := findWeakScenarios(dims.ByScenario)
	for _, ws := range scenarioWeakness {
		insights.WeakScenarios = append(insights.WeakScenarios, contracts.InsightItem{
			Category: "weak_scenario",
			Title:    fmt.Sprintf("%s 场景表现较弱", ws.scenario),
			Detail:   ws.detail,
			Icon:     "warning",
			Priority: 2,
		})
	}

	// 强项场景洞察
	scenarioStrengths := findStrongScenarios(dims.ByScenario)
	for _, ss := range scenarioStrengths {
		insights.StrongScenarios = append(insights.StrongScenarios, contracts.InsightItem{
			Category: "strong_scenario",
			Title:    fmt.Sprintf("%s 场景表现优秀", ss.scenario),
			Detail:   ss.detail,
			Icon:     "star",
			Priority: 3,
		})
	}

	// 语言差异洞察
	langGaps := findLanguageGaps(dims.ByLanguage)
	for _, lg := range langGaps {
		insights.LanguageGaps = append(insights.LanguageGaps, contracts.InsightItem{
			Category: "language_gap",
			Title:    lg.title,
			Detail:   lg.detail,
			Icon:     "compare",
			Priority: 2,
		})
	}

	// 改进建议
	if len(failures) > 0 {
		topFailure := failures[0]
		insights.Recommendations = append(insights.Recommendations, contracts.InsightItem{
			Category: "recommendation",
			Title:    fmt.Sprintf("关注 %s 类错误", topFailure.ErrorType),
			Detail:   fmt.Sprintf("该类错误出现 %d 次，主要发生在 %s 阶段。建议检查生成的测试代码是否正确处理该场景。", topFailure.Count, topFailure.Stage),
			Icon:     "fix",
			Priority: 3,
		})
	}
	if summary.CompilePassRate < 0.8 {
		insights.Recommendations = append(insights.Recommendations, contracts.InsightItem{
			Category: "recommendation",
			Title:    "提升编译通过率",
			Detail:   fmt.Sprintf("当前编译通过率 %.1f%%，建议检查模型生成的代码语法和导入是否正确。", summary.CompilePassRate*100),
			Icon:     "fix",
			Priority: 2,
		})
	}

	// 评测说明
	insights.BenchmarkNotes = append(insights.BenchmarkNotes, contracts.InsightItem{
		Category: "benchmark_note",
		Title:    "评分公式说明",
		Detail:   "综合得分 = 编译通过率×0.3 + 样本测试通过率×0.3 + 行覆盖率×0.2 + 变异分数×0.2。变异分数反映测试检测代码缺陷的能力。",
		Icon:     "info",
		Priority: 4,
	})

	return insights
}

type scenarioInsight struct {
	scenario string
	detail   string
}

func findWeakScenarios(scenarios []contracts.ScenarioDim) []scenarioInsight {
	var result []scenarioInsight
	for _, s := range scenarios {
		if s.AvgMutationScore < 0.3 || s.CompilePassRate < 0.6 {
			detail := fmt.Sprintf("语言=%s，变异分数=%.1f%%，编译通过率=%.1f%%。该场景对模型生成能力要求较高。",
				s.Language, s.AvgMutationScore*100, s.CompilePassRate*100)
			result = append(result, scenarioInsight{scenario: s.Scenario, detail: detail})
		}
	}
	return result
}

func findStrongScenarios(scenarios []contracts.ScenarioDim) []scenarioInsight {
	var result []scenarioInsight
	for _, s := range scenarios {
		if s.AvgMutationScore > 0.7 && s.CompilePassRate > 0.9 {
			detail := fmt.Sprintf("语言=%s，变异分数=%.1f%%，编译通过率=%.1f%%。模型在该场景下表现稳定。",
				s.Language, s.AvgMutationScore*100, s.CompilePassRate*100)
			result = append(result, scenarioInsight{scenario: s.Scenario, detail: detail})
		}
	}
	return result
}

type langGapInsight struct {
	title  string
	detail string
}

func findLanguageGaps(languages []contracts.LanguageDim) []langGapInsight {
	var result []langGapInsight
	if len(languages) < 2 {
		return result
	}
	var bestLang, worstLang contracts.LanguageDim
	bestCov := -1.0
	worstCov := 100.0
	for _, l := range languages {
		if l.AvgLineCoverage > bestCov {
			bestCov = l.AvgLineCoverage
			bestLang = l
		}
		if l.AvgLineCoverage < worstCov {
			worstCov = l.AvgLineCoverage
			worstLang = l
		}
	}
	if bestLang.Language != worstLang.Language && bestCov-worstCov > 0.15 {
		result = append(result, langGapInsight{
			title: fmt.Sprintf("%s 与 %s 覆盖率差距较大", strings.ToUpper(bestLang.Language), strings.ToUpper(worstLang.Language)),
			detail: fmt.Sprintf("%s 行覆盖率 %.1f%%，%s 行覆盖率 %.1f%%，差距 %.1f%%。可能原因：语言特性差异、测试框架复杂度不同。",
				strings.ToUpper(bestLang.Language), bestCov*100, strings.ToUpper(worstLang.Language), worstCov*100, (bestCov-worstCov)*100),
		})
	}
	return result
}

// buildEfficiencyStats 生成效率统计
func buildEfficiencyStats(topModels []contracts.ModelRank, rows []contracts.EvaluationResult) contracts.EfficiencyStats {
	stats := contracts.EfficiencyStats{}

	// Token效率：得分/千Token
	var tokenRows []contracts.TokenEfficiencyRow
	for _, m := range topModels {
		if m.AvgTotalTokens > 0 {
			scorePerToken := m.CompositeScore / (m.AvgTotalTokens / 1000)
			tokenRows = append(tokenRows, contracts.TokenEfficiencyRow{
				Model:          m.Model,
				ScorePerToken:  scorePerToken,
				AvgTokens:      m.AvgTotalTokens,
				CompositeScore: m.CompositeScore,
			})
		}
	}
	sort.Slice(tokenRows, func(i, j int) bool {
		return tokenRows[i].ScorePerToken > tokenRows[j].ScorePerToken
	})
	for i := range tokenRows {
		tokenRows[i].Rank = i + 1
	}
	stats.TokenEfficiency = tokenRows

	// 时间效率：得分/秒
	var timeRows []contracts.TimeEfficiencyRow
	for _, m := range topModels {
		if m.AvgLatencyMS > 0 {
			scorePerSecond := m.CompositeScore / (m.AvgLatencyMS / 1000)
			timeRows = append(timeRows, contracts.TimeEfficiencyRow{
				Model:          m.Model,
				ScorePerSecond: scorePerSecond,
				AvgLatencyMS:   m.AvgLatencyMS,
				CompositeScore: m.CompositeScore,
			})
		}
	}
	sort.Slice(timeRows, func(i, j int) bool {
		return timeRows[i].ScorePerSecond > timeRows[j].ScorePerSecond
	})
	for i := range timeRows {
		timeRows[i].Rank = i + 1
	}
	stats.TimeEfficiency = timeRows

	// 成本估算
	totalTokens := 0
	totalCost := 0.0
	pricingConfigured := false
	pricedSamples := 0
	actualTokenSamples := 0
	estimatedTokenSamples := 0
	missingTokenSamples := 0
	type costAgg struct {
		tokens           int
		costUSD          float64
		sampleCount      int
		pricedSamples    int
		actualSamples    int
		estimatedSamples int
		missingSamples   int
	}
	modelCostMap := map[string]*costAgg{}
	for _, row := range rows {
		if row.TotalTokens != nil {
			totalTokens += *row.TotalTokens
		}
		agg := modelCostMap[row.Model]
		if agg == nil {
			agg = &costAgg{}
			modelCostMap[row.Model] = agg
		}
		agg.sampleCount++
		if row.TotalTokens != nil {
			agg.tokens += *row.TotalTokens
		}
		switch strings.ToLower(strings.TrimSpace(row.TokenSource)) {
		case "actual":
			actualTokenSamples++
			agg.actualSamples++
		case "estimated":
			estimatedTokenSamples++
			agg.estimatedSamples++
		default:
			if row.TotalTokens == nil && row.PromptTokens == nil && row.CompletionTokens == nil {
				missingTokenSamples++
				agg.missingSamples++
			}
		}
		if row.EstimatedCostUSD != nil {
			totalCost += *row.EstimatedCostUSD
			agg.costUSD += *row.EstimatedCostUSD
			pricedSamples++
			agg.pricedSamples++
			pricingConfigured = true
		}
	}
	stats.CostEstimate = contracts.CostEstimate{
		TotalTokens:           totalTokens,
		EstimatedCostUSD:      totalCost,
		PricingConfigured:     pricingConfigured,
		PricedSamples:         pricedSamples,
		ActualTokenSamples:    actualTokenSamples,
		EstimatedTokenSamples: estimatedTokenSamples,
		MissingTokenSamples:   missingTokenSamples,
	}
	for model, agg := range modelCostMap {
		stats.CostEstimate.ModelCostBreakdown = append(stats.CostEstimate.ModelCostBreakdown, contracts.ModelCostRow{
			Model:                 model,
			TotalTokens:           agg.tokens,
			EstimatedCostUSD:      agg.costUSD,
			AvgCostPerSample:      agg.costUSD / float64(maxInt(1, agg.pricedSamples)),
			PricedSamples:         agg.pricedSamples,
			ActualTokenSamples:    agg.actualSamples,
			EstimatedTokenSamples: agg.estimatedSamples,
		})
	}
	sort.Slice(stats.CostEstimate.ModelCostBreakdown, func(i, j int) bool {
		return stats.CostEstimate.ModelCostBreakdown[i].EstimatedCostUSD > stats.CostEstimate.ModelCostBreakdown[j].EstimatedCostUSD
	})

	return stats
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// buildErrorDiagnosis 生成错误诊断
func buildErrorDiagnosis(rows []contracts.EvaluationResult) contracts.ErrorDiagnosis {
	diagnosis := contracts.ErrorDiagnosis{}

	compileErrorTypes := map[string]*errorTypeAgg{}
	testErrorTypes := map[string]*errorTypeAgg{}

	for _, row := range rows {
		if !row.CompilePass && row.CompileError != "" {
			cat := categorizeCompileError(row.CompileError)
			if _, ok := compileErrorTypes[cat]; !ok {
				compileErrorTypes[cat] = &errorTypeAgg{type_: cat, models: map[string]bool{}, langs: map[string]bool{}}
			}
			compileErrorTypes[cat].count++
			compileErrorTypes[cat].models[row.Model] = true
			compileErrorTypes[cat].langs[row.Language] = true
			if compileErrorTypes[cat].exampleMsg == "" {
				compileErrorTypes[cat].exampleMsg = row.CompileError
			}
		}
		if row.TestPass != nil && !*row.TestPass && row.TestError != "" {
			cat := categorizeTestError(row.TestError)
			if _, ok := testErrorTypes[cat]; !ok {
				testErrorTypes[cat] = &errorTypeAgg{type_: cat, models: map[string]bool{}, langs: map[string]bool{}}
			}
			testErrorTypes[cat].count++
			testErrorTypes[cat].models[row.Model] = true
			testErrorTypes[cat].langs[row.Language] = true
			if testErrorTypes[cat].exampleMsg == "" {
				testErrorTypes[cat].exampleMsg = row.TestError
			}
		}
	}

	totalCompileErrors := 0
	for _, agg := range compileErrorTypes {
		totalCompileErrors += agg.count
	}
	for _, agg := range compileErrorTypes {
		models := make([]string, 0, len(agg.models))
		for m := range agg.models {
			models = append(models, m)
		}
		langs := make([]string, 0, len(agg.langs))
		for l := range agg.langs {
			langs = append(langs, l)
		}
		diagnosis.CompileErrors = append(diagnosis.CompileErrors, contracts.ErrorCategory{
			Type:           agg.type_,
			Count:          agg.count,
			Rate:           float64(agg.count) / float64(maxInt(1, totalCompileErrors)),
			ExampleMsg:     shortErrText(agg.exampleMsg),
			AffectedModels: models,
			AffectedLangs:  langs,
		})
	}
	sort.Slice(diagnosis.CompileErrors, func(i, j int) bool {
		return diagnosis.CompileErrors[i].Count > diagnosis.CompileErrors[j].Count
	})

	totalTestErrors := 0
	for _, agg := range testErrorTypes {
		totalTestErrors += agg.count
	}
	for _, agg := range testErrorTypes {
		models := make([]string, 0, len(agg.models))
		for m := range agg.models {
			models = append(models, m)
		}
		langs := make([]string, 0, len(agg.langs))
		for l := range agg.langs {
			langs = append(langs, l)
		}
		diagnosis.TestErrors = append(diagnosis.TestErrors, contracts.ErrorCategory{
			Type:           agg.type_,
			Count:          agg.count,
			Rate:           float64(agg.count) / float64(maxInt(1, totalTestErrors)),
			ExampleMsg:     shortErrText(agg.exampleMsg),
			AffectedModels: models,
			AffectedLangs:  langs,
		})
	}
	sort.Slice(diagnosis.TestErrors, func(i, j int) bool {
		return diagnosis.TestErrors[i].Count > diagnosis.TestErrors[j].Count
	})

	if len(diagnosis.CompileErrors) > 0 {
		top := diagnosis.CompileErrors[0]
		diagnosis.CommonPatterns = append(diagnosis.CommonPatterns, contracts.ErrorPattern{
			Pattern: fmt.Sprintf("编译错误: %s (%d次)", top.Type, top.Count),
			Count:   top.Count,
			Advice:  getErrorAdvice(top.Type, "compile"),
		})
	}
	if len(diagnosis.TestErrors) > 0 {
		top := diagnosis.TestErrors[0]
		diagnosis.CommonPatterns = append(diagnosis.CommonPatterns, contracts.ErrorPattern{
			Pattern: fmt.Sprintf("测试错误: %s (%d次)", top.Type, top.Count),
			Count:   top.Count,
			Advice:  getErrorAdvice(top.Type, "test"),
		})
	}

	diagnosis.Recommendations = generateErrorRecommendations(diagnosis)

	return diagnosis
}

type errorTypeAgg struct {
	type_      string
	count      int
	models     map[string]bool
	langs      map[string]bool
	exampleMsg string
}

func categorizeCompileError(err string) string {
	errLower := strings.ToLower(err)
	if strings.Contains(errLower, "cannot find symbol") || strings.Contains(errLower, "undefined") || strings.Contains(errLower, "not found") || strings.Contains(errLower, "import") {
		return "import_error"
	}
	if strings.Contains(errLower, "syntax") || strings.Contains(errLower, "unexpected") || strings.Contains(errLower, "parse") {
		return "syntax_error"
	}
	if strings.Contains(errLower, "type") || strings.Contains(errLower, "cannot convert") || strings.Contains(errLower, "mismatch") {
		return "type_error"
	}
	if strings.Contains(errLower, "package") || strings.Contains(errLower, "module") || strings.Contains(errLower, "dependency") {
		return "dependency_error"
	}
	return "other_compile_error"
}

func categorizeTestError(err string) string {
	errLower := strings.ToLower(err)
	if strings.Contains(errLower, "assertion") || strings.Contains(errLower, "expected") || strings.Contains(errLower, "assert") {
		return "assertion_failure"
	}
	if strings.Contains(errLower, "timeout") || strings.Contains(errLower, "timed out") {
		return "timeout"
	}
	if strings.Contains(errLower, "exception") || strings.Contains(errLower, "error") || strings.Contains(errLower, "panic") {
		return "runtime_exception"
	}
	if strings.Contains(errLower, "fixture") || strings.Contains(errLower, "setup") || strings.Contains(errLower, "teardown") {
		return "fixture_error"
	}
	return "other_test_error"
}

func getErrorAdvice(errorType, stage string) string {
	if stage == "compile" {
		switch errorType {
		case "import_error":
			return "检查模型是否正确生成了 import/require 语句，确保引用了被测代码的包和类。"
		case "syntax_error":
			return "检查生成代码的语法正确性，可能需要调整提示词以强调语法规范。"
		case "type_error":
			return "检查模型是否理解了被测代码的类型定义，建议在提示词中包含类型信息。"
		case "dependency_error":
			return "确保测试环境的依赖包已正确配置，或提示模型使用正确的依赖导入方式。"
		default:
			return "检查具体的编译错误信息，针对性优化提示词或模型配置。"
		}
	}
	if stage == "test" {
		switch errorType {
		case "assertion_failure":
			return "检查测试断言是否符合被测代码的预期行为，可能需要调整提示词以更准确地描述期望。"
		case "timeout":
			return "测试可能包含死循环或耗时操作，建议增加超时设置或检查生成的测试逻辑。"
		case "runtime_exception":
			return "检查测试是否正确处理了异常情况，建议在提示词中强调异常处理测试。"
		default:
			return "检查具体的测试错误信息，针对性优化测试生成策略。"
		}
	}
	return "根据具体错误信息进行分析和优化。"
}

func generateErrorRecommendations(diagnosis contracts.ErrorDiagnosis) []string {
	var recs []string
	if len(diagnosis.CompileErrors) > 0 && diagnosis.CompileErrors[0].Rate > 0.5 {
		rec := fmt.Sprintf("编译错误主要集中在 %s 类型，建议优先解决。", diagnosis.CompileErrors[0].Type)
		recs = append(recs, rec)
	}
	if len(diagnosis.TestErrors) > 0 && diagnosis.TestErrors[0].Rate > 0.3 {
		rec := fmt.Sprintf("测试失败主要原因是 %s，可考虑调整测试策略。", diagnosis.TestErrors[0].Type)
		recs = append(recs, rec)
	}
	if len(recs) == 0 {
		recs = append(recs, "整体错误分布较为均匀，建议综合优化模型提示词和测试环境配置。")
	}
	return recs
}
