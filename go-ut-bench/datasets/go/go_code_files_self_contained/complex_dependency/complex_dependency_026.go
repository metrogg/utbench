package main

import (
	"os"
	"path/filepath"
)

func DataDir() string {
	dataHome := filepath.Join(os.Getenv("HOME"), DefaultDataDir)
	if xdg := os.Getenv("XDG_DATA_HOME"); len(xdg) > 0 {
		dataHome = xdg
	}

	return filepath.Join(dataHome, DataFolder)
}
