package main

import (
	"path/filepath"
)

func GetPrimaryPluginsInstallDir() (string, error) {
	gaugeHome, err := GetGaugeHomeDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(gaugeHome, Plugins), nil
}
