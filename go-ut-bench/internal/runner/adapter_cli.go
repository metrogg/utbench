package runner

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-ut-bench/internal/contracts"
)

// generateCLIAgent 在沙箱中执行 CLI Agent 并收集丰富化的 trace。
func generateCLIAgent(ctx context.Context, sandboxRunner SandboxRunner, req AgentGenerateRequest) AgentGenerateResult {
	subjectID := req.Subject.Spec.ID
	framework := req.Subject.Framework
	skill := req.Subject.Skill
	sample := req.Sample

	// 1. 准备独立 workspace
	workRoot := filepath.Join(req.OutputRoot, "runs", req.RunID, "agent_workspaces", subjectID, sample.Language, sample.ID)
	if err := os.RemoveAll(workRoot); err != nil {
		return agentError("workspace_error", err)
	}
	if err := os.MkdirAll(workRoot, 0o755); err != nil {
		return agentError("workspace_error", err)
	}

	// 2. 复制样本到 workspace
	sourceFile, err := prepareAgentWorkspace(workRoot, sample)
	if err != nil {
		return agentError("workspace_error", err)
	}

	// 3. 注入 skill 文件
	outputFile := filepath.Join(workRoot, "generated_test"+languageExt(sample.Language))
	var skillDir string
	if strings.EqualFold(defaultString(skill.InjectMode, "prompt_append"), "agent_native") {
		// agent_native: 写入 Agent 框架原生 skill 目录（CodeBuddy: .codebuddy/skills/, OpenCode: .opencode/skills/）
		nativeDir, nativeErr := injectAgentNativeSkill(workRoot, framework.Name, skill)
		if nativeErr != nil {
			return agentError("skill_injection_error", nativeErr)
		}
		skillDir = nativeDir
	} else {
		// prompt_append / workspace_mount: 复制到 .utbench/skills/
		utbenchDir, err := injectSkillWorkspace(workRoot, skill)
		if err != nil {
			return agentError("skill_injection_error", err)
		}
		skillDir = utbenchDir
	}

	// 4. 构建容器内路径提示
	// 注意：如果 sandbox_mode=docker 但当前环境没有 Docker daemon，会降级为 local。
	effectiveSandboxMode := framework.SandboxMode
	if strings.EqualFold(effectiveSandboxMode, "docker") && !isDockerAvailable() {
		effectiveSandboxMode = "local"
	}
	outputHint := outputFile
	sourceHint := sourceFile
	skillHint := skillDir
	if !strings.EqualFold(effectiveSandboxMode, "local") {
		outputHint = "/workspace/generated_test" + languageExt(sample.Language)
		sourceHint = "/workspace/" + filepath.ToSlash(mustRel(workRoot, sourceFile))
		if skillDir != "" {
			skillHint = "/workspace/" + filepath.ToSlash(mustRel(workRoot, skillDir))
		}
	}

	// 5. 写入 Agent 任务说明
	agentPrompt := buildAgentPrompt(req.Prompt, sample, sourceHint, outputHint, skillHint, req.Subject.Spec.Framework, skill.Name)
	promptFile := filepath.Join(workRoot, "utbench_agent_prompt.md")
	if err := os.WriteFile(promptFile, []byte(agentPrompt), 0o644); err != nil {
		return agentError("workspace_error", err)
	}

	// 6. 执行前快照
	before, _ := snapshotWorkspace(workRoot)

	// 7. 准备 trace 输出目录
	traceDir := filepath.Join(req.MetaRoot, "agent_traces", subjectID, sample.Language)
	tracePath := filepath.Join(traceDir, sample.ID+".trace.jsonl")
	diffPath := filepath.Join(traceDir, sample.ID+".diff.json")
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		return agentError("trace_error", err)
	}

	// 8. 渲染命令和环境变量
	templateData := commandTemplateData{
		Workspace:         workRoot,
		PromptFile:        promptFile,
		OutputFile:        outputFile,
		SkillDir:          skillDir,
		Model:             req.Model.Name,
		ModelID:           req.Model.Model,
		ModelProvider:     req.Model.Provider,
		ModelEndpoint:     req.Model.Endpoint,
		AnthropicEndpoint: req.Model.AnthropicEndpoint,
		ModelAPIKeyEnv:    req.Model.APIKeyEnv,
		Framework:         req.Subject.Spec.Framework,
		SubjectID:         subjectID,
		Skill:             req.Subject.Spec.Skill,
		Language:          sample.Language,
		SampleID:          sample.ID,
		SourceFile:        sourceFile,
		ContainerWorkdir:  containerPath("/workspace", workRoot, effectiveSandboxMode),
		ContainerPrompt:   containerPath("/workspace/utbench_agent_prompt.md", promptFile, effectiveSandboxMode),
		ContainerOutput:   containerPath("/workspace/generated_test"+languageExt(sample.Language), outputFile, effectiveSandboxMode),
		ContainerSkillDir: containerPath("/workspace/.utbench/skills/"+safePathName(skill.Name), skillDir, effectiveSandboxMode),
	}
	cmdText, err := renderTemplateText("agent-command", framework.Command, templateData)
	if err != nil {
		return agentError("command_template_error", err)
	}
	if strings.TrimSpace(cmdText) == "" {
		return agentError("command_template_error", fmt.Errorf("cli agent command is empty"))
	}
	envMap, envFromHost, err := buildAgentEnv(framework, req.Model, templateData)
	if err != nil {
		return agentError("agent_env_error", err)
	}
	// 注入 run ID 以便 sandbox 容器能被 label 标记和清理
	envMap["UTBENCH_RUN_ID"] = req.RunID

	// 9. 平台托管的样本依赖准备
	sandboxReq := buildSandboxRunRequest(req.OutputRoot, framework, sample.Language, workRoot, cmdText, envMap, envFromHost)
	// 将源文件加入只读挂载，防止 Agent 意外修改被测源码
	if sourceFile != "" {
		sandboxReq.ReadOnlyMounts = append(sandboxReq.ReadOnlyMounts, sourceFile)
	}
	environmentSetup, setupErr := runSandboxPreflight(ctx, sandboxRunner, sandboxReq, buildSampleEnvironmentSetupCommands(sample, workRoot))
	if setupErr != nil {
		trace := AgentTrace{
			SubjectID:          subjectID,
			Framework:          req.Subject.Spec.Framework,
			Model:              req.Subject.Spec.Model,
			Skill:              req.Subject.Spec.Skill,
			SampleID:           sample.ID,
			Language:           sample.Language,
			Command:            cmdText,
			StartedAt:          time.Now().UTC(),
			FinishedAt:         time.Now().UTC(),
			EnvironmentSetup:   environmentSetup,
			SandboxProvider:    sandboxReq.Provider,
			SandboxImage:       sandboxReq.DockerImage,
			SandboxFingerprint: sandboxFingerprintForRequest(sandboxReq),
			TracePath:          tracePath,
			WorkspaceDiffPath:  diffPath,
		}
		_ = writeAgentTrace(tracePath, trace)
		return AgentGenerateResult{
			RawResponse: map[string]any{
				"adapter":             "cli_agent",
				"subject_id":          subjectID,
				"framework":           req.Subject.Spec.Framework,
				"model":               req.Subject.Spec.Model,
				"skill":               req.Subject.Spec.Skill,
				"command":             cmdText,
				"environment_setup":   environmentSetup,
				"sandbox_provider":    sandboxReq.Provider,
				"sandbox_image":       sandboxReq.DockerImage,
				"sandbox_mode":        sandboxReq.Mode,
				"sandbox_workspace":   workRoot,
				"sandbox_fingerprint": trace.SandboxFingerprint,
			},
			Trace: trace,
			Error: &contracts.ErrorInfo{
				Kind:      "sample_env_prepare_error",
				Message:   setupErr.Error(),
				Retryable: false,
			},
		}
	}

	// 10. 执行沙箱预检
	preflightChecks, preflightErr := runSandboxPreflight(ctx, sandboxRunner, sandboxReq, frameworkPreflightCommands(framework, sample.Language))
	if preflightErr != nil {
		trace := AgentTrace{
			SubjectID:          subjectID,
			Framework:          req.Subject.Spec.Framework,
			Model:              req.Subject.Spec.Model,
			Skill:              req.Subject.Spec.Skill,
			SampleID:           sample.ID,
			Language:           sample.Language,
			Command:            cmdText,
			StartedAt:          time.Now().UTC(),
			FinishedAt:         time.Now().UTC(),
			EnvironmentSetup:   environmentSetup,
			PreflightChecks:    preflightChecks,
			SandboxProvider:    sandboxReq.Provider,
			SandboxImage:       sandboxReq.DockerImage,
			SandboxFingerprint: sandboxFingerprintForRequest(sandboxReq),
			TracePath:          tracePath,
			WorkspaceDiffPath:  diffPath,
		}
		_ = writeAgentTrace(tracePath, trace)
		return AgentGenerateResult{
			RawResponse: map[string]any{
				"adapter":             "cli_agent",
				"subject_id":          subjectID,
				"framework":           req.Subject.Spec.Framework,
				"model":               req.Subject.Spec.Model,
				"skill":               req.Subject.Spec.Skill,
				"command":             cmdText,
				"environment_setup":   environmentSetup,
				"preflight_checks":    preflightChecks,
				"sandbox_provider":    sandboxReq.Provider,
				"sandbox_image":       sandboxReq.DockerImage,
				"sandbox_mode":        sandboxReq.Mode,
				"sandbox_workspace":   workRoot,
				"sandbox_fingerprint": trace.SandboxFingerprint,
			},
			Trace: trace,
			Error: &contracts.ErrorInfo{
				Kind:      "sandbox_preflight_error",
				Message:   preflightErr.Error(),
				Retryable: false,
			},
		}
	}

	// 11. 执行主命令
	started := time.Now()
	runOutput, runErr := sandboxRunner.Run(ctx, sandboxReq)
	finished := time.Now()
	latency := int(finished.Sub(started).Milliseconds())

	// 12. 执行后快照 + diff
	after, _ := snapshotWorkspace(workRoot)
	changes := diffSnapshots(before, after)

	// 13. 解析 Agent 输出，提取丰富化 trace 信息
	trace := AgentTrace{
		SubjectID:          subjectID,
		Framework:          req.Subject.Spec.Framework,
		Model:              req.Subject.Spec.Model,
		Skill:              req.Subject.Spec.Skill,
		SampleID:           sample.ID,
		Language:           sample.Language,
		Command:            cmdText,
		ExitCode:           runOutput.ExitCode,
		DurationMS:         latency,
		StartedAt:          started.UTC(),
		FinishedAt:         finished.UTC(),
		Stdout:             trimText(runOutput.Stdout, 8000),
		Stderr:             trimText(runOutput.Stderr, 8000),
		WorkspaceDiff:      changes,
		EnvironmentSetup:   environmentSetup,
		PreflightChecks:    preflightChecks,
		SandboxProvider:    sandboxReq.Provider,
		SandboxImage:       sandboxReq.DockerImage,
		SandboxFingerprint: sandboxFingerprintForRequest(sandboxReq),
		TracePath:          tracePath,
		WorkspaceDiffPath:  diffPath,
	}

	// 从 Agent 输出中解析结构化信息
	parseAgentOutput(&trace, runOutput.Stdout, runOutput.Stderr)
	collectOpenCodeSessionExport(ctx, sandboxRunner, sandboxReq, workRoot, traceDir, sample.ID, &trace)
	finalizeAgentAccounting(&trace, req.Prompt, "", req.Model)

	// 14. 写入 trace 文件（完整结构化数据）
	_ = writeAgentTrace(tracePath, trace)
	_ = contracts.WriteJSON(diffPath, map[string]any{
		"workspace":  workRoot,
		"subject_id": subjectID,
		"changes":    changes,
	})

	// 15. 构建原始响应
	rawResponse := map[string]any{
		"adapter":             "cli_agent",
		"subject_id":          subjectID,
		"framework":           req.Subject.Spec.Framework,
		"model":               req.Subject.Spec.Model,
		"skill":               req.Subject.Spec.Skill,
		"command":             cmdText,
		"exit_code":           runOutput.ExitCode,
		"latency_ms":          latency,
		"trace_path":          tracePath,
		"workspace_diff_path": diffPath,
		"stdout":              trimText(runOutput.Stdout, 4000),
		"stderr":              trimText(runOutput.Stderr, 4000),
		"interaction_count":   trace.InteractionCount,
		"tool_call_count":     len(trace.ToolCalls),
		"files_read":          trace.FilesRead,
		"files_written":       trace.FilesWritten,
		"commands_executed":   trace.CommandsExecuted,
		"environment_setup":   environmentSetup,
		"preflight_checks":    preflightChecks,
		"sandbox_provider":    sandboxReq.Provider,
		"sandbox_image":       sandboxReq.DockerImage,
		"sandbox_fingerprint": trace.SandboxFingerprint,
		"token_source":        trace.TokenSource,
		"estimated_cost_usd":  trace.EstimatedCost,
		"cost_source":         trace.CostSource,
		"usage_source_detail": trace.UsageSourceDetail,
		"session_id":          trace.SessionID,
		"session_export_path": trace.SessionExportPath,
	}
	if trace.SessionExportError != "" {
		rawResponse["session_export_error"] = trace.SessionExportError
	}

	// 16. 处理执行错误
	if runErr != nil {
		return AgentGenerateResult{
			RawResponse: rawResponse,
			Trace:       trace,
			LatencyMS:   latency,
			Error: &contracts.ErrorInfo{
				Kind:      "agent_execution_error",
				Message:   fmt.Sprintf("agent command failed: %s", trimText(runOutput.Stderr+"\n"+runErr.Error(), 1000)),
				Retryable: false,
			},
		}
	}

	// 17. 拦截环境漂移行为
	if violation := detectSandboxPolicyViolation(trace.CommandsExecuted, frameworkForbiddenCommandPatterns(framework)); violation != "" {
		rawResponse["policy_violation"] = violation
		return AgentGenerateResult{
			RawResponse: rawResponse,
			Trace:       trace,
			LatencyMS:   latency,
			Error: &contracts.ErrorInfo{
				Kind:      "sandbox_policy_error",
				Message:   violation,
				Retryable: false,
			},
		}
	}

	// 18. 查找生成的测试文件
	generatedPath := findGeneratedTest(workRoot, outputFile, framework.OutputGlobs, changes, sample.Language)
	if generatedPath == "" {
		return AgentGenerateResult{
			RawResponse: rawResponse,
			Trace:       trace,
			LatencyMS:   latency,
			Error: &contracts.ErrorInfo{
				Kind:      "agent_output_error",
				Message:   "agent did not produce a test file",
				Retryable: false,
			},
		}
	}

	raw, err := os.ReadFile(generatedPath)
	if err != nil {
		return AgentGenerateResult{
			RawResponse: rawResponse,
			Trace:       trace,
			LatencyMS:   latency,
			Error:       &contracts.ErrorInfo{Kind: "agent_output_error", Message: err.Error(), Retryable: false},
		}
	}
	code := strings.TrimSpace(string(raw))
	finalizeAgentAccounting(&trace, req.Prompt, code, req.Model)
	rawResponse["token_source"] = trace.TokenSource
	rawResponse["estimated_cost_usd"] = trace.EstimatedCost
	rawResponse["cost_source"] = trace.CostSource
	rawResponse["usage_source_detail"] = trace.UsageSourceDetail
	rawResponse["session_id"] = trace.SessionID
	rawResponse["session_export_path"] = trace.SessionExportPath
	if trace.SessionExportError != "" {
		rawResponse["session_export_error"] = trace.SessionExportError
	}
	_ = writeAgentTrace(tracePath, trace)
	if err := validateGeneratedTest(code, sample.Language); err != nil {
		return AgentGenerateResult{
			Code:             code,
			RawResponse:      rawResponse,
			Trace:            trace,
			LatencyMS:        latency,
			PromptTokens:     trace.PromptTokens,
			CompletionTokens: trace.CompletionTokens,
			TotalTokens:      trace.TotalTokens,
			TokenSource:      trace.TokenSource,
			EstimatedCostUSD: trace.EstimatedCost,
			CostSource:       trace.CostSource,
			Error:            &contracts.ErrorInfo{Kind: "quality_error", Message: err.Error(), Retryable: false},
		}
	}

	rawResponse["generated_test_source_path"] = generatedPath
	rawResponse["generated_test_path"] = req.TestPath

	return AgentGenerateResult{
		Code:             code,
		RawResponse:      rawResponse,
		Trace:            trace,
		LatencyMS:        latency,
		PromptTokens:     trace.PromptTokens,
		CompletionTokens: trace.CompletionTokens,
		TotalTokens:      trace.TotalTokens,
		TokenSource:      trace.TokenSource,
		EstimatedCostUSD: trace.EstimatedCost,
		CostSource:       trace.CostSource,
	}
}

