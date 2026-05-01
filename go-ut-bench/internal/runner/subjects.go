package runner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"go-ut-bench/internal/agentconfig"
	"go-ut-bench/internal/contracts"
)

type subjectTrace struct {
	TracePath          string
	WorkspaceDiffPath  string
	SandboxFingerprint string
}

type commandTemplateData struct {
	Workspace         string
	PromptFile        string
	OutputFile        string
	SkillDir          string
	Model             string
	ModelID           string
	ModelProvider     string
	ModelEndpoint     string
	ModelAPIKeyEnv    string
	Framework         string
	SubjectID         string
	Skill             string
	Language          string
	SampleID          string
	SourceFile        string
	ContainerWorkdir  string
	ContainerPrompt   string
	ContainerOutput   string
	ContainerSkillDir string
}

func loadSubjectTargets(spec contracts.RunSpec, models []modelConfig) ([]subjectTarget, error) {
	modelNames := getModelNames(models)
	subjects, err := agentconfig.Load(spec.AgentsConfigPath, modelNames, spec.Subjects)
	if err != nil {
		return nil, err
	}
	modelByName := make(map[string]modelConfig, len(models))
	for _, model := range models {
		modelByName[model.Name] = model
	}
	out := make([]subjectTarget, 0, len(subjects))
	for _, subject := range subjects {
		modelCfg, ok := modelByName[subject.Spec.Model]
		if !ok {
			return nil, fmt.Errorf("subject %s references unavailable model %s", subject.Spec.ID, subject.Spec.Model)
		}
		out = append(out, subjectTarget{subject: subject, model: modelCfg})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no subjects selected")
	}
	return out, nil
}

func countEligibleSubjectTasks(subjects []subjectTarget, samples []contracts.SampleRef) int {
	total := 0
	for _, subject := range subjects {
		for _, sample := range samples {
			if subjectSupportsLanguage(subject, sample.Language) {
				total++
			}
		}
	}
	return total
}

func subjectSupportsLanguage(subject subjectTarget, language string) bool {
	return stringInAllowList(language, subject.subject.Framework.CompatibleLangs) &&
		stringInAllowList(language, subject.subject.Skill.CompatibleLanguages)
}

func stringInAllowList(value string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	value = strings.ToLower(strings.TrimSpace(value))
	for _, item := range allowed {
		if strings.ToLower(strings.TrimSpace(item)) == value {
			return true
		}
	}
	return false
}

func (s *Service) generateWithSubject(
	ctx context.Context,
	spec contracts.RunSpec,
	target subjectTarget,
	sample contracts.SampleRef,
	prompt string,
	testPath string,
	metaRoot string,
) (string, map[string]any, subjectTrace, int, *int, *int, *int, bool, *contracts.ErrorInfo) {
	prompt = appendSkillInstruction(prompt, target.subject.Skill)
	switch target.subject.Spec.Kind {
	case agentconfig.KindModelAPI, "":
		return s.generateWithModelAPI(ctx, target, sample, prompt)
	case agentconfig.KindCLIAgent:
		return s.generateWithCLIAgent(ctx, spec, target, sample, prompt, testPath, metaRoot)
	default:
		return "", nil, subjectTrace{}, 0, nil, nil, nil, false, &contracts.ErrorInfo{
			Kind:      "unsupported_subject_kind",
			Message:   fmt.Sprintf("unsupported subject kind: %s", target.subject.Spec.Kind),
			Retryable: false,
		}
	}
}

