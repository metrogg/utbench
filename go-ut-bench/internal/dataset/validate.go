package dataset

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go-ut-bench/internal/contracts"
)

type ValidateOptions struct {
	DatasetRoot string
	Languages   []string
	Classes     []string
	Scenario    string
	Strict      bool
}

type ValidationReport struct {
	DatasetRoot string            `json:"dataset_root"`
	Total       int               `json:"total_samples"`
	Counts      []ValidationCount `json:"counts"`
	Errors      []ValidationIssue `json:"errors,omitempty"`
	Warnings    []ValidationIssue `json:"warnings,omitempty"`
	OK          bool              `json:"ok"`
}

type ValidationCount struct {
	Language string `json:"language"`
	Class    string `json:"class"`
	Scenario string `json:"scenario"`
	Count    int    `json:"count"`
}

type ValidationIssue struct {
	Severity string `json:"severity"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Language string `json:"language,omitempty"`
	SampleID string `json:"sample_id,omitempty"`
	Path     string `json:"path,omitempty"`
}

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
