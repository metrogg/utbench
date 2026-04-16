package main

import (
	"io/ioutil"
	"strings"
)

func Read() (string, error) {
	keyPath, err := kiteKeyPath()
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
