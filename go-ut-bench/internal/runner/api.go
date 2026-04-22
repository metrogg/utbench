package runner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"
)

type apiClient struct {
	client  *http.Client
	retries int
	backoff time.Duration
}

func newAPIClient() *apiClient {
	return &apiClient{
		client:  &http.Client{Timeout: 120 * time.Second},
		retries: 3,
		backoff: 2 * time.Second,
	}
}

func (c *apiClient) generateTest(
	ctx context.Context,
	model modelConfig,
	language string,
	samplePath string,
	sourceCode string,
) (string, map[string]any, int, *int, *int, *int, bool, *contracts.ErrorInfo) {
	apiKey := strings.TrimSpace(os.Getenv(model.APIKeyEnv))
	if apiKey == "" {
		return "", nil, 0, nil, nil, nil, false, &contracts.ErrorInfo{
			Kind:      "auth_config_error",
			Message:   fmt.Sprintf("missing API key env var: %s (model=%s)", model.APIKeyEnv, model.Name),
			Retryable: false,
		}
	}

	prompt := buildPrompt(language, samplePath, sourceCode)
	payload := buildPayload(model, prompt)
	body, err := json.Marshal(payload)
	if err != nil {
		return "", nil, 0, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "payload_error", Message: err.Error(), Retryable: false}
	}

	endpoint := resolveEndpoint(model)
	var lastErr *contracts.ErrorInfo
	var lastTruncated bool
	for attempt := 1; attempt <= c.retries; attempt++ {
		started := time.Now()
		code, rawResp, p, cm, total, truncated, errInfo := c.doOnce(ctx, endpoint, apiKey, model.Provider, body)
		latency := int(time.Since(started).Milliseconds())
		if errInfo == nil {
			san := sanitizeModelOutput(code, language)
			extracted := extractCode(san, language)
			if vErr := validateGeneratedTest(extracted, language); vErr != nil {
				lastErr = &contracts.ErrorInfo{Kind: "quality_error", Message: vErr.Error(), Retryable: false}
				break
			}
			return extracted, rawResp, latency, p, cm, total, truncated, nil
		}

		lastErr = errInfo
		lastTruncated = truncated
		if !errInfo.Retryable || attempt >= c.retries {
			break
		}
		sleep := float64(c.backoff) * math.Pow(2, float64(attempt-1))
		sleep += float64(time.Duration(rand.Int63n(int64(200 * time.Millisecond))))
		time.Sleep(time.Duration(sleep))
	}

	if lastErr == nil {
		lastErr = &contracts.ErrorInfo{Kind: "unknown_error", Message: "unknown generation error", Retryable: false}
	}
	return "", nil, 0, nil, nil, nil, lastTruncated, lastErr
}

func (c *apiClient) doOnce(
	ctx context.Context,
	endpoint string,
	apiKey string,
	provider string,
	body []byte,
) (string, map[string]any, *int, *int, *int, bool, *contracts.ErrorInfo) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", nil, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "request_build_error", Message: err.Error(), Retryable: false}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		msg := err.Error()
		kind := "network_error"
		retryable := true
		if strings.Contains(strings.ToLower(msg), "timeout") {
			kind = "timeout"
		}
		return "", nil, nil, nil, nil, false, &contracts.ErrorInfo{Kind: kind, Message: msg, Retryable: retryable}
	}
	defer resp.Body.Close()

	rawBytes, _ := io.ReadAll(resp.Body)
	rawText := string(rawBytes)

	if resp.StatusCode >= 400 {
		retryable := resp.StatusCode == 408 || resp.StatusCode == 429 || resp.StatusCode >= 500
		code := resp.StatusCode
		return "", nil, nil, nil, nil, false, &contracts.ErrorInfo{
			Kind:       "http_error",
			Message:    fmt.Sprintf("http %d: %s", resp.StatusCode, trimText(rawText, 500)),
			Retryable:  retryable,
			StatusCode: &code,
		}
	}

	var payload map[string]any
	if err := json.Unmarshal(rawBytes, &payload); err != nil {
		return "", nil, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "response_parse_error", Message: err.Error(), Retryable: false}
	}

	text, err := extractResponseText(payload, provider)
	if err != nil {
		return "", payload, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "response_extract_error", Message: err.Error(), Retryable: false}
	}
	promptTokens, completionTokens, totalTokens := extractUsage(payload)
	truncated := extractFinishReason(payload, provider)
	return text, payload, promptTokens, completionTokens, totalTokens, truncated, nil
}

func resolveEndpoint(model modelConfig) string {
	base := strings.TrimSuffix(model.Endpoint, "/")
	if model.Provider == "dashscope" {
		if strings.Contains(base, "compatible-mode") {
			return base + "/chat/completions"
		}
		if strings.HasSuffix(base, "/api/v1") {
			return base + "/services/aigc/text-generation/generation"
		}
		return base + "/services/aigc/text-generation/generation"
	}
	return base + "/chat/completions"
}

