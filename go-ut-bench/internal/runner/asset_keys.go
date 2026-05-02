package runner

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go-ut-bench/internal/agentconfig"
	"go-ut-bench/internal/contracts"
)

type generationIdentity struct {
	SampleUID                string
	SubjectVersionID         string
	GenerationKey            string
	DependencyFingerprint    string
	GenerationEnvFingerprint string
	PromptRenderingSHA256    string
}

func buildGenerationIdentity(target subjectTarget, model modelConfig, sample contracts.SampleRef, sourceCode []byte, renderedPrompt, promptVersionID string) generationIdentity {
	sourceSHA := sha256Bytes(sourceCode)
	repoManifestSHA := ""
	sampleContentSHA := sourceSHA
	if meta := loadRepoLevelMetaForRunner(sample.Path); meta != nil {
		if root := resolveRepoWorkspaceRoot(sample.Path, meta.WorkspaceRoot); root != "" {
			if sha := hashDirectory(root); sha != "" {
				sampleContentSHA = sha
			}
		}
		if sha := hashExistingFile(filepath.Join(filepath.Dir(sample.Path), "meta.json")); sha != "" {
			repoManifestSHA = sha
		}
	}

	dependencyFingerprint := dependencyFingerprintForSample(sample)
	subjectVersionID := subjectVersionIDForTarget(target, model, sample.Language)
	generationEnvFingerprint := generationEnvFingerprintForTarget(target, sample.Language)
	promptSHA := sha256String(renderedPrompt)
	sampleUID := assetStableID(
		"dataset_sample",
		sample.Language,
		string(sample.Category),
		sample.Scenario,
		sample.ID,
		sampleContentSHA,
		repoManifestSHA,
	)
	generationKey := assetStableID(
		"generation",
		subjectVersionID,
		sampleUID,
		promptSHA,
		promptVersionID,
		sample.Language,
		string(sample.Category),
		dependencyFingerprint,
		generationEnvFingerprint,
	)
	return generationIdentity{
		SampleUID:                sampleUID,
		SubjectVersionID:         subjectVersionID,
		GenerationKey:            generationKey,
		DependencyFingerprint:    dependencyFingerprint,
		GenerationEnvFingerprint: generationEnvFingerprint,
		PromptRenderingSHA256:    promptSHA,
	}
}

func subjectVersionIDForTarget(target subjectTarget, model modelConfig, language string) string {
	payload := map[string]any{
		"subject":             target.subject.Spec,
		"model":               modelVersionPayload(model),
		"framework":           frameworkVersionPayload(target.subject.Framework, language),
		"skill":               skillVersionPayload(target.subject.Skill),
		"docker_image":        frameworkDockerImage(target.subject.Framework, language),
		"sandbox_fingerprint": generationEnvFingerprintForTarget(target, language),
	}
	return assetStableID("subject_version", sha256JSON(payload))
}

func generationEnvFingerprintForTarget(target subjectTarget, language string) string {
	fw := target.subject.Framework
	payload := map[string]any{
		"sandbox_mode":     fw.SandboxMode,
		"docker_image":     frameworkDockerImage(fw, language),
		"network_disabled": fw.NetworkDisabled,
		"cpu":              fw.CPU,
		"memory":           fw.Memory,
		"env_keys":         sortedStringMapKeys(fw.Env),
		"env_from_host":    uniqueSortedStrings(fw.EnvFromHost),
		"preflight":        fw.Preflight,
	}
	return sha256JSON(payload)
}

func modelVersionPayload(model modelConfig) map[string]any {
	return map[string]any{
		"name":        model.Name,
		"provider":    model.Provider,
		"endpoint":    model.Endpoint,
		"model":       model.Model,
		"api_key_env": model.APIKeyEnv,
		"params":      model.Params,
		"pricing":     model.Pricing,
	}
}

