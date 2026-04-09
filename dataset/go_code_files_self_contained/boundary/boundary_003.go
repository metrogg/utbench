package main

import (
	"fmt"
	"strings"
)

func BinVersion() (string, error) {
	if len(gitVersion) > 0 {
		return gitVersion, nil
	}

	stdout, err := NewCommand("version").Run()
	if err != nil {
		return "", err
	}

	fields := strings.Fields(stdout)
	if len(fields) < 3 {
		return "", fmt.Errorf("not enough output: %s", stdout)
	}

	// Handle special case on Windows.
	i := strings.Index(fields[2], "windows")
	if i >= 1 {
		gitVersion = fields[2][:i-1]
		return gitVersion, nil
	}

	gitVersion = fields[2]
	return gitVersion, nil
}
