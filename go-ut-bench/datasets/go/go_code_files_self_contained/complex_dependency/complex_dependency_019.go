package main

import (
	"io/ioutil"
	"os"
)

func setupFiles() (*os.File, *os.File) {
	outFile, _ := ioutil.TempFile("", "outputmock")
	errFile, _ := ioutil.TempFile("", "errormock")
	return outFile, errFile
}
