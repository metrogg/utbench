package main

import (
	"os"
	"path/filepath"
)

func UserConfigDir() string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome != "" {
		return filepath.Join(xdgConfigHome, "buildkit")
	}
	home := os.Getenv("HOME")
	if home != "" {
		return filepath.Join(home, ".config", "buildkit")
	}
	return ConfigDir
}
