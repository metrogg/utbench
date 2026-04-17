package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func createConfigFile() error {
	filePath := filepath.Join(ConfigDir(), ConfigFileName)

	var err error

	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		if err = ioutil.WriteFile(filePath, []byte(defaultConfigFile), 0644); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
