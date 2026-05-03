package runner

import (
	"context"
	"time"
)

// generateModelAPI 调用模型 API 生成测试（纯模型基线）。
func generateModelAPI(ctx context.Context, req AgentGenerateRequest) AgentGenerateResult {
	client := newAPIClient()
	waitModelInterval(req.Model.Name)
	started := time.Now()

	generated, response, latency, pTok, cTok, tTok, truncated, genErr := client.generateTest(
		ctx,
		req.Model,
		req.Sample.Language,
		req.Prompt,
	)
	finished := started.Add(time.Duration(latency) * time.Millisecond)

	trace := AgentTrace{
		SubjectID:        req.Subject.Spec.ID,
		Framework:        "model_api",
		Model:            req.Model.Name,
		Skill:            "no_skill",
		SampleID:         req.Sample.ID,
		Language:         req.Sample.Language,
		ExitCode:         0,
		DurationMS:       latency,
		StartedAt:        started.UTC(),
		FinishedAt:       finished.UTC(),
		PromptTokens:     pTok,
		CompletionTokens: cTok,
		TotalTokens:      tTok,
		InteractionCount: 1, // API 调用固定 1 轮
	}
	finalizeAgentAccounting(&trace, req.Prompt, generated, req.Model)

	if genErr != nil {
		trace.ExitCode = 1
		trace.Stderr = genErr.Message
	}

	return AgentGenerateResult{
		Code:             generated,
		RawResponse:      response,
		Trace:            trace,
		LatencyMS:        latency,
		PromptTokens:     pTok,
		CompletionTokens: cTok,
		TotalTokens:      tTok,
		TokenSource:      trace.TokenSource,
		EstimatedCostUSD: trace.EstimatedCost,
		CostSource:       trace.CostSource,
		Truncated:        truncated,
		Error:            genErr,
	}
}
