// prompt 包提供提示词构建和管理功能
// 负责为不同语言生成 LLM 提示词、管理提示词版本和快照
// 支持三种提示词模式：完整文件模式、续写模式、仓库级模式
package runner

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// PromptMode 提示词模式类型
// 定义三种不同的提示词构建方式
type PromptMode string

// 三种提示词模式常量
const (
	PromptModeFullFile   PromptMode = "full_file"  // 完整文件模式：提供完整源码，用于 self_contained 类型样本
	PromptModeCompletion PromptMode = "completion" // 续写模式：基于已生成内容继续补全，用于截断后续写
	PromptModeRepoLevel  PromptMode = "repo_level" // 仓库级模式：包含 workspace 和 module_import 上下文，用于 repo_level 类型样本
)

// systemMessage 系统消息，强调只生成可运行的测试代码
const systemMessage = "You are a senior unit test generation model. " +
	"Return only runnable test code. Do not include explanations, markdown fences, or placeholder tests. " +
	"Use only the provided source and context. Do not invent APIs, imports, or behavior."

// promptStrategy 提示词策略名称
const promptStrategy = "structured-v1"

// promptLanguages 支持提示词的语言列表
var promptLanguages = []string{"python", "go", "java", "cpp"}

// PromptCatalog 提示词目录
// 包含版本信息、系统消息和各语言各模式的提示词模板
type PromptCatalog struct {
	Strategy      string                           `json:"strategy"`       // 提示词策略名称
	VersionID     string                           `json:"version_id"`     // 版本ID，由 SHA1 哈希生成
	SystemMessage string                           `json:"system_message"` // 系统消息
	Modes         []PromptMode                     `json:"modes"`          // 支持的提示词模式列表
	Templates     map[string]map[PromptMode]string `json:"templates"`      // 各语言各模式的模板内容
}

// PromptRequest 提示词构建请求参数
// 包含构建提示词所需的所有输入信息
type PromptRequest struct {
	Mode            PromptMode              // 提示词模式
	Language        string                  // 编程语言
	SamplePath      string                  // 样本文件路径
	SourceCode      string                  // 源代码内容
	ExistingTestSrc string                  // 已有测试代码（用于续写模式）
	ModuleMeta      *repoLevelMetaForRunner // 仓库级样本元数据（用于 repo_level 模式）
}

// buildPrompt 构建提示词（简化版本）
// 根据样本路径自动检测是否为 repo_level 类型，选择合适的模式
//
// 参数:
//   - language: 编程语言
//   - samplePath: 样本文件路径
//   - sourceCode: 源代码内容
//
// 返回值:
//   - string: 构建完成的提示词
func buildPrompt(language, samplePath, sourceCode string) string {
	req := PromptRequest{
		Mode:       PromptModeFullFile,
		Language:   language,
		SamplePath: samplePath,
		SourceCode: sourceCode,
	}
	if meta := loadRepoLevelMetaForRunner(samplePath); meta != nil {
		req.Mode = PromptModeRepoLevel
		req.ModuleMeta = meta
	}
	return BuildPrompt(req)
}

// BuildPrompt 根据请求参数构建提示词
// 根据模式选择对应的构建函数
//
// 参数:
//   - req: 提示词构建请求参数
//
// 返回值:
//   - string: 构建完成的提示词
func BuildPrompt(req PromptRequest) string {
	switch req.Mode {
	case PromptModeCompletion:
		return buildCompletionPrompt(req)
	case PromptModeRepoLevel:
		if req.ModuleMeta != nil {
			return buildRepoLevelPrompt(req)
		}
		return buildFullFilePrompt(req)
	default:
		return buildFullFilePrompt(req)
	}
}

// PromptStrategy 返回提示词策略名称
func PromptStrategy() string {
	return promptStrategy
}

// PromptVersionID 返回提示词版本ID
// 通过 BuildPromptCatalog 计算哈希生成
func PromptVersionID() string {
	return BuildPromptCatalog().VersionID
}