// parseAgentOutput 从 Agent 的 stdout/stderr 中解析结构化信息。
// 支持解析 OpenCode、Claude Code、CodeBuddy 等 CLI Agent 的日志格式。
func parseAgentOutput(trace *AgentTrace, stdout, stderr string) {
	fullOutput := stdout + "\n" + stderr

	trace.CommandsExecuted = parseCommands(fullOutput)

	// 解析工具调用（根据框架选择不同解析策略）
	switch {
	case strings.EqualFold(trace.Framework, "codebuddy"):
		trace.ToolCalls = parseCodeBuddyToolCalls(stdout, stderr)
	case strings.EqualFold(trace.Framework, "claudecode"), strings.EqualFold(trace.Framework, "claude_code"), strings.EqualFold(trace.Framework, "claude-code"):
		trace.ToolCalls = parseClaudeCodeToolCalls(stdout, stderr)
	default:
		trace.ToolCalls = parseToolCalls(fullOutput)
	}

	// 解析文件读写
	trace.FilesRead = parseFileReads(trace.CommandsExecuted, trace.Command)
	trace.FilesWritten = parseFileWrites(trace.CommandsExecuted, trace.WorkspaceDiff)

	// 尝试从结构化输出/日志中解析 token 用量与 session 信息
	parseUsageAndSession(trace, stdout, stderr)

	// 计算交互轮次（基于工具调用数量）
	trace.InteractionCount = len(trace.ToolCalls)
	if trace.InteractionCount == 0 {
		// 如果没有解析到工具调用，至少算 1 轮（Agent 生成了一次）
		trace.InteractionCount = 1
	}
}

