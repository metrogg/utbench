// dataset 包提供数据集发现和验证功能
// 负责扫描数据集目录、验证配置参数、分类样本、计算源码哈希
// 支持从目录扫描和从清单文件加载两种样本发现方式
package dataset

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"go-ut-bench/internal/contracts"
)

// Service 数据集服务
// 提供样本发现、验证、依赖检查等功能
type Service struct{}

// NewService 创建数据集服务实例
// 返回值:
//   - *Service: 数据集服务实例
func NewService() *Service {
	return &Service{}
}

// ValidateSpec 验证运行规格参数的有效性
// 检查必填字段、支持的语言/类别/场景、参数范围等
//
// 参数:
//   - spec: 运行规格说明
//
// 返回值:
//   - error: 参数无效时的错误信息，nil表示验证通过
//
// 验证项:
//  1. DatasetRoot 必填
//  2. 如果有模型，则 OutputRoot、RunID、ConfigPath 必填
//  3. DatasetClasses 必须为 self_contained 或 repo_level
//  4. Languages 必须为支持的语言
//  5. MaxSamples 和 MutationTimeout 不能为负数
func (s *Service) ValidateSpec(spec contracts.RunSpec) error {
	if strings.TrimSpace(spec.DatasetRoot) == "" {
		return errors.New("dataset root is required")
	}
	if len(spec.Models) > 0 || len(spec.Subjects) > 0 {
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
	if len(spec.DatasetClasses) > 0 {
		validClasses := map[string]bool{
			"self_contained": true,
			"repo_level":     true,
		}
		for _, c := range spec.DatasetClasses {
			if !validClasses[c] {
				return fmt.Errorf("unsupported dataset class: %s (valid: self_contained, repo_level)", c)
			}
		}
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
		if policy != "warn" && policy != "fail" && policy != "skip" {
			return fmt.Errorf("unsupported mutation policy: %s", spec.MutationPolicy)
		}
	}
	return nil
}

// DiscoverSamples 发现并返回符合条件的数据集样本列表
// 支持两种发现方式：目录扫描和清单文件加载
//
// 参数:
//   - spec: 运行规格说明，包含数据集根目录、语言、类别、场景等过滤条件
//
// 返回值:
//   - []contracts.SampleRef: 样本引用列表
//   - error: 发现过程中的错误
//
// 发现逻辑:
//  1. 如果指定了 DatasetManifest 或 DatasetLevel，从清单文件加载
//  2. 否则扫描数据集目录，根据目录结构推断样本属性
//  3. 应用语言、类别、场景过滤条件
//  4. 计算源码 MD5 哈希
//  5. 如果指定了 MaxSamples，按语言+场景分组后限制样本数
//
// 目录结构推断规则:
//  - Language: 第一级目录名（python/go/java/cpp）
//  - Category: 第二级目录名（self_contained/repo_level）
//  - Scenario: 第三级目录名（boundary/simple_function等）
//  - SampleID: 文件名（去掉扩展名）
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
			if !d.IsDir() {
				if strings.HasPrefix(d.Name(), ".") {
					return nil
				}
				if !matchLanguageExt(path, lang) {
					return nil
				}
				rel, _ := filepath.Rel(langDir, path)
				if strings.Contains(rel, "/workspace/") || strings.Contains(rel, "\\workspace\\") {
					return nil
				}

				id := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
				cat := classifySampleClass(id, rel)
				scenario := classifySampleScenario(id, rel)
				if !matchDatasetClassFilter(spec.DatasetClasses, cat) {
					return nil
				}
				if spec.DatasetScenario != "" {
					allowedScenarios := make(map[string]struct{})
					for _, sc := range strings.Split(spec.DatasetScenario, ",") {
						if n := normalizeScenario(strings.TrimSpace(sc)); n != "" {
							allowedScenarios[n] = struct{}{}
						}
					}
					if _, ok := allowedScenarios[scenario]; !ok {
						return nil
					}
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
			}

			if d.Name() != "." && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}

			if d.Name() == "workspace" {
				return filepath.SkipDir
			}

			entryPath := filepath.Join(path, "entry.py")
			metaPath := filepath.Join(path, "meta.json")
			if _, err1 := os.Stat(entryPath); err1 == nil {
				if _, err2 := os.Stat(metaPath); err2 == nil {
					id := filepath.Base(path)
					rel, _ := filepath.Rel(langDir, path)
					cat := classifySampleClass(id, rel)
					scenario := classifySampleScenario(id, rel)
					if !matchDatasetClassFilter(spec.DatasetClasses, cat) {
						return filepath.SkipDir
					}
					if spec.DatasetScenario != "" {
						allowedScenarios := make(map[string]struct{})
						for _, sc := range strings.Split(spec.DatasetScenario, ",") {
							if n := normalizeScenario(strings.TrimSpace(sc)); n != "" {
								allowedScenarios[n] = struct{}{}
							}
						}
						if _, ok := allowedScenarios[scenario]; !ok {
							return filepath.SkipDir
						}
					}

					hash, err := fileMD5(entryPath)
					if err != nil {
						return err
					}

					all = append(all, contracts.SampleRef{
						ID:        id,
						Language:  lang,
						Category:  cat,
						Scenario:  scenario,
						Path:      entryPath,
						SourceMD5: hash,
					})
					return filepath.SkipDir
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	sortSampleRefs(all)
	if spec.MaxSamples > 0 {
		all = applyMaxSamplesPerLanguageScenario(all, spec.MaxSamples)
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("no dataset samples found (langs=%v classes=%v)", langs, spec.DatasetClasses)
	}

	return all, nil
}

// discoverFromManifest 从清单文件加载样本列表
// 当 spec.DatasetManifest 或 spec.DatasetLevel 指定时使用此方式
//
// 参数:
//   - spec: 运行规格说明
//
// 返回值:
//   - []contracts.SampleRef: 样本引用列表
//   - error: 加载过程中的错误
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
		if !matchDatasetClassFilter(spec.DatasetClasses, item.Category) {
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

	sortSampleRefs(all)
	if spec.MaxSamples > 0 {
		all = applyMaxSamplesPerLanguageScenario(all, spec.MaxSamples)
	}
	if len(all) == 0 {
		return nil, fmt.Errorf("no dataset samples found in manifest=%s", manifestPath)
	}
	return all, nil
}

// ValidateLayout 验证数据集目录结构是否符合要求
// 检查数据集根目录是否存在，以及各语言目录是否完整
//
// 参数:
//   - datasetRoot: 数据集根目录路径
//
// 返回值:
//   - error: 验证失败的错误信息，nil表示验证通过
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

// CheckDependencies 检查 Python 样本的依赖是否可用
// 通过嵌入的 Python 脚本扫描数据集中的 import 语句，验证是否可导入
//
// 参数:
//   - datasetRoot: 数据集根目录路径
//   - languages: 要检查的语言列表
//
// 返回值:
//   - string: 检查结果输出（ALL_DEPS_AVAILABLE 或缺失依赖列表）
//   - error: 检查过程中的错误
//
// 检查逻辑:
//  1. 扫描所有 Python 文件的 import 语句
//  2. 排除标准库模块
//  3. 尝试导入每个第三方模块
//  4. 对于缺失的模块，提供安装建议
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

// ValidateLayout 验证数据集目录结构（公开版本）
// 检查数据集根目录和各语言子目录是否存在
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

// classifySampleClass 从样本ID和相对路径推断数据集类别
// 根据路径中的 self_contained/repo_level 关键字判断
func classifySampleClass(sampleID string, relPath string) contracts.DatasetClass {
	lower := strings.ToLower(sampleID + "|" + relPath)
	lower = strings.ReplaceAll(lower, "\\", "/")
	if strings.Contains(lower, "self_contained") {
		return contracts.DatasetClassSelfContained
	}
	if strings.Contains(lower, "repo_level") {
		return contracts.DatasetClassRepoLevel
	}
	for _, scenario := range contracts.SupportedScenarios {
		if strings.Contains(lower, scenario) {
			return contracts.DatasetClassSelfContained
		}
	}
	return contracts.DatasetClassSelfContained
}

// classifySampleScenario 从样本ID和相对路径推断场景类型
// 根据路径中的 boundary/simple_function/complex_dependency/interface_mock 关键字判断
func classifySampleScenario(sampleID string, relPath string) string {
	lower := strings.ToLower(sampleID + "|" + relPath)
	lower = strings.ReplaceAll(lower, "\\", "/")
	for _, scenario := range contracts.SupportedScenarios {
		if strings.Contains(lower, scenario) {
			return scenario
		}
	}
	return "unknown"
}

// normalizeScenario 规范化场景名称
// 返回有效的场景名称，无效输入返回空字符串
func normalizeScenario(raw string) string {
	val := strings.ToLower(strings.TrimSpace(raw))
	valid := map[string]bool{"unknown": true}
	for _, s := range contracts.SupportedScenarios {
		valid[s] = true
	}
	if valid[val] {
		return val
	}
	return ""
}

// matchDatasetClassFilter 检查样本类别是否匹配过滤条件
// 如果过滤条件为空，则匹配所有类别
func matchDatasetClassFilter(filters []string, sample contracts.DatasetClass) bool {
	if len(filters) == 0 {
		return true
	}
	for _, f := range filters {
		if f == string(sample) {
			return true
		}
	}
	return false
}

// isSupportedLanguage 检查语言是否为支持的语言
// 支持的语言：python、go、java、cpp
func isSupportedLanguage(lang string) bool {
	for _, v := range contracts.SupportedLanguages {
		if lang == v {
			return true
		}
	}
	return false
}

// matchLanguageExt 检查文件扩展名是否匹配指定语言
// 根据语言返回对应的标准扩展名进行比较
func matchLanguageExt(path, lang string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == languageExtByName(lang)
}

// languageExtByName 根据语言名称返回对应的文件扩展名
// python -> .py, java -> .java, go -> .go, cpp -> .cpp
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

// fileMD5 计算文件的 MD5 哈希值
// 用于样本源码的唯一标识和变更检测
func fileMD5(path string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	h := md5.Sum(raw)
	return fmt.Sprintf("%x", h[:]), nil
}

// sortSampleRefs 对样本引用列表进行排序
// 排序优先级：Language > Scenario > ID > Category > Path
func sortSampleRefs(samples []contracts.SampleRef) {
	sort.Slice(samples, func(i, j int) bool {
		if samples[i].Language != samples[j].Language {
			return samples[i].Language < samples[j].Language
		}
		if samples[i].Scenario != samples[j].Scenario {
			return samples[i].Scenario < samples[j].Scenario
		}
		if samples[i].ID != samples[j].ID {
			return samples[i].ID < samples[j].ID
		}
		if samples[i].Category != samples[j].Category {
			return samples[i].Category < samples[j].Category
		}
		return samples[i].Path < samples[j].Path
	})
}

// applyMaxSamplesPerLanguageScenario 按语言+场景分组后限制样本数量
// 对每个 Language+Scenario 组合，最多保留 maxSamples 个样本
func applyMaxSamplesPerLanguageScenario(samples []contracts.SampleRef, maxSamples int) []contracts.SampleRef {
	if maxSamples <= 0 || len(samples) == 0 {
		return samples
	}

	grouped := make(map[string][]contracts.SampleRef)
	for _, sample := range samples {
		key := sample.Language + "|" + sample.Scenario
		grouped[key] = append(grouped[key], sample)
	}

	keys := make([]string, 0, len(grouped))
	for key := range grouped {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := make([]contracts.SampleRef, 0, len(samples))
	for _, key := range keys {
		group := grouped[key]
		sortSampleRefs(group)
		if len(group) > maxSamples {
			group = group[:maxSamples]
		}
		result = append(result, group...)
	}

	sortSampleRefs(result)
	return result
}

// loadRepoLevelMeta 加载 repo_level 类型样本的元数据
// 从样本同目录的 .meta.json 文件加载 RepoLevelMeta 信息
func loadRepoLevelMeta(samplePath string) (*contracts.RepoLevelMeta, bool) {
	dir := filepath.Dir(samplePath)
	base := filepath.Base(samplePath)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]
	metaPath := filepath.Join(dir, name+".meta.json")
	if _, err := os.Stat(metaPath); err != nil {
		return nil, false
	}
	raw, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, false
	}
	var meta contracts.RepoLevelMeta
	if err := json.Unmarshal(raw, &meta); err != nil {
		return nil, false
	}
	return &meta, true
}