// BuildPromptCatalog 构建提示词目录
// 包含所有语言所有模式的提示词模板，以及版本信息
// 版本ID 由整个目录内容的 SHA1 哈希生成，确保可追溯
func BuildPromptCatalog() PromptCatalog {
	templates := make(map[string]map[PromptMode]string, len(promptLanguages))
	modes := []PromptMode{PromptModeFullFile, PromptModeCompletion, PromptModeRepoLevel}
	for _, language := range promptLanguages {
		templates[language] = map[PromptMode]string{
			PromptModeFullFile:   buildPromptPreview(previewPromptRequest(language, PromptModeFullFile)),
			PromptModeCompletion: buildPromptPreview(previewPromptRequest(language, PromptModeCompletion)),
			PromptModeRepoLevel:  buildPromptPreview(previewPromptRequest(language, PromptModeRepoLevel)),
		}
	}

	payload := struct {
		Strategy      string                           `json:"strategy"`
		SystemMessage string                           `json:"system_message"`
		Modes         []PromptMode                     `json:"modes"`
		Templates     map[string]map[PromptMode]string `json:"templates"`
	}{
		Strategy:      promptStrategy,
		SystemMessage: systemMessage,
		Modes:         modes,
		Templates:     templates,
	}
	raw, _ := json.Marshal(payload)
	sum := sha1.Sum(raw)

	return PromptCatalog{
		Strategy:      promptStrategy,
		VersionID:     hex.EncodeToString(sum[:])[:12],
		SystemMessage: systemMessage,
		Modes:         modes,
		Templates:     templates,
	}
}

// WritePromptCatalog 将提示词目录写入目录
// 生成 system.txt、各语言各模式的 .prompt.txt 文件、prompt_catalog.json
//
// 参数:
//   - dir: 目标目录路径
//
// 返回值:
//   - PromptCatalog: 写入的目录内容
//   - error: 写入过程中的错误
func WritePromptCatalog(dir string) (PromptCatalog, error) {
	catalog := BuildPromptCatalog()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return PromptCatalog{}, err
	}
	if err := os.WriteFile(filepath.Join(dir, "system.txt"), []byte(catalog.SystemMessage), 0o644); err != nil {
		return PromptCatalog{}, err
	}
	for _, language := range promptLanguages {
		modeTemplates := catalog.Templates[language]
		for _, mode := range catalog.Modes {
			filename := fmt.Sprintf("%s_%s.prompt.txt", language, string(mode))
			if err := os.WriteFile(filepath.Join(dir, filename), []byte(modeTemplates[mode]), 0o644); err != nil {
				return PromptCatalog{}, err
			}
		}
	}
	raw, err := json.MarshalIndent(catalog, "", "  ")
	if err != nil {
		return PromptCatalog{}, err
	}
	if err := os.WriteFile(filepath.Join(dir, "prompt_catalog.json"), raw, 0o644); err != nil {
		return PromptCatalog{}, err
	}
	return catalog, nil
}

// LoadPromptCatalog 从目录加载提示词目录
// 读取 prompt_catalog.json 文件
//
// 参数:
//   - dir: 目录路径
//
// 返回值:
//   - PromptCatalog: 加载的目录内容
//   - error: 加载过程中的错误
func LoadPromptCatalog(dir string) (PromptCatalog, error) {
	var catalog PromptCatalog
	raw, err := os.ReadFile(filepath.Join(dir, "prompt_catalog.json"))
	if err != nil {
		return PromptCatalog{}, err
	}
	if err := json.Unmarshal(raw, &catalog); err != nil {
		return PromptCatalog{}, err
	}
	return catalog, nil
}

// PromptTemplatePreview 返回指定语言的完整文件模式提示词模板预览
// 用于查看提示词内容
func PromptTemplatePreview(language string) string {
	catalog := BuildPromptCatalog()
	lang := normalizeLanguage(language)
	if templates, ok := catalog.Templates[lang]; ok {
		return templates[PromptModeFullFile]
	}
	return catalog.Templates["python"][PromptModeFullFile]
}

