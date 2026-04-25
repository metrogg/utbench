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

type PromptMode string

const (
	PromptModeFullFile    PromptMode = "full_file"
	PromptModeCompletion  PromptMode = "completion"
	PromptModeModuleLevel PromptMode = "module_level"
)

const systemMessage = "You are a senior unit test generation model. " +
	"Return only runnable test code. Do not include explanations, markdown fences, or placeholder tests. " +
	"Use only the provided source and context. Do not invent APIs, imports, or behavior."

const promptStrategy = "structured-v1"

var promptLanguages = []string{"python", "go", "java", "cpp"}

type PromptCatalog struct {
	Strategy      string                           `json:"strategy"`
	VersionID     string                           `json:"version_id"`
	SystemMessage string                           `json:"system_message"`
	Modes         []PromptMode                     `json:"modes"`
	Templates     map[string]map[PromptMode]string `json:"templates"`
}

type PromptRequest struct {
	Mode            PromptMode
	Language        string
	SamplePath      string
	SourceCode      string
	ExistingTestSrc string
	ModuleMeta      *moduleLevelMetaForRunner
}

func buildPrompt(language, samplePath, sourceCode string) string {
	req := PromptRequest{
		Mode:       PromptModeFullFile,
		Language:   language,
		SamplePath: samplePath,
		SourceCode: sourceCode,
	}
	if meta := loadModuleLevelMetaForRunner(samplePath); meta != nil {
		req.Mode = PromptModeModuleLevel
		req.ModuleMeta = meta
	}
	return BuildPrompt(req)
}

func BuildPrompt(req PromptRequest) string {
	switch req.Mode {
	case PromptModeCompletion:
		return buildCompletionPrompt(req)
	case PromptModeModuleLevel:
		if req.ModuleMeta != nil {
			return buildModuleLevelPrompt(req)
		}
		return buildFullFilePrompt(req)
	default:
		return buildFullFilePrompt(req)
	}
}

func PromptStrategy() string {
	return promptStrategy
}

func PromptVersionID() string {
	return BuildPromptCatalog().VersionID
}

func BuildPromptCatalog() PromptCatalog {
	templates := make(map[string]map[PromptMode]string, len(promptLanguages))
	modes := []PromptMode{PromptModeFullFile, PromptModeCompletion, PromptModeModuleLevel}
	for _, language := range promptLanguages {
		templates[language] = map[PromptMode]string{
			PromptModeFullFile:    buildPromptPreview(previewPromptRequest(language, PromptModeFullFile)),
			PromptModeCompletion:  buildPromptPreview(previewPromptRequest(language, PromptModeCompletion)),
			PromptModeModuleLevel: buildPromptPreview(previewPromptRequest(language, PromptModeModuleLevel)),
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

func PromptTemplatePreview(language string) string {
	catalog := BuildPromptCatalog()
	lang := normalizeLanguage(language)
	if templates, ok := catalog.Templates[lang]; ok {
		return templates[PromptModeFullFile]
	}
	return catalog.Templates["python"][PromptModeFullFile]
}

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
	moduleSymbols := extractModuleLevelSymbols(req.SourceCode, lang)
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

func buildModuleLevelPrompt(req PromptRequest) string {
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
	fmt.Fprintf(&b, "Mode: %s\n\n", PromptModeModuleLevel)

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
			"Use JUnit 4 with `@Test` and standard `Assert` methods.",
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

func normalizeLanguage(language string) string {
	return strings.ToLower(strings.TrimSpace(language))
}

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
	case PromptModeModuleLevel:
		req.ModuleMeta = &moduleLevelMetaForRunner{
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
