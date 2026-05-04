package runner

import (
	"fmt"
	"strings"
)

// OpenAICompatibleProvider OpenAI 兼容提供商
// 适用于 deepseek、minimax、volcengine、dashscope-compatible-mode 等
// 使用标准 OpenAI Chat Completions 格式
type OpenAICompatibleProvider struct{}

func (p *OpenAICompatibleProvider) ResolveEndpoint(model modelConfig) string {
	return strings.TrimSuffix(model.Endpoint, "/") + "/chat/completions"
}

func (p *OpenAICompatibleProvider) BuildPayload(model modelConfig, prompt string) map[string]any {
	params := map[string]any{}
	for k, v := range model.Params {
		params[k] = v
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

func (p *OpenAICompatibleProvider) BuildContinuationPayload(model modelConfig, originalPrompt, generatedSoFar string) map[string]any {
	continuationPrompt := "Continue generating the unit test code from where you left off. " +
		"Output only the remaining code without any explanations or markdown fences. " +
		"Do not repeat what was already generated."

	params := map[string]any{}
	for k, v := range model.Params {
		params[k] = v
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
}

func (p *OpenAICompatibleProvider) ExtractResponseText(response map[string]any) (string, error) {
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

func (p *OpenAICompatibleProvider) ExtractFinishReason(response map[string]any) bool {
	if choices, ok := response["choices"].([]any); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]any); ok {
			if fr, ok := choice["finish_reason"].(string); ok {
				return fr == "length"
			}
		}
	}
	return false
}

// DashscopeNativeProvider Dashscope 原生 API 提供商
// 使用 Dashscope 专有的 input/parameters 格式，不包含 system message
type DashscopeNativeProvider struct{}

func (p *DashscopeNativeProvider) ResolveEndpoint(model modelConfig) string {
	base := strings.TrimSuffix(model.Endpoint, "/")
	if strings.HasSuffix(base, "/api/v1") {
		return base + "/services/aigc/text-generation/generation"
	}
	return base + "/services/aigc/text-generation/generation"
}

func (p *DashscopeNativeProvider) BuildPayload(model modelConfig, prompt string) map[string]any {
	params := map[string]any{}
	for k, v := range model.Params {
		params[k] = v
	}
	return map[string]any{
		"model":      model.Model,
		"input":      map[string]any{"messages": []map[string]any{{"role": "user", "content": prompt}}},
		"parameters": params,
	}
}

func (p *DashscopeNativeProvider) BuildContinuationPayload(model modelConfig, originalPrompt, generatedSoFar string) map[string]any {
	continuationPrompt := "Continue generating the unit test code from where you left off. " +
		"Output only the remaining code without any explanations or markdown fences. " +
		"Do not repeat what was already generated."

	params := map[string]any{}
	for k, v := range model.Params {
		params[k] = v
	}
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

func (p *DashscopeNativeProvider) ExtractResponseText(response map[string]any) (string, error) {
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
	// 回退到 OpenAI 格式
	if choices, ok := response["choices"].([]any); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]any); ok {
			if msg, ok := choice["message"].(map[string]any); ok {
				if content, ok := msg["content"].(string); ok {
					return content, nil
				}
			}
		}
	}
	return "", fmt.Errorf("unable to extract text from dashscope response")
}

func (p *DashscopeNativeProvider) ExtractFinishReason(response map[string]any) bool {
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

// init 注册默认提供商
// dashscope 通过 isDashscopeNative 判断使用哪个实现
func init() {
	RegisterProvider("deepseek", &OpenAICompatibleProvider{})
	RegisterProvider("minimax", &OpenAICompatibleProvider{})
	RegisterProvider("volcengine", &OpenAICompatibleProvider{})
	// dashscope 在 resolveProvider 中特殊处理
}

// isDashscopeNative 判断 dashscope 模型是否使用原生 API（非 compatible-mode）
func isDashscopeNative(model modelConfig) bool {
	return model.Provider == "dashscope" && !strings.Contains(model.Endpoint, "compatible-mode")
}
