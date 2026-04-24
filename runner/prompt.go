package runner

// Prompt building module for go-ut-bench.
//
// Design inspiration: facebookresearch/testgeneval (inference/configs/*).
// Key ideas borrowed:
//   - Separate prompt-building from transport (api.go stays HTTP-only).
//   - Dispatch multiple prompt "modes" (fullfile / completion / modulelevel).
//   - Keep output-format contract strict so downstream parsing is cheap.
//   - Inject structured context (dependencies, critical conditions, symbols)
//     extracted from the source to reduce hallucination.
//
// See docs/prompt-design.md for the full reference.

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// PromptMode selects which prompt template is rendered.
type PromptMode string

const (
	// PromptModeFullFile: generate a complete test file from scratch.
	PromptModeFullFile PromptMode = "fullfile"

	// PromptModeCompletion: continue an existing test file by producing the
	// next test function only (inspired by testgeneval's PROMPT_COMPLETION).
	PromptModeCompletion PromptMode = "completion"

	// PromptModeModuleLevel: multi-file package samples, driven by meta.json.
	PromptModeModuleLevel PromptMode = "modulelevel"
)

// systemMessage is used as the chat system prompt for all modes.
// It is intentionally terse: role + strict output contract.
const systemMessage = "You are an expert unit testing engineer. " +
	"Output only executable test code. No explanations, no markdown fences, no placeholder tests."

// PromptRequest carries all inputs required to render a prompt.
// Most callers go through buildPrompt() which infers fields from samplePath.
type PromptRequest struct {
	Mode       PromptMode
	Language   string
	SamplePath string
	SourceCode string

	// Optional: for PromptModeCompletion, the partial test file content
	// that the model should continue from.
	ExistingTestSrc string

	// Optional: module-level metadata (populated by loadModuleLevelMetaForRunner).
	ModuleMeta *moduleLevelMetaForRunner
}

// buildPrompt is the thin entry point kept for backwards compatibility with
// existing callers (api.go and api_test.go). It auto-detects module-level
// samples from a neighbouring meta.json and dispatches accordingly.
func buildPrompt(language, samplePath, sourceCode string) string {
	req := PromptRequest{
		Language:   language,
		SamplePath: samplePath,
		SourceCode: sourceCode,
	}
	if meta := loadModuleLevelMetaForRunner(samplePath); meta != nil {
		req.Mode = PromptModeModuleLevel
		req.ModuleMeta = meta
	} else {
		req.Mode = PromptModeFullFile
	}
	return BuildPrompt(req)
}

// BuildPrompt dispatches on req.Mode. It is the single public entry point.
func BuildPrompt(req PromptRequest) string {
	switch req.Mode {
	case PromptModeCompletion:
		return buildCompletionPrompt(req)
	case PromptModeModuleLevel:
		if req.ModuleMeta == nil {
			// Fallback to fullfile if caller forgot to load meta.
			return buildFullFilePrompt(req)
		}
		return buildModuleLevelPrompt(req.Language, req.SamplePath, req.SourceCode, req.ModuleMeta)
	default:
		return buildFullFilePrompt(req)
	}
}

// -----------------------------------------------------------------------------
// Mode 1: Full-file generation
// -----------------------------------------------------------------------------

