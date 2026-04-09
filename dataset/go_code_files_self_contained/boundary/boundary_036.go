package main

import (
	"bytes"
	"io/ioutil"
)

func Token() ([]byte, error) {
	data, err := ioutil.ReadFile(TokenFile)
	if err != nil {
		return data, err
	}
	data = bytes.TrimSpace(data)
	return data, nil
}
