package main

import (
	"os"
	"path"
	"time"
)

func GetTempDir() string {
	return path.Join(os.TempDir(), ToStr(time.Now().Nanosecond()))
}
