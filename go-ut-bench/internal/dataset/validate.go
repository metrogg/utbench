// dataset/validate.go 提供数据集验证功能
// 检查数据集完整性、样本结构、潜在风险模式
package dataset

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go-ut-bench/internal/contracts"
)

// ValidateOptions 验证选项
// 定义验证范围和严格程度
type ValidateOptions struct {
	DatasetRoot string   // 数据集根目录
	Languages   []string // 要验证的语言列表
	Classes     []string // 数据集类别过滤
	Scenario    string   // 场景过滤
	Strict      bool     // 是否严格模式
}

// ValidationReport 验证报告
// 包含验证结果统计和问题列表
type ValidationReport struct {
	DatasetRoot string            `json:"dataset_root"` // 数据集根目录
	Total       int               `json:"total_samples"` // 总样本数
	Counts      []ValidationCount `json:"counts"`       // 分类统计
	Errors      []ValidationIssue `json:"errors,omitempty"` // 错误列表
	Warnings    []ValidationIssue `json:"warnings,omitempty"` // 警告列表
	OK          bool              `json:"ok"`           // 是否通过验证
}

// ValidationCount 分类统计
// 按语言、类别、场景统计样本数
type ValidationCount struct {
	Language string `json:"language"` // 编程语言
	Class    string `json:"class"`    // 数据集类别
	Scenario string `json:"scenario"` // 场景名称
	Count    int    `json:"count"`    // 样本数量
}

// ValidationIssue 验证问题
// 描述单个错误或警告
type ValidationIssue struct {
	Severity string `json:"severity"`          // 严重程度（error/warning）
	Code     string `json:"code"`              // 问题代码
	Message  string `json:"message"`           // 问题消息
	Language string `json:"language,omitempty"` // 相关语言
	SampleID string `json:"sample_id,omitempty"` // 相关样本ID
	Path     string `json:"path,omitempty"`     // 相关路径
}

// ValidateReadiness 验证数据集是否就绪
// 检查目录结构、样本完整性、潜在风险
//
// 参数:
//   - opts: 验证选项
//
// 返回值:
//   - ValidationReport: 验证报告
func (s *Service) ValidateReadiness(opts ValidateOptions) ValidationReport {
	root := strings.TrimSpace(opts.DatasetRoot)
	report := ValidationReport{DatasetRoot: root}
	if root == "" {
		report.Errors = append(report.Errors, validationIssue("error", "dataset_root_required", "dataset root is required", "", "", ""))
		report.OK = false
		return report
	}
	if _, err := os.Stat(root); err != nil {
		report.Errors = append(report.Errors, validationIssue("error", "dataset_root_not_accessible", err.Error(), "", "", root))
		report.OK = false
		return report
	}

	langs := normalizeLangs(opts.Languages)
	classFilters := normalizeClassFilters(opts.Classes)
	scenarioFilters := normalizeScenarioFilters(opts.Scenario)
	counts := map[string]*ValidationCount{}
	seen := map[string]string{}

	for _, lang := range langs {
		langDir := filepath.Join(root, lang)
		if _, err := os.Stat(langDir); err != nil {
			report.Errors = append(report.Errors, validationIssue("error", "missing_language_dir", "missing language directory", lang, "", langDir))
			continue
		}

		_ = filepath.WalkDir(langDir, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				report.Errors = append(report.Errors, validationIssue("error", "path_not_readable", walkErr.Error(), lang, "", path))
				return nil
			}
			if d.IsDir() {
				if d.Name() != "." && strings.HasPrefix(d.Name(), ".") {
					return filepath.SkipDir
				}
				if d.Name() == "workspace" || d.Name() == "__pycache__" {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasPrefix(d.Name(), ".") {
				return nil
			}

			rel, _ := filepath.Rel(langDir, path)
			if strings.Contains(rel, string(filepath.Separator)+"workspace"+string(filepath.Separator)) {
				return nil
			}
			if !matchLanguageExt(path, lang) {
				report.Warnings = append(report.Warnings, validationIssue("warning", "extension_mismatch", "file extension does not match language", lang, "", path))
				return nil
			}

			id := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
			class := classifySampleClass(id, rel)
			scenario := classifySampleScenario(id, rel)
			if !matchDatasetClassFilter(classFilters, class) || !matchScenarioFilter(scenarioFilters, scenario) {
				return nil
			}
			if scenario == "unknown" {
				report.Errors = append(report.Errors, validationIssue("error", "unknown_scenario", "sample scenario could not be inferred", lang, id, path))
			}

			dupKey := lang + "|" + id
			if prior, ok := seen[dupKey]; ok {
				report.Errors = append(report.Errors, validationIssue("error", "duplicate_sample_id", "duplicate sample id; first seen at "+prior, lang, id, path))
			} else {
				seen[dupKey] = path
			}

			raw, err := os.ReadFile(path)
			if err != nil {
				report.Errors = append(report.Errors, validationIssue("error", "path_not_readable", err.Error(), lang, id, path))
				return nil
			}
			report.Total++
			countKey := lang + "|" + string(class) + "|" + scenario
			if _, ok := counts[countKey]; !ok {
				counts[countKey] = &ValidationCount{Language: lang, Class: string(class), Scenario: scenario}
			}
			counts[countKey].Count++
			report.Warnings = append(report.Warnings, scanRiskWarnings(lang, id, path, string(raw))...)
			return nil
		})
	}

	for _, count := range counts {
		report.Counts = append(report.Counts, *count)
	}
	sort.Slice(report.Counts, func(i, j int) bool {
		if report.Counts[i].Language != report.Counts[j].Language {
			return report.Counts[i].Language < report.Counts[j].Language
		}
		if report.Counts[i].Class != report.Counts[j].Class {
			return report.Counts[i].Class < report.Counts[j].Class
		}
		return report.Counts[i].Scenario < report.Counts[j].Scenario
	})
	sortValidationIssues(report.Errors)
	sortValidationIssues(report.Warnings)
	report.OK = len(report.Errors) == 0
	return report
}

