package runner

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"go-ut-bench/internal/agentconfig"
	"go-ut-bench/internal/contracts"
)

// subjectTrace 保留用于向后兼容，内部逻辑已迁移到 AgentTrace。
// 新代码应直接使用 AgentTrace。
type subjectTrace struct {
	TracePath          string
	WorkspaceDiffPath  string
	SandboxProvider    string
	SandboxFingerprint string
	TokenSource        string
	EstimatedCostUSD   *float64
	CostSource         string
}

// agentTraceSummary 封装从 AgentTrace 中提取的摘要信息，用于运行时显示。
type agentTraceSummary struct {
	InteractionCount int
	ToolCallCount    int
	FilesRead        int
	FilesWritten     int
	CommandsExecuted int
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
) (string, map[string]any, subjectTrace, int, *int, *int, *int, bool, *contracts.ErrorInfo, agentTraceSummary) {
	prompt = appendSkillInstruction(prompt, target.subject.Skill)

	// 选择 adapter
	var adapter AgentAdapter
	switch target.subject.Spec.Kind {
	case agentconfig.KindModelAPI, "":
		adapter = newModelAPIAdapter()
	case agentconfig.KindCLIAgent:
		adapter = newCLIAgentAdapter(s.sandboxRunner)
	default:
		return "", nil, subjectTrace{}, 0, nil, nil, nil, false, &contracts.ErrorInfo{
			Kind:      "unsupported_subject_kind",
			Message:   fmt.Sprintf("unsupported subject kind: %s", target.subject.Spec.Kind),
			Retryable: false,
		}, agentTraceSummary{}
	}

	// 记录 API 请求日志（model_api 场景）
	if target.subject.Spec.Kind == agentconfig.KindModelAPI || target.subject.Spec.Kind == "" {
		s.logger.LogAPIRequest(target.subject.Spec.ID, sample.Language, sample.ID, 0, 0)
	}

	// 调用 adapter
	result := adapter.Generate(ctx, AgentGenerateRequest{
		Subject:    target.subject,
		Model:      target.model,
		Sample:     sample,
		Prompt:     prompt,
		TestPath:   testPath,
		MetaRoot:   metaRoot,
		OutputRoot: spec.OutputRoot,
		RunID:      spec.RunID,
	})

	// 记录 API 响应日志
	s.logger.LogAPIResponse(target.subject.Spec.ID, sample.Language, sample.ID, result.Error == nil, result.Truncated, errorMsgSafe(result.Error))
	s.logger.ToFile("runner").Trace("generate_response",
		"subject_id", target.subject.Spec.ID,
		"framework", target.subject.Spec.Framework,
		"model", target.model.Name,
		"language", sample.Language,
		"sample_id", sample.ID,
		"prompt_tokens", result.PromptTokens,
		"completion_tokens", result.CompletionTokens,
		"total_tokens", result.TotalTokens,
		"latency_ms", result.LatencyMS,
		"truncated", result.Truncated,
		"interaction_count", result.Trace.InteractionCount,
		"tool_call_count", len(result.Trace.ToolCalls),
		"files_read", len(result.Trace.FilesRead),
		"files_written", len(result.Trace.FilesWritten),
		"commands_executed", len(result.Trace.CommandsExecuted),
		"success", result.Error == nil,
	)

	// 转换为旧的 subjectTrace 格式（向后兼容）
	trace := subjectTrace{
		TracePath:          result.Trace.TracePath,
		WorkspaceDiffPath:  result.Trace.WorkspaceDiffPath,
		SandboxProvider:    result.Trace.SandboxProvider,
		SandboxFingerprint: result.Trace.SandboxFingerprint,
		TokenSource:        result.TokenSource,
		EstimatedCostUSD:   result.EstimatedCostUSD,
		CostSource:         result.CostSource,
	}

	summary := agentTraceSummary{
		InteractionCount: result.Trace.InteractionCount,
		ToolCallCount:    len(result.Trace.ToolCalls),
		FilesRead:        len(result.Trace.FilesRead),
		FilesWritten:     len(result.Trace.FilesWritten),
		CommandsExecuted: len(result.Trace.CommandsExecuted),
	}

	return result.Code, result.RawResponse, trace, result.LatencyMS,
		result.PromptTokens, result.CompletionTokens, result.TotalTokens,
		result.Truncated, result.Error, summary
}

