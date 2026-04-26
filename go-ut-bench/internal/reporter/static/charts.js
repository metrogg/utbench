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

const reportTopModels = __TOP_MODELS_JSON_PLACEHOLDER__;
const evaluationRows = __ROWS_JSON_PLACEHOLDER__;

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
  return '<div class="metric"><div class="bar"><span class="' + fillClass + '" style="width:' + width + '%"></span></div><span class="val">' + width + '%</span></div>';
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
  const scenarioBarCanvas = document.getElementById('scenarioBarChart');
  const scenarioTrendCanvas = document.getElementById('scenarioTrendChart');
  if (!scenarioBarCanvas || !scenarioTrendCanvas) return;

  const labels = items.map(item => item.scenario + ' / ' + item.language.toUpperCase());
  const compileRates = items.map(item => item.compilePassRate);
  const testRates = items.map(item => item.testPassRate);
  const mutationRates = items.map(item => item.mutationScore);

  if (scenarioBarChart) scenarioBarChart.destroy();
  scenarioBarChart = new Chart(scenarioBarCanvas, {
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
      scales: { y: { beginAtZero: true, max: 1, ticks: { callback: value => Math.round(value * 100) + '%' } } }
    }
  });

  if (scenarioTrendChart) scenarioTrendChart.destroy();
  scenarioTrendChart = new Chart(scenarioTrendCanvas, {
    type: 'line',
    data: {
      labels,
      datasets: [
        { label: '行覆盖率', data: items.map(item => item.lineCoverage), borderColor: '#3b82f6', backgroundColor: 'rgba(59,130,246,.12)', tension: .3, fill: true },
        { label: '变异分数', data: mutationRates, borderColor: '#f59e0b', backgroundColor: 'rgba(245,158,11,.12)', tension: .3, fill: true }
      ]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { position: 'top' } },
      scales: { y: { beginAtZero: true, max: 1, ticks: { callback: value => Math.round(value * 100) + '%' } } }
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
      const bg = 'hsla(' + hue + ', 75%, 85%, 1)';
      html += '<div class="heatmap-cell" style="background:' + bg + ';">' + safeText(row.sample_id) + '<br>' + Math.round(value * 100) + '%</div>';
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
  // 渲染综合排名主图
  renderOverallChart();
  // 渲染雷达图
  renderRadarChart();
  // 渲染维度分解
  renderDimensionBreakdown();
  // 渲染热力图
  const aggregated = aggregateRows('all', 'all');
  renderHeatmap('coverageHeatmap', aggregated.filtered, 'line_coverage');
  renderHeatmap('mutationHeatmap', aggregated.filtered, 'mutation_score');
}

// ========== 图表切换功能 ==========
let currentChartTab = 'overall';
let mainCompareChart = null;

function switchChartTab(tab) {
  currentChartTab = tab;
  // 更新 tab 样式
  document.querySelectorAll('.chart-tab').forEach(btn => {
    if (btn.dataset.tab === tab) {
      btn.style.background = '#1e40af';
      btn.style.color = '#fff';
      btn.classList.add('active');
    } else {
      btn.style.background = '#f1f5f9';
      btn.style.color = '#475569';
      btn.classList.remove('active');
    }
  });
  // 显示对应的筛选面板
  document.querySelectorAll('.filter-panel').forEach(panel => panel.style.display = 'none');
  const filterPanel = document.getElementById('filter-' + tab);
  if (filterPanel) filterPanel.style.display = 'block';
  // 刷新图表
  if (tab === 'overall') renderOverallChart();
  else if (tab === 'language') refreshLanguageChart();
  else if (tab === 'scenario') refreshScenarioChart();
  else if (tab === 'model') refreshModelDetailChart();
}

function renderOverallChart() {
  const titleEl = document.getElementById('main-chart-title');
  if (titleEl) titleEl.textContent = '模型综合得分对比';

  const modelNames = reportTopModels.map(item => item.model);
  const scores = reportTopModels.map(item => item.composite_score);
  const colors = reportTopModels.map((item, idx) => {
    if (item.rank === 1) return '#22c55e'; // green for #1
    if (item.rank === 2) return '#3b82f6'; // blue for #2
    if (item.rank === 3) return '#f59e0b'; // amber for #3
    return '#94a3b8'; // gray for others
  });

  if (mainCompareChart) mainCompareChart.destroy();
  mainCompareChart = new Chart(document.getElementById('mainCompareChart'), {
    type: 'bar',
    data: {
      labels: modelNames,
      datasets: [{
        label: '综合得分',
        data: scores,
        backgroundColor: colors,
        borderRadius: 6
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: {
        legend: { display: false },
        tooltip: {
          callbacks: {
            label: ctx => '综合得分: ' + (ctx.raw * 100).toFixed(1) + '%'
          }
        }
      },
      scales: {
        y: { beginAtZero: true, max: 1, ticks: { callback: v => Math.round(v * 100) + '%' } },
        x: { }
      }
    }
  });

  // 更新维度分解
  renderDimensionBreakdownOverall();
}

function renderDimensionBreakdownOverall() {
  const tbody = document.getElementById('dimensionBreakdownBody');
  if (!tbody) return;
  let html = '';
  reportTopModels.forEach(item => {
    const rankBadge = item.rank <= 3 ? '<span style="background:#1e40af;color:#fff;padding:2px 6px;border-radius:4px;font-size:11px;margin-right:4px;">#' + item.rank + '</span>' : '';
    html += '<tr><td style="padding:8px;border-bottom:1px solid #e2e8f0;">' + rankBadge + safeText(item.model) + '</td>';
    html += '<td style="padding:8px;border-bottom:1px solid #e2e8f0;text-align:center;">' + metricCell(item.compile_pass_rate) + '</td>';
    html += '<td style="padding:8px;border-bottom:1px solid #e2e8f0;text-align:center;">' + metricCell(item.avg_test_pass_rate) + '</td>';
    html += '<td style="padding:8px;border-bottom:1px solid #e2e8f0;text-align:center;">' + metricCell(item.avg_line_coverage) + '</td>';
    html += '<td style="padding:8px;border-bottom:1px solid #e2e8f0;text-align:center;">' + metricCell(item.avg_mutation_score) + '</td>';
    html += '<td style="padding:8px;border-bottom:1px solid #e2e8f0;text-align:center;font-weight:700;">' + Math.round(item.composite_score * 100) + '%</td></tr>';
  });
  tbody.innerHTML = html;
}

function refreshLanguageChart() {
  const langSelect = document.getElementById('chart-language-select');
  const lang = langSelect ? langSelect.value : 'all';

  const titleEl = document.getElementById('main-chart-title');
  if (titleEl) {
    titleEl.textContent = lang === 'all' ? '各语言下模型综合得分对比' : safeText(lang.toUpperCase()) + ' 语言下模型得分对比';
  }

  // 按语言筛选数据
  const filteredRows = evaluationRows.filter(row => lang === 'all' || row.language === lang);
  const byModel = new Map();
  filteredRows.forEach(row => {
    if (!byModel.has(row.model)) byModel.set(row.model, { compile: 0, test: 0, line: 0, mutation: 0, count: 0 });
    const agg = byModel.get(row.model);
    agg.count++;
    if (row.compile_pass) agg.compile++;
    if (row.test_pass) agg.test++;
    if (row.line_coverage) agg.line += row.line_coverage;
    if (row.mutation_score) agg.mutation += row.mutation_score;
  });

  const data = Array.from(byModel.entries()).map(([model, agg]) => ({
    model,
    compileRate: agg.count ? agg.compile / agg.count : 0,
    testRate: agg.count ? agg.test / agg.count : 0,
    lineRate: agg.line / agg.count,
    mutationRate: agg.mutation / agg.count,
    composite: agg.count ? (agg.compile/agg.count)*0.3 + (agg.test/agg.count)*0.3 + (agg.line/agg.count)*0.2 + (agg.mutation/agg.count)*0.2 : 0
  })).sort((a, b) => b.composite - a.composite);

  if (mainCompareChart) mainCompareChart.destroy();
  mainCompareChart = new Chart(document.getElementById('mainCompareChart'), {
    type: 'bar',
    data: {
      labels: data.map(d => d.model),
      datasets: [{
        label: '综合得分',
        data: data.map(d => d.composite),
        backgroundColor: data.map((_, idx) => idx === 0 ? '#22c55e' : idx === 1 ? '#3b82f6' : idx === 2 ? '#f59e0b' : '#94a3b8'),
        borderRadius: 6
      }]
    },
    options: {
      
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false }, tooltip: { callbacks: { label: ctx => '综合得分: ' + (ctx.raw * 100).toFixed(1) + '%' } } },
      scales: { x: { beginAtZero: true, max: 1, ticks: { callback: v => Math.round(v * 100) + '%' } } }
    }
  });

  // 更新维度分解
  renderDimensionBreakdownFiltered(data, '语言: ' + (lang === 'all' ? '全部' : lang.toUpperCase()));
}

function refreshScenarioChart() {
  const scenarioSelect = document.getElementById('chart-scenario-select');
  const scenario = scenarioSelect ? scenarioSelect.value : 'all';

  const titleEl = document.getElementById('main-chart-title');
  if (titleEl) {
    titleEl.textContent = scenario === 'all' ? '各场景下模型综合得分对比' : safeText(scenario) + ' 场景下模型得分对比';
  }

  // 按场景筛选数据
  const filteredRows = evaluationRows.filter(row => {
    const rowScenario = getScenarioFromSample(row.sample_id);
    return scenario === 'all' || rowScenario === scenario;
  });
  const byModel = new Map();
  filteredRows.forEach(row => {
    if (!byModel.has(row.model)) byModel.set(row.model, { compile: 0, test: 0, line: 0, mutation: 0, count: 0 });
    const agg = byModel.get(row.model);
    agg.count++;
    if (row.compile_pass) agg.compile++;
    if (row.test_pass) agg.test++;
    if (row.line_coverage) agg.line += row.line_coverage;
    if (row.mutation_score) agg.mutation += row.mutation_score;
  });

  const data = Array.from(byModel.entries()).map(([model, agg]) => ({
    model,
    compileRate: agg.count ? agg.compile / agg.count : 0,
    testRate: agg.count ? agg.test / agg.count : 0,
    lineRate: agg.line / agg.count,
    mutationRate: agg.mutation / agg.count,
    composite: agg.count ? (agg.compile/agg.count)*0.3 + (agg.test/agg.count)*0.3 + (agg.line/agg.count)*0.2 + (agg.mutation/agg.count)*0.2 : 0,
    samples: agg.count
  })).sort((a, b) => b.composite - a.composite);

  if (mainCompareChart) mainCompareChart.destroy();
  mainCompareChart = new Chart(document.getElementById('mainCompareChart'), {
    type: 'bar',
    data: {
      labels: data.map(d => d.model),
      datasets: [{
        label: '综合得分',
        data: data.map(d => d.composite),
        backgroundColor: data.map((_, idx) => idx === 0 ? '#22c55e' : idx === 1 ? '#3b82f6' : idx === 2 ? '#f59e0b' : '#94a3b8'),
        borderRadius: 6
      }]
    },
    options: {
      
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false }, tooltip: { callbacks: { label: ctx => '综合得分: ' + (ctx.raw * 100).toFixed(1) + '%' } } },
      scales: { x: { beginAtZero: true, max: 1, ticks: { callback: v => Math.round(v * 100) + '%' } } }
    }
  });

  renderDimensionBreakdownFiltered(data, '场景: ' + (scenario === 'all' ? '全部' : scenario));
}

