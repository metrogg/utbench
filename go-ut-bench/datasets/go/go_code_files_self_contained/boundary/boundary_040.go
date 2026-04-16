package main

import (
	"os"
)

func createConfigPath() error {
	path, err := GetConfigFilePath("")
	if err != nil {
		return err
	}
	return os.MkdirAll(path, 0700)
}
