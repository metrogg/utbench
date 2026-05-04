package reporter

import (
	"sort"
	"strings"

	"go-ut-bench/internal/contracts"
)

func buildRuntimeSummary(manifest *contracts.GeneratedManifest, set contracts.EvaluationResultSet) contracts.RuntimeSummary {
	summary := contracts.RuntimeSummary{
		EvaluatorEnvFingerprint: strings.TrimSpace(set.EnvironmentFingerprint),
	}
	providers := map[string]struct{}{}
	images := map[string]struct{}{}
	frameworks := map[string]struct{}{}

	if manifest != nil {
		for _, item := range manifest.Cases {
			if provider := strings.TrimSpace(item.SandboxProvider); provider != "" {
				providers[provider] = struct{}{}
				switch strings.ToLower(provider) {
				case "docker":
					summary.DockerBackedSubjects++
				case "local":
					summary.LocalBackedSubjects++
				}
			}
			if image := strings.TrimSpace(item.DockerImage); image != "" {
				images[image] = struct{}{}
			}
			if fw := strings.TrimSpace(item.AgentFramework); fw != "" {
				frameworks[fw] = struct{}{}
			}
		}
	}

	summary.SandboxProviders = sortedKeysFromSet(providers)
	summary.SandboxImages = sortedKeysFromSet(images)
	summary.AgentFrameworks = sortedKeysFromSet(frameworks)

	if summary.EvaluatorEnvFingerprint != "" {
		summary.Notes = append(summary.Notes, "评测环境指纹已写入本次报告，可用于跨 run 对比。")
	}
	if len(summary.SandboxImages) > 0 {
		summary.Notes = append(summary.Notes, "Agent 沙箱镜像与评测环境已分离管理；镜像变化应视为新的实验条件。")
	}
	return summary
}

func sortedKeysFromSet(values map[string]struct{}) []string {
	if len(values) == 0 {
		return nil
	}
	out := make([]string, 0, len(values))
	for key := range values {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}
