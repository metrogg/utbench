package reporter

import (
	"fmt"
	"strings"

	"go-ut-bench/internal/contracts"
)

// buildInsightsSection 生成洞察区域HTML
func buildInsightsSection(insights contracts.Insights) string {
	if insights.BestModel.Title == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<div class="section" id="insights">
  <h2>核心洞察 Key Insights</h2>
  <div class="insights-grid-3">`)

	// SVG icons (inline, consistent stroke width 2, size 24x24)
	const trophySVG = `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"/><path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"/><path d="M4 22h16"/><path d="M10 14.66V17c0 .55-.47.98-.97 1.21C7.85 18.75 7 20.24 7 22"/><path d="M14 14.66V17c0 .55.47.98.97 1.21C16.15 18.75 17 20.24 17 22"/><path d="M18 2H6v7a6 6 0 0 0 12 0V2Z"/></svg>`
	const warningSVG = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><line x1="12" x2="12" y1="9" y2="13"/><line x1="12" x2="12.01" y1="17" y2="17"/></svg>`
	const starSVG = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg>`
	const chartSVG = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" x2="12" y1="20" y2="10"/><line x1="18" x2="18" y1="20" y2="4"/><line x1="6" x2="6" y1="20" y2="16"/></svg>`
	const wrenchSVG = `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>`

	// 最佳模型
	if insights.BestModel.Title != "" {
		b.WriteString(fmt.Sprintf(`
    <div class="insight-card insight-best" style="background:#0f766e;color:#fff;padding:16px;border-radius:12px;">
      <div style="margin-bottom:8px;">%s</div>
      <div style="font-size:18px;font-weight:700;margin-bottom:6px;">%s</div>
      <div style="font-size:14px;line-height:1.5;">%s</div>
    </div>`, trophySVG, escapeHTML(insights.BestModel.Title), escapeHTML(insights.BestModel.Detail)))
	}

	// 弱项场景
	for _, ws := range insights.WeakScenarios {
		b.WriteString(fmt.Sprintf(`
    <div class="insight-card" style="background:#fef3c7;border:1px solid #f59e0b;padding:16px;border-radius:12px;">
      <div style="color:#b45309;margin-bottom:6px;">%s</div>
      <div style="font-size:16px;font-weight:600;color:#92400e;">%s</div>
      <div style="font-size:13px;color:#78350f;margin-top:4px;">%s</div>
    </div>`, warningSVG, escapeHTML(ws.Title), escapeHTML(ws.Detail)))
	}

	// 强项场景
	for _, ss := range insights.StrongScenarios {
		b.WriteString(fmt.Sprintf(`
    <div class="insight-card" style="background:#d1fae5;border:1px solid #10b981;padding:16px;border-radius:12px;">
      <div style="color:#047857;margin-bottom:6px;">%s</div>
      <div style="font-size:16px;font-weight:600;color:#065f46;">%s</div>
      <div style="font-size:13px;color:#047857;margin-top:4px;">%s</div>
    </div>`, starSVG, escapeHTML(ss.Title), escapeHTML(ss.Detail)))
	}

	// 语言差异
	for _, lg := range insights.LanguageGaps {
		b.WriteString(fmt.Sprintf(`
    <div class="insight-card" style="background:#e6f4f1;border:1px solid #0f766e;padding:16px;border-radius:12px;">
      <div style="color:#0f766e;margin-bottom:6px;">%s</div>
      <div style="font-size:16px;font-weight:600;color:#115e59;">%s</div>
      <div style="font-size:13px;color:#0f766e;margin-top:4px;">%s</div>
    </div>`, chartSVG, escapeHTML(lg.Title), escapeHTML(lg.Detail)))
	}

	// 改进建议
	for _, rec := range insights.Recommendations {
		b.WriteString(fmt.Sprintf(`
    <div class="insight-card" style="background:#fce7f3;border:1px solid #ec4899;padding:16px;border-radius:12px;">
      <div style="color:#be185d;margin-bottom:6px;">%s</div>
      <div style="font-size:16px;font-weight:600;color:#9d174d;">%s</div>
      <div style="font-size:13px;color:#831843;margin-top:4px;">%s</div>
    </div>`, wrenchSVG, escapeHTML(rec.Title), escapeHTML(rec.Detail)))
	}

	b.WriteString(`
  </div>
</div>`)
	return b.String()
}

