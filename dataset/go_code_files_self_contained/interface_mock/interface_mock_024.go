package main

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

func getVersion() string {
	slurp, err := ioutil.ReadFile(filepath.Join(camliDir, "VERSION"))
	if err == nil {
		return strings.TrimSpace(string(slurp))
	}
	return gitVersion()
}
