// evaluator/mutation_types.go 提供变异测试相关类型定义
// 定义统计结构、状态枚举和检查结果
package evaluator

import "fmt"

// mutationStats 变异测试统计信息
// 这个类型需要在所有平台可用
type mutationStats struct {
	Total      int // 总变异体数量
	Killed     int // 被测试杀死的变异体
	Survived   int // 存活的变异体（测试未检测到）
	NoTests    int // 无测试覆盖的变异体
	NotChecked int // 未检查的变异体
	Duplicated int // 重复的变异体
	Timeout    int // 执行超时的变异体
	Skipped    int // 被跳过的变异体
	Suspicious int // 可疑的变异体
}

// mutationResultStatus 变异测试结果状态枚举
type mutationResultStatus int

// 变异测试结果状态常量
const (
	MutationStatusNotRun                  mutationResultStatus = iota // 未运行
	MutationStatusSuccess                                             // 成功完成
	MutationStatusFailed                                              // 执行失败
	MutationStatusSkippedLowPassRate                                  // 因测试通过率过低跳过
	MutationStatusSkippedNoCoverage                                   // 因无测试覆盖跳过
	MutationStatusSkippedToolNotAvailable                             // 因工具不可用跳过
)

func (s mutationResultStatus) String() string {
	switch s {
	case MutationStatusNotRun:
		return "not_run"
	case MutationStatusSuccess:
		return "success"
	case MutationStatusFailed:
		return "failed"
	case MutationStatusSkippedLowPassRate:
		return "skipped_low_pass_rate"
	case MutationStatusSkippedNoCoverage:
		return "skipped_no_coverage"
	case MutationStatusSkippedToolNotAvailable:
		return "skipped_tool_not_available"
	default:
		return "unknown"
	}
}

// MutationCheckResult 变异测试检查结果
// 包含是否应运行变异测试的判断信息
type MutationCheckResult struct {
	ShouldRun   bool                 // 是否应该运行变异测试
	Status      mutationResultStatus // 结果状态
	PassRate    float64              // 测试通过率
	TotalTests  int                  // 总测试数
	PassedTests int                  // 通过测试数
	Message     string               // 结果消息
}

// CheckTestPassRate 检查测试通过率是否满足变异测试运行条件
// 根据最小通过率阈值判断是否应运行变异测试
//
// 参数:
//   - passed: 通过测试数
//   - total: 总测试数
//   - toolName: 工具名称
//   - minPassRate: 最小通过率阈值
//
// 返回值:
//   - MutationCheckResult: 检查结果
func CheckTestPassRate(passed, total int, toolName string, minPassRate float64) MutationCheckResult {
	if total <= 0 {
		return MutationCheckResult{
			ShouldRun:   false,
			Status:      MutationStatusSkippedNoCoverage,
			PassRate:    0,
			TotalTests:  0,
			PassedTests: 0,
			Message:     fmt.Sprintf("%s: no tests found, skipping mutation", toolName),
		}
	}

	passRate := float64(passed) / float64(total)

	if passRate >= minPassRate {
		return MutationCheckResult{
			ShouldRun:   true,
			Status:      MutationStatusSuccess,
			PassRate:    passRate,
			TotalTests:  total,
			PassedTests: passed,
			Message:     fmt.Sprintf("%s: pass rate %.1f%% (>=%.0f%%), proceeding with mutation", toolName, passRate*100, minPassRate*100),
		}
	}

	skipStatus := MutationStatusSkippedLowPassRate
	skipMessage := fmt.Sprintf("%s: pass rate %.1f%% (<%.0f%%), skipping mutation", toolName, passRate*100, minPassRate*100)

	if passRate == 0 {
		skipMessage = fmt.Sprintf("%s: all tests failed, skipping mutation", toolName)
	}

	return MutationCheckResult{
		ShouldRun:   false,
		Status:      skipStatus,
		PassRate:    passRate,
		TotalTests:  total,
		PassedTests: passed,
		Message:     skipMessage,
	}
}

// GetMinPassRateForTool 获取指定变异测试工具的最小通过率阈值
// 不同工具有不同的阈值要求
//
// 参数:
//   - toolName: 工具名称
//
// 返回值:
//   - float64: 最小通过率（0-1）
func GetMinPassRateForTool(toolName string) float64 {
	switch toolName {
	case "mull", "cpp":
		return 1.0
	case "gremlins", "go-mutesting", "go":
		return 1.0
	case "pitest", "java":
		return 0.5
	case "mutmut", "python":
		return 0.8
	default:
		return 0.8
	}
}