// buildCompareSection 生成对比分析区域HTML
func buildCompareSection(topModels []contracts.ModelRank) string {
	if len(topModels) < 2 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<div class="section" id="compare">
  <h2>对比分析 Model Compare</h2>
  <p class="muted" style="margin-bottom:12px;">选择两个模型进行直接对比，查看各维度差距。</p>
  <div style="display:flex;gap:12px;margin-bottom:16px;">
    <select id="compare-model-a" style="padding:10px 14px;border:1px solid #cbd5e1;border-radius:8px;background:#fff;min-height:44px;">
      <option value="">选择模型A</option>`)
	for _, m := range topModels {
		b.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, escapeHTML(m.Model), escapeHTML(m.Model)))
	}
	b.WriteString(`    </select>
    <select id="compare-model-b" style="padding:10px 14px;border:1px solid #cbd5e1;border-radius:8px;background:#fff;min-height:44px;">
      <option value="">选择模型B</option>`)
	for _, m := range topModels {
		b.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, escapeHTML(m.Model), escapeHTML(m.Model)))
	}
	b.WriteString(`    </select>
    <button onclick="runCompare()" style="padding:10px 16px;background:#1e40af;color:#fff;border:none;border-radius:8px;cursor:pointer;min-height:44px;">开始对比</button>
  </div>
  <div id="compare-result" style="background:#f8fafc;border:1px solid #e2e8f0;border-radius:12px;padding:16px;">
    <div style="color:#64748b;text-align:center;">选择两个模型后点击"开始对比"</div>
  </div>
</div>`)
	return b.String()
}

// buildEfficiencySection 生成效率分析区域HTML
func buildEfficiencySection(stats contracts.EfficiencyStats) string {
	if len(stats.TokenEfficiency) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<div class="section" id="efficiency">
  <h2>效率分析 Efficiency Analysis</h2>
  <div class="grid-2">
    <div class="panel">
      <h3>Token 效率排名</h3>
      <p class="muted" style="font-size:12px;">得分/千Token，数值越高表示用更少的Token获得更高的得分</p>
      <table style="width:100%%;font-size:13px;border-collapse:collapse;margin-top:8px;">
        <thead><tr style="background:#f1f5f9;"><th style="padding:6px;">排名</th><th style="padding:6px;">模型</th><th style="padding:6px;">效率值</th><th style="padding:6px;">平均Token</th></tr></thead>
        <tbody>`)
	for _, row := range stats.TokenEfficiency {
		b.WriteString(fmt.Sprintf(`<tr><td style="padding:6px;border-bottom:1px solid #e2e8f0;">#%d</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%s</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%.3f</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%.0f</td></tr>`,
			row.Rank, escapeHTML(row.Model), row.ScorePerToken, row.AvgTokens))
	}
	b.WriteString(`        </tbody>
      </table>
    </div>
    <div class="panel">
      <h3>时间效率排名</h3>
      <p class="muted" style="font-size:12px;">得分/秒，数值越高表示响应更快且得分更高</p>
      <table style="width:100%%;font-size:13px;border-collapse:collapse;margin-top:8px;">
        <thead><tr style="background:#f1f5f9;"><th style="padding:6px;">排名</th><th style="padding:6px;">模型</th><th style="padding:6px;">效率值</th><th style="padding:6px;">平均耗时</th></tr></thead>
        <tbody>`)
	for _, row := range stats.TimeEfficiency {
		latencyStr := fmt.Sprintf("%.1fs", row.AvgLatencyMS/1000)
		b.WriteString(fmt.Sprintf(`<tr><td style="padding:6px;border-bottom:1px solid #e2e8f0;">#%d</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%s</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%.3f</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%s</td></tr>`,
			row.Rank, escapeHTML(row.Model), row.ScorePerSecond, latencyStr))
	}
	b.WriteString(`        </tbody>
      </table>
    </div>
  </div>
  <div class="panel" style="margin-top:16px;">
    <h3>成本估算</h3>
    <div style="display:flex;gap:24px;margin-top:8px;font-size:14px;">
      <div><strong>总Token消耗:</strong> ` + fmt.Sprintf("%d", stats.CostEstimate.TotalTokens) + `</div>
      <div><strong>估算成本:</strong> ~¥` + fmt.Sprintf("%.2f", stats.CostEstimate.EstimatedCostUSD*7.2) + ` (按 ¥0.0072/1K Token 估算)</div>
    </div>
    <table style="width:100%%;font-size:13px;border-collapse:collapse;margin-top:12px;">
      <thead><tr style="background:#f1f5f9;"><th style="padding:6px;">模型</th><th style="padding:6px;">总Token</th><th style="padding:6px;">每样本成本(¥)</th></tr></thead>
      <tbody>`)
	for _, row := range stats.CostEstimate.ModelCostBreakdown {
		b.WriteString(fmt.Sprintf(`<tr><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%s</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%d</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;">¥%.4f</td></tr>`,
			escapeHTML(row.Model), row.TotalTokens, row.AvgCostPerSample*7.2))
	}
	b.WriteString(`      </tbody>
    </table>
  </div>
</div>`)
	return b.String()
}

