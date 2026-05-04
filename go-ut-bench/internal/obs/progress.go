// progress 包提供终端进度报告功能
// 在终端显示实时进度和统计面板，与 Logger 解耦
package obs

import (
	"fmt"
	"io"
	"os"
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

	// Agent 专属字段
	SubjectID        string // 被测对象 ID，如 opencode__deepseek-v4-flash__unit_test_skill
	SubjectKind      string // model_api / cli_agent / http_agent
	AgentFramework   string // Agent 框架名，如 opencode、model_api
	SkillName        string // skill 名称，如 unit_test_skill、no_skill
	InteractionCount int    // Agent 交互轮次
	ToolCallCount    int    // 工具调用次数
	FilesRead        int    // 读取的文件数
	FilesWritten     int    // 写入的文件数
	CommandsExecuted int    // 执行的命令数
}

// StageStats 阶段统计
type StageStats struct {
	Total    int
	Success  int
	Failed   int
	Skipped  int
	Duration time.Duration
}

// frameworkStats 单个 framework 的运行统计
type frameworkStats struct {
	success    int
	fail       int
	truncated  int
	totalTok   int
	totalDurMS int
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
	frameworks       map[string]*frameworkStats // framework -> stats
	writer           io.Writer                  // 输出目标，默认 os.Stdout
	mu               sync.Mutex
}

// NewProgressReporter 创建新的进度报告器（输出到 stdout）
func NewProgressReporter(total int, stage string) *ProgressReporter {
	return NewProgressReporterWithWriter(total, stage, os.Stdout)
}

// NewProgressReporterWithWriter 创建新的进度报告器，输出到指定 writer
func NewProgressReporterWithWriter(total int, stage string, w io.Writer) *ProgressReporter {
	if w == nil {
		w = os.Stdout
	}
	now := time.Now()
	return &ProgressReporter{
		total:          total,
		startTime:      now,
		stageStartTime: now,
		stage:          stage,
		frameworks:     make(map[string]*frameworkStats),
		writer:         w,
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

	// 追踪 per-framework 统计
	if result.AgentFramework != "" {
		fw, ok := pr.frameworks[result.AgentFramework]
		if !ok {
			fw = &frameworkStats{}
			pr.frameworks[result.AgentFramework] = fw
		}
		if result.Success {
			fw.success++
		} else {
			fw.fail++
		}
		if result.Truncated {
			fw.truncated++
		}
		fw.totalTok += result.Tokens
		fw.totalDurMS += result.LatencyMS
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

	w := pr.writer
	fmt.Fprintln(w, "========================================")
	fmt.Fprintln(w, "[STAT] 实时统计")
	fmt.Fprintf(w, "   已处理: %d/%d (%d%%)\n", pr.completed, pr.total, pr.completed*100/pr.total)

	if pr.stage == "generate" {
		fmt.Fprintf(w, "   成功: %d | 失败: %d | 截断: %d\n",
			pr.successCount, pr.failCount, pr.truncatedCount)

		// per-framework 分组统计
		if len(pr.frameworks) > 1 {
			fmt.Fprintln(w, "   ---")
			for fw, st := range pr.frameworks {
				total := st.success + st.fail
				avgDur := 0
				if total > 0 {
					avgDur = st.totalDurMS / total / 1000 // 秒
				}
				avgTok := 0
				if total > 0 {
					avgTok = st.totalTok / total
				}
				fmt.Fprintf(w, "   [%s] %d/%d 成功 | 截断 %d | 平均 %ds | 平均 %d tok\n",
					fw, st.success, total, st.truncated, avgDur, avgTok)
			}
		}
	} else if pr.stage == "evaluate" {
		fmt.Fprintf(w, "   编译通过: %d/%d (%d%%)\n",
			pr.compilePassCount, pr.completed, pr.compilePassCount*100/pr.completed)
		fmt.Fprintf(w, "   测试通过: %d/%d (%d%%)\n",
			pr.testPassCount, pr.completed,
			pr.testPassCount*100/pr.completed)

		if pr.coverageCount > 0 {
			fmt.Fprintf(w, "   平均行覆盖: %.1f%%\n", pr.coverageSum/float64(pr.coverageCount)*100)
		}
		if pr.mutationCount > 0 {
			fmt.Fprintf(w, "   平均变异分数: %.1f%%\n", pr.mutationSum/float64(pr.mutationCount)*100)
		}
	}

	fmt.Fprintf(w, "   阶段已耗时: %s\n", formatDuration(elapsed))
	fmt.Fprintf(w, "   预计剩余: ~%s\n", formatDuration(remaining))
	fmt.Fprintln(w, "========================================")
}

// PrintStageStart 打印阶段开始
func (pr *ProgressReporter) PrintStageStart(stageName string, details string) {
	pr.mu.Lock()
	pr.stage = stageName
	pr.stageStartTime = time.Now()
	pr.mu.Unlock()

	w := pr.writer
	fmt.Fprintln(w)
	fmt.Fprintln(w, "========================================")
	fmt.Fprintf(w, "[STAGE] %s\n", stageName)
	fmt.Fprintln(w, "========================================")
	if details != "" {
		fmt.Fprintln(w, details)
		fmt.Fprintln(w)
	}
}

// PrintStageDone 打印阶段结束
func (pr *ProgressReporter) PrintStageDone(stageName string, stats StageStats) {
	elapsed := time.Since(pr.stageStartTime)

	w := pr.writer
	fmt.Fprintln(w)
	fmt.Fprintln(w, "========================================")
	fmt.Fprintf(w, "[DONE] %s完成\n", stageName)
	fmt.Fprintf(w, "   总计: %d | 成功: %d | 失败: %d\n", stats.Total, stats.Success, stats.Failed)
	fmt.Fprintf(w, "   耗时: %s\n", formatDuration(elapsed))
	fmt.Fprintln(w, "========================================")
}

// PrintTaskLine 打印单行任务进度（旧版兼容签名）
func (pr *ProgressReporter) PrintTaskLine(idx, total int, model, lang, sampleID, status string, extras ...string) {
	extra := ""
	if len(extras) > 0 {
		extra = " | " + extras[0]
	}
	fmt.Fprintf(pr.writer, "[%d/%d]  %s | %s | %s | %s%s\n", idx, total, model, lang, sampleID, status, extra)
}

// PrintAgentTaskLine 打印 Agent 感知的单行任务进度
// 格式: [idx/total] framework | model | skill | lang | sampleID | status | duration | tokens | agent_info
func (pr *ProgressReporter) PrintAgentTaskLine(idx, total int, result TaskResult, status string) {
	w := pr.writer
	duration := formatDuration(time.Duration(result.LatencyMS) * time.Millisecond)

	// 基础信息行
	fmt.Fprintf(w, "[%d/%d]  %s | %s | %s | %s | %s | %s | %s",
		idx, total,
		result.AgentFramework,
		result.Model,
		result.SkillName,
		result.Language,
		result.SampleID,
		status,
		duration,
	)

	// token 信息
	if result.Tokens > 0 {
		fmt.Fprintf(w, " | %d tok", result.Tokens)
	}

	// Agent 追踪摘要（仅 cli_agent 类型）
	if result.SubjectKind == "cli_agent" {
		agentParts := []string{}
		if result.InteractionCount > 0 {
			agentParts = append(agentParts, fmt.Sprintf("%d轮", result.InteractionCount))
		}
		if result.ToolCallCount > 0 {
			agentParts = append(agentParts, fmt.Sprintf("%d工具", result.ToolCallCount))
		}
		if result.FilesWritten > 0 {
			agentParts = append(agentParts, fmt.Sprintf("%d写入", result.FilesWritten))
		}
		if result.FilesRead > 0 {
			agentParts = append(agentParts, fmt.Sprintf("%d读取", result.FilesRead))
		}
		if result.CommandsExecuted > 0 {
			agentParts = append(agentParts, fmt.Sprintf("%d命令", result.CommandsExecuted))
		}
		if len(agentParts) > 0 {
			fmt.Fprintf(w, " | %s", joinParts(agentParts))
		}
	}

	fmt.Fprintln(w)
}

func joinParts(parts []string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += ","
		}
		result += p
	}
	return result
}

