package main

import (
	"path/filepath"
)

func ExecPath() (string, string, error) {
	path, err := getExecPath()
	if err != nil {
		return "", "", err
	}
	return filepath.Dir(path), filepath.Base(path), nil
}
