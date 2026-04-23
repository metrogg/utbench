package contracts

// ErrorInfo 定义错误信息的结构
// 用于统一描述各种操作中的错误详情，便于日志记录和错误追踪
type ErrorInfo struct {
	Kind       string `json:"kind"`        // 错误类型/分类，如"network_error"、"validation_error"等
	Message    string `json:"message"`     // 错误消息，描述具体错误内容
	Retryable  bool   `json:"retryable"`  // 是否可重试，true表示该错误可能是临时性的，可以重试
	StatusCode *int   `json:"status_code,omitempty"` // HTTP状态码（如果有），用于网络相关错误
}