function refreshModelDetailChart() {
  const modelSelect = document.getElementById('chart-model-select');
  const model = modelSelect ? modelSelect.value : reportTopModels[0]?.model;

  const titleEl = document.getElementById('main-chart-title');
  if (titleEl) titleEl.textContent = safeText(model) + ' 模型各维度得分';

  // 获取该模型的数据
  const modelData = reportTopModels.find(m => m.model === model);
  if (!modelData) return;

  if (mainCompareChart) mainCompareChart.destroy();
  mainCompareChart = new Chart(document.getElementById('mainCompareChart'), {
    type: 'bar',
    data: {
      labels: ['编译通过率', '样本测试通过率', '行覆盖率', '变异分数'],
      datasets: [{
        label: model,
        data: [modelData.compile_pass_rate, modelData.avg_test_pass_rate, modelData.avg_line_coverage, modelData.avg_mutation_score],
        backgroundColor: ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6'],
        borderRadius: 6
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      plugins: { legend: { display: false } },
      scales: { y: { beginAtZero: true, max: 1, ticks: { callback: v => Math.round(v * 100) + '%' } } }
    }
  });

  // 更新维度分解 - 显示该模型在各语言/场景下的表现
  renderModelDetailBreakdown(model);
}

function renderModelDetailBreakdown(model) {
  const container = document.getElementById('dimensionBreakdown');
  if (!container) return;

  // 按语言统计
  const byLang = new Map();
  evaluationRows.filter(row => row.model === model).forEach(row => {
    if (!byLang.has(row.language)) byLang.set(row.language, { compile: 0, test: 0, line: 0, mutation: 0, count: 0 });
    const agg = byLang.get(row.language);
    agg.count++;
    if (row.compile_pass) agg.compile++;
    if (row.test_pass) agg.test++;
    if (row.line_coverage) agg.line += row.line_coverage;
    if (row.mutation_score) agg.mutation += row.mutation_score;
  });

  let html = '<div style="font-size:13px;"><strong>' + safeText(model) + ' 模型各语言表现</strong><table style="width:100%;margin-top:8px;border-collapse:collapse;">';
  html += '<thead><tr style="background:#f1f5f9;"><th style="padding:6px;text-align:left;">语言</th><th style="padding:6px;">样本</th><th style="padding:6px;">编译</th><th style="padding:6px;">测试</th><th style="padding:6px;">覆盖</th><th style="padding:6px;">变异</th></tr></thead><tbody>';
  Array.from(byLang.entries()).forEach(([lang, agg]) => {
    html += '<tr><td style="padding:6px;border-bottom:1px solid #e2e8f0;">' + safeText(lang.toUpperCase()) + '</td>';
    html += '<td style="padding:6px;border-bottom:1px solid #e2e8f0;text-align:center;">' + agg.count + '</td>';
    html += '<td style="padding:6px;border-bottom:1px solid #e2e8f0;text-align:center;">' + metricCell(agg.compile/agg.count) + '</td>';
    html += '<td style="padding:6px;border-bottom:1px solid #e2e8f0;text-align:center;">' + metricCell(agg.test/agg.count) + '</td>';
    html += '<td style="padding:6px;border-bottom:1px solid #e2e8f0;text-align:center;">' + metricCell(agg.line/agg.count) + '</td>';
    html += '<td style="padding:6px;border-bottom:1px solid #e2e8f0;text-align:center;">' + metricCell(agg.mutation/agg.count) + '</td></tr>';
  });
  html += '</tbody></table></div>';
  container.innerHTML = html;
}

function renderRadarChart() {
  if (radarChart) radarChart.destroy();
  radarChart = new Chart(document.getElementById('radarChart'), {
    type: 'radar',
    data: {
      labels: ['编译通过率', '样本测试通过率', '行覆盖率', '变异分数'],
      datasets: reportTopModels.slice(0, 5).map((item, index) => ({
        label: item.model,
        data: [item.compile_pass_rate, item.avg_test_pass_rate, item.avg_line_coverage, item.avg_mutation_score],
        fill: true,
        backgroundColor: ['rgba(59,130,246,0.18)', 'rgba(16,185,129,0.18)', 'rgba(245,158,11,0.18)', 'rgba(139,92,246,0.18)', 'rgba(20,184,166,0.18)'][index % 5],
        borderColor: ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6', '#14b8a6'][index % 5],
        pointBackgroundColor: ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6', '#14b8a6'][index % 5]
      }))
    },
    options: { responsive: true, maintainAspectRatio: false, scales: { r: { beginAtZero: true, max: 1, ticks: { callback: value => Math.round(value * 100) + '%' } } } }
  });
}

function renderDimensionBreakdown() {
  renderDimensionBreakdownOverall();
}

function renderFilteredSections() {
  // 渲染表格部分
  const aggregated = aggregateRows('all', 'all');
  renderLanguageTable(aggregated.languages);
  renderScenarioTable(aggregated.scenarios);
  renderErrorTable(aggregated.failureRows);
  // 渲染错误分布图表
  renderErrorCharts(aggregated.failureRows);
}

(async function init() {
  try {
    if (window.__utBenchEnsureChartJS) {
      await window.__utBenchEnsureChartJS();
    }
    renderModelCharts();
    renderFilteredSections();
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

// ========== 原始数据筛选功能 ==========
const RAW_DATA_DEFAULT_LIMIT = 20;

function initRawDataFilters() {
  const tbody = document.getElementById('raw-data-body');
  if (!tbody) return;
  const rows = tbody.querySelectorAll('tr');
  let shown = 0;
  rows.forEach((row, idx) => {
    if (idx < RAW_DATA_DEFAULT_LIMIT) {
      row.style.display = '';
      shown++;
    } else {
      row.style.display = 'none';
    }
  });
  updateFilterResultCount(shown, rows.length);
}

function applyRawDataFilters() {
  const modelFilter = document.getElementById('filter-model');
  const langFilter = document.getElementById('filter-language');
  const scenarioFilter = document.getElementById('filter-scenario');
  const statusFilter = document.getElementById('filter-status');
  const searchInput = document.getElementById('filter-search');

  const modelVal = modelFilter ? modelFilter.value : '';
  const langVal = langFilter ? langFilter.value : '';
  const scenarioVal = scenarioFilter ? scenarioFilter.value : '';
  const statusVal = statusFilter ? statusFilter.value : '';
  const searchVal = searchInput ? searchInput.value.trim().toLowerCase() : '';

  const tbody = document.getElementById('raw-data-body');
  if (!tbody) return;
  const rows = tbody.querySelectorAll('tr');

  let matched = 0;
  rows.forEach((row) => {
    const rowModel = row.getAttribute('data-model') || '';
    const rowLang = row.getAttribute('data-language') || '';
    const rowScenario = row.getAttribute('data-scenario') || '';
    const rowSample = row.getAttribute('data-sample') || '';
    const rowCompilePass = row.getAttribute('data-compile-pass') || '';
    const rowTestPass = row.getAttribute('data-test-pass') || '';
    const rowMutationZero = row.getAttribute('data-mutation-zero') || '';
    const rowAllPass = row.getAttribute('data-all-pass') || '';
    const rowHasFail = row.getAttribute('data-has-fail') || '';

    // 状态筛选逻辑
    let statusMatch = true;
    if (statusVal === 'pass') {
      statusMatch = rowAllPass === 'true';
    } else if (statusVal === 'fail') {
      statusMatch = rowHasFail === 'true';
    } else if (statusVal === 'compile_fail') {
      statusMatch = rowCompilePass === 'false';
    } else if (statusVal === 'test_fail') {
      statusMatch = rowCompilePass === 'true' && rowTestPass === 'false';
    } else if (statusVal === 'mutation_zero') {
      statusMatch = rowMutationZero === 'true';
    }

    // 模型、语言、场景、搜索筛选
    const modelMatch = !modelVal || rowModel === modelVal;
    const langMatch = !langVal || rowLang === langVal;
    const scenarioMatch = !scenarioVal || rowScenario === scenarioVal;
    const searchMatch = !searchVal || rowSample.toLowerCase().includes(searchVal);

    if (modelMatch && langMatch && scenarioMatch && statusMatch && searchMatch) {
      matched++;
    }
  });

  // 显示匹配的行（最多 RAW_DATA_DEFAULT_LIMIT 条）
  let shown = 0;
  rows.forEach((row) => {
    const rowModel = row.getAttribute('data-model') || '';
    const rowLang = row.getAttribute('data-language') || '';
    const rowScenario = row.getAttribute('data-scenario') || '';
    const rowSample = row.getAttribute('data-sample') || '';
    const rowCompilePass = row.getAttribute('data-compile-pass') || '';
    const rowTestPass = row.getAttribute('data-test-pass') || '';
    const rowMutationZero = row.getAttribute('data-mutation-zero') || '';
    const rowAllPass = row.getAttribute('data-all-pass') || '';
    const rowHasFail = row.getAttribute('data-has-fail') || '';

    let statusMatch = true;
    if (statusVal === 'pass') {
      statusMatch = rowAllPass === 'true';
    } else if (statusVal === 'fail') {
      statusMatch = rowHasFail === 'true';
    } else if (statusVal === 'compile_fail') {
      statusMatch = rowCompilePass === 'false';
    } else if (statusVal === 'test_fail') {
      statusMatch = rowCompilePass === 'true' && rowTestPass === 'false';
    } else if (statusVal === 'mutation_zero') {
      statusMatch = rowMutationZero === 'true';
    }

    const modelMatch = !modelVal || rowModel === modelVal;
    const langMatch = !langVal || rowLang === langVal;
    const scenarioMatch = !scenarioVal || rowScenario === scenarioVal;
    const searchMatch = !searchVal || rowSample.toLowerCase().includes(searchVal);

    if (modelMatch && langMatch && scenarioMatch && statusMatch && searchMatch && shown < RAW_DATA_DEFAULT_LIMIT) {
      row.style.display = '';
      shown++;
    } else {
      row.style.display = 'none';
    }
  });

  updateFilterResultCount(shown, matched);
}

function resetRawDataFilters() {
  const modelFilter = document.getElementById('filter-model');
  const langFilter = document.getElementById('filter-language');
  const scenarioFilter = document.getElementById('filter-scenario');
  const statusFilter = document.getElementById('filter-status');
  const searchInput = document.getElementById('filter-search');

  if (modelFilter) modelFilter.value = '';
  if (langFilter) langFilter.value = '';
  if (scenarioFilter) scenarioFilter.value = '';
  if (statusFilter) statusFilter.value = '';
  if (searchInput) searchInput.value = '';

  initRawDataFilters();
}

function updateFilterResultCount(shown, total) {
  const countEl = document.getElementById('filter-result-count');
  if (countEl) {
    countEl.textContent = '显示 ' + shown + ' / ' + total + ' 条';
  }
}
