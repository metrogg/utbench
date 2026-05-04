package runner

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go-ut-bench/internal/contracts"
)

type generationTaskPlan struct {
	PromptMode     string
	PromptPath     string
	RenderedPrompt string
	Identity       generationIdentity
	ReadError      error
	Reused         *contracts.ReusableGeneratedCase
}

type generationReusePlan struct {
	ByTaskKey    map[string]generationTaskPlan
	ReusableHits int
}

func prepareGenerationReusePlan(
	ctx context.Context,
	spec contracts.RunSpec,
	subjects []subjectTarget,
	samples []contracts.SampleRef,
	promptRoot string,
	promptVersionID string,
	reuseStore GenerationReuseStore,
) generationReusePlan {
	out := generationReusePlan{ByTaskKey: make(map[string]generationTaskPlan)}
	if spec.DryRun {
		return out
	}
	for _, target := range subjects {
		for _, sample := range samples {
			if !subjectSupportsLanguage(target, sample.Language) {
				continue
			}
			taskID := taskKey(target.subject.Spec.ID, sample.Language, sample.ID)
			promptMode := string(PromptModeFullFile)
			if loadRepoLevelMetaForRunner(sample.Path) != nil {
				promptMode = string(PromptModeRepoLevel)
			}
			plan := generationTaskPlan{PromptMode: promptMode}
			sourceCode, err := os.ReadFile(sample.Path)
			if err != nil {
				plan.ReadError = err
				out.ByTaskKey[taskID] = plan
				continue
			}
			renderedPrompt := buildPrompt(sample.Language, sample.Path, string(sourceCode))
			promptPath := filepath.Join(promptRoot, "rendered", target.subject.Spec.ID, sample.Language, fmt.Sprintf("%s.prompt.txt", sample.ID))
			if mkErr := os.MkdirAll(filepath.Dir(promptPath), 0o755); mkErr == nil {
				if writeErr := os.WriteFile(promptPath, []byte(renderedPrompt), 0o644); writeErr == nil {
					plan.PromptPath = promptPath
				}
			}
			plan.RenderedPrompt = renderedPrompt
			plan.Identity = buildGenerationIdentity(target, target.model, sample, sourceCode, renderedPrompt, promptVersionID)
			if reuseStore != nil {
				if reused, ok, err := reuseStore.FindReusableGeneratedAsset(ctx, plan.Identity.GenerationKey); err == nil && ok {
					plan.Reused = &reused
					out.ReusableHits++
				}
			}
			out.ByTaskKey[taskID] = plan
		}
	}
	return out
}
