package main

import (
	"os"
	"path/filepath"
)

func ConfigDir() string {
	configHome := filepath.Join(os.Getenv("HOME"), DefaultConfigDir)
	if xdg := os.Getenv("XDG_CONFIG_HOME"); len(xdg) > 0 {
		configHome = xdg
	}

	return filepath.Join(configHome, ConfigFolder)
}
