package main

import (
	"path"
)

func DefaultLocalConfigPath() (string, error) {
	dir, err := DefaultConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(dir, "config"), nil
}
