package main

import (
	"os"
	"path/filepath"
	"strings"
)

func EnsureUserAddressDir() error {
	xdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")
	if xdgRuntimeDir != "" {
		dirs := strings.Split(xdgRuntimeDir, ":")
		dir := filepath.Join(dirs[0], "buildkit")
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
		return os.Chmod(dir, 0700|os.ModeSticky)
	}
	return nil
}
