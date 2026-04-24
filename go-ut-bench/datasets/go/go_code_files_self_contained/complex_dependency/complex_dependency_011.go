package main

import (
	"os"
	"path/filepath"
)

func Geted() string {
	if dir, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		return dir
	} else {
		panic("Failed to retrieve executable directory")
	}
}