// buildErrorDiagnosisSection 生成错误诊断区域HTML
func buildErrorDiagnosisSection(diagnosis contracts.ErrorDiagnosis) string {
	if len(diagnosis.CompileErrors) == 0 && len(diagnosis.TestErrors) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<div class="section" id="error-diagnosis">
  <h2>错误诊断 Error Diagnosis</h2>
  <p class="muted" style="margin-bottom:12px;">错误类型分类与改进建议</p>
  <div class="grid-2">`)

	// 编译错误
	if len(diagnosis.CompileErrors) > 0 {
		b.WriteString(`    <div class="panel">
      <h3 style="color:#dc2626;">编译错误分布</h3>
      <table style="width:100%%;font-size:13px;border-collapse:collapse;margin-top:8px;">
        <thead><tr style="background:#fee2e2;"><th style="padding:6px;">类型</th><th style="padding:6px;">次数</th><th style="padding:6px;">占比</th><th style="padding:6px;">影响模型</th></tr></thead>
        <tbody>`)
		for _, cat := range diagnosis.CompileErrors {
			modelsStr := strings.Join(cat.AffectedModels, ", ")
			if len(modelsStr) > 30 {
				modelsStr = modelsStr[:30] + "..."
			}
			b.WriteString(fmt.Sprintf(`<tr><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%s</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%d</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%.1f%%</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;font-size:11px;">%s</td></tr>`,
				escapeHTML(cat.Type), cat.Count, cat.Rate*100, escapeHTML(modelsStr)))
		}
		b.WriteString(`        </tbody>
      </table>
    </div>`)
	}

	// 测试错误
	if len(diagnosis.TestErrors) > 0 {
		b.WriteString(`    <div class="panel">
      <h3 style="color:#ea580c;">测试错误分布</h3>
      <table style="width:100%%;font-size:13px;border-collapse:collapse;margin-top:8px;">
        <thead><tr style="background:#fed7aa;"><th style="padding:6px;">类型</th><th style="padding:6px;">次数</th><th style="padding:6px;">占比</th><th style="padding:6px;">影响模型</th></tr></thead>
        <tbody>`)
		for _, cat := range diagnosis.TestErrors {
			modelsStr := strings.Join(cat.AffectedModels, ", ")
			if len(modelsStr) > 30 {
				modelsStr = modelsStr[:30] + "..."
			}
			b.WriteString(fmt.Sprintf(`<tr><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%s</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%d</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;">%.1f%%</td><td style="padding:6px;border-bottom:1px solid #e2e8f0;font-size:11px;">%s</td></tr>`,
				escapeHTML(cat.Type), cat.Count, cat.Rate*100, escapeHTML(modelsStr)))
		}
		b.WriteString(`        </tbody>
      </table>
    </div>`)
	}

	b.WriteString(`  </div>`)

	// 常见错误模式和改进建议
	if len(diagnosis.CommonPatterns) > 0 {
		b.WriteString(`  <div class="panel" style="margin-top:16px;">
    <h3>常见错误模式与建议</h3>`)
		for _, pattern := range diagnosis.CommonPatterns {
			b.WriteString(fmt.Sprintf(`
    <div style="background:#f8fafc;border:1px solid #e2e8f0;border-radius:8px;padding:12px;margin-bottom:8px;">
      <div style="font-weight:600;margin-bottom:4px;">%s</div>
      <div style="font-size:13px;color:#64748b;">%s</div>
    </div>`, escapeHTML(pattern.Pattern), escapeHTML(pattern.Advice)))
		}
		b.WriteString(`  </div>`)
	}

	// 总体建议
	if len(diagnosis.Recommendations) > 0 {
		b.WriteString(`  <div class="panel" style="margin-top:16px;background:#fce7f3;border:1px solid #ec4899;">
    <h3 style="color:#be185d;">改进建议</h3>`)
		for _, rec := range diagnosis.Recommendations {
			b.WriteString(fmt.Sprintf(`<div style="padding:8px;border-bottom:1px dashed #f9a8d4;">%s</div>`, escapeHTML(rec)))
		}
		b.WriteString(`  </div>`)
	}

	b.WriteString(`
</div>`)
	return b.String()
}

