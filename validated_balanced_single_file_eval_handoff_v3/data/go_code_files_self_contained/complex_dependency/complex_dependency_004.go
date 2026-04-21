package main

import (
	"os"
	"path/filepath"
)

func GetwdEvalSymLinks() (string, error) {
	// get wd
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// resolve any symlinks in path
	physicalWd, err := filepath.EvalSymlinks(wd)
	if err != nil {
		return "", err
	}
	return physicalWd, nil
}
