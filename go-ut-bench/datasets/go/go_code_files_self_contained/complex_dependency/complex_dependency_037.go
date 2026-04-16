package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func Slaves() ([]string, error) {
	var slaves []string
	fis, err := ioutil.ReadDir(basePath)
	if err != nil {
		return nil, err
	}
	for i := range fis {
		if strings.HasPrefix(fis[i].Name(), slavePrefix) {
			if (fis[i].Mode() & os.ModeSymlink) == os.ModeSymlink {
				slaves = append(slaves, fis[i].Name())
			}
		}
	}
	return slaves, nil
}