func (s *Service) generateWithModelAPI(ctx context.Context, target subjectTarget, sample contracts.SampleRef, prompt string) (string, map[string]any, subjectTrace, int, *int, *int, *int, bool, *contracts.ErrorInfo) {
	client := newAPIClient()
	waitModelInterval(target.model.Name)
	generated, response, latency, pTok, cTok, tTok, truncated, genErr := client.generateTest(
		ctx,
		target.model,
		sample.Language,
		prompt,
	)
	subjectID := target.subject.Spec.ID
	s.logger.LogAPIRequest(subjectID, sample.Language, sample.ID, 0, latency)
	s.logger.LogAPIResponse(subjectID, sample.Language, sample.ID, genErr == nil, truncated, errorMsgSafe(genErr))
	s.logger.ToFile("runner").Trace("generate_response",
		"subject_id", subjectID,
		"model", target.model.Name,
		"language", sample.Language,
		"sample_id", sample.ID,
		"prompt_tokens", pTok,
		"completion_tokens", cTok,
		"total_tokens", tTok,
		"latency_ms", latency,
		"truncated", truncated,
		"success", genErr == nil,
	)
	return generated, response, subjectTrace{}, latency, pTok, cTok, tTok, truncated, genErr
}

func (s *Service) generateWithCLIAgent(
	ctx context.Context,
	spec contracts.RunSpec,
	target subjectTarget,
	sample contracts.SampleRef,
	prompt string,
	testPath string,
	metaRoot string,
) (string, map[string]any, subjectTrace, int, *int, *int, *int, bool, *contracts.ErrorInfo) {
	subjectID := target.subject.Spec.ID
	workRoot := filepath.Join(spec.OutputRoot, "runs", spec.RunID, "agent_workspaces", subjectID, sample.Language, sample.ID)
	if err := os.RemoveAll(workRoot); err != nil {
		return "", nil, subjectTrace{}, 0, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "workspace_error", Message: err.Error(), Retryable: false}
	}
	if err := os.MkdirAll(workRoot, 0o755); err != nil {
		return "", nil, subjectTrace{}, 0, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "workspace_error", Message: err.Error(), Retryable: false}
	}
	sourceFile, err := prepareAgentWorkspace(workRoot, sample)
	if err != nil {
		return "", nil, subjectTrace{}, 0, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "workspace_error", Message: err.Error(), Retryable: false}
	}
	outputFile := filepath.Join(workRoot, "generated_test"+languageExt(sample.Language))
	skillDir, err := injectSkillWorkspace(workRoot, target.subject.Skill)
	if err != nil {
		return "", nil, subjectTrace{}, 0, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "skill_injection_error", Message: err.Error(), Retryable: false}
	}
	outputHint := outputFile
	sourceHint := sourceFile
	skillHint := skillDir
	if !strings.EqualFold(target.subject.Framework.SandboxMode, "local") {
		outputHint = "/workspace/generated_test" + languageExt(sample.Language)
		sourceHint = "/workspace/" + filepath.ToSlash(mustRel(workRoot, sourceFile))
		if skillDir != "" {
			skillHint = "/workspace/" + filepath.ToSlash(mustRel(workRoot, skillDir))
		}
	}
	agentPrompt := buildAgentPrompt(prompt, sample, sourceHint, outputHint, skillHint)
	promptFile := filepath.Join(workRoot, "utbench_agent_prompt.md")
	if err := os.WriteFile(promptFile, []byte(agentPrompt), 0o644); err != nil {
		return "", nil, subjectTrace{}, 0, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "workspace_error", Message: err.Error(), Retryable: false}
	}
	before, _ := snapshotWorkspace(workRoot)
	traceDir := filepath.Join(metaRoot, "agent_traces", subjectID, sample.Language)
	tracePath := filepath.Join(traceDir, sample.ID+".trace.jsonl")
	diffPath := filepath.Join(traceDir, sample.ID+".diff.json")
	if err := os.MkdirAll(traceDir, 0o755); err != nil {
		return "", nil, subjectTrace{}, 0, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "trace_error", Message: err.Error(), Retryable: false}
	}

	started := time.Now()
	templateData := commandTemplateData{
		Workspace:         workRoot,
		PromptFile:        promptFile,
		OutputFile:        outputFile,
		SkillDir:          skillDir,
		Model:             target.model.Name,
		ModelID:           target.model.Model,
		ModelProvider:     target.model.Provider,
		ModelEndpoint:     target.model.Endpoint,
		ModelAPIKeyEnv:    target.model.APIKeyEnv,
		Framework:         target.subject.Spec.Framework,
		SubjectID:         subjectID,
		Skill:             target.subject.Spec.Skill,
		Language:          sample.Language,
		SampleID:          sample.ID,
		SourceFile:        sourceFile,
		ContainerWorkdir:  "/workspace",
		ContainerPrompt:   "/workspace/utbench_agent_prompt.md",
		ContainerOutput:   "/workspace/generated_test" + languageExt(sample.Language),
		ContainerSkillDir: "/workspace/.utbench/skills/" + safePathName(target.subject.Skill.Name),
	}
	cmdText, err := renderTemplateText("agent-command", target.subject.Framework.Command, templateData)
	if err != nil {
		return "", nil, subjectTrace{}, 0, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "command_template_error", Message: err.Error(), Retryable: false}
	}
	if strings.TrimSpace(cmdText) == "" {
		return "", nil, subjectTrace{}, 0, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "command_template_error", Message: "cli agent command is empty", Retryable: false}
	}
	envMap, envFromHost, err := buildAgentEnv(target.subject.Framework, target.model, templateData)
	if err != nil {
		return "", nil, subjectTrace{}, 0, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "agent_env_error", Message: err.Error(), Retryable: false}
	}
	req := buildSandboxRunRequest(spec.OutputRoot, target.subject.Framework, workRoot, cmdText, envMap, envFromHost)
	runOutput, runErr := s.sandboxRunner.Run(ctx, req)
	latency := int(time.Since(started).Milliseconds())
	after, _ := snapshotWorkspace(workRoot)
	changes := diffSnapshots(before, after)
	_ = writeTrace(tracePath, map[string]any{
		"ts_utc":      time.Now().UTC(),
		"subject_id":  subjectID,
		"framework":   target.subject.Spec.Framework,
		"model":       target.subject.Spec.Model,
		"skill":       target.subject.Spec.Skill,
		"sample_id":   sample.ID,
		"language":    sample.Language,
		"command":     cmdText,
		"exit_code":   runOutput.ExitCode,
		"duration_ms": latency,
		"stdout":      trimText(runOutput.Stdout, 4000),
		"stderr":      trimText(runOutput.Stderr, 4000),
	})
	_ = contracts.WriteJSON(diffPath, map[string]any{
		"workspace":  workRoot,
		"subject_id": subjectID,
		"changes":    changes,
	})
	trace := subjectTrace{
		TracePath:          tracePath,
		WorkspaceDiffPath:  diffPath,
		SandboxFingerprint: sandboxFingerprintForRequest(req),
	}
	rawResponse := map[string]any{
		"adapter":             "cli_agent",
		"subject_id":          subjectID,
		"framework":           target.subject.Spec.Framework,
		"model":               target.subject.Spec.Model,
		"skill":               target.subject.Spec.Skill,
		"command":             cmdText,
		"exit_code":           runOutput.ExitCode,
		"latency_ms":          latency,
		"trace_path":          tracePath,
		"workspace_diff_path": diffPath,
		"stdout":              trimText(runOutput.Stdout, 4000),
		"stderr":              trimText(runOutput.Stderr, 4000),
	}
	if runErr != nil {
		return "", rawResponse, trace, latency, nil, nil, nil, false, &contracts.ErrorInfo{
			Kind:      "agent_execution_error",
			Message:   fmt.Sprintf("agent command failed: %s", trimText(runOutput.Stderr+"\n"+runErr.Error(), 1000)),
			Retryable: false,
		}
	}
	generatedPath := findGeneratedTest(workRoot, outputFile, target.subject.Framework.OutputGlobs, changes, sample.Language)
	if generatedPath == "" {
		return "", rawResponse, trace, latency, nil, nil, nil, false, &contracts.ErrorInfo{
			Kind:      "agent_output_error",
			Message:   "agent did not produce a test file",
			Retryable: false,
		}
	}
	raw, err := os.ReadFile(generatedPath)
	if err != nil {
		return "", rawResponse, trace, latency, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "agent_output_error", Message: err.Error(), Retryable: false}
	}
	code := strings.TrimSpace(string(raw))
	if err := validateGeneratedTest(code, sample.Language); err != nil {
		return code, rawResponse, trace, latency, nil, nil, nil, false, &contracts.ErrorInfo{Kind: "quality_error", Message: err.Error(), Retryable: false}
	}
	rawResponse["generated_test_source_path"] = generatedPath
	rawResponse["generated_test_path"] = testPath
	return code, rawResponse, trace, latency, nil, nil, nil, false, nil
}

