// progress 包提供终端进度报告功能
// 在终端显示实时进度和统计面板，与 Logger 解耦
package obs

import (
	"fmt"
	"sync"
	"time"
)

// TaskResult 单个任务的结果
type TaskResult struct {
	Model            string
	Language         string
	SampleID         string
	Success          bool
	Truncated        bool
	Error            string
	LatencyMS        int
	Tokens           int
	CompilePass      bool
	TestPass         bool
	TestPassCount    *int
	TestTotalCount   *int
	LineCoverage     float64
	MutationScore    float64
	MutationTool     string
	MutationTotal    int
	MutationKilled   int
	MutationSurvived int
}

// StageStats 阶段统计
type StageStats struct {
	Total    int
	Success  int
	Failed   int
	Skipped  int
	Duration time.Duration
}

// ProgressReporter 进度报告器
type ProgressReporter struct {
	total            int
	completed        int
	successCount     int
	failCount        int
	truncatedCount   int
	skipCount        int
	compilePassCount int
	testPassCount    int
	coverageSum      float64
	coverageCount    int
	mutationSum      float64
	mutationCount    int
	startTime        time.Time
	stageStartTime   time.Time
	stage            string
	mu               sync.Mutex
}

// NewProgressReporter 创建新的进度报告器
func NewProgressReporter(total int, stage string) *ProgressReporter {
	now := time.Now()
	return &ProgressReporter{
		total:          total,
		startTime:      now,
		stageStartTime: now,
		stage:          stage,
	}
}

// GetStartTime 获取开始时间
func (pr *ProgressReporter) GetStartTime() time.Time {
	pr.mu.Lock()
	defer pr.mu.Unlock()
	return pr.startTime
}

// OnTaskDone 处理任务完成
func (pr *ProgressReporter) OnTaskDone(result TaskResult) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	pr.completed++

	if result.Success {
		pr.successCount++
	} else if result.Error != "" {
		pr.failCount++
	}

	if result.Truncated {
		pr.truncatedCount++
	}

	if result.CompilePass {
		pr.compilePassCount++
	}
	if result.TestPass {
		pr.testPassCount++
	}

	if result.LineCoverage > 0 {
		pr.coverageSum += result.LineCoverage
		pr.coverageCount++
	}

	if result.MutationScore > 0 {
		pr.mutationSum += result.MutationScore
		pr.mutationCount++
	}
}

// PrintStats 打印统计面板
func (pr *ProgressReporter) PrintStats() {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	if pr.completed == 0 {
		return
	}

	elapsed := time.Since(pr.startTime)
	avgTime := elapsed / time.Duration(pr.completed)
	remaining := avgTime * time.Duration(pr.total-pr.completed)

	fmt.Println("========================================")
	fmt.Println("[STAT] 实时统计")
	fmt.Printf("   已处理: %d/%d (%d%%)\n", pr.completed, pr.total, pr.completed*100/pr.total)

	if pr.stage == "generate" {
		fmt.Printf("   成功: %d | 失败: %d | 截断: %d\n",
			pr.successCount, pr.failCount, pr.truncatedCount)
	} else if pr.stage == "evaluate" {
		fmt.Printf("   编译通过: %d/%d (%d%%)\n",
			pr.compilePassCount, pr.completed, pr.compilePassCount*100/pr.completed)
		fmt.Printf("   测试通过: %d/%d (%d%%)\n",
			pr.testPassCount, pr.completed,
			pr.testPassCount*100/pr.completed)

		if pr.coverageCount > 0 {
			fmt.Printf("   平均行覆盖: %.1f%%\n", pr.coverageSum/float64(pr.coverageCount)*100)
		}
		if pr.mutationCount > 0 {
			fmt.Printf("   平均变异分数: %.1f%%\n", pr.mutationSum/float64(pr.mutationCount)*100)
		}
	}

	fmt.Printf("   阶段已耗时: %s\n", formatDuration(elapsed))
	fmt.Printf("   预计剩余: ~%s\n", formatDuration(remaining))
	fmt.Println("========================================")
}

// PrintStageStart 打印阶段开始
func (pr *ProgressReporter) PrintStageStart(stageName string, details string) {
	pr.mu.Lock()
	pr.stage = stageName
	pr.stageStartTime = time.Now()
	pr.mu.Unlock()

	fmt.Println()
	fmt.Println("========================================")
	fmt.Printf("[STAGE] %s\n", stageName)
	fmt.Println("========================================")
	if details != "" {
		fmt.Println(details)
		fmt.Println()
	}
}

// PrintStageDone 打印阶段结束
func (pr *ProgressReporter) PrintStageDone(stageName string, stats StageStats) {
	elapsed := time.Since(pr.stageStartTime)

	fmt.Println()
	fmt.Println("========================================")
	fmt.Printf("[DONE] %s完成\n", stageName)
	fmt.Printf("   总计: %d | 成功: %d | 失败: %d\n", stats.Total, stats.Success, stats.Failed)
	fmt.Printf("   耗时: %s\n", formatDuration(elapsed))
	fmt.Println("========================================")
}

// PrintTaskLine 打印单行任务进度
func (pr *ProgressReporter) PrintTaskLine(idx, total int, model, lang, sampleID, status string, extras ...string) {
	extra := ""
	if len(extras) > 0 {
		extra = " | " + extras[0]
	}
	fmt.Printf("[%d/%d]  %s | %s | %s | %s%s\n", idx, total, model, lang, sampleID, status, extra)
}

// PrintMutationResult 打印变异测试结果
func (pr *ProgressReporter) PrintMutationResult(model, lang, sampleID, tool string, total, killed, survived int, score float64, elapsed time.Duration, skipReason string) {
	if skipReason != "" {
		fmt.Printf("        [SKIP] 跳过变异(%s) | %s\n", skipReason, elapsed)
		return
	}
	fmt.Printf("        [MUTATION] %s | %d/%d | 存活: %d | 杀死: %d | 分数: %.0f%% | %s\n",
		tool, total, total, survived, killed, score*100, elapsed)
}

// PrintFinalSummary 打印最终汇总
func (pr *ProgressReporter) PrintFinalSummary(runID string, stages map[string]StageStats, totalCompilePass, totalTestPass, totalSamples int, avgCoverage, avgMutation float64) {
	totalElapsed := time.Since(pr.startTime)

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("[COMPLETE] 全部完成!")
	fmt.Printf("   运行ID: %s\n", runID)
	fmt.Printf("   总耗时: %s\n", formatDuration(totalElapsed))
	fmt.Println()

	for stageName, stats := range stages {
		fmt.Printf("   %s: %d/%d (%d%%) | 耗时: %s\n",
			stageName, stats.Success, stats.Total,
			stats.Success*100/stats.Total,
			formatDuration(stats.Duration))
	}

	fmt.Println()
	fmt.Printf("   编译通过率: %d%% (%d/%d)\n", totalCompilePass*100/totalSamples, totalCompilePass, totalSamples)
	fmt.Printf("   测试通过率: %d%% (%d/%d)\n", totalTestPass*100/totalSamples, totalTestPass, totalSamples)
	if avgCoverage > 0 {
		fmt.Printf("   平均行覆盖: %.1f%%\n", avgCoverage*100)
	}
	if avgMutation > 0 {
		fmt.Printf("   平均变异分数: %.1f%%\n", avgMutation*100)
	}
	fmt.Println("========================================")
}

// formatDuration 格式化持续时间
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}
