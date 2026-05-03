// runner/api.go 提供 LLM API 调用功能
// 封装 HTTP 客户端，处理请求构建、响应解析、错误处理和自动续写
// 支持多种 API 提供商：OpenAI、Dashscope、Volcengine 等
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
	"sync"
	"time"

	"go-ut-bench/internal/contracts"
)

// apiClient LLM API 客户端
// 管理 HTTP 连接、重试策略和退避时间
type apiClient struct {
	client  *http.Client  // HTTP 客户端，设置超时时间
	retries int           // 最大重试次数
	backoff time.Duration // 基础退避时间
}

// 全局模型速率限制器
// 用于错开对同一模型提供者/端点的调用，避免触发速率限制
var (
	globalRateLimiter sync.Mutex                   // 全局互斥锁，保护模型调用时间记录
	modelLastCall     = make(map[string]time.Time) // 各模型最后调用时间
	modelMinInterval  = 200 * time.Millisecond     // 同一模型调用最小间隔
	modelJitter       = 100 * time.Millisecond     // 最大随机抖动时间
)

// newAPIClient 创建新的 API 客户端
// 使用默认配置：300秒超时、3次重试、2秒退避
func newAPIClient() *apiClient {
	return &apiClient{
		client:  &http.Client{Timeout: 300 * time.Second},
		retries: 3,
		backoff: 2 * time.Second,
	}
}

// waitModelInterval 确保对同一模型的调用是错开的
// 如果同一模型刚被调用，会等待直到最小间隔时间过去
//
// 参数:
//   - modelName: 模型名称，用于查找调用记录
func waitModelInterval(modelName string) {
	if modelMinInterval <= 0 {
		return
	}
	globalRateLimiter.Lock()
	lastCall, ok := modelLastCall[modelName]
	if !ok {
		modelLastCall[modelName] = time.Now()
		globalRateLimiter.Unlock()
		return
	}
	elapsed := time.Since(lastCall)
	globalRateLimiter.Unlock()

	if elapsed < modelMinInterval {
		jitter := time.Duration(rand.Int63n(int64(modelJitter)))
		sleep := modelMinInterval - elapsed + jitter
		time.Sleep(sleep)
	}

	globalRateLimiter.Lock()
	modelLastCall[modelName] = time.Now()
	globalRateLimiter.Unlock()
}

