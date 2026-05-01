package agentconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go-ut-bench/internal/contracts"

	"gopkg.in/yaml.v3"
)

const (
	KindModelAPI = "model_api"
	KindCLIAgent = "cli_agent"
	NoSkill      = "no_skill"
)

type FrameworkSpec struct {
	Name             string            `json:"name"`
	Kind             string            `json:"kind"`
	Enabled          bool              `json:"enabled"`
	Command          string            `json:"command,omitempty"`
	DockerImage      string            `json:"docker_image,omitempty"`
	SandboxMode      string            `json:"sandbox_mode,omitempty"`
	TimeoutSeconds   int               `json:"timeout_seconds,omitempty"`
	OutputGlobs      []string          `json:"output_globs,omitempty"`
	Env              map[string]string `json:"env,omitempty"`
	EnvFromHost      []string          `json:"env_from_host,omitempty"`
	CompatibleModels []string          `json:"compatible_models,omitempty"`
	CompatibleLangs  []string          `json:"compatible_languages,omitempty"`
	DisableNoSkill   bool              `json:"disable_no_skill,omitempty"`
	NetworkDisabled  bool              `json:"network_disabled,omitempty"`
	CPU              string            `json:"cpu,omitempty"`
	Memory           string            `json:"memory,omitempty"`
}

type SubjectEntry struct {
	ID        string   `yaml:"id"`
	Enabled   *bool    `yaml:"enabled"`
	Kind      string   `yaml:"kind"`
	Framework string   `yaml:"framework"`
	Model     string   `yaml:"model"`
	Skill     string   `yaml:"skill"`
	Labels    []string `yaml:"labels"`
	Tags      []string `yaml:"tags"`
}

type ResolvedSubject struct {
	Spec      contracts.SubjectSpec
	Framework FrameworkSpec
	Skill     contracts.SkillSpec
}

type fileConfig struct {
	Models     []string `yaml:"models"`
	Frameworks map[string]struct {
		Enabled             *bool             `yaml:"enabled"`
		Kind                string            `yaml:"kind"`
		Command             string            `yaml:"command"`
		DockerImage         string            `yaml:"docker_image"`
		SandboxMode         string            `yaml:"sandbox_mode"`
		TimeoutSeconds      int               `yaml:"timeout_seconds"`
		OutputGlobs         []string          `yaml:"output_globs"`
		Env                 map[string]string `yaml:"env"`
		EnvFromHost         []string          `yaml:"env_from_host"`
		CompatibleModels    []string          `yaml:"compatible_models"`
		CompatibleLanguages []string          `yaml:"compatible_languages"`
		DisableNoSkill      bool              `yaml:"disable_no_skill"`
		NetworkDisabled     *bool             `yaml:"network_disabled"`
		CPU                 string            `yaml:"cpu"`
		Memory              string            `yaml:"memory"`
	} `yaml:"frameworks"`
	Skills map[string]struct {
		Enabled              *bool    `yaml:"enabled"`
		Version              string   `yaml:"version"`
		Description          string   `yaml:"description"`
		InstructionPath      string   `yaml:"instruction_path"`
		Files                []string `yaml:"files"`
		InjectMode           string   `yaml:"inject_mode"`
		CompatibleFrameworks []string `yaml:"compatible_frameworks"`
		CompatibleLanguages  []string `yaml:"compatible_languages"`
	} `yaml:"skills"`
	Subjects []SubjectEntry `yaml:"subjects"`
}

func Load(path string, modelNames []string, selected []string) ([]ResolvedSubject, error) {
	modelNames = uniqueNonEmpty(modelNames)
	if len(modelNames) == 0 {
		return nil, fmt.Errorf("at least one model is required to build subjects")
	}
	if strings.TrimSpace(path) == "" {
		return filterSubjects(defaultModelAPISubjects(modelNames), selected)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg fileConfig
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.Models) > 0 {
		modelNames = intersect(modelNames, uniqueNonEmpty(cfg.Models))
	}
	if len(modelNames) == 0 {
		return nil, fmt.Errorf("agents config selected no models")
	}
	resolved := defaultModelAPISubjects(modelNames)
	frameworks := normalizeFrameworks(cfg)
	skills := normalizeSkills(cfg, filepath.Dir(path))
	if len(cfg.Subjects) > 0 {
		for _, entry := range cfg.Subjects {
			if entry.Enabled != nil && !*entry.Enabled {
				continue
			}
			rs, ok, err := resolveSubjectEntry(entry, frameworks, skills)
			if err != nil {
				return nil, err
			}
			if ok {
				resolved = append(resolved, rs)
			}
		}
	} else {
		resolved = append(resolved, expandCartesian(modelNames, frameworks, skills)...)
	}
	resolved = dedupeSubjects(resolved)
	return filterSubjects(resolved, selected)
}