func buildFullFilePrompt(req PromptRequest) string {
	lang := strings.ToLower(strings.TrimSpace(req.Language))
	sampleID := strings.TrimSuffix(filepath.Base(req.SamplePath), filepath.Ext(req.SamplePath))
	framework := languageFramework(lang)
	dependencies := extractDependencies(req.SourceCode, lang)
	scenario, complexity := parseSampleMeta(sampleID, req.SamplePath)
	moduleName := moduleImportName(sampleID)
	mockReq := mockRequirement(req.SourceCode)
	criticalConditions := extractCriticalConditions(req.SourceCode, lang)
	moduleSymbols := extractModuleLevelSymbols(req.SourceCode, lang)

	dependencyText := "none detected"
	if len(dependencies) > 0 {
		dependencyText = strings.Join(dependencies, ", ")
	}
	moduleSymbolsText := ""
	if moduleSymbols != "" {
		moduleSymbolsText = "\n" + moduleSymbols
	}
	if scenario == "" {
		scenario = "general"
	}
	if complexity == "" {
		complexity = "unknown"
	}

	criticalText := "- none detected"
	if len(criticalConditions) > 0 {
		lines := make([]string, 0, len(criticalConditions))
		for _, item := range criticalConditions {
			lines = append(lines, fmt.Sprintf("- `%s`", item))
		}
		criticalText = strings.Join(lines, "\n")
	}

	var b strings.Builder
	b.WriteString("You are an expert unit testing engineer.\n")
	b.WriteString("你是一名资深单元测试工程师。\n")
	b.WriteString("Generate high-quality unit tests based on the following specification.\n")
	b.WriteString("请基于以下规范生成高质量单元测试。\n\n")

	b.WriteString("## Role & Objective（角色与目标）\n")
	b.WriteString("- Goal: produce executable tests that match source behavior exactly.\n")
	b.WriteString("- 目标：生成可执行且与源码行为严格一致的测试。\n\n")

	b.WriteString("## Language & Framework（语言与框架）\n")
	fmt.Fprintf(&b, "- Language（语言）: %s\n", lang)
	fmt.Fprintf(&b, "- Test Framework（测试框架）: %s\n\n", framework)

	b.WriteString("## Step-by-Step Workflow（分步流程）\n")
	b.WriteString("1) Identify callable symbols and input/output contracts from source.\n")
	b.WriteString("2) Build a test matrix: normal, boundary, and exception paths.\n")
	b.WriteString("3) Derive expected values only from implementation semantics.\n")
	b.WriteString("4) Write deterministic, runnable tests with clear assertions.\n")
	b.WriteString("5) Self-check syntax/imports/assertions before final output.\n\n")

	b.WriteString("## Test Requirements（测试要求）\n")
	b.WriteString("- Cover normal paths, boundary conditions, and error/exception behavior.\n")
	b.WriteString("- 覆盖正常路径、边界条件和异常行为。\n")
	b.WriteString("- Keep tests deterministic and runnable.\n")
	b.WriteString("- 保持测试可重复、可执行（避免随机性）。\n")
	b.WriteString("- Use clear assertions with meaningful expected values.\n")
	b.WriteString("- 使用清晰断言和有意义的期望值。\n")
	b.WriteString(buildLanguageSpecificRules(lang, moduleName))
	b.WriteString("- Avoid importing unrelated third-party packages by default.\n")
	b.WriteString("- 默认不要引入无关第三方包。\n")
	fmt.Fprintf(&b, "- Coverage targets（覆盖率目标，供参考）: %s\n", coverageTargetsText())
	fmt.Fprintf(&b, "- Mock requirements（Mock 要求）: %s\n\n", mockReq)

	b.WriteString("## Semantic Alignment Hard Rules（语义对齐硬约束）\n")
	b.WriteString(buildSemanticAlignmentRules(lang))
	b.WriteString("\n")

	b.WriteString("## Critical Conditions Extracted（关键逻辑条件）\n")
	b.WriteString(criticalText)
	b.WriteString("\n\n")

	b.WriteString("## Error Prevention Checklist（错误预防清单，仅内部执行）\n")
	b.WriteString(buildErrorPreventionChecklist(lang))
	b.WriteString("\n")

	b.WriteString("## Context Information（上下文信息）\n")
	fmt.Fprintf(&b, "- Sample ID（样本ID）: %s\n", sampleID)
	fmt.Fprintf(&b, "- Scenario（场景）: %s\n", scenario)
	fmt.Fprintf(&b, "- Complexity（复杂度）: %s\n", complexity)
	fmt.Fprintf(&b, "- Dependencies detected（检测到依赖）: %s\n%s\n\n", dependencyText, moduleSymbolsText)

	b.WriteString("## Output Format（输出格式）\n")
	b.WriteString("- Return raw test code only (no Markdown fences).\n")
	b.WriteString("- 仅输出原始测试代码，不要 Markdown 代码块。\n")
	b.WriteString("- Do not include explanations.\n")
	b.WriteString("- 不要输出解释文字。\n\n")

	b.WriteString("## Source Code Under Test（被测源码）\n")
	fmt.Fprintf(&b, "```%s\n%s\n```", lang, req.SourceCode)
	return b.String()
}

