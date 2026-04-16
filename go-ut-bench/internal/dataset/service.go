package dataset

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go-ut-bench/internal/contracts"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ValidateSpec(spec contracts.RunSpec) error {
	if strings.TrimSpace(spec.DatasetRoot) == "" {
		return errors.New("dataset root is required")
	}
	if len(spec.Models) > 0 {
		if strings.TrimSpace(spec.OutputRoot) == "" {
			return errors.New("output root is required")
		}
		if strings.TrimSpace(spec.RunID) == "" {
			return errors.New("run id is required")
		}
		if strings.TrimSpace(spec.ConfigPath) == "" {
			return errors.New("config path is required")
		}
	}
	if spec.DatasetClass != "" && spec.DatasetClass != contracts.DatasetClassSelfContained && spec.DatasetClass != contracts.DatasetClassModuleLevel && spec.DatasetClass != contracts.DatasetClassComplexDependency {
		return fmt.Errorf("unsupported dataset class: %s", spec.DatasetClass)
	}
	if spec.DatasetScenario != "" {
		scenario := normalizeScenario(spec.DatasetScenario)
		if scenario == "" {
			return fmt.Errorf("unsupported dataset scenario: %s", spec.DatasetScenario)
		}
	}
	if spec.Mode != "" && spec.Mode != contracts.RunModeFull && spec.Mode != contracts.RunModeIncremental {
		return fmt.Errorf("unsupported mode: %s", spec.Mode)
	}
	for _, lang := range spec.Languages {
		norm := strings.ToLower(strings.TrimSpace(lang))
		if !isSupportedLanguage(norm) {
			return fmt.Errorf("unsupported language: %s", lang)
		}
	}
	if spec.MaxSamples < 0 {
		return errors.New("max samples cannot be negative")
	}
	if spec.MutationTimeout < 0 {
		return errors.New("mutation timeout cannot be negative")
	}
	if spec.MutationPolicy != "" {
		policy := strings.ToLower(strings.TrimSpace(spec.MutationPolicy))
		if policy != "warn" && policy != "fail" {
			return fmt.Errorf("unsupported mutation policy: %s", spec.MutationPolicy)
		}
	}
	return nil
}

