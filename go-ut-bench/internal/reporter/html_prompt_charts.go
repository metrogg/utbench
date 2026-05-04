// reporter 包 - Prompt 展示与图表脚本生成
// 包含 Prompt 展示区、中文说明、交互脚本和图表脚本
package reporter

import (
	"fmt"
	"strings"

	"go-ut-bench/internal/contracts"
)

// buildPromptHTMLNew 生成新的 Prompt 展示区域
func buildPromptHTMLNew(strategy, versionID string, prompts map[string]string, promptProvider contracts.PromptMetaProvider) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<div class="section prompt-compact" id="prompt">
  <h2>Prompt 策略与实际提示词</h2>
  <div class="prompt-meta">
    <span>Strategy: <strong>%s</strong></span>
    <span>Version: <strong>%s</strong></span>
    <span>展示优先级：实际 rendered prompt 快照 &gt; 模板预览</span>
  </div>
  <div class="prompt-tags">
    <span class="prompt-tag">真实运行快照</span>
    <span class="prompt-tag">跨模型一致</span>
    <span class="prompt-tag">四语言对照</span>
    <span class="prompt-tag">中文说明</span>
  </div>
  <ul class="prompt-bullets">
    <li>这里展示的是生成阶段发送给模型的提示词。若本次 run 保存了 rendered prompt，则直接读取样本快照；否则展示同版本模板预览。</li>
    <li>Prompt 不显式暴露 benchmark 场景名和复杂度，避免用题型标签提示模型。</li>
    <li>Python 在正式评测前会把中性导入名 <code>module_under_test</code> 重写到真实源码模块，减少文件名泄露。</li>
    <li>所有语言都约束小而有代表性的输入、禁止真实网络/认证/外部服务，并优先使用 fake/stub。</li>
  </ul>`, escapeHTML(emptyDash(strategy)), escapeHTML(emptyDash(versionID))))

	languages := contracts.SupportedLanguages
	b.WriteString(`<div class="prompt-language-grid">`)
	for _, lang := range languages {
		prompt := prompts[lang]
		if strings.TrimSpace(prompt) == "" {
			prompt = promptProvider.PromptTemplatePreview(lang)
		}
		mode := promptField(prompt, "Mode")
		framework := promptField(prompt, "Framework")
		sourceKind := "模板预览"
		if !strings.Contains(prompt, "preview_sample") && !strings.Contains(prompt, "PreviewSample") {
			sourceKind = "实际渲染快照"
		}
		displayPrompt := truncateText(prompt, 12000)
		truncatedNote := ""
		if len(displayPrompt) < len(prompt) {
			truncatedNote = fmt.Sprintf(`<div class="prompt-note">原文较长，当前仅展示前 %d 字符；完整内容已保存在 rendered prompt 快照文件中。</div>`, len(displayPrompt))
		}
		b.WriteString(fmt.Sprintf(`
  <details class="prompt-language-card">
    <summary>
      <span>%s Prompt</span>
      <small>%s · %s · %s</small>
    </summary>
    <div class="prompt-translation">
      <h3>中文说明</h3>
      %s
    </div>
    %s
    <div class="prompt-template-box">%s</div>
  </details>`,
			escapeHTML(strings.ToUpper(lang)),
			escapeHTML(emptyDash(mode)),
			escapeHTML(emptyDash(framework)),
			escapeHTML(sourceKind),
			promptChineseSummary(lang, mode),
			truncatedNote,
			escapeHTML(displayPrompt)))
	}
	b.WriteString(`</div>`)

	b.WriteString(`
