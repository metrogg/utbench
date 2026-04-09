package main

import (
	"os"
	"path/filepath"
)

func AssetsPath() string {
	if caddyPath := os.Getenv("CADDYPATH"); caddyPath != "" {
		return caddyPath
	}
	return filepath.Join(userHomeDir(), ".caddy")
}