type usageRecord struct {
	Prompt     *int
	Completion *int
	Total      *int
}

// parseToolCalls 从输出中解析工具调用记录。
func parseToolCalls(output string) []ToolCall {
	var calls []ToolCall
	seen := map[string]struct{}{}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		match := toolStartPattern.FindStringSubmatch(line)
		if len(match) != 2 {
			continue
		}
		tool := strings.TrimSpace(match[1])
		if tool == "" {
			continue
		}
		key := tool + "|" + line
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		calls = append(calls, ToolCall{
			Tool:    tool,
			Input:   trimText(line, 500),
			Success: true,
		})
	}
	return calls
}

// parseFileReads 从输出中解析 Agent 读取的文件列表。
func parseFileReads(commands []string, topLevelCommand string) []string {
	seen := map[string]struct{}{}
	var files []string
	for _, line := range append([]string{topLevelCommand}, commands...) {
		for _, candidate := range extractRelevantPaths(line) {
			if _, ok := seen[candidate]; ok {
				continue
			}
			seen[candidate] = struct{}{}
			files = append(files, candidate)
		}
	}
	sort.Strings(files)
	return files
}

// parseFileWrites 从输出中解析 Agent 写入的文件列表。
func parseFileWrites(commands []string, workspaceDiff []string) []string {
	seen := map[string]struct{}{}
	var files []string

	// 从 workspace diff 中获取（最可靠）
	for _, f := range workspaceDiff {
		if !isRelevantWorkspacePath(f) {
			continue
		}
		if _, ok := seen[f]; !ok {
			seen[f] = struct{}{}
			files = append(files, f)
		}
	}

	// 从命令文本中补充工作区内的高置信度路径。
	for _, line := range commands {
		for _, candidate := range extractRelevantPaths(line) {
			if !isRelevantWorkspaceWritePath(candidate) {
				continue
			}
			if _, ok := seen[candidate]; !ok {
				seen[candidate] = struct{}{}
				files = append(files, candidate)
			}
		}
	}
	sort.Strings(files)
	return files
}