// -----------------------------------------------------------------------------
// Mode 2: Completion (continue an existing test file)
// Inspired by testgeneval's PROMPT_COMPLETION.
// -----------------------------------------------------------------------------

func buildCompletionPrompt(req PromptRequest) string {
	lang := strings.ToLower(strings.TrimSpace(req.Language))
	framework := languageFramework(lang)
	existing := strings.TrimSpace(req.ExistingTestSrc)
	if existing == "" {
		existing = "// (no existing tests yet — generate the first one)"
	}

	var b strings.Builder
	b.WriteString("You are an expert unit testing engineer continuing an existing test file.\n")
	b.WriteString("你是一名资深单元测试工程师，需要在已有测试文件上追加下一条测试。\n\n")

	b.WriteString("## Objective（目标）\n")
	b.WriteString("- Write ONLY the next test function, designed to improve coverage.\n")
	b.WriteString("- 只编写下一个测试函数，目标是提升覆盖率。\n\n")

	b.WriteString("## Language & Framework（语言与框架）\n")
	fmt.Fprintf(&b, "- Language（语言）: %s\n", lang)
	fmt.Fprintf(&b, "- Test Framework（测试框架）: %s\n\n", framework)

	b.WriteString("## Rules（规则）\n")
	b.WriteString("- Do NOT restate imports, setup or tests that already exist in the context.\n")
	b.WriteString("- 不要重复已存在的 import / setup / 既有测试。\n")
	b.WriteString("- Preserve indentation and match the style of existing tests.\n")
	b.WriteString("- 保留缩进，风格与已有测试保持一致。\n")
	b.WriteString("- Choose a code path NOT yet covered by the existing tests.\n")
	b.WriteString("- 选择尚未覆盖的代码路径。\n")
	b.WriteString(buildLanguageSpecificRules(lang, moduleImportName(filepath.Base(req.SamplePath))))
	b.WriteString("\n")

	b.WriteString("## Source Code Under Test（被测源码）\n")
	fmt.Fprintf(&b, "```%s\n%s\n```\n\n", lang, req.SourceCode)

	b.WriteString("## Existing Test File（已有测试文件）\n")
	fmt.Fprintf(&b, "```%s\n%s\n```\n\n", lang, existing)

	b.WriteString("## Output Format（输出格式）\n")
	b.WriteString("- Output ONLY the next test function, no fences, no prose.\n")
	b.WriteString("- 仅输出下一条测试函数，不要代码围栏，不要解释。\n")
	return b.String()
}

// -----------------------------------------------------------------------------
// Mode 3: Module-level (multi-file package)
// -----------------------------------------------------------------------------

