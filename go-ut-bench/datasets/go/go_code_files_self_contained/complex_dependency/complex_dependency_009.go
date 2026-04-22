package main

import (
	"os"
	"path/filepath"
)

func Getwd() string {
	if cwd, err := os.Getwd(); err == nil {
		cwd, _ = filepath.Abs(cwd)
		return cwd
	} else {
		panic("Failed to retrieve current working directory")
	}
}