// parseCommands 从输出中解析 Agent 执行的 shell 命令。
func parseCommands(output string) []string {
	var commands []string
	seen := map[string]struct{}{}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if match := bashPattern.FindStringSubmatch(trimmed); len(match) == 2 {
			cmd := strings.TrimSpace(match[1])
			if !looksLikeRealCommand(cmd) {
				continue
			}
			if _, ok := seen[cmd]; ok {
				continue
			}
			seen[cmd] = struct{}{}
			commands = append(commands, trimText(cmd, 500))
		}
	}
	return commands
}

func runSandboxPreflight(ctx context.Context, sandboxRunner SandboxRunner, baseReq SandboxRunRequest, commands []string) ([]PreflightCheck, error) {
	checks := make([]PreflightCheck, 0, len(commands))
	for _, command := range commands {
		command = strings.TrimSpace(command)
		if command == "" {
			continue
		}
		req := baseReq
		req.Command = command
		started := time.Now()
		out, err := sandboxRunner.Run(ctx, req)
		check := PreflightCheck{
			Command:    command,
			ExitCode:   out.ExitCode,
			DurationMS: int(time.Since(started).Milliseconds()),
			Stdout:     trimText(out.Stdout, 1000),
			Stderr:     trimText(out.Stderr, 1000),
			Passed:     err == nil && out.ExitCode == 0,
		}
		checks = append(checks, check)
		if !check.Passed {
			message := trimText(strings.TrimSpace(out.Stderr+"\n"+out.Stdout), 500)
			if message == "" && err != nil {
				message = err.Error()
			}
			return checks, fmt.Errorf("sandbox preflight failed for %q: %s", command, message)
		}
	}
	return checks, nil
}

func detectSandboxPolicyViolation(commands []string, patterns []string) string {
	if len(patterns) == 0 {
		return ""
	}
	for _, pattern := range patterns {
		pattern = strings.ToLower(strings.TrimSpace(pattern))
		if pattern == "" {
			continue
		}
		for _, command := range commands {
			if strings.Contains(strings.ToLower(command), pattern) {
				return fmt.Sprintf("sandbox policy violation: attempted forbidden environment mutation command %q", command)
			}
		}
	}
	return ""
}

// looksLikeFilePath 判断一个字符串是否像文件路径。
func looksLikeFilePath(s string) bool {
	if len(s) < 3 || len(s) > 500 {
		return false
	}
	if strings.ContainsAny(s, "\r\n\t ") {
		return false
	}
	// 包含路径分隔符
	if strings.Contains(s, "/") || strings.Contains(s, "\\") {
		// 包含文件扩展名
		ext := filepath.Ext(s)
		return ext != "" && len(ext) <= 6
	}
	return false
}

// writeAgentTrace 将完整 trace 写入 JSONL 文件。
func writeAgentTrace(path string, trace AgentTrace) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.Marshal(trace)
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(path, raw, 0o644)
}

// agentError 构建一个包含错误的 AgentGenerateResult。
func agentError(kind string, err error) AgentGenerateResult {
	return AgentGenerateResult{
		Error: &contracts.ErrorInfo{
			Kind:      kind,
			Message:   err.Error(),
			Retryable: false,
		},
	}
}

var (
	sessionIDPattern = regexp.MustCompile(`session id=([A-Za-z0-9_\-]+)`)
	pathPattern      = regexp.MustCompile("(?:path=|to\\s+|from\\s+|at\\s+|`)([/\\\\][^\\s\"'`,:;()]+(?:\\.[A-Za-z0-9]{1,8})?)")
	toolStartPattern = regexp.MustCompile(`service=tool\.registry\s+status=started\s+([A-Za-z0-9._-]+)`)
	bashPattern      = regexp.MustCompile(`permission=bash pattern=(.+?)(?:\s+ruleset=|\s+action=|\s+evaluated|$)`)
	extTokenPattern  = regexp.MustCompile("(?:^|[\\s\"'`=])((?:/workspace/)?[A-Za-z0-9_./-]+\\.(?:py|go|java|cpp|cc|cxx|h|hpp|md|txt|json|ya?ml))(?:$|[\\s\"'`,:;()])")
)

