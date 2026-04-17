package dataset

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
		for _, s := range strings.Split(spec.DatasetScenario, ",") {
			normalized := normalizeScenario(strings.TrimSpace(s))
			if normalized == "" {
				return fmt.Errorf("unsupported dataset scenario: %s", strings.TrimSpace(s))
			}
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
		if spec.DatasetScenario != "" {
			allowed := make(map[string]struct{})
			for _, s := range strings.Split(spec.DatasetScenario, ",") {
				allowed[normalizeScenario(strings.TrimSpace(s))] = struct{}{}
			}
			if _, ok := allowed[scenario]; !ok {
				continue
			}
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
		if all[i].Scenario == all[j].Scenario {
			if all[i].Language == all[j].Language {
				return all[i].ID < all[j].ID
			}
			return all[i].Language < all[j].Language
		}
		return all[i].Scenario < all[j].Scenario
	})

	if spec.MaxSamples > 0 && len(all) > spec.MaxSamples {
		groups := make(map[string][]contracts.SampleRef)
		for _, s := range all {
			groups[s.Scenario] = append(groups[s.Scenario], s)
		}
		scenarios := make([]string, 0, len(groups))
		for k := range groups {
			scenarios = append(scenarios, k)
		}
		sort.Strings(scenarios)
		result := make([]contracts.SampleRef, 0, spec.MaxSamples)
		indices := make(map[string]int)
		for len(result) < spec.MaxSamples {
			took := false
			for _, sc := range scenarios {
				if indices[sc] >= len(groups[sc]) {
					continue
				}
				result = append(result, groups[sc][indices[sc]])
				indices[sc]++
				took = true
				if len(result) >= spec.MaxSamples {
					break
				}
			}
			if !took {
				break
			}
		}
		all = result
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

func (s *Service) CheckDependencies(datasetRoot string, languages []string) (string, error) {
	if strings.TrimSpace(datasetRoot) == "" {
		return "", errors.New("dataset root is required")
	}

	script := `import sys
import os
import re
import importlib
import itertools

missing = {}
std_lib = frozenset(sys.stdlib_module_names)

def extract_imports(path):
    try:
        with open(path, 'r', encoding='utf-8') as f:
            content = f.read()
    except:
        return []
    imports = set()
    for m in re.finditer(r'^(?:from\s+([\w.]+)|import\s+([\w.]+))', content, re.MULTILINE):
        mod = m.group(1) or m.group(2)
        if not mod or mod.strip() == '':
            continue
        mod = mod.split('.')[0]
        if mod and mod not in std_lib and not mod.startswith('_') and not mod.startswith('.'):
            imports.add(mod)
    return list(imports)

visited = set()
langs_set = set()
for l in sys.argv[2:]:
    langs_set.add(l.strip().lower())
if not langs_set:
    langs_set = {'python', 'go', 'java', 'cpp'}

skip_dirs = {'.git', 'artifacts', 'storage', 'configs', '__pycache__', 'node_modules', '.mutmut_shim', '.pytest_cache'}

for root, dirs, files in os.walk('.'):
    dirs[:] = [d for d in dirs if d not in skip_dirs]
    for f in files:
        if not f.endswith('.py'):
            continue
        path = os.path.join(root, f)
        rel = os.path.relpath(path, '.')
        lang = None
        if '/python/' in rel or '\\python\\' in rel or rel.startswith('python\\') or rel.startswith('python/'):
            lang = 'python'
        elif '/go/' in rel or '\\go\\' in rel or rel.startswith('go\\') or rel.startswith('go/'):
            lang = 'go'
        elif '/java/' in rel or '\\java\\' in rel or rel.startswith('java\\') or rel.startswith('java/'):
            lang = 'java'
        elif '/cpp/' in rel or '\\cpp\\' in rel or rel.startswith('cpp\\') or rel.startswith('cpp/'):
            lang = 'cpp'
        if lang is None or lang not in langs_set:
            continue
        key = rel
        if key in visited:
            continue
        path = os.path.join(root, f)
        rel = os.path.relpath(path, '.')
        lang = None
        if '/python/' in rel or '\\\\python\\\\' in rel or rel.startswith('python\\\\') or rel.startswith('python/'):
            lang = 'python'
        elif '/go/' in rel or '\\\\go\\\\' in rel or rel.startswith('go\\\\') or rel.startswith('go/'):
            lang = 'go'
        elif '/java/' in rel or '\\\\java\\\\' in rel or rel.startswith('java\\\\') or rel.startswith('java/'):
            lang = 'java'
        elif '/cpp/' in rel or '\\\\cpp\\\\' in rel or rel.startswith('cpp\\\\') or rel.startswith('cpp/'):
            lang = 'cpp'
        if lang is None:
            continue
        key = rel
        if key in visited:
            continue
        visited.add(key)
        for mod in extract_imports(path):
            try:
                importlib.import_module(mod)
            except ImportError as e:
                msg = str(e).split('"')[0].strip()
                if mod not in missing:
                    missing[mod] = {'samples': [], 'error': msg}
                missing[mod]['samples'].append(rel)

install_hints = {
    'numpy': 'pip install numpy',
    'pandas': 'pip install pandas',
    'scipy': 'pip install scipy',
    'scikit-learn': 'pip install scikit-learn',
    'PIL': 'pip install Pillow',
    'cv2': 'pip install opencv-python',
    'matplotlib': 'pip install matplotlib',
    'requests': 'pip install requests',
    'yaml': 'pip install pyyaml',
    'turtle': 'sudo apt install python3-tk',
    'tkinter': 'sudo apt install python3-tk',
    'h2': 'pip install h2',
    'httpcore': 'pip install httpcore',
    'httpx': 'pip install httpx',
    'tqdm': 'pip install tqdm',
    'huggingface_hub': 'pip install huggingface_hub',
    'transformers': 'pip install transformers',
    'datasets': 'pip install datasets',
    'dask': 'pip install dask',
    'distributed': 'pip install distributed',
    'fsspec': 'pip install fsspec',
    's3fs': 'pip install s3fs',
    'smbclient': 'pip install smbprotocol',
    'smbprotocol': 'pip install smbprotocol',
}

if not missing:
    print('ALL_DEPS_AVAILABLE')
else:
    for mod, info in sorted(missing.items()):
        print(f'MISSING:{mod}')
        print(f"  Error: {info['error']}")
        for s in info['samples'][:3]:
            print(f'  Sample: {s}')
        if len(info['samples']) > 3:
            print(f"  ... and {len(info['samples'])-3} more")
        hint = install_hints.get(mod, f'pip install {mod}')
        print(f'  Install: {hint}')
`
	cmd := exec.Command("python3", "-c", script)
	cmd.Dir = datasetRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("dependency check failed: %w\nOutput: %s", err, string(out))
	}
	return string(out), nil
}

func ValidateLayout(datasetRoot string) error {
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