func frameworkVersionPayload(fw agentconfig.FrameworkSpec, language string) map[string]any {
	return map[string]any{
		"name":                       fw.Name,
		"kind":                       fw.Kind,
		"command_sha256":             sha256String(fw.Command),
		"docker_image":               frameworkDockerImage(fw, language),
		"sandbox_mode":               fw.SandboxMode,
		"timeout_seconds":            fw.TimeoutSeconds,
		"output_globs":               fw.OutputGlobs,
		"env_keys":                   sortedStringMapKeys(fw.Env),
		"env_from_host":              uniqueSortedStrings(fw.EnvFromHost),
		"preflight":                  fw.Preflight,
		"forbidden_command_patterns": fw.ForbiddenCommandPatterns,
		"network_disabled":           fw.NetworkDisabled,
		"cpu":                        fw.CPU,
		"memory":                     fw.Memory,
	}
}

func skillVersionPayload(skill contracts.SkillSpec) map[string]any {
	fileHashes := map[string]string{}
	if skill.InstructionPath != "" {
		fileHashes[filepath.ToSlash(skill.InstructionPath)] = hashExistingFile(skill.InstructionPath)
	}
	for _, path := range skill.Files {
		fileHashes[filepath.ToSlash(path)] = hashExistingFile(path)
	}
	return map[string]any{
		"name":                  skill.Name,
		"version":               skill.Version,
		"description":           skill.Description,
		"inject_mode":           skill.InjectMode,
		"compatible_frameworks": skill.CompatibleFrameworks,
		"compatible_languages":  skill.CompatibleLanguages,
		"files":                 fileHashes,
	}
}

func dependencyFingerprintForSample(sample contracts.SampleRef) string {
	root := filepath.Dir(sample.Path)
	if meta := loadRepoLevelMetaForRunner(sample.Path); meta != nil {
		if resolved := resolveRepoWorkspaceRoot(sample.Path, meta.WorkspaceRoot); resolved != "" {
			root = resolved
		}
	}
	names := []string{
		"requirements.txt",
		"requirements-dev.txt",
		"pyproject.toml",
		"poetry.lock",
		"go.mod",
		"go.sum",
		"pom.xml",
		"build.gradle",
		"gradle.lockfile",
		"CMakeLists.txt",
		"vcpkg.json",
		"conanfile.txt",
	}
	parts := make([]string, 0, len(names))
	for _, name := range names {
		path := filepath.Join(root, name)
		if sha := hashExistingFile(path); sha != "" {
			parts = append(parts, filepath.ToSlash(name)+"="+sha)
		}
	}
	sort.Strings(parts)
	return sha256String(strings.Join(parts, "\n"))
}

func resolveRepoWorkspaceRoot(samplePath, workspaceRoot string) string {
	workspaceRoot = strings.TrimSpace(workspaceRoot)
	if workspaceRoot == "" {
		return ""
	}
	if filepath.IsAbs(workspaceRoot) {
		return workspaceRoot
	}
	return filepath.Join(filepath.Dir(samplePath), workspaceRoot)
}

func hashDirectory(root string) string {
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return ""
	}
	var parts []string
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if d.IsDir() {
			switch name {
			case ".git", ".utbench", "__pycache__", "node_modules", "target", "build", "dist":
				return filepath.SkipDir
			}
			return nil
		}
		sha := hashExistingFile(path)
		if sha == "" {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}
		parts = append(parts, filepath.ToSlash(rel)+"="+sha)
		return nil
	})
	sort.Strings(parts)
	return sha256String(strings.Join(parts, "\n"))
}

func hashExistingFile(path string) string {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return sha256Bytes(raw)
}

func assetStableID(prefix string, parts ...string) string {
	return prefix + "_" + sha256String(strings.Join(parts, "\x00"))[:16]
}

func sha256String(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func sha256JSON(value any) string {
	raw, _ := json.Marshal(value)
	return sha256String(string(raw))
}

func sortedStringMapKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