func normalizeFrameworks(cfg fileConfig) map[string]FrameworkSpec {
	out := map[string]FrameworkSpec{}
	for name, item := range cfg.Frameworks {
		enabled := true
		if item.Enabled != nil {
			enabled = *item.Enabled
		}
		if !enabled {
			continue
		}
		kind := strings.TrimSpace(item.Kind)
		if kind == "" {
			kind = KindCLIAgent
		}
		networkDisabled := true
		if item.NetworkDisabled != nil {
			networkDisabled = *item.NetworkDisabled
		}
		out[name] = FrameworkSpec{
			Name:             name,
			Kind:             kind,
			Enabled:          true,
			Command:          item.Command,
			DockerImage:      item.DockerImage,
			SandboxMode:      defaultString(item.SandboxMode, "docker"),
			TimeoutSeconds:   item.TimeoutSeconds,
			OutputGlobs:      item.OutputGlobs,
			Env:              item.Env,
			EnvFromHost:      uniqueNonEmpty(item.EnvFromHost),
			CompatibleModels: item.CompatibleModels,
			CompatibleLangs:  item.CompatibleLanguages,
			DisableNoSkill:   item.DisableNoSkill,
			NetworkDisabled:  networkDisabled,
			CPU:              item.CPU,
			Memory:           item.Memory,
		}
	}
	return out
}

func normalizeSkills(cfg fileConfig, baseDir string) map[string]contracts.SkillSpec {
	out := map[string]contracts.SkillSpec{
		NoSkill: {Name: NoSkill, Enabled: true, InjectMode: "none"},
	}
	for name, item := range cfg.Skills {
		enabled := true
		if item.Enabled != nil {
			enabled = *item.Enabled
		}
		if !enabled {
			continue
		}
		instructionPath := resolveRelative(baseDir, item.InstructionPath)
		files := make([]string, 0, len(item.Files))
		for _, f := range item.Files {
			files = append(files, resolveRelative(baseDir, f))
		}
		out[name] = contracts.SkillSpec{
			Name:                 name,
			Version:              item.Version,
			Description:          item.Description,
			InstructionPath:      instructionPath,
			Files:                files,
			InjectMode:           defaultString(item.InjectMode, "prompt_append"),
			CompatibleFrameworks: item.CompatibleFrameworks,
			CompatibleLanguages:  item.CompatibleLanguages,
			Enabled:              true,
		}
	}
	return out
}

func defaultModelAPISubjects(models []string) []ResolvedSubject {
	out := make([]ResolvedSubject, 0, len(models))
	for _, model := range models {
		spec := contracts.SubjectSpec{
			ID:        SubjectID(KindModelAPI, model, NoSkill),
			Kind:      KindModelAPI,
			Framework: KindModelAPI,
			Model:     model,
			Skill:     NoSkill,
		}
		out = append(out, ResolvedSubject{
			Spec:      spec,
			Framework: FrameworkSpec{Name: KindModelAPI, Kind: KindModelAPI, Enabled: true},
			Skill:     contracts.SkillSpec{Name: NoSkill, Enabled: true, InjectMode: "none"},
		})
	}
	return out
}

func expandCartesian(models []string, frameworks map[string]FrameworkSpec, skills map[string]contracts.SkillSpec) []ResolvedSubject {
	var out []ResolvedSubject
	frameworkNames := sortedKeys(frameworks)
	skillNames := sortedSkillKeys(skills)
	for _, frameworkName := range frameworkNames {
		fw := frameworks[frameworkName]
		if fw.Kind == KindModelAPI || fw.Name == KindModelAPI {
			continue
		}
		for _, model := range models {
			if !isCompatible(model, fw.CompatibleModels) {
				continue
			}
			for _, skillName := range skillNames {
				if skillName == NoSkill && fw.DisableNoSkill {
					continue
				}
				skill := skills[skillName]
				if !isSkillCompatible(skill, fw.Name) {
					continue
				}
				out = append(out, buildResolvedSubject(fw, model, skill, nil, nil, ""))
			}
		}
	}
	return out
}

