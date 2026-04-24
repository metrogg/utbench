package main

import (
	"os"
	"path/filepath"
)

func GoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		home := os.Getenv("HOME")
		gopath = filepath.Join(home, "go")
	}
	return gopath
}
