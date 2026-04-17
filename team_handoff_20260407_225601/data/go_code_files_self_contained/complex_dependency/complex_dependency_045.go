package main

import (
	"path/filepath"
)

func ExecutableFolder() (string, error) {
	p, err := Executable()
	if err != nil {
		return "", err
	}
	folder, _ := filepath.Split(p)
	return folder, nil
}