// generateTest 调用 LLM API 生成单元测试代码
// 处理重试、自动续写（截断时）、代码提取和验证
//
// 参数:
//   - ctx: 上下文，用于超时控制
//   - model: 模型配置信息
//   - language: 编程语言
//   - prompt: 提示词内容
//
// 返回值:
//   - string: 生成的测试代码
//   - map[string]any: 原始响应数据
//   - int: 请求耗时（毫秒）
//   - *int: prompt token 数量
//   - *int: completion token 数量
//   - *int: 总 token 数量
//   - bool: 是否被截断
//   - *ErrorInfo: 错误信息（成功时为 nil）
func (c *apiClient) generateTest(
	ctx context.Context,
	model modelConfig,
	language string,
	prompt string,
) (string, map[string]any, int, *int, *int, *int, bool, *contracts.ErrorInfo) {
	apiKey := strings.TrimSpace(os.Getenv(model.APIKeyEnv))
	if apiKey == "" {
		return "", nil, 0, nil, nil, nil, false, &contracts.ErrorInfo{
			Kind:      "auth_config_error",
			Message:   fmt.Sprintf("missing API key env var: %s (model=%s)", model.APIKeyEnv, model.Name),
			Retryable: false,
		}
	}

	payload := buildPayload(model, prompt)
	body, err := json.Marshal(payload)
	if err != nil {
		return "", nil, 0, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "payload_error", Message: err.Error(), Retryable: false}
	}

	endpoint := resolveEndpoint(model)
	var lastErr *contracts.ErrorInfo
	var lastTruncated bool
	var allResponses []map[string]any
	var accumulatedCode strings.Builder
	var totalPromptTokens, totalCompletionTokens, totalTokens int
	var tokenCountsSet bool
	maxContinuationAttempts := 3

	for attempt := 1; attempt <= c.retries; attempt++ {
		started := time.Now()
		code, rawResp, p, cm, total, truncated, errInfo := c.doOnce(ctx, endpoint, apiKey, model.Provider, body)
		latency := int(time.Since(started).Milliseconds())

		if errInfo != nil {
			lastErr = errInfo
			lastTruncated = truncated
			if !errInfo.Retryable || attempt >= c.retries {
				break
			}
			sleep := float64(c.backoff) * math.Pow(2, float64(attempt-1))
			sleep += float64(time.Duration(rand.Int63n(int64(200 * time.Millisecond))))
			time.Sleep(time.Duration(sleep))
			continue
		}

		allResponses = append(allResponses, rawResp)
		accumulatedCode.WriteString(code)

		if !tokenCountsSet && p != nil && cm != nil && total != nil {
			totalPromptTokens = *p
			totalCompletionTokens = *cm
			totalTokens = *total
			tokenCountsSet = true
		} else if tokenCountsSet && cm != nil {
			totalCompletionTokens += *cm
			totalTokens += *cm
		}

		if !truncated {
			san := sanitizeModelOutput(accumulatedCode.String(), language)
			extracted := extractCode(san, language)
			if vErr := validateGeneratedTest(extracted, language); vErr != nil {
				lastErr = &contracts.ErrorInfo{Kind: "quality_error", Message: vErr.Error(), Retryable: false}
				break
			}
			finalTruncated := false
			return extracted, mergeResponses(allResponses), latency, &totalPromptTokens, &totalCompletionTokens, &totalTokens, finalTruncated, nil
		}

		if maxContinuationAttempts <= 0 {
			lastTruncated = true
			break
		}
		maxContinuationAttempts--

		continuationPayload := buildContinuationPayload(model, prompt, accumulatedCode.String())
		body, err = json.Marshal(continuationPayload)
		if err != nil {
			return accumulatedCode.String(), mergeResponses(allResponses), latency, &totalPromptTokens, &totalCompletionTokens, &totalTokens, true, &contracts.ErrorInfo{
				Kind:      "continuation_payload_error",
				Message:   err.Error(),
				Retryable: false,
			}
		}
	}

	if lastErr == nil && accumulatedCode.Len() > 0 {
		san := sanitizeModelOutput(accumulatedCode.String(), language)
		extracted := extractCode(san, language)
		return extracted, mergeResponses(allResponses), 0, &totalPromptTokens, &totalCompletionTokens, &totalTokens, lastTruncated, nil
	}

	if lastErr == nil {
		lastErr = &contracts.ErrorInfo{Kind: "unknown_error", Message: "unknown generation error", Retryable: false}
	}
	return accumulatedCode.String(), mergeResponses(allResponses), 0, &totalPromptTokens, &totalCompletionTokens, &totalTokens, lastTruncated, lastErr
}

// doOnce 执行单次 API 调用
// 发送请求并解析响应，提取文本内容和 token 使用量
//
// 参数:
//   - ctx: 上下文
//   - endpoint: API 端点 URL
//   - apiKey: API 密钥
//   - provider: 提供商类型
//   - body: 请求体 JSON
//
// 返回值:
//   - string: 响应文本
//   - map[string]any: 原始响应
//   - *int, *int, *int: prompt/completion/total token 数量
//   - bool: 是否因长度截断
//   - *ErrorInfo: 错误信息
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
	if len(rawBytes) == 0 {
		return "", nil, nil, nil, nil, false, &contracts.ErrorInfo{
			Kind:      "response_parse_error",
			Message:   "empty response body",
			Retryable: true,
		}
	}

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
		return "", nil, nil, nil, nil, false, &contracts.ErrorInfo{
			Kind:      "response_parse_error",
			Message:   fmt.Sprintf("%s | body=%s", err.Error(), trimText(rawText, 500)),
			Retryable: true,
		}
	}

	text, err := extractResponseText(payload, provider)
	if err != nil {
		return "", payload, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "response_extract_error", Message: err.Error(), Retryable: false}
	}
	promptTokens, completionTokens, totalTokens := extractUsage(payload)
	truncated := extractFinishReason(payload, provider)
	return text, payload, promptTokens, completionTokens, totalTokens, truncated, nil
}

