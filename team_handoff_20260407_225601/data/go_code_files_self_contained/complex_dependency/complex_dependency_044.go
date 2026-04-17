package main

import (
	"path/filepath"
)

func GetOpenFileDescriptorCount() (int64, error) {
	files, err := filepath.Glob("/proc/self/fd/*")
	if err != nil {
		return 0, err
	}
	return int64(len(files)), nil
}