func buildPayload(model modelConfig, prompt string) map[string]any {
	params := map[string]any{}
	for k, v := range model.Params {
		params[k] = v
	}

	if model.Provider == "dashscope" && !strings.Contains(model.Endpoint, "compatible-mode") {
		return map[string]any{
			"model":      model.Model,
			"input":      map[string]any{"messages": []map[string]any{{"role": "user", "content": prompt}}},
			"parameters": params,
		}
	}

	payload := map[string]any{
		"model": model.Model,
		"messages": []map[string]any{
			{"role": "system", "content": "You generate high-quality unit tests."},
			{"role": "user", "content": prompt},
		},
		"stream": false,
	}
	for k, v := range params {
		payload[k] = v
	}
	return payload
}

func extractResponseText(response map[string]any, provider string) (string, error) {
	if provider == "dashscope" {
		if output, ok := response["output"].(map[string]any); ok {
			if text, ok := output["text"].(string); ok && text != "" {
				return text, nil
			}
			if choices, ok := output["choices"].([]any); ok && len(choices) > 0 {
				if item, ok := choices[0].(map[string]any); ok {
					if msg, ok := item["message"].(map[string]any); ok {
						if content, ok := msg["content"].(string); ok {
							return content, nil
						}
					}
				}
			}
		}
	}

	if choices, ok := response["choices"].([]any); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]any); ok {
			if msg, ok := choice["message"].(map[string]any); ok {
				if content, ok := msg["content"].(string); ok {
					return content, nil
				}
			}
			if text, ok := choice["text"].(string); ok {
				return text, nil
			}
		}
	}
	return "", fmt.Errorf("unable to extract text from response")
}

func extractUsage(response map[string]any) (*int, *int, *int) {
	usage, ok := response["usage"].(map[string]any)
	if !ok {
		return nil, nil, nil
	}
	p := toIntPtr(usage["prompt_tokens"])
	if p == nil {
		p = toIntPtr(usage["input_tokens"])
	}
	c := toIntPtr(usage["completion_tokens"])
	if c == nil {
		c = toIntPtr(usage["output_tokens"])
	}
	t := toIntPtr(usage["total_tokens"])
	return p, c, t
}

func extractFinishReason(response map[string]any, provider string) bool {
	if provider == "dashscope" {
		if output, ok := response["output"].(map[string]any); ok {
			if choices, ok := output["choices"].([]any); ok && len(choices) > 0 {
				if item, ok := choices[0].(map[string]any); ok {
					if fr, ok := item["finish_reason"].(string); ok {
						return fr == "length"
					}
				}
			}
		}
		return false
	}

	if choices, ok := response["choices"].([]any); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]any); ok {
			if fr, ok := choice["finish_reason"].(string); ok {
				return fr == "length"
			}
		}
	}
	return false
}

func toIntPtr(v any) *int {
	switch x := v.(type) {
	case int:
		return &x
	case int64:
		v := int(x)
		return &v
	case float64:
		v := int(x)
		return &v
	default:
		return nil
	}
}

func sanitizeModelOutput(content string, language string) string {
	text := strings.TrimSpace(content)
	re := regexp.MustCompile(`(?is)<think>.*?</think>`)
	text = re.ReplaceAllString(text, "")
	re2 := regexp.MustCompile(`(?is)<analysis>.*?</analysis>`)
	text = re2.ReplaceAllString(text, "")
	if language == "python" {
		text = trimNonCodePrefix(text)
	}
	return strings.TrimSpace(text)
}

func trimNonCodePrefix(text string) string {
	lines := strings.Split(text, "\n")
	re := regexp.MustCompile(`^\s*(from\s+\w|import\s+\w|def\s+\w|class\s+\w|@|if\s+__name__|#|\"\"\"|''')`)
	for i, line := range lines {
		if re.MatchString(line) {
			return strings.TrimSpace(strings.Join(lines[i:], "\n"))
		}
	}
	return text
}

func extractCode(content, language string) string {
	re := regexp.MustCompile("```([a-zA-Z0-9_+\\-]*)\\s*\\n([\\s\\S]*?)```")
	blocks := re.FindAllStringSubmatch(content, -1)
	if len(blocks) == 0 {
		return stripMarkdownFence(content)
	}
	aliases := map[string]bool{strings.ToLower(language): true}
	if strings.EqualFold(language, "python") {
		aliases["py"] = true
	}
	if strings.EqualFold(language, "cpp") {
		aliases["c++"] = true
	}
	for _, block := range blocks {
		tag := strings.ToLower(strings.TrimSpace(block[1]))
		if aliases[tag] {
			return strings.TrimSpace(block[2])
		}
	}
	return strings.TrimSpace(blocks[0][2])
}

