package main

import (
	"os"
	"path/filepath"
)

func cwdPeers() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Join(cwd, filename)
}
