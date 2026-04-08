from __future__ import annotations

import json
from html import escape
from typing import Any

from ..contracts import (
    _COLUMN_LABELS_ZH,
    _ERROR_TYPE_LABELS_ZH,
    _GEN_METRICS_REASON_LABELS_ZH,
    _GEN_METRICS_SOURCE_LABELS_ZH,
    _STAGE_LABELS_ZH,
)
from ..common import to_float as _to_float


def build_html_report(payload: dict[str, Any]) -> str:
    overall = payload.get("overall", {})
    by_model = payload.get("dimensions", {}).get("by_model", [])
    by_language = payload.get("dimensions", {}).get("by_language", [])
    failures = payload.get("failure_breakdown", [])
    sample_rows = payload.get("sample_rows", [])

    model_options = sorted({str(row.get("model") or "unknown") for row in sample_rows})
    language_options = sorted({str(row.get("language") or "unknown") for row in sample_rows})

    generated_at = escape(str(payload.get("generated_at_utc") or ""))
    source_file = escape(str(payload.get("source_evaluator_summary") or ""))

    by_complexity = payload.get("dimensions", {}).get("by_complexity", [])
    by_model_language = payload.get("dimensions", {}).get("by_model_language", [])
    top_models = sorted(
        by_model,
        key=lambda row: (
            -float(row.get("test_pass_rate") or 0.0),
            -float(row.get("avg_line_coverage") or 0.0),
            -float(row.get("avg_mutation_score") or 0.0),
        ),
    )
    model_names = [str(row.get("model") or "unknown") for row in by_model]
    test_pass_rates = [float(row.get("test_pass_rate") or 0.0) for row in by_model]
    line_covs = [float(row.get("avg_line_coverage") or 0.0) for row in by_model]
    mutation_scores = [float(row.get("avg_mutation_score") or 0.0) for row in by_model]
    compile_rates = [float(row.get("compile_pass_rate") or 0.0) for row in by_model]
    branch_covs = [float(row.get("avg_branch_coverage") or 0.0) for row in by_model]

    complexity_names = [str(row.get("complexity") or "unknown") for row in by_complexity]
    complexity_test_rates = [float(row.get("test_pass_rate") or 0.0) for row in by_complexity]

    heatmap_models = sorted({str(row.get("model") or "unknown") for row in by_model_language})
    heatmap_languages = sorted({str(row.get("language") or "unknown") for row in by_model_language})
    heatmap_values = _build_heatmap_matrix(by_model_language, heatmap_models, heatmap_languages)

    threshold_map = payload.get("thresholds", {})
    threshold_compile = float(threshold_map.get("compile_pass_rate", 1.0))
    threshold_test = float(threshold_map.get("test_pass_rate", 0.7))
    threshold_line = float(threshold_map.get("line_coverage", 0.7))
    threshold_branch = float(threshold_map.get("branch_coverage", 0.6))
    threshold_mutation = float(threshold_map.get("mutation_score", 0.85))

    palette = _palette(payload.get("chart_style"))

    summary_cards = [
        _summary_card("样本总数", overall.get("total_samples")),
        _summary_card("编译通过率", _format_pct(overall.get("compile_pass_rate"))),
        _summary_card("测试通过率", _format_pct(overall.get("test_pass_rate"))),
        _summary_card("平均行覆盖率", _format_pct(overall.get("avg_line_coverage"))),
        _summary_card("平均分支覆盖率", _format_pct(overall.get("avg_branch_coverage"))),
        _summary_card("平均变异得分", _format_pct(overall.get("avg_mutation_score"))),
        _summary_card("有效变异杀死率", _format_pct(overall.get("mutation_effective_kill_rate"))),
        _summary_card("平均生成耗时(ms)", _format_num(overall.get("avg_generation_latency_ms"))),
        _summary_card("平均总 Token", _format_num(overall.get("avg_total_tokens"))),
    ]

    rank_rows = []
    for index, row in enumerate(top_models, start=1):
        rank_rows.append(
            {
                "排名": index,
                "模型": row.get("model"),
                "测试通过率": _format_pct(row.get("test_pass_rate")),
                "平均行覆盖率": _format_pct(row.get("avg_line_coverage")),
                "平均变异得分": _format_pct(row.get("avg_mutation_score")),
                "平均生成耗时(ms)": _format_num(row.get("avg_generation_latency_ms")),
                "平均总Token": _format_num(row.get("avg_total_tokens")),
            }
        )

    appendix_rows = [
        {
            "model": row.get("model"),
            "language": row.get("language"),
            "sample_id": row.get("sample_id"),
            "complexity": row.get("complexity"),
            "scenario": row.get("scenario"),
            "compile_pass": row.get("compile_pass"),
            "test_pass": row.get("test_pass"),
            "line_coverage": row.get("line_coverage"),
            "branch_coverage": row.get("branch_coverage"),
            "mutation_score": row.get("mutation_score"),
            "mutation_total": row.get("mutation_total"),
            "mutation_killed": row.get("mutation_killed"),
            "mutation_survived": row.get("mutation_survived"),
            "mutation_no_tests": row.get("mutation_no_tests"),
            "mutation_not_checked": row.get("mutation_not_checked"),
            "mutation_timeout": row.get("mutation_timeout"),
            "mutation_skipped": row.get("mutation_skipped"),
            "mutation_suspicious": row.get("mutation_suspicious"),
            "generation_latency_ms": row.get("generation_latency_ms"),
            "total_tokens": row.get("total_tokens"),
            "generation_metrics_source": row.get("generation_metrics_source"),
            "generation_metrics_reason": row.get("generation_metrics_reason"),
            "test_error": row.get("test_error"),
            "coverage_error": row.get("coverage_error"),
            "mutation_error": row.get("mutation_error"),
        }
        for row in sample_rows
    ]

    return f"""<!doctype html>
<html lang=\"zh-CN\">
<head>
  <meta charset=\"utf-8\" />
  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\" />
  <title>ut-bench 可视化评测报告</title>
  <script src=\"https://cdn.jsdelivr.net/npm/chart.js@4.4.1/dist/chart.umd.min.js\"></script>
  <style>
    :root {{
      --bg: linear-gradient(180deg, #f6fbff 0%, #f7f8fc 100%);
      --card: #ffffff;
      --text: #102a43;
      --muted: #6b7280;
      --accent: {palette['accent']};
      --accent-soft: {palette['accent_soft']};
      --accent-2: {palette['accent_2']};
      --warn: {palette['warn']};
      --ok: {palette['ok']};
      --line: #e5e7eb;
    }}
    body {{ margin: 0; background: var(--bg); color: var(--text); font-family: "Noto Sans SC", "Microsoft YaHei", "PingFang SC", sans-serif; }}
    .wrap {{ max-width: 1200px; margin: 0 auto; padding: 24px; }}
    .hero {{ background: radial-gradient(circle at right top, var(--accent-soft), #fff 55%); border: 1px solid var(--line); border-radius: 16px; padding: 18px 20px; margin-bottom: 16px; }}
    .hero h1 {{ margin: 0; font-size: 28px; letter-spacing: 0.3px; }}
    .hero .sub {{ margin-top: 6px; color: var(--muted); font-size: 13px; line-height: 1.6; }}
    .meta {{ color: var(--muted); font-size: 12px; margin-top: 10px; }}
    .cards {{ display: grid; gap: 12px; grid-template-columns: repeat(auto-fill, minmax(210px, 1fr)); margin-bottom: 24px; }}
    .card {{ background: var(--card); border: 1px solid var(--line); border-radius: 10px; padding: 12px 14px; }}
    .card .k {{ color: var(--muted); font-size: 12px; }}
    .card .v {{ margin-top: 4px; font-size: 22px; font-weight: 700; color: var(--accent); }}
    h2 {{ margin: 24px 0 10px; font-size: 20px; }}
    h3 {{ margin: 16px 0 10px; font-size: 16px; color: #334e68; }}
    .section {{ margin-bottom: 26px; }}
    .grid-2 {{ display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }}
    .panel {{ background: var(--card); border: 1px solid var(--line); border-radius: 10px; padding: 12px; }}
    .chart-box {{ height: 300px; }}
    .chart-box canvas {{ width: 100% !important; height: 100% !important; }}
    .conclusion {{ margin: 0; padding-left: 18px; color: #243b53; line-height: 1.7; }}
    .conclusion li strong {{ color: var(--accent); }}
    table {{ width: 100%; border-collapse: collapse; background: var(--card); border: 1px solid var(--line); border-radius: 10px; overflow: hidden; }}
    th, td {{ border-bottom: 1px solid var(--line); padding: 8px 10px; font-size: 13px; text-align: left; vertical-align: top; }}
    th {{ background: #f8fbff; }}
    tr:last-child td {{ border-bottom: none; }}
    .filters {{ display: flex; gap: 12px; margin: 10px 0; flex-wrap: wrap; }}
    select, input {{ padding: 6px 8px; border: 1px solid var(--line); border-radius: 8px; background: #fff; }}
    .small {{ color: var(--muted); font-size: 12px; }}
    .pill {{ display: inline-block; padding: 2px 8px; border-radius: 999px; font-size: 11px; margin-left: 6px; background: #eef2ff; color: #3730a3; }}
    @media (max-width: 900px) {{
      .grid-2 {{ grid-template-columns: 1fr; }}
      .wrap {{ padding: 14px; }}
      .hero h1 {{ font-size: 22px; }}
      .chart-box {{ height: 260px; }}
    }}
  </style>
</head>
<body>
  <div class=\"wrap\">
    <div class=\"hero\">
      <h1>ut-bench 多模型单测生成可视化报告</h1>
      <div class=\"sub\">摘要 + 详细分析 + 附录三段式呈现。支持按模型、语言、复杂度进行多维分析，并输出 JSON/CSV/HTML 产物。</div>
      <div class=\"meta\">生成时间：{generated_at} | 数据源：{source_file}</div>
    </div>

    <div class=\"section\">
      <h2>摘要</h2>
      <div class=\"cards\">{''.join(summary_cards)}</div>
      <h3>模型综合排名 <span class=\"pill\">按测试通过率→覆盖率→变异得分</span></h3>
      {_table_html(rank_rows, ['排名', '模型', '测试通过率', '平均行覆盖率', '平均变异得分', '平均生成耗时(ms)', '平均总Token'])}
    </div>

    <div class=\"section\">
      <h2>详细分析</h2>
      <div class=\"grid-2\">
        <div class=\"panel\">
          <h3>模型关键指标柱状图</h3>
          <div class=\"chart-box\"><canvas id=\"barModel\"></canvas></div>
        </div>
        <div class=\"panel\">
          <h3>复杂度-测试通过率折线图</h3>
          <div class=\"chart-box\"><canvas id=\"lineComplexity\"></canvas></div>
        </div>
        <div class=\"panel\">
          <h3>五维能力雷达图（模型均值）</h3>
          <div class=\"chart-box\"><canvas id=\"radarDimensions\"></canvas></div>
        </div>
        <div class=\"panel\">
          <h3>模型×语言热力图（测试通过率）</h3>
          <div class=\"chart-box\"><canvas id=\"heatmapModelLang\"></canvas></div>
        </div>
      </div>
      <h3>失败类型统计</h3>
      {_table_html(failures, ['stage', 'error_type', 'count', 'example_model', 'example_sample', 'example_message'])}
      <h3>结论摘要</h3>
      <ul class=\"conclusion\">
        <li><strong>正确性：</strong>整体编译通过率 {_format_pct(overall.get('compile_pass_rate'))}，测试通过率 {_format_pct(overall.get('test_pass_rate'))}。</li>
        <li><strong>覆盖率：</strong>行覆盖率 {_format_pct(overall.get('avg_line_coverage'))}，分支覆盖率 {_format_pct(overall.get('avg_branch_coverage'))}。</li>
        <li><strong>有效性：</strong>变异得分均值 {_format_pct(overall.get('avg_mutation_score'))}，有效变异杀死率 {_format_pct(overall.get('mutation_effective_kill_rate'))}，与目标 {_format_pct(threshold_mutation)} 对比可持续优化。</li>
        <li><strong>效率：</strong>平均生成耗时 {_format_num(overall.get('avg_generation_latency_ms'))} ms，平均总 Token {_format_num(overall.get('avg_total_tokens'))}。</li>
      </ul>
      <h3>数据质量说明</h3>
      {_build_data_quality_panel(sample_rows)}
    </div>

    <div class=\"section\">
      <h2>附录</h2>
      <h3>按模型统计表</h3>
      {_table_html(by_model, ['model', 'total_samples', 'compile_pass_rate', 'test_pass_rate', 'avg_line_coverage', 'avg_branch_coverage', 'avg_mutation_score', 'avg_generation_latency_ms', 'avg_total_tokens'])}
      <h3>按语言统计表</h3>
      {_table_html(by_language, ['language', 'total_samples', 'compile_pass_rate', 'test_pass_rate', 'avg_line_coverage', 'avg_branch_coverage', 'avg_mutation_score', 'avg_generation_latency_ms', 'avg_total_tokens'])}
      <h3>原始样本明细（支持筛选）</h3>
    <div class=\"filters\">
      <label>模型
        <select id=\"modelFilter\">
          <option value=\"all\">全部</option>
          {''.join(f'<option value="{escape(opt)}">{escape(opt)}</option>' for opt in model_options)}
        </select>
      </label>
      <label>语言
        <select id=\"languageFilter\">
          <option value=\"all\">全部</option>
          {''.join(f'<option value="{escape(opt)}">{escape(opt)}</option>' for opt in language_options)}
        </select>
      </label>
      <label>测试状态
        <select id=\"testStatusFilter\">
          <option value=\"pass\" selected>仅通过</option>
          <option value=\"all\">全部</option>
          <option value=\"fail\">仅失败</option>
          <option value=\"none\">未执行</option>
        </select>
      </label>
      <label>覆盖率
        <select id=\"coverageFilter\">
          <option value=\"all\">全部</option>
          <option value=\"has\">有覆盖率</option>
          <option value=\"none\">无覆盖率</option>
        </select>
      </label>
      <label>变异状态
        <select id=\"mutationFilter\">
          <option value=\"all\">全部</option>
          <option value=\"has_score\">有变异得分</option>
          <option value=\"has_stats\">有变异统计</option>
          <option value=\"stats_no_score\">有统计但无得分</option>
          <option value=\"none\">无变异统计</option>
          <option value=\"error\">变异阶段报错</option>
        </select>
      </label>
      <label>样本关键词
        <input id=\"sampleKeywordFilter\" type=\"text\" placeholder=\"输入 sample_id 关键字\" />
      </label>
    </div>
    <p id=\"filterStats\" class=\"small\"></p>
    {_sample_table_html(appendix_rows)}
    <p class=\"small\">提示：默认仅显示测试通过样本；可切换筛选查看失败样本和异常样本。</p>
  </div>

  <script>
    const styleConfig = {json.dumps(palette, ensure_ascii=False)};
    const modelNames = {json.dumps(model_names, ensure_ascii=False)};
    const testPassRates = {json.dumps(test_pass_rates)};
    const lineCoverages = {json.dumps(line_covs)};
    const mutationScores = {json.dumps(mutation_scores)};
    const complexityNames = {json.dumps(complexity_names, ensure_ascii=False)};
    const complexityTestRates = {json.dumps(complexity_test_rates)};
    const radarCompile = {json.dumps(compile_rates)};
    const radarBranch = {json.dumps(branch_covs)};
    const heatmapModels = {json.dumps(heatmap_models, ensure_ascii=False)};
    const heatmapLanguages = {json.dumps(heatmap_languages, ensure_ascii=False)};
    const heatmapValues = {json.dumps(heatmap_values)};

    const thresholdCompile = {threshold_compile};
    const thresholdTest = {threshold_test};
    const thresholdLine = {threshold_line};
    const thresholdBranch = {threshold_branch};
    const thresholdMutation = {threshold_mutation};

    function avg(arr) {{
      if (!arr || arr.length === 0) return 0;
      return arr.reduce((sum, value) => sum + value, 0) / arr.length;
    }}

    const barCtx = document.getElementById('barModel');
    if (barCtx && modelNames.length > 0) {{
      new Chart(barCtx, {{
        type: 'bar',
        data: {{
          labels: modelNames,
          datasets: [
            {{ label: '测试通过率', data: testPassRates.map(v => v * 100), backgroundColor: styleConfig.accent }},
            {{ label: '行覆盖率', data: lineCoverages.map(v => v * 100), backgroundColor: styleConfig.accentSoft }},
            {{ label: '变异得分', data: mutationScores.map(v => v * 100), backgroundColor: styleConfig.accent2 }}
          ]
        }},
        options: {{
          responsive: true,
          maintainAspectRatio: false,
          scales: {{ y: {{ beginAtZero: true, max: 100, ticks: {{ callback: (v) => v + '%' }} }} }},
          plugins: {{ legend: {{ position: 'bottom' }} }}
        }}
      }});
    }}

    const lineCtx = document.getElementById('lineComplexity');
    if (lineCtx && complexityNames.length > 0) {{
      new Chart(lineCtx, {{
        type: 'line',
        data: {{
          labels: complexityNames,
          datasets: [
            {{
              label: '测试通过率',
              data: complexityTestRates.map(v => v * 100),
              borderColor: styleConfig.accent,
              backgroundColor: styleConfig.accentSoft,
              tension: 0.3,
              fill: true,
            }},
            {{
              label: '测试通过率阈值',
              data: complexityNames.map(() => thresholdTest * 100),
              borderColor: styleConfig.warn,
              borderDash: [6, 4],
              pointRadius: 0,
              tension: 0,
            }},
          ]
        }},
        options: {{
          responsive: true,
          maintainAspectRatio: false,
          scales: {{ y: {{ beginAtZero: true, max: 100, ticks: {{ callback: (v) => v + '%' }} }} }},
          plugins: {{ legend: {{ position: 'bottom' }} }}
        }}
      }});
    }}

    const radarCtx = document.getElementById('radarDimensions');
    if (radarCtx) {{
      new Chart(radarCtx, {{
        type: 'radar',
        data: {{
          labels: ['编译通过率', '测试通过率', '行覆盖率', '分支覆盖率', '变异得分'],
          datasets: [
            {{
              label: '模型均值',
              data: [
                avg(radarCompile) * 100,
                avg(testPassRates) * 100,
                avg(lineCoverages) * 100,
                avg(radarBranch) * 100,
                avg(mutationScores) * 100,
              ],
              borderColor: styleConfig.accent,
              backgroundColor: styleConfig.accentSoft,
            }},
            {{
              label: '目标阈值',
              data: [
                thresholdCompile * 100,
                thresholdTest * 100,
                thresholdLine * 100,
                thresholdBranch * 100,
                thresholdMutation * 100,
              ],
              borderColor: styleConfig.warn,
              backgroundColor: 'rgba(0,0,0,0)',
            }}
          ]
        }},
        options: {{
          responsive: true,
          maintainAspectRatio: false,
          scales: {{ r: {{ min: 0, max: 100, ticks: {{ callback: (v) => v + '%' }} }} }},
          plugins: {{ legend: {{ position: 'bottom' }} }}
        }}
      }});
    }}

    function drawHeatmap() {{
      const canvas = document.getElementById('heatmapModelLang');
      if (!canvas || heatmapModels.length === 0 || heatmapLanguages.length === 0) return;
      const ctx = canvas.getContext('2d');
      const width = canvas.width;
      const height = canvas.height;
      ctx.clearRect(0, 0, width, height);

      const paddingLeft = 90;
      const paddingTop = 30;
      const paddingRight = 20;
      const paddingBottom = 40;
      const gridWidth = width - paddingLeft - paddingRight;
      const gridHeight = height - paddingTop - paddingBottom;
      const cellW = gridWidth / heatmapLanguages.length;
      const cellH = gridHeight / heatmapModels.length;

      function colorFor(value) {{
        const v = Math.max(0, Math.min(1, value));
        const r = Math.round(240 - 90 * v);
        const g = Math.round(249 - 120 * v);
        const b = Math.round(255 - 180 * v);
        return `rgb(${{r}}, ${{g}}, ${{b}})`;
      }}

      ctx.font = '12px "Noto Sans SC", "Microsoft YaHei", sans-serif';
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';

      for (let i = 0; i < heatmapModels.length; i++) {{
        for (let j = 0; j < heatmapLanguages.length; j++) {{
          const x = paddingLeft + j * cellW;
          const y = paddingTop + i * cellH;
          const value = heatmapValues[i][j];
          ctx.fillStyle = colorFor(value);
          ctx.fillRect(x, y, cellW, cellH);
          ctx.strokeStyle = '#dbe7f5';
          ctx.strokeRect(x, y, cellW, cellH);
          ctx.fillStyle = value >= 0.6 ? '#102a43' : '#334e68';
          ctx.fillText((value * 100).toFixed(1) + '%', x + cellW / 2, y + cellH / 2);
        }}
      }}

      ctx.fillStyle = '#334e68';
      ctx.textAlign = 'right';
      for (let i = 0; i < heatmapModels.length; i++) {{
        const y = paddingTop + i * cellH + cellH / 2;
        ctx.fillText(heatmapModels[i], paddingLeft - 8, y);
      }}

      ctx.textAlign = 'center';
      for (let j = 0; j < heatmapLanguages.length; j++) {{
        const x = paddingLeft + j * cellW + cellW / 2;
        ctx.fillText(heatmapLanguages[j], x, height - 18);
      }}
    }}

    drawHeatmap();
    window.addEventListener('resize', drawHeatmap);

    function applyFilters() {{
      const model = document.getElementById('modelFilter').value;
      const lang = document.getElementById('languageFilter').value;
      const testStatus = document.getElementById('testStatusFilter').value;
      const coverageStatus = document.getElementById('coverageFilter').value;
      const mutationStatus = document.getElementById('mutationFilter').value;
      const sampleKeyword = (document.getElementById('sampleKeywordFilter').value || '').trim().toLowerCase();
      const rows = document.querySelectorAll('#sampleTable tbody tr');
      let visible = 0;
      rows.forEach((row) => {{
        const rowModel = row.getAttribute('data-model');
        const rowLang = row.getAttribute('data-language');
        const rowSampleId = (row.getAttribute('data-sample-id') || '').toLowerCase();
        const rowTestPass = row.getAttribute('data-test-pass');
        const rowHasCoverage = row.getAttribute('data-has-coverage');
        const rowHasMutationScore = row.getAttribute('data-has-mutation-score');
        const rowHasMutationStats = row.getAttribute('data-has-mutation-stats');
        const rowHasMutationError = row.getAttribute('data-has-mutation-error');

        const modelOk = (model === 'all') || (rowModel === model);
        const langOk = (lang === 'all') || (rowLang === lang);
        const sampleOk = (!sampleKeyword) || rowSampleId.includes(sampleKeyword);

        let testOk = true;
        if (testStatus === 'pass') testOk = rowTestPass === 'true';
        else if (testStatus === 'fail') testOk = rowTestPass === 'false';
        else if (testStatus === 'none') testOk = rowTestPass === 'none';

        let coverageOk = true;
        if (coverageStatus === 'has') coverageOk = rowHasCoverage === 'true';
        else if (coverageStatus === 'none') coverageOk = rowHasCoverage !== 'true';

        let mutationOk = true;
        if (mutationStatus === 'has_score') mutationOk = rowHasMutationScore === 'true';
        else if (mutationStatus === 'has_stats') mutationOk = rowHasMutationStats === 'true';
        else if (mutationStatus === 'stats_no_score') mutationOk = (rowHasMutationStats === 'true') && (rowHasMutationScore !== 'true');
        else if (mutationStatus === 'none') mutationOk = rowHasMutationStats !== 'true';
        else if (mutationStatus === 'error') mutationOk = rowHasMutationError === 'true';

        const show = modelOk && langOk && sampleOk && testOk && coverageOk && mutationOk;
        if (show) visible += 1;
        row.style.display = show ? '' : 'none';
      }});
      const stats = document.getElementById('filterStats');
      if (stats) {{
        stats.textContent = `当前筛选：${{visible}} / ${{rows.length}} 条样本`;
      }}
    }}
    document.getElementById('modelFilter').addEventListener('change', applyFilters);
    document.getElementById('languageFilter').addEventListener('change', applyFilters);
    document.getElementById('testStatusFilter').addEventListener('change', applyFilters);
    document.getElementById('coverageFilter').addEventListener('change', applyFilters);
    document.getElementById('mutationFilter').addEventListener('change', applyFilters);
    document.getElementById('sampleKeywordFilter').addEventListener('input', applyFilters);
    applyFilters();
  </script>
</body>
</html>
"""