// generateWithModelAPI 已迁移到 adapter_model.go 中的 modelAPIAdapter。
// 保留此函数签名用于向后兼容测试。

// generateWithCLIAgent 已迁移到 adapter_cli.go 中的 cliAgentAdapter。
// 保留 helper 函数供 adapter 使用。

func prepareAgentWorkspace(workRoot string, sample contracts.SampleRef) (string, error) {
	if meta := loadRepoLevelMetaForRunner(sample.Path); meta != nil && meta.WorkspaceRoot != "" {
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
	// 将源文件设为只读，防止 agent 意外截断或修改源代码
	_ = os.Chmod(dst, 0444)
	return dst, nil
}

func buildSampleEnvironmentSetupCommands(sample contracts.SampleRef, workRoot string) []string {
	var commands []string
	switch strings.ToLower(strings.TrimSpace(sample.Language)) {
	case "python":
		if meta := loadRepoLevelMetaForRunner(sample.Path); meta != nil && len(meta.Requirements) > 0 {
			requirements := shellJoinArgs(meta.Requirements)
			if requirements != "" {
				commands = append(commands, "python3 -m pip install --disable-pip-version-check "+requirements)
			}
		}
		for _, rel := range []string{"requirements.txt", "requirements-dev.txt"} {
			if fileExists(filepath.Join(workRoot, rel)) {
				commands = append(commands, "python3 -m pip install --disable-pip-version-check -r "+rel)
			}
		}
	case "go":
		if fileExists(filepath.Join(workRoot, "go.mod")) {
			commands = append(commands, "go mod download")
		}
	case "java":
		if fileExists(filepath.Join(workRoot, "pom.xml")) {
			commands = append(commands, "mvn -q -DskipTests dependency:go-offline")
		}
	}
	return uniqueSortedStrings(commands)
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

func buildSandboxRunRequest(outputRoot string, fw agentconfig.FrameworkSpec, language, workspace, command string, env map[string]string, envFromHost []string) SandboxRunRequest {
	mode := strings.TrimSpace(fw.Sandbox.Mode)
	if mode == "" {
		mode = fw.SandboxMode
	}
	// 如果配置要求 docker 模式但当前环境没有 Docker daemon 可用，自动降级为 local。
	// 典型场景：Web 在 Windows 上运行，启动外层容器执行 utbench run，
	// 但外层容器没有挂载 docker.sock。
	if strings.EqualFold(mode, "docker") && !isDockerAvailable() {
		mode = "local"
	}
	provider := strings.TrimSpace(fw.Sandbox.Provider)
	if provider == "" {
		if strings.EqualFold(mode, "docker") {
			provider = "docker"
		} else {
			provider = "local"
		}
	}
	return SandboxRunRequest{
		Provider:            provider,
		Mode:                mode,
		Workspace:           workspace,
		ContainerOutputRoot: outputRoot,
		Command:             command,
		Env:                 env,
		EnvFromHost:         envFromHost,
		DockerImage:         frameworkSandboxImage(fw, language),
		NetworkDisabled:     frameworkSandboxNetworkDisabled(fw),
		CPU:                 frameworkSandboxCPU(fw),
		Memory:              frameworkSandboxMemory(fw),
		TimeoutSeconds:      frameworkSandboxTimeout(fw),
	}
}

func frameworkSandboxImage(fw agentconfig.FrameworkSpec, language string) string {
	// 优先使用统一镜像（sandbox.image），不再按语言拆分。
	// 统一镜像包含所有语言运行时，适合仓库级多语言样本。
	if image := strings.TrimSpace(fw.Sandbox.Image); image != "" {
		return image
	}
	if image := strings.TrimSpace(fw.DockerImage); image != "" {
		return image
	}
	// 向后兼容：按语言查找 sandbox.images.{lang}
	language = normalizeFrameworkLookupKey(language)
	images := fw.Sandbox.Images
	if len(images) == 0 {
		images = fw.DockerImages
	}
	if images != nil {
		if image := strings.TrimSpace(images[language]); image != "" {
			return image
		}
		if image := strings.TrimSpace(images["default"]); image != "" {
			return image
		}
		if image := strings.TrimSpace(images["*"]); image != "" {
			return image
		}
	}
	return ""
}

func frameworkDockerImage(fw agentconfig.FrameworkSpec, language string) string {
	return frameworkSandboxImage(fw, language)
}

func frameworkSandboxTimeout(fw agentconfig.FrameworkSpec) int {
	if fw.Sandbox.TimeoutSeconds > 0 {
		return fw.Sandbox.TimeoutSeconds
	}
	return fw.TimeoutSeconds
}

func frameworkSandboxProvider(fw agentconfig.FrameworkSpec) string {
	if provider := strings.TrimSpace(fw.Sandbox.Provider); provider != "" {
		return provider
	}
	mode := strings.TrimSpace(fw.SandboxMode)
	if strings.EqualFold(mode, "docker") {
		return "docker"
	}
	return "local"
}

func frameworkSandboxNetworkDisabled(fw agentconfig.FrameworkSpec) bool {
	if fw.Sandbox.Mode != "" || fw.Sandbox.Provider != "" || fw.Sandbox.Image != "" || len(fw.Sandbox.Images) > 0 || fw.Sandbox.TimeoutSeconds > 0 || fw.Sandbox.CPU != "" || fw.Sandbox.Memory != "" || fw.Sandbox.NetworkDisabled != fw.NetworkDisabled {
		return fw.Sandbox.NetworkDisabled
	}
	return fw.NetworkDisabled
}

func frameworkSandboxCPU(fw agentconfig.FrameworkSpec) string {
	if strings.TrimSpace(fw.Sandbox.CPU) != "" {
		return strings.TrimSpace(fw.Sandbox.CPU)
	}
	return strings.TrimSpace(fw.CPU)
}

func frameworkSandboxMemory(fw agentconfig.FrameworkSpec) string {
	if strings.TrimSpace(fw.Sandbox.Memory) != "" {
		return strings.TrimSpace(fw.Sandbox.Memory)
	}
	return strings.TrimSpace(fw.Memory)
}

func frameworkPreflightCommands(fw agentconfig.FrameworkSpec, language string) []string {
	var out []string
	language = normalizeFrameworkLookupKey(language)
	if fw.Preflight != nil {
		for _, key := range []string{"default", "*", language} {
			for _, cmd := range fw.Preflight[key] {
				cmd = strings.TrimSpace(cmd)
				if cmd != "" {
					out = append(out, cmd)
				}
			}
		}
	}
	if len(out) > 0 {
		return uniqueSortedStrings(out)
	}
	switch language {
	case "python":
		return []string{"python3 --version", "pytest --version"}
	case "go":
		return []string{"go version"}
	case "java":
		return []string{"java -version", "mvn -version"}
	case "cpp", "c++", "cc":
		return []string{"g++ --version", "cmake --version"}
	default:
		return nil
	}
}

func frameworkForbiddenCommandPatterns(fw agentconfig.FrameworkSpec) []string {
	patterns := uniqueSortedStrings(fw.ForbiddenCommandPatterns)
	if len(patterns) > 0 {
		return patterns
	}
	return []string{
		"apt-get update",
		"apt-get install",
		"apt install",
		"apk add",
		"yum install",
		"dnf install",
		"zypper install",
		"pacman -s",
		"pip install",
		"pip3 install",
		"python -m pip install",
		"python3 -m pip install",
		"uv pip install",
		"poetry add",
		"npm install",
		"pnpm add",
		"yarn add",
		"go install ",
		"cargo install ",
	}
}

func normalizeFrameworkLookupKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func shellJoinArgs(values []string) string {
	var out []string
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if strings.ContainsAny(value, " \t\"'") {
			value = "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
		}
		out = append(out, value)
	}
	return strings.Join(out, " ")
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

// writeTrace 已迁移至 adapter_cli.go 中的 writeAgentTrace。
// 保留此函数签名供测试使用。

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
