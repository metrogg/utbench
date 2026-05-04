package runner

import "go-ut-bench/internal/contracts"

// DefaultPromptMetaProvider 默认的提示词元数据提供者
// 委托给 runner 包的现有函数
type DefaultPromptMetaProvider struct{}

func (p *DefaultPromptMetaProvider) PromptStrategy() string {
	return PromptStrategy()
}

func (p *DefaultPromptMetaProvider) PromptVersionID() string {
	return PromptVersionID()
}

func (p *DefaultPromptMetaProvider) PromptTemplatePreview(language string) string {
	return PromptTemplatePreview(language)
}

// 确保 DefaultPromptMetaProvider 实现了 PromptMetaProvider 接口
var _ contracts.PromptMetaProvider = (*DefaultPromptMetaProvider)(nil)
