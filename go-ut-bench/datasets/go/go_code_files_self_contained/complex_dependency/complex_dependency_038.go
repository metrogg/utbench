package main

import (
	"fmt"
	"os"
	"path"
)

func GodelHomePath() (string, error) {
	// check the environment variable
	if godelHomeDir := os.Getenv(godelHomeEnvVar); godelHomeDir != "" {
		return godelHomeDir, nil
	}
	// if not present, create from home directory
	if userHomeDir := os.Getenv("HOME"); userHomeDir != "" {
		return path.Join(userHomeDir, defaultGodelHome), nil
	}
	return "", fmt.Errorf("failed to get %s home directory", AppName)
}