// buildFullFilePrompt 构建完整文件模式提示词
// 用于 self_contained 类型样本，提供完整源码和测试要求
//
// 参数:
//   - req: 提示词构建请求参数
//
// 返回值:
//   - string: 构建完成的提示词
func buildFullFilePrompt(req PromptRequest) string {
	lang := normalizeLanguage(req.Language)
	sampleID := strings.TrimSuffix(filepath.Base(req.SamplePath), filepath.Ext(req.SamplePath))
	framework := languageFramework(lang)
	moduleName := moduleImportName(sampleID)
	if lang == "python" {
		moduleName = "module_under_test"
	}
	dependencies := formatDependencyText(extractDependencies(req.SourceCode, lang))
	criticalConditions := formatBulletList(extractCriticalConditions(req.SourceCode, lang), "none detected")
	moduleSymbols := extractRepoLevelSymbols(req.SourceCode, lang)
	mockReq := mockRequirement(req.SourceCode)

	var b strings.Builder
	b.WriteString("Task: Generate one complete test file for the provided source code.\n")
	fmt.Fprintf(&b, "Language: %s\n", lang)
	fmt.Fprintf(&b, "Framework: %s\n", framework)
	fmt.Fprintf(&b, "Mode: %s\n\n", PromptModeFullFile)

	b.WriteString("Hard requirements:\n")
	b.WriteString("- Return one complete runnable test file.\n")
	b.WriteString("- Tests must be deterministic.\n")
	b.WriteString("- Cover normal, boundary, and error paths when they exist in the source.\n")
	b.WriteString("- Derive assertions from implementation behavior, not comments or common sense.\n")
	b.WriteString("- Use only symbols that appear in the source or explicit context.\n")
	b.WriteString("- Avoid unrelated third-party packages unless the source already depends on them.\n")
	b.WriteString("- Prefer meaningful assertions over placeholder tests.\n")
	b.WriteString("- Keep test inputs small and representative; do not create stress tests or huge inputs.\n")
	b.WriteString("- Never access real networks, real credentials, or real external services.\n")
	b.WriteString("- If external I/O exists, use mocks, stubs, or fakes instead of real services.\n")
	b.WriteString("- If the source exposes injection parameters for dependencies, prefer those fakes over patching globals.\n")
	for _, rule := range promptLanguageRules(lang, moduleName) {
		b.WriteString("- ")
		b.WriteString(rule)
		b.WriteString("\n")
	}

	b.WriteString("\nContext:\n")
	fmt.Fprintf(&b, "- Sample ID: %s\n", sampleID)
	fmt.Fprintf(&b, "- Detected dependencies: %s\n", dependencies)
	for _, line := range promptSourceContext(lang, req.SourceCode) {
		fmt.Fprintf(&b, "- %s\n", line)
	}
	if moduleSymbols != "" {
		fmt.Fprintf(&b, "- %s\n", moduleSymbols)
	}
	fmt.Fprintf(&b, "- Mock guidance: %s\n", mockReq)
	fmt.Fprintf(&b, "- Coverage target (reference only): %s\n", coverageTargetsText())
	fmt.Fprintf(&b, "- Critical conditions:\n%s\n", criticalConditions)

	b.WriteString("\nOutput contract:\n")
	b.WriteString("- Output raw code only.\n")
	b.WriteString("- No markdown fences.\n")
	b.WriteString("- No explanations.\n\n")

	fmt.Fprintf(&b, "Source code:\n```%s\n%s\n```", lang, req.SourceCode)
	return b.String()
}

// buildCompletionPrompt 构建续写模式提示词
// 用于截断后续写，基于已生成内容继续补全
//
// 参数:
//   - req: 提示词构建请求参数，需包含 ExistingTestSrc
//
// 返回值:
//   - string: 构建完成的提示词
func buildCompletionPrompt(req PromptRequest) string {
	lang := normalizeLanguage(req.Language)
	framework := languageFramework(lang)
	existing := strings.TrimSpace(req.ExistingTestSrc)
	if existing == "" {
		existing = "// no existing tests"
	}

	var b strings.Builder
	b.WriteString("Task: Continue an existing test file by adding the next useful test.\n")
	fmt.Fprintf(&b, "Language: %s\n", lang)
	fmt.Fprintf(&b, "Framework: %s\n", framework)
	fmt.Fprintf(&b, "Mode: %s\n\n", PromptModeCompletion)

	b.WriteString("Hard requirements:\n")
	b.WriteString("- Output only the next test function or test block.\n")
	b.WriteString("- Do not repeat imports, helpers, or existing tests.\n")
	b.WriteString("- Choose a path not yet covered by the existing tests.\n")
	b.WriteString("- Keep the existing style and indentation.\n")
	for _, rule := range promptLanguageRules(lang, moduleImportName(filepath.Base(req.SamplePath))) {
		b.WriteString("- ")
		b.WriteString(rule)
		b.WriteString("\n")
	}

	fmt.Fprintf(&b, "\nSource code:\n```%s\n%s\n```\n\n", lang, req.SourceCode)
	fmt.Fprintf(&b, "Existing test file:\n```%s\n%s\n```", lang, existing)
	return b.String()
}

