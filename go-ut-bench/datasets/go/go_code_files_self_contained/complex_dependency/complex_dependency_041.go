package main

import (
	"os"
)

func IsDocker() bool {
	_, e := os.Stat("/.dockerenv")
	if os.IsNotExist(e) {
		return false
	}

	return e == nil
}
