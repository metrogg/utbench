package main

import (
	"os"
	"strings"
)

func WorkDir() (string, error) {
	wd := os.Getenv("GOGS_WORK_DIR")
	if len(wd) > 0 {
		return wd, nil
	}

	i := strings.LastIndex(AppPath, "/")
	if i == -1 {
		return AppPath, nil
	}
	return AppPath[:i], nil
}