// resolveEndpoint 根据模型配置解析 API 端点 URL
// 处理不同提供商的端点格式差异
//
// 参数:
//   - model: 模型配置
//
// 返回值:
//   - string: 完整的 API 端点 URL
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

// buildPayload 构建 API 请求体
// 根据提供商类型生成不同格式的请求结构
//
// 参数:
//   - model: 模型配置
//   - prompt: 提示词内容
//
// 返回值:
//   - map[string]any: 请求体 JSON 结构
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
			{"role": "system", "content": systemMessage},
			{"role": "user", "content": prompt},
		},
		"stream": false,
	}
	for k, v := range params {
		payload[k] = v
	}
	return payload
}

// buildContinuationPayload 构建续写请求体
// 当响应被截断时，使用此函数构建继续生成的请求
//
// 参数:
//   - model: 模型配置
//   - originalPrompt: 原始提示词
//   - generatedSoFar: 已生成的内容
//
// 返回值:
//   - map[string]any: 续写请求体 JSON 结构
func buildContinuationPayload(model modelConfig, originalPrompt string, generatedSoFar string) map[string]any {
	continuationPrompt := "Continue generating the unit test code from where you left off. " +
		"Output only the remaining code without any explanations or markdown fences. " +
		"Do not repeat what was already generated."

	params := map[string]any{}
	for k, v := range model.Params {
		params[k] = v
	}

	switch model.Provider {
	case "dashscope":
		if !strings.Contains(model.Endpoint, "compatible-mode") {
			return map[string]any{
				"model": model.Model,
				"input": map[string]any{
					"messages": []map[string]any{
						{"role": "user", "content": originalPrompt},
						{"role": "assistant", "content": generatedSoFar},
						{"role": "user", "content": continuationPrompt},
					},
				},
				"parameters": params,
			}
		}
		return map[string]any{
			"model": model.Model,
			"messages": []map[string]any{
				{"role": "system", "content": systemMessage},
				{"role": "user", "content": originalPrompt},
				{"role": "assistant", "content": generatedSoFar},
				{"role": "user", "content": continuationPrompt},
			},
			"stream": false,
		}
	case "volcengine":
		return map[string]any{
			"model": model.Model,
			"messages": []map[string]any{
				{"role": "system", "content": systemMessage},
				{"role": "user", "content": originalPrompt},
				{"role": "assistant", "content": generatedSoFar},
				{"role": "user", "content": continuationPrompt},
			},
			"stream": false,
		}
	default:
		return map[string]any{
			"model": model.Model,
			"messages": []map[string]any{
				{"role": "system", "content": systemMessage},
				{"role": "user", "content": originalPrompt},
				{"role": "assistant", "content": generatedSoFar},
				{"role": "user", "content": continuationPrompt},
			},
			"stream": false,
		}
	}
}

// mergeResponses 合并多次响应（用于续写场景）
// 将多个 API 响应合并为一个，累加 token 使用量
//
// 参数:
//   - responses: 响应列表
//
// 返回值:
//   - map[string]any: 合并后的响应数据
func mergeResponses(responses []map[string]any) map[string]any {
	if len(responses) == 0 {
		return map[string]any{"merged": true, "count": 0}
	}
	if len(responses) == 1 {
		responses[0]["merged"] = true
		responses[0]["continuation_count"] = 1
		return responses[0]
	}

	mergedText := ""
	var totalPromptTokens, totalCompletionTokens, totalTokens int64
	continuationCount := len(responses)

	for _, resp := range responses {
		if text, err := extractResponseTextFromAny(resp); err == nil {
			mergedText += text
		}
		if usage, ok := resp["usage"].(map[string]any); ok {
			if p, ok := usage["prompt_tokens"].(float64); ok {
				totalPromptTokens += int64(p)
			}
			if c, ok := usage["completion_tokens"].(float64); ok {
				totalCompletionTokens += int64(c)
			}
			if t, ok := usage["total_tokens"].(float64); ok {
				totalTokens += int64(t)
			}
		}
	}

	return map[string]any{
		"merged":             true,
		"continuation_count": continuationCount,
		"merged_text":        mergedText,
		"prompt_tokens":      totalPromptTokens,
		"completion_tokens":  totalCompletionTokens,
		"total_tokens":       totalTokens,
		"first_response":     responses[0],
	}
}

