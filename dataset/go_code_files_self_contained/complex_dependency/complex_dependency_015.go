package main

import (
	"io/ioutil"
	"os"
	"path"
)

func InstallDefaultUserConfig() error {
	err := os.MkdirAll(path.Dir(GetUserConfigFilePath()), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(GetUserConfigFilePath(), []byte(DefaultConfigFileContent), 0644)
}