func prepareAgentWorkspace(workRoot string, sample contracts.SampleRef) (string, error) {
	if meta := loadModuleLevelMetaForRunner(sample.Path); meta != nil && meta.WorkspaceRoot != "" {
		sourceRoot := meta.WorkspaceRoot
		if !filepath.IsAbs(sourceRoot) {
			sourceRoot = filepath.Join(filepath.Dir(sample.Path), sourceRoot)
		}
		if err := copyDir(sourceRoot, workRoot); err != nil {
			return "", err
		}
		if meta.TargetFile != "" {
			return filepath.Join(workRoot, meta.TargetFile), nil
		}
		return filepath.Join(workRoot, filepath.Base(sample.Path)), nil
	}
	dst := filepath.Join(workRoot, filepath.Base(sample.Path))
	if err := copyFile(sample.Path, dst); err != nil {
		return "", err
	}
	return dst, nil
}

func injectSkillWorkspace(workRoot string, skill contracts.SkillSpec) (string, error) {
	if skill.Name == "" || skill.Name == agentconfig.NoSkill {
		return "", nil
	}
	skillDir := filepath.Join(workRoot, ".utbench", "skills", safePathName(skill.Name))
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return "", err
	}
	if strings.EqualFold(skill.InjectMode, "workspace_mount") || len(skill.Files) > 0 {
		for _, src := range skill.Files {
			if strings.TrimSpace(src) == "" {
				continue
			}
			info, err := os.Stat(src)
			if err != nil {
				return "", err
			}
			dst := filepath.Join(skillDir, filepath.Base(src))
			if info.IsDir() {
				if err := copyDir(src, dst); err != nil {
					return "", err
				}
			} else if err := copyFile(src, dst); err != nil {
				return "", err
			}
		}
	}
	return skillDir, nil
}

