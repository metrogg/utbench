package main

import (
	"errors"
	"math/rand"
)

func generateSecretKey() ([]byte, error) {
	var secretKey [32]byte
	if _, err := rand.Read(secretKey[:]); err != nil {
		return nil, errors.Trace(err)
	}
	return secretKey[:], nil
}