func buildModuleLevelPrompt(language, samplePath, sourceCode string, meta *moduleLevelMetaForRunner) string {
	lang := strings.ToLower(strings.TrimSpace(language))
	framework := languageFramework(lang)
	sampleID := meta.SampleID
	if sampleID == "" {
		sampleID = strings.TrimSuffix(filepath.Base(samplePath), filepath.Ext(samplePath))
	}

	moduleImport := meta.ModuleImport
	packageName := meta.PackageName
	if packageName == "" {
		packageName = moduleImport
	}
	targetFile := meta.TargetFile
	requirements := meta.Requirements

	requirementsText := "none specified"
	if len(requirements) > 0 {
		requirementsText = strings.Join(requirements, ", ")
	}

	var b strings.Builder
	b.WriteString("You are an expert unit testing engineer.\n")
	b.WriteString("你是一名资深单元测试工程师。\n")
	b.WriteString("Generate high-quality unit tests for a module-level (multi-file) package.\n")
	b.WriteString("请为一个模块级（多文件）包生成高质量单元测试。\n\n")

	b.WriteString("## Role & Objective（角色与目标）\n")
	b.WriteString("- Goal: produce executable tests that match source behavior exactly.\n")
	b.WriteString("- 目标：生成可执行且与源码行为严格一致的测试。\n")
	b.WriteString("- This is a MODULE-LEVEL sample: the target code is part of a larger package.\n")
	b.WriteString("- 这是模块级样本：被测代码是更大包的一部分。\n\n")

	b.WriteString("## Language & Framework（语言与框架）\n")
	fmt.Fprintf(&b, "- Language（语言）: %s\n", lang)
	fmt.Fprintf(&b, "- Test Framework（测试框架）: %s\n\n", framework)

	b.WriteString("## Module-Level Context（模块级上下文）\n")
	fmt.Fprintf(&b, "- Sample ID（样本ID）: %s\n", sampleID)
	fmt.Fprintf(&b, "- Target Module（目标模块）: `%s`\n", moduleImport)
	fmt.Fprintf(&b, "- Package Name（包名）: `%s`\n", packageName)
	fmt.Fprintf(&b, "- Target File（目标文件）: `%s`\n", targetFile)
	fmt.Fprintf(&b, "- Requirements（依赖）: %s\n\n", requirementsText)

	b.WriteString("## Test Requirements（测试要求）\n")
	fmt.Fprintf(&b, "- Import the target module using: `from %s import ...`\n", moduleImport)
	fmt.Fprintf(&b, "- 使用 `from %s import ...` 导入被测模块。\n", moduleImport)
	b.WriteString("- Write tests for the main classes/functions in the target module.\n")
	b.WriteString("- 为目标模块中的主要类/函数编写测试。\n")
	b.WriteString("- Cover normal paths, boundary conditions, and error/exception behavior.\n")
	b.WriteString("- 覆盖正常路径、边界条件和异常行为。\n")
	b.WriteString("- Keep tests deterministic and runnable.\n")
	b.WriteString("- 保持测试可重复、可执行（避免随机性）。\n")
	b.WriteString("- Use clear assertions with meaningful expected values.\n")
	b.WriteString("- 使用清晰断言和有意义的期望值。\n")
	b.WriteString("- Test function names must start with `test_`.\n")
	b.WriteString("- 测试函数命名必须以 `test_` 开头。\n")
	b.WriteString("- Use plain `assert` and `pytest.raises` for failure paths.\n")
	b.WriteString("- 断言使用 `assert`，异常路径使用 `pytest.raises`。\n")
	b.WriteString("- Use function-based tests (e.g., `def test_xxx()`) instead of class-based.\n")
	b.WriteString("- 使用函数式测试（如 `def test_xxx()`），不要用 class-based。\n")
	fmt.Fprintf(&b, "- Coverage targets（覆盖率目标，供参考）: %s\n\n", coverageTargetsText())

	b.WriteString("## Important Notes（重要提示）\n")
	b.WriteString("- DO NOT import from relative paths like `from ._common import ...`.\n")
	b.WriteString("- 不要从相对路径导入，如 `from ._common import ...`。\n")
	b.WriteString("- Always use the full module import path provided above.\n")
	b.WriteString("- 始终使用上面提供的完整模块导入路径。\n")
	b.WriteString("- The test file will be placed in a `tests/` subdirectory of the workspace.\n")
	b.WriteString("- 测试文件将放在 workspace 的 `tests/` 子目录中。\n\n")

	b.WriteString("## Output Format（输出格式）\n")
	b.WriteString("- Return raw test code only (no Markdown fences).\n")
	b.WriteString("- 仅输出原始测试代码，不要 Markdown 代码块。\n")
	b.WriteString("- Do not include explanations.\n")
	b.WriteString("- 不要输出解释文字。\n\n")

	b.WriteString("## Sample Entry File（样本入口文件）\n")
	b.WriteString("This file shows the import structure:\n")
	fmt.Fprintf(&b, "```%s\n%s\n```", lang, sourceCode)
	return b.String()
}

// -----------------------------------------------------------------------------
// Per-language rule adapters
// -----------------------------------------------------------------------------

