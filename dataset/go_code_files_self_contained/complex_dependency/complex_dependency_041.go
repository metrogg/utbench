package main

import (
	"path/filepath"
)

func UpdaterBinPath() (string, error) {
	path, err := BinPath()
	if err != nil {
		return "", err
	}
	name, err := updaterBinName()
	if err != nil {
		return "", err
	}
	return filepath.Join(filepath.Dir(path), name), nil
}