</div>`)
	return b.String()
}

func promptField(prompt, field string) string {
	prefix := field + ":"
	for _, line := range strings.Split(prompt, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			return strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), prefix))
		}
	}
	return ""
}

func promptChineseSummary(language, mode string) string {
	modeText := map[string]string{
		"full_file":  "完整文件模式：模型需要直接返回一个完整、可运行的测试文件。",
		"repo_level": "仓库级模式：模型基于多文件模块上下文生成目标模块的完整测试文件。",
		"completion": "续写模式：模型只补充下一个有价值的测试函数或测试块。",
	}
	if modeText[mode] == "" {
		modeText[mode] = "当前模式由 prompt 原文中的 Mode 字段决定。"
	}
	rules := []string{
		"必须只输出原始测试代码，不允许 Markdown 代码块、解释文字或占位测试。",
		"测试必须可重复，断言必须从源码实现行为推导，不能根据注释或常识猜测。",
		"覆盖正常路径、边界路径和错误路径，但输入要小而有代表性，不能做压力测试。",
		"不得访问真实网络、真实凭证或真实外部服务；涉及外部 I/O 时使用 mock、stub 或 fake。",
	}
	switch strings.ToLower(language) {
	case "python":
		rules = append(rules,
			"使用 pytest 函数式测试；从 module_under_test 导入目标符号，后续评测准备阶段会重写到真实模块。",
			"使用 plain assert 和 pytest.raises；需要 mock 时 patch 目标模块实际引用的符号。")
	case "go":
		rules = append(rules,
			"使用 testing 包和 TestXxx 函数；测试包名与源码包保持一致。",
			"适合时使用表驱动测试；不要假设标准库内部可以 monkey patch。")
	case "java":
		rules = append(rules,
			"使用 JUnit 5 Jupiter；测试类命名为 ClassNameTest，package 声明与源码保持一致。",
			"严格使用源码中声明的类名、方法名和 static/instance 调用方式。")
	case "cpp":
		rules = append(rules,
			"使用 GoogleTest 的 TEST、EXPECT_*、ASSERT_*。",
			"只包含测试所需头文件，不重复声明源码中已有的类或函数。")
	}
	var b strings.Builder
	fmt.Fprintf(&b, `<p>%s</p><ul>`, escapeHTML(modeText[mode]))
	for _, rule := range rules {
		fmt.Fprintf(&b, `<li>%s</li>`, escapeHTML(rule))
	}
	b.WriteString(`</ul>`)
	return b.String()
}

// buildChartScripts 生成图表脚本
// buildInteractiveScripts 生成图表脚本
func buildInteractiveScripts(payload contracts.ReportPayload, rows []contracts.EvaluationResult) string {
	topModelsJSON := marshalJSONSimple(payload.TopModels)
	rowsJSON := marshalJSONSimple(rows)
	return "<script>\n" + BuildChartsJS(topModelsJSON, rowsJSON) + "\n</script>"
}

func buildChartScripts(models []contracts.ModelRank) string {
	if len(models) == 0 {
		return ""
	}

	modelNames, compileRates, testRates, lineCovs, mutScores := extractChartDataSimple(models)

	return fmt.Sprintf(`
<script>
const modelNames = %s;
const compileRates = %s;
const testPassRates = %s;
const lineCoverages = %s;
const mutationScores = %s;

new Chart(document.getElementById('modelBarChart'), {
  type: 'bar',
  data: {
    labels: modelNames,
    datasets: [
      { label: '编译', data: compileRates, backgroundColor: '#3b82f6' },
      { label: '测试通过', data: testPassRates, backgroundColor: '#10b981' },
      { label: '行覆盖', data: lineCoverages, backgroundColor: '#f59e0b' },
      { label: '变异', data: mutationScores, backgroundColor: '#8b5cf6' }
    ]
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: 'top' },
      tooltip: {
        callbacks: {
          label: function(ctx) {
            return ctx.dataset.label + ': ' + (ctx.raw * 100).toFixed(1) + '%%';
          }
        }
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        max: 1,
        ticks: {
          callback: function(value) {
            return (value * 100).toFixed(0) + '%%';
          }
        }
      }
    }
  }
});

new Chart(document.getElementById('radarChart'), {
  type: 'radar',
  data: {
    labels: ['编译', '测试通过', '行覆盖', '变异'],
    datasets: modelNames.map((name, i) => ({
      label: name,
      data: [compileRates[i], testPassRates[i], lineCoverages[i], mutationScores[i]],
      fill: true,
      backgroundColor: ['rgba(59, 130, 246, 0.2)', 'rgba(16, 185, 129, 0.2)', 'rgba(245, 158, 11, 0.2)', 'rgba(139, 92, 246, 0.2)'][i %% 4],
      borderColor: ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6'][i %% 4],
      pointBackgroundColor: ['#3b82f6', '#10b981', '#f59e0b', '#8b5cf6'][i %% 4],
    }))
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: 'top' },
      tooltip: {
        callbacks: {
          label: function(ctx) {
            return ctx.dataset.label + ': ' + (ctx.raw * 100).toFixed(1) + '%%';
          }
        }
      }
    },
    scales: {
      r: {
        beginAtZero: true,
        max: 1,
        ticks: {
          callback: function(value) {
            return (value * 100).toFixed(0) + '%%';
          }
        }
      }
    }
  }
});
</script>
`, modelNames, compileRates, testRates, lineCovs, mutScores)
}
