package evaluator

import "time"

// 超时常量集中管理
// 消除各语言评测函数中分散的硬编码超时值
const (
	// DefaultTestTimeoutSeconds 默认测试执行超时时间
	DefaultTestTimeoutSeconds = 180

	// MutationTimeoutSeconds 变异测试默认超时时间（Go/Java/Python 共用）
	MutationTimeoutSeconds = 120

	// CppMutationTimeoutSeconds C++ 变异测试超时时间（Mull 需要更长时间）
	CppMutationTimeoutSeconds = 300

	// PythonCompileTimeoutSeconds Python 编译检查超时时间
	PythonCompileTimeoutSeconds = 30
)

// defaultTestTimeout 默认测试超时时间（兼容旧代码）
const defaultTestTimeoutSeconds = DefaultTestTimeoutSeconds

func normalizedTimeout(seconds int) time.Duration {
	if seconds <= 0 {
		seconds = DefaultTestTimeoutSeconds
	}
	return time.Duration(seconds) * time.Second
}
