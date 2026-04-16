package main

import (
	"os"
	"path/filepath"
	"strings"
)

func DefaultSourceTree() string {
	paths := strings.Split(os.Getenv("GOPATH"), string(filepath.ListSeparator))
	if len(paths) > 0 && len(paths[0]) > 0 {
		return filepath.Join(paths[0], "src")
	}
	return "./"
}
