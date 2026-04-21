package main

import (
	"os"
)

func homeDirectory() string {
	var homeDir string
	switch {
	case os.Getenv("CF_HOME") != "":
		homeDir = os.Getenv("CF_HOME")
	case os.Getenv("HOMEDRIVE")+os.Getenv("HOMEPATH") != "":
		homeDir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	default:
		homeDir = os.Getenv("USERPROFILE")
	}
	return homeDir
}
