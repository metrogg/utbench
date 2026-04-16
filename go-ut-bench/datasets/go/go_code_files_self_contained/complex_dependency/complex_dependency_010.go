package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func GetGaugeHomeDirectory() (string, error) {
	customPluginRoot := os.Getenv(GaugeHome)
	if customPluginRoot != "" {
		return customPluginRoot, nil
	}
	if isWindows() {
		appDataDir := os.Getenv(appData)
		if appDataDir == "" {
			return "", fmt.Errorf("Failed to find plugin installation path. Could not get APPDATA")
		}
		return filepath.Join(appDataDir, ProductName), nil
	}
	userHome := getUserHomeFromEnv()
	return filepath.Join(userHome, DotGauge), nil
}
