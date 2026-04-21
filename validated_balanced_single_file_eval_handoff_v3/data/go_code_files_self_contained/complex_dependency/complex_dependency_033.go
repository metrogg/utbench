package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func GetMachineID() (string, error) {
	machineID, err := ioutil.ReadFile("/etc/machine-id")
	if err != nil {
		return "", fmt.Errorf("failed to read /etc/machine-id: %v", err)
	}
	return strings.TrimSpace(string(machineID)), nil
}
