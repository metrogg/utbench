package config

import (
	"errors"
	"strings"

	"go-ut-bench/internal/contracts"
)

type AppConfig struct {
	DefaultDatasetRoot string
	DefaultOutputRoot  string
	DefaultDBPath      string
	DefaultModels      []string
	DefaultLanguages   []string
	DefaultClass       contracts.DatasetClass
	DefaultMode        contracts.RunMode
}

func Default() AppConfig {
	return AppConfig{
		DefaultDatasetRoot: "./datasets",
		DefaultOutputRoot:  "./artifacts",
		DefaultDBPath:      "./storage/utbench.db",
		DefaultModels:      []string{"deepseek"},
		DefaultLanguages:   []string{"python"},
		DefaultClass:       contracts.DatasetClassSelfContained,
		DefaultMode:        contracts.RunModeFull,
	}
}

func ValidateClass(v string) error {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	validClasses := map[string]bool{
		"self_contained": true,
		"module_level":    true,
	}
	for _, c := range strings.Split(v, ",") {
		c = strings.TrimSpace(c)
		if !validClasses[c] {
			return errors.New("dataset class must be self_contained or module_level (comma-separated allowed)")
		}
	}
	return nil
}
