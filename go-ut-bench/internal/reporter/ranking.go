// reporter 包 - 排名、比较、错误分类、工具函数、HTML 辅助工具
// 包含模型排名、Agent/Skill 对比、变异测试分解、失败分析、
// 评分筛选、零变异体分析，以及所有共享的数值/HTML 工具函数
package reporter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"go-ut-bench/internal/contracts"
)

// mutationBreakdown 变异测试统计分布
// 包含各状态（killed/survived/no_tests等）的计数
type mutationBreakdown struct {
	Total      int                     `json:"total"`             // 变异体总数
	Killed     int                     `json:"killed"`            // 被杀死数
	Survived   int                     `json:"survived"`          // 存活数
	NoTests    int                     `json:"no_tests"`          // 无测试数
	Timeouts   int                     `json:"timeouts"`          // 超时数
	Skipped    int                     `json:"skipped"`           // 跳过数
	Suspicious int                     `json:"suspicious"`        // 可疑数
	ByTool     []mutationToolBreakdown `json:"by_tool,omitempty"` // 按工具分解
}

// mutationToolBreakdown 按变异工具的统计分布
type mutationToolBreakdown struct {
	Tool       string `json:"tool"`       // 工具名称（如 mutmut、go-mutesting、pitest）
	Total      int    `json:"total"`      // 变异体总数
	Killed     int    `json:"killed"`     // 被杀死数
	Survived   int    `json:"survived"`   // 存活数
	NoTests    int    `json:"no_tests"`   // 无测试数
	Timeouts   int    `json:"timeouts"`   // 超时数
	Skipped    int    `json:"skipped"`    // 跳过数
	Suspicious int    `json:"suspicious"` // 可疑数
}

type comparisonAgg struct {
	subjectID         string
	baselineSubjectID string
	framework         string
	model             string
	skill             string
	skillVersion      string
	language          string
	count             int
	compileDelta      float64
	testDelta         float64
	lineDelta         float64
	mutationDelta     float64
	latencyDelta      float64
	latencyCount      int
	tokensDelta       float64
	tokensCount       int
}

func buildAgentComparisons(rows []contracts.EvaluationResult) []contracts.AgentComparisonRow {
	byKey := resultLookup(rows)
	aggs := map[string]*comparisonAgg{}
	for _, row := range rows {
		if row.AgentFramework == "" || row.AgentFramework == "model_api" || row.SkillName != "no_skill" {
			continue
		}
		model := firstNonEmpty(row.AgentModel, row.Model)
		baselineKey := comparisonKey("model_api", model, "no_skill", row.Language, row.SampleID)
		baseline, ok := byKey[baselineKey]
		if !ok {
			continue
		}
		key := strings.Join([]string{row.Model, baseline.Model, row.AgentFramework, model, row.Language}, "|")
		agg := getComparisonAgg(aggs, key)
		agg.subjectID = firstNonEmpty(row.SubjectID, row.Model)
		agg.baselineSubjectID = firstNonEmpty(baseline.SubjectID, baseline.Model)
		agg.framework = row.AgentFramework
		agg.model = model
		agg.skill = "no_skill"
		agg.language = row.Language
		addComparisonDelta(agg, row, baseline)
	}
	return agentComparisonRows(aggs)
}

func buildSkillUplifts(rows []contracts.EvaluationResult) []contracts.SkillUpliftRow {
	byKey := resultLookup(rows)
	aggs := map[string]*comparisonAgg{}
	for _, row := range rows {
		skill := firstNonEmpty(row.SkillName, "no_skill")
		if skill == "no_skill" {
			continue
		}
		framework := firstNonEmpty(row.AgentFramework, "model_api")
		model := firstNonEmpty(row.AgentModel, row.Model)
		baselineKey := comparisonKey(framework, model, "no_skill", row.Language, row.SampleID)
		baseline, ok := byKey[baselineKey]
		if !ok {
			continue
		}
		key := strings.Join([]string{row.Model, baseline.Model, framework, model, skill, row.Language}, "|")
		agg := getComparisonAgg(aggs, key)
		agg.subjectID = firstNonEmpty(row.SubjectID, row.Model)
		agg.baselineSubjectID = firstNonEmpty(baseline.SubjectID, baseline.Model)
		agg.framework = framework
		agg.model = model
		agg.skill = skill
		agg.skillVersion = row.SkillVersion
		agg.language = row.Language
		addComparisonDelta(agg, row, baseline)
	}
	return skillUpliftRows(aggs)
}