func buildLanguageSpecificRules(lang, moduleName string) string {
	switch lang {
	case "python":
		return "- Test function names must start with `test_` (e.g., `def test_xxx()`).\n" +
			"- 测试函数命名必须以 `test_` 开头（如 `def test_xxx()`）。\n" +
			"- Use plain `assert` and `pytest.raises` for failure paths.\n" +
			"- 断言使用 `assert`，异常路径使用 `pytest.raises`。\n" +
			"- Use function-based tests (e.g., `def test_xxx()`) instead of class-based.\n" +
			"- 使用函数式测试（如 `def test_xxx()`），不要用 class-based。\n" +
			fmt.Sprintf("- MUST import target symbols from local module `%s` before writing tests.\n", moduleName) +
			fmt.Sprintf("- 必须先从同目录模块 `%s` 导入被测对象，再编写测试。\n", moduleName)

	case "go":
		return "- Test function names MUST start with `Test` (capital T) followed by a name (e.g., `func TestAdd(t *testing.T)`).\n" +
			"- 测试函数命名必须以 `Test`（大写T）开头（如 `func TestAdd(t *testing.T)`）。\n" +
			"- The test file MUST declare the EXACT same package as the source file. If source is `package main`, test MUST also be `package main`.\n" +
			"- 测试文件必须与源文件声明完全相同的 package。若源文件为 `package main`，测试文件也必须是 `package main`。\n" +
			"- Use `t.Errorf` for non-fatal failures and `t.Fatalf` when the test cannot continue. Do NOT use plain `assert`.\n" +
			"- 非致命失败用 `t.Errorf`，无法继续执行时用 `t.Fatalf`，不要用 `assert`。\n" +
			"- PREFER table-driven tests: define a `tests := []struct{ name string; input ...; want ... }` slice and iterate with `for _, tc := range tests { t.Run(tc.name, func(t *testing.T) { ... }) }`.\n" +
			"- 优先使用表格驱动测试：定义含 name/input/want 字段的结构体切片，用 `t.Run(tc.name, ...)` 遍历执行。\n" +
			"- Only import `\"testing\"` and standard library packages. Do NOT import third-party packages (e.g., testify, gomock) unless the source code already imports them.\n" +
			"- 只导入 `\"testing\"` 和标准库包。除非源码本身依赖，否则不要引入第三方包（如 testify、gomock）。\n" +
			"- Unexported functions (lowercase names) are accessible in tests since the test is in the same package — test them directly.\n" +
			"- 未导出函数（小写名称）因测试文件在同一包中可直接调用，无需特殊处理。\n" +
			"- For functions using stdin/stdout/files/network: use dependency injection or patch the dependency; do NOT require real I/O in tests.\n" +
			"- 对依赖 stdin/stdout/文件/网络的函数：通过依赖注入或替换依赖的方式测试，不要要求真实 I/O。\n" +
			"- CRITICAL: Before calling ANY function in your test, verify it EXISTS in the source code.\n" +
			"  Check the exact function signature. Do NOT guess or assume a function exists.\n" +
			"- 重要：调用任何函数前，必须验证它存在于提供的源码中。检查确切函数签名，不要猜测或假设函数存在。\n"

	case "java":
		return "- Test class name MUST be `ClassNameTest` where `ClassName` matches the source class name.\n" +
			"- 测试类命名必须为 `ClassNameTest`，其中 `ClassName` 与源类名匹配。\n" +
			"- Test file MUST be named `ClassNameTest.java` (e.g., `EmployeeValidatorTest.java`).\n" +
			"- 测试文件命名必须为 `ClassNameTest.java`（如 `EmployeeValidatorTest.java`）。\n" +
			"- Use JUnit 4 annotations: `@Test` for test methods.\n" +
			"- 使用 JUnit 4 注解：测试方法使用 `@Test`。\n" +
			"- Test methods must be `public void` and can be named `testXxx()` or descriptively.\n" +
			"- 测试方法必须是 `public void`，命名可以是 `testXxx()` 或有描述性的名称。\n" +
			"- Import JUnit: `import org.junit.Test; import static org.junit.Assert.*;`.\n" +
			"- 导入 JUnit：`import org.junit.Test; import static org.junit.Assert.*;`。\n" +
			"- Use `assertEquals(expected, actual)`, `assertTrue(condition)`, `assertFalse(condition)`.\n" +
			"- 断言使用 `assertEquals(expected, actual)`、`assertTrue(condition)`、`assertFalse(condition)`。\n" +
			"- For exceptions, use `@Test(expected = IllegalArgumentException.class)` or try-catch.\n" +
			"- 异常测试使用 `@Test(expected = IllegalArgumentException.class)` 或 try-catch。\n" +
			"- For static methods, call them directly: `ClassName.methodName(args)`.\n" +
			"- 静态方法直接调用：`ClassName.methodName(args)`。\n" +
			"- Test class must be `public` and in the same package as source (no package declaration for simple files).\n" +
			"- 测试类必须是 `public`，与源类在同一包中（简单文件无需包声明）。\n" +
			"- For nested/inner classes in source: import them or reference as `OuterClass.InnerClass`.\n" +
			"- 对于源码中的嵌套类：导入它们或使用 `OuterClass.InnerClass` 引用。\n" +
			"- DO NOT access private fields/methods directly. Use public API or reflection only if necessary.\n" +
			"- 不要直接访问私有字段/方法。使用公共 API，必要时才用反射。\n" +
			"- Output COMPLETE test code. Ensure all braces `{}` are properly closed.\n" +
			"- 输出完整的测试代码。确保所有 `{}` 大括号正确闭合。\n" +
			"- Check your output: if source has inner classes, test must handle them correctly.\n" +
			"- 检查输出：如果源码有嵌套类，测试必须正确处理。\n"

	case "cpp":
		return "- Use GoogleTest framework macros: `TEST()`, `EXPECT_EQ`, `ASSERT_EQ`.\n" +
			"- 使用 GoogleTest 框架宏：`TEST()`、`EXPECT_EQ`、`ASSERT_EQ`。\n" +
			"- Test function naming: `TEST(TestSuiteName, TestName)`.\n" +
			"- 测试函数命名：`TEST(TestSuiteName, TestName)`。\n" +
			"- Include `<gtest/gtest.h>` header.\n" +
			"- 包含 `<gtest/gtest.h>` 头文件。\n"

	case "javascript":
		return "- Test function names should use `test()` or `it()` from Jest.\n" +
			"- 测试函数使用 Jest 的 `test()` 或 `it()`。\n" +
			"- Use `expect()` for assertions (e.g., `expect(result).toBe(expected)`).\n" +
			"- 断言使用 `expect()`（如 `expect(result).toBe(expected)`）。\n" +
			"- Test file should be named `xxx.test.js` or `xxx.spec.js`.\n" +
			"- 测试文件命名应为 `xxx.test.js` 或 `xxx.spec.js`。\n"

	default:
		return "- Follow standard testing conventions for the language.\n" +
			"- 遵循该语言的标准测试规范。\n"
	}
}

