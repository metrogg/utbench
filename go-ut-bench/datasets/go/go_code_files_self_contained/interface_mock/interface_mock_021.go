package main

import (
	"io"
	"math/rand"
)

func GenerateSalt() ([]byte, error) {
	var salt []byte = make([]byte, SALT_SIZE)

	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	return salt, nil
}