func resultLookup(rows []contracts.EvaluationResult) map[string]contracts.EvaluationResult {
	out := map[string]contracts.EvaluationResult{}
	for _, row := range rows {
		framework := firstNonEmpty(row.AgentFramework, inferFramework(row.Model))
		model := firstNonEmpty(row.AgentModel, inferAgentModel(row.Model))
		skill := firstNonEmpty(row.SkillName, inferSkill(row.Model))
		out[comparisonKey(framework, model, skill, row.Language, row.SampleID)] = row
	}
	return out
}

func comparisonKey(framework, model, skill, language, sampleID string) string {
	return strings.Join([]string{framework, model, skill, language, sampleID}, "|")
}

func getComparisonAgg(aggs map[string]*comparisonAgg, key string) *comparisonAgg {
	if agg, ok := aggs[key]; ok {
		return agg
	}
	agg := &comparisonAgg{}
	aggs[key] = agg
	return agg
}

func addComparisonDelta(agg *comparisonAgg, row, baseline contracts.EvaluationResult) {
	agg.count++
	agg.compileDelta += boolMetric(row.CompilePass) - boolMetric(baseline.CompilePass)
	agg.testDelta += ptrBoolMetric(row.TestPass) - ptrBoolMetric(baseline.TestPass)
	agg.lineDelta += ptrFloatMetric(row.LineCoverage) - ptrFloatMetric(baseline.LineCoverage)
	agg.mutationDelta += ptrFloatMetric(row.MutationScore) - ptrFloatMetric(baseline.MutationScore)
	if row.LatencyMS != nil && baseline.LatencyMS != nil {
		agg.latencyDelta += float64(*row.LatencyMS - *baseline.LatencyMS)
		agg.latencyCount++
	}
	if row.TotalTokens != nil && baseline.TotalTokens != nil {
		agg.tokensDelta += float64(*row.TotalTokens - *baseline.TotalTokens)
		agg.tokensCount++
	}
}

func agentComparisonRows(aggs map[string]*comparisonAgg) []contracts.AgentComparisonRow {
	keys := make([]string, 0, len(aggs))
	for key := range aggs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]contracts.AgentComparisonRow, 0, len(keys))
	for _, key := range keys {
		agg := aggs[key]
		out = append(out, contracts.AgentComparisonRow{
			SubjectID:          agg.subjectID,
			BaselineSubjectID:  agg.baselineSubjectID,
			Framework:          agg.framework,
			Model:              agg.model,
			Skill:              agg.skill,
			Language:           agg.language,
			SampleCount:        agg.count,
			CompilePassDelta:   avg(agg.compileDelta, agg.count),
			TestPassDelta:      avg(agg.testDelta, agg.count),
			LineCoverageDelta:  avg(agg.lineDelta, agg.count),
			MutationScoreDelta: avg(agg.mutationDelta, agg.count),
			LatencyMSDelta:     avg(agg.latencyDelta, agg.latencyCount),
			TotalTokensDelta:   avg(agg.tokensDelta, agg.tokensCount),
		})
	}
	return out
}