// extractResponseTextFromAny 从响应中提取文本内容
// 支持多种响应格式（OpenAI、Dashscope 等）
//
// 参数:
//   - response: 响应数据
//
// 返回值:
//   - string: 提取的文本
//   - error: 错误信息
func extractResponseTextFromAny(response map[string]any) (string, error) {
	if choices, ok := response["choices"].([]any); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]any); ok {
			if msg, ok := choice["message"].(map[string]any); ok {
				if content, ok := msg["content"].(string); ok {
					return content, nil
				}
			}
		}
	}
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
	return "", fmt.Errorf("unable to extract text from response")
}

// extractResponseText 从响应中提取文本内容（带提供商参数）
// 根据提供商类型选择正确的提取路径
//
// 参数:
//   - response: 响应数据
//   - provider: 提供商类型
//
// 返回值:
//   - string: 提取的文本
//   - error: 错误信息
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

// extractUsage 从响应中提取 token 使用量
// 解析 usage 字段，返回 prompt/completion/total token 数量
//
// 参数:
//   - response: 响应数据
//
// 返回值:
//   - *int: prompt tokens
//   - *int: completion tokens
//   - *int: total tokens
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

// extractFinishReason 从响应中提取结束原因
// 检查 finish_reason 是否为 "length"（表示被截断）
//
// 参数:
//   - response: 响应数据
//   - provider: 提供商类型
//
// 返回值:
//   - bool: 是否因长度截断
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

// toIntPtr 将任意类型转换为 int 指针
// 支持 int、int64、float64 类型
//
// 参数:
//   - v: 输入值
//
// 返回值:
//   - *int: 转换后的 int 指针（失败时为 nil）
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

// sanitizeModelOutput 清理模型输出内容
// 移除推理标签（<think>、<analysis>）和非代码前缀
//
// 参数:
//   - content: 原始输出内容
//   - language: 编程语言
//
// 返回值:
//   - string: 清理后的内容
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

// trimNonCodePrefix 移除 Python 输出中的非代码前缀
// 从第一个代码行开始截取内容
//
// 参数:
//   - text: 输出文本
//
// 返回值:
//   - string: 去除前缀后的文本
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

// extractCode 从输出中提取代码块
// 解析 markdown 代码块，返回指定语言的代码
//
// 参数:
//   - content: 输出内容
//   - language: 编程语言
//
// 返回值:
//   - string: 提取的代码
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

// stripMarkdownFence 去除 markdown 代码围栏标记
// 处理以 ``` 开头但没有语言标记的情况
//
// 参数:
//   - content: 输出内容
//
// 返回值:
//   - string: 去除围栏后的内容
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

// dropTrailingFenceLines 去除末尾的 markdown 围栏行
// 从文本末尾移除 ``` 行
//
// 参数:
//   - text: 输出文本
//
// 返回值:
//   - string: 清理后的文本
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

// validateGeneratedTest 验证生成的测试代码
// 检查输出是否为空、是否包含泄露的推理标签、是否有测试结构
//
// 参数:
//   - code: 生成的代码
//   - language: 编程语言
//
// 返回值:
//   - error: 验证错误（成功时为 nil）
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

// trimText 截断文本到指定最大长度
//
// 参数:
//   - v: 输入文本
//   - max: 最大长度
//
// 返回值:
//   - string: 截断后的文本
func trimText(v string, max int) string {
	if len(v) <= max {
		return v
	}
	return v[:max]
}