// buildMetaSection 生成报告元信息区域HTML
func buildMetaSection(runID string, specInfo string) string {
	return fmt.Sprintf(`<div class="section" id="meta">
  <h2>报告信息 Report Meta</h2>
  <div class="meta-grid-3" style="background:#f8fafc;padding:16px;border-radius:12px;font-size:14px;">
    <div><strong>Run ID:</strong> <span style="color:#64748b;">%s</span></div>
    <div><strong>报告版本:</strong> <span style="color:#64748b;">%s</span></div>
    <div><strong>提示词策略:</strong> <span style="color:#64748b;">structured-v1</span></div>
  </div>
  <div style="margin-top:12px;padding:12px;background:#e0e7ff;border-radius:8px;font-size:13px;">
    <strong>评分规则说明：</strong> 综合得分 = 编译通过率×0.3 + 样本测试通过率×0.3 + 行覆盖率×0.2 + 变异分数×0.2。
    其中变异分数反映测试用例检测代码缺陷的能力，通过变异测试工具注入缺陷来验证测试的有效性。
  </div>
</div>`, escapeHTML(runID), escapeHTML(specInfo))
}

// buildDimensionBreakdownSection 生成维度得分分解区域HTML
func buildDimensionBreakdownSection() string {
	return `<div class="section" id="dimension-breakdown">
  <h2>维度得分分解 Dimension Breakdown</h2>
  <p class="muted" style="margin-bottom:12px;">各模型在编译、测试、覆盖率、变异四个维度的详细得分对比。</p>
  <div class="table-wrap">
    <table id="dimensionBreakdownTable">
      <thead>
        <tr style="background:#f1f5f9;">
          <th style="padding:10px;text-align:left;">模型</th>
          <th style="padding:10px;text-align:center;">编译通过率</th>
          <th style="padding:10px;text-align:center;">样本测试通过率</th>
          <th style="padding:10px;text-align:center;">行覆盖率</th>
          <th style="padding:10px;text-align:center;">变异分数</th>
          <th style="padding:10px;text-align:center;">综合得分</th>
        </tr>
      </thead>
      <tbody id="dimensionBreakdownBody"></tbody>
    </table>
  </div>
</div>`
}