func resolveSubjectEntry(entry SubjectEntry, frameworks map[string]FrameworkSpec, skills map[string]contracts.SkillSpec) (ResolvedSubject, bool, error) {
	frameworkName := defaultString(entry.Framework, KindModelAPI)
	model := strings.TrimSpace(entry.Model)
	if model == "" {
		return ResolvedSubject{}, false, fmt.Errorf("subject %q missing model", entry.ID)
	}
	skillName := defaultString(entry.Skill, NoSkill)
	if frameworkName == KindModelAPI {
		return buildResolvedSubject(FrameworkSpec{Name: KindModelAPI, Kind: KindModelAPI, Enabled: true}, model, contracts.SkillSpec{Name: NoSkill, Enabled: true, InjectMode: "none"}, entry.Labels, entry.Tags, entry.ID), true, nil
	}
	fw, ok := frameworks[frameworkName]
	if !ok {
		return ResolvedSubject{}, false, fmt.Errorf("subject %q references unknown framework %q", entry.ID, frameworkName)
	}
	skill, ok := skills[skillName]
	if !ok {
		return ResolvedSubject{}, false, fmt.Errorf("subject %q references unknown skill %q", entry.ID, skillName)
	}
	return buildResolvedSubject(fw, model, skill, entry.Labels, entry.Tags, entry.ID), true, nil
}

func buildResolvedSubject(fw FrameworkSpec, model string, skill contracts.SkillSpec, labels, tags []string, explicitID string) ResolvedSubject {
	if skill.Name == "" {
		skill = contracts.SkillSpec{Name: NoSkill, Enabled: true, InjectMode: "none"}
	}
	id := strings.TrimSpace(explicitID)
	if id == "" {
		id = SubjectID(fw.Name, model, skill.Name)
	}
	return ResolvedSubject{
		Spec: contracts.SubjectSpec{
			ID:        id,
			Kind:      fw.Kind,
			Framework: fw.Name,
			Model:     model,
			Skill:     skill.Name,
			Labels:    labels,
			Tags:      tags,
		},
		Framework: fw,
		Skill:     skill,
	}
}

func SubjectID(framework, model, skill string) string {
	framework = defaultString(framework, KindModelAPI)
	skill = defaultString(skill, NoSkill)
	return sanitize(framework) + "__" + sanitize(model) + "__" + sanitize(skill)
}

func filterSubjects(subjects []ResolvedSubject, selected []string) ([]ResolvedSubject, error) {
	selected = uniqueNonEmpty(selected)
	if len(selected) == 0 {
		sort.Slice(subjects, func(i, j int) bool { return subjects[i].Spec.ID < subjects[j].Spec.ID })
		return subjects, nil
	}
	want := map[string]struct{}{}
	for _, id := range selected {
		want[id] = struct{}{}
	}
	var out []ResolvedSubject
	for _, subject := range subjects {
		if _, ok := want[subject.Spec.ID]; ok {
			out = append(out, subject)
			delete(want, subject.Spec.ID)
		}
	}
	if len(want) > 0 {
		missing := sortedSet(want)
		return nil, fmt.Errorf("unknown selected subjects: %s", strings.Join(missing, ","))
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Spec.ID < out[j].Spec.ID })
	return out, nil
}

func dedupeSubjects(subjects []ResolvedSubject) []ResolvedSubject {
	seen := map[string]struct{}{}
	out := make([]ResolvedSubject, 0, len(subjects))
	for _, subject := range subjects {
		if subject.Spec.ID == "" {
			continue
		}
		if _, ok := seen[subject.Spec.ID]; ok {
			continue
		}
		seen[subject.Spec.ID] = struct{}{}
		out = append(out, subject)
	}
	return out
}

func isCompatible(value string, allowed []string) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, item := range allowed {
		if strings.EqualFold(strings.TrimSpace(item), value) {
			return true
		}
	}
	return false
}

func isSkillCompatible(skill contracts.SkillSpec, framework string) bool {
	return isCompatible(framework, skill.CompatibleFrameworks)
}

func resolveRelative(baseDir, path string) string {
	path = strings.TrimSpace(path)
	if path == "" || filepath.IsAbs(path) {
		return path
	}
	return filepath.Clean(filepath.Join(baseDir, path))
}

func uniqueNonEmpty(items []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	sort.Strings(out)
	return out
}

func intersect(left, right []string) []string {
	rightSet := map[string]struct{}{}
	for _, item := range right {
		rightSet[item] = struct{}{}
	}
	var out []string
	for _, item := range left {
		if _, ok := rightSet[item]; ok {
			out = append(out, item)
		}
	}
	return uniqueNonEmpty(out)
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func sortedSkillKeys(m map[string]contracts.SkillSpec) []string {
	return sortedKeys(m)
}

func sortedSet(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func defaultString(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func sanitize(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == '-' || r == '.' {
			b.WriteRune(r)
			continue
		}
		b.WriteRune('_')
	}
	out := strings.Trim(b.String(), "_")
	if out == "" {
		return "unknown"
	}
	return out
}
