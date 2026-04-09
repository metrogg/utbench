package main

import (
	"path/filepath"
)

func LocalConfigFileAbsolutePath() (string, error) {
	configDir, err := LocalConfigDirectoryAbsolutePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, LocalConfigFilename), nil
}