// buildRepoLevelPrompt 构建仓库级模式提示词
// 用于 repo_level 类型样本，包含 workspace 和 module_import 上下文
//
// 参数:
//   - req: 提示词构建请求参数，需包含 ModuleMeta
//
// 返回值:
//   - string: 构建完成的提示词
func buildRepoLevelPrompt(req PromptRequest) string {
	lang := normalizeLanguage(req.Language)
	framework := languageFramework(lang)
	meta := req.ModuleMeta
	sampleID := meta.SampleID
	if sampleID == "" {
		sampleID = strings.TrimSuffix(filepath.Base(req.SamplePath), filepath.Ext(req.SamplePath))
	}

	packageName := meta.PackageName
	if packageName == "" {
		packageName = meta.ModuleImport
	}

	requirements := "none specified"
	if len(meta.Requirements) > 0 {
		requirements = strings.Join(meta.Requirements, ", ")
	}

	var b strings.Builder
	b.WriteString("Task: Generate one complete test file for the target module in this multi-file package.\n")
	fmt.Fprintf(&b, "Language: %s\n", lang)
	fmt.Fprintf(&b, "Framework: %s\n", framework)
	fmt.Fprintf(&b, "Mode: %s\n\n", PromptModeRepoLevel)

	b.WriteString("Hard requirements:\n")
	b.WriteString("- Return one complete runnable test file.\n")
	b.WriteString("- Tests must be deterministic.\n")
	b.WriteString("- Cover normal, boundary, and error paths when they exist in the target module.\n")
	b.WriteString("- Derive assertions from implementation behavior, not comments or common sense.\n")
	b.WriteString("- Use only symbols that appear in the provided package context.\n")
	b.WriteString("- Do not import from relative helper paths unless they are explicitly shown in the context.\n")
	b.WriteString("- Keep test inputs small and representative; do not create stress tests or huge inputs.\n")
	b.WriteString("- Never access real networks, real credentials, or real external services.\n")
	b.WriteString("- If external I/O exists, use mocks, stubs, or fakes instead of real services.\n")
	b.WriteString("- If the source exposes injection parameters for dependencies, prefer those fakes over patching globals.\n")
	for _, rule := range promptLanguageRules(lang, meta.ModuleImport) {
		b.WriteString("- ")
		b.WriteString(rule)
		b.WriteString("\n")
	}

	b.WriteString("\nModule context:\n")
	fmt.Fprintf(&b, "- Sample ID: %s\n", sampleID)
	fmt.Fprintf(&b, "- Target module: %s\n", meta.ModuleImport)
	fmt.Fprintf(&b, "- Package name: %s\n", packageName)
	fmt.Fprintf(&b, "- Target file: %s\n", meta.TargetFile)
	fmt.Fprintf(&b, "- Declared requirements: %s\n", requirements)
	fmt.Fprintf(&b, "- Coverage target (reference only): %s\n", coverageTargetsText())

	b.WriteString("\nOutput contract:\n")
	b.WriteString("- Output raw code only.\n")
	b.WriteString("- No markdown fences.\n")
	b.WriteString("- No explanations.\n\n")

	fmt.Fprintf(&b, "Package entry file:\n```%s\n%s\n```", lang, req.SourceCode)
	return b.String()
}