// buildSemanticAlignmentRules produces language-specific hard constraints that
// anchor the test semantics to the source implementation.
func buildSemanticAlignmentRules(lang string) string {
	common := "- Derive expected values ONLY from the shown source code; never invent behaviour.\n" +
		"- 期望值必须严格源自所展示的源码，禁止凭空捏造行为。\n" +
		"- If a behaviour is ambiguous, test the implemented branch, not the ideal one.\n" +
		"- 行为不明确时测试实际实现的分支，而非“理想实现”。\n" +
		"- Do NOT test behaviour of symbols that are not present in the source.\n" +
		"- 不要测试源码中不存在的符号或函数。\n"
	switch lang {
	case "python":
		return common +
			"- Match exception TYPES exactly (TypeError vs ValueError vs custom).\n" +
			"- 异常类型必须完全匹配（TypeError / ValueError / 自定义异常）。\n" +
			"- Preserve integer vs float return types in assertions.\n" +
			"- 断言保留整型 / 浮点型的真实返回类型。\n"
	case "go":
		return common +
			"- When comparing errors, prefer `errors.Is` / sentinel values that exist in source.\n" +
			"- 对比 error 时优先使用源码中真实存在的 sentinel 或 `errors.Is`。\n" +
			"- Respect pointer vs value receiver semantics when constructing inputs.\n" +
			"- 构造入参时尊重指针 / 值接收器语义。\n"
	case "java":
		return common +
			"- Match declared checked exceptions; do not broaden to `Exception`.\n" +
			"- 匹配声明的 checked 异常，不要泛化为 `Exception`。\n"
	default:
		return common
	}
}