// parseCodeBuddyJSONOutput 尝试从 CodeBuddy Code 的 --output-format json 输出中
// 解析 session_id 和 token usage。
// CodeBuddy 在 -p 模式下输出 JSON（可能是对象或数组），包含消息和 usage 等字段。
func parseCodeBuddyJSONOutput(trace *AgentTrace, stdout string) {
	parseStructuredAgentJSONOutput(trace, stdout, "codebuddy_json_output")
}

// accumulateUsage 将 usageRecord 累加到 trace 的 token 字段。
func accumulateUsage(trace *AgentTrace, records []usageRecord) {
	var promptSum, completionSum, totalSum int
	var promptSeen, completionSeen, totalSeen bool
	for _, record := range records {
		if record.Prompt != nil {
			promptSum += *record.Prompt
			promptSeen = true
		}
		if record.Completion != nil {
			completionSum += *record.Completion
			completionSeen = true
		}
		if record.Total != nil {
			totalSum += *record.Total
			totalSeen = true
		}
	}
	if promptSeen && trace.PromptTokens == nil {
		trace.PromptTokens = intPtr(promptSum)
	}
	if completionSeen && trace.CompletionTokens == nil {
		trace.CompletionTokens = intPtr(completionSum)
	}
	if totalSeen && trace.TotalTokens == nil {
		trace.TotalTokens = intPtr(totalSum)
	}
	if trace.TotalTokens == nil && trace.PromptTokens != nil && trace.CompletionTokens != nil {
		trace.TotalTokens = intPtr(*trace.PromptTokens + *trace.CompletionTokens)
	}
}

// parseCodeBuddyJSONLUsage 从 CodeBuddy 的 stream-json (JSONL) 输出行中累加 token usage。
func parseCodeBuddyJSONLUsage(trace *AgentTrace, output string) {
	parseStructuredAgentJSONLUsage(trace, output, "codebuddy_jsonl_lines")
}

// parseCodeBuddyToolCalls 从 CodeBuddy 的 stream-json (JSONL) 输出中解析工具调用。
// 每行是一个 JSON 对象，助手消息的 content 数组中包含 type:"tool_use" 的块。
func parseCodeBuddyToolCalls(stdout, stderr string) []ToolCall {
	return parseStructuredJSONToolCalls(stdout, stderr)
}

func parseClaudeCodeJSONOutput(trace *AgentTrace, stdout string) {
	parseStructuredAgentJSONOutput(trace, stdout, "claudecode_json_output")
}

func parseClaudeCodeJSONLUsage(trace *AgentTrace, output string) {
	parseStructuredAgentJSONLUsage(trace, output, "claudecode_jsonl_lines")
}

func parseClaudeCodeToolCalls(stdout, stderr string) []ToolCall {
	return parseStructuredJSONToolCalls(stdout, stderr)
}

func parseStructuredAgentJSONOutput(trace *AgentTrace, stdout, usageDetail string) {
	jsonStr := extractFinalJSON(stdout)
	if jsonStr == "" {
		return
	}

	var payload any
	if err := json.Unmarshal([]byte(jsonStr), &payload); err != nil {
		return
	}

	assignSessionID(trace, payload)
	records := extractUsageRecords(payload)
	if len(records) == 0 {
		return
	}
	accumulateUsage(trace, records)
	trace.TokenSource = "actual"
	trace.UsageSourceDetail = usageDetail
}

func parseStructuredAgentJSONLUsage(trace *AgentTrace, output, usageDetail string) {
	var promptSum, completionSum, totalSum int
	var promptSeen, completionSeen, totalSeen bool

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var payload map[string]any
		if err := json.Unmarshal([]byte(line), &payload); err != nil {
			continue
		}
		assignSessionID(trace, payload)
		if typ, _ := payload["type"].(string); typ == "result" {
			if usage, ok := payload["usage"].(map[string]any); ok {
				for _, r := range extractUsageRecords(usage) {
					if r.Prompt != nil {
						promptSum += *r.Prompt
						promptSeen = true
					}
					if r.Completion != nil {
						completionSum += *r.Completion
						completionSeen = true
					}
					if r.Total != nil {
						totalSum += *r.Total
						totalSeen = true
					}
				}
				continue
			}
		}
		for _, r := range extractUsageRecords(payload) {
			if r.Prompt != nil {
				promptSum += *r.Prompt
				promptSeen = true
			}
			if r.Completion != nil {
				completionSum += *r.Completion
				completionSeen = true
			}
			if r.Total != nil {
				totalSum += *r.Total
				totalSeen = true
			}
		}
	}

	if promptSeen && trace.PromptTokens == nil {
		trace.PromptTokens = intPtr(promptSum)
	}
	if completionSeen && trace.CompletionTokens == nil {
		trace.CompletionTokens = intPtr(completionSum)
	}
	if totalSeen && trace.TotalTokens == nil {
		trace.TotalTokens = intPtr(totalSum)
	}
	if trace.TotalTokens == nil && trace.PromptTokens != nil && trace.CompletionTokens != nil {
		trace.TotalTokens = intPtr(*trace.PromptTokens + *trace.CompletionTokens)
	}
	if promptSeen || completionSeen || totalSeen {
		trace.TokenSource = "actual"
		trace.UsageSourceDetail = usageDetail
	}
}