// normalizeLangs 规范化语言列表
// 空列表返回所有支持的语言
//
// 参数:
//   - langs: 输入语言列表
//
// 返回值:
//   - []string: 规范化后的语言列表
func normalizeLangs(langs []string) []string {
	if len(langs) == 0 {
		return append([]string{}, contracts.SupportedLanguages...)
	}
	out := make([]string, 0, len(langs))
	seen := map[string]struct{}{}
	for _, lang := range langs {
		lang = strings.ToLower(strings.TrimSpace(lang))
		if lang == "" {
			continue
		}
		if _, ok := seen[lang]; ok {
			continue
		}
		seen[lang] = struct{}{}
		out = append(out, lang)
	}
	return out
}

func normalizeClassFilters(classes []string) []string {
	out := make([]string, 0, len(classes))
	for _, class := range classes {
		class = strings.TrimSpace(class)
		if class != "" {
			out = append(out, class)
		}
	}
	return out
}

func normalizeScenarioFilters(raw string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, part := range strings.Split(raw, ",") {
		scenario := normalizeScenario(strings.TrimSpace(part))
		if scenario != "" {
			out[scenario] = struct{}{}
		}
	}
	return out
}

func matchScenarioFilter(filters map[string]struct{}, scenario string) bool {
	if len(filters) == 0 {
		return true
	}
	_, ok := filters[scenario]
	return ok
}

// scanRiskWarnings 扫描源代码中的潜在风险模式
// 检测网络IO、进程调用、随机性、文件IO等
//
// 参数:
//   - lang: 编程语言
//   - id: 样本ID
//   - path: 文件路径
//   - source: 源代码内容
//
// 返回值:
//   - []ValidationIssue: 风险警告列表
func scanRiskWarnings(lang, id, path, source string) []ValidationIssue {
	lower := strings.ToLower(source)
	checks := []struct {
		code    string
		message string
		needles []string
	}{
		{code: "external_network_io", message: "sample appears to use network, SMTP, HTTP, or sockets", needles: []string{"smtplib", "requests.", "http.client", "net/http", "socket", "smtp", "urlopen", "fetch("}},
		{code: "external_process", message: "sample appears to spawn external processes", needles: []string{"subprocess", "os.system", "exec.command", "processbuilder", "system("}},
		{code: "nondeterministic_time_random", message: "sample appears to use time or randomness", needles: []string{"random.", "time.", "datetime.now", "date.now", "system.currenttimemillis", "std::chrono"}},
		{code: "filesystem_io", message: "sample appears to use filesystem I/O", needles: []string{"open(", "os.open", "os.readfile", "ioutil.", "files.", "ifstream", "ofstream"}},
		{code: "exponential_complexity", message: "sample appears to contain combinations/permutations or powerset-style logic", needles: []string{"itertools.combinations", "itertools.permutations", "combinations(", "permutations("}},
	}
	var out []ValidationIssue
	for _, check := range checks {
		for _, needle := range check.needles {
			if strings.Contains(lower, strings.ToLower(needle)) {
				out = append(out, validationIssue("warning", check.code, check.message, lang, id, path))
				break
			}
		}
	}
	return out
}

// validationIssue 创建验证问题结构
//
// 参数:
//   - severity: 严重程度
//   - code: 问题代码
//   - message: 问题消息
//   - lang: 编程语言
//   - sampleID: 样本ID
//   - path: 文件路径
//
// 返回值:
//   - ValidationIssue: 问题结构
func validationIssue(severity, code, message, lang, sampleID, path string) ValidationIssue {
	return ValidationIssue{
		Severity: severity,
		Code:     code,
		Message:  message,
		Language: lang,
		SampleID: sampleID,
		Path:     path,
	}
}

func sortValidationIssues(items []ValidationIssue) {
	sort.Slice(items, func(i, j int) bool {
		if items[i].Severity != items[j].Severity {
			return items[i].Severity < items[j].Severity
		}
		if items[i].Code != items[j].Code {
			return items[i].Code < items[j].Code
		}
		if items[i].Language != items[j].Language {
			return items[i].Language < items[j].Language
		}
		return items[i].Path < items[j].Path
	})
}
