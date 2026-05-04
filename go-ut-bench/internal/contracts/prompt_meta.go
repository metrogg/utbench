package contracts

// PromptMetaProvider 提示词元数据提供者接口
// 用于解耦 reporter 对 runner 的直接依赖
// 实现者负责加载或生成提示词策略、版本ID和模板预览
type PromptMetaProvider interface {
	// PromptStrategy 返回提示词策略名称（如 "structured-v1"）
	PromptStrategy() string

	// PromptVersionID 返回提示词版本ID（SHA1 哈希前12位）
	PromptVersionID() string

	// PromptTemplatePreview 返回指定语言的完整文件模式提示词模板预览
	PromptTemplatePreview(language string) string
}