def _build_heatmap_matrix(
    rows: list[dict[str, Any]],
    models: list[str],
    languages: list[str],
) -> list[list[float]]:
    value_map: dict[tuple[str, str], float] = {}
    for row in rows:
        model = str(row.get("model") or "unknown")
        language = str(row.get("language") or "unknown")
        value_map[(model, language)] = float(row.get("test_pass_rate") or 0.0)

    matrix: list[list[float]] = []
    for model in models:
        matrix.append([value_map.get((model, language), 0.0) for language in languages])
    return matrix


def _palette(style: Any) -> dict[str, str]:
    name = str(style or "teal").strip().lower()
    palettes = {
        "teal": {
            "accent": "#0f766e",
            "accent_soft": "rgba(15, 118, 110, 0.20)",
            "accent_2": "rgba(12, 74, 110, 0.75)",
            "warn": "#b45309",
            "ok": "#166534",
        },
        "orange": {
            "accent": "#c2410c",
            "accent_soft": "rgba(194, 65, 12, 0.20)",
            "accent_2": "rgba(146, 64, 14, 0.75)",
            "warn": "#9a3412",
            "ok": "#166534",
        },
        "blue": {
            "accent": "#1d4ed8",
            "accent_soft": "rgba(29, 78, 216, 0.20)",
            "accent_2": "rgba(30, 64, 175, 0.75)",
            "warn": "#b45309",
            "ok": "#15803d",
        },
    }
    return palettes.get(name, palettes["teal"])