func (s *Service) DiscoverSamples(spec contracts.RunSpec) ([]contracts.SampleRef, error) {
	if err := s.ValidateSpec(spec); err != nil {
		return nil, err
	}

	if strings.TrimSpace(spec.DatasetManifest) != "" || strings.TrimSpace(spec.DatasetLevel) != "" {
		return s.discoverFromManifest(spec)
	}

	langs := spec.Languages
	if len(langs) == 0 {
		langs = append([]string{}, contracts.SupportedLanguages...)
	}

	var all []contracts.SampleRef
	for _, langRaw := range langs {
		lang := strings.ToLower(strings.TrimSpace(langRaw))
		if lang == "" {
			continue
		}
		langDir := filepath.Join(spec.DatasetRoot, lang)
		if _, err := os.Stat(langDir); err != nil {
			continue
		}

		err := filepath.WalkDir(langDir, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if d.IsDir() {
				if d.Name() != "." && strings.HasPrefix(d.Name(), ".") {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasPrefix(d.Name(), ".") {
				return nil
			}
			if !matchLanguageExt(path, lang) {
				return nil
			}

			id := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
			rel, _ := filepath.Rel(langDir, path)
			cat := classifySampleClass(id, rel)
			scenario := classifySampleScenario(id, rel)
			if !matchDatasetClassFilter(spec.DatasetClass, cat) {
				return nil
			}
			if spec.DatasetScenario != "" && scenario != normalizeScenario(spec.DatasetScenario) {
				return nil
			}

			hash, err := fileMD5(path)
			if err != nil {
				return err
			}

			all = append(all, contracts.SampleRef{
				ID:        id,
				Language:  lang,
				Category:  cat,
				Scenario:  scenario,
				Path:      path,
				SourceMD5: hash,
			})
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	sort.Slice(all, func(i, j int) bool {
		if all[i].Language == all[j].Language {
			return all[i].ID < all[j].ID
		}
		return all[i].Language < all[j].Language
	})

	if spec.MaxSamples > 0 && len(all) > spec.MaxSamples {
		all = all[:spec.MaxSamples]
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("no dataset samples found (langs=%v class=%s)", langs, spec.DatasetClass)
	}

	return all, nil
}

func (s *Service) discoverFromManifest(spec contracts.RunSpec) ([]contracts.SampleRef, error) {
	manifestPath := strings.TrimSpace(spec.DatasetManifest)
	if manifestPath == "" {
		manifestPath = filepath.Join("configs", "dataset_"+strings.ToLower(strings.TrimSpace(spec.DatasetLevel))+".json")
	}
	mf, err := loadManifest(manifestPath, spec.DatasetRoot)
	if err != nil {
		return nil, err
	}

	langFilter := map[string]struct{}{}
	for _, lang := range spec.Languages {
		lang = strings.ToLower(strings.TrimSpace(lang))
		if lang != "" {
			langFilter[lang] = struct{}{}
		}
	}

	all := make([]contracts.SampleRef, 0, len(mf.Samples))
	for _, item := range mf.Samples {
		if len(langFilter) > 0 {
			if _, ok := langFilter[item.Language]; !ok {
				continue
			}
		}
		if !matchDatasetClassFilter(spec.DatasetClass, item.Category) {
			continue
		}
		scenario := item.Scenario
		if scenario == "" {
			scenario = classifySampleScenario(item.ID, item.Path)
		}
		if spec.DatasetScenario != "" && scenario != normalizeScenario(spec.DatasetScenario) {
			continue
		}
		if !matchLanguageExt(item.Path, item.Language) {
			continue
		}
		hash, err := fileMD5(item.Path)
		if err != nil {
			continue
		}
		all = append(all, contracts.SampleRef{
			ID:        item.ID,
			Language:  item.Language,
			Category:  item.Category,
			Scenario:  scenario,
			Path:      item.Path,
			SourceMD5: hash,
		})
	}

	sort.Slice(all, func(i, j int) bool {
		if all[i].Language == all[j].Language {
			return all[i].ID < all[j].ID
		}
		return all[i].Language < all[j].Language
	})

	if spec.MaxSamples > 0 && len(all) > spec.MaxSamples {
		all = all[:spec.MaxSamples]
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("no dataset samples found in manifest=%s", manifestPath)
	}
	return all, nil
}

func (s *Service) ValidateLayout(datasetRoot string) error {
	if strings.TrimSpace(datasetRoot) == "" {
		return errors.New("dataset root is required")
	}
	if _, err := os.Stat(datasetRoot); err != nil {
		return fmt.Errorf("dataset root not accessible: %w", err)
	}
	for _, lang := range contracts.SupportedLanguages {
		langDir := filepath.Join(datasetRoot, lang)
		if _, err := os.Stat(langDir); err != nil {
			return fmt.Errorf("missing language directory: %s", langDir)
		}
	}
	return nil
}

func classifySampleClass(sampleID string, relPath string) contracts.DatasetClass {
	lower := strings.ToLower(sampleID + "|" + relPath)
	lower = strings.ReplaceAll(lower, "\\", "/")
	if strings.Contains(lower, "self_contained") {
		return contracts.DatasetClassSelfContained
	}
	if strings.Contains(lower, "module_level") {
		return contracts.DatasetClassModuleLevel
	}
	if strings.Contains(lower, "complex_dependency") {
		return contracts.DatasetClassComplexDependency
	}
	if strings.Contains(lower, "interface_mock") {
		return contracts.DatasetClassComplexDependency
	}
	return contracts.DatasetClassSelfContained
}

func classifySampleScenario(sampleID string, relPath string) string {
	lower := strings.ToLower(sampleID + "|" + relPath)
	lower = strings.ReplaceAll(lower, "\\", "/")
	for _, scenario := range []string{"boundary", "simple_function", "complex_dependency", "interface_mock"} {
		if strings.Contains(lower, scenario) {
			return scenario
		}
	}
	return "unknown"
}

func normalizeScenario(raw string) string {
	val := strings.ToLower(strings.TrimSpace(raw))
	switch val {
	case "boundary", "simple_function", "complex_dependency", "interface_mock", "unknown":
		return val
	default:
		return ""
	}
}

func matchDatasetClassFilter(filter contracts.DatasetClass, sample contracts.DatasetClass) bool {
	if filter == "" {
		return true
	}
	if filter == sample {
		return true
	}
	if filter == contracts.DatasetClassComplexDependency && sample == contracts.DatasetClassModuleLevel {
		return true
	}
	return false
}

func isSupportedLanguage(lang string) bool {
	for _, v := range contracts.SupportedLanguages {
		if lang == v {
			return true
		}
	}
	return false
}

func matchLanguageExt(path, lang string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == languageExtByName(lang)
}

func languageExtByName(lang string) string {
	switch lang {
	case "python":
		return ".py"
	case "java":
		return ".java"
	case "go":
		return ".go"
	case "cpp":
		return ".cpp"
	default:
		return ""
	}
}

func fileMD5(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := md5.Sum(raw)
	return fmt.Sprintf("%x", h[:]), nil
}