// promptLanguageRules 返回各语言特有的测试规则
// 为不同语言提供特定的测试框架、命名规范、导入规则等要求
//
// 参数:
//   - lang: 编程语言
//   - moduleName: 模块名称（用于导入语句）
//
// 返回值:
//   - []string: 语言特定的规则列表
func promptLanguageRules(lang, moduleName string) []string {
	switch lang {
	case "python":
		return []string{
			fmt.Sprintf("Use pytest function-based tests and import target symbols from `%s`.", moduleName),
			"Name every test function with the `test_` prefix.",
			"Use plain `assert` and `pytest.raises` for failure paths.",
			"If the test code uses a module such as `csv`, `json`, `os`, `tempfile`, or `xml.etree.ElementTree`, import it explicitly in the test file even if the source imports it.",
			"Do not introduce third-party modules that are not already present in the source context.",
			"When mocking, patch the symbol as imported by the target module and avoid building a mock spec from another mock object.",
		}
	case "go":
		return []string{
			"Use `func TestXxx(t *testing.T)` and import `testing`.",
			"Keep the test package consistent with the source package.",
			"Prefer table-driven tests when they make the cases clearer.",
			"Import every package referenced anywhere in the test file, including helper types and helper functions.",
			"Do not assume replaceable globals or identifiers exist inside standard library packages unless they are visible in the provided source context.",
			"Do not try to monkey-patch standard library internals to force error paths when the source code exposes no seam for doing so.",
		}
	case "java":
		return []string{
			"Use JUnit 5 (Jupiter) with `@Test` from `org.junit.jupiter.api.Test`.",
			"Use `Assertions.*` methods from `org.junit.jupiter.api.Assertions`.",
			"Test methods should be `void` (no need for `public` modifier in JUnit 5).",
			"Name the test class `ClassNameTest` and keep package declarations consistent with source.",
			"Do not access private members directly unless the source makes that the intended API surface.",
			"If the source file has no `package` declaration, the test file must also have no `package` declaration.",
			"Instantiate the exact class declared in the source before calling instance methods; only call methods statically when the source declares them as `static`.",
			"Use the exact class names and method names shown in the source; do not derive package names or type names from the sample id or file path.",
		}
	case "cpp":
		return []string{
			"Use GoogleTest with `TEST`, `EXPECT_*`, and `ASSERT_*` macros.",
			"Include only the headers needed by the generated tests.",
			"Include every standard header required by constants or helpers used in the test code, such as `<climits>` for `INT_MAX` and `INT_MIN`.",
			"Do not invent extra helper headers or duplicate declarations for classes and functions that are already defined in the provided source.",
		}
	default:
		return []string{
			"Follow the standard unit testing conventions for the target language and framework.",
		}
	}
}

// normalizeLanguage 规范化语言名称
// 转为小写并去除空白
func normalizeLanguage(language string) string {
	return strings.ToLower(strings.TrimSpace(language))
}

// formatDependencyText 格式化依赖列表文本
// 空列表返回 "none detected"
func formatDependencyText(deps []string) string {
	if len(deps) == 0 {
		return "none detected"
	}
	return strings.Join(deps, ", ")
}

func formatBulletList(items []string, fallback string) string {
	if len(items) == 0 {
		return "- " + fallback
	}
	lines := make([]string, 0, len(items))
	for _, item := range items {
		lines = append(lines, "- "+item)
	}
	return strings.Join(lines, "\n")
}

func buildPromptPreview(req PromptRequest) string {
	var b strings.Builder
	b.WriteString("System\n")
	b.WriteString(systemMessage)
	b.WriteString("\n\nUser\n")
	b.WriteString(BuildPrompt(req))
	return b.String()
}

func previewPromptRequest(language string, mode PromptMode) PromptRequest {
	lang := normalizeLanguage(language)
	samplePath := previewSamplePath(lang)
	req := PromptRequest{
		Mode:       mode,
		Language:   lang,
		SamplePath: samplePath,
		SourceCode: previewSourceCode(lang),
	}
	switch mode {
	case PromptModeCompletion:
		req.ExistingTestSrc = previewExistingTest(lang)
	case PromptModeRepoLevel:
		req.ModuleMeta = &repoLevelMetaForRunner{
			SampleID:     "preview_sample",
			ModuleImport: previewModuleImport(lang),
			PackageName:  previewPackageName(lang),
			TargetFile:   filepath.Base(samplePath),
			Requirements: []string{"preview dependency"},
		}
	}
	return req
}