func parseStructuredJSONToolCalls(stdout, stderr string) []ToolCall {
	var calls []ToolCall
	seen := map[string]struct{}{}

	for _, output := range []string{stdout, stderr} {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "{") {
				continue
			}
			var event map[string]any
			if err := json.Unmarshal([]byte(line), &event); err != nil {
				continue
			}

			// stream-json 格式：助手消息包含 content 数组
			// {"type":"assistant","message":{"role":"assistant","content":[{"type":"tool_use","name":"bash","input":{...}}]}}
			msg, _ := event["message"].(map[string]any)
			if msg == nil {
				msg = event // 兼容非嵌套格式
			}
			content, _ := msg["content"].([]any)
			for _, block := range content {
				blockMap, ok := block.(map[string]any)
				if !ok {
					continue
				}
				blockType, _ := blockMap["type"].(string)
				if blockType != "tool_use" && blockType != "function_call" {
					continue
				}
				tool, _ := blockMap["name"].(string)
				if tool == "" {
					continue
				}
				key := tool + "|" + line
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}

				input := ""
				if args, ok := blockMap["input"]; ok {
					if b, err := json.Marshal(args); err == nil {
						input = trimText(string(b), 500)
					}
				}
				calls = append(calls, ToolCall{
					Tool:    tool,
					Input:   input,
					Success: true,
				})
			}

			// 兼容：直接在 event 层级的工具字段
			if tool := extractToolName(event); tool != "" {
				key := tool + "|" + line
				if _, ok := seen[key]; !ok {
					seen[key] = struct{}{}
					calls = append(calls, ToolCall{Tool: tool, Success: true})
				}
			}
		}
	}
	return calls
}

func assignSessionID(trace *AgentTrace, value any) {
	if trace == nil || trace.SessionID != "" {
		return
	}
	switch v := value.(type) {
	case map[string]any:
		for _, key := range []string{"session_id", "sessionId"} {
			if sid, ok := v[key]; ok {
				if s, ok := sid.(string); ok && s != "" {
					trace.SessionID = s
					return
				}
			}
		}
		for _, child := range v {
			assignSessionID(trace, child)
			if trace.SessionID != "" {
				return
			}
		}
	case []any:
		for _, child := range v {
			assignSessionID(trace, child)
			if trace.SessionID != "" {
				return
			}
		}
	}
}

// extractToolName 从 CodeBuddy JSON 事件中提取工具名称。
// 支持多种常见字段名和嵌套结构。
func extractToolName(event map[string]any) string {
	// 直接字段
	for _, key := range []string{"tool", "tool_name", "name"} {
		if v, ok := event[key]; ok {
			if s, ok := v.(string); ok && s != "" {
				// 过滤掉非工具类型的事件
				lower := strings.ToLower(s)
				if lower == "result" || lower == "message" || lower == "error" || lower == "status" {
					continue
				}
				return s
			}
		}
	}
	// type 字段中提取（如 "tool_use", "tool_call"）
	if typ, ok := event["type"]; ok {
		if s, ok := typ.(string); ok {
			lower := strings.ToLower(s)
			if strings.Contains(lower, "tool") {
				// 尝试从 function.name 或 name 获取
				if fn, ok := event["function"]; ok {
					if fnMap, ok := fn.(map[string]any); ok {
						if n, ok := fnMap["name"]; ok {
							if name, ok := n.(string); ok {
								return name
							}
						}
					}
				}
			}
		}
	}
	return ""
}

// extractFinalJSON 从可能包含混合文本+JSON的输出中提取最后一个完整 JSON 值（对象或数组）。
func extractFinalJSON(output string) string {
	output = strings.TrimSpace(output)
	if output == "" {
		return ""
	}
	// 快速路径：整个输出就是合法 JSON
	if json.Valid([]byte(output)) {
		return output
	}
	// 从末尾向前搜索，尝试找到最后一个可解析的 JSON 对象或数组
	lastObj := strings.LastIndex(output, "{")
	lastArr := strings.LastIndex(output, "[")
	start := lastObj
	if lastArr > start {
		start = lastArr
	}
	if start < 0 {
		return ""
	}
	for i := start; i >= 0; i-- {
		c := output[i]
		if c != '{' && c != '[' {
			continue
		}
		candidate := output[i:]
		if json.Valid([]byte(candidate)) {
			return candidate
		}
		var payload any
		if json.Unmarshal([]byte(candidate), &payload) == nil {
			return candidate
		}
	}
	return ""
}

func parseUsageAndSession(trace *AgentTrace, stdout, stderr string) {
	fullOutput := stdout + "\n" + stderr

	// CodeBuddy / Claude Code: 优先从 JSON 输出中提取精确数据。
	// 同时检查 stdout 和 stderr，因为框架可能将结果输出到任一流。
	if strings.EqualFold(trace.Framework, "codebuddy") {
		parseCodeBuddyJSONOutput(trace, stdout)
		if trace.PromptTokens == nil && trace.CompletionTokens == nil && trace.TotalTokens == nil {
			parseCodeBuddyJSONOutput(trace, stderr)
		}
		// 如果 session 或 usage 仍不完整，再尝试从 JSONL 行中补齐
		if trace.SessionID == "" || trace.PromptTokens == nil && trace.CompletionTokens == nil && trace.TotalTokens == nil {
			parseCodeBuddyJSONLUsage(trace, fullOutput)
		}
	} else if strings.EqualFold(trace.Framework, "claudecode") || strings.EqualFold(trace.Framework, "claude_code") || strings.EqualFold(trace.Framework, "claude-code") {
		parseClaudeCodeJSONOutput(trace, stdout)
		if trace.PromptTokens == nil && trace.CompletionTokens == nil && trace.TotalTokens == nil {
			parseClaudeCodeJSONOutput(trace, stderr)
		}
		if trace.SessionID == "" || trace.PromptTokens == nil && trace.CompletionTokens == nil && trace.TotalTokens == nil {
			parseClaudeCodeJSONLUsage(trace, fullOutput)
		}
	}

	if trace.SessionID == "" {
		if match := sessionIDPattern.FindStringSubmatch(fullOutput); len(match) == 2 {
			trace.SessionID = match[1]
		}
	}

	records := collectUsageRecords(fullOutput)
	if len(records) == 0 {
		return
	}

	var promptSum, completionSum, totalSum int
	var promptSeen, completionSeen, totalSeen bool
	for _, record := range records {
		if record.Prompt != nil {
			promptSum += *record.Prompt
			promptSeen = true
		}
		if record.Completion != nil {
			completionSum += *record.Completion
			completionSeen = true
		}
		if record.Total != nil {
			totalSum += *record.Total
			totalSeen = true
		}
	}
	if promptSeen && trace.PromptTokens == nil {
		trace.PromptTokens = intPtr(promptSum)
	}
	if completionSeen && trace.CompletionTokens == nil {
		trace.CompletionTokens = intPtr(completionSum)
	}
	if totalSeen && trace.TotalTokens == nil {
		trace.TotalTokens = intPtr(totalSum)
	}
	if trace.TotalTokens == nil && trace.PromptTokens != nil && trace.CompletionTokens != nil {
		trace.TotalTokens = intPtr(*trace.PromptTokens + *trace.CompletionTokens)
	}
	if trace.TokenSource == "" {
		trace.TokenSource = "actual"
	}
	if strings.TrimSpace(trace.UsageSourceDetail) == "" {
		trace.UsageSourceDetail = "log_json"
	}
}

