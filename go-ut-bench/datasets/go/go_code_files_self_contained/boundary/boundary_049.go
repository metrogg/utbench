package main

import (
	"strings"
)

func GoRoot() (string, error) {
	output, err := execGo("go", nil, "", "env", "GOROOT")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output), nil
}
