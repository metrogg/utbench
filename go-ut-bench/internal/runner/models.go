package runner

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type modelConfig struct {
	Name      string
	Provider  string
	Endpoint  string
	Model     string
	APIKeyEnv string
	Enabled   bool
	Params    map[string]any
}

type modelsFile struct {
	Models map[string]struct {
		Enabled  bool   `yaml:"enabled"`
		Provider string `yaml:"provider"`
		Config   struct {
			APIEndpoint string         `yaml:"api_endpoint"`
			Model       string         `yaml:"model"`
			APIKeyEnv   string         `yaml:"api_key_env"`
			Parameters  map[string]any `yaml:"parameters"`
		} `yaml:"config"`
	} `yaml:"models"`
}

func loadModelConfigs(configPath string, selected []string) ([]modelConfig, error) {
	raw, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg modelsFile
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, err
	}
	if len(cfg.Models) == 0 {
		return nil, errors.New("models config is empty")
	}

	selectedSet := map[string]struct{}{}
	for _, item := range selected {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			selectedSet[trimmed] = struct{}{}
		}
	}

	out := make([]modelConfig, 0)
	for name, item := range cfg.Models {
		if len(selectedSet) > 0 {
			if _, ok := selectedSet[name]; !ok {
				continue
			}
		}
		if !item.Enabled {
			continue
		}
		if item.Config.APIKeyEnv == "" {
			return nil, fmt.Errorf("model %s missing api_key_env", name)
		}
		out = append(out, modelConfig{
			Name:      name,
			Provider:  strings.ToLower(strings.TrimSpace(item.Provider)),
			Endpoint:  strings.TrimSuffix(strings.TrimSpace(item.Config.APIEndpoint), "/"),
			Model:     strings.TrimSpace(item.Config.Model),
			APIKeyEnv: strings.TrimSpace(item.Config.APIKeyEnv),
			Enabled:   item.Enabled,
			Params:    item.Config.Parameters,
		})
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no enabled models selected: %v", selected)
	}
	return out, nil
}
