package evaluator

import "fmt"

// mutationStats 变异测试统计信息
// 这个类型需要在所有平台可用
type mutationStats struct {
	Total      int
	Killed     int
	Survived   int
	NoTests    int
	NotChecked int
	Duplicated int
	Timeout    int
	Skipped    int
	Suspicious int
}

type mutationResultStatus int

const (
	MutationStatusNotRun mutationResultStatus = iota
	MutationStatusSuccess
	MutationStatusFailed
	MutationStatusSkippedLowPassRate
	MutationStatusSkippedNoCoverage
	MutationStatusSkippedToolNotAvailable
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

type MutationCheckResult struct {
	ShouldRun   bool
	Status      mutationResultStatus
	PassRate    float64
	TotalTests  int
	PassedTests int
	Message     string
}

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
