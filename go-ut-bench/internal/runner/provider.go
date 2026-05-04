package runner

import "strings"

// LLMProvider LLM 提供商接口
// 封装不同 API 提供商的请求构建和响应解析差异
//
// 实现要求：
//   - ResolveEndpoint: 根据模型配置构建完整的 API 端点 URL
//   - BuildPayload: 构建初始请求体
//   - BuildContinuationPayload: 构建截断后续写请求体
//   - ExtractResponseText: 从响应 JSON 中提取生成的文本
//   - ExtractFinishReason: 检查响应是否因长度截断
//
// 当前实现：
//   - OpenAICompatibleProvider: 适用于 deepseek、minimax、volcengine、dashscope-compatible
//   - DashscopeNativeProvider: 适用于 dashscope 原生 API（非 compatible-mode）
type LLMProvider interface {
	// ResolveEndpoint 构建完整的 API 端点 URL
	ResolveEndpoint(model modelConfig) string

	// BuildPayload 构建初始 API 请求体
	BuildPayload(model modelConfig, prompt string) map[string]any

	// BuildContinuationPayload 构建截断后的续写请求体
	BuildContinuationPayload(model modelConfig, originalPrompt, generatedSoFar string) map[string]any

	// ExtractResponseText 从响应中提取生成的文本
	ExtractResponseText(response map[string]any) (string, error)

	// ExtractFinishReason 检查响应是否因长度截断
	ExtractFinishReason(response map[string]any) bool
}

// providerRegistry 提供商注册表
var providerRegistry = map[string]LLMProvider{}

// RegisterProvider 注册 LLM 提供商
func RegisterProvider(name string, p LLMProvider) {
	providerRegistry[strings.ToLower(name)] = p
}

// getProvider 获取指定提供商
// 返回 nil 表示未注册
func getProvider(name string) LLMProvider {
	return providerRegistry[strings.ToLower(name)]
}

// resolveProvider 根据模型配置解析提供商
// 优先使用注册表，未注册时回退到 OpenAI 兼容模式
// dashscope 特殊处理：根据 endpoint 是否含 "compatible-mode" 选择实现
func resolveProvider(model modelConfig) LLMProvider {
	if isDashscopeNative(model) {
		return &DashscopeNativeProvider{}
	}
	if p := getProvider(model.Provider); p != nil {
		return p
	}
	// 未注册的提供商默认使用 OpenAI 兼容模式
	return &OpenAICompatibleProvider{}
}