def _normalize_formats(formats: list[str] | None) -> set[str]:
    if not formats:
        return {"json", "csv", "html"}

    allowed = {"json", "csv", "html"}
    output = {item.strip().lower() for item in formats if item and item.strip()}
    unknown = output - allowed
    if unknown:
        raise ValueError(f"Unsupported report format(s): {sorted(unknown)}")
    return output or {"json", "csv", "html"}


def _normalize_thresholds(thresholds: dict[str, float] | None) -> dict[str, float]:
    defaults = {
        "compile_pass_rate": 1.0,
        "test_pass_rate": 0.7,
        "line_coverage": 0.7,
        "branch_coverage": 0.6,
        "mutation_score": 0.85,
    }
    if not thresholds:
        return defaults

    merged = dict(defaults)
    for key, value in thresholds.items():
        if key not in merged:
            continue
        merged[key] = float(value)
    return merged


def _summary_card(key: str, value: Any) -> str:
    return (
        "<div class=\"card\">"
        f"<div class=\"k\">{escape(str(key))}</div>"
        f"<div class=\"v\">{escape(str(value if value is not None else 'N/A'))}</div>"
        "</div>"
    )


def _table_html(rows: list[dict[str, Any]], columns: list[str]) -> str:
    if not rows:
        return "<p class=\"small\">No data.</p>"

    head = "".join(f"<th>{escape(_column_title(col))}</th>" for col in columns)
    body_rows: list[str] = []
    for row in rows:
        cells = "".join(f"<td>{escape(_display_value(col, row.get(col)))}</td>" for col in columns)
        body_rows.append(f"<tr>{cells}</tr>")
    body = "".join(body_rows)
    return f"<table><thead><tr>{head}</tr></thead><tbody>{body}</tbody></table>"


