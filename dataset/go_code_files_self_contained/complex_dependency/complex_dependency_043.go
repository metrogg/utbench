package main

import (
	"os"
	"path"
)

func RcPath() (string, error) {
	homedrive := os.Getenv("HOMEDRIVE")
	if homedrive == "" {
		return "", ErrMissingHome
	}

	homedir := os.Getenv("HOMEDIR")
	if homedir == "" {
		return "", ErrMissingHome
	}

	return path.Join(homedrive, homedir, rcFilename), nil
}