func previewSamplePath(language string) string {
	switch language {
	case "python":
		return "/tmp/preview_sample.py"
	case "go":
		return "/tmp/preview_sample.go"
	case "java":
		return "/tmp/PreviewSample.java"
	case "cpp":
		return "/tmp/preview_sample.cpp"
	default:
		return "/tmp/preview_sample.txt"
	}
}

func previewSourceCode(language string) string {
	switch language {
	case "python":
		return "import math\n\ndef preview(value):\n    if value <= 0:\n        return 0\n    return value + 1\n"
	case "go":
		return "package preview\n\nfunc Add(a int, b int) int {\n\treturn a + b\n}\n"
	case "java":
		return "public class PreviewSample {\n    public int add(int a, int b) {\n        return a + b;\n    }\n}\n"
	case "cpp":
		return "#include <string>\n\nint add(int a, int b) {\n    return a + b;\n}\n"
	default:
		return "preview source"
	}
}

func previewExistingTest(language string) string {
	switch language {
	case "python":
		return "def test_preview_positive():\n    assert preview(1) == 2\n"
	case "go":
		return "func TestAdd_Positive(t *testing.T) {\n\tif got := Add(1, 2); got != 3 {\n\t\tt.Fatalf(\"got %d\", got)\n\t}\n}\n"
	case "java":
		return "@Test\npublic void testAddPositive() {\n    assertEquals(3, new PreviewSample().add(1, 2));\n}\n"
	case "cpp":
		return "TEST(PreviewSample, Positive) {\n    EXPECT_EQ(add(1, 2), 3);\n}\n"
	default:
		return "existing test"
	}
}

func previewModuleImport(language string) string {
	switch language {
	case "python":
		return "preview_sample"
	case "go":
		return "preview"
	case "java":
		return "preview.PreviewSample"
	case "cpp":
		return "preview_sample"
	default:
		return "preview_sample"
	}
}

func previewPackageName(language string) string {
	switch language {
	case "go":
		return "preview"
	case "java":
		return "preview"
	default:
		return previewModuleImport(language)
	}
}

func sortedPromptLanguages() []string {
	languages := append([]string(nil), promptLanguages...)
	sort.Strings(languages)
	return languages
}

func promptSourceContext(lang, sourceCode string) []string {
	switch lang {
	case "java":
		return javaPromptContext(sourceCode)
	default:
		return nil
	}
}

func javaPromptContext(sourceCode string) []string {
	var lines []string
	pkg := extractJavaPackageName(sourceCode)
	if pkg == "" {
		lines = append(lines, "Declared package: default package (no package declaration in source).")
	} else {
		lines = append(lines, fmt.Sprintf("Declared package: %s.", pkg))
	}

	classes := extractJavaClassNames(sourceCode)
	if len(classes) > 0 {
		lines = append(lines, fmt.Sprintf("Declared classes: %s.", strings.Join(classes, ", ")))
	}

	staticMethods := extractJavaStaticMethodNames(sourceCode)
	if len(staticMethods) > 0 {
		lines = append(lines, fmt.Sprintf("Static methods declared in source: %s.", strings.Join(staticMethods, ", ")))
	}
	return lines
}

func extractJavaPackageName(sourceCode string) string {
	re := regexp.MustCompile(`(?m)^\s*package\s+([A-Za-z_][A-Za-z0-9_\.]*)\s*;`)
	match := re.FindStringSubmatch(sourceCode)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func extractJavaClassNames(sourceCode string) []string {
	re := regexp.MustCompile(`(?m)^\s*(?:public\s+)?(?:final\s+|abstract\s+)?class\s+([A-Za-z_][A-Za-z0-9_]*)\b`)
	matches := re.FindAllStringSubmatch(sourceCode, -1)
	seen := map[string]struct{}{}
	var out []string
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		name := strings.TrimSpace(m[1])
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}

func extractJavaStaticMethodNames(sourceCode string) []string {
	re := regexp.MustCompile(`(?m)^\s*(?:public|protected|private)?\s*static\s+[A-Za-z0-9_<>\[\], ?]+\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`)
	matches := re.FindAllStringSubmatch(sourceCode, -1)
	seen := map[string]struct{}{}
	var out []string
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		name := strings.TrimSpace(m[1])
		if name == "" {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		out = append(out, name)
	}
	return out
}
