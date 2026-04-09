package main

import (
	"os"
)

func init() {
	var err error

	_SETTINGS, err = getSettingsForUid(os.Getuid())
	if err != nil {
		panic(err)
	}
}
