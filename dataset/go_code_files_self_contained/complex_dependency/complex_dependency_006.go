package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func BoltFile() (string, error) {
	dir, err := InfluxDir()
	if err != nil {
		return "", err
	}
	var file string
	filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if file != "" {
			return fmt.Errorf("bolt file found")
		}

		if strings.Contains(p, ".bolt") {
			file = p
		}

		return nil
	})

	if file == "" {
		return "", fmt.Errorf("bolt file not found")
	}

	return file, nil
}
