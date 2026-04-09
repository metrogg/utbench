package main

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

func PluginSearchPath() ([]string, error) {
	// Search all PATH directories
	paths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))

	// Search all bin/ directories in GOPATH
	for _, gopath := range strings.Split(os.Getenv("GOPATH"), string(os.PathListSeparator)) {
		paths = append(paths, path.Join(gopath, gopathPluginDir))
	}

	// Search the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	paths = append(paths, wd)

	// Search the directory of the current executable
	executableDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return nil, err
	}
	paths = append(paths, executableDir)

	return RemoveDuplicates(paths), nil
}