def _sample_table_html(rows: list[dict[str, Any]]) -> str:
    columns = [
        "model",
        "language",
        "sample_id",
        "complexity",
        "scenario",
        "compile_pass",
        "test_pass",
        "line_coverage",
        "mutation_score",
        "mutation_total",
        "mutation_killed",
        "mutation_survived",
        "mutation_no_tests",
        "mutation_not_checked",
        "mutation_timeout",
        "mutation_skipped",
        "mutation_suspicious",
        "generation_latency_ms",
        "total_tokens",
        "generation_metrics_source",
        "generation_metrics_reason",
    ]
    if not rows:
        return "<p class=\"small\">No sample rows.</p>"

    head = "".join(f"<th>{escape(_column_title(col))}</th>" for col in columns)
    body_rows: list[str] = []
    for row in rows:
        model = escape(str(row.get("model") or "unknown"))
        language = escape(str(row.get("language") or "unknown"))
        sample_id = escape(str(row.get("sample_id") or ""))
        raw_test_pass = row.get("test_pass")
        if raw_test_pass is True:
            test_pass_attr = "true"
        elif raw_test_pass is False:
            test_pass_attr = "false"
        else:
            test_pass_attr = "none"
        has_coverage = "true" if row.get("line_coverage") is not None else "false"
        has_mutation_score = "true" if row.get("mutation_score") is not None else "false"
        has_mutation_stats = "true" if row.get("mutation_total") is not None else "false"
        has_mutation_error = "true" if row.get("mutation_error") else "false"
        cells = "".join(f"<td>{escape(_display_value(col, row.get(col)))}</td>" for col in columns)
        body_rows.append(
            f"<tr data-model=\"{model}\" data-language=\"{language}\" "
            f"data-sample-id=\"{sample_id}\" data-test-pass=\"{test_pass_attr}\" "
            f"data-has-coverage=\"{has_coverage}\" "
            f"data-has-mutation-score=\"{has_mutation_score}\" "
            f"data-has-mutation-stats=\"{has_mutation_stats}\" "
            f"data-has-mutation-error=\"{has_mutation_error}\">{cells}</tr>"
        )
    body = "".join(body_rows)
    return (
        f"<table id=\"sampleTable\"><thead><tr>{head}</tr></thead>"
        f"<tbody>{body}</tbody></table>"
    )


