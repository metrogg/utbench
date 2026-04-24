package main

import (
	"io/ioutil"
	"os"
)

func tempfile() (*os.File, error) {
	f, err := ioutil.TempFile("", "rqlilte-snap-")
	if err != nil {
		return nil, err
	}
	return f, nil
}
