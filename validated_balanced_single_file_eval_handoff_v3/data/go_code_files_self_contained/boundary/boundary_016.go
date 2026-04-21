package main

import (
	"math/rand"
)

func randomInprocEndpoint() (string, error) {
	rnd := make([]byte, 10)
	if _, err := rand.Read(rnd); err != nil {
		return "", err
	}
	return "inproc://rnd" + string(rnd), nil
}
