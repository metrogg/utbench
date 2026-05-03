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
let efficiencyQualityChart;
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
  if (value === null || value === undefined || Number.isNaN(Number(value))) return '<span class="badge">-</span>';
  const width = Math.max(0, Math.min(100, Math.round(value * 100)));
  const fillClass = value >= 0.7 ? 'ok' : 'bad';
  return '<div class="metric"><div class="bar"><span class="' + fillClass + '" style="width:' + width + '%"></span></div><span class="val">' + width + '%</span></div>';
}

function getScenarioFromSample(sampleID) {
  if (!sampleID) return 'unknown';
  for (const prefix of ['complex_dependency', 'interface_mock', 'simple_function', 'boundary']) {
    if (sampleID === prefix || sampleID.startsWith(prefix + '_')) return prefix;
  }
  const idx = sampleID.indexOf('_');
  return idx > 0 ? sampleID.slice(0, idx) : sampleID;
}

function scenarioLabel(value) {
  const labels = {
    boundary: '边界值 / Boundary',
    simple_function: '简单函数 / Simple',
    complex_dependency: '复杂依赖 / Complex',
    interface_mock: '接口 Mock / Interface',
    unknown: '未知 / Unknown'
  };
  return labels[value] || value || '未知 / Unknown';
}

function stageLabel(value) {
  const labels = {
    generate: '生成 / Generate',
    compile: '编译 / Compile',
    test: '测试 / Test',
    coverage: '覆盖率 / Coverage',
    mutation: '变异测试 / Mutation'
  };
  return labels[value] || value || '-';
}

