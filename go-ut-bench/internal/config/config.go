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
	if v != string(contracts.DatasetClassSelfContained) && v != string(contracts.DatasetClassModuleLevel) && v != string(contracts.DatasetClassComplexDependency) {
		return errors.New("dataset class must be self_contained, module_level or complex_dependency")
	}
	return nil
}
