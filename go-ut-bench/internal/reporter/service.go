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
	Total      int `json:"total"`
	Killed     int `json:"killed"`
	Survived   int `json:"survived"`
	NoTests    int `json:"no_tests"`
	Timeouts   int `json:"timeouts"`
	Skipped    int `json:"skipped"`
	Suspicious int `json:"suspicious"`
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

	summary := buildSummary(set.Results)
	dims := buildDimensions(set.Results)
	topModels := buildTopModels(dims.ByModel)
	failures := buildFailureRows(set.Results)
	breakdown := buildMutationBreakdown(set.Results)

	payload := contracts.ReportPayload{
		SchemaVersion:    contracts.SchemaVersion,
		RunID:            spec.RunID,
		GeneratedAtUTC:   time.Now().UTC(),
		SourceEvaluation: evaluationPath,
		Summary:          summary,
		Dimensions:       dims,
		TopModels:        topModels,
		Failures:         failures,
		Thresholds: contracts.Thresholds{
			CompilePassRate: 1.0,
			TestPassRate:    0.7,
			LineCoverage:    0.7,
			BranchCoverage:  0.6,
			MutationScore:   0.85,
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

func buildSummary(rows []contracts.EvaluationResult) contracts.ReportSummary {
	total := len(rows)
	compilePass := 0
	testPassTotal := 0
	testTotal := 0
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
		if row.TestPassCount != nil {
			testPassTotal += *row.TestPassCount
		}
		if row.TestTotalCount != nil {
			testTotal += *row.TestTotalCount
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

func buildDimensions(rows []contracts.EvaluationResult) contracts.Dimensions {
	modelMap := map[string]*modelAgg{}
	langMap := map[string]*modelAgg{}

	for _, row := range rows {
		agg := getOrCreateModelAgg(modelMap, row.Model)
		mergeModelAgg(agg, row)

		langAgg := getOrCreateModelAgg(langMap, row.Language)
		mergeModelAgg(langAgg, row)
	}

	var byModel []contracts.ModelDim
	for _, agg := range modelMap {
		byModel = append(byModel, contracts.ModelDim{
			Model:             agg.key,
			TotalSamples:      agg.count,
			CompilePassRate:   rate(agg.compilePass, agg.count),
			AvgTestPassRate:   rate(agg.testPassTotal, agg.testTotal),
			AvgLineCoverage:   avg(agg.lineSum, agg.lineCnt),
			AvgBranchCoverage: avg(agg.branchSum, agg.branchCnt),
			AvgMutationScore:  avg(agg.mutationSum, agg.mutationCnt),
			AvgLatencyMS:      avgFloat(agg.latencySum, agg.latencyCnt),
			AvgTokens:         avgFloat(agg.tokenSum, agg.tokenCnt),
		})
	}
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

	return contracts.Dimensions{ByModel: byModel, ByLanguage: byLanguage}
}

type modelAgg struct {
	key           string
	count         int
	compilePass   int
	testPassTotal int
	testTotal     int
	lineSum       float64
	lineCnt       int
	branchSum     float64
	branchCnt     int
	mutationSum   float64
	mutationCnt   int
	latencySum    float64
	latencyCnt    int
	tokenSum      float64
	tokenCnt      int
}

func getOrCreateModelAgg(m map[string]*modelAgg, key string) *modelAgg {
	if a, ok := m[key]; ok {
		return a
	}
	a := &modelAgg{key: key}
	m[key] = a
	return a
}

func mergeModelAgg(a *modelAgg, row contracts.EvaluationResult) {
	a.count++
	if row.CompilePass {
		a.compilePass++
	}
	if row.TestPassCount != nil {
		a.testPassTotal += *row.TestPassCount
	}
	if row.TestTotalCount != nil {
		a.testTotal += *row.TestTotalCount
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
}

func buildTopModels(models []contracts.ModelDim) []contracts.ModelRank {
	var sorted []contracts.ModelDim
	for _, m := range models {
		sorted = append(sorted, m)
	}
	sort.Slice(sorted, func(i, j int) bool {
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
			Rank:             i + 1,
			Model:            m.Model,
			AvgTestPassRate:  m.AvgTestPassRate,
			AvgLineCoverage:  m.AvgLineCoverage,
			AvgMutationScore: m.AvgMutationScore,
			AvgLatencyMS:     m.AvgLatencyMS,
			AvgTokens:        m.AvgTokens,
		})
	}
	return out
}

func buildFailureRows(rows []contracts.EvaluationResult) []contracts.FailureRow {
	m := map[failureKey]*failureAgg{}

	for _, row := range rows {
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
			k := failureKey{stage: "mutation", errType: "mutation_error"}
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
	}
	return out
}

func buildHTML(payload contracts.ReportPayload, breakdown mutationBreakdown, rows []contracts.EvaluationResult) string {
	var b strings.Builder

	b.WriteString(`<!doctype html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>ut-bench 多模型单测生成可视化报告</title>
<script src="https://cdn.jsdelivr.net/npm/chart.js@4.4.1/dist/chart.umd.min.js"></script>
<style>
:root {
  --bg: linear-gradient(180deg, #f6fbff 0%, #f7f8fc 100%);
  --card: #ffffff;
  --text: #102a43;
  --muted: #6b7280;
  --accent: #0ea5a4;
  --accent-soft: #e6f8f7;
  --accent-2: #7c3aed;
  --warn: #f59e0b;
  --ok: #059669;
  --fail: #dc2626;
  --line: #e5e7eb;
}
* { box-sizing: border-box; }
body { margin: 0; background: var(--bg); color: var(--text); font-family: "Noto Sans SC", "Microsoft YaHei", "PingFang SC", sans-serif; }
.wrap { max-width: 1200px; margin: 0 auto; padding: 24px; }
.hero { background: radial-gradient(circle at right top, var(--accent-soft), #fff 55%); border: 1px solid var(--line); border-radius: 16px; padding: 18px 20px; margin-bottom: 16px; }
.hero h1 { margin: 0; font-size: 28px; letter-spacing: 0.3px; }
.hero .sub { margin-top: 6px; color: var(--muted); font-size: 13px; line-height: 1.6; }
.meta { color: var(--muted); font-size: 12px; margin-top: 10px; }
.cards { display: grid; gap: 12px; grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); margin-bottom: 24px; }
.card { background: var(--card); border: 1px solid var(--line); border-radius: 10px; padding: 12px 14px; }
.card .k { color: var(--muted); font-size: 12px; }
.card .v { margin-top: 4px; font-size: 22px; font-weight: 700; color: var(--accent); }
h2 { margin: 24px 0 10px; font-size: 20px; }
h3 { margin: 16px 0 10px; font-size: 16px; color: #334e68; }
.section { margin-bottom: 26px; }
.grid-2 { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.grid-3 { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }
.panel { background: var(--card); border: 1px solid var(--line); border-radius: 10px; padding: 12px; }
.chart-box { height: 280px; position: relative; }
.chart-box canvas { width: 100% !important; height: 100% !important; }
table { width: 100%; border-collapse: collapse; background: var(--card); border: 1px solid var(--line); border-radius: 10px; overflow: hidden; }
th, td { border-bottom: 1px solid var(--line); padding: 8px 10px; font-size: 13px; text-align: left; vertical-align: top; }
th { background: #f8fbff; font-weight: 600; }
tr:last-child td { border-bottom: none; }
tr:hover td { background: #fafcff; }
.rank-badge { display: inline-block; width: 22px; height: 22px; border-radius: 50%; background: var(--accent); color: #fff; text-align: center; line-height: 22px; font-size: 12px; font-weight: 700; margin-right: 6px; }
.rank-1 { background: #fbbf24; color: #78350f; }
.rank-2 { background: #94a3b8; color: #fff; }
.rank-3 { background: #cd7c2f; color: #fff; }
.rate-bar { display: flex; align-items: center; gap: 8px; }
.rate-bar .bar { flex: 1; height: 6px; background: #e5e7eb; border-radius: 3px; overflow: hidden; }
.rate-bar .fill { height: 100%; border-radius: 3px; background: var(--accent); }
.rate-bar .fill.warn { background: var(--warn); }
.rate-bar .fill.fail { background: var(--fail); }
.rate-bar .pct { font-size: 12px; color: var(--muted); min-width: 40px; }
.pill { display: inline-block; padding: 2px 8px; border-radius: 999px; font-size: 11px; margin-left: 4px; background: #eef2ff; color: #3730a3; }
.pill-ok { color: var(--ok); font-weight: 600; }
.pill-fail { color: var(--fail); font-weight: 600; }
.pill-none { color: var(--muted); }
.small { color: var(--muted); font-size: 12px; }
.conclusion { margin: 0; padding-left: 18px; color: #243b53; line-height: 1.8; }
.conclusion li strong { color: var(--accent); }
.error-cell { max-width: 200px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; font-size: 11px; color: var(--fail); }
.mutation-total { font-size: 11px; color: var(--muted); }
@media (max-width: 900px) {
  .grid-2, .grid-3 { grid-template-columns: 1fr; }
  .wrap { padding: 14px; }
  .hero h1 { font-size: 22px; }
  .chart-box { height: 220px; }
}
</style>
</head>
<body>
<div class="wrap">
<div class="hero">
  <h1>ut-bench 多模型单测生成可视化报告</h1>
  <div class="sub">摘要 + 详细分析 + 附录三段式呈现。支持按模型、语言进行多维分析，并输出 JSON/HTML 产物。</div>
  <div class="meta">run_id: ` + escapeHTML(payload.RunID) + ` | 生成时间：` + payload.GeneratedAtUTC.Format("2006-01-02 15:04:05") + ` | 数据源：` + escapeHTML(payload.SourceEvaluation) + `</div>
</div>

<div class="section">
  <h2>一、摘要</h2>
  <div class="cards">`)

	b.WriteString(summaryCard("样本总数", fmt.Sprintf("%d", payload.Summary.TotalSamples)))
	b.WriteString(summaryCard("编译通过率", fmtPct(payload.Summary.CompilePassRate)))
	b.WriteString(summaryCard("测试通过率", fmtPct(payload.Summary.TestPassRate)))
	b.WriteString(summaryCard("平均行覆盖率", fmtPct(payload.Summary.AvgLineCoverage)))
	b.WriteString(summaryCard("平均变异得分", fmtPct(payload.Summary.AvgMutationScore)))
	b.WriteString(summaryCard("有效杀死率", fmtPct(effectiveKillRate(breakdown))))
	b.WriteString(summaryCard("变异体总数", fmt.Sprintf("%d", breakdown.Total)))
	b.WriteString(summaryCard("平均断言密度", fmt.Sprintf("%.2f", payload.Summary.AvgAssertionDensity)))
	b.WriteString(summaryCard("有效样本数", fmt.Sprintf("%d", payload.Summary.TestPassCount)))

	b.WriteString(`</div>
  <h3>模型综合排名 <span class="pill">按测试通过率→覆盖率→变异得分</span></h3>
  <table>
  <tr><th>#</th><th>模型</th><th>测试通过率</th><th>行覆盖率</th><th>分支覆盖率</th><th>变异得分</th><th>生成耗时(ms)</th><th>总Token</th></tr>`)

	for _, r := range payload.TopModels {
		bgClass := ""
		if r.Rank <= 3 {
			bgClass = fmt.Sprintf(" class=\"rank-%d\"", r.Rank)
		}
		b.WriteString(fmt.Sprintf("<tr%s><td><span class=\"rank-badge%s\">%d</span></td><td><strong>%s</strong></td>",
			bgClass, rankClass(r.Rank), r.Rank, escapeHTML(r.Model)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", rateCell(r.AvgTestPassRate, payload.Thresholds.TestPassRate)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", rateCell(r.AvgLineCoverage, payload.Thresholds.LineCoverage)))
		b.WriteString(fmt.Sprintf("<td>-</td>"))
		b.WriteString(fmt.Sprintf("<td>%s</td>", rateCell(r.AvgMutationScore, payload.Thresholds.MutationScore)))
		b.WriteString(fmt.Sprintf("<td>%s</td></tr>", numCell(r.AvgLatencyMS)))
	}
	b.WriteString("</table>")
	b.WriteString("</div>")

	b.WriteString(`<div class="section">
  <h2>二、详细分析</h2>
  <div class="grid-2">
    <div class="panel">
      <h3>模型关键指标柱状图</h3>
      <div class="chart-box"><canvas id="barModel"></canvas></div>
    </div>
    <div class="panel">
      <h3>多维能力雷达图（模型均值）</h3>
      <div class="chart-box"><canvas id="radarDimensions"></canvas></div>
    </div>
  </div>
  <h3>失败类型统计</h3>`)
	b.WriteString(failureTable(payload.Failures))
	b.WriteString(`<h3>结论摘要</h3><ul class="conclusion">`)
	b.WriteString(fmt.Sprintf("<li><strong>正确性：</strong>整体编译通过率 %s，平均测试通过率 %s。</li>",
		fmtPct(payload.Summary.CompilePassRate), fmtPct(payload.Summary.TestPassRate)))
	b.WriteString(fmt.Sprintf("<li><strong>覆盖率：</strong>行覆盖率 %s，分支覆盖率 %s。</li>",
		fmtPct(payload.Summary.AvgLineCoverage), fmtPct(0.0)))
	effectiveKR := effectiveKillRate(breakdown)
	b.WriteString(fmt.Sprintf("<li><strong>有效性：</strong>变异得分均值 %s，有效变异杀死率 %s，与目标 %s 对比可持续优化。</li>",
		fmtPct(payload.Summary.AvgMutationScore), fmtPct(effectiveKR), fmtPct(payload.Thresholds.MutationScore)))
	b.WriteString("</ul>")
	b.WriteString("</div>")

	b.WriteString(`<div class="section">
  <h2>三、附录</h2>
  <h3>按模型统计</h3>`)
	b.WriteString(dimModelTable(payload.Dimensions.ByModel, payload.Thresholds))
	b.WriteString(`<h3>按语言统计</h3>`)
	b.WriteString(dimLanguageTable(payload.Dimensions.ByLanguage, payload.Thresholds))
	b.WriteString(`<h3>变异统计详情</h3><table>
  <tr><th>指标</th><th>值</th></tr>
  <tr><td>mutation_total</td><td>` + fmt.Sprintf("%d", breakdown.Total) + `</td></tr>
  <tr><td>mutation_killed</td><td>` + fmt.Sprintf("%d", breakdown.Killed) + `</td></tr>
  <tr><td>mutation_survived</td><td>` + fmt.Sprintf("%d", breakdown.Survived) + `</td></tr>
  <tr><td>mutation_no_tests</td><td>` + fmt.Sprintf("%d", breakdown.NoTests) + `</td></tr>
  <tr><td>mutation_timeouts</td><td>` + fmt.Sprintf("%d", breakdown.Timeouts) + `</td></tr>
  <tr><td>mutation_skipped</td><td>` + fmt.Sprintf("%d", breakdown.Skipped) + `</td></tr>
  <tr><td>mutation_suspicious</td><td>` + fmt.Sprintf("%d", breakdown.Suspicious) + `</td></tr>
  </table>`)

	b.WriteString(`<h3>原始样本明细</h3>
  <p class="small">共 ` + fmt.Sprintf("%d", len(rows)) + ` 条记录</p>
  <table id="sampleTable">
  <thead><tr>
    <th>模型</th><th>语言</th><th>样本ID</th><th>编译</th><th>测试</th><th>通过率</th>
    <th>行覆盖</th><th>分支覆盖</th><th>变异得分</th><th>变异体</th><th>断言密度</th><th>错误</th>
  </tr></thead><tbody>`)

	for _, row := range rows {
		testStatus := `<span class="pill-none">-</span>`
		if row.TestPass != nil {
			if *row.TestPass {
				testStatus = `<span class="pill pill-ok">pass</span>`
			} else {
				testStatus = `<span class="pill pill-fail">fail</span>`
			}
		}
		compileStatus := "✗"
		if row.CompilePass {
			compileStatus = "✓"
		}

		errText := shortErrText(row.CompileError)
		if errText == "" {
			errText = shortErrText(row.TestError)
		}
		if errText == "" {
			errText = shortErrText(row.CoverageError)
		}
		if errText == "" {
			errText = shortErrText(row.MutationError)
		}

		b.WriteString("<tr>")
		b.WriteString(fmt.Sprintf("<td>%s</td>", escapeHTML(row.Model)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", escapeHTML(row.Language)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", escapeHTML(row.SampleID)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", compileStatus))
		b.WriteString(fmt.Sprintf("<td>%s</td>", testStatus))
		b.WriteString(fmt.Sprintf("<td>%s</td>", optPctCell(row.TestPassRate)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", optPctCell(row.LineCoverage)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", optPctCell(row.BranchCoverage)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", optPctCell(row.MutationScore)))
		b.WriteString(fmt.Sprintf("<td class=\"mutation-total\">%s</td>", optIntCell(row.MutationTotal)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", optFloatCell(row.AssertionDensity)))
		b.WriteString(fmt.Sprintf("<td class=\"error-cell\" title=\"%s\">%s</td>", escapeHTML(errText), escapeHTML(errText)))
		b.WriteString("</tr>")
	}
	b.WriteString("</tbody></table>")
	b.WriteString("</div></div>")

	modelNamesJSON, testRatesJSON, lineCovsJSON, mutScoresJSON, compileRatesJSON := chartData(payload.Dimensions.ByModel)
	b.WriteString(fmt.Sprintf(`<script>
var styleConfig = {"accent":"#0ea5a4","accent_soft":"#e6f8f7","accent_2":"#7c3aed","warn":"#f59e0b","ok":"#059669"};
var modelNames = %s;
var testPassRates = %s;
var lineCoverages = %s;
var mutationScores = %s;
var compileRates = %s;
var thresholds = {testPassRate:%f,lineCoverage:%f,mutationScore:%f};

document.addEventListener("DOMContentLoaded", function() {
  new Chart(document.getElementById("barModel"), {
    type: "bar",
    data: {
      labels: modelNames,
      datasets: [
        {label:"测试通过率",data:testPassRates,backgroundColor:"rgba(14,165,164,0.7)"},
        {label:"行覆盖率",data:lineCoverages,backgroundColor:"rgba(124,58,237,0.7)"},
        {label:"变异得分",data:mutationScores,backgroundColor:"rgba(245,158,11,0.7)"}
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {legend:{position:"top"},tooltip:{callbacks:{
        label: function(ctx){return ctx.dataset.label+": "+(ctx.raw*100).toFixed(1)+"%%";}
      }}},
      scales: {y:{beginAtZero:true,max:1}}
    }
  });

  var radarLabels = ["测试通过率","行覆盖率","变异得分"];
  var radarData = modelNames.map(function(_, i){
    return [testPassRates[i]||0,lineCoverages[i]||0,mutationScores[i]||0];
  });
  var radarDatasets = radarData.map(function(d,i){
    return {label:modelNames[i],data:d,fill:true};
  });
  new Chart(document.getElementById("radarDimensions"), {
    type: "radar",
    data: {labels:radarLabels,datasets:radarDatasets},
    options: {
      responsive:true,maintainAspectRatio:false,
      plugins:{tooltip:{callbacks:{label:function(ctx){return ctx.dataset.label+": "+(ctx.raw*100).toFixed(1)+"%%";}}}},
      scales:{r:{beginAtZero:true,max:1}}
    }
  });
});
</script>`, modelNamesJSON, testRatesJSON, lineCovsJSON, mutScoresJSON, compileRatesJSON,
		payload.Thresholds.TestPassRate, payload.Thresholds.LineCoverage, payload.Thresholds.MutationScore))

	b.WriteString("</body></html>")
	return b.String()
}

func chartData(models []contracts.ModelDim) (names, testRates, lineCovs, mutScores, compileRates string) {
	var ns []string
	var trs, lcs, mss, crs []float64
	for _, m := range models {
		ns = append(ns, m.Model)
		trs = append(trs, m.AvgTestPassRate)
		lcs = append(lcs, m.AvgLineCoverage)
		mss = append(mss, m.AvgMutationScore)
		crs = append(crs, m.CompilePassRate)
	}
	names = marshalJSON(ns)
	testRates = marshalJSONF(trs)
	lineCovs = marshalJSONF(lcs)
	mutScores = marshalJSONF(mss)
	compileRates = marshalJSONF(crs)
	return
}

func marshalJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func marshalJSONF(v []float64) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func dimModelTable(models []contracts.ModelDim, thresh contracts.Thresholds) string {
	var b strings.Builder
	b.WriteString(`<table>
  <tr><th>模型</th><th>样本数</th><th>编译通过率</th><th>测试通过率</th><th>行覆盖</th><th>分支覆盖</th><th>变异得分</th><th>平均耗时</th></tr>`)
	for _, m := range models {
		b.WriteString(fmt.Sprintf("<tr><td><strong>%s</strong></td>", escapeHTML(m.Model)))
		b.WriteString(fmt.Sprintf("<td>%d</td>", m.TotalSamples))
		b.WriteString(fmt.Sprintf("<td>%s</td>", rateCell(m.CompilePassRate, thresh.TestPassRate)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", rateCell(m.AvgTestPassRate, thresh.TestPassRate)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", rateCell(m.AvgLineCoverage, thresh.LineCoverage)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", rateCell(m.AvgBranchCoverage, thresh.BranchCoverage)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", rateCell(m.AvgMutationScore, thresh.MutationScore)))
		b.WriteString(fmt.Sprintf("<td>%s</td></tr>", numCell(m.AvgLatencyMS)))
	}
	b.WriteString("</table>")
	return b.String()
}

func dimLanguageTable(langs []contracts.LanguageDim, thresh contracts.Thresholds) string {
	var b strings.Builder
	b.WriteString(`<table>
  <tr><th>语言</th><th>样本数</th><th>编译通过率</th><th>测试通过率</th><th>行覆盖</th><th>分支覆盖</th><th>变异得分</th></tr>`)
	for _, l := range langs {
		b.WriteString(fmt.Sprintf("<tr><td><strong>%s</strong></td>", escapeHTML(l.Language)))
		b.WriteString(fmt.Sprintf("<td>%d</td>", l.TotalSamples))
		b.WriteString(fmt.Sprintf("<td>%s</td>", rateCell(l.CompilePassRate, thresh.TestPassRate)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", rateCell(l.AvgTestPassRate, thresh.TestPassRate)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", rateCell(l.AvgLineCoverage, thresh.LineCoverage)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", rateCell(l.AvgBranchCoverage, thresh.BranchCoverage)))
		b.WriteString(fmt.Sprintf("<td>%s</td></tr>", rateCell(l.AvgMutationScore, thresh.MutationScore)))
	}
	b.WriteString("</table>")
	return b.String()
}

func failureTable(failures []contracts.FailureRow) string {
	if len(failures) == 0 {
		return `<p class="small">无失败记录。</p>`
	}
	var b strings.Builder
	b.WriteString(`<table>
  <tr><th>阶段</th><th>错误类型</th><th>数量</th><th>示例模型</th><th>示例样本</th><th>示例信息</th></tr>`)
	for _, f := range failures {
		b.WriteString(fmt.Sprintf("<tr><td>%s</td>", escapeHTML(f.Stage)))
		b.WriteString(fmt.Sprintf("<td><span class=\"pill pill-fail\">%s</span></td>", escapeHTML(f.ErrorType)))
		b.WriteString(fmt.Sprintf("<td>%d</td>", f.Count))
		b.WriteString(fmt.Sprintf("<td>%s</td>", escapeHTML(f.ExampleModel)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", escapeHTML(f.ExampleSample)))
		b.WriteString(fmt.Sprintf("<td class=\"error-cell\">%s</td></tr>", escapeHTML(f.ExampleMessage)))
	}
	b.WriteString("</table>")
	return b.String()
}

func rateCell(v, threshold float64) string {
	if v == 0 {
		return `<span class="pill-none">-</span>`
	}
	cls := "ok"
	if v < threshold {
		cls = "fail"
	} else if v < threshold*1.1 {
		cls = "warn"
	}
	filled := int(v * 100)
	return fmt.Sprintf(`<div class="rate-bar"><div class="bar"><div class="fill %s" style="width:%d%%"></div></div><span class="pct %s">%s</span></div>`,
		cls, min(100, filled), cls, fmtPct(v))
}

func numCell(v float64) string {
	if v == 0 {
		return `<span class="pill-none">-</span>`
	}
	return fmt.Sprintf("%.0f", v)
}

func optPctCell(v *float64) string {
	if v == nil {
		return `<span class="pill-none">-</span>`
	}
	return fmtPct(*v)
}

func optIntCell(v *int) string {
	if v == nil {
		return `<span class="pill-none">-</span>`
	}
	return fmt.Sprintf("%d", *v)
}

func optFloatCell(v *float64) string {
	if v == nil {
		return `<span class="pill-none">-</span>`
	}
	return fmt.Sprintf("%.2f", *v)
}

func summaryCard(key, value string) string {
	return fmt.Sprintf("<div class=\"card\"><div class=\"k\">%s</div><div class=\"v\">%s</div></div>", escapeHTML(key), escapeHTML(value))
}

func fmtPct(v float64) string {
	return fmt.Sprintf("%.1f%%", v*100)
}

func escapeHTML(v string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
	)
	return replacer.Replace(v)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func effectiveKillRate(b mutationBreakdown) float64 {
	total := b.Killed + b.Survived
	if total == 0 {
		return 0
	}
	return round(float64(b.Killed)/float64(total), 6)
}

func rankClass(rank int) string {
	switch rank {
	case 1:
		return " rank-1"
	case 2:
		return " rank-2"
	case 3:
		return " rank-3"
	default:
		return ""
	}
}
