package main

import (
	"os"
	"path/filepath"
	"strings"
)

func setRoot() {
	wdStr, err := os.Getwd()
	if err == nil {
		gofpdfDir = ""
		list := strings.Split(filepath.ToSlash(wdStr), "/")
		for j := len(list) - 1; j >= 0 && list[j] != "gofpdf"; j-- {
			gofpdfDir = filepath.Join(gofpdfDir, "..")
		}
	} else {
		panic(err)
	}
}
