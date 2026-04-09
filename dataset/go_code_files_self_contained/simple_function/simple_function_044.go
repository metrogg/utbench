package main

import (
	"strings"
)

func InstanceName() (string, error) {
	host, err := Hostname()
	if err != nil {
		return "", err
	}
	return strings.Split(host, ".")[0], nil
}
