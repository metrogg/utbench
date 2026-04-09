package main

import (
	"path/filepath"
)

func getLastCmdStatus() (string, error) {
	rootPath, err := getIntelRdtRoot()
	if err != nil {
		return "", err
	}

	path := filepath.Join(rootPath, "info")
	lastCmdStatus, err := getIntelRdtParamString(path, "last_cmd_status")
	if err != nil {
		return "", err
	}

	return lastCmdStatus, nil
}