func appendSkillInstruction(prompt string, skill contracts.SkillSpec) string {
	if skill.Name == "" || skill.Name == agentconfig.NoSkill || !strings.EqualFold(defaultString(skill.InjectMode, "prompt_append"), "prompt_append") {
		return prompt
	}
	var b strings.Builder
	b.WriteString(prompt)
	b.WriteString("\n\nAdditional skill package: ")
	b.WriteString(skill.Name)
	if skill.Version != "" {
		b.WriteString(" v")
		b.WriteString(skill.Version)
	}
	b.WriteString("\n")
	if skill.Description != "" {
		b.WriteString(skill.Description)
		b.WriteString("\n")
	}
	if skill.InstructionPath != "" {
		if raw, err := os.ReadFile(skill.InstructionPath); err == nil && len(raw) > 0 {
			b.WriteString("\nSkill instructions:\n")
			b.Write(raw)
			b.WriteString("\n")
		}
	}
	return b.String()
}

func buildAgentPrompt(prompt string, sample contracts.SampleRef, sourceFile, outputFile, skillDir string) string {
	var b strings.Builder
	b.WriteString(prompt)
	b.WriteString("\n\nAgent execution contract:\n")
	b.WriteString("- Work only inside the provided workspace.\n")
	b.WriteString("- Do not modify the original source behavior.\n")
	b.WriteString("- Generate one complete unit test file.\n")
	b.WriteString("- Write the final test file to: ")
	b.WriteString(outputFile)
	b.WriteString("\n")
	b.WriteString("- Target language: ")
	b.WriteString(sample.Language)
	b.WriteString("\n- Source file in workspace: ")
	b.WriteString(sourceFile)
	b.WriteString("\n")
	if skillDir != "" {
		b.WriteString("- Skill files are available at: ")
		b.WriteString(skillDir)
		b.WriteString("\n")
	}
	return b.String()
}

