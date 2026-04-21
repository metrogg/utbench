package main

import (
	"os"
	"path/filepath"
)

func BinPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(exePath)
}
