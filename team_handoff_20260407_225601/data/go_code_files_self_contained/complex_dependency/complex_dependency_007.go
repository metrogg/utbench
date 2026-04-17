package main

import (
	"os"
	"path/filepath"
)

func SetGoEnvVariables() error {
	// set GOPATH environment variable to current value of GOPATH environment variable with symlinks resolved
	gopath := os.Getenv("GOPATH")
	if resolvedGoPath, err := filepath.EvalSymlinks(gopath); err == nil {
		gopath = resolvedGoPath
	}
	if err := os.Setenv("GOPATH", gopath); err != nil {
		return err
	}

	// set GOROOT environment variable to be current value of Go root determined by GoRoot() with symlinks resolved
	goroot, err := GoRoot()
	if err != nil {
		return err
	}
	if resolvedGoRoot, err := filepath.EvalSymlinks(goroot); err == nil {
		goroot = resolvedGoRoot
	}
	return os.Setenv("GOROOT", goroot)
}
