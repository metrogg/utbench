package main

import (
	"fmt"
	"path/filepath"
)

func GetExecutablePath() (string, error) {
	exePath, err := Readlink("/proc/self/exe")

	if err != nil {
		err = fmt.Errorf("can't read /proc/self/exe: %v", err)
	}

	return filepath.Clean(exePath), err
}