// buildErrorPreventionChecklist lists self-check items the model should run
// mentally before emitting the final answer.
func buildErrorPreventionChecklist(lang string) string {
	common := "- [ ] All imports resolve to symbols that exist in the provided source.\n" +
		"- [ ] 所有 import 均指向源码中真实存在的符号。\n" +
		"- [ ] No placeholder tests like `assert True`.\n" +
		"- [ ] 不要占位测试（如 `assert True`）。\n" +
		"- [ ] No reliance on real network / files / env variables.\n" +
		"- [ ] 不依赖真实网络 / 文件 / 环境变量。\n" +
		"- [ ] Every test has at least one meaningful assertion.\n" +
		"- [ ] 每个测试至少包含一条有意义的断言。\n"
	switch lang {
	case "python":
		return common +
			"- [ ] `def test_*` naming is respected.\n" +
			"- [ ] 函数名以 `test_` 开头。\n"
	case "go":
		return common +
			"- [ ] `func TestXxx(t *testing.T)` signature is correct; package matches source.\n" +
			"- [ ] 函数签名 `func TestXxx(t *testing.T)` 正确，package 与源码一致。\n"
	case "java":
		return common +
			"- [ ] Class is `ClassNameTest`, methods annotated with `@Test`, braces balanced.\n" +
			"- [ ] 测试类为 `ClassNameTest`，方法带 `@Test`，大括号配对。\n"
	default:
		return common
	}
}

// -----------------------------------------------------------------------------
// Small helpers (framework name, coverage targets, module name, scenario)
// -----------------------------------------------------------------------------

func coverageTargetsText() string {
	// Keep in sync with benchmark/config/models.yaml default thresholds.
	const line, branch, function = 0.7, 0.6, 0.8
	return fmt.Sprintf(
		"line >= %.0f%%, branch >= %.0f%%, function >= %.0f%%",
		line*100, branch*100, function*100,
	)
}

func languageFramework(language string) string {
	switch language {
	case "java":
		return "JUnit 4"
	case "python":
		return "pytest"
	case "go":
		return "Go testing package"
	case "cpp":
		return "GoogleTest"
	case "javascript":
		return "Jest"
	default:
		return "the standard test framework"
	}
}

func moduleImportName(sampleID string) string {
	normalized := regexp.MustCompile(`[^a-zA-Z0-9_]`).ReplaceAllString(sampleID, "_")
	normalized = strings.Trim(normalized, "_")
	if normalized == "" {
		return "solution"
	}
	if regexp.MustCompile(`^[0-9]`).MatchString(normalized) {
		return "sample_" + normalized
	}
	return normalized
}

func parseSampleMeta(sampleID, samplePath string) (string, string) {
	parts := strings.Split(sampleID, "_")
	if len(parts) >= 4 {
		complexity := strings.ToLower(parts[0])
		scenario := strings.ToLower(strings.Join(parts[2:len(parts)-1], "_"))
		return scenario, complexity
	}
	parent := strings.ToLower(filepath.Base(filepath.Dir(samplePath)))
	grand := strings.ToLower(filepath.Base(filepath.Dir(filepath.Dir(samplePath))))
	if parent != "" && parent != grand {
		return parent, ""
	}
	return "", ""
}

func mockRequirement(sourceCode string) string {
	lower := strings.ToLower(sourceCode)
	markers := []string{
		"http", "request", "socket", "open(", "file", "database",
		"sql", "redis", "grpc", "client", "os.environ", "subprocess",
	}
	for _, marker := range markers {
		if strings.Contains(lower, marker) {
			return "Use mocks/stubs/fakes for network, file, database, subprocess, or env dependencies."
		}
	}
	return "Mock only when necessary; avoid over-mocking pure functions."
}