func collectOpenCodeSessionExport(
	ctx context.Context,
	sandboxRunner SandboxRunner,
	sandboxReq SandboxRunRequest,
	workRoot string,
	traceDir string,
	sampleID string,
	trace *AgentTrace,
) {
	if trace == nil || !strings.EqualFold(trace.Framework, "opencode") {
		return
	}
	if strings.TrimSpace(trace.SessionID) == "" {
		return
	}

	exportWorkspaceRel := filepath.Join(".utbench", "opencode", "session_export.json")
	exportHostPath := filepath.Join(workRoot, exportWorkspaceRel)
	exportContainerPath := "/workspace/" + filepath.ToSlash(exportWorkspaceRel)
	sessionExportPath := filepath.Join(traceDir, sampleID+".session_export.json")

	exportReq := sandboxReq
	exportReq.Command = fmt.Sprintf(
		"mkdir -p /workspace/.utbench/opencode && opencode export %s > %s",
		shQuote(trace.SessionID),
		shQuote(exportContainerPath),
	)
	if exportReq.TimeoutSeconds <= 0 || exportReq.TimeoutSeconds > 120 {
		exportReq.TimeoutSeconds = 120
	}

	result, err := sandboxRunner.Run(ctx, exportReq)
	if err != nil {
		trace.SessionExportError = fmt.Sprintf("export command failed: %v", err)
		return
	}
	if result.ExitCode != 0 {
		trace.SessionExportError = strings.TrimSpace(result.Stderr)
		if trace.SessionExportError == "" {
			trace.SessionExportError = fmt.Sprintf("export command exited with code %d", result.ExitCode)
		}
		return
	}

	raw, err := os.ReadFile(exportHostPath)
	if err != nil {
		trace.SessionExportError = fmt.Sprintf("read session export: %v", err)
		return
	}
	if err := os.WriteFile(sessionExportPath, raw, 0o644); err != nil {
		trace.SessionExportError = fmt.Sprintf("persist session export: %v", err)
		return
	}
	trace.SessionExportPath = sessionExportPath

	var payload any
	if err := json.Unmarshal(raw, &payload); err != nil {
		trace.SessionExportError = fmt.Sprintf("parse session export json: %v", err)
		return
	}

	records := extractUsageRecords(payload)
	if len(records) == 0 {
		return
	}

	var promptSum, completionSum, totalSum int
	var promptSeen, completionSeen, totalSeen bool
	for _, record := range records {
		if record.Prompt != nil {
			promptSum += *record.Prompt
			promptSeen = true
		}
		if record.Completion != nil {
			completionSum += *record.Completion
			completionSeen = true
		}
		if record.Total != nil {
			totalSum += *record.Total
			totalSeen = true
		}
	}
	if promptSeen {
		trace.PromptTokens = intPtr(promptSum)
	}
	if completionSeen {
		trace.CompletionTokens = intPtr(completionSum)
	}
	if totalSeen {
		trace.TotalTokens = intPtr(totalSum)
	}
	if trace.TotalTokens == nil && trace.PromptTokens != nil && trace.CompletionTokens != nil {
		trace.TotalTokens = intPtr(*trace.PromptTokens + *trace.CompletionTokens)
	}
	trace.TokenSource = "actual"
	trace.UsageSourceDetail = "opencode_session_export"
}

func collectUsageRecords(output string) []usageRecord {
	lines := strings.Split(output, "\n")
	seen := map[string]struct{}{}
	var out []usageRecord
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "{") {
			continue
		}
		var payload any
		if err := json.Unmarshal([]byte(line), &payload); err != nil {
			continue
		}
		records := extractUsageRecords(payload)
		for _, record := range records {
			keyBytes, _ := json.Marshal(record)
			key := string(keyBytes)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			out = append(out, record)
		}
	}
	return out
}

func extractUsageRecords(value any) []usageRecord {
	switch v := value.(type) {
	case map[string]any:
		var out []usageRecord
		if record, ok := usageRecordFromMap(v); ok {
			out = append(out, record)
		}
		for _, child := range v {
			out = append(out, extractUsageRecords(child)...)
		}
		return out
	case []any:
		var out []usageRecord
		for _, child := range v {
			out = append(out, extractUsageRecords(child)...)
		}
		return out
	default:
		return nil
	}
}

func usageRecordFromMap(v map[string]any) (usageRecord, bool) {
	prompt := firstIntValue(v, "prompt_tokens", "input_tokens", "inputTokens", "promptTokens")
	completion := firstIntValue(v, "completion_tokens", "output_tokens", "outputTokens", "completionTokens")
	total := firstIntValue(v, "total_tokens", "totalTokens")
	if (prompt == nil || completion == nil || total == nil) && v["tokens"] != nil {
		if nested, ok := v["tokens"].(map[string]any); ok {
			if prompt == nil {
				prompt = firstIntValue(nested, "input", "input_tokens", "prompt", "prompt_tokens")
			}
			if completion == nil {
				completion = firstIntValue(nested, "output", "output_tokens", "completion", "completion_tokens")
			}
			if total == nil {
				total = firstIntValue(nested, "total", "total_tokens")
			}
		}
	}
	if prompt == nil && completion == nil && total == nil {
		return usageRecord{}, false
	}
	if total == nil && prompt != nil && completion != nil {
		total = intPtr(*prompt + *completion)
	}
	return usageRecord{Prompt: prompt, Completion: completion, Total: total}, true
}

func firstIntValue(v map[string]any, keys ...string) *int {
	for _, key := range keys {
		raw, ok := v[key]
		if !ok {
			continue
		}
		if out, ok := asInt(raw); ok {
			return intPtr(out)
		}
	}
	return nil
}