func stripMarkdownFence(content string) string {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "```") {
		return dropTrailingFenceLines(trimmed)
	}
	lines := strings.Split(trimmed, "\n")
	if len(lines) <= 2 {
		return trimmed
	}
	lines = lines[1:]
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	if len(lines) > 0 && strings.HasPrefix(strings.TrimSpace(lines[len(lines)-1]), "```") {
		lines = lines[:len(lines)-1]
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func dropTrailingFenceLines(text string) string {
	lines := strings.Split(text, "\n")
	for len(lines) > 0 {
		last := strings.TrimSpace(lines[len(lines)-1])
		if strings.HasPrefix(last, "```") {
			lines = lines[:len(lines)-1]
			continue
		}
		break
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func validateGeneratedTest(code, language string) error {
	stripped := strings.TrimSpace(code)
	if stripped == "" {
		return fmt.Errorf("empty output")
	}
	lower := strings.ToLower(stripped)
	if strings.Contains(lower, "<think>") || strings.Contains(lower, "</think>") {
		return fmt.Errorf("contains leaked reasoning tags")
	}
	if strings.EqualFold(language, "python") {
		if !strings.Contains(stripped, "def test_") && !strings.Contains(stripped, "import pytest") {
			return fmt.Errorf("invalid python test structure")
		}
	}
	return nil
}

func trimText(v string, max int) string {
	if len(v) <= max {
		return v
	}
	return v[:max]
}

func buildPrompt(language, samplePath, sourceCode string) string {
	meta := loadModuleLevelMetaForRunner(samplePath)
	if meta != nil {
		return buildModuleLevelPrompt(language, samplePath, sourceCode, meta)
	}

	lang := strings.ToLower(strings.TrimSpace(language))
	sampleID := strings.TrimSuffix(filepath.Base(samplePath), filepath.Ext(samplePath))
	framework := languageFramework(lang)
	dependencies := extractDependencies(sourceCode, lang)
	scenario, complexity := parseSampleMeta(sampleID, samplePath)
	moduleName := moduleImportName(sampleID)
	coverageTargets := coverageTargetsText()
	mockReq := mockRequirement(sourceCode)
	criticalConditions := extractCriticalConditions(sourceCode, lang)
	moduleSymbols := extractModuleLevelSymbols(sourceCode, lang)

	dependencyText := "none detected"
	if len(dependencies) > 0 {
		dependencyText = strings.Join(dependencies, ", ")
	}
	moduleSymbolsText := ""
	if moduleSymbols != "" {
		moduleSymbolsText = "\n" + moduleSymbols
	}
	scenarioText := scenario
	if scenarioText == "" {
		scenarioText = "general"
	}
	complexityText := complexity
	if complexityText == "" {
		complexityText = "unknown"
	}

	criticalText := "- none detected"
	if len(criticalConditions) > 0 {
		lines := make([]string, 0, len(criticalConditions))
		for _, item := range criticalConditions {
			lines = append(lines, fmt.Sprintf("- `%s`", item))
		}
		criticalText = strings.Join(lines, "\n")
	}

	return "You are an expert unit testing engineer.\n" +
		"你是一名资深单元测试工程师。\n" +
		"Generate high-quality unit tests based on the following specification.\n" +
		"请基于以下规范生成高质量单元测试。\n\n" +
		"## Role & Objective（角色与目标）\n" +
		"- Goal: produce executable tests that match source behavior exactly.\n" +
		"- 目标：生成可执行且与源码行为严格一致的测试。\n\n" +
		"## Language & Framework（语言与框架）\n" +
		fmt.Sprintf("- Language（语言）: %s\n", lang) +
		fmt.Sprintf("- Test Framework（测试框架）: %s\n\n", framework) +
		"## Step-by-Step Workflow（分步流程）\n" +
		"1) Identify callable symbols and input/output contracts from source.\n" +
		"2) Build a test matrix: normal, boundary, and exception paths.\n" +
		"3) Derive expected values only from implementation semantics.\n" +
		"4) Write deterministic, runnable tests with clear assertions.\n" +
		"5) Self-check syntax/imports/assertions before final output.\n\n" +
		"## Test Requirements（测试要求）\n" +
		"- Cover normal paths, boundary conditions, and error/exception behavior.\n" +
		"- 覆盖正常路径、边界条件和异常行为。\n" +
		"- Keep tests deterministic and runnable.\n" +
		"- 保持测试可重复、可执行（避免随机性）。\n" +
		"- Use clear assertions with meaningful expected values.\n" +
		"- 使用清晰断言和有意义的期望值。\n" +
		buildLanguageSpecificRules(lang, moduleName) +
		"- Avoid importing unrelated third-party packages by default.\n" +
		"- 默认不要引入无关第三方包。\n" +
		fmt.Sprintf("- Coverage targets（覆盖率目标，供参考）: %s\n", coverageTargets) +
		fmt.Sprintf("- Mock requirements（Mock 要求）: %s\n\n", mockReq) +
		"## Semantic Alignment Hard Rules（语义对齐硬约束）\n" +
		"- Derive expected values strictly from the given source code behavior.\n" +
		"- 期望值必须严格依据给定源码行为推导，不要按题型常识脑补。\n" +
		"- Respect exact comparison semantics in code (`<`, `<=`, `>`, `>=`, `==`).\n" +
		"- 必须严格遵守源码比较符号语义（尤其阈值边界等于时）。\n" +
		"- Include explicit boundary-equality assertions when threshold/limit checks exist.\n" +
		"- 当存在阈值/边界判断时，必须包含“等于边界”的断言样例。\n" +
		"- If implementation looks counter-intuitive, still assert implementation behavior.\n" +
		"- 若实现与常识不一致，也必须以源码实现为准。\n" +
		"- If docstring/comment conflicts with implementation, trust implementation.\n" +
		"- 若注释/文档示例与实现冲突，以实现为准。\n" +
		"- Do NOT assume implicit coercion not present in code (e.g., str->number).\n" +
		"- 不要假设源码未实现的隐式转换（例如字符串自动转数字）。\n" +
		"- For sliding-window or two-pointer counting algorithms: when the loop condition is\n" +
		"  `while left < right and sorted[right] - sorted[left] >= threshold`,\n" +
		"  a threshold of 0 with duplicate values produces a count of 0 — the window only\n" +
		"  advances when `>` threshold, NOT when `==` threshold.\n" +
		"- 对于滑窗/双指针计数算法：当循环条件是 `>= threshold` 时，\n" +
		"  threshold=0 且有重复值的情况下会计数为 0 —— 只有 `>` threshold 时窗口才右移。\n" +
		"- Trace through the algorithm by hand for boundary values before writing assertions.\n" +
		"- 写断言前，必须手工推导一遍算法在边界值上的执行过程。\n" +
		"- For functions returning structured results (e.g., namedtuple, dataclass): always assert\n" +
		"  each field individually. Never assert the whole object equality without field-level checks.\n" +
		"- 对于返回结构体结果的函数，必须逐字段断言，切勿直接做整体相等判断而不验证字段值。\n" +
		"- For `nearest_pair` or similar: the returned index fields are `left_index=min(original_indices)`\n" +
		"  and `right_index=max(original_indices)` — trace the sorted enumeration carefully.\n" +
		"- 对于类似 `nearest_pair` 的函数：返回的索引字段是 `left_index=min(原始索引)` 和\n" +
		"  `right_index=max(原始索引)`，必须仔细追踪排序后的枚举过程。\n" +
		"- For functions returning collections where order may be non-deterministic (e.g., `most_common`\n" +
		"  with ties, `sorted` with equal keys, `dict.items()`, set iterations): either (a) use test\n" +
		"  inputs that produce a deterministic ordering, or (b) assert on properties rather than exact\n" +
		"  order/position, or (c) check that the result CONTAINS the expected elements without\n" +
		"  assuming their order.\n" +
		"- 对于返回集合且顺序不确定的函数（如 ties 时的 `most_common`、equal keys 时的 `sorted`、\n" +
		"  `dict.items()`、set 迭代器）：(a) 使用能产生确定性顺序的输入，或 (b) 断言属性而非精确\n" +
		"  顺序/位置，或 (c) 检查结果包含期望元素而不假设其顺序。\n" +
		"- For file/directory operations: trace ALL branches of the conditional logic. If the code\n" +
		"  has two distinct handling paths (e.g., 'invalid filename' vs 'valid filename'), your tests\n" +
		"  must cover BOTH paths with appropriate expected outcomes for each.\n" +
		"- 对于文件/目录操作类函数：必须追踪所有条件分支。如果代码有两条处理路径（如\"无效文件名\"\n" +
		"  vs \"有效文件名\"），测试必须覆盖两条路径，并为每条路径设置正确的期望值。\n" +
		"- DO NOT assume file extension classification behavior without reading the source code\n" +
		"  carefully. The actual behavior may differ from intuition.\n" +
		"- 不要凭直觉假设文件扩展名分类行为，必须仔细阅读源码后才能确定实际行为。\n" +
		"- For functions with side effects (file moves, directory creation): verify the EXACT\n" +
		"  destination path and naming logic from the source, not from assumptions.\n" +
		"- 对于有副作用的函数（文件移动、目录创建）：必须从源码确认精确的目标路径和命名逻辑，不要假设。\n" +
		"- When source code has `if/elif/else` chains: create at least one test case for EACH branch.\n" +
		"- 当源码有 `if/elif/else` 链时：必须为每个分支至少创建一个测试用例。\n\n" +
		"## Critical Conditions Extracted（关键逻辑条件）\n" +
		criticalText +
		"\n\n" +
		"## Error Prevention Checklist（错误预防清单，仅内部执行）\n" +
		"- No placeholder tests like `assert True`.\n" +
		"- No assertions for behavior that cannot be inferred from source.\n" +
		"- Ensure every referenced symbol exists in source imports/definitions.\n" +
		"- Ensure generated file is directly runnable by the target test framework.\n" +
		"- If you use any Python standard library modules (string, zipfile, json, os, sys, tempfile, etc.) or third-party packages in your test code, you MUST explicitly import them.\n" +
		"- 如果测试中使用了任何 Python 标准库模块（string、zipfile、json、os、sys、tempfile 等）或第三方包，必须显式 import。\n\n" +
		"## Context Information（上下文信息）\n" +
		fmt.Sprintf("- Sample ID（样本ID）: %s\n", sampleID) +
		fmt.Sprintf("- Scenario（场景）: %s\n", scenarioText) +
		fmt.Sprintf("- Complexity（复杂度）: %s\n", complexityText) +
		fmt.Sprintf("- Dependencies detected（检测到依赖）: %s\n%s\n\n", dependencyText, moduleSymbolsText) +
		"## Output Format（输出格式）\n" +
		"- Return raw test code only (no Markdown fences).\n" +
		"- 仅输出原始测试代码，不要 Markdown 代码块。\n" +
		"- Do not include explanations.\n" +
		"- 不要输出解释文字。\n\n" +
		"## Source Code Under Test（被测源码）\n" +
		fmt.Sprintf("```%s\n%s\n```", lang, sourceCode)
}

func buildLanguageSpecificRules(lang, moduleName string) string {
	switch lang {
	case "python":
		return "- Test function names must start with `test_` (e.g., `def test_xxx()`).\n" +
			"- 测试函数命名必须以 `test_` 开头（如 `def test_xxx()`）。\n" +
			"- Use plain `assert` and `pytest.raises` for failure paths.\n" +
			"- 断言使用 `assert`，异常路径使用 `pytest.raises`。\n" +
			"- Use function-based tests (e.g., `def test_xxx()`) instead of class-based.\n" +
			"- 使用函数式测试（如 `def test_xxx()`），不要用 class-based。\n" +
			fmt.Sprintf("- MUST import target symbols from local module `%s` before writing tests.\n", moduleName) +
			fmt.Sprintf("- 必须先从同目录模块 `%s` 导入被测对象，再编写测试。\n", moduleName)

	case "go":
		return "- Test function names MUST start with `Test` (capital T) followed by a name (e.g., `func TestAdd(t *testing.T)`).\n" +
			"- 测试函数命名必须以 `Test`（大写T）开头（如 `func TestAdd(t *testing.T)`）。\n" +
			"- The test file must be named `xxx_test.go` where `xxx` matches the source file name.\n" +
			"- 测试文件命名必须为 `xxx_test.go`，其中 `xxx` 与源文件名匹配。\n" +
			"- Use `t.Error`, `t.Errorf`, `t.Fatal`, `t.Fatalf` for assertions, NOT plain `assert`.\n" +
			"- 断言使用 `t.Error`、`t.Errorf`、`t.Fatal`、`t.Fatalf`，不要用 `assert`。\n" +
			"- The test function parameter must be `t *testing.T`.\n" +
			"- 测试函数参数必须是 `t *testing.T`。\n" +
			"- Use `testing` package import: `import \"testing\"`.\n" +
			"- 导入测试包：`import \"testing\"`。\n" +
			"- Package name in test file should match source file package (usually same directory).\n" +
			"- 测试文件包名应与源文件包名一致（通常在同目录）。\n" +
			"- CRITICAL: Before calling ANY function in your test, verify it EXISTS in the source code.\n" +
			"  Check the source code provided above for the exact function signature. Do NOT guess or\n" +
			"  assume a function exists if you cannot see its definition in the source.\n" +
			"- 重要：调用任何函数前，必须验证它存在于源码中。检查上面提供的源码中的确切函数签名。\n" +
			"  如果源码中找不到函数定义，不要猜测或假设函数存在。\n"

	case "java":
		return "- Test class name MUST be `ClassNameTest` where `ClassName` matches the source class name.\n" +
			"- 测试类命名必须为 `ClassNameTest`，其中 `ClassName` 与源类名匹配。\n" +
			"- Test file MUST be named `ClassNameTest.java` (e.g., `EmployeeValidatorTest.java`).\n" +
			"- 测试文件命名必须为 `ClassNameTest.java`（如 `EmployeeValidatorTest.java`）。\n" +
			"- Use JUnit 4 annotations: `@Test` for test methods.\n" +
			"- 使用 JUnit 4 注解：测试方法使用 `@Test`。\n" +
			"- Test methods must be `public void` and can be named `testXxx()` or descriptively.\n" +
			"- 测试方法必须是 `public void`，命名可以是 `testXxx()` 或有描述性的名称。\n" +
			"- Import JUnit: `import org.junit.Test; import static org.junit.Assert.*;`.\n" +
			"- 导入 JUnit：`import org.junit.Test; import static org.junit.Assert.*;`。\n" +
			"- Use `assertEquals(expected, actual)`, `assertTrue(condition)`, `assertFalse(condition)`.\n" +
			"- 断言使用 `assertEquals(expected, actual)`、`assertTrue(condition)`、`assertFalse(condition)`。\n" +
			"- For exceptions, use `@Test(expected = IllegalArgumentException.class)` or try-catch.\n" +
			"- 异常测试使用 `@Test(expected = IllegalArgumentException.class)` 或 try-catch。\n" +
			"- For static methods, call them directly: `ClassName.methodName(args)`.\n" +
			"- 静态方法直接调用：`ClassName.methodName(args)`。\n" +
			"- Test class must be `public` and in the same package as source (no package declaration for simple files).\n" +
			"- 测试类必须是 `public`，与源类在同一包中（简单文件无需包声明）。\n" +
			"- For nested/inner classes in source: import them or reference as `OuterClass.InnerClass`.\n" +
			"- 对于源码中的嵌套类：导入它们或使用 `OuterClass.InnerClass` 引用。\n" +
			"- DO NOT access private fields/methods directly. Use public API or reflection only if necessary.\n" +
			"- 不要直接访问私有字段/方法。使用公共 API，必要时才用反射。\n" +
			"- Output COMPLETE test code. Ensure all braces `{}` are properly closed.\n" +
			"- 输出完整的测试代码。确保所有 `{}` 大括号正确闭合。\n" +
			"- Check your output: if source has inner classes, test must handle them correctly.\n" +
			"- 检查输出：如果源码有嵌套类，测试必须正确处理。\n"

	case "cpp":
		return "- Use GoogleTest framework macros: `TEST()`, `EXPECT_EQ`, `ASSERT_EQ`.\n" +
			"- 使用 GoogleTest 框架宏：`TEST()`、`EXPECT_EQ`、`ASSERT_EQ`。\n" +
			"- Test function naming: `TEST(TestSuiteName, TestName)`.\n" +
			"- 测试函数命名：`TEST(TestSuiteName, TestName)`。\n" +
			"- Include `<gtest/gtest.h>` header.\n" +
			"- 包含 `<gtest/gtest.h>` 头文件。\n"

	case "javascript":
		return "- Test function names should use `test()` or `it()` from Jest.\n" +
			"- 测试函数使用 Jest 的 `test()` 或 `it()`。\n" +
			"- Use `expect()` for assertions (e.g., `expect(result).toBe(expected)`).\n" +
			"- 断言使用 `expect()`（如 `expect(result).toBe(expected)`）。\n" +
			"- Test file should be named `xxx.test.js` or `xxx.spec.js`.\n" +
			"- 测试文件命名应为 `xxx.test.js` 或 `xxx.spec.js`。\n"

	default:
		return "- Follow standard testing conventions for the language.\n" +
			"- 遵循该语言的标准测试规范。\n"
	}
}

func coverageTargetsText() string {
	// Keep in sync with benchmark/config/models.yaml default thresholds.
	line := 0.7
	branch := 0.6
	function := 0.8
	return fmt.Sprintf(
		"line >= %.0f%%, branch >= %.0f%%, function >= %.0f%%",
		line*100,
		branch*100,
		function*100,
	)
}

func languageFramework(language string) string {
	switch language {
	case "java":
		return "JUnit 4"
	case "python":
		return "pytest"
	case "go":
		return "Go testing package"
	case "cpp":
		return "GoogleTest"
	case "javascript":
		return "Jest"
	default:
		return "the standard test framework"
	}
}

func moduleImportName(sampleID string) string {
	normalized := regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(sampleID, "_")
	normalized = strings.Trim(normalized, "_")
	if normalized == "" {
		return "solution"
	}
	if regexp.MustCompile(`^[0-9]`).MatchString(normalized) {
		return "sample_" + normalized
	}
	return normalized
}

func parseSampleMeta(sampleID, samplePath string) (string, string) {
	parts := strings.Split(sampleID, "_")
	if len(parts) >= 4 {
		complexity := strings.ToLower(parts[0])
		scenario := strings.ToLower(strings.Join(parts[2:len(parts)-1], "_"))
		return scenario, complexity
	}
	parent := strings.ToLower(filepath.Base(filepath.Dir(samplePath)))
	grand := strings.ToLower(filepath.Base(filepath.Dir(filepath.Dir(samplePath))))
	if parent != "" && parent != grand {
		return parent, ""
	}
	return "", ""
}

func extractDependencies(sourceCode, language string) []string {
	var deps []string
	switch language {
	case "python":
		re := regexp.MustCompile(`(?m)^\s*(?:from\s+([a-zA-Z0-9_\.]+)\s+import|import\s+([a-zA-Z0-9_\.]+))`)
		matches := re.FindAllStringSubmatch(sourceCode, -1)
		for _, m := range matches {
			dep := strings.TrimSpace(m[1])
			if dep == "" {
				dep = strings.TrimSpace(m[2])
			}
			if dep != "" {
				deps = append(deps, dep)
			}
		}
	case "java":
		re := regexp.MustCompile(`(?m)^\s*import\s+([^;]+);`)
		matches := re.FindAllStringSubmatch(sourceCode, -1)
		for _, m := range matches {
			if len(m) > 1 && strings.TrimSpace(m[1]) != "" {
				deps = append(deps, strings.TrimSpace(m[1]))
			}
		}
	case "go":
		reSingle := regexp.MustCompile(`(?m)^\s*import\s+"([^"]+)"`)
		singleMatches := reSingle.FindAllStringSubmatch(sourceCode, -1)
		for _, m := range singleMatches {
			if len(m) > 1 && strings.TrimSpace(m[1]) != "" {
				deps = append(deps, strings.TrimSpace(m[1]))
			}
		}
		reBlock := regexp.MustCompile(`(?s)import\s*\((.*?)\)`)
		block := reBlock.FindStringSubmatch(sourceCode)
		if len(block) > 1 {
			reQuoted := regexp.MustCompile(`"([^"]+)"`)
			quoted := reQuoted.FindAllStringSubmatch(block[1], -1)
			for _, m := range quoted {
				if len(m) > 1 && strings.TrimSpace(m[1]) != "" {
					deps = append(deps, strings.TrimSpace(m[1]))
				}
			}
		}
	case "cpp":
		re := regexp.MustCompile(`(?m)^\s*#include\s*[<"]([^>"]+)[>"]`)
		matches := re.FindAllStringSubmatch(sourceCode, -1)
		for _, m := range matches {
			if len(m) > 1 && strings.TrimSpace(m[1]) != "" {
				deps = append(deps, strings.TrimSpace(m[1]))
			}
		}
	}

	seen := make(map[string]struct{}, len(deps))
	out := make([]string, 0, len(deps))
	for _, dep := range deps {
		if _, ok := seen[dep]; ok {
			continue
		}
		seen[dep] = struct{}{}
		out = append(out, dep)
		if len(out) >= 12 {
			break
		}
	}
	return out
}

func mockRequirement(sourceCode string) string {
	lower := strings.ToLower(sourceCode)
	markers := []string{
		"http",
		"request",
		"socket",
		"open(",
		"file",
		"database",
		"sql",
		"redis",
		"grpc",
		"client",
		"os.environ",
		"subprocess",
	}
	for _, marker := range markers {
		if strings.Contains(lower, marker) {
			return "Use mocks/stubs/fakes for network, file, database, subprocess, or env dependencies."
		}
	}
	return "Mock only when necessary; avoid over-mocking pure functions."
}

func extractCriticalConditions(sourceCode, language string) []string {
	if language != "python" {
		return nil
	}

	lines := strings.Split(sourceCode, "\n")
	candidates := make([]string, 0, 16)
	loopLines := map[int]string{}

	for idx, raw := range lines {
		lineno := idx + 1
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if containsComparator(line) {
			if strings.HasPrefix(line, "if ") ||
				strings.HasPrefix(line, "while ") ||
				strings.HasPrefix(line, "return ") {
				candidates = append(candidates, line)
			} else {
				pos := strings.Index(sourceCode, line)
				prefix := sourceCode
				if pos > 0 {
					prefix = sourceCode[:pos]
				}
				if strings.Contains(prefix, "while ") && strings.Contains(line, "count") {
					candidates = append(candidates, line)
				}
			}
		}

		if (strings.HasPrefix(line, "while ") || strings.HasPrefix(line, "for ")) && containsComparator(line) {
			loopLines[lineno] = line
		}
	}

	if len(loopLines) > 0 && len(candidates) == 0 {
		lineNos := make([]int, 0, len(loopLines))
		for lineNo := range loopLines {
			lineNos = append(lineNos, lineNo)
		}
		sort.Ints(lineNos)
		for _, lineNo := range lineNos {
			candidates = append(candidates, loopLines[lineNo])
		}
	}

	dedup := make([]string, 0, 12)
	seen := map[string]struct{}{}
	for _, item := range candidates {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		dedup = append(dedup, item)
		if len(dedup) >= 12 {
			break
		}
	}

	return dedup
}

func containsComparator(line string) bool {
	comparators := []string{"<=", ">=", "==", "!=", "<", ">"}
	for _, item := range comparators {
		if strings.Contains(line, item) {
			return true
		}
	}
	return false
}

func extractModuleLevelSymbols(sourceCode, language string) string {
	var symbols []string
	seen := make(map[string]struct{})

	addSymbol := func(name string) {
		if name == "" || strings.HasPrefix(name, "_") {
			return
		}
		if _, ok := seen[name]; ok {
			return
		}
		seen[name] = struct{}{}
		symbols = append(symbols, name)
	}

	switch language {
	case "go":
		funcPattern := regexp.MustCompile(`(?m)^\s*func\s+(\([a-zA-Z\s]+\*?[a-zA-Z_][a-zA-Z0-9_]*\)\s+)?([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`)
		for _, m := range funcPattern.FindAllStringSubmatch(sourceCode, -1) {
			if len(m) > 2 {
				addSymbol(m[2])
			}
		}

		typePattern := regexp.MustCompile(`(?m)^\s*type\s+([a-zA-Z_][a-zA-Z0-9_]*)\s+(?:struct|interface|int|string|bool|float|byte|rune|error|map|chan)\b`)
		for _, m := range typePattern.FindAllStringSubmatch(sourceCode, -1) {
			if len(m) > 1 {
				addSymbol(m[1])
			}
		}

		typeBlockPattern := regexp.MustCompile(`(?m)^\s*type\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\{`)
		for _, m := range typeBlockPattern.FindAllStringSubmatch(sourceCode, -1) {
			if len(m) > 1 {
				addSymbol(m[1])
			}
		}

		constBlock := regexp.MustCompile(`(?s)const\s*\(([^)]+)\)`)
		for _, m := range constBlock.FindAllStringSubmatch(sourceCode, -1) {
			lineConst := regexp.MustCompile(`(?m)^\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*=`)
			for _, c := range lineConst.FindAllStringSubmatch(m[1], -1) {
				if len(c) > 1 {
					addSymbol(c[1])
				}
			}
		}

		varPattern := regexp.MustCompile(`(?m)^\s*var\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?:\[\]|map|chan|\*)`)
		for _, m := range varPattern.FindAllStringSubmatch(sourceCode, -1) {
			if len(m) > 1 {
				addSymbol(m[1])
			}
		}

	default:
		re := regexp.MustCompile("(?m)^([A-Za-z_][A-Za-z0-9_]*)\\s*=")
		for _, m := range re.FindAllStringSubmatch(sourceCode, -1) {
			if len(m) > 1 {
				addSymbol(m[1])
			}
		}
	}

	if len(symbols) == 0 {
		return ""
	}
	return "Module-level symbols available to import: " + strings.Join(symbols, ", ") + "."
}

type moduleLevelMetaForRunner struct {
	SampleID      string   `json:"sample_id"`
	ModuleImport  string   `json:"module_import"`
	PackageName   string   `json:"package_name"`
	TargetFile    string   `json:"target_file"`
	WorkspaceRoot string   `json:"workspace_root"`
	Requirements  []string `json:"requirements,omitempty"`
}

func loadModuleLevelMetaForRunner(samplePath string) *moduleLevelMetaForRunner {
	dir := filepath.Dir(samplePath)
	base := filepath.Base(samplePath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	var metaPath string
	if name == "entry" {
		metaPath = filepath.Join(dir, "meta.json")
	} else {
		metaPath = filepath.Join(dir, name+".meta.json")
	}

	if _, err := os.Stat(metaPath); err != nil {
		return nil
	}
	raw, err := os.ReadFile(metaPath)
	if err != nil {
		return nil
	}
	var meta moduleLevelMetaForRunner
	if err := json.Unmarshal(raw, &meta); err != nil {
		return nil
	}
	if meta.ModuleImport == "" {
		return nil
	}
	return &meta
}

func buildModuleLevelPrompt(language, samplePath, sourceCode string, meta *moduleLevelMetaForRunner) string {
	lang := strings.ToLower(strings.TrimSpace(language))
	framework := languageFramework(lang)
	coverageTargets := coverageTargetsText()
	sampleID := meta.SampleID
	if sampleID == "" {
		sampleID = strings.TrimSuffix(filepath.Base(samplePath), filepath.Ext(samplePath))
	}

	moduleImport := meta.ModuleImport
	packageName := meta.PackageName
	if packageName == "" {
		packageName = moduleImport
	}
	targetFile := meta.TargetFile
	requirements := meta.Requirements

	requirementsText := "none specified"
	if len(requirements) > 0 {
		requirementsText = strings.Join(requirements, ", ")
	}

	return "You are an expert unit testing engineer.\n" +
		"你是一名资深单元测试工程师。\n" +
		"Generate high-quality unit tests for a module-level (multi-file) package.\n" +
		"请为一个模块级（多文件）包生成高质量单元测试。\n\n" +
		"## Role & Objective（角色与目标）\n" +
		"- Goal: produce executable tests that match source behavior exactly.\n" +
		"- 目标：生成可执行且与源码行为严格一致的测试。\n" +
		"- This is a MODULE-LEVEL sample: the target code is part of a larger package.\n" +
		"- 这是模块级样本：被测代码是更大包的一部分。\n\n" +
		"## Language & Framework（语言与框架）\n" +
		fmt.Sprintf("- Language（语言）: %s\n", lang) +
		fmt.Sprintf("- Test Framework（测试框架）: %s\n\n", framework) +
		"## Module-Level Context（模块级上下文）\n" +
		fmt.Sprintf("- Target Module（目标模块）: `%s`\n", moduleImport) +
		fmt.Sprintf("- Package Name（包名）: `%s`\n", packageName) +
		fmt.Sprintf("- Target File（目标文件）: `%s`\n", targetFile) +
		fmt.Sprintf("- Requirements（依赖）: %s\n\n", requirementsText) +
		"## Test Requirements（测试要求）\n" +
		"- Import the target module using: `from " + moduleImport + " import ...`\n" +
		"- 使用 `from " + moduleImport + " import ...` 导入被测模块。\n" +
		"- Write tests for the main classes/functions in the target module.\n" +
		"- 为目标模块中的主要类/函数编写测试。\n" +
		"- Cover normal paths, boundary conditions, and error/exception behavior.\n" +
		"- 覆盖正常路径、边界条件和异常行为。\n" +
		"- Keep tests deterministic and runnable.\n" +
		"- 保持测试可重复、可执行（避免随机性）。\n" +
		"- Use clear assertions with meaningful expected values.\n" +
		"- 使用清晰断言和有意义的期望值。\n" +
		"- Test function names must start with `test_`.\n" +
		"- 测试函数命名必须以 `test_` 开头。\n" +
		"- Use plain `assert` and `pytest.raises` for failure paths.\n" +
		"- 断言使用 `assert`，异常路径使用 `pytest.raises`。\n" +
		"- Use function-based tests (e.g., `def test_xxx()`) instead of class-based.\n" +
		"- 使用函数式测试（如 `def test_xxx()`），不要用 class-based（如 `class Test:`）。\n" +
		fmt.Sprintf("- Coverage targets（覆盖率目标，供参考）: %s\n\n", coverageTargets) +
		"## Important Notes（重要提示）\n" +
		"- DO NOT import from relative paths like `from ._common import ...`.\n" +
		"- 不要从相对路径导入，如 `from ._common import ...`。\n" +
		"- Always use the full module import path provided above.\n" +
		"- 始终使用上面提供的完整模块导入路径。\n" +
		"- The test file will be placed in a `tests/` subdirectory of the workspace.\n" +
		"- 测试文件将放在 workspace 的 `tests/` 子目录中。\n\n" +
		"## Output Format（输出格式）\n" +
		"- Return raw test code only (no Markdown fences).\n" +
		"- 仅输出原始测试代码，不要 Markdown 代码块。\n" +
		"- Do not include explanations.\n" +
		"- 不要输出解释文字。\n\n" +
		"## Sample Entry File（样本入口文件）\n" +
		"This file shows the import structure:\n" +
		fmt.Sprintf("```%s\n%s\n```", lang, sourceCode)
}
