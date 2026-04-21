package main

import (
	"errors"
	"os"
)

func ListenAddress() (string, error) {
	listenAddr := os.Getenv("CAMLI_APP_LISTEN")
	if listenAddr == "" {
		return "", errors.New("CAMLI_APP_LISTEN is undefined")
	}
	return listenAddr, nil
}
