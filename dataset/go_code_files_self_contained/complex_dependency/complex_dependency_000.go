package main

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func GetTempDir() string {
	tempGaugeDir := filepath.Join(os.TempDir(), "gauge_temp")
	tempGaugeDir += strconv.FormatInt(time.Now().UnixNano(), 10)
	if !exists(tempGaugeDir) {
		os.MkdirAll(tempGaugeDir, NewDirectoryPermissions)
	}
	return tempGaugeDir
}
