package reporter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
	payload := contracts.ReportPayload{
		SchemaVersion:    contracts.SchemaVersion,
		RunID:            spec.RunID,
		GeneratedAtUTC:   time.Now().UTC(),
		SourceEvaluation: evaluationPath,
		Summary:          summary,
	}

	breakdown := buildMutationBreakdown(set.Results)

	jsonPath := filepath.Join(reportRoot, "report_summary.json")
	htmlPath := filepath.Join(reportRoot, "report.html")
	if err := contracts.WriteJSON(jsonPath, map[string]any{
		"schema_version":     payload.SchemaVersion,
		"run_id":             payload.RunID,
		"generated_at_utc":   payload.GeneratedAtUTC,
		"source_evaluation":  payload.SourceEvaluation,
		"summary":            payload.Summary,
		"mutation_breakdown": breakdown,
	}); err != nil {
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
	testPass := 0
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
		if row.TestPass != nil && *row.TestPass {
			testPass++
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
	lineAvg := avg(lineSum, lineCnt)
	mutationAvg := avg(mutationSum, mutationCnt)
	densityAvg := avg(assertDensitySum, assertDensityCnt)
	if lineCnt == 0 {
		lineAvg = 0
	}
	if mutationCnt == 0 {
		mutationAvg = 0
	}
	if assertDensityCnt == 0 {
		densityAvg = 0
	}

	return contracts.ReportSummary{
		TotalSamples:        total,
		CompilePassCount:    compilePass,
		CompilePassRate:     rate(compilePass, total),
		TestPassCount:       testPass,
		TestPassRate:        rate(testPass, compilePass),
		AvgLineCoverage:     lineAvg,
		AvgMutationScore:    mutationAvg,
		AvgAssertionDensity: densityAvg,
	}
}

func buildHTML(payload contracts.ReportPayload, breakdown mutationBreakdown, rows []contracts.EvaluationResult) string {
	var b strings.Builder
	b.WriteString("<!doctype html><html lang=\"zh-CN\"><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><title>utbench report</title>")
	b.WriteString(`<style>
    :root { --bg:#f4f8fb; --card:#fff; --line:#e5e7eb; --text:#102a43; --muted:#6b7280; --accent:#0ea5a4; --accent-soft:#e6f8f7; --ok:#059669; --fail:#dc2626; }
    body { margin:0; background:linear-gradient(180deg,#f3fbfb 0%,#f7f8fc 100%); color:var(--text); font-family:"Noto Sans SC","PingFang SC","Microsoft YaHei",sans-serif; }
    .wrap { max-width:1160px; margin:0 auto; padding:20px; }
    .hero { background:radial-gradient(circle at right top, var(--accent-soft), #fff 58%); border:1px solid var(--line); border-radius:14px; padding:16px 18px; }
    h1 { margin:0; font-size:26px; }
    .sub { margin-top:8px; font-size:13px; color:var(--muted); }
    .cards { display:grid; gap:10px; grid-template-columns:repeat(auto-fill,minmax(180px,1fr)); margin:16px 0 22px; }
    .card { background:var(--card); border:1px solid var(--line); border-radius:10px; padding:10px 12px; }
    .k { font-size:12px; color:var(--muted); }
    .v { margin-top:4px; font-size:22px; font-weight:700; color:var(--accent); }
    h2 { margin:20px 0 10px; font-size:20px; }
    table { width:100%; border-collapse:collapse; background:#fff; border:1px solid var(--line); border-radius:10px; overflow:hidden; }
    th,td { border-bottom:1px solid var(--line); padding:8px 10px; font-size:13px; text-align:left; vertical-align:top; }
    th { background:#f8fbff; }
    tr:last-child td { border-bottom:none; }
    .pill-ok{color:var(--ok);font-weight:700;} .pill-fail{color:var(--fail);font-weight:700;} .pill-none{color:var(--muted);font-weight:700;}
    .small{font-size:12px;color:var(--muted);} .section{margin-bottom:22px;}
  </style>`)
	b.WriteString("</head><body><div class=\"wrap\">")
	b.WriteString("<div class=\"hero\"><h1>utbench 评测报告</h1>")
	b.WriteString("<div class=\"sub\">摘要 + 变异统计 + 样本明细。该样式参考旧版 Python reporter 的可读性布局。</div>")
	b.WriteString(fmt.Sprintf("<div class=\"small\">run_id: %s | generated_at: %s</div></div>", payload.RunID, payload.GeneratedAtUTC.Format(time.RFC3339)))

	b.WriteString("<div class=\"cards\">")
	b.WriteString(summaryCard("样本总数", fmt.Sprintf("%d", payload.Summary.TotalSamples)))
	b.WriteString(summaryCard("编译通过率", fmtPct(payload.Summary.CompilePassRate)))
	b.WriteString(summaryCard("测试通过率", fmtPct(payload.Summary.TestPassRate)))
	b.WriteString(summaryCard("平均行覆盖率", fmtPct(payload.Summary.AvgLineCoverage)))
	b.WriteString(summaryCard("平均变异得分", fmtPct(payload.Summary.AvgMutationScore)))
	b.WriteString(summaryCard("平均断言密度", fmt.Sprintf("%.2f", payload.Summary.AvgAssertionDensity)))
	b.WriteString("</div>")

	b.WriteString("<div class=\"section\"><h2>变异统计</h2><table>")
	b.WriteString("<tr><th>指标</th><th>值</th></tr>")
	b.WriteString(fmt.Sprintf("<tr><td>mutation_total</td><td>%d</td></tr>", breakdown.Total))
	b.WriteString(fmt.Sprintf("<tr><td>mutation_killed</td><td>%d</td></tr>", breakdown.Killed))
	b.WriteString(fmt.Sprintf("<tr><td>mutation_survived</td><td>%d</td></tr>", breakdown.Survived))
	b.WriteString(fmt.Sprintf("<tr><td>mutation_no_tests</td><td>%d</td></tr>", breakdown.NoTests))
	b.WriteString(fmt.Sprintf("<tr><td>mutation_timeouts</td><td>%d</td></tr>", breakdown.Timeouts))
	b.WriteString(fmt.Sprintf("<tr><td>mutation_skipped</td><td>%d</td></tr>", breakdown.Skipped))
	b.WriteString(fmt.Sprintf("<tr><td>mutation_suspicious</td><td>%d</td></tr>", breakdown.Suspicious))
	b.WriteString("</table></div>")

	b.WriteString("<div class=\"section\"><h2>样本明细（前20条）</h2><table>")
	b.WriteString("<tr><th>model</th><th>lang</th><th>sample_id</th><th>compile</th><th>test</th><th>test_rate</th><th>line_cov</th><th>branch_cov</th><th>error</th></tr>")
	limit := len(rows)
	if limit > 20 {
		limit = 20
	}
	for i := 0; i < limit; i++ {
		row := rows[i]
		status := "<span class=\"pill-none\">none</span>"
		if row.TestPass != nil && *row.TestPass {
			status = "<span class=\"pill-ok\">pass</span>"
		} else if row.TestPass != nil {
			status = "<span class=\"pill-fail\">fail</span>"
		}
		errText := shortErr(row.CompileError)
		if errText == "" {
			errText = shortErr(row.TestError)
		}
		if errText == "" {
			errText = shortErr(row.CoverageError)
		}
		if errText == "" {
			errText = shortErr(row.MutationError)
		}
		b.WriteString("<tr>")
		b.WriteString(fmt.Sprintf("<td>%s</td>", escapeHTML(row.Model)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", escapeHTML(row.Language)))
		b.WriteString(fmt.Sprintf("<td>%s</td>", escapeHTML(row.SampleID)))
		b.WriteString(fmt.Sprintf("<td>%v</td>", row.CompilePass))
		b.WriteString(fmt.Sprintf("<td>%s</td>", status))
		if row.TestPassRate != nil {
			b.WriteString(fmt.Sprintf("<td>%s</td>", fmtPct(*row.TestPassRate)))
		} else {
			b.WriteString("<td>-</td>")
		}
		if row.LineCoverage != nil {
			b.WriteString(fmt.Sprintf("<td>%s</td>", fmtPct(*row.LineCoverage)))
		} else {
			b.WriteString("<td>-</td>")
		}
		if row.BranchCoverage != nil {
			b.WriteString(fmt.Sprintf("<td>%s</td>", fmtPct(*row.BranchCoverage)))
		} else {
			b.WriteString("<td>-</td>")
		}
		b.WriteString(fmt.Sprintf("<td>%s</td>", escapeHTML(errText)))
		b.WriteString("</tr>")
	}
	b.WriteString("</table>")
	b.WriteString("<p class=\"small\">完整明细见 evaluation_result.json。</p></div>")
	b.WriteString("</div></body></html>")
	return b.String()
}

func summaryCard(key, value string) string {
	return fmt.Sprintf("<div class=\"card\"><div class=\"k\">%s</div><div class=\"v\">%s</div></div>", escapeHTML(key), escapeHTML(value))
}

func fmtPct(v float64) string {
	return fmt.Sprintf("%.2f%%", v*100)
}

func shortErr(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ReplaceAll(v, "\n", " ")
	if len(v) > 200 {
		return v[:200] + "..."
	}
	return v
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