// coverageTargetsText 返回覆盖率目标说明文本
// 用于提示词中告知模型覆盖率要求
//
// 返回值:
//   - string: 覆盖率目标说明
func coverageTargetsText() string {
	// 与 benchmark/config/models.yaml 默认阈值保持同步
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

// languageFramework 返回指定语言的测试框架名称
// 用于提示词中告知模型使用正确的测试框架
//
// 参数:
//   - language: 编程语言
//
// 返回值:
//   - string: 测试框架名称
func languageFramework(language string) string {
	switch language {
	case "java":
		return "JUnit 5 (Jupiter)"
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

// moduleImportName 将样本 ID 转换为合法的模块导入名
// 规范化特殊字符，处理数字开头的情况
//
// 参数:
//   - sampleID: 样本 ID
//
// 返回值:
//   - string: 合法的模块导入名
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

// parseSampleMeta 从样本 ID 或路径解析场景和复杂度信息
// 用于提示词构建时提供样本上下文
//
// 参数:
//   - sampleID: 样本 ID
//   - samplePath: 样本文件路径
//
// 返回值:
//   - string: 场景名称
//   - string: 复杂度级别
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

// extractDependencies 从源代码中提取依赖列表
// 根据语言类型使用正则匹配 import/include 语句
//
// 参数:
//   - sourceCode: 源代码内容
//   - language: 编程语言
//
// 返回值:
//   - []string: 依赖列表（最多12个）
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

// mockRequirement 根据源代码内容推断 mock 策略
// 检测外部依赖（网络、文件、数据库等）并给出 mock 建议
//
// 参数:
//   - sourceCode: 源代码内容
//
// 返回值:
//   - string: mock 策略说明
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

// extractCriticalConditions 从 Python 源代码中提取关键条件语句
// 用于提示词中强调需要测试的边界条件
//
// 参数:
//   - sourceCode: 源代码内容
//   - language: 编程语言（仅处理 Python）
//
// 返回值:
//   - []string: 关键条件语句列表（最多12个）
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

// containsComparator 检查行是否包含比较运算符
// 用于识别条件语句
//
// 参数:
//   - line: 代码行
//
// 返回值:
//   - bool: 是否包含比较运算符
func containsComparator(line string) bool {
	comparators := []string{"<=", ">=", "==", "!=", "<", ">"}
	for _, item := range comparators {
		if strings.Contains(line, item) {
			return true
		}
	}
	return false
}

// extractRepoLevelSymbols 从源代码中提取仓库级别的符号
// 用于提示词中告知模型可导入的符号
//
// 参数:
//   - sourceCode: 源代码内容
//   - language: 编程语言
//
// 返回值:
//   - string: 符号列表说明文本
func extractRepoLevelSymbols(sourceCode, language string) string {
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
	return "Repo-level symbols available to import: " + strings.Join(symbols, ", ") + "."
}

// repoLevelMetaForRunner 仓库级别样本的元数据结构
// 用于 Runner 加载和处理仓库级别样本
type repoLevelMetaForRunner struct {
	SampleID      string   `json:"sample_id"`              // 样本 ID
	ModuleImport  string   `json:"module_import"`          // 模块导入路径
	PackageName   string   `json:"package_name"`           // 包名
	TargetFile    string   `json:"target_file"`            // 目标文件路径
	WorkspaceRoot string   `json:"workspace_root"`         // 工作区根目录
	Requirements  []string `json:"requirements,omitempty"` // 依赖包列表
}

// loadRepoLevelMetaForRunner 从样本路径加载仓库级别元数据
// 查找并解析 meta.json 或 {name}.meta.json 文件
//
// 参数:
//   - samplePath: 样本文件路径
//
// 返回值:
//   - *repoLevelMetaForRunner: 元数据（失败时为 nil）
func loadRepoLevelMetaForRunner(samplePath string) *repoLevelMetaForRunner {
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
	var meta repoLevelMetaForRunner
	if err := json.Unmarshal(raw, &meta); err != nil {
		return nil
	}
	if meta.ModuleImport == "" {
		return nil
	}
	return &meta
}