func asInt(v any) (int, bool) {
	switch value := v.(type) {
	case int:
		return value, true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	case json.Number:
		i, err := value.Int64()
		return int(i), err == nil
	case string:
		value = strings.TrimSpace(value)
		if value == "" {
			return 0, false
		}
		i, err := strconv.Atoi(value)
		return i, err == nil
	default:
		return 0, false
	}
}

func finalizeAgentAccounting(trace *AgentTrace, _ string, _ string, model modelConfig) {
	switch {
	case trace.PromptTokens != nil && trace.CompletionTokens != nil:
		if trace.TotalTokens == nil {
			trace.TotalTokens = intPtr(*trace.PromptTokens + *trace.CompletionTokens)
		}
		if strings.TrimSpace(trace.TokenSource) == "" {
			trace.TokenSource = "actual"
		}
		if strings.TrimSpace(trace.UsageSourceDetail) == "" {
			trace.UsageSourceDetail = "agent_usage"
		}
	case trace.PromptTokens != nil || trace.CompletionTokens != nil || trace.TotalTokens != nil:
		if trace.TotalTokens == nil && trace.PromptTokens != nil && trace.CompletionTokens != nil {
			trace.TotalTokens = intPtr(*trace.PromptTokens + *trace.CompletionTokens)
		}
		if strings.TrimSpace(trace.TokenSource) == "" {
			trace.TokenSource = "partial"
		}
		if strings.TrimSpace(trace.UsageSourceDetail) == "" {
			trace.UsageSourceDetail = "agent_usage_partial"
		}
	default:
		// For CLI agents, final code size is only a weak lower bound and badly
		// undercounts multi-turn/tool-heavy sessions. If we do not have usage from
		// the agent itself, prefer marking the token data as missing rather than
		// publishing a misleadingly low number.
		trace.TokenSource = "missing"
		if strings.TrimSpace(trace.UsageSourceDetail) == "" {
			trace.UsageSourceDetail = "unavailable"
		}
	}

	cost, costSource := estimateCostUSD(model, trace.PromptTokens, trace.CompletionTokens)
	trace.EstimatedCost = cost
	trace.CostSource = costSource
}

func shQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'\''`) + "'"
}

func estimateCostUSD(model modelConfig, promptTokens, completionTokens *int) (*float64, string) {
	if promptTokens == nil || completionTokens == nil {
		return nil, "unavailable"
	}
	if model.Pricing.PromptPer1KUSD <= 0 && model.Pricing.CompletionPer1KUSD <= 0 {
		return nil, "unavailable"
	}
	cost := (float64(*promptTokens) / 1000.0 * model.Pricing.PromptPer1KUSD) +
		(float64(*completionTokens) / 1000.0 * model.Pricing.CompletionPer1KUSD)
	return floatPtr(cost), "configured_pricing"
}

func extractLikelyPaths(line string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, match := range pathPattern.FindAllStringSubmatch(line, -1) {
		if len(match) != 2 {
			continue
		}
		candidate := strings.Trim(match[1], "\"'`,:;()[]{}")
		if !looksLikeFilePath(candidate) {
			continue
		}
		if _, ok := seen[candidate]; ok {
			continue
		}
		seen[candidate] = struct{}{}
		out = append(out, candidate)
	}
	sort.Strings(out)
	return out
}

func extractRelevantPaths(line string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, candidate := range extractLikelyPaths(line) {
		if !isRelevantWorkspacePath(candidate) {
			continue
		}
		if _, ok := seen[candidate]; ok {
			continue
		}
		seen[candidate] = struct{}{}
		out = append(out, candidate)
	}
	for _, match := range extTokenPattern.FindAllStringSubmatch(line, -1) {
		if len(match) != 2 {
			continue
		}
		candidate := strings.TrimSpace(strings.Trim(match[1], "\"'`"))
		if !isRelevantWorkspacePath(candidate) {
			continue
		}
		if _, ok := seen[candidate]; ok {
			continue
		}
		seen[candidate] = struct{}{}
		out = append(out, candidate)
	}
	sort.Strings(out)
	return out
}

func isRelevantWorkspacePath(candidate string) bool {
	candidate = strings.TrimSpace(filepath.ToSlash(candidate))
	if candidate == "" {
		return false
	}
	candidate = strings.TrimPrefix(candidate, "./")
	if strings.HasPrefix(candidate, "/workspace/.utbench/") ||
		strings.HasPrefix(candidate, ".utbench/") ||
		strings.HasPrefix(candidate, "/workspace/test_runner") ||
		candidate == "test_runner" ||
		strings.HasPrefix(candidate, "/workspace/go.mod") ||
		candidate == "go.mod" {
		return false
	}
	base := filepath.Base(candidate)
	switch {
	case strings.HasPrefix(candidate, "/workspace/"):
		return true
	case base == "utbench_agent_prompt.md":
		return true
	case strings.HasPrefix(base, "generated_test."):
		return true
	default:
		return false
	}
}

func isRelevantWorkspaceWritePath(candidate string) bool {
	candidate = strings.TrimSpace(filepath.ToSlash(candidate))
	if !isRelevantWorkspacePath(candidate) {
		return false
	}
	return filepath.Base(candidate) != "utbench_agent_prompt.md"
}

func looksLikeRealCommand(cmd string) bool {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" || len(cmd) > 500 {
		return false
	}
	if strings.HasPrefix(strings.ToLower(cmd), "info ") {
		return false
	}
	fields := strings.Fields(cmd)
	if len(fields) == 2 {
		if _, err := strconv.Atoi(fields[0]); err == nil && strings.HasPrefix(fields[1], "/tmp/") {
			return false
		}
	}
	return true
}

func intPtr(v int) *int {
	return &v
}

func floatPtr(v float64) *float64 {
	return &v
}

// containerPath 根据沙箱模式返回路径。
// docker 模式返回容器内路径（如 /workspace/...），local 模式返回实际本地路径。
func containerPath(dockerPath, localPath, mode string) string {
	if strings.EqualFold(mode, "local") {
		return localPath
	}
	return dockerPath
}