def _display_value(key: str, value: Any) -> str:
    if value is None:
        return "N/A"
    if isinstance(value, bool):
        return "是" if value else "否"
    if key == "stage":
        return _STAGE_LABELS_ZH.get(str(value), str(value))
    if key == "error_type":
        return _ERROR_TYPE_LABELS_ZH.get(str(value), str(value))
    if key == "generation_metrics_source":
        return _GEN_METRICS_SOURCE_LABELS_ZH.get(str(value), str(value))
    if key == "generation_metrics_reason":
        return _GEN_METRICS_REASON_LABELS_ZH.get(str(value), str(value))
    if isinstance(value, (int, float)):
        if key.endswith("_rate") or "coverage" in key or key.endswith("_score"):
            return _format_pct(float(value))
        return _format_num(float(value))
    return str(value)


def _build_data_quality_panel(rows: list[dict[str, Any]]) -> str:
    total = len(rows)
    if total <= 0:
        return "<p class=\"small\">无样本数据。</p>"

    missing_latency = sum(1 for row in rows if row.get("generation_latency_ms") is None)
    missing_tokens = sum(1 for row in rows if row.get("total_tokens") is None)
    pass_rows = [row for row in rows if row.get("test_pass") is True]
    pass_count = len(pass_rows)
    missing_cov_on_pass = sum(1 for row in pass_rows if row.get("line_coverage") is None)
    has_metadata = sum(1 for row in rows if row.get("generation_metrics_source") == "metadata")
    has_coverage = sum(1 for row in rows if row.get("line_coverage") is not None)
    has_mutation = sum(1 for row in rows if row.get("mutation_total") is not None)
    mut_not_checked_like = sum(
        1
        for row in rows
        if row.get("mutation_total") not in (None, 0)
        and (row.get("mutation_killed") or 0) == 0
        and (row.get("mutation_survived") or 0) == 0
    )

    return (
        "<table><thead><tr>"
        "<th>指标</th><th>数量</th><th>占比</th><th>解释</th>"
        "</tr></thead><tbody>"
        f"<tr><td>Metadata 完整样本</td><td>{has_metadata}</td><td>{_format_pct(has_metadata/total)}</td>"
        "<td>来自 runner metadata 的样本比例，决定耗时/Token 指标可信度。</td></tr>"
        f"<tr><td>覆盖率有效样本</td><td>{has_coverage}</td><td>{_format_pct(has_coverage/total)}</td>"
        "<td>line_coverage 非空样本比例，反映测试是否命中真实源码。</td></tr>"
        f"<tr><td>变异统计有效样本</td><td>{has_mutation}</td><td>{_format_pct(has_mutation/total)}</td>"
        "<td>mutation_total 非空样本比例，反映 mutmut 是否完成统计写出。</td></tr>"
        f"<tr><td>生成耗时缺失</td><td>{missing_latency}</td><td>{_format_pct(missing_latency/total)}</td>"
        "<td>历史样本缺 metadata，当前仅新一轮 runner 样本有完整时延与 token。</td></tr>"
        f"<tr><td>Token 缺失</td><td>{missing_tokens}</td><td>{_format_pct(missing_tokens/total)}</td>"
        "<td>同上，属于历史数据源不完整，不是 reporter 计算错误。</td></tr>"
        f"<tr><td>测试通过但无覆盖率</td><td>{missing_cov_on_pass}</td><td>{_format_pct(missing_cov_on_pass/max(pass_count,1))}</td>"
        "<td>多数是测试未命中真实被测源码（导入占位模块/测到测试内函数）。</td></tr>"
        f"<tr><td>变异体总数>0但杀死/存活都为0</td><td>{mut_not_checked_like}</td><td>{_format_pct(mut_not_checked_like/total)}</td>"
        "<td>通常为 mutmut 未真正执行到变异体（not checked），需结合 mutation_error 诊断。</td></tr>"
        "</tbody></table>"
    )


def _column_title(column: str) -> str:
    return _COLUMN_LABELS_ZH.get(column, column)


def _format_pct(value: Any) -> str:
    number = _to_float(value)
    if number is None:
        return "N/A"
    return f"{number * 100:.2f}%"


def _format_num(value: Any) -> str:
    number = _to_float(value)
    if number is None:
        return "N/A"
    if abs(number - int(number)) < 1e-9:
        return str(int(number))
    return f"{number:.6f}".rstrip("0").rstrip(".")
