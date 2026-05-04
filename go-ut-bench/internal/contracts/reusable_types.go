package contracts

// ReusableGeneratedCase 可复用的生成结果
// 用于跨运行复用已生成的测试，避免重复 LLM 调用
type ReusableGeneratedCase struct {
	GeneratedCaseID       string   `json:"generated_case_id"`
	RunID                 string   `json:"run_id"`
	Model                 string   `json:"model"`
	Language              string   `json:"language"`
	SampleID              string   `json:"sample_id"`
	GeneratedTestPath     string   `json:"generated_test_path"`
	GeneratedTestSHA256   string   `json:"generated_test_sha256,omitempty"`
	ResponsePath          string   `json:"response_path,omitempty"`
	MetadataPath          string   `json:"metadata_path,omitempty"`
	TracePath             string   `json:"trace_path,omitempty"`
	WorkspaceDiffPath     string   `json:"workspace_diff_path,omitempty"`
	SandboxFingerprint    string   `json:"sandbox_fingerprint,omitempty"`
	TokenSource           string   `json:"token_source,omitempty"`
	CostSource            string   `json:"cost_source,omitempty"`
	EstimatedCostUSD      *float64 `json:"estimated_cost_usd,omitempty"`
	PromptVersionID       string   `json:"prompt_version_id,omitempty"`
	PromptMode            string   `json:"prompt_mode,omitempty"`
	LatencyMS             int      `json:"latency_ms,omitempty"`
	PromptTokens          *int     `json:"prompt_tokens,omitempty"`
	CompletionTokens      *int     `json:"completion_tokens,omitempty"`
	TotalTokens           *int     `json:"total_tokens,omitempty"`
	GeneratedAtUTC        string   `json:"generated_at_utc,omitempty"`
	GeneratedTestArtifact string   `json:"generated_test_artifact_id,omitempty"`
}
