package runner

import (
	"context"
	"time"

	"go-ut-bench/internal/agentconfig"
	"go-ut-bench/internal/contracts"
)

// AgentAdapter 是所有 Agent 生成器的统一接口。
// 每种 agent 类型（model_api、cli_agent、http_agent 等）实现此接口。
type AgentAdapter interface {
	Generate(ctx context.Context, req AgentGenerateRequest) AgentGenerateResult
}

// AgentGenerateRequest 封装一次 Agent 生成任务的全部输入。
type AgentGenerateRequest struct {
	Subject    agentconfig.ResolvedSubject
	Model      modelConfig
	Sample     contracts.SampleRef
	Prompt     string
	TestPath   string
	MetaRoot   string
	OutputRoot string
	RunID      string
}

// AgentGenerateResult 封装一次 Agent 生成任务的全部输出。
type AgentGenerateResult struct {
	Code             string
	RawResponse      map[string]any
	Trace            AgentTrace
	LatencyMS        int
	PromptTokens     *int
	CompletionTokens *int
	TotalTokens      *int
	TokenSource      string
	EstimatedCostUSD *float64
	CostSource       string
	Truncated        bool
	Error            *contracts.ErrorInfo
}

// AgentTrace 记录 Agent 执行的完整操作轨迹。
type AgentTrace struct {
	// 基础标识
	SubjectID string `json:"subject_id"`
	Framework string `json:"framework"`
	Model     string `json:"model"`
	Skill     string `json:"skill"`
	SampleID  string `json:"sample_id"`
	Language  string `json:"language"`

	// 执行信息
	Command    string    `json:"command"`
	ExitCode   int       `json:"exit_code"`
	DurationMS int       `json:"duration_ms"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`

	// Token 消耗（从 Agent 日志解析，不一定所有框架都有）
	PromptTokens      *int     `json:"prompt_tokens,omitempty"`
	CompletionTokens  *int     `json:"completion_tokens,omitempty"`
	TotalTokens       *int     `json:"total_tokens,omitempty"`
	TokenSource       string   `json:"token_source,omitempty"`
	EstimatedCost     *float64 `json:"estimated_cost,omitempty"`
	CostSource        string   `json:"cost_source,omitempty"`
	UsageSourceDetail string   `json:"usage_source_detail,omitempty"`

	// Agent 交互过程
	InteractionCount int              `json:"interaction_count"`           // Agent 交互轮次
	ToolCalls        []ToolCall       `json:"tool_calls,omitempty"`        // 工具调用序列
	FilesRead        []string         `json:"files_read,omitempty"`        // Agent 读取的文件
	FilesWritten     []string         `json:"files_written,omitempty"`     // Agent 写入/修改的文件
	CommandsExecuted []string         `json:"commands_executed,omitempty"` // Agent 执行的 shell 命令
	EnvironmentSetup []PreflightCheck `json:"environment_setup,omitempty"` // 平台托管的依赖准备步骤
	PreflightChecks  []PreflightCheck `json:"preflight_checks,omitempty"`  // 执行前环境检查

	// 输出
	Stdout        string   `json:"stdout"`
	Stderr        string   `json:"stderr"`
	WorkspaceDiff []string `json:"workspace_diff"` // workspace 文件变更列表

	// 沙箱与产物路径
	SandboxProvider    string `json:"sandbox_provider,omitempty"`
	SandboxImage       string `json:"sandbox_image,omitempty"`
	SandboxFingerprint string `json:"sandbox_fingerprint"`
	SessionID          string `json:"session_id,omitempty"`
	SessionExportPath  string `json:"session_export_path,omitempty"`
	SessionExportError string `json:"session_export_error,omitempty"`
	TracePath          string `json:"trace_path"`
	WorkspaceDiffPath  string `json:"workspace_diff_path"`
}

// ToolCall 记录一次 Agent 工具调用。
type ToolCall struct {
	Tool       string `json:"tool"`                  // 工具名：read_file, write_file, bash, search...
	Input      string `json:"input,omitempty"`       // 输入参数（截断）
	Output     string `json:"output,omitempty"`      // 输出摘要（截断）
	DurationMS int    `json:"duration_ms,omitempty"` // 耗时
	Success    bool   `json:"success"`               // 是否成功
}

// PreflightCheck 记录一次沙箱预检命令的执行结果。
type PreflightCheck struct {
	Command    string `json:"command"`
	ExitCode   int    `json:"exit_code"`
	DurationMS int    `json:"duration_ms"`
	Stdout     string `json:"stdout,omitempty"`
	Stderr     string `json:"stderr,omitempty"`
	Passed     bool   `json:"passed"`
}

// cliAgentAdapter 实现 AgentAdapter 接口，负责在沙箱中执行 CLI Agent。
type cliAgentAdapter struct {
	sandboxRunner SandboxRunner
}

func newCLIAgentAdapter(runner SandboxRunner) AgentAdapter {
	return &cliAgentAdapter{sandboxRunner: runner}
}

func (a *cliAgentAdapter) Generate(ctx context.Context, req AgentGenerateRequest) AgentGenerateResult {
	return generateCLIAgent(ctx, a.sandboxRunner, req)
}

// modelAPIAdapter 实现 AgentAdapter 接口，负责调用模型 API 生成测试。
type modelAPIAdapter struct{}

func newModelAPIAdapter() AgentAdapter {
	return &modelAPIAdapter{}
}

func (a *modelAPIAdapter) Generate(ctx context.Context, req AgentGenerateRequest) AgentGenerateResult {
	return generateModelAPI(ctx, req)
}
