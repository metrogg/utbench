package main

import (
	"os"
	"path/filepath"
)

func GetBuildpackDir() (string, error) {
	var err error

	bpDir := os.Getenv("BUILDPACK_DIR")

	if bpDir == "" {
		bpDir, err = filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), ".."))

		if err != nil {
			return "", err
		}
	}

	return bpDir, nil
}