function errorTypeLabel(value) {
  const labels = {
    truncated: '输出截断 / Truncated',
    module_not_found: '模块缺失 / Module not found',
    name_error: '名称错误 / Name error',
    assertion_failure: '断言失败 / Assertion failure',
    syntax_error: '语法错误 / Syntax error',
    indentation_error: '缩进错误 / Indentation',
    timeout: '超时 / Timeout',
    permission_error: '权限错误 / Permission',
    mutation_skipped_baseline_failed: '变异前测试失败 / Mutation baseline failed',
    mutation_target_not_exercised: '未覆盖变异目标 / Target not exercised',
    mutation_no_results: '无变异结果 / No mutation results',
    mutation_no_coverage: '无覆盖数据 / No coverage',
    mutation_no_effective_mutants: '无有效变异体 / No effective mutants',
    mutation_timeout: '变异超时 / Mutation timeout',
    mutation_tool_error: '变异工具错误 / Mutation tool error',
    mutation_error: '变异错误 / Mutation error',
    coverage_error: '覆盖率错误 / Coverage error',
    other: '其他 / Other'
  };
  return labels[value] || value || '-';
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

function filterMatch(value, selectedValues) {
  return !selectedValues || selectedValues.length === 0 || selectedValues.includes(value);
}

function isScoreEligibleRow(row) {
  return row.score_eligible === undefined || row.score_eligible === null || row.score_eligible === true;
}

function createAgg(key, extra = {}) {
  return Object.assign({ key, total: 0, compilePass: 0, testPass: 0, testTotal: 0, lineSum: 0, lineCnt: 0, branchSum: 0, branchCnt: 0, mutationSum: 0, mutationCnt: 0 }, extra);
}

function aggRow(agg, row) {
  agg.total += 1;
  if (row.compile_pass) agg.compilePass += 1;
  if (row.test_pass_count !== null && row.test_pass_count !== undefined && row.test_total_count !== null && row.test_total_count !== undefined) {
    agg.testPass += row.test_pass_count;
    agg.testTotal += row.test_total_count;
  } else if (row.test_pass !== null && row.test_pass !== undefined) {
    agg.testTotal += 1;
    if (row.test_pass) agg.testPass += 1;
  }
  if (row.line_coverage !== null && row.line_coverage !== undefined) { agg.lineSum += row.line_coverage; agg.lineCnt += 1; }
  if (row.branch_coverage !== null && row.branch_coverage !== undefined) { agg.branchSum += row.branch_coverage; agg.branchCnt += 1; }
  if (row.mutation_score !== null && row.mutation_score !== undefined) { agg.mutationSum += row.mutation_score; agg.mutationCnt += 1; }
}

function aggToMetrics(item) {
  return {
    total: item.total,
    compilePassRate: item.total ? item.compilePass / item.total : 0,
    testPassRate: item.testTotal ? item.testPass / item.testTotal : 0,
    lineCoverage: item.lineCnt ? item.lineSum / item.lineCnt : 0,
    branchCoverage: item.branchCnt ? item.branchSum / item.branchCnt : 0,
    mutationScore: item.mutationCnt ? item.mutationSum / item.mutationCnt : 0
  };
}

function aggregateRows(selectedModels, selectedLanguages, selectedScenarios) {
  const filtered = evaluationRows.filter(row => {
    const scenario = getScenarioFromSample(row.sample_id);
    return isScoreEligibleRow(row) &&
      filterMatch(row.model || '', selectedModels) &&
      filterMatch(row.language || 'unknown', selectedLanguages) &&
      filterMatch(scenario, selectedScenarios);
  });
  const byModel = new Map();
  const byLanguage = new Map();
  const byScenario = new Map();
  const failures = new Map();

  for (const row of filtered) {
    const modelKey = row.model || 'unknown';
    if (!byModel.has(modelKey)) byModel.set(modelKey, createAgg(modelKey, { model: modelKey }));
    aggRow(byModel.get(modelKey), row);

    const langKey = row.language || 'unknown';
    if (!byLanguage.has(langKey)) byLanguage.set(langKey, createAgg(langKey, { language: langKey }));
    aggRow(byLanguage.get(langKey), row);

    const scenario = getScenarioFromSample(row.sample_id);
    const scenKey = langKey + '|' + scenario;
    if (!byScenario.has(scenKey)) byScenario.set(scenKey, createAgg(scenKey, { scenario, language: langKey }));
    aggRow(byScenario.get(scenKey), row);

    const stage = getStageFromRow(row);
    const errorType = getErrorTypeFromRow(row);
    if (stage && errorType) {
      const key = stage + '|' + errorType;
      if (!failures.has(key)) failures.set(key, { stage, errorType, count: 0, exampleModel: row.model || '', exampleSample: row.sample_id || '' });
      failures.get(key).count += 1;
    }
  }

  const models = Array.from(byModel.values()).map(item => Object.assign({ model: item.model }, aggToMetrics(item))).sort((a, b) => b.mutationScore - a.mutationScore || b.testPassRate - a.testPassRate);
  const languages = Array.from(byLanguage.values()).map(item => Object.assign({ language: item.language }, aggToMetrics(item))).sort((a, b) => a.language.localeCompare(b.language));
  const scenarios = Array.from(byScenario.values()).map(item => Object.assign({ scenario: item.scenario, language: item.language }, aggToMetrics(item))).sort((a, b) => (a.language + a.scenario).localeCompare(b.language + b.scenario));

  const failureRows = Array.from(failures.values()).sort((a, b) => b.count - a.count);
  return { models, languages, scenarios, failureRows, filtered };
}

function renderModelTable(items) {
  const body = document.getElementById('by-model-body');
  const empty = document.getElementById('by-model-empty');
  if (!body || !empty) return;
  body.innerHTML = items.map(item => '<tr><td><strong>' + safeText(item.model) + '</strong></td><td>' + item.total + '</td><td>' + metricCell(item.compilePassRate) + '</td><td>' + metricCell(item.testPassRate) + '</td><td>' + metricCell(item.lineCoverage) + '</td><td>' + metricCell(item.branchCoverage) + '</td><td>' + metricCell(item.mutationScore) + '</td></tr>').join('');
  empty.style.display = items.length ? 'none' : 'block';
}

function renderLanguageTable(items) {
  const body = document.getElementById('by-language-body');
  const empty = document.getElementById('by-language-empty');
  if (!body || !empty) return;
  body.innerHTML = items.map(item => '<tr><td><strong>' + safeText(String(item.language).toUpperCase()) + '</strong></td><td>' + item.total + '</td><td>' + metricCell(item.compilePassRate) + '</td><td>' + metricCell(item.testPassRate) + '</td><td>' + metricCell(item.lineCoverage) + '</td><td>' + metricCell(item.branchCoverage) + '</td><td>' + metricCell(item.mutationScore) + '</td></tr>').join('');
  empty.style.display = items.length ? 'none' : 'block';
}

function renderScenarioTable(items) {
  const body = document.getElementById('by-scenario-body');
  const empty = document.getElementById('by-scenario-empty');
  if (!body || !empty) return;
  body.innerHTML = items.map(item => '<tr><td><strong>' + safeText(scenarioLabel(item.scenario)) + '</strong></td><td>' + safeText(String(item.language).toUpperCase()) + '</td><td>' + item.total + '</td><td>' + metricCell(item.compilePassRate) + '</td><td>' + metricCell(item.testPassRate) + '</td><td>' + metricCell(item.lineCoverage) + '</td><td>' + metricCell(item.branchCoverage) + '</td><td>' + metricCell(item.mutationScore) + '</td></tr>').join('');
  empty.style.display = items.length ? 'none' : 'block';
}

function renderErrorTable(items) {
  const body = document.getElementById('error-analysis-body');
  const empty = document.getElementById('error-analysis-empty');
  body.innerHTML = items.map(item => '<tr><td>' + safeText(stageLabel(item.stage)) + '</td><td>' + safeText(errorTypeLabel(item.errorType)) + '</td><td>' + item.count + '</td><td>' + safeText(item.exampleModel) + '</td><td>' + safeText(item.exampleSample) + '</td></tr>').join('');
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
  const typeLabels = typeData.labels.length ? typeData.labels.map(errorTypeLabel) : ['无错误 / No Errors'];
  const stageLabels = stageData.labels.length ? stageData.labels.map(stageLabel) : ['无错误 / No Errors'];
  errorTypeChart = upsertChart(errorTypeChart, 'errorTypeChart', 'doughnut', typeLabels, typeData.values.length ? typeData.values : [1], ['#c43d2f', '#d5902f', '#2f7d72', '#315f9c', '#6b7280', '#111827']);
  stageChart = upsertChart(stageChart, 'stageChart', 'pie', stageLabels, stageData.values.length ? stageData.values : [1], ['#c43d2f', '#d5902f', '#2f7d72', '#315f9c', '#6b7280']);
}

function renderScenarioCharts(items) {
  const labels = items.map(item => scenarioLabel(item.scenario) + ' / ' + item.language.toUpperCase());
  const compileRates = items.map(item => item.compilePassRate);
  const testRates = items.map(item => item.testPassRate);
  const mutationRates = items.map(item => item.mutationScore);

  if (scenarioBarChart) scenarioBarChart.destroy();
  scenarioBarChart = new Chart(document.getElementById('scenarioBarChart'), {
    type: 'bar',
    data: {
      labels,
      datasets: [
        { label: '编译通过率', data: compileRates, backgroundColor: '#315f9c' },
        { label: '测试通过率', data: testRates, backgroundColor: '#2f7d72' }
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
  scenarioTrendChart = new Chart(document.getElementById('scenarioTrendChart'), {
    type: 'bar',
    data: {
      labels,
      datasets: [
        { label: '行覆盖率', data: items.map(item => item.lineCoverage), borderColor: '#315f9c', backgroundColor: 'rgba(49,95,156,.28)' },
        { label: '变异分数', data: mutationRates, borderColor: '#d5902f', backgroundColor: 'rgba(213,144,47,.32)' }
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

function initRawColumnToggles() {
  document.querySelectorAll('[data-raw-column]').forEach(input => {
    const apply = () => {
      const name = input.getAttribute('data-raw-column');
      document.querySelectorAll('.raw-col-' + name).forEach(cell => {
        cell.style.display = input.checked ? 'table-cell' : 'none';
      });
    };
    input.addEventListener('change', apply);
    apply();
  });
}

function normalizedAssertionDensity(value) {
  const raw = Number(value || 0);
  return Math.max(0, Math.min(1, raw / 5));
}

function qualityScore(item) {
  return (
    Number(item.compile_pass_rate || 0) * 0.25 +
    Number(item.avg_test_pass_rate || 0) * 0.30 +
    Number(item.avg_line_coverage || 0) * 0.15 +
    Number(item.avg_mutation_score || 0) * 0.25 +
    normalizedAssertionDensity(item.avg_assertion_density) * 0.05
  ) * 100;
}

function modelColor(index) {
  return ['#315f9c', '#2f7d72', '#c43d2f', '#d5902f', '#5b6472', '#0f766e', '#7c5f35', '#475569', '#4d7c0f', '#9f4f3b', '#111827', '#287f9c'][index % 12];
}

function shortModelName(value) {
  return String(value || '');
}

function renderEfficiencyQualityChart(axis = 'latency') {
  const canvas = document.getElementById('efficiencyQualityChart');
  if (!canvas) return;

  const axisIsTokens = axis === 'tokens';
  const denseMode = reportTopModels.length > 8;
  const points = reportTopModels.map((item, index) => {
    const xRaw = axisIsTokens ? Number(item.avg_total_tokens || 0) : Number(item.avg_latency_ms || 0) / 1000;
    const yRaw = qualityScore(item);
    const sampleSize = denseMode
      ? Math.max(4, Math.min(8, 4 + Math.log2(Number(item.total_samples || 1) + 1) * 0.6))
      : Math.max(5, Math.min(10, 5 + Math.log2(Number(item.total_samples || 1) + 1) * 0.8));
    return {
      label: item.model,
      data: [{ x: xRaw, y: yRaw }],
      pointRadius: sampleSize,
      pointHoverRadius: sampleSize + 2,
      pointHitRadius: 10,
      backgroundColor: modelColor(index) + (denseMode ? '99' : 'bb'),
      borderColor: modelColor(index),
      borderWidth: 1.5
    };
  });

  if (efficiencyQualityChart) efficiencyQualityChart.destroy();
  efficiencyQualityChart = new Chart(canvas, {
    type: 'scatter',
    data: { datasets: points },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      layout: { padding: { top: 18, right: 28, bottom: 10, left: 8 } },
      interaction: { mode: 'nearest', intersect: false },
      plugins: {
        legend: {
          position: 'bottom',
          labels: { boxWidth: 10, boxHeight: 10, usePointStyle: true }
        },
        tooltip: {
          callbacks: {
            label: ctx => {
              const item = reportTopModels[ctx.datasetIndex] || {};
              const x = ctx.parsed.x;
              const unit = axisIsTokens ? ' tokens' : 's';
              return [
                item.model,
                '质量分: ' + ctx.parsed.y.toFixed(1),
                (axisIsTokens ? '平均 Token: ' : '平均耗时: ') + x.toFixed(axisIsTokens ? 0 : 1) + unit,
                '样本数: ' + safeText(item.total_samples)
              ];
            }
          }
        }
      },
      scales: {
        x: {
          beginAtZero: true,
          title: { display: true, text: axisIsTokens ? '平均 Token Avg Total Tokens（越低越省）' : '平均耗时 Avg Latency（秒，越低越快）' },
          ticks: { callback: value => axisIsTokens ? Math.round(value) : Number(value).toFixed(0) + 's' },
          grid: { color: '#d8e0e7', tickLength: 0 },
          border: { color: '#7d8793', width: 2 }
        },
        y: {
          beginAtZero: true,
          suggestedMax: 100,
          title: { display: true, text: '质量分 Quality Score（越高越好）' },
          ticks: { callback: value => value + '' },
          grid: { color: '#d8e0e7', tickLength: 0 },
          border: { color: '#7d8793', width: 2 }
        }
      }
    }
  });

  renderEfficiencyNotes(axis);
}

function renderEfficiencyNotes(axis) {
  const box = document.getElementById('efficiencyQualityNotes');
  if (!box) return;
  const axisIsTokens = axis === 'tokens';
  const ranked = reportTopModels
    .map(item => ({ item, score: qualityScore(item) }))
    .sort((a, b) => b.score - a.score);
  const visible = ranked.slice(0, 6);
  const extra = ranked.length > visible.length ? '<div class="muted">其余 ' + (ranked.length - visible.length) + ' 个模型请悬停图中点查看。</div>' : '';
  box.innerHTML = visible.map(({ item, score }) => {
    const x = axisIsTokens ? Number(item.avg_total_tokens || 0).toFixed(0) + ' tokens' : (Number(item.avg_latency_ms || 0) / 1000).toFixed(1) + 's';
    return '<div><strong>' + safeText(item.model) + '</strong>：质量分 ' + score.toFixed(1) + '，' + (axisIsTokens ? '平均 Token ' : '平均耗时 ') + x + '。</div>';
  }).join('') + extra;
}

function initEfficiencyAxisToggle() {
  const buttons = document.querySelectorAll('[data-efficiency-axis]');
  if (!buttons.length) return;
  buttons.forEach(button => {
    button.addEventListener('click', () => {
      buttons.forEach(item => item.classList.remove('active'));
      button.classList.add('active');
      renderEfficiencyQualityChart(button.getAttribute('data-efficiency-axis') || 'latency');
    });
  });
}

function renderModelCharts() {
  const modelNames = reportTopModels.map(item => item.model);
  const compileRates = reportTopModels.map(item => item.compile_pass_rate);
  const testRates = reportTopModels.map(item => item.avg_test_pass_rate);
  const lineRates = reportTopModels.map(item => item.avg_line_coverage);
  const mutationRates = reportTopModels.map(item => item.avg_mutation_score);

  modelBarChart = new Chart(document.getElementById('modelBarChart'), {
    type: 'bar',
    data: { labels: modelNames, datasets: [
      { label: '编译', data: compileRates, backgroundColor: '#315f9c' },
      { label: '测试', data: testRates, backgroundColor: '#2f7d72' },
      { label: '覆盖', data: lineRates, backgroundColor: '#d5902f' },
      { label: '变异', data: mutationRates, backgroundColor: '#c43d2f' }
    ]},
    options: { responsive: true, maintainAspectRatio: false, scales: { y: { beginAtZero: true, max: 1, ticks: { callback: value => Math.round(value * 100) + '%' } } } }
  });

  radarChart = new Chart(document.getElementById('radarChart'), {
    type: 'radar',
    data: {
      labels: ['编译', '测试', '覆盖', '变异'],
      datasets: reportTopModels.map((item, index) => ({
        label: item.model,
        data: [item.compile_pass_rate, item.avg_test_pass_rate, item.avg_line_coverage, item.avg_mutation_score],
        fill: true,
        backgroundColor: ['rgba(49,95,156,0.16)', 'rgba(47,125,114,0.16)', 'rgba(213,144,47,0.18)', 'rgba(196,61,47,0.16)'][index % 4],
        borderColor: ['#315f9c', '#2f7d72', '#d5902f', '#c43d2f'][index % 4],
        pointBackgroundColor: ['#315f9c', '#2f7d72', '#d5902f', '#c43d2f'][index % 4]
      }))
    },
    options: { responsive: true, maintainAspectRatio: false, scales: { r: { beginAtZero: true, max: 1, ticks: { callback: value => Math.round(value * 100) + '%' } } } }
  });

  renderEfficiencyQualityChart('latency');
  initEfficiencyAxisToggle();
}

function selectedFilterValues(kind) {
  const all = document.querySelector('[data-filter-all="' + kind + '"]');
  if (all && all.checked) return [];
  return Array.from(document.querySelectorAll('[data-filter-' + kind + ']:checked')).map(item => item.value);
}

function formatFilterSummary(kind, values, allText) {
  if (!values.length) return allText;
  if (values.length <= 3) return values.map(value => kind === 'scenario' ? scenarioLabel(value) : value).join(', ');
  return values.slice(0, 3).map(value => kind === 'scenario' ? scenarioLabel(value) : value).join(', ') + ' 等' + values.length + '项';
}

function updateFilterSummary(models, languages, scenarios) {
  const box = document.getElementById('filter-summary');
  if (!box) return;
  box.textContent = '当前筛选：' +
    formatFilterSummary('model', models, '全部模型') + ' · ' +
    formatFilterSummary('language', languages.map(item => item.toUpperCase()), '全部语言') + ' · ' +
    formatFilterSummary('scenario', scenarios, '全部场景');
}

function initAnalysisFilters() {
  document.querySelectorAll('[data-filter-all]').forEach(all => {
    const kind = all.getAttribute('data-filter-all');
    all.addEventListener('change', () => {
      if (all.checked) {
        document.querySelectorAll('[data-filter-' + kind + ']').forEach(item => { item.checked = false; });
      }
      renderFilteredSections();
    });
  });
  ['model', 'language', 'scenario'].forEach(kind => {
    document.querySelectorAll('[data-filter-' + kind + ']').forEach(item => {
      item.addEventListener('change', () => {
        const all = document.querySelector('[data-filter-all="' + kind + '"]');
        if (all) all.checked = document.querySelectorAll('[data-filter-' + kind + ']:checked').length === 0;
        renderFilteredSections();
      });
    });
  });
  const reset = document.getElementById('reset-analysis-filters');
  if (reset) {
    reset.addEventListener('click', () => {
      document.querySelectorAll('[data-filter-model], [data-filter-language], [data-filter-scenario]').forEach(item => { item.checked = false; });
      document.querySelectorAll('[data-filter-all]').forEach(item => { item.checked = true; });
      renderFilteredSections();
    });
  }
}

function renderFilteredSections() {
  const selectedModels = selectedFilterValues('model');
  const selectedLanguages = selectedFilterValues('language');
  const selectedScenarios = selectedFilterValues('scenario');
  updateFilterSummary(selectedModels, selectedLanguages, selectedScenarios);
  const aggregated = aggregateRows(selectedModels, selectedLanguages, selectedScenarios);
  renderModelTable(aggregated.models);
  renderLanguageTable(aggregated.languages);
  renderScenarioTable(aggregated.scenarios);
  renderErrorTable(aggregated.failureRows);
  renderErrorCharts(aggregated.failureRows);
  renderScenarioCharts(aggregated.scenarios);
  const exportBtn = document.getElementById('export-scenario-csv');
  if (exportBtn) exportBtn.onclick = () => exportScenarioCSV(aggregated.scenarios);
}

function initReportSidebar() {
  const links = Array.from(document.querySelectorAll('.report-sidebar a[href^="#"]'));
  if (!links.length) return;
  const sections = links
    .map(link => ({ link, section: document.querySelector(link.getAttribute('href')) }))
    .filter(item => item.section);
  if (!sections.length) return;

  function setActive() {
    const currentY = window.scrollY + 120;
    let active = sections[0];
    sections.forEach(item => {
      if (item.section.offsetTop <= currentY) active = item;
    });
    links.forEach(link => link.classList.toggle('active', link === active.link));
  }

  links.forEach(link => {
    link.addEventListener('click', () => {
      links.forEach(item => item.classList.remove('active'));
      link.classList.add('active');
    });
  });
  setActive();
  window.addEventListener('scroll', setActive, { passive: true });
}

(async function init() {
  try {
    if (window.__utBenchEnsureChartJS) {
      await window.__utBenchEnsureChartJS();
    }
    renderModelCharts();
    renderFilteredSections();
    initAnalysisFilters();
  } catch (e) {
    if (window.__utBenchShowBanner) {
      window.__utBenchShowBanner('图表库 Chart.js 加载失败，通常是网络或企业代理拦截了 CDN。请在联网环境打开，或让报告改为本地内置 Chart.js。');
    }
    if (window.console && console.error) console.error(e);
  }

  initRawColumnToggles();
  initReportSidebar();

  const backToTop = document.getElementById('back-to-top');
  window.addEventListener('scroll', () => {
    if (!backToTop) return;
    if (window.scrollY > 400) backToTop.classList.add('visible'); else backToTop.classList.remove('visible');
  });
  if (backToTop) backToTop.addEventListener('click', () => window.scrollTo({ top: 0, behavior: 'smooth' }));
})();
