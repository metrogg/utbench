package main

import (
	"errors"
	"os"
)

func GetRuntimeDir() (string, error) {
	if xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR"); xdgRuntimeDir != "" {
		return xdgRuntimeDir, nil
	}
	return "", errors.New("could not get XDG_RUNTIME_DIR")
}
