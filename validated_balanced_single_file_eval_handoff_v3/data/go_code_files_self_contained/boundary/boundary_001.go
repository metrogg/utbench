package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
)

func RandomCryptoString() (string, error) {
	buf := make([]byte, 32)
	n, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	if n != len(buf) {
		return "", fmt.Errorf("not enough random bytes read")
	}

	return hex.EncodeToString(buf), nil
}
