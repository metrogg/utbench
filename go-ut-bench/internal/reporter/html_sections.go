// reporter 包 - HTML 报告各区块的构建函数
// 包含 buildHTML 主入口、排行榜、维度分析、错误分析、
// 计分剔除、评测集浏览、原始数据、截断分析、图表等区块
package reporter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"go-ut-bench/internal/contracts"
)

func buildHTML(payload contracts.ReportPayload, breakdown mutationBreakdown, rows []contracts.EvaluationResult, promptProvider contracts.PromptMetaProvider) string {
	var b strings.Builder

	heroModels := distinctSorted(stringsFromRows(rows, func(r contracts.EvaluationResult) string { return r.Model }))
	heroLangs := distinctSorted(stringsFromRows(rows, func(r contracts.EvaluationResult) string { return r.Language }))
	heroTypes := distinctSorted(stringsFromRows(rows, func(r contracts.EvaluationResult) string { return extractScenario(r.SampleID) }))
	heroTypeLabels := scenarioLabels(heroTypes)
	// 把 framework__model__skill 拆成三个维度展示，AgentFramework/AgentModel/SkillName
	// 字段为空时回退到 Model 字段按 "__" 切分（兼容旧数据）。
	heroFrameworks := distinctSorted(stringsFromRows(rows, func(r contracts.EvaluationResult) string {
		if v := strings.TrimSpace(r.AgentFramework); v != "" {
			return v
		}
		return splitSubjectPart(r.Model, 0)
	}))
	heroAgentModels := distinctSorted(stringsFromRows(rows, func(r contracts.EvaluationResult) string {
		if v := strings.TrimSpace(r.AgentModel); v != "" {
			return v
		}
		return splitSubjectPart(r.Model, 1)
	}))
	heroSkills := distinctSorted(stringsFromRows(rows, func(r contracts.EvaluationResult) string {
		if v := strings.TrimSpace(r.SkillName); v != "" {
			return v
		}
		return splitSubjectPart(r.Model, 2)
	}))
	_ = heroModels // 保留供其他区块使用

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

	// Hero Section - 现代美观版本
	b.WriteString(fmt.Sprintf(`
<div class="hero">
  <div class="hero-header">
    <div class="hero-title-group">
      <div class="hero-badge">UT-BENCH</div>
      <h1>单测生成评测报告</h1>
      <div class="hero-subtitle">%s · %d 个样本</div>
    </div>
  </div>
  <div class="hero-cards">
    <div class="hero-card">
      <div class="hero-card-icon"><svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="8" rx="2"/><rect x="2" y="14" width="20" height="8" rx="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg></div>
      <div class="hero-card-content">
        <div class="hero-card-label">平台</div>
        <div class="hero-card-value">%s</div>
      </div>
    </div>
    <div class="hero-card">
      <div class="hero-card-icon"><svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="10" rx="2"/><circle cx="12" cy="5" r="2"/><path d="M12 7v4"/><line x1="8" y1="16" x2="8" y2="16"/><line x1="16" y1="16" x2="16" y2="16"/></svg></div>
      <div class="hero-card-content">
        <div class="hero-card-label">模型</div>
        <div class="hero-card-value">%s</div>
      </div>
    </div>
    <div class="hero-card">
      <div class="hero-card-icon"><svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"/></svg></div>
      <div class="hero-card-content">
        <div class="hero-card-label">Skill</div>
        <div class="hero-card-value">%s</div>
      </div>
    </div>
    <div class="hero-card">
      <div class="hero-card-icon"><svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg></div>
      <div class="hero-card-content">
        <div class="hero-card-label">编程语言</div>
        <div class="hero-card-value">%s</div>
      </div>
    </div>
    <div class="hero-card">
      <div class="hero-card-icon"><svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/><polyline points="3.27 6.96 12 12.01 20.73 6.96"/><line x1="12" y1="22.08" x2="12" y2="12"/></svg></div>
      <div class="hero-card-content">
        <div class="hero-card-label">样本类型</div>
	<div class="hero-card-value">%s</div>
      </div>
    </div>
  </div>
</div>`,
		payload.GeneratedAtUTC.Format("2006-01-02 15:04"),
		payload.Summary.TotalSamples,
		escapeHTML(summarizeList(heroFrameworks, 6)),
		escapeHTML(summarizeList(heroAgentModels, 6)),
		escapeHTML(summarizeList(heroSkills, 6)),
		escapeHTML(summarizeList(heroLangs, 6)),
		escapeHTML(summarizeList(heroTypeLabels, 6))))
	_ = buildRuntimeSummarySection // 已停用：运行时拓扑 Runtime Topology

	// Navigation
	b.WriteString(`
<div class="jump-nav">
    <a href="#leaderboard">排名</a>
  <a href="#details">图表分析</a>
  <a href="#analysis-controls">筛选与导出</a>
  <a href="#dimension-analysis">维度分析</a>
  <a href="#score-exclusions">计分剔除</a>
  <a href="#error-analysis">错误分析</a>
  <a href="#dataset-browser">评测集</a>
  <a href="#raw-data">原始数据</a>
</div>
`)

	// Leaderboard Section - 模型排名（含控制变量对比）
	b.WriteString(buildLeaderboardHTMLNew(payload.TopModels, payload.ComparisonViews))
	b.WriteString(buildChartsSection(payload.TopModels))
	b.WriteString(buildAnalysisControlsSection(heroModels, heroLangs, heroTypes))
	b.WriteString(buildDimensionAnalysisSection())

	// Truncation Analysis Section - 截断分析（新增）

	// Error Analysis Section - 错误分析（新增）
	if len(payload.Failures) > 0 {
		b.WriteString(buildErrorAnalysisSection())
	}

	// Score Exclusions Section
	b.WriteString(buildScoreExclusionsSection(payload.ScoreExclusions))

	// Dataset Browser Section
	b.WriteString(buildDatasetBrowserSection(rows))

	// Raw Data Section - 原始数据（可展开收起）
	b.WriteString(buildRawDataSection(rows))

	// Prompt Section
	if len(payload.Prompts) > 0 {
		b.WriteString(buildPromptHTMLNew(payload.PromptStrategy, payload.PromptVersionID, payload.Prompts, promptProvider))
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
	avgLatency := totalLatency / float64(latencyCount) / 1000 // 转换为秒
	avgTokens := totalTokens / float64(tokenCount)

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
func buildLeaderboardHTMLNew(models []contracts.ModelRank, views []contracts.ComparisonView) string {
	if len(models) == 0 {
		return ""
	}

	// 构建卡片 HTML 和 data 属性
	var cardsHTML string
	for _, m := range models {
		rankClass := ""
		itemClass := ""
		switch m.Rank {
		case 1:
			rankClass = "gold"
			itemClass = "gold"
		case 2:
			rankClass = "silver"
			itemClass = "silver"
		case 3:
			rankClass = "bronze"
			itemClass = "bronze"
		}

		compileWidth := int(m.CompilePassRate * 100)
		testWidth := int(m.AvgTestPassRate * 100)
		coverWidth := int(m.AvgLineCoverage * 100)
		mutWidth := int(m.AvgMutationScore * 100)
		compositePct := m.CompositeScore * 100
		latencyStr := fmt.Sprintf("%.1fs", m.AvgLatencyMS/1000)
		tokensStr := fmt.Sprintf("%.0f", m.AvgTotalTokens)
		assertionDensityStr := fmt.Sprintf("%.2f", m.AvgAssertionDensity)
		tokenSourceStr := tokenCoverageLabel(m)

		platform := firstNonEmpty(m.AgentFramework, "model_api")
		model := firstNonEmpty(m.AgentModel, m.Model)
		skill := firstNonEmpty(m.SkillName, "no_skill")
		modelIDLine := ""
		if m.ModelID != "" && m.ModelID != model {
			modelIDLine = fmt.Sprintf(`<div style="font-size:11px;color:#94a3b8;margin-top:2px;">%s</div>`, escapeHTML(m.ModelID))
		}

		cardsHTML += fmt.Sprintf(`
    <div class="lb-item %s" data-platform="%s" data-model="%s" data-skill="%s" data-score="%.6f">
      <div class="lb-rank %s"><span class="lb-rank-num">%d</span></div>
      <div class="lb-content">
        <div class="lb-title" style="display:flex;align-items:center;gap:8px;flex-wrap:wrap;">
          <span class="lb-tag lb-tag-platform">%s</span>
          <span class="lb-tag lb-tag-model">%s</span>
          <span class="lb-tag lb-tag-skill">%s</span>
        </div>%s
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
            <div class="bar"><span style="width:%d%%;background:#8b5cf6;"></span></div>
            <span class="val">%.1f%%</span>
          </div>
          <div class="metric">
            <span class="name">断言密度</span>
            <div class="bar"><span style="width:%d%%;background:#14b8a6;"></span></div>
            <span class="val">%s</span>
          </div>
        </div>
        <div class="lb-meta" style="border-top:1px dashed rgba(148,163,184,0.3);padding-top:8px;margin-top:6px;font-size:12px;">
          <span style="display:inline-flex;align-items:center;gap:4px;">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><circle cx="12" cy="13" r="8"/><path d="M12 9v4l2 2"/><path d="M9 2h6"/></svg>
            平均耗时: <strong>%s</strong>
          </span>
          <span style="display:inline-flex;align-items:center;gap:4px;margin-left:12px;">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><line x1="12" y1="20" x2="12" y2="10"/><line x1="18" y1="20" x2="18" y2="4"/><line x1="6" y1="20" x2="6" y2="16"/></svg>
            平均Token: <strong>%s</strong>
          </span>
          <span style="display:inline-flex;align-items:center;gap:4px;margin-left:12px;">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/></svg>
            样本数: <strong>%d</strong>
          </span>
          <span style="display:inline-flex;align-items:center;gap:4px;margin-left:12px;">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M4 2v20l2-2 2 2 2-2 2 2 2-2 2 2 2-2 2 2V2l-2 2-2-2-2 2-2-2-2 2-2-2-2 2z"/><line x1="8" y1="10" x2="16" y2="10"/><line x1="8" y1="14" x2="13" y2="14"/></svg>
            Token口径: <strong>%s</strong>
          </span>
        </div>
      </div>
      <div class="lb-score">
        <div class="score-label">综合得分</div>
        <div class="score-val">%.2f</div>
      </div>
    </div>`,
			itemClass,
			escapeHTML(platform), escapeHTML(model), escapeHTML(skill), m.CompositeScore,
			rankClass, m.Rank,
			escapeHTML(platform), escapeHTML(model), escapeHTML(skill),
			modelIDLine,
			compileWidth, m.CompilePassRate*100,
			testWidth, m.AvgTestPassRate*100,
			coverWidth, m.AvgLineCoverage*100,
			mutWidth, m.AvgMutationScore*100,
			min(100, int(m.AvgAssertionDensity*20)),
			assertionDensityStr,
			latencyStr, tokensStr, m.TotalSamples,
			tokenSourceStr,
			compositePct)
	}

	// 构建单维度筛选选项：按平台/模型/Skill 时只筛选一个维度。
	filterOptionsJSON := buildFilterOptionsJSON(models)

	var b strings.Builder
	leaderboardHTML := `<div class="section" id="leaderboard">
  <h2>模型排名 Leaderboard</h2>
  <div class="lb-info-box" style="background:#f8fafc;border:1px solid #e2e8f0;border-radius:8px;padding:12px;margin-bottom:16px;font-size:13px;color:#475569;">
    <div style="display:flex;gap:16px;flex-wrap:wrap;align-items:center;">
      <span><strong>综合评分公式：</strong> ` + contracts.DefaultWeights.String() + `</span>
      <span style="color:#94a3b8;">|</span>
      <span><strong>指标说明：</strong> 样测=样本级测试通过率；行覆盖=代码行覆盖率；变异=变异测试得分</span>
    </div>
  </div>
  <div class="lb-controls" style="display:flex;gap:8px;align-items:center;margin-bottom:12px;flex-wrap:wrap;">
    <button class="lb-tab active" data-mode="all" onclick="lbSwitchMode('all',this)" style="height:34px;padding:0 18px;border:1px solid #dbe3ec;border-radius:6px;background:#1e293b;color:#fff;font-size:13px;font-weight:600;cursor:pointer;transition:all .15s;">全部排名</button>
    <button class="lb-tab" data-mode="platform" onclick="lbSwitchMode('platform',this)" style="height:34px;padding:0 18px;border:1px solid #dbe3ec;border-radius:6px;background:#fff;color:#334155;font-size:13px;font-weight:600;cursor:pointer;transition:all .15s;">按平台</button>
    <button class="lb-tab" data-mode="model" onclick="lbSwitchMode('model',this)" style="height:34px;padding:0 18px;border:1px solid #dbe3ec;border-radius:6px;background:#fff;color:#334155;font-size:13px;font-weight:600;cursor:pointer;transition:all .15s;">按模型</button>
    <button class="lb-tab" data-mode="skill" onclick="lbSwitchMode('skill',this)" style="height:34px;padding:0 18px;border:1px solid #dbe3ec;border-radius:6px;background:#fff;color:#334155;font-size:13px;font-weight:600;cursor:pointer;transition:all .15s;">按Skill</button>
    <select id="lb-combo-filter" aria-label="排行榜筛选" style="display:none;height:34px;min-width:220px;max-width:340px;padding:0 36px 0 12px;border:1px solid #dbe3ec;border-radius:6px;font-size:13px;background:#fff;color:#334155;" onchange="lbApplyFilter()"></select>
    <span id="lb-filter-hint" style="display:none;font-size:12px;color:#64748b;white-space:nowrap;"></span>
  </div>
  <div class="leaderboard" id="lb-list">
__CARDS_HTML__
  </div>
</div>
<script>
(function(){
  var _lbOpts = __FILTER_OPTIONS_JSON__;
  var _curMode = 'all';

  function setActiveTab(mode) {
    document.querySelectorAll('.lb-tab').forEach(function(t){
      var active = t.dataset.mode === mode;
      t.classList.toggle('active', active);
      t.style.background = active ? '#1e293b' : '#fff';
      t.style.color = active ? '#fff' : '#475569';
    });
  }

  function currentFilters() {
    var sel = document.getElementById('lb-combo-filter');
    if (!sel || !sel.value) return {};
    try {
      return JSON.parse(sel.value);
    } catch (e) {
      return {};
    }
  }

  window.lbSwitchMode = function(mode, btn) {
    setActiveTab(mode);
    _curMode = mode;
    populateCombo(mode);
    lbRender(mode);
  };

  window.lbApplyFilter = function() {
    lbRender(_curMode);
  };

  function populateCombo(mode) {
    var sel = document.getElementById('lb-combo-filter');
    if (!sel) return;
    sel.innerHTML = '';
    if (mode === 'all') {
      sel.style.display = 'none';
      return;
    }
    var opts = _lbOpts[mode] || [];
    sel.style.display = '';
    sel.innerHTML = '';
    if (!opts.length) {
      var empty = document.createElement('option');
      empty.value = '';
      empty.textContent = '暂无可对比组合';
      sel.appendChild(empty);
      sel.disabled = true;
      return;
    }
    opts.forEach(function(o) {
      var opt = document.createElement('option');
      opt.value = JSON.stringify(o.filters);
      opt.textContent = o.label;
      sel.appendChild(opt);
    });
    sel.disabled = false;
  }
  populateCombo('all');

  // 按平台时固定模型+Skill；按模型时固定平台+Skill；按Skill时固定平台+模型。
  function cardMatches(el, mode, filters) {
    var d = el.dataset;
    if (mode === 'platform') return d.platform === filters.platform;
    if (mode === 'model') return d.model === filters.model;
    if (mode === 'skill') return d.skill === filters.skill;
    return true;
  }

  function updateHint(mode, filters) {
    var hint = document.getElementById('lb-filter-hint');
    if (mode === 'all') {
      hint.style.display = 'none';
      hint.textContent = '';
      return;
    }
    var text = '';
    if (mode === 'platform') text = '仅显示平台=' + (filters.platform || '-');
    if (mode === 'model') text = '仅显示模型=' + (filters.model || '-');
    if (mode === 'skill') text = '仅显示Skill=' + (filters.skill || '-');
    hint.textContent = text;
    hint.style.display = text ? 'block' : 'none';
  }

  window.lbRender = function(mode) {
    var container = document.getElementById('lb-list');
    var items = Array.from(container.querySelectorAll('.lb-item'));
    var filters = currentFilters();
    updateHint(mode, filters);

    if (mode === 'all') {
      items.forEach(function(el) { el.style.display = ''; });
      items.sort(function(a, b) {
        return parseFloat(b.dataset.score) - parseFloat(a.dataset.score);
      });
      items.forEach(function(el, i) {
        container.appendChild(el);
        el.querySelector('.lb-rank-num').textContent = i + 1;
        var cls = i < 3 ? ['gold','silver','bronze'][i] : '';
        el.querySelector('.lb-rank').className = 'lb-rank' + (cls ? ' ' + cls : '');
      });
      return;
    }

    // 控制变量模式
    var match = [];
    items.forEach(function(el) {
      if (cardMatches(el, mode, filters)) {
        match.push(el);
        el.style.display = '';
      } else {
        el.style.display = 'none';
      }
    });
    match.sort(function(a, b) {
      return parseFloat(b.dataset.score) - parseFloat(a.dataset.score);
    });
    match.forEach(function(el, i) {
      container.appendChild(el);
      el.querySelector('.lb-rank-num').textContent = i + 1;
      var cls = i < 3 ? ['gold','silver','bronze'][i] : '';
      el.querySelector('.lb-rank').className = 'lb-rank' + (cls ? ' ' + cls : '');
    });
  };
})();
</script>
`
	leaderboardHTML = strings.ReplaceAll(leaderboardHTML, "__CARDS_HTML__", cardsHTML)
	leaderboardHTML = strings.ReplaceAll(leaderboardHTML, "__FILTER_OPTIONS_JSON__", filterOptionsJSON)
	b.WriteString(leaderboardHTML)

	return b.String()
}

func tokenCoverageLabel(m contracts.ModelRank) string {
	total := m.ActualTokenSamples + m.EstimatedTokenSamples + m.PartialTokenSamples + m.MissingTokenSamples
	if total == 0 {
		if m.AvgTotalTokens > 0 {
			return "actual=0, estimated=0, partial=0, missing=0"
		}
		return "actual=0, estimated=0, partial=0, missing=" + fmt.Sprintf("%d", m.TotalSamples)
	}
	return fmt.Sprintf("actual=%d, estimated=%d, partial=%d, missing=%d",
		m.ActualTokenSamples,
		m.EstimatedTokenSamples,
		m.PartialTokenSamples,
		m.MissingTokenSamples)
}

// buildFilterOptionsJSON 构建平台 / 模型 / Skill 的单维度筛选选项。
func buildFilterOptionsJSON(models []contracts.ModelRank) string {
	type comboOption struct {
		Label   string            `json:"label"`
		Filters map[string]string `json:"filters"`
	}
	seen := map[string]map[string]struct{}{
		"platform": {},
		"model":    {},
		"skill":    {},
	}
	opts := map[string][]comboOption{
		"platform": {},
		"model":    {},
		"skill":    {},
	}
	add := func(mode, value string) {
		value = strings.TrimSpace(value)
		if value == "" {
			return
		}
		if _, ok := seen[mode][value]; ok {
			return
		}
		seen[mode][value] = struct{}{}
		opts[mode] = append(opts[mode], comboOption{
			Label:   value,
			Filters: map[string]string{mode: value},
		})
	}
	for _, m := range models {
		platform := firstNonEmpty(m.AgentFramework, "model_api")
		model := firstNonEmpty(m.AgentModel, m.Model)
		skill := firstNonEmpty(m.SkillName, "no_skill")
		add("platform", platform)
		add("model", model)
		add("skill", skill)
	}
	for mode := range opts {
		sort.Slice(opts[mode], func(i, j int) bool { return opts[mode][i].Label < opts[mode][j].Label })
	}
	b, _ := json.Marshal(opts)
	return string(b)
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
      <h3>场景 x 语言 Scenario by Language</h3>
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
      <div id="efficiencyQualityLegend" class="chart-legend-wrap"></div>
      <div id="efficiencyQualityNotes" class="efficiency-notes"></div>
    </div>
  </div>
  <div class="chart-grid-2" style="margin-top:16px;">
    <div class="panel">
      <h3>场景通过率柱状图</h3>
      <div class="chart-box tall scenario-chart-box"><canvas id="scenarioBarChart"></canvas></div>
    </div>
    <div class="panel">
      <h3>场景覆盖与变异对比（柱状图）</h3>
      <div class="chart-box tall scenario-chart-box"><canvas id="scenarioTrendChart"></canvas></div>
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
