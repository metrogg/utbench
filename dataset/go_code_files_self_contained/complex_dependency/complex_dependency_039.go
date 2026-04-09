package main

import (
	"fmt"
	"path/filepath"
)

func GetDefaultPropertiesFile() (string, error) {
	envDir, err := GetDirInProject(EnvDirectoryName, "")
	if err != nil {
		return "", err
	}
	defaultEnvFile := filepath.Join(envDir, DefaultEnvDir, DefaultEnvFileName)
	if !FileExists(defaultEnvFile) {
		return "", fmt.Errorf("Default environment file does not exist: %s \n", defaultEnvFile)
	}
	return defaultEnvFile, nil
}
