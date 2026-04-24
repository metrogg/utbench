package main

import (
	"os"
	"path/filepath"
)

func GetCurrentDir() (string, error) {
	d, err := filepath.Abs(filepath.Dir(os.Args[0]))
	return d, err
}