// PrintMutationResult 打印变异测试结果
func (pr *ProgressReporter) PrintMutationResult(model, lang, sampleID, tool string, total, killed, survived int, score float64, elapsed time.Duration, skipReason string) {
	w := pr.writer
	if skipReason != "" {
		fmt.Fprintf(w, "        [SKIP] 跳过变异(%s) | %s\n", skipReason, elapsed)
		return
	}
	fmt.Fprintf(w, "        [MUTATION] %s | %d/%d | 存活: %d | 杀死: %d | 分数: %.0f%% | %s\n",
		tool, total, total, survived, killed, score*100, elapsed)
}

// PrintFinalSummary 打印最终汇总
func (pr *ProgressReporter) PrintFinalSummary(runID string, stages map[string]StageStats, totalCompilePass, totalTestPass, totalSamples int, avgCoverage, avgMutation float64) {
	totalElapsed := time.Since(pr.startTime)
	w := pr.writer

	fmt.Fprintln(w)
	fmt.Fprintln(w, "========================================")
	fmt.Fprintln(w, "[COMPLETE] 全部完成!")
	fmt.Fprintf(w, "   运行ID: %s\n", runID)
	fmt.Fprintf(w, "   总耗时: %s\n", formatDuration(totalElapsed))
	fmt.Fprintln(w)

	for stageName, stats := range stages {
		fmt.Fprintf(w, "   %s: %d/%d (%d%%) | 耗时: %s\n",
			stageName, stats.Success, stats.Total,
			stats.Success*100/stats.Total,
			formatDuration(stats.Duration))
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "   编译通过率: %d%% (%d/%d)\n", totalCompilePass*100/totalSamples, totalCompilePass, totalSamples)
	fmt.Fprintf(w, "   测试通过率: %d%% (%d/%d)\n", totalTestPass*100/totalSamples, totalTestPass, totalSamples)
	if avgCoverage > 0 {
		fmt.Fprintf(w, "   平均行覆盖: %.1f%%\n", avgCoverage*100)
	}
	if avgMutation > 0 {
		fmt.Fprintf(w, "   平均变异分数: %.1f%%\n", avgMutation*100)
	}
	fmt.Fprintln(w, "========================================")
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