func renderTemplateText(name, content string, data commandTemplateData) (string, error) {
	tpl, err := template.New(name).Parse(content)
	if err != nil {
		return "", err
	}
	var b strings.Builder
	if err := tpl.Execute(&b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}

func buildAgentEnv(fw agentconfig.FrameworkSpec, model modelConfig, data commandTemplateData) (map[string]string, []string, error) {
	out := make(map[string]string, len(fw.Env))
	keys := make([]string, 0, len(fw.Env))
	for key := range fw.Env {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value, err := renderTemplateText("agent-env-"+key, fw.Env[key], data)
		if err != nil {
			return nil, nil, fmt.Errorf("render env %s: %w", key, err)
		}
		out[key] = value
	}
	envFromHost := append([]string{}, fw.EnvFromHost...)
	if model.APIKeyEnv != "" {
		envFromHost = append(envFromHost, model.APIKeyEnv)
	}
	return out, uniqueSortedStrings(envFromHost), nil
}

func buildSandboxRunRequest(outputRoot string, fw agentconfig.FrameworkSpec, workspace, command string, env map[string]string, envFromHost []string) SandboxRunRequest {
	return SandboxRunRequest{
		Mode:                fw.SandboxMode,
		Workspace:           workspace,
		ContainerOutputRoot: outputRoot,
		Command:             command,
		Env:                 env,
		EnvFromHost:         envFromHost,
		DockerImage:         fw.DockerImage,
		NetworkDisabled:     fw.NetworkDisabled,
		CPU:                 fw.CPU,
		Memory:              fw.Memory,
		TimeoutSeconds:      fw.TimeoutSeconds,
	}
}

func snapshotWorkspace(root string) (map[string]string, error) {
	out := map[string]string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		sum := sha256.Sum256(raw)
		out[filepath.ToSlash(rel)] = hex.EncodeToString(sum[:])
		return nil
	})
	return out, err
}

func diffSnapshots(before, after map[string]string) []string {
	var out []string
	for path, hash := range after {
		if before[path] != hash {
			out = append(out, path)
		}
	}
	sort.Strings(out)
	return out
}

func findGeneratedTest(workRoot, preferred string, globs, changes []string, language string) string {
	if raw, err := os.ReadFile(preferred); err == nil && strings.TrimSpace(string(raw)) != "" {
		return preferred
	}
	for _, pattern := range globs {
		matches, _ := filepath.Glob(filepath.Join(workRoot, filepath.FromSlash(pattern)))
		sort.Strings(matches)
		for _, match := range matches {
			if isTestFile(match, language) {
				return match
			}
		}
	}
	for _, rel := range changes {
		path := filepath.Join(workRoot, filepath.FromSlash(rel))
		if isTestFile(path, language) {
			return path
		}
	}
	return ""
}

func isTestFile(path, language string) bool {
	base := strings.ToLower(filepath.Base(path))
	switch strings.ToLower(language) {
	case "python":
		return strings.HasSuffix(base, ".py") && (strings.HasPrefix(base, "test_") || strings.Contains(base, "_test"))
	case "go":
		return strings.HasSuffix(base, "_test.go")
	case "java":
		return strings.HasSuffix(base, "test.java") || strings.Contains(base, "test")
	case "cpp":
		return strings.HasSuffix(base, ".cpp") && strings.Contains(base, "test")
	default:
		return strings.HasPrefix(base, "generated_test")
	}
}

func writeTrace(path string, row map[string]any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.Marshal(row)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(raw); err != nil {
		return err
	}
	_, err = f.WriteString("\n")
	return err
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func safePathName(raw string) string {
	return sanitizeIdentifier(raw)
}

func uniqueSortedStrings(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func mustRel(base, target string) string {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return filepath.Base(target)
	}
	return rel
}