func skillUpliftRows(aggs map[string]*comparisonAgg) []contracts.SkillUpliftRow {
	keys := make([]string, 0, len(aggs))
	for key := range aggs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]contracts.SkillUpliftRow, 0, len(keys))
	for _, key := range keys {
		agg := aggs[key]
		out = append(out, contracts.SkillUpliftRow{
			SubjectID:          agg.subjectID,
			BaselineSubjectID:  agg.baselineSubjectID,
			Framework:          agg.framework,
			Model:              agg.model,
			Skill:              agg.skill,
			SkillVersion:       agg.skillVersion,
			Language:           agg.language,
			SampleCount:        agg.count,
			CompilePassDelta:   avg(agg.compileDelta, agg.count),
			TestPassDelta:      avg(agg.testDelta, agg.count),
			LineCoverageDelta:  avg(agg.lineDelta, agg.count),
			MutationScoreDelta: avg(agg.mutationDelta, agg.count),
			LatencyMSDelta:     avg(agg.latencyDelta, agg.latencyCount),
			TotalTokensDelta:   avg(agg.tokensDelta, agg.tokensCount),
		})
	}
	return out
}

func buildTopModels(models []contracts.ModelDim) []contracts.ModelRank {
	var sorted []contracts.ModelDim
	for _, m := range models {
		// 计算综合得分：编译30% + 测试30% + 覆盖20% + 变异20%
		composite := m.CompilePassRate*contracts.DefaultWeights.Compile + m.AvgTestPassRate*contracts.DefaultWeights.Test + m.AvgLineCoverage*contracts.DefaultWeights.Coverage + m.AvgMutationScore*contracts.DefaultWeights.Mutation
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
			Rank:                  i + 1,
			Model:                 m.Model,
			SubjectID:             m.SubjectID,
			SubjectKind:           m.SubjectKind,
			AgentFramework:        m.AgentFramework,
			AgentModel:            m.AgentModel,
			SkillName:             m.SkillName,
			SkillVersion:          m.SkillVersion,
			ModelID:               m.ModelID,
			Provider:              m.Provider,
			CompilePassRate:       m.CompilePassRate,
			AvgTestPassRate:       m.AvgTestPassRate,
			AvgTestCasePassRate:   m.AvgTestCasePassRate,
			AvgLineCoverage:       m.AvgLineCoverage,
			AvgMutationScore:      m.AvgMutationScore,
			CompositeScore:        m.CompositeScore,
			AvgLatencyMS:          m.AvgLatencyMS,
			AvgPromptTokens:       m.AvgPromptTokens,
			AvgCompletionTokens:   m.AvgCompletionTokens,
			AvgTotalTokens:        m.AvgTotalTokens,
			AvgAssertionDensity:   m.AvgAssertionDensity,
			TotalSamples:          m.TotalSamples,
			ActualTokenSamples:    m.ActualTokenSamples,
			EstimatedTokenSamples: m.EstimatedTokenSamples,
			PartialTokenSamples:   m.PartialTokenSamples,
			MissingTokenSamples:   m.MissingTokenSamples,
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

func boolMetric(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

func ptrBoolMetric(v *bool) float64 {
	if v != nil && *v {
		return 1
	}
	return 0
}

func ptrFloatMetric(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func inferFramework(subjectID string) string {
	parts := strings.Split(subjectID, "__")
	if len(parts) >= 1 && parts[0] != "" {
		return parts[0]
	}
	return "model_api"
}

func inferAgentModel(subjectID string) string {
	parts := strings.Split(subjectID, "__")
	if len(parts) >= 2 && parts[1] != "" {
		return parts[1]
	}
	return subjectID
}

func inferSkill(subjectID string) string {
	parts := strings.Split(subjectID, "__")
	if len(parts) >= 3 && parts[2] != "" {
		return parts[2]
	}
	return "no_skill"
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

// ---------------------------------------------------------------------------
// HTML 辅助工具函数
// ---------------------------------------------------------------------------

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
	if label, ok := contracts.ScenarioLabels[scenario]; ok {
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

func truncateText(value string, limit int) string {
	if limit <= 0 || len(value) <= limit {
		return value
	}
	return value[:limit]
}
